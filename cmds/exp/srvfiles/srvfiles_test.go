// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSRVFiles(t *testing.T) {
	dir := t.TempDir()
	content := []byte("hello world")
	err := os.WriteFile(filepath.Join(dir, "hello"), content, 0o644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	ts := httptest.NewServer(maxAgeHandler(http.FileServer(http.Dir(dir))))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/hello", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("If-Modified-Since", "Mon, 02 Jan 2006 15:04:05 MST")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if !bytes.Equal(b, content) {
		t.Errorf("Expected %q, got %q", content, b)
	}
}
