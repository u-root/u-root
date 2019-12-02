// Copyright -2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unit implements parsing for string values with units.
package unit

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

const (
	_ = iota
	// K is 1024 byte
	K int64 = 1 << (10 * iota)
	// M is 1024 K
	M
	// G is 1024 M
	G
	// T is 1024 G
	T
	// P is 1024 T
	P
	// E is 1024 P
	E
)

// DefaultUnits is the default unit mapping as used by many standard cli-tools.
var DefaultUnits = map[string]int64{
	"B":  1,
	"K":  K,
	"M":  M,
	"G":  G,
	"T":  T,
	"P":  P,
	"E":  E,
	"kB": 1000,
	"KB": 1000,
	"MB": 1000 * 1000,
	"GB": 1000 * 1000 * 1000,
	"TB": 1000 * 1000 * 1000 * 1000,
	"PB": 1000 * 1000 * 1000 * 1000 * 1000,
	"EB": 1000 * 1000 * 1000 * 1000 * 1000 * 1000,
}

// Sign is the sign associated with a unit's value.
type Sign uint8

const (
	// None signals that no explicit sign is set.
	None Sign = iota
	// Negative signals that an explicit negative sign is set.
	Negative
	// Positive signals that an explicit positive sign is set.
	Positive
)

// Value is any value that can be represented by a unit.
//
// Value implements flag.Value and flag.Getter.
type Value struct {
	// unit is the associated unit.
	unit *Unit

	// value is the integer value.
	Value int64

	// sign is the explicit sign given by the string converted to the
	// integer.
	ExplicitSign Sign

	// set to false if this is the default value, true if the the option was given
	IsSet bool
}

// Unit is a map of unit names to conversion multipliers.
//
// There must be a unit that maps to 1.
type Unit struct {
	mapping map[string]int64
}

// NewUnit returns a new Unit given a mapping 'm'.
func NewUnit(m map[string]int64) (*Unit, error) {
	var found bool
	for _, mult := range m {
		if mult == 1 {
			found = true
		} else if mult == 0 {
			return nil, fmt.Errorf("mapping contains unit that maps to illegal multiplier 0 for %v", m)
		}
	}
	if !found {
		return nil, fmt.Errorf("could not find unit that maps to multiplier 1 for %v", m)
	}
	return &Unit{m}, nil
}

// MustNewUnit is like NewUnit but panics if the mapping 'm' is not valid.
func MustNewUnit(m map[string]int64) *Unit {
	u, err := NewUnit(m)
	if err != nil {
		panic(fmt.Sprintf("MustNewUnit error: %v", err))
	}
	return u
}

// NewValue returns a Value based on a Unit. It's value is set to "value".
// "value" is the initial value, an explicit sign can be set, but is usually
// "None".
func (u *Unit) NewValue(value int64, explicitSign Sign) (*Value, error) {
	if (value < 0 && explicitSign == Positive) || (value > 0 && explicitSign == Negative) {
		return nil, errors.New("Invalid value/explicitSign combination")
	}

	return &Value{
		unit:         u,
		Value:        value,
		ExplicitSign: explicitSign,
	}, nil
}

// MustNewValue is like NewValue but panics if the new Value could not be
// generated.
func (u *Unit) MustNewValue(value int64, explicitSign Sign) *Value {
	v, err := u.NewValue(value, explicitSign)
	if err != nil {
		panic(fmt.Sprintf("MustNewValue error: %v", err))
	}
	return v
}

// ValueFromString converts the given string to a Value.
// The string is allowed to have '+'/'-' as prefix, followed by a number, and
// an optional unit name as defined in its mapping.
func (u *Unit) ValueFromString(str string) (*Value, error) {
	s := &Value{unit: u}

	if err := s.Set(str); err != nil {
		return nil, err
	}
	return s, nil
}

// String implements flag.Value.String and fmt.Stringer.
func (s Value) String() string {
	var bestName string
	bestMult := int64(1)
	if s.unit == nil {
		return ""
	}
	for name, mult := range s.unit.mapping {
		if s.Value%mult == 0 && mult >= bestMult {
			bestName = name
			bestMult = mult
		}
	}
	var sign string
	if s.ExplicitSign == Positive {
		sign = "+"
	}
	if bestName == "" {
		return fmt.Sprintf("%s%d (no unit)", sign, s.Value)
	}
	return fmt.Sprintf("%s%d%s", sign, s.Value/bestMult, bestName)
}

// Get implements flag.Getter.Get.
func (s Value) Get() interface{} {
	return s
}

// Set implements flag.Value.Set.
func (s *Value) Set(str string) error {
	if len(str) == 0 {
		return fmt.Errorf("invalid size %q", str)
	}

	start, end := 0, len(str)
	if str[0] == '+' {
		s.ExplicitSign = Positive
		start++
	} else if str[0] == '-' {
		s.ExplicitSign = Negative
		start++
	}

	for i, r := range str[start:] {
		if unicode.IsLetter(r) {
			end = start + i
			break
		}
	}

	value, err := strconv.ParseInt(str[:end], 10, 64)
	if err != nil {
		return fmt.Errorf("could not convert %q to size: %v", str, err)
	}

	unitName := str[end:]
	mult, ok := s.unit.mapping[unitName]
	if !ok {
		if len(unitName) != 0 {
			return fmt.Errorf("unit %q is not valid", unitName)
		}
		mult = 1
	}
	if (value > 0 && (math.MaxInt64/mult < value)) || (value < 0 && (math.MinInt64/mult > value)) {
		return fmt.Errorf("size %d with unit %q is out of range", value, unitName)
	}
	s.Value = value * mult
	s.IsSet = true
	return nil
}

// Type implements pflag.Value.Type.
func (s *Value) Type() string { return "unit" }
