// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/u-root/u-root/pkg/boot/ibft"
	"github.com/u-root/u-root/pkg/boot/multiboot"
)

// MultibootImage is a multiboot-formated OSImage, such as ESXi, Xen, Akaros,
// tboot.
type MultibootImage struct {
	Name string

	Kernel  io.ReaderAt
	Cmdline string
	Modules []multiboot.Module
	IBFT    *ibft.IBFT
}

// JSONMap is implemented only in order to compare MultibootImages in tests.
//
// It should be json-encodable and decodable.
func (mi *MultibootImage) JSONMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["image_type"] = "multiboot"
	m["name"] = mi.Name
	m["cmdline"] = mi.Cmdline
	if mi.Kernel != nil {
		m["kernel"] = module(mi.Kernel)
	}

	var modules []map[string]interface{}
	for _, mod := range mi.Modules {
		mmod := module(mod.Module)
		mmod["cmdline"] = mod.CmdLine
		mmod["name"] = mod.Name
		modules = append(modules, mmod)
	}
	m["modules"] = modules
	return m
}

func (mi *MultibootImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(mi.JSONMap())
}

var _ OSImage = &MultibootImage{}

// Label returns either Name or a short description.
func (mi *MultibootImage) Label() string {
	if len(mi.Name) > 0 {
		return mi.Name
	}
	return fmt.Sprintf("Multiboot(kernel=%s, cmdline=%s, iBFT=%s)", mi.Kernel, mi.Cmdline, mi.IBFT)
}

// Load implements OSImage.Load.
func (mi *MultibootImage) Load(verbose bool) error {
	return multiboot.Load(verbose, mi.Kernel, mi.Cmdline, mi.Modules, mi.IBFT)
}

// String implements fmt.Stringer.
func (mi *MultibootImage) String() string {
	modules := make([]string, len(mi.Modules))
	for i, mod := range mi.Modules {
		modules[i] = mod.CmdLine
	}
	return fmt.Sprintf("MultibootImage(\n  Name: %s\n  Kernel: %s\n  Cmdline: %s\n  iBFT: %s\n  Modules: %s\n)",
		mi.Name, mi.Kernel, mi.Cmdline, mi.IBFT, strings.Join(modules, ", "))
}
