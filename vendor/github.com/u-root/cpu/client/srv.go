// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/u-root/pkg/ulog"
)

// Made harder as you can't set a read deadline on ssh.Conn
func (c *Cmd) srv(l net.Listener) error {
	// We only accept once
	defer l.Close()
	var (
		errs = make(chan error)
		s    net.Conn
		err  error
	)
	go func() {
		verbose("srv: try to accept l %v", l)
		s, err = l.Accept()
		verbose("Accept: %v %v", s, err)
		if err != nil {
			errs <- fmt.Errorf("accept 9p socket: %v", err)
			return
		}
		verbose("srv got %v", s)
		var rn nonce
		if _, err := io.ReadAtLeast(s, rn[:], len(rn)); err != nil {
			errs <- fmt.Errorf("Reading nonce from remote: %v", err)
			return
		}
		verbose("srv: read the nonce back got %s", rn)
		if c.nonce.String() != rn.String() {
			errs <- fmt.Errorf("nonce mismatch: got %s but want %s", rn, c.nonce)
			return
		}
		errs <- nil
	}()

	// We block here on the mount. If the user wanted the 9p mount, and it never
	// occurs, we don't want to continue; files they may want to use might be
	// aliased by local files. Similarly, on the cpud side, if the mount
	// has an error, cpud will exit now. It used to soldier on, but we've
	// realized that's a very bad idea; now that we have the -9p switch,
	// we can now do a cpu session without the 9p server. The timeout
	// is no longer important, since not all cpu sessions need 9p.
	if err := <-errs; err != nil {
		return fmt.Errorf("srv: %v", err)
	}
	// If we are debugging, add the option to trace records.
	verbose("Start serving on %v", c.Root)
	var opts []p9.ServerOpt
	if Debug9p {
		if Dump9p {
			log.SetOutput(DumpWriter)
			log.SetFlags(log.Ltime | log.Lmicroseconds)
			ulog.Log = log.New(DumpWriter, "9p", log.Ltime|log.Lmicroseconds)
		}
		opts = append(opts, p9.WithServerLogger(ulog.Log))
	}

	if err := p9.NewServer(c.fileServer, opts...).Handle(s, s); err != nil {
		if err != io.EOF {
			log.Printf("Serving cpu remote: %v", err)
			return err
		}
	}
	return nil
}
