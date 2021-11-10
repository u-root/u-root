// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestSeek(t *testing.T) {
	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "seek_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			memPath = tmpFile.Name()
			defer func() { memPath = "/dev/mem" }()

			if err := pathWrite(memPath, tt.addr, tt.writeData); err != nil {
				t.Fatal(err)
			}
			if err := pathRead(memPath, tt.addr, tt.readData); err != nil {
				t.Fatal(err)
			}

			want := tt.writeData
			got := tt.readData
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("Write(%#016x, %v) = %v; want %v",
					tt.addr, want, got, want)
			}

		})
	}
}

func TestErrors(t *testing.T) {
	for _, tt := range testsInvalid {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {

			memPath = tt.path
			defer func() { memPath = "/dev/mem" }()

			if err := pathWrite(memPath, tt.addr, tt.writeData); err != nil {
				want := os.ErrNotExist
				if !errors.Is(err, want) {
					t.Fatal(err)
				}
			}
			if err := pathRead(memPath, tt.addr, tt.readData); err != nil {
				want := os.ErrNotExist
				if !errors.Is(err, want) {
					t.Fatal(err)
				}
			}
		})
	}
}
