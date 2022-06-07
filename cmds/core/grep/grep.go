// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// grep searches file contents using regular expressions.
//
// Synopsis:
//     grep [-vrlq] [FILE]...
//
// Options:
//     -v: print only non-matching lines
//     -r: recursive
//     -l: list only files
//     -q: don't print matches; exit on first match
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

type grepResult struct {
	match   bool
	c       *grepCommand
	line    *string
	lineNum int
}

type grepCommand struct {
	name string
	*os.File
}

type oneGrep struct {
	c chan *grepResult
}

var (
	expr            = flag.StringP("regexp", "e", "", "Pattern to match")
	headers         = flag.BoolP("no-filename", "h", false, "Suppress file name prefixes on output")
	invert          = flag.BoolP("invert-match", "v", false, "Print only non-matching lines")
	recursive       = flag.BoolP("recursive", "r", false, "recursive")
	noshowmatch     = flag.BoolP("files-with-matches", "l", false, "list only files")
	count           = flag.BoolP("count", "c", false, "Just show counts")
	caseinsensitive = flag.BoolP("ignore-case", "i", false, "case-insensitive matching")
	number          = flag.BoolP("line-number", "n", false, "Show line numbers")
	showname        bool
	allGrep         = make(chan *oneGrep)
	nGrep           int
	matchCount      int
)

// grep reads data from the os.File embedded in grepCommand.
// It creates a chan of grepResults and pushes a pointer to it into allGrep.
// It matches each line against the re and pushes the matching result
// into the chan.
// Bug: this chan should be created by the caller and passed in
// to preserve file name order. Oops.
// If we are only looking for a match, we exit as soon as the condition is met.
// "match" means result of re.Match == match flag.
func grep(f *grepCommand, re *regexp.Regexp) {
	r := bufio.NewReader(f)
	res := make(chan *grepResult, 1)
	allGrep <- &oneGrep{res}
	var lineNum int
	for {
		if i, err := r.ReadString('\n'); err == nil {
			m := re.Match([]byte(i))
			if m == !*invert {
				res <- &grepResult{re.Match([]byte(i)), f, &i, lineNum}
				if *noshowmatch {
					break
				}
			}
		} else {
			break
		}
		lineNum++
	}
	close(res)
	f.Close()
}

func printmatch(r *grepResult) {
	var prefix string
	if r.match == !*invert {
		matchCount++
	}
	if *count {
		return
	}
	if showname {
		fmt.Printf("%v", r.c.name)
		prefix = ":"
	}
	if *noshowmatch {
		fmt.Printf("\n")
		return
	}
	if *number {
		prefix = fmt.Sprintf(":%d:", r.lineNum)
	}
	if r.match == !*invert {
		fmt.Printf("%v%v", prefix, *r.line)
	}
}

func main() {
	r := ".*"
	flag.Parse()
	a := flag.Args()
	if *expr != "" {
		// they used the -e flag
		a = append([]string{*expr}, a...)
	}
	if len(a) > 0 {
		r = a[0]
	}
	if *caseinsensitive && !strings.HasPrefix(r, "(?i)") {
		r = "(?i)" + r
	}
	re := regexp.MustCompile(r)
	// very special case, just stdin ...
	if len(a) < 2 {
		nGrep++
		go grep(&grepCommand{"<stdin>", os.Stdin}, re)
	} else {
		showname = (len(a[1:]) > 1 || *recursive) && !*headers
		// generate a chan of file names, bounded by the size of the chan. This in turn
		// throttles the opens.
		treenames := make(chan string, 128)
		go func() {
			defer close(treenames)
			for _, v := range a[1:] {
				// we could parallelize the open part but people might want
				// things to be in order. I don't care but who knows.
				// just ignore the errors. If there is not a single one that works,
				// then all the sizes will be 0 and we'll just fall through.
				filepath.Walk(v, func(name string, fi os.FileInfo, err error) error {
					if err != nil {
						// This is non-fatal because grep searches through
						// all the files it has access to.
						log.Print(err)
						return nil
					}
					if fi.IsDir() && !*recursive {
						fmt.Fprintf(os.Stderr, "grep: %v: Is a directory\n", name)
						return filepath.SkipDir
					}
					if err != nil {
						fmt.Fprintf(os.Stderr, "%v: %v\n", name, err)
						return err
					}
					treenames <- name
					return nil
				})
			}
		}()

		files := make(chan *grepCommand)
		// convert the file names to a stream of os.File
		go func() {
			for i := range treenames {
				fp, err := os.Open(i)
				if err != nil {
					fmt.Fprintf(os.Stderr, "can't open %s: %v\n", i, err)
					continue
				}
				files <- &grepCommand{i, fp}
			}
			close(files)
		}()
		// now kick off the greps
		// bug: file name order is not preserved here. Darn.

		for f := range files {
			nGrep++
			go grep(f, re)
		}
	}

	if nGrep > 0 {
		for c := range allGrep {
			for r := range c.c {
				// exit on first match.
				if *quiet {
					os.Exit(0)
				}
				printmatch(r)
			}
			nGrep--
			if nGrep == 0 {
				break
			}
		}
	}
	if *quiet {
		os.Exit(1)
	}
	if *count {
		fmt.Printf("%d\n", matchCount)
	}
}
