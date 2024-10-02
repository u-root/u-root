// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

func handler(s ssh.Session) {
	var a []string
	if len(s.Command()) > 0 {
		a = append([]string{"-c"}, strings.Join(s.Command(), " "))
	}
	cmd := exec.Command("/bin/sh", a...)
	cmd.Env = append(cmd.Env, s.Environ()...)
	ptyReq, winCh, isPty := s.Pty()
	if isPty {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
		f, err := pty.Start(cmd)
		if err != nil {
			log.Print(err)
			return
		}
		go func() {
			for win := range winCh {
				pty.Setsize(f, &pty.Winsize{Rows: uint16(win.Height), Cols: uint16(win.Width)})
			}
		}()
		go func() {
			io.Copy(f, s) // stdin
		}()
		io.Copy(s, f) // stdout
	} else {
		cmd.Stdin, cmd.Stdout, cmd.Stderr = s, s, s
		if err := cmd.Run(); err != nil {
			log.Print(err)
			return
		}
	}
}

type cmd struct {
	hostKeyFile string
	pubKeyFile  string
	port        string
}

func (c *cmd) run() error {
	publicKeyOption := func(ctx ssh.Context, key ssh.PublicKey) bool {
		// Glob the users's home directory for all the
		// possible keys?
		data, err := os.ReadFile(c.pubKeyFile)
		if err != nil {
			fmt.Print(err)
			return false
		}
		allowed, _, _, _, _ := ssh.ParseAuthorizedKey(data)
		return ssh.KeysEqual(key, allowed)
	}

	server := ssh.Server{
		LocalPortForwardingCallback: ssh.LocalPortForwardingCallback(func(ctx ssh.Context, dhost string, dport uint32) bool {
			log.Println("Accepted forward", dhost, dport)
			return true
		}),
		Addr:             ":" + c.port,
		PublicKeyHandler: publicKeyOption,
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
			log.Println("attempt to bind", host, port, "granted")
			return true
		}),
		Handler: handler,
	}

	server.SetOption(ssh.HostKeyFile(c.hostKeyFile))
	log.Println("starting ssh server on port " + c.port)

	return (server.ListenAndServe())
}

func command(args []string) *cmd {
	c := &cmd{}
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&c.hostKeyFile, "hostkeyfile", "/etc/ssh_host_rsa_key", "file for host key")
	f.StringVar(&c.hostKeyFile, "h", "/etc/ssh_host_rsa_key", "file for host key (shorthand)")

	f.StringVar(&c.pubKeyFile, "pubkeyfile", "key.pub", "file for public key")
	f.StringVar(&c.pubKeyFile, "k", "key.pub", "file for public key (shorthand)")

	f.StringVar(&c.port, "port", "2222", "default port")
	f.StringVar(&c.port, "p", "2222", "default port (shorthand)")

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))

	return c
}

func main() {
	if err := command(os.Args).run(); err != nil {
		log.Fatal(err)
	}
}
