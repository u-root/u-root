// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
	func() {
		f, err := os.OpenFile(
			"some-random-acyclic-graph.txt",
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0o600)
		if err != nil {
			panic(err)
		}

		rnd := rand.New(rand.NewPCG(1, 1))
		// Produces an acyclic graph through a fixed RNG seed, a carefully
		// chosen node range and sorted edge endpoints.
		nodeRange := 10_000
		for range 100 * nodeRange {
			x := rnd.IntN(nodeRange + 1)
			y := rnd.IntN(nodeRange + 1)
			_, _ = fmt.Fprintln(f, min(x, y), max(x, y))
		}
	}()

	func() {
		f, err := os.OpenFile(
			"some-random-cyclic-graph.txt",
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0o600)
		if err != nil {
			panic(err)
		}

		rnd := rand.New(rand.NewPCG(1, 1))
		// Produces a cyclic graph through a fixed RNG seed and sheer probability.
		nodeRange := uint(200)
		for range 100 * nodeRange {
			x := rnd.UintN(nodeRange + 1)
			y := rnd.UintN(nodeRange + 1)
			_, _ = fmt.Fprintln(f, x, y)
		}
	}()
}
