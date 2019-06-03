// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"log"
	"strings"

	"github.com/u-root/u-root/pkg/multiboot"
)

// MultibootImage is a multiboot-formated OSImage, such as ESXi, Xen, Akaros,
// tboot.
type MultibootImage struct {
	Debug   bool
	Path    string
	Cmdline string
	Modules []string
}

var _ OSImage = &MultibootImage{}

// ExecutionInfo implements OSImage.ExecutionInfo.
func (MultibootImage) ExecutionInfo(log *log.Logger) {
	log.Printf("Multiboot images are unsupported")
}

// Load implements OSImage.Load.
func (mi *MultibootImage) Load() error {
	return multiboot.Load(mi.Debug, mi.Path, mi.Cmdline, mi.Modules)
}

// String implements fmt.Stringer.
func (mi *MultibootImage) String() string {
	return fmt.Sprintf("MultibootImage(\n  KernelPath: %s\n  Cmdline: %s\n  Modules: %s\n)",
		mi.Path, mi.Cmdline, strings.Join(mi.Modules, ", "))
}
