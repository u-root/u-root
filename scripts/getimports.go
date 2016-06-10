package main

import (
	"go/build"
	"flag"
	"github.com/u-root/u-root/uroot"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	urootPath = "src/github.com/u-root/u-root"
)

var (
	config struct {
		Goroot          string
		Godot           string
		Arch            string
		Goos            string
		Gopath          string
		Urootpath       string
		Go              string
		Debug           bool
		Fail            bool
	}
	Imports        map[string]bool
	debug     = nodebug
)

func nodebug(string, ...interface{}) {}

// TODO: put this in a common place. Or, merge ramfs and getimports, which almost makes sense.
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
		return
	}
	log.Print("Goroot is not set, trying to find a go binary")
	p := os.Getenv("PATH")
	paths := strings.Split(p, ":")
	for _, v := range paths {
		g := path.Join(v, "go")
		if _, err := os.Stat(g); err == nil {
			config.Goroot = path.Dir(path.Dir(v))
			log.Printf("Guessing that goroot is %v from $PATH", config.Goroot)
			return
		}
	}
	log.Printf("GOROOT is not set and can't find a go binary in %v", p)
	config.Fail = true
}

func guessgopath() {
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

// listPackages computes the packages needed by all the u-root commands.
// It walks the cmds/ directory, and for each directory, adds all packages used.
func listPackages() error {
	var pkgList []string

	err := filepath.Walk(path.Join(config.Urootpath, "cmds"), func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf(" WALK FAIL%v: %v\n", name, err)
			// That's ok, sometimes things are not there.
			return filepath.SkipDir
		}
		if fi.Name() == "cmds" {
			return nil
		}
		if !fi.IsDir() {
			return nil
		}
		pkgList = append(pkgList, path.Join("github.com/u-root/u-root/cmds", fi.Name()))
		return filepath.SkipDir
	})
	if err != nil {
		debug("Walking cmds/: %v\n", err)
	}
	// It would be nice to run go list -json with lots of package names but it produces invalid JSON.
	// It produces a stream thatis {}{}{} at the top level and the decoders don't like that.
	// TODO: fix it later. Maybe use template after all. For now this is more than adequate.
	for _, v := range pkgList {
		p, err := build.Default.Import(v, "", 0)
		if err != nil {
			log.Printf("Error on %v: %v", v, err)
			continue
		}
		debug("v %v  Goroot %v %v ", v, p.Goroot, p.Imports)
		if err != nil {
			log.Fatalf("%v", err)
		}
		debug("cmd p is %v", p)
		for _,v := range p.Imports {
			debug("Check %v", v)
			if ! p.Goroot {
				Imports[v] = true
			}
		}
	}

	return nil
}

func main() {
	flag.BoolVar(&config.Debug, "d", false, "Debugging")
	flag.Parse()
	if config.Debug {
		debug = log.Printf
	}

	Imports = make(map[string]bool)
	guessgoarch()
	config.Go = ""
	config.Goos = "linux"
	guessgoroot()
	guessgopath()
	if config.Fail {
		log.Fatal("Setup failed")
	}

	if err := listPackages(); err != nil {
		log.Fatalf("%v", err)
	}

	for i := range Imports {
		debug("Dep: %v", i)
		_, err := build.Default.Import(i, "", build.FindOnly)
		if err == nil {
			debug("Package %v exists, not getting it", i)
			continue
		}
		// TODO: consider making a substitution possible if we ever need to lock down a
		// package version/tag/whatever.
		debug("go get %v", i)
		cmd := exec.Command("go", "get", "-a", i)
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		debug("Run %v @ %v", cmd, cmd.Dir)
		j, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Go get failed: err %v, output \n%v\n", err, string(j))
		}
		debug("We got %v", i)
	}

}
