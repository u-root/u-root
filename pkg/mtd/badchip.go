// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import "fmt"

type badChip string

// ID returns the ChipID.
func (c badChip) ID() ChipID {
	return 0
}

// Name returns the canonical chip name.
func (c badChip) Name() ChipName {
	return ChipName(c)
}

// Synonyms returns all synonyms for a chip.
func (c badChip) Synonyms() []ChipName {
	return []ChipName{c.Name()}
}

// Size returns a ChipSize in bytes.
func (c badChip) Size() ChipSize {
	return 0
}

// String is a stringer for a badChip.
func (c badChip) String() string {
	return fmt.Sprintf("Unknown(%s)", string(c))
}

func newBadChip(c string) ChipInfo {
	return badChip(c)
}
