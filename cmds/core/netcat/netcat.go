// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// netcat creates arbitrary TCP and UDP connections and listens and sends arbitrary data.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const usage = "netcat host port"

var errMissingHostnameOrPort = fmt.Errorf("missing hostname or port")

type params struct {
	network string
	listen  bool
	verbose bool
	host    string
	port    string
}

func parseParams() params {
	netType := flag.String("net", "tcp", "What net type to use, e.g. tcp, unix, etc.")
	listen := flag.Bool("l", false, "Listen for connections.")
	verbose := flag.Bool("v", false, "Verbose output.")
	flag.Parse()

	return params{
		network: *netType,
		listen:  *listen,
		verbose: *verbose,
	}
}

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	addr   string
	params
}

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string, p params) (*cmd, error) {
	if len(args) != 2 {
		return nil, errMissingHostnameOrPort
	}

	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		params: p,
		addr:   net.JoinHostPort(args[0], args[1]),
	}, nil
}

func init() {
	flag.Usage = util.Usage(flag.Usage, usage)
}

func main() {
	p := parseParams()
	c, err := command(os.Stdin, os.Stdout, os.Stderr, flag.Args(), p)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	if err := c.run(); err != nil {
		log.Fatalf("netcat: %v", err)
	}
}

func (c *cmd) run() error {
	var conn net.Conn
	var err error

	if c.listen {
		ln, err := net.Listen(c.network, c.addr)
		if err != nil {
			return err
		}
		if c.verbose {
			fmt.Fprintln(c.stderr, "Listening on", ln.Addr())
		}

		conn, err = ln.Accept()
		if err != nil {
			return err
		}
	} else {
		if conn, err = net.Dial(c.network, c.addr); err != nil {
			return err
		}
	}
	if c.verbose {
		fmt.Fprintln(c.stderr, "Connected to", conn.RemoteAddr())
	}

	go func() {
		if _, err := io.Copy(conn, c.stdin); err != nil {
			fmt.Fprintln(c.stderr, err)
		}
	}()
	if _, err = io.Copy(c.stdout, conn); err != nil {
		fmt.Fprintln(c.stderr, err)
	}
	if c.verbose {
		fmt.Fprintln(c.stderr, "Disconnected")
	}

	return nil
}
