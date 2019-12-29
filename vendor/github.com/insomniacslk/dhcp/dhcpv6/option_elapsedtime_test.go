package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptElapsedTime(t *testing.T) {
	opt, err := ParseOptElapsedTime([]byte{0xaa, 0xbb})
	if err != nil {
		t.Fatal(err)
	}
	if elapsedTime := opt.ElapsedTime; elapsedTime != 0xaabb {
		t.Fatalf("Invalid elapsed time. Expected 0xaabb, got %v", elapsedTime)
	}
}

func TestOptElapsedTimeToBytes(t *testing.T) {
	opt := OptElapsedTime{}
	expected := []byte{0, 0}
	if toBytes := opt.ToBytes(); !bytes.Equal(expected, toBytes) {
		t.Fatalf("Invalid ToBytes output. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptElapsedTimeSetGetElapsedTime(t *testing.T) {
	opt := OptElapsedTime{}
	opt.ElapsedTime = 10
	if elapsedTime := opt.ElapsedTime; elapsedTime != 10 {
		t.Fatalf("Invalid elapsed time. Expected 10, got %v", elapsedTime)
	}
}

func TestOptElapsedTimeString(t *testing.T) {
	opt := OptElapsedTime{}
	opt.ElapsedTime = 10
	expected := "OptElapsedTime{elapsedtime=10}"
	if optString := opt.String(); optString != expected {
		t.Fatalf("Invalid elapsed time string. Expected %v, got %v", expected, optString)
	}
}

func TestOptElapsedTimeParseInvalidOption(t *testing.T) {
	_, err := ParseOptElapsedTime([]byte{0xaa})
	require.Error(t, err, "A short option should return an error")

	_, err = ParseOptElapsedTime([]byte{0xaa, 0xbb, 0xcc})
	require.Error(t, err, "An option with too many bytes should return an error")
}
