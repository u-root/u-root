// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty" // TODO: get rid of krpty
	flag "github.com/spf13/pflag"
)

var (
	pubKeyFile = flag.StringP("pubkeyfile", "k", "key.pub", "file for public key")
	port       = flag.StringP("port", "p", "2222", "default port")
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func handler(s ssh.Session) {
	var a []string
	if len(s.Command()) > 0 {
		a = append([]string{"-c"}, s.Command()...)
	}
	cmd := exec.Command("/bin/sh", a...)
	cmd.Env = append(cmd.Env, s.Environ()...)
	ptyReq, winCh, isPty := s.Pty()
	if isPty {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
		f, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}
		go func() {
			for win := range winCh {
				setWinsize(f, win.Width, win.Height)
			}
		}()
		go func() {
			io.Copy(f, s) // stdin
		}()
		io.Copy(s, f) // stdout
	} else {
		cmd.Stdin, cmd.Stdout, cmd.Stderr = s, s, s
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
func main() {
	flag.Parse()
	publicKeyOption := func(ctx ssh.Context, key ssh.PublicKey) bool {
		// Glob the users's home directory for all the
		// possible keys?
		data, err := ioutil.ReadFile(*pubKeyFile)
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

	log.Println("starting ssh server on port " + *port)
	log.Fatal(server.ListenAndServe())
}
