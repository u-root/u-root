// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A parity test can be run:
//     go test
//     EXECPATH="wget -O -" go test
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
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
	retCode int      // out
}{
	{
		name:    "basic",
		flags:   []string{},
		url:     "http://localhost:%[1]d/200",
		content: content,
		retCode: 0,
	},
	{
		name:    "ipv4",
		flags:   []string{},
		url:     "http://127.0.0.1:%[1]d/200",
		content: content,
		retCode: 0,
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
		flags:   []string{},
		url:     "http://localhost:%[1]d/302",
		content: "",
		retCode: 0,
	},
	{
		name:    "4xx error",
		flags:   []string{},
		url:     "http://localhost:%[1]d/404",
		content: "",
		retCode: 1,
	},
	{
		name:    "5xx error",
		flags:   []string{},
		url:     "http://localhost:%[1]d/500",
		content: "",
		retCode: 1,
	},
	{
		name:    "no server",
		flags:   []string{},
		url:     "http://localhost:%[2]d/200",
		content: "",
		retCode: 1,
	},
}

func getFreePort(t *testing.T) int {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Cannot create free port: %v", err)
	}
	l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// TestWget implements a table-driven test.
func TestWget(t *testing.T) {
	// Start a webserver on a free port.
	unusedPort := getFreePort(t)

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Cannot create free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port

	h := handler{}
	go func() {
		log.Fatal(http.Serve(l, h))
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change the working directory to a temporary directory, so we can
			// delete the temporary files after the test runs.
			tmpDir := t.TempDir()

			fileName := filepath.Base(tt.url)

			args := append(tt.flags,
				"-O", filepath.Join(tmpDir, fileName),
				fmt.Sprintf(tt.url, port, unusedPort))
			cmd := testutil.Command(t, args...)
			output, err := cmd.CombinedOutput()

			// Check return code.
			if err := testutil.IsExitCode(err, tt.retCode); err != nil {
				t.Errorf("exit code: %v, output: %s", err, string(output))
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

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
