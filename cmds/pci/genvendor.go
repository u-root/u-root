
// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	debug = func(s string, arg ...interface{})  {}//{log.Printf(s, arg...)}
	inTable, inVID, inDEVS, inDID, inSUB bool
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

func closeSub() {
	if !inSUB {
		return
	}
	fmt.Printf("\t\t}, // subdevice\n")
	inSUB = false
}

func closeDev() {
	closeSub()
	if ! inDID {
		return
	}
	fmt.Printf("\t}, // device \n")
	inDID = false
}

func closeDevs() {
	closeDev()
	if ! inDEVS {
		return
	}
	fmt.Printf("\t}, // devices\n")
	inDEVS = false
}

func closeVendor() {
	closeDevs()
	if ! inVID {
		return
	}
	inVID = false
	fmt.Printf("}, // vendor\n")
}

func closeTable() {
	if ! inTable {
		return
	}
	closeVendor()
	fmt.Printf("} // table\n")
}
	
func newSub(s string) {
	if ! inDID {
		log.Fatalf("%s: found a sub but not in a device", s)
	}
	if (! inSUB) {
		fmt.Printf("Sub: []SubVendor{\n")
	}
	inSUB = true
	
	fmt.Printf("\t\tSubVendor{Ven: 0x%s, Dev: 0x%s, Name: %q},\n", s[3:7], s[7:11], s[13:])
}

func newDev(s string) {
	if ! inVID {
		log.Fatalf("%s: found a dev but not in a vendor", s)
	}
	if ! inDEVS {
		fmt.Printf("Devs: map[DID]Device {\n")
	}
	closeDev()
	inDEVS = true
	inDID = true
	
	fmt.Printf("\n\t0x%s: Device{Name: %q,", s[1:5], s[6:])
}

func newVendor(v string) {
	closeVendor()
	fmt.Printf("0x%v: Vendor{Name: \"%v\", ", v[0:4], v[6:])
	inVID = true
}

func main() {
	var line string
	var lineno int
	defer func() {
		log.Printf("well, we're leaving on line %v %v", lineno, line)

		switch err := recover().(type) {
		case nil:
		case error:
			log.Fatalf("Bummer: %v", err)
		default:
			log.Fatalf("unexpected panic value: %T(%v)", err, err)
		}
	}()

	fmt.Printf("package main\n")
	s := bufio.NewScanner(os.Stdin)

	// Simple state machine.
	for s.Scan() {
		line = s.Text()
		lineno++
		debug("Line %d: %v\n", lineno, line)
		switch {
		case len(line) == 0, line[0] == '#':
			debug("commend\n")
			continue
		case isHex(line[0]) && isHex(line[1]):
			debug("Vendor %v, inVID %v", line, inVID)
			if ! inTable {
				fmt.Printf("var vendor = map[VID]Vendor {\n")
				inTable = true
			}
			newVendor(line)
		case line[0] == '\t' && isHex(line[1]):
			debug("Device %v", line)
			newDev(line)
		case line[0:2] == "\t\t" && isHex(line[2]):
			debug("SubDevice %v", line)
			newSub(line)
		case line[0] == 'C':
			debug("Class %v", line)
			closeTable()
			os.Exit(0)
		}
	}
	closeTable()
}

