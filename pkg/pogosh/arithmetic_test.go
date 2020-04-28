// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import (
	"math/big"
	"reflect"
	"testing"
)

// Fake variables for the test.
func getVar(name string) *big.Int {
	fakeVars := map[string]*big.Int{
		"x":   big.NewInt(1),
		"y":   big.NewInt(123),
		"z":   big.NewInt(-1),
		"xyz": big.NewInt(42),
	}
	val, ok := fakeVars[name]
	if !ok {
		return big.NewInt(0)
	}
	return val
}

// The positive tests are expected to pass parsing.
var arithmeticPositiveTests = []struct {
	name string
	in   string
	out  *big.Int
}{
	// Constants
	{"Zero",
		"0",
		big.NewInt(0),
	},
	{"Decimal",
		"123",
		big.NewInt(123),
	},
	{"Octal",
		"0123",
		big.NewInt(0123),
	},
	{"Hex",
		"0x123",
		big.NewInt(0x123),
	},

	// Order of operation
	{"BEDMAS1",
		"1+2*3",
		big.NewInt(7),
	},
	{"BEDMAS2",
		"(1+2)*3",
		big.NewInt(9),
	},
	{"BEDMAS3",
		"1*2+3",
		big.NewInt(5),
	},
	{"BEDMAS4",
		"-1*-2+-+-3",
		big.NewInt(5),
	},

	// Associativity
	{"Associativity1",
		"1-2-3",
		big.NewInt(-4),
	},
	{"Associativity2",
		"1-2*3-4",
		big.NewInt(-9),
	},

	// Spacing
	{"Spacing",
		" 1 + 3 * 4 / 3 - 8 ",
		big.NewInt(-3),
	},

	// Conditional Operator
	{"Conditional",
		" 0 ? 4 : 1 ? 5 : 6 ",
		big.NewInt(5),
	},

	// Variables
	{"Variables1",
		"x",
		big.NewInt(1),
	},
	{"Variables2",
		"x + y",
		big.NewInt(124),
	},
	{"Variables3",
		"xyz * 2",
		big.NewInt(84),
	},
}

func TestArithmeticPositive(t *testing.T) {
	for _, tt := range arithmeticPositiveTests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arithmetic{
				getVar: getVar,
				input:  tt.in,
			}
			got := a.evalExpression()

			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("got %v, want %v", got, tt.out)
			}
		})
	}
}
