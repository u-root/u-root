// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"errors"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
)

func imageEqual(li1, li2 *LinuxImage) bool {
	return uio.ReaderAtEqual(li1.Kernel, li2.Kernel) &&
		uio.ReaderAtEqual(li1.Kernel, li2.Kernel) &&
		li1.Cmdline == li2.Cmdline
}

var errSkip = errors.New("foo")

type errorReaderAt struct {
	err error
}

func (e *errorReaderAt) ReadAt([]byte, int64) (int, error) {
	return 0, e.err
}

func TestLinuxImage(t *testing.T) {
	for _, tt := range []struct {
		li  *LinuxImage
		err error
	}{
		{
			li: &LinuxImage{
				Kernel:  strings.NewReader("foo"),
				Initrd:  strings.NewReader("bar"),
				Cmdline: "foo=bar",
			},
			err: nil,
		},
		{
			li: &LinuxImage{
				Kernel:  strings.NewReader("foo"),
				Initrd:  nil,
				Cmdline: "foo=bar",
			},
			err: nil,
		},
		{
			li: &LinuxImage{
				Kernel:  nil,
				Initrd:  nil,
				Cmdline: "foo=bar",
			},
			err: ErrKernelMissing,
		},
		{
			li: &LinuxImage{
				Kernel:  &errorReaderAt{err: errSkip},
				Initrd:  nil,
				Cmdline: "foo=bar",
			},
			err: errSkip,
		},
		{
			li: &LinuxImage{
				Kernel:  strings.NewReader("foo"),
				Initrd:  &errorReaderAt{err: errSkip},
				Cmdline: "foo=bar",
			},
			err: errSkip,
		},
		{
			li: &LinuxImage{
				Kernel:  strings.NewReader("foo"),
				Initrd:  nil,
				Cmdline: "",
			},
			err: nil,
		},
	} {
		a := cpio.InMemArchive()
		sw := NewSigningWriter(a)
		if err := tt.li.Pack(sw); err != tt.err {
			t.Errorf("Pack(%v) = %v, want %v", tt.li, err, tt.err)
		} else if err == nil {
			li, err := NewLinuxImageFromArchive(a)
			if err != nil {
				t.Errorf("Linux image from %v: %v", a, err)
			}
			if !imageEqual(tt.li, li) {
				t.Errorf("Images are not equal: got %v\nwant %v", li, tt.li)
			}
		}
	}
}
