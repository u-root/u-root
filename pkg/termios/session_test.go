// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"os"
	"testing"
)

func TestNewSession_NotATerminal(t *testing.T) {
	// Create a regular file (not a terminal)
	tmpfile, err := os.CreateTemp("", "session_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// NewSession should fail for non-terminal files
	_, err = NewSession(tmpfile)
	if err == nil {
		t.Error("Expected error for non-terminal file, got nil")
	}
}

func TestSession_RestoreIdempotent(t *testing.T) {
	// We can't create a real terminal in tests, so we'll test
	// the idempotent behavior with a mock
	s := &Session{
		file:     nil,
		oldState: nil,
		restored: false,
	}

	// First restore should succeed
	err := s.Restore()
	if err != nil {
		t.Errorf("First Restore() failed: %v", err)
	}
	if !s.restored {
		t.Error("Session should be marked as restored")
	}

	// Second restore should be a no-op
	err = s.Restore()
	if err != nil {
		t.Errorf("Second Restore() failed: %v", err)
	}
	if !s.restored {
		t.Error("Session should still be marked as restored")
	}
}

func TestSession_CloseIsRestore(t *testing.T) {
	// Test that Close() is an alias for Restore()
	s := &Session{
		file:     nil,
		oldState: nil,
		restored: false,
	}

	err := s.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
	if !s.restored {
		t.Error("Session should be marked as restored after Close()")
	}
}

func TestSession_MakeCanonicalFailsWhenRestored(t *testing.T) {
	s := &Session{
		file:     nil,
		oldState: nil,
		restored: true,
	}

	_, err := s.MakeCanonical()
	if err == nil {
		t.Error("MakeCanonical() should fail on restored session")
	}
}

func TestSession_File(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "session_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	s := &Session{
		file: tmpfile,
	}

	if s.File() != tmpfile {
		t.Error("File() should return the underlying file")
	}
}

func TestSession_Read(t *testing.T) {
	// Create a temporary file with some content
	tmpfile, err := os.CreateTemp("", "session_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := []byte("test data")
	if _, err := tmpfile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek: %v", err)
	}

	s := &Session{
		file: tmpfile,
	}

	buf := make([]byte, len(testData))
	n, err := s.Read(buf)
	if err != nil {
		t.Errorf("Read() failed: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Read() returned %d bytes, expected %d", n, len(testData))
	}
	if string(buf) != string(testData) {
		t.Errorf("Read() returned %q, expected %q", buf, testData)
	}
}
