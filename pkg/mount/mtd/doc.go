// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Chips are made by vendors, and an individual vendor is
// defined by a 1 to 8 byte vendor id stored in the chip.
// An instance of a type of chip, with, e.g., a particular size,
// is defined by a 1 to 3 or so byte device id stored in the chip.
// The vendor, device id pair can be used to find both unique
// properties of a chip (e.g. size) and the common properties of a chip
// (e.g. erase value, voltages, erase blocks. etc.)
// It is not at all unusual to know a vendor but not a device id.
// Hence we need to be able to get a vendor given a vendor id,
// and then a chip given a chip id. It's ok, however, to know
// the vendor and fail to know the chip.
//
// Sadly, device ids are not unique; they are reused per vendor.
// And, as mentioned, both vendor and device id are variable length.
// In a not uncommon failure of vision, they started out as 1 byte
// each, grew to 2, then 3 in some cases, 7 in other. Good times.
// Life would be easier if everybody just made these things strings
// in the beginning.
//
// An ID identifies a vendor, but not the same vendor over time.
// Vendor names for a given ID change over time, due to buyouts,
// bankruptcies, and the occasional near depression.
// For example, what was AMD is now Spansion.
// This name changing complicates the picture a bit,
// so we maintain a list of vendor names for a given part, with
// the first name in the list being the current name. This will allow
// us to accomodate scripts that might have the wrong vendor name.
// As time goes by, and bankruptcies accumulate, this first name
// can change.
//
// Hence, it is useful to have 3 bits of knowledge
// o list of vendor names given a vendor id
// o list of chips and their unique properties given a device id
// o list of common properties which can be referenced from a chip
//
// We wish to embed this code in FLASH so if needed we can burn
// a chip from FLASH-embedded u-root.
//
// This code uses strings, not integers,
// since device and vendor IDs are now variable length, depending
// on year of manufacture. Further, it is just nicer to work with
// strings.
//
// In most cases, we will walk these tables once, so we design for
// exhaustive search.  The tables are short and are traversed in
// microseconds, you only do it once, and it's important to keep data
// as compact as possible.
//
// A note on flashing.  Writing is not zero cost: each erase/write
// cycle reduces chip lifetime. Data in the chip need not be erased to
// be written: 0xee can be changed to 0xcc without an erase cycle in
// many parts.  Code can make a guess a guess at an optimal
// erase/write pattern based on the size of the regions to be written,
// the content of regions, and the size of the blocks available.
// Getting this calculation right has proven to be tricky, as it has
// to balance time costs of writing, expected costs of too many erase
// cycles, and several other factors I can not recall just now. Watch
// this space.
//
// TODO: figure out some minimum set of config options for Linux, with
// the proviso that this will be very kernel version dependent.
package mtd
