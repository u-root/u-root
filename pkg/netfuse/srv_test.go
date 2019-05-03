// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netfuse

import (
	"net"
	"net/rpc"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	ln, err := net.Listen("tcp", ":")
	if err != nil {
		t.Fatalf("Listen: got %v, want nil", err)
	}
	id := FSID("root")
	go func() {
		a := ln.Addr().String()
		// If you get lost, enabled this. t.Logf("Now dial %v", a)
		cl, err := rpc.Dial("tcp", a)
		if err != nil {
			t.Fatalf("Dial: got %v, want nil", err)
		}

		defer cl.Close()
		t.Logf("c is %v", cl)

		var r PingResp
		if err := cl.Call("FSID.Ping", &PingReq{id}, &r); err != nil {
			t.Fatalf("Ping: got %v, want nil", err)
		}
		t.Logf("%v(%v): %v\n", "FSID.Ping", nil, r.R)
		if err := cl.Call("FSID.PPing", &PingReq{id}, &r); err == nil {
			t.Fatalf("Ping: got nil, want err")
		}

	}()

	t.Logf("Listening on %v at %v", ln.Addr(), time.Now())
	s, err := ln.Accept()
	if err != nil {
		t.Fatalf("Listen failed: %v at %v", err, time.Now())
	}
	t.Logf("Accepted %v", s)

	if err := Serve("/tmp", s); err != nil {
		t.Fatalf("TestNewServer: Got %v, want nil", err)
	}

}
