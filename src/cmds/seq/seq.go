// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	flags struct {
		format     string
		separator  string
		widthEqual bool
	}
	cmd = "seq [-f format] [-w] [-s separator string] <start> [step] <end>"
)

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: %v\n", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.StringVar(&flags.format, "f", "%v", "use printf style floating-point FORMAT")
	flag.StringVar(&flags.separator, "s", "\n", "use STRING to separate numbers")
	flag.BoolVar(&flags.widthEqual, "w", false, "equalize width by padding with leading zeroes")
	flag.Parse()
	flag.Usage = usage
	flags.format = strings.Replace(flags.format, "%", "%0*", 1) // support widthEqual
}

func seq(w io.Writer, args []string) error {
	var (
		stt, end float64
		width    int
		stp      float64 = 1.0
	)
	argv, argc := args, len(args)
	// loading step value
	if argc == 3 {
		_, err := fmt.Sscanf(argv[1], "%v", &stp)
		if err != nil {
			return err
		}
	} else if argc < 2 || argc > 4 {
		return errors.New(fmt.Sprintf("Mismatch n args; got %v, wants 2 > n args > 4", argc))
	}

	if _, err := fmt.Sscanf(argv[0]+" "+argv[argc-1], "%v %v", &stt, &end); err != nil {
		return err
	}

	if flags.widthEqual {
		width = len(fmt.Sprintf(flags.format, 0, end))
	}
	for stt <= end {
		fmt.Fprintf(w, flags.format+flags.separator, width, stt)
		stt += stp
	}

	return nil
}

func main() {
	if err := seq(os.Stdout, flag.Args()); err != nil {
		log.Println(err)
		flag.Usage()
	}
}
