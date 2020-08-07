package tpm2

import (
	"reflect"
	"testing"

	"github.com/google/go-tpm/tpmutil"
)

func TestError(t *testing.T) {
	tests := []struct {
		response tpmutil.ResponseCode
		expected error
	}{
		{0x501, VendorError{Code: 0x501}},
		{0x922, Warning{Code: RCRetry}},
		{0x100, Error{Code: RCInitialize}},
		{0xfc1, ParameterError{Code: RCAsymmetric, Parameter: RCF}},
		{0x7a3, HandleError{Code: RCExpired, Handle: RC7}},
		{0xfa2, SessionError{Code: RCBadAuth, Session: RC7}},
	}

	for _, test := range tests {
		err := decodeResponse(test.response)
		if !reflect.DeepEqual(err, test.expected) {
			t.Fatalf("decodeResponse(0x%x) = %#v, want %#v", test.response, err, test.expected)
		}
	}
}

// nil ReadWriter handle causes tpmutil.RunCommand to return an error.
func TestRunCommandErr(t *testing.T) {
	if _, err := runCommand(nil, TagSessions, CmdSign); err == nil {
		t.Error("runCommand returned nil error on error from tpmutil.RunCommand")
	}
}
