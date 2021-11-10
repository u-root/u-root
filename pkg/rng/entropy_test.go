// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rng

import (
	"os"
	"testing"
)

func TestSetAvailableTRNG(t *testing.T) {
	for _, tt := range []struct {
		name  string
		setup func(*testing.T)
	}{
		{
			name: "set trng",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Fail()
				}
				_, err = f.WriteString("tpm-rng")
				if err != nil {
					t.Fail()
				}
				f.Close()
			},
		},
		{
			name: "no rng available",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Fail()
				}
				_, err = f.WriteString("none")
				if err != nil {
					t.Fail()
				}
				f.Close()
			},
		},
		{
			name:  "no available file",
			setup: func(t *testing.T) {},
		},
		{
			name: "no write access",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Fail()
				}
				_, err = f.WriteString("tpm-rng")
				if err != nil {
					t.Fail()
				}
				f.Close()
				f, err = os.Create(HwRandomCurrentFile)
				if err != nil {
					t.Fail()
				}
				f.Close()
			},
		},
		{
			name: "no read access",
			setup: func(t *testing.T) {
				f, err := os.Create(HwRandomAvailableFile)
				if err != nil {
					t.Fail()
				}
				_, err = f.WriteString("tpm-rng")
				if err != nil {
					t.Fail()
				}
				f.Close()
				f, err = os.Create(HwRandomCurrentFile)
				if err != nil {
					t.Fail()
				}
				f.Close()
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			HwRandomAvailableFile = "/tmp/testdata/rng_available"
			HwRandomCurrentFile = "/tmp/testdata/rng_current"
			err := os.MkdirAll("/tmp/testdata", 0775)
			if err != nil {
				t.Fail()
			}
			tt.setup(t)

			if tt.name == "set trng" {
				err = setAvailableTRNG()
				if err != nil {
					t.Fail()
				}
			}
			if tt.name == "no rng available" {
				err = setAvailableTRNG()
				if err == nil {
					t.Fail()
				}
			}
			if tt.name == "no available file" {
				err = setAvailableTRNG()
				if err == nil {
					t.Fail()
				}
			}
			if tt.name == "no write access" {
				err = os.Chmod(HwRandomCurrentFile, 0)
				if err != nil {
					t.Fail()
				}
				err = setAvailableTRNG()
				if err == nil {
					t.Fail()
				}
				err = os.Chmod(HwRandomCurrentFile, 0755)
				if err != nil {
					t.Fail()
				}
			}
			if tt.name == "no write access" {
				err = os.Chmod(HwRandomCurrentFile, 0333)
				if err != nil {
					t.Fail()
				}
				err = setAvailableTRNG()
				if err == nil {
					t.Fail()
				}
				err = os.Chmod(HwRandomCurrentFile, 0755)
				if err != nil {
					t.Fail()
				}
			}
			err = os.RemoveAll("/tmp/testdata")
			if err != nil {
				t.Fail()
			}
		})
	}
}

func TestUpdateLinuxRandomness(t *testing.T) {

}
