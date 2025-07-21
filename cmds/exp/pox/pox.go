// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pox packages dynamic executable into an archive.
//
// Synopsis:
//
//	pox [-v] [-f tcz-file] -c FILE [...FILE]
//	pox [-v] [-f tcz-file] -r PROGRAM -- [...ARGS]
//	pox [-v] [-f tcz-file] -c -r PROGRAM -- [...ARGS]
//
// Description:
//
//	pox packages a dynamic executable into an archive for use on another
//	machine. By default, it uses the tcz format compatible with tinycore.
//
//	pox supports 3 archive formats:
//	1) squashfs (default): The tcz is a squashfs. This requires mksquashfs.
//	2) zip
//	3) elf+zip: Self-extracting.
//
// Options:
//
//	c: create the TCZ file.
//	d: verbose
//	f FILE: file name (default /tmp/pox.tcz)
//	r: Runs the first non-flag argument to pox.  Remaining arguments will
//	       be passed to the program.  Use '--' before any flag-like arguments
//	       to prevent pox from interpretting the flags.
//	s: Create a self-extracting elf. This implies -z.
//	z: Use zip and unzip instead of a loopback mounted squashfs. Be sure
//	       to use -z for both creation and running, or not at all.
//	For convenience and testing, you can create and run a pox in one command.
//
// Example:
//
//	$ pox -c /bin/bash /bin/cat /bin/ls /etc/hosts
//	Will build a squashfs, which will be /tmp/pox.tcz
//
//	$ sudo pox -r /bin/bash
//	Will drop you into the /tmp/pox.tcz running bash
//	You can use ls and cat on /etc/hosts.
//
//	Simpler example, with arguments:
//	$ sudo pox -r /bin/ls -- -la
//	will run `ls -la` and exit.
//
//	$ sudo pox -r -- /bin/ls -la
//	Syntactically easier: the program name can come after '--'
//
//	$ sudo pox -c -r /bin/bash
//	Create a pox with a bash and run it.
//
//	$ pox -cvsf date /bin/date
//	Creates a self-executing pox called "date".
//	$ ./date --utc
//
// Notes:
//
//   - When running a pox, you likely need sudo to chroot
//
//   - Binaries run out of a chroot often need files you are unaware of. For
//     instance, if bash can't find terminfo files, it won't know to handle
//     backspaces properly. (They occur, but are not shown). To fix this, pass
//     pox all of the files you need.  For bash: `find /lib/terminfo -type f`.
//
//   - Other programs rely on helper functions, such as '/bin/man'. If your
//     program has built-in help commands that trigger man pages, e.g. "git
//     help foo", you'll want to include /bin/man too. But you'll also need
//     everything that man uses, such as /etc/manpath.config. My advice: skip
//     it.
//
//   - When adding all files in a directory, the easiest thing to do is:
//     `find $DIR -type f` (Note the ticks: this is a bash command execution).
//
//   - When creating a pox with an executable with shared libraries that are
//     not installed on your system, such as for a project installed in your
//     home directory, run pox from the installation prefix directory, such
//     that the lib/ and bin/ are in pox's working directory. Pox will strip
//     its working directory from the paths of the files it builds. Having bin/
//     in the root of the pox file helps with PATH lookups, and not having the
//     full path from your machine in the pox file makes it easier to extract a
//     pox file to /usr/local/.
//
//   - pox is not a security boundary. chroot is well known to have holes.
//     Pox is about enabling execution. Don't expect it to "wall things off".
//     In fact, we mount /dev, /proc, and /sys; and you can add more things.
//     Commands run under pox are just as dangerous as anything else.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/core/cp"
	"github.com/u-root/u-root/pkg/ldd"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/loop"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
	"github.com/u-root/u-root/pkg/uzip"
)

var ErrUsage = errors.New("pox [-v] [-r] [-c]  [-f tcz-file] file [...file]")

type mp struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
	perm   os.FileMode // for target in the chroot
}

type cmd struct {
	// flag
	verbose bool
	run     bool
	create  bool
	zip     bool
	self    bool
	file    string
	extra   string
	// positional args
	arg0       string
	args       []string
	pathInRoot string
	// verbose output
	debug func(string, ...interface{})
	// io
	in  io.Reader
	out io.Writer
	err io.Writer
}

// When chrooting, programs often want to access various system directories:
var chrootMounts = []mp{
	// mount --bind /sys /chroot/sys
	{"/sys", "/sys", "", mount.MS_BIND, "", 0o555},
	// mount -t proc /proc /chroot/proc
	{"/proc", "/proc", "proc", 0, "", 0o555},
	// mount --bind /dev /chroot/dev
	{"/dev", "/dev", "", mount.MS_BIND, "", 0o755},
}

func (c cmd) poxCreate() error {
	if len(c.args) == 0 {
		return ErrUsage
	}
	names, err := ldd.List(c.args...)
	if err != nil {
		var stderr []byte
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr = eerr.Stderr
		}
		return fmt.Errorf("running ldd on %v: %w %s", c.args, err, stderr)
	}
	// At some point the ldd API changed and it no longer includes the
	// bins.
	names = append(names, c.args...)

	sort.Strings(names)
	// Now we need to make a template file hierarchy and put
	// the stuff we want in there.
	dir, err := os.MkdirTemp("", "pox")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// We don't use defer() here to close files as
	// that can cause open failures with a large enough number.
	for _, f := range names {
		c.debug("Adding %q", f)
		fi, err := os.Stat(f)
		if err != nil {
			return err
		}
		in, err := os.Open(f)
		if err != nil {
			return err
		}
		f = strings.TrimPrefix(f, pwd)
		dfile := filepath.Join(dir, f)
		d := filepath.Dir(dfile)
		if err := os.MkdirAll(d, 0o755); err != nil {
			in.Close()
			return err
		}
		out, err := os.OpenFile(dfile, os.O_WRONLY|os.O_CREATE, fi.Mode().Perm())
		if err != nil {
			in.Close()
			return err
		}
		_, err = io.Copy(out, in)
		in.Close()
		out.Close()
		if err != nil {
			return err
		}
	}

	for _, m := range chrootMounts {
		d := filepath.Join(dir, m.target)
		c.debug("Adding mount %q, perm %s", d, m.perm.String())
		if err := os.MkdirAll(d, m.perm); err != nil {
			return err
		}
	}
	if err := os.Remove(c.file); err != nil && !os.IsNotExist(err) {
		return err
	}

	if c.self {
		// Make a copy of the exe and append the zip file.
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		if err := cp.Copy(exe, c.file); err != nil {
			return err
		}
		if err := uzip.AppendZip(dir, c.file, c.args[0]); err != nil {
			return err
		}
	} else if c.zip {
		if err := uzip.ToZip(dir, c.file, ""); err != nil {
			return err
		}
	} else {
		ec := exec.Command("mksquashfs", dir, c.file, "-noappend")
		o, err := ec.CombinedOutput()
		c.debug("%v", string(o))
		if err != nil {
			return fmt.Errorf("%v: %v: %w", ec.Args, string(o), err)
		}
	}

	c.debug("Done, your pox is %q", c.file)
	return nil
}

func (c cmd) poxRun() error {
	if os.Getuid() != 0 {
		return fmt.Errorf("this primitive kernel requires root permissions to run pox:%w", os.ErrPermission)
	}
	v := c.debug
	dir, err := os.MkdirTemp("", "pox")
	if err != nil {
		return err
	}

	if c.zip {
		v("uzip...")
		if err := uzip.FromZip(c.file, dir); err != nil {
			return err
		}
		v("...ok")
	} else {
		lo, err := loop.New(c.file, "squashfs", "")
		if err != nil {
			return err
		}
		defer lo.Free() //nolint:errcheck

		mountPoint, err := lo.Mount(dir, 0)
		if err != nil {
			return err
		}
		defer mountPoint.Unmount(0) //nolint:errcheck
	}
	for _, m := range chrootMounts {
		c.debug("mount(%q, %q, %q, %q, %#x)", m.source, filepath.Join(dir, m.target), m.fstype, m.data, m.flags)
		mp, err := mount.Mount(m.source, filepath.Join(dir, m.target), m.fstype, m.data, m.flags)
		if err != nil {
			return err
		}
		defer mp.Unmount(0) //nolint:errcheck
	}

	// We called ourselves with name -rz -f name -- command args
	// So the args to exec.Command should just be c.args...
	ec := exec.Command(c.pathInRoot, c.args[1:]...)
	v("cmd q args %q", c.args)
	// Go does the chroot first, then the chdir.
	// So set dir to /
	// If that is not done, pwd won't work at all, since
	// Dir will not be set and the process cwd will be outside
	// the chroot. We ran for years with this bug, but apptainer
	// in a pox revealed. And, yes, sometimes you have to put
	// apptainer in a pox; it's dynamically linked!
	ec.Dir = "/"
	ec.Stdin, ec.Stdout, ec.Stderr = c.in, c.out, c.err
	ec.SysProcAttr = &syscall.SysProcAttr{
		Chroot: dir,
	}
	v("chroot %q", dir)
	ec.Env = append(os.Environ(), "PWD=.")

	v("cmd %v", ec)
	if err = ec.Run(); err != nil {
		v("pox command exited with: %v", err)
		return err
	}
	v("POX ran it; did anything go")

	return nil
}

func (c cmd) extraMounts(mountList string) error {
	if mountList == "" {
		return nil
	}
	c.debug("Extra mounts: %q", mountList)
	// We have to specify the extra directories and do the create here b/c it is a squashfs. Sorry.
	for _, e := range strings.Split(mountList, ",") {
		m := mp{flags: mount.MS_BIND, perm: 0o755}
		mp := strings.Split(e, ":")
		switch len(mp) {
		case 1:
			m.source, m.target = mp[0], mp[0]
		case 2:
			m.source, m.target = mp[0], mp[1]
		default:
			return fmt.Errorf("%q is not in the form src:target", mp)
		}
		c.debug("Extra mounts: append %q to chrootMounts", m)
		chrootMounts = append(chrootMounts, m)
	}
	return nil
}

func (c cmd) start() error {
	if c.verbose {
		c.debug = log.Printf
	}
	v := c.debug
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	v("%q: %q, %v", exe, c.args, c)
	if comment, err := uzip.Comment(exe); err == nil {
		if len(comment) == 0 {
			return fmt.Errorf("zip comment is empty, no path to command:%w", os.ErrNotExist)
		}
		v("self-running zip file, comment %q, arg0 %q, args: %q", comment, c.arg0, c.args)
		c.args = append([]string{c.arg0}, c.args...)
		c.run = true
		c.zip = true
		c.file = exe
		c.pathInRoot = comment
		v("re-written args: %q %q", c.pathInRoot, c.args)

	}

	if err := c.extraMounts(c.extra); err != nil {
		return err
	}
	if err := c.extraMounts(os.Getenv("POX_EXTRA")); err != nil {
		return err
	}
	if !c.create && !c.run {
		return ErrUsage
	}
	if c.create {
		if err := c.poxCreate(); err != nil {
			return err
		}
	}
	if c.run {
		return c.poxRun()
	}
	return nil
}

func command(in io.Reader, out, err io.Writer, args []string) *cmd {
	c := cmd{arg0: args[0], in: in, out: out, err: err}

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	f.BoolVar(&c.verbose, "v", false, "enable verbose prints")

	f.BoolVar(&c.run, "r", false, "Run the first file argument")

	f.BoolVar(&c.create, "c", false, "create it")

	f.BoolVar(&c.zip, "z", false, "use zip instead of squashfs")

	f.BoolVar(&c.self, "s", false, "use self-extracting zip")

	f.StringVar(&c.file, "f", "/tmp/pox.tcz", "Output file")

	const extraUsage = `comma-separated list of extra directories to add (on create) and binds to do (on run).
You can specify what directories to add, and when you run, specify what directories are bound over them, e.g.:
pox -c -e /tmp,/etc commands ....
pox -r -e /a/b/c/tmp:/tmp,/etc:/etc commands ...
This can also be passed in with the POX_EXTRA variable.
`
	f.StringVar(&c.extra, "e", "", extraUsage)

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))

	c.args = f.Args()
	c.debug = func(string, ...interface{}) {}

	return &c
}

func main() {
	if err := command(os.Stdin, os.Stdout, os.Stderr, os.Args).start(); err != nil {
		log.Fatal(err)
	}
}
