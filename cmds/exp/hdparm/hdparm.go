// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hdparm perform control operations on disks
// modeled after linux command of the same name.
// The hdparm command is a mess. It is modeled on the
// command [switches] verb model but many of the verbs
// are switches which can conflict. To add to the
// fun they decided a conf file in /etc/ would be
// a good idea; we have no plans to support that.
//
// We also have no plans to support ata12. It's 2019.
//
// Synopsis:
//     hdparm [not all the switches hdparm supports] [device ...]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/scuzz"
)

// Because all the verbs and options are switches, hence global,
// we can have all the functions take one Disk as a parameter and
// return an error. This simplifies other aspects of this program.
type op = func(scuzz.Disk) error

var (
	verbose = flag.Bool("v", false, "verbose log")
	debug   = func(string, ...interface{}) {}
	unlock  = flag.String("security-unlock", "", "Unlock the drive")
	master  = flag.Bool("user-master", false, "Unlock master (true) or user (false)")
	timeout = flag.Uint("timeout", 15000, "Timeout for operations")
	verbs   = map[string]op{
		"--security-unlock": unlockop,
	}
	verb op
)

// The hdparm switches can conflict. This function returns nil if there is no conflict, and a (hopefully)
// helpful error message otherwise. As a side effect it assigns verb.
func checkConflict() error {
	var v []string

	for _, a := range os.Args {
		if f, ok := verbs[a]; ok {
			v = append(v, a)
			verb = f
		}
	}

	if len(v) > 1 {
		return fmt.Errorf("%v verbs were invoked and only one is allowed", v)
	}
	if len(v) < 1 {
		return fmt.Errorf("no verbs were invoked and one of %v is required", verbs)
	}
	return nil
}

func unlockop(d scuzz.Disk) error {
	return scuzz.Unlock(d, *unlock, *timeout, *master)
}

func main() {
	flag.Parse()

	if err := checkConflict(); err != nil {
		log.Fatal(err)
	}

	if *verbose {
		debug = log.Printf
	}
	scuzz.Debug = debug

	for _, n := range flag.Args() {
		d, err := scuzz.NewSGDisk(n)
		if err != nil {
			log.Printf("%v: %v", n, err)
		}
		if err := verb(d); err != nil {
			log.Printf("%v: %v", n, err)
		}
	}

}
