// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"io"
	"strings"
)

func (l *NewerLineReader) updateline(prev string, w io.Writer) error {
	if prev != l.Line {
		if _, err := w.Write([]byte(strings.Repeat("\b", len(l.Prompt+prev)) + l.Prompt + l.Line)); err != nil {
			return err
		}
	}
	return nil
}
