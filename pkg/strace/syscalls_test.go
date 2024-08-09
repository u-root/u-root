// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package strace

import (
	"fmt"
	"testing"

	"golang.org/x/sys/unix"
)

func TestByName(t *testing.T) {
	for _, tt := range []struct {
		name string
		val  uintptr
		ret  error
	}{
		{name: "read", val: unix.SYS_READ, ret: nil},
		{name: "Xread", val: unix.SYS_READ, ret: fmt.Errorf("Xread:not found")},
	} {
		n, err := ByName(tt.name)
		if err != nil && tt.ret == nil {
			t.Errorf("ByName(%s): %v != %v", tt.name, err, tt.ret)
		}
		if err == nil && tt.ret != nil {
			t.Errorf("ByName(%s): %v != %v", tt.name, err, tt.ret)
		}
		if err == nil && n != tt.val {
			t.Errorf("ByName(%s): %v != %v", tt.name, n, tt.val)
		}
	}
}

func TestByNum(t *testing.T) {
	for _, tt := range []struct {
		name string
		val  uintptr
		ret  error
	}{
		{name: "read", val: unix.SYS_READ, ret: nil},
		{name: "bogus", val: 10000000, ret: fmt.Errorf("Xread:not found")},
	} {
		n, err := ByNumber(tt.val)
		if err != nil && tt.ret == nil {
			t.Errorf("Bynumber(%d): %v != %v", tt.val, err, tt.ret)
		}
		if err == nil && tt.ret != nil {
			t.Errorf("Bynumber(%d): %v != %v", tt.val, err, tt.ret)
		}
		if err == nil && n != tt.name {
			t.Errorf("Bynumber(%d): %v != %v", tt.val, n, tt.name)
		}
	}
}
