// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	flags struct{ universal bool }
	cmd   = "date [-u]"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flags.universal, "u", false, "print or set Coordinated Universal Time (UTC)")
	flag.Usage = usage
	flag.Parse()
}

func date() (string, error) {
	t := time.Now()
	if flags.universal {
		t = t.UTC()
	}
	presentation := t.Format(time.UnixDate)
	return fmt.Sprintf("%v", presentation), nil
}

func main() {
	msg, err := date()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", msg)
}
