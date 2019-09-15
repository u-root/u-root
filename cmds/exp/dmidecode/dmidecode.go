// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/smbios"
)

var (
	flagFromDump = flag.String("from-dump", "", `Read the DMI data from a binary file previously generated using --dump-bin.`)
	flagType     = flag.StringSliceP("type", "t", nil, `Only  display  the  entries of type TYPE. TYPE can be either a DMI type number, or a comma-separated list of type numbers, or a keyword from the following list: bios, system, baseboard, chassis, processor, memory, cache, connector, slot. If this option is used more than once, the set of displayed entries will be the union of all the given types. If TYPE is not provided or not valid, a list of all valid keywords is printed and dmidecode exits with an error.`)
	// NB: When adding flags, update resetFlags in dmidecode_test.
)

var (
	typeGroups = map[string][]uint8{
		"bios":      {0, 13},
		"system":    {1, 12, 15, 23, 32},
		"baseboard": {2, 10, 41},
		"chassis":   {3},
		"processor": {4},
		"memory":    {5, 6, 16, 17},
		"cache":     {7},
		"connector": {8},
		"slot":      {9},
	}
)

type dmiDecodeError struct {
	error
	code int
}

// parseTypeFilter parses the --type argument(s) and returns a set of types taht should be included.
func parseTypeFilter(typeStrings []string) (map[smbios.TableType]bool, error) {
	types := map[smbios.TableType]bool{}
	for _, ts := range typeStrings {
		if tg, ok := typeGroups[strings.ToLower(ts)]; ok {
			for _, t := range tg {
				types[smbios.TableType(t)] = true
			}
		} else {
			u, err := strconv.ParseUint(ts, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("Invalid type: %s", ts)
			}
			types[smbios.TableType(uint8(u))] = true
		}
	}
	return types, nil
}

func dmiDecode(textOut io.Writer) *dmiDecodeError {
	typeFilter, err := parseTypeFilter(*flagType)
	if err != nil {
		return &dmiDecodeError{code: 2, error: fmt.Errorf("invalid --type: %s", err)}
	}
	fmt.Fprintf(textOut, "# dmidecode-go\n") // TODO: version.
	entry, data, err := getData(textOut, *flagFromDump, "/sys/firmware/dmi/tables")
	if err != nil {
		return &dmiDecodeError{code: 1, error: fmt.Errorf("error parsing loading data: %s", err)}
	}
	si, err := smbios.ParseInfo(entry, data)
	if err != nil {
		return &dmiDecodeError{code: 1, error: fmt.Errorf("error parsing data: %s", err)}
	}
	if si.Entry64 != nil {
		fmt.Fprintf(textOut, "SMBIOS %d.%d.%d present.\n\n", si.GetSMBIOSMajorVersion(), si.GetSMBIOSMinorVersion(), si.GetSMBIOSDocRev())
	} else {
		fmt.Fprintf(textOut, "SMBIOS %d.%d present.\n", si.GetSMBIOSMajorVersion(), si.GetSMBIOSMinorVersion())
	}
	if si.Entry32 != nil {
		fmt.Fprintf(textOut, "%d structures occupying %d bytes.\n", si.Entry32.NumberOfStructs, si.Entry32.StructTableLength)
	}
	fmt.Fprintf(textOut, "\n")
	for _, t := range si.Tables {
		if len(typeFilter) != 0 && !typeFilter[t.Type] {
			continue
		}
		pt, err := smbios.ParseTypedTable(t)
		if err != nil {
			if err != smbios.ErrUnsupportedTableType {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			// Print as raw table
			pt = t
		}
		fmt.Fprintf(textOut, "%s\n\n", pt)
	}
	return nil
}

func main() {
	flag.Parse()
	err := dmiDecode(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(err.code)
	}
}
