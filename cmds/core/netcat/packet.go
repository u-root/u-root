// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

// buffer represents a fixed size (64 KB) byte array, with the number of bytes
// used (at the front) given by n.
type buffer struct {
	b [65536]byte
	n int
}

// transferPackets copies data from c.stdin to packetConn, and from packetConn
// to output. Packet size is limited to 64 KB.
//
// If listenMode is true, then transferPackets doesn't start reading c.stdin
// until such a packet is received from packetConn whose sender passes the
// allow list. Once that happens, said sender address determines where
// transferPackets sends data from c.stdin. transferPackets only returns upon
// errors; reading EOF from c.stdin is not considered an exit condition.
//
// If listenMode is false, then packetConn is expected to implement net.Conn,
// and to have a default destination address. transferPackets returns upon
// errors, plus it returns nil if EOF is seen on c.stdin.
func (c *cmd) transferPackets(output io.Writer, packetConn net.PacketConn, listenMode bool) error {
	var connToOutput buffer
	var clientAddr net.Addr
	var wg sync.WaitGroup
	var stdinToConnError error
	var connToOutputError error

	if listenMode {
		for {
			var err error

			connToOutput.n, clientAddr, err = packetConn.ReadFrom(connToOutput.b[:])
			if err != nil {
				return err
			}
			if c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_UDP &&
				!c.config.AccessControl.IsAllowed(parseRemoteAddr(netcat.SOCKET_TYPE_UDP, clientAddr.String())) {
				continue
			}
			log.Printf("receiving packets from %v", clientAddr)
			break
		}
	}

	// Upon reading EOF from c.stdin, the goroutine that copies c.stdin to
	// packetConn aborts the goroutine that copies packetConn to output, by
	// sending an empty message over a channel. However, the latter goroutine may
	// have exited meanwhile, due to an independent error. The former goroutine
	// should not block in its abort attempt indefinitely in this case (just
	// because nothing is reading the channel anymore); thus, make the channel
	// buffered (with buffer size 1). If the latter goroutine is no longer there,
	// the empty message in the channel buffer is ignored.
	abortConnToOutput := make(chan struct{}, 1)

	wg.Add(2)

	// copy c.stdin to packetConn
	go func() {
		defer wg.Done()

		for {
			var stdinToConn buffer
			var readError error
			var sendError error

			stdinToConn.n, readError = c.stdin.Read(stdinToConn.b[:])
			if stdinToConn.n > 0 {
				if listenMode {
					_, sendError = packetConn.WriteTo(stdinToConn.b[:stdinToConn.n], clientAddr)
				} else {
					_, sendError = packetConn.(net.Conn).Write(stdinToConn.b[:stdinToConn.n])
				}
			}

			if readError != nil {
				if readError == io.EOF {
					if !listenMode {
						abortConnToOutput <- struct{}{}
					}
				} else {
					abortConnToOutput <- struct{}{}
					stdinToConnError = readError
				}
				return
			}
			if sendError != nil {
				abortConnToOutput <- struct{}{}
				stdinToConnError = sendError
				return
			}
		}
	}()

	// copy packetConn to output
	go func() {
		defer wg.Done()

		var receiveError error
		for {
			var writeError error
			var deadlineError error
			var addr net.Addr

			if connToOutput.n > 0 {
				_, writeError = output.Write(connToOutput.b[:connToOutput.n])
			}
			if receiveError != nil && !errors.Is(receiveError, os.ErrDeadlineExceeded) {
				connToOutputError = receiveError
				return
			}
			if writeError != nil {
				connToOutputError = writeError
				return
			}

			select {
			case <-abortConnToOutput:
				return
			default:
			}

			deadlineError = packetConn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			if deadlineError != nil {
				connToOutputError = deadlineError
				return
			}

			connToOutput.n, addr, receiveError = packetConn.ReadFrom(connToOutput.b[:])
			if listenMode && connToOutput.n > 0 && addr.String() != clientAddr.String() {
				log.Printf("ignoring packet from %v", addr)
				connToOutput.n = 0
			}
		}
	}()

	wg.Wait()
	if stdinToConnError != nil {
		return stdinToConnError
	}
	return connToOutputError
}
