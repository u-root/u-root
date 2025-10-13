// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type set map[string]struct{}

func makeSet() set {
	return make(set)
}

func (s set) add(value string) {
	s[value] = struct{}{}
}

func (s set) has(value string) bool {
	_, ok := s[value]
	return ok
}

func (s set) remove(value string) {
	delete(s, value)
}
