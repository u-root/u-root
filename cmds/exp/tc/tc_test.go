// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	trafficctl "github.com/u-root/u-root/pkg/tc"
)

type DummyTctl struct{}

func (d *DummyTctl) ShowQdisc(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) AddQdisc(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) DeleteQdisc(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) ReplaceQdisc(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) ChangeQdisc(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) LinkQdisc(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) ShowClass(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) AddClass(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) DeleteClass(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) ReplaceClass(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) ChangeClass(io.Writer, *trafficctl.Args) error {
	return nil
}

func (d *DummyTctl) ShowFilter(io.Writer, *trafficctl.FArgs) error {
	return nil
}

func (d *DummyTctl) AddFilter(io.Writer, *trafficctl.FArgs) error {
	return nil
}

func (d *DummyTctl) DeleteFilter(io.Writer, *trafficctl.FArgs) error {
	return nil
}

func (d *DummyTctl) ReplaceFilter(io.Writer, *trafficctl.FArgs) error {
	return nil
}

func (d *DummyTctl) ChangeFilter(io.Writer, *trafficctl.FArgs) error {
	return nil
}

func (d *DummyTctl) GetFilter(io.Writer, *trafficctl.FArgs) error {
	return nil
}

func TestRun(t *testing.T) {
	d := &DummyTctl{}

	for _, tt := range []struct {
		name   string
		args   []string
		err    error
		outStr string
	}{
		{
			name: "Show help",
			args: []string{
				"help",
			},
			outStr: cmdHelp,
		},
		{
			name:   "Show help no args",
			args:   nil,
			outStr: cmdHelp,
		},
		{
			name: "Show Qdisc",
			args: []string{
				"qdisc",
				"show",
			},
		},
		{
			name: "Add Qdisc",
			args: []string{
				"qdisc",
				"add",
				"dev",
				"eth0",
			},
		},
		{
			name: "Delete Qdisc",
			args: []string{
				"qdisc",
				"del",
			},
		},
		{
			name: "replace Qdisc",
			args: []string{
				"qdisc",
				"replace",
			},
		},
		{
			name: "change Qdisc",
			args: []string{
				"qdisc",
				"change",
			},
		},
		{
			name: "link Qdisc",
			args: []string{
				"qdisc",
				"link",
			},
		},
		{
			name: "help Qdisc",
			args: []string{
				"qdisc",
				"help",
			},
			outStr: trafficctl.QdiscHelp,
		},
		{
			name: "Show Class",
			args: []string{
				"class",
				"show",
			},
		},
		{
			name: "Add Class",
			args: []string{
				"class",
				"add",
				"dev",
				"eth0",
			},
		},
		{
			name: "Delete Class",
			args: []string{
				"class",
				"del",
			},
		},
		{
			name: "replace Class",
			args: []string{
				"class",
				"replace",
			},
		},
		{
			name: "change Class",
			args: []string{
				"class",
				"change",
			},
		},
		{
			name: "help Class",
			args: []string{
				"class",
				"help",
			},
			outStr: trafficctl.ClassHelp,
		},
		{
			name: "Show Filter",
			args: []string{
				"filter",
				"show",
			},
		},
		{
			name: "Add Filter",
			args: []string{
				"filter",
				"add",
				"dev",
				"eth0",
			},
		},
		{
			name: "Delete Filter",
			args: []string{
				"filter",
				"del",
			},
		},
		{
			name: "replace Filter",
			args: []string{
				"filter",
				"replace",
			},
		},
		{
			name: "change Filter",
			args: []string{
				"filter",
				"change",
			},
		},
		{
			name: "help Filter",
			args: []string{
				"filter",
				"help",
			},
			outStr: trafficctl.FilterHelp,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var outbuf bytes.Buffer
			if err := run(&outbuf, d, tt.args); !errors.Is(err, tt.err) {
				t.Errorf("run() = %v", err)
			}

			if tt.outStr != "" {
				if tt.outStr != outbuf.String() {
					t.Errorf("output: \n%s\nnot equal expectation: \n%s\n", outbuf.String(), tt.outStr)
				}
			}
		})
	}
}
