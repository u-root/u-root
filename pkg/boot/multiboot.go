// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot"
)

// MultibootImage is a multiboot-formated OSImage.
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

// Execute implements OSImage.Execute.
func (mi *MultibootImage) Execute() error {
	m, err := multiboot.New(mi.Path, mi.Cmdline, mi.Modules)
	if err != nil {
		return err
	}
	if err := m.Load(mi.Debug); err != nil {
		return err
	}
	if err := kexec.Load(m.EntryPoint, m.Segments(), 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}
	return kexec.Reboot()
}

// String implements fmt.Stringer.
func (MultibootImage) String() string {
	return fmt.Sprintf("multiboot images unimplemented")
}
