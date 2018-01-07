// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// free displays information about total, used and available physical memory and
// swap space.
//
// Synopsis:
//     free [-k] [-m] [-g] [-t]
//
// Description:
//     Read memory information from /proc/meminfo and display a summary for
//     physical memory and swap space. The unit options use powers of 1024.
//
// Options:
//     -k: display the values in kibibytes
//     -m: display the values in mebibytes
//     -g: display the values in gibibytes
//     -t: display the values in tebibytes

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

const MEMINFO_FILE = "/proc/meminfo"

var inBytes = flag.Bool("b", false, "Express the values in bytes")
var inKB = flag.Bool("k", false, "Express the values in kibibytes (default)")
var inMB = flag.Bool("m", false, "Express the values in mebibytes")
var inGB = flag.Bool("g", false, "Express the values in gibibytes")
var inTB = flag.Bool("t", false, "Express the values in tebibytes")

type Unit uint

const (
	B  Unit = 0
	KB      = 10
	MB      = 20
	GB      = 30
	TB      = 40
)

func meminfo() (map[string]int, error) {
	buf, err := ioutil.ReadFile(MEMINFO_FILE)
	if err != nil {
		return nil, err
	}
	return meminfoFromBytes(buf)
}

func meminfoFromBytes(buf []byte) (map[string]int, error) {
	ret := make(map[string]int, 0)
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		kv := bytes.SplitN(line, []byte{':'}, 2)
		if len(kv) != 2 {
			// invalid line?
			continue
		}
		key := string(kv[0])
		tokens := bytes.SplitN(bytes.TrimSpace(kv[1]), []byte{' '}, 2)
		if len(tokens) > 0 {
			value, err := strconv.Atoi(string(tokens[0]))
			if err != nil {
				return nil, err
			}
			ret[key] = value
		}
	}
	return ret, nil
}

func missingRequiredFields(m map[string]int, fields []string) bool {
	for _, f := range fields {
		if _, ok := m[f]; !ok {
			log.Printf("Missing field '%v'", f)
			return true
		}
	}
	return false
}

func printMem(m map[string]int, unit Unit) error {
	fields := []string{
		"MemTotal",
		"MemFree",
		"Buffers",
		"Cached",
		"SReclaimable",
		"MemAvailable",
	}
	if missingRequiredFields(m, fields) {
		return fmt.Errorf("Missing required fields from meminfo")
	}

	// These values are expressed in kibibytes, convert to the desired unit
	memTotal := (m["MemTotal"] << KB) >> unit
	memFree := (m["MemFree"] << KB) >> unit
	memShared := (m["Shmem"] << KB) >> unit
	memCached := ((m["Cached"] + m["SReclaimable"]) << KB) >> unit
	memBuffers := ((m["Buffers"]) << KB) >> unit
	memUsed := memTotal - memFree - memCached - memBuffers
	if memUsed < 0 {
		memUsed = memTotal - memFree
	}
	memAvailable := (m["MemAvailable"] << KB) >> unit

	fmt.Printf("%-7s %11d %11d %11d %11d %11d %11d\n",
		"Mem:",
		memTotal, memUsed, memFree, memShared, memBuffers+memCached, memAvailable)
	return nil
}

func printSwap(m map[string]int, unit Unit) error {
	fields := []string{
		"SwapTotal",
		"SwapFree",
	}
	if missingRequiredFields(m, fields) {
		return fmt.Errorf("Missing required fields from meminfo")
	}
	// These values are expressed in kibibytes, convert to the desired unit
	swapTotal := (m["SwapTotal"] << KB) >> unit
	swapUsed := ((m["SwapTotal"] - m["SwapFree"]) << KB) >> unit
	swapFree := (m["SwapFree"] << KB) >> unit
	fmt.Printf("%-7s %11d %11d %11d\n", "Swap:", swapTotal, swapUsed, swapFree)
	return nil
}

func Free(unit Unit) error {
	m, err := meminfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("              total        used        free      shared  buff/cache   available\n")
	if err := printMem(m, unit); err != nil {
		return err
	}
	if err := printSwap(m, unit); err != nil {
		return err
	}
	return nil
}

func validateUnits() bool {
	count := 0
	if *inBytes {
		count++
	}
	if *inKB {
		count++
	}
	if *inMB {
		count++
	}
	if *inGB {
		count++
	}
	if *inTB {
		count++
	}
	if count > 1 {
		return false
	}
	return true
}

func main() {
	flag.Parse()
	if !validateUnits() {
		log.Fatal("Options -k, -m, -g and -t are mutually exclusive")
	}
	var unit Unit
	if *inBytes {
		unit = B
	} else if *inKB {
		unit = KB
	} else if *inMB {
		unit = MB
	} else if *inGB {
		unit = GB
	}

	if err := Free(unit); err != nil {
		log.Fatal(err)
	}
}
