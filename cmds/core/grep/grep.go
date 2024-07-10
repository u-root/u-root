// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// grep searches file contents using regular expressions.
//
// Synopsis:
//
//	grep [-clFivnhqre] [FILE]...
//
// Options:
//
//  -c, --count                Just show counts
//  -l, --files-with-matches   list only files
//  -F, --fixed-strings        Match using fixed strings
//  -i, --ignore-case          case-insensitive matching
//  -v, --invert-match         Print only non-matching lines
//  -n, --line-number          Show line numbers
//  -h, --no-filename          Suppress file name prefixes on output
//  -q, --quiet                Don't print matches; exit on first match
//  -r, --recursive            recursive
//  -e, --regexp string        Pattern to match

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errQuiet = fmt.Errorf("not found")

type params struct {
	expr string
	headers, invert, recursive, caseInsensitive, fixed,
	noShowMatch, quiet, count, number bool
}

type grepCommand struct {
	rc   io.ReadCloser
	name string
}

// run runs a command. args are as from os.Args, i.e., args[0] is the command name.
func run(stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, args []string) error {
	var c cmd

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&c.params.expr, "regexp", "", "Pattern to match")
	f.StringVar(&c.expr, "e", "", "Pattern to match (shorthand)")

	f.BoolVar(&c.params.headers, "no-filename", false, "Suppress file name prefixes on output")
	f.BoolVar(&c.params.headers, "h", false, "Suppress file name prefixes on output (shorthand)")

	f.BoolVar(&c.params.invert, "invert-match", false, "Print only non-matching lines")
	f.BoolVar(&c.params.invert, "v", false, "Print only non-matching lines (shorthand)")

	f.BoolVar(&c.params.recursive, "recursive", false, "recursive")
	f.BoolVar(&c.params.recursive, "r", false, "recursive")

	f.BoolVar(&c.params.noShowMatch, "files-with-matches", false, "list only files")
	f.BoolVar(&c.params.noShowMatch, "l", false, "list only files (shorthand)")

	f.BoolVar(&c.params.count, "count", false, "Just show counts")
	f.BoolVar(&c.params.count, "c", false, "Just show counts")

	f.BoolVar(&c.params.caseInsensitive, "ignore-case", false, "case-insensitive matching")
	f.BoolVar(&c.params.caseInsensitive, "i", false, "case-insensitive matching (shorthand)")

	f.BoolVar(&c.params.number, "line-number", false, "Show line numbers")
	f.BoolVar(&c.params.number, "n", false, "Show line numbers (shorthand)")

	f.BoolVar(&c.params.fixed, "fixed-strings", false, "Match using fixed strings")
	f.BoolVar(&c.params.fixed, "F", false, "Match using fixed strings (shorthand)")

	f.BoolVar(&c.params.quiet, "quiet", false, "Don't print matches; exit on first match")
	f.BoolVar(&c.params.quiet, "q", false, "Don't print matches; exit on first match (shorthand)")

	f.BoolVar(&c.params.quiet, "silent", false, "Don't print matches; exit on first match")
	f.BoolVar(&c.params.quiet, "s", false, "Don't print matches; exit on first match (shorthand)")

	f.Usage = func() {
		fmt.Fprint(f.Output(), "Usage: grep [-clFivnhqre] [FILE]...\n\n")
		f.PrintDefaults()
	}

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))

	c.args = f.Args()
	c.stdin = stdin
	c.stdout = bufio.NewWriter(stdout)
	c.stderr = stderr

	return c.run()
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr, os.Args); err != nil {
		if err == errQuiet {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}

// cmd contains the actually business logic of grep
type cmd struct {
	stdin  io.ReadCloser
	stdout *bufio.Writer
	stderr io.Writer
	args   []string
	params
	matchCount int
	showName   bool
}

// grep reads data from the os.File embedded in grepCommand.
// It matches each line against the re and prints the matching result
// If we are only looking for a match, we exit as soon as the condition is met.
// "match" means result of re.Match == match flag.
func (c *cmd) grep(f *grepCommand, re *regexp.Regexp) (ok bool) {
	r := bufio.NewScanner(f.rc)
	defer f.rc.Close()
	var lineNum int
	for r.Scan() {
		line := r.Text()
		var m bool
		switch {
		case c.fixed && c.caseInsensitive:
			m = strings.Contains(strings.ToLower(line), strings.ToLower(c.expr))
		case c.fixed && !c.caseInsensitive:
			m = strings.Contains(line, c.expr)
		default:
			m = re.MatchString(line)
		}
		if m != c.invert {
			// in quiet mode, exit before the first match
			if c.quiet {
				return false
			}
			c.printMatch(f, line, lineNum+1, m)
			if c.noShowMatch {
				break
			}
		}
		lineNum++
	}
	c.stdout.Flush()
	return true
}

func (c *cmd) printMatch(cmd *grepCommand, line string, lineNum int, match bool) {
	if match == !c.invert {
		c.matchCount++
	}
	if c.count {
		return
	}
	// at this point, we have committed to writing a line
	defer func() {
		c.stdout.WriteByte('\n')
	}()
	// if showName, write name to stdout
	if c.showName {
		c.stdout.WriteString(cmd.name)
	}
	// if dont show match, then newline and return, we are done
	if c.noShowMatch {
		return
	}
	if match == !c.invert {
		// if showName, need a :
		if c.showName {
			c.stdout.WriteByte(':')
		}
		// if showing line number, print the line number then a :
		if c.number {
			c.stdout.Write(strconv.AppendUint(nil, uint64(lineNum), 10))
			c.stdout.WriteByte(':')
		}
		// now write the line to stdout
		c.stdout.WriteString(line)
	}
}

func (c *cmd) run() error {
	defer c.stdout.Flush()
	// parse the expression into valid regex
	if c.expr != "" {
		c.args = append([]string{c.expr}, c.args...)
	}
	r := ".*"
	if len(c.args) > 0 {
		r = c.args[0]
	}
	if c.caseInsensitive && !bytes.HasPrefix([]byte(r), []byte("(?i)")) && !c.fixed {
		r = "(?i)" + r
	}
	var re *regexp.Regexp
	if !c.fixed {
		re = regexp.MustCompile(r)
	} else if c.expr == "" {
		c.expr = c.args[0]
	}

	// if len(c.args) < 2, then we read from stdin
	if len(c.args) < 2 {
		if !c.grep(&grepCommand{c.stdin, "<stdin>"}, re) {
			return nil
		}
	} else {
		c.showName = (len(c.args[1:]) > 1 || c.recursive || c.noShowMatch) && !c.headers
		var ok bool
		for _, v := range c.args[1:] {
			err := filepath.Walk(v, func(name string, fi os.FileInfo, err error) error {
				if err != nil {
					fmt.Fprintf(c.stderr, "grep: %v: %v\n", name, err)
					return nil
				}
				if fi.IsDir() && !c.recursive {
					fmt.Fprintf(c.stderr, "grep: %v: Is a directory\n", name)
					return filepath.SkipDir
				}
				fp, err := os.Open(name)
				if err != nil {
					fmt.Fprintf(c.stderr, "can't open %s: %v\n", name, err)
					return nil
				}
				defer fp.Close()
				if !c.grep(&grepCommand{fp, name}, re) {
					ok = true
					return nil
				}
				return nil
			})
			if ok {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
	if c.quiet {
		return errQuiet
	}
	if c.count {
		c.stdout.Write(strconv.AppendUint(nil, uint64(c.matchCount), 10))
		c.stdout.WriteByte('\n')
	}
	return nil
}
