// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// free reports usage information for physical memory and swap space.
//
// Synopsis:
//     free [-k] [-m] [-g] [-t] [-h] [-json]
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
//     -h: display the values in human-readable form
//     -json: use JSON output
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

var (
	humanOutput = flag.Bool("h", false, "Human output: show automatically the shortest three-digits unit")
	inBytes     = flag.Bool("b", false, "Express the values in bytes")
	inKB        = flag.Bool("k", false, "Express the values in kibibytes (default)")
	inMB        = flag.Bool("m", false, "Express the values in mebibytes")
	inGB        = flag.Bool("g", false, "Express the values in gibibytes")
	inTB        = flag.Bool("t", false, "Express the values in tebibytes")
	toJSON      = flag.Bool("json", false, "Use JSON for output")
)

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

// FreeConfig is a structure used to configure the behaviour of Free()
type FreeConfig struct {
	Unit        unit
	HumanOutput bool
	ToJSON      bool
}

// the following types are used for JSON serialization
type mainMemInfo struct {
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
	Free      uint64 `json:"free"`
	Shared    uint64 `json:"shared"`
	Cached    uint64 `json:"cached"`
	Buffers   uint64 `json:"buffers"`
	Available uint64 `json:"available"`
}

type swapInfo struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}

// MemInfo represents the main memory and swap space information in a structured
// manner, suitable for JSON encoding.
type MemInfo struct {
	Mem  mainMemInfo `json:"mem"`
	Swap swapInfo    `json:"swap"`
}

type meminfomap map[string]uint64

// missingRequiredFields checks if any of the specified fields are present in
// the input map.
func missingRequiredFields(m meminfomap, fields []string) bool {
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
// woud return the string "10 kB". Note that the decimal part is truncated, not
// rounded, so the values are guaranteed to be "at least X"
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
		decimal = ((value - (value >> shift << shift)) >> (shift - 10)) * 1000 / 1024 / 100
	}
	return fmt.Sprintf("%v.%v%v",
		value>>shift,
		decimal,
		units[shift/10],
	)
}

// formatValueByConfig formats a size in bytes in the appropriate unit,
// depending on whether FreeConfig specifies a human-readable format or a
// specific unit
func formatValueByConfig(value uint64, config *FreeConfig) string {
	if config.HumanOutput {
		return humanReadableValue(value)
	}
	// units and decimal part are not printed when a unit is explicitly specified
	return fmt.Sprintf("%v", value>>config.Unit)
}

// Free prints physical memory and swap space information. The fields will be
// expressed with the specified unit (e.g. KB, MB)
func Free(config *FreeConfig) error {
	m, err := meminfo()
	if err != nil {
		log.Fatal(err)
	}
	mmi, err := getMainMemInfo(m, config)
	if err != nil {
		return err
	}
	si, err := getSwapInfo(m, config)
	if err != nil {
		return err
	}
	mi := MemInfo{Mem: *mmi, Swap: *si}
	if config.ToJSON {
		jsonData, err := json.Marshal(mi)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Printf("              total        used        free      shared  buff/cache   available\n")
		fmt.Printf("%-7s %11v %11v %11v %11v %11v %11v\n",
			"Mem:",
			formatValueByConfig(mmi.Total, config),
			formatValueByConfig(mmi.Used, config),
			formatValueByConfig(mmi.Free, config),
			formatValueByConfig(mmi.Shared, config),
			formatValueByConfig(mmi.Buffers+mmi.Cached, config),
			formatValueByConfig(mmi.Available, config),
		)
		fmt.Printf("%-7s %11v %11v %11v\n",
			"Swap:",
			formatValueByConfig(si.Total, config),
			formatValueByConfig(si.Used, config),
			formatValueByConfig(si.Free, config),
		)
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
		log.Fatal("Options -k, -m, -g, -t and -h are mutually exclusive")
	}
	config := FreeConfig{ToJSON: *toJSON}
	if *humanOutput {
		config.HumanOutput = true
	} else {
		switch {
		case *inBytes:
			config.Unit = B
		case *inKB:
			config.Unit = KB
		case *inMB:
			config.Unit = MB
		case *inGB:
			config.Unit = GB
		case *inTB:
			config.Unit = TB
		}
	}

	if err := Free(&config); err != nil {
		log.Fatal(err)
	}
}
