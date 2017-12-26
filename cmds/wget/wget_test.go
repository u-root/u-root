// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A parity test can be run:
//     go test
//     EXECPATH="wget -O -" go test
package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const content = "Very simple web server"

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/200":
		w.WriteHeader(200)
		w.Write([]byte(content))
	case "/302":
		http.Redirect(w, r, "/200", 302)
	case "/500":
		w.WriteHeader(500)
		w.Write([]byte(content))
	default:
		w.WriteHeader(404)
		w.Write([]byte(content))
	}
}

var tests = []struct {
	url    string
	stdout string
	isErr  bool
}{
	{
		// basic
		url:    "http://localhost:%[1]d/200",
		stdout: content,
	}, {
		// ipv4
		url:    "http://127.0.0.1:%[1]d/200",
		stdout: content,
	}, {
		// ipv6
		url:    "http://[::1]:%[1]d/200",
		stdout: content,
	}, {
		// redirect
		url:    "http://localhost:%[1]d/302",
		stdout: content,
	}, {
		// 4xx error
		url:   "http://localhost:%[1]d/404",
		isErr: true,
	}, {
		// 5xx error
		url:   "http://localhost:%[1]d/500",
		isErr: true,
	}, {
		// no server
		url:   "http://localhost:%[2]d/200",
		isErr: true,
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
	port := getFreePort(t)
	unusedPort := getFreePort(t)
	go func() {
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler{}))
	}()

	time.Sleep(500 * time.Millisecond) // TODO: better synchronization
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test [%02d]", i), func(t *testing.T) {
			execPath := os.Getenv("EXECPATH")
			uri := fmt.Sprintf(tt.url, port, unusedPort)

			var err error
			var out string

			if len(execPath) > 0 {
				// Arguments inherited from the environment.
				execArgs := strings.Fields(execPath)
				args := append(execArgs[1:], uri)
				cmd := exec.Command(execArgs[0], args...)

				var byteOut []byte
				byteOut, err = cmd.Output()
				out = string(byteOut)
			} else {
				var stdout bytes.Buffer
				err = wget(uri, &stdout)
				out = stdout.String()
			}

			// Check return code.
			if tt.isErr && err == nil {
				t.Errorf("wget(%s) got no error, but expected error", uri)
			} else if !tt.isErr && err != nil {
				t.Errorf("wget(%s) got error %v, but expected none", uri, err)
			}

			// Check stdout.
			if out != tt.stdout {
				t.Errorf("wget(%s) = %v, want %v", uri, out, tt.stdout)
			}
		})
	}
}
