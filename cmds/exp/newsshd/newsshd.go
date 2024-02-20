// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	flag "github.com/spf13/pflag"
)

var (
	hostKeyFile = flag.StringP("hostkeyfile", "h", "/etc/ssh_host_rsa_key", "file for host key")
	pubKeyFile  = flag.StringP("pubkeyfile", "k", "key.pub", "file for public key")
	port        = flag.StringP("port", "p", "2222", "default port")
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

func main() {
	flag.Parse()
	publicKeyOption := func(ctx ssh.Context, key ssh.PublicKey) bool {
		// Glob the users's home directory for all the
		// possible keys?
		data, err := os.ReadFile(*pubKeyFile)
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
		Addr:             ":" + *port,
		PublicKeyHandler: publicKeyOption,
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
			log.Println("attempt to bind", host, port, "granted")
			return true
		}),
		Handler: handler,
	}

	server.SetOption(ssh.HostKeyFile(*hostKeyFile))
	log.Println("starting ssh server on port " + *port)
	log.Fatal(server.ListenAndServe())
}
