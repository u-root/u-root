// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cbmem prints out coreboot mem table information in JSON by default,
// and also implements the basic cbmem -list and -console commands.
// TODO: checksum tables.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

// The C version of cbmem has a complex function to list
// numbers in xx,xxx,xxx form. I personally think this
// is a terrible idea (what about the EU among other things?)
// If you decide yor really want this, then don't write a function,
// do this.
//    "golang.org/x/text/language"
//    "golang.org/x/text/message"
//    p := message.NewPrinter(language.English)
//    p.Printf("%d\n", 1000)

var (
	mem                 = flag.String("mem", "/dev/mem", "file for coreboot image")
	debug               = func(string, ...interface{}) {}
	addr                int64
	size                int
	console             bool
	coverage            bool
	list                bool
	hexdump             bool
	timestamps          bool
	parseabletimestamps bool
	verbose             bool
	version             bool
	dumpJSON            = flag.Bool("json", true, "Output tables in JSON format")
)

func init() {
	flag.BoolVar(&console, "console", false, "print cbmem console")
	flag.BoolVar(&coverage, "coverage", false, "dump coverage information")
	flag.BoolVar(&list, "list", false, "print cbmem table of contents")
	flag.BoolVar(&hexdump, "hexdump", false, "print hexdump of cbmem area")
	flag.BoolVar(&timestamps, "timestamps", false, "print timestamp information")
	flag.BoolVar(&parseabletimestamps, "parseable-timestamps", false, "print parseable timestamps")
	flag.BoolVar(&verbose, "verbose", false, "verbose (debugging) output")
	flag.BoolVar(&version, "version", false, "print the version")
	flag.BoolVar(&console, "c", false, "print cbmem console")
	flag.BoolVar(&coverage, "C", false, "dump coverage information")
	flag.BoolVar(&list, "l", false, "print cbmem table of contents")
	flag.BoolVar(&hexdump, "x", false, "print hexdump of cbmem area")
	flag.BoolVar(&timestamps, "t", false, "print timestamp information")
	flag.BoolVar(&parseabletimestamps, "T", false, "print parseable timestamps")
	flag.BoolVar(&verbose, "v", false, "verbose (debugging) output")
	flag.BoolVar(&version, "V", false, "print the version")
}

func parseCBtable(r io.ReaderAt, address int64, sz int) (*CBmem, error) {
	debug("Looking for coreboot table at %v %v bytes", address, sz)
	var (
		i     int64
		lbh   Header
		found = fmt.Errorf("No cb table found")
		cbmem = &CBmem{StringVars: make(map[string]string)}
	)

	for i = address; i < address+0x1000 && found != nil; i += 0x10 {
		readOne(r, &lbh, i)
		debug("header is %q", lbh)
		if string(lbh.Signature[:]) != "LBIO" {
			debug("no LBIO at %v", i)
			continue
		}
		if lbh.HeaderSz == 0 {
			debug("HeaderSz is 0 at %v", i)
		}
		// TODO: checksum the header.
		// Although I know of no case in 10 years where that
		// was useful.
		addr = i + int64(lbh.HeaderSz)
		found = nil
		debug("Found!\n")

		/* Keep reference to lbtable. */
		size = int(lbh.TableSz)
		j := addr
		for j < addr+int64(lbh.TableSz) {
			var rec Record
			debug("  coreboot table entry 0x%02x\n", rec.Tag)
			readOne(r, &rec, j)
			debug("Found Tag %s (%v) Size %v", tagNames[rec.Tag], rec.Tag, rec.Size)
			start := j
			j += int64(reflect.TypeOf(r).Size())
			n := tagNames[rec.Tag]
			switch rec.Tag {
			case LB_TAG_BOARD_ID:
				readOne(r, &cbmem.BoardID, start)
			case
				LB_TAG_VERSION,
				LB_TAG_EXTRA_VERSION,
				LB_TAG_BUILD,
				LB_TAG_COMPILE_TIME,
				LB_TAG_COMPILE_BY,
				LB_TAG_COMPILE_HOST,
				LB_TAG_COMPILE_DOMAIN,
				LB_TAG_COMPILER,
				LB_TAG_LINKER,
				LB_TAG_ASSEMBLER:
				s, err := bufio.NewReader(io.NewSectionReader(r, j, 65536)).ReadString(0)
				if err != nil {
					log.Fatalf("Trying to read string for %s: %v", n, err)
				}
				cbmem.StringVars[n] = s[:len(s)-1]
			case LB_TAG_SERIAL:
				var s serialEntry
				readOne(r, &s, start)
				cbmem.UART = append(cbmem.UART, s)

			case LB_TAG_CONSOLE:
				var c uint32
				readOne(r, &c, j)
				cbmem.Consoles = append(cbmem.Consoles, consoleNames[c])
			case LB_TAG_VERSION_TIMESTAMP:
				readOne(r, &cbmem.VersionTimeStamp, start)
			case LB_TAG_BOOT_MEDIA_PARAMS:
				readOne(r, &cbmem.BootMediaParams, start)
			case LB_TAG_CBMEM_ENTRY:
				var c cbmemEntry
				readOne(r, &c, start)
				cbmem.CBMemory = append(cbmem.CBMemory, c)
			case LB_TAG_MEMORY:
				debug("    Found memory map.\n")
				cbmem.Memory = &memoryEntry{Record: rec}
				nel := (int64(cbmem.Memory.Size) - (j - start)) / int64(reflect.TypeOf(memoryRange{}).Size())
				cbmem.Memory.Maps = make([]memoryRange, nel)
				readOne(r, cbmem.Memory.Maps, j)
			case LB_TAG_TIMESTAMPS:
				debug("    Found timestamp table.\n")
				var t timestampEntry
				readOne(r, &t, start)
				cbmem.TimeStamps = append(cbmem.TimeStamps, t)
			case LB_TAG_MAINBOARD:
				// The mainboard entry is a bit weird.
				// There is a byte after the Record
				// for the Vendor Index and a byte after
				// that for the Part Number Index.
				// In general, the vx is 0, and it's also
				// null terminated. The struct is a bit
				// over-general, actually, and the indexes
				// can be safely ignored.
				cbmem.MainBoard.Record = rec
				v, err := bufio.NewReader(io.NewSectionReader(r, j+2, 65536)).ReadString(0)
				if err != nil {
					log.Fatalf("Trying to read string for %s: %v", n, err)
				}
				p, err := bufio.NewReader(io.NewSectionReader(r, j+2+int64(len(v)), 65536)).ReadString(0)
				if err != nil {
					log.Fatalf("Trying to read string for %s: %v", n, err)
				}
				cbmem.MainBoard.Vendor = v[:len(v)-1]
				cbmem.MainBoard.PartNumber = p[:len(p)-1]
			case LB_TAG_HWRPB:
				readOne(r, &cbmem.Hwrpb, start)
			case LB_TAG_CBMEM_CONSOLE:
				var c = &memconsoleEntry{Record: rec}
				debug("    Found cbmem console, %d byte record.\n", c.Size)
				readOne(r, &c.CSize, j)
				j += int64(reflect.TypeOf(c.CSize).Size())
				readOne(r, &c.Cursor, j)
				j += int64(reflect.TypeOf(c.Cursor).Size())
				if c.CSize > c.Cursor {
					c.CSize = c.Cursor
				}
				debug("CSize is %d, and Cursor is at %d", c.CSize, c.Cursor)
				c.Data = make([]byte, c.CSize)
				readOne(r, c.Data, j)
				cbmem.MemConsole = c
			case LB_TAG_FORWARD:
				var newTable int64
				readOne(r, &newTable, j)
				debug("Forward to %08x", newTable)
				return parseCBtable(r, newTable, 1048576)
			default:
				if n, ok := tagNames[rec.Tag]; ok {
					log.Printf("Ignoring record %v", n)
					cbmem.Ignored = append(cbmem.Ignored, n)
					break
				}
				log.Printf("Unknown tag record %v", r)
				cbmem.Unknown = append(cbmem.Unknown, rec.Tag)
				break

			}
			j = start + int64(rec.Size)
		}
	}
	return cbmem, found
}

func DumpMem(cbmem *CBmem) {
	if cbmem.Memory == nil {
		fmt.Printf("No cbmem table name")
	}
	m := cbmem.Memory.Maps
	if len(m) == 0 {
		fmt.Printf("No cbmem map entries")
	}
	fmt.Printf("%19s %8s %8s\n", "Name", "Start", "Size")
	for _, e := range m {
		fmt.Printf("%19s %08x %08x\n", memTags[e.Mtype], e.Start, e.Size)
	}
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("cbmem in Go, a superset of cbmem v1.1 from coreboot")
		os.Exit(0)
	}
	if verbose {
		debug = log.Printf
	}

	mf, err := os.Open(*mem)
	if err != nil {
		log.Fatal(err)
	}

	var cbmem *CBmem
	for _, addr := range []int64{0, 0xf0000} {
		if cbmem, err = parseCBtable(mf, addr, TableSize); err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("Reading coreboot table: %v", err)
	}
	if *dumpJSON {
		b, err := json.MarshalIndent(cbmem, "", "\t")
		if err != nil {
			log.Fatalf("json marshal: %v", err)
		}
		fmt.Printf("%s\n", b)
	}
	// list is kind of misnamed I think. It really just prints
	// memory table entries.
	if list {
		DumpMem(cbmem)
	}
	if console && cbmem.MemConsole != nil {
		fmt.Printf("Console is %d bytes and cursor is at %d\n", len(cbmem.MemConsole.Data), cbmem.MemConsole.Cursor)
		fmt.Printf("%s%s", cbmem.MemConsole.Data[cbmem.MemConsole.Cursor:], cbmem.MemConsole.Data[0:cbmem.MemConsole.Cursor])
	}

}
