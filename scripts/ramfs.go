// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/u-root/u-root/uroot"
)

type copyfiles struct {
	dir  string
	spec string
}

type goDirs struct {
	Dir        string
	Deps       []string
	GoFiles    []string
	SFiles     []string
	HFiles     []string
	Goroot     bool
	ImportPath string
}

const (
	devcpio   = "scripts/dev.cpio"
	urootPath = "src/github.com/u-root/u-root"
	urootCmds = "github.com/u-root/u-root/cmds"
)

var (
	// be VERY CAREFUL with these. If you have an empty line here it will
	// result in cpio copying the whole tree.
	goList = `{{.Goroot}}
go
pkg/include
VERSION.cache`
	urootList = `{{.Gopath}}
`
	config struct {
		Goroot          string
		Godotdot        string
		Godot           string
		Arch            string
		Goos            string
		Gopath          string
		Urootpath       string
		TempDir         string
		Go              string
		Debug           bool
		Fail            bool
		TestChroot      bool
		RemoveDir       bool
		InitialCpio     string
		UseExistingInit bool
	}
	pkgList     []string
	dirs        map[string]bool
	deps        map[string]bool
	gorootFiles map[string]bool
	urootFiles  map[string]bool
	// the whitelist is a list of u-root tools that we feel
	// can replace existing tools. It is, sadly, a very short
	// list at present.
	whitelist = []string{"date"}
	debug     = nodebug
)

func nodebug(string, ...interface{}) {}

func getenvOrDefault(e, defaultValue string) string {
	v := os.Getenv(e)
	if v == "" {
		v = defaultValue
	}
	return v
}

func lsr(n string, w *os.File) error {
	n = n + "/"
	err := filepath.Walk(n, func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cn := strings.TrimPrefix(name, n)
		fmt.Fprintf(w, "%v\n", cn)
		return nil
	})
	return err
}

// cpio copies a tree from one place to another, defined by a template.
func cpiop(c string) error {

	t := template.Must(template.New("filelist").Parse(c))
	var b bytes.Buffer
	if err := t.Execute(&b, config); err != nil {
		log.Fatalf("spec %v: %v\n", c, err)
	}

	n := strings.Split(b.String(), "\n")
	debug("cpiop: from %v, to %v, :%v:\n", n[0], n[1], n[2:])

	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	cmd := exec.Command("sudo", "cpio", "--make-directories", "-p", path.Join(config.TempDir, n[1]))
	d := path.Clean(n[0])
	cmd.Dir = d
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	if config.Debug {
		cmd.Stderr = os.Stderr
	}
	debug("Run %v @ %v", cmd, cmd.Dir)
	err = cmd.Start()
	if err != nil {
		log.Printf("%v\n", err)
	}

	for _, v := range n[2:] {
		debug("%v\n", v)
		err := filepath.Walk(path.Join(d, v), func(name string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Printf(" WALK FAIL%v: %v\n", name, err)
				// That's ok, sometimes things are not there.
				return filepath.SkipDir
			}
			cn := strings.TrimPrefix(name, d+"/")
			if cn == ".git" {
				return filepath.SkipDir
			}
			fmt.Fprintf(w, "%v\n", cn)
			//log.Printf("c.dir %v %v %v\n", d, name, cn)
			return nil
		})
		if err != nil {
			log.Printf("%s: %v\n", v, err)
		}
	}
	w.Close()
	debug("Done sending files to external")
	err = cmd.Wait()
	if err != nil {
		log.Printf("%v\n", err)
	}
	debug("External cpio is done")
	return nil
}

// buildToolChain builds the four binaries needed for the go toolchain:
// go, compile, link, and asm. We do this to ensure we get smaller binaries.
// Smaller, in this, meaning 25M instead of 33M. What a world!
func buildToolChain() {
	goBin := path.Join(config.TempDir, "go/bin/go")
	cmd := exec.Command("go", "build", "-x", "-a", "-installsuffix", "cgo", "-ldflags", "-s -w", "-o", goBin)
	cmd.Dir = path.Join(config.Goroot, "src/cmd/go")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Building statically linked go tool info %v: %v, %v\n", goBin, string(o), err)
	}

	toolDir := path.Join(config.TempDir, fmt.Sprintf("go/pkg/tool/%v_%v", config.Goos, config.Arch))

	for _, pkg := range []string{"compile", "link", "asm"} {
		c := path.Join(toolDir, pkg)
		cmd = exec.Command("go", "build", "-x", "-a", "-installsuffix", "cgo", "-ldflags", "-s -w", "-o", c, "cmd/"+pkg)
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		if o, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("Building statically linked %v: %v, %v\n", pkg, string(o), err)
		}
	}
}

// It's annoying asking them to set lots of things. So let's try to figure it out.
func guessgoarch() {
	config.Arch = os.Getenv("GOARCH")
	if config.Arch != "" {
		config.Arch = path.Clean(config.Arch)
		return
	}
	log.Printf("GOARCH is not set, trying to guess")
	u, err := uroot.Uname()
	if err != nil {
		log.Printf("uname failed, using default amd64")
		config.Arch = "amd64"
	} else {
		switch {
		case u.Machine == "i686" || u.Machine == "i386" || u.Machine == "x86":
			config.Arch = "386"
		case u.Machine == "x86_64" || u.Machine == "amd64":
			config.Arch = "amd64"
		case u.Machine == "armv7l" || u.Machine == "armv6l":
			config.Arch = "arm"
		case u.Machine == "ppc" || u.Machine == "ppc64":
			config.Arch = "ppc64"
		default:
			log.Printf("Unrecognized arch")
			config.Fail = true
		}
	}
}
func guessgoroot() {
	config.Goroot = os.Getenv("GOROOT")
	if config.Goroot != "" {
		config.Goroot = path.Clean(config.Goroot)
		log.Printf("Using %v from the environment as the GOROOT", config.Goroot)
		config.Godotdot = path.Dir(config.Goroot)
		return
	}
	log.Print("Goroot is not set, trying to find a go binary")
	p := os.Getenv("PATH")
	paths := strings.Split(p, ":")
	for _, v := range paths {
		g := path.Join(v, "go")
		log.Printf("Try %s as the Go binary", g)
		if _, err := os.Stat(g); err == nil {
			config.Goroot = path.Dir(v)
			config.Godotdot = path.Dir(config.Goroot)
			log.Printf("Guessing that goroot is %v from $PATH", config.Goroot)
			return
		}
	}
	log.Printf("GOROOT is not set and can't find a go binary in %v", p)
	config.Fail = true
}

func guessgopath() {
	defer func() {
		config.Godotdot = path.Dir(config.Goroot)
	}()
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		config.Gopath = gopath
		config.Urootpath = path.Join(gopath, urootPath)
		return
	}
	// It's a good chance they're running this from the u-root source directory
	log.Fatalf("Fix up guessgopath")
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("GOPATH was not set and I can't get the wd: %v", err)
		config.Fail = true
		return
	}
	// walk up the cwd until we find a u-root entry. See if cmds/init/init.go exists.
	for c := cwd; c != "/"; c = path.Dir(c) {
		if path.Base(c) != "u-root" {
			continue
		}
		check := path.Join(c, "cmds/init/init.go")
		if _, err := os.Stat(check); err != nil {
			//log.Printf("Could not stat %v", check)
			continue
		}
		config.Gopath = c
		log.Printf("Guessing %v as GOPATH", c)
		os.Setenv("GOPATH", c)
		return
	}
	config.Fail = true
	log.Printf("GOPATH was not set, and I can't see a u-root-like name in %v", cwd)
	return
}

// goListPkg takes one package name, and computes all the files it needs to build,
// separating them into Go tree files and uroot files. For now we just 'go list'
// but hopefully later we can do this programmatically.
func goListPkg(name string) (*goDirs, error) {
	cmd := exec.Command("go", "list", "-json", name)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	debug("Run %v @ %v", cmd, cmd.Dir)
	j, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var p goDirs
	if err := json.Unmarshal([]byte(j), &p); err != nil {
		return nil, err
	}

	debug("%v, %v %v %v", p, p.GoFiles, p.SFiles, p.HFiles)
	for _, v := range append(append(p.GoFiles, p.SFiles...), p.HFiles...) {
		if p.Goroot {
			gorootFiles[path.Join(p.ImportPath, v)] = true
		} else {
			urootFiles[path.Join(p.ImportPath, v)] = true
		}
	}

	return &p, nil
}

// addGoFiles Computes the set of Go files to be added to the initramfs.
func addGoFiles() error {
	// For each directory in pkgList, add its files and all its
	// dependencies.  It would be nice to run go list -json with
	// lots of package names but it produces invalid JSON.  It
	// produces a stream thatis {}{}{} at the top level and the
	// decoders don't like that.
	for _, v := range pkgList {
		p, err := goListPkg(v)
		if err != nil {
			log.Printf("Can't do go list in %v, ignoring\n", v)
			continue
		}
		debug("cmd p is %v", p)
		for _, v := range p.Deps {
			deps[v] = true
		}
	}

	for v := range deps {
		if _, err := goListPkg(v); err != nil {
			log.Fatalf("%v", err)
		}
	}
	for v := range gorootFiles {
		goList += "\n" + path.Join("src", v)
	}
	for v := range urootFiles {
		urootList += "\n" + path.Join("src", v)
	}
	return nil
}

func globlist(s ...string) []string {
	// For each arg, use it as a Glob pattern and add any matches to the
	// package list. If there are no arguments, use [a-zA-Z]* as the glob pattern.
	var pat []string
	for _, v := range s {
		pat = append(pat, path.Join(config.Urootpath, "cmds", v))
	}
	if len(s) == 0 {
		pat = []string{path.Join(config.Urootpath, "cmds", "[a-zA-Z]*")}
	}
	return pat
}

// sad news. If I concat the Go cpio with the other cpios, for reasons I don't understand,
// the kernel can't unpack it. Don't know why, don't care. Need to create one giant cpio and unpack that.
// It's not size related: if the go archive is first or in the middle it still fails.
func main() {
	flag.BoolVar(&config.Debug, "d", false, "Debugging")
	flag.BoolVar(&config.TestChroot, "test", false, "test the directory by chrooting to it")
	flag.BoolVar(&config.UseExistingInit, "useinit", false, "If there is an existing init, don't replace it")
	flag.BoolVar(&config.RemoveDir, "removedir", true, "remove the directory when done -- cleared if test fails")
	flag.StringVar(&config.InitialCpio, "cpio", "", "An initial cpio image to build on")
	flag.StringVar(&config.TempDir, "tmpdir", "", "tmpdir to use instead of ioutil.TempDir")
	flag.Parse()
	if config.Debug {
		debug = log.Printf
	}

	var err error
	dirs = make(map[string]bool)
	deps = make(map[string]bool)
	gorootFiles = make(map[string]bool)
	urootFiles = make(map[string]bool)
	guessgoarch()
	config.Go = ""
	config.Goos = "linux"
	guessgoroot()
	guessgopath()
	if config.Fail {
		log.Fatal("Setup failed")
	}

	pat := globlist(flag.Args()...)

	debug("Initial glob is %v", pat)
	for _, v := range pat {
		g, err := filepath.Glob(v)
		if err != nil {
			log.Fatalf("Glob error: %v", err)
		}
		// We have a set of absolute paths in g.  We can not
		// use absolute paths in go list, however, so we have
		// to adjust them.
		for i := range g {
			g[i] = path.Join(urootCmds, path.Base(g[i]))
		}
		pkgList = append(pkgList, g...)
	}

	debug("Initial pkgList is %v", pkgList)

	if err := addGoFiles(); err != nil {
		log.Fatalf("%v", err)
	}

	if config.TempDir == "" {
		config.TempDir, err = ioutil.TempDir("", "u-root")
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	defer func() {
		if config.RemoveDir {
			log.Printf("Removing %v\n", config.TempDir)
			// Wow, this one is *scary*
			cmd := exec.Command("sudo", "rm", "-rf", config.TempDir)
			cmd.Stderr, cmd.Stdout = os.Stderr, os.Stdout
			err = cmd.Run()
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
	}()

	buildToolChain()

	if config.InitialCpio != "" {
		f, err := ioutil.ReadFile(config.InitialCpio)
		if err != nil {
			log.Fatalf("%v", err)
		}

		cmd := exec.Command("sudo", "cpio", "-i", "-v")
		cmd.Dir = config.TempDir
		// Note: if you print Cmd out with %v after assigning cmd.Stdin, it will print
		// the whole cpio; so don't do that.
		if config.Debug {
			cmd.Stdout = os.Stdout
		}
		debug("Run %v @ %v", cmd, cmd.Dir)

		// There's a bit of a tough problem here. There's lots of stuff owned by root in
		// these directories. They probably have to stay that way. But how do we create init
		// and do other things? For now, we're going to set the modes of select places to
		// 666 and remove a few things we know need to be removed.
		// It's hard to say what else to do.
		cmd.Stdin = bytes.NewBuffer(f)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Printf("Unpacking %v: %v", config.InitialCpio, err)
		}
	}

	if !config.UseExistingInit {
		init := path.Join(config.TempDir, "init")
		// Must move config.TempDir/init to inito if one is not there.
		inito := path.Join(config.TempDir, "inito")
		if _, err := os.Stat(inito); err != nil {
			// WTF? did Ron forget about rename? Yuck!
			if err := syscall.Rename(init, inito); err != nil {
				log.Printf("%v", err)
			}
		} else {
			log.Printf("Not replacing %v because there is already one there.", inito)
		}

		// Build init
		cmd := exec.Command("go", "build", "-x", "-a", "-installsuffix", "cgo", "-ldflags", "'-s'", "-o", init, ".")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = path.Join(config.Urootpath, "cmds/init")

		err = cmd.Run()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}

	// These produce arrays of strings, the first element being the
	// directory to walk from.
	cpio := []string{
		goList,
		urootList,
	}

	for _, c := range cpio {
		if err := cpiop(c); err != nil {
			log.Printf("Things went south. TempDir is %v", config.TempDir)
			log.Fatalf("Bailing out near line 666")
		}
	}

	debug("Done all cpio operations")

	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// First create the archive and put the device cpio in it.
	dev, err := ioutil.ReadFile(path.Join(config.Urootpath, devcpio))
	if err != nil {
		log.Fatalf("%v %v\n", dev, err)
	}

	debug("Creating initramf file")

	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", config.Goos, config.Arch)
	if err := ioutil.WriteFile(oname, dev, 0600); err != nil {
		log.Fatalf("%v\n", err)
	}

	// Now use the append option for cpio to append to it.
	// That way we get one cpio.
	// We need sudo as there may be files created from an initramfs that
	// can only be read by root.
	cmd := exec.Command("sudo", "cpio", "-H", "newc", "-o", "-A", "-F", oname)
	cmd.Dir = config.TempDir
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	debug("Run %v @ %v", cmd, cmd.Dir)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := lsr(config.TempDir, w); err != nil {
		log.Fatalf("%v\n", err)
	}
	w.Close()
	debug("Finished sending file list for initramfs cpio")
	err = cmd.Wait()
	if err != nil {
		log.Printf("%v\n", err)
	}
	debug("cpio for initramfs is done")
	defer func() {
		log.Printf("Output file is in %v\n", oname)
	}()

	if !config.TestChroot {
		return
	}

	// We need to populate the temp directory with dev.cpio. It's a chicken and egg thing;
	// we can't run init without, e.g., /dev/console and /dev/null.
	cmd = exec.Command("sudo", "cpio", "-i")
	cmd.Dir = config.TempDir
	// We have it in memory. Get a better way to do this!
	r, err = os.Open(path.Join(config.Urootpath, devcpio))
	if err != nil {
		log.Fatalf("%v", err)
	}

	// OK, at this point, we know we can run as root. And, we're going to create things
	// we can only remove as root. So, we'll have to remove the directory with
	// extreme measures.
	reallyRemoveDir := config.RemoveDir
	config.RemoveDir = false
	cmd.Stdin, cmd.Stderr, cmd.Stdout = r, os.Stderr, os.Stdout
	debug("Run %v @ %v", cmd, cmd.Dir)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("%v", err)
	}
	// Arrange to start init in the directory in a new namespace.
	// That should make all mounts go away when we're done.
	// On real kernels you can unshare without being root. Not on Linux.
	cmd = exec.Command("sudo", "unshare", "-m", "chroot", config.TempDir, "/init")
	cmd.Dir = config.TempDir
	cmd.Stdin, cmd.Stderr, cmd.Stdout = os.Stdin, os.Stderr, os.Stdout
	debug("Run %v @ %v", cmd, cmd.Dir)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Test failed, not removing %v: %v", config.TempDir, err)
		config.RemoveDir = false
	}
	config.RemoveDir = reallyRemoveDir
}
