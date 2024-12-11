// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/boot/kexec"
)

// KexecLoad loads arm64 Image, with the given ramfs and kernel cmdline.
//
// reservedRanges are additional pieces of physical memory that are not used
// for kexec segment allocation. They are not transmitted to the next kernel to
// be considered reserved.
func KexecLoad(kernel, ramfs *os.File, cmdline string, dtb io.ReaderAt, reservedRanges kexec.Ranges) error {
	img, err := kexecLoadImage(kernel, ramfs, cmdline, dtb, reservedRanges)
	if err != nil {
		return err
	}
	defer img.clean()
	if err = kexec.Load(img.entry, img.segments, 0); err != nil {
		return fmt.Errorf("kexec Load(%v, %v, %d) = %w", img.entry, img.segments, 0, err)
	}
	return nil
}
