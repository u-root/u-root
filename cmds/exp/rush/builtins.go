// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// builtins: cd, exit, pwd, echo
// extras: rushinfo

package main

import (
	"runtime"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Error messages
var errCdUsage = errors.New("usage: cd <directory-to-change-to>")
var errExitUsage = errors.New("usage: exit <|0-255|>")

// Add builtins to the shell
func init() {
	_ = addBuiltIn("rushinfo", infocmd)
	_ = addBuiltIn("cd", cd)
	_ = addBuiltIn("exit", exit)
	_ = addBuiltIn("pwd", pwd)
	_ = addBuiltIn("echo", echo)
}

// rushinfo command: print info about the shell and its environment
func infocmd(c *Command) error {
	_, err := fmt.Fprintf(c.Stdout, "%s %s %s %q: builtins %v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH, os.Args, builtins)
	return err
}

// cd command: change directory
func cd(c *Command) error {
	if len(c.argv) != 1 {
		return errCdUsage
	}

	err := os.Chdir(c.argv[0])
	return err
}

// exit command: exit the shell with a given exit code
func exit(c *Command) error {
	code := 0 // Default exit code is 0

	if len(c.argv) > 1 {
		return errExitUsage
	}
	if len(c.argv) == 1 {
		// Attempt to convert argument to integer
		var err error
		code, err = strconv.Atoi(c.argv[0])
		if err != nil || code < 0 || code > 255 {
			return fmt.Errorf("exit: invalid exit code")
		}
	}

	// Exit with the specified code
	os.Exit(code)
	return nil
}

// pwd command: print the current working directory
func pwd(c *Command) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: %v", err)
	}
	fmt.Println(dir)
	return nil
}

// echo command: print the arguments with space separation
func echo(c *Command) error {
	// Join the arguments with spaces and print them
	fmt.Println(strings.Join(c.argv, " "))
	return nil
}
