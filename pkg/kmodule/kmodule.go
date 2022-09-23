// Copyright 2017-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"errors"
	"io"
	"os"
)

// Loader is the interface for module setup, loading, and listing.
type Loader interface {
	Init(image []byte, opts string) error
	FileInit(f *os.File, opts string, flags uintptr) error
	Probe(name string, modParams string) error
	Delete(name string, flags uintptr) error
	Read([]byte) (int, error)
}

var (
	// ErrNoLoader means there is no kmodule loader
	// for this kernel.
	ErrNoLoader = errors.New("no kmodule loader")
	// ErrNotSupported means the requested operation
	// is not supported.
	ErrNotSupported = errors.New("not supported")
)

// List lists the modules in use by the kernel.
func List(l Loader, w io.Writer) error {
	b, err := io.ReadAll(l)
	if err != nil {
		return err
	}
	return Pretty(w, string(b))
}
