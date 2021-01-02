// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"bytes"
	"fmt"
	"log"
	"os"
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
	// Dryrun indicates that Load should not Exec
	Dryrun bool
}

var _ = boot.OSImage(&Image{})
var v = func(string, ...interface{}) {}

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
	return fmt.Sprintf("FDT %s from %s, kernel %q, initrd %q", i.Root, i.name, i.Kernel, i.InitRAMFS)
}

// Label returns an Image Label.
func (i *Image) Label() string {
	return i.name
}

// Edit edits the Image cmdline using a func.
func (i *Image) Edit(f func(s string) string) {
	i.Cmdline = f(i.Cmdline)
}

// Load LinuxImage to memory
func loadLinuxImage(i *boot.LinuxImage, verbose bool) error {
	return i.Load(verbose)
}

// provide chance to mock in test
var loadImage = loadLinuxImage
var kexecReboot = kexec.Reboot

// Load loads an image and reboots
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

	if err := loadImage(image, verbose); err != nil {
		return err
	}

	if !i.Dryrun {
		if err := kexecReboot(); err != nil {
			return err
		}
	} else {
		v("Not trying to boot since this is a dry run")
	}

	return nil
}

// LoadBzConfig loads a configuration from a FIT image
// Returns <kernel_name>, <ramdisk_name>, error
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
