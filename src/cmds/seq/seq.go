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
	cmd = "seq [-f format] [-w] [-s separator] [start] [step] <end>"
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
		stt   float64 = 1.0
		stp   float64 = 1.0
		end   float64
		width int
	)

	argv, argc := args, len(args)
	if argc < 1 || argc > 4 {
		return errors.New(fmt.Sprintf("Mismatch n args; got %v, wants 1 >= n args >= 3", argc))
	}

	// loading step value if args is <start> <step> <end>
	if argc == 3 {
		_, err := fmt.Sscanf(argv[1], "%v", &stp)
		if err != nil {
			return err
		}
	}

	if argc >= 2 { // start + end or start + step + end
		if _, err := fmt.Sscanf(argv[0]+" "+argv[argc-1], "%v %v", &stt, &end); err != nil {
			return err
		}
	} else { // only end
		if _, err := fmt.Sscanf(argv[0], "%v", &end); err != nil {
			return err
		}
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
