package shlex_test

import (
	"fmt"
	"log"

	"github.com/anmitsu/go-shlex"
	flynn_shlex "github.com/flynn/go-shlex"
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

func ExampleSplit_compareFlynn() {
	cmd := `English and 日本語`

	// Split for github.com/flynn/go-shlex imported as flynn_shlex
	words_flynn, err1 := flynn_shlex.Split(cmd)

	// Split for github.com/anmitsu/go-shlex
	words_anmitsu, err2 := shlex.Split(cmd, true)

	fmt.Println("Source string:")
	fmt.Println(cmd)
	fmt.Println()

	fmt.Println("Result of github.com/flynn/go-shlex:")
	for _, word := range words_flynn {
		fmt.Println(word)
	}
	fmt.Println(err1.Error())

	fmt.Println()
	fmt.Println("Result of github.com/anmitsu/go-shlex:")
	for _, word := range words_anmitsu {
		fmt.Println(word)
	}
	if err2 != nil {
		fmt.Println(err2.Error())
	}

	// Output:
	// Source string:
	// English and 日本語
	//
	// Result of github.com/flynn/go-shlex:
	// English
	// and
	// Unknown rune: 26085
	//
	// Result of github.com/anmitsu/go-shlex:
	// English
	// and
	// 日本語
}
