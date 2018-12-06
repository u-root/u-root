// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import "bytes"

// wc embeds a bytes.Buffer, and we add a Close function
// to make qemu package happy.
// It can be used for SerialOutput.
type wc struct {
	bytes.Buffer
}

// Close implements close on the wc.
// There's not much to do, if we cared, we could
// just put a chan in here that someone could
// wait on to confirm it was really closed.
func (w *wc) Close() error {
	return nil
}
