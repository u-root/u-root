// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package vector implements persistent vector.
package vector

// Vectors from elvish. I'm not convinced a lot of the claims were correct
// so I stripped the comment. I'm not sure how much of what this does we need,
// e.g. SubVector, but we'll see.
type Vector interface {
	// Len returns the length of the vector.
	Len() int
	// Index returns the i-th element of the vector, if it exists. The second
	// return value indicates whether the element exists.
	Index(i int) (interface{}, bool)
	// Assoc returns an almost identical Vector, with the i-th element
	// replaced. If the index is smaller than 0 or greater than the length of
	// the vector, it returns nil. If the index is equal to the size of the
	// vector, it is equivalent to Cons.
	Assoc(i int, val interface{}) Vector
	// Cons returns an almost identical Vector, with an additional element
	// appended to the end.
	Cons(val interface{}) Vector
	// Pop returns an almost identical Vector, with the last element removed. It
	// returns nil if the vector is already empty.
	Pop() Vector
	// SubVector returns a subvector containing the elements from i up to but
	// not including j.
	SubVector(i, j int) Vector
	// Iterator returns an iterator over the vector.
	Iterator() Iterator
}

// Iterator is an iterator over vector elements. It can be used like this:
//
//     for it := v.Iterator(); it.HasElem(); it.Next() {
//         elem := it.Elem()
//         // do something with elem...
//     }
type Iterator interface {
	// Elem returns the element at the current position.
	Elem() interface{}
	// HasElem returns whether the iterator is pointing to an element.
	HasElem() bool
	// Next moves the iterator to the next position.
	Next()
}

type vector struct {
	pos   int
	nodes []interface{}
}

// Empty is an empty Vector.
var Empty Vector = &vector{}

// Count returns the number of elements in a Vector.
func (v *vector) Len() int {
	return len(v.nodes)
}

func (v *vector) Index(i int) (interface{}, bool) {
	if i < 0 || i >= len(v.nodes) {
		return nil, false
	}

	return v.nodes[i], true
}

func (v *vector) Assoc(i int, val interface{}) Vector {
	if i < 0 || i > len(v.nodes) {
		return nil
	}
	n := append(append(v.nodes[:i], val), v.nodes[i:]...)
	return &vector{nodes: n}
}

func (v *vector) Cons(val interface{}) Vector {
	return &vector{nodes: append(v.nodes, val)}
}

func (v *vector) Pop() Vector {
	switch len(v.nodes) {
	case 0:
		return nil
	case 1:
		return Empty
	}
	return &vector{nodes: v.nodes[:len(v.nodes)-1]}
}

func (v *vector) SubVector(begin, end int) Vector {
	if begin < 0 || begin > end || end > len(v.nodes) {
		return nil
	}
	return &subVector{v, begin, end}
}

func (v *vector) Iterator() Iterator {
	return newIterator(v)
}

type subVector struct {
	v     *vector
	begin int
	end   int
}

func (s *subVector) Len() int {
	return s.end - s.begin
}

func (s *subVector) Index(i int) (interface{}, bool) {
	if i < 0 || s.begin+i >= s.end {
		return nil, false
	}
	return s.v.Index(s.begin + i)
}

func (s *subVector) Assoc(i int, val interface{}) Vector {
	if i < 0 || s.begin+i > s.end {
		return nil
	} else if s.begin+i == s.end {
		return s.Cons(val)
	}
	return s.v.Assoc(s.begin+i, val).SubVector(s.begin, s.end)
}

func (s *subVector) Cons(val interface{}) Vector {
	return s.v.Assoc(s.end, val).SubVector(s.begin, s.end+1)
}

func (s *subVector) Pop() Vector {
	switch s.Len() {
	case 0:
		return nil
	case 1:
		return Empty
	default:
		return s.v.SubVector(s.begin, s.end-1)
	}
}

func (s *subVector) SubVector(i, j int) Vector {
	return s.v.SubVector(s.begin+i, s.begin+j)
}

func (s *subVector) Iterator() Iterator {
	return newIteratorWithRange(s.v, s.begin, s.end)
}

type iterator struct {
	v     *vector
	index int
	end   int
}

func newIterator(v *vector) *iterator {
	return newIteratorWithRange(v, 0, v.Len())
}

func newIteratorWithRange(v *vector, begin, end int) *iterator {
	it := &iterator{v, begin, end}
	return it
}

func (it *iterator) Elem() interface{} {
	if it.index >= len(it.v.nodes) {
		return nil
	}
	return it.v.nodes[it.index]
}

func (it *iterator) HasElem() bool {
	return it.index < it.end
}

func (it *iterator) Next() {
	if it.index+1 >= len(it.v.nodes) {
		it.index++
		return
	}
	it.index++
}
