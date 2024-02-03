// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securelaunch

import (
	"path/filepath"
	"testing"
)

func TestAddToPersistQueue(t *testing.T) {
	desc := "test"

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "out")

	dataStr := "Hello World!"

	if err := AddToPersistQueue(desc, []byte(dataStr), tempFile, ""); err != nil {
		t.Fatalf(`AddToPersistQueue(desc, []byte(dataStr), tempFile, "") = %v, not nil`, err)
	}
}

func TestGetBlkInfo(t *testing.T) {
	if err := GetBlkInfo(); err != nil {
		t.Fatalf("GetBlkInfo() = %v, not nil", err)
	}
}

func TestGetStorageDevice(t *testing.T) {
	if err := GetBlkInfo(); err != nil {
		t.Fatalf("GetBlkInfo() = %v, not nil", err)
	}

	if len(StorageBlkDevices) == 0 {
		t.Fatal("len(StorageBlockDevices) = 0, not > 0")
	}

	deviceName := StorageBlkDevices[0].Name

	device, err := GetStorageDevice(deviceName)
	if err != nil {
		t.Fatalf("GetStorageDevice(deviceName) = %v, not nil", err)
	}
	if device.Name != deviceName {
		t.Fatalf("GetStorageDevice(deviceName) = %q, not %q", device.Name, deviceName)
	}
}
