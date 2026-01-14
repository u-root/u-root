package main

import (
	_ "embed"
	"fmt"
)

//go:embed foo/*.txt
var s string

func main() {
	fmt.Printf(s)
}
