// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/u-root/u-root/pkg/pci"
)

func TestRun(t *testing.T) {
	// Cover the switch case
	log.SetOutput(io.Discard)
	for _, tt := range []struct {
		name string
		args []string
	}{
		{
			name: "switch hexdump case 1",
			args: []string{"-x", "1"},
		},
		{
			name: "switch hexdump case 2",
			args: []string{"-x", "2"},
		},
		{
			name: "switch hexdump case 3",
			args: []string{"-x", "3"},
		},
		{
			name: "switch hexdump case 4",
			args: []string{"-x", "4"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := command(io.Discard, tt.args...)
			c.run()
		})
	}
	// Cover the rest
	for _, tt := range []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "readJSON true, without error",
			args: []string{"-J", "testdata/testfile1.json"},
		},
		{
			name: "readJSON true, error in os.ReadFile",
			args: []string{"-J", "testdata/testfilex.json"},
			err:  os.ErrNotExist,
		},
		{
			name: "readJSON true, error in json.Unmarshal",
			args: []string{"-J", "testdata/testfile2.json"},
			err:  errBadJSON,
		},
		{
			name: "dumpJSON",
			args: []string{"-J", "testdata/testfile1.json", "-j"},
		},
		{
			name: "invoke registers",
			args: []string{"examplearg"},
			err:  strconv.ErrSyntax,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := command(io.Discard, tt.args...)
			if err := c.run(); !errors.Is(err, tt.err) {
				t.Errorf("run() got %v, want %v", err, tt.err)
			}
		})
	}
}

func TestRegisters(t *testing.T) {
	configBytes := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "config"))
	if err != nil {
		t.Errorf("Creating file failed: %v", err)
	}
	_, err = f.Write(configBytes)
	if err != nil {
		t.Errorf("Writing to file failed: %v", err)
	}
	for _, tt := range []struct {
		name    string
		devices pci.Devices
		cmds    []string
		err     error
	}{
		{
			name: "trigger first log.Printf",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"cmd=cmd=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "trigger second log.Printf",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"c.m.d=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "error in first strconv.ParseUint satisfying l case",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"cmd.l=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "error in first strconv.ParseUint satisfying w case",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"cmd.w=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "error in first strconv.ParseUint satisfying b case",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"cmd.b=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "triggers Bad size log and satisfies the justCheck check",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"cmd.cmd=cmd", "cmd.b=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "triggers error, reading out of bounds, EOF",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"10.b"},
			err:  io.EOF,
		},
		{
			name: "reading works and write in PCI.ExtraInfo",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"0.w"},
		},
		{
			name: "error in second strconv.ParseUint",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"0.w=cmd"},
			err:  strconv.ErrSyntax,
		},
		{
			name: "writing successful",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds: []string{"0.w=10"},
		},
		{
			name: "writing failes because config file does not exist",
			devices: []*pci.PCI{
				{
					Config: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77},
				},
			},
			cmds: []string{"0.w=10"},
			err:  os.ErrNotExist,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			log.SetOutput(buf)
			if err := registers(tt.devices, tt.cmds...); !errors.Is(err, tt.err) {
				t.Errorf("registers(): got %v, want %v", err, tt.err)
			}
		})
	}
}

// This test is here because of very strange encoding/json behavior.
// Right now it works.
// It stopped working with the JSON file that has been in use for 8 years.
// How to see it yourself.
// It seems the json package makes an effort to parse "Primary": "00" to a Uint8, but some usages
// of the package cause a failure? it's very mysterious.
//
// go test, see it pass.
// go test -test.run TestJSON
// PASS
// ok  	github.com/u-root/u-root/cmds/core/pci	0.002s
// in the JSON file, change the JSON for Primary:
// +               "Primary": "00",
// -               "Primary": 55,
// Watch it fail, even though it has worked for years.
// --- FAIL: TestJSONWeirdness (0.00s)
//
//	pci_test.go:240: json: cannot unmarshal string into Go struct field PCI.Primary of type uint8
//	pci_test.go:243: json: cannot unmarshal string into Go struct field PCI.Primary of type uint8
//	pci_test.go:247: run() got json: cannot unmarshal string into Go struct field PCI.Primary of type uint8:JSON parsing failed, want nil
//
// Change the true to false on the if. So the code is removed at link time.
// rminnich@pop-os:~/go/src/github.com/u-root/u-root/cmds/core/pci$ go test -test.run TestJSONWeirdness
// PASS
// ok  	github.com/u-root/u-root/cmds/core/pci	0.002s
// NOTE, the JSON did not change between the bad and good run; only the way in which you called
// Unmarshal.
// So, just enabling those lines cause the PREVIOUS Unmarshal's to fail! Precrime!
// tinygo and go see the same issue.
// It's very hard to get this to repro, much less make a standalone repro, but hopefully this
// test will help protect us in future.
// And, the weirdest part:
func TestJSONWeirdness(t *testing.T) {
	var d pci.Devices
	b, err := os.ReadFile("testdata/testfile1.json")
	if err != nil {
		t.Log(err)
	}
	if err := json.Unmarshal(b, &d); err != nil {
		t.Log(err)
	}
	if err := json.Unmarshal(b, &d); err != nil {
		t.Log(err)
	}
	if true {
		c := command(io.Discard, "-J", "testdata/testfile1.json")
		if err := c.run(); err != nil {
			t.Errorf("run() got %v, want nil", err)
		}
	}
}
