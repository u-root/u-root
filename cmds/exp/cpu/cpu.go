// Copyright 2018-2019 the u-root Authors. All rights reserved
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

	// We use this ssh because it implements port redirection.
	// It can not, however, unpack password-protected keys yet.
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty" // TODO: get rid of krpty
	"github.com/u-root/u-root/pkg/termios"
	"github.com/u-root/u-root/pkg/uroot/util"
	// We use this ssh because it can unpack password-protected private keys.
	ossh "golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"
)

var (
	// For the ssh server part
	hostKeyFile = flag.String("hk", "" /*"/etc/ssh/ssh_host_rsa_key"*/, "file for host key")
	pubKeyFile  = flag.String("pk", "key.pub", "file for public key")
	port        = flag.String("sp", "2222", "ssh default port")

	debug     = flag.Bool("d", false, "enable debug prints")
	runAsInit = flag.Bool("init", false, "run as init (Debug only; normal test is if we are pid 1")
	v         = func(string, ...interface{}) {}
	remote    = flag.Bool("remote", false, "indicates we are the remote side of the cpu session")
	network   = flag.String("network", "tcp", "network to use")
	host      = flag.String("h", "localhost", "host to use")
	keyFile   = flag.String("key", filepath.Join(os.Getenv("HOME"), ".ssh/cpu_rsa"), "key file")
	srv9p     = flag.String("srv", "unpfs", "what server to run")
	bin       = flag.String("bin", "cpu", "path of cpu binary")
	port9p    = flag.String("port9p", "", "port9p # on remote machine for 9p mount")
	dbg9p     = flag.Bool("dbg9p", false, "show 9p io")
	root      = flag.String("root", "/", "9p root")
	bindover  = flag.String("bindover", "/lib:/lib64:/lib32:/usr:/bin:/etc", ": separated list of directories in /tmp/cpu to bind over /")
)

func verbose(f string, a ...interface{}) {
	v(f+"\r\n", a...)
}

func dial(n, a string, config *ossh.ClientConfig) (*ossh.Client, error) {
	client, err := ossh.Dial(n, a, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %v", err)
	}
	return client, nil
}

func config(kf string) (*ossh.ClientConfig, error) {
	cb := ossh.InsecureIgnoreHostKey()
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
	signer, err := ossh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("ParsePrivateKey %v: %v", kf, err)
	}
	if *hostKeyFile != "" {
		hk, err := ioutil.ReadFile(*hostKeyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read host key %v: %v", *hostKeyFile, err)
		}
		pk, err := ossh.ParsePublicKey(hk)
		if err != nil {
			return nil, fmt.Errorf("host key %v: %v", string(hk), err)
		}
		cb = ossh.FixedHostKey(pk)
	}
	config := &ossh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ossh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ossh.PublicKeys(signer),
		},
		HostKeyCallback: cb,
	}
	return config, nil
}

func cmd(client *ossh.Client, s string) ([]byte, error) {
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
// We enter here as uid 0 and once the mount is done, back down.
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
	for _, n := range []string{"/tmp/cpu", "/tmp/local", "/tmp/merge", "/tmp/root"} {
		if err := os.Mkdir(n, 0666); err != nil && !os.IsExist(err) {
			log.Println(err)
		}
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

	// Further, bind / onto /tmp/local so a non-hacked-on version may be visible.
	if err := unix.Mount("/", "/tmp/local", "", syscall.MS_BIND, ""); err != nil {
		log.Printf("Warning: binding / over /tmp/cpu did not work: %v, continuing anyway", err)
	}

	var overlaid bool
	if util.FindFileSystem("overlay") == nil {
		if err := unix.Mount("overlay", "/tmp/root", "overlay", unix.MS_MGC_VAL, "lowerdir=/tmp/cpu,upperdir=/tmp/local,workdir=/tmp/merge"); err == nil {
			//overlaid = true
		} else {
			log.Printf("Overlayfs mount failed: %v. Proceeding with selective mounts from /tmp/cpu into /", err)
		}
	}
	if !overlaid {
		// We could not get an overlayfs mount.
		// There are lots of cases where binaries REQUIRE that ld.so be in the right place.
		// In some cases if you set LD_LIBRARY_PATH it is ignored.
		// This is disappointing to say the least. We just bind a few things into /
		// bind *may* hide local resources but for now it's the least worst option.
		dirs := strings.Split(*bindover, ":")
		for _, n := range dirs {
			t := filepath.Join("/tmp/cpu", n)
			if err := unix.Mount(t, n, "", syscall.MS_BIND, ""); err != nil {
				log.Printf("Warning: mounting %v on %v failed: %v", t, n, err)
			} else {
				log.Printf("Mounted %v on %v", t, n)
			}

		}
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
	c := exec.CommandContext(ctx, "unpfs", "tcp!localhost!5641", *root)
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
func env(s *ossh.Session) {
	e := []string{"HOME", "PATH", "LD_LIBRARY_PATH"}
	// HOME and PATH are not allowed to be set by many sshds. Annoying.
	for _, v := range e {
		if err := s.Setenv(v, os.Getenv(v)); err != nil {
			log.Printf("Warning: s.Setenv(%q, %q): %v", v, os.Getenv(v), err)
		}
	}
}

func shell(client *ossh.Client, a, port9p string) error {
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
	modes := ossh.TerminalModes{
		ossh.ECHO:          0,     // disable echoing
		ossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
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
	//	cmd := fmt.Sprintf("PATH=$PATH:%s %s", os.Getenv("PATH"), a)
	cmd := a
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
	if os.Getpid() == 1 {
		*runAsInit, *debug = true, true
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

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

// adjust adjusts environment variables containing paths from
// / to /tmp/cpu. On Plan 9 this is done with a union mount.
// PATH variables are the union mounts of Unix so we use them
// instead.
func adjust(env []string) []string {
	var res []string
	for _, e := range env {
		n := strings.SplitN(e, "=", 2)
		if len(n) < 2 {
			res = append(res, e)
			continue
		}
		v := strings.Split(n[1], ":")
		for i := range v {
			if filepath.IsAbs(v[i]) {
				v[i] = filepath.Join("/tmp/cpu", v[i])
			}
		}
		res = append(res, n[0]+"="+strings.Join(v, ":"))
	}
	return res
}

func handler(s ssh.Session) {
	a := s.Command()
	verbose("the handler is here, cmd is %v", a)
	cmd := exec.Command(a[0], a[1:]...)
	log.Printf("cmd.Env ius %v", cmd.Env)
	adj := adjust(s.Environ())
	log.Printf("s.Environt is %v, adjusted is %v", s.Environ(), adj)
	cmd.Env = append(cmd.Env, adj...)
	ptyReq, winCh, isPty := s.Pty()
	verbose("the command is %v", *cmd)
	if isPty {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
		f, err := pty.Start(cmd)
		verbose("command started with pty")
		if err != nil {
			log.Print(err)
			return
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
		verbose("running command without pty")
		if err := cmd.Run(); err != nil {
			log.Print(err)
			return
		}
	}
	verbose("handler exits")
}

func doInit() error {
	if err := cpuSetup(); err != nil {
		log.Printf("CPU setup error with cpu running as init: %v", err)
	}
	cmds := [][]string{{"/bin/defaultsh"}, {"/bbin/dhclient", "-verbose"}}
	verbose("Try to run %v", cmds)

	for _, v := range cmds {
		verbose("Let's try to run %v", v)
		if _, err := os.Stat(v[0]); os.IsNotExist(err) {
			verbose("it's not there")
			continue
		}

		// I *love* special cases. Evaluate just the top-most symlink.
		//
		// In source mode, this would be a symlink like
		// /buildbin/defaultsh -> /buildbin/elvish ->
		// /buildbin/installcommand.
		//
		// To actually get the command to build, argv[0] has to end
		// with /elvish, so we resolve one level of symlink.
		if filepath.Base(v[0]) == "defaultsh" {
			s, err := os.Readlink(v[0])
			if err == nil {
				v[0] = s
			}
			verbose("readlink of %v returns %v", v[0], s)
			// and, well, it might be a relative link.
			// We must go deeper.
			d, b := filepath.Split(v[0])
			d = filepath.Base(d)
			v[0] = filepath.Join("/", os.Getenv("UROOT_ROOT"), d, b)
			verbose("is now %v", v[0])
		}

		cmd := exec.Command(v[0], v[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
		verbose("Run %v", cmd)
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", v, err)
			continue
		}
	}
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

	// Now we run as an ssh server, and each time we get a connection,
	// we run that command after setting things up for it.
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

	// start the process reaper
	procs := make(chan int)
	go cpuDone(procs)

	server.SetOption(ssh.HostKeyFile(*hostKeyFile))
	log.Println("starting ssh server on port " + *port)
	if err := server.ListenAndServe(); err != nil {
		log.Print(err)
	}
	verbose("server.ListenAndServer returned")

	numprocs := <-procs
	verbose("Reaped %d procs", numprocs)
	return nil
}

func main() {
	verbose("Args %v pid %d *runasinit %v *remote %v", os.Args, os.Getpid(), *runAsInit, *remote)
	a := strings.Join(flag.Args(), " ")
	switch {
	case *runAsInit:
		verbose("Running as Init")
		if err := doInit(); err != nil {
			log.Fatal(err)
		}
	case *remote:
		verbose("Running as remote")
		if err := runRemote(a, *port9p); err != nil {
			log.Fatal(err)
		}
	default:
		verbose("Running as client")
		if a == "" {
			a = os.Getenv("SHELL")
		}
		if err := runClient(a); err != nil {
			log.Fatal(err)
		}
	}
}
