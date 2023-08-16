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

const usage = "netcat [go-style network address]"
const bufSize = 1500

var errMissingHostnamePort = fmt.Errorf("missing hostname:port")

type params struct {
	network string
	listen  bool
	verbose bool
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

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, p params, args []string) (*cmd, error) {
	if len(args) < 1 {
		return nil, errMissingHostnamePort
	}

	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		params: p,
		addr:   args[0],
	}, nil
}

func init() {
	flag.Usage = util.Usage(flag.Usage, usage)
}

func main() {
	p := parseParams()
	c, err := command(os.Stdin, os.Stdout, os.Stderr, p, flag.Args())
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	if err = c.run(); err != nil {
		log.Fatalf("netcat: %v", err)
	}
}

func (c *cmd) run() error {
	switch c.network {
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
		return c.run(func() (net.Conn, error) {
			l, err :=  net.Listen(c.network, c.addr)
			if err != nil {
				return nil, err
			}
		return l.Accept()
		})
	case "udp", "udp4", "udp6":
		return c.run(func() (net.Conn, error) {
			addr, err := net.ResolveUDPAddr(c.network, c.addr)
			if err != nil {
				return nil, err
			}
			return net.ListenUDP(c.network, addr)
		})
	default:
		return fmt.Errorf("unsupported network type %q", c.network)
	}
}

func (c *cmd) run(listen func() (net.Conn, error)) error {
	var conn net.Conn
	var err error

	if c.listen {
		conn, err = listen()
		if err != nil {
			return err
		}
		if c.verbose {
			fmt.Fprintln(c.stderr, "Listening on", conn.LocalAddr())
		}

	} else {
		if conn, err = net.Dial(c.network, c.addr); err != nil {

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
