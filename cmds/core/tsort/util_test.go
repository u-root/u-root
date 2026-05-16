// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func orderInsensitiveDiff(a []string, b []string) string {
	return cmp.Diff(
		a, b, cmpopts.SortSlices(func(x, y string) bool { return x < y }))
}
