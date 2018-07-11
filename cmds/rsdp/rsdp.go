// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rsdp allows to determine the ACPI RSDP structure address which could
// be pass to the boot command later on
// It must be executed at the system init as it relies on scanning
// the kernel messages which could be quickly filled up in some cases

//
// Synopsis:
//	rsdp
//
// Description:
//	Look for rsdp value into kernel messages
//
// Example:
//	rsdp 	- Start the script and save the founded value into /tmp/rsdp

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func getRSDP(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var returnValue string
	channel := make(chan string)

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			mystring := scanner.Text()
			channel <- mystring
		}
		close(channel)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	var dataRead int
	var exit int
	dataRead = 1
	exit = 0
	for dataRead == 1 && exit == 0 {
		select {
		case res := <-channel:
			if strings.Contains(res, "RSDP") {
				returnValue = strings.Split(res, " ")[2]
				exit = 0
			}

		case <-time.After(1 * time.Second):
			dataRead = 0
		}
	}
	return returnValue, err
}

func main() {
	rsdp_value, _ := getRSDP("/dev/kmsg")
	f, err := os.Create("/tmp/rsdp")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, _ = fmt.Fprintf(f, " acpi_rsdp=%s ", rsdp_value)
}
