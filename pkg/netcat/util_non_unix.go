// Copyright 2024-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9 || windows || tamago

package netcat

// Close implements io.WriteCloser.Close.
func (swc *StdoutWriteCloser) Close() error {
	return nil
}
