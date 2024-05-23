// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"testing"
)

func TestEvalFlags(t *testing.T) {
	for _, tt := range []struct {
		name   string
		ops    string
		arg    []string
		af     string
		expErr error
	}{
		{
			name:   "SucessSocketsDefault",
			ops:    "sockets",
			expErr: nil,
		},
		{
			name:   "SucessSocketsIPv4",
			ops:    "sockets",
			af:     "ipv4",
			expErr: nil,
		},
		{
			name:   "SucessSocketsIPv6",
			ops:    "sockets",
			af:     "ipv6",
			expErr: nil,
		},
		{
			name:   "SucessRoute",
			ops:    "route",
			expErr: nil,
		},
		{
			name:   "SuccessInterfaces",
			ops:    "interfaces",
			expErr: nil,
		},
		{
			name:   "SuccessIface_th0",
			ops:    "iface",
			arg:    []string{"eth0"},
			expErr: nil,
		},
		{
			name:   "SuccessStats",
			ops:    "stats",
			expErr: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.ops {
			case "route":
				*routeFlag = true
			case "interfaces":
				*interfacesFlag = true
			case "stats":
				*statsFlag = true
				*ipv4Flag = true
				*ipv6Flag = true
			case "sockets":
				switch tt.af {
				case "ipv4":
					*ipv4Flag = true
				case "ipv6":
					*ipv6Flag = true
				}
			case "groups":
				*groupsFlag = true
			case "iface":
				*ifFlag = tt.arg[0]
			}
			if err := evalFlags(); !errors.Is(err, tt.expErr) {
				t.Errorf("evalFlags() failed: %v, want: %v", err, tt.expErr)
			}

			resetFlags()
		})
	}
}

func resetFlags() {
	// Info source flags
	*routeFlag = false
	*interfacesFlag = false
	*ifFlag = ""
	*groupsFlag = false
	*statsFlag = false

	// Socket flags
	*tcpFlag = false
	*udpFlag = false
	*udpLFlag = false
	*rawFlag = false
	*unixFlag = false

	// AF Flags
	*ipv4Flag = false
	*ipv6Flag = false

	// Route type flag
	*routecacheFalg = false

	// Format flags
	*wideFlag = false
	*numericFlag = false
	*numHostFlag = false
	*numPortsFlag = false
	*numUsersFlag = false
	*symbolicFlag = false
	*extendFlag = false
	*programsFlag = false
	*timersFlag = false
	*continFlag = false
	*listeningFlag = false
	*allFlag = false
}
