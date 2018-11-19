// Copyright 2016-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// Simple Test trying execute the ps
// If no errors returns, it's okay
func TestPsExecution(t *testing.T) {
	t.Logf("TestPsExecution")
	d, err := ioutil.TempDir("", fmt.Sprintf("ps"))
	if err != nil {
		t.Fatal(err)
	}
	//defer os.Rmdirall(d)
	var tests = []struct {
		n     string
		pid   string
		files map[string]string
		o     string
		err   error
	}{
		{n: "missing files", pid: "1", files: map[string]string{"stat": "bad status file"}},
		{n: "one process", pid: "1", files: map[string]string{"stat": "bad status file",
			"status": `Name:	systemd
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
FDSize:	128
`,
			"cmdline": "/sbin/init"},
			o: "PID PGRP SID TTY    STAT        TIME  COMMAND \n1     ?    file    00:00:00   status \n",
		},
		// Fix things up
		{n: "correct pid 1", pid: "1", files: map[string]string{"stat": "1 (systemd) S 0 1 1 0 -1 4194560 82923 51272244 88 3457 153 671 103226 39563 20 0 1 0 2 230821888 2325 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 1 0 0 69 0 0 0 0 0 0 0 0 0 0",
			"status": `Name:	systemd
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
FDSize:	128
`,
			"cmdline": "/sbin/init"},
			o: "PID PGRP SID TTY    STAT         TIME  COMMAND \n1 1 1 ?    S        00:00:08  systemd \n",
		},
		{n: "second process", pid: "1996", files: map[string]string{"stat": "1996 (dnsmasq) S 1 1995 1995 0 -1 4194624 64 0 0 0 1 10 0 0 20 0 1 0 1208 51163136 91 18446744073709551615 1 1 0 0 0 0 0 4096 92675 0 0 0 17 2 0 0 0 0 0 0 0 0 0 0 0 0 0",

			"status": `Name:	dnsmasq
Umask:	0022
Uid:	110	110	110	110
`,
			"cmdline": "/usr/sbin/dnsmasq\000--conf-file=/var/lib/libvirt/dnsmasq/default.conf\000--leasefile-ro\000--dhcp-script=/usr/lib/libvirt/libvirt_leaseshelper\000"},
			o: " PID PGRP  SID TTY    STAT         TIME  COMMAND \n   1    1    1 ?    S        00:00:08  systemd \n1996 1995 1995 ?    S        00:00:00  dnsmasq \n",
		},
		{n: "nethost process", pid: "srv/1996", files: map[string]string{"stat": "1996 (dnsmasq) S 1 1995 1995 0 -1 4194624 64 0 0 0 1 10 0 0 20 0 1 0 1208 51163136 91 18446744073709551615 1 1 0 0 0 0 0 4096 92675 0 0 0 17 2 0 0 0 0 0 0 0 0 0 0 0 0 0",

			"status": `Name:	dnsmasq
Umask:	0022
Uid:	110	110	110	110
`,
			"cmdline": "/usr/sbin/dnsmasq\000--conf-file=/var/lib/libvirt/dnsmasq/default.conf\000--leasefile-ro\000--dhcp-script=/usr/lib/libvirt/libvirt_leaseshelper\000"},
			o: "     PID     PGRP      SID TTY    STAT         TIME  COMMAND \n       1        1        1 ?    S        00:00:08  systemd \n    1996     1995     1995 ?    S        00:00:00  dnsmasq \nsrv/1996     1995     1995 ?    S        00:00:00  dnsmasq \n",
		},
	}

	for _, tt := range tests {
		pd := filepath.Join(d, tt.pid)
		t.Logf("Create %v", pd)
		if err := os.MkdirAll(pd, 0777); err != nil {
			t.Fatalf("Make proc dir: %v", err)
		}
		for n, f := range tt.files {
			procf := filepath.Join(pd, n)
			t.Logf("Write %v", procf)
			if err := ioutil.WriteFile(procf, []byte(f), 0666); err != nil {
				t.Fatal(err)
			}
		}
		c := testutil.Command(t, "aux")
		psp := fmt.Sprintf("UROOT_PSPATH=%s:%s", d, filepath.Join(d, "srv"))
		c.Env = append(c.Env, psp)
		o, err := c.CombinedOutput()
		t.Logf("%s: %s %v", tt.n, string(o), err)
		if string(o) != tt.o {
			t.Errorf("%v: got %q, want %q", tt.n, string(o), tt.o)
		}
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("%v: got %v, want %v", tt.n, err, tt.err)
		}
	}

}

// Test Parsing of stat
func TestParse(t *testing.T) {
	var tests = []struct {
		name string
		p    *Process
		out  string
		err  error
	}{
		{
			name: "no status file",
			p: &Process{
				stat: "1 (systemd) S 0 1 1 0 -1 4194560 45535 23809816 88 2870 76 378 35944 9972 20 0 1 0 2 230821888 2325 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 2 0 0 69 0 0 0 0 0 0 0 0 0 0",
			},
			err: fmt.Errorf("no Uid string in "),
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
`},
			err: nil,
		},
	}

	flags.all = true
	for _, tt := range tests {
		t.Logf("%v", tt.name)
		err := tt.p.Parse()
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Check %v: got %v, want %v", tt.p.stat, err, tt.err)
			continue
		}
		if err != nil {
			continue
		}
		pT := NewProcessTable()
		pT.table = []*Process{tt.p}
		pT.mProc = tt.p
		var b bytes.Buffer
		err = ps(pT, &b)
		t.Logf("ps out is %s", b.String())
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("%v: got %v, want %v", pT, err, tt.err)
			continue
		}
		if b.String() != tt.out {
			t.Errorf("%s: got %q, want %q", tt.p.stat, b.String(), tt.out)
			continue
		}

	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
