package tbs

import (
	"bytes"
	"os"
	"testing"
)

// Encodes a call to Getrandom() with a buffer of length zero, making this
// command an effective no-op.
var getRandomRawCommand = []byte{128, 1, 0, 0, 0, 12, 0, 0, 1, 123, 0, 0}

func getContext(t *testing.T) Context {
	ctx, err := CreateContext(TPMVersion20, IncludeTPM12|IncludeTPM20)
	if err != nil {
		t.Skipf("Skipping test as we couldn't access the TPM: %v", err)
	}
	return ctx
}

// Get the log by passing in progressively larger buffers
func TestGetLogLargeBuffer(t *testing.T) {
	ctx := getContext(t)
	defer ctx.Close()

	log := make([]byte, os.Getpagesize())
	for {
		logLen, err := ctx.GetTCGLog(log)
		if err == nil {
			if logLen == 0 {
				t.Fatal("Expected positive TCGLog length")
			}
			return
		}
		if err != ErrInsufficientBuffer {
			t.Fatalf("GetTCGLog failed: %v", err)
		}
		log = make([]byte, 2*len(log))
	}
}

// Get the log by passing in nil, checking the size, and then getting the log.
func TestGetLogWithNilSlice(t *testing.T) {
	ctx := getContext(t)
	defer ctx.Close()

	logLen, err := ctx.GetTCGLog(nil)
	if err != nil {
		t.Fatalf("First GetTCGLog failed: %v", err)
	}
	if logLen == 0 {
		t.Fatal("Expected positive TCGLog length")
	}

	log := make([]byte, logLen)
	if _, err := ctx.GetTCGLog(log); err != nil {
		t.Fatalf("Second GetTCGLog failed: %v", err)
	}
}

// SubmitCommand can handle a nil command buffer.
func TestSubmitCommandNilCommand(t *testing.T) {
	ctx := getContext(t)
	defer ctx.Close()

	rawResponse := make([]byte, os.Getpagesize())
	_, err := ctx.SubmitCommand(NormalPriority, nil, rawResponse)
	if err != ErrBadParameter {
		t.Fatalf("SubmitCommand failed with %v: expected ErrBadParameter", err)
	}
}

// SubmitCommand can handle a nil response buffer.
func TestSubmitCommandNilResponse(t *testing.T) {
	ctx := getContext(t)
	defer ctx.Close()

	_, err := ctx.SubmitCommand(NormalPriority, getRandomRawCommand, nil)
	if err != ErrInvalidOutputPointer {
		t.Fatalf("SubmitCommand failed with %v: expected ErrInvalidOutputPointer", err)
	}
}

// SubmitCommand can handle a response buffer that is shorter than necessary.
func TestSubmitCommandShortResponse(t *testing.T) {
	ctx := getContext(t)
	defer ctx.Close()

	rawResponse := make([]byte, 1)
	_, err := ctx.SubmitCommand(NormalPriority, getRandomRawCommand, rawResponse)
	if err != ErrInsufficientBuffer {
		t.Fatalf("SubmitCommand failed with %v: expected ErrInsufficientBuffer", err)
	}
}

// SubmitCommand can handle a response buffer that is longer than necessary.
func TestSubmitCommandLongResponse(t *testing.T) {
	ctx := getContext(t)
	defer ctx.Close()

	rawResponse := make([]byte, os.Getpagesize())
	responseLen, err := ctx.SubmitCommand(NormalPriority, getRandomRawCommand, rawResponse)
	if err != nil {
		t.Fatalf("SubmitCommand failed: %v", err)
	}
	rawResponse = rawResponse[:responseLen]

	// Expected response buffer for getRandomRawCommand
	expectedGetRandomRawResponse := []byte{128, 1, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(rawResponse, expectedGetRandomRawResponse) {
		t.Fatalf("Got response of %v, expected %v", rawResponse, expectedGetRandomRawResponse)
	}
}
