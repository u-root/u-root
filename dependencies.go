// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tools

package main

// List vmtest commands that need to be in go.mod & go.sum to be buildable as
// dependencies. This way, they aren't eliminated by `go mod tidy`.
//
// But obviously aren't actually importable, since they are main packages.
import (
	_ "github.com/hugelgupf/vmtest/vminit/gouinit"
	_ "github.com/hugelgupf/vmtest/vminit/shelluinit"
	_ "github.com/hugelgupf/vmtest/vminit/shutdownafter"
	_ "github.com/hugelgupf/vmtest/vminit/vmmount"
)
