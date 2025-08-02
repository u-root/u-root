// Copyright 2017-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sort lines.
//
// Synopsis:
//
//	sort [OPTIONS]... [INPUT]...
//
// Description:
//
//	Sort copies lines from the input to the output, sorting them in the
//	process. This does nothing fancy (no multi-threading, compression,
//	optiminzations, ...); it simply uses Go's sort.Sort function.
//
// Options:
//
//	-r:      Reverse the result of comparisons
//	-C:      Check that the single input file is ordered. No warnings.
//	-u:	     Unique keys. Suppress all lines that have a key that is equal to an already processed one.
//	-f: 	 Fold lower case to upper case character.
//	-b: 	 Ignore leading blank characters when comparing lines.
//	-n:      Compare according to string numerical value.
//	-o FILE: Specify the name of an output file to be used instead of the standard output.
package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	reverse      = flag.Bool("r", false, "Reverse the result of comparisons.")
	ordered      = flag.Bool("C", false, "Check that the single input file is ordered. No warnings.")
	unique       = flag.Bool("u", false, "Unique keys. Suppress all lines that have a key that is equal to an already processed one.")
	ignoreCase   = flag.Bool("f", false, "Fold lower case to upper case character.")
	ignoreBlanks = flag.Bool("b", false, "Ignore leading blank characters when comparing lines.")
	numeric      = flag.Bool("n", false, "Compare according to string numerical value.")
	outputFile   = flag.String("o", "", "Specify the name of an output file to be used instead of the standard output.")
)

type ignoreCaseSort []string

func (a ignoreCaseSort) Len() int           { return len(a) }
func (a ignoreCaseSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ignoreCaseSort) Less(i, j int) bool { return strings.ToUpper(a[i]) < strings.ToUpper(a[j]) }

type ignoreBlanksSort []string

func (a ignoreBlanksSort) Len() int      { return len(a) }
func (a ignoreBlanksSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ignoreBlanksSort) Less(i, j int) bool {
	l := strings.TrimLeftFunc(a[i], unicode.IsSpace)
	r := strings.TrimLeftFunc(a[j], unicode.IsSpace)
	if l == r {
		return len(a[i]) >= len(a[j])
	}
	return l < r
}

type ignoreBlanksCaseSort []string

func (a ignoreBlanksCaseSort) Len() int      { return len(a) }
func (a ignoreBlanksCaseSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ignoreBlanksCaseSort) Less(i, j int) bool {
	l := strings.ToUpper(strings.TrimLeftFunc(a[i], unicode.IsSpace))
	r := strings.ToUpper(strings.TrimLeftFunc(a[j], unicode.IsSpace))
	if l == r {
		return len(a[i]) >= len(a[j])
	}
	return l < r
}

type numericSort []string

func (a numericSort) Len() int      { return len(a) }
func (a numericSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a numericSort) Less(i, j int) bool {
	// consider removing thousands separator and parsing LC_NUMERIC
	l := strings.ToUpper(strings.TrimLeftFunc(a[i], unicode.IsSpace))
	r := strings.ToUpper(strings.TrimLeftFunc(a[j], unicode.IsSpace))

	// treat all non-numeric characters as zeros
	ln, _ := strconv.ParseFloat(l, 64)
	rn, _ := strconv.ParseFloat(r, 64)

	return ln < rn
}

var errNotOrdered = errors.New("not ordered")

type params struct {
	outputFile   string
	reverse      bool
	ordered      bool
	unique       bool
	ignoreCase   bool
	ignoreBlanks bool
	numeric      bool
}

type cmd struct {
	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer
	params params
	args   []string
}

func command(stdin io.ReadCloser, stdout, stderr io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		params: p,
		args:   args,
	}
}

func (c *cmd) run() error {
	// Input files
	from := []io.ReadCloser{}
	for _, v := range c.args {
		f, err := os.Open(v)
		if err != nil {
			return err
		}
		defer f.Close()
		from = append(from, f)
	}

	if len(c.args) == 0 {
		from = append(from, c.stdin)
	}

	// Read unicode string from input
	fileContents := []string{}
	for _, f := range from {
		bytes, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		s := string(bytes)
		fileContents = append(fileContents, s)
		// Add a separator between files if the file is not newline
		// terminated. Prevents concatenating lines between files.
		if len(s) > 0 && s[len(s)-1] != '\n' {
			fileContents = append(fileContents, "\n")
		}
	}

	s := strings.Join(fileContents, "")
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1] // remove newline terminator
	}

	if c.params.ordered {
		// if ordered is true, set ignoreBlanks to false to be consistent with coreutils
		// see https://github.com/coreutils/coreutils/blob/d53190ed46a55f599800ebb2d8ddfe38205dbd24/src/sort.c#L4147
		c.params.ignoreBlanks = false
		lines := strings.Split(s, "\n")
		if !sort.IsSorted(c.sortInterface(lines)) {
			return errNotOrdered
		}

		if c.params.unique && len(lines) > 1 {
			for i := 1; i < len(lines); i++ {
				if c.params.ignoreCase && strings.EqualFold(lines[i], lines[i-1]) {
					return errNotOrdered
				} else if lines[i] == lines[i-1] {
					return errNotOrdered
				}
			}
		}

		return nil
	}

	if err := c.writeOutput(c.stdout, c.sortAlgorithm(s)); err != nil {
		return err
	}
	return nil
}

func (c *cmd) sortInterface(lines []string) sort.Interface {
	var si sort.Interface
	switch {
	case c.params.ignoreBlanks && c.params.ignoreCase:
		si = ignoreBlanksCaseSort(lines)
	case c.params.ignoreBlanks:
		si = ignoreBlanksSort(lines)
	case c.params.ignoreCase:
		si = ignoreCaseSort(lines)
	case c.params.numeric:
		si = numericSort(lines)
	default:
		si = sort.StringSlice(lines)
	}

	return si
}

func (c *cmd) sortAlgorithm(s string) string {
	if len(s) == 0 {
		return "" // edge case mimics coreutils
	}
	lines := strings.Split(s, "\n")
	si := c.sortInterface(lines)

	if c.params.reverse {
		sort.Sort(sort.Reverse(si))
	} else {
		sort.Sort(si)
	}

	if c.params.unique && len(lines) > 1 {
		j := 1
		for i := 1; i < len(lines); i++ {
			if c.params.ignoreCase && c.params.ignoreBlanks {
				l1 := strings.TrimLeftFunc(lines[i], unicode.IsSpace)
				l2 := strings.TrimLeftFunc(lines[i-1], unicode.IsSpace)
				if strings.EqualFold(l1, l2) {
					continue
				}
			} else if c.params.ignoreCase && strings.EqualFold(lines[i], lines[j]) {
				continue
			} else if lines[i] == lines[i-1] {
				continue
			}
			lines[j] = lines[i]
			j++
		}

		lines = lines[:j]
	}

	return strings.Join(lines, "\n") + "\n" // append newline terminator
}

func (c *cmd) writeOutput(w io.Writer, s string) error {
	to := w
	if c.params.outputFile != "" {
		f, err := os.Create(c.params.outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		to = f
	}

	_, err := to.Write([]byte(s))
	return err
}

func main() {
	flag.Parse()
	p := params{
		reverse: *reverse, ordered: *ordered, outputFile: *outputFile, unique: *unique,
		ignoreCase: *ignoreCase, ignoreBlanks: *ignoreBlanks, numeric: *numeric,
	}
	if err := command(os.Stdin, os.Stdout, os.Stderr, p, flag.Args()).run(); err != nil {
		if err == errNotOrdered {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
