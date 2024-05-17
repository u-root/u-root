// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netstat"
)

var (
	interfacesFlag = flag.BoolP("interfaces", "i", false, "display interface table")
	ifFlag         = flag.StringP("interface", "I", "", "Display interface table for interface <if>")

	continFlag    = flag.BoolP("continuous", "c", false, "continuous listing")
)

func evalFlags() error {
	flag.Parse()

	if *interfacesFlag {
		if err := netstat.PrintInterfaceTable(*ifFlag, *continFlag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *ifFlag != "" {
		if err := netstat.PrintInterfaceTable(*ifFlag, *continFlag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	return nil
}
func main() {
	if err := evalFlags(); err != nil {
		log.Fatal(err)
	}
}
