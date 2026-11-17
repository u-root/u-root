// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pin/tftp/v3"
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
	server := tftp.NewServer(tftpReadHandler(root.FS()), nil)
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
		{name: "pxelinux.0", want: "boot"},
		{name: "/pxelinux.0", want: "boot"},
		{name: "../secret"},
		{name: "link"},
	} {
		t.Run(tt.name, func(t *testing.T) {
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
}

func TestHTTPFileServer(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "http")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "kernel"), []byte("kernel"), 0o644); err != nil {
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

	ts := httptest.NewServer(http.FileServerFS(root.FS()))
	defer ts.Close()

	for _, tt := range []struct {
		name string
		want string
	}{
		{name: "kernel", want: "kernel"},
		{name: "link"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + "/" + tt.name)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tt.want == "" {
				if resp.StatusCode == http.StatusOK {
					t.Fatalf("GET %q = %d %q, want error", tt.name, resp.StatusCode, got)
				}
				return
			}
			if resp.StatusCode != http.StatusOK || string(got) != tt.want {
				t.Fatalf("GET %q = %d %q, want %d %q", tt.name, resp.StatusCode, got, http.StatusOK, tt.want)
			}
		})
	}
}
