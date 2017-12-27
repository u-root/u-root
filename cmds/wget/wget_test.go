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
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"testing"
	"time"

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
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Start a webserver on a free port.
	port := getFreePort(t)
	unusedPort := getFreePort(t)
	h := handler{}
	go func() {
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))
	}()

	time.Sleep(500 * time.Millisecond) // TODO: better synchronization
	for i, tt := range tests {
		// Arguments inherited from the environment.
		execArgs := strings.Split(os.Getenv("EXECPATH"), " ")[1:]

		args := append(append(execArgs, tt.flags...), fmt.Sprintf(tt.url, port, unusedPort))
		_, err := exec.Command(execPath, args...).Output()

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

		if tt.content != "" {
			fileName := path.Base(tt.url)
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
