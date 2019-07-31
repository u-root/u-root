// Copyright 2018 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package abi

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// The Flag interface has three functions:
// Match, to see if a given value is covered by the Flag.
// Mask, to return which bits in the flag are covered.
// String, to print a string determined by the flag and possibly the value.
type Flag interface {
	Match(uint64) bool
	Mask() uint64
	String(uint64) string
}

// A FlagSet is a slice of Flags
type FlagSet []Flag

// A BitFlag is one or more bits and a name.
// The name is printed if all bits are matched
// (i.e. enabled) in a value. Typically, the Value
// is only one bit.
type BitFlag struct {
	Name  string
	Value uint64
}

// Match returns true if the val and the BitFlag Value
// are the same, i.e. the flag is set.
func (f *BitFlag) Match(val uint64) bool {
	return f.Value&val == f.Value
}

// String returns the string value of the BitFlag.
// The parameter is ignored.
func (f *BitFlag) String(_ uint64) string {
	return f.Name
}

// Mask returns a mask for the BitFlag, namely, the value.
func (f *BitFlag) Mask() uint64 {
	return f.Value
}

// Field is used to extract named fields, e.g. the
// Tries field in the ChromeOS GPT. The value
// is masked and shifted right. We may at some
// point want a format string.
type Field struct {
	Name    string
	BitMask uint64
	Shift   uint64
}

// Match always matches in a Field, regardless of value.
func (f *Field) Match(val uint64) bool {
	return true
}

// String returns the part or the val covered by this Field.
// The bits are extracted, shifted, and printed as hex.
func (f *Field) String(val uint64) string {
	return fmt.Sprintf("%s=%#x", f.Name, (val&f.BitMask)>>f.Shift)
}

// Mask returns the bits covered by the Field.
func (f *Field) Mask() uint64 {
	return f.BitMask
}

// Value is the simplest implementation for Flags.
// If an entire uint64 matches the Value, then Match
// will be true. Note that Value could be implemented
// by a Field with a Mask of MaxUint64, but Value
// is more convenient to use.
type Value struct {
	Name  string
	Value uint64
}

// String returns the name of the value associated with `val`.
func (e *Value) String(_ uint64) string {
	return e.Name
}

// Match determines if a Value matches the argument, meaning the
// Value and the arg are the same.
func (e *Value) Match(val uint64) bool {
	return e.Value == val
}

// Mask returns the bits covered by a Value, in this case MaxUint64.
func (e *Value) Mask() uint64 {
	return math.MaxUint64
}

// Parse returns a pretty version of a FlagSet, using the flag names for known flags.
// Unknown flags are printed as numeric if the Flagset did not cover all the bits
// in the argument.
func (s FlagSet) Parse(val uint64) string {
	var flags []string
	var clr uint64

	for _, f := range s {
		if f.Match(val) {
			flags = append(flags, f.String(val))
			val &^= f.Mask()
			clr |= f.Mask()
		}
	}

	// If there are bits we did not check, then print out
	// whatever is left, *if it is non-zero*. This is a bit
	// hokey, but at the same time, it seems the most usable.
	// We may, later, want to print out the bits not covered.
	// It's possible to miss bits and have val be 0 if the bits
	// we missed are 0, as in the earlier code.
	if clr != math.MaxUint64 && val != 0 {
		flags = append(flags, "0x"+strconv.FormatUint(val, 16))
	}

	return strings.Join(flags, "|")
}
