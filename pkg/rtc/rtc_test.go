// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtc

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
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
				// As f.Name() is now the only object in the devs array
				// we can close f here an don't need to run a defer on
				// os.Remove here but rather do that manually a few lines
				// further down
				f.Close()
				rtc, err := OpenRTC()
				if err != nil {
					t.Error(err)
				}
				rtc.Close()
				if err := os.Remove(devs[0]); err != nil {
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
				//TODO(MDr164): This only works locally, no idea why
				testutil.SkipIfInVMTest(t)
				f, err := ioutil.TempFile("", "rtc-*")
				if err != nil {
					t.Error(err)
				}
				devs = []string{f.Name()}
				// As f.Name() is now the only object in the devs array
				// we can close f here an don't need to run a defer on
				// os.Remove here but rather do that manually a few lines
				// further down
				if err := f.Chmod(0); err != nil {
					t.Error(err)
				}
				f.Close()
				_, err = OpenRTC()
				if err == nil {
					t.Error(err)
				}
				if err := os.Chmod(devs[0], 0777); err != nil {
					t.Error(err)
				}
				if err := os.Remove(devs[0]); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
