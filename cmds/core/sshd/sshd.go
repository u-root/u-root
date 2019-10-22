// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/u-root/u-root/pkg/pty"
	"golang.org/x/crypto/ssh"
)

// The ssh package does not define these things so we will
type (
	ptyReq struct {
		TERM   string //TERM environment variable value (e.g., vt100)
		Col    uint32
		Row    uint32
		Xpixel uint32
		Ypixel uint32
		Modes  string //encoded terminal modes
	}
	execReq struct {
		Command string
	}
	exitStatusReq struct {
		ExitStatus uint32
	}
)

var (
	shells  = [...]string{"bash", "zsh", "elvish"}
	shell   = "/bin/sh"
	debug   = flag.Bool("d", false, "Enable debug prints")
	keys    = flag.String("keys", "authorized_keys", "Path to the authorized_keys file")
	privkey = flag.String("privatekey", "id_rsa", "Path of private key")
	ip      = flag.String("ip", "0.0.0.0", "ip address to listen on")
	port    = flag.String("port", "2022", "port to listen on")
	dprintf = func(string, ...interface{}) {}
)

// start a command
// TODO: use /etc/passwd, but the Go support for that is incomplete
func runCommand(c ssh.Channel, p *pty.Pty, cmd string, args ...string) error {
	var ps *os.ProcessState
	defer c.Close()

	if p != nil {
		log.Printf("Executing PTY command %s %v", cmd, args)
		p.Command(cmd, args...)
		if err := p.C.Start(); err != nil {
			dprintf("Failed to execute: %v", err)
			return err
		}
		defer p.C.Wait()
		go io.Copy(p.Ptm, c)
		go io.Copy(c, p.Ptm)
		ps, _ = p.C.Process.Wait()
	} else {
		e := exec.Command(cmd, args...)
		e.Stdin, e.Stdout, e.Stderr = c, c, c
		log.Printf("Executing non-PTY command %s %v", cmd, args)
		if err := e.Start(); err != nil {
			dprintf("Failed to execute: %v", err)
			return err
		}
		ps, _ = e.Process.Wait()
	}

	ws := ps.Sys().(syscall.WaitStatus)
	// TODO(bluecmd): If somebody wants we can send exit-signal to return
	// information about signal termination, but leave it until somebody needs
	// it.
	// if ws.Signaled() {
	// }
	if ws.Exited() {
		code := uint32(ws.ExitStatus())
		dprintf("Exit status %v", code)
		c.SendRequest("exit-status", false, ssh.Marshal(exitStatusReq{code}))
	}
	return nil
}

func newPTY(b []byte) (*pty.Pty, error) {
	ptyReq := &ptyReq{}
	err := ssh.Unmarshal(b, ptyReq)
	dprintf("newPTY: %q", ptyReq)
	if err != nil {
		return nil, err
	}
	p, err := pty.New()
	if err != nil {
		return nil, err
	}
	ws, err := p.TTY.GetWinSize()
	if err != nil {
		return nil, err
	}
	ws.Row = uint16(ptyReq.Row)
	ws.Ypixel = uint16(ptyReq.Ypixel)
	ws.Col = uint16(ptyReq.Col)
	ws.Xpixel = uint16(ptyReq.Xpixel)
	dprintf("newPTY: Set winsizes to %v", ws)
	if err := p.TTY.SetWinSize(ws); err != nil {
		return nil, err
	}
	dprintf("newPTY: set TERM to %q", ptyReq.TERM)
	if err := os.Setenv("TERM", ptyReq.TERM); err != nil {
		return nil, err
	}
	return p, nil
}

func init() {
	for _, s := range shells {
		if _, err := exec.LookPath(s); err == nil {
			shell = s
		}
	}
}

func session(chans <-chan ssh.NewChannel) {
	var p *pty.Pty
	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "shell" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				dprintf("Request %v", req.Type)
				switch req.Type {
				case "shell":
					err := runCommand(channel, p, shell)
					req.Reply(true, []byte(fmt.Sprintf("%v", err)))
				case "exec":
					e := &execReq{}
					if err := ssh.Unmarshal(req.Payload, e); err != nil {
						log.Printf("sshd: %v", err)
						break
					}
					// Execute command using user's shell. This is what OpenSSH does
					// so it's the least surprising to the user.
					err := runCommand(channel, p, shell, "-c", e.Command)
					req.Reply(true, []byte(fmt.Sprintf("%v", err)))
				case "pty-req":
					p, err = newPTY(req.Payload)
					req.Reply(err == nil, nil)
				default:
					log.Printf("Not handling req %v %q", req, string(req.Payload))
					req.Reply(false, nil)
				}
			}
		}(requests)

	}
}

func main() {
	flag.Parse()
	if *debug {
		dprintf = log.Printf
	}
	// Public key authentication is done by comparing
	// the public key of a received connection
	// with the entries in the authorized_keys file.
	authorizedKeysBytes, err := ioutil.ReadFile(*keys)
	if err != nil {
		log.Fatal(err)
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			log.Fatal(err)
		}

		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		// Remove to disable public key auth.
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKeysMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	privateBytes, err := ioutil.ReadFile(*privkey)
	if err != nil {
		log.Fatal(err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal(err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", net.JoinHostPort(*ip, *port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept incoming connection: %s", err)
			continue
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Printf("failed to handshake: %v", err)
			continue
		}
		log.Printf("%v logged in with key %s", conn.RemoteAddr(), conn.Permissions.Extensions["pubkey-fp"])

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		go session(chans)
	}
}
