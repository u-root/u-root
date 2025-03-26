// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/brctl"
	"golang.org/x/sys/unix"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name    string
		argv    []string
		expErr  error
		skipVM  bool
		forceVM bool
	}{
		{
			name: "help",
			argv: []string{"help"},
		},
		{
			name:   "addbr",
			argv:   []string{"addbr", "bridge0"},
			skipVM: true,
		},
		{
			name:    "addbr_VM",
			argv:    []string{"addbr", "bridge0"},
			expErr:  unix.EEXIST,
			forceVM: true,
		},
		{
			name:   "addbr_fewArgs",
			argv:   []string{"addbr", "bridge0", "bridge1"},
			expErr: errFewArgs,
		},
		{
			name:   "delbr",
			argv:   []string{"delbr", "bridge0"},
			skipVM: true,
		},
		{
			name:    "delbr_VM",
			argv:    []string{"delbr", "bridge0"},
			expErr:  unix.ENXIO,
			forceVM: true,
		},
		{
			name:   "delbr_fewArgs",
			argv:   []string{"delbr", "bridge0", "bridge1"},
			expErr: errFewArgs,
		},
		{
			name:   "addif",
			argv:   []string{"addif", "bridge0", "eth0"},
			expErr: unix.ENODEV,
			skipVM: true,
		},
		{
			name:    "addif_VM",
			argv:    []string{"addif", "bridge0", "eth0"},
			expErr:  unix.ENODEV,
			forceVM: true,
		},
		{
			name:   "addif_fewArgs",
			argv:   []string{"addif", "bridge0", "eth0", "eth1"},
			expErr: errFewArgs,
		},
		// see https://github.com/hugelgupf/vmtest/issues/130
		// {
		// 	name:   "delif",
		// 	argv:   []string{"delif", "bridge0", "eth0"},
		// 	expErr: unix.ENODEV,
		// 	skipVM: true,
		// },
		{
			name:    "delif_VM",
			argv:    []string{"delif", "bridge0", "eth0"},
			expErr:  unix.ENODEV,
			forceVM: true,
		},
		{
			name:   "delif_fewArgs",
			argv:   []string{"delif", "bridge0", "eth0", "eth1"},
			expErr: errFewArgs,
		},
		{
			name:   "show",
			argv:   []string{"show"},
			expErr: nil,
		},
		{
			name:   "showstp",
			argv:   []string{"showstp", "bridge0"},
			expErr: os.ErrNotExist,
		},
		{
			name:   "showstp",
			argv:   []string{"showstp"},
			expErr: errFewArgs,
		},
		{
			name:   "showmacs",
			argv:   []string{"showmacs", "eth0"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "showmacs_fewArgs",
			argv:   []string{"showmacs", "eth0", "eth1"},
			expErr: errFewArgs,
		},
		{
			name:   "setageing",
			argv:   []string{"setageing", "bridge0", "10"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "setageing",
			argv:   []string{"setageing", "bridge0", "10", "garbage"},
			expErr: errFewArgs,
		},
		{
			name:   "stp",
			argv:   []string{"stp", "bridge0", "10"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "stp_fewArgs",
			argv:   []string{"stp", "bridge0", "10", "garbage"},
			expErr: errFewArgs,
		},
		{
			name:   "setbridgeprio",
			argv:   []string{"setbridgeprio", "bridge0", "10"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "setbridgeprio_fewArgs",
			argv:   []string{"setbridgeprio", "bridge0", "10", "garbage"},
			expErr: errFewArgs,
		},
		{
			name:   "setfd",
			argv:   []string{"setfd", "bridge0", "10"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "setfd_fewArgs",
			argv:   []string{"setfd", "bridge0", "10", "garbage"},
			expErr: errFewArgs,
		},
		{
			name:   "sethello",
			argv:   []string{"sethello", "bridge0", "10"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "sethello_fewArgs",
			argv:   []string{"sethello", "bridge0", "10", "garbage"},
			expErr: errFewArgs,
		},
		{
			name:   "setmaxage",
			argv:   []string{"setmaxage", "bridge0", "10"},
			expErr: brctl.ErrBridgeNotExist,
		},
		{
			name:   "setmaxage_fewArgs",
			argv:   []string{"setmaxage", "bridge0", "10", "garbage"},
			expErr: errFewArgs,
		},
		// see https://github.com/hugelgupf/vmtest/issues/130
		// {
		// 	name:   "setpathcost",
		// 	argv:   []string{"setpathcost", "bridge0", "eth0", "10"},
		// 	expErr: brctl.ErrPortNotExist,
		// 	skipVM: true,
		// },
		{
			name:   "setpathcost_fewArgs",
			argv:   []string{"setpathcost", "bridge0", "eth0", "10", "garbage"},
			expErr: errFewArgs,
		},
		// see https://github.com/hugelgupf/vmtest/issues/130
		// {
		// 	name:   "setportprio",
		// 	argv:   []string{"setportprio", "bridge0", "eth0", "10"},
		// 	expErr: brctl.ErrPortNotExist,
		// 	skipVM: true,
		// },
		{
			name:   "setportprio_fewArgs",
			argv:   []string{"setportprio", "bridge0", "eth0", "10", "garbage"},
			expErr: errFewArgs,
		},
		// see https://github.com/hugelgupf/vmtest/issues/130
		// {
		// 	name:   "hairpin",
		// 	argv:   []string{"hairpin", "bridge0", "eth0", "enable"},
		// 	expErr: brctl.ErrPortNotExist,
		// 	skipVM: true,
		// },
		{
			name:   "hairpin_fewArgs",
			argv:   []string{"hairpin", "bridge0", "eth0", "enable", "garbage"},
			expErr: errFewArgs,
		},
		{
			name:   "garbage_command",
			argv:   []string{"garbage_command"},
			expErr: errInvalidCmd,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipVM {
				guest.SkipIfInVM(t)
			}

			if tt.forceVM {
				guest.SkipIfNotInVM(t)
			}

			var outbuf bytes.Buffer
			if err := run(&outbuf, tt.argv); !errors.Is(err, tt.expErr) {
				t.Errorf("run(): %v, not: %v", err, tt.expErr)
			}
		})
	}
}
