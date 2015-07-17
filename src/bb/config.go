package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
)

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
	config.Uroot = os.Getenv("UROOT")
	if config.Uroot == "" {
		/* let's tr to guess. If we see bb.go, then we're in u-root/src/bb */
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
	config.Arch = getenv("GOARCH")
	config.Goroot = getenv("GOROOT")
	config.Gopath = getenv("GOPATH")
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
