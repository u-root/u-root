// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import "testing"

func TestNewServer(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		addr string
		opts []ServerOpt

		expectedAddrStr    string
		expectedNet        string
		expectedRetransmit int
		expectedError      error
	}{
		"default": {
			addr: "",

			expectedNet:        "udp",
			expectedRetransmit: 10,
		},
		"net udp6": {
			addr: "",
			opts: []ServerOpt{
				ServerNet("udp6"),
			},

			expectedNet:        "udp6",
			expectedRetransmit: 10,
		},
		"net, invalid": {
			addr: "",
			opts: []ServerOpt{
				ServerNet("tcp"),
			},

			expectedError: ErrInvalidNetwork,
		},
		"retransmit, valid": {
			addr: "",
			opts: []ServerOpt{
				ServerRetransmit(2),
			},

			expectedNet:        "udp",
			expectedRetransmit: 2,
		},
		"retransmit, invalid": {
			addr: "",
			opts: []ServerOpt{
				ServerRetransmit(-1),
			},

			expectedError: ErrInvalidRetransmit,
		},
	}

	for label, c := range cases {
		server, err := NewServer(c.addr, c.opts...)

		// Error
		if err != c.expectedError {
			t.Errorf("%s: Expected %#v to be %#v", label, err, c.expectedError)
		}

		if err != nil {
			continue // Skip remaining test if error, avoid nil dereference
		}

		// Addr
		if server.addrStr != c.expectedAddrStr {
			t.Errorf("%s: Expected addr to be %q, but it was %q", label, c.expectedAddrStr, server.addrStr)
		}

		// Net
		if server.net != c.expectedNet {
			t.Errorf("%s: Expected net to be %q, but it was %q", label, c.expectedNet, server.net)
		}

		// Retransmit
		if server.retransmit != c.expectedRetransmit {
			t.Errorf("%s: Expected retransmit to be %d, but it was %d", label, c.expectedRetransmit, server.retransmit)
		}
	}
}
