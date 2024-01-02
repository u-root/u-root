// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtc

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/guest"
	"golang.org/x/sys/unix"
)

type testSyscalls struct{}

func (tsc testSyscalls) ioctlGetRTCTime(fd int) (*unix.RTCTime, error) {
	return &unix.RTCTime{}, nil
}

func (tsc testSyscalls) ioctlSetRTCTime(fd int, time *unix.RTCTime) error {
	return nil
}

func OpenMockRTC() (*RTC, error) {
	var tsc testSyscalls
	f, err := os.CreateTemp(os.TempDir(), "rtc-*")
	if err != nil {
		return nil, err
	}
	return &RTC{
		f,
		tsc,
	}, nil
}

func TestOpenRTC(t *testing.T) {
	rtc, err := OpenMockRTC()
	if err != nil {
		t.Errorf("OpenRTC got: %v; want nil", err)
	}
	if err := rtc.Close(); err != nil {
		t.Errorf("Failed to close RTC: %v", err)
	}
}

func TestOpenRealRTC(t *testing.T) {
	guest.SkipIfNotInVM(t)

	// The u-root amd64 VM does not seem to have a RTC device
	if runtime.GOARCH == "amd64" {
		t.Skip("Test not supported in amd64 Qemu")
	}
	rtc, err := OpenRTC()
	if err != nil {
		t.Errorf("OpenRTC got: %v; want nil", err)
	}
	if err := rtc.Close(); err != nil {
		t.Errorf("Failed to close RTC: %v", err)
	}
}

func TestSet(t *testing.T) {
	rtc, err := OpenMockRTC()
	if err != nil {
		t.Errorf("Error opening RTC: %v", err)
	}
	if err := rtc.Set(time.Now()); err != nil {
		t.Errorf("rtc.Set got: %v; want nil", err)
	}
	if err := rtc.Close(); err != nil {
		t.Errorf("Failed to close RTC: %v", err)
	}
}

func TestRead(t *testing.T) {
	rtc, err := OpenMockRTC()
	if err != nil {
		t.Errorf("Error opening RTC: %v", err)
	}
	if _, err = rtc.Read(); err != nil {
		t.Errorf("rtc.Read got: %v; want nil", err)
	}
	if err := rtc.Close(); err != nil {
		t.Errorf("Failed to close RTC: %v", err)
	}
}
