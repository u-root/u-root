// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watchdog

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type mockDog struct {
	forceErrNo unix.Errno
	forceErr   error

	option        Option
	currentStatus uint32
	bootStatus    uint32
	timeout       uint32
	preTimeout    uint32
	timeLeft      uint32
	info          unix.WatchdogInfo
}

func (f *mockDog) unixSyscall(trap, a1, a2 uintptr, a3 unsafe.Pointer) (uintptr, uintptr, unix.Errno) {
	if f.forceErrNo != 0 {
		return 0, 0, f.forceErrNo
	}
	if trap != unix.SYS_IOCTL {
		return 0, 0, unix.EINVAL
	}
	switch a2 {
	case wdiocGetSupport:
		*(*unix.WatchdogInfo)(a3) = f.info
	case wdiocSetOptions:
		f.option = *(*Option)(a3)
	case wdiocSetTimeout:
		f.timeout = *(*uint32)(a3)
	case wdiocSetPreTimeout:
		f.preTimeout = *(*uint32)(a3)
	default:
		return 0, 0, unix.EINVAL
	}
	return 0, 0, 0
}

func (f *mockDog) unixIoctlGetUint32(fd int, req uint) (uint32, error) {
	if f.forceErr != nil {
		return 0, f.forceErr
	}
	if fd < 0 {
		return 0, fmt.Errorf("invalid file descriptor")
	}
	switch req {
	case wdiocGetStatus:
		return f.currentStatus, nil
	case wdiocGetBootStatus:
		return f.bootStatus, nil
	case wdiocGetTimeout:
		return f.timeout, nil
	case wdiocGetPreTimeout:
		return f.preTimeout, nil
	case wdiocGetTimeLeft:
		return f.timeLeft, nil
	default:
		return 0, fmt.Errorf("no valid value passed to unixIoctlGetUnit32 for req")
	}
}

func TestWatchdogSyscallFunctions(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	wd, err := Open(tmpFile.Name())
	if err != nil {
		t.Errorf("Could not open watchdog : %v", err)
	}
	defer wd.Close()

	m := mockDog{}

	wd.syscalls = &m

	for _, tt := range []struct {
		name       string
		forceErrno unix.Errno
		wantErr    error
		option     Option
		timeout    time.Duration
	}{
		{
			name:    "NoError",
			option:  OptionDisableCard,
			timeout: 0,
		},
		{
			name:       "FroceSyscallError",
			forceErrno: unix.EINVAL,
			wantErr:    unix.EINVAL,
		},
	} {
		m.forceErrNo = tt.forceErrno
		t.Run("Support"+tt.name, func(t *testing.T) {
			_, err := wd.Support()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Want: %q Got: %q", tt.name, tt.wantErr, err)
			}
		})

		t.Run("SetOption"+tt.name, func(t *testing.T) {
			if err := wd.SetOptions(tt.option); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
		})

		t.Run("SetTimeout"+tt.name, func(t *testing.T) {
			if err := wd.SetTimeout(tt.timeout); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
		})

		t.Run("SetPreTimeout"+tt.name, func(t *testing.T) {
			if err := wd.SetPreTimeout(tt.timeout); !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestWatchdogIoctlGetUint32Functions(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	wd, err := Open(tmpFile.Name())
	if err != nil {
		t.Errorf("Could not open watchdog : %v", err)
	}
	defer wd.Close()

	m := mockDog{}

	wd.syscalls = &m

	for _, tt := range []struct {
		name       string
		forceErrno unix.Errno
		forceErr   error
		wantErr    error
		option     Option
		timeout    time.Duration
		pretimeout time.Duration
		timeleft   time.Duration
		status     Status
	}{
		{
			name:       "NoError",
			forceErr:   nil,
			wantErr:    nil,
			option:     OptionDisableCard,
			pretimeout: 0,
			timeleft:   0,
			timeout:    0,
			status:     0,
		},
		{
			name:       "FroceSyscallError",
			forceErr:   errors.New("Duh"),
			wantErr:    nil,
			option:     OptionDisableCard,
			pretimeout: 0,
			timeleft:   0,
			timeout:    0,
			status:     0,
		},
	} {

		m.forceErr = tt.forceErr
		tt.wantErr = tt.forceErr
		t.Run("Status"+tt.name, func(t *testing.T) {
			s, err := wd.Status()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if s != tt.status {
				t.Errorf("Test %q failed. Got: %d Want %d", tt.name, s, tt.status)
			}
		})

		t.Run("BootStatus"+tt.name, func(t *testing.T) {
			s, err := wd.BootStatus()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if s != tt.status {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, s, tt.status)
			}
		})

		t.Run("GetTimeout"+tt.name, func(t *testing.T) {
			s, err := wd.Timeout()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if s != tt.timeout {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, s, tt.timeout)
			}
		})

		t.Run("GetPreTimeout"+tt.name, func(t *testing.T) {
			s, err := wd.PreTimeout()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if s != tt.timeout {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, s, tt.pretimeout)
			}
		})

		t.Run("GetTimeLeft"+tt.name, func(t *testing.T) {
			s, err := wd.TimeLeft()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if s != tt.timeout {
				t.Errorf("Test %q failed. Got: %q Want %q", tt.name, s, tt.timeleft)
			}
		})
	}
}

func TestSetTimeoutError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	wd, err := Open(tmpFile.Name())
	if err != nil {
		t.Errorf("Could not open watchdog : %v", err)
	}
	defer wd.Close()

	m := mockDog{}

	wd.syscalls = &m
	wantErr := errors.New("watchdog timeout set to 0s, wanted 5ns")

	if err := wd.SetTimeout(5); err != nil {
		if !strings.Contains(err.Error(), wantErr.Error()) {
			t.Errorf("SetTimeout failed. Want: %q Got: %q", wantErr, err)
		}
		return
	}
	t.Error("TestSetTimeout succeeded but shouldnt")
}

func TestSetPreTimeoutError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	wd, err := Open(tmpFile.Name())
	if err != nil {
		t.Errorf("Could not open watchdog : %v", err)
	}
	defer wd.Close()

	m := mockDog{}

	wd.syscalls = &m
	wantErr := errors.New("watchdog pretimeout set to 0s, wanted 5ns")

	if err := wd.SetPreTimeout(5); err != nil {
		if !strings.Contains(err.Error(), wantErr.Error()) {
			t.Errorf("SetTimeout failed. Want: %q Got: %q", wantErr, err)
		}
		return
	}
	t.Error("TestSetTimeout succeeded but shouldnt")
}
