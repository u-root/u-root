// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package kmodule interfaces with Linux kernel modules.
//
// kmodule allows loading and unloading kernel modules with dependencies, as
// well as locating them through probing.
package kmodule

import (
	"io"
	"os"
	"os/exec"
)

// DarwinLoader holds options and information for Darwin modules.
type DarwinLoader struct {
}

// Init loads the kernel module given by image with the given options.
func (l *DarwinLoader) Init(_ []byte, _ string) error {
	return ErrNotSupported
}

// Read implements io.Reader.
func (l *DarwinLoader) Read(b []byte) (int, error) {
	o, err := exec.Command("kmutil", "showloaded").CombinedOutput()
	if err != nil {
		return -1, err
	}
	copy(b, o)
	return len(o), err
}

// FileInit returns ErrNotSupported.
func (l *DarwinLoader) FileInit(_ *os.File, _ string, _ uintptr) error {
	return ErrNotSupported
}

// Delete returns ErrNotSupported.
func (l *DarwinLoader) Delete(_ string, _ uintptr) error {
	return ErrNotSupported
}

// New returns a new DarwinLoader given a path
func NewPath(_ string) (*DarwinLoader, error) {
	return &DarwinLoader{}, nil
}

// New returns a new DarwinLoader using a default path.
func New() (*DarwinLoader, error) {
	return &DarwinLoader{}, nil
}

var _ Loader = &DarwinLoader{}

// ProbeOpts are options for probing modules.
type ProbeOpts struct {
}

// Probe loads the given kernel module and its dependencies.
// It is calls ProbeOptions with the default ProbeOpts.
func (l *DarwinLoader) Probe(name string, modParams string) error {
	return l.ProbeOptions(name, modParams, ProbeOpts{})
}

// ProbeOptions returns ErrNotSupported.
func (l *DarwinLoader) ProbeOptions(_, _ string, _ ProbeOpts) error {
	return ErrNotSupported
}

// Pretty prints the string from /proc/modules in a pretty format.
func Pretty(w io.Writer, s string) error {
	_, err := w.Write([]byte(s))
	return err
}
