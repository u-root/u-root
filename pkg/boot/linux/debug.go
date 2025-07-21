// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

// Debug is called to print out verbose debug info.
//
// Set this to appropriate output stream for display
// of useful debug info.
var Debug = func(string, ...interface{}) {}
