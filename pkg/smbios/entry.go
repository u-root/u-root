// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
)

func calcChecksum(data []byte, skipIndex int) uint8 {
	var cs uint8
	for i, b := range data {
		if i == skipIndex {
			continue
		}
		cs += b
	}
	return uint8(0x100 - int(cs))
}

// ParseEntry parses SMBIOS 32 or 64-bit entrypoint structure.
func ParseEntry(data []byte) (*Entry32, *Entry64, error) {
	// Entry is either 32 or 64-bit, try them both.
	var e32 Entry32
	if err32 := e32.UnmarshalBinary(data); err32 != nil {
		var e64 Entry64
		if err64 := e64.UnmarshalBinary(data); err64 != nil {
			return nil, nil, fmt.Errorf("%w / %w", err32, err64)
		}
		return nil, &e64, nil
	}
	return &e32, nil, nil
}
