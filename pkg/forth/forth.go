// Copyright 2010-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package forth implements Forth parsing, which allows
// programs to use forth-like syntax to manipulate a stack
// of Cells.
// It is designed for use by programs
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
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type (
	// Op is an opcode type. It does not return an error value,
	// instead using panic when parsing issues occur, to ease
	// the programming annoyance of propagating errors up the
	// stack (following common Go practice for parsers).
	// If you write an op it can use panic as well.
	// Lest you get upset about the use of panic, be aware
	// I've talked to the folks in Go core about this and
	// they feel it's fine for parsers.
	Op func(f Forth)
	// Cell is a stack element.
	Cell interface{}

	stack struct {
		stack []Cell
	}
)

var (
	// Debug is an empty function that can be set to, e.g.,
	// fmt.Printf or log.Printf for debugging.
	Debug   = func(string, ...interface{}) {}
	opmap   map[string]Op
	mapLock sync.Mutex
)

func init() {
	opmap = map[string]Op{
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
		"drop":     drop,
		"newword":  newword,
		"words":    words,
	}
}

// Forth is an interface used by the package. The interface
// requires definition of Push, Pop, Length, Empty (convenience function
// meaning Length is 0), Newop (insert a new or replacement operator),
// and Reset (clear the stack, mainly diagnostic)
type Forth interface {
	Push(Cell)
	Pop() Cell
	Length() int
	Empty() bool
	Reset()
	Stack() []Cell
}

// New creates a new stack
func New() Forth {
	f := new(stack)
	return f
}

// Getop gets an op from the map.
func Getop(n string) Op {
	mapLock.Lock()
	defer mapLock.Unlock()
	op, ok := opmap[n]
	if !ok {
		return nil
	}
	return op
}

// Putop creates a new operation. We considered having
// an opmap per stack but don't feel the package requires it
func Putop(n string, op Op) {
	mapLock.Lock()
	defer mapLock.Unlock()
	if _, ok := opmap[n]; ok {
		panic("Putting %s: op already assigned")
	}
	opmap[n] = op
}

// Ops returns the operator map.
func Ops() map[string]Op {
	return opmap
}

// Reset resets the stack to empty
func (f *stack) Reset() {
	f.stack = f.stack[0:0]
}

// Return the stack
func (f *stack) Stack() []Cell {
	return f.stack
}

// Push pushes the interface{} on the stack.
func (f *stack) Push(c Cell) {
	f.stack = append(f.stack, c)
	Debug("push: %v: stack: %v\n", c, f.stack)
}

// Pop pops the stack. If the stack is Empty Pop will panic.
// Eval recovers() the panic.
func (f *stack) Pop() (ret Cell) {
	if len(f.stack) < 1 {
		panic(errors.New("Empty stack"))
	}
	ret = f.stack[len(f.stack)-1]
	f.stack = f.stack[0 : len(f.stack)-1]
	Debug("Pop: %v stack %v\n", ret, f.stack)
	return ret
}

// Length returns the stack length.
func (f *stack) Length() int {
	return len(f.stack)
}

// Empty is a convenience function synonymous with Length == 0
func (f *stack) Empty() bool {
	return len(f.stack) == 0
}

// errRecover converts panics to errstr iff it is an os.Error, panics
// otherwise.
func errRecover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			Debug("errRecover panics with a runtime error")
			panic(e)
		}
		Debug("errRecover returns %v", e)
		*errp = fmt.Errorf("%v", e)
	}
}

// Eval takes a Forth and []Cell, pushing each element on the stack or invoking the
// operator if it is found in the opmap.
func eval(f Forth, cells ...Cell) {
	Debug("eval cells %v", cells)
	for _, c := range cells {
		Debug("eval %v(%T) stack %v", c, c, f.Stack())
		switch s := c.(type) {
		case string:
			fun := Getop(s)
			if fun != nil {
				Debug("eval ... %v:", f.Stack())
				fun(f)
				Debug("eval: Stack now %v", f.Stack())
				break
			}
			if s[0] == '\'' {
				s = s[1:]
			}
			f.Push(s)
			Debug("push %v(%T), stack %v", s, s, f.Stack())
		default:
			Debug("push %v(%T), stack %v", s, s, f.Stack())
			f.Push(s)
		}
	}
}

// Eval calls eval and catches panics.
func Eval(f Forth, cells ...Cell) (err error) {
	defer errRecover(&err)
	eval(f, cells...)
	return
}

// EvalString takes a Forth and string and splits the string on space
// characters, calling Eval for each one.
func EvalString(f Forth, s string) (err error) {
	for _, c := range strings.Fields(s) {
		if err = Eval(f, c); err != nil {
			return
		}
	}
	Debug("EvalString err %v", err)
	return
}

// EvalPop takes a Forth and string, calls EvalString, and
// returns TOS and an error, if any.
// For EvalPop it is an error to leave the stack non-Empty.
// EvalPop is typically used for programs that want to
// parse forth contained in, e.g., flag.Args(), and return
// a result. In most use cases, we want the stack to be empty.
func EvalPop(f Forth, s string) (ret Cell, err error) {
	defer errRecover(&err)
	if err = EvalString(f, s); err != nil {
		return
	}
	if f.Length() != 1 {
		panic(fmt.Sprintf("%v: length is not 1", f.Stack()))
	}
	ret = f.Pop()
	Debug("EvalPop ret %v err %v", ret, err)
	return
}

// String returns the Top Of Stack if it is a string
// or panics.
func String(f Forth) string {
	c := f.Pop()
	switch s := c.(type) {
	case string:
		return s
	default:
		panic(fmt.Sprintf("Can't convert %v to a string", c))
	}
}

// toInt converts to int64.
func toInt(f Forth) int64 {
	Debug("toint %v", f.Stack())
	c := f.Pop()
	Debug("%T", c)
	switch s := c.(type) {
	case string:
		i, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			panic(err)
		}
		return i
	case int64:
		return s
	default:
		panic(fmt.Sprintf("NaN: %T", c))
	}
}

func plus(f Forth) {
	x := toInt(f)
	y := toInt(f)
	z := x + y
	f.Push(strconv.FormatInt(z, 10))
}

func words(f Forth) {
	mapLock.Lock()
	defer mapLock.Unlock()
	var w []string
	for i := range opmap {
		w = append(w, i)
	}
	f.Push(w)
}

func newword(f Forth) {
	s := String(f)
	n := toInt(f)
	// Pop <n> Cells.
	if f.Length() < int(n) {
		panic(fmt.Sprintf("newword %s: stack is %d elements, need %d", s, f.Length(), n))
	}
	var c = make([]Cell, n)
	for i := n; i > 0; i-- {
		c[i-1] = f.Pop()
	}
	Putop(s, func(f Forth) {
		Debug("c %v", c)
		eval(f, c...)
	})
}

func drop(f Forth) {
	_ = f.Pop()
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
	x := String(f)
	y := String(f)
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
	host := String(f)
	f.Push(strings.TrimLeft(host, "abcdefghijklmnopqrstuvwxyz -"))
}

// NewWord allows for definition of new operators from strings.
// For example, should you wish to create a word which adds a number
// to itself twice, you can call:
// NewWord(f, "d3d", "dup dup + +")
// which does two dups, and two adds.
func NewWord(f Forth, name string, cell Cell, cells ...Cell) {
	cmd := append([]Cell{cell}, cells...)
	newword := func(f Forth) {
		eval(f, cmd...)
	}
	Putop(name, newword)
}
