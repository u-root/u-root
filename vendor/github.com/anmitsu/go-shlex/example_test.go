package shlex_test

import (
	"fmt"
	"log"

	"github.com/anmitsu/go-shlex"
)

func ExampleSplit() {
	cmd := `cp -Rdp "file name" 'file name2' dir\ name`

	// Split of cmd with POSIX mode.
	words1, err := shlex.Split(cmd, true)
	if err != nil {
		log.Fatal(err)
	}
	// Split of cmd with Non-POSIX mode.
	words2, err := shlex.Split(cmd, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Source command:")
	fmt.Println(`cp -Rdp "file name" 'file name2' dir\ name`)
	fmt.Println()

	fmt.Println("POSIX mode:")
	for _, word := range words1 {
		fmt.Println(word)
	}
	fmt.Println()
	fmt.Println("Non-POSIX mode:")
	for _, word := range words2 {
		fmt.Println(word)
	}

	// Output:
	// Source command:
	// cp -Rdp "file name" 'file name2' dir\ name
	//
	// POSIX mode:
	// cp
	// -Rdp
	// file name
	// file name2
	// dir name
	//
	// Non-POSIX mode:
	// cp
	// -Rdp
	// "file name"
	// 'file name2'
	// dir\
	// name
}
