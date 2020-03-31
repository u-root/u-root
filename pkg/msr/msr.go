// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msr

type MSRVal struct {
	Addr  MSR
	Name  string
	Clear uint64
	Set   uint64
}

func (m MSRVal) String() string {
	return m.Name
}

var Debug = func(string, ...interface{}) {}
