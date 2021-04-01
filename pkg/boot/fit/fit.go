// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot"
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
	// ConfigOverride is the optional FIT config to use instead of default
	ConfigOverride string
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

// ParseConfig reads r for a FIT image and returns a OSImage for each
// configuration parsed.
func ParseConfig(r io.ReadSeeker, verbose bool) ([]Image, error) {
	fdt, err := dt.ReadFDT(r)
	if err != nil {
		return nil, err
	}

	var images []Image
	configs := fdt.Root().Walk("configurations")
	cn, _ := configs.ListChildNodes()

	for _, n := range cn {
		var i = Image{name: n, Root: fdt, ConfigOverride: n}

		kn, in, err := i.LoadConfig(false, verbose)

		if err == nil {
			i.Kernel, i.InitRAMFS = kn, in
			images = append(images, i)
		}
		if err != nil && verbose {
			log.Printf("Failed to load config for config name %s: %v", n, err)
		}
	}

	return images, nil
}

// String is a Stringer for Image.
func (i *Image) String() string {
	return fmt.Sprintf("FDT %s, kernel %q, initrd %q", i.name, i.Kernel, i.InitRAMFS)
}

// Label returns an Image Label.
func (i *Image) Label() string {
	if i.Kernel == "" {
		return i.name
	}
	return i.name + " (kernel:" + i.Kernel + ")"
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

	return nil
}

// GetConfigName finds the name of the default configuration or returns the
// override config if available
func (i *Image) GetConfigName(verbose bool) (string, error) {
	if len(i.ConfigOverride) != 0 {
		return i.ConfigOverride, nil
	}

	configs := i.Root.Root().Walk("configurations")
	dc, err := configs.Property("default").AsString()
	if err != nil {
		return "", err
	}

	return dc, nil
}

// LoadConfig loads a configuration from a FIT image
// Returns <kernel_name>, <ramdisk_name>, error
func (i *Image) LoadConfig(skipinitfs bool, verbose bool) (string, string, error) {
	tc, err := i.GetConfigName(verbose)
	if err != nil {
		return "", "", err
	}

	configs := i.Root.Root().Walk("configurations")
	config := configs.Walk(tc)
	_, err = config.AsString()

	if err != nil {
		if verbose {
			cs, _ := configs.ListChildNodes()
			log.Printf("Config options: %#v", cs)
		}
		return "", "", err
	}

	var kn, rn string

	kn, err = config.Property("kernel").AsString()
	if err != nil {
		return "", "", err
	}

	if !skipinitfs {
		rn, err = config.Property("ramdisk").AsString()
		if err != nil {
			return "", "", err
		}
	}

	return kn, rn, nil
}
