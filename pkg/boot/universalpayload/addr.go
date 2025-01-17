// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (amd64 || arm64) && !tinygo

package universalpayload

func addrOfStart() uintptr
func addrOfStackTop() uintptr
func addrOfHobAddr() uintptr
