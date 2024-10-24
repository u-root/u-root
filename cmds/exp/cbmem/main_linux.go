// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cbmem prints out coreboot mem table information in JSON by default,
// and also implements the basic cbmem -list and -console commands.
// TODO: checksum tables.
package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"text/tabwriter"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
	dumpJSON            bool
)

//
// usage: /home/rminnich/bin/cbmem [-cCltTLxVvh?]
//   -c | --console:                   print cbmem console

func init() {
	const longfmt = "-%s | --%s:%s%s (default %v)\n"
	var (
		ushort = "Usage: cbmem [h?"
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
		{&timestamps, true, "t", "timestamps", "print timestamp information (default)", "\t\t"},
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

// parseCBtable looks for a coreboot table in the range address, address + size - 1
// If it finds one it tries to parse it.
// If it found a table it returns true.
// If the parsing had an error, it returns the error.
func parseCBtable(f *os.File, address int64, sz int) (*CBmem, bool, error) {
	var found bool
	r, err := newOffsetReader(f, address, sz)
	if err != nil {
		return nil, found, err
	}
	debug("Looking for coreboot table at %#08x %d bytes", address, sz)
	var (
		i     int64
		lbh   Header
		cbmem = &CBmem{StringVars: make(map[string]string)}
	)

	for i = address; i < address+0x1000 && !found; i += 0x10 {
		if err := readOne(r, &lbh, i); err != nil {
			return nil, found, err
		}
		debug("header is %s", lbh.String())
		if string(lbh.Signature[:]) != "LBIO" {
			debug("no LBIO at %#08x", i)
			continue
		}
		if lbh.HeaderSz == 0 {
			debug("HeaderSz is 0 at %#08x", i)
		}
		debug("Found at %#08x!", i)

		// TODO: checksum the header.
		// Although I know of no case in 10 years where that
		// was useful.
		addr = i + int64(lbh.HeaderSz)
		found = true

		/* Keep reference to lbtable. */
		size = int(lbh.TableSz)
		j := addr
		debug("Process %d entires", lbh.TableEntries)
		for j < addr+int64(lbh.TableSz) {
			var rec Record
			debug("\tcoreboot table entry 0x%02x\n", rec.Tag)
			if err := readOne(r, &rec, j); err != nil {
				return nil, found, err
			}
			debug("\tFound Tag %s (%v)@%#08x Size %v", tagNames[rec.Tag], rec.Tag, j, rec.Size)
			start := j
			j += int64(reflect.TypeOf(r).Size())
			n := tagNames[rec.Tag]
			switch rec.Tag {
			case LB_TAG_BOARD_ID:
				if err := readOne(r, &cbmem.BoardID, start); err != nil {
					return nil, found, err
				}
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
					return nil, false, fmt.Errorf("trying to read string for %s: %w", n, err)
				}
				cbmem.StringVars[n] = s[:len(s)-1]
			case LB_TAG_SERIAL:
				var s serialEntry
				if err := readOne(r, &s, start); err != nil {
					return nil, found, err
				}
				cbmem.UART = append(cbmem.UART, s)

			case LB_TAG_CONSOLE:
				var c uint32
				if err := readOne(r, &c, j); err != nil {
					return nil, found, err
				}
				cbmem.Consoles = append(cbmem.Consoles, consoleNames[c])
			case LB_TAG_VERSION_TIMESTAMP:
				if err := readOne(r, &cbmem.VersionTimeStamp, start); err != nil {
					return nil, found, err
				}
			case LB_TAG_BOOT_MEDIA_PARAMS:
				if err := readOne(r, &cbmem.BootMediaParams, start); err != nil {
					return nil, found, err
				}
			case LB_TAG_CBMEM_ENTRY:
				var c cbmemEntry
				if err := readOne(r, &c, start); err != nil {
					return nil, found, err
				}
				cbmem.CBMemory = append(cbmem.CBMemory, c)
			case LB_TAG_MEMORY:
				debug("    Found memory map.\n")
				cbmem.Memory = &memoryEntry{Record: rec}
				nel := (int64(cbmem.Memory.Size) - (j - start)) / int64(reflect.TypeOf(memoryRange{}).Size())
				cbmem.Memory.Maps = make([]memoryRange, nel)
				if err := readOne(r, cbmem.Memory.Maps, j); err != nil {
					return nil, found, err
				}
			case LB_TAG_TIMESTAMPS:
				if err := readOne(r, &cbmem.TimeStampsTable, start); err != nil {
					return nil, found, err
				}
				if cbmem.TimeStampsTable.Addr == 0 {
					continue
				}
				if cbmem.TimeStamps, err = cbmem.readTimeStamps(f); err != nil {
					log.Printf("TimeStampAddress is %#x but ReadTimeStamps failed: %v", cbmem.TimeStampsTable, err)
					return nil, found, err
				}
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
					return nil, false, fmt.Errorf("trying to read string for %s: %w", n, err)
				}
				p, err := bufio.NewReader(io.NewSectionReader(r, j+2+int64(len(v)), 65536)).ReadString(0)
				if err != nil {
					return nil, false, fmt.Errorf("trying to read string for %s: %w", n, err)
				}
				cbmem.MainBoard.Vendor = v[:len(v)-1]
				cbmem.MainBoard.PartNumber = p[:len(p)-1]
			case LB_TAG_HWRPB:
				if err := readOne(r, &cbmem.Hwrpb, start); err != nil {
					return nil, found, err
				}

				// "Nobody knew consoles could be so hard."
			case LB_TAG_CBMEM_CONSOLE:
				c := &memconsoleEntry{Record: rec}
				debug("    Found cbmem console(%#x), %d byte record.\n", rec, c.Size)
				if err := readOne(r, &c.Address, j); err != nil {
					return nil, found, err
				}
				debug("    console data is at %#x", c.Address)
				cbcons := int64(c.Address)
				// u32 size;
				// u32 cursor;
				// u8  body[0];
				// The cbmem size is a guess.
				cr, err := newOffsetReader(f, cbcons, 8)
				if err != nil {
					return nil, found, err
				}
				if err := readOne(cr, &c.Size, cbcons); err != nil {
					return nil, found, err
				}

				cbcons += int64(reflect.TypeOf(c.Size).Size())
				if err := readOne(cr, &c.Cursor, cbcons); err != nil {
					return nil, found, err
				}
				cbcons += int64(reflect.TypeOf(c.Cursor).Size())
				debug("CSize is %#x, and Cursor is at %#x", c.CSize, c.Cursor)
				// p.cur f8b4 p.si 1fff8 curs f8b4 size f8b4
				sz := int(c.Size)

				cr, err = newOffsetReader(f, cbcons, sz)
				if err != nil {
					return nil, found, err
				}

				curse := int(c.Cursor & CBMC_CURSOR_MASK)
				data := make([]byte, sz)
				// This one is easy. Read from 0 to the cursor.
				if c.Cursor&CBMC_OVERFLOW == 0 {
					if curse < int(c.Size) {
						sz = curse
						data = data[:sz]
					}

					debug("CSize is %d, and Cursor is at %d", c.CSize, c.Cursor)

					if n, err := cr.ReadAt(data, cbcons); err != nil || n != len(data) {
						return nil, found, err
					}
				} else {
					debug("CSize is %#x, and Cursor is at %#x", curse, sz)
					// This should not happen, but that means that it WILL happen
					// some day ...
					if curse > sz {
						curse = 0
					}
					off := cbcons + int64(curse)
					if n, err := cr.ReadAt(data[:curse], off); err != nil || n != len(data[:curse]) {
						return nil, found, err
					}
					if n, err := cr.ReadAt(data[curse:], cbcons); err != nil || n != len(data[curse:]) {
						debug("2nd read: %v", err)
						return nil, found, err
					}
				}

				c.Data = string(data)
				cbmem.MemConsole = c

			case LB_TAG_FORWARD:
				var newTable int64
				if err := readOne(r, &newTable, j); err != nil {
					return nil, found, err
				}
				debug("Forward to %08x", newTable)
				return parseCBtable(f, newTable, 1048576)
			default:
				if n, ok := tagNames[rec.Tag]; ok {
					debug("Ignoring record %v", n)
					cbmem.Ignored = append(cbmem.Ignored, n)
					j = start + int64(rec.Size)
					continue
				}
				log.Printf("Unknown tag record %v %#x", rec, rec.Tag)
				cbmem.Unknown = append(cbmem.Unknown, rec.Tag)

			}
			j = start + int64(rec.Size)
		}
	}
	return cbmem, found, nil
}

// DumpMem prints the memory areas. If hexdump is set, it will hexdump
// LB tables.
func DumpMem(f *os.File, cbmem *CBmem, hexdump bool, w io.Writer) error {
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
			r, err := newOffsetReader(f, int64(e.Start), int(e.Size))
			if err != nil {
				log.Print(err)
				continue
			}
			// The hexdump does a lot of what we want, but not all of
			// what we want. In particular, we'd like better control of
			// what is printed with the offset. So ... hackery.
			out := ""
			same := 0
			var line [16]byte
			for i := e.Start; i < e.Start+e.Size; i += 16 {
				n, err := r.ReadAt(line[:], int64(i))
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				s := hex.Dump(line[:n])[10:]
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
	return nil
}

func cbMem(w io.Writer) error {
	var err error
	if version {
		fmt.Fprintln(w, "cbmem in Go, including JSON output")
		return err
	}
	if verbose {
		debug = log.Printf
	}

	f, err := os.Open(*mem)
	if err != nil {
		return err
	}

	var cbmem *CBmem
	var found bool
	for _, addr := range []int64{0, 0xf0000} {
		cbmem, found, err = parseCBtable(f, addr, 0x10000)
		if err != nil {
			return err
		}
		if found {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("reading coreboot table: %w", err)
	}
	if !found {
		return fmt.Errorf("no coreboot table found")
	}

	if timestamps {
		ts := cbmem.TimeStamps

		// Format in tab-separated columns with a tab stop of 8.
		tw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		freq := uint64(ts.TickFreqMHZ)
		debug("ts %#x freq %#x stamps %#x\n", ts, freq, ts.TS)
		prev := ts.TS[0].EntryStamp
		p := message.NewPrinter(language.English)
		p.Fprintf(tw, "%d entries total:\n\n", len(ts.TS))
		for _, t := range ts.TS {
			n, ok := TimeStampNames[t.EntryID]
			if !ok {
				n = fmt.Sprintf("[%#x]", t.EntryID)
			}
			cur := t.EntryStamp
			debug("cur %#x cur / freq %#x", cur, cur/freq)
			p.Fprintf(tw, "\t%d:%s\t%d (%d)\n", t.EntryID, n, cur/freq, (cur-prev)/freq)
			prev = cur

		}
		tw.Flush()
	}
	if dumpJSON {
		b, err := json.MarshalIndent(cbmem, "", "\t")
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}
		fmt.Fprintf(w, "%s\n", b)
	}
	// list is kind of misnamed I think. It really just prints
	// memory table entries.
	if list || hexdump {
		DumpMem(f, cbmem, hexdump, os.Stdout)
	}
	if console && cbmem.MemConsole != nil {
		fmt.Fprintf(w, "%s%s", cbmem.MemConsole.Data[cbmem.MemConsole.Cursor:], cbmem.MemConsole.Data[0:cbmem.MemConsole.Cursor])
	}
	return err
}

//go:generate go run gen/gen.go -apu2

func main() {
	flag.Parse()
	if err := cbMem(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
