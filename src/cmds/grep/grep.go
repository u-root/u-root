// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	nGrep		= 0
)

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
				if (*noshowmatch) {
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
						fmt.Printf("grep: %v: Is a directory\n", name)
						return filepath.SkipDir
					}
					if err != nil {
						fmt.Printf("%v: %v\n", name, err)
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
