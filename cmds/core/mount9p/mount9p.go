// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
//
// Mount9p mounts a 9P server with the transport=fd method, optionally using SSH
// for the transport.
//
//	mount9p -ssh -private_key ~/.ssh/id_rsa /path/to/mount [::1]:1234

//go:build linux

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/u-root/u-root/pkg/ssh9p"
	"github.com/u-root/u-root/pkg/sshstream"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"
)

var (
	onto        string // Where to mount onto.
	addr        string // Address:port to dial.
	useSSH      = flag.Bool("ssh", false, "Use ssh for transport")
	privateKey  = flag.String("private_key", "id_rsa", "Path to your private key")
	fastTimeout = flag.Bool("fast_timeout", false, "Quickly timeout server connections")
	cache       = flag.String("cache", "none", "Cache mode")
)

// Run does the mount and blocks if necessary.  Cancel the context to unblock
// it and kill the client side of the mount server.
func run(ctx context.Context) error {
	var cfg *ssh.ClientConfig
	if *useSSH {
		ncfg, err := sshstream.NewClientConfig(*privateKey, ssh.InsecureIgnoreHostKey())
		if err != nil {
			return err
		}
		cfg = ncfg
	}

	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer c.Close()

	ready := make(chan bool)
	go func() {
		if <-ready {
			fmt.Println("Mount ready at", onto)
		}
	}()
	// If we blocked for ssh, err will be non-nil.
	return ssh9p.Mount9P(ctx, ready, c, onto, ssh9p.WithSSHClient(cfg), ssh9p.WithUnixFlags(unix.MS_NOSUID|unix.MS_NODEV), ssh9p.WithFastTCPTimeout(*fastTimeout), ssh9p.WithCache(*cache))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Printf("Usage: %s [options] <onto> <addr>\n\n", os.Args[0])
		flag.Usage()
		return
	}
	onto = flag.Arg(0)
	addr = flag.Arg(1)

	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
