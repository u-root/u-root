package passphrase

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type RunTestCase struct {
	name  string
	essid string
	pass  string
	out   []byte
	err   error
}

var (
	essidStub     = "stub"
	shortPass     = "aaaaaaa"                                                          // 7 chars
	longPass      = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 64 chars
	validPass     = "aaaaaaaaaaaaaaaa"                                                 // 16 chars
	correctOutput = []byte(
		`network={
	ssid="stub"
	#psk="aaaaaaaaaaaaaaaa"
	psk=e270ba95a72c6d922e902f65dfa23315f7ba43b69debc75167254acd778f2fe9
}
`)

	runTestCases = []RunTestCase{
		{
			name:  "No essid",
			essid: "",
			pass:  validPass,
			out:   nil,
			err:   fmt.Errorf("essid cannot be empty"),
		},
		{
			name:  "pass length is less than 8 chars",
			essid: essidStub,
			pass:  shortPass,
			out:   nil,
			err:   fmt.Errorf("Passphrase must be 8..63 characters"),
		},
		{
			name:  "pass length is more than 63 chars",
			essid: essidStub,
			pass:  longPass,
			out:   nil,
			err:   fmt.Errorf("Passphrase must be 8..63 characters"),
		},
		{
			name:  "Correct Input",
			essid: essidStub,
			pass:  validPass,
			out:   correctOutput,
			err:   nil,
		},
	}
)

func TestRun(t *testing.T) {
	for _, test := range runTestCases {
		out, err := Run(test.essid, test.pass)
		if !reflect.DeepEqual(err, test.err) || bytes.Compare(out, test.out) != 0 {
			t.Logf("TEST %v", test.name)
			t.Errorf("\ngot:[%v, %v]\nwant:[%v, %v]", err, string(out), test.err, string(test.out))
		}
	}
}
