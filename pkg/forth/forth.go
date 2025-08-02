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
// Get the hostbase, if it is 0 % 20, return the hostbase / 20,
// else return hostbase % 20
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
	"math/big"
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
	Cell any

	stack struct {
		stack []Cell
	}
)

var (
	// Debug is an empty function that can be set to, e.g.,
	// fmt.Printf or log.Printf for debugging.
	Debug   = func(string, ...any) {}
	opmap   map[string]Op
	mapLock sync.Mutex
	// EmptyStack means we wanted something and ... nothing there.
	ErrEmptyStack = errors.New("empty stack")
	// NotEnoughElements means the stack is not deep enough for whatever operator we have.
	ErrNotEnoughElements = errors.New("not enough elements on stack")
	// ErrWordExist is the error for trying to create a word that's already in use.
	ErrWordExist = errors.New("word already exists")
)

func init() {
	opmap = map[string]Op{
		"+":        plus,
		"-":        sub,
		"*":        times,
		"/":        div,
		"%":        rem,
		"mod":      mod,
		"printf":   fmtprintf,
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
		"typeof":   typeOf,
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
		panic(fmt.Errorf("putting %s: %w", n, ErrWordExist))
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
		panic(ErrEmptyStack)
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
	Debug("pkg/forth:e is %v:%T", e, e)
	if e != nil {
		switch err := e.(type) {
		case runtime.Error:
			Debug("pkg/forth:errRecover panics with a runtime error")
			panic(e)
		case error:
			*errp = err
		default:
			*errp = fmt.Errorf("pkg/forth:%v", err)
		}
		Debug("pkg/forth:errRecover returns %v:%T", *errp, *errp)
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
	for c := range strings.FieldsSeq(s) {
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
		panic(fmt.Errorf("%v: length is not 1;%w", f.Stack(), strconv.ErrRange))
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
		panic(fmt.Errorf("%v:%w", c, strconv.ErrSyntax))
	}
}

// typeOf pops the stack, and replaces it with the
// type as a string.
func typeOf(f Forth) {
	Debug("toint %v", f.Stack())
	c := f.Pop()
	Debug("%T", c)
	f.Push(fmt.Sprintf("%T", c))
}

// toRat converts TOS to big.Rat.
func toRat(f Forth) *big.Rat {
	Debug("toint %v", f.Stack())
	c := f.Pop()
	Debug("%T", c)
	r := new(big.Rat)
	var num string
	switch s := c.(type) {
	case string:
		num = s
	case *big.Rat:
		return s
	default:
		num = fmt.Sprintf("%d", s)
	}
	if _, err := fmt.Sscan(num, r); err != nil {
		// Older go versions of go big don't wrap the error
		// So we must do it ourselves.
		panic(strconv.ErrSyntax)
	}
	return r
}

// toInt converts TOS to a big.Int
func toInt(f Forth) *big.Int {
	r := toRat(f)
	if !r.IsInt() {
		panic(fmt.Errorf("%v: not an int:%w", r.String(), strconv.ErrSyntax))
	}
	return r.Num()
}

func plus(f Forth) {
	x := toRat(f)
	y := toRat(f)
	x.Add(x, y)
	f.Push(x)
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
	n := toInt(f).Int64()
	// Pop <n> Cells.
	if int64(f.Length()) < n {
		panic(fmt.Errorf("newword %s: stack is %d elements, need %d:%w", s, f.Length(), n, ErrNotEnoughElements))
	}
	c := make([]Cell, n)
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
	x := toRat(f)
	y := toRat(f)
	x.Mul(x, y)
	f.Push(x)
}

func sub(f Forth) {
	x := toRat(f)
	y := toRat(f)
	x.Sub(x, y)
	f.Push(x)
}

func div(f Forth) {
	x := toRat(f)
	y := toRat(f)
	x.Quo(x, y)
	f.Push(x)
}

func mod(f Forth) {
	x := toInt(f)
	y := toInt(f)
	x.Mod(x, y)
	f.Push((&big.Rat{}).SetInt(x))
}

func rem(f Forth) {
	x := toInt(f)
	y := toInt(f)
	x.Rem(x, y)
	f.Push((&big.Rat{}).SetInt(x))
}

func roundup(f Forth) {
	rnd := toRat(f)
	v := toRat(f)
	v = v.Add(v, rnd)
	v = v.Sub(v, big.NewRat(1, 1))
	v = v.Quo(v, rnd)
	v = v.Mul(v, rnd)
	f.Push(v)
}

func fmtprintf(f Forth) {
	x := f.Pop()
	s := x.(string)
	y := f.Pop()
	f.Push(fmt.Sprintf(s, y))
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
	if x.Cmp(big.NewInt(0)) == 0 {
		f.Push(y)
	} else {
		f.Push(z)
	}
}

func hostname(f Forth) {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
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
