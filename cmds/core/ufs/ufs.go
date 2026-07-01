// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
//
// Ufs runs a 9p server.
//
// Start the server:
//
//	ufs :3333 -root /path/to/share
//
// Then, connect using the Linux 9P filesystem:
//
//	mount -t 9p -o trans=tcp,port=3333 127.0.0.1 /mnt
//
// IPv6 mounts either require a fairly recent Linux, containing a22a29655c42
// ("net/9p/fd: support ipv6 for trans=tcp"), or the use of a helper program
// using 9p_fd, e.g. mount9p.
//
// To serve using ssh:
//
//	ufs :3333 -ssh \
//		-authorized_keys ~/.ssh/authorized_keys \
//		-hostkey ~/.ssh/id_rsa
//
// You'll need to use the mount9p client to connect to the ssh-enabled server.

//go:build linux

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hugelgupf/p9/fsimpl/localfs"
	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/u-root/pkg/p9dev"
	"github.com/u-root/u-root/pkg/ssh9p"
	"github.com/u-root/u-root/pkg/sshstream"
	"golang.org/x/crypto/ssh"
)

var (
	addr           string // IP address:port to listen on.
	root           = flag.String("root", "/", "Root of file system to serve")
	useSSH         = flag.Bool("ssh", false, "Use ssh for transport")
	authorizedKeys = flag.String("authorized_keys", "authorized_keys", "Path to the authorized_keys file")
	hostkey        = flag.String("hostkey", "id_rsa", "Path to the host's private key")
	// Flag is nodevmap, just like in 9p2000.u.
	noDevMap = flag.Bool("nodevmap", false, "Pretend special device files are regular files")
)

func run(ctx context.Context) error {
	var cfg *ssh.ServerConfig
	if *useSSH {
		_cfg, err := sshstream.NewServerConfig(*authorizedKeys, *hostkey)
		if err != nil {
			return fmt.Errorf("cfg: %w", err)
		}
		cfg = _cfg
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("bind: %w", err)
	}
	defer l.Close()

	var attacher p9.Attacher
	attacher = localfs.Attacher(*root)
	if *noDevMap {
		attacher = p9dev.New(attacher)
	}

	server := ssh9p.NewServer(p9.NewServer(attacher), ssh9p.WithSSHServer(cfg))
	if err := server.Serve(ctx, l); err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [options] <addr>\n\n", os.Args[0])
		flag.Usage()
		return
	}
	addr = flag.Arg(0)

	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
