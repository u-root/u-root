// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/complete"
	"github.com/u-root/u-root/pkg/termios"
)

var (
	debug = flag.Bool("d", true, "enable debug prints")
	v = log.Printf
)

func verbose(f string, a ...interface{}) {
	v(f+"\r\n", a...)
}

func main() {
	flag.Parse()
	if *debug {
		complete.Debug = verbose
	}
	t, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := t.Raw()
	defer t.Set(r)
	for {
		p, err := complete.NewPathCompleter()
		if err != nil {
			log.Fatal(err)
		}
		c := complete.NewMultiCompleter(complete.NewStringCompleter([]string{"exit"}), p)
		l := complete.NewLineReader(c, t, t)
		s, err := l.ReadOne()
		v("Readone: %v, %v", s, err)
		if err != nil && err != complete.EOL {
			log.Print(err)
			continue
		}
		if len(s) == 0 {
			continue
		}
		if s[0] == "exit" {
			break
		}
		// s[0] is either the match or what they typed so far.
		t.Write([]byte(" "))
		bin := s[0]
		var args []string
		for err == nil {
			c := complete.NewFileCompleter(".")
			l := complete.NewLineReader(c, t, t)
			s, err := l.ReadOne()
			args = append(args, s...)
			v("add %v", s)
			if err != nil {
				log.Print(err)
				break
			}
		}
		v("Done reading args")
		cmd := exec.Command(bin, args...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}
	}
}
