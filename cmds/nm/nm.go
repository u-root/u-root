// The nm command prints out information about binaries.
// At this point, there is very little flexibility in what it
// does and it only does PE files, but more will come.
package main

import (
	"debug/pe"
	"fmt"
	"os"
)

func syms(n string) error {
	f, err := pe.Open(n)
	if err != nil {
		return err
	}
	fmt.Printf("%v: %v", n, f)
	return nil
}

func main() {
	for _, n := range os.Args {
		if err := syms(n); err != nil {
			fmt.Printf("%v: %v\n", n, err)
		}
	}
}
