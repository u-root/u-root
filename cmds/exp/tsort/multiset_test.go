// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestMultisetIsEmpty(t *testing.T) {
	m := newMultiset()
	if !m.isEmpty() {
		t.Fatalf("multiset %v: want empty, got non-empty", m)
	}

	m.add("a", 1)
	if m.isEmpty() {
		t.Fatalf("multiset %v: want non-empty, got empty", m)
	}

	m.add("b", 2)
	if m.isEmpty() {
		t.Fatalf("multiset %v: want non-empty, got empty", m)
	}

	m.removeOne("a")
	if m.isEmpty() {
		t.Fatalf("multiset %v: want non-empty, got empty", m)
	}

	m.removeOne("b")
	if m.isEmpty() {
		t.Fatalf("multiset %v: want non-empty, got empty", m)
	}

	m.removeOne("b")
	if !m.isEmpty() {
		t.Fatalf("multiset %v: want empty, got non-empty", m)
	}
}

func TestMultisetHas(t *testing.T) {
	m := newMultiset()
	if m.has("a") {
		t.Fatalf(`multiset %v: want "a" to be absent, but it was present`, m)
	}

	m.add("a", 1)
	if !m.has("a") {
		t.Fatalf(`multiset %v: want "a" to be present, but it was absent`, m)
	}

	m.add("b", 2)
	if !m.has("b") {
		t.Fatalf(`multiset %v: want "b" to be present, but it was absent`, m)
	}

	m.removeOne("a")
	if m.has("a") {
		t.Fatalf(`multiset %v: want "a" to be absent, but it was present`, m)
	}

	m.removeOne("b")
	if !m.has("b") {
		t.Fatalf(`multiset %v: want "b" to be present, but it was absent`, m)
	}

	m.removeOne("b")
	if m.has("b") {
		t.Fatalf(`multiset %v: want "b" to be absent, but it was present`, m)
	}
}

func TestMultisetCount(t *testing.T) {
	m := newMultiset()
	if c := m.count("a"); c != 0 {
		t.Fatalf(`multiset %v: want count of "a" to be 0, got %d`, m, c)
	}

	m.add("a", 1)
	if c := m.count("a"); c != 1 {
		t.Fatalf(`multiset %v: want count of "a" to be 1, got %d`, m, c)
	}

	m.add("b", 2)
	if c := m.count("b"); c != 2 {
		t.Fatalf(`multiset %v: want count of "b" to be 2, got %d`, m, c)
	}

	m.removeOne("a")
	if c := m.count("a"); c != 0 {
		t.Fatalf(`multiset %v: want count of "a" to be 0, got %d`, m, c)
	}

	m.removeOne("b")
	if c := m.count("b"); c != 1 {
		t.Fatalf(`multiset %v: want count of "b" to be 1, got %d`, m, c)
	}

	m.removeOne("b")
	if c := m.count("b"); c != 0 {
		t.Fatalf(`multiset %v: want count of "b" to be 0, got %d`, m, c)
	}
}

func TestMultisetAdd(t *testing.T) {
	m := newMultiset()

	caughtPanic := catchPanic(func() { m.add("a", 0) })
	if caughtPanic == nil ||
		caughtPanic.Error() != "count is non-positive: 0" {
		t.Fatalf(
			`multiset %v: want add to panic with "count is non-positive: 0", got %#v`,
			m, caughtPanic)
	}

	caughtPanic = catchPanic(func() { m.add("a", -1) })
	if caughtPanic == nil || caughtPanic.Error() != "count is non-positive: -1" {
		t.Fatalf(
			`multiset %v: want add to panic with "count is non-positive: -1", got %#v`,
			m, caughtPanic)
	}
}

func TestMultisetRemoveOne(t *testing.T) {
	m := newMultiset()
	m.add("a", 1)

	caughtPanic := catchPanic(func() { m.removeOne("b") })
	if caughtPanic == nil ||
		caughtPanic.Error() != "multiset does not have value" {
		t.Fatalf(
			`multiset %v: want removeOne to panic with "multiset does not have value", got %#v`,
			m, caughtPanic)
	}
}

func TestMultisetForEachUnique(t *testing.T) {
	m := newMultiset()
	m.add("a", 1)
	m.add("b", 2)
	m.add("c", 3)

	var actual []string
	m.forEachUnique(func(value string) bool {
		actual = append(actual, value)
		return true
	})
	expected := []string{"a", "b", "c"}
	if diff := orderInsensitiveDiff(actual, expected); diff != "" {
		t.Fatalf(
			"forEachUnique mismatch (-actual +expected):\n%s",
			diff)
	}

	var values []string
	m.forEachUnique(func(value string) bool {
		values = append(values, value)
		return false
	})
	if len(values) != 1 {
		t.Fatalf("expected forEachUnique to break when false is returned")
	}
}
