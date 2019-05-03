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
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// FSID is the File System ID.
type FSID string

var (
	// SrvDebug is used to print server debug messages.
	SrvDebug = func(string, ...interface{}) {}
	// servers maps FSID to file systems.
	servers = map[FSID]*FS{}
	smu     sync.Mutex
)

// PingReq is a ping request for an FSID at a server.
type PingReq struct {
	R FSID
}

// PingResp is the response to a PingReq.
type PingResp struct {
	R   string
	Err error
}

// Ping responds to a Ping Request.
func (id FSID) Ping(req *PingReq, resp *PingResp) error {
	fs, ok := servers[id]
	if !ok {
		resp.Err = fmt.Errorf("No such fs %v", id)
		return nil
	}
	resp.R = fmt.Sprintf("%v", fs)
	resp.Err = nil
	return nil
}

func addFS(id FSID, fs *FS) error {
	smu.Lock()
	defer smu.Unlock()
	if _, ok := servers[id]; ok {
		return fmt.Errorf("Server %q already exists", id)
	}
	servers[id] = fs
	return nil
}

func getFS(id FSID) (*FS, error) {
	smu.Lock()
	defer smu.Unlock()
	fs, ok := servers[id]
	if !ok {
		return nil, fmt.Errorf("Server %q does not exist", id)
	}
	return fs, nil
}

// Serve is a server for netfuse Client requests. Its only use of the
// FUSE package is for structure definitions.
func Serve(root string, c net.Conn) error {
	fs := NewFS(root)
	id := FSID(root)
	if err := addFS(id, fs); err != nil {
		return err
	}

	s := rpc.NewServer()
	if err := s.Register(&id); err != nil {
		log.Printf("register failed: %v", err)
		return err
	}
	s.ServeConn(c)
	return nil
}
