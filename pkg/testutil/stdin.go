// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

// FakeStdin implements io.Reader and returns predefined list of answers,
// suitable for mocking standard input.
type FakeStdin struct {
	answers    []string
	pos        int
	overflowed bool
}

// NewFakeStdin creates new FakeStdin value with given answers.
func NewFakeStdin(answers ...string) *FakeStdin {
	fs := FakeStdin{answers: make([]string, len(answers))}
	for i, a := range answers {
		fs.answers[i] = a + "\n"
	}
	return &fs
}

// Read answers one by one and keep record whether the stdin
// has overflowed.
func (fs *FakeStdin) Read(p []byte) (int, error) {
	if len(fs.answers) <= fs.pos {
		fs.overflowed = true
		fs.pos = 0
	}
	n := copy(p, fs.answers[fs.pos])
	fs.pos++
	return n, nil
}

// Count returns how many answers have been read.
//
// This counter overflows and never returns value bigger than the count of
// given answers.
func (fs *FakeStdin) Count() int {
	return fs.pos
}

// Overflowed reports whether more reads happened than expected.
func (fs *FakeStdin) Overflowed() bool {
	return fs.overflowed
}
