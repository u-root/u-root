// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"os"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/ulog"
)

func TestNewUDPListener(t *testing.T) {
	logger := ulog.Null
	tests := []struct {
		network string
		addr    string
	}{
		{"udp", "127.0.0.1:0"},
		{"udp4", "127.0.0.1:0"},
		{"unixgram", "/tmp/test_unixgram"},
	}

	for _, tt := range tests {
		t.Run(tt.network, func(t *testing.T) {
			if tt.network == "unixgram" {
				// Ensure the unixgram socket does not exist before testing
				os.Remove(tt.addr)
				defer os.Remove(tt.addr)
			}

			l, err := NewUDPListener(tt.network, tt.addr, logger)
			if err != nil {
				t.Fatalf("NewUDPListener() error = %v, wantErr %v", err, false)
			}
			defer l.Close()

			if l == nil {
				t.Fatal("Expected non-nil UDPListener")
			}

			if addr := l.Addr().String(); strings.HasPrefix(addr, tt.addr) && addr != tt.addr {
				t.Errorf("Expected %v, got %v", tt.addr, addr)
			}
		})
	}
}
