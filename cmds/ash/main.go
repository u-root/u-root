package main

import (
	"log"

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
		var data [1]byte
		for {
			if _, err := t.Read(data[:]); err != nil {
				log.Fatal(err)
			}
			// Log the error but it may be transient.
			if _, err := t.Write(data[:]); err != nil {
				log.Fatal(err)
			}
		}

	}
}
