// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

// Session manages a raw terminal session with automatic cleanup.
// The session must be closed with Close() or Restore() to restore
// the terminal to its original state.
type Session struct {
	file     *os.File
	oldState *term.State
	restored bool
	mu       sync.Mutex
}

// NewSession creates a new terminal session in raw mode.
// The session must be closed with Close() or Restore() to restore
// the terminal to its original state.
func NewSession(f *os.File) (*Session, error) {
	if !term.IsTerminal(int(f.Fd())) {
		return nil, fmt.Errorf("not a terminal")
	}

	oldState, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to set raw mode: %w", err)
	}

	s := &Session{
		file:     f,
		oldState: oldState,
	}

	return s, nil
}

// Restore restores the terminal to its original state.
// This function is idempotent - calling it multiple times is safe.
func (s *Session) Restore() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.restored {
		return nil
	}

	var err error
	if s.oldState != nil {
		err = term.Restore(int(s.file.Fd()), s.oldState)
	}
	s.restored = true
	return err
}

// Close is an alias for Restore for io.Closer compatibility.
func (s *Session) Close() error {
	return s.Restore()
}

// File returns the underlying file.
func (s *Session) File() *os.File {
	return s.file
}

// Read implements io.Reader by reading from the underlying file.
func (s *Session) Read(p []byte) (n int, err error) {
	return s.file.Read(p)
}

// MakeCanonical temporarily restores canonical mode.
// It returns a function to re-enter raw mode.
// This is useful for operations that need line-buffered input,
// such as timeout prompts.
func (s *Session) MakeCanonical() (restore func() error, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.restored {
		return nil, fmt.Errorf("session already restored")
	}

	// Restore to canonical (original state)
	err = term.Restore(int(s.file.Fd()), s.oldState)
	if err != nil {
		return nil, err
	}

	// Return function to re-enter raw mode
	// Note: We do NOT update s.oldState here - it must remain as the
	// original canonical state so that Close() can properly restore it.
	return func() error {
		s.mu.Lock()
		defer s.mu.Unlock()

		_, err := term.MakeRaw(int(s.file.Fd()))
		return err
	}, nil
}

// Ensure Session implements io.ReadCloser
var _ io.ReadCloser = (*Session)(nil)
