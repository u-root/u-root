// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/uio/uio"
)

type ArchiveValidator interface {
	Validate(a *cpio.Archive) error
}

type HasRecord struct {
	R cpio.Record
}

func (hr HasRecord) Validate(a *cpio.Archive) error {
	r, ok := a.Get(hr.R.Name)
	if !ok {
		return fmt.Errorf("archive does not contain %v", hr.R)
	}
	if !cpio.Equal(r, hr.R) {
		return fmt.Errorf("archive does not contain %v; instead has %v", hr.R, r)
	}
	return nil
}

type HasFile struct {
	Path string
}

func (hf HasFile) Validate(a *cpio.Archive) error {
	if _, ok := a.Get(hf.Path); ok {
		return nil
	}
	return fmt.Errorf("archive does not contain %s, but should", hf.Path)
}

type HasContent struct {
	Path    string
	Content string
}

func (hc HasContent) Validate(a *cpio.Archive) error {
	r, ok := a.Get(hc.Path)
	if !ok {
		return fmt.Errorf("archive does not contain %s, but should", hc.Path)
	}
	if c, err := uio.ReadAll(r); err != nil {
		return fmt.Errorf("reading record %s failed: %w", hc.Path, err)
	} else if string(c) != hc.Content {
		return fmt.Errorf("content of %s is %s, want %s", hc.Path, string(c), hc.Content)
	}
	return nil
}

type MissingFile struct {
	Path string
}

func (mf MissingFile) Validate(a *cpio.Archive) error {
	if _, ok := a.Get(mf.Path); ok {
		return fmt.Errorf("archive contains %s, but shouldn't", mf.Path)
	}
	return nil
}

type IsEmpty struct{}

func (IsEmpty) Validate(a *cpio.Archive) error {
	if empty := a.Empty(); !empty {
		return fmt.Errorf("expected archive to be empty")
	}
	return nil
}

func ReadArchive(path string) (*cpio.Archive, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return cpio.ArchiveFromReader(cpio.Newc.Reader(f))
}
