// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boottest

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/uio/uio"
)

func mustReadAll(r io.ReaderAt) string {
	if r == nil {
		return ""
	}
	b, err := uio.ReadAll(r)
	if err != nil {
		return fmt.Sprintf("read error: %s", err)
	}
	return string(b)
}

// SameBootImage compares the contents of given boot images, but not the
// underlying URLs.
//
// Works for Linux and Multiboot images.
func SameBootImage(got, want boot.OSImage) error {
	if got.Label() != want.Label() {
		return fmt.Errorf("got image label %s, want %s", got.Label(), want.Label())
	}

	if gotLinux, ok := got.(*boot.LinuxImage); ok {
		wantLinux, ok := want.(*boot.LinuxImage)
		if !ok {
			return fmt.Errorf("got image %s is Linux image, but %s is not", got, want)
		}

		// Same kernel?
		if !uio.ReaderAtEqual(gotLinux.Kernel, wantLinux.Kernel) {
			return fmt.Errorf("got kernel %s, want %s", mustReadAll(gotLinux.Kernel), mustReadAll(wantLinux.Kernel))
		}

		// Same initrd?
		if !uio.ReaderAtEqual(gotLinux.Initrd, wantLinux.Initrd) {
			return fmt.Errorf("got initrd %s, want %s", mustReadAll(gotLinux.Initrd), mustReadAll(wantLinux.Initrd))
		}

		// Same cmdline?
		if gotLinux.Cmdline != wantLinux.Cmdline {
			return fmt.Errorf("got cmdline %s, want %s", gotLinux.Cmdline, wantLinux.Cmdline)
		}
		return nil
	}

	if gotMB, ok := got.(*boot.MultibootImage); ok {
		wantMB, ok := want.(*boot.MultibootImage)
		if !ok {
			return fmt.Errorf("got image %s is Multiboot image, but %s is not", got, want)
		}

		// Same kernel?
		if !uio.ReaderAtEqual(gotMB.Kernel, wantMB.Kernel) {
			return fmt.Errorf("got kernel %s, want %s", mustReadAll(gotMB.Kernel), mustReadAll(wantMB.Kernel))
		}

		// Same cmdline?
		if gotMB.Cmdline != wantMB.Cmdline {
			return fmt.Errorf("got cmdline %s, want %s", gotMB.Cmdline, wantMB.Cmdline)
		}

		if len(gotMB.Modules) != len(wantMB.Modules) {
			return fmt.Errorf("got %d modules, want %d modules", len(gotMB.Modules), len(wantMB.Modules))
		}

		for i := range gotMB.Modules {
			g := gotMB.Modules[i]
			w := wantMB.Modules[i]
			if g.Cmdline != w.Cmdline {
				return fmt.Errorf("module %d got name %s, want %s", i, g.Cmdline, w.Cmdline)
			}
			if !uio.ReaderAtEqual(g.Module, w.Module) {
				return fmt.Errorf("got kernel %s, want %s", mustReadAll(g.Module), mustReadAll(w.Module))
			}
		}
		return nil
	}

	return fmt.Errorf("image not supported")
}
