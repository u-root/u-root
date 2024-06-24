// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pox packages dynamic executable into an archive.
//
// Synopsis:
//
//	pox [-[-verbose]|v] [-[-file]|f tcz-file] -[-create]|c FILE [...FILE]
//	pox [-[-verbose]|v] [-[-file]|f tcz-file] -[-run|r] PROGRAM -- [...ARGS]
//	pox [-[-verbose]|v] [-[-file]|f tcz-file] -[-create]|c -[-run|r] PROGRAM -- [...ARGS]
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
//	create|c: create the TCZ file.
//	verbose|d: verbose
//	file|f FILE: file name (default /tmp/pox.tcz)
//	run|r: Runs the first non-flag argument to pox.  Remaining arguments will
//	       be passed to the program.  Use '--' before any flag-like arguments
//	       to prevent pox from interpretting the flags.
//	self|s: Create a self-extracting elf. This implies -z.
//	zip|z: Use zip and unzip instead of a loopback mounted squashfs. Be sure
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
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/ldd"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/loop"
	"github.com/u-root/u-root/pkg/uzip"
)

const usage = "pox [-[-verbose]|v] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]"

type mp struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
	perm   os.FileMode // for target in the chroot
}

var (
	verbose = flag.BoolP("verbose", "v", false, "enable verbose prints")
	run     = flag.BoolP("run", "r", false, "Run the first file argument")
	create  = flag.BoolP("create", "c", false, "create it")
	zip     = flag.BoolP("zip", "z", false, "use zip instead of squashfs")
	self    = flag.BoolP("self", "s", false, "use self-extracting zip")
	file    = flag.StringP("output", "f", "/tmp/pox.tcz", "Output file")
	extra   = flag.StringP("extra", "e", "", `comma-separated list of extra directories to add (on create) and binds to do (on run).
You can specify what directories to add, and when you run, specify what directories are bound over them, e.g.:
pox -c -e /tmp,/etc commands ....
pox -r -e /a/b/c/tmp:/tmp,/etc:/etc commands ...
This can also be passed in with the POX_EXTRA variable.
`)
	v = func(string, ...interface{}) {}
)

// When chrooting, programs often want to access various system directories:
var chrootMounts = []mp{
	// mount --bind /sys /chroot/sys
	{"/sys", "/sys", "", mount.MS_BIND, "", 0o555},
	// mount -t proc /proc /chroot/proc
	{"/proc", "/proc", "proc", 0, "", 0o555},
	// mount --bind /dev /chroot/dev
	{"/dev", "/dev", "", mount.MS_BIND, "", 0o755}}

func poxCreate(bin ...string) error {
	if len(bin) == 0 {
		return fmt.Errorf(usage)
	}
	names, err := ldd.List(bin...)
	if err != nil {
		var stderr []byte
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr = eerr.Stderr
		}
		return fmt.Errorf("running ldd on %v: %v %s", bin, err, stderr)
	}
	// At some point the ldd API changed and it no longer includes the
	// bins.
	names = append(names, bin...)

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
		v("Adding %q", f)
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
		v("Adding mount %q, perm %s", d, m.perm.String())
		if err := os.MkdirAll(d, m.perm); err != nil {
			return err
		}
	}
	if err := os.Remove(*file); err != nil && !os.IsNotExist(err) {
		return err
	}

	if *self {
		// Make a copy of the exe and append the zip file.
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		if err := cp.Copy(exe, *file); err != nil {
			return err
		}
		if err := uzip.AppendZip(dir, *file, bin[0]); err != nil {
			return err
		}
	} else if *zip {
		if err := uzip.ToZip(dir, *file, ""); err != nil {
			return err
		}
	} else {
		c := exec.Command("mksquashfs", dir, *file, "-noappend")
		o, err := c.CombinedOutput()
		v("%v", string(o))
		if err != nil {
			return fmt.Errorf("%v: %v: %v", c.Args, string(o), err)
		}
	}

	v("Done, your pox is %q", *file)
	return nil
}

func poxRun(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf(usage)
	}
	dir, err := os.MkdirTemp("", "pox")
	if err != nil {
		return err
	}

	if *zip {
		if err := uzip.FromZip(*file, dir); err != nil {
			return err
		}
	} else {
		lo, err := loop.New(*file, "squashfs", "")
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
		v("mount(%q, %q, %q, %q, %#x)", m.source, filepath.Join(dir, m.target), m.fstype, m.data, m.flags)
		mp, err := mount.Mount(m.source, filepath.Join(dir, m.target), m.fstype, m.data, m.flags)
		if err != nil {
			return err
		}
		defer mp.Unmount(0) //nolint:errcheck
	}

	// If you pass Command a path with no slashes, it'll use PATH from the
	// parent to resolve the path to exec.  Once we chroot, whatever path we
	// picked is undoubtably wrong.  Let's help them out: if they give us a
	// program with no /, let's look in /bin/.  If they want the root of the
	// chroot, they can use "./"
	if filepath.Base(args[0]) == args[0] {
		args[0] = filepath.Join(string(os.PathSeparator), "bin", args[0])
	}
	c := exec.Command(args[0], args[1:]...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{
		Chroot: dir,
	}
	c.Env = append(os.Environ(), "PWD=.")

	if err = c.Run(); err != nil {
		v("pox command exited with: %v", err)
	}

	return nil
}

func extraMounts(mountList string) error {
	if mountList == "" {
		return nil
	}
	v("Extra mounts: %q", mountList)
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
		v("Extra mounts: append %q to chrootMounts", m)
		chrootMounts = append(chrootMounts, m)
	}
	return nil
}

func pox(args ...string) error {
	// If the current executable is a zip file, extract and run.
	// Sneakily re-write os.Args to include a "-rzf" before flag parsing.
	// The zip comment contains the executable path once extracted.
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if comment, err := uzip.Comment(exe); err == nil {
		if comment == "" {
			return errors.New("expected zip comment on self-extracting pox")
		}
		os.Args = append([]string{
			os.Args[0],
			"-rzf", exe,
			"--",
			comment,
		}, os.Args[1:]...)
	}

	if *verbose {
		v = log.Printf
	}
	if err := extraMounts(*extra); err != nil {
		return err
	}
	if err := extraMounts(os.Getenv("POX_EXTRA")); err != nil {
		return err
	}
	if !*create && !*run {
		return fmt.Errorf(usage)
	}
	if *create {
		if err := poxCreate(args...); err != nil {
			return err
		}
	}
	if *run {
		return poxRun(args...)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := pox(flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
