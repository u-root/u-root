// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"bytes"
	"fmt"
	"io"
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
	// SkipInitRAMFS skips the search for an ramdisk entry in the config
	SkipInitRAMFS bool
	// BootRank ranks the priority of the images in boot menu
	BootRank int
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

// ParseConfig reads r for a FIT image and returns a OSImage for each
// configuration parsed.
func ParseConfig(r io.ReadSeeker) ([]Image, error) {
	fdt, err := dt.ReadFDT(r)
	if err != nil {
		return nil, err
	}

	var images []Image
	configs := fdt.Root().Walk("configurations")
	cn, _ := configs.ListChildNodes()

	for _, n := range cn {
		var i = Image{name: n, Root: fdt, ConfigOverride: n}

		kn, in, err := i.LoadConfig()

		if err == nil {
			i.Kernel, i.InitRAMFS = kn, in
			images = append(images, i)
		}
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("failed to find valid usable config")
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
	return fmt.Sprintf("%s (kernel: %s)", i.name, i.Kernel)
}

// Rank returns an Image Rank.
func (i *Image) Rank() int {
	return i.BootRank
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
func (i *Image) GetConfigName() (string, error) {
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
func (i *Image) LoadConfig() (string, string, error) {
	tc, err := i.GetConfigName()
	if err != nil {
		return "", "", err
	}

	configs := i.Root.Root().Walk("configurations")
	config := configs.Walk(tc)
	_, err = config.AsString()

	if err != nil {
		return "", "", err
	}

	var kn, rn string

	kn, err = config.Property("kernel").AsString()
	if err != nil {
		return "", "", err
	}

	// Allow missing initram nodes
	rn, _ = config.Property("ramdisk").AsString()

	return kn, rn, nil
}
