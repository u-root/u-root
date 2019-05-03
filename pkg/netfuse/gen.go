// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"
)

type RPC struct {
	N  string
	OT string
}

var ops = []RPC{
	{N: "ReleaseFileHandle"},
	{N: "MkDir"},
	{N: "StatFS"},
	{N: "SetInodeAttributes"},
	{N: "MkNode"},
	{N: "CreateFile"},
	{N: "CreateSymlink"},
	{N: "CreateLink"},
	{N: "Rename"},
	{N: "RmDir"},
	{N: "Unlink"},
	{N: "OpenFile"},
	{N: "ReadSymlink"},
	{N: "OpenDir"},
	{N: "ForgetInode"},
	{N: "ReleaseDirHandle"},
	{N: "GetInodeAttributes"},
	{N: "LookUpInode"},
}

// ops that folllow a pattern of sending a []byte and getting
// a length back
var writeOps = []RPC{
	{N: "WriteFile"},
}

// ops that folllow a pattern of sending an int and a []byte back
var readOps = []RPC{
	{N: "ReadFile", OT: "int64"},
	{N: "ReadDir", OT: "fuseops.DirOffset"},
}

var commonOp = `// Resp{{.N}} is used to transmit {{.N}} responses to the client.
type Resp{{.N}} struct {
   fuseops.{{.N}}Op
   // Err is the error represented as a string. Why a string?
   // https://groups.google.com/forum/#!topic/golang-dev/Cua1Av1J8Nc
   // Errors must be sent as strings.

   Err string
}

// {{.N}} implements {{.N}} on the RPC server side.
func (id FSID) {{.N}}(req *fuseops.{{.N}}Op, resp *Resp{{.N}}) error {
	fs, ok := servers[id]
	if !ok {
		return fmt.Errorf("no server for %v", fs)
	}
	resp.Err = ErrToString(fs.{{.N}}(context.TODO(), req))
        resp.{{.N}}Op = *req
	return nil
}

// {{.N}} implements {{.N}} on the RPC client side.
// It processes FUSE requests to forward to the server.
func (fs *Clnt) {{.N}}(ctx context.Context, op *fuseops.{{.N}}Op) (error) {
        var resp = &Resp{{.N}}{}
        err := fs.Call("FSID.{{.N}}", op, resp)
        if err != nil {
              return err
        }
        *op = resp.{{.N}}Op
        return StringToErr(resp.Err)
}

`

// the writeOp sends a slice and expects int in return.
var writeOp = `// Resp{{.N}} is used to transmit {{.N}} responses to the client.
type Resp{{.N}} struct {
   // BytesWritten is the number of bytes written
   BytesWritten int

   // Err is the error represented as a string. Why a string?
   // https://groups.google.com/forum/#!topic/golang-dev/Cua1Av1J8Nc
   // Errors must be sent as strings.
   Err string
}

// {{.N}} implements {{.N}} on the RPC server side.
func (id FSID) {{.N}}(req *fuseops.{{.N}}Op, resp *Resp{{.N}}) error {
	fs, ok := servers[id]
	if !ok {
		return fmt.Errorf("no server for %v", fs)
	}
	resp.Err = ErrToString(fs.{{.N}}(context.TODO(), req))
        resp.BytesWritten = len(req.Data)
	return nil
}

// {{.N}} implements {{.N}} on the RPC client side.
// It processes FUSE requests to forward to the server.
func (fs *Clnt) {{.N}}(ctx context.Context, op *fuseops.{{.N}}Op) (error) {
        var resp = &Resp{{.N}}{}
        err := fs.Call("FSID.{{.N}}", op, resp)
        if err != nil {
              return err
        }
	if resp.BytesWritten < 0 {
		op.Data = []byte{}
	} else {
		op.Data = op.Data[:resp.BytesWritten]
	}

        return StringToErr(resp.Err)
}

`

// the readOp sends an int and gets a [] in return
var readOp = `// Req{{.N}} is used to transmit {{.N}} resquests to the server.
type Req{{.N}} struct {
	// The file inode that we are reading, and the handle previously returned by
	// CreateFile or OpenFile when opening that inode.
	Inode  fuseops.InodeID
	Handle fuseops.HandleID

	// The offset within the file at which to read.
	Offset int64
        BytesRead int
}

// Resp{{.N}} is used to transmit {{.N}} responses to the client.
type Resp{{.N}} struct {
   // Dst is the data from the read.
   Dst []byte

   // Err is the error represented as a string. Why a string?
   // https://groups.google.com/forum/#!topic/golang-dev/Cua1Av1J8Nc
   // Errors must be sent as strings.
   Err string
}

// {{.N}} implements {{.N}} on the RPC server side.
func (id FSID) {{.N}}(req *Req{{.N}}, resp *Resp{{.N}}) error {
	fs, ok := servers[id]
	if !ok {
		return fmt.Errorf("no server for %v", fs)
	}
        var r = &fuseops.{{.N}}Op{
               Inode: req.Inode,
               Handle: req.Handle,
               Offset: {{.OT}}(req.Offset),
               Dst: make([]byte, req.BytesRead),
        }

	resp.Err = ErrToString(fs.{{.N}}(context.TODO(), r))
        if r.BytesRead > 0 {
           resp.Dst = r.Dst[:r.BytesRead]
        }
	return nil
}

// {{.N}} implements {{.N}} on the RPC client side.
// It processes FUSE requests to forward to the server.
func (fs *Clnt) {{.N}}(ctx context.Context, op *fuseops.{{.N}}Op) (error) {
        var resp = &Resp{{.N}}{}
        var req = &Req{{.N}}{
               Inode: op.Inode,
               Handle: op.Handle,
               Offset: int64(op.Offset),
               BytesRead: len(op.Dst),
        }
        err := fs.Call("FSID.{{.N}}", req, resp)
        if err != nil {
              return err
        }
        copy(op.Dst, resp.Dst)
        op.BytesRead = len(resp.Dst)
        return StringToErr(resp.Err)
}

`

func main() {
	c, err := template.New("crpc").Parse(commonOp)
	if err != nil {
		log.Fatalf("new failed: %v", err)
	}
	w, err := template.New("wrpc").Parse(writeOp)
	if err != nil {
		log.Fatalf("new failed: %v", err)
	}
	r, err := template.New("rrpc").Parse(readOp)
	if err != nil {
		log.Fatalf("new failed: %v", err)
	}
	var b = bytes.NewBufferString(`// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is generated by gen.go. Only edit it for debugging purposes.

package netfuse

import (
	"context"
	"fmt"
	"github.com/u-root/fuse/fuseops"
)

// https://groups.google.com/forum/#!topic/golang-dev/Cua1Av1J8Nc
// Errors must be sent as strings.


`)

	for _, i := range ops {
		if err := c.Execute(b, i); err != nil {
			log.Fatalf("execution failed: %s", err)
		}
	}
	for _, i := range writeOps {
		if err := w.Execute(b, i); err != nil {
			log.Fatalf("execution failed: %s", err)
		}
	}
	for _, i := range readOps {
		if err := r.Execute(b, i); err != nil {
			log.Fatalf("execution failed: %s", err)
		}
	}
	if err := ioutil.WriteFile("methods.go", b.Bytes(), 0600); err != nil {
		log.Fatalf("%v", err)
	}

}
