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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
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

func parseParams() params {
	p := params{}
	flag.StringVarP(&p.expr, "regexp", "e", "", "Pattern to match")
	flag.BoolVarP(&p.headers, "no-filename", "h", false, "Suppress file name prefixes on output")
	flag.BoolVarP(&p.invert, "invert-match", "v", false, "Print only non-matching lines")
	flag.BoolVarP(&p.recursive, "recursive", "r", false, "recursive")
	flag.BoolVarP(&p.noShowMatch, "files-with-matches", "l", false, "list only files")
	flag.BoolVarP(&p.count, "count", "c", false, "Just show counts")
	flag.BoolVarP(&p.caseInsensitive, "ignore-case", "i", false, "case-insensitive matching")
	flag.BoolVarP(&p.number, "line-number", "n", false, "Show line numbers")
	flag.BoolVarP(&p.fixed, "fixed-strings", "F", false, "Match using fixed strings")
	flag.BoolVarP(&p.quiet, "quiet", "q", false, "Don't print matches; exit on first match")
	flag.BoolVarP(&p.quiet, "silent", "s", false, "Don't print matches; exit on first match")
	flag.Parse()

	return p
}

func main() {
	if err := command(os.Stdin, os.Stdout, os.Stderr, parseParams(), flag.Args()).run(); err != nil {
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

func command(stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: bufio.NewWriter(stdout),
		stderr: stderr,
		params: p,
		args:   args,
	}
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
