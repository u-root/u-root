// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package watchdog provides functions for interacting with the Linux watchdog.
//
// The basic usage is:
//
//	wd, err := watchdog.Open(watchdog.Dev)
//	while running {
//	    wd.KeepAlive()
//	}
//	wd.MagicClose()
//
// Open() arms the watchdog. MagicClose() disarms the watchdog.
//
// Note not every watchdog driver supports every function!
//
// For more, see:
// https://www.kernel.org/doc/Documentation/watchdog/watchdog-api.txt
package watchdog

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Dev is the name of the first watchdog. If there are multiple watchdogs, they
// are named /dev/watchdog0, /dev/watchdog1, ...
const Dev = "/dev/watchdog"

// Various ioctl numbers.
const (
	wdiocGetSupport    = 0x80285700
	wdiocGetStatus     = 0x80045701
	wdiocGetBootStatus = 0x80045702
	wdiocGetTemp       = 0x80045703
	wdiocSetOptions    = 0x80045704
	wdiocKeepAlive     = 0x80045705
	wdiocSetTimeout    = 0xc0045706
	wdiocGetTimeout    = 0x80045707
	wdiocSetPreTimeout = 0xc0045708
	wdiocGetPreTimeout = 0x80045709
	wdiocGetTimeLeft   = 0x8004570a
)

// Status contains flags returned by Status() and BootStatus(). These are the
// same flags used for Support()'s options field.
type Status int32

// Bitset of possible flags for the Status() type.
const (
	// Unknown flag error
	StatusUnknown Status = -1
	// Reset due to CPU overheat
	StatusOverheat Status = 0x0001
	// Fan failed
	StatusFanFault Status = 0x0002
	// External relay 1
	StatusExtern1 Status = 0x0004
	// ExStatusl relay 2
	StatusExtern2 Status = 0x0008
	// Power bad/power fault
	StatusPowerUnder Status = 0x0010
	// Card previously reset the CPU
	StatusCardReset Status = 0x0020
	// Power over voltage
	StatusPowerOver Status = 0x0040
	// Set timeout (in seconds)
	StatusSetTimeout Status = 0x0080
	// Supports magic close char
	StatusMagicClose Status = 0x0100
	// Pretimeout (in seconds), get/set
	StatusPreTimeout Status = 0x0200
	// Watchdog triggers a management or other external alarm not a reboot
	StatusAlarmOnly Status = 0x0400
	// Keep alive ping reply
	StatusKeepAlivePing Status = 0x8000
)

// Option are options passed to SetOptions().
type Option int32

// Bitset of possible flags for the Option type.
const (
	// Unknown status error
	OptionUnknown Option = -1
	// Turn off the watchdog timer
	OptionDisableCard Option = 0x0001
	// Turn on the watchdog timer
	OptionEnableCard Option = 0x0002
	// Kernel panic on temperature trip
	OptionTempPanic Option = 0x0004
)

type syscalls interface {
	unixSyscall(uintptr, uintptr, uintptr, unsafe.Pointer) (uintptr, uintptr, unix.Errno)
	unixIoctlGetUint32(int, uint) (uint32, error)
}

// Watchdog holds the descriptor of an open watchdog driver.
type Watchdog struct {
	*os.File
	syscalls
}

type realSyscalls struct{}

func (r *realSyscalls) unixSyscall(trap, a1, a2 uintptr, a3 unsafe.Pointer) (uintptr, uintptr, unix.Errno) {
	return unix.Syscall(trap, a1, a2, uintptr(a3))
}

func (r *realSyscalls) unixIoctlGetUint32(fd int, req uint) (uint32, error) {
	return unix.IoctlGetUint32(fd, req)
}

// Open arms the watchdog.
func Open(dev string) (*Watchdog, error) {
	f, err := os.OpenFile(dev, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &Watchdog{
			File:     f,
			syscalls: &realSyscalls{},
		},
		nil
}

// Close closes the device without disarming the watchdog.
func (w *Watchdog) Close() error {
	return w.File.Close()
}

// MagicClose disarms the watchdog. However if the kernel is compiled with
// CONFIG_WATCHDOG_NOWAYOUT=y, there may be no way to disarm the watchdog.
func (w *Watchdog) MagicClose() error {
	if _, err := w.File.Write([]byte("V")); err != nil {
		w.File.Close()
		return err
	}
	return w.File.Close()
}

// Support returns the WatchdogInfo struct.
func (w *Watchdog) Support() (*unix.WatchdogInfo, error) {
	var wi unix.WatchdogInfo
	if _, _, err := w.unixSyscall(unix.SYS_IOCTL, w.File.Fd(), wdiocGetSupport, unsafe.Pointer(&wi)); err != 0 {
		return nil, err
	}
	return &wi, nil
}

// Status returns the current status.
func (w *Watchdog) Status() (Status, error) {
	flags, err := w.unixIoctlGetUint32(int(w.File.Fd()), wdiocGetStatus)
	if err != nil {
		return StatusUnknown, err
	}
	return Status(flags), nil
}

// BootStatus returns the status at the last reboot.
func (w *Watchdog) BootStatus() (Status, error) {
	flags, err := w.unixIoctlGetUint32(int(w.File.Fd()), wdiocGetBootStatus)
	if err != nil {
		return StatusUnknown, err
	}
	return Status(flags), nil
}

// SetOptions can be used to control some aspects of the cards operation.
func (w *Watchdog) SetOptions(options Option) error {
	if _, _, err := w.unixSyscall(unix.SYS_IOCTL, w.File.Fd(), wdiocSetOptions, unsafe.Pointer(&options)); err != 0 {
		return err
	}
	return nil
}

// KeepAlive pets the watchdog.
func (w *Watchdog) KeepAlive() error {
	_, err := w.File.WriteString("1")
	return err
}

// SetTimeout sets the watchdog timeout on the fly. It returns an error if the
// timeout gets set to the wrong value. timeout must be a multiple of seconds;
// otherwise, an error is returned.
func (w *Watchdog) SetTimeout(timeout time.Duration) error {
	to := timeout / time.Second
	if _, _, err := w.unixSyscall(unix.SYS_IOCTL, w.File.Fd(), wdiocSetTimeout, unsafe.Pointer(&timeout)); err != 0 {
		return err
	}
	gotTimeout := to * time.Second
	if gotTimeout != timeout {
		return fmt.Errorf("watchdog timeout set to %v, wanted %v", gotTimeout, timeout)
	}
	return nil
}

// Timeout returns the current watchdog timeout.
func (w *Watchdog) Timeout() (time.Duration, error) {
	timeout, err := w.unixIoctlGetUint32(int(w.File.Fd()), wdiocGetTimeout)
	if err != nil {
		return 0, err
	}
	return time.Duration(timeout) * time.Second, nil
}

// SetPreTimeout sets the watchdog pretimeout on the fly. The pretimeout is the
// duration before triggering the preaction (such as an NMI, interrupt, ...).
// timeout must be a multiple of seconds; otherwise, an error is returned.
func (w *Watchdog) SetPreTimeout(timeout time.Duration) error {
	to := timeout / time.Second
	if _, _, err := w.unixSyscall(unix.SYS_IOCTL, w.File.Fd(), wdiocSetPreTimeout, unsafe.Pointer(&timeout)); err != 0 {
		return err
	}
	gotTimeout := to * time.Second
	if gotTimeout != timeout {
		return fmt.Errorf("watchdog pretimeout set to %v, wanted %v", gotTimeout, timeout)
	}
	return nil
}

// PreTimeout returns the current watchdog pretimeout.
func (w *Watchdog) PreTimeout() (time.Duration, error) {
	timeout, err := w.unixIoctlGetUint32(int(w.File.Fd()), wdiocGetPreTimeout)
	if err != nil {
		return 0, err
	}
	return time.Duration(timeout) * time.Second, nil
}

// TimeLeft returns the duration before the reboot (to the nearest second).
func (w *Watchdog) TimeLeft() (time.Duration, error) {
	left, err := w.unixIoctlGetUint32(int(w.File.Fd()), wdiocGetTimeLeft)
	if err != nil {
		return 0, err
	}
	return time.Duration(left) * time.Second, nil
}
