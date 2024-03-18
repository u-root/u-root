// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A slice of bytes of the literal plain text of pci.ids has been found
// to produce the smallest binary compared to a native Go map, marshalled
// JSON, and go-bindata (gzip'ed bytes). Further runtime of parsing the plain
// text pci.ids is lower than all options compared. The pciids in this package
// is stripped of all comments, empty lines, sub-devices, and classes to save
// on binary size.

package pci

import (
	_ "embed"
)

type idMap map[uint16]Vendor

var (
	//go:embed pci.ids
	pciids []byte
	ids    idMap
)

// newIDs contains the plain text contents of pci.ids. It returns
// a map to be used as lookup from hex ID to human readable label.
// We do not admit of the possibility of error, any failure
// should be caught by the test. We might just want to just always
// create ids since the most common use of pci will be with names,
// not numbers.
func newIDs() idMap {
	if ids != nil {
		return ids
	}

	ids = parse(pciids)
	return ids
}
