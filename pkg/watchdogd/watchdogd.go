//go:build !tinygo

// Copyright 2021 the u-root Authors. All rights reserved
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

package watchdogd

import (
	"context"
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

const (
	defaultUDS = "/tmp/watchdogd"
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
	PettingOp chan WatchdogOperation

	// PettingOn indicate if there is an active petting session.
	PettingOn bool
}

// DaemonOpts contain operating parameters for bootstrapping a watchdog daemon.
type DaemonOpts struct {
	// Dev is the watchdog device. Ex: /dev/watchdog
	Dev string

	// nil uses the preset values. 0 disables the timeout.
	Timeout, PreTimeout *time.Duration

	// KeepAlive is the length of the keep alive interval.
	KeepAlive time.Duration

	// Monitors are called before each keepalive interval. If any monitor
	// function returns an error, the .
	Monitors []func() error

	// UDS is the name of daemon's unix domain socket.
	UDS string
}

// MonitorOops return an error if the kernel logs contain an oops.
func MonitorOops() error {
	dmesg := make([]byte, OopsBuffSize)
	n, err := unix.Klogctl(unix.SYSLOG_ACTION_READ_ALL, dmesg)
	if err != nil {
		return fmt.Errorf("syslog failed: %v", err)
	}
	if strings.Contains(string(dmesg[:n]), "Oops:") {
		return fmt.Errorf("founds Oops in dmesg")
	}
	return nil
}

// StartServing enters a loop of accepting and processing next incoming watchdogd operation call.
func (d *Daemon) StartServing(l *net.UnixListener) error {
	for {
		// All requests are processed sequentially.
		c, err := l.AcceptUnix()
		if err != nil {
			log.Printf("Failed to accept new request: %v", err)
			continue
		}
		b := make([]byte, 1)
		if _, err := io.ReadAtLeast(c, b, 1); err != nil {
			log.Printf("Failed to read operation bit, err: %v", err)
		}
		op, err := NewWatchdogOperation(b[0])
		if err != nil {
			return err
		}

		log.Printf("New op received: %c", op)
		var r error
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
func (d *Daemon) ArmWatchdog() error {
	if d.CurrentOpts == nil {
		return NewWatchdogError("Current daemon opts is nil, don't know how to arm Watchdog")
	}
	wd, err := watchdog.Open(d.CurrentOpts.Dev)
	if err != nil {
		// Most likely cause is /dev/watchdog does not exist.
		// Second most likely cause is another process (perhaps
		// another watchdogd?) has the file open.
		return NewWatchdogErrorf("Failed to arm: %v", err)
	}
	if d.CurrentOpts.Timeout != nil {
		if err := wd.SetTimeout(*d.CurrentOpts.Timeout); err != nil {
			d.CurrentWd.Close()
			return NewWatchdogErrorf("Failed to set timeout: %v", err)
		}
	}
	if d.CurrentOpts.PreTimeout != nil {
		if err := wd.SetPreTimeout(*d.CurrentOpts.PreTimeout); err != nil {
			d.CurrentWd.Close()
			return NewWatchdogErrorf("Failed to set pretimeout: %v", err)
		}
	}
	d.CurrentWd = wd
	log.Printf("Watchdog armed")
	return nil
}

// disarmWatchdog disarm the watchdog if already armed.
func (d *Daemon) DisarmWatchdog() error {
	if d.CurrentWd == nil {
		log.Printf("No armed Watchdog")
		return nil
	}
	if err := d.CurrentWd.MagicClose(); err != nil {
		NewWatchdogErrorf("Failed to disarm watchdog: %v", err)
	}
	log.Println("Watchdog disarming request went through (Watchdog will not be disabled if CONFIG_WATCHDOG_NOWAYOUT is enabled).")
	return nil
}

// doPetting sends keepalive signal to Watchdog when necessary.
//
// If at least one of the custom monitors failed check(s), it won't send a keepalive
// signal.
func (d *Daemon) DoPetting() error {
	if d.CurrentWd == nil {
		return NewWatchdogError("no reference to any Watchdog")
	}
	if err := doMonitors(d.CurrentOpts.Monitors); err != nil {
		return NewWatchdogErrorf("won't keepalive since at least one of the custom monitors failed: %v", err)
	}
	if err := d.CurrentWd.KeepAlive(); err != nil {
		return NewWatchdogErrorf("failed to keepalive: %w", err)
	}
	return nil
}

// startPetting starts Watchdog petting in a new goroutine.
func (d *Daemon) StartPetting() error {
	if d.PettingOn {
		return NewWatchdogError("Petting ongoing")
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
	return nil
}

// stopPetting stops an ongoing petting process if there is.
func (d *Daemon) StopPetting() error {
	if !d.PettingOn {
		// No petting on, simply return.
		return nil
	}

	var ret error = nil
	erredOut := func() {
		<-d.PettingOp
		ret = NewWatchdogErrorf("Stop petting times out after %d seconds", opStopPettingTimeoutSeconds)
	}
	// It will time out when there is no active petting.
	t := time.AfterFunc(opStopPettingTimeoutSeconds*time.Second, erredOut)
	defer t.Stop()
	d.PettingOp <- OpStop
	return ret
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
		return fmt.Errorf("failed to setup server: %v", err)
	}
	go func() {
		log.Println("Start serving.")
		err := d.StartServing(l)
	}()

	log.Println("Start arming watchdog initially.")
	if r := d.ArmWatchdog(); r != OpResultOk {
		return fmt.Errorf("initial arm failed")
	}

	if r := d.StartPetting(); r != OpResultOk {
		return fmt.Errorf("start petting failed")
	}

	for {
		select {
		case <-ctx.Done():
			cleanup()
		}
	}
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
func sendAndCheckResult(con *net.UnixConn, op WatchdogOperation) error {
	n, err := con.Write([]byte{byte(op)})
	if err != nil {
		return err
	}

	if n != 1 {
		return NewWatchdogError("no error; but message not delivered neither")
	}

	b := make([]byte, 1)
	if _, err := io.ReadAtLeast(con, b, 1); err != nil {
		log.Printf("Failed to read operation bit from server: %v", err)
	}
	ret, err := NewWatchdogOperation(b[0])
	if err != nil {
		return err
	}

	if ret != OpResultOk {
		return fmt.Errorf("non-Ok op result: %c", ret)
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

func NewClient() (*client, error) {
	return NewClientFromUDS(defaultUDS)
}
