// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"pack.ag/tftp"
)

type Flags struct {
	cmd     string
	mode    string
	rPort   string
	literal bool
	verbose bool
}

func main() {
	f := Flags{}
	flag.StringVarP(&f.cmd, "c", "c", "", "Execute command as if it had been entered on the tftp prompt.  Must be specified last on the command line.")
	flag.StringVarP(&f.mode, "m", "m", "netascii", "Set the default transfer mode to mode.  This is usually used with -c.")
	flag.StringVarP(&f.rPort, "R", "R", "", "Force the originating port number to be in the specified range of port numbers.")
	flag.BoolVarP(&f.literal, "l", "l", false, "Default to literal mode. Used to avoid special processing of ':' in a file name.")
	flag.BoolVarP(&f.verbose, "v", "v", false, "Default to verbose mode.")

	flag.Parse()

	if err := run(f, os.Args[1:], flag.Args(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(f Flags, cmdline, args []string, stdin io.Reader, stdout io.Writer) error {
	// If we have IP/Host/Port supplied before command, ipPort holds this information.
	cmdArgs, ipPort := splitArgs(cmdline, args)

	if len(ipPort) < 1 || f.cmd == "" {
		return runInteractive(f, ipPort, stdin, stdout)
	}

	// Deconstruct files and look for hosts in the supplied cmdArgs
	// Only if "put" or "get"
	files := make([]string, 0)
	if f.cmd == "put" || f.cmd == "get" {
		hosts := make([]string, 0)
		for _, file := range cmdArgs {
			if !strings.Contains(file, ":") || f.literal {
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

type clientCfg struct {
	host    string
	port    string
	client  ClientIf
	mode    tftp.TransferMode
	rexmt   tftp.ClientOpt
	timeout tftp.ClientOpt
	trace   bool
	literal bool
	verbose bool
}

func runInteractive(f Flags, ipPort []string, stdin io.Reader, stdout io.Writer) error {
	const defaultPort = "69"
	var ipHost string
	var port string
	inScan := bufio.NewScanner(stdin)

	if len(ipPort) == 0 {
		ipHost = readHostInteractive(inScan, stdout)
	} else {
		ipHost = ipPort[0]

		if len(ipPort) > 1 {
			port = ipPort[1]
		} else {
			port = defaultPort
		}
	}

	clientcfg := &clientCfg{
		host:    ipHost,
		port:    port,
		mode:    tftp.ModeNetASCII,
		rexmt:   tftp.ClientRetransmit(10),
		timeout: tftp.ClientTimeout(1),
		trace:   false,
		literal: f.literal,
	}

	for {
		input := readInputInteractive(inScan, stdout)
		exit, err := executeOp(input, clientcfg, stdout)
		if err != nil {
			fmt.Fprintf(stdout, "%v", err)
		}
		if exit {
			return nil
		}
	}
}

func executeOp(input []string, clientcfg *clientCfg, stdout io.Writer) (bool, error) {
	var err error

	switch input[0] {
	case "q", "quit":
		return true, nil
	case "h", "help", "?":
		fmt.Fprintf(stdout, "%s", printHelp())
	case "ascii":
		clientcfg.mode, _ = validateMode("ascii")
	case "binary":
		clientcfg.mode, _ = validateMode("binary")
	case "mode":
		if len(input) > 1 {
			clientcfg.mode, err = validateMode(input[1])
			if err != nil {
				fmt.Fprintf(stdout, "%v", err)

			}
		}
		fmt.Fprintf(stdout, "Using %s mode to transfer files.\n", clientcfg.mode)
	case "get":
		clientcfg.client, err = NewClient(clientcfg)
		if err != nil {
			return false, err
		}

		err = executeGet(clientcfg.client, clientcfg.host, clientcfg.port, input[1:])
	case "put":
		clientcfg.client, err = NewClient(clientcfg)
		if err != nil {
			return false, err
		}

		err = executePut(clientcfg.client, clientcfg.host, clientcfg.port, input[1:])
	case "connect":
		if len(input) > 1 {
			clientcfg.port = input[2]
		}
		clientcfg.host = input[1]
	case "literal":
		clientcfg.literal = !clientcfg.literal
		fmt.Fprintf(stdout, "Literal mode is %s\n", statusString(clientcfg.literal))
	case "rexmt":
		var val int
		val, err = strconv.Atoi(input[1])

		clientcfg.rexmt = tftp.ClientRetransmit(val)
	case "status":
		fmt.Fprintf(stdout, "Connected to %s\n", clientcfg.host)
		fmt.Fprintf(stdout, "Mode: %s Verbose: %s Tracing: %s Literal: %s\n",
			clientcfg.mode,
			statusString(clientcfg.verbose),
			statusString(clientcfg.trace),
			statusString(clientcfg.literal),
		)
	case "timeout":
		var val int
		val, err = strconv.Atoi(input[1])

		clientcfg.timeout = tftp.ClientTimeout(val)
	case "trace":
		clientcfg.trace = !clientcfg.trace
		fmt.Fprintf(stdout, "Packet tracing %s.\n", statusString(clientcfg.trace))
	case "verbose":
		clientcfg.verbose = !clientcfg.verbose
		fmt.Fprintf(stdout, "Verbose mode %s.\n", statusString(clientcfg.verbose))
	}
	if err != nil {
		fmt.Fprintf(stdout, "%v\n", err)
	}
	return false, nil
}

func constructURL(host, port, dir string, file string) string {
	var s strings.Builder
	fmt.Fprintf(&s, "tftp://%s:%s/", host, port)
	if dir != "" {
		fmt.Fprintf(&s, "%s/", dir)
	}
	fmt.Fprintf(&s, "%s", file)

	return s.String()
}

func statusString(state bool) string {
	if state {
		return "on"
	}
	return "off"
}

func printHelp() string {
	var s strings.Builder
	fmt.Fprintf(&s, "not implemented yet\n")
	return s.String()
}

func readInputInteractive(in *bufio.Scanner, out io.Writer) []string {
	fmt.Fprint(out, "tftp:> ")
	in.Scan()
	return strings.Split(in.Text(), " ")
}

func readHostInteractive(in *bufio.Scanner, out io.Writer) string {
	fmt.Fprint(out, "(to): ")
	in.Scan()
	return in.Text()
}

var errInvalidTransferMode = errors.New("invalid transfer mode")

func validateMode(mode string) (tftp.TransferMode, error) {
	var ret tftp.TransferMode
	switch tftp.TransferMode(mode) {
	case "ascii":
		ret = tftp.ModeNetASCII
	case "binary":
		ret = tftp.ModeOctet
	default:
		return ret, errInvalidTransferMode
	}
	return ret, nil
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
