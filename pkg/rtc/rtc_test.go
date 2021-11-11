// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtc

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestOpenRTC(t *testing.T) {
	for _, tt := range []struct {
		name string
	}{
		{
			name: "open rtc",
		},
		{
			name: "no rtc",
		},
		{
			name: "no permission",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "open rtc" {
				f, err := ioutil.TempFile("", "rtc-*")
				if err != nil {
					t.Error(err)
				}
				devs = []string{f.Name()}
				f.Close()
				rtc, err := OpenRTC()
				if err != nil {
					t.Error(err)
				}
				rtc.Close()
				err = os.Remove(devs[0])
				if err != nil {
					t.Error(err)
				}
			}
			if tt.name == "no rtc" {
				devs = []string{"bogusfile"}
				_, err := OpenRTC()
				if err == nil {
					t.Error(err)
				}
			}
			if tt.name == "no permission" {
				f, err := ioutil.TempFile("", "rtc-*")
				if err != nil {
					t.Error(err)
				}
				devs = []string{f.Name()}
				err = f.Chmod(0)
				if err != nil {
					t.Error(err)
				}
				f.Close()
				_, err = OpenRTC()
				if err == nil {
					t.Error(err)
				}
				err = os.Chmod(devs[0], 0777)
				if err != nil {
					t.Error(err)
				}
				err = os.Remove(devs[0])
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
