// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"pack.ag/tftp"
)

type Flags struct {
	Cmd       string
	Mode      string
	PortRange string
	Literal   bool
	Verbose   bool
}

type ClientCfg struct {
	Host    string
	Port    string
	Client  ClientIf
	Mode    tftp.TransferMode
	Rexmt   tftp.ClientOpt
	Timeout tftp.ClientOpt
	Trace   bool
	Literal bool
	Verbose bool
}

func RunInteractive(f Flags, ipPort []string, stdin io.Reader, stdout io.Writer) error {
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

	clientcfg := &ClientCfg{
		Host:    ipHost,
		Port:    port,
		Mode:    tftp.ModeNetASCII,
		Rexmt:   tftp.ClientRetransmit(10),
		Timeout: tftp.ClientTimeout(1),
		Trace:   false,
		Literal: f.Literal,
	}

	for {
		input := readInputInteractive(inScan, stdout)
		exit, err := ExecuteOp(input, clientcfg, stdout)
		if err != nil {
			fmt.Fprintf(stdout, "%v", err)
		}
		if exit {
			return nil
		}
	}
}

func ExecuteOp(input []string, clientcfg *ClientCfg, stdout io.Writer) (bool, error) {
	var err error

	switch input[0] {
	case "q", "quit":
		return true, nil
	case "h", "help", "?":
		fmt.Fprintf(stdout, "%s", printHelp())
	case "ascii":
		clientcfg.Mode, _ = ValidateMode("ascii")
	case "binary":
		clientcfg.Mode, _ = ValidateMode("binary")
	case "mode":
		if len(input) > 1 {
			clientcfg.Mode, err = ValidateMode(input[1])
			if err != nil {
				fmt.Fprintf(stdout, "%v", err)

			}
		}
		fmt.Fprintf(stdout, "Using %s mode to transfer files.\n", clientcfg.Mode)
	case "get":
		clientcfg.Client, err = NewClient(clientcfg)
		if err != nil {
			return false, err
		}

		err = executeGet(clientcfg.Client, clientcfg.Host, clientcfg.Port, input[1:])
	case "put":
		clientcfg.Client, err = NewClient(clientcfg)
		if err != nil {
			return false, err
		}

		err = executePut(clientcfg.Client, clientcfg.Host, clientcfg.Port, input[1:])
	case "connect":
		if len(input) > 2 {
			clientcfg.Port = input[2]
		}
		clientcfg.Host = input[1]
	case "literal":
		clientcfg.Literal = !clientcfg.Literal
		fmt.Fprintf(stdout, "Literal mode is %s\n", statusString(clientcfg.Literal))
	case "rexmt":
		var val int
		val, err = strconv.Atoi(input[1])

		clientcfg.Rexmt = tftp.ClientRetransmit(val)
	case "status":
		fmt.Fprintf(stdout, "Connected to %s\n", clientcfg.Host)
		fmt.Fprintf(stdout, "Mode: %s Verbose: %s Tracing: %s Literal: %s\n",
			clientcfg.Mode,
			statusString(clientcfg.Verbose),
			statusString(clientcfg.Trace),
			statusString(clientcfg.Literal),
		)
	case "timeout":
		var val int
		val, err = strconv.Atoi(input[1])

		clientcfg.Timeout = tftp.ClientTimeout(val)
	case "trace":
		clientcfg.Trace = !clientcfg.Trace
		fmt.Fprintf(stdout, "Packet tracing %s.\n", statusString(clientcfg.Trace))
	case "verbose":
		clientcfg.Verbose = !clientcfg.Verbose
		fmt.Fprintf(stdout, "Verbose mode %s.\n", statusString(clientcfg.Verbose))
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

var ErrInvalidTransferMode = errors.New("invalid transfer mode")

func ValidateMode(mode string) (tftp.TransferMode, error) {
	var ret tftp.TransferMode
	switch tftp.TransferMode(mode) {
	case "ascii":
		ret = tftp.ModeNetASCII
	case "binary":
		ret = tftp.ModeOctet
	default:
		return ret, ErrInvalidTransferMode
	}
	return ret, nil
}

type PutCmd struct {
	localfiles []string
	remotefile string
	remotedir  string
}

func executePut(client ClientIf, host, port string, files []string) error {
	ret := &PutCmd{}
	switch len(files) {
	case 1:
		// Put file to server
		ret.localfiles = append(ret.localfiles, files...)
	case 2:
		// files[0] == localfile
		ret.localfiles = append(ret.localfiles, files[0])
		// files[1] == remotefile
		ret.remotefile = files[1]
	default:
		// files[:len(files)-2] == localfiles,
		ret.localfiles = append(ret.localfiles, files[:len(files)-2]...)
		// files[len(files)-1] == remote-directory
		ret.remotedir = files[len(files)-1]
	}

	for _, file := range ret.localfiles {
		url := constructURL(host, port, "", file)

		if len(ret.localfiles) == 1 && ret.remotefile != "" {
			url = constructURL(host, port, "", ret.remotefile)
		} else if len(ret.localfiles) > 1 {
			url = constructURL(host, port, ret.remotedir, file)
		}

		locFile, err := os.Open(file)
		if err != nil {
			return err
		}

		fs, err := locFile.Stat()
		if err != nil {
			return err
		}
		if err := client.Put(url, locFile, fs.Size()); err != nil {
			return err
		}
	}

	return nil
}

type GetCmd struct {
	remotefiles []string
	localfile   string
}

var errSizeNoMatch = errors.New("data size of read and write mismatch")

func executeGet(client ClientIf, host, port string, files []string) error {
	ret := &GetCmd{}
	switch len(files) {
	case 1:
		// files[0] == remotefile
		ret.remotefiles = append(ret.remotefiles, files[0])
	case 2:
		// files[0] == remotefile
		ret.remotefiles = append(ret.remotefiles, files[0])
		// files[1] == localfile
		ret.localfile = files[1]
	default:
		// files... == remotefiles
		ret.remotefiles = append(ret.remotefiles, files...)
	}

	for _, file := range ret.remotefiles {
		resp, err := client.Get(constructURL(host, port, "", file))
		if err != nil {
			return err
		}

		localfile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0o666)
		if err != nil {
			return nil
		}
		defer localfile.Close()

		if ret.localfile != "" && len(ret.remotefiles) == 1 {
			localfile, err = os.OpenFile(ret.localfile, os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				return err
			}
		}

		datalen, err := resp.Size()
		if err != nil {
			return err
		}

		data := make([]byte, datalen)
		nR, err := resp.Read(data)
		if err != nil {
			return err
		}

		nW, err := localfile.Write(data)
		if err != nil {
			return err
		}

		if nR != nW {
			return errSizeNoMatch
		}
	}

	return nil
}
