// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

// ReadHandler responds to a TFTP read request.
type ReadHandler interface {
	ServeTFTP(ReadRequest)
}

// WriteHandler responds to a TFTP write request.
type WriteHandler interface {
	ReceiveTFTP(WriteRequest)
}

// ReadWriteHandler combines ReadHandler and WriteHandler.
type ReadWriteHandler interface {
	ReadHandler
	WriteHandler
}

// WriteRequest is provided to a WriteHandler's ReceiveTFTP method.
type WriteRequest interface {
	// Addr is the network address of the client.
	Addr() *net.UDPAddr

	// Name is the file name provided by the client.
	Name() string

	// Read reads the request data from the client.
	Read([]byte) (int, error)

	// Size returns the transfer size (tsize) as provided by the client.
	// If the tsize option was not negotiated, an error will be returned.
	Size() (int64, error)

	// WriteError sends an error to the client and terminates the
	// connection. WriteError can only be called once. Read cannot
	// be called after an error has been written.
	WriteError(ErrorCode, string)

	// TransferMode returns the TFTP transfer mode requested by the client.
	TransferMode() TransferMode
}

// writeRequest implements WriteRequest.
type writeRequest struct {
	conn *conn

	name string
}

func (w *writeRequest) Addr() *net.UDPAddr {
	return w.conn.remoteAddr.(*net.UDPAddr)
}

func (w *writeRequest) Name() string {
	return w.name
}

func (w *writeRequest) Read(p []byte) (int, error) {
	return w.conn.Read(p)
}

func (w *writeRequest) Size() (int64, error) {
	if w.conn.tsize == nil {
		return 0, ErrSizeNotReceived
	}
	return *w.conn.tsize, nil
}

func (w *writeRequest) WriteError(c ErrorCode, s string) {
	w.conn.sendError(c, s)
}

func (w *writeRequest) TransferMode() TransferMode {
	return w.conn.mode
}

// ReadRequest is provided to a ReadHandler's ServeTFTP method.
type ReadRequest interface {
	// Addr is the network address of the client.
	Addr() *net.UDPAddr

	// Name is the file name requested by the client.
	Name() string

	// Write write's data to the client.
	Write([]byte) (int, error)

	// WriteError sends an error to the client and terminates the
	// connection. WriteError can only be called once. Write cannot
	// be called after an error has been written.
	WriteError(ErrorCode, string)

	// WriteSize sets the transfer size (tsize) value to be sent to
	// the client. It must be called before any calls to Write.
	WriteSize(int64)

	// TransferMode returns the TFTP transfer mode requested by the client.
	TransferMode() TransferMode
}

// readRequest implements ReadRequest.
type readRequest struct {
	conn *conn

	name string
}

func (w *readRequest) Addr() *net.UDPAddr {
	return w.conn.remoteAddr.(*net.UDPAddr)
}

func (w *readRequest) Name() string {
	return w.name
}

func (w *readRequest) Write(p []byte) (int, error) {
	return w.conn.Write(p)
}

func (w *readRequest) WriteError(c ErrorCode, s string) {
	w.conn.sendError(c, s)
}

func (w *readRequest) WriteSize(i int64) {
	w.conn.tsize = &i
}

func (w *readRequest) TransferMode() TransferMode {
	return w.conn.mode
}

// FileServer creates a handler for sending and reciving files on the filesystem.
func FileServer(dir string) ReadWriteHandler {
	return &fileServer{path: dir, log: newLogger("fileserver")}
}

type fileServer struct {
	log  *logger
	path string
}

// ServeTFTP serves files rooted at the configured directory.
//
// If the file does not exist or otherwise cannot be opened, a File Not Found
// error will be sent.
func (f *fileServer) ServeTFTP(w ReadRequest) {
	path := filepath.Join(f.path, filepath.Clean(w.Name()))

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		w.WriteError(ErrCodeFileNotFound, fmt.Sprintf("File %q does not exist", w.Name()))
		return
	}
	defer errorDefer(file.Close, f.log, "error closing file")

	finfo, _ := file.Stat()
	w.WriteSize(finfo.Size())
	if _, err = io.Copy(w, file); err != nil {
		log.Println(err)
	}
}

// ReceiveTFTP writes received files to the configured directory.
//
// If the file cannot be created an Access Violation error will be sent.
func (f *fileServer) ReceiveTFTP(r WriteRequest) {
	path := filepath.Join(f.path, filepath.Clean(r.Name()))

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		r.WriteError(ErrCodeAccessViolation, fmt.Sprintf("Cannot create file %q", filepath.Clean(r.Name())))
	}
	defer errorDefer(file.Close, f.log, "error closing file")

	_, err = io.Copy(file, r)
	if err != nil {
		log.Println(err)
	}
}

// ReadHandlerFunc is an adapter type to allow a function to serve as a ReadHandler.
type ReadHandlerFunc func(ReadRequest)

// ServeTFTP calls the ReadHandlerFunc function.
func (h ReadHandlerFunc) ServeTFTP(w ReadRequest) {
	h(w)
}

// WriteHandlerFunc is an adapter type to allow a function to serve as a WriteHandler.
type WriteHandlerFunc func(WriteRequest)

// ReceiveTFTP calls the WriteHandlerFunc function.
func (h WriteHandlerFunc) ReceiveTFTP(w WriteRequest) {
	h(w)
}
