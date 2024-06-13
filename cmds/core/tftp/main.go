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

	flag "github.com/spf13/pflag"
	tftppkg "github.com/u-root/u-root/pkg/tftp"
	"pack.ag/tftp"
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

	ip, port := getIPPort(ipPort)

	m, err := tftppkg.ValidateMode(f.Mode)
	if err != nil {
		return err
	}

	clientcfg := &tftppkg.ClientCfg{
		Host:    ip,
		Port:    port,
		Mode:    m,
		Rexmt:   tftp.ClientRetransmit(10),
		Timeout: tftp.ClientTimeout(1),
		Trace:   false,
		Literal: f.Literal,
		Verbose: f.Verbose,
	}

	input := make([]string, 0)
	input = append(input, f.Cmd)
	input = append(input, cmdArgs...)

	if _, err := tftppkg.ExecuteOp(input, clientcfg, os.Stdout); err != nil {
		return err
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

func getIPPort(ipPort []string) (string, string) {
	const defaultPort = "69"
	if len(ipPort) == 2 {
		return ipPort[0], ipPort[0]
	}
	return ipPort[0], defaultPort
}
