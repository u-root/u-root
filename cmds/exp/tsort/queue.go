// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type queue struct {
	q []string
}

func (q *queue) enqueue(value string) {
	q.q = append(q.q, value)
}

func (q *queue) dequeue() string {
	if q.isEmpty() {
		panic("queue is empty")
	}

	result := q.q[0]

	q.q = q.q[1:]

	return result
}

func (q *queue) isEmpty() bool {
	return len(q.q) == 0
}
