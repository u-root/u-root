// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import (
	"errors"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewConn(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:65000")
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		net  string
		mode transferMode
		addr *net.UDPAddr

		expectedAddr  *net.UDPAddr
		expectedMode  transferMode
		expectedError string
	}{
		"success": {
			net:  "udp",
			mode: ModeOctet,
			addr: addr,

			expectedAddr: addr,
			expectedMode: ModeOctet,
		},
		"error": {
			net:  "udp7",
			mode: ModeOctet,
			addr: addr,

			expectedError: "listen udp7 :0: unknown network udp7",
		},
	}

	for label, c := range cases {
		conn, err := newConn(c.net, c.mode, c.addr)

		// Errorf
		if err != nil && ErrorCause(err).Error() != c.expectedError {
			t.Errorf("%s: expected error %q, got %q", label, c.expectedError, ErrorCause(err).Error())
		}
		if err != nil {
			continue
		}

		// Addr
		if c.expectedAddr != conn.remoteAddr {
			t.Errorf("%s: Expected addr %#v, but it was %#v", label, c.expectedAddr, conn.remoteAddr)
		}

		// Mode
		if c.expectedMode != conn.mode {
			t.Errorf("%s: Expected mode %q, but it was %q", label, c.expectedMode, conn.mode)
		}
		conn.Close()

		// Defaults
		if conn.blksize != 512 {
			t.Errorf("%s: Expected blocksize to be default 512, but it was %d", label, conn.blksize)
		}
		if conn.timeout != time.Second {
			t.Errorf("%s: Expected timeout to be default 1s, but it was %s", label, conn.timeout)
		}
		if conn.windowsize != 1 {
			t.Errorf("%s: Expected window to be default 1, but it was %d", label, conn.windowsize)
		}
		if conn.retransmit != 10 {
			t.Errorf("%s: Expected retransmit to be default 1, but it was %d", label, conn.retransmit)
		}
		if len(conn.rx.buf) != 516 {
			t.Errorf("%s: Expected buf len to be default 516, but it was %d", label, len(conn.buf))
		}
	}
}

func testWriteConn(t *testing.T, conn *net.UDPConn, addr *net.UDPAddr, dg datagram) {
	if _, err := conn.WriteTo(dg.bytes(), addr); err != nil {
		t.Fatal(err)
	}
}

func TestConn_getAck(t *testing.T) {
	tDG := datagram{}

	cases := map[string]struct {
		timeout  time.Duration
		block    uint16
		window   uint16
		connFunc func(string, *net.UDPConn, *net.UDPAddr)

		expectedBlock   uint16
		expectedWindow  uint16
		expectedRingBuf int
		expectedError   string
	}{
		"success": {
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeAck(14)
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  "^$",
		},
		"timeout": {
			timeout:  time.Millisecond,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {},

			expectedError: "read udp .*: i/o timeout",
		},
		"wrong client": {
			timeout: time.Millisecond * 10,
			block:   67,
			window:  4,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				dg := datagram{buf: make([]byte, 516)}

				// Create and send a packet from a different port
				otherConn, err := net.ListenUDP("udp", nil)
				if err != nil {
					t.Fatal(err)
				}
				_, err = otherConn.WriteTo([]byte("anything"), sAddr)
				if err != nil {
					t.Fatal(err)
				}
				n, _, err := otherConn.ReadFrom(dg.buf)
				dg.offset = n
				if err != nil {
					t.Fatal(err)
				}

				// Result should be Unexpected TID
				if err := dg.validate(); err != nil {
					t.Errorf("wrong client: expected valid datagram: %v", err)
				}

				if dg.opcode() != opCodeERROR {
					t.Errorf("wrong client: expected opcode to be %s", opCodeERROR)
				}

				if dg.errorCode() != ErrCodeUnknownTransferID {
					t.Errorf("wrong client: expected error code to be %q", ErrCodeUnknownTransferID)
				}

				// Send correct ACK, the server should try again for a datagram from the correct client
				dg.writeAck(67)
				conn.WriteTo(dg.bytes(), sAddr)
			},

			expectedBlock:  67,
			expectedWindow: 4,
			expectedError:  "^$",
		},
		"invalid datagram": {
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeError(13, "error")
				tDG.offset = 5
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  `ACK validation`,
		},
		"error datagram": {
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeError(ErrCodeDiskFull, "error")
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  "error receiving ACK",
		},
		"other datagram": {
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeWriteReq("file", ModeNetASCII, nil)
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  "error receiving ACK.*unexpected datagram",
		},
		"incorrect block": {
			timeout: time.Second * 1,
			block:   18,
			window:  5,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeAck(14)
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:   14,
			expectedWindow:  0,
			expectedRingBuf: -4,
			expectedError:   "^$",
		},
		"incorrect block, ahead": {
			timeout: time.Second * 1,
			block:   18,
			window:  5,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeAck(20)
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  18,
			expectedWindow: 5,
			expectedError:  "^$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.timeout = c.timeout
		tConn.block = c.block
		tConn.window = c.window
		tConn.rx.buf = make([]byte, 516)
		tConn.txBuf = newRingBuffer(100, 100)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			c.connFunc(label, cNetConn, sAddr)
			wg.Done()
		}()

		err := tConn.getAck()
		wg.Wait()
		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}
		if err != nil {
			continue
		}

		// Block number
		if tConn.block != c.expectedBlock {
			t.Errorf("%s: Expected block %d, got %d", label, c.expectedBlock, tConn.block)
		}

		// Window number
		if tConn.window != c.expectedWindow {
			t.Errorf("%s: Expected window %d, got %d", label, c.expectedWindow, tConn.window)
		}

		// ringBuf
		if tConn.txBuf.current != c.expectedRingBuf {
			t.Errorf("%s: Expected ringBuf current %d, got %d", label, c.expectedRingBuf, tConn.txBuf.current)
		}
	}
}

func TestConn_sendWriteRequest(t *testing.T) {
	tDG := datagram{}

	cases := map[string]struct {
		timeout  time.Duration
		connFunc func(string, *net.UDPConn, *net.UDPAddr)

		expectedBlksize    uint16
		expectedTimeout    time.Duration
		expectedWindowsize uint16
		expectedTsize      *int64
		expectedBufLen     int
		expectedError      string
	}{
		"ACK": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeAck(0)
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     512,
			expectedError:      "^$",
		},
		"OACK, blksize 600": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optBlocksize: "600"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    600,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     600,
			expectedError:      "^$",
		},
		"OACK, timeout 2s": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optTimeout: "2"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second * 2,
			expectedWindowsize: 1,
			expectedBufLen:     512,
			expectedError:      "^$",
		},
		"OACK, window 10": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optWindowSize: "10"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 10,
			expectedBufLen:     512,
			expectedError:      "^$",
		},
		"OACK, tsize 1024": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optTransferSize: "1024"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     512,
			expectedTsize:      ptrInt64(1024),
			expectedError:      "^$",
		},
		"ERROR": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeError(ErrCodeFileNotFound, "error")
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^WRQ OACK response: remote error",
		},
		"OACK, invalid": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optTransferSize: "three"})
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^parsing OACK to WRQ",
		},
		"invalid datagram": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeReadReq("file", "error", nil)
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^sending WRQ: validating request response",
		},
		"other datagram": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeReadReq("file", ModeNetASCII, nil)
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^WRQ OACK response: unexpected datagram",
		},
		"no ack": {
			timeout:  time.Millisecond * 50,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {},

			expectedError: "^sending WRQ.*i/o timeout$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.timeout = c.timeout

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			c.connFunc(label, cNetConn, sAddr)
			wg.Done()
		}()

		err := tConn.sendWriteRequest("file", options{})
		wg.Wait()
		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}
		if err != nil {
			continue
		}

		if tConn.blksize != c.expectedBlksize {
			t.Errorf("%s: Expected blocksize to be %d, but it was %d", label, c.expectedBlksize, tConn.blksize)
		}
		if tConn.timeout != c.expectedTimeout {
			t.Errorf("%s: Expected timeout to be %s, but it was %s", label, c.expectedTimeout, tConn.timeout)
		}
		if tConn.windowsize != c.expectedWindowsize {
			t.Errorf("%s: Expected window to be %d, but it was %d", label, c.expectedWindowsize, tConn.windowsize)
		}
		if tConn.tsize != c.expectedTsize {
			if tConn.tsize == nil || c.expectedTsize == nil {
				t.Errorf("%s: Expected tsize to be %d, but it was %d", label, c.expectedTsize, tConn.tsize)
			} else if *tConn.tsize != *c.expectedTsize {
				t.Errorf("%s: Expected tsize to be %d, but it was %d", label, *c.expectedTsize, *tConn.tsize)
			}
		}
		if len(tConn.buf) != c.expectedBufLen {
			t.Errorf("%s: Expected buf len to be %d, but it was %d", label, c.expectedBufLen, len(tConn.buf))
		}
	}
}

func TestConn_sendReadRequest(t *testing.T) {
	tDG := datagram{}

	data := getTestData(t, "1MB-random")

	cases := map[string]struct {
		timeout  time.Duration
		mode     transferMode
		connFunc func(string, *net.UDPConn, *net.UDPAddr)

		expectedBuf        string
		expectDone         bool
		expectNetascii     bool
		expectedBlksize    uint16
		expectedTimeout    time.Duration
		expectedWindowsize uint16
		expectedTsize      *int64
		expectedBufLen     int
		expectedError      string
	}{
		"DATA, small": {
			timeout: time.Second,
			mode:    ModeOctet,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeData(1, []byte("data"))
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBuf:        "data",
			expectDone:         true,
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		"DATA, 512": {
			timeout: time.Second,
			mode:    ModeOctet,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeData(1, data[:512])
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBuf:        string(data[:512]),
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		"DATA, netascii": {
			timeout: time.Second,
			mode:    ModeNetASCII,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeData(1, []byte("data\ndata"))
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBuf:        "data\r\ndata",
			expectDone:         true,
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		"OACK, blksize 2048": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optBlocksize: "2048"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    2048,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     2052,
			expectedError:      "^$",
		},
		"OACK, timeout 2s": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optTimeout: "2"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second * 2,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		"OACK, window 10": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optWindowSize: "10"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 10,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		"OACK, tsize 1024": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optTransferSize: "1024"})
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedTsize:      ptrInt64(1024),
			expectedError:      "^$",
		},
		"OACK, invalid": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeOptionAck(options{optTransferSize: "three"})
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^got OACK, read setup",
		},
		"invalid datagram": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeReadReq("file", "error", nil)
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^sending RRQ: validating request response",
		},
		"other datagram": {
			timeout: time.Second,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeReadReq("file", ModeNetASCII, nil)
				testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^RRQ OACK response: unexpected datagram",
		},
		"no ack": {
			timeout:  time.Millisecond * 50,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {},

			expectedError: "^sending RRQ.*i/o timeout$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.timeout = c.timeout
		tConn.mode = c.mode

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			c.connFunc(label, cNetConn, sAddr)
			wg.Done()
		}()

		err := tConn.sendReadRequest("file", options{})
		wg.Wait()
		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}
		if err != nil {
			continue
		}

		// Flush buffer
		if tConn.netasciiEnc != nil {
			tConn.netasciiEnc.Flush()
		}

		if buf := tConn.rxBuf.String(); buf != c.expectedBuf {
			t.Errorf("%s: Expected buf to contain %q, but it was %q", label, c.expectedBuf, buf)
		}
		if tConn.done != c.expectDone {
			t.Errorf("%s: Expected done %t, but it wasn't", label, c.expectDone)
		}
		if tConn.blksize != c.expectedBlksize {
			t.Errorf("%s: Expected blocksize to be %d, but it was %d", label, c.expectedBlksize, tConn.blksize)
		}
		if tConn.timeout != c.expectedTimeout {
			t.Errorf("%s: Expected timeout to be %s, but it was %s", label, c.expectedTimeout, tConn.timeout)
		}
		if tConn.windowsize != c.expectedWindowsize {
			t.Errorf("%s: Expected window to be %d, but it was %d", label, c.expectedWindowsize, tConn.windowsize)
		}
		if tConn.tsize != c.expectedTsize {
			if tConn.tsize == nil || c.expectedTsize == nil {
				t.Errorf("%s: Expected tsize to be %d, but it was %d", label, c.expectedTsize, tConn.tsize)
			} else if *tConn.tsize != *c.expectedTsize {
				t.Errorf("%s: Expected tsize to be %d, but it was %d", label, *c.expectedTsize, *tConn.tsize)
			}
		}
		if len(tConn.rx.buf) != c.expectedBufLen {
			t.Errorf("%s: Expected buf len to be %d, but it was %d", label, c.expectedBufLen, len(tConn.rx.buf))
		}
	}
}

func TestConn_readData(t *testing.T) {
	tDG := datagram{}

	data := getTestData(t, "1MB-random")

	cases := map[string]struct {
		timeout  time.Duration
		window   uint16
		connFunc func(string, *net.UDPConn, *net.UDPAddr)

		expectedBlock  uint16
		expectedData   []byte
		expectedWindow uint16
		expectedError  string
	}{
		"success": {
			timeout: time.Second,
			window:  1,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeData(13, data[:512])
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  13,
			expectedWindow: 1,
			expectedData:   data[:512],
			expectedError:  "^$",
		},
		"1 retry": {
			timeout: time.Millisecond * 100,
			window:  56,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				time.Sleep(110 * time.Millisecond)
				tDG.writeData(13, data[:512])
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  13,
			expectedWindow: 0, // reset to 0, +1
			expectedData:   data[:512],
			expectedError:  "^$",
		},
		"invalid": {
			timeout: time.Millisecond * 100,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeData(13, data[:512])
				tDG.offset = 3
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^validating read data: Corrupt block number$",
		},
		"error datagram": {
			timeout: time.Millisecond * 100,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeError(ErrCodeDiskFull, "error")
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^reading data: remote error:",
		},
		"other datagram": {
			timeout: time.Millisecond * 100,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				tDG.writeAck(12)
				testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^read data response: unexpected datagram:",
		},
		"no data": {
			timeout:  time.Millisecond * 10,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {},

			expectedError: "^reading data.*i/o timeout$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.timeout = c.timeout
		tConn.window = c.window

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			c.connFunc(label, cNetConn, sAddr)
			wg.Done()
		}()

		err := tConn.readData()
		wg.Wait()
		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
			continue
		}

		// Data
		if string(tConn.rx.data()) != string(c.expectedData) {
			// t.Errorf("%s: Expected data %q, got %q", label, string(c.expectedData), string(data))
		}

		// Block number
		if tConn.rx.block() != c.expectedBlock {
			t.Errorf("%s: Expected block %d, got %d", label, c.expectedBlock, tConn.block)
		}

		// Window number
		if tConn.window != c.expectedWindow {
			t.Errorf("%s: Expected window %d, got %d", label, c.expectedWindow, tConn.window)
		}
	}
}

func TestConn_ackData(t *testing.T) {
	tDG := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := map[string]struct {
		timeout    time.Duration
		rx         datagram
		block      uint16
		window     uint16
		windowsize uint16
		catchup    bool
		connFunc   func(string, *net.UDPConn, *net.UDPAddr)

		expectCatchup  bool
		expectedBlock  uint16
		expectedWindow uint16
		expectedError  string
	}{
		"success, reached window": {
			timeout:    time.Second,
			block:      12,
			windowsize: 1,
			window:     0,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(13, data[:512])
				return dg
			}(),

			expectedBlock:  13,
			expectedWindow: 0,
			expectedError:  "^$",
		},
		"success, reset catchup": {
			timeout:    time.Second,
			block:      12,
			windowsize: 4,
			window:     0,
			catchup:    true,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(13, data[:512])
				return dg
			}(),

			expectedBlock:  13,
			expectedWindow: 1,
			expectedError:  "^$",
		},
		"repeat block": {
			timeout:    time.Second,
			block:      12,
			windowsize: 2,
			window:     1,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(12, data[:512])
				return dg
			}(),

			expectedBlock:  12,
			expectedWindow: 1,
			expectedError:  errBlockSequence.Error(),
		},
		"future block, no catchup": {
			timeout:    time.Second,
			block:      12,
			windowsize: 2,
			window:     1,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(14, data[:512])
				return dg
			}(),
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, _, err := conn.ReadFrom(tDG.buf)
				if err != nil {
					t.Errorf("future block, no catchup: expected ACK %v", err)
					return
				}

				if tDG.block() != 12 {
					t.Errorf("future block, no catchup: expected ACK with block 12, got %d", tDG.block())
				}
			},

			expectCatchup:  true,
			expectedBlock:  12,
			expectedWindow: 0,
			expectedError:  errBlockSequence.Error(),
		},
		"future block, rollover": {
			timeout:    time.Second,
			block:      65534,
			windowsize: 4,
			window:     1,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(0, data[:512])
				return dg
			}(),
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, _, err := conn.ReadFrom(tDG.buf)
				if err != nil {
					t.Errorf("future block, no catchup: expected ACK %v", err)
					return
				}

				if tDG.block() != 65534 {
					t.Errorf("future block, no catchup: expected ACK with block 65534, got %d", tDG.block())
				}
			},

			expectCatchup:  true,
			expectedBlock:  65534,
			expectedWindow: 0,
			expectedError:  errBlockSequence.Error(),
		},
		"future block, catchup": {
			timeout:    time.Second,
			block:      12,
			windowsize: 32,
			window:     1,
			catchup:    true,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(24, data[:512])
				return dg
			}(),

			expectCatchup:  true,
			expectedBlock:  12,
			expectedWindow: 1,
			expectedError:  errBlockSequence.Error(),
		},
		"past block": {
			timeout:    time.Second,
			block:      12,
			windowsize: 1,
			window:     1,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(1, data[:512])
				return dg
			}(),

			expectedBlock:  12,
			expectedWindow: 1,
			expectedError:  errBlockSequence.Error(),
		},
		"success, below window": {
			timeout:    time.Second,
			block:      12,
			window:     1,
			windowsize: 4,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(13, data[:512])
				return dg
			}(),

			expectedBlock:  13,
			expectedWindow: 2,
			expectedError:  "^$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.rx = c.rx
		tConn.timeout = c.timeout
		tConn.block = c.block
		tConn.window = c.window
		tConn.windowsize = c.windowsize
		tConn.catchup = c.catchup

		err := tConn.ackData()
		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}

		if c.connFunc != nil {
			c.connFunc(label, cNetConn, sAddr)
		}

		// Block number
		if tConn.block != c.expectedBlock {
			t.Errorf("%s: Expected block %d, got %d", label, c.expectedBlock, tConn.block)
		}

		// Window number
		if tConn.window != c.expectedWindow {
			t.Errorf("%s: Expected window %d, got %d", label, c.expectedWindow, tConn.window)
		}
		// Catchup
		if tConn.catchup != c.expectCatchup {
			t.Errorf("%s: Expected catchup %t, but it wasn't", label, c.expectCatchup)
		}
	}
}

func TestConn_parseOptions(t *testing.T) {
	dg := datagram{}

	cases := map[string]struct {
		rx       func() datagram
		tsize    *int64
		isSender bool

		expectOptionsParsed bool
		expectedOptions     options
		expectedBlksize     uint16
		expectedTimeout     time.Duration
		expectedWindowsize  uint16
		expectedTsize       *int64
		expectedError       string
	}{
		"blocksize, valid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optBlocksize: "234"})
				return dg
			},

			expectOptionsParsed: true,
			expectedOptions:     options{optBlocksize: "234"},
			expectedBlksize:     234,
			expectedError:       "^$",
		},
		"blocksize, invalid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optBlocksize: "a"})
				return dg
			},

			expectOptionsParsed: false,
			expectedBlksize:     0,
			expectedError:       `error parsing .* for option "blksize"`,
		},
		"timeout, valid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optTimeout: "3"})
				return dg
			},

			expectedOptions:     options{optTimeout: "3"},
			expectOptionsParsed: true,
			expectedTimeout:     3 * time.Second,
			expectedError:       `^$`,
		},
		"timeout, invalid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optTimeout: "three"})
				return dg
			},

			expectOptionsParsed: false,
			expectedTimeout:     0,
			expectedError:       `error parsing .* for option "timeout"`,
		},
		"tsize, valid, sending side": {
			rx: func() datagram {
				dg.writeOptionAck(options{optTransferSize: "0"})
				return dg
			},
			tsize:    ptrInt64(1000),
			isSender: true,

			expectedOptions:     options{optTransferSize: "1000"},
			expectedTsize:       ptrInt64(1000),
			expectOptionsParsed: true,
			expectedError:       `^$`,
		},
		"tsize, valid, receive side": {
			rx: func() datagram {
				dg.writeOptionAck(options{optTransferSize: "42"})
				return dg
			},

			expectedOptions:     options{},
			expectOptionsParsed: true,
			expectedTsize:       ptrInt64(42),
			expectedError:       `^$`,
		},
		"tsize, invalid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optTransferSize: "large"})
				return dg
			},

			expectedError: `^error parsing .* for option "tsize"$`,
		},
		"windowsize, valid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optWindowSize: "32"})
				return dg
			},

			expectedOptions:     options{optWindowSize: "32"},
			expectOptionsParsed: true,
			expectedWindowsize:  32,
			expectedError:       `^$`,
		},
		"windowsize, invalid": {
			rx: func() datagram {
				dg.writeOptionAck(options{optWindowSize: "x"})
				return dg
			},

			expectedError: `^error parsing .* for option "windowsize"$`,
		},
		"all options, sending side": {
			rx: func() datagram {
				dg.writeOptionAck(options{
					optBlocksize:    "1024",
					optTimeout:      "3",
					optTransferSize: "0",
					optWindowSize:   "16",
				})
				return dg
			},
			tsize:    ptrInt64(1234567890),
			isSender: true,

			expectedOptions: options{
				optBlocksize:    "1024",
				optTimeout:      "3",
				optTransferSize: "1234567890",
				optWindowSize:   "16",
			},
			expectOptionsParsed: true,
			expectedBlksize:     1024,
			expectedTimeout:     3 * time.Second,
			expectedTsize:       ptrInt64(1234567890),
			expectedWindowsize:  16,
		},
		"all options, receive side": {
			rx: func() datagram {
				dg.writeOptionAck(options{
					optBlocksize:    "1024",
					optTimeout:      "3",
					optTransferSize: "1234567890",
					optWindowSize:   "16",
				})
				return dg
			},

			expectedOptions: options{
				optBlocksize:  "1024",
				optTimeout:    "3",
				optWindowSize: "16",
			},
			expectOptionsParsed: true,
			expectedBlksize:     1024,
			expectedTimeout:     3 * time.Second,
			expectedTsize:       ptrInt64(1234567890),
			expectedWindowsize:  16,
		},
	}

	for label, c := range cases {
		tConn := conn{rx: c.rx()}
		tConn.tsize = c.tsize
		tConn.isSender = c.isSender

		opts, err := tConn.parseOptions()

		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}

		// Options
		if !reflect.DeepEqual(c.expectedOptions, opts) {
			t.Errorf("%s: Expected options %q, got %q", label, c.expectedOptions, opts)
		}

		// OptionsParsed
		if c.expectOptionsParsed != tConn.optionsParsed {
			t.Errorf("%s: Expected optionsParsed %t, but it wasn't", label, c.expectOptionsParsed)
		}

		if tConn.blksize != c.expectedBlksize {
			t.Errorf("%s: Expected blocksize to be %d, but it was %d", label, c.expectedBlksize, tConn.blksize)
		}
		if tConn.timeout != c.expectedTimeout {
			t.Errorf("%s: Expected timeout to be %s, but it was %s", label, c.expectedTimeout, tConn.timeout)
		}
		if tConn.windowsize != c.expectedWindowsize {
			t.Errorf("%s: Expected window to be %d, but it was %d", label, c.expectedWindowsize, tConn.windowsize)
		}
		if tConn.tsize != c.expectedTsize {
			if tConn.tsize == nil || c.expectedTsize == nil {
				t.Errorf("%s: Expected tsize to be *%d, but it was *%d", label, c.expectedTsize, tConn.tsize)
			} else if *tConn.tsize != *c.expectedTsize {
				t.Errorf("%s: Expected tsize to be %d, but it was %d", label, *c.expectedTsize, *tConn.tsize)
			}
		}
	}
}

func TestConn_write(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := map[string]struct {
		bytes         []byte
		optionsParsed bool
		blksize       uint16
		window        uint16
		windowsize    uint16
		rx            func() datagram
		timeout       time.Duration
		connFunc      func(label string, conn *net.UDPConn, sAddr *net.UDPAddr)
		connErr       error

		expectedCount  int
		expectedError  string
		expectedWindow uint64
	}{
		"writeSetup fails": {
			timeout:       time.Millisecond,
			optionsParsed: false,
			rx: func() datagram {
				dg.writeOptionAck(options{optBlocksize: "234"})
				return dg
			},

			expectedError: "parsing options before write: write setup: network read failed:",
		},
		"success, buf < blksize": {
			timeout:       time.Millisecond,
			bytes:         data[:300],
			blksize:       512,
			optionsParsed: true,

			expectedCount: 300,
			expectedError: "^$",
		},
		"success, buf > blksize, window 1": {
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    1,
			optionsParsed: true,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				dg.writeAck(1)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, _ = conn.WriteTo(dg.buf, sAddr)
				dg.writeAck(2)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, _ = conn.WriteTo(dg.buf, sAddr)
			},

			expectedCount: 1024,
			expectedError: "^$",
		},
		"success, buf > blksize, window 2": {
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    2,
			optionsParsed: true,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				dg.writeAck(1)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, _ = conn.WriteTo(dg.buf, sAddr)
			},

			expectedCount: 1024,
			expectedError: "^$",
		},
		"fail to ack": {
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    1,
			optionsParsed: true,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
			},

			expectedCount: 1024,
			expectedError: "^receiving ACK after writing data: network read failed",
		},
		"conn err": {
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    1,
			optionsParsed: true,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				dg.writeAck(1)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, _ = conn.WriteTo(dg.buf, sAddr)
			},
			connErr: errors.New("conn error"),

			expectedCount: 1024,
			expectedError: "writing data: conn error",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		if c.rx != nil {
			tConn.rx = c.rx()
		}
		tConn.blksize = c.blksize
		tConn.window = c.window
		tConn.windowsize = c.windowsize
		tConn.optionsParsed = c.optionsParsed
		tConn.timeout = c.timeout
		tConn.buf = make([]byte, c.blksize)
		tConn.txBuf = newRingBuffer(int(c.windowsize), int(c.blksize))
		tConn.err = c.connErr

		var wg sync.WaitGroup
		if c.connFunc != nil {
			wg.Add(1)
			go func() {
				c.connFunc(label, cNetConn, sAddr)
				wg.Done()
			}()
		}

		count, err := tConn.write(c.bytes)
		wg.Wait()

		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}

		// Count
		if c.expectedCount != count {
			t.Errorf("%s: expected count %d, got %d", label, c.expectedCount, count)
		}
	}
}

func TestConn_Close(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := map[string]struct {
		bytes    []byte
		blksize  uint16
		timeout  time.Duration
		connFunc func(label string, conn *net.UDPConn, sAddr *net.UDPAddr)
		connErr  error

		expectedError string
	}{
		"conn err": {
			connErr: errors.New("conn error"),

			expectedError: "checking conn err before Close: conn error",
		},
		"success, no data": {
			timeout: time.Millisecond * 100,
			bytes:   []byte{},
			blksize: 512,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				n, _, err := conn.ReadFrom(dg.buf)
				if err != nil {
					t.Fatal(err)
				}
				dg.offset = n

				if dg.opcode() != opCodeDATA {
					t.Errorf("%s: Expected opcode %s, got %s", label, opCodeDATA, dg.opcode())
				}

				if l := len(dg.data()); l != 0 {
					t.Errorf("%s: Expected data len to be 0, but they were %d", label, l)
				}

				dg.writeAck(1)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, err = conn.WriteTo(dg.buf, sAddr)
				if err != nil {
					t.Fatal(err)
				}
			},

			expectedError: "^$",
		},
		"success, with data": {
			timeout: time.Millisecond * 100,
			bytes:   data[:384],
			blksize: 512,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				n, _, err := conn.ReadFrom(dg.buf)
				if err != nil {
					t.Fatal(err)
				}
				dg.offset = n

				if dg.opcode() != opCodeDATA {
					t.Errorf("%s: Expected opcode %s, got %s", label, opCodeDATA, dg.opcode())
				}

				if l := len(dg.data()); l != 384 {
					t.Errorf("%s: Expected data len to be 384, but they were %d", label, l)
				}

				dg.writeAck(1)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				_, err = conn.WriteTo(dg.buf, sAddr)
				if err != nil {
					t.Fatal(err)
				}
			},

			expectedError: "^$",
		},
		"timeout": {
			timeout:  time.Millisecond * 100,
			blksize:  512,
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {},

			expectedError: "^writing final data ACK before Close: network read failed: .* i/o timeout$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.blksize = c.blksize
		tConn.timeout = c.timeout
		tConn.buf = make([]byte, c.blksize)
		tConn.txBuf = newRingBuffer(1, int(c.blksize))
		tConn.txBuf.Write(c.bytes)
		tConn.err = c.connErr

		var wg sync.WaitGroup
		if c.connFunc != nil {
			wg.Add(1)
			go func() {
				c.connFunc(label, cNetConn, sAddr)
				wg.Done()
			}()
		}

		err := tConn.Close()
		wg.Wait()

		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}
	}
}

func TestConn_read(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := map[string]struct {
		bytes         []byte
		blksize       uint16
		windowsize    uint16
		optionsParsed bool
		timeout       time.Duration
		connFunc      func(label string, conn *net.UDPConn, sAddr *net.UDPAddr)
		connErr       error

		expectedRead  int
		expectedError string
	}{
		"success": {
			timeout:       time.Millisecond * 100,
			optionsParsed: true,
			blksize:       512,
			bytes:         make([]byte, 512),
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				dg.writeData(1, data[:512])
				conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				if err != nil {
					t.Fatal(err)
				}
			},

			expectedRead:  512,
			expectedError: "^$",
		},
		"success, EOF": {
			timeout:       time.Millisecond * 100,
			optionsParsed: true,
			blksize:       512,
			bytes:         make([]byte, 512),
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				dg.writeData(1, data[:300])
				conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				if err != nil {
					t.Fatal(err)
				}
			},

			expectedRead:  300,
			expectedError: "^EOF$",
		},
		"block sequence error": {
			timeout:       time.Millisecond * 100,
			optionsParsed: true,
			blksize:       512,
			windowsize:    2,
			bytes:         make([]byte, 512),
			connFunc: func(label string, conn *net.UDPConn, sAddr *net.UDPAddr) {
				// Write wrong block
				dg.writeData(2, data[:512])
				conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				if err != nil {
					t.Fatal(err)
				}

				// Receive ACK for previous
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				n, _, err := conn.ReadFrom(dg.buf)
				if err != nil {
					t.Fatal(err)
				}
				dg.offset = n
				if dg.block() != 0 {
					t.Errorf("%s: expected block 0 again, got %d", label, dg.block())
				}

				// Write correct data
				dg.writeData(1, data[:300])
				conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
				_, err = conn.WriteTo(dg.bytes(), sAddr)
				if err != nil {
					t.Fatal(err)
				}
			},

			expectedRead:  300,
			expectedError: "^EOF$",
		},
	}

	for label, c := range cases {
		tConn, sAddr, cNetConn, closer := testConns(t)
		defer closer()
		tConn.optionsParsed = c.optionsParsed
		tConn.blksize = c.blksize
		tConn.timeout = c.timeout
		tConn.windowsize = c.windowsize
		tConn.buf = make([]byte, c.blksize+4)
		tConn.err = c.connErr

		var wg sync.WaitGroup
		if c.connFunc != nil {
			wg.Add(1)
			go func() {
				c.connFunc(label, cNetConn, sAddr)
				wg.Done()
			}()
		}

		read, err := tConn.read(c.bytes)
		wg.Wait()

		// Error
		if err != nil {
			if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
				t.Errorf("%s: expected error %q, got %q", label, c.expectedError, err.Error())
			}
		}

		// Read Count
		if c.expectedRead != read {
			t.Errorf("%s: Expected read bytes to be %d, but it was %d", label, c.expectedRead, read)
		}
	}
}

func TestConn_sendError(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	cases := map[string]struct {
		code    ErrorCode
		msg     string
		blksize uint16
		timeout time.Duration

		expectedCode  ErrorCode
		expectedError string
	}{
		"message, undersize": {
			timeout: time.Millisecond * 100,
			blksize: 512,
			code:    ErrCodeNoSuchUser,
			msg:     "foo",

			expectedCode:  ErrCodeNoSuchUser,
			expectedError: "foo",
		},
		"message, oversize": {
			timeout: time.Millisecond * 100,
			blksize: 10,
			code:    ErrCodeNoSuchUser,
			msg:     "there was a long error",

			expectedCode:  ErrCodeNoSuchUser,
			expectedError: "there was",
		},
	}

	for label, c := range cases {
		tConn, _, cNetConn, closer := testConns(t)
		defer closer()
		tConn.blksize = c.blksize
		tConn.timeout = c.timeout
		tConn.buf = make([]byte, c.blksize+4)

		tConn.sendError(c.code, c.msg)

		// Receive Error
		cNetConn.SetReadDeadline(time.Now().Add(c.timeout))
		n, _, err := cNetConn.ReadFrom(dg.buf)
		if err != nil {
			t.Fatal(err)
		}
		dg.offset = n

		// Error Code
		if c.expectedCode != dg.errorCode() {
			t.Errorf("%s: expected errorCode %s, got %s", label, c.expectedCode, dg.errorCode())
		}

		// Error Message
		if c.expectedError != dg.errMsg() {
			t.Errorf("%s: expected message %q, got %q", label, c.expectedError, dg.errMsg())
		}
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}

func testConns(t *testing.T) (tConn *conn, sAddr *net.UDPAddr, cNetConn *net.UDPConn, closer func()) {
	cNetConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		t.Fatal(err)
	}
	cPort := cNetConn.LocalAddr().(*net.UDPAddr).Port
	cAddr, _ := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(cPort))

	tConn, err = newConn("udp", ModeOctet, cAddr)
	if err != nil {
		t.Fatal(err)
	}

	sPort := tConn.netConn.LocalAddr().(*net.UDPAddr).Port
	sAddr, _ = net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(sPort))

	closer = func() {
		cNetConn.Close()
		tConn.netConn.Close()
	}

	return
}
