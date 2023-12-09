// Copyright 2017-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"errors"
	"fmt"

	"github.com/u-root/u-root/pkg/crypto"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/vpd"
)

// Get, Set and GetAll are defined here as variables so they can be overridden
// for testing, or for using a key-value store other than VPD.
var (
	Get    = vpd.Get
	Set    = vpd.Set
	GetAll = vpd.GetAll
)

// BootEntry represents a boot entry, with its name, configuration, and Booter
// instance. It can map to existing key-value stores like VPD or EFI vars.
type BootEntry struct {
	Name   string
	Config []byte
	Booter Booter
}

var supportedBooterParsers = []func([]byte, ulog.Logger) (Booter, error){
	NewPxeBooter,
	NewBootBooter,
	NewNetBooter,
	NewLocalBooter,
}

var errNoBooterFound = errors.New("no booter found for entry")

// GetBooterFor looks for a supported Booter implementation and returns it, if
// found. If not found, error errNoBooterFound is returned.
func GetBooterFor(entry BootEntry, l ulog.Logger) (Booter, error) {
	var (
		booter Booter
		err    error
	)
	for idx, booterParser := range supportedBooterParsers {
		l.Printf("Trying booter #%d", idx)
		booter, err = booterParser(entry.Config, l)
		if err != nil {
			l.Printf("This config is not valid for this booter (#%d)", idx)
			l.Printf("  Error: %v", err.Error())
			continue
		}
		break
	}
	if booter == nil {
		return booter, fmt.Errorf("%w: %s: %s", errNoBooterFound, entry.Name, string(entry.Config))
	}
	return booter, nil
}

// GetBootEntries returns a list of BootEntry objects stored in the VPD
// partition of the flash chip
func GetBootEntries(l ulog.Logger) []BootEntry {
	var bootEntries []BootEntry

	for idx := 0; idx < 9999; idx++ {
		key := fmt.Sprintf("Boot%04d", idx)
		// try the RW entries first
		value, err := Get(key, false)
		if err == nil {
			crypto.TryMeasureData(crypto.NvramVarsPCR, value, key)
			bootEntries = append(bootEntries, BootEntry{Name: key, Config: value})
			// WARNING WARNING WARNING this means that read-write boot entries
			// have priority over read-only ones
			continue
		}
		// try the RO entries then
		value, err = Get(key, true)
		if err == nil {
			crypto.TryMeasureData(crypto.NvramVarsPCR, value, key)
			bootEntries = append(bootEntries, BootEntry{Name: key, Config: value})
		}
	}
	var err error
	// look for a Booter that supports the given configuration
	for idx, entry := range bootEntries {
		entry.Booter, err = GetBooterFor(entry, l)
		if err != nil {
			l.Printf("No booter found for entry: %s: %s", entry.Name, string(entry.Config))
		}
		bootEntries[idx] = entry
	}
	return bootEntries
}
