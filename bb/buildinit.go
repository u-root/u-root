package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

func buildinit() {
	e := os.Environ()
	e = append(e, "CGO_ENABLED=0")
	e = append(e, "GO15VENDOREXPERIMENT=1")
	for i := range e {
		if e[i][0:6] == "GOPATH" {
			e[i] = e[i] + ":" + path.Join(config.Uroot, "bb/bbsh")
		}
	}
	cmd := exec.Command("go", "build", "-o", "init", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = config.Bbsh
	cmd.Env = e

	if err := cmd.Run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
