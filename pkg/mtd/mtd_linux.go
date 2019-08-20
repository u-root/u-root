// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
)

// LinuxChip contains all the information the Linux
// MTD driver provides for a single chip.
// We ignore oobsize and oobavailable.
// They are mentioned as obsolete in library and driver.
type LinuxChip struct {
	// ChipInfo is filled in assuming mtd provides info we can use
	// about a chip that is in our tables. If no such table is found
	// ChipInfo will be a badChipInfo. Bad ChipInfo! Bad! Bad!
	ChipInfo
	// Device name, e.g. mtd0
	Dev string `mtd:"dev"`
	// Chip name as reported by MTD
	// If this is a synonym then the ChipInfo may return
	// something else.
	MTDName string `mtd:"name"`
	// Chip type
	ChipType string `mtd:"type"`
	// Erase size. MTD presents one erase size.
	// It is hard to tell if that is going to work out.
	EraseSize string `mtd:"erasesize"`
	// Part Size as reported by MTD. This should be the same as the
	// ChipInfo, but you never know.
	MTDSize string `mtd:"size"`
	// Write Size, which, yes, can differ
	// from EraseSize.
	WriteSize string `mtd:"writesize"`
	// SubPageSize. Is this the writeable fragment
	// of an erase area?
	SubPageSize string `mtd:"subpagesize"`
	// Regions. Generally Size could equal subpagesize * regions.
	Regions string `mtd:"numeraseregions"`
	// Flags. Which of these we care about remains to be seen.
	Flags string `mtd:"flags"`
}

// DevName is the path name, minus the unit number, of the Linux MTD device.
var DevName = "/sys/devices/virtual/mtd/mtd"

func (l *LinuxChip) String() string {
	return l.MTDName
}

// NewChipInfoFromDev creates a LinuxChip given a file name.
func NewChipInfoFromDev(name string) (ChipInfo, error) {
	var l LinuxChip
	Debug("NewChipInfoFromDev(%s)", name)
	v := reflect.TypeOf(l)
	for ix := 0; ix < v.NumField(); ix++ {
		f := v.Field(ix)
		n := f.Tag.Get("mtd")
		Debug("Field %v, n %v", f, n)
		if n == "" {
			continue
		}
		mf := filepath.Join(name, n)
		s, err := ioutil.ReadFile(mf)
		Debug("Contents of %s is %s", mf, string(s))
		if err != nil {
			return nil, err
		}
		reflect.ValueOf(&l).Elem().Field(ix).SetString(string(s))
	}
	l.ChipInfo = newBadChip(name)
	Debug("Chip is %v", l)
	return &l, nil
}
