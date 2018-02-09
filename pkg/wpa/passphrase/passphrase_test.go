package passphrase

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type RunTestCase struct {
	name       string
	essid      string
	pass       string
	out        []byte
	err_exists bool
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
`) // indentation matters

	runTestCases = []RunTestCase{
		{
			name:       "No essid",
			essid:      "",
			pass:       validPass,
			out:        nil,
			err_exists: true,
		},
		{
			name:       "pass length is less than 8 chars",
			essid:      essidStub,
			pass:       shortPass,
			out:        nil,
			err_exists: true,
		},
		{
			name:       "pass length is more than 63 chars",
			essid:      essidStub,
			pass:       longPass,
			out:        nil,
			err_exists: true,
		},
		{
			name:       "Correct Input",
			essid:      essidStub,
			pass:       validPass,
			out:        correctOutput,
			err_exists: false,
		},
	}
)

func outEqualsExp(out []byte, exp []byte) bool {
	switch {
	case out == nil && exp == nil:
		return true
	case out != nil && exp != nil:
		return bytes.Compare(out, exp) == 0
	default:
		return false
	}
}

func TestRun(t *testing.T) {
	for _, test := range runTestCases {
		t.Logf("TEST %v", test.name)
		out, err := Run(test.essid, test.pass)
		if test.err_exists != testutil.ErrorExists(err) || !outEqualsExp(out, test.out) {
			testutil.PrintError(t, string(test.out), test.err_exists, string(out), err)
		}
	}
}
