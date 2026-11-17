// Copyright 2012-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// tftpd is a basic TFTP server that serves files from the current directory
//
// Synopsis:
//
//	tftpd [-port PORT] [-root DIRECTORY]
//
// Description:
//
//	tftpd runs a simple TFTP server that can handle both read and write requests.
//	It uses the github.com/pin/tftp library.
//
// Options:
//
//	-port: Port to listen on (default 69)
//	-root: Root directory to serve files from (default current directory)

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pin/tftp/v3"
)

var (
	port    = flag.Int("port", 69, "Port to listen on")
	rootDir = flag.String("root", ".", "Root directory to serve files from")
	verbose = flag.Bool("v", false, "Enable verbose logging")
)

var (
	errInvalidPath     = errors.New("invalid path")
	errAccessViolation = errors.New("access violation")
	errFileNotFound    = errors.New("file not found")
)

// isValidPath checks if the provided path is valid and secure
// using the approach recommended in https://go.dev/blog/osroot
func isValidPath(rootDir, path string) (string, error) {
	// Clean the path to remove any ".." elements
	cleanPath := filepath.Clean(path)

	// Create an absolute path by joining with root
	targetPath := filepath.Join(rootDir, cleanPath)

	// Resolve any symlinks
	realPath, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			// A dangling symlink also ends up here. Creating the
			// file would follow it to wherever it points.
			if _, err := os.Lstat(targetPath); err == nil {
				return "", errInvalidPath
			}

			// For non-existent paths (e.g., for file creation),
			// we need to check the parent directory
			parentDir := filepath.Dir(targetPath)
			realParent, err := filepath.EvalSymlinks(parentDir)
			if err != nil {
				return "", err
			}

			// Check if the parent is still under root
			rootAbs, err := filepath.EvalSymlinks(rootDir)
			if err != nil {
				return "", err
			}

			// Make paths canonical for comparison
			realParent = filepath.Clean(realParent)
			rootAbs = filepath.Clean(rootAbs)

			// Check if parent directory is under root
			rel, err := filepath.Rel(rootAbs, realParent)
			if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
				return "", errInvalidPath
			}

			return targetPath, nil
		}
		return "", err
	}

	// Get the canonical path for the root dir
	rootAbs, err := filepath.EvalSymlinks(rootDir)
	if err != nil {
		return "", err
	}

	// Make paths canonical for comparison
	realPath = filepath.Clean(realPath)
	rootAbs = filepath.Clean(rootAbs)

	// Check if the path is contained within the root
	rel, err := filepath.Rel(rootAbs, realPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		return "", errInvalidPath
	}

	return realPath, nil
}

// FileReadHandler serves TFTP read requests from a directory
type FileReadHandler struct {
	Root string
}

// ServeTFTP handles TFTP read requests by serving files from the root directory
func (h *FileReadHandler) ServeTFTP(requestedPath string, rf io.ReaderFrom) error {
	if *verbose {
		log.Printf("Received read request for %s", requestedPath)
	}

	// Validate the path securely
	realPath, err := isValidPath(h.Root, requestedPath)
	if err != nil {
		if errors.Is(err, errInvalidPath) {
			return errAccessViolation
		} else if os.IsNotExist(err) {
			return errFileNotFound
		}
		return errors.New("cannot access file")
	}

	file, err := os.Open(realPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errFileNotFound
		} else if os.IsPermission(err) {
			return errors.New("permission denied")
		}
		return errors.New("cannot open file")
	}
	defer file.Close()

	_, err = rf.ReadFrom(file)
	if err != nil {
		log.Printf("Error sending file: %v", err)
	}
	return err
}

// FileWriteHandler stores files received by TFTP write requests in a directory
type FileWriteHandler struct {
	Root string
}

// ReceiveTFTP handles TFTP write requests by storing files in the root directory
func (h *FileWriteHandler) ReceiveTFTP(requestedPath string, wt io.WriterTo) error {
	if *verbose {
		log.Printf("Received write request for %s", requestedPath)
	}

	// Validate the path securely
	realPath, err := isValidPath(h.Root, requestedPath)
	if err != nil {
		if errors.Is(err, errInvalidPath) {
			return errAccessViolation
		}
		return errors.New("cannot access path")
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(realPath), 0o755)
	if err != nil {
		return errors.New("cannot create directory")
	}

	file, err := os.Create(realPath)
	if err != nil {
		if os.IsPermission(err) {
			return errors.New("permission denied")
		}
		return errors.New("cannot create file")
	}
	defer file.Close()

	_, err = wt.WriteTo(file)
	if err != nil {
		log.Printf("Error receiving file: %v", err)
	}
	return err
}

func main() {
	flag.Parse()

	// Resolve the absolute path of root directory
	absRoot, err := filepath.Abs(*rootDir)
	if err != nil {
		log.Fatalf("Failed to resolve root directory path: %v", err)
	}

	// Create the server
	server := tftp.NewServer(
		(&FileReadHandler{Root: absRoot}).ServeTFTP,
		(&FileWriteHandler{Root: absRoot}).ReceiveTFTP,
	)

	log.Printf("TFTP server listening on port %d, serving files from %s", *port, absRoot)
	err = server.ListenAndServe(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
