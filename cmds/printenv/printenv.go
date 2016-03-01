package main

import (
	"fmt"
	"io"
	"os"
)

func printenv(w io.Writer) {
	e := os.Environ()

	for _, v := range e {
		fmt.Fprintf(w, "%v\n", v)
	}
}

func main() {
	printenv(os.Stdout)
}
