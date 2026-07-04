// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteDNSSettings(t *testing.T) {
	ns := []net.IP{net.ParseIP("10.0.0.1")}

	for _, tt := range []struct {
		name    string
		domain  string
		search  []string
		wantErr bool
	}{
		{
			name:   "clean",
			domain: "corp.example",
			search: []string{"corp.example", "example.com"},
		},
		{
			name:    "newline in domain",
			domain:  "corp.example\nnameserver 6.6.6.6",
			wantErr: true,
		},
		{
			name:    "carriage return in domain",
			domain:  "corp.example\rnameserver 6.6.6.6",
			wantErr: true,
		},
		{
			name:    "newline in search label",
			search:  []string{"a.example\noptions ndots:0"},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "resolv.conf")
			err := WriteDNSSettings(ns, tt.search, tt.domain, path)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("WriteDNSSettings(%q, %q) = nil, want error", tt.domain, tt.search)
				}
				// The malformed field must never reach the file.
				if b, rerr := os.ReadFile(path); rerr == nil {
					if strings.Contains(string(b), "6.6.6.6") || strings.Contains(string(b), "ndots") {
						t.Fatalf("injected directive written to resolv.conf:\n%s", b)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("WriteDNSSettings() = %v, want nil", err)
			}
			b, rerr := os.ReadFile(path)
			if rerr != nil {
				t.Fatalf("ReadFile: %v", rerr)
			}
			if got := string(b); !strings.Contains(got, "nameserver 10.0.0.1\n") {
				t.Fatalf("resolv.conf missing nameserver line:\n%s", got)
			}
		})
	}
}
