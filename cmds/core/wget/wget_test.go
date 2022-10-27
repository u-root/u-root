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

var tests = []struct {
	name    string
	flags   []string // in, %[1]d is the server's port, %[2] is an unopen port
	url     string   // in
	content string   // out
	err     error
}{
	{
		name:    "basic",
		flags:   nil,
		url:     "http://localhost:%[1]d/200",
		content: content,
		err:     nil,
	},
	{
		name:    "ipv4",
		flags:   nil,
		url:     "http://127.0.0.1:%[1]d/200",
		content: content,
		err:     nil,
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
		name:    "redirect",
		flags:   nil,
		url:     "http://localhost:%[1]d/302",
		content: "",
		err:     nil,
	},
	{
		name:    "4xx error",
		flags:   nil,
		url:     "http://localhost:%[1]d/404",
		content: "",
		err:     curl.ErrStatusNotOk,
	},
	{
		name:    "5xx error",
		flags:   nil,
		url:     "http://localhost:%[1]d/500",
		content: "",
		err:     curl.ErrStatusNotOk,
	},
	{
		name:    "no server",
		flags:   nil,
		url:     "http://localhost:%[2]d/200",
		content: "",
		err:     io.EOF,
	},
	// {
	// 	name:    "empty url",
	// 	flags:   nil,
	// 	url:     "",
	// 	content: "",
	// 	err:     errEmptyURL,
	// },
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change the working directory to a temporary directory, so we can
			// delete the temporary files after the test runs.
			tmpDir := t.TempDir()
			fileName := filepath.Base(tt.url)
			*outPath = filepath.Join(tmpDir, fileName)
			err := run(fmt.Sprintf(tt.url, port, unusedPort))

			if tt.err == nil && err != nil {
				t.Errorf("expect nil, got %v", err)
			}
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("expect: %v, got: %v", tt.err, err)
			}

			if tt.content != "" {
				content, err := os.ReadFile(filepath.Join(tmpDir, fileName))
				if err != nil {
					t.Errorf("File %s was not created: %v", fileName, err)
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
