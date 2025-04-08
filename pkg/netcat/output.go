// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/u-root/u-root/pkg/ulog"
)

type OutputOptions struct {
	OutFilePath     string      // Dump session data to a file
	OutFileMutex    sync.Mutex  // Mutex for the file
	OutFileHexPath  string      // Dump session data in hex to a file
	OutFileHexMutex sync.Mutex  // Mutex for the hex file
	AppendOutput    bool        // Append the resulted output rather than truncating
	Logger          ulog.Logger // Verbose output
}

// Write writes the data to the file specified in the options
// If Netcat is not configured to write to a file, it will return 0, nil
// https://go.dev/src/io/io.go
func (n *OutputOptions) Write(data []byte) (int, error) {
	if n.OutFilePath == "" && n.OutFileHexPath == "" {
		return 0, nil
	}

	fileOpts := os.O_CREATE | os.O_WRONLY
	if n.AppendOutput {
		fileOpts |= os.O_APPEND
	} else {
		fileOpts |= os.O_TRUNC
	}

	if n.OutFilePath != "" {
		n.OutFileMutex.Lock()
		defer n.OutFileMutex.Unlock()

		f, err := os.OpenFile(n.OutFilePath, fileOpts, 0o644)
		if err != nil {
			return 0, fmt.Errorf("netcat file open: %w", err)
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			return 0, fmt.Errorf("netcat file write: %w", err)
		}
	}

	if n.OutFileHexPath != "" {
		n.OutFileHexMutex.Lock()
		defer n.OutFileHexMutex.Unlock()

		f, err := os.OpenFile(n.OutFileHexPath, fileOpts, 0o644)
		if err != nil {
			return 0, fmt.Errorf("netcat hex file open: %w", err)
		}
		defer f.Close()

		_, err = f.Write([]byte(hex.Dump(data)))
		if err != nil {
			return 0, fmt.Errorf("netcat hex file write: %w", err)
		}

	}

	return len(data), nil
}

// Close does nothing: Write already closes any underlying files. Close exists
// only to implement the io.WriteCloser interface with OutputOptions.
func (n *OutputOptions) Close() error {
	return nil
}
