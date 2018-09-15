// Copyright 2010-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// the forth package is designed for use by programs
// needing to evaluate command-line arguments or simple
// expressions to set program variables. It is designed
// to map host names to numbers. We use it to
// easily convert host names and IP addresses into
// parameters.
// The language
// is a Forth-like postfix notation. Elements are
// either commands or strings. Strings are
// immediately pushed. Commands consume stack variables
// and produce new ones.
// Simple examples:
// push hostname, strip alpha characters to produce a number. If your
// hostname is sb47, top of stack will be left with 47.
// hostname  hostbase
// Get the hostbase, if it is 0 mod 20, return the hostbase / 20,
// else return hostbase mod 20
//
// hostname hostbase dup 20 / swap 20 % dup ifelse
//
// At the end of the evaluation the stack should have one element
// left; that element is popped and returned. It is an error (currently)
// to return with a non-empty stack.
// This package was used for real work at Sandia National Labs from 2010 to 2012 and possibly later.
// Some of the use of error may seem a bit weird but the creation of this package predates the
// creation of the error type (it was still an os thing back then).
package forth

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type ForthOp func(f Forth)

type forthstack struct {
	stack []string
}

var opmap = map[string]ForthOp{
	"+":        plus,
	"-":        sub,
	"*":        times,
	"/":        div,
	"%":        mod,
	"swap":     swap,
	"ifelse":   ifelse,
	"hostname": hostname,
	"hostbase": hostbase,
	"strcat":   strcat,
	"roundup":  roundup,
	"dup":      dup,
}

// Forth is an interface used by the package. The interface
// requires definition of Push, Pop, Length, Empty (convenience function
// meaning Length is 0), Newop (insert a new or replacement operator),
// and Reset (clear the stack, mainly diagnostic)
type Forth interface {
	Push(string)
	Pop() string
	Length() int
	Empty() bool
	Newop(string, ForthOp)
	Reset()
	Stack() []string
}

// New creates a new stack
func New() Forth {
	f := new(forthstack)
	return f
}

// Newop creates a new operation. We considered having
// an opmap per stack but don't feel the package requires it
func (f *forthstack) Newop(n string, op ForthOp) {
	opmap[n] = op
}

func Ops() map[string]ForthOp {
	return opmap
}

// Reset resets the stack to empty
func (f *forthstack) Reset() {
	f.stack = f.stack[0:0]
}

// Return the stack as a []string
func (f *forthstack) Stack() []string {
	return f.stack
}

// Push pushes the string on the stack.
func (f *forthstack) Push(s string) {
	f.stack = append(f.stack, s)
	//fmt.Printf("push: %v: stack: %v\n", s, f.stack)
}

// Pop pops the stack. If the stack is Empty Pop will panic.
// Eval recovers() the panic.
func (f *forthstack) Pop() (ret string) {

	if len(f.stack) < 1 {
		panic(errors.New("Empty stack"))
	}
	ret = f.stack[len(f.stack)-1]
	f.stack = f.stack[0 : len(f.stack)-1]
	//fmt.Printf("Pop: %v stack %v\n", ret, f.stack)
	return ret
}

// Length returns the stack length.
func (f *forthstack) Length() int {
	return len(f.stack)
}

// Empty is a convenience function synonymous with Length == 0
func (f *forthstack) Empty() bool {
	return len(f.stack) == 0
}

// errRecover converts panics to errstr iff it is an os.Error, panics
// otherwise.
func errRecover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		*errp = e.(error)
	}
}

/* iEval takes a Forth and strings and splits the string on space
 * characters, pushing each element on the stack or invoking the
 * operator if it is found in the opmap.
 */
func iEval(f Forth, s string) {
	// TODO: create a separator list based on isspace and all the
	// non alpha numberic characters in the opmap.
	for _, val := range strings.Fields(s) {
		//fmt.Printf("eval %s stack %v", val, f.Stack())
		fun := opmap[val]
		if fun != nil {
			//fmt.Printf("Eval ...:")
			fun(f)
			//fmt.Printf("Stack now %v", f.Stack())
		} else {
			f.Push(val)
			//fmt.Printf("push %s, stack %v", val, f.Stack())
		}
	}
	return
}

/* Eval takes a Forth and strings and splits the string on space
 * characters, pushing each element on the stack or invoking the
 * operator if it is found in the opmap. It returns TOS when it is done.
 * it is an error to leave the stack non-Empty.
 */
func Eval(f Forth, s string) (ret string, err error) {
	defer errRecover(&err)
	iEval(f, s)
	ret = f.Pop()
	return

}

// toInt converts to int64.
func toInt(f Forth) int64 {
	i, err := strconv.ParseInt(f.Pop(), 0, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func plus(f Forth) {
	x := toInt(f)
	y := toInt(f)
	z := x + y
	f.Push(strconv.FormatInt(z, 10))
}

func times(f Forth) {
	x := toInt(f)
	y := toInt(f)
	z := x * y
	f.Push(strconv.FormatInt(z, 10))
}

func sub(f Forth) {
	x := toInt(f)
	y := toInt(f)
	z := y - x
	f.Push(strconv.FormatInt(z, 10))
}

func div(f Forth) {
	x := toInt(f)
	y := toInt(f)
	z := y / x
	f.Push(strconv.FormatInt(z, 10))
}

func mod(f Forth) {
	x := toInt(f)
	y := toInt(f)
	z := y % x
	f.Push(strconv.FormatInt(z, 10))
}

func roundup(f Forth) {
	rnd := toInt(f)
	v := toInt(f)
	v = ((v + rnd - 1) / rnd) * rnd
	f.Push(strconv.FormatInt(v, 10))
}

func swap(f Forth) {
	x := f.Pop()
	y := f.Pop()
	f.Push(x)
	f.Push(y)
}

func strcat(f Forth) {
	x := f.Pop()
	y := f.Pop()
	f.Push(y + x)
}

func dup(f Forth) {
	x := f.Pop()
	f.Push(x)
	f.Push(x)
}

func ifelse(f Forth) {
	x := toInt(f)
	y := f.Pop()
	z := f.Pop()
	if x != 0 {
		f.Push(y)
	} else {
		f.Push(z)
	}
}

func hostname(f Forth) {
	h, err := os.Hostname()
	if err != nil {
		panic("No hostname")
	}
	f.Push(h)
}

func hostbase(f Forth) {
	host := f.Pop()
	f.Push(strings.TrimLeft(host, "abcdefghijklmnopqrstuvwxyz -"))
}

func NewWord(f Forth, name, command string) {
	newword := func(f Forth) {
		iEval(f, command)
		return
	}
	opmap[name] = newword
	return
}
