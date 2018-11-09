// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}
