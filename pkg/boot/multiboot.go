// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/ibft"
	"github.com/u-root/u-root/pkg/boot/multiboot"
)

// MultibootImage is a multiboot-formated OSImage, such as ESXi, Xen, Akaros,
// tboot.
type MultibootImage struct {
	Path    string
	Cmdline string
	Modules []string
	IBFT    *ibft.IBFT
}

var _ OSImage = &MultibootImage{}

// Load implements OSImage.Load.
func (mi *MultibootImage) Load(verbose bool) error {
	mbkernel, err := os.Open(mi.Path)
	if err != nil {
		return err
	}
	defer mbkernel.Close()
	modules := make([]multiboot.Module, len(mi.Modules))
	for i, cmd := range mi.Modules {
		modules[i].CmdLine = cmd
		name := strings.Fields(cmd)[0]
		f, err := os.Open(name)
		if err != nil {
			return fmt.Errorf("error opening module %v: %v", name, err)
		}
		defer f.Close()
		modules[i].Module = f
	}

	return multiboot.Load(verbose, mbkernel, mi.Cmdline, modules, mi.IBFT)
}

// String implements fmt.Stringer.
func (mi *MultibootImage) String() string {
	return fmt.Sprintf("MultibootImage(\n  KernelPath: %s\n  Cmdline: %s\n  Modules: %s\n)",
		mi.Path, mi.Cmdline, strings.Join(mi.Modules, ", "))
}
