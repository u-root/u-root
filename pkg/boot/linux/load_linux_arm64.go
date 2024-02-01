// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/boot/kexec"
)

// KexecLoad loads arm64 Image, with the given ramfs and kernel cmdline.
func KexecLoad(kernel, ramfs *os.File, cmdline string, opts KexecOptions) error {
	img, err := kexecLoadImage(kernel, ramfs, cmdline, opts)
	if err != nil {
		return err
	}
	defer img.clean()
	if err = kexec.Load(img.entry, img.segments, 0); err != nil {
		return fmt.Errorf("kexec Load(%v, %v, %d) = %v", img.entry, img.segments, 0, err)
	}
	return nil
}
