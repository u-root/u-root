// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Concurrent, parallel grep.
// It has to deal with the EMFILE limit.
// To do so we have one chan that is bounded.
// From args, we use filepath.Walk to generate a chan of names.
// From that, we create a chan of grepCommands.
// From that, we create a chan of grepResults.
// The grepResults contain matches or not-matches only.
// if we are in -l mode, the goprocs handling the grep bail out as soon as the condition is met.
// This grep is about 2x faster than GNU grep for simple non-recursive greps and slower
// as soon as filepath.Walk enters the picture. Let's fix this.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type grepResult struct {
	match bool
	c     *grepCommand
	line  *string
}

type grepCommand struct {
	name string
	*os.File
}

type oneGrep struct {
	c chan *grepResult
}

var (
	Match       = flag.Bool("v", true, "Print only non-matching lines")
	recursive   = flag.Bool("r", false, "recursive")
	noshowmatch = flag.Bool("l", false, "list only files")
	showname    = false
	allGrep     = make(chan *oneGrep)
	nGrep       = 0
)

// grep reads data from the os.File embedded in grepCommand.
// It creates a chan of grepResults and pushes a pointer to it into allGrep.
// It matches each line against the re and pushes the matching result
// into the chan.
// Bug: this chan should be created by the caller and passed in
// to preserve file name order. Oops.
// If we are only looking for a match, we exit as soon as the condition is met.
// "match" means result of re.Match == Match flag.
// bug: Match should be match, I expect.
func grep(f *grepCommand, re *regexp.Regexp) {
	nGrep++
	r := bufio.NewReader(f)
	res := make(chan *grepResult, 1)
	allGrep <- &oneGrep{res}
	for {
		if i, err := r.ReadString('\n'); err == nil {
			m := re.Match([]byte(i))
			if m == *Match {
				res <- &grepResult{re.Match([]byte(i)), f, &i}
				if *noshowmatch {
					break
				}
			}
		} else {
			break
		}
	}
	close(res)
	f.Close()
}

func printmatch(r *grepResult) {
	if showname {
		fmt.Printf("%v", r.c.name)
	}
	if *noshowmatch {
		return
	} else if showname {
		fmt.Printf(":")
	}
	if r.match == *Match {
		fmt.Printf("%v", *r.line)
	}
}

func main() {
	r := ".*"
	flag.Parse()
	a := flag.Args()
	if len(a) > 0 {
		r = a[0]
	}
	re := regexp.MustCompile(r)
	// very special case, just stdin ...
	if len(a) < 2 {
		go grep(&grepCommand{"<stdin>", os.Stdin}, re)
	} else {
		showname = len(a[1:]) > 1
		// generate a chan of file names, bounded by the size of the chan. This in turn
		// throttles the opens.
		treenames := make(chan string, 128)
		go func() {
			for _, v := range a[1:] {
				// we could parallelize the open part but people might want
				// things to be in order. I don't care but who knows.
				// just ignore the errors. If there is not a single one that works,
				// then all the sizes will be 0 and we'll just fall through.
				filepath.Walk(v, func(name string, fi os.FileInfo, err error) error {
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
			close(treenames)
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
			go grep(f, re)
		}
	}

	for c := range allGrep {
		for r := range c.c {
			printmatch(r)
		}
		nGrep--
		if nGrep == 0 {
			break
		}
	}
}
