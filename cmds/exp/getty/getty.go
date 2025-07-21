// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// getty Open a TTY and invoke a shell
// There are no special options and no login support
// Also getty exits after starting the shell so if one exits the shell, there
// is no more shell!
//
// Synopsis:
//
//	getty <port> <baud> [term]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/u-root/u-root/pkg/termios"
	"github.com/u-root/u-root/pkg/upath"
)

var (
	verbose = flag.Bool("v", false, "verbose log")
	debug   = func(string, ...any) {}
	cmdList []string
	envs    []string
)

func init() {
	r := upath.UrootPath
	cmdList = []string{
		r("/bin/defaultsh"),
		r("/bin/sh"),
	}
}

func main() {
	flag.Parse()

	if *verbose {
		debug = log.Printf
	}

	port := flag.Arg(0)
	baud, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		baud = 0
	}
	term := flag.Arg(2)

	log.SetPrefix("getty: ")

	ttyS, err := termios.NewTTYS(port)
	if err != nil {
		log.Fatalf("Unable to open port %s: %v", port, err)
	}

	if _, err := ttyS.Serial(baud); err != nil {
		log.Printf("Unable to configure port %s and set baudrate %d: %v", port, baud, err)
	}

	// Output the u-root banner
	log.New(ttyS, "", log.LstdFlags).Printf("Welcome to u-root!")
	fmt.Fprintln(ttyS, `                              _`)
	fmt.Fprintln(ttyS, `   _   _      _ __ ___   ___ | |_`)
	fmt.Fprintln(ttyS, `  | | | |____| '__/ _ \ / _ \| __|`)
	fmt.Fprintln(ttyS, `  | |_| |____| | | (_) | (_) | |_`)
	fmt.Fprintln(ttyS, `   \__,_|    |_|  \___/ \___/ \__|`)
	fmt.Fprintln(ttyS)

	if term != "" {
		err = os.Setenv("TERM", term)
		if err != nil {
			debug("Unable to set 'TERM=%s': %v", port, err)
		}
	}
	envs = os.Environ()
	debug("envs %v", envs)

	for _, v := range cmdList {
		debug("Trying to run %v", v)
		if _, err := os.Stat(v); os.IsNotExist(err) {
			debug("%v", err)
			continue
		}

		cmd := exec.Command(v)
		cmd.Env = envs
		ttyS.Ctty(cmd)
		debug("running %v", cmd)
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", v, err)
			continue
		}
		if err := cmd.Process.Release(); err != nil {
			log.Printf("Error releasing process %v:%v", v, err)
		}
		// stop after first valid command
		return
	}
	log.Printf("No suitable executable found in %+v", cmdList)
}
