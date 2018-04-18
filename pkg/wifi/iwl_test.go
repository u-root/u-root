// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wifi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/wpa/passphrase"
)

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

func TestGenerateConfig(t *testing.T) {
	for _, test := range generateConfigTestcases {
		out, err := generateConfig(test.args...)
		if !reflect.DeepEqual(err, test.err) || !bytes.Equal(out, test.exp) {
			t.Logf("TEST %v", test.name)
			fncCall := fmt.Sprintf("genrateConfig(%v)", test.args)
			t.Errorf("%s\ngot:[%v, %v]\nwant:[%v, %v]", fncCall, string(out), err, string(test.exp), test.err)
		}
	}
}

func TestCellRE(t *testing.T) {
	testcases := []struct {
		s   string
		exp bool
	}{
		{"blahblahblah\n   Cell 01:", true},
		{"blahblahblah\n   Cell 01: blah blah", true},
		{"\"Cell\"", false},
		{"\"blah blah Cell blah blah\"", false},
	}
	for _, test := range testcases {
		if out := cellRE.MatchString(test.s); out != test.exp {
			t.Errorf("%s\ngot:%v\nwant:%v", test.s, out, test.exp)
		}
	}
}

func TestEssidRE(t *testing.T) {
	testcases := []struct {
		s   string
		exp bool
	}{
		{"blahblahblah\n    ESSID:\"stub\"", true},
		{"blahblahblah\n    ESSID:\"stub\"\n", true},
		{"blahblahblah\n    ESSID:\"stub-stub\"", true},
		{"blahblahblah\n    ESSID:\"stub-stub\"\n", true},
		{"blah blah ESSID blah", false},
	}
	for _, test := range testcases {
		if out := essidRE.MatchString(test.s); out != test.exp {
			t.Errorf("%s\ngot:%v\nwant:%v", test.s, out, test.exp)
		}
	}
}

func TestEncKeyOptRE(t *testing.T) {
	testcases := []struct {
		s   string
		exp bool
	}{
		{"blahblahblah\n      Encryption key:on\n", true},
		{"blahblahblah\n      Encryption key:on", true},
		{"blahblahblah\n      Encryption key:off\n", true},
		{"blahblahblah\n      Encryption key:off", true},
		{"blah blah Encryption key blah blah", false},
		{"blah blah Encryption key:on  blah blah", false},
		{"blah blah Encryption key:off blah blah", false},
	}
	for _, test := range testcases {
		if out := encKeyOptRE.MatchString(test.s); out != test.exp {
			t.Errorf("%s\ngot:%v\nwant:%v", test.s, out, test.exp)
		}
	}
}

func TestWpa2RE(t *testing.T) {
	testcases := []struct {
		s   string
		exp bool
	}{
		{"blahblahblah\n            IE: IEEE 802.11i/WPA2 Version 1\n", true},
		{"blahblahblah\n            IE: IEEE 802.11i/WPA2 Version 1", true},
		{"blah blah IE: IEEE 802.11i/WPA2 Version 1", false},
	}
	for _, test := range testcases {
		if out := wpa2RE.MatchString(test.s); out != test.exp {
			t.Errorf("%s\ngot:%v\nwant:%v", test.s, out, test.exp)
		}
	}
}

func TestAuthSuitesRE(t *testing.T) {
	testcases := []struct {
		s   string
		exp bool
	}{
		{"blahblahblah\n            Authentication Suites (1) : 802.1x\n", true},
		{"blahblahblah\n            Authentication Suites (1) : 802.1x", true},
		{"blahblahblah\n            Authentication Suites (1) : PSK\n", true},
		{"blahblahblah\n            Authentication Suites (1) : PSK\n", true},
		{"blahblahblah\n            Authentication Suites (2) : blah, blah\n", true},
		{"blahblahblah\n            Authentication Suites (1) : other protocol\n", true},
		{"blahblahblah\n            Authentication Suites (1) : other protocol", true},
		{"blah blah Authentication Suites : blah blah", false},
	}
	for _, test := range testcases {
		if out := authSuitesRE.MatchString(test.s); out != test.exp {
			t.Errorf("%s\ngot:%v\nwant:%v", test.s, out, test.exp)
		}
	}
}

func TestParseIwlistOutput(t *testing.T) {
	var (
		o        []byte
		exp, out []Option
		err      error
	)

	// No WiFi present
	o = nil
	exp = nil
	out = parseIwlistOut(o)
	if !reflect.DeepEqual(out, exp) {
		t.Errorf("\ngot:[%v]\nwant:[%v]", out, exp)
	}

	// Only 1 WiFi present
	o = []byte(`
wlan0    Scan completed :
          Cell 01 - Address: 00:00:00:00:00:01
                    Channel:001
                    Frequency:5.58 GHz (Channel 001)
                    Quality=1/2  Signal level=-23 dBm  
                    Encryption key:on
                    ESSID:"stub-wpa-eap-1"
                    Bit Rates:36 Mb/s; 48 Mb/s; 54 Mb/s
                    Mode:Master
                    Extra:tsf=000000000000000000
                    Extra: Last beacon: 1260ms ago
                    IE: Unknown: 000000000000000000
                    IE: Unknown: 000000000000000000
                    IE: Unknown: 000000000000000000
                    IE: IEEE 802.11i/WPA2 Version 1
                        Group Cipher : CCMP
                        Pairwise Ciphers (1) : CCMP
                        Authentication Suites (1) : 802.1x
                    IE: Unknown: 000000000000000000
                    IE: Unknown: 000000000000000000
                    IE: Unknown: 000000000000000000
                    IE: Unknown: 000000000000000000
                    IE: Unknown: 000000000000000000
`)
	exp = []Option{
		{"stub-wpa-eap-1", WpaEap},
	}
	out = parseIwlistOut(o)
	if !reflect.DeepEqual(out, exp) {
		t.Errorf("\ngot:[%v]\nwant:[%v]", out, exp)
	}

	// Regular scenarios (many choices)
	exp = []Option{
		{"stub-wpa-eap-1", WpaEap},
		{"stub-rsa-1", NoEnc},
		{"stub-wpa-psk-1", WpaPsk},
		{"stub-rsa-2", NoEnc},
		{"stub-wpa-psk-2", WpaPsk},
	}
	o, err = ioutil.ReadFile("iwlistStubOutput.txt")
	if err != nil {
		t.Errorf("error reading iwlistStubOutput.txt: %v", err)
	}
	out = parseIwlistOut(o)
	if !reflect.DeepEqual(out, exp) {
		t.Errorf("\ngot:[%v]\nwant:[%v]", out, exp)
	}
}

func BenchmarkParseIwlistOutput(b *testing.B) {
	// Set Up
	o, err := ioutil.ReadFile("iwlistStubOutput.txt")
	if err != nil {
		b.Errorf("error reading iwlistStubOutput.txt: %v", err)
	}
	for i := 0; i < b.N; i++ {
		parseIwlistOut(o)
	}
}
