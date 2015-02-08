package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
	"strings"
	//"user"

	_ "pty"
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
	chdir = flag.String("chroot", "", "where to chrdir to")
	dest = flag.String("dest", "", "where the root is")
	console = flag.String("console", "", "where the root is")
	keepenv = flag.Bool("keepenv", false, "Keep the environment")
	env = flag.String("env", "", "other environment variables")
	user = flag.String("user", "root"/*user.User.Username*/, "User name")
)
	
func main() {
	// option is not an option.
	fmt.Printf("greetings\n")
	//a := flag.Args()

	// Just create the container and run with it for now.
	c := cgroup(*cg)
	ppid := 1048576
	m := NewMlist()
	// child code.
	if true {
		c.Do(ppid)

		m.Do(*dest, *console)

		if (dest != nil) {
			copy_nodes(*dest, *console)

			//make_ptmx(dest);

			make_symlinks(*dest)

			//make_console(dest, master_name);

			do_chroot(*dest)
		}

		//umask(0022);

		/* TODO: drop capabilities */

		//do_user(user);
/*
		if (change != NULL) {
			rc = chdir(change);
			if (rc < 0) sysf_printf("chdir()");
*/

		
		if dest != nil {
			term := os.Getenv("TERM")
			
			if ! *keepenv {
				os.Clearenv()
			}
			os.Setenv("PATH", "/usr/sbin:/usr/bin:/sbin:/bin")
			os.Setenv("USER", *user)
			os.Setenv("LOGNAME", *user)
			os.Setenv("TERM", term)
		}
		
		if env != nil {
			for _, c := range strings.Split(*env, ",") {
				k := strings.SplitN(c, "=", 2)
				if len(k) != 2 {
					log.Printf("Bogus environment string %v", c)
				}
				if err := os.Setenv(k[0], k[1]); err != nil {
					log.Printf(err.Error())
				}
			}
		}
		os.Setenv("container", "pflask")
		/*
		if (argc > optind)
			rc = execvpe(argv[optind], argv + optind, environ);
		else
			rc = execle("/bin/bash", "-bash", NULL, environ);

		if (rc < 0) sysf_printf("exec()");
*/
	}
	// end child code.
}
