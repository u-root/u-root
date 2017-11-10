package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	mountPoint := "/go/src/github.com/u-root/u-root"
	goVersion := "1.9.2"

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	args := []string{
		"run",
		"--privileged",
		"-v", fmt.Sprintf("%s:%s", cwd, mountPoint),
		"-t",
		"-e", "USER=root",
		"-w", mountPoint,
		fmt.Sprintf("golang:%s", goVersion),
		"./travis.sh",
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
