// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hdparm performs control operations on disks.
//
// hdparm is modeled after linux command of the same name,
// using a command [switches] verb model, but many of the verbs
// are also switches which can conflict.
//
// Upstream hdparm also decided a conf file in /etc/ would be
// a good idea; we have no plans to support that.
//
// We also have no plans to support ata12. It's 2019.
//
// Unlike the standard hdparm command, we do not allow empty
// arguments for commands requiring a password.
//
// Synopsis:
//
//	hdparm [--i] [--security-unlock[=password]] [--user-master|--timeout] [device ...]
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/u-root/u-root/pkg/mount/scuzz"
)

// Because all the verbs and options are switches, hence global,
// we can have all the functions take one Disk as a parameter and
// return a string and error. This simplifies other aspects of this program.
type op func(scuzz.Disk) (string, error)

var (
	verbose         = flag.Bool("v", false, "verbose log")
	debug           = func(string, ...interface{}) {}
	unlock          = flag.String("security-unlock", "", "Unlock the drive with a password")
	identify        = flag.Bool("i", false, "Get drive identifying information")
	admin           = flag.Bool("user-master", false, "Unlock admin (true) or user (false)")
	timeoutDuration = flag.String("timeout", "15s", "Timeout for operations expressed as a Go duration (e.g. 15s)")
	verbs           = []string{"security-unlock", "i"}
)

// The hdparm switches can conflict. This function returns nil if there is no conflict, and a (hopefully)
// helpful error message otherwise. As a side effect it assigns verb.
// TODO: use the visitor pattern to make this better. The next flag that's added will include that fix.
func checkVerbs() (op, error) {
	var v []string
	var verb op

	if len(*unlock) > 0 {
		verb = unlockop
		v = append(v, "security-unlock")
	}
	if *identify {
		verb = identifyop
		v = append(v, "i")
	}

	if len(v) > 1 {
		return nil, fmt.Errorf("%v verbs were invoked and only one is allowed", v)
	}
	if len(v) < 1 {
		return nil, fmt.Errorf("no verbs were invoked and one of %v is required", verbs)
	}
	return verb, nil
}

func unlockop(d scuzz.Disk) (string, error) {
	return "", d.Unlock(*unlock, *admin)
}

func identifyop(d scuzz.Disk) (string, error) {
	i, err := d.Identify()
	if err != nil {
		return "", err
	}
	return i.String(), nil
}

func main() {
	flag.Parse()

	verb, err := checkVerbs()
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		debug = log.Printf
	}
	scuzz.Debug = debug

	timeout, err := time.ParseDuration(*timeoutDuration)
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range flag.Args() {
		d, err := scuzz.NewSGDisk(n, scuzz.WithTimeout(timeout))
		if err != nil {
			log.Printf("%v: %v", n, err)
			continue
		}
		s, err := verb(d)
		if err != nil {
			log.Printf("%v: %v", n, err.Error())
		}
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
}
