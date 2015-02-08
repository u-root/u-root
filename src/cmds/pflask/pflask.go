package main

import (
	"os/exec"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
	"strings"
	//"user"

	"pty"
)

type cgroup string

func (c cgroup) apply(s string, f func(s string)) {
	for _,g := range strings.Split(s, ",") {
		p := path.Join(string(c), g)
		f(p)
	}
}

func (c cgroup) Validate(s string) {
	c.apply(s, func(s string) {
		if st, err := os.Stat(path.Join(string(c), s)); err != nil {
			log.Fatal("%v", err)
		} else if ! st.IsDir() {
				log.Fatal("%s: not a directory", s)
		}})
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

func (c cgroup) Do(pid int) {
	cgn := fmt.Sprintf("pflask.%d", pid)
	c.apply(string(c), func(s string) {
		c.Create(s, cgn)
		c.Attach(s, cgn, 1)
	})
}

type mount struct {
	src, dst, mtype, opts string
	flags uintptr
	dir bool
	mounted bool
}

type mlist struct {
	mounts[]*mount
}

func NewMlist(base string) (*mlist, error){
	m := &mlist{}
	if err := syscall.Mount("", "/", "", syscall.MS_SLAVE|syscall.MS_REC, ""); err != nil {
		err := fmt.Errorf("Mount :%s: on :%s: type :%s: flags %x: opts :%v: %v\n", 
			"", "/", "", syscall.MS_SLAVE|syscall.MS_REC, "", err)
		return nil, err
	}

	return m, nil
}

func (m *mlist) Add(base, src, dst, mtype, opts string, flags uintptr, dir bool) {
	m.mounts = append(m.mounts, &mount{src: src, dst: path.Join(base, dst), mtype: mtype, flags: flags, opts: opts, dir: dir})

}

func (m* mount) One() error {
	if m.dir {
		if err := os.MkdirAll(m.dst, 0755); err != nil {
			return fmt.Errorf("One: mkdirall %v: %v", m.dst, err)
		}
	}
	if err := syscall.Mount(m.src, m.dst, m.mtype, m.flags, m.opts); err != nil {
		return fmt.Errorf("Mount :%s: on :%s: type :%s: flags %x: opts :%v: %v\n", 
			m.src, m.dst, m.mtype, m.flags, m.opts, err)
	}
	m.mounted = true
	return nil
}
func (m *mlist) Do(base, console string) {
	ok := true
	if base != "" {
		m.Add(base, "proc", "/proc", "proc", "", 
			syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV, true)

		m.Add(base, "/proc/sys", "/proc/sys", "", "", 
			syscall.MS_BIND, true)

		m.Add(base, "", "/proc/sys", "", "", 
			syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REMOUNT, true)

		m.Add(base, "sysfs", "/sys", "sysfs", "",
			syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_RDONLY, true)

		m.Add(base, "tmpfs", "/dev", "tmpfs", "mode=755", 
			syscall.MS_NOSUID | syscall.MS_STRICTATIME, true)

		m.Add(base, "devpts", "/dev/pts", "devpts","newinstance,ptmxmode=000,mode=620,gid=5",
			syscall.MS_NOSUID | syscall.MS_NOEXEC, true)

		// OOPS! have to mknod first. Now I see why they did it the way they
		// did.
		//m.Add(base, console, "/dev/console", "", "", syscall.MS_BIND, false)

		m.Add(base, "tmpfs", "/dev/shm", "tmpfs", "mode=1777", 
			syscall.MS_NOSUID | syscall.MS_STRICTATIME | syscall.MS_NODEV, true)

		m.Add(base, "tmpfs", "/run", "tmpfs", "mode=755",
			syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_STRICTATIME, true)


	}

	for _, m := range m.mounts[1:] {
		err := m.One()
		if err != nil {
			log.Printf(err.Error())
			ok = false
		}
	}
	if ! ok {
		m.Undo()
		log.Fatal("Not all mounts succeeded.")
	}
}

func (m *mlist) Undo() {
	for i := range m.mounts {
		m := m.mounts[len(m.mounts)-i-1]
		if ! m.mounted {
			continue
		}
		if err := syscall.Unmount(m.dst, 0); err != nil {
			log.Printf("Unmounting %v: %v", m, err)
		}
		m.mounted = false
	}
}

func copy_nodes(base, console string) {
	nodes := []string {
		console,
		"/dev/tty",
		"/dev/full",
		"/dev/null",
		"/dev/zero",
		"/dev/random",
		"/dev/urandom", }

	if err := os.Chmod(console, 0600); err != nil {
		log.Printf("%v", err)
	}
	if err := os.Chown(console, 0, 0); err != nil {
		log.Printf("%v", err)
	}

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

func make_symlinks(base string) {
	linkit := []struct {
		src, dst string
	} {
		{"/dev/pts/ptmx", "/dev/ptmx"},
		{"/proc/kcore",	"/dev/core"},
		{"/proc/self/fd", "/dev/fd"},
		{"/proc/self/fd/0", "/dev/stdin"},
		{"/proc/self/fd/1","/dev/stdout"},
		{"/proc/self/fd/2", "/dev/stderr"},

	}

	for i := range linkit {
		dst := path.Join(base, linkit[i].dst)

		if err := os.Symlink(linkit[i].src, dst); err != nil {
			log.Printf("%v", err)
		}
	}
}

func do_chroot(dest string) {

	if err := os.Chdir(dest); err != nil {
		log.Fatalf("%v", err)
	}
	if err := syscall.Chroot(dest); err != nil {
		log.Fatalf("%v", err)
	}
	if err := os.Chdir("/"); err != nil {
		log.Fatalf("%v", err)
	}
}



var (
	cg = flag.String("cgroup", "/sys/fs/cgroup", "set the cgroups")
	mnt = flag.String("mount", "", "define mounts")
	chroot = flag.String("chroot", "", "where to chroot to")
	chdir = flag.String("chdir", "", "where to chrdir to")
	dest = flag.String("dest", "", "where the root is")
	console = flag.String("console", "/dev/console", "where the root is")
	keepenv = flag.Bool("keepenv", false, "Keep the environment")
	env = flag.String("env", "", "other environment variables")
	user = flag.String("user", "root"/*user.User.Username*/, "User name")
)
	
func main() {
	// note the unshare system call worketh not for Go. You have to run
	// this under the unshare command. Good times.
	flag.Parse()
	if *dest == "" {
		log.Fatalf("you are required to set the dest via --dest")
	}
	fmt.Printf("greetings\n")
	a := flag.Args()

	// Just create the container and run with it for now.
	c := cgroup(*cg)
	ppid := 1048576
	m, err := NewMlist(*dest)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// child code. Not really. What really happens here is we set
	// ourselves into the container, and spawn the child. It's a bit odd
	// but we're the master, but we'll run in the container? I don't know
	// how else to do it. This may require we set some things up first,
	// esp. the network. But, it's all fun and games until someone loses
	// an eye.
	if true {
		c.Do(ppid)

		m.Do(*dest, *console)

			copy_nodes(*dest, *console)

			//make_ptmx(dest);

			make_symlinks(*dest)

			//make_console(dest, master_name);

			do_chroot(*dest)

		//umask(0022);

		/* TODO: drop capabilities */

		//do_user(user);
/*
		if (change != NULL) {
			rc = chdir(change);
			if (rc < 0) sysf_printf("chdir()");
*/

		e := make(map[string]string)
		if *keepenv {
			for _, v := range os.Environ() {
				k := strings.SplitN(v, "=", 2)
				e[k[0]] = k[1]
			}
		}
		
			term := os.Getenv("TERM")
			e["TERM"] = term
			e["PATH"] =  "/usr/sbin:/usr/bin:/sbin:/bin"
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
			c.Env = append(c.Env, k + "=" + v)
		}

		f, err := pty.Start(c)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, f)
	}
	// end child code.
	m.Undo()
}
