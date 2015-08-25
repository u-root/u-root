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
	cmd := exec.Command("go", "build", "-o", "init", "init.go")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = config.Bb
	cmd.Env = e

	if err := cmd.Run(); err != nil {
		log.Fatalf("%v\n", err)
	}

	for i := range e {
		if e[i][0:6] == "GOPATH" {
			e[i] = e[i] + ":" + path.Join(config.Uroot, "src/bb/bbsh")
		}
	}
	cmd = exec.Command("go", "build", "-o", "sh", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = config.Bbcmdsh
	cmd.Env = e

	if err := cmd.Run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
