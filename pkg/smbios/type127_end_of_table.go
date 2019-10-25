// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
)

// EndOfTable is Defined in DSP0134 7.46.
type EndOfTable struct {
	Table
}

// NewEndOfTable parses a generic Table into EndOfTable.
func NewEndOfTable(t *Table) (*EndOfTable, error) {
	if t.Type != TableTypeEndOfTable {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	return &EndOfTable{Table: *t}, nil
}

func (eot *EndOfTable) String() string {
	return eot.Header.String()
}
