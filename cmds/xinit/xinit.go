package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func x11(n string, args ...string) error {
	cmd := exec.Command(n, args...)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("X11 start %v %v: %v", n, args, err)
	}
	return nil
}

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
	for _, f := range []string{"wingo", "flwm", "chrome"} {
		log.Printf("Run %v", f)
		go x11(f)
	}

	// we block on the aterm. When the aterm exits, we do too.
	return x11("aterm")
}

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
}
