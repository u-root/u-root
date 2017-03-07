// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run lua code or call a lua function
//
// Synopsis:
//     lua CMD [ARG]...
//     luacall FUNCTION [ARG]...
//
// Description:
//     lua assemble the named args into a string and pass the string to the interpreter
//     luacall will assembled the args and string into this: FUNCTION(ARG...) and pass that string to the interpreter.
//
// Bugs:
//     luacall is a lot to type. What's a better name?
package main

import (
	"strings"

	"github.com/Shopify/go-lua"
)

var (
	l = lua.NewState()
)

func init() {
	lua.OpenLibraries(l)
	addBuiltIn("lua", runlua)
	addBuiltIn("calllua", calllua)
}

func runlua(c *Command) error {
	return lua.DoString(l, strings.Join(c.argv, " "))
}

func calllua(c *Command) error {
	return lua.DoString(l, c.argv[0]+"("+strings.Join(c.argv[1:], ",")+")")
}
