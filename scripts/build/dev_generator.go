// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"os"
)

func init() {
	buildGenerators["dev"] = devGenerator{}
}

type devGenerator struct {
}

// There are a few dev nodes which are required to run an initramfs.
func (g devGenerator) generate(config Config) ([]file, error) {
	// TODO: there are probably some files here we don't actually need
	return []file{
		{"dev/console", nil, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(5, 1)},
		{"dev/null", nil, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(1, 3)},
		{"dev/ttyS0", nil, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(7, 2)},
		{"dev/zero", nil, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(1, 5)},
	}, nil
}
