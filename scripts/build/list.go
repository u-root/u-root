// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
	"fmt"
	"io/ioutil"
)

func init() {
	archivers["list"] = listArchiver{}
}

type listArchiver struct {
}

// Rather than creating an archive, this simply outputs a list for debugging.
func (a listArchiver) generate(config Config, files []file) error {
	totalSize := 0
	// TODO: use "text/tabwriter" for nicer alignment
	for _, f := range files {
		if f.data != nil {
			// TODO: Can we get the length of a Reader without reading its entirety into memory?
			data, err := ioutil.ReadAll(f.data)
			if err != nil {
				return err
			}
			fmt.Printf("%v\t%d\t%q\n", f.mode, len(data), f.path)
			totalSize += len(data)
		} else if f.rdev != 0 {
			fmt.Printf("%v\t%d, %d\t%q\n", f.mode, major(f.rdev), minor(f.rdev), f.path)
		} else {
			fmt.Printf("%v\t\t%q\n", f.mode, f.path)
		}
	}
	fmt.Println("Number of files:", len(files))
	fmt.Printf("Total size: %.1f MiB (%d bytes)\n", float64(totalSize)/1024/1024, totalSize)
	return nil
}

func (a listArchiver) run(config Config) error {
	return errors.New("not implemented")
}
