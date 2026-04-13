// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type queue struct {
	q []string
}

func (q *queue) enqueue(value string) {
	q.q = append(q.q, value)
}

func (q *queue) dequeue() (string, bool) {
	if len(q.q) == 0 {
		return "", false
	}

	result := q.q[0]

	q.q = q.q[1:]

	return result, true
}
