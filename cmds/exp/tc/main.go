// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/florianl/go-tc"
	trafficctl "github.com/u-root/u-root/pkg/tc"
)

func main() {
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}
	if err := run(os.Stdout, os.Args[1:], tctl); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func run(stdout io.Writer, args []string, tctl trafficctl.Tctl) error {
	cursor := 0
	want := []string{
		"qdisc",
		"class",
		"filter",
	}

	switch one(args[cursor], want) {
	case "qdisc":
		return runQdisc(stdout, args[cursor+1:], tctl)
	case "class":
		return runClass(stdout, args[cursor+1:], tctl)
	case "filter":
		return runFilter(stdout, args[cursor+1:], tctl)
	}

	return nil
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

func runQdisc(stdout io.Writer, args []string, tctl trafficctl.Tctl) error {
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
		qArgs, err = trafficctl.ParseQDiscArgs(os.Stdout, args[1:])
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
		return tctl.ChangeQDisc(stdout, qArgs)
	case "link":
		return tctl.LinkQDisc(stdout, qArgs)
	case "help":
		fmt.Fprintf(stdout, "%s", trafficctl.QdiscHelp)
	}

	return nil
}

func runClass(stdout io.Writer, args []string, tctl trafficctl.Tctl) error {
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
		fmt.Fprintf(stdout, "%s", trafficctl.ClassHelp)
		return nil
	}

	return nil
}

func runFilter(stdout io.Writer, args []string, tctl trafficctl.Tctl) error {
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
		fmt.Fprintf(stdout, "%s", trafficctl.Filterhelp)
	}

	return nil
}
