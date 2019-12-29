// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import "testing"

func TestIsUnexpectedDatagram(t *testing.T) {
	cases := []struct {
		name string
		err  error

		expected bool
	}{
		{
			name:     "true",
			err:      &errUnexpectedDatagram{},
			expected: true,
		},
		{
			name:     "true, wrapped",
			err:      wrapError(&errUnexpectedDatagram{}, "testing"),
			expected: true,
		},
		{
			name:     "false",
			err:      errBlockSequence,
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := IsUnexpectedDatagram(c.err)
			if result != c.expected {
				t.Errorf("expected to IsUnexpectedDatagram %t, but it wasn't", c.expected)
			}
		})
	}
}

func TestIsRemoteError(t *testing.T) {
	cases := []struct {
		name string
		err  error

		expected bool
	}{
		{
			name:     "true",
			err:      &errRemoteError{},
			expected: true,
		},
		{
			name:     "true, wrapped",
			err:      wrapError(&errRemoteError{}, "testing"),
			expected: true,
		},
		{
			name:     "false",
			err:      errBlockSequence,
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := IsRemoteError(c.err)
			if result != c.expected {
				t.Errorf("expected to IsUnexpectedDatagram %t, but it wasn't", c.expected)
			}
		})
	}
}

func TestIsOptionParsingError(t *testing.T) {
	cases := []struct {
		name string
		err  error

		expected bool
	}{
		{
			name:     "true",
			err:      &errParsingOption{},
			expected: true,
		},
		{
			name:     "true, wrapped",
			err:      wrapError(&errParsingOption{}, "testing"),
			expected: true,
		},
		{
			name:     "false",
			err:      errBlockSequence,
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := IsOptionParsingError(c.err)
			if result != c.expected {
				t.Errorf("expected to IsUnexpectedDatagram %t, but it wasn't", c.expected)
			}
		})
	}
}

func TestErrorStrings(t *testing.T) {
	dg := datagram{}
	dg.writeAck(68)

	cases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "unexpected datagram",
			err:      &errUnexpectedDatagram{dg: dg.String()},
			expected: `unexpected datagram: ACK[Block: 68]`,
		},
		{
			name:     "remote error",
			err:      &errRemoteError{dg: dg.String()},
			expected: `remote error: ACK[Block: 68]`,
		},
		{
			name:     "parse error",
			err:      &errParsingOption{option: "timeout", value: "a"},
			expected: `error parsing "a" for option "timeout"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.err.Error() != c.expected {
				t.Errorf("Expected %q to be %q", c.err.Error(), c.expected)
			}
		})
	}
}
