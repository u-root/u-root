// Copyright 2013-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestPrintf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple String",
			format:   "Hello, %s!",
			args:     []string{"World"},
			expected: "Hello, World!",
		},
		{
			name:     "Multiple Strings",
			format:   "Hello, %s and %s!",
			args:     []string{"Alice", "Bob"},
			expected: "Hello, Alice and Bob!",
		},
		{
			name:     "Integer",
			format:   "Number: %d",
			args:     []string{"42"},
			expected: "Number: 42",
		},
		{
			name:     "Float",
			format:   "Pi: %f",
			args:     []string{"3.14159"},
			expected: "Pi: 3.141590",
		},
		{
			name:     "Escape Sequences",
			format:   "Newline: \\nTab: \\tBackslash: \\\\",
			args:     []string{},
			expected: "Newline: \nTab: \tBackslash: \\",
		},
		{
			name:     "Percent Sign",
			format:   "Percent: %%",
			args:     []string{},
			expected: "Percent: %",
		},
		{
			name:    "Invalid Format",
			format:  "Invalid: %z",
			args:    []string{},
			wantErr: true,
		},
		{
			name:     "Insufficient Arguments",
			format:   "Hello, %s!",
			args:     []string{},
			expected: "Hello, !",
		},
		{
			name:     "Invalid Number",
			format:   "Number: %d",
			args:     []string{"invalid"},
			expected: "Number: 0",
		},
		{
			name:     "Hexadecimal Escape",
			format:   "Hex: \\x41",
			args:     []string{},
			expected: "Hex: A",
		},
		{
			name:     "Octal Escape",
			format:   "Octal: \\101",
			args:     []string{},
			expected: "Octal: A",
		},
		{
			name:     "Terminate Output",
			format:   "Terminate: \\c",
			args:     []string{},
			expected: "Terminate: ",
		},
		{
			name:     "Unescaped String",
			format:   "Unescaped: %b",
			args:     []string{"Hello\\nWorld"},
			expected: "Unescaped: Hello\nWorld",
		},
		{
			name:     "Variable Precision",
			format:   "Precision: %.2f",
			args:     []string{"3.14159"},
			expected: "Precision: 3.14",
		},
		{
			name:     "Bell Character",
			format:   "Bell: \\a",
			args:     []string{},
			expected: "Bell: \a",
		},
		{
			name:     "Backspace Character",
			format:   "Backspace: \\b",
			args:     []string{},
			expected: "Backspace: \b",
		},
		{
			name:     "Escape Character",
			format:   "Escape: \\e",
			args:     []string{},
			expected: "Escape: \033",
		},
		{
			name:     "Form Feed Character",
			format:   "Form Feed: \\f",
			args:     []string{},
			expected: "Form Feed: \f",
		},
		{
			name:     "Carriage Return Character",
			format:   "Carriage Return: \\r",
			args:     []string{},
			expected: "Carriage Return: \r",
		},
		{
			name:     "Vertical Tab Character",
			format:   "Vertical Tab: \\v",
			args:     []string{},
			expected: "Vertical Tab: \v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr for this test
			oldStderr := os.Stderr
			_, w, _ := os.Pipe()
			os.Stderr = w

			output, err := printf(tt.format, tt.args)

			// Restore stderr
			w.Close()
			os.Stderr = oldStderr

			if (err != nil) != tt.wantErr {
				t.Errorf("printf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if output != tt.expected {
				t.Errorf("printf() = %v, want %v", output, tt.expected)
			} else {
				t.Logf("Test %s passed", tt.name)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Capture standard output for testing
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the tests
	code := m.Run()

	// Restore standard output
	w.Close()
	os.Stdout = old

	// Print captured output
	out, _ := io.ReadAll(r)
	fmt.Printf("%s", out)

	os.Exit(code)
}
