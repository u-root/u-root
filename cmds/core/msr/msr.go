// Copyright 2010-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// msr reads and writes msrs using a Forth interpreter on argv
//
// Synopsis:
//     To see what is available:
//     msr words
//
// Description:
//     msr provides a set of Forth words that let you manage MSRs.
//     You can add new ones of your own.
//     For a start, it provides some pre-defined words for well-known MSRs
//
//     push a [] of MSR names and the 0x3a register on the stack
//     IA32_FEATURE_CONTROL -- equivalent to * msr 0x3a reg
//     The next two commands use IA32_FEATURE_CONTROL:
//     READ_IA32_FEATURE_CONTROL -- equivalent to IA32_FEATURE_CONTROL rd
//     LOCK IA32_FEATURE_CONTROL -- equivalent to IA32_FEATURE_CONTROL rd IA32_FEATURE_CONTROL 1 val or wr
//     e.g.
//./msr IA32_FEATURE_CONTROL
// [[/dev/cpu/0/msr /dev/cpu/1/msr /dev/cpu/2/msr /dev/cpu/3/msr] 58]
//
//     As a special convenience, we have two useful cases:
//     r glob register -- read the MSR 'register' from cores matching 'glob'
//     w glob register value -- write the value to 'register' on all cores matching 'glob'
//
// Examples:
//     Show the IA32 feature MSR on all cores
//     sudo fio READ_IA32_FEATURE_CONTROL
//     [[5 5 5 5]]
//     lock the registers
//     sudo fio LOCK_IA32_FEATURE_CONTROL
//     Just see it one core 0 and 1
//     sudo ./fio '[01]' msr 0x3a reg rd
//     [[5 5]]
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/u-root/u-root/pkg/forth"
)

// let's just do MSRs for now

var (
	debug = flag.Bool("d", false, "debug messages")
	words = []struct {
		name string
		w    []forth.Cell
	}{
		{name: "IA32_FEATURE_CONTROL", w: []forth.Cell{"'*", "msr", "0x3a", "reg"}},
		{name: "READ_IA32_FEATURE_CONTROL", w: []forth.Cell{"IA32_FEATURE_CONTROL", "rd"}},
		{name: "LOCK_IA32_FEATURE_CONTROL", w: []forth.Cell{"IA32_FEATURE_CONTROL", "READ_IA32_FEATURE_CONTROL", "1", "val", "or", "wr"}},
	}
	ops = []struct {
		name string
		op   forth.Op
	}{
		{name: "msr", op: msr},
		{name: "reg", op: reg},
		{name: "val", op: val},
		{name: "rd", op: rd},
		{name: "wr", op: wr},
		{name: "swr", op: swr},
		{name: "and", op: and},
		{name: "or", op: or},
	}
)

// The great panic discussion.
// Rob has given talks on using panic for parsers.
// I have talked to Russ about using panic for parsers.
// Short form: it's ok. In general, don't panic.
// But parsers are special: using panic
// in a parser makes the code tons cleaner.

// Note that if any type asserts fail the forth interpret loop catches
// it. It also catches stack underflow, all that stuff.
func msr(f forth.Forth) {
	forth.Debug("msr")
	g := f.Pop().(string)
	m := msrList(g)
	forth.Debug("MSRs are %v", m)
	f.Push(m)
}

func reg(f forth.Forth) {
	n, err := strconv.ParseUint(f.Pop().(string), 0, 32)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	f.Push(uint32(n))
}

func val(f forth.Forth) {
	n, err := strconv.ParseUint(f.Pop().(string), 0, 64)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	f.Push(uint64(n))
}

func rd(f forth.Forth) {
	r := f.Pop().(uint32)
	m := f.Pop().([]string)
	forth.Debug("rd: %v %v", m, r)
	data, errs := rdmsr(m, r)
	forth.Debug("data %v errs %v", data, errs)
	if errs != nil {
		panic(fmt.Sprintf("%v", errs))
	}
	f.Push(data)
}

// the standard msr api takes one value for all msrs.
// That's arguably substandard. We're going to require
// a []uint64. It's naive to expect every core to have
// exactly the same msr values for each msr in our
// modern world.
// If you're determined to write a fixed value, the same
// for all, it's easy:
// fio "'"* msr 0x3a reg rd 0 val and your-value new-val val or wr
// Then you'll have a fixed value.
func wr(f forth.Forth) {
	v := f.Pop().([]uint64)
	r := f.Pop().(uint32)
	m := f.Pop().([]string)
	forth.Debug("wr: %v %v %v", m, r, v)
	errs := wrmsr(m, r, v)
	forth.Debug("errs %v", errs)
	if errs != nil {
		f.Push(errs)
	}
}

// We had been counting on doing a rd, which would produce a nice
// []u64 at TOS which we could use in a write. Problem is, some MSRs
// can not be read. There are write-only MSRs.  This really
// complicates the picture: we can't just read them, change them, and
// write them; we would not even know if reading is side-effect free.
//
// Our beautiful scheme was murdered by a brutal gang of facts.
//
// To keep it simple, we have swr, for scalar write, which turns the
// scalar value into a []u64.  This function pops the args, and turns
// the scalar at TOS into a u64 for wrmsr.
func swr(f forth.Forth) {
	v := f.Pop().(uint64)
	r := f.Pop().(uint32)
	m := f.Pop().([]string)
	var vv = make([]uint64, len(m))
	for i := range vv {
		vv[i] = v
	}
	forth.Debug("swr: %v %v %v", m, r, vv)
	errs := wrmsr(m, r, vv)
	forth.Debug("errs %v", errs)
	if errs != nil {
		f.Push(errs)
	}
}

func and(f forth.Forth) {
	v := f.Pop().(uint64)
	m := f.Pop().([]uint64)
	forth.Debug("and: %v(%T) %v(%T)", m, m, v, v)
	for i := range m {
		m[i] &= v
	}
	f.Push(m)
}

func or(f forth.Forth) {
	v := f.Pop().(uint64)
	m := f.Pop().([]uint64)
	forth.Debug("or: %v(%T) %v(%T)", m, m, v, v)
	for i := range m {
		m[i] |= v
	}
	f.Push(m)
}

func main() {
	flag.Parse()
	if *debug {
		forth.Debug = log.Printf
	}

	// TODO: probably can do this by just having two passes, and write
	// in the first pass is a no op. Which will fail to catch the problem
	// of read-only and write-only MSRs but there's only so much you can do.
	//
	// To avoid the command list from being partially executed when the
	// args fail to parse, queue them up and run all at once at the end.
	//queue := []func(){} etc. etc.
	f := forth.New()
	for _, o := range ops {
		forth.Putop(o.name, o.op)
	}
	for _, w := range words {
		forth.NewWord(f, w.name, w.w[0], w.w[1:]...)
	}
	a := flag.Args()
	// If the first arg is r or w, we're going to assume they're not doing Forth.
	// It is too confusing otherwise if they type a wrong r or w command and
	// see the Forth stack and nothing else.
	switch a[0] {
	case "r":
		if len(a) != 3 {
			log.Fatal("Usage for r: r <msr-glob> <register>")
		}
		// Because the msr arg is a glob and may have things like * in it (* being the
		// most common) gratuitiously add a Forth ' before it (i.e. quote it).
		if err := forth.EvalString(f, fmt.Sprintf("'%s msr %s reg rd", a[1], a[2])); err != nil {
			log.Fatal(err)
		}
	case "w":
		if len(a) != 4 {
			log.Fatal("Usage for w: w <msr-glob> <register> <value>")
		}
		// Because the msr arg is a glob and may have things like * in it (* being the
		// most common) gratuitiously add a Forth ' before it (i.e. quote it).
		if err := forth.EvalString(f, fmt.Sprintf("'%s msr %s reg %s val swr", a[1], a[2], a[3])); err != nil {
			log.Fatal(err)
		}
	default:
		for _, a := range flag.Args() {
			if err := forth.EvalString(f, a); err != nil {
				log.Fatal(err)
			}
			forth.Debug("%vOK\n", f.Stack())
		}
	}
	// special case: if the length of stack is 1, just print out stack[0].
	s := f.Stack()
	if len(s) == 1 {
		fmt.Printf("%v\n", s[0])
	} else {
		fmt.Printf("%v\n", s)
	}
}
