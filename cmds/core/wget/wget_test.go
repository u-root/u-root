// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A parity test can be run:
//
//	go test
//	EXECPATH="wget -O -" go test
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/curl"
)

const content = "Very simple web server"

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/200":
		w.WriteHeader(200)
		w.Write([]byte(content))
	case "/302":
		http.Redirect(w, r, "/200", http.StatusFound /* 302 */)
	case "/500":
		w.WriteHeader(500)
		w.Write([]byte(content))
	default:
		w.WriteHeader(404)
		w.Write([]byte(content))
	}
}

func getListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("error setting up TCP listener: %v", err)
	}
	return l, l.Addr().(*net.TCPAddr).Port
}

// TestWget implements a table-driven test.
func TestWget(t *testing.T) {
	// Start a webserver on a free port.
	l, port := getListener(t)
	defer l.Close()
	ul, unusedPort := getListener(t)
	defer ul.Close()
	go func() {
		for {
			conn, err := ul.Accept()
			if err != nil {
				// End of test.
				return
			}
			conn.Close()
		}
	}()

	h := handler{}
	go func() {
		log.Print(http.Serve(l, h))
	}()

	tmpDir := t.TempDir()

	var tests = []struct {
		name       string
		url        string // in
		content    string // out
		outputPath string
		err        error
	}{
		{
			name:       "basic",
			url:        fmt.Sprintf("http://localhost:%d/200", port),
			content:    content,
			outputPath: filepath.Join(tmpDir, "basic"),
			err:        nil,
		},
		{
			name:       "ipv4",
			url:        fmt.Sprintf("http://127.0.0.1:%d/200", port),
			content:    content,
			outputPath: filepath.Join(tmpDir, "ipv4"),
			err:        nil,
		},
		// TODO: CircleCI does not support ipv6
		// {
		// 	name:    "ipv6",
		// 	flags:   []string{},
		// 	url:     "http://[::1]:%[1]d/200",
		// 	content: content,
		// 	retCode: 0,
		// },
		{
			name:       "redirect",
			url:        fmt.Sprintf("http://localhost:%[1]d/302", port),
			content:    "",
			outputPath: filepath.Join(tmpDir, "redirect"),
			err:        nil,
		},
		{
			name:    "4xx error",
			url:     fmt.Sprintf("http://localhost:%d/404", port),
			content: "",
			err:     curl.ErrStatusNotOk,
		},
		{
			name:    "5xx error",
			url:     fmt.Sprintf("http://localhost:%d/500", port),
			content: "",
			err:     curl.ErrStatusNotOk,
		},
		{
			name:       "no server",
			url:        fmt.Sprintf("http://localhost:%d/200", unusedPort),
			outputPath: filepath.Join(tmpDir, "no-server"),
			content:    "",
			err:        io.EOF,
		},
		{
			name:    "empty url",
			url:     "",
			content: "",
			err:     errEmptyURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.outputPath, tt.url).run()

			if tt.err == nil && err != nil {
				t.Errorf("expect nil, got %v", err)
			}
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("expect: %v, got: %v", tt.err, err)
			}

			if tt.content != "" {
				content, err := os.ReadFile(tt.outputPath)
				if err != nil {
					t.Errorf("File %s was not created: %v", tt.outputPath, err)
				}

				// Check content.
				if string(content) != tt.content {
					t.Errorf("Want:\n%#v\nGot:\n%#v", tt.content, string(content))
				}
			}
		})
	}
}

func TestDefaultOutputPath(t *testing.T) {
	tests := []struct {
		path   string
		output string
	}{
		{
			path:   "/",
			output: "index.html",
		},
		{
			path:   "",
			output: "index.html",
		},
		{
			path:   "file",
			output: "file",
		},
	}

	for _, test := range tests {
		r := defaultOutputPath(test.path)
		if r != test.output {
			t.Errorf("expect: %s, got: %s", test.output, r)
		}
	}
}
