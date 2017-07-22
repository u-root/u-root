// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// console implements a basic console. It establishes a pair of files
// to read from, the default being a UART at 0x3f8, but an alternative
// being just stdin and stdout. It will also set up a root file system
// using uroot.Rootfs, although this can be disabled as well.
// Console uses a Go version of fork_pty to start up a shell, default
// /ubin/rush. Console runs until the shell exits and then exits itself.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/pty"
	"github.com/u-root/u-root/uroot"
)

var (
	serial    = flag.String("serial", "0x3f8", "use IO for stdin")
	setupRoot = flag.Bool("setuproot", true, "Set up a root file system")
)

func main() {
	fmt.Printf("console -- starting")
	flag.Parse()

	a := flag.Args()
	if len(a) == 0 {
		a = []string{"/ubin/rush"}
	}

	p, err := pty.New(a[0], a[1:]...)
	if err != nil {
		log.Fatalf("Can't open pty: %v", err)
	}
	// Make a good faith effort to set up root. This being
	// a kind of init program, we do our best and keep going.
	if *setupRoot {
		uroot.Rootfs()
	}

	in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

	if *serial != "" {
		u, err := openUART(*serial)
		if err != nil {
			log.Fatalf("Sorry, can't get a uart: %v", err)
		}
		in, out = u, u
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("Can't start %v: %v", a, err)
		os.Exit(1)
	}
	kid := p.C.Process.Pid

	// You need the \r\n as we are now in raw mode!
	fmt.Printf("Started %d\r\n", kid)

	go io.Copy(out, p.Ptm)

	go func() {
		var data = make([]byte, 1)
		for {
			if _, err := in.Read(data); err != nil {
				fmt.Printf("kid stdin: done\n")
			}
			if data[0] == '\r' {
				if _, err := out.Write(data); err != nil {
					log.Printf("error on echo %v: %v", data, err)
				}
				data[0] = '\n'
			}
			if _, err := p.Ptm.Write(data); err != nil {
				log.Printf("Error writing input to ptm: %v: give up\n", err)
				break
			}
		}
	}()

	if err := p.Wait(); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}
