package passphrase

import (
	"bytes"
	"fmt"
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
`) // indentation matters

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

// Helper function to craft the message
func craftPrintMsg(err error, out []byte) string {
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("Error Status: %v\n", err))
	msg.WriteString("Output:\n")
	msg.Write(out)
	return msg.String()
}

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

func errEqualsExp(err error, exp error) bool {
	switch {
	case err == nil && exp == nil:
		return true
	case err != nil && exp != nil:
		return err.Error() == exp.Error()
	default:
		return false
	}
}

func TestRun(t *testing.T) {
	for _, test := range runTestCases {
		out, err := Run(test.essid, test.pass)
		if !errEqualsExp(err, test.err) || !outEqualsExp(out, test.out) {
			t.Logf("TEST %v", test.name)
			actualMsg := craftPrintMsg(err, out)
			expectMsg := craftPrintMsg(test.err, test.out)
			t.Errorf("\ngot:\n%s\n\nwant:\n%s", actualMsg, expectMsg)
		}
	}
}
