// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && !tamago

// Package align provides helpers for doing uint alignment.
package align

import "os"

var pageSize = uint(os.Getpagesize())
