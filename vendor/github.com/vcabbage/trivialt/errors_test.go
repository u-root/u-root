// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import "testing"

func TestIsUnexpectedDatagram(t *testing.T) {
	cases := map[string]struct {
		err error

		expected bool
	}{
		"true": {
			err:      &errUnexpectedDatagram{},
			expected: true,
		},
		"true, wrapped": {
			err:      wrapError(&errUnexpectedDatagram{}, "testing"),
			expected: true,
		},
		"false": {
			err:      errBlockSequence,
			expected: false,
		},
	}

	for label, c := range cases {
		result := IsUnexpectedDatagram(c.err)
		if result != c.expected {
			t.Errorf("%s: Expected to IsUnexpectedDatagram %t, but it wasn't", label, c.expected)
		}
	}
}

func TestIsRemoteError(t *testing.T) {
	cases := map[string]struct {
		err error

		expected bool
	}{
		"true": {
			err:      &errRemoteError{},
			expected: true,
		},
		"true, wrapped": {
			err:      wrapError(&errRemoteError{}, "testing"),
			expected: true,
		},
		"false": {
			err:      errBlockSequence,
			expected: false,
		},
	}

	for label, c := range cases {
		result := IsRemoteError(c.err)
		if result != c.expected {
			t.Errorf("%s: Expected to IsUnexpectedDatagram %t, but it wasn't", label, c.expected)
		}
	}
}

func TestIsOptionParsingError(t *testing.T) {
	cases := map[string]struct {
		err error

		expected bool
	}{
		"true": {
			err:      &errParsingOption{},
			expected: true,
		},
		"true, wrapped": {
			err:      wrapError(&errParsingOption{}, "testing"),
			expected: true,
		},
		"false": {
			err:      errBlockSequence,
			expected: false,
		},
	}

	for label, c := range cases {
		result := IsOptionParsingError(c.err)
		if result != c.expected {
			t.Errorf("%s: Expected to IsUnexpectedDatagram %t, but it wasn't", label, c.expected)
		}
	}
}

func TestErrorStrings(t *testing.T) {
	dg := datagram{}
	dg.writeAck(68)

	cases := []struct {
		err      error
		expected string
	}{
		{
			err:      &errUnexpectedDatagram{dg: dg.String()},
			expected: `unexpected datagram: ACK[Block: 68]`,
		},
		{
			err:      &errRemoteError{dg: dg.String()},
			expected: `remote error: ACK[Block: 68]`,
		},
		{
			err:      &errParsingOption{option: "timeout", value: "a"},
			expected: `error parsing "a" for option "timeout"`,
		},
	}

	for _, c := range cases {
		if c.err.Error() != c.expected {
			t.Errorf("Expected %q to be %q", c.err.Error(), c.expected)
		}
	}
}
