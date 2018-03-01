// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/wpa/passphrase"
)

type WifiErrorTestCase struct {
	name   string
	args   []string
	expect string
}

type GenerateConfigTestCase struct {
	name string
	args []string
	exp  []byte
	err  error
}

var (
	EssidStub       = "stub"
	IdStub          = "stub"
	PassStub        = "123456789"
	BadWpaPskPass   = "123"
	expWpaPsk, _    = passphrase.Run(EssidStub, PassStub)
	_, expWpaPskErr = passphrase.Run(EssidStub, BadWpaPskPass)

	errorTestcases = []WifiErrorTestCase{
		{
			name:   "More elements than needed",
			args:   []string{"a", "a", "a", "a"},
			expect: "Usage",
		},
		{
			name:   "Flags, More elements than needed",
			args:   []string{"-i=123", "a", "a", "a", "a"},
			expect: "Usage",
		},
	}

	generateConfigTestcases = []GenerateConfigTestCase{
		{
			name: "No Pass Phrase",
			args: []string{EssidStub},
			exp:  []byte(fmt.Sprintf(nopassphrase, EssidStub)),
			err:  nil,
		},
		{
			name: "WPA-PSK",
			args: []string{EssidStub, PassStub},
			exp:  expWpaPsk,
			err:  nil,
		},
		{
			name: "WPA-EAP",
			args: []string{EssidStub, PassStub, IdStub},
			exp:  []byte(fmt.Sprintf(eap, EssidStub, IdStub, PassStub)),
			err:  nil,
		},
		{
			name: "WPA-PSK Error",
			args: []string{EssidStub, BadWpaPskPass},
			exp:  nil,
			err:  fmt.Errorf("essid: %v, pass: %v : %v", EssidStub, BadWpaPskPass, expWpaPskErr),
		},
		{
			name: "Invalid Args Length Error",
			args: nil,
			exp:  nil,
			err:  fmt.Errorf("generateConfig needs 1, 2, or 3 args"),
		},
	}
)

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestWifiErrors(t *testing.T) {
	// Set up
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Tests
	for _, test := range errorTestcases {
		c := exec.Command(execPath, test.args...)
		_, e, _ := run(c)
		if !strings.Contains(e, test.expect) {
			t.Logf("TEST %v", test.name)
			execStatement := fmt.Sprintf("exec(wifi %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
			t.Errorf("%s\ngot:%s\nwant:%s", execStatement, e, test.expect)
		}
	}
}

func TestWifiGenerateConfig(t *testing.T) {
	for _, test := range generateConfigTestcases {
		out, err := generateConfig(test.args...)
		if !reflect.DeepEqual(err, test.err) || !bytes.Equal(out, test.exp) {
			t.Logf("TEST %v", test.name)
			fncCall := fmt.Sprintf("genrateConfig(%v)", test.args)
			t.Errorf("%s\ngot:[%v, %v]\nwant:[%v, %v]", fncCall, string(out), err, string(test.exp), test.err)

		}
	}
}
