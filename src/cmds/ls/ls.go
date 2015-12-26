// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Ls reads the directories in the command line and prints out the names.

The options are:
	-l		Long form.
	-r		raw (%v) form
	-R		recurse
*/

package main

import (
	"flag"
	"fmt"
	"github.com/proxypoke/group.go"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
)

var (
	long      = flag.Bool("l", false, "Long form")
	raw       = flag.Bool("r", false, "raw struct")
	recursive = flag.Bool("R", false, "Recurse")
)

func show(fullpath string, fi os.FileInfo) error {
	switch {
	case *raw == true:
		fmt.Printf("%v\n", fi)
	case *long == false:
		fmt.Printf("%v\n", fi.Name())
	// -rw-r--r-- 1 root root 174 Aug 18 17:18 /etc/hosts
	case *long == true:
		usr, err := user.LookupId(fmt.Sprintf("%v", fi.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			return err
		}
		grp, err := group.LookupId(fmt.Sprintf("%v", fi.Sys().(*syscall.Stat_t).Gid))
		if err != nil {
			return err
		}
		fmt.Printf("%v %v %v %v %v %v", fi.Mode(), usr.Username, grp.Name, fi.Size(), fi.ModTime().Format("Jan _2 15:4"), fi.Name())
		if link, err := os.Readlink(fullpath); err == nil {
			fmt.Printf(" -> %v", link)
		}
		fmt.Printf("\n")
	}
	return nil

}

func main() {
	flag.Parse()

	dirs := flag.Args()

	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	for _, v := range dirs {
		if len(dirs) > 1 {
			fmt.Printf("%v:\n", v)
		}
		err := filepath.Walk(v, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%v: %v\n", path, err)
				return err
			}
			if err := show(path, fi); err != nil {
				return err
			}
			if fi.IsDir() && !*recursive && path != v {
				return filepath.SkipDir
			}

			return err
		})
		if err != nil {
			fmt.Printf("%s: %v\n", v, err)
		}
	}
}
