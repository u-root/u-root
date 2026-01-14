// Package exthello has one external dependency.
package exthello

import (
	hello1 "github.com/u-root/gobusybox/test/diamonddep/mod1/pkg/hello"
	hello3 "github.com/u-root/gobusybox/test/diamonddep/mod3/pkg/hello"
)

func Hello() string {
	return "test/diamonddep/mod2/exthello: " + hello1.Hello() + " and " + hello3.Hello()
}
