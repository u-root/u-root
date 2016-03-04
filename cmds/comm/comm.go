// Copyright 2013-2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

//Comm reads file1 and file2, which are in lexicographical order, and
//produces a three column output: lines only in file1; lines only in
//file2; and lines in both files. The file name – means the standard
//input.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	s1          = flag.Bool("1", false, "suppress printing of column 1")
	s2          = flag.Bool("2", false, "suppress printing of column 2")
	s3          = flag.Bool("3", false, "suppress printing of column 3")
	help        = flag.Bool("h", false, "print this help message and exit")
	insensitive = flag.Bool("i", false, "case insensitive comparison of lines")
)

func usage() {
	log.Printf("Usage: %s [-123i] file1 file2\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func reader(f *os.File, c chan string) {
	b := bufio.NewReader(f)

	for {
		s, err := b.ReadString('\n')
		c <- strings.TrimRight(s, "\r\n")
		if err != nil {
			break
		}
	}
	close(c)
}

type out struct {
	s1, s2, s3 string
}

func outer(c1, c2 chan string, c chan out) {
	s1, ok1 := <-c1
	s2, ok2 := <-c2
	for {
		if ok1 && ok2 {
			if *insensitive && strings.EqualFold(s1, s2) {
				s2 = s1
			}
			switch {
			case s1 < s2:
				c <- out{s1, "", ""}
				s1, ok1 = <-c1
			case s1 > s2:
				c <- out{"", s2, ""}
				s2, ok2 = <-c2
			default:
				c <- out{"", "", s2}
				s1, ok1 = <-c1
				s2, ok2 = <-c2
			}
		} else if ok1 {
			c <- out{s1, "", ""}
			s1, ok1 = <-c1
		} else if ok2 {
			c <- out{"", s2, ""}
			s2, ok2 = <-c2
		} else {
			break
		}
	}
	close(c)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 || *help {
		flag.Usage()
	}

	c1 := make(chan string, 100)
	c2 := make(chan string, 100)
	c := make(chan out, 100)

	f1, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalf("Can't open %s: %v", flag.Args()[0], err)
	}

	f2, err := os.Open(flag.Args()[1])
	if err != nil {
		log.Fatalf("Can't open %s: %v", flag.Args()[1], err)
	}
	go reader(f1, c1)
	go reader(f2, c2)
	go outer(c1, c2, c)

	for {
		out, ok := <-c
		if !ok {
			break
		}

		line := ""
		if !*s1 {
			line += out.s1
		}
		line += "\t"
		if !*s2 {
			line += out.s2
		}
		line += "\t"
		if !*s3 {
			line += out.s3
		}
		if line != "\t\t" {
			fmt.Println(strings.TrimRight(line, "\t")) // the unix comm utility does this
		}
	}
}
