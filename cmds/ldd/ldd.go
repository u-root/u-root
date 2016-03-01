// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// ldd prints the full path of dependencies.
// unlike the standard one, you can use it in a script, e.g.
// i=`ldd whatever`
// leaves you with a list of files you can usefully copy.
// You can also feed it a long list of files (/bin/*) and get
// a short list of libraries; further, it will read stdin.
package main

import (
	"debug/elf"
	"errors"
	"fmt"
	"os"
)

var list map[string]bool

func process(file *os.File, name string) error {

	if f, err := elf.NewFile(file); err != nil {
		return err
	} else {
		if fl, err := f.ImportedLibraries(); err != nil {
			return err
		} else {
			if s := f.Section(".interp"); s == nil {
				return errors.New("No interpreter")
			} else {
				if interp, err := s.Data(); err != nil {
					return err
				} else {
					// We could just append the interp but people
					// expect to see that first.
					fl = append([]string{string(interp)}, fl...)
					for _, i := range fl {
						list[i] = true
					}
				}

			}
		}
	}
	return nil
}

func main() {
	list = make(map[string]bool)
	if len(os.Args) < 2 {
		process(os.Stdin, "stdin")
	} else {
		for _, i := range os.Args[1:] {
			if f, err := os.Open(i); err == nil {
				process(f, i)
			} else {
				fmt.Fprintf(os.Stderr, "%v: %v\n", i, err)
			}
		}
	}
	for n := range list {
		fmt.Printf("%v\n", n)
	}
}
