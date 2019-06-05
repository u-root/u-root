// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sort lines.
//
// Synopsis:
//     sort [OPTIONS]... [INPUT]...
//
// Description:
//     Sort copies lines from the input to the output, sorting them in the
//     process. This does nothing fancy (no multi-threading, compression,
//     optiminzations, ...); it simply uses Go's sort.Sort function.
//
// Options:
//     -r:      reverse
//     -o FILE: output file
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

func readInput() string {
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

	// Read unicode string from input
	fileContents := []string{}
	for _, f := range from {
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		s := string(bytes)
		fileContents = append(fileContents, s)
		// Add a separator between files if the file is not newline
		// terminated. Prevents concatenating lines between files.
		if len(s) > 0 && s[len(s)-1] != '\n' {
			fileContents = append(fileContents, "\n")
		}
	}
	return strings.Join(fileContents, "")
}

func sortAlgorithm(s string) string {
	if len(s) == 0 {
		return "" // edge case mimics coreutils
	}
	if s[len(s)-1] == '\n' {
		s = s[:len(s)-1] // remove newline terminator
	}
	lines := strings.Split(string(s), "\n")
	if *reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	} else {
		sort.Strings(lines)
	}
	return strings.Join(lines, "\n") + "\n" // append newline terminator
}

func writeOutput(s string) {
	to := os.Stdout
	if *outputFile != "" {
		if f, err := os.Create(*outputFile); err == nil {
			to = f
			defer f.Close()
		} else {
			log.Fatal(err)
		}
	}
	to.Write([]byte(s))
}

func main() {
	flag.Parse()

	// Input files must be closed before writing to output files to solve
	// the situtation in which the output file is the same as an input.
	writeOutput(sortAlgorithm(readInput()))
}
