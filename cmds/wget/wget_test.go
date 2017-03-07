// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/u-root/u-root/shared/testutil"
)

const content = "Very simple web server"

type handler struct {
	status chan int
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/redirect":
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	default:
		w.WriteHeader(<-h.status)
		w.Write([]byte(content))
	}
}

var tests = []struct {
	status  int      // in
	flags   []string // in, %[1]d is the server's port, %[2] is an unopen port
	url     string   // in
	stdout  string   // out
	retCode int      // out
}{
	{
		// basic
		status:  http.StatusOK,
		flags:   []string{"-O", "-"},
		url:     "http://localhost:%[1]d/",
		stdout:  content,
		retCode: 0,
	}, {
		// ipv4
		status:  http.StatusOK,
		flags:   []string{"-O", "-"},
		url:     "http://127.0.0.1:%[1]d/",
		stdout:  content,
		retCode: 0,
	}, {
		// ipv6
		status:  http.StatusOK,
		flags:   []string{"-O", "-"},
		url:     "http://[::1]:%[1]d/",
		stdout:  content,
		retCode: 0,
	}, {
		// redirect
		status:  http.StatusOK,
		flags:   []string{"-O", "-"},
		url:     "http://localhost:%[1]d/redirect",
		stdout:  content,
		retCode: 0,
	}, {
		// 4xx error
		status:  http.StatusNotFound,
		flags:   []string{"-O", "-"},
		url:     "http://localhost:%[1]d/",
		stdout:  "",
		retCode: 1,
	}, {
		// 5xx error
		status:  http.StatusInternalServerError,
		flags:   []string{"-O", "-"},
		url:     "http://localhost:%[1]d/",
		stdout:  "",
		retCode: 1,
	}, {
		// no server
		status:  http.StatusOK,
		flags:   []string{"-O", "-"},
		url:     "http://localhost:%[2]d/",
		stdout:  "",
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
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Start a webserver on a free port.
	port := getFreePort(t)
	unusedPort := getFreePort(t)
	h := handler{make(chan int, 1)}
	go func() {
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))
	}()

	time.Sleep(500 * time.Millisecond) // TODO: better synchronization
	for i, tt := range tests {
		h.status <- tt.status

		args := append(tt.flags, fmt.Sprintf(tt.url, port, unusedPort))
		out, err := exec.Command(execPath, args...).Output()

		// Check return code.
		retCode := 0
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Errorf("%d. Error running wget: %v", i, err)
				continue
			}
			retCode = exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		}
		if retCode != tt.retCode {
			t.Errorf("%d. Want: %d; Got: %d", i, tt.retCode, retCode)
		}

		// Check stdout.
		if string(out) != tt.stdout {
			t.Errorf("%d. Want:\n%#v\nGot:\n%#v", i, tt.stdout, string(out))
		}
	}
}
