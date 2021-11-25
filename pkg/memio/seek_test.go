// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestSeek(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "seek_test")
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	memPath = tmpFile.Name()
	defer func() { memPath = "/dev/mem" }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := pathWrite(memPath, tt.addr, tt.writeData); err != nil {
				t.Errorf("pathwrite failed: %v", err)
			}
			if err := pathRead(memPath, tt.addr, tt.readData); err != nil {
				t.Errorf("pathRead failed: %v", err)
			}

			want := tt.writeData
			got := tt.readData
			if !reflect.DeepEqual(want, got) {
				t.Errorf("Write(%#016x, %v) = %v; want %v",
					tt.addr, want, got, want)
			}

		})
	}
}

func TestErrors(t *testing.T) {
	for _, tt := range testsInvalid {
		t.Run(tt.name, func(t *testing.T) {

			memPath = tt.path
			defer func() { memPath = "/dev/mem" }()

			if err := pathWrite(memPath, tt.addr, tt.writeData); err != nil {
				want := os.ErrNotExist
				if !errors.Is(err, want) {
					t.Errorf("Want %v, got %v", want, err)
				}
			}
			if err := pathRead(memPath, tt.addr, tt.readData); err != nil {
				want := os.ErrNotExist
				if !errors.Is(err, want) {
					t.Errorf("Want %v, got %v", want, err)
				}
			}
		})
	}
}
