// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tftp gets and puts files from and to a tftp server
//
// Synopsis: tftp [ options... ] [host [port]] [-c command]
//

package main

import (
	"io"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	tftppkg "github.com/u-root/u-root/pkg/tftp"
)

func main() {
	f := tftppkg.Flags{}
	flag.StringVarP(&f.Cmd, "c", "c", "", "Execute command as if it had been entered on the tftp prompt.  Must be specified last on the command line.")
	flag.StringVarP(&f.Mode, "m", "m", "netascii", "Set the default transfer mode to mode.  This is usually used with -c.")
	flag.StringVarP(&f.PortRange, "R", "R", "", "Force the originating port number to be in the specified range of port numbers.")
	flag.BoolVarP(&f.Literal, "l", "l", false, "Default to literal mode. Used to avoid special processing of ':' in a file name.")
	flag.BoolVarP(&f.Verbose, "v", "v", false, "Default to verbose mode.")

	flag.Parse()

	if err := run(f, os.Args[1:], flag.Args(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(f tftppkg.Flags, cmdline, args []string, stdin io.Reader, stdout io.Writer) error {
	// If we have IP/Host/Port supplied before command, ipPort holds this information.
	cmdArgs, ipPort := splitArgs(cmdline, args)

	if len(ipPort) < 1 || f.Cmd == "" {
		return tftppkg.RunInteractive(f, ipPort, stdin, stdout)
	}

	// Deconstruct files and look for hosts in the supplied cmdArgs
	// Only if "put" or "get"
	files := make([]string, 0)
	if f.Cmd == "put" || f.Cmd == "get" {
		hosts := make([]string, 0)
		for _, file := range cmdArgs {
			if !strings.Contains(file, ":") || f.Literal {
				files = append(files, file)
				continue
			}

			splitFile := strings.Split(file, ":")
			hosts = append(hosts, splitFile[0])
			files = append(files, splitFile[1])
		}

		if len(hosts) > 0 {
			// Use the last host/ip from host as stated in the man page of tftp
			if len(ipPort) > 0 {
				ipPort[0] = hosts[len(hosts)-1]
			} else {
				ipPort = append(ipPort, hosts[len(hosts)-1])
			}
		}
	}

	return nil
}

func splitArgs(cmdline, args []string) ([]string, []string) {
	retCmdArgs := make([]string, 0)
	retIPPort := make([]string, 0)
	for i := len(cmdline) - 1; i > 0; i-- {
		if cmdline[i] == "-c" {
			retCmdArgs = append(retCmdArgs, cmdline[i+2:]...)
		}
	}

	retIPPort = append(retIPPort, args[:len(args)-len(retCmdArgs)]...)

	return retCmdArgs, retIPPort
}
