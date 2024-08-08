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
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func run(stdout io.Writer, args []string) error {
	cursor := 0
	want := []string{
		"qdisc",
		"class",
		"filter",
	}

	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		return err
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}

	switch one(args[cursor], want) {
	case "qdisc":
		return runQdisc(args[cursor+1:], tctl, stdout)
	case "class":
		return runClass(args[cursor+1:], tctl, stdout)
	case "filter":
		return runFilter(args[cursor+1:], tctl, stdout)
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

func runQdisc(args []string, tctl *trafficctl.Trafficctl, stdout io.Writer) error {
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

	qArgs := &trafficctl.QArgs{}
	var err error
	if len(args[1:]) > 1 {
		qArgs, err = trafficctl.ParseQDiscArgs(args[1:], os.Stdout)
		if err != nil {
			return err
		}
	}

	switch one(args[cursor], want) {
	case "show", "list":
		return tctl.ShowQdisc(qArgs, stdout)
	case "add":
		return tctl.AddQdisc(qArgs, stdout)
	case "del":
		return tctl.DelQdisc(qArgs, stdout)
	case "replace":
		return tctl.ReplaceQdisc(qArgs, stdout)
	case "change":
		return tctl.ChangeQDisc(qArgs, stdout)
	case "link":
		return tctl.LinkQDisc(qArgs, stdout)
	case "help":
		trafficctl.PrintQdiscHelp(stdout)
	}

	return nil
}

func runClass(args []string, tctl *trafficctl.Trafficctl, stdout io.Writer) error {
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

	cArgs := &trafficctl.CArgs{}
	var err error
	if len(args[1:]) > 1 {
		cArgs, err = trafficctl.ParseClassArgs(args[1:], stdout)
		if err != nil {
			return err
		}
	}

	switch one(args[cursor], want) {
	case "show", "list":
		return tctl.ShowClass(cArgs, stdout)
	case "add":
		return tctl.AddClass(cArgs, stdout)
	case "delete", "del":
		return tctl.DeleteClass(cArgs, stdout)
	case "change":
		return tctl.ChangeClass(cArgs, stdout)
	case "replace":
		return tctl.ReplaceClass(cArgs, stdout)
	case "help":
		trafficctl.PrintClassHelp(stdout)
		return nil
	}

	return nil
}

func runFilter(args []string, tctl *trafficctl.Trafficctl, stdout io.Writer) error {
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
		fArgs, err = trafficctl.ParseFilterArgs(args[1:], stdout)
		if err != nil {
			return err
		}
	}

	switch one(args[cursor], want) {
	case "show", "list":
		return tctl.ShowFilter(fArgs, stdout)
	case "add":
		return tctl.AddFilter(fArgs, stdout)
	case "del":
		return tctl.DeleteFilter(fArgs, stdout)
	case "change":
		return tctl.ChangeFilter(fArgs, stdout)
	case "replace":
		return tctl.ReplaceFilter(fArgs, stdout)
	case "get":
		return tctl.GetFilter(fArgs, stdout)
	case "help":
		trafficctl.PrintFilterHelp(stdout)
	}

	return nil
}
