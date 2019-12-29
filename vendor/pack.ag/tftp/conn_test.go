// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"regexp"
	"runtime"
	"testing"
	"time"
)

const testConnTimeout = 500 * time.Millisecond

func TestNewConn(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:65000")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		net  string
		mode TransferMode
		addr *net.UDPAddr

		expectedAddr  *net.UDPAddr
		expectedMode  TransferMode
		expectedError string
	}{
		{
			name: "success",
			net:  "udp",
			mode: ModeOctet,
			addr: addr,

			expectedAddr: addr,
			expectedMode: ModeOctet,
		},
		{
			name: "error",
			net:  "udp7",
			mode: ModeOctet,
			addr: addr,

			expectedError: "listen udp7 :0: unknown network udp7",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, err := newConn(c.net, c.mode, c.addr)

			// Errorf
			if err != nil && ErrorCause(err).Error() != c.expectedError {
				t.Errorf("expected error %q, got %q", c.expectedError, ErrorCause(err).Error())
			}
			if err != nil {
				return
			}

			// Addr
			if c.expectedAddr != conn.remoteAddr {
				t.Errorf("expected addr %#v, but it was %#v", c.expectedAddr, conn.remoteAddr)
			}

			// Mode
			if c.expectedMode != conn.mode {
				t.Errorf("expected mode %q, but it was %q", c.expectedMode, conn.mode)
			}
			conn.Close()

			// Defaults
			if conn.blksize != 512 {
				t.Errorf("expected blocksize to be default 512, but it was %d", conn.blksize)
			}
			if conn.timeout != time.Second {
				t.Errorf("expected timeout to be default 1s, but it was %s", conn.timeout)
			}
			if conn.windowsize != 1 {
				t.Errorf("expected window to be default 1, but it was %d", conn.windowsize)
			}
			if conn.retransmit != 10 {
				t.Errorf("expected retransmit to be default 1, but it was %d", conn.retransmit)
			}
			if len(conn.rx.buf) != 516 {
				t.Errorf("expected buf len to be default 516, but it was %d", len(conn.buf))
			}
		})
	}
}

func testWriteConn(t *testing.T, conn *net.UDPConn, addr *net.UDPAddr, dg datagram) error {
	conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
	_, err := conn.WriteTo(dg.bytes(), addr)
	return err
}

func testConnFunc(conn *net.UDPConn, addr *net.UDPAddr, connFunc func(*net.UDPConn, *net.UDPAddr) error) chan error {
	errChan := make(chan error)
	if connFunc != nil {
		go func() {
			errChan <- connFunc(conn, addr)
		}()
	} else {
		close(errChan)
	}
	return errChan
}

func TestConn_getAck(t *testing.T) {
	tDG := datagram{}

	cases := []struct {
		name     string
		timeout  time.Duration
		block    uint16
		window   uint16
		connFunc func(*net.UDPConn, *net.UDPAddr) error

		expectedBlock   uint16
		expectedWindow  uint16
		expectedRingBuf int
		expectedError   string
	}{
		{
			name:    "success",
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeAck(14)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  "^$",
		},
		{
			name:    "timeout",
			timeout: time.Millisecond,

			expectedError: "read .*: i/o timeout",
		},
		{
			name:    "wrong client",
			timeout: time.Millisecond * 10,
			block:   67,
			window:  4,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				dg := datagram{buf: make([]byte, 516)}

				// Create and send a packet from a different port
				otherConn, err := net.ListenUDP("udp", nil)
				if err != nil {
					return err
				}
				otherConn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err = otherConn.WriteTo([]byte("anything"), sAddr)
				if err != nil {
					return err
				}
				otherConn.SetReadDeadline(time.Now().Add(testConnTimeout))
				n, _, err := otherConn.ReadFrom(dg.buf)
				if err != nil {
					return err
				}
				dg.offset = n

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
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err = conn.WriteTo(dg.bytes(), sAddr)
				return err
			},

			expectedBlock:  67,
			expectedWindow: 4,
			expectedError:  "^$",
		},
		{
			name:    "invalid datagram",
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeError(13, "error")
				tDG.offset = 5
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  `ACK validation`,
		},
		{
			name:    "error datagram",
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeError(ErrCodeDiskFull, "error")
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  "error receiving ACK",
		},
		{
			name:    "other datagram",
			timeout: time.Second * 1,
			block:   14,
			window:  5,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeWriteReq("file", ModeNetASCII, nil)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  14,
			expectedWindow: 5,
			expectedError:  "error receiving ACK.*unexpected datagram",
		},
		{
			name:    "incorrect block",
			timeout: time.Second * 1,
			block:   18,
			window:  5,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeAck(14)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:   14,
			expectedWindow:  0,
			expectedRingBuf: -4,
			expectedError:   "^$",
		},
		{
			name:    "incorrect block, ahead",
			timeout: time.Second * 1,
			block:   18,
			window:  5,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeAck(20)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  18,
			expectedWindow: 5,
			expectedError:  "^$",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.timeout = c.timeout
			tConn.block = c.block
			tConn.window = c.window
			tConn.rx.buf = make([]byte, 516)
			tConn.txBuf = newRingBuffer(100, 100)
			tConn.tx.writeAck(1) // TODO: set prev opcode in test, needs to be done when checking for OACK

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			_ = tConn.getAck() // TODO: check return func
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if tConn.err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, tConn.err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, tConn.err.Error())
				}
			}
			if tConn.err != nil {
				return
			}

			// Block number
			if tConn.block != c.expectedBlock {
				t.Errorf("expected block %d, got %d", c.expectedBlock, tConn.block)
			}

			// Window number
			if tConn.window != c.expectedWindow {
				t.Errorf("expected window %d, got %d", c.expectedWindow, tConn.window)
			}

			// ringBuf
			if tConn.txBuf.current != c.expectedRingBuf {
				t.Errorf("expected ringBuf current %d, got %d", c.expectedRingBuf, tConn.txBuf.current)
			}
		})
	}
}

func TestConn_sendWriteRequest(t *testing.T) {
	tDG := datagram{}

	cases := []struct {
		name     string
		timeout  time.Duration
		connFunc func(*net.UDPConn, *net.UDPAddr) error

		expectedBlksize    uint16
		expectedTimeout    time.Duration
		expectedWindowsize uint16
		expectedTsize      *int64
		expectedBufLen     int
		expectedError      string
	}{
		{
			name:    "ACK",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeAck(0)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     512,
			expectedError:      "^$",
		},
		{
			name:    "OACK, blksize 600",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optBlocksize: "600"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    600,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     600,
			expectedError:      "^$",
		},
		{
			name:    "OACK, timeout 2s",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optTimeout: "2"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second * 2,
			expectedWindowsize: 1,
			expectedBufLen:     512,
			expectedError:      "^$",
		},
		{
			name:    "OACK, window 10",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optWindowSize: "10"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 10,
			expectedBufLen:     512,
			expectedError:      "^$",
		},
		{
			name:    "OACK, tsize 1024",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optTransferSize: "1024"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     512,
			expectedTsize:      ptrInt64(1024),
			expectedError:      "^$",
		},
		{
			name:    "ERROR",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeError(ErrCodeFileNotFound, "error")
				return testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^WRQ OACK response: remote error",
		},
		{
			name:    "OACK, invalid",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optTransferSize: "three"})
				return testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^parsing options",
		},
		{
			name:    "invalid datagram",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeReadReq("file", "error", nil)
				return testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^validating request response",
		},
		{
			name:    "other datagram",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeReadReq("file", ModeNetASCII, nil)
				return testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^WRQ OACK response: unexpected datagram",
		},
		{
			name:    "no ack",
			timeout: time.Millisecond * 50,

			expectedError: "^receiving request response:.*i/o timeout$",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.timeout = c.timeout

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			err := tConn.sendWriteRequest("file", options{})
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}
			if err != nil {
				return
			}

			if tConn.blksize != c.expectedBlksize {
				t.Errorf("expected blocksize to be %d, but it was %d", c.expectedBlksize, tConn.blksize)
			}
			if tConn.timeout != c.expectedTimeout {
				t.Errorf("expected timeout to be %s, but it was %s", c.expectedTimeout, tConn.timeout)
			}
			if tConn.windowsize != c.expectedWindowsize {
				t.Errorf("expected window to be %d, but it was %d", c.expectedWindowsize, tConn.windowsize)
			}
			if tConn.tsize != c.expectedTsize {
				if tConn.tsize == nil || c.expectedTsize == nil {
					t.Errorf("expected tsize to be %d, but it was %d", c.expectedTsize, tConn.tsize)
				} else if *tConn.tsize != *c.expectedTsize {
					t.Errorf("expected tsize to be %d, but it was %d", *c.expectedTsize, *tConn.tsize)
				}
			}
			if len(tConn.buf) != c.expectedBufLen {
				t.Errorf("expected buf len to be %d, but it was %d", c.expectedBufLen, len(tConn.buf))
			}
		})
	}
}

func TestConn_sendReadRequest(t *testing.T) {
	tDG := datagram{}

	data := getTestData(t, "1MB-random")

	cases := []struct {
		name        string
		timeout     time.Duration
		mode        TransferMode
		connFunc    func(*net.UDPConn, *net.UDPAddr) error
		windowsOnly bool
		nixOnly     bool

		skip string

		expectedBuf        string
		expectNetascii     bool
		expectedBlksize    uint16
		expectedTimeout    time.Duration
		expectedWindowsize uint16
		expectedTsize      *int64
		expectedBufLen     int
		expectedError      string
	}{
		{
			name:    "DATA, small",
			timeout: time.Second,
			mode:    ModeOctet,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeData(1, []byte("data"))
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBuf:        "data",
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		{
			name:    "DATA, 512",
			timeout: time.Second,
			mode:    ModeOctet,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeData(1, data[:512])
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBuf:        string(data[:512]),
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		{
			name:    "DATA, netascii",
			timeout: time.Second,
			mode:    ModeNetASCII,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeData(1, []byte("data\r\ndata"))
				return testWriteConn(t, conn, sAddr, tDG)
			},
			nixOnly: true,

			expectedBuf:        "data\ndata", // Writes in as netascii, read out normal
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		{
			name:    "DATA, netascii",
			timeout: time.Second,
			mode:    ModeNetASCII,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeData(1, []byte("data\r\ndata"))
				return testWriteConn(t, conn, sAddr, tDG)
			},
			windowsOnly: true,

			expectedBuf:        "data\r\ndata", // Writes in as netascii, read out normal
			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		{
			name:    "OACK, blksize 2048",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optBlocksize: "2048"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    2048,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     2052,
			expectedError:      "^$",
		},
		{
			name:    "OACK, timeout 2s",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optTimeout: "2"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second * 2,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		{
			name:    "OACK, window 10",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optWindowSize: "10"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 10,
			expectedBufLen:     516,
			expectedError:      "^$",
		},
		{
			name:    "OACK, tsize 1024",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optTransferSize: "1024"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlksize:    512,
			expectedTimeout:    time.Second,
			expectedWindowsize: 1,
			expectedBufLen:     516,
			expectedTsize:      ptrInt64(1024),
			expectedError:      "^$",
		},
		{
			name:    "OACK, invalid",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeOptionAck(options{optTransferSize: "three"})
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "read setup: error parsing \"three\" for option \"tsize\"",
		},
		{
			name:    "invalid datagram",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeReadReq("file", "error", nil)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^validating request response",
		},
		{
			name:    "other datagram",
			timeout: time.Second,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeReadReq("file", ModeNetASCII, nil)
				return testWriteConn(t, conn, sAddr, tDG)
			},
			expectedError: "^RRQ OACK response: unexpected datagram",
		},
		{
			name:    "no ack",
			timeout: time.Millisecond * 50,

			expectedError: "^receiving request response:.*i/o timeout$",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.skip != "" {
				t.Skip(c.skip)
			}
			if c.windowsOnly && runtime.GOOS != "windows" {
				t.Skip("widows only")
			}
			if c.nixOnly && runtime.GOOS == "windows" {
				t.Skip("*nix only")
			}

			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.timeout = c.timeout
			tConn.mode = c.mode

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			err := tConn.sendReadRequest("file", options{})
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}
			if err != nil {
				return
			}

			// Flush buffer
			tConn.Close()

			if buf, _ := ioutil.ReadAll(tConn.reader); string(buf) != c.expectedBuf {
				t.Errorf("expected buf to contain %q, but it was %q", c.expectedBuf, buf)
			}
			if tConn.blksize != c.expectedBlksize {
				t.Errorf("expected blocksize to be %d, but it was %d", c.expectedBlksize, tConn.blksize)
			}
			if tConn.timeout != c.expectedTimeout {
				t.Errorf("expected timeout to be %s, but it was %s", c.expectedTimeout, tConn.timeout)
			}
			if tConn.windowsize != c.expectedWindowsize {
				t.Errorf("expected window to be %d, but it was %d", c.expectedWindowsize, tConn.windowsize)
			}
			if tConn.tsize != c.expectedTsize {
				if tConn.tsize == nil || c.expectedTsize == nil {
					t.Errorf("expected tsize to be %d, but it was %d", c.expectedTsize, tConn.tsize)
				} else if *tConn.tsize != *c.expectedTsize {
					t.Errorf("expected tsize to be %d, but it was %d", *c.expectedTsize, *tConn.tsize)
				}
			}
			if len(tConn.rx.buf) != c.expectedBufLen {
				t.Errorf("expected buf len to be %d, but it was %d", c.expectedBufLen, len(tConn.rx.buf))
			}
		})
	}
}

func TestConn_readData(t *testing.T) {
	tDG := datagram{}

	data := getTestData(t, "1MB-random")

	cases := []struct {
		name     string
		timeout  time.Duration
		window   uint16
		connFunc func(*net.UDPConn, *net.UDPAddr) error

		skip string

		expectedBlock  uint16
		expectedData   []byte
		expectedWindow uint16
		expectedError  string
	}{
		{
			name:    "success",
			timeout: time.Second,
			window:  1,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeData(13, data[:512])
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedBlock:  13,
			expectedWindow: 1,
			expectedData:   data[:512],
			expectedError:  "^$",
		},
		{
			name:    "1 retry",
			timeout: time.Millisecond * 100,
			window:  56,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				time.Sleep(110 * time.Millisecond)
				tDG.writeData(13, data[:512])
				return testWriteConn(t, conn, sAddr, tDG)
			},

			skip: "need to cycle state",

			expectedBlock:  13,
			expectedWindow: 0, // reset to 0, +1
			expectedData:   data[:512],
			expectedError:  "^$",
		},
		{
			name:    "invalid",
			timeout: time.Millisecond * 100,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeData(13, data[:512])
				tDG.offset = 3
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^validating read data: Corrupt block number$",
		},
		{
			name:    "error datagram",
			timeout: time.Millisecond * 100,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeError(ErrCodeDiskFull, "error")
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^reading data: remote error:",
		},
		{
			name:    "other datagram",
			timeout: time.Millisecond * 100,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				tDG.writeAck(12)
				return testWriteConn(t, conn, sAddr, tDG)
			},

			expectedError: "^read data response: unexpected datagram:",
		},
		{
			name:    "no data",
			timeout: time.Millisecond * 10,

			skip: "need to cycle state",

			expectedError: "^reading data.*i/o timeout$",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.skip != "" {
				t.Skip(c.skip)
			}

			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.timeout = c.timeout
			tConn.window = c.window

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			_ = tConn.readData()
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if tConn.err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, tConn.err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, tConn.err.Error())
				}
				return
			}

			// Data
			if string(tConn.rx.data()) != string(c.expectedData) {
				t.Errorf("expected data %q, got %q", string(c.expectedData), string(data))
			}

			// Block number
			if tConn.rx.block() != c.expectedBlock {
				t.Errorf("expected block %d, got %d", c.expectedBlock, tConn.block)
			}

			// Window number
			if tConn.window != c.expectedWindow {
				t.Errorf("expected window %d, got %d", c.expectedWindow, tConn.window)
			}
		})
	}
}

func TestConn_ackData(t *testing.T) {
	tDG := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := []struct {
		name       string
		timeout    time.Duration
		rx         datagram
		block      uint16
		window     uint16
		windowsize uint16
		catchup    bool
		connFunc   func(*net.UDPConn, *net.UDPAddr) error

		expectCatchup  bool
		expectedBlock  uint16
		expectedWindow uint16
		expectedError  string
	}{
		{
			name:       "success, reached window",
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
		{
			name:       "success, reset catchup",
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
		{
			name:       "repeat block",
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
		{
			name:       "future block, no catchup",
			timeout:    time.Second,
			block:      12,
			windowsize: 2,
			window:     1,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(14, data[:512])
				return dg
			}(),
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				conn.SetReadDeadline(time.Now().Add(testConnTimeout))
				_, _, err := conn.ReadFrom(tDG.buf)
				if err != nil {
					t.Errorf("future block, no catchup: expected ACK %v", err)
					return nil
				}

				if tDG.block() != 12 {
					t.Errorf("future block, no catchup: expected ACK with block 12, got %d", tDG.block())
				}
				return nil
			},

			expectCatchup:  true,
			expectedBlock:  12,
			expectedWindow: 0,
			expectedError:  errBlockSequence.Error(),
		},
		{
			name:       "future block, rollover",
			timeout:    time.Second,
			block:      65534,
			windowsize: 4,
			window:     1,
			rx: func() datagram {
				dg := datagram{}
				dg.writeData(0, data[:512])
				return dg
			}(),
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				conn.SetReadDeadline(time.Now().Add(testConnTimeout))
				_, _, err := conn.ReadFrom(tDG.buf)
				if err != nil {
					t.Errorf("future block, no catchup: expected ACK %v", err)
					return nil
				}

				if tDG.block() != 65534 {
					t.Errorf("future block, no catchup: expected ACK with block 65534, got %d", tDG.block())
				}
				return nil
			},

			expectCatchup:  true,
			expectedBlock:  65534,
			expectedWindow: 0,
			expectedError:  errBlockSequence.Error(),
		},
		{
			name:       "future block, catchup",
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
		{
			name:       "past block",
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
		{
			name:       "success, below window",
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.rx = c.rx
			tConn.timeout = c.timeout
			tConn.block = c.block
			tConn.window = c.window
			tConn.windowsize = c.windowsize
			tConn.catchup = c.catchup

			_ = tConn.ackData() // TODO: check return func
			// Error
			if tConn.err != nil {
				err := tConn.err
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}

			if c.connFunc != nil {
				if err := c.connFunc(cNetConn, sAddr); err != nil {
					t.Fatal(err)
				}
			}

			// Block number
			if tConn.block != c.expectedBlock {
				t.Errorf("expected block %d, got %d", c.expectedBlock, tConn.block)
			}

			// Window number
			if tConn.window != c.expectedWindow {
				t.Errorf("expected window %d, got %d", c.expectedWindow, tConn.window)
			}
			// Catchup
			if tConn.catchup != c.expectCatchup {
				t.Errorf("expected catchup %t, but it wasn't", c.expectCatchup)
			}
		})
	}
}

func TestConn_parseOptions(t *testing.T) {
	dg := datagram{}

	cases := []struct {
		name     string
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
		{
			name: "blocksize, valid",
			rx: func() datagram {
				dg.writeOptionAck(options{optBlocksize: "234"})
				return dg
			},

			expectOptionsParsed: true,
			expectedOptions:     options{optBlocksize: "234"},
			expectedBlksize:     234,
			expectedError:       "^$",
		},
		{
			name: "blocksize, invalid",
			rx: func() datagram {
				dg.writeOptionAck(options{optBlocksize: "a"})
				return dg
			},

			expectOptionsParsed: false,
			expectedBlksize:     0,
			expectedError:       `error parsing .* for option "blksize"`,
		},
		{
			name: "timeout, valid",
			rx: func() datagram {
				dg.writeOptionAck(options{optTimeout: "3"})
				return dg
			},

			expectedOptions:     options{optTimeout: "3"},
			expectOptionsParsed: true,
			expectedTimeout:     3 * time.Second,
			expectedError:       `^$`,
		},
		{
			name: "timeout, invalid",
			rx: func() datagram {
				dg.writeOptionAck(options{optTimeout: "three"})
				return dg
			},

			expectOptionsParsed: false,
			expectedTimeout:     0,
			expectedError:       `error parsing .* for option "timeout"`,
		},
		{
			name: "tsize, valid, sending side",
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
		{
			name: "tsize, valid, receive side",
			rx: func() datagram {
				dg.writeOptionAck(options{optTransferSize: "42"})
				return dg
			},

			expectedOptions:     options{},
			expectOptionsParsed: true,
			expectedTsize:       ptrInt64(42),
			expectedError:       `^$`,
		},
		{
			name: "tsize, invalid",
			rx: func() datagram {
				dg.writeOptionAck(options{optTransferSize: "large"})
				return dg
			},

			expectedError: `^error parsing .* for option "tsize"$`,
		},
		{
			name: "windowsize, valid",
			rx: func() datagram {
				dg.writeOptionAck(options{optWindowSize: "32"})
				return dg
			},

			expectedOptions:     options{optWindowSize: "32"},
			expectOptionsParsed: true,
			expectedWindowsize:  32,
			expectedError:       `^$`,
		},
		{
			name: "windowsize, invalid",
			rx: func() datagram {
				dg.writeOptionAck(options{optWindowSize: "x"})
				return dg
			},

			expectedError: `^error parsing .* for option "windowsize"$`,
		},
		{
			name: "all options, sending side",
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
		{
			name: "all options, receive side",
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tConn := conn{rx: c.rx()}
			tConn.tsize = c.tsize
			tConn.isSender = c.isSender

			opts, err := tConn.parseOptions()

			// Error
			if err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}

			// Options
			if !reflect.DeepEqual(c.expectedOptions, opts) {
				t.Errorf("expected options %q, got %q", c.expectedOptions, opts)
			}

			// OptionsParsed
			if c.expectOptionsParsed != tConn.optionsParsed {
				t.Errorf("expected optionsParsed %t, but it wasn't", c.expectOptionsParsed)
			}

			if tConn.blksize != c.expectedBlksize {
				t.Errorf("expected blocksize to be %d, but it was %d", c.expectedBlksize, tConn.blksize)
			}
			if tConn.timeout != c.expectedTimeout {
				t.Errorf("expected timeout to be %s, but it was %s", c.expectedTimeout, tConn.timeout)
			}
			if tConn.windowsize != c.expectedWindowsize {
				t.Errorf("expected window to be %d, but it was %d", c.expectedWindowsize, tConn.windowsize)
			}
			if tConn.tsize != c.expectedTsize {
				if tConn.tsize == nil || c.expectedTsize == nil {
					t.Errorf("expected tsize to be *%d, but it was *%d", c.expectedTsize, tConn.tsize)
				} else if *tConn.tsize != *c.expectedTsize {
					t.Errorf("expected tsize to be %d, but it was %d", *c.expectedTsize, *tConn.tsize)
				}
			}
		})
	}
}

func TestConn_write(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := []struct {
		name          string
		bytes         []byte
		optionsParsed bool
		blksize       uint16
		window        uint16
		windowsize    uint16
		rx            func() datagram
		timeout       time.Duration
		connFunc      func(conn *net.UDPConn, sAddr *net.UDPAddr) error
		connErr       error

		skip bool

		expectedCount  int
		expectedError  string
		expectedWindow uint64
	}{
		{
			name:          "success, buf < blksize",
			timeout:       time.Millisecond,
			bytes:         data[:300],
			blksize:       512,
			optionsParsed: true,

			expectedCount: 300,
			expectedError: "^$",
		},
		{
			name:          "success, buf > blksize, window 1",
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    1,
			optionsParsed: true,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				dg.writeAck(1)
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				if _, err := conn.WriteTo(dg.bytes(), sAddr); err != nil {
					return err
				}

				dg.writeAck(2)
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				return err
			},

			expectedCount: 1024,
			expectedError: "^$",
		},
		{
			name:          "success, buf > blksize, window 2",
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    2,
			optionsParsed: true,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				dg.writeAck(1)
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				return err
			},

			expectedCount: 1024,
			expectedError: "^$",
		},
		{
			name:          "fail to ack",
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    1,
			optionsParsed: true,

			skip: true,

			expectedCount: 1024,
			expectedError: "receiving ACK after writing data: network read failed",
		},
		{
			name:          "conn err",
			timeout:       time.Millisecond * 100,
			bytes:         data[:1024],
			blksize:       512,
			windowsize:    1,
			optionsParsed: true,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				dg.writeAck(1)
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				if err != nil {
					fmt.Printf("Buf: %v\n", dg.buf)
					fmt.Printf("Addr: %v\n", sAddr)
					fmt.Printf("Conn: %#v\n", conn)
				}
				return err
			},
			connErr: errors.New("conn error"),

			expectedCount: 0,
			expectedError: "conn error",
		},
		{
			name:          "writeSetup fails",
			timeout:       time.Millisecond,
			optionsParsed: false,
			rx: func() datagram {
				dg.writeOptionAck(options{optBlocksize: "234"})
				return dg
			},

			skip: true,

			expectedError: "parsing options before write: write setup: network read failed:",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.skip {
				t.Skip()
			}

			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.rx.writeAck(1)
			if c.rx != nil {
				tConn.rx = c.rx()
			}
			tConn.blksize = c.blksize
			tConn.window = c.window
			tConn.windowsize = c.windowsize
			tConn.optionsParsed = false
			tConn.timeout = c.timeout
			tConn.buf = make([]byte, c.blksize)
			tConn.txBuf = newRingBuffer(int(c.windowsize), int(c.blksize))
			tConn.err = c.connErr

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			count, err := tConn.Write(c.bytes)
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}

			// Count
			if c.expectedCount != count {
				t.Errorf("expected count %d, got %d", c.expectedCount, count)
			}
		})
	}
}

func TestConn_Close(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := []struct {
		name     string
		bytes    []byte
		blksize  uint16
		timeout  time.Duration
		connFunc func(conn *net.UDPConn, sAddr *net.UDPAddr) error
		connErr  error

		expectedError string
	}{
		{
			name:    "conn err",
			connErr: errors.New("conn error"),

			expectedError: "checking conn err before Close: conn error",
		},
		{
			name:    "success, no data",
			timeout: time.Millisecond * 100,
			bytes:   []byte{},
			blksize: 512,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				conn.SetReadDeadline(time.Now().Add(testConnTimeout))
				n, _, err := conn.ReadFrom(dg.buf)
				if err != nil {
					return err
				}
				dg.offset = n

				if dg.opcode() != opCodeDATA {
					t.Errorf("expected opcode %s, got %s", opCodeDATA, dg.opcode())
				}

				if l := len(dg.data()); l != 0 {
					t.Errorf("expected data len to be 0, but they were %d", l)
				}

				dg.writeAck(1)
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err = conn.WriteTo(dg.buf, sAddr)
				return err
			},

			expectedError: "^$",
		},
		{
			name:    "success, with data",
			timeout: time.Millisecond * 100,
			bytes:   data[:384],
			blksize: 512,
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				conn.SetReadDeadline(time.Now().Add(testConnTimeout))
				n, _, err := conn.ReadFrom(dg.buf)
				if err != nil {
					return err
				}
				dg.offset = n

				if dg.opcode() != opCodeDATA {
					t.Errorf("expected opcode %s, got %s", opCodeDATA, dg.opcode())
				}

				if l := len(dg.data()); l != 384 {
					t.Errorf("expected data len to be 384, but they were %d", l)
				}

				dg.writeAck(1)
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err = conn.WriteTo(dg.buf, sAddr)
				return err
			},

			expectedError: "^$",
		},
		{
			name:    "timeout",
			timeout: time.Millisecond * 100,
			blksize: 512,

			expectedError: "^reading ack: max retries reached$",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			tConn.blksize = c.blksize
			tConn.timeout = c.timeout
			tConn.buf = make([]byte, c.blksize)
			tConn.txBuf = newRingBuffer(1, int(c.blksize))
			tConn.txBuf.Write(c.bytes)
			tConn.writer = tConn.txBuf
			tConn.err = c.connErr
			tConn.optionsParsed = true

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			err := tConn.Close()
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}
		})
	}
}

func TestConn_read(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	data := getTestData(t, "1MB-random")

	cases := []struct {
		name          string
		bytes         []byte
		blksize       uint16
		windowsize    uint16
		optionsParsed bool
		timeout       time.Duration
		connFunc      func(conn *net.UDPConn, sAddr *net.UDPAddr) error
		connErr       error

		expectedRead  int
		expectedError string
	}{
		{
			name:          "success",
			timeout:       time.Millisecond * 100,
			optionsParsed: true,
			blksize:       512,
			bytes:         make([]byte, 512),
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				dg.writeData(1, data[:512])
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				return err
			},

			expectedRead:  512,
			expectedError: "^$",
		},
		{
			name:          "success, EOF",
			timeout:       time.Millisecond * 100,
			optionsParsed: true,
			blksize:       512,
			bytes:         make([]byte, 512),
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				dg.writeData(1, data[:300])
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				return err
			},

			expectedRead:  300,
			expectedError: "^EOF$",
		},
		{
			name:          "block sequence error",
			timeout:       time.Millisecond * 100,
			optionsParsed: true,
			blksize:       512,
			windowsize:    2,
			bytes:         make([]byte, 512),
			connFunc: func(conn *net.UDPConn, sAddr *net.UDPAddr) error {
				// Write wrong block
				dg.writeData(2, data[:512])
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err := conn.WriteTo(dg.bytes(), sAddr)
				if err != nil {
					return err
				}

				// Receive ACK for previous
				conn.SetReadDeadline(time.Now().Add(testConnTimeout))
				n, _, err := conn.ReadFrom(dg.buf)
				if err != nil {
					return err
				}
				dg.offset = n
				if dg.block() != 0 {
					t.Errorf("expected block 0 again, got %d", dg.block())
				}

				// Write correct data
				dg.writeData(1, data[:300])
				conn.SetWriteDeadline(time.Now().Add(testConnTimeout))
				_, err = conn.WriteTo(dg.bytes(), sAddr)
				return err
			},

			expectedRead:  300,
			expectedError: "^EOF$",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tConn, sAddr, cNetConn, closer := testConns(t)
			defer closer()
			// tConn.optionsParsed = c.optionsParsed
			tConn.blksize = c.blksize
			tConn.timeout = c.timeout
			tConn.windowsize = c.windowsize
			tConn.err = c.connErr
			tConn.rx.writeAck(0)

			errChan := testConnFunc(cNetConn, sAddr, c.connFunc)
			read, err := tConn.Read(c.bytes)
			if err := <-errChan; err != nil {
				t.Fatal(err)
			}

			// Error
			if err != nil {
				if ok, _ := regexp.MatchString(c.expectedError, err.Error()); !ok {
					t.Errorf("expected error %q, got %q", c.expectedError, err.Error())
				}
			}

			// Read Count
			if c.expectedRead != read {
				t.Errorf("expected read bytes to be %d, but it was %d", c.expectedRead, read)
			}
		})
	}
}

func TestConn_sendError(t *testing.T) {
	dg := datagram{buf: make([]byte, 512)}

	cases := []struct {
		name    string
		code    ErrorCode
		msg     string
		blksize uint16
		timeout time.Duration

		expectedCode  ErrorCode
		expectedError string
	}{
		{
			name:    "message, undersize",
			timeout: time.Millisecond * 100,
			blksize: 512,
			code:    ErrCodeNoSuchUser,
			msg:     "foo",

			expectedCode:  ErrCodeNoSuchUser,
			expectedError: "foo",
		},
		{
			name:    "message, oversize",
			timeout: time.Millisecond * 100,
			blksize: 10,
			code:    ErrCodeNoSuchUser,
			msg:     "there was a long error",

			expectedCode:  ErrCodeNoSuchUser,
			expectedError: "there was",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
				t.Errorf("expected errorCode %s, got %s", c.expectedCode, dg.errorCode())
			}

			// Error Message
			if c.expectedError != dg.errMsg() {
				t.Errorf("expected message %q, got %q", c.expectedError, dg.errMsg())
			}
		})
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}

func testConns(t *testing.T) (*conn, *net.UDPAddr, *net.UDPConn, func()) {
	// Statically chose port, letting system assign results in an error on Linux w/ nf_conntrack
	// related to this bug http://marc.info/?l=linux-netdev&s=Possible+race+condition+in+conntracking
	cAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 54321}
	cNetConn, err := net.ListenUDP("udp4", cAddr)
	if err != nil {
		t.Fatal(err)
	}

	sAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 54322}
	sNetConn, err := net.ListenUDP("udp4", sAddr)
	if err != nil {
		t.Fatal(err)
	}

	tConn, err := newConn("udp4", ModeOctet, cAddr)
	if err != nil {
		t.Fatal(err)
	}
	// Replace auto assigned
	tConn.netConn.Close()
	tConn.netConn = sNetConn

	closer := func() {
		cNetConn.Close()
		tConn.netConn.Close()
	}

	return tConn, sAddr, cNetConn, closer
}
