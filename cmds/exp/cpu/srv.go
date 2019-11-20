// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/hugelgupf/p9/p9"
)

func srv(l net.Listener, root string, n nonce, deadline time.Time) {
	// We only accept once
	defer l.Close()
	v("srv: try to accept")
	c, err := l.Accept()
	if err != nil {
		log.Fatalf("accept 9p socket: %v", err)
	}
	v("srv got %v", c)
	if err := c.SetDeadline(deadline); err != nil {
		log.Fatalf("Set deadline for nonce: %v", err)
	}
	var rn nonce
	if _, err := io.ReadAtLeast(c, rn[:], len(rn)); err != nil {
		log.Fatalf("Reading nonce from remote: %v", err)
	}
	v("srv: read the nonce back got %s", rn)
	if n.String() != rn.String() {
		log.Fatalf("nonce mismatch: got %s but want %s", rn, n)
	}

	// There is no deadline once we are set up.
	if err := c.SetDeadline(time.Time{}); err != nil {
		log.Printf("Warning: clearing deadline for socket: %v, things may fail", err)
	}
	if err := p9.NewServer(&cpu9p{path: root}).Handle(c); err != nil {
		log.Fatalf("Serving cpu remote: %v", err)
	}
}
