package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bb"
)

func run() {
	name := filepath.Base(os.Args[0])
	if err := bb.Run(name); err != nil {
		log.Fatalf("%q: %v", name, err)
	}
}

func main() {
	run()
}

func init() {
	bb.Register("bb", bb.Noop, func() {
		if len(os.Args) <= 1 {
			log.Fatalf("You need to specify which command to invoke.")
		}
		os.Args = os.Args[1:]
		run()
	})
}
