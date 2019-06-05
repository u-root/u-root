package eval

// Builtin functions.

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
	"unicode/utf8"

	"github.com/u-root/u-root/cmds/core/elvish/eval/vals"
	"github.com/u-root/u-root/cmds/core/elvish/parse"
)

// Builtins that have not been put into their own groups go here.

func init() {
	addBuiltinFns(map[string]interface{}{
		"-source": source,
	})
}

func nop(opts RawOptions, args ...interface{}) {
	// Do nothing
}

func kindOf(fm *Frame, args ...interface{}) {
	out := fm.ports[1].Chan
	for _, a := range args {
		out <- vals.Kind(a)
	}
}

func constantly(args ...interface{}) Callable {
	// XXX Repr of this fn is not right
	return NewBuiltinFn(
		"created by constantly",
		func(fm *Frame) {
			out := fm.ports[1].Chan
			for _, v := range args {
				out <- v
			}
		},
	)
}

func resolve(fm *Frame, head string) string {
	// Emulate static resolution of a command head. This needs to be kept in
	// sync with (*compiler).form.

	_, special := builtinSpecials[head]
	if special {
		return "special"
	} else {
		explode, ns, name := ParseVariableRef(head)
		if !explode && fm.ResolveVar(ns, name+FnSuffix) != nil {
			return "$" + head + FnSuffix
		} else {
			return "(external " + parse.Quote(head) + ")"
		}
	}
}

func source(fm *Frame, fname string) error {
	path, err := filepath.Abs(fname)
	if err != nil {
		return err
	}
	code, err := readFileUTF8(path)
	if err != nil {
		return err
	}
	n, err := parse.Parse(fname, code)
	if err != nil {
		return err
	}
	scriptGlobal := fm.local.static()
	for name := range fm.up.static() {
		scriptGlobal.set(name)
	}
	op, err := compile(fm.Builtin.static(),
		scriptGlobal, n, NewScriptSource(fname, path, code))
	if err != nil {
		return err
	}
	return fm.Eval(op)
}

func readFileUTF8(fname string) (string, error) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(bytes) {
		return "", fmt.Errorf("%s: source is not valid UTF-8", fname)
	}
	return string(bytes), nil
}

func sleep(fm *Frame, t float64) error {
	d := time.Duration(float64(time.Second) * t)
	select {
	case <-fm.Interrupts():
		return ErrInterrupted
	case <-time.After(d):
		return nil
	}
}

func _time(fm *Frame, f Callable) error {
	t0 := time.Now()
	err := f.Call(fm, NoArgs, NoOpts)
	t1 := time.Now()

	dt := t1.Sub(t0)
	fmt.Fprintln(fm.ports[1].File, dt)

	return err
}
