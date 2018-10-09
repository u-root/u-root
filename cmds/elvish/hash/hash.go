// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package hash contains some common hash functions suitable for use in hash
// maps.
package hash

import "fmt"

const DJBInit uint32 = 5381

func DJBCombine(acc, h uint32) uint32 {
	return (acc<<5 + acc) + h
}

func DJB(hs ...uint32) uint32 {
	acc := DJBInit
	for _, h := range hs {
		acc = DJBCombine(acc, h)
	}
	return acc
}

func Hash(i interface{}) uint32 {
	s := fmt.Sprintf("%x", i)
	h := DJBInit
	for i := 0; i < len(s); i++ {
		h = DJBCombine(h, uint32(s[i]))
	}
	return h
}
