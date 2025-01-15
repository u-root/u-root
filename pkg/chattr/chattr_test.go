// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chattr

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"unsafe"

	"github.com/hugelgupf/vmtest/guest"
	"golang.org/x/sys/unix"
)

func TestSetAttr(t *testing.T) {
	guest.SkipIfNotInVM(t)

	// Create a temporary file
	file, err := ioutil.TempFile("", "chattr_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Test setting immutable flag
	err = SetAttr(file, "+i")
	if err != nil {
		t.Errorf("Error setting immutable flag: %v", err)
	}
	attr, err := GetAttr(file)
	if err != nil {
		t.Errorf("Error getting attributes: %v", err)
	}

	if attr&FS_IMMUTABLE_FL == 0 {
		t.Error("Immutable flag not set")
	}

	// Test unsetting immutable flag
	err = SetAttr(file, "-i")
	if err != nil {
		t.Errorf("Error unsetting immutable flag: %v", err)
	}
	attr, err = GetAttr(file)
	if err != nil {
		t.Errorf("Error getting attributes: %v", err)
	}
	if attr&FS_IMMUTABLE_FL != 0 {
		t.Error("Immutable flag not unset")
	}

	// Test setting append-only flag
	err = SetAttr(file, "+a")
	if err != nil {
		t.Errorf("Error setting append-only flag: %v", err)
	}
	attr, err = GetAttr(file)
	if err != nil {
		t.Errorf("Error getting attributes: %v", err)
	}
	if attr&FS_APPEND_FL == 0 {
		t.Error("Append-only flag not set")
	}

	// Test unsetting append-only flag
	err = SetAttr(file, "-a")
	if err != nil {
		t.Errorf("Error unsetting append-only flag: %v", err)
	}
	attr, err = GetAttr(file)
	if err != nil {
		t.Errorf("Error getting attributes: %v", err)
	}
	if attr&FS_APPEND_FL != 0 {
		t.Error("Append-only flag not unset")
	}
}

func TestGetAttr(t *testing.T) {
	guest.SkipIfNotInVM(t)
	// Create a temporary file
	file, err := ioutil.TempFile("", "chattr_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Get initial attributes
	attr, err := GetAttr(file)
	if err != nil {
		t.Errorf("Error getting attributes: %v", err)
	}
	log.Printf("Chattr test: current attr %v", attr)
	// Set immutable flag using unix.IoctlSetInt directly
	attr = attr | FS_IMMUTABLE_FL
	ptr := (uintptr)(unsafe.Pointer(&attr))
	err = unix.IoctlSetInt(int(file.Fd()), unix.FS_IOC_SETFLAGS, int(ptr))
	if err != nil {
		t.Errorf("Error setting immutable flag using ioctl: %v", err)
	}

	// Get attributes again and check if immutable flag is set
	attr, err = GetAttr(file)
	if err != nil {
		t.Errorf("Error getting attributes: %v", err)
	}
	if attr&FS_IMMUTABLE_FL == 0 {
		t.Error("Immutable flag not set after setting it with ioctl")
	}
}
