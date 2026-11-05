// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64 || tinygo

package crabefi

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/boot/kexec"
)

// PrepareLoad is supported on amd64 Linux with the standard Go toolchain.
func PrepareLoad(io.ReaderAt) (uintptr, kexec.Segments, error) {
	return 0, nil, fmt.Errorf("CrabEFI boot requires amd64 Linux and the standard Go toolchain")
}
