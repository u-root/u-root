package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBootFileURL(t *testing.T) {
	expected := "https://insomniac.slackware.it"
	opt, err := parseOptBootFileURL([]byte(expected))
	if err != nil {
		t.Fatal(err)
	}
	if string(opt) != expected {
		t.Fatalf("Invalid boot file URL. Expected %v, got %v", expected, opt)
	}
	require.Contains(t, opt.String(), "https://insomniac.slackware.it", "String() should contain the correct BootFileUrl output")
}

func TestOptBootFileURLToBytes(t *testing.T) {
	urlString := "https://insomniac.slackware.it"
	opt := OptBootFileURL(urlString)
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, []byte(urlString)) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", urlString, toBytes)
	}
}
