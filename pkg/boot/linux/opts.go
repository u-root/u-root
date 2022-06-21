// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

// KexecOptions abstract a collection of options to be passed in KexecLoad.
//
// Arch agnostic. Each arch knows to just look for options they care about.
// Alternatively, we could introduce arch specific options, so irrelevant options
// won't be compiled. But for simplification, have one shared struct to begin
// with, we can split when time comes.
type KexecOptions struct {
	// DTB is used as the device tree blob, if specified.
	DTB string

	// Mmap kernel and initramfs, so virtual pages are directly mapped
	// to page cache. Here it is agnostic to whether original kernel and
	// initramfs file is in tmpfs, or other devices.
	//
	// *) If in tmpfs, file objects are backed by page cache.
	// *) If on disk, kernel cache it during first disk I/O. In case when
	//    we are mmapping a file on disk, pages frames backing the page cache
	//    are only allocated when bytes are accessed, as kernel is lazy, so disk
	//    I/O happens at kexec load time.
	//
	// MmapKernel indicates if mmap kernel kernel into virtual memory.
	MmapKernel bool
	// MmapRamfs indicates if mmap initramfs into virtual memory.
	MmapRamfs bool
}
