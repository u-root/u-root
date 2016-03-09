package main

import (
	"bytes"
	"debug/elf"
	"flag"
	"fmt"
	"github.com/u-root/u-root/uroot"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
)

type copyfiles struct {
	dir  string
	spec string
}

const (
	devcpio   = "scripts/dev.cpio"
	urootPath = "src/github.com/u-root/u-root"
	// huge suckage here. the 'old' usage is going away but it not gone yet. Just suck in old6a for now.
	// I don't want to revive the 'letter' stuff.
	// This has gotten kind of ugly. But [0] is source, [1] is dest, and [2..] is the list.
	// FIXME. this is ugly.
	goList = `{{.Goroot}}
go
{{.Go}}
pkg/include
src
VERSION.cache
misc
pkg/tool/{{.Goos}}_{{.Arch}}/compile
pkg/tool/{{.Goos}}_{{.Arch}}/link
pkg/tool/{{.Goos}}_{{.Arch}}/asm
pkg/tool/{{.Goos}}_{{.Arch}}/old6a`
	urootList = `{{.Gopath}}

src/github.com/u-root/u-root/cmds
src/github.com/u-root/u-root/uroot
src/github.com/u-root/u-root/vendor`
)

var (
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
		TmpDir          string
		UseExistingInit bool
	}
	letter = map[string]string{
		"amd64": "6",
		"386":   "8",
		"arm":   "5",
		"ppc":   "9",
	}
	// the whitelist is a list of u-root tools that we feel
	// can replace existing tools. It is, sadly, a very short
	// list at present.
	whitelist = []string{"date"}
)

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
	if config.Debug {
		log.Printf("cpiop: from %v, to %v, :%v:\n", n[0], n[1], n[2:])
	}

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
		log.Printf("Run %v @ %v", cmd, cmd.Dir)
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("%v\n", err)
	}

	for _, v := range n[2:] {
		if config.Debug {
			log.Printf("%v\n", v)
		}
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
	if config.Debug {
		log.Printf("Done sending files to external")
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("%v\n", err)
	}
	if config.Debug {
		log.Printf("External cpio is done")
	}
	return nil
}

func sanity() {
	binGo := path.Join(config.Goroot, "bin/go")
	log.Printf("check %v as the go binary", binGo)
	_, err := os.Stat(binGo)
	if err == nil {
		config.Go = "bin/go"
	}
	log.Printf("%v exists, but check go/bin/OS_ARCH too", config.Go)
	// but does the one in go/bin/OS_ARCH exist too?
	archgo := fmt.Sprintf("bin/%s_%s/go", config.Goos, config.Arch)
	OsArchBinGo := path.Join(config.Goroot, archgo)
	log.Printf("check %v as the go binary", OsArchBinGo)
	_, err = os.Stat(OsArchBinGo)
	if err == nil {
		config.Go = archgo
		binGo = OsArchBinGo
	}
	log.Printf("Using %v as the go command", binGo)
	if config.Go == "" {
		log.Fatalf("Can't find a go binary! Is GOROOT set correctly?")
	}
	f, err := elf.Open(binGo)
	if err != nil {
		log.Fatalf("%v is not an ELF file; don't know what to do", binGo)
	}
	ds := f.SectionByType(elf.SHT_DYNAMIC)
	if ds == nil {
		return
	}

	log.Printf("*************************************************************************")
	log.Printf("U-root requires a staticically built go command. %v is dynamic.", binGo)
	log.Printf("This is ok; u-root is all source, but we have to rebuild  the go binary")
	log.Printf("Another way to  fix this:\ncd %v/src\nexport CGO_ENABLED=0\nGOARCH=%v ./make.bash", config.Goroot, config.Arch)
	log.Printf("*************************************************************************")

	goBin := path.Join(config.TempDir, "go/bin/go")
	cmd := exec.Command("go", "build", "-x", "-a", "-installsuffix", "cgo", "-ldflags", "'-s'", "-o", goBin)
	cmd.Dir = path.Join(config.Goroot, "src/cmd/go")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Building statically linked go tool info %v: %v, %v\n", goBin, string(o), err)
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
		if _, err := os.Stat(g); err == nil {
			config.Goroot = path.Dir(path.Dir(v))
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

// sad news. If I concat the Go cpio with the other cpios, for reasons I don't understand,
// the kernel can't unpack it. Don't know why, don't care. Need to create one giant cpio and unpack that.
// It's not size related: if the go archive is first or in the middle it still fails.
func main() {
	flag.BoolVar(&config.Debug, "d", false, "Debugging")
	flag.BoolVar(&config.TestChroot, "test", false, "test the directory by chrooting to it")
	flag.BoolVar(&config.UseExistingInit, "useinit", false, "If there is an existing init, don't replace it")
	flag.BoolVar(&config.RemoveDir, "removedir", true, "remove the directory when done -- cleared if test fails")
	flag.StringVar(&config.InitialCpio, "cpio", "", "An initial cpio image to build on")
	flag.StringVar(&config.TmpDir, "tmpdir", "", "tmpdir to use instead of ioutil.TempDir")
	flag.Parse()
	var err error
	guessgoarch()
	config.Go = ""
	config.Goos = "linux"
	guessgoroot()
	guessgopath()
	if config.Fail {
		log.Fatal("Setup failed")
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

	// sanity checking: do $GROOT/bin/go, and some basic source files exist?
	sanity()

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
			log.Printf("Run %v @ %v", cmd, cmd.Dir)
			cmd.Stdout = os.Stdout
		}

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

	if config.Debug {
		log.Printf("Done all cpio operations")
	}

	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// First create the archive and put the device cpio in it.
	dev, err := ioutil.ReadFile(path.Join(config.Urootpath, devcpio))
	if err != nil {
		log.Fatal("%v %v\n", dev, err)
	}

	if config.Debug {
		log.Printf("Creating initramf file")
	}

	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", config.Goos, config.Arch)
	if err := ioutil.WriteFile(oname, dev, 0600); err != nil {
		log.Fatal("%v\n", err)
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
	if config.Debug {
		log.Printf("Run %v @ %v", cmd, cmd.Dir)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := lsr(config.TempDir, w); err != nil {
		log.Fatal("%v\n", err)
	}
	w.Close()
	if config.Debug {
		log.Printf("Finished sending file list for initramfs cpio")
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("%v\n", err)
	}
	if config.Debug {
		log.Printf("cpio for initramfs is done")
	}
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
	if config.Debug {
		log.Printf("Run %v @ %v", cmd, cmd.Dir)
	}
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
	if config.Debug {
		log.Printf("Run %v @ %v", cmd, cmd.Dir)
	}
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Test failed, not removing %v: %v", config.TempDir, err)
		config.RemoveDir = false
	}
	config.RemoveDir = reallyRemoveDir
}
