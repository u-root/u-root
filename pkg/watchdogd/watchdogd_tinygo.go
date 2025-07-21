// Copyright 2021-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// watchdogd implements a background process which periodically issues a
// keepalive.
// This implementation revives an old implementation, that was replaced by a more
// robust implementation utilizing a unix domain socket.
// As of now (24/08) `tinygo` does not support linux networking.
// Until the implementation of this feature is finished, this implementation will be used.
//
// It starts in the running+armed state:
//
//	         | watchdogd Running     | watchdogd Stopped
//	---------+-----------------------+--------------------------
//	Watchdog | watchdogd is actively | machine will soon reboot
//	Armed    | keeping machine alive |
//	---------+-----------------------+--------------------------
//	Watchdog | a hang will not       | a hang will not reboot
//	Disarmed | reboot the machine    | the machine
//
// The following signals control changing state:
//
//   - STOP: running -> stopped
//   - CONT: stopped -> running
//   - USR1: armed -> disarmed
//   - USR2: disarmed -> armed

//go:build tinygo

package watchdogd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
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

const (
	// DaemonProcess is the name of the daemon process.
	DaemonProcess = "watchdogd"

	// ProcessGlob is the glob pattern for finding the watchdogd process.
	ProcessGlob = "/proc/*/comm"

	// oopsBuffSize is the size of the buffer for reading the kernel logs.
	OopsBuffSize = 256 * 1024
)

// DaemonOpts contain options for the watchdog daemon.
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
}

// Daemon represents a daemon running in a separate process.
type Daemon os.Process

// Abstract flag initialization to the DaemonOpts struct so we can separately define it for tinygo and non-tinygo builds.
func (d *DaemonOpts) InitFlags() (fs *flag.FlagSet) {
	fs = flag.NewFlagSet("run", flag.PanicOnError)
	fs.StringVar(&d.Dev, "dev", watchdog.Dev, "device")
	fs.DurationVar(&d.Timeout, "timeout", timeoutIgnore, "duration before timing out")
	fs.DurationVar(&d.PreTimeout, "pre_timeout", timeoutIgnore, "duration for pretimeout")
	fs.DurationVar(&d.KeepAlive, "keep_alive", 5*time.Second, "duration between issuing keepalive")
	return
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

// Run runs the watchdog on the current goroutine. The USR1, USR2, STOP and
// CONT signals are used to control, so consider using a dedicated process.
// Consider using the watchdogd command in u-root. Cancelling the context will
// leave with the armed/disarmed state as is.
func Run(ctx context.Context, opts *DaemonOpts) error {
	defer log.Println("watchdogd: Daemon quit")

	signals := make(chan os.Signal, 5)
	signal.Notify(signals, unix.SIGUSR1, unix.SIGUSR2)
	defer signal.Stop(signals)

	for {
		wd, err := watchdog.Open(opts.Dev)
		if err != nil {
			// Most likely cause is /dev/watchdog does not exist.
			// Second most likely cause is another process (perhaps
			// another watchdogd?) has the file open.
			return fmt.Errorf("watchdog: Failed to arm: %v", err)
		}
		if opts.Timeout != timeoutIgnore {
			if err := wd.SetTimeout(opts.Timeout); err != nil {
				wd.Close()
				return fmt.Errorf("watchdog: Failed to set timeout: %v", err)
			}
		}
		if opts.PreTimeout != timeoutIgnore {
			if err := wd.SetPreTimeout(opts.PreTimeout); err != nil {
				wd.Close()
				return fmt.Errorf("watchdog: Failed to set pretimeout: %v", err)
			}
		}
		log.Println("watchdog: Armed")

	armed: // Loop while armed. SIGUSR1 to break.
		for {
			select {
			case <-time.After(opts.KeepAlive):
				doMonitors(opts.Monitors)
				if err := wd.KeepAlive(); err != nil {
					log.Printf("watchdog: Failed to keepalive: %v", err)
					// Keep trying to pet until the watchdog times out.
				}
			case s := <-signals:
				if s == unix.SIGUSR1 {
					break armed
				}
			case <-ctx.Done():
				return wd.Close()
			}
		}
		if err := wd.MagicClose(); err != nil {
			log.Printf("watchdog: Failed to disarm: %v", err)
		} else {
			log.Println("watchdog: Disarmed")
		}

	disarmed: // Loop while disarmed. SIGUSR2 to break.
		for {
			select {
			case s := <-signals:
				if s == unix.SIGUSR2 {
					break disarmed
				}
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// doMonitors is a helper function to run the monitors.
func doMonitors(monitors []func() error) {
	for _, m := range monitors {
		if err := m(); err != nil {
			log.Printf("watchdog: Stopping keepalive due to: %v", err)
			// Stop the current process.
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				// We can't stop the process, so take it down.
				log.Fatalf("watchdog: Error stopping: %v", err)
			}
			if err := (*Daemon)(p).Stop(); err != nil {
				// We can't stop the process, so take it down.
				log.Fatalf("watchdog: Error stopping: %v", err)
			}

			// Someone intentionally sent a SIGCONT (probably via a
			// `watchdogd continue`). They probably have some
			// unfinished business with the machine, so continue
			// petting.
			break
		}
	}
}

// Find returns the process id of the daemon called watchdogd.
func New() (*Daemon, error) {
	files, err := filepath.Glob(ProcessGlob)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		// Ignore errors since /proc changes frequently.
		comm, _ := ioutil.ReadFile(f)
		// Skip matches where the wildcard is not an integer.
		pidStr := filepath.Base(filepath.Dir(f))
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}
		// Skip matches for the current process.
		if pid == os.Getpid() {
			continue
		}
		// Skip kernel workers. This is the same mechanism used by ps.
		cmdline, err := ioutil.ReadFile(filepath.Join(filepath.Dir(f), "cmdline"))
		if err != nil || len(cmdline) == 0 {
			continue
		}
		// /proc files have a gratuitous newline.
		if string(comm) == DaemonProcess+"\n" {
			p, err := os.FindProcess(pid)
			return (*Daemon)(p), err
		}
	}
	return nil, fmt.Errorf("could not find %q", DaemonProcess)
}

// Stop stops the daemon. It can be resumed with Continue().
func (d *Daemon) Stop() error {
	return (*os.Process)(d).Signal(unix.SIGSTOP)
}

// Continue continues the daemon from a previous Stop().
func (d *Daemon) Continue() error {
	return (*os.Process)(d).Signal(unix.SIGCONT)
}

// Disarm sends a signal to the watchdog daemon to disarm.
func (d *Daemon) Disarm() error {
	return (*os.Process)(d).Signal(unix.SIGUSR1)
}

// Arm sends a signal to the watchdog deamon to arm.
func (d *Daemon) Arm() error {
	return (*os.Process)(d).Signal(unix.SIGUSR2)
}
