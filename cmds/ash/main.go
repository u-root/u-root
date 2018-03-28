package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/complete"
	"github.com/u-root/u-root/pkg/termios"
)

func main() {
	t, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := t.Raw()
	defer t.Set(r)
	for {
		p, err := complete.NewPathCompleter()
		if err != nil {
			log.Fatal(err)
		}
		c := complete.NewMultiCompleter(complete.NewStringCompleter([]string{"exit"}), p)
		l := complete.NewLineReader(c, t, t)
		s, err := l.ReadOne()
		if err != nil {
			log.Print(err)
			continue
		}
		if len(s) == 0 {
			continue
		}
		if s[0] == "exit" {
			break
		}
		// s[0] is either the match or what they typed so far.
		cmd := exec.Command(s[0])
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}
	}
}
