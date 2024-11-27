// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	func() {
		f, err := os.OpenFile(
			"some-random-acyclic-graph.txt",
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0o666)
		if err != nil {
			panic(err)
		}

		rnd := rand.New(rand.NewSource(1))
		n := 10_000
		for range 100 * n {
			x := rnd.Intn(n + 1)
			y := rnd.Intn(n + 1)
			_, _ = fmt.Fprintln(f, min(x, y), max(x, y))
		}
	}()

	func() {
		f, err := os.OpenFile(
			"some-random-cyclic-graph.txt",
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0o666)
		if err != nil {
			panic(err)
		}

		rnd := rand.New(rand.NewSource(1))
		n := 200
		for range 100 * n {
			x := rnd.Intn(n + 1)
			y := rnd.Intn(n + 1)
			_, _ = fmt.Fprintln(f, x, y)
		}
	}()
}
