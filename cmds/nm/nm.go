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
	fmt.Printf("%v: FileHeader %v *OptionalHeader %v\n", n, f.FileHeader, f.OptionalHeader)
	fmt.Printf("%d Symbols: %v:\n", len(f.Symbols), f.Symbols)
	for _, s := range f.Symbols {
		fmt.Printf("\t%v\n", *s)
	}
	fmt.Printf("%d COFFSymbols: %v:\n", len(f.COFFSymbols), f.COFFSymbols)
	for _, s := range f.COFFSymbols {
		fmt.Printf("\t%v\n", s)
	}
	fmt.Printf("%d Sections: %v:\n", len(f.Sections), f.Sections)
	for _, s := range f.Sections {
		fmt.Printf("\t%v\n", *s)
	}
	return nil
}

func main() {
	for _, n := range os.Args {
		if err := syms(n); err != nil {
			fmt.Printf("%v: %v\n", n, err)
		}
	}
}
