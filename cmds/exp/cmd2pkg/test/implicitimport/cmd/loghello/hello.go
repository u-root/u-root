package main

import (
	"github.com/u-root/gobusybox/test/implicitimport/pkg/defaultlog"
)

// Default returns a *log.Logger, but "log" is not imported in this package.
//
// The busybox build must add "log" to the import statements.
var l = defaultlog.Default()

// Call it twice to make sure we do not add the new import twice.
var l2 = defaultlog.Default()

func main() {
	l.Printf("Log Hello")
}
