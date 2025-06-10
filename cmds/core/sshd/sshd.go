// Copyright 2018-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !windows

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/pkg/sftp"
	"github.com/u-root/u-root/pkg/pty"
	"golang.org/x/crypto/ssh"
)

// The ssh package does not define these things so we will
type (
	ptyReq struct {
		TERM   string // TERM environment variable value (e.g., vt100)
		Col    uint32
		Row    uint32
		Xpixel uint32
		Ypixel uint32
		Modes  string // encoded terminal modes
	}
	execReq struct {
		Command string
	}
	exitStatusReq struct {
		ExitStatus uint32
	}
)

var (
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
		// execute command and wait for response
		if err := e.Run(); err != nil {
			dprintf("Failed to execute: %v", err)
			return err
		}

		ps = e.ProcessState
	}

	// TODO(bluecmd): If somebody wants we can send exit-signal to return
	// information about signal termination, but leave it until somebody needs
	// it.
	// if ws.Signaled() {
	// }
	if ps.Exited() {
		code := uint32(ps.ExitCode())
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

func handleSftp(req *ssh.Request, channel ssh.Channel) error {
	s, err := sftp.NewServer(channel, sftp.WithServerWorkingDirectory("/"))
	if err != nil {
		req.Reply(false, nil)
		return err
	}
	if err := s.Serve(); err != nil {
		req.Reply(false, nil)
		return err
	}

	req.Reply(true, nil)

	// Need to tell the client that the operations was a success (0) and
	// kill the connection by any means necessary.  You may see stuff like
	// "read failed... Broken pipe", but that's OK!  (probably).
	//
	// An openssh server will send a msgChannelEOF (96),
	// msgChannelRequest(exit-status) (98), and then msgChannelClose (97).
	channel.CloseWrite() // EOF
	channel.SendRequest("exit-status", false, ssh.Marshal(exitStatusReq{0}))
	channel.Close()

	return nil
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
				case "subsystem":
					switch {
					case string(req.Payload[4:]) == "sftp":
						// This handles the req.Reply
						// and may close the channel
						err := handleSftp(req, channel)
						if err != nil {
							log.Printf("sshd: sftp: %v", err)
						}
					default:
						log.Printf("Not handling subsystem req %v %q",
							req, string(req.Payload))
						req.Reply(false, nil)
					}
				default:
					log.Printf("Not handling req %v %q", req, string(req.Payload))
					req.Reply(false, nil)
				}
			}
		}(requests)

	}
}

type params struct {
	keys    string
	privkey string
	ip      string
	port    string
	debug   bool
}

func parseParams() params {
	return params{
		debug:   *debug,
		keys:    *keys,
		privkey: *privkey,
		ip:      *ip,
		port:    *port,
	}
}

type cmd struct {
	params
}

func command(p params) *cmd {
	return &cmd{
		params: p,
	}
}

func (c *cmd) run() error {
	if c.debug {
		dprintf = log.Printf
	}
	// Public key authentication is done by comparing
	// the public key of a received connection
	// with the entries in the authorized_keys file.
	authorizedKeysBytes, err := os.ReadFile(c.keys)
	if err != nil {
		return err
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return err
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

	privateBytes, err := os.ReadFile(c.privkey)
	if err != nil {
		return err
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return err
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", net.JoinHostPort(c.ip, c.port))
	if err != nil {
		return err
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

func main() {
	flag.Parse()
	if err := command(parseParams()).run(); err != nil {
		log.Fatal(err)
	}
}
