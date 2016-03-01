// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Wget reads one file from the argument and writes it on the standard output.
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type fdata struct {
	err error
	data []byte
	url string
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b := bytes.NewBuffer([]byte{})
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func get(url string, result chan<- *fdata) {
	b, err := fetch(url)
	result <- &fdata{data: b, err: err, url: url}
}


func main() {
	// synchronously get the head. If that's not there, there's no
	// point in doing anything else.
	h, err := fetch(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("Head %v\n", h)
}
