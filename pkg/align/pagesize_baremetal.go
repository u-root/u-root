// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tamago || tinygo

// Package align provides helpers for doing uint alignment.
package align

// While there is no page size in bare metal, but it is probably sensible
// to align it to a common cache line size boundary.
var pageSize = uint(64)
