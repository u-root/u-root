// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cbmem prints out coreboot mem table information in JSON by default,
// and also implements the basic cbmem -list and -console commands.
// TODO: checksum tables.
package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"syscall"
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
	memFile             *os.File
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
	dumpJSON            bool
)

//
// usage: /home/rminnich/bin/cbmem [-cCltTLxVvh?]
//   -c | --console:                   print cbmem console

func init() {
	const longfmt = "-%s | --%s:%s%s (default %v)\n"
	var (
		ushort = "cbmem [h?"
		ulong  string
	)

	for _, f := range []struct {
		b     *bool
		def   bool
		short string
		long  string
		help  string
		tab   string
	}{
		{&console, false, "c", "console", "print cbmem console", "\t\t\t"},
		{&coverage, false, "C", "coverage", "dump coverage information", "\t\t"},
		{&list, false, "l", "list", "print cbmem table of contents", "\t\t\t"},
		{&hexdump, false, "x", "hexdump", "print hexdump of cbmem area", "\t\t\t"},
		{&timestamps, false, "t", "timestamps", "print timestamp information", "\t\t"},
		{&parseabletimestamps, false, "p", "parseable-timestamps", "print parseable timestamps", "\t"},
		{&verbose, false, "v", "verbose", "verbose (debugging) output", "\t\t\t"},
		{&version, false, "V", "version", "print version information", "\t\t\t"},
		{&dumpJSON, false, "j", "json", "Output tables in JSON format", "\t\t\t"},
	} {
		flag.BoolVar(f.b, f.short, f.def, f.help)
		flag.BoolVar(f.b, f.long, f.def, f.help)
		ushort += f.short
		ulong += fmt.Sprintf(longfmt, f.short, f.long, f.tab, f.help, f.def)
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s]\n\n%s", ushort, ulong)
		os.Exit(1)
	}

}

func mapit(addr int64, sz int) ([]byte, error) {
	b, err := syscall.Mmap(int(memFile.Fd()), 0, int(addr)+sz, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap %d bytes at %#x: %v", sz, addr, err)
	}
	return b, nil
}

func parseCBtable(address int64, sz int) (*CBmem, error) {
	// Note:
	// Code uses mmmap
	// address is 0-relative
	// potential for a large mmap is "large"
	// it's nice if we don't have to mmap a slice starting
	// at address and then tweak address everywhere.
	// This is easy: always mmap at 0, and then we can use address in the returned
	// slice directly.
	// mmap is backed by a VMA, and it is not populated until an address is used.
	// Linux VMAs are very cheap for this sort of thing.
	// So if you're worried about using offset size with what might be a 2G range,
	// no need to worry.
	b, err := mapit(0, int(address)+sz)
	if err != nil {
		return nil, err
	}
	var r io.ReaderAt = bytes.NewReader(b)
	debug("Looking for coreboot table at %#08x %d bytes", address, sz)
	var (
		i     int64
		lbh   Header
		found = fmt.Errorf("No cb table found")
		cbmem = &CBmem{StringVars: make(map[string]string)}
	)

	for i = address; i < address+0x1000 && found != nil; i += 0x10 {
		readOne(r, &lbh, i)
		debug("header is %s", lbh.String())
		if string(lbh.Signature[:]) != "LBIO" {
			debug("no LBIO at %#08x", i)
			continue
		}
		if lbh.HeaderSz == 0 {
			debug("HeaderSz is 0 at %#08x", i)
		}
		// TODO: checksum the header.
		// Although I know of no case in 10 years where that
		// was useful.
		addr = i + int64(lbh.HeaderSz)
		found = nil
		debug("Found at %#08x!", addr)

		/* Keep reference to lbtable. */
		size = int(lbh.TableSz)
		j := addr
		debug("Process %d entires", lbh.TableEntries)
		for j < addr+int64(lbh.TableSz) {
			var rec Record
			debug("\tcoreboot table entry 0x%02x\n", rec.Tag)
			readOne(r, &rec, j)
			debug("\tFound Tag %s (%v)@%#08x Size %v", tagNames[rec.Tag], rec.Tag, j, rec.Size)
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
				LB_TAG_ASSEMBLER,
				LB_TAG_PLATFORM_BLOB_VERSION:
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

				// "Nobody knew consoles could be so hard."
			case LB_TAG_CBMEM_CONSOLE:
				var c = &memconsoleEntry{Record: rec}
				debug("    Found cbmem console(%#x), %d byte record.\n", rec, c.Size)
				readOne(r, &c.Address, j)
				j += int64(reflect.TypeOf(c.Address).Size())
				debug("    console data is at %#x", c.Address)
				cbcons := int64(c.Address)
				// u32 size;
				// u32 cursor;
				// u8  body[0];
				readOne(r, &c.Size, cbcons)
				cbcons += int64(reflect.TypeOf(c.Size).Size())
				readOne(r, &c.Cursor, cbcons)
				cbcons += int64(reflect.TypeOf(c.Cursor).Size())
				debug("CSize is %d, and Cursor is at %d", c.CSize, c.Cursor)

				curse := c.Cursor & CBMC_CURSOR_MASK
				sz := c.Size
				if (c.Cursor&CBMC_OVERFLOW) == 0 && curse < c.Size {
					sz = curse
				}

				debug("CSize is %d, and Cursor is at %d", c.CSize, c.Cursor)
				data := make([]byte, sz)
				// TODO: deal with wrap.
				readOne(r, data, cbcons)
				c.Data = string(data)
				cbmem.MemConsole = c

			case LB_TAG_FORWARD:
				var newTable int64
				readOne(r, &newTable, j)
				debug("Forward to %08x", newTable)
				return parseCBtable(newTable, 1048576)
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

func DumpMem(cbmem *CBmem, w io.Writer) {
	if cbmem.Memory == nil {
		fmt.Fprintf(w, "No cbmem table name")
	}
	m := cbmem.Memory.Maps
	if len(m) == 0 {
		fmt.Fprintf(w, "No cbmem map entries")
	}
	fmt.Fprintf(w, "%19s %8s %8s\n", "Name", "Start", "Size")
	for _, e := range m {
		fmt.Fprintf(w, "%19s %08x %08x\n", memTags[e.Mtype], e.Start, e.Size)
		if hexdump && e.Mtype == LB_MEM_TABLE {
			b, err := mapit(int64(e.Start), int(e.Size))
			if err != nil {
				log.Print(err)
				continue
			}
			// The hexdump does a lot of what we want, but not all of
			// what we want. In particular, we'd like better control of
			// what is printed with the offset. So ... hackery.
			out := ""
			same := 0
			for i := e.Start; i < e.Start+e.Size; i += 16 {
				s := hex.Dump(b[i : i+16])[10:]
				// If it's the same as the previous, increment same
				if s == out {
					if same == 0 {
						fmt.Fprintf(w, "...\n")
					}
					same++
					continue
				}
				same = 0
				out = s
				fmt.Fprintf(w, "%08x: %s", i, s)
			}
		}
	}
}

//go:generate go run gen/gen.go -apu2

func main() {
	var err error
	flag.Parse()
	if version {
		fmt.Println("cbmem in Go, including JSON output")
		os.Exit(0)
	}
	if verbose {
		debug = log.Printf
	}

	if memFile, err = os.Open(*mem); err != nil {
		log.Fatal(err)
	}

	var cbmem *CBmem
	for _, addr := range []int64{0, 0xf0000} {
		if cbmem, err = parseCBtable(addr, 0x10000); err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("Reading coreboot table: %v", err)
	}
	if dumpJSON {
		b, err := json.MarshalIndent(cbmem, "", "\t")
		if err != nil {
			log.Fatalf("json marshal: %v", err)
		}
		fmt.Printf("%s\n", b)
	}
	// list is kind of misnamed I think. It really just prints
	// memory table entries.
	if list || hexdump {
		DumpMem(cbmem, os.Stdout)
	}
	if console && cbmem.MemConsole != nil {
		//fmt.Printf("Console is %d bytes and cursor is at %d\n", len(cbmem.MemConsole.Data), cbmem.MemConsole.Cursor)
		fmt.Printf("%s%s", cbmem.MemConsole.Data[cbmem.MemConsole.Cursor:], cbmem.MemConsole.Data[0:cbmem.MemConsole.Cursor])
	}

}
