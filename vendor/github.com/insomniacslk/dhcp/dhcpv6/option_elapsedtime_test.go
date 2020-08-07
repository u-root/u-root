package dhcpv6

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptElapsedTime(t *testing.T) {
	opt, err := parseOptElapsedTime([]byte{0xaa, 0xbb})
	if err != nil {
		t.Fatal(err)
	}
	if elapsedTime := opt.ElapsedTime; elapsedTime != 0xaabb*10*time.Millisecond {
		t.Fatalf("Invalid elapsed time. Expected 0xaabb, got %v", elapsedTime)
	}
}

func TestOptElapsedTimeToBytes(t *testing.T) {
	opt := OptElapsedTime(0)
	expected := []byte{0, 0}
	if toBytes := opt.ToBytes(); !bytes.Equal(expected, toBytes) {
		t.Fatalf("Invalid ToBytes output. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptElapsedTimeString(t *testing.T) {
	opt := OptElapsedTime(100 * time.Millisecond)
	expected := "ElapsedTime: 100ms"
	if optString := opt.String(); optString != expected {
		t.Fatalf("Invalid elapsed time string. Expected %v, got %v", expected, optString)
	}
}

func TestOptElapsedTimeParseInvalidOption(t *testing.T) {
	_, err := parseOptElapsedTime([]byte{0xaa})
	require.Error(t, err, "A short option should return an error")

	_, err = parseOptElapsedTime([]byte{0xaa, 0xbb, 0xcc})
	require.Error(t, err, "An option with too many bytes should return an error")
}
