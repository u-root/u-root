// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"io"
	"strings"

	"github.com/u-root/u-root/pkg/boot/ibft"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/multiboot"
)

// MultibootImage is a multiboot-formated OSImage, such as ESXi, Xen, Akaros,
// tboot.
type MultibootImage struct {
	Name string

	Kernel   io.ReaderAt
	Cmdline  string
	Modules  []multiboot.Module
	IBFT     *ibft.IBFT
	BootRank int
}

var _ OSImage = &MultibootImage{}

// Label returns either Name or a short description.
func (mi *MultibootImage) Label() string {
	if len(mi.Name) > 0 {
		return mi.Name
	}
	return fmt.Sprintf("Multiboot(kernel=%s cmdline=%s iBFT=%s)", stringer(mi.Kernel), mi.Cmdline, mi.IBFT)
}

// Rank for the boot menu order
func (mi *MultibootImage) Rank() int {
	return mi.BootRank
}

// Edit the kernel command line.
func (mi *MultibootImage) Edit(f func(cmdline string) string) {
	mi.Cmdline = f(mi.Cmdline)
}

// Load implements OSImage.Load.
func (mi *MultibootImage) Load(opts ...LoadOption) error {
	loadOpts := defaultLoadOptions()
	for _, opt := range opts {
		opt(loadOpts)
	}

	entryPoint, segments, err := multiboot.PrepareLoad(loadOpts.verbose, mi.Kernel, mi.Cmdline, mi.Modules, mi.IBFT)
	if err != nil {
		return err
	}
	if !loadOpts.callKexecLoad {
		return nil
	}
	if err := kexec.Load(entryPoint, segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %w", err)
	}
	return nil
}

// String implements fmt.Stringer.
func (mi *MultibootImage) String() string {
	modules := make([]string, len(mi.Modules))
	for i, mod := range mi.Modules {
		modules[i] = mod.Cmdline
	}
	return fmt.Sprintf("MultibootImage(\n  Name: %s\n  Kernel: %s\n  Cmdline: %s\n  iBFT: %s\n  Modules: %s\n)",
		mi.Name, stringer(mi.Kernel), mi.Cmdline, mi.IBFT, strings.Join(modules, ", "))
}
