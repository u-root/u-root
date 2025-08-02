// Copyright 2021-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// watchdogd implements a background process which periodically issues a
// keepalive.
//
// It starts in the running+armed state:
//
//              | watchdogd Running     | watchdogd Stopped
//     ---------+-----------------------+--------------------------
//     Watchdog | watchdogd is actively | machine will soon reboot
//     Armed    | keeping machine alive |
//     ---------+-----------------------+--------------------------
//     Watchdog | a hang will not       | a hang will not reboot
//     Disarmed | reboot the machine    | the machine
//

//go:build !tinygo

package watchdogd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/watchdog"
	"golang.org/x/sys/unix"
)

var timeoutIgnore = time.Duration(-1)

var (
	ErrInvalidMonitor     = errors.New("unrecognized monitor")
	ErrNoCommandSpecified = errors.New("no command specified")
	ErrMissingArgument    = errors.New("missing argument")
	ErrTooManyArguments   = errors.New("too many arguments")
)

const defaultUDS = "/tmp/watchdogd"

const (
	OpStop     = 'S' // Stop the watchdogd petting.
	OpContinue = 'C' // Continue the watchdogd petting.
	OpDisarm   = 'D' // Disarm the watchdog.
	OpArm      = 'A' // Arm the watchdog.
)

const (
	OpResultOk        = 'O' // Ok.
	OpResultError     = 'E' // Error.
	OpResultInvalidOp = 'I' // Invalid Op.
)

const (
	opStopPettingTimeoutSeconds = 10
)

// Daemon contains running states of an instance of the daemon.
type Daemon struct {
	// CurrentOpts is current operating parameters for the daemon.
	//
	// It is assigned at the first call of Run and updated on each subsequent call of it.
	CurrentOpts *DaemonOpts

	// CurrentWd is an open file descriptor to the watchdog device specified in the daemon options.
	CurrentWd *watchdog.Watchdog

	// PettingOp syncs the signal to continue or stop petting the watchdog.
	PettingOp chan int

	// PettingOn indicate if there is an active petting session.
	PettingOn bool
}

// DaemonOpts contain operating parameters for bootstrapping a watchdog daemon.
type DaemonOpts struct {
	// Dev is the watchdog device. Ex: /dev/watchdog
	Dev string

	// nil uses the preset values. 0 disables the timeout.
	Timeout    time.Duration
	PreTimeout time.Duration

	// KeepAlive is the length of the keep alive interval.
	KeepAlive time.Duration

	// Monitors are called before each keepalive interval. If any monitor
	// function returns an error, the .
	Monitors []func() error

	// UDS is the name of daemon's unix domain socket.
	UDS string
}

// Abstract flag initialization to the DaemonOpts struct so we can separately define it for tinygo and non-tinygo builds.
func (d *DaemonOpts) InitFlags() (fs *flag.FlagSet) {
	fs = flag.NewFlagSet("run", flag.PanicOnError)
	fs.StringVar(&d.Dev, "dev", watchdog.Dev, "device")
	fs.DurationVar(&d.Timeout, "timeout", timeoutIgnore, "duration before timing out")
	fs.DurationVar(&d.PreTimeout, "pre_timeout", timeoutIgnore, "duration for pretimeout")
	fs.DurationVar(&d.KeepAlive, "keep_alive", 5*time.Second, "duration between issuing keepalive")
	fs.StringVar(&d.UDS, "uds", defaultUDS, "unix domain socket")
	return
}

// MonitorOops return an error if the kernel logs contain an oops.
func MonitorOops() error {
	dmesg := make([]byte, 256*1024)
	n, err := unix.Klogctl(unix.SYSLOG_ACTION_READ_ALL, dmesg)
	if err != nil {
		return fmt.Errorf("syslog failed: %w", err)
	}
	if strings.Contains(string(dmesg[:n]), "Oops:") {
		return fmt.Errorf("founds Oops in dmesg")
	}
	return nil
}

// StartServing enters a loop of accepting and processing next incoming watchdogd operation call.
func (d *Daemon) StartServing(l *net.UnixListener) {
	for { // All requests are processed sequentially.
		c, err := l.AcceptUnix()
		if err != nil {
			log.Printf("Failed to accept new request: %v", err)
			continue
		}
		b := make([]byte, 1) // Expect single byte operation instruction.
		if _, err := io.ReadAtLeast(c, b, 1); err != nil {
			log.Printf("Failed to read operation bit, err: %v", err)
		}
		op := int(b[0])
		log.Printf("New op received: %c", op)
		var r rune
		switch op {
		case OpStop:
			r = d.StopPetting()
		case OpContinue:
			r = d.StartPetting()
		case OpArm:
			r = d.ArmWatchdog()
		case OpDisarm:
			r = d.DisarmWatchdog()
		default:
			r = OpResultInvalidOp
		}
		c.Write([]byte{byte(r)})
		c.Close()
	}
}

// setupListener sets up a new "unix" network listener for the daemon.
func setupListener(uds string) (*net.UnixListener, func(), error) {
	os.Remove(uds)

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: uds, Net: "unix"})
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		os.Remove(uds)
	}
	return l, cleanup, nil
}

// armWatchdog starts watchdog timer.
func (d *Daemon) ArmWatchdog() rune {
	if d.CurrentOpts == nil {
		log.Printf("Current daemon opts is nil, don't know how to arm Watchdog")
		return OpResultError
	}
	wd, err := watchdog.Open(d.CurrentOpts.Dev)
	if err != nil {
		// Most likely cause is /dev/watchdog does not exist.
		// Second most likely cause is another process (perhaps
		// another watchdogd?) has the file open.
		log.Printf("Failed to arm: %v", err)
		return OpResultError
	}
	if d.CurrentOpts.Timeout != timeoutIgnore {
		if err := wd.SetTimeout(d.CurrentOpts.Timeout); err != nil {
			d.CurrentWd.Close()
			log.Printf("Failed to set timeout: %v", err)
			return OpResultError
		}
	}
	if d.CurrentOpts.PreTimeout != timeoutIgnore {
		if err := wd.SetPreTimeout(d.CurrentOpts.PreTimeout); err != nil {
			d.CurrentWd.Close()
			log.Printf("Failed to set pretimeout: %v", err)
			return OpResultError
		}
	}
	d.CurrentWd = wd
	log.Printf("Watchdog armed")
	return OpResultOk
}

// disarmWatchdog disarm the watchdog if already armed.
func (d *Daemon) DisarmWatchdog() rune {
	if d.CurrentWd == nil {
		log.Printf("No armed Watchdog")
		return OpResultOk
	}
	if err := d.CurrentWd.MagicClose(); err != nil {
		log.Printf("Failed to disarm watchdog: %v", err)
		return OpResultError
	}
	log.Println("Watchdog disarming request went through (Watchdog will not be disabled if CONFIG_WATCHDOG_NOWAYOUT is enabled).")
	return OpResultOk
}

// doPetting sends keepalive signal to Watchdog when necessary.
//
// If at least one of the custom monitors failed check(s), it won't send a keepalive
// signal.
func (d *Daemon) DoPetting() error {
	if d.CurrentWd == nil {
		return fmt.Errorf("no reference to any Watchdog")
	}
	if err := doMonitors(d.CurrentOpts.Monitors); err != nil {
		return fmt.Errorf("won't keepalive since at least one of the custom monitors failed: %w", err)
	}
	if err := d.CurrentWd.KeepAlive(); err != nil {
		return err
	}
	return nil
}

// startPetting starts Watchdog petting in a new goroutine.
func (d *Daemon) StartPetting() rune {
	if d.PettingOn {
		log.Printf("Petting ongoing")
		return OpResultError
	}

	go func() {
		d.PettingOn = true
		defer func() { d.PettingOn = false }()
		for {
			select {
			case op := <-d.PettingOp:
				if op == OpStop {
					log.Println("Petting stopped.")
					return
				}
			case <-time.After(d.CurrentOpts.KeepAlive):
				if err := d.DoPetting(); err != nil {
					log.Printf("Failed to keep alive: %v", err)
					// Keep trying to pet until the watchdog times out.
				}
			}
		}
	}()

	log.Println("Start petting watchdog.")
	return OpResultOk
}

// stopPetting stops an ongoing petting process if there is.
func (d *Daemon) StopPetting() rune {
	if !d.PettingOn {
		return OpResultOk
	} // No petting on, simply return.
	r := OpResultOk
	erredOut := func() {
		<-d.PettingOp
		log.Printf("Stop petting times out after %d seconds", opStopPettingTimeoutSeconds)
		r = OpResultError
	}
	// It will time out when there is no active petting.
	t := time.AfterFunc(opStopPettingTimeoutSeconds*time.Second, erredOut)
	defer t.Stop()
	d.PettingOp <- OpStop
	return r
}

// Run starts up the daemon.
//
// That includes:
// 1) Starts listening for watchdog(d) operation requests over unix network.
// 2) Arms the watchdog timer if it is not already armed.
// 3) Starts petting the watchdog timer.
func Run(ctx context.Context, opts *DaemonOpts) error {
	log.SetPrefix("watchdogd: ")
	defer log.Printf("Daemon quit")
	d := NewDaemon(opts)
	l, cleanup, err := setupListener(d.CurrentOpts.UDS)
	if err != nil {
		return fmt.Errorf("failed to setup server: %w", err)
	}
	go func() {
		log.Println("Start serving.")
		d.StartServing(l)
	}()

	log.Println("Start arming watchdog initially.")
	if r := d.ArmWatchdog(); r != OpResultOk {
		return fmt.Errorf("initial arm failed")
	}

	if r := d.StartPetting(); r != OpResultOk {
		return fmt.Errorf("start petting failed")
	}

	<-ctx.Done()
	cleanup()
	return nil
}

// doMonitors is a helper function to run the monitors.
//
// If there is anything wrong identified, it serves as a signal to stop
// petting Watchdog.
func doMonitors(monitors []func() error) error {
	for _, m := range monitors {
		if err := m(); err != nil {
			return err
		}
	}
	// All monitors return normal.
	return nil
}

func NewDaemon(opts *DaemonOpts) *Daemon {
	d := &Daemon{
		CurrentOpts: opts,
		PettingOp:   make(chan int),
		PettingOn:   false,
	}
	return d
}

type client struct {
	Conn *net.UnixConn
}

func (c *client) Stop() error {
	return sendAndCheckResult(c.Conn, OpStop)
}

func (c *client) Continue() error {
	return sendAndCheckResult(c.Conn, OpContinue)
}

func (c *client) Disarm() error {
	return sendAndCheckResult(c.Conn, OpDisarm)
}

func (c *client) Arm() error {
	return sendAndCheckResult(c.Conn, OpArm)
}

// sendAndCheckResult sends operation bit and evaluates result.
func sendAndCheckResult(c *net.UnixConn, op int) error {
	n, err := c.Write([]byte{byte(op)})
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("no error; but message not delivered neither")
	}
	b := make([]byte, 1)
	if _, err := io.ReadAtLeast(c, b, 1); err != nil {
		log.Printf("Failed to read operation bit from server: %v", err)
	}
	r := int(b[0])
	if r != OpResultOk {
		return fmt.Errorf("non-Ok op result: %c", r)
	}
	return nil
}

func NewClientFromUDS(uds string) (*client, error) {
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: uds, Net: "unix"})
	if err != nil {
		return nil, err
	}
	return &client{Conn: conn}, nil
}

// Create a new client to communicate with the watchdog daemon.
// In the previous implementation, the watchdog was created by finding the process id of the daemon called watchdogd.
func New() (*client, error) {
	return NewClientFromUDS(defaultUDS)
}
