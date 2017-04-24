// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
	"fmt"
)

func init() {
	archiveGenerators["list"] = listGenerator{}
}

type listGenerator struct {
}

// Rather than creating an archive, this simply outputs a list for debugging.
func (g listGenerator) generate(config Config, files []file) error {
	totalSize := 0
	for _, f := range files {
		if f.rdev == 0 {
			fmt.Printf("%v\t%d\t%q\n", f.mode, len(f.data), f.path)
			totalSize += len(f.data)
		} else {
			fmt.Printf("%v\t%d, %d\t%q\n", f.mode, major(f.rdev), minor(f.rdev), f.path)
		}
	}
	fmt.Println("Number of files:", len(files))
	fmt.Printf("Total size: %.1f MiB (%d bytes)\n", float64(totalSize)/1024/1024, totalSize)
	return nil
}

func (g listGenerator) run(config Config) error {
	return errors.New("not implemented")
}
