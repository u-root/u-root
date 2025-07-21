// Copyright 2010-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9

// msr reads and writes msrs using a Forth interpreter on argv
//
// Synopsis:
//
//	To see what is available:
//	msr words
//
// Description:
//
//	msr provides a set of Forth words that let you manage MSRs.
//	You can add new ones of your own.
//	If you are going to use forth, the general pattern of arguments
//	looks something like this:
//	msr <glob pattern> cpu <msr-number> msr <opcode>
//	e.g.,
//	msr '*' 0x3a rd
//	Will read the 3a msr on all CPUs.
//	The two first items will remain at TOS. They are implicitly converted to
//	msr.CPUs and msr.MSR when used. You can show how the will be converted with
//	the cpu and reg works.
//	You can build up the expressions bit by bit:
//	For a start, it provides some pre-defined words for well-known MSRs
//
//	push a [] of MSR names and the 0x3a register on the stack
//	IA32_FEATURE_CONTROL -- equivalent to * cpu 0x3a msr
//	The next two commands use IA32_FEATURE_CONTROL:
//	READ_IA32_FEATURE_CONTROL -- equivalent to IA32_FEATURE_CONTROL rd
//	LOCK IA32_FEATURE_CONTROL -- equivalent to IA32_FEATURE_CONTROL rd IA32_FEATURE_CONTROL 1 u64 or wr
//	e.g.
//
// ./msr IA32_FEATURE_CONTROL
// [[/dev/cpu/0/msr /dev/cpu/1/msr /dev/cpu/2/msr /dev/cpu/3/msr] 58]
//
//	As a special convenience, we have two useful cases:
//	r glob register -- read the MSR 'register' from cores matching 'glob'
//	w glob register value -- write the value to 'register' on all cores matching 'glob'
//
// Examples (NOTE: single ', since it is a forth literal! '* NOT '*'):
//
//	Show the IA32 feature MSR on all cores
//	sudo msr READ_IA32_FEATURE_CONTROL
//	[[5 5 5 5]]
//	lock the registers
//	sudo msr LOCK_IA32_FEATURE_CONTROL
//	Just see it one core 0 and 1
//	sudo ./msr '[01]' 0x3a rd
//	[[5 5]]
//	Debug your cpu mask
//	sudo ./msr "'1*" "cpu"
//	1,10-19
//	Debug your command stack
//	sudo ./msr "'[01]" cpu 0x3a reg
//	[0-1 0x3a]
//
//	For rd, if you want to see it in hex (which should be the default
//	but it's complicated
//	sudo ./msr "'[01]" cpu 0x10 reg rd %#x printf
//	[0xf41c40e682bb8 0xf41c40e69876b]
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/u-root/u-root/pkg/forth"
	"github.com/u-root/u-root/pkg/msr"
)

// let's just do MSRs for now
// NOTE: all these defines now take advantage of the fact that wr and commands
// leave the msr.CPUs and msr.MSR on the stack.

var (
	debug = flag.Bool("d", false, "debug messages")
	words = []struct {
		name string
		w    []forth.Cell
	}{
		// Architectural MSR. All systems.
		// Enables features like VMX.
		{name: "MSR_IA32_FEATURE_CONTROL", w: []forth.Cell{"'*", "cpu", "0x3a", "reg"}},
		{name: "READ_MSR_IA32_FEATURE_CONTROL", w: []forth.Cell{"MSR_IA32_FEATURE_CONTROL", "rd"}},
		{name: "LOCK_MSR_IA32_FEATURE_CONTROL", w: []forth.Cell{"MSR_IA32_FEATURE_CONTROL", "READ_MSR_IA32_FEATURE_CONTROL", "1", "or", "wr"}},
		// PM ENABLE
		{name: "MSR_IA32_PM_ENABLE", w: []forth.Cell{"'*", "0x770"}},
		{name: "READ_MSR_IA32_PM_ENABLE", w: []forth.Cell{"MSR_IA32_PM_ENABLE", "rd"}},
		// Silvermont, Airmont, Nehalem...
		// Controls Processor C States.
		{name: "MSR_PKG_CST_CONFIG_CONTROL", w: []forth.Cell{"'*", "cpu", "0xe2", "reg"}},
		{name: "READ_MSR_PKG_CST_CONFIG_CONTROL", w: []forth.Cell{"MSR_PKG_CST_CONFIG_CONTROL", "rd"}},
		{name: "LOCK_MSR_PKG_CST_CONFIG_CONTROL", w: []forth.Cell{"READ_MSR_PKG_CST_CONFIG_CONTROL", uint64(1 << 15), "or", "wr"}},
		// Westmere onwards.
		// Note that this turns on AES instructions, however
		// 3 will turn off AES until reset.
		{name: "MSR_FEATURE_CONFIG", w: []forth.Cell{"'*", "cpu", "0x13c", "reg"}},
		{name: "READ_MSR_FEATURE_CONFIG", w: []forth.Cell{"MSR_FEATURE_CONFIG", "rd"}},
		{name: "LOCK_MSR_FEATURE_CONFIG", w: []forth.Cell{"READ_MSR_FEATURE_CONFIG", uint64(1 << 0), "or", "wr"}},
		// Goldmont, SandyBridge
		// Controls DRAM power limits. See Intel SDM
		{name: "MSR_DRAM_POWER_LIMIT", w: []forth.Cell{"'*", "cpu", "0x618", "reg"}},
		{name: "READ_MSR_DRAM_POWER_LIMIT", w: []forth.Cell{"MSR_DRAM_POWER_LIMIT", "rd"}},
		{name: "LOCK_MSR_DRAM_POWER_LIMIT", w: []forth.Cell{"READ_MSR_DRAM_POWER_LIMIT", uint64(1 << 31), "or", "wr"}},
		// IvyBridge Onwards.
		// Not much information in the SDM, seems to control power limits
		{name: "MSR_CONFIG_TDP_CONTROL", w: []forth.Cell{"'*", "cpu", "0xe2", "reg"}},
		{name: "READ_MSR_CONFIG_TDP_CONTROL", w: []forth.Cell{"MSR_CONFIG_TDP_CONTROL", "rd"}},
		{name: "LOCK_MSR_CONFIG_TDP_CONTROL", w: []forth.Cell{"READ_MSR_CONFIG_TDP_CONTROL", uint64(1 << 31), "or", "wr"}},
		// Architectural MSR. All systems.
		// This is the actual spelling of the MSR in the manual.
		// Controls availability of silicon debug interfaces
		{name: "IA32_DEBUG_INTERFACE", w: []forth.Cell{"'*", "cpu", "0xe2", "reg"}},
		{name: "READ_IA32_DEBUG_INTERFACE", w: []forth.Cell{"IA32_DEBUG_INTERFACE", "rd"}},
		{name: "LOCK_IA32_DEBUG_INTERFACE", w: []forth.Cell{"READ_IA32_DEBUG_INTERFACE", uint64(1 << 15), "or", "wr"}},
		// Locks all known msrs to lock
		{name: "LOCK_KNOWN_MSRS", w: []forth.Cell{"LOCK_MSR_IA32_FEATURE_CONTROL", "LOCK_MSR_PKG_CST_CONFIG_CONTROL", "LOCK_MSR_FEATURE_CONFIG", "LOCK_MSR_DRAM_POWER_LIMIT", "LOCK_MSR_CONFIG_TDP_CONTROL", "LOCK_IA32_DEBUG_INTERFACE"}},
	}
	ops = []struct {
		name string
		op   forth.Op
	}{
		{name: "cpu", op: evalCPUs},
		{name: "reg", op: evalMSR},
		{name: "u64", op: u64},
		{name: "u64slice", op: u64slice},
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

func parseCPUs(s string) msr.CPUs {
	c, errs := msr.GlobCPUs(s)
	if errs != nil {
		panic(fmt.Sprintf("%q:%v", s, errs))
	}
	return c
}

// Note that if any type asserts fail the forth interpret loop catches
// it. It also catches stack underflow, all that stuff.
func evalCPUs(f forth.Forth) {
	forth.Debug("cpu")
	r := f.Pop()
	var c msr.CPUs
	switch v := r.(type) {
	case msr.CPUs:
		c = v
	case string:
		c = parseCPUs(v)
	default:
		panic(fmt.Sprintf("%v(%T): can not convert to msr.MSR", r, f))
	}

	forth.Debug("CPUs are %v", c)
	f.Push(c)
}

func CPUs(f forth.Forth) msr.CPUs {
	evalCPUs(f)
	return f.Pop().(msr.CPUs)
}

func parseMSR(s string) msr.MSR {
	n, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	return msr.MSR(n)
}

func evalMSR(f forth.Forth) {
	r := f.Pop()
	var m msr.MSR
	switch v := r.(type) {
	case msr.MSR:
		m = v
	case string:
		m = parseMSR(v)
	default:
		panic(fmt.Sprintf("%v(%T): can not convert to msr.MSR", r, f))
	}
	f.Push(m)
}

func MSR(f forth.Forth) msr.MSR {
	evalMSR(f)
	return f.Pop().(msr.MSR)
}

func tou64slice(a any) []uint64 {
	forth.Debug("tou64slice: %v:%T", a, a)
	switch v := a.(type) {
	case string:
		n, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("%q:%T:%v", v, v, err))
		}
		return []uint64{uint64(n)}
	case []string:
		u := make([]uint64, len(v))
		for i, s := range v {
			n, err := strconv.ParseUint(s, 0, 64)
			if err != nil {
				panic(fmt.Sprintf("%q:%T:%v", s, s, err))
			}
			u[i] = n
		}
	case uint64:
		// no idea how long it should be, so ...
		return []uint64{uint64(v)}
	case uint32:
		return []uint64{uint64(v)}[:]
	case uint16:
		return []uint64{uint64(v)}[:]
	case uint8:
		return []uint64{uint64(v)}[:]
	case []uint64:
		return v
	default:
		panic(fmt.Sprintf("can not convert %v:%T to uint64", a, a))
	}
	return nil
}

func u64slice(f forth.Forth) {
	f.Push(tou64slice(f.Pop()))
}

func tou64(a any) uint64 {
	var u uint64
	switch v := a.(type) {
	case string:
		n, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		u = n
	case uint64:
		u = v
	case uint32:
		u = uint64(v)
	case uint16:
		u = uint64(v)
	case uint8:
		u = uint64(v)
	}
	return u
}

func u64(f forth.Forth) {
	f.Push(tou64(f.Pop()))
}

func cpumsr(f forth.Forth) (msr.CPUs, msr.MSR) {
	m := MSR(f)
	c := CPUs(f)
	// It is proving to be much more convenient to always leave these
	// at TOS. This allows sequences like
	// msr 0 0x3a rd 1 and wr
	// without having to repeat the cpu and msr all the time.
	f.Push(c)
	f.Push(m)
	forth.Debug("cpumsr: cpu %v msr %v", c, m)
	return c, m
}

func cpumsrval(f forth.Forth) (msr.CPUs, msr.MSR, uint64) {
	u := tou64(f.Pop())
	c, m := cpumsr(f)
	return c, m, u
}

func rd(f forth.Forth) {
	c, m := cpumsr(f)
	forth.Debug("rd: cpus %v, msr %v", c, m)
	data, errs := m.Read(c)
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
// msr "'"* 0x3a rd 0 and new-val val or val wr
// Then you'll have a fixed value.
func wr(f forth.Forth) {
	v := tou64slice(f.Pop())
	c, r := cpumsr(f)
	forth.Debug("wr: cpus %v, msr %v, values %v", c, r, v)
	errs := r.Write(c, v...)
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
// Write now accepts a single value
func swr(f forth.Forth) {
	c, r, v := cpumsrval(f)

	forth.Debug("swr: cpus %v, msr %v, %v", c, r, v)
	errs := r.Write(c, v)
	forth.Debug("errs %v", errs)
	if errs != nil {
		f.Push(errs)
	}
}

// Not needed after go1.23
func clone(u []uint64) []uint64 {
	n := make([]uint64, len(u))
	copy(n, u)
	return n
}

func and(f forth.Forth) {
	v := tou64(f.Pop())
	m := tou64slice(f.Pop())
	m = clone(m)
	forth.Debug("and: %v(%T) %v(%T)", m, m, v, v)
	for i := range m {
		m[i] &= v
	}
	forth.Debug("Result:%v:%T", m, m)
	f.Push(m)
}

func or(f forth.Forth) {
	v := tou64(f.Pop())
	m := tou64slice(f.Pop())
	m = clone(m)
	forth.Debug("or: %v(%T) %v(%T)", m, m, v, v)
	for i := range m {
		m[i] |= v
	}
	forth.Debug("Result:%v:%T", m, m)
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
	if len(a) == 0 {
		flag.Usage()
		return
	}
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
		if err := forth.EvalString(f, fmt.Sprintf("'%s %s rd", a[1], a[2])); err != nil {
			log.Fatal(err)
		}
	case "w":
		if len(a) != 4 {
			log.Fatal("Usage for w: w <msr-glob> <register> <value>")
		}
		// Because the msr arg is a glob and may have things like * in it (* being the
		// most common) gratuitiously add a Forth ' before it (i.e. quote it).
		if err := forth.EvalString(f, fmt.Sprintf("'%s %s %s swr", a[1], a[2], a[3])); err != nil {
			log.Fatal(err)
		}
	case "lock":
		if len(a) != 4 {
			log.Fatal("Usage for lock: lock <msr-glob> <register> <bit>")
		}
		if err := forth.EvalString(f, fmt.Sprintf("'%s %s rd %s or wr", a[1], a[2], a[3])); err != nil {
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
	// special case: if the length of stack is 3, just print out stack[2].
	// The reason being that TOS is always going to be cpus, msr, and value.
	// We tried just printing out the whole stack but it's annoying.
	// If you really want to see the whole stack you can force the issue
	// by adding an extraneous word that makes the end result not 3, e.g.
	// msr TOS 0 0x3a rd
	// will show you all the stack.
	// $ msr TOS [10] 0x3a rd
	// [TOS 0-1 0x3a [5 5]]
	s := f.Stack()
	if len(s) == 3 {
		fmt.Printf("%v\n", s[2])
	} else {
		fmt.Printf("%v\n", s)
	}
}
