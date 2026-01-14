// getppid is a package that has one external dependency.
package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println(unix.Getppid())
}
