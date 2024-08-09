// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux

package ulog

// KernelLog prints to stderr log on non-Linux systems.
var KernelLog = Log
