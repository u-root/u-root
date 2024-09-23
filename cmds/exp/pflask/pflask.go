// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// pty support. We used to import github.com/kr/pty but what we need is not that complex.
// Thanks to keith rarick for these functions.

func ptsopen() (controlPTY, processTTY *os.File, ttyname string, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return
	}

	ttyname, err = ptsname(p)
	if err != nil {
		return
	}

	err = ptsunlock(p)
	if err != nil {
		return
	}

	v("OpenFile %v %x\n", ttyname, os.O_RDWR|syscall.O_NOCTTY)
	t, err := os.OpenFile(ttyname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return
	}
	return p, t, ttyname, nil
}

func ptsname(f *os.File) (string, error) {
	n, err := unix.IoctlGetInt(int(f.Fd()), unix.TIOCGPTN)
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(n), nil
}

func ptsunlock(f *os.File) error {
	var u uintptr
	// use TIOCSPTLCK with a zero valued arg to clear the pty lock
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&u)))
	if err != 0 {
		return err
	}
	return nil
}

type cgroupname string

func (c cgroupname) apply(s string, f func(s string)) {
	// range of strings.Split("",",") is 1.
	// not exactly what we might expect.
	if s == "" {
		return
	}
	for _, g := range strings.Split(s, ",") {
		p := filepath.Join(g)
		f(p)
	}
}

func (c cgroupname) Validate(s string) {
	c.apply(s, func(s string) {
		if st, err := os.Stat(filepath.Join(string(c), s)); err != nil {
			log.Fatalf("%v", err)
		} else if !st.IsDir() {
			log.Fatalf("%s: not a directory", s)
		}
	})
}

func (c cgroupname) Create(s, name string) {
	if err := os.MkdirAll(filepath.Join(string(c), s, name), 0o755); err != nil {
		log.Fatal(err)
	}
}

func (c cgroupname) Attach(s, name string, pid int) {
	t := filepath.Join(string(c), s, name, "tasks")
	b := []byte(fmt.Sprintf("%v", pid))
	if err := os.WriteFile(t, b, 0o600); err != nil {
		log.Fatal(err)
	}
}

func (c cgroupname) Destroy(s, n string) {
	if err := os.RemoveAll(filepath.Join(string(c), s, n)); err != nil {
		log.Fatal(err)
	}
}

func (c cgroupname) Do(groups string, pid int) {
	cgn := fmt.Sprintf("pflask.%d", pid)
	c.apply(groups, func(s string) {
		c.Create(s, cgn)
		c.Attach(s, cgn, pid)
	})
}

type mount struct {
	src, dst, mtype, opts string
	flags                 uintptr
	dir                   bool
	needPrivilege         bool
}

// Add adds a mount to the global mountlist. Don't know if we need it, but we might for additional volumes?
func Add(src, dst, mtype, opts string, flags uintptr, dir bool) {
	mounts = append(mounts, mount{src: src, dst: dst, mtype: mtype, flags: flags, opts: opts, dir: dir})
}

// One mounts one mountpoint, using base as a prefix for the destination.
// If anything goes wrong, we just bail out; we've privatized the namespace
// so there is no cleanup we need to do.
func (m *mount) One(base string) {
	dst := filepath.Join(base, m.dst)
	if m.dir {
		if err := os.MkdirAll(dst, 0o755); err != nil {
			log.Fatalf("One: mkdirall %v: %v", m.dst, err)
		}
	}
	if err := syscall.Mount(m.src, dst, m.mtype, m.flags, m.opts); err != nil {
		log.Fatalf("Mount :%s: on :%s: type :%s: flags %x: opts :%v: %v\n",
			m.src, m.dst, m.mtype, m.flags, m.opts, err)
	}
}

// MountAll mounts all the mount points. root is a bit special in that it just sets
// needed flags for non-shared mounts.
func MountAll(base string, unprivileged bool) {
	root.One("")
	for _, m := range mounts {
		if m.needPrivilege && unprivileged {
			continue
		}
		m.One(base)
	}
}

// modedev returns a mode and dev suitable for use in mknod.
// It's very odd, but the Dev either needs to be byteswapped
// or comes back byteswapped. I just love it that the world
// has fixed on a 45-year-old ABI (stat in this case)
// that was abandoned by its designers 30 years ago.
// Oh well.
func modedev(st os.FileInfo) (uint32, int) {
	// Weird. The Dev is byte-swapped for some reason.
	dev := int(st.Sys().(*syscall.Stat_t).Dev)
	devlo := dev & 0xff
	dev >>= 8
	dev |= (devlo << 8)
	return uint32(st.Sys().(*syscall.Stat_t).Mode), dev
}

// makeConsole sets the right modes for the real console, then creates
// a /dev/console in the chroot.
func makeConsole(base, console string, unprivileged bool) {
	if err := os.Chmod(console, 0o600); err != nil {
		log.Printf("%v", err)
	}
	if err := os.Chown(console, 0, 0); err != nil {
		log.Printf("%v", err)
	}

	st, err := os.Stat(console)
	if err != nil {
		log.Printf("%v", err)
	}

	nn := filepath.Join(base, "/dev/console")
	mode, dev := modedev(st)
	if unprivileged {
		// In unprivileged uses, we can't mknod /dev/console, however,
		// we can just create a file /dev/console and use bind mount on file.
		if _, err := os.Stat(nn); err != nil {
			os.WriteFile(nn, []byte{}, 0o600) // best effort, ignore error
		}
	} else {
		if err := syscall.Mknod(nn, mode, dev); err != nil {
			log.Printf("%v", err)
		}
	}

	// if any previous steps failed, this one will too, so we can bail here.
	if err := syscall.Mount(console, nn, "", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Mount :%s: on :%s: flags %v: %v",
			console, nn, syscall.MS_BIND, err)
	}
}

// copyNodes makes copies of needed nodes in the chroot.
func copyNodes(base string) {
	nodes := []string{
		"/dev/tty",
		"/dev/full",
		"/dev/null",
		"/dev/zero",
		"/dev/random",
		"/dev/urandom",
	}

	for _, n := range nodes {
		st, err := os.Stat(n)
		if err != nil {
			log.Printf("%v", err)
		}
		nn := filepath.Join(base, n)
		mode, dev := modedev(st)
		if err := syscall.Mknod(nn, mode, dev); err != nil {
			log.Printf("%v", err)
		}
	}
}

// makePtmx creates /dev/ptmx in the root. Because of order of operations
// it has to happen at a different time than copyNodes.
func makePtmx(base string) {
	dst := filepath.Join(base, "/dev/ptmx")

	if _, err := os.Stat(dst); err == nil {
		return
	}

	if err := os.Symlink("/dev/pts/ptmx", dst); err != nil {
		log.Printf("%v", err)
	}
}

// makeSymlinks sets up standard symlinks as found in /dev.
func makeSymlinks(base string) {
	linkit := []struct {
		src, dst string
	}{
		{"/dev/pts/ptmx", "/dev/ptmx"},
		{"/proc/kcore", "/dev/core"},
		{"/proc/self/fd", "/dev/fd"},
		{"/proc/self/fd/0", "/dev/stdin"},
		{"/proc/self/fd/1", "/dev/stdout"},
		{"/proc/self/fd/2", "/dev/stderr"},
	}

	for i := range linkit {
		dst := filepath.Join(base, linkit[i].dst)

		if _, err := os.Stat(dst); err == nil {
			continue
		}

		if err := os.Symlink(linkit[i].src, dst); err != nil {
			log.Printf("%v", err)
		}
	}
}

var (
	cgpath  = flag.String("cgpath", "/sys/fs/cgroup", "set the cgroups")
	cgroup  = flag.String("cgroup", "", "set the cgroups")
	chroot  = flag.String("chroot", "", "where to chroot to")
	chdir   = flag.String("chdir", "/", "where to chrdir to in the chroot")
	keepenv = flag.Bool("keepenv", false, "Keep the environment")
	debug   = flag.Bool("d", false, "Enable debug logs")
	env     = flag.String("env", "", "other environment variables")
	user    = flag.String("user", "root" /*user.User.Username*/, "User name")
	root    = &mount{"", "/", "", "", syscall.MS_SLAVE | syscall.MS_REC, false, false}
	mounts  = []mount{
		{"proc", "/proc", "proc", "", syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV, true, false},
		{"/proc/sys", "/proc/sys", "", "", syscall.MS_BIND, true, true},
		{"", "/proc/sys", "", "", syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REMOUNT, true, true},
		{"sysfs", "/sys", "sysfs", "", syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_RDONLY, true, true},
		{"tmpfs", "/dev", "tmpfs", "mode=755", syscall.MS_NOSUID | syscall.MS_STRICTATIME, true, true}, // unprivileged system needs a pre-populated /dev
		{"devpts", "/dev/pts", "devpts", "newinstance,ptmxmode=0660,mode=0620", syscall.MS_NOSUID | syscall.MS_NOEXEC, true, false},
		{"tmpfs", "/dev/shm", "tmpfs", "mode=1777", syscall.MS_NOSUID | syscall.MS_STRICTATIME | syscall.MS_NODEV, true, false},
		{"tmpfs", "/run", "tmpfs", "mode=755", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_STRICTATIME, true, false},
	}
	v = func(string, ...interface{}) {}
)

func main() {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	v("pflask: Let's go!")

	if len(flag.Args()) < 1 {
		v("pflask: no args given")
		os.Exit(1)
	}

	// note the unshare system call worketh not for Go.
	// So do it ourselves. We have to start ourselves up again,
	// after having spawned ourselves with lots of clone
	// flags sets. To know that we spawned ourselves we add '#'
	// as the last arg. # was chosen because shells normally filter
	// it out, so its presence as our last arg is highly indicative
	// that we really spawned us. Also, for testing, you can always
	// pass it by hand to see what the namespace looks like.
	a := os.Args
	if a[len(a)-1][0] != '#' {
		a = append(a, "#")
		euid := syscall.Geteuid()
		v("Running as user %v\n", euid)
		if euid != 0 {
			a[len(a)-1] = "#u"
		}
		if *debug {
			testc := exec.Command("/bbin/echo", "    ===== cmd test")
			testc.Stdout = os.Stdout
			testc.Run()
			testc = exec.Command("/bbin/ls", a[0])
			testc.Stdout = os.Stdout
			testc.SysProcAttr = &syscall.SysProcAttr{Cloneflags: 0}
			testc.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
			testc.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUTS
			testc.SysProcAttr.Cloneflags |= syscall.CLONE_NEWIPC
			testc.SysProcAttr.Cloneflags |= syscall.CLONE_NEWPID
			if err := testc.Run(); err != nil {
				log.Printf("Could not run:\n   %v\n    %v\n", testc, err.Error())
			}
		}
		// spawn ourselves with the right unsharing settings.
		c := exec.Command(a[0], a[1:]...)
		c.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID}
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET

		if euid != 0 {
			c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER
			c.SysProcAttr.UidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: syscall.Getuid(), Size: 1}}
			c.SysProcAttr.GidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: syscall.Getgid(), Size: 1}}
		}
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		v("pflask: respawning...")
		if err := c.Run(); err != nil {
			log.Printf("Could not run:\n   %v\n    %v\n", c, err.Error())
			if strings.Contains(err.Error(), "invalid argument") {
				log.Println("Ensure that your kernel is configured for CGROUPs and NS.")
				log.Println("The following are needed: IPC, PID, USER, UTS")
			}
			if strings.Contains(err.Error(), "device or resource busy") {
				log.Println("No clue...")
			}
		}
		//if err := termios.SetTermios(1, t); err != nil {
		//	log.Printf("Can't reset termios on fd1: %v", err)
		//}
		os.Exit(1)
	}

	unprivileged := a[len(a)-1] == "#u"

	// unlike the original pflask, we require that you set a chroot.
	// If you make it /, strange things are bound to happen.
	// if that is too limiting we'll have to change this.
	if *chroot == "" {
		log.Fatal("you are required to set the chroot via -chroot")
	}
	if *chroot == "/" {
		log.Println("[WARN]: chroot set to /: strange things are bound to happen")
	}

	a = flag.Args()
	v("greetings %v\n", a)
	a = a[:len(a)-1]

	v("pflask: ptsopen")
	controlPTY, processTTY, sname, err := ptsopen()
	if err != nil {
		log.Fatal(err)
	}

	// child code. Not really. What really happens here is we set
	// ourselves into the container, and spawn the child. It's a bit odd
	// but we're the parent, but we'll run in the container? I don't know
	// how else to do it. This may require we set some things up first,
	// esp. the network. But, it's all fun and games until someone loses
	// an eye.
	v("MountAll")
	MountAll(*chroot, unprivileged)

	if !unprivileged {
		v("copyNodes")
		copyNodes(*chroot)
	}

	v("makePtmx")
	makePtmx(*chroot)

	v("makeSymlinks")
	makeSymlinks(*chroot)

	v("makeConsole")
	makeConsole(*chroot, sname, unprivileged)

	// umask(0022);
	/* TODO: drop capabilities */
	// do_user(user);

	e := make(map[string]string)
	if *keepenv {
		for _, v := range os.Environ() {
			k := strings.SplitN(v, "=", 2)
			e[k[0]] = k[1]
		}
	}

	term := os.Getenv("TERM")
	e["TERM"] = term
	e["PATH"] = "/usr/sbin:/usr/bin:/sbin:/bin"
	e["USER"] = *user
	e["LOGNAME"] = *user
	e["HOME"] = "/root"

	if *env != "" {
		for _, c := range strings.Split(*env, ",") {
			k := strings.SplitN(c, "=", 2)
			if len(k) != 2 {
				log.Printf("Bogus environment string %v", c)
				continue
			}
			e[k[0]] = k[1]
		}
	}
	e["container"] = "pflask"

	if *cgroup == "" {
		var envs []string
		for k, v := range e {
			envs = append(envs, k+"="+v)
		}
		v("envs\n  %v\n", e)
		v("-- chroot --")
		if err := syscall.Chroot(*chroot); err != nil {
			log.Fatal(err)
		}
		v("--- chdir --")
		if err := syscall.Chdir(*chdir); err != nil {
			log.Fatal(err)
		}
		v("---- exec --")
		log.Fatal(syscall.Exec(a[0], a[1:], envs))
	}

	v("exec.Command")
	c := exec.Command(a[0], a[1:]...)
	c.Env = nil
	for k, v := range e {
		c.Env = append(c.Env, k+"="+v)
	}

	c.SysProcAttr = &syscall.SysProcAttr{
		Chroot:  *chroot,
		Setctty: true,
		Setsid:  true,
	}
	c.Stdout = processTTY
	c.Stdin = processTTY
	c.Stderr = c.Stdout
	c.SysProcAttr.Setctty = true
	c.SysProcAttr.Setsid = true
	c.SysProcAttr.Ptrace = true
	c.Dir = *chdir
	err = c.Start()
	if err != nil {
		panic(err)
	}
	kid := c.Process.Pid
	log.Printf("Started %d\n", kid)

	// set up the containers, then resume the process.
	// Its children will get the containers as it clones.

	cg := cgroupname(*cgpath)
	cg.Do(*cgroup, kid)

	// sometimes the detach fails. Looks like a race condition: we're
	// sending the detach before the child has hit the TRACE_ME point.
	// Experimentally, when it fails, even one seconds it too short to
	// sleep. Sleep for 5 seconds.
	// Oh well it's not that. It's that there is some one of these
	// processes not in the PID namespace of the child? Who knows, sigh.
	// This is an aspect of the Go runtime that is seriously broken.

	for i := 0; ; i++ {
		if err = syscall.PtraceDetach(kid); err != nil {
			log.Printf("Could not detach %v, sleeping 250 milliseconds", kid)
			time.Sleep(250 * time.Millisecond)
			continue
		}
		if i > 100 {
			log.Fatalf("Tried for 10 seconds to get a DETACH. Let's fix the go runtime someday")
		}
		break
	}

	raw()

	go func() {
		io.Copy(os.Stdout, controlPTY)
		os.Exit(1)
	}()
	io.Copy(controlPTY, os.Stdin)
}
