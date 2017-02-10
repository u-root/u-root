/*
 * ifdtool - dump Intel Firmware Descriptor information
 *
 * Copyright (C) 2011 The ChromiumOS Authors.  All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */
package main

import (
	"fmt"
	"io"
)

const (
	v1 = 1
	v2
	linelen = 80
	magic   = 0x0ff0a55a
)

const (
	f20MHZ = 0
	f33MHZ
	f48MHZ
	f50_30MHZ = 4
	f17MHZ    = 6
)

const (
	s512KB = 0
	s1MB
	s2MB
	s4MB
	s8MB
	s16M
	s32MB
	s64MB
	sUNUSED = 0xf
)

type fdbar struct {
	Map0  uint32
	Map1  uint32
	Map2  uint32
	_     [0xefc - 0x20]uint8
	UMap1 uint32
}

func (f *fdbar) String() string {
	return fmt.Sprintf("Map[%#x, %#x, %#x] Umap %#x", f.Map0, f.Map1, f.Map2, f.UMap1)
}

/*
 * WR / RD bits start at different locations within the flmstr regs, but
 * otherwise have identical meaning.
 */
const (
	flmstr_wr_v1  = 24
	flmastr_wr_v2 = 20
	flmastr_rd_v1 = 16
	flmast_rd_vw  = 8
)

// regions
type oldregions [5]uint32

type regions [9]uint32

type fcba struct {
	flcomp uint32
	flill  uint32
	flpb   uint32
}

type fpsba [17]uint32

// master
type fmba [5]uint32

// processor strap
type fmsba [8]uint32

// ME VSCC
type vscc struct {
	jid  uint32
	vscc uint32
}

// Actual number of entries specified in vtl
type vtba [8]vscc

func (r *region) String() string {
	return fmt.Sprintf("%#08x %#08x %#08x",
		r.Base, r.Limit, r.Size)
}

func (i *image) String() string {
	return fmt.Sprintf("%s", i.fdbar.String())

}

type region struct {
	Base  int32
	Limit int32
	Size  int32
}

type regionName struct {
	pretty string
	terse  string
}

type image struct {
	fdbar
}
type chip struct {
	io.Reader
	Data image
}
