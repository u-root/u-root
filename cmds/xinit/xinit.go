package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func setup() error {
	if err := os.Symlink("/usr/local/bin/bash", "/bin/bash"); err != nil {
		return err
	}
	if err := os.Symlink("/lib/ld-linux-x86-64.so.2", "/lib64/ld-linux-x86-64.so.2"); err != nil {
		return err
	}
	go func() {
		cmd := exec.Command("Xfbdev")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("X11 startup: %v", err)
		}
	}()
	for {
		s, err := filepath.Glob("/tmp/.X*/X?")
		if err != nil {
			return err
		}
		if len(s) > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	for _, f := range []string{"flwm", "opera-12"} {
		log.Printf("Run %v", f)
		go func(n string) {
			cmd := exec.Command(n, "-display", ":0")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("%v: %v", n, err)
			}
		}(f)
	}
	cmd := exec.Command("aterm", "-display", ":0")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
}
