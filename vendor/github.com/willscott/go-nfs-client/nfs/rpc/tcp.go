// Copyright © 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause
//
package rpc

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type tcpTransport struct {
	r       io.Reader
	wc      net.Conn
	timeout time.Duration
}

// Get the response from the conn, buffer the contents, and return a reader to
// it.
func (t *tcpTransport) recv() (io.ReadSeeker, error) {
	var hdr uint32
	if err := binary.Read(t.r, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}

	buf := make([]byte, hdr&0x7fffffff)
	if _, err := io.ReadFull(t.r, buf); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func (t *tcpTransport) Write(buf []byte) (int, error) {
	var hdr uint32 = uint32(len(buf)) | 0x80000000
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, hdr)
	if t.timeout != 0 {
		deadline := time.Now().Add(t.timeout)
		t.wc.SetWriteDeadline(deadline)
	}
	n, err := t.wc.Write(append(b, buf...))

	return n, err
}

func (t *tcpTransport) Close() error {
	return t.wc.Close()
}

func (t *tcpTransport) SetTimeout(d time.Duration) {
	t.timeout = d
	if d == 0 {
		var zeroTime time.Time
		t.wc.SetDeadline(zeroTime)
	}
}
