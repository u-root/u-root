// Copyright 2012-2017 the u-root Authors. All rights reserved
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
// Notes:
//	The code is looking for RSDP address into /dev/kmsg 
//
// Example:
//	rsdp 	- Start the script and save the founded value into /tmp/rsdp

package main

import (
        "strings"
        "log"
	"os"
	"bufio"
	"time"
)

func getRSDP(path string) (string,error) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var returnValue string
	channel := make(chan  string)

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

	var data_read int
	var exit int
	data_read = 1
	exit = 0
	for (data_read == 1 && exit == 0) {
		select {
		    case res := <-channel:
			if strings.Contains(res, "RSDP") {
				returnValue = strings.Split(res," ")[2]
				exit = 0
	                }

		    case <-time.After(1 * time.Second):
			data_read = 0
		}
	}
        return returnValue,err
}

func main(){
	rsdp_value,_:=getRSDP("/dev/kmsg")
	f,err := os.Create("/tmp/rsdp")
	if err != nil {
                log.Fatal(err)
        }
	defer f.Close()
	writer := bufio.NewWriter(f)
	writer.WriteString(rsdp_value)
	writer.Flush()
}

