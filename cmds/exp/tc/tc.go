// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/florianl/go-tc"
	trafficctl "github.com/u-root/u-root/pkg/tc"

	// To build the dependencies of this package with TinyGo, we need to include
	// the cpuid package, since tinygo does not support the asm code in the
	// cpuid package. The cpuid package will use the tinygo bridge to get the
	// CPU information. For further information see
	// github.com/u-root/cpuid/cpuid_amd64_tinygo_bridge.go
	_ "github.com/u-root/cpuid"
)

var cmdHelp = `Usage:	tc OBJECT { COMMAND | help }
where  OBJECT := { qdisc | class | filter }
`

func main() {
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}
	if err := run(os.Stdout, tctl, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(stdout io.Writer, tctl trafficctl.Tctl, args []string) error {
	if len(args) == 0 {
		fmt.Fprint(stdout, cmdHelp)
		return nil
	}

	cursor := 0
	want := []string{
		"qdisc",
		"class",
		"filter",
		"help",
	}

	var err error

	switch one(args[cursor], want) {
	case "qdisc":
		err = runQdisc(stdout, tctl, args[cursor+1:])
	case "class":
		err = runClass(stdout, tctl, args[cursor+1:])
	case "filter":
		err = runFilter(stdout, tctl, args[cursor+1:])
	case "help":
		fmt.Fprint(stdout, cmdHelp)
	default:
		fmt.Fprint(stdout, cmdHelp)
	}
	if errors.Is(err, trafficctl.ErrExitAfterHelp) {
		return nil
	}

	return err
}

func one(cmd string, cmds []string) string {
	var x, n int
	for i, v := range cmds {
		if strings.HasPrefix(v, cmd) {
			n++
			x = i
		}
	}
	if n == 1 {
		return cmds[x]
	}
	return ""
}

func runQdisc(stdout io.Writer, tctl trafficctl.Tctl, args []string) error {
	cursor := 0
	want := []string{
		"show",
		"list",
		"add",
		"del",
		"replace",
		"change",
		"link",
		"help",
	}

	qArgs := &trafficctl.Args{}
	var err error
	if len(args[1:]) > 1 {
		qArgs, err = trafficctl.ParseQdiscArgs(os.Stdout, args[1:])
		if err != nil {
			return err
		}
	}

	switch one(args[cursor], want) {
	case "show", "list":
		return tctl.ShowQdisc(stdout, qArgs)
	case "add":
		return tctl.AddQdisc(stdout, qArgs)
	case "del":
		return tctl.DeleteQdisc(stdout, qArgs)
	case "replace":
		return tctl.ReplaceQdisc(stdout, qArgs)
	case "change":
		return tctl.ChangeQdisc(stdout, qArgs)
	case "help":
		fmt.Fprint(stdout, trafficctl.QdiscHelp)
	}

	return nil
}

func runClass(stdout io.Writer, tctl trafficctl.Tctl, args []string) error {
	cursor := 0
	want := []string{
		"show",
		"list",
		"add",
		"del",
		"change",
		"replace",
		"help",
	}

	cArgs := &trafficctl.Args{}
	var err error
	if len(args[1:]) > 1 {
		cArgs, err = trafficctl.ParseClassArgs(stdout, args[1:])
		if err != nil {
			return err
		}
	}

	switch one(args[cursor], want) {
	case "show", "list":
		return tctl.ShowClass(stdout, cArgs)
	case "add":
		return tctl.AddClass(stdout, cArgs)
	case "delete", "del":
		return tctl.DeleteClass(stdout, cArgs)
	case "change":
		return tctl.ChangeClass(stdout, cArgs)
	case "replace":
		return tctl.ReplaceClass(stdout, cArgs)
	case "help":
		fmt.Fprint(stdout, trafficctl.ClassHelp)
		return nil
	}

	return nil
}

func runFilter(stdout io.Writer, tctl trafficctl.Tctl, args []string) error {
	cursor := 0
	want := []string{
		"show",
		"list",
		"add",
		"del",
		"change",
		"replace",
		"get",
		"help",
	}

	fArgs := &trafficctl.FArgs{}
	var err error
	if len(args[1:]) > 1 {
		fArgs, err = trafficctl.ParseFilterArgs(stdout, args[1:])
		if err != nil {
			return err
		}
	}

	switch one(args[cursor], want) {
	case "show", "list":
		return tctl.ShowFilter(stdout, fArgs)
	case "add":
		return tctl.AddFilter(stdout, fArgs)
	case "del":
		return tctl.DeleteFilter(stdout, fArgs)
	case "change":
		return tctl.ChangeFilter(stdout, fArgs)
	case "replace":
		return tctl.ReplaceFilter(stdout, fArgs)
	case "get":
		return tctl.GetFilter(stdout, fArgs)
	case "help":
		fmt.Fprint(stdout, trafficctl.FilterHelp)
	}

	return nil
}
