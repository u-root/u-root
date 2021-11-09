// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"testing"
)

func TestWatchdogRunning(t *testing.T) {
	i := GetMockIPMI()
	_, err := i.WatchdogRunning()
	if err != nil {
		t.Error(err)
	}

}
func TestShutiffWatchdog(t *testing.T) {
	i := GetMockIPMI()
	if err := i.ShutoffWatchdog(); err != nil {
		t.Error(err)
	}
}
func TestGetDeviceID(t *testing.T) {
	i := GetMockIPMI()
	_, err := i.GetDeviceID()
	if err != nil {
		t.Error(err)
	}
}
func TestEnableSEL(t *testing.T) {
	i := GetMockIPMI()
	_, err := i.EnableSEL()
	if err != nil {
		t.Error(err)
	}
}

func TestGetSELInfo(t *testing.T) {
	i := GetMockIPMI()
	_, err := i.GetSELInfo()
	if err != nil {
		t.Error(err)
	}
}
func TestGetLanConfig(t *testing.T) {
	i := GetMockIPMI()
	_, err := i.GetLanConfig(1, 1)
	if err != nil {
		t.Error(err)
	}
}
func TestRawCmd(t *testing.T) {
	i := GetMockIPMI()
	data := []byte{0x1, 0x2, 0x3}
	_, err := i.RawCmd(data)
	if err != nil {
		t.Error(err)
	}
}
func TestSetSystemFWVersion(t *testing.T) {
	i := GetMockIPMI()
	if err := i.SetSystemFWVersion("TestTest"); err != nil {
		t.Error(err)
	}
}
func TestLogSystemEvent(t *testing.T) {
	i := GetMockIPMI()
	e := &Event{}
	if err := i.LogSystemEvent(e); err != nil {
		t.Error(err)
	}
}
