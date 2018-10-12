// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rsdp allows to determine the ACPI RSDP structure address which could
// be passed to the boot command later on
// It must be executed at the system init as it relies on scanning
// the kernel messages which could be quickly filled up in some cases
//
// Synopsis:
//	rsdp [-f file]
//
// Description:
//	Look for rsdp value in a file, default /dev/kmsg
//
// Example:
//	rsdp
//	rsdp -f /path/to/file
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"golang.org/x/sys/unix"
)

var (
	cmdUsage = "Usage: rsdp [-f file]"
	file     = flag.StringP("file", "f", "/dev/kmsg", "File to read from")
)

func usage() {
	log.Fatalf(cmdUsage)
}

func getRSDP(path string) (string, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_NONBLOCK, 0)
	if err != nil {
		log.Fatal(err)
	}
	file := os.NewFile(uintptr(fd), "kernel messages")
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}
		res := s.Text()
		if strings.Contains(res, "RSDP") {
			rv := strings.Split(res, " ")
			if len(res) < 3 {
				continue
			}
			return rv[2], nil
		}
	}
	return "", fmt.Errorf("Could not find RSDP")
}

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		usage()
	}
	rsdp_value, err := getRSDP(*file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" acpi_rsdp=%s \n", rsdp_value)
}
