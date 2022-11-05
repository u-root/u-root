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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

// TestWget implements a table-driven test.
func TestWget(t *testing.T) {
	srv := httptest.NewServer(handler{})
	defer srv.Close()

	var tests = []struct {
		name        string
		url         string // in
		wantContent string // out
		outputPath  string
		wantErr     error
	}{
		{
			name:        "ipv4",
			url:         fmt.Sprintf("%s/200", srv.URL),
			wantContent: content,
			outputPath:  "basic",
			wantErr:     nil,
		},
		{
			name:        "localhost",
			url:         strings.Replace(srv.URL, "127.0.0.1", "localhost", 1) + "/200",
			wantContent: content,
			outputPath:  "ipv4",
			wantErr:     nil,
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
			url:        fmt.Sprintf("%s/302", srv.URL),
			outputPath: "redirect",
			wantErr:    nil,
		},
		{
			name:    "4xx error",
			url:     fmt.Sprintf("%s/404", srv.URL),
			wantErr: curl.ErrStatusNotOk,
		},
		{
			name:    "5xx error",
			url:     fmt.Sprintf("%s/500", srv.URL),
			wantErr: curl.ErrStatusNotOk,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: errEmptyURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.outputPath = filepath.Join(t.TempDir(), tt.outputPath)
			err := newCommand(tt.outputPath, tt.url).run()

			if tt.wantErr == nil && err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantContent != "" {
				content, err := os.ReadFile(tt.outputPath)
				if err != nil {
					t.Fatalf("file %s was not created: %v", tt.outputPath, err)
				}

				// Check content.
				if string(content) != tt.wantContent {
					t.Errorf("wanted:\n%#v\ngot:\n%#v", tt.wantContent, string(content))
				}
			}
		})
	}
}

func TestDefaultOutputPath(t *testing.T) {
	tests := []struct {
		path       string
		wantOutput string
	}{
		{
			path:       "/",
			wantOutput: "index.html",
		},
		{
			path:       "",
			wantOutput: "index.html",
		},
		{
			path:       "file",
			wantOutput: "file",
		},
	}

	for _, test := range tests {
		got := defaultOutputPath(test.path)
		if got != test.wantOutput {
			t.Errorf("defaultOutputPath(%q) returned: (%q) expected: %q", test.path, got, test.wantOutput)
		}
	}
}
