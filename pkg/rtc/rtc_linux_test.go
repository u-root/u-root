// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtc

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestSet(t *testing.T) {
	for _, tt := range []struct {
		name string
	}{
		{
			name: "success",
		},
		{
			name: "error",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ioutil.TempFile("", "rtc-*")
			if err != nil {
				t.Error(err)
			}
			devs = []string{f.Name()}
			// As f.Name() is now the only object in the devs array
			// we can close f here an don't need to run a defer on
			// os.Remove here but rather do that manually a few lines
			// further down
			f.Close()
			rtc, err := OpenRTC()
			if err != nil {
				t.Error(err)
			}
			if tt.name == "success" {
				unixIoctlSetRTCTime = func(fd int, value *unix.RTCTime) error {
					return nil
				}
				if err := rtc.Set(time.Now()); err != nil {
					t.Error(err)
				}
			}
			if tt.name == "error" {
				unixIoctlSetRTCTime = func(fd int, value *unix.RTCTime) error {
					return errors.New("error")
				}
				if err := rtc.Set(time.Now()); err == nil {
					t.Error(err)
				}
			}
			rtc.Close()
			if err := os.Remove(devs[0]); err != nil {
				t.Error(err)
			}
			unixIoctlSetRTCTime = unix.IoctlSetRTCTime
		})
	}
}

func TestRead(t *testing.T) {
	for _, tt := range []struct {
		name string
	}{
		{
			name: "success",
		},
		{
			name: "error",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ioutil.TempFile("", "rtc-*")
			if err != nil {
				t.Error(err)
			}
			devs = []string{f.Name()}
			// As f.Name() is now the only object in the devs array
			// we can close f here an don't need to run a defer on
			// os.Remove here but rather do that manually a few lines
			// further down
			f.Close()
			rtc, err := OpenRTC()
			if err != nil {
				t.Error(err)
			}
			if tt.name == "success" {
				unixIoctlGetRTCTime = func(fd int) (*unix.RTCTime, error) {
					return &unix.RTCTime{}, nil
				}
				_, err = rtc.Read()
				if err != nil {
					t.Error(err)
				}
			}
			if tt.name == "error" {
				unixIoctlGetRTCTime = func(fd int) (*unix.RTCTime, error) {
					return nil, errors.New("error")
				}
				_, err = rtc.Read()
				if err == nil {
					t.Error(err)
				}
			}
			rtc.Close()
			if err := os.Remove(devs[0]); err != nil {
				t.Error(err)
			}
			unixIoctlGetRTCTime = unix.IoctlGetRTCTime
		})
	}
}
