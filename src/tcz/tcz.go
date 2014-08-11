// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Wget reads one file from the argument and writes it on the standard output.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

const tcz = "/tinycorelinux.net/5.x/x86/tcz"

// consider making this a goroutine which pushes the string down the channel.
func findloop() (name string, err error) {
     for i := 0; i < 1024; i++ {
     	 l := fmt.Sprintf("/dev/loop%d", i)
	 f, err := os.Open(l)
	 if err != nil {
	    continue
	    }
	    _, err = f.Stat()
	    if err != nil {
	       continue
	       }
	       // if we can get the status, it's in use, too bad.
		// so do this:
//ioctl(3, LOOP_GET_STATUS, {number=0, offset=0, flags=0, name="encstateful", ...}) = 0
//	   and if is 0, you failed; if it fails, you succeed.
	   }
	   return name, nil
}
func linkone(p string, i os.FileInfo, err error) error {
	log.Printf("symtree: p %v\n", p)
	if err != nil {
		return err
	}

	// the tree of symlinks starts at /tcz
	l := filepath.SplitList(p)
	// surely there's a better way.
	n := append([]string{"/"}, l[2:]...)
	to := path.Join(n...)

	log.Printf("symtree: symlink %v to %v\n", p, to)
	return os.Symlink(p, to)
}
func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	cmdName := os.Args[1]
	l := log.New(os.Stdout, "tcz: ", 0)

	if err := os.MkdirAll(tcz, 0600); err != nil {
		l.Fatal(err)
	}

	packagePath := path.Join("/tcz", cmdName)
	if err := os.MkdirAll(packagePath, 0600); err != nil {
		l.Fatal(err)
	}

	// path.Join doesn't quite work here.
	filepath := path.Join(tcz, cmdName+".tcz")
	cmd := "http:/" + filepath

	resp, err := http.Get(cmd)
	if err != nil {
		l.Fatalf("Get of %v failed: %v\n", cmd, err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		l.Fatalf("Not OK! %v\n", resp.Status)
	}

	l.Printf("resp %v err %v\n", resp, err)
	// we've to the whole tcz in resp.Body.
	// First, save it to /tcz/name
	f, err := os.Create(filepath)
	if err != nil {
		l.Fatal("Create of :%v: failed: %v\n", filepath, err)
	} else {
		l.Printf("created %v f %v\n", filepath, f)
	}

	if c, err := io.Copy(f, resp.Body); err != nil {
		l.Fatal(err)
	} else {
		/* OK, these are compressed tars ... */
		l.Printf("c %v err %v\n", c, err)
	}
	/* now mount it. The convention is the mount is in /tcz/packagename */
	if err := syscall.Mount(filepath, packagePath, "squashfs", 0, ""); err != nil {
		l.Fatalf("Mount %s on %s: %v\n", filepath, packagePath, err)
	}
	/* now walk it, and fetch anything we need. */

}
