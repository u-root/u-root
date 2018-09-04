// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"fmt"
)

// Region contains the start and end of a region in flash. This can be a BIOS, ME, PDR or GBE region.
// This value seems to index blocks of block size 0x1000
// TODO: figure out of block sizes are read from some location on flash or fixed.
// Right now we assume they're fixed

const (
	// RegionBlockSize assumes the region struct values correspond to blocks of 0x1000 in size
	RegionBlockSize = 0x1000
)

// Region holds the base and limit of every type of region. Each region such as the bios region
// should point back to it.
type Region struct {
	Base  uint16 // Index of first 4k block
	Limit uint16 // Index of last block
}

// Valid checks to see if a region is valid
func (r *Region) Valid() bool {
	return r.Limit > 0 && r.Limit >= r.Base
}

func (r *Region) String() string {
	return fmt.Sprintf("[%#x, %#x)", r.Base, r.Limit)
}

// BaseOffset calculates the offset into the flash image where the Region begins
func (r *Region) BaseOffset() uint32 {
	return uint32(r.Base) * RegionBlockSize
}

// EndOffset calculates the offset into the flash image where the Region ends
func (r *Region) EndOffset() uint32 {
	return (uint32(r.Limit) + 1) * RegionBlockSize
}
