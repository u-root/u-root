// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	q := queue{}
	if !q.isEmpty() {
		t.Fatalf("queue %v: want empty, got non-empty", q)
	}

	q.enqueue("a")
	if q.isEmpty() {
		t.Fatalf(`queue %v: want non-empty, got empty`, q)
	}

	q.enqueue("b")
	q.enqueue("c")

	next := q.dequeue()
	if next != "a" {
		t.Fatalf(`queue %v: want dequeued element to be "a", got %q`, q, next)
	}
	next = q.dequeue()
	if next != "b" {
		t.Fatalf(`queue %v: want dequeued element to be "b", got %q`, q, next)
	}
	next = q.dequeue()
	if next != "c" {
		t.Fatalf(`queue %v: want dequeued element to be "c", got %q`, q, next)
	}

	if !q.isEmpty() {
		t.Errorf("queue %v: want empty, got non-empty", q)
	}
	caughtPanic := catchPanic(func() { q.dequeue() })
	if caughtPanic == nil {
		t.Fatalf(`queue %v: want dequeue to panic, got no panic`, q)
	}
	if caughtPanic.Error() != "queue is empty" {
		t.Fatalf(
			`queue %v: want dequeue to panic with "queue is empty", got %q`,
			q, caughtPanic)
	}
}

func catchPanic(f func()) (caughtPanic error) {
	defer func() {
		if e := recover(); e != nil {
			caughtPanic = fmt.Errorf("%v", e)
		}
	}()

	f()
	return
}
