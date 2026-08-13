// Copyright 2026 The HuaTuo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

// primitives a slice of primitive with some methods
type primitives []primitive

// combine combines multiple primitives until nothing is left to combine
func (p *primitives) combine() *primitives {
	// nothing to combine
	if p == nil || len(*p) == 0 || len(*p) == 1 {
		return p
	}
	// The simplest first cut is to have each one combine with its neighbor.
	// It isn't perfect - e.g. we will miss A combining with D - but it is a good start.
	pd := *p
	list := make(primitives, 0)
	var (
		prev, elm primitive
		lastMatch bool
		i         int
	)
	for i, elm = range pd {
		// do not bother combining with myself
		if i == 0 {
			prev = elm
			continue
		}
		if n := prev.Combine(&elm); n != nil {
			lastMatch = true
			list = append(list, *n)
			prev = *n
		} else {
			lastMatch = false
			list = append(list, prev)
			prev = elm
		}
	}
	// add the last element if it was not merged with the previous
	if !lastMatch {
		list = append(list, elm)
	}
	return &list
}

func (p *primitives) equal(o *primitives) bool {
	if o == nil {
		return false
	}
	pd, od := *p, *o
	if len(pd) != len(od) {
		return false
	}
	for i, p1 := range pd {
		if !p1.Equal(od[i]) {
			return false
		}
	}
	return true
}
