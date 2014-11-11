// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type grepResult struct {
	match bool
	line  *string
}

var (
	Match = flag.Bool("v", true, "Print only non-matching lines")
)

func grep(f *os.File, re *regexp.Regexp, res chan *grepResult) {
	r := bufio.NewReader(f)
	for {
		if i, err := r.ReadString('\n'); err == nil {
			res <- &grepResult{re.Match([]byte(i)), &i}
		} else {
			break
		}
	}
	close(res)
}

func main() {
	r := ".*"
	flag.Parse()
	a := flag.Args()
	if len(a) > 0 {
		r = a[0]
	}
	re := regexp.MustCompile(r)
	if len(a) < 2 {
		res := make(chan *grepResult)
		go grep(os.Stdin, re, res)
		for i := range res {
			if i.match == *Match {
				fmt.Printf("%v", *i.line)
			}
		}
	}
}
