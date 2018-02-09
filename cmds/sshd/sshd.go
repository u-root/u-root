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
)

var (
	shells  = [...]string{"bash", "zsh", "rush"}
	shell   = "/bin/sh"
	debug   = flag.Bool("d", false, "Enable debug prints")
	keys    = flag.String("keys", "authorized_keys", "Path to the authorized_keys file")
	privkey = flag.String("privatekey", "id_rsa", "Path of private key")
	ip      = flag.String("ip", "0.0.0.0", "ip address to listen on")
	port    = flag.String("port", "2022", "port to listen on")
	dprintf = func(string, ...interface{}) {}
)

func echoCopy(w io.Writer, r io.Reader) (int64, error) {
	var b [8192]byte
	var err error
	var tot int64
	for err == nil {
		var amt int
		if amt, err = r.Read(b[:]); err != nil || amt < 1 {
			fmt.Printf("Read: %v", err)
			break
		}
		log.Printf("Read %d bytes: %q\n", amt, b[:amt])
		if _, err = w.Write(b[:amt]); err != nil {
			fmt.Printf("Write: %v", err)
		}
		tot += int64(amt)
	}
	return tot, err
}

// start a shell
// TODO: use /etc/passwd, but the Go support for that is incomplete
func runShell(c ssh.Channel, p *pty.Pty, shell string) error {
	copy := io.Copy
	defer c.Close()

	p.Command(shell)
	if err := p.C.Start(); err != nil {
		return err
	}
	defer p.C.Wait()
	if *debug {
		copy = echoCopy
	}
	go copy(p.Ptm, c)
	go copy(c, p.Ptm)
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
			log.Fatalf("Could not accept channel: %v", err)
		}

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "shell" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				dprintf("Request %v", req.Type)
				switch req.Type {
				case "shell":
					if p == nil {
						p, err = pty.New()
						if err != nil {
							log.Printf("sshd: pty.New failed, not running a shell")
							break
						}
					}
					err := runShell(channel, p, shell)
					req.Reply(true, []byte(fmt.Sprintf("%v", err)))
				case "pty-req":
					p, err = newPTY(req.Payload)
					req.Reply(err == nil, nil)
				default:
					fmt.Printf("Not handling req %v %q", req, string(req.Payload))
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
	listener, err := net.Listen("tcp", *ip+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept incoming connection: %s", err)
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Fatalf("failed to handshake: %v", err)
		}
		log.Printf("logged in with key %s", conn.Permissions.Extensions["pubkey-fp"])

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		go session(chans)
	}
}
