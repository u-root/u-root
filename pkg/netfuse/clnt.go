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
	"net/rpc"

	"github.com/u-root/fuse"
	"github.com/u-root/fuse/fuseutil"
)

// Clnt is an RPC client which transmits FUSE requests to an RPC server.
type Clnt struct {
	fuseutil.NotImplementedFileSystem
	*rpc.Client
}

// ClntDebug is used to print debug messages for clients.
var ClntDebug = func(string, ...interface{}) {}

// NewClntFS accepts an rpc.Client and creates a FUSE server which can be
// used to transmit FUSE requests to a server.
func NewClntFS(c *rpc.Client) fuse.Server {
	fs := &Clnt{
		Client: c,
	}
	return fuseutil.NewFileSystemServer(fs)
}

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func (fs *Clnt) checkInvariants() {
}
