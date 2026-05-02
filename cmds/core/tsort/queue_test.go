// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := queue{}
	if next, ok := q.dequeue(); ok {
		t.Fatalf(`queue %v: want no element to be dequeued, got %q`, q, next)
	}

	q.enqueue("a")
	q.enqueue("b")
	q.enqueue("c")

	next, ok := q.dequeue()
	if !ok {
		t.Fatalf(`queue %v: want non-empty queue, got empty queue`, q)
	}
	if next != "a" {
		t.Fatalf(`queue %v: want dequeued element to be "a", got %q`, q, next)
	}
	if next, _ := q.dequeue(); next != "b" {
		t.Fatalf(`queue %v: want dequeued element to be "b", got %q`, q, next)
	}
	if next, _ := q.dequeue(); next != "c" {
		t.Fatalf(`queue %v: want dequeued element to be "c", got %q`, q, next)
	}

	if next, ok := q.dequeue(); ok {
		t.Fatalf(`queue %v: want no element to be dequeued, got %q`, q, next)
	}
}
