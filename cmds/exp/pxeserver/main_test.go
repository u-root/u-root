// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"pack.ag/tftp"
)

func TestTFTPFileServer(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "tftp")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "pxelinux.0"), []byte("boot"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "secret"), []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(base, "secret"), filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	server, err := tftp.NewServer("")
	if err != nil {
		t.Fatal(err)
	}
	server.ReadHandler(&tftpFileServer{root: root})
	go server.Serve(conn)
	defer server.Close()

	client, err := tftp.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name string
		want string
	}{
		{name: "pxelinux.0", want: "boot"},
		{name: "/pxelinux.0", want: "boot"},
		{name: "../secret"},
		{name: "link"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(fmt.Sprintf("tftp://%s/%s", conn.LocalAddr(), tt.name))
			var got []byte
			if err == nil {
				got, err = io.ReadAll(resp)
			}
			if tt.want == "" {
				if err == nil {
					t.Fatalf("Get(%q) = %q, want error", tt.name, got)
				}
				return
			}
			if err != nil || string(got) != tt.want {
				t.Fatalf("Get(%q) = %q, %v, want %q, nil", tt.name, got, err, tt.want)
			}
		})
	}
}
