// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"cmp"

	gocmp "github.com/google/go-cmp/cmp"
	gocmpopts "github.com/google/go-cmp/cmp/cmpopts"
)

func orderInsensitiveDiff[T cmp.Ordered](a []T, b []T) string {
	return gocmp.Diff(
		a,
		b,
		gocmpopts.SortSlices(func(x, y T) bool { return x < y }),
		gocmpopts.EquateEmpty(),
	)
}
