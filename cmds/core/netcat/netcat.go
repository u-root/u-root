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
	"sync"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const usage = "netcat [go-style network address]"

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

// udpRemoteConn saves raddr from first connection and implement io.ReadWriter
// interface, so io.Copy will work
type udpRemoteConn struct {
	raddr   *net.UDPAddr
	conn    *net.UDPConn
	wg      *sync.WaitGroup
	once    *sync.Once
	stderr  io.Writer
	verbose bool
}

func (u *udpRemoteConn) Read(b []byte) (int, error) {
	n, raddr, err := u.conn.ReadFromUDP(b)
	if err != nil {
		return n, err
	}
	setRaddr := func() {
		u.raddr = raddr
		if u.verbose {
			fmt.Fprintln(u.stderr, "Connected to", raddr)
		}
		u.wg.Done()
	}
	u.once.Do(setRaddr)
	return n, nil
}

func (u *udpRemoteConn) Write(b []byte) (int, error) {
	// we can't answer without raddr, so waiting for incomming request
	u.wg.Wait()
	return u.conn.WriteToUDP(b, u.raddr)
}

func (c *cmd) connection() (io.ReadWriter, error) {
	switch c.network {
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
		if c.listen {
			ln, err := net.Listen(c.network, c.addr)
			if err != nil {
				return nil, err
			}
			if c.verbose {
				fmt.Fprintln(c.stderr, "Listening on", ln.Addr())
			}
			return ln.Accept()
		}
		return net.Dial(c.network, c.addr)
	case "udp", "udp4", "udp6":
		addr, err := net.ResolveUDPAddr(c.network, c.addr)
		if err != nil {
			return nil, err
		}
		if c.listen {
			conn, err := net.ListenUDP(c.network, addr)
			if err != nil {
				return nil, err
			}
			if c.verbose {
				fmt.Fprintln(c.stderr, "Listening on", conn.LocalAddr())
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			return &udpRemoteConn{conn: conn, wg: wg, once: &sync.Once{}, stderr: c.stderr, verbose: c.verbose}, nil
		}
		return net.DialUDP(c.network, nil, addr)
	default:
		return nil, fmt.Errorf("unsupported network type %q", c.network)
	}
}

func (c *cmd) run() error {
	conn, err := c.connection()
	if err != nil {
		return err
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
