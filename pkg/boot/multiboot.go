// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/multiboot"
)

// MultibootImage is a multiboot-formated OSImage, such as ESXi, Xen, Akaros,
// tboot.
type MultibootImage struct {
	Path    string
	Cmdline string
	Modules []string
}

var _ OSImage = &MultibootImage{}

// Load implements OSImage.Load.
func (mi *MultibootImage) Load(verbose bool) error {
	return multiboot.Load(verbose, mi.Path, mi.Cmdline, mi.Modules, nil)
}

// String implements fmt.Stringer.
func (mi *MultibootImage) String() string {
	return fmt.Sprintf("MultibootImage(\n  KernelPath: %s\n  Cmdline: %s\n  Modules: %s\n)",
		mi.Path, mi.Cmdline, strings.Join(mi.Modules, ", "))
}
