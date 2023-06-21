// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// grep searches file contents using regular expressions.
//
// Synopsis:
//
//	grep [-vrlq] [FILE]...
//
// Options:
//
//	-v: print only non-matching lines
//	-r: recursive
//	-l: list only files
//	-q: don't print matches; exit on first match
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

var errQuite = fmt.Errorf("not found")

type grepResult struct {
	c       *grepCommand
	line    *string
	lineNum int
	match   bool
}

type grepCommand struct {
	rc   io.ReadCloser
	name string
}

type oneGrep struct {
	c chan *grepResult
}

var (
	expr            = flag.StringP("regexp", "e", "", "Pattern to match")
	headers         = flag.BoolP("no-filename", "h", false, "Suppress file name prefixes on output")
	invert          = flag.BoolP("invert-match", "v", false, "Print only non-matching lines")
	recursive       = flag.BoolP("recursive", "r", false, "recursive")
	noShowMatch     = flag.BoolP("files-with-matches", "l", false, "list only files")
	count           = flag.BoolP("count", "c", false, "Just show counts")
	caseInsensitive = flag.BoolP("ignore-case", "i", false, "case-insensitive matching")
	number          = flag.BoolP("line-number", "n", false, "Show line numbers")
	fixed           = flag.BoolP("fixed-strings", "F", false, "Match using fixed strings")
)

// grep reads data from the os.File embedded in grepCommand.
// It creates a chan of grepResults and pushes a pointer to it into allGrep.
// It matches each line against the re and pushes the matching result
// into the chan.
// If we are only looking for a match, we exit as soon as the condition is met.
// "match" means result of re.Match == match flag.
func (c *cmd) grep(f *grepCommand, re *regexp.Regexp, res chan *grepResult) {
	r := bufio.NewReader(f.rc)
	var lineNum int
	for {
		if i, err := r.ReadString('\n'); err == nil {
			var m bool
			if c.fixed {
				if c.caseInsensitive {
					m = strings.Contains(strings.ToLower(i), strings.ToLower(c.expr))
				} else {
					m = strings.Contains(i, c.expr)
				}
			} else {
				m = re.Match([]byte(i))
			}
			if m == !c.invert {
				res <- &grepResult{
					match:   m,
					c:       f,
					line:    &i,
					lineNum: lineNum + 1,
				}
				if c.noShowMatch {
					break
				}
			}
		} else {
			break
		}
		lineNum++
	}
	close(res)
	_ = f.rc.Close()
}

func (c *cmd) printMatch(r *grepResult) {
	var prefix string
	if r.match == !c.invert {
		c.matchCount++
	}
	if c.count {
		return
	}
	if c.showName {
		fmt.Fprintf(c.stdout, "%v", r.c.name)
		prefix = ":"
	}
	if c.noShowMatch {
		fmt.Fprintf(c.stdout, "\n")
		return
	}
	if c.number {
		prefix = fmt.Sprintf("%d:", r.lineNum)
	}
	if r.match == !c.invert {
		fmt.Fprintf(c.stdout, "%v%v", prefix, *r.line)
	}
}

type params struct {
	expr            string
	headers         bool
	invert          bool
	recursive       bool
	noShowMatch     bool
	count           bool
	caseInsensitive bool
	number          bool
	quiet           bool
	fixed           bool
}

type cmd struct {
	stdin   io.ReadCloser
	stdout  io.Writer
	stderr  io.Writer
	allGrep chan *oneGrep
	args    []string
	params
	matchCount int
	showName   bool
}

func command(stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		params:  p,
		args:    args,
		allGrep: make(chan *oneGrep),
	}
}

func main() {
	flag.Parse()
	p := params{
		expr:            *expr,
		headers:         *headers,
		invert:          *invert,
		recursive:       *recursive,
		noShowMatch:     *noShowMatch,
		count:           *count,
		caseInsensitive: *caseInsensitive,
		number:          *number,
		quiet:           *quiet,
		fixed:           *fixed,
	}

	if err := command(os.Stdin, os.Stdout, os.Stderr, p, flag.Args()).run(); err != nil {
		if err == errQuite {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
func (c *cmd) run() error {
	// parse the expression into valid regex
	if c.expr != "" {
		c.args = append([]string{c.expr}, c.args...)
	}
	r := ".*"
	if len(c.args) > 0 {
		r = c.args[0]
	}
	if c.caseInsensitive && !strings.HasPrefix(r, "(?i)") && !c.fixed {
		r = "(?i)" + r
	}
	var re *regexp.Regexp
	if !c.fixed {
		re = regexp.MustCompile(r)
	} else if c.expr == "" {
		c.expr = c.args[0]
	}

	// start producing greps in a goroutine
	go func() {
		// very special case, just stdin
		if len(c.args) < 2 {
			c.run_stdin(re)
		} else {
			c.run_files(re)
		}
		// the allGrep channel is closed when work is done
		close(c.allGrep)
	}()

	return c.run_print()
}

func (c *cmd) run_stdin(re *regexp.Regexp) {
	res := make(chan *grepResult, 1)
	go c.grep(&grepCommand{c.stdin, "<stdin>"}, re, res)
	c.allGrep <- &oneGrep{res}

}

func (c *cmd) run_files(re *regexp.Regexp) {
	c.showName = (len(c.args[1:]) > 1 || c.recursive || c.noShowMatch) && !c.headers
	// generate a chan of file names, bounded by the size of the chan. This in turn
	// throttles the opens.
	treeNames := make(chan string, 128)
	go func() {
		defer close(treeNames)
		for _, v := range c.args[1:] {
			filepath.Walk(v, func(name string, fi os.FileInfo, err error) error {
				if err != nil {
					fmt.Fprintf(c.stderr, "grep: %v: %v\n", name, err)
					return nil
				}
				if fi.IsDir() && !c.recursive {
					fmt.Fprintf(c.stderr, "grep: %v: Is a directory\n", name)
					return filepath.SkipDir
				}
				treeNames <- name
				return nil
			})
		}
	}()

	files := make(chan *grepCommand)
	// convert the file names to a stream of os.File
	go func() {
		for i := range treeNames {
			fp, err := os.Open(i)
			if err != nil {
				fmt.Fprintf(c.stderr, "can't open %s: %v\n", i, err)
				continue
			}
			files <- &grepCommand{fp, i}
		}
		close(files)
	}()
	// now kick off the greps
	for f := range files {
		res := make(chan *grepResult, 1)
		c.allGrep <- &oneGrep{res}
		go c.grep(f, re, res)
	}
}

func (c *cmd) run_print() error {
	for {
		og, ok := <-c.allGrep
		if !ok {
			break
		}
		for r := range og.c {
			// exit on first match if running in quiet mode.
			if c.quiet {
				return nil
			}
			c.printMatch(r)
		}
	}

	if c.quiet {
		return errQuite
	}
	if c.count {
		fmt.Fprintf(c.stdout, "%d\n", c.matchCount)
	}
	return nil
}
