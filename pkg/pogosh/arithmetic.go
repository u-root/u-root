// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import (
	"fmt"
	"math/big"
)

// After variable substitution has occured, this
// evaluates the arithmetic expression.

// These rules come from ISO C standard, section 6.5 Expression and section
// 6.4.4.1 Integer Constants.
// LL(3) grammar

// Params:
//	 map[string]Variable
// Returns:
//   big.Int: value
//   map[string]big.Int: assignments
// Panics: on parse error

// Arithmetic computes value in $((...)) expression.
// TODO: why public?
type Arithmetic struct {
	getVar func(string) *big.Int
	setVar func(string, *big.Int)

	// Initial string
	input string

	// Remaining unparsed string
	rem string
}

// if x == 0, 0
// otherwise, 1
func bigBool(x *big.Int) *big.Int {
	if x.BitLen() == 0 {
		return big.NewInt(0)
	}
	return big.NewInt(1)
}

func asBig(b bool) *big.Int {
	if b {
		return big.NewInt(1)
	}
	return big.NewInt(0)
}

// Wrapper which should be used outside this file.
func (a *Arithmetic) evalExpression() *big.Int {
	// TODO: this might not be standard and might be inefficient
	// Augment for LL(3)
	a.rem = a.input + "\000\000\000"
	val := a.evalAssignmentExpression()
	a.evalSpaces()
	if a.rem != "\000\000\000" {
		panic("Expected EOF at " + a.rem)
	}
	return val
}

func (a *Arithmetic) evalSpaces() {
	for a.rem[0] == ' ' || a.rem[0] == '\t' || a.rem[0] == '\n' { // TODO: more space characters
		a.rem = a.rem[1:]
	}
}

// [_0-9a-zA-Z]
func isIdentifierChar(char byte) bool {
	return char == '_' || ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9')
}

func isDecimal(char byte) bool {
	return '0' <= char && char <= '9'
}

func isHex(char byte) bool {
	return ('0' <= char && char <= '9') || ('a' <= char && char <= 'f') || ('A' <= char && char <= 'F')
}

func isOctal(char byte) bool {
	return ('0' <= char && char <= '7')
}

// Identifier ::= [_a-zA-Z][_0-9a-zA-Z]*
func (a *Arithmetic) evalIdentifier() *big.Int {
	// This does not check the first character is not numeric because it has
	// already been done as part of the callee's FIRST set.
	i := 0
	for isIdentifierChar(a.rem[i]) {
		i++
	}

	identifier := a.rem[:i]
	a.rem = a.rem[i:]
	return a.getVar(identifier)
}

// Constant ::= DecimalConstant | OctalConstant | HexadecimalConstant
// DecimalConstant ::= [1-9][0-9]*
// OctalConstant ::= 0[0-9]*
// HexadecimalConstant ::= 0x[0-9]* | 0X[0-9]*
func (a *Arithmetic) evalConstant() *big.Int {
	// Get length of constant.
	i := 0
	if a.rem[i] == '0' {
		i++
		if a.rem[i] == 'x' || a.rem[i] == 'X' {
			i++
			for isHex(a.rem[i]) {
				i++
			}
		} else {
			for isOctal(a.rem[i]) {
				i++
			}
		}
	} else {
		for isDecimal(a.rem[i]) {
			i++
		}
	}

	var val big.Int
	_, ok := val.SetString(a.rem[:i], 0)
	if !ok {
		panic("Not a valid constant")
	}
	a.rem = a.rem[i:]
	return &val
}

// PrimaryExpression ::= Identifier | Constant | '(' AssignmentExpression ')'
func (a *Arithmetic) evalPrimaryExpression() *big.Int {
	a.evalSpaces()
	char := a.rem[0]
	switch {
	case '0' <= char && char <= '9':
		return a.evalConstant()
	case isIdentifierChar(char):
		return a.evalIdentifier()
	case char == '(':
		a.rem = a.rem[1:]
		val := a.evalAssignmentExpression()
		a.evalSpaces()
		if a.rem[0] != ')' {
			panic("No matching closing parenthesis")
		}
		a.rem = a.rem[1:]
		return val
	default:
		panic(fmt.Sprintf("Expected identifier or constant at %c", char))
	}
}

// UnaryExpression ::= PrimaryExpression
// UnaryExpression ::= UnaryOperator UnaryExpression
// UnaryOperator ::= '+' | '-' | '~' | '!'
func (a *Arithmetic) evalUnaryExpression() *big.Int {
	val := big.NewInt(0)
	a.evalSpaces()
	switch a.rem[0] {
	case '+':
		a.rem = a.rem[1:]
		val = a.evalUnaryExpression()
	case '-':
		a.rem = a.rem[1:]
		val.Neg(a.evalUnaryExpression())
	case '~':
		a.rem = a.rem[1:]
		val.Not(a.evalUnaryExpression())
	case '!':
		a.rem = a.rem[1:]
		val.Xor(big.NewInt(1), bigBool(a.evalUnaryExpression()))
	default:
		val = a.evalPrimaryExpression()
	}
	return val
}

// MultiplicativeExpression ::= UnaryExpression MultiplicativeExpression2
// MultiplicativeExpression2 ::= MultiplicativeOperator MultiplicativeExpression |
// MultiplicativeOperator ::= '*' | '/' | '%'
func (a *Arithmetic) evalMultiplicativeExpression() *big.Int {
	val := a.evalUnaryExpression()
	for {
		a.evalSpaces()
		switch a.rem[0] {
		case '*':
			a.rem = a.rem[1:]
			val.Mul(val, a.evalUnaryExpression())
		case '/':
			a.rem = a.rem[1:]
			val.Div(val, a.evalUnaryExpression())
		case '%':
			a.rem = a.rem[1:]
			val.Mod(val, a.evalUnaryExpression())
		default:
			return val
		}
	}
}

// AdditiveExpression ::= MultiplicativeExpression AdditiveExpression2
// AdditiveExpression2 ::= AdditiveOperator AdditiveExpression |
// AdditiveOperator ::= '+' | '-'
func (a *Arithmetic) evalAdditiveExpression() *big.Int {
	val := a.evalMultiplicativeExpression()
	for {
		a.evalSpaces()
		switch a.rem[0] {
		case '+':
			a.rem = a.rem[1:]
			val.Add(val, a.evalMultiplicativeExpression())
		case '-':
			a.rem = a.rem[1:]
			val.Sub(val, a.evalMultiplicativeExpression())
		default:
			return val
		}
	}
}

// ShiftExpression ::= AdditiveExpression ShiftExpression2
// ShiftExpression2 ::= ShiftOperator ShiftExpression |
// ShiftOperator ::= '<<' | '>>'
func (a *Arithmetic) evalShiftExpression() *big.Int {
	val := a.evalAdditiveExpression()
	for {
		a.evalSpaces()
		switch a.rem[:2] {
		case "<<":
			a.rem = a.rem[2:]
			// TODO: might be undefined if > UINT64_MAX
			val.Lsh(val, uint(a.evalAdditiveExpression().Uint64()))
		case ">>":
			a.rem = a.rem[2:]
			// TODO: might be undefined if > UINT64_MAX
			val.Rsh(val, uint(a.evalAdditiveExpression().Uint64()))
		default:
			return val
		}
	}
}

// RelationalExpression ::= ShiftExpression RelationalExpression2
// RelationalExpression2 ::= RelationalOperator RelationalExpression |
// RelationalOperator ::= '<' | '>' | '<=' | '>='
func (a *Arithmetic) evalRelationalExpression() *big.Int {
	val := a.evalShiftExpression()
	for {
		a.evalSpaces()
		switch {
		case a.rem[:2] == "<=":
			a.rem = a.rem[2:]
			val = asBig(val.Cmp(a.evalShiftExpression()) <= 0)
		case a.rem[:2] == ">=":
			a.rem = a.rem[2:]
			val = asBig(val.Cmp(a.evalShiftExpression()) >= 0)
		case a.rem[0] == '<' && a.rem[1] != '<':
			a.rem = a.rem[2:]
			val = asBig(val.Cmp(a.evalShiftExpression()) < 0)
		case a.rem[0] == '>' && a.rem[1] != '>':
			a.rem = a.rem[2:]
			val = asBig(val.Cmp(a.evalShiftExpression()) > 0)
		default:
			return val
		}
	}
}

// EqualityExpression ::= RelationalExpression EqualityExpression2
// EqualityExpression2 ::= EqualityOperator EqualityExpression
// EqualityOperator ::= '==' | '!='
func (a *Arithmetic) evalEqualityExpression() *big.Int {
	val := a.evalRelationalExpression()
	for {
		a.evalSpaces()
		switch a.rem[:2] {
		case "==":
			a.rem = a.rem[2:]
			val = asBig(val.Cmp(a.evalRelationalExpression()) == 0)
		case "!=":
			a.rem = a.rem[2:]
			val = asBig(val.Cmp(a.evalRelationalExpression()) != 0)
		default:
			return val
		}
	}
}

// ANDExpression ::= EqualityExpression AndExpression2
// ANDExpression2 ::= '&' ANDExpression |
func (a *Arithmetic) evalANDExpression() *big.Int {
	val := a.evalEqualityExpression()
	for {
		a.evalSpaces()
		switch {
		case a.rem[0] == '&' && a.rem[1] != '&':
			a.rem = a.rem[2:]
			val.And(val, bigBool(a.evalEqualityExpression()))
		default:
			return val
		}
	}
}

// ExclusiveORExpression ::= ANDExpression ExclusiveORExpression2
// ExclusiveORExpression2 ::= '^' ExclusiveORExpression |
func (a *Arithmetic) evalExclusiveORExpression() *big.Int {
	val := a.evalANDExpression()
	for {
		a.evalSpaces()
		switch a.rem[0] {
		case '^':
			a.rem = a.rem[1:]
			val.Xor(val, bigBool(a.evalANDExpression()))
		default:
			return val
		}
	}
}

// InclusiveORExpression ::= ExclusiveORExpression InclusiveORExpression2
// InclusiveORExpression2 ::= '|' InclusiveORExpression |
func (a *Arithmetic) evalInclusiveORExpression() *big.Int {
	val := a.evalExclusiveORExpression()
	for {
		a.evalSpaces()
		switch {
		case a.rem[0] == '|' && a.rem[1] != '|':
			a.rem = a.rem[2:]
			val.Or(val, bigBool(a.evalExclusiveORExpression()))
		default:
			return val
		}
	}
}

// LogicalANDExpression ::= InclusiveORExpression LogicalANDExpression2
// LogicalANDExpression2 ::= '&&' InclusiveORExpression |
func (a *Arithmetic) evalLogicalANDExpression() *big.Int {
	val := a.evalInclusiveORExpression()
	for {
		a.evalSpaces()
		switch a.rem[:2] {
		case "&&":
			a.rem = a.rem[2:]
			val.And(bigBool(val), bigBool(a.evalInclusiveORExpression()))
		default:
			return val
		}
	}
}

// LogicalORExpression ::= LogicalANDExpression LogicalORExpression2
// LogicalORExpression2 ::= '||' LogicalORExpression |
func (a *Arithmetic) evalLogicalORExpression() *big.Int {
	val := a.evalLogicalANDExpression()
	for {
		a.evalSpaces()
		switch a.rem[:2] {
		case "||":
			a.rem = a.rem[2:]
			val.Or(bigBool(val), bigBool(a.evalLogicalANDExpression()))
		default:
			return val
		}
	}
}

// ConditionalExpression ::= LogicalORExpression ConditionalExpression2
// ConditionalExpression2 ::= '?' AssignmentExpression ':' ConditionalExpression |
func (a *Arithmetic) evalConditionalExpression() *big.Int {
	val := a.evalLogicalORExpression()

	a.evalSpaces()
	if a.rem[0] != '?' {
		return val
	}
	a.rem = a.rem[1:]
	trueVal := a.evalAssignmentExpression()

	if a.rem[0] != ':' {
		panic("Bad conditional expression")
	}
	a.rem = a.rem[1:]
	falseVal := a.evalConditionalExpression()

	if val.BitLen() == 0 {
		return falseVal
	}
	return trueVal
}

// AssignmentExpression ::= ConditionalExpression
// AssignmentExpression ::= Identifier AssignmentOperator AssignmentExpression
// AssignmentOperator ::= '=' | '*=' | '/=' | '%=' | '+=' | '-=' | '<<='
//                      | '>>=' | '&=' | '^=' | '|='
func (a *Arithmetic) evalAssignmentExpression() *big.Int {
	val := a.evalConditionalExpression()
	// TODO: assignment
	// TODO: some other follow sets need to be updated
	return val
}
