// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Sort copies lines from the input to the output, sorting them in the process.
This does nothing fancy (no multi-threading, compression, tables, ...); it
simply uses Go's sort.Sort function.

sort [OPTION]... [FILE]...

The options are:
	-r		reverse
	-o string	output file
*/

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	reverse    = flag.Bool("r", false, "Reverse")
	outputFile = flag.String("o", "", "Output file")
)

// Sort from file a to file b.
func sortFiles(from []*os.File, to *os.File) {
	// Read unicode string from input
	fileContents := []string{}
	for _, f := range from {
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		fileContents = append(fileContents, string(bytes))
	}
	s := strings.Join(fileContents, "")

	// Sorting algorithm
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1] // remove terminating newline
	}
	lines := strings.Split(string(s), "\n")
	if *reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	} else {
		sort.Strings(lines)
	}
	s = strings.Join(lines, "\n") + "\n" // append newline terminated

	// Write to output
	to.Write([]byte(s))
}

func main() {
	flag.Parse()

	// Input files
	from := []*os.File{}
	if flag.NArg() > 0 {
		for _, v := range flag.Args() {
			if f, err := os.Open(v); err == nil {
				from = append(from, f)
				defer f.Close()
			} else {
				log.Fatal(err)
			}
		}
	} else {
		from = []*os.File{os.Stdin}
	}

	// Output file
	var to *os.File = os.Stdout
	if *outputFile != "" {
		log.Println(*outputFile)
		if f, err := os.Create(*outputFile); err == nil {
			to = f
			defer f.Close()
		} else {
			log.Fatal(err)
		}
	}

	sortFiles(from, to)
}
