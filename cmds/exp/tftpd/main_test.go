// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"pack.ag/tftp"
)

func TestHandlersStayInRoot(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "tftp")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "file"), []byte("file"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "secret"), []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(base, "secret"), filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	// A dangling symlink: the write must not create its target.
	if err := os.Symlink(filepath.Join(base, "planted"), filepath.Join(dir, "dangling")); err != nil {
		t.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	server, err := tftp.NewServer("")
	if err != nil {
		t.Fatal(err)
	}
	server.ReadHandler(&FileReadHandler{Root: dir})
	server.WriteHandler(&FileWriteHandler{Root: dir})
	go server.Serve(conn)
	defer server.Close()

	client, err := tftp.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	url := func(name string) string {
		return fmt.Sprintf("tftp://%s/%s", conn.LocalAddr(), name)
	}

	for _, tt := range []struct {
		name string
		want string
	}{
		{name: "file", want: "file"},
		{name: "/file", want: "file"},
		{name: "../secret"},
		{name: "link"},
	} {
		t.Run("read "+tt.name, func(t *testing.T) {
			resp, err := client.Get(url(tt.name))
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

	data := []byte("data")
	for _, tt := range []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "new", path: filepath.Join(dir, "new")},
		{name: "../escape", path: filepath.Join(base, "escape"), wantErr: true},
		{name: "dangling", path: filepath.Join(base, "planted"), wantErr: true},
	} {
		t.Run("write "+tt.name, func(t *testing.T) {
			err := client.Put(url(tt.name), bytes.NewReader(data), int64(len(data)))
			got, rerr := os.ReadFile(tt.path)
			if tt.wantErr {
				if err == nil || rerr == nil {
					t.Fatalf("Put(%q) = %v, %s: %q, %v; want error and no file", tt.name, err, tt.path, got, rerr)
				}
				return
			}
			if err != nil || rerr != nil || !bytes.Equal(got, data) {
				t.Fatalf("Put(%q) = %v, %s: %q, %v; want nil, %q", tt.name, err, tt.path, got, rerr, data)
			}
		})
	}
}
