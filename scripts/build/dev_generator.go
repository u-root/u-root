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
		{"bin", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"dev", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"dev/console", []byte{}, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(5, 1)},
		{"dev/loop-control", []byte{}, 0600 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(10, 237)},
		{"dev/loop0", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 0)},
		{"dev/loop1", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 1)},
		{"dev/loop2", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 2)},
		{"dev/loop3", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 3)},
		{"dev/loop4", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 4)},
		{"dev/loop5", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 5)},
		{"dev/loop6", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 6)},
		{"dev/loop7", []byte{}, 0660 | os.ModeDevice, 0, 0, dev(7, 7)},
		{"dev/null", []byte{}, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(1, 3)},
		{"dev/ttyS0", []byte{}, 0644 | os.ModeDevice | os.ModeCharDevice, 0, 0, dev(7, 2)},
		{"etc", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"etc/localtime", []byte{}, 0644, 0, 0, 0},   // TODO: data
		{"etc/resolv.conf", []byte{}, 0644, 0, 0, 0}, // TODO: data
		{"lib64", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"tcz", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"tmp", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"usr", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
		{"usr/lib", []byte{}, 0755 | os.ModeDir, 0, 0, 0},
	}, nil
}
