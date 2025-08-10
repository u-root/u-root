// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// flash reads and writes to a flash chip.
//
// Synopsis:
//
//	flash -p PROGRAMMER[:parameter[,parameter[...]]] [-r FILE|-w FILE]
//
// Options:
//
//	-o offset: Offset at which to start.
//	-s size: Number of bytes to read or write.
//	-p PROGRAMMER: Specify the programmer with zero or more parameters (see
//	               below).
//	-r FILE: Read flash data into the file.
//	-w FILE: Write the file to the flash chip. First, the flash chip is read
//	         and then diffed against the file. The differing blocks are
//	         erased and written. Finally, the contents are verified.
//
// Programmers:
//
//	dummy
//	  Virtual flash programmer for testing in a memory buffer.
//
//	  dummy:image=image.rom
//	    File to memmap for the memory buffer.
//
//	linux_spi:dev=/dev/spidev0.0
//	  Use Linux's spidev driver. This is only supported on Linux. The dev
//	  parameter is required.
//
//	  linux_spi:dev=/dev/spidev0.0,spispeed=5000
//	    Set the SPI controller's speed. The frequency is in kilohertz.
//
// Description:
//
//	flash is u-root's implementation of flashrom. It has a very limited
//	feature set and depends on the the flash chip implementing the SFDP.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

type programmer interface {
	io.ReaderAt
	io.WriterAt
	EraseAt(int64, int64) (int64, error)
	Size() int64
	Close() error
}

type (
	programmerParams map[string]string
	programmerInit   func(programmerParams) (programmer, error)
)

// supportedProgrammers is populated by the other files in this package.
var supportedProgrammers = map[string]programmerInit{}

func parseProgrammerParams(arg string) (string, map[string]string) {
	params := map[string]string{}

	colon := strings.IndexByte(arg, ':')
	if colon == -1 {
		return arg, params
	}
	for _, p := range strings.Split(arg[colon+1:], ",") {
		equal := strings.IndexByte(p, '=')
		if equal == -1 {
			params[p] = ""
			continue
		}
		params[p[:equal]] = p[equal+1:]
	}
	return arg[:colon], params
}

func run(args []string, supportedProgrammers map[string]programmerInit) (reterr error) {
	// Make a human readable list of supported programmers.
	programmerList := []string{}
	for k := range supportedProgrammers {
		programmerList = append(programmerList, k)
	}
	sort.Strings(programmerList)

	var (
		e    bool
		p    string
		r    string
		w    string
		off  int64
		size int64 = math.MaxInt64
	)

	fs := flag.NewFlagSet("flash", flag.ContinueOnError)
	fs.BoolVar(&e, "erase", false, "erase the flash part")
	fs.BoolVar(&e, "e", false, "erase the flash part")
	fs.StringVar(&p, "programmer", "", fmt.Sprintf("programmer (%s)", strings.Join(programmerList, ",")))
	fs.StringVar(&p, "p", "", fmt.Sprintf("programmer (%s)", strings.Join(programmerList, ",")))
	fs.StringVar(&r, "read", "", "read flash data into the file")
	fs.StringVar(&r, "r", "", "read flash data into the file")
	fs.StringVar(&w, "write", "", "write the file to flash")
	fs.StringVar(&w, "w", "", "write the file to flash")
	fs.Int64Var(&off, "offset", 0, "off at which to write")
	fs.Int64Var(&off, "o", 0, "off at which to write")
	fs.Int64Var(&size, "size", math.MaxInt64, "number of bytes")
	fs.Int64Var(&size, "s", math.MaxInt64, "number of bytes")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		flag.Usage()
		return errors.New("unexpected positional arguments")
	}

	if p == "" {
		return errors.New("-p needs to be set")
	}

	if r == "" && w == "" && !e {
		return errors.New("at least one of -e, -r or -w need to be set")
	}
	if r != "" && w != "" {
		return errors.New("both -r and -w cannot be set")
	}

	programmerName, params := parseProgrammerParams(p)
	init, ok := supportedProgrammers[programmerName]
	if !ok {
		return fmt.Errorf("unrecognized programmer %q", programmerName)
	}

	programmer, err := init(params)
	if err != nil {
		return err
	}
	defer func() {
		err := programmer.Close()
		if reterr == nil {
			reterr = err
		}
	}()

	if e {
		n, err := programmer.EraseAt(min(size, programmer.Size()), off)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Erased %#x bytes @ %#x", n, off)
	}

	// Create a buffer to hold the contents of the image.

	if r != "" {
		buf := make([]byte, min(size, programmer.Size()))
		f, err := os.Create(r)
		if err != nil {
			return err
		}
		defer func() {
			err := f.Close()
			if reterr == nil {
				reterr = err
			}
		}()
		if _, err := programmer.ReadAt(buf, off); err != nil {
			return err
		}
		if _, err := f.Write(buf); err != nil {
			return err
		}
	} else if w != "" {
		buf, err := os.ReadFile(w)
		if err != nil {
			return err
		}
		buf = buf[:min(int64(len(buf)), size)]
		amt, err := programmer.WriteAt(buf, off)
		if err != nil {
			return fmt.Errorf("writing %d bytes to dev %v:%w", len(buf), programmer, err)
		}
		if amt != len(buf) {
			return fmt.Errorf("only flashed %d of %d bytes", amt, len(buf))
		}

		return nil
	}

	return nil
}

func main() {
	if err := run(os.Args[1:], supportedProgrammers); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
