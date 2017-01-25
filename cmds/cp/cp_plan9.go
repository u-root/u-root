// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func sameFile(sys1, sys2 interface{}) bool {
	a := sys1.(*dir)
	b := sys2.(*dir)
	return a.Qid.Path == b.Qid.Path && a.Type == b.Type && a.Dev == b.Dev
}
