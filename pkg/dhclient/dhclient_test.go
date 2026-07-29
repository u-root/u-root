// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteDNSSettingsPreservesExistingSettings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "resolv.conf")
	initial := "# managed by another service\nnameserver 192.0.2.1\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := WriteDNSSettings([]net.IP{net.ParseIP("192.0.2.2")}, nil, "", path); err != nil {
		t.Fatal(err)
	}
	if err := WriteDNSSettings([]net.IP{net.ParseIP("192.0.2.3")}, nil, "", path); err != nil {
		t.Fatal(err)
	}

	want := initial + "nameserver 192.0.2.2\nnameserver 192.0.2.3\n"
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("resolv.conf = %q, want %q", got, want)
	}
}

func TestWriteDNSSettingsCreatesFileAndPreservesEmptyUpdate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "resolv.conf")

	if err := WriteDNSSettings(nil, nil, "", path); err != nil {
		t.Fatal(err)
	}
	if err := WriteDNSSettings(nil, []string{"example.test", "internal.example.test"}, "example.test", path); err != nil {
		t.Fatal(err)
	}
	if err := WriteDNSSettings(nil, nil, "", path); err != nil {
		t.Fatal(err)
	}

	want := "domain example.test\nsearch example.test internal.example.test\n"
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("resolv.conf = %q, want %q", got, want)
	}
}
