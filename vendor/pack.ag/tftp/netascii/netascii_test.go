// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package netascii // import "pack.ag/tftp/netascii"

import (
	"bytes"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping non-windows tests")
	}
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "A string with no encoding",
			expected: "A string with no encoding",
		},
		{
			input:    "A string \r\x00 with \r\n encoding",
			expected: "A string \r with \n encoding",
		},
		{
			input:    "A string with incorrect \r encoding",
			expected: "A string with incorrect \r encoding",
		},
	}

	for _, c := range cases {
		reader := NewReader(strings.NewReader(c.input))

		result, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}

		if string(result) != c.expected {
			t.Errorf("Expected %q to be %q, but it was %q", c.input, c.expected, result)
		}
	}
}

func TestReader_windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows only tests")
	}
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "A string with no encoding",
			expected: "A string with no encoding",
		},
		{
			input:    "A string \r\x00 with \r\n encoding",
			expected: "A string \r with \r\n encoding",
		},
		{
			input:    "A string with incorrect \r encoding",
			expected: "A string with incorrect \r encoding",
		},
	}

	for _, c := range cases {
		reader := NewReader(strings.NewReader(c.input))

		result, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}

		if string(result) != c.expected {
			t.Errorf("Expected %q to be %q, but it was %q", c.input, c.expected, result)
		}
	}
}

func TestWriter(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "A string with no encoding",
			expected: "A string with no encoding",
		},
		{
			input:    "A string \r with \n encoding",
			expected: "A string \r\x00 with \r\n encoding",
		},
		{
			input:    "A string \r\x00 with existing \r\n encoding",
			expected: "A string \r\x00 with existing \r\n encoding",
		},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		writer := NewWriter(&buf)

		_, err := io.Copy(writer, strings.NewReader(c.input))
		if err != nil {
			t.Fatal(err)
		}

		if err := writer.Flush(); err != nil {
			t.Fatal(err)
		}

		if result := buf.String(); result != c.expected {
			t.Errorf("Expected %q to be %q, but it was %q", c.input, c.expected, result)
		}
	}
}
