// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import "testing"

func TestNewServer(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		addr string
		opts []ServerOpt

		expectedAddrStr    string
		expectedNet        string
		expectedRetransmit int
		expectedError      error
	}{
		{
			name: "default",
			addr: "",

			expectedNet:        "udp",
			expectedRetransmit: 10,
		},
		{
			name: "net udp6",
			addr: "",
			opts: []ServerOpt{
				ServerNet("udp6"),
			},

			expectedNet:        "udp6",
			expectedRetransmit: 10,
		},
		{
			name: "net, invalid",
			addr: "",
			opts: []ServerOpt{
				ServerNet("tcp"),
			},

			expectedError: ErrInvalidNetwork,
		},
		{
			name: "retransmit, valid",
			addr: "",
			opts: []ServerOpt{
				ServerRetransmit(2),
			},

			expectedNet:        "udp",
			expectedRetransmit: 2,
		},
		{
			name: "retransmit, invalid",
			addr: "",
			opts: []ServerOpt{
				ServerRetransmit(-1),
			},

			expectedError: ErrInvalidRetransmit,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server, err := NewServer(c.addr, c.opts...)

			// Error
			if err != c.expectedError {
				t.Errorf("expected %#v to be %#v", err, c.expectedError)
			}

			if err != nil {
				return // Skip remaining test if error, avoid nil dereference
			}

			// Addr
			if server.addrStr != c.expectedAddrStr {
				t.Errorf("expected addr to be %q, but it was %q", c.expectedAddrStr, server.addrStr)
			}

			// Net
			if server.net != c.expectedNet {
				t.Errorf("expected net to be %q, but it was %q", c.expectedNet, server.net)
			}

			// Retransmit
			if server.retransmit != c.expectedRetransmit {
				t.Errorf("expected retransmit to be %d, but it was %d", c.expectedRetransmit, server.retransmit)
			}
		})
	}
}
