package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
)

var outputPath = flag.String("o", "bb", "Path to busybox binary")

func main() {
	flag.Parse()

	env := golang.Default()
	if env.CgoEnabled {
		log.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)

	pkgs := flag.Args()
	var err error
	if len(pkgs) == 0 {
		pkgs, err = uroot.DefaultPackageImports(env)
	} else {
		pkgs, err = uroot.ResolvePackagePaths(env, pkgs)
	}
	if err != nil {
		log.Fatal(err)
	}

	o, err := filepath.Abs(*outputPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := uroot.BuildBusybox(env, pkgs, o); err != nil {
		log.Fatal(err)
	}
}
