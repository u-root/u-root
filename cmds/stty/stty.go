// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// stty is an stty command in Go.
// It follows many of the conventions of standard stty.
// However, it can produce JSON output, for later use, and can
// read that JSON later to configure it.
//
// stty has always had an odd set of flags. -flag means turn flag off;
// flag means turn flag on. Except for those flags which make an argument;
// in that case they look like flag <arg>
// To make the flag package continue to work, we've changed the - to a ~.
//
// Programmatically, the options are set with a []string, not lots of magic numbers that
// are not portable across kernels.
//
// The default action is to print in the model of standard stty, which is all most
// people ever do anyway.

// The command works like this:
// stty [verb] [options]
// Verbs are:
// dump -- dump the json of the struct to stdout
// load -- read a json file from stdin and use it to set
// raw -- convenience command to set raw
// cooked -- convenience command to set cooked
// In common stty usage, options may be specified without a verb.
//
// any other verb, with a ~ or without, is taken to mean standard stty args, e.g.
// stty ~echo
// turns off echo. Flags with arguments work too:
// stty intr 1
// sets the interrupt character to ^A.
//
// The JSON encoding lets you do things like this:
// stty dump | sed whatever > file
// stty load file
// Further, one can easily push and pop state in by storing the current
// state in a file in JSON, making changes, and restoring it later. This has
// always been inconvenient in standard stty.
//
// While GNU stty can do some of this, its way of doing it is harder to read and not
// as portable, since the format they use is not self-describing:
// stty -g
// 4500:5:bf:8a3b:3:1c:7f:15:4:0:1:0:11:13:1a:0:12:f:17:16:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0
//
// We always do our operations on fd 0, as that is standard, and we always do an initial
// gtty to ensure we have access to fd 0.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	t, err := gtty(0)

	if err != nil {
		log.Fatalf("gtty: %v", err)
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "pretty")
	}

	switch os.Args[1] {
	case "pretty":
		pretty(os.Stdout, t)
	case "dump":
		b, err := json.MarshalIndent(t, "", "\t")

		if err != nil {
			log.Fatalf("json marshal: %v", err)
		}
		fmt.Printf("%s\n", b)
	case "load":
		if len(os.Args) != 3 {
			log.Fatalf("arg count")
		}
		b, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			log.Fatalf("stty load: %v", err)
		}
		if err := json.Unmarshal(b, t); err != nil {
			log.Fatalf("stty load: %v", err)
		}
		if t, err = stty(0, t); err != nil {
			log.Fatalf("stty: %v", err)
		}
		pretty(os.Stdout, t)
	case "raw":
		if _, err := setRaw(0); err != nil {
			log.Fatalf("raw: %v", err)
		}
	default:
		if err := setOpts(t, os.Args[1:]); err != nil {
			log.Fatalf("setting opts: %v", err)
		}
		t, err = stty(0, t)
		if err != nil {
			log.Fatalf("stty: %v", err)
		}
		pretty(os.Stdout, t)
	}
}
