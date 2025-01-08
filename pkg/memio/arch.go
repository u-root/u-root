// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && linux && (amd64 || 386)

package memio

// functions exported by the architecture-specific code in ioprt_linux_amd64.s
// and ioprt_linux_386.s
func archInl(uint16) uint32
func archInw(uint16) uint16
func archInb(uint16) uint8

func archOutl(uint16, uint32)
func archOutw(uint16, uint16)
func archOutb(uint16, uint8)
