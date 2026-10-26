// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/boot/crabefi"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

// CrabEFIImage boots a fresh UEFI environment using a CrabEFI coreboot payload.
// CrabEFI then discovers EFI boot managers on disks using its own drivers.
type CrabEFIImage struct {
	Name     string
	Payload  io.ReaderAt
	BootRank int
}

var _ OSImage = &CrabEFIImage{}

func (ci *CrabEFIImage) Label() string {
	if ci.Name != "" {
		return ci.Name
	}
	return "UEFI (CrabEFI)"
}

func (ci *CrabEFIImage) String() string {
	return fmt.Sprintf("CrabEFI(payload=%s)", stringer(ci.Payload))
}
func (ci *CrabEFIImage) Rank() int { return ci.BootRank }

// Edit does nothing: this is firmware, not a Linux command line.
func (ci *CrabEFIImage) Edit(func(string) string) {}

// Load implements OSImage.Load. The boot menu still owns cleanup and execution.
func (ci *CrabEFIImage) Load(opts ...LoadOption) error {
	options := defaultLoadOptions()
	for _, opt := range opts {
		opt(options)
	}
	entry, segments, err := crabefi.PrepareLoad(ci.Payload)
	if err != nil {
		return fmt.Errorf("prepare CrabEFI: %w", err)
	}
	options.logger.Printf("CrabEFI entry %#x, segments: %s", entry, segments)
	if !options.callKexecLoad {
		return nil
	}
	return kexec.Load(entry, segments, 0)
}
