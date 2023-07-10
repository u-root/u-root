// Copyright 2017-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"bytes"
	"errors"
	"testing"

	"github.com/u-root/u-root/pkg/ulog"
)

func TestGetBooterForNetBooter(t *testing.T) {
	validConfig := BootEntry{
		Name:   "Boot0000",
		Config: []byte(`{"type": "netboot", "method": "dhcpv6", "mac": "aa:bb:cc:dd:ee:ff"}`),
	}
	booter, err := GetBooterFor(validConfig, ulog.Null)
	if err != nil {
		t.Fatalf(`GetBooterFor(validConfig) error = %v`, err)
	}
	if booter.TypeName() != "netboot" {
		t.Errorf(`GetBooterFor(validConfig).TypeName() = %q, want "netboot"`, booter.TypeName())
	}
	if booter.(*NetBooter) == nil {
		t.Errorf(`booter.(*NetBooter) = %v, want not nil`, booter.(*NetBooter))
	}
}

func TestGetBooterForNullBooter(t *testing.T) {
	validConfig := BootEntry{
		Name:   "Boot0000",
		Config: []byte(`{"type": "null"}`),
	}
	booter, err := GetBooterFor(validConfig, ulog.Null)
	if !errors.Is(err, errNoBooterFound) {
		t.Fatalf(`GetBooterFor(validConfig) error = %v, wantErr %v`, err, errNoBooterFound)
	}
	if booter != nil {
		t.Fatalf(`GetBooterFor(validConfig) = %v, want nil`, booter)
	}
}

func TestGetBooterForInvalidBooter(t *testing.T) {
	invalidConfig := BootEntry{
		Name:   "Boot0000",
		Config: []byte(`{"type": "invalid"`),
	}
	booter, err := GetBooterFor(invalidConfig, ulog.Null)
	if !errors.Is(err, errNoBooterFound) {
		t.Fatalf(`GetBooterFor(invalidConfig) error = %v, wantErr %v`, err, errNoBooterFound)
	}
	if booter != nil {
		t.Fatalf(`GetBooterFor(invalidConfig) = %v, want nil`, booter)
	}
}

func TestGetBootEntries(t *testing.T) {
	var (
		bootConfig0000 = []byte(`{"type": "netboot", "method": "dhcpv6", "mac": "aa:bb:cc:dd:ee:ff"}`)
		bootConfig0001 = []byte(`{"type": "localboot", "uuid": "blah-bleh", "kernel": "/path/to/kernel"}`)
	)
	// Override the package-level variable Get so it will use our test getter
	// instead of VPD
	Get = func(key string, readOnly bool) ([]byte, error) {
		switch key {
		case "Boot0000":
			return bootConfig0000, nil
		case "Boot0001":
			return bootConfig0001, nil
		default:
			return nil, errors.New("No such key")
		}
	}
	entries := GetBootEntries(ulog.Null)
	if len(entries) != 2 {
		t.Errorf(`len(entries) = %d, want "2"`, len(entries))
	}
	if entries[0].Name != "Boot0000" {
		t.Errorf(`entries[0].Name = %q, want "Boot0000"`, entries[0].Name)
	}
	if !bytes.Equal(entries[0].Config, bootConfig0000) {
		t.Errorf(`entries[0].Config = %v, want %v`, entries[0].Config, bootConfig0000)
	}
	if entries[1].Name != "Boot0001" {
		t.Errorf(`entries[1].Name = %q, want "Boot0001"`, entries[1].Name)
	}
	if !bytes.Equal(entries[1].Config, bootConfig0001) {
		t.Errorf(`entries[1].Config = %v, want %v`, entries[1].Config, bootConfig0001)
	}
}

func TestGetBootEntriesOnlyRO(t *testing.T) {
	// Override the package-level variable Get so it will use our test getter
	// instead of VPD
	Get = func(key string, readOnly bool) ([]byte, error) {
		if !readOnly || key != "Boot0000" {
			return nil, errors.New("No such key")
		}
		return []byte(`{"type": "netboot", "method": "dhcpv6", "mac": "aa:bb:cc:dd:ee:ff"}`), nil
	}
	entries := GetBootEntries(ulog.Null)
	if len(entries) != 1 {
		t.Errorf(`len(entries) = %d, want "1"`, len(entries))
	}
}
