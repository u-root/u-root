// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A parity test can be run:
//     go test
//     EXECPATH="wget -O -" go test
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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
	flags   []string // in, %[1]d is the server's port, %[2] is an unopen port
	url     string   // in
	content string   // out
	retCode int      // out
}{
	{
		// basic
		flags:   []string{},
		url:     "http://localhost:%[1]d/200",
		content: content,
		retCode: 0,
	}, {
		// ipv4
		flags:   []string{},
		url:     "http://127.0.0.1:%[1]d/200",
		content: content,
		retCode: 0,
	}, /*{ TODO: travis does not support ipv6
		// ipv6
		flags:   []string{},
		url:     "http://[::1]:%[1]d/200",
		content:  content,
		retCode: 0,
	},*/{
		// redirect
		flags:   []string{},
		url:     "http://localhost:%[1]d/302",
		content: "",
		retCode: 0,
	}, {
		// 4xx error
		flags:   []string{},
		url:     "http://localhost:%[1]d/404",
		content: "",
		retCode: 1,
	}, {
		// 5xx error
		flags:   []string{},
		url:     "http://localhost:%[1]d/500",
		content: "",
		retCode: 1,
	}, {
		// no server
		flags:   []string{},
		url:     "http://localhost:%[2]d/200",
		content: "",
		retCode: 1,
	}, {
		// output file
		flags:   []string{"-O", "/dev/null"},
		url:     "http://localhost:%[1]d/200",
		content: "",
		retCode: 0,
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

	for i, tt := range tests {
		args := append(tt.flags, fmt.Sprintf(tt.url, port, unusedPort))
		output, err := testutil.Command(t, args...).CombinedOutput()

		// Check return code.
		if err := testutil.IsExitCode(err, tt.retCode); err != nil {
			t.Errorf("exit code: %v, output: %s", err, string(output))
		}

		if tt.content != "" {
			fileName := filepath.Base(tt.url)
			content, err := ioutil.ReadFile(fileName)
			if err != nil {
				t.Errorf("%d. File %s was not created: %v", i, fileName, err)
			}

			// Check content.
			if string(content) != tt.content {
				t.Errorf("%d. Want:\n%#v\nGot:\n%#v", i, tt.content, string(content))
			}
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
