// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
)

// FIT is a Flattened Image Tree implementation for OSImage.
type Image struct {
	name       string
	mem        kexec.Memory
	Cmdline    string
	Root       *dt.FDT
	KernelName string
	RootFS     string
	InitRAMFS  string
	entryPoint uintptr
}

var _ = boot.OSImage(&Image{})

// New returns a new image initialized with a file containing an FDT.
func New(n string) (*Image, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fdt, err := dt.ReadFDT(f)
	if err != nil {
		return nil, err
	}
	return &Image{name: n, Root: fdt, mem: kexec.Memory{}}, nil
}

// String is a Stringer for Image.
func (i *Image) String() string {
	return fmt.Sprintf("FDT %s from %s, kernel %s, RootFS %s, InitRamfs %s", i.Root, i.name, i.KernelName, i.RootFS, i.InitRAMFS)
}

// Label returns an Image Label.
func (i *Image) Label() string {
	return i.name
}

// Edit edits the Image cmdline using a func.
func (i *Image) Edit(f func(s string) string) {
	i.Cmdline = f(i.Cmdline)
}

// Load loads an Image.
func (i *Image) Load(verbose bool) error {
	if err := i.mem.ParseMemoryMap(); err != nil {
		return err
	}
	// Find the kernel name and its needed data area.
	if verbose {
		log.Printf("Find kernel %s", i.KernelName)
	}
	kb, err := i.Root.Bytes(i.KernelName, "data")
	if err != nil {
		return err
	}
	i.mem.Segments.Insert(kexec.NewSegment(kb, kexec.Range{Start: 0x100000, Size: uint(len(kb))}))
	// Technically, we can have both an initramfs and a rootfs.
	// The current question: how do we provided the rootfs from a file
	// to the kernel? We do not know.
	if len(i.RootFS) > 0 {
		return fmt.Errorf("No way to provide a rootfs yet")
	}

	if len(i.InitRAMFS) > 0 {
		if verbose {
			log.Printf("Find initramfs %s", i.InitRAMFS)
		}
		ib, err := i.Root.Bytes(i.InitRAMFS, "data")
		if err != nil {
			return err
		}
		i.mem.Segments.Insert(kexec.NewSegment(kb, kexec.Range{Start: 32 * 0x100000, Size: uint(len(ib))}))
	}
	i.entryPoint = 0x100000

	if verbose {
		log.Printf("segments cmdline %v %q", i.mem.Segments, i.Cmdline)
	}

	if err := kexec.Load(i.entryPoint, i.mem.Segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}

	return nil
}
