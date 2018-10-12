// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rsdp allows to determine the ACPI RSDP structure address which could
// be pass to the boot command later on
// It must be executed at the system init as it relies on scanning
// the kernel messages which could be quickly filled up in some cases

//
// Synopsis:
//	rsdp [-d] [-f file]
//
// Description:
//	Look for rsdp value into kernel messages
//
// Example:
//	rsdp 	- Start the script and print rsdp to stdout

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	d        = flag.BoolP("debug", "d", false, "Print debug messages")
	errTimeOut  = fmt.Errorf("Timeout scanning for rsdp")
	cmdUsage = "Usage: rsdp [-f file]"
	file     = flag.StringP("file", "f", "/dev/kmsg", "File to read from")
	debug    = func(string, ...interface{}) {}
)

func usage() {
	log.Fatalf(cmdUsage)
}

func getRSDP(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	channel := make(chan string)

	go func() {
		debug("Read from scanner")
		s := bufio.NewScanner(file)
		for s.Scan() {
			debug("Read '%q'", s.Text())
			channel <- s.Text()
		}
		close(channel)
		if err := s.Err(); err != nil {
			log.Print(err)
		}
	}()

	for {
		select {
		case res := <-channel:
			debug("Read '%q' from chan", res)
			// The Scanner works fine for /dev/kmsg
			// and doesn't seem to know how to deliver
			// EOF on files. It hangs. If we get a ""
			// then quit.
			if res == "" {
				return "", fmt.Errorf("Could not find RSDP")
			}
			if strings.Contains(res, "RSDP") {
				s := strings.Split(res, " ")
				if len(s) < 3 {
					continue
				}
				return s[2], nil
			}

		case <-time.After(1 * time.Second):
			log.Print(errTimeOut)
			return "", errTimeOut
		}
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		usage()
	}
	if *d {
		debug = log.Printf
	}
	rsdp_value, err := getRSDP(*file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" acpi_rsdp=%s \n", rsdp_value)
}
