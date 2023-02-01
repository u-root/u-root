package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	osi, err := boot.PEImageFromFile(f)
	if err != nil {
		log.Fatal(err)
	}

	if err := osi.Execute(); err != nil {
		log.Fatal(err)
	}
}
