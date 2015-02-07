package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
	"strings"
)

type cgroup string

func (c *cgroup) apply(s string, f func(s string)) {
	for _,g := range strings.Split(s, ",") {
		p := path.Join(string(*c), g)
		f(p)
	}
}

func (c *cgroup) Validate(s string) {
	c.apply(s, func(s string) {
		if st, err := os.Stat(s); err != nil {
			log.Fatal("%v", err)
		} else if ! st.IsDir() {
				log.Fatal("%s: not a directory", s)
		}})
}

func (c *cgroup) Create(s string) {
	c.apply(s, func(s string) {
		if err := os.MkdirAll(s, 0700); err != nil {
			log.Fatal(err)
		}})
}

func (c *cgroup) Attach(s string, pid int) {
	c.apply(s, func(s string) {
		t := path.Join(s, "tasks")
		b := []byte(fmt.Sprintf("%v", pid))
		if err := ioutil.WriteFile(t, b, 0600); err != nil {
			log.Fatal(err)
		}})
}

func (c *cgroup) Destroy(s string) {
	c.apply(s, func(s string) {
		if err := os.RemoveAll(s); err != nil {
			log.Fatal(err)
		}})
}

type mount struct {
	src, dst, mtype, opts string
	flags uintptr
}

type mlist struct {
	mounts[]*mount
}

func NewMlist() *mlist {
	m := &mlist{}
	m.Add("", "", "/", "", "", syscall.MS_SLAVE|syscall.MS_REC)
	return m
}

func (m *mlist) Add(base, src, dst, mtype, opts string, flags uintptr) {
	m.mounts = append(m.mounts, &mount{src: src, dst: path.Join(base, dst), mtype: mtype, flags: flags, opts: opts})

}

func (m* mount) One() error {
	if err := syscall.Mount(m.src, m.dst, m.mtype, m.flags, m.opts); err != nil {
		err := fmt.Errorf("Mount :%s: on :%s: type :%s: flags %x: %v\n", 
			m.src, m.dst, m.mtype, m.flags, m.opts, err)
		return err
	}
	return nil
}
func (m *mlist) Do(base, console string) {
	// Accumulate all the errors
	// Do the first one to test.
	e := ""
	if err := m.mounts[0].One(); err != nil {
		e = e + "\n" + err.Error()
	}
	
	
	if base != "" {
		m.Add(base, "proc", "/proc", "proc", "", 
			syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV)

		m.Add(base, "/proc/sys", "/proc/sys", "", "", 
			syscall.MS_BIND)

		m.Add(base, "", "/proc/sys", "", "", 
			syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REMOUNT)

		m.Add(base, "sysfs", "/sys", "sysfs", "",
			syscall.MS_NOSUID | syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_RDONLY)

		m.Add(base, "tmpfs", "/dev", "tmpfs", "mode=755", 
			syscall.MS_NOSUID | syscall.MS_STRICTATIME)

		m.Add(base, "devpts", "/dev/pts", "devpts","newinstance,ptmxmode=000,mode=620,gid=5",
			syscall.MS_NOSUID | syscall.MS_NOEXEC)

		m.Add(base, console, "/dev/console", "", "", syscall.MS_BIND)

		m.Add(base, "tmpfs", "/dev/shm", "tmpfs", "mode=1777", 
			syscall.MS_NOSUID | syscall.MS_STRICTATIME | syscall.MS_NODEV)

		m.Add(base, "tmpfs", "/run", "tmpfs", "mode=755",
			syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_STRICTATIME)


	}

	for _, m := range m.mounts[1:] {
		err := m.One()
		if err != nil {
			e = e + "\n" + err.Error()
		}
	}

	if e == "" {
		return
	}

	log.Printf("%v", e)
	log.Fatal("Not all mounts succeeded.")
}

func setup_nodes(base, console string) {
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

	for i, n := range nodes {
		st, err := os.Stat(n)
		if err != nil {
			log.Printf("%v", err)
		}
		nn := path.Join(base, n)
		// special case.
		if i == 0 {
			nn = path.Join(base, "/dev/console")
		}
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


func main() {
	fmt.Printf("greetings\n")
}
