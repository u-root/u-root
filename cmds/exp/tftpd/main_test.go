// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/pin/tftp/v3"
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
	server := tftp.NewServer((&FileReadHandler{Root: dir}).ServeTFTP, (&FileWriteHandler{Root: dir}).ReceiveTFTP)
	go server.Serve(conn)
	defer server.Shutdown()

	client, err := tftp.NewClient(conn.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
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
			var got bytes.Buffer
			wt, err := client.Receive(tt.name, "octet")
			if err == nil {
				_, err = wt.WriteTo(&got)
			}
			if tt.want == "" {
				if err == nil {
					t.Fatalf("Receive(%q) = %q, want error", tt.name, got.String())
				}
				return
			}
			if err != nil || got.String() != tt.want {
				t.Fatalf("Receive(%q) = %q, %v, want %q, nil", tt.name, got.String(), err, tt.want)
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
			rf, err := client.Send(tt.name, "octet")
			if err == nil {
				_, err = rf.ReadFrom(bytes.NewReader(data))
			}
			got, rerr := os.ReadFile(tt.path)
			if tt.wantErr {
				if err == nil || rerr == nil {
					t.Fatalf("Send(%q) = %v, %s: %q, %v; want error and no file", tt.name, err, tt.path, got, rerr)
				}
				return
			}
			if err != nil || rerr != nil || !bytes.Equal(got, data) {
				t.Fatalf("Send(%q) = %v, %s: %q, %v; want nil, %q", tt.name, err, tt.path, got, rerr, data)
			}
		})
	}
}
