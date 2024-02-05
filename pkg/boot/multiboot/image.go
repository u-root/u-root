// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"fmt"
	"io"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/ibft"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

// Image is a multiboot-formated boot.OSImage, such as ESXi, Xen, Akaros,
// tboot.
type Image struct {
	Name string

	Kernel   io.ReaderAt
	Cmdline  string
	Modules  []Module
	IBFT     *ibft.IBFT
	BootRank int
}

var _ boot.OSImage = &Image{}

// Opt modifies Image.
type Opt func(*Image)

// Append appends to the kernel command-line.
func Append(cmdline string) Opt {
	return func(img *Image) {
		if img.Cmdline == "" {
			img.Cmdline = cmdline
		} else {
			img.Cmdline += " " + cmdline
		}
	}
}

// WithModule appends modules.
func WithModule(m ...Module) Opt {
	return func(img *Image) {
		img.Modules = append(img.Modules, m...)
	}
}

func WithName(name string) Opt {
	return func(img *Image) {
		img.Name = name
	}
}

// NewImage creates a new multiboot image with the given options.
func NewImage(kernel io.ReaderAt, opts ...Opt) *Image {
	img := &Image{
		Kernel: kernel,
	}
	for _, opt := range opts {
		opt(img)
	}
	return img
}

// named is satisifed by *os.File.
type named interface {
	Name() string
}

func stringer(mod interface{}) string {
	if s, ok := mod.(fmt.Stringer); ok {
		return s.String()
	}
	if f, ok := mod.(named); ok {
		return f.Name()
	}
	return fmt.Sprintf("%T", mod)
}

// Label returns either Name or a short description.
func (mi *Image) Label() string {
	if len(mi.Name) > 0 {
		return mi.Name
	}
	return fmt.Sprintf("Multiboot(kernel=%s cmdline=%s iBFT=%s)", stringer(mi.Kernel), mi.Cmdline, mi.IBFT)
}

// Rank for the boot menu order
func (mi *Image) Rank() int {
	return mi.BootRank
}

// Edit the kernel command line.
func (mi *Image) Edit(f func(cmdline string) string) {
	mi.Cmdline = f(mi.Cmdline)
}

// Load implements OSImage.Load.
func (mi *Image) Load(opts ...boot.LoadOption) error {
	loadOpts := boot.DefaultLoadOptions()
	for _, opt := range opts {
		opt(loadOpts)
	}

	entryPoint, segments, err := PrepareLoad(loadOpts.Verbose, mi.Kernel, mi.Cmdline, mi.Modules, mi.IBFT)
	if err != nil {
		return err
	}
	if !loadOpts.CallKexecLoad {
		return nil
	}
	if err := kexec.Load(entryPoint, segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}
	return nil
}

// String implements fmt.Stringer.
func (mi *Image) String() string {
	modules := make([]string, len(mi.Modules))
	for i, mod := range mi.Modules {
		modules[i] = mod.Cmdline
	}
	return fmt.Sprintf("MultibootImage(\n  Name: %s\n  Kernel: %s\n  Cmdline: %s\n  iBFT: %s\n  Modules: %s\n)",
		mi.Name, stringer(mi.Kernel), mi.Cmdline, mi.IBFT, strings.Join(modules, ", "))
}
