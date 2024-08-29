// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/florianl/go-tc"
	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	trafficctl "github.com/u-root/u-root/pkg/tc"
)

const (
	DummyInterface0 = "eth0"
	DummyInterface1 = "eth1"
)

func TestVM(t *testing.T) {
	govmtest.Run(t, "tc integration",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/tc"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute*2),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", DummyInterface0)),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", DummyInterface1)),
		),
	)
}

func TestQDisc(t *testing.T) {
	guest.SkipIfNotInVM(t)

	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		t.Error(err)
	}
	defer rtnl.Close()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}

	argsStr := []string{
		"dev",
		DummyInterface0,
		"ingress",
	}

	var outbuf bytes.Buffer
	args, err := trafficctl.ParseQdiscArgs(&outbuf, argsStr)
	if err != nil {
		t.Errorf("ParseQDiscArgs() = %v, not nil", err)
	}

	if err := tctl.AddQdisc(&outbuf, args); err != nil {
		t.Errorf("AddQdisc() = %v, not nil", err)
	}

	if err := tctl.ReplaceQdisc(&outbuf, args); err != nil {
		t.Errorf("AddQdisc() = %v, not nil", err)
	}

	if err := tctl.ChangeQdisc(&outbuf, args); err != nil {
		t.Errorf("AddQdisc() = %v, not nil", err)
	}

	if err := tctl.DeleteQdisc(&outbuf, args); err != nil {
		t.Errorf("AddQdisc() = %v, not nil", err)
	}
}

func TestClass(t *testing.T) {
	guest.SkipIfNotInVM(t)

	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		t.Error(err)
	}
	defer rtnl.Close()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}

	htbQdiscStr := []string{
		"dev",
		DummyInterface0,
		"root",
		"handle",
		"1:",
		"htb",
		"default",
		"30",
	}

	var outbuf bytes.Buffer
	qargs, err := trafficctl.ParseQdiscArgs(&outbuf, htbQdiscStr)
	if err != nil {
		t.Errorf("ParseQDiscArgs() = %v, not nil", err)
	}

	if err := tctl.AddQdisc(&outbuf, qargs); err != nil {
		t.Errorf("AddQdisc() = %v, not nil", err)
	}

	htbClassAddStr := []string{
		"dev",
		DummyInterface0,
		"parent",
		"1:1",
		"classid",
		"1:10",
		"htb",
		"rate",
		"5mbit",
		"burst",
		"15k",
	}

	cargs, err := trafficctl.ParseClassArgs(&outbuf, htbClassAddStr)
	if err != nil {
		t.Errorf("ParseClassArgs() = %v, not nil", err)
	}

	if err := tctl.AddClass(&outbuf, cargs); err != nil {
		t.Errorf("AddClass() = %v, not nil", err)
	}

	htbClassDelStr := []string{
		"dev",
		DummyInterface0,
		"classid",
		"1:10",
	}

	cargs, err = trafficctl.ParseClassArgs(&outbuf, htbClassDelStr)
	if err != nil {
		t.Errorf("ParseClassArgs() = %v, not nil", err)
	}

	if err := tctl.DeleteClass(&outbuf, cargs); err != nil {
		t.Errorf("DelClass() = %v, not nil", err)
	}

	delQdiscStr := []string{
		"dev",
		DummyInterface0,
		"root",
	}

	delArgs, err := trafficctl.ParseClassArgs(&outbuf, delQdiscStr)
	if err != nil {
		t.Errorf("ParseClassArgs() = %v, not nil", err)
	}

	if err := tctl.DeleteQdisc(&outbuf, delArgs); err != nil {
		t.Errorf("DelClass() = %v, not nil", err)
	}
}

func TestFilter(t *testing.T) {
	guest.SkipIfNotInVM(t)

	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		t.Error(err)
	}
	defer rtnl.Close()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}

	htbQdiscStr := []string{
		"dev",
		DummyInterface0,
		"root",
		"handle",
		"1:",
		"htb",
		"default",
		"30",
	}

	var outbuf bytes.Buffer
	qargs, err := trafficctl.ParseQdiscArgs(&outbuf, htbQdiscStr)
	if err != nil {
		t.Errorf("ParseQDiscArgs() = %v, not nil", err)
	}

	if err := tctl.AddQdisc(&outbuf, qargs); err != nil {
		t.Errorf("AddQdisc() = %v, not nil", err)
	}

	filterArgsStr := []string{
		"dev",
		DummyInterface0,
		"parent",
		"1:",
		"protocol",
		"ip",
		"basic",
		"action",
		"drop",
	}

	fArgs, err := trafficctl.ParseFilterArgs(&outbuf, filterArgsStr)
	if err != nil {
		t.Errorf("ParseFilterArgs() = %v, not nil", err)
	}

	if err := tctl.AddFilter(&outbuf, fArgs); err != nil {
		t.Errorf("AddFilter() = %v, not nil", err)
	}

	delFilterStr := []string{
		"dev",
		DummyInterface0,
		"parent",
		"1:",
	}

	fArgs, err = trafficctl.ParseFilterArgs(&outbuf, delFilterStr)
	if err != nil {
		t.Errorf("ParseFilterArgs() = %v, not nil", err)
	}

	if err := tctl.DeleteFilter(&outbuf, fArgs); err != nil {
		t.Errorf("DeleteFilter() = %v, not nil", err)
	}
}
