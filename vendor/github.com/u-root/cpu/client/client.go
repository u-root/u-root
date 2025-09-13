// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hugelgupf/p9/p9"
	"github.com/mdlayher/vsock"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

const (
	// From setting up the forward to having the nonce written back to us,
	// we would like to default to 100ms. This is a lot, considering that at this point,
	// the sshd has forked a server for us and it's waiting to be
	// told what to do.
	defaultTimeOut = time.Duration(100 * time.Millisecond)

	// DefaultNameSpace is the default used if the user does not request
	// something else.
	DefaultNameSpace = "/lib:/lib64:/usr:/bin:/etc:/home"
)

// v allows debug printing.
// Do not call it directly, call verbose instead.
var v = func(string, ...interface{}) {}

// Cmd is a cpu client.
// It implements as much of exec.Command as makes sense.
type Cmd struct {
	config  ssh.ClientConfig
	client  *ssh.Client
	session *ssh.Session
	// CPU-specific options.
	// As in exec.Command, these controls are exposed and can
	// be set directly.
	Host string
	// HostName as found in .ssh/config; set to Host if not found
	HostName          string
	Args              []string
	Root              string
	HostKeyFile       string
	PrivateKeyFile    string
	DisablePrivateKey bool
	Port              string
	Timeout           time.Duration
	Env               []string
	SessionIn         io.WriteCloser
	SessionOut        io.Reader
	SessionErr        io.Reader
	Stdin             io.Reader
	Stdout            io.WriteCloser
	Stderr            io.WriteCloser
	Row               int
	Col               int
	hasTTY            bool // Set if we have a TTY
	// NameSpace is a string as defined in the cpu documentation.
	NameSpace string
	// FSTab is an fstab(5)-format string
	FSTab string
	// Ninep determines if client will run a 9P server
	Ninep bool

	nonce      nonce
	network    string // This is a variable but we expect it will always be tcp
	port9p     uint16 // port on which we serve 9p
	cmd        string // The command is built up, bit by bit, as we configure the client
	closers    []func() error
	fileServer p9.Attacher
}

// SetOptions sets various options into the Command.
func (c *Cmd) SetOptions(opts ...Set) error {
	for _, o := range opts {
		if err := o(c); err != nil {
			return err
		}
	}
	return nil
}

// SetVerbose sets the verbose printing function.
// e.g., one might call SetVerbose(log.Printf)
func SetVerbose(f func(string, ...interface{})) {
	v = f
}

// Listen implements net.Listen on the ssh socket.
func (c *Cmd) Listen(n, addr string) (net.Listener, error) {
	return c.client.Listen(n, addr)
}

func sameFD(w io.WriteCloser, std *os.File) bool {
	if file, ok := w.(*os.File); ok {
		return file.Fd() == std.Fd()
	}
	return false
}

// Command implements exec.Command. The required parameter is a host.
// The args arg args to $SHELL. If there are no args, then starting $SHELL
// is assumed.
func Command(host string, args ...string) *Cmd {
	var hasTTY bool
	if len(args) == 0 {
		shell, ok := os.LookupEnv("SHELL")
		// We've found in some cases SHELL is not set!
		if !ok {
			shell = "/bin/sh"
		}
		args = []string{shell}
	}

	col, row := 80, 40
	if c, r, err := term.GetSize(int(os.Stdin.Fd())); err != nil {
		verbose("Can not get winsize: %v; assuming %dx%d and non-interactive", err, col, row)
	} else {
		hasTTY = true
		col, row = c, r
	}

	h, u := GetHostUser(host)
	return &Cmd{
		Host:     host,
		HostName: h,
		Args:     args,
		Port:     DefaultPort,
		Timeout:  defaultTimeOut,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Row:      row,
		Col:      col,
		config: ssh.ClientConfig{
			User:            u,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		hasTTY:  hasTTY,
		network: "tcp",
		// Safety first: if they want a namespace, they must say so
		Root: "",
	}
}

// Set is the type of function used to set options in SetOptions.
type Set func(*Cmd) error

// WithServer allows setting custom 9P servers.
// One use: should users wish to serve from a flattened
// docker container saved as a cpio or tar.
func WithServer(a p9.Attacher) Set {
	return func(c *Cmd) error {
		c.fileServer = a
		return nil
	}
}

// With9P enables the 9P2000 server in cpu. The server is by default disabled.
func With9P(p9 bool) Set {
	return func(c *Cmd) error {
		c.Ninep = p9
		return nil
	}
}

// WithNameSpace sets the namespace to Cmd.There is no default: having some default
// violates the principle of least surprise for package users.
func WithNameSpace(ns string) Set {
	return func(c *Cmd) error {
		c.NameSpace = ns
		return nil
	}
}

// WithFSTab reads a file for the FSTab member.
func WithFSTab(fstab string) Set {
	return func(c *Cmd) error {
		if len(fstab) == 0 {
			return nil
		}
		b, err := os.ReadFile(fstab)
		if err != nil {
			return fmt.Errorf("Reading fstab: %w", err)
		}
		c.FSTab = string(b)
		return nil
	}
}

// WithTimeout sets the 9p timeout.
func WithTimeout(timeout string) Set {
	return func(c *Cmd) error {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}

		c.Timeout = d
		return nil
	}
}

// WithPrivateKeyFile adds a private key file to a Cmd
func WithPrivateKeyFile(key string) Set {
	return func(c *Cmd) error {
		c.PrivateKeyFile = key
		return nil
	}
}

// WithDisablePrivateKey disables using private keys to encrypt the SSH
// connection.
func WithDisablePrivateKey(disable bool) Set {
	return func(c *Cmd) error {
		c.DisablePrivateKey = disable
		return nil
	}
}

// WithHostKeyFile adds a host key to a Cmd
func WithHostKeyFile(key string) Set {
	return func(c *Cmd) error {
		c.HostKeyFile = key
		return nil
	}
}

// WithRoot adds a root to a Cmd
func WithRoot(root string) Set {
	return func(c *Cmd) error {
		c.Root = root
		return nil
	}
}

// WithNetwork sets the network. This almost never needs
// to be set, save for vsock.
func WithNetwork(network string) Set {
	return func(c *Cmd) error {
		if len(network) > 0 {
			c.network = network
		}
		return nil
	}
}

// WithPort sets the port in the Cmd.
// It calls GetPort with the passed-in port
// before assigning it.
func WithPort(port string) Set {
	return func(c *Cmd) error {
		if len(port) == 0 {
			p, err := GetPort(c.HostName, c.Port)
			if err != nil {
				return err
			}
			port = p
		}

		c.Port = port
		return nil
	}

}

// It's a shame vsock is not in the net package (yet ... or ever?)
func vsockDial(host, port string) (net.Conn, string, error) {
	id, portid, err := vsockIDPort(host, port)
	verbose("vsock(%v, %v) = %v, %v, %v", host, port, id, portid, err)
	if err != nil {
		return nil, "", err
	}
	addr := fmt.Sprintf("%#x:%d", id, portid)
	conn, err := vsock.Dial(id, portid, nil)
	verbose("vsock id %#x port %s addr %#x conn %v err %v", id, port, addr, conn, err)
	return conn, addr, err

}

// https://github.com/firecracker-microvm/firecracker/blob/main/docs/vsock.md#host-initiated-connections
func unixVsockDial(path, port string) (net.Conn, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	connectMsg := fmt.Sprintf("CONNECT %s\n", port)
	if _, err := io.WriteString(conn, connectMsg); err != nil {
		return nil, fmt.Errorf("sending connect request: %w", err)
	}
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("reading connect request: %w", err)
	}
	if string(buf) != "OK" {
		return nil, fmt.Errorf("vsock: expect OK, got %s", buf)
	}
	buf = make([]byte, 1)
	for buf[0] != '\n' {
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

// Dial implements ssh.Dial for cpu.
// Additionaly, if Cmd.Root is not "", it
// starts up a server for 9p requests.
// Note that any bind parsing is deferred until this point,
// to avoid callers getting ordering of setting variables
// in the Cmd wrong.
func (c *Cmd) Dial() error {
	fstab, err := ParseBinds(c.NameSpace)
	if err != nil {
		return err
	}
	c.FSTab = JoinFSTab(c.FSTab, fstab)

	if err := c.UserKeyConfig(); err != nil {
		return err
	}
	// Sadly, no vsock in net package.
	var (
		conn net.Conn
		addr string
	)

	switch c.network {
	case "vsock":
		conn, addr, err = vsockDial(c.HostName, c.Port)
	case "unix", "unixgram", "unixpacket":
		// There is not port on a unix domain socket.
		addr = c.HostName
		conn, err = net.Dial(c.network, c.HostName)
	case "unix-vsock":
		addr = c.HostName
		conn, err = unixVsockDial(c.HostName, c.Port)
	default:
		addr = net.JoinHostPort(c.HostName, c.Port)
		conn, err = net.Dial(c.network, addr)
	}
	verbose("connect: err %v", err)
	if err != nil {
		return err
	}
	sshconn, chans, reqs, err := ssh.NewClientConn(conn, addr, &c.config)
	if err != nil {
		return err
	}
	cl := ssh.NewClient(sshconn, chans, reqs)
	verbose("cpu:ssh.Dial(%s, %s, %v): (%v, %v)", c.network, addr, c.config, cl, err)
	if err != nil {
		return fmt.Errorf("Failed to dial: %v", err)
	}

	c.client = cl
	// Specifying a root is required for a remote namespace.
	if len(c.Root) == 0 {
		return nil
	}

	// Arrange port forwarding from remote ssh to our server.
	// Note: cl.Listen returns a TCP listener with network "tcp"
	// or variants. This lets us use a listen deadline.
	if c.Ninep {
		l, err := cl.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			// If ipv4 isn't available, try ipv6.  It's not enough
			// to use Listen("tcp", "localhost:0a)", since we (the
			// cpu client) might have v4 (which the runtime will
			// use if we say "localhost"), but the server (cpud)
			// might not.
			l, err = cl.Listen("tcp", "[::1]:0")
			if err != nil {
				return fmt.Errorf("cpu client listen for forwarded 9p port %v", err)
			}
		}
		verbose("ssh.listener %v", l.Addr().String())
		ap := strings.Split(l.Addr().String(), ":")
		if len(ap) == 0 {
			return fmt.Errorf("Can't find a port number in %v", l.Addr().String())
		}
		port9p, err := strconv.ParseUint(ap[len(ap)-1], 0, 16)
		if err != nil {
			return fmt.Errorf("Can't find a 16-bit port number in %v", l.Addr().String())
		}
		c.port9p = uint16(port9p)

		verbose("listener %T %v addr %v port %v", l, l, l.Addr().String(), port9p)

		nonce, err := generateNonce()
		if err != nil {
			log.Fatalf("Getting nonce: %v", err)
		}
		c.nonce = nonce
		c.Env = append(c.Env, "CPUNONCE="+nonce.String())
		verbose("Set NONCE to %q", nonce.String())
		go func(l net.Listener) {
			if err := c.srv(l); err != nil {
				log.Printf("9p server error: %v", err)
			}
		}(l)
	}
	if len(c.FSTab) > 0 {
		c.Env = append(c.Env, "CPU_FSTAB="+c.FSTab)
		c.Env = append(c.Env, "LC_GLENDA_CPU_FSTAB="+c.FSTab)
	}

	return nil
}

func quoteArg(arg string) string {
	return "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
}

// Start implements exec.Start for CPU.
func (c *Cmd) Start() error {
	var err error
	if c.client == nil {
		return fmt.Errorf("Cmd has no client")
	}
	if c.session, err = c.client.NewSession(); err != nil {
		return err
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // do not disable echoing!
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if c.hasTTY {
		verbose("c.session.RequestPty(\"ansi\", %v, %v, %#x", c.Row, c.Col, modes)
		if err := c.session.RequestPty("ansi", c.Row, c.Col, modes); err != nil {
			return fmt.Errorf("request for pseudo terminal failed: %v", err)
		}
	}

	c.closers = append(c.closers, func() error {
		if err := c.session.Close(); err != nil && err != io.EOF {
			return fmt.Errorf("closing session: %v", err)
		}
		return nil
	})

	// The rules for the environment follow those of os/exec:
	// if c.Env is nil, os.Environ is used.
	if c.Env == nil {
		c.Env = os.Environ()
	}

	if err := c.SetEnv(c.Env...); err != nil {
		return err
	}

	// if they did not set an attacher, provide a default one
	if c.fileServer == nil {
		c.fileServer = &CPU9P{path: c.Root}
	}

	if c.SessionIn, err = c.session.StdinPipe(); err != nil {
		return err
	}
	c.closers = append([]func() error{func() error {
		c.SessionIn.Close()
		return nil
	}}, c.closers...)

	if c.SessionOut, err = c.session.StdoutPipe(); err != nil {
		return err
	}
	if c.SessionErr, err = c.session.StderrPipe(); err != nil {
		return err
	}

	// Unlike the cpu command source, which assumes an SSH-like stdin,
	// but very much like es/exec, users of Stdin in this package
	// will need to set the IO.
	// e.g.,
	// go c.SSHStdin(i, c.Stdin)
	// N.B.: if a 9p server was needed, it was started in Dial.

	cmd := c.cmd
	if c.port9p != 0 {
		cmd += fmt.Sprintf("-port9p=%v", c.port9p)
	}
	// The ABI for ssh.Start uses a string, not a []string
	// On the other end, it splits the string back up
	// as needed, claiming to do proper unquote handling.
	// This means we have to take care about quotes on
	// our side.
	quotedArgs := make([]string, len(c.Args))
	for i, arg := range c.Args {
		quotedArgs[i] = quoteArg(arg)
	}
	cmd += " " + strings.Join(quotedArgs, " ")

	verbose("call session.Start(%s)", cmd)
	if err := c.session.Start(cmd); err != nil {
		return fmt.Errorf("Failed to run %v: %v", c, err.Error())
	}
	if c.hasTTY {
		verbose("Setup interactive input")
		if err := c.SetupInteractive(); err != nil {
			return err
		}
		go c.TTYIn(c.session, c.SessionIn, c.Stdin)
	} else {
		verbose("Setup batch input")
		go func() {
			if _, err := io.Copy(c.SessionIn, c.Stdin); err != nil && !errors.Is(err, io.EOF) {
				log.Printf("copying stdin: %v", err)
			}
			if err := c.SessionIn.Close(); err != nil {
				log.Printf("Closing stdin: %v", err)
			}
		}()
	}
	go func() {
		verbose("set up copying to c.Stdout")
		if _, err := io.Copy(c.Stdout, c.SessionOut); err != nil && !errors.Is(err, io.EOF) {
			log.Printf("copying stdout: %v", err)
		}

		// If the file is NOT stdout, close it.
		// This is needed when programmers have
		// set c.Stdout to be some other WriteCloser, e.g. a pipe.
		if !sameFD(c.Stdout, os.Stdout) {
			c.Stdout.Close()
		}
	}()
	go func() {
		verbose("set up copying to c.Stderr")
		if _, err := io.Copy(c.Stderr, c.SessionErr); err != nil && !errors.Is(err, io.EOF) {
			log.Printf("copying stderr: %v", err)
		}
		if !sameFD(c.Stdout, os.Stderr) {
			c.Stderr.Close()
		}
	}()

	return nil
}

// Wait waits for a Cmd to finish.
func (c *Cmd) Wait() error {
	err := c.session.Wait()
	return err
}

// Run runs a command with Start, and waits for it to finish with Wait.
func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *Cmd) CombinedOutput() ([]byte, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	c.Stdout, c.Stderr = w, w

	cpuerr := c.Run()

	b, err := io.ReadAll(r)

	return b, errors.Join(cpuerr, err)
}

// TTYIn manages tty input for a cpu session.
// It exists mainly to deal with ~.
func (c *Cmd) TTYIn(s *ssh.Session, w io.WriteCloser, r io.Reader) {
	var newLine, tilde bool
	var t = []byte{'~'}
	var b [1]byte
	for {
		if _, err := r.Read(b[:]); err != nil {
			return
		}
		switch b[0] {
		default:
			newLine = false
			if tilde {
				if _, err := w.Write(t[:]); err != nil {
					return
				}
				tilde = false
			}
			if _, err := w.Write(b[:]); err != nil {
				return
			}
		case '\n', '\r':
			newLine = true
			if _, err := w.Write(b[:]); err != nil {
				return
			}
		case '~':
			if newLine {
				newLine = false
				tilde = true
				break
			}
			if _, err := w.Write(t[:]); err != nil {
				return
			}
		case '.':
			if tilde {
				s.Close()
				return
			}
			if _, err := w.Write(b[:]); err != nil {
				return
			}
		}
	}
}

// SetupInteractive sets up a cpu client for interactive access.
// It adds a function to c.Closers to clean up the terminal.
func (c *Cmd) SetupInteractive() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	c.closers = append(c.closers, func() error {
		term.Restore(int(os.Stdin.Fd()), oldState)
		return nil
	})

	return nil
}

// Close ends a cpu session, doing whatever is needed.
func (c *Cmd) Close() error {
	var err error
	for _, f := range c.closers {
		if e := f(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
