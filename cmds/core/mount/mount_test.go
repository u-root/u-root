// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"bytes"
	"errors"
	"os"
	"slices"
	"sort"
	"testing"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	procMounts, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	procMountContent := `sysfs /sys sysfs rw,nosuid,nodev,noexec,relatime 0 0
proc /proc proc rw,nosuid,nodev,noexec,relatime 0 0
`
	_, err = procMounts.WriteString(procMountContent)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedErr    error
		name           string
		fsType         string
		fileSystemPath string
		expectedStdout string
		opts           mountOptions
		args           []string
		mountPath      []string
		ro             bool
	}{
		{
			name:        "wrong usage",
			args:        []string{"arg1"},
			expectedErr: errUsage,
		},
		{
			name:        "no mount path error",
			expectedErr: errMountPath,
		},
		{
			name:           "no args",
			mountPath:      []string{procMounts.Name()},
			expectedStdout: procMountContent,
		},
		{
			name:        "mount not exist",
			args:        []string{"/errNotExitPath", "/mount/somewhere"},
			ro:          true,
			expectedErr: os.ErrNotExist,
		},
	}

	for _, test := range tests {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		cmd := command(stdout, stderr, test.ro, test.fsType, test.opts)
		cmd.mountsPath = test.mountPath
		cmd.fileSystemsPath = test.fileSystemPath

		err := cmd.run(test.args...)
		if !errors.Is(err, test.expectedErr) {
			t.Fatalf("expected %v go %v", test.expectedErr, err)
		}

		if test.expectedStdout != stdout.String() {
			t.Errorf("expected %v got %v", test.expectedStdout, stdout.String())
		}
	}
}

func TestGetSupportedFilesystem(t *testing.T) {
	dir := t.TempDir()
	pfs, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	procFSContent := `nodev	sysfs
nodev	tmpfs
nodev	bdev
nodev	proc
nodev	cgroup`

	_, err = pfs.WriteString(procFSContent)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		originFS      string
		fsPath        string
		expectedErr   error
		expectedList  []string
		expectedKnown bool
	}{
		{
			originFS:      "tmpfs",
			fsPath:        pfs.Name(),
			expectedKnown: true,
			expectedList:  []string{"sysfs", "tmpfs", "bdev", "proc", "cgroup"},
		},
		{
			originFS:      "tmpfs",
			fsPath:        "filenotfound",
			expectedKnown: false,
			expectedErr:   os.ErrNotExist,
		},
	}

	for _, test := range tests {
		c := &cmd{fileSystemsPath: test.fsPath}
		xs, known, err := c.getSupportedFilesystem(test.originFS)
		if !errors.Is(err, test.expectedErr) {
			t.Fatalf("expected error %v got %v", test.expectedErr, err)
		}
		if known != test.expectedKnown {
			t.Errorf("expected known %t got %t", test.expectedKnown, known)
		}

		sort.Strings(xs)
		sort.Strings(test.expectedList)
		if !slices.Equal(xs, test.expectedList) {
			t.Errorf("expected %v, got %v", test.expectedList, xs)
		}
	}
}
