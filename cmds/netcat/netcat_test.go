// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

var tableDriven = []struct {
	host, port, input string
}{
	{"127.0.0.1", ":9991", "An unicode²³£øĸøþ stream using IPV4"},
	{"[::1]", ":9992", "An unicode²³£øĸøþ stream using IPV6"},
}

func TestTCP(t *testing.T) {
	// Send test data to listener from goroutine and wait for potentials errors at the end of the test
	for _, test := range tableDriven {

		go func() {
			// Wait for main thread starts listener
			time.Sleep(200 * time.Millisecond)
			con, err := net.Dial("tcp", test.host+test.port)

			if err != nil {
				t.Fatalf("Connection using tcp %v%v fails: %v", test.host, test.port, err)
			}

			// Transfer data
			c1 := readAndWrite(strings.NewReader(test.input), con)

			// Wait for data will be transferred
			time.Sleep(200 * time.Millisecond)
			select {
			case progress := <-c1:
				t.Logf("Remote connection is closed: %+v\n", progress)
			default:
				t.Fatal("handle() must write to result channel")
			}
		}()

		ln, err := net.Listen("tcp", test.port)
		if err != nil {
			t.Errorf("Listen Port %q fails using TCP: %v", test.port, err)
		}

		con, err := ln.Accept()
		if err != nil {
			t.Errorf("Connecting accept fails: %v", err)
		}

		buf := make([]byte, 1024)
		n, err := con.Read(buf)
		if err != nil {
			t.Errorf("Reading from connection fails: %v", err)
		}

		output := string(buf[0:n])
		if test.input != output {
			t.Errorf("Message passing between connections mismatch; wants %v, got %v", test.input, output)
		}
	}
}

func TestUDP(t *testing.T) {
	// Send test data to listener from goroutine and wait
	// for potentials errors at the end of the test
	for _, test := range tableDriven {
		go func() {
			// Wait for main thread starts listener
			time.Sleep(200 * time.Millisecond)
			con, err := net.Dial("udp", test.host+test.port)
			if err != nil {
				t.Fatalf("Connection using udp %v%v fails: %v", test.host, test.port, err)
			}

			// Transfer data
			addr, err := net.ResolveUDPAddr("udp", test.host+test.port)
			fmt.Println(con.RemoteAddr())
			c1 := readAndWriteToAddr(strings.NewReader(test.input), con, addr)

			// Wait for data will be transferred
			time.Sleep(200 * time.Millisecond)
			select {
			case progress := <-c1:
				t.Logf("Remote connection is closed: %+v\n", progress)
			default:
				t.Fatal("handle() must write to result channel")
			}
		}()

		con, err := net.ListenPacket("udp", test.port)
		if err != nil {
			t.Errorf("Listen port %q fails using TCP: %v", test.port, err)
		}

		buf := make([]byte, 1024)
		n, _, err := con.ReadFrom(buf)

		if err != nil {
			t.Errorf("Reading from connection fails: %v", err)
		}

		output := string(buf[0:n])
		if test.input != output {
			t.Errorf("Message passing between connections mismatch; wants %v, got %v", test.input, output)
		}
	}
}
