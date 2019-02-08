// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"
)

var (
	debug       = flag.Bool("d", false, "enable debug prints")
	v           = func(string, ...interface{}) {}
	remote      = flag.Bool("remote", false, "Indicates we are the remote side of the cpu session")
	network     = flag.String("network", "tcp", "network to use")
	host        = flag.String("h", "localhost", "host to use")
	port        = flag.String("p", "22", "port to use")
	keyFile     = flag.String("key", filepath.Join(os.Getenv("HOME"), ".ssh/cpu_rsa"), "key file")
	hostKeyFile = flag.String("hostkey", "" /*filepath.Join(os.Getenv("HOME"), ".ssh/hostkeyfile"), */, "host key file")
	srv9p       = flag.String("srv", "unpfs", "what server to run")
	bin         = flag.String("bin", "cpu", "path of cpu binary")
	port9p      = flag.String("port9p", "", "port9p # on remote machine for 9p mount")
	dbg9p       = flag.Bool("dbg9p", false, "show 9p io")
)

func verbose(f string, a ...interface{}) {
	v(f+"\r\n", a...)
}

func dial(n, a string, config *ssh.ClientConfig) (*ssh.Client, error) {
	client, err := ssh.Dial(n, a, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %v", err)
	}
	return client, nil
}

func config(kf string) (*ssh.ClientConfig, error) {
	cb := ssh.InsecureIgnoreHostKey()
	//var hostKey ssh.PublicKey
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(kf)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key %v: %v", kf, err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("ParsePrivateKey %v: %v", kf, err)
	}
	if *hostKeyFile != "" {
		return nil, fmt.Errorf("no support for hostkeyfile arg yet")
		//cb = ssh.FixedHostKey(hostKeyFile)
	}
	config := &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: cb,
	}
	return config, nil
}

func cmd(client *ssh.Client, s string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(s); err != nil {
		return nil, fmt.Errorf("Failed to run %v: %v", s, err.Error())
	}
	return b.Bytes(), nil
}

func dropPrivs() error {
	uid := unix.Getuid()
	v("dropPrives: uid is %v", uid)
	if uid == 0 {
		v("dropPrivs: not dropping privs")
		return nil
	}
	gid := unix.Getgid()
	v("dropPrivs: gid is %v", gid)
	if err := unix.Setreuid(-1, uid); err != nil {
		return err
	}
	return unix.Setregid(-1, gid)
}

// start up a namespace. We must
// mkdir /tmp/cpu on the remote machine
// issue the mount command
// test via an ls of /tmp/cpu
// TODO: unshare first
// We enter here are uid 0 and once the mount is done, back down.
func runRemote(cmd, port9p string) error {
	// for some reason echo is not set.
	t, err := termios.New()
	if err != nil {
		log.Printf("can't get a termios; oh well; %v", err)
	} else {
		term, err := t.Get()
		if err != nil {
			log.Printf("can't get a termios; oh well; %v", err)
		} else {
			term.Lflag |= unix.ECHO | unix.ECHONL
			if err := t.Set(term); err != nil {
				log.Printf("can't set a termios; oh well; %v", err)
			}
		}
	}

	// It's true we are making this directory while still root.
	// This ought to be safe as it is a private namespace mount.
	if err := os.Mkdir("/tmp/cpu", 0666); err != nil && !os.IsExist(err) {
		log.Println(err)
	}

	user := os.Getenv("USER")
	if user == "" {
		user = "nouser"
	}
	flags := uintptr(unix.MS_NODEV | unix.MS_NOSUID)
	opts := fmt.Sprintf("version=9p2000.L,trans=tcp,port=%v,uname=%v", port9p, user)
	if err := unix.Mount("127.0.0.1", "/tmp/cpu", "9p", flags, opts); err != nil {
		return fmt.Errorf("9p mount %v", err)
	}
	// We don't want to run as the wrong uid.
	if err := dropPrivs(); err != nil {
		return err
	}
	// The unmount happens for free since we unshared.
	v("runRemote: command is %q", cmd)
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}

// srv on 5641.
// TODO: make it more private, and also, have server only take
// one connection or use stdin/stdout
func srv(ctx context.Context) (net.Conn, *exec.Cmd, error) {
	c := exec.CommandContext(ctx, "unpfs", "tcp!localhost!5641", os.Getenv("HOME"))
	o, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	c.Stderr = c.Stdout
	if err := c.Start(); err != nil {
		return nil, nil, err
	}
	// Wait for the ready message.
	var b = make([]byte, 8192)
	n, err := o.Read(b)
	if err != nil {
		return nil, nil, err
	}
	v("Server says: %q", string(b[:n]))

	srvSock, err := net.Dial("tcp", "localhost:5641")
	if err != nil {
		return nil, nil, err
	}
	return srvSock, c, nil
}

// We only do one accept for now.
func forward(l net.Listener, s net.Conn) error {
	//if err := l.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
	//return fmt.Errorf("Can't set 9p client listen deadline: %v", err)
	//}
	c, err := l.Accept()
	v("forward: c %v err %v", c, err)
	if err != nil {
		v("forward: accept: %v", err)
		return err
	}
	go io.Copy(s, c)
	go io.Copy(c, s)
	return nil
}

// To make sure defer gets run and you tty is sane on exit
func runClient(a string) error {
	c, err := config(*keyFile)
	if err != nil {
		return err
	}
	cl, err := dial(*network, *host+":"+*port, c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	srvSock, p, err := srv(ctx)
	if err != nil {
		cancel()
		return err
	}
	defer func() {
		cancel()
		p.Wait()
	}()
	// Arrange port forwarding from remote ssh to our server.

	// Request the remote side to open port 5640 on all interfaces.
	// Note: cl.Listen returns a TCP listener with network is "tcp"
	// or variants. This lets us use a listen deadline.
	l, err := cl.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("First cl.Listen %v", err)
	}
	ap := strings.Split(l.Addr().String(), ":")
	if len(ap) == 0 {
		return fmt.Errorf("Can't find a port number in %v", l.Addr().String())
	}
	port := ap[len(ap)-1]
	v("listener %T %v addr %v port %v", l, l, l.Addr().String(), port)

	go forward(l, srvSock)
	v("Connected to %v", cl)

	// now run stuff.
	if err := shell(cl, a, port); err != nil {
		return err
	}
	return nil
}

// env sets environment variables. While we might think we ought to set
// HOME and PATH, it's possibly not a great idea. We leave them here as markers
// to remind ourselves not to try it later.
// We don't just grab all environment variables because complex bash functions
// will have no meaning to elvish. If there are simpler environment variables
// you want to set, add them here. Note however that even basic ones like TERM
// don't work either.
func env(s *ssh.Session) {
	e := []string{"HOME", "PATH"}
	// HOME and PATH are not allowed to be set by many sshds. Annoying.
	for _, v := range e[2:] {
		if err := s.Setenv(v, os.Getenv(v)); err != nil {
			log.Printf("Warning: s.Setenv(%q, %q): %v", v, os.Getenv(v), err)
		}
	}
}
func shell(client *ssh.Client, a, port9p string) error {
	t, err := termios.New()
	if err != nil {
		return err
	}
	r, err := t.Raw()
	if err != nil {
		return err
	}
	defer t.Set(r)
	if *bin == "" {
		if *bin, err = exec.LookPath("cpu"); err != nil {
			return err
		}
	}
	a = fmt.Sprintf("%v -remote -port9p %v -bin %v %v", *bin, port9p, *bin, a)
	v("command is %q", a)
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	env(session)
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("ansi", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	i, err := session.StdinPipe()
	if err != nil {
		return err
	}
	o, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	e, err := session.StderrPipe()
	if err != nil {
		return err
	}

	// sshd doesn't want to set us set the HOME and PATH via the normal
	// request route. So we do this nasty hack to ensure we can find
	// the cpu binary. We append our paths to the one the shell has.
	// This should suffice for u-root systems with paths including
	// /bbin and /ubin as well as more conventional systems.
	// The only possible flaw in this approach is elvish, which
	// has a very odd PATH syntax. For elvish, the PATH= is ignored,
	// so does no harm. Our use case for elvish is u-root, and
	// we will have the right path anyway, so it will still work.
	// It is working well in testing.
	cmd := fmt.Sprintf("PATH=$PATH:%s %s", os.Getenv("PATH"), a)
	v("Start remote with command %q", cmd)
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("Failed to run %v: %v", a, err.Error())
	}
	go io.Copy(i, os.Stdin)
	go io.Copy(os.Stdout, o)
	go io.Copy(os.Stderr, e)
	return session.Wait()
}

// We do flag parsing in init so we can
// Unshare if needed while we are still
// single threaded.
func init() {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	if *remote {
		// The unshare system call in Linux doesn't unshare mount points
		// mounted with --shared. Systemd mounts / with --shared. For a
		// long discussion of the pros and cons of this see debian bug 739593.
		// The Go model of unsharing is more like Plan 9, where you ask
		// to unshare and the namespaces are unconditionally unshared.
		// To make this model work we must further mark / as MS_PRIVATE.
		// This is what the standard unshare command does.
		var (
			none  = [...]byte{'n', 'o', 'n', 'e', 0}
			slash = [...]byte{'/', 0}
			flags = uintptr(unix.MS_PRIVATE | unix.MS_REC) // Thanks for nothing Linux.
		)
		if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
			log.Printf("bad Unshare: %v", err)
		}
		_, _, err1 := syscall.RawSyscall6(unix.SYS_MOUNT, uintptr(unsafe.Pointer(&none[0])), uintptr(unsafe.Pointer(&slash[0])), 0, flags, 0, 0)
		if err1 != 0 {
			log.Printf("Warning: unshare failed (%v). There will be no private 9p mount", err1)
		}
		flags = 0
		if err := unix.Mount("cpu", "/tmp", "tmpfs", flags, ""); err != nil {
			log.Printf("Warning: tmpfs mount on /tmp (%v) failed. There will be no 9p mount", err)
		}
	}
}
func main() {
	a := strings.Join(flag.Args(), " ")
	if *remote {
		if err := runRemote(a, *port9p); err != nil {
			log.Fatal(err)
		}
		return
	}
	if a == "" {
		a = os.Getenv("SHELL")
	}
	if err := runClient(a); err != nil {
		log.Fatal(err)
	}

}
