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

const (
	_SECURITY_SUPPORTED     = 0x1
	_SECURITY_ENABLED       = 0x2
	_SECURITY_LOCKED        = 0x4
	_SECURITY_FROZEN        = 0x8
	_SECURITY_COUNT_EXPIRED = 0x10
	_SECURITY_LEVEL_MAX     = 0x100
)

// Info is information about a SCSI disk device.
type Info struct {
	NumberSectors           uint64
	ECCBytes                uint
	MasterPasswordRev       uint16
	SecurityStatus          DiskSecurityStatus
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

// DiskSecurityStatus is information about how the disk is secured.
type DiskSecurityStatus uint16

// SecuritySupported returns true if the disk has security.
func (d DiskSecurityStatus) SecuritySupported() bool {
	return (d & _SECURITY_SUPPORTED) != 0
}

// SecurityEnabled returns true if security is enabled on the disk.
func (d DiskSecurityStatus) SecurityEnabled() bool {
	return (d & _SECURITY_ENABLED) != 0
}

// SecurityLocked returns true if the disk is locked.
func (d DiskSecurityStatus) SecurityLocked() bool {
	return (d & _SECURITY_LOCKED) != 0
}

// SecurityFrozen returns true if the disk is frozen and its security state
// cannot be changed.
func (d DiskSecurityStatus) SecurityFrozen() bool {
	return (d & _SECURITY_FROZEN) != 0
}

// SecurityCountExpired returns true if all attempts to unlock the disk have
// been used up.
func (d DiskSecurityStatus) SecurityCountExpired() bool {
	return (d & _SECURITY_COUNT_EXPIRED) != 0
}

func (d DiskSecurityStatus) String() string {
	return fmt.Sprintf(`
	Security Status:
		Supported: %d,
		Enabled: %d,
		Locked: %d,
		Frozen: %d,
		Count Expired: %d,
		Level Max: %d
	`, d&_SECURITY_SUPPORTED, d&_SECURITY_ENABLED, d&_SECURITY_LOCKED,
		d&_SECURITY_FROZEN, d&_SECURITY_COUNT_EXPIRED, d&_SECURITY_LEVEL_MAX)
}

// Disk is the interface to a disk, with operations to create packets and
// operate on them.
type Disk interface {
	// Unlock unlocks the drive, given a password and an indication of
	// whether it is the admin (true) or user (false) password.
	Unlock(password string, admin bool) error

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
