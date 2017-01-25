// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

const webSite = "http://example.com"

func TestWget(t *testing.T) {

	var buf, buf2 bytes.Buffer

	if err := wget(webSite, &buf); err != nil {
		t.Fatalf("%v", err)
	}

	resp, err := http.Get(webSite)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(&buf2, resp.Body); err != nil {
		t.Fatalf("%v", err)
	}

	if bytes.Compare(buf.Bytes(), buf2.Bytes()) != 0 {
		t.Fatalf("Fetching %v: want %v got %v", webSite, buf2, buf)
	}
}
