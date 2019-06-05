// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package hashmap implements persistent hashmap.
package hashmap

import "log"

// Equal is the type of a function that reports whether two keys are equal.
type Equal func(k1, k2 interface{}) bool

// Hash is the type of a function that returns the hash code of a key.
type Hash func(k interface{}) uint32

// New takes an equality function and a hash function, and returns an empty
// Map.
func New(e Equal, h Hash) Map {
	return &hashMap{m: make(map[interface{}]interface{})}
}

type hashMap struct {
	m map[interface{}]interface{}
}

func (m *hashMap) Len() int {
	return len(m.m)
}

func (m *hashMap) Index(k interface{}) (interface{}, bool) {
	i, ok := m.m[k]
	return i, ok
}

func (m *hashMap) Assoc(k, v interface{}) Map {
	item, ok := m.m[k]
	if ok {
		log.Panicf("Storing %v at %v want ! ok, got %v", v, k, item)
	}
	n := &hashMap{m: make(map[interface{}]interface{})}
	for k, v := range m.m {
		n.m[k] = v
	}
	n.m[k] = v
	return n
}

func (m *hashMap) Dissoc(del interface{}) Map {
	n := &hashMap{m: make(map[interface{}]interface{})}
	for k, v := range m.m {
		if k == del {
			continue
		}
		n.m[k] = v
	}
	return n
}

type kv struct {
	k interface{}
	v interface{}
}

type iter struct {
	m   *hashMap
	cur *kv
	kvs chan (*kv)
}

func (i *iter) Elem() (interface{}, interface{}) {
	return i.cur.k, i.cur.v
}

func (i *iter) HasElem() bool {
	return i.cur != nil
}

func (i *iter) Next() {
	i.cur = <-i.kvs
}
func (m *hashMap) Iterator() Iterator {
	c := make(chan *kv)
	go func() {
		for k, v := range m.m {
			c <- &kv{k, v}
		}
		close(c)
	}()
	return &iter{m: m, kvs: c, cur: <-c}
}
