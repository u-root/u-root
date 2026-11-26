// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

// NewWinsize builds a window size from the dimensions a terminal request
// carries. It exists so that callers can build one without knowing whether the
// platform's Winsize embeds unix.Winsize or spells the fields out.
func NewWinsize(row, col, xpixel, ypixel uint16) *Winsize {
	return &Winsize{Row: row, Col: col, Xpixel: xpixel, Ypixel: ypixel}
}
