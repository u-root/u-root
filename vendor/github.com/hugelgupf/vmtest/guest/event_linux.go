// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package guest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ports = "/sys/class/virtio-ports"

// VirtioSerialDevice looks up the device path for the given virtio-serial
// name.
//
// The name would be configured in the QEMU command-line (or e.g. with
// qemu.EventChannel).
func VirtioSerialDevice(name string) (string, error) {
	entries, err := os.ReadDir(ports)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		entryName, err := os.ReadFile(filepath.Join(ports, entry.Name(), "name"))
		if err != nil {
			continue
		}
		if strings.TrimRight(string(entryName), "\n") == name {
			return filepath.Join("/dev", entry.Name()), nil
		}
	}
	return "", fmt.Errorf("no virtio-serial device with name %s", name)
}

// SerialEventChannel opens an event channel to the host over virtio-serial
// with the given virtio-serial port name.
//
// Callers must call Close on Emitter to publish a final "done" event to signal
// the host no more events are coming. If the "done" event is not published,
// qemu.EventChannel is configured to return an error on VM exit on the host.
//
// T should be the type of a JSON event being sent, matching the host
// configuration on qemu.EventChannel reading from this channel.
//
// The name should match the qemu.EventChannel configuration on the host as
// well.
func SerialEventChannel[T any](name string) (*Emitter[T], error) {
	dev, err := VirtioSerialDevice(name)
	if err != nil {
		return nil, err
	}
	return EventChannel[T](dev)
}
