// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/pci"
)

func TestPCIExecution(t *testing.T) {
	// Cover the switch case
	for _, tt := range []struct {
		name    string
		hexdump int
	}{
		{
			name:    "switch hexdump case 1",
			hexdump: 1,
		},
		{
			name:    "switch hexdump case 2",
			hexdump: 2,
		},
		{
			name:    "switch hexdump case 3",
			hexdump: 3,
		},
		{
			name:    "switch hexdump case 4",
			hexdump: 4,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*hexdump = tt.hexdump
			pciExecution(io.Discard, []string{}...)
		})
	}
	// Cover the rest
	for _, tt := range []struct {
		name      string
		args      []string
		numbers   bool
		devs      string
		dumpJSON  bool
		verbosity int
		readJSON  string
		wantErr   string
	}{
		{
			name:     "readJSON true, without error",
			readJSON: "testdata/testfile1.json",
		},
		{
			name:     "readJSON true, error in os.ReadFile",
			readJSON: "testdata/testfile.json",
			wantErr:  "no such file or directory",
		},
		{
			name:     "readJSON true, error in json.Unmarshal",
			readJSON: "testdata/testfile2.json",
			wantErr:  "unexpected end of JSON input",
		},
		{
			name:     "dumpJSON",
			readJSON: "testdata/testfile1.json",
			dumpJSON: true,
		},
		{
			name: "invoke registers",
			args: []string{"examplearg"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*numbers = tt.numbers
			*devs = tt.devs
			*dumpJSON = tt.dumpJSON
			*verbosity = tt.verbosity
			*readJSON = tt.readJSON
			if got := pciExecution(io.Discard, tt.args...); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("pciExecution() = %q, should contain: %q", got, tt.wantErr)
				}
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
		wantErr string
	}{
		{
			name: "trigger first log.Printf",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"cmd=cmd=cmd"},
			wantErr: "only one = allowed",
		},
		{
			name: "trigger second log.Printf",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"c.m.d=cmd"},
			wantErr: "only one . allowed",
		},
		{
			name: "error in first strconv.ParseUint satisfying l case",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"cmd.l=cmd"},
			wantErr: "parsing \"cmd\": invalid syntax",
		},
		{
			name: "error in first strconv.ParseUint satisfying w case",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"cmd.w=cmd"},
			wantErr: "parsing \"cmd\": invalid syntax",
		},
		{
			name: "error in first strconv.ParseUint satisfying b case",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"cmd.b=cmd"},
			wantErr: "parsing \"cmd\": invalid syntax",
		},
		{
			name: "triggers Bad size log and satisfies the justCheck check",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"cmd.cmd=cmd", "cmd.b=cmd"},
			wantErr: "Bad size",
		},
		{
			name: "triggers error, reading out of bounce, EOF",
			devices: []*pci.PCI{
				{
					FullPath: dir,
				},
			},
			cmds:    []string{"10.b"},
			wantErr: "EOF",
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
			cmds:    []string{"0.w=cmd"},
			wantErr: "parsing \"cmd\": invalid syntax",
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
			cmds:    []string{"0.w=10"},
			wantErr: "open config: no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			log.SetOutput(buf)
			registers(tt.devices, tt.cmds...)
			if !strings.Contains(buf.String(), tt.wantErr) {
				t.Errorf("registers() = %q, should contain: %q", buf.String(), tt.wantErr)
			}
		})
	}
}
