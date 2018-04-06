// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

// Completer is an interface for completion functions.
// It is passed a string and returns a []string with all
// matches and an error.
type Completer interface {
	Complete(s string) ([]string, error)
}
