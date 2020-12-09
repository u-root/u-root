// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"bytes"
	"fmt"
	"os"
	"log"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
)

// Image is a Flattened Image Tree implementation for OSImage.
type Image struct {
	name string
	// Cmdline is the command line for the new kernel.
	Cmdline string
	// Root is the FDT.
	Root *dt.FDT
	// Kernel is the name of the kernel node.
	Kernel string
	// InitRAMFS is the name of the initramfs node.
	InitRAMFS string
	// RootFS is a root file system image file.
	// It is unclear how to use this in kexec.
	RootFS string
	// Dryrun indicates that Load should not Exec
	Dryrun bool
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
	return &Image{name: n, Root: fdt}, nil
}

// String is a Stringer for Image.
func (i *Image) String() string {
	return fmt.Sprintf("FDT %s from %s, kernel %q, initrd %q, RootFS %q", i.Root, i.name, i.Kernel, i.InitRAMFS, i.RootFS)
}

// Label returns an Image Label.
func (i *Image) Label() string {
	return i.name
}

// Edit edits the Image cmdline using a func.
func (i *Image) Edit(f func(s string) string) {
	i.Cmdline = f(i.Cmdline)
}

var v = func(string, ...interface{}) {}

// Load loads an Image.
func (i *Image) Load(verbose bool) error {
	if verbose {
		v = log.Printf
	}

	w := i.Root.Root().Walk("images").Walk(i.Kernel)
	b, err := w.Property("data").AsBytes()
	if err != nil {
		return err
	}

	image := &boot.LinuxImage{
		Kernel:  bytes.NewReader(b),
		Cmdline: i.Cmdline,
	}

	if len(i.InitRAMFS) != 0 {
		w := i.Root.Root().Walk("images").Walk(i.InitRAMFS)
		b, err := w.Property("data").AsBytes()
		if err != nil {
			return err
		}
		image.Initrd = bytes.NewReader(b)
	}
	// Technically, we can have both an initramfs and a rootfs.
	// The current question: how do we provided the rootfs from a file
	// to the kernel? We do not know.
	if len(i.RootFS) > 0 {
		return fmt.Errorf("No way to provide a rootfs yet")
	}
	if err := image.Load(verbose); err != nil {
		return err
	}
	if !i.Dryrun {
		if err := kexec.Reboot(); err != nil {
			return err
		}
	} else {
		v("Not trying to boot since this is a dry run")
	}

	return nil
}

// Load configuration of FIT Image.
func (i *Image) LoadBzConfig(verbose bool) (string, string, error) {
	configs := i.Root.Root().Walk("configurations")
	dc, err := configs.Property("default").AsString()
	if err != nil {
		return "", "", err
	}
	n := strings.Split(dc, "@")
	if len(n) != 2 {
		return "", "", fmt.Errorf("Invalid default configuration naming: %v", dc)
	}
	bzn := n[0] + "_bz@" + n[1]
	bzconfig := i.Root.Root().Walk("configurations").Walk(bzn)

	kn, err := bzconfig.Property("kernel").AsString()
	if err != nil {
		return "", "", err
	}
	rn, err := bzconfig.Property("ramdisk").AsString()
	if err != nil {
		return "", "", err
	}

	return kn, rn, nil
}
