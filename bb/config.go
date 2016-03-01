package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
}

// TODO: put this in the uroot package
// It's annoying asking them to set lots of things. So let's try to figure it out.
func guessgoroot() {
	config.Goroot = os.Getenv("GOROOT")
	if config.Goroot != "" {
		log.Printf("Using %v as GOROOT from environment variable", config.Goroot)
		config.Gosrcroot = path.Dir(config.Goroot)
		return
	}
	log.Print("Goroot is not set, trying to find a go binary")
	p := os.Getenv("PATH")
	paths := strings.Split(p, ":")
	for _, v := range paths {
		g := path.Join(v, "go")
		if _, err := os.Stat(g); err == nil {
			config.Goroot = path.Dir(path.Dir(v))
			config.Gosrcroot = path.Dir(config.Goroot)
			log.Printf("Guessing that goroot is %v", config.Goroot)
			return
		}
	}
	log.Printf("GOROOT is not set and can't find a go binary in %v", p)
	config.Fail = true
}

func guessgopath() {
	defer func() {
		config.Gosrcroot = path.Dir(config.Goroot)
	}()
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		config.Gopath = path.Clean(gopath)
		return
	}
	// It's a good chance they're running this from the u-root source directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("GOPATH was not set and I can't get the wd: %v", err)
		config.Fail = true
		return
	}
	// walk up the cwd until we find a u-root entry. See if src/cmds/init/init.go exists.
	for c := cwd; c != "/"; c = path.Dir(c) {
		if path.Base(c) != "u-root" {
			continue
		}
		check := path.Join(c, "src/cmds/init/init.go")
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

func guessuroot() {
	config.Uroot = os.Getenv("UROOT")
	if config.Uroot == "" {
		/* let's try to guess. If we see bb.go, then we're in u-root/src/bb */
		if _, err := os.Stat("bb.go"); err == nil {
			dir := path.Dir(config.Cwd)
			config.Uroot = path.Dir(dir)
		} else if _, err := os.Stat("src/bb/bb.go"); err == nil {
			// Maybe they're at top level? If there is a srb/bb/bb.go, that's it.
			config.Uroot = config.Cwd
		} else {
			log.Fatalf("UROOT was not set and I don't seem to be in u-root/src/bb/bb.go or u-root")
		}
	}

}

func doConfig() {
	var err error
	flag.BoolVar(&config.Debug, "d", false, "Debugging")
	flag.Parse()
	if config.Debug {
		debug = debugPrint
	}
	if config.Cwd, err = os.Getwd(); err != nil {
		log.Fatalf("Getwd: %v", err)
	}
	guessgoroot()
	guessgopath()
	guessuroot()
	config.Arch = getenv("GOARCH", "amd64")
	if config.Fail {
		os.Exit(1)
	}
	config.Gosrcroot = path.Dir(config.Goroot)
	config.Goos = "linux"
	config.TempDir, err = ioutil.TempDir("", "u-root")
	config.Go = ""
	if err != nil {
		log.Fatalf("%v", err)
	}
	config.Bbsh = path.Join(config.Cwd, "bbsh")
	os.RemoveAll(config.Bbsh)
	config.Args = flag.Args()
	if len(config.Args) == 0 {
		config.Args = defaultCmd
	}

}
