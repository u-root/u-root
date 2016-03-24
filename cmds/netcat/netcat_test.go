package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

var Host = "127.0.0.1"
var Port = ":9991"
var Input = "Input from other side, пока, £, 语汉"

func TestTCP(t *testing.T) {
	// Send test data to listener from goroutine and wait for potentials errors at the end of the test
	go func() {
		// Wait for main thread starts listener
		time.Sleep(200 * time.Millisecond)
		con, err := net.Dial("tcp", Host+Port)

		if err != nil {
			t.Fatalf("Connection using tcp %v%v fails: %v", Host, Port, err)
		}

		// Transfer data
		c1 := readAndWrite(strings.NewReader(Input), con)

		// Wait for data will be transferred
		time.Sleep(200 * time.Millisecond)
		select {
		case progress := <-c1:
			t.Logf("Remote connection is closed: %+v\n", progress)
		default:
			t.Fatal("handle() must write to result channel")
		}
	}()

	ln, err := net.Listen("tcp", Port)
	if err != nil {
		t.Errorf("Listen Port %q fails using TCP: %v", Port, err)
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
	if Input != output {
		t.Errorf("Message passing between connections mismatch; wants %v, got %v", Input, output)
	}

}

func TestUDP(t *testing.T) {
	// Send test data to listener from goroutine and wait for potentials errors at the end of the test
	go func() {
		// Wait for main thread starts listener
		time.Sleep(200 * time.Millisecond)
		con, err := net.Dial("udp", Host+Port)
		if err != nil {
			t.Fatalf("Connection using udp %v%v fails: %v", Host, Port, err)
		}

		// Transfer data
		addr, err := net.ResolveUDPAddr("udp", Host+Port)
		fmt.Println(con.RemoteAddr())
		c1 := readAndWriteToAddr(strings.NewReader(Input), con, addr)

		// Wait for data will be transferred
		time.Sleep(200 * time.Millisecond)
		select {
		case progress := <-c1:
			t.Logf("Remote connection is closed: %+v\n", progress)
		default:
			t.Fatal("handle() must write to result channel")
		}
	}()

	con, err := net.ListenPacket("udp", Port)
	if err != nil {
		t.Errorf("Listen Port %q fails using TCP: %v", Port, err)
	}

	buf := make([]byte, 1024)
	n, _, err := con.ReadFrom(buf)

	if err != nil {
		t.Errorf("Reading from connection fails: %v", err)
	}

	output := string(buf[0:n])
	if Input != output {
		t.Errorf("Message passing between connections mismatch; wants %v, got %v", Input, output)
	}

}
