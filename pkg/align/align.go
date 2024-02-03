// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package align provides helpers for doing uint alignment.
//
// alignment is done via bit operation at the moment, so alignment
// size need be a power of 2.
package align

import "golang.org/x/exp/constraints"

// Up aligns v up to next multiple of alignSize.
//
// alignSize need be a power of 2.
func Up[T constraints.Unsigned](v T, alignSize T) T {
	mask := alignSize - 1
	return (v + mask) &^ mask
}

// Down aligns v down to a previous multiple of alignSize.
//
// alignSize need be a power of 2.
func Down[T constraints.Unsigned](v T, alignSize T) T {
	return Up(v-(alignSize-1), alignSize)
}

// UpPage aligns v up by system page size.
func UpPage[T constraints.Unsigned](v T) T {
	return Up(v, T(pageSize))
}

// DownPage aligns v down by system page size.
func DownPage[T constraints.Unsigned](v T) T {
	return Down(v, T(pageSize))
}

// IsAligned checks whether v is aligned to alignSize.
//
// alignSize need be a power of 2.
func IsAligned[T constraints.Unsigned](v T, alignSize T) bool {
	return v%alignSize == 0
}
