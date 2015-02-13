package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"pty"
)

type cgroupname string

func (c cgroupname) apply(s string, f func(s string)) {
	// range of strings.Split("",",") is 1.
	// not exactly what we might expect.
	if s == "" {
		return
	}
	for _, g := range strings.Split(s, ",") {
		p := path.Join(g)
		f(p)
	}
}

func (c cgroupname) Validate(s string) {
	c.apply(s, func(s string) {
		if st, err := os.Stat(path.Join(string(c), s)); err != nil {
			log.Fatal("%v", err)
		} else if !st.IsDir() {
			log.Fatal("%s: not a directory", s)
		}
	})
}

func (c cgroupname) Create(s, name string) {
	if err := os.MkdirAll(path.Join(string(c), s, name), 0755); err != nil {
		log.Fatal(err)
	}
}

func (c cgroupname) Attach(s, name string, pid int) {
	t := path.Join(string(c), s, name, "tasks")
	b := []byte(fmt.Sprintf("%v", pid))
	if err := ioutil.WriteFile(t, b, 0600); err != nil {
		log.Fatal(err)
	}
}

func (c cgroupname) Destroy(s, n string) {
	if err := os.RemoveAll(path.Join(string(c), s, n)); err != nil {
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
}

// Add adds a mount to the global mountlist. Don't know if we need it, but we might for additional volumes?
func Add(src, dst, mtype, opts string, flags uintptr, dir bool) {
	mounts = append(mounts, mount{src: src, dst: dst, mtype: mtype, flags: flags, opts: opts, dir: dir})

}

// One mounts one mountpoint, using base as a prefix for the destination.
// If anything goes wrong, we just bail out; we've privatized the namespace
// so there is no cleanup we need to do.
func (m *mount) One(base string) {
	dst := path.Join(base, m.dst)
	if m.dir {
		if err := os.MkdirAll(dst, 0755); err != nil {
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
func MountAll(base string) {
	root.One("")
	for _, m := range mounts {
		m.One(base)
	}
}

// make_console sets the right modes for the real console, then creates
// a /dev/console in the chroot.
func make_console(base, console string) {
	if err := os.Chmod(console, 0600); err != nil {
		log.Printf("%v", err)
	}
	if err := os.Chown(console, 0, 0); err != nil {
		log.Printf("%v", err)
	}

	st, err := os.Stat(console)
	if err != nil {
		log.Printf("%v", err)
	}

	nn := path.Join(base, "/dev/console")
	if err := syscall.Mknod(nn, uint32(st.Mode()), int(st.Sys().(*syscall.Stat_t).Dev)); err != nil {
		log.Printf("%v", err)
	}

	// if any previous steps failed, this one will too, so we can bail here.
	if err := syscall.Mount(console, nn, "", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Mount :%s: on :%s: flags %v: %v",
			console, nn, syscall.MS_BIND, err)
	}

}

// copy_nodes makes copies of needed nodes in the chroot.
func copy_nodes(base string) {
	nodes := []string{
		"/dev/tty",
		"/dev/full",
		"/dev/null",
		"/dev/zero",
		"/dev/random",
		"/dev/urandom"}

	for _, n := range nodes {
		st, err := os.Stat(n)
		if err != nil {
			log.Printf("%v", err)
		}
		nn := path.Join(base, n)
		if err := syscall.Mknod(nn, uint32(st.Mode()), int(st.Sys().(*syscall.Stat_t).Dev)); err != nil {
			log.Printf("%v", err)
		}

	}
}

// make_ptmx creates /dev/ptmx in the root. Because of order of operations
// it has to happen at a different time than copy_nodes.
func make_ptmx(base string) {
	dst := path.Join(base, "/dev/ptmx")

	if err := os.Symlink("/dev/ptmx", dst); err != nil {
		log.Printf("%v", err)
	}
}

// make_symlinks sets up standard symlinks as found in /dev.
func make_symlinks(base string) {
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
		dst := path.Join(base, linkit[i].dst)

		if err := os.Symlink(linkit[i].src, dst); err != nil {
			log.Printf("%v", err)
		}
	}
}

var (
	cgpath  = flag.String("cgpath", "/sys/fs/cgroup", "set the cgroups")
	cgroup  = flag.String("cgroup", "", "set the cgroups")
	mnt     = flag.String("mount", "", "define mounts")
	chroot  = flag.String("chroot", "", "where to chroot to")
	chdir   = flag.String("chdir", "/", "where to chrdir to in the chroot")
	console = flag.String("console", "/dev/console", "where the root is")
	keepenv = flag.Bool("keepenv", false, "Keep the environment")
	env     = flag.String("env", "", "other environment variables")
	user    = flag.String("user", "root" /*user.User.Username*/, "User name")
	root    = &mount{"", "/", "", "", syscall.MS_SLAVE | syscall.MS_REC, false}
	mounts  = []mount{
		{"proc", "/proc", "proc", "", syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV, true},
		{"/proc/sys", "/proc/sys", "", "", syscall.MS_BIND, true},
		{"", "/proc/sys", "", "", syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REMOUNT, true},
		{"sysfs", "/sys", "sysfs", "", syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_RDONLY, true},
		{"tmpfs", "/dev", "tmpfs", "mode=755", syscall.MS_NOSUID | syscall.MS_STRICTATIME, true},
		{"devpts", "/dev/pts", "devpts", "newinstance,ptmxmode=000,mode=620,gid=5", syscall.MS_NOSUID | syscall.MS_NOEXEC, true},
		{"tmpfs", "/dev/shm", "tmpfs", "mode=1777", syscall.MS_NOSUID | syscall.MS_STRICTATIME | syscall.MS_NODEV, true},
		{"tmpfs", "/run", "tmpfs", "mode=755", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_STRICTATIME, true},
	}
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
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
	if a[len(a)-1] != "#" {

		a = append(a, "#")
		// spawn ourselves with the right unsharing settings.
		c := exec.Command(a[0], a[1:]...)
		c.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID}
		//		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET

		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		t, err := getTermios(1)
		if err != nil {
			log.Fatalf("Can't get termios on fd 1: %v", err)
		}
		if err := c.Run(); err != nil {
			log.Printf(err.Error())
		}
		if err := t.set(1); err != nil {
			log.Printf("Can't reset termios on fd1: %v", err)
		}
		os.Exit(1)
	}

	// unlike pflask, we require that you set a chroot.
	// If you make it /, strange things are bound to happen.
	// if that is too limiting we'll have to change this.
	if *chroot == "" {
		log.Fatalf("you are required to set the chroot via --chroot")
	}

	a = flag.Args()
	//log.Printf("greetings %v\n", a)
	a = a[:len(a)-1]

	ptm, pts, sname, err := pty.Open()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// child code. Not really. What really happens here is we set
	// ourselves into the container, and spawn the child. It's a bit odd
	// but we're the master, but we'll run in the container? I don't know
	// how else to do it. This may require we set some things up first,
	// esp. the network. But, it's all fun and games until someone loses
	// an eye.
	MountAll(*chroot)

	copy_nodes(*chroot)

	make_ptmx(*chroot)

	make_symlinks(*chroot)

	make_console(*chroot, sname)

	//umask(0022);

	/* TODO: drop capabilities */

	//do_user(user);

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
	c.Stdout = pts
	c.Stdin = pts
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

	for {
		if err = syscall.PtraceDetach(kid); err != nil {
			log.Printf("Could not detach %v, sleeping five seconds", kid)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	raw()

	go func() {
		io.Copy(os.Stdout, ptm)
		os.Exit(1)
	}()
	io.Copy(ptm, os.Stdin)
}
