// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uiotest

import (
	"io"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uio"
)

// TestLineWriter is an io.Writer that logs full lines of serial to tb.
func TestLineWriter(tb testing.TB, prefix string) io.WriteCloser {
	if len(prefix) > 0 {
		return uio.FullLineWriter(&testLinePrefixWriter{tb: tb, prefix: prefix})
	}
	return uio.FullLineWriter(&testLineWriter{tb: tb})
}

// testLinePrefixWriter is an io.Writer that logs full lines of serial to tb.
type testLinePrefixWriter struct {
	tb     testing.TB
	prefix string
}

func (tsw *testLinePrefixWriter) OneLine(p []byte) {
	tsw.tb.Logf("%s %s: %s", testutil.NowLog(), tsw.prefix, p)
}

// testLineWriter is an io.Writer that logs full lines of serial to tb.
type testLineWriter struct {
	tb testing.TB
}

func (tsw *testLineWriter) OneLine(p []byte) {
	tsw.tb.Logf("%s: %s", testutil.NowLog(), p)
}
