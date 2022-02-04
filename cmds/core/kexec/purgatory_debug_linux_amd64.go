// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/* C to go purgatory */

// TODO(10000TB): Add purgatory impl in go.

// Purgatory is a compiled purgatory buffer ("bin-to-hex"-ed) that can be loaded as trampoline.
//
// Compiled from kexec-tools, taken from generated kexec/purgatory.c, which is made by
// https://github.com/horms/kexec-tools/blob/main/purgatory/Makefile#L10
// kexec-tools
// v2.0.23
// HEAD: 91ff1e713733c0f9e6503a29d3f468ac9cc8f97f
//
// Easily questionable purgatory "impl" (for now :) -- To be replaced by go impl.
//
// This is what kexec userspace would pass to elf_rel_build_load during trampline load step:
//
//  elf_rel_build_load(info, &info->rhdr, purgatory, purgatory_size,
//      0x3000, 0x7fffffff, -1, 0);
//
//  OR
//
//  elf_rel_build_load(info, &info->rhdr, purgatory, purgatory_size,
//      0x3000, 640*1024, -1, 0);
//
// This was produced by ld -N -e entry64 -Ttext=0x3000 -o x purgatory/purgatory.ro

/* command line parameter may be appended by purgatory */
const PURGATORY_CMDLINE_SIZE = 64
