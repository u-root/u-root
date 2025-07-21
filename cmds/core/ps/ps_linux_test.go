// Copyright 2016-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPs(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []string
		all     bool
		every   bool
		x       bool
		nSidTty bool
		aux     bool
		want    []string
		wantErr string
	}{
		{
			name: "aux",
			args: []string{"aux"},
			x:    true,
			want: []string{"PID", "PGRP", "SID", "TTY", "STAT", "TIME", "COMMAND"},
		},
		{
			name: "flag x",
			x:    true,
			want: []string{"PID", "TTY", "STAT", "TIME", "COMMAND"},
		},
		{
			name: "switch case 2 default case",
			x:    false,
			want: []string{"PID", "TTY", "TIME", "CMD"},
		},
		{
			name: "usage()",
			args: []string{"test"},
			x:    true,
			want: []string{},
		},
		{
			name:    "flag nSidTty",
			nSidTty: true,
			want:    []string{"PID", "TTY", "TIME", "CMD"},
		},
	} {
		all = tt.all
		every = tt.every
		x = tt.x
		nSidTty = tt.nSidTty
		aux = tt.aux
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			if err := ps(buf, tt.args...); err != nil {
				t.Errorf("ps() = %q, want: %q", err.Error(), tt.wantErr)
			} else {
				for _, want := range tt.want {
					if !strings.Contains(buf.String(), want) {
						t.Errorf("ps() = %q, want to contain: %q", buf.String(), tt.want)
					}
				}
			}
		})
	}
}

// Test Parsing of stat
func TestParse(t *testing.T) {
	for _, tt := range []struct {
		name string
		p    *Process
		out  string
		err  string
	}{
		{
			name: "no status file",
			p: &Process{
				stat: "1 (systemd) S 0 1 1 0 -1 4194560 45535 23809816 88 2870 76 378 35944 9972 20 0 1 0 2 230821888 2325 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 2 0 0 69 0 0 0 0 0 0 0 0 0 0",
			},
			err: "no Uid string in ",
		},
		{
			name: "Valid output",
			out:  "PID TTY        TIME CMD     \n1 ?    00:00:04 systemd \n",
			p: &Process{
				stat: "1 (systemd) S 0 1 1 0 -1 4194560 45535 23809816 88 2870 76 378 35944 9972 20 0 1 0 2 230821888 2325 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 2 0 0 69 0 0 0 0 0 0 0 0 0 0",
				status: `Name:	systemd
Umask:	0000
State:	S (sleeping)
Tgid:	1
Ngid:	0
Pid:	1
PPid:	0
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
FDSize:	128
Groups:
NStgid:	1
NSpid:	1
NSpgid:	1
NSsid:	1
VmPeak:	  290768 kB
VmSize:	  225412 kB
VmLck:	       0 kB
VmPin:	       0 kB
VmHWM:	    9308 kB
VmRSS:	    9300 kB
RssAnon:	    2524 kB
RssFile:	    6776 kB
RssShmem:	       0 kB
VmData:	   18696 kB
VmStk:	     132 kB
VmExe:	    1336 kB
VmLib:	   10008 kB
VmPTE:	     204 kB
VmSwap:	       0 kB
HugetlbPages:	       0 kB
CoreDumping:	0
Threads:	1
SigQ:	0/31573
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	7be3c0fe28014a03
SigIgn:	0000000000001000
SigCgt:	00000001800004ec
CapInh:	0000000000000000
CapPrm:	0000003fffffffff
CapEff:	0000003fffffffff
CapBnd:	0000003fffffffff
CapAmb:	0000000000000000
NoNewPrivs:	0
Seccomp:	0
Speculation_Store_Bypass:	thread vulnerable
Cpus_allowed:	ffffffff,ffffffff,ffffffff,ffffffff
Cpus_allowed_list:	0-127
Mems_allowed:	00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	10168
nonvoluntary_ctxt_switches:	3746
`,
			},
			err: "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Parse(); err != nil {
				if err.Error() != tt.err {
					t.Errorf("Parse() = %q, want: %q", err, tt.err)
				}
			}
		})
	}
}
