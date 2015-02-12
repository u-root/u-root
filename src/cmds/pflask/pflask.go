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
	//"user"

	"pty"
)

type cgroup string

func (c cgroup) apply(s string, f func(s string)) {
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

func (c cgroup) Validate(s string) {
	c.apply(s, func(s string) {
		if st, err := os.Stat(path.Join(string(c), s)); err != nil {
			log.Fatal("%v", err)
		} else if !st.IsDir() {
			log.Fatal("%s: not a directory", s)
		}
	})
}

func (c cgroup) Create(s, name string) {
	if err := os.MkdirAll(path.Join(string(c), s, name), 0755); err != nil {
		log.Fatal(err)
	}
}

func (c cgroup) Attach(s, name string, pid int) {
	t := path.Join(string(c), s, name, "tasks")
	b := []byte(fmt.Sprintf("%v", pid))
	if err := ioutil.WriteFile(t, b, 0600); err != nil {
		log.Fatal(err)
	}
}

func (c cgroup) Destroy(s, n string) {
	if err := os.RemoveAll(path.Join(string(c), s, n)); err != nil {
		log.Fatal(err)
	}
}

func (c cgroup) Do(groups string, pid int) {
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
	mounted               bool
}

type mlist struct {
	mounts []*mount
}

func NewMlist(base string) (*mlist, error) {
	m := &mlist{}
	if err := syscall.Mount("", "/", "", syscall.MS_SLAVE|syscall.MS_REC, ""); err != nil {
		err := fmt.Errorf("Mount :%s: on :%s: type :%s: flags %x: opts :%v: %v\n",
			"", "/", "", syscall.MS_SLAVE|syscall.MS_REC, "", err)
		return nil, err
	}

	return m, nil
}

func (m *mlist) Add(src, dst, mtype, opts string, flags uintptr, dir bool) {
	m.mounts = append(m.mounts, &mount{src: src, dst: dst, mtype: mtype, flags: flags, opts: opts, dir: dir})

}

func (m *mount) One(base string) error {
	dst := path.Join(base, m.dst)
	if m.dir {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("One: mkdirall %v: %v", m.dst, err)
		}
	}
	if err := syscall.Mount(m.src, dst, m.mtype, m.flags, m.opts); err != nil {
		return fmt.Errorf("Mount :%s: on :%s: type :%s: flags %x: opts :%v: %v\n",
			m.src, m.dst, m.mtype, m.flags, m.opts, err)
	}
	m.mounted = true
	return nil
}
func (m *mlist) Do(base string) {
	ok := true
	if base != "" {
		m.Add("proc", "/proc", "proc", "",
			syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV, true)

		m.Add("/proc/sys", "/proc/sys", "", "",
			syscall.MS_BIND, true)

		m.Add("", "/proc/sys", "", "",
			syscall.MS_BIND|syscall.MS_RDONLY|syscall.MS_REMOUNT, true)

		m.Add("sysfs", "/sys", "sysfs", "",
			syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV|syscall.MS_RDONLY, true)

		m.Add("tmpfs", "/dev", "tmpfs", "mode=755",
			syscall.MS_NOSUID|syscall.MS_STRICTATIME, true)

		m.Add("devpts", "/dev/pts", "devpts", "newinstance,ptmxmode=000,mode=620,gid=5",
			syscall.MS_NOSUID|syscall.MS_NOEXEC, true)

		m.Add("tmpfs", "/dev/shm", "tmpfs", "mode=1777",
			syscall.MS_NOSUID|syscall.MS_STRICTATIME|syscall.MS_NODEV, true)

		m.Add("tmpfs", "/run", "tmpfs", "mode=755",
			syscall.MS_NOSUID|syscall.MS_NODEV|syscall.MS_STRICTATIME, true)

	}

	for _, m := range m.mounts {
		err := m.One(base)
		if err != nil {
			log.Printf(err.Error())
			ok = false
		}
	}
	if !ok {
		m.Undo(base)
		log.Fatal("Not all mounts succeeded.")
	}
}

func (m *mlist) Undo(base string) {
	for i := range m.mounts {
		m := m.mounts[len(m.mounts)-i-1]
		if !m.mounted {
			continue
		}
		dst := path.Join(base, m.dst)
		if err := syscall.Unmount(dst, 0); err != nil {
			log.Printf("Unmounting %v: %v", m, err)
		}
		m.mounted = false
	}
}

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

func make_ptmx(base string) {
	dst := path.Join(base, "/dev/ptmx")
	
	if err := os.Symlink("/dev/ptmx", dst); err != nil {
		log.Printf("%v", err)
	}
}

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

func do_chroot(chroot, chdir string) {
	if chdir == "" {
		chdir = "/"
	}
	if err := syscall.Chroot(chroot); err != nil {
		log.Fatalf("%v", err)
	}
	if err := os.Chdir(chdir); err != nil {
		log.Fatalf("%v", err)
	}
}

var (
	cgpath  = flag.String("cgpath", "/sys/fs/cgroup", "set the cgroups")
	cg      = flag.String("cgroup", "", "set the cgroups")
	mnt     = flag.String("mount", "", "define mounts")
	chroot  = flag.String("chroot", "", "where to chroot to")
	chdir   = flag.String("chdir", "", "where to chrdir to in the chroot")
	console = flag.String("console", "/dev/console", "where the root is")
	keepenv = flag.Bool("keepenv", false, "Keep the environment")
	env     = flag.String("env", "", "other environment variables")
	user    = flag.String("user", "root" /*user.User.Username*/, "User name")
)

func r() {
    var (
        s *winsize
        t *termios
        err error
    )

    defer func() {
        if err != nil { fmt.Println(err) }
//        if err = defaultTermios.set(); err != nil { fmt.Println(err) }
        fmt.Print("\n")
	os.Exit(1)
    }()

    if s, err = getWinsize(1); err != nil { return }
    fmt.Printf("Window Size:\n\tLines: %d\n\tColumns: %d\n", s.ws_row, s.ws_col)

    fmt.Println("Entering Raw mode. . .")
    fmt.Println("Type some characters!  Press 'q' to quit!")

    if t, err = getTermios(1); err != nil { return } ;fmt.Printf("%v\n", t)
    if err = t.setRaw(1); err != nil { return } ;fmt.Printf("%v\n", t)

    var buff [4]byte
    for buff[0] != 'q' {
        if n, err := os.Stdin.Read(buff[:]); err != nil { return 
	}else {
			fmt.Printf(">>>>  %c\n", buff[:n])
	}
    }
}

func main() {
	//r()
	// note the unshare system call worketh not for Go. You have to run
	// this under the unshare command. Good times.
	flag.Parse()
	if *chroot == "" {
		log.Fatalf("you are required to set the chroot via --chroot")
	}
	fmt.Printf("greetings\n")
	a := flag.Args()

	// Just create the container and run with it for now.
	c := cgroup(*cgpath)
	ppid := os.Getpid()
	m, err := NewMlist(*chroot)
	if err != nil {
		log.Fatalf("%v", err)
	}
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
	if true {
		c.Do(*cg, ppid)

		m.Do(*chroot)

		copy_nodes(*chroot)

		make_ptmx(*chroot);

		make_symlinks(*chroot)

		make_console(*chroot, sname)

		//do_chroot(*chroot, *chdir)

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

		if len(a) == 0 {
			a = []string{"/bin/bash", "bash"}
		}
		c := exec.Command(a[0]) // , a[1:]...)
		c.Env = nil
		for k, v := range e {
			c.Env = append(c.Env, k+"="+v)
		}

		c.SysProcAttr = &syscall.SysProcAttr{
				Chroot: *chroot, 
				Setctty: true, Setsid: true, Cloneflags:
				syscall.CLONE_NEWIPC |
				syscall.CLONE_NEWPID |
				syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWNS  |
					0}
//			SIGCHLD      |
//				syscall.CLONE_NEWUSER|


		
			c.Stdout = pts
			c.Stdin = pts
			c.Stderr = c.Stdout
		c.SysProcAttr.Setctty = true
		c.SysProcAttr.Setsid = true
		c.SysProcAttr.Ctty = 0
		t, err := getTermios(2)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if err = t.setRaw(1); err != nil {log.Fatalf(err.Error())}
		
		err = c.Start()
		if err != nil {
			panic(err)
		}
		go io.Copy(os.Stdout, ptm)
		io.Copy(ptm, os.Stdin)
		// end child code.
		// Just be lazy, in case we screw the order up again.
		m.Undo("")
		m.Undo("")
		m.Undo("")
	}
}
