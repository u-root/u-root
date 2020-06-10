// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

import (
	"encoding/json"
	"fmt"
	"time"
)

// DefaultTimeout is the default timeout for disk operations.
const DefaultTimeout time.Duration = 15 * time.Second

// Info is information about a SCSI disk device.
type Info struct {
	NumberSectors           uint64
	ECCBytes                uint
	MasterPasswordRev       uint16
	SecurityStatus          uint16
	TrustedComputingSupport uint16

	Serial           string
	Model            string
	FirmwareRevision string

	// These are the pair-byte-swapped space-padded versions of serial,
	// model, and firmware revision as originally returned by the SCSI
	// device.
	OrigSerial           string
	OrigModel            string
	OrigFirmwareRevision string
}

// Disk is the interface to a disk, with operations to create packets and
// operate on them.
type Disk interface {
	// Unlock unlocks the drive, given a password and an indication of
	// whether it is the master (true) or user (false) password.
	Unlock(password string, master bool) error

	// Identify returns drive identity information
	Identify() (*Info, error)
}

// String prints a nice JSON-formatted info.
func (i *Info) String() string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(s)
}
