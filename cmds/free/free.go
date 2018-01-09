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

const meminfoFile = "/proc/meminfo"

var humanOutput = flag.Bool("h", false, "Human output: show automatically the shortest three-digits unit")
var inBytes = flag.Bool("b", false, "Express the values in bytes")
var inKB = flag.Bool("k", false, "Express the values in kibibytes (default)")
var inMB = flag.Bool("m", false, "Express the values in mebibytes")
var inGB = flag.Bool("g", false, "Express the values in gibibytes")
var inTB = flag.Bool("t", false, "Express the values in tebibytes")

type unit uint

const (
	// B is bytes
	B unit = 0
	// KB is kibibytes
	KB = 10
	// MB is mebibytes
	MB = 20
	// GB is gibibytes
	GB = 30
	// TB is tebibytes
	TB = 40
)

var units = [...]string{"B", "K", "M", "G", "T"}

type freeConfig struct {
	Unit        unit
	HumanOutput bool
}

type memInfo map[string]uint64

// meminfo returns a mapping that represents the fields contained in
// /proc/meminfo
func meminfo() (memInfo, error) {
	buf, err := ioutil.ReadFile(meminfoFile)
	if err != nil {
		return nil, err
	}
	return meminfoFromBytes(buf)
}

// meminfoFromBytes returns a mapping that represents the fields contained in a
// byte stream with a content compatible with /proc/meminfo
func meminfoFromBytes(buf []byte) (memInfo, error) {
	ret := make(memInfo, 0)
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		kv := bytes.SplitN(line, []byte{':'}, 2)
		if len(kv) != 2 {
			// invalid line?
			continue
		}
		key := string(kv[0])
		tokens := bytes.SplitN(bytes.TrimSpace(kv[1]), []byte{' '}, 2)
		if len(tokens) > 0 {
			value, err := strconv.ParseUint(string(tokens[0]), 10, 64)
			if err != nil {
				return nil, err
			}
			ret[key] = value
		}
	}
	return ret, nil
}

// missingRequiredFields checks if any of the specified fields are present in
// the input map.
func missingRequiredFields(m memInfo, fields []string) bool {
	for _, f := range fields {
		if _, ok := m[f]; !ok {
			log.Printf("Missing field '%v'", f)
			return true
		}
	}
	return false
}

// humanReadableValue returns a string representing the input value, treated as
// a size in bytes, interpreted in a human readable form. E.g. the number 10240
// woud return the string "10 kB"
func humanReadableValue(value uint64) string {
	v := value
	// bits to shift. 0 means bytes, 10 means kB, and so on. 40 is the highest
	// and it means tB
	var shift uint
	for {
		if shift >= uint(len(units)*10) {
			// 4 means tebibyte, we don't go further
			break
		}
		if v/1024 < 1 {
			break
		}
		v /= 1024
		shift += 10
	}
	var decimal uint64
	if shift > 0 {
		// no rounding. Is there a better way to do this?
		decimal = ((value - (value >> shift << shift)) >> (shift - 10)) / 100
	}
	return fmt.Sprintf("%v.%v%v",
		value>>shift,
		decimal,
		units[shift/10],
	)
}

// formatValueByConfig formats a size in bytes in the appropriate unit,
// depending on whether freeConfig specifies a human-readable format or a
// specific unit
func formatValueByConfig(value uint64, config *freeConfig) string {
	if config.HumanOutput {
		return humanReadableValue(value)
	}
	// units and decimal part are not printed when a unit is explicitly specified
	return fmt.Sprintf("%v", value>>config.Unit)
}

// printMem prints the physical memory information in the specified units. Only
// the relevant fields will be used from the input map.
func printMem(m memInfo, config *freeConfig) error {
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
	memTotal := m["MemTotal"] << KB
	memFree := m["MemFree"] << KB
	memShared := m["Shmem"] << KB
	memCached := (m["Cached"] + m["SReclaimable"]) << KB
	memBuffers := (m["Buffers"]) << KB
	memUsed := memTotal - memFree - memCached - memBuffers
	if memUsed < 0 {
		memUsed = memTotal - memFree
	}
	memAvailable := m["MemAvailable"] << KB

	fmt.Printf("%-7s %11v %11v %11v %11v %11v %11v\n",
		"Mem:",
		formatValueByConfig(memTotal, config),
		formatValueByConfig(memUsed, config),
		formatValueByConfig(memFree, config),
		formatValueByConfig(memShared, config),
		formatValueByConfig(memBuffers+memCached, config),
		formatValueByConfig(memAvailable, config),
	)
	return nil
}

// printSwap prints the swap space information in the specified units. Only the
// relevant fields will be used from the input map.
func printSwap(m memInfo, config *freeConfig) error {
	fields := []string{
		"SwapTotal",
		"SwapFree",
	}
	if missingRequiredFields(m, fields) {
		return fmt.Errorf("Missing required fields from meminfo")
	}
	// These values are expressed in kibibytes, convert to the desired unit
	swapTotal := m["SwapTotal"] << KB
	swapUsed := (m["SwapTotal"] - m["SwapFree"]) << KB
	swapFree := m["SwapFree"] << KB
	fmt.Printf("%-7s %11v %11v %11v\n",
		"Swap:",
		formatValueByConfig(swapTotal, config),
		formatValueByConfig(swapUsed, config),
		formatValueByConfig(swapFree, config),
	)
	return nil
}

// Free prints physical memory and swap space information. The fields will be
// expressed with the specified unit (e.g. KB, MB)
func Free(config *freeConfig) error {
	m, err := meminfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("              total        used        free      shared  buff/cache   available\n")
	if err := printMem(m, config); err != nil {
		return err
	}
	if err := printSwap(m, config); err != nil {
		return err
	}
	return nil
}

// validateUnits checks that only one option of -b, -k, -m, -g, -t or -h has been
// specified on the command line
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
	if *humanOutput {
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
	var config freeConfig
	if *humanOutput {
		config.HumanOutput = true
	} else {
		var unit unit = KB
		if *inBytes {
			unit = B
		} else if *inKB {
			unit = KB
		} else if *inMB {
			unit = MB
		} else if *inGB {
			unit = GB
		} else if *inTB {
			unit = TB
		}
		config.Unit = unit
	}

	if err := Free(&config); err != nil {
		log.Fatal(err)
	}
}
