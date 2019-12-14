package eval

import (
	"errors"
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/eval/vals"
)

func TestReflectBuiltinFnCall(t *testing.T) {
	called := false
	theFrame := new(Frame)
	theOptions := map[string]interface{}{}
	theError := errors.New("the error")

	var f Callable
	callGood := func(fm *Frame, args []interface{}, opts map[string]interface{}) {
		err := f.Call(fm, args, opts)
		if err != nil {
			t.Errorf("Failed to call f: %v", err)
		}
	}
	callBad := func(fm *Frame, args []interface{}, opts map[string]interface{}) {
		err := f.Call(fm, args, opts)
		if err == nil {
			t.Errorf("Calling f didn't return error")
		}
	}
	callBad2 := func(fm *Frame, args []interface{}, opts map[string]interface{}) {
		err := f.Call(fm, args, opts)
		if err == nil {
			t.Errorf("Calling f didn't return error")
		}
		if err != theError {
			t.Errorf("Calling f didn't return the right error")
		}
	}

	// func()
	called = false
	f = NewBuiltinFn("f1", func() {
		called = true
	})
	callGood(theFrame, nil, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(float64)
	called = false
	f = NewBuiltinFn("f2", func(a float64) {
		called = true
		if a != 123.456 {
			t.Errorf("arg1 not passed")
		}
	})
	callGood(theFrame, []interface{}{"123.456"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(int)
	called = false
	f = NewBuiltinFn("f3", func(a int) {
		called = true
		if a != 123 {
			t.Errorf("arg1 not passed")
		}
	})
	callGood(theFrame, []interface{}{"123"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(...int) error
	f = NewBuiltinFn("f4", func(va ...int) error {
		if va[0] != 123 {
			t.Errorf("arg1 not passed")
		}
		if va[1] != 456 {
			t.Errorf("arg2 not passed")
		}
		return theError
	})
	callBad2(theFrame, []interface{}{"123", "456"}, theOptions)

	// func() error
	f = NewBuiltinFn("f5", func() error {
		called = true
		return theError
	})
	callBad2(theFrame, nil, theOptions)

	// func(*eval.Frame)
	called = false
	f = NewBuiltinFn("f10", func(f *Frame) {
		called = true
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
	})
	callGood(theFrame, nil, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(*eval.Frame, ...interface {}) error
	f = NewBuiltinFn("f12", func(f *Frame, va ...interface{}) error {
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if va[0].(int) != 123 {
			t.Errorf("arg1 not passed")
		}
		if va[1].(string) != "abc" {
			t.Errorf("arg2 not passed")
		}
		if va[2].(error) != theError {
			t.Errorf("arg3 not passed")
		}
		return theError
	})
	callBad2(theFrame, []interface{}{123, "abc", theError}, theOptions)

	// func(*eval.Frame, interface {}, interface {}, interface {})
	called = false
	f = NewBuiltinFn("f12", func(f *Frame, a1 interface{}, a2 interface{}, a3 interface{}) {
		called = true
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if a1.(int) != 123 {
			t.Errorf("arg1 not passed")
		}
		if a2.(string) != "abc" {
			t.Errorf("arg2 not passed")
		}
		if a3.(error) != theError {
			t.Errorf("arg3 not passed")
		}
	})
	callGood(theFrame, []interface{}{123, "abc", theError}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(*eval.Frame, ...int) error
	f = NewBuiltinFn("f13", func(f *Frame, va ...int) error {
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if va[0] != 123 {
			t.Errorf("arg1 not passed")
		}
		return theError
	})
	callBad2(theFrame, []interface{}{"123"}, theOptions)

	// func(*eval.Frame, ...string) error
	f = NewBuiltinFn("f13", func(f *Frame, va ...string) error {
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if va[0] != "abc" {
			t.Errorf("arg1 not passed")
		}
		return theError
	})
	callBad2(theFrame, []interface{}{"abc"}, theOptions)

	// func(*eval.Frame, string)
	called = false
	f = NewBuiltinFn("f14", func(f *Frame, s string) {
		called = true
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if s != "abc" {
			t.Errorf("arg1 not passed")
		}
	})
	callGood(theFrame, []interface{}{"abc"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(*eval.Frame, string) error
	f = NewBuiltinFn("f15", func(f *Frame, s string) error {
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if s != "abc" {
			t.Errorf("arg1 not passed")
		}
		return theError
	})
	callBad2(theFrame, []interface{}{"abc"}, theOptions)

	// Options parameter gets options.
	called = false
	f = NewBuiltinFn("f20", func(opts RawOptions) {
		called = true
		if opts["foo"] != "bar" {
			t.Errorf("Options parameter doesn't get options")
		}
	})
	callGood(theFrame, nil, RawOptions{"foo": "bar"})
	if !called {
		t.Errorf("not called")
	}

	// Combination of Frame and Options.
	called = false
	f = NewBuiltinFn("f30", func(f *Frame, opts RawOptions) {
		called = true
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if opts["foo"] != "bar" {
			t.Errorf("Options parameter doesn't get options")
		}
	})
	callGood(theFrame, nil, RawOptions{"foo": "bar"})
	if !called {
		t.Errorf("not called")
	}

	// func(*eval.Frame, eval.RawOptions, eval.Callable, eval.Callable)
	called = false
	theCallable1 := &BuiltinFn{}
	theCallable2 := &BuiltinFn{}
	f = NewBuiltinFn("f35", func(f *Frame, opts RawOptions, a1 Callable, a2 Callable) {
		called = true
		if f != theFrame {
			t.Errorf("*Frame parameter doesn't get current frame")
		}
		if opts["foo"] != "bar" {
			t.Errorf("Options parameter doesn't get options")
		}
		if a1 != theCallable1 {
			t.Errorf("arg1 not passed")
		}
		if a2 != theCallable2 {
			t.Errorf("arg2 not passed")
		}
	})
	callGood(theFrame, []interface{}{theCallable1, theCallable2}, RawOptions{"foo": "bar"})
	if !called {
		t.Errorf("not called")
	}

	// Argument passing.
	called = false
	f = NewBuiltinFn("f40", func(x string) {
		called = true
		if x != "lorem" {
			t.Errorf("Argument x not passed")
		}
	})
	callGood(theFrame, []interface{}{"lorem"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// Variadic arguments.
	called = false
	f = NewBuiltinFn("f50", func(f *Frame, x ...int) {
		called = true
		if len(x) != 2 || x[0] != 123 || x[1] != 456 {
			t.Errorf("Variadic argument not passed")
		}
	})
	callGood(theFrame, []interface{}{"123", "456"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// Conversion into int and float64.
	called = false
	f = NewBuiltinFn("f60", func(i int, f float64) {
		called = true
		if i != 314 {
			t.Errorf("Integer argument i not passed")
		}
		if f != 1.25 {
			t.Errorf("Float argument f not passed")
		}
	})
	callGood(theFrame, []interface{}{"314", "1.25"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// func(string, ...string)
	called = false
	f = NewBuiltinFn("f65", func(a string, va ...string) {
		called = true
		if a != "lorem" {
			t.Errorf("arg1 not passed")
		}
		if va[0] != "ipsum" {
			t.Errorf("arg2 not passed")
		}
	})
	callGood(theFrame, []interface{}{"lorem", "ipsum"}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// Conversion of supplied inputs.
	called = false
	f = NewBuiltinFn("f70", func(i Inputs) {
		called = true
		var values []interface{}
		i(func(x interface{}) {
			values = append(values, x)
		})
		if len(values) != 2 || values[0] != "foo" || values[1] != "bar" {
			t.Errorf("Inputs parameter didn't get supplied inputs")
		}
	})
	callGood(theFrame, []interface{}{vals.MakeList("foo", "bar")}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// Conversion of implicit inputs.
	inFrame := &Frame{ports: make([]*Port, 3)}
	ch := make(chan interface{}, 10)
	ch <- "foo"
	ch <- "bar"
	close(ch)
	inFrame.ports[0] = &Port{Chan: ch}
	called = false
	f = NewBuiltinFn("f80", func(f *Frame, opts RawOptions, s string, i Inputs) {
		called = true
		var values []interface{}
		i(func(x interface{}) {
			values = append(values, x)
		})
		if s != "s" {
			t.Errorf("Explicit argument not passed")
		}
		if len(values) != 2 || values[0] != "foo" || values[1] != "bar" {
			t.Errorf("Inputs parameter didn't get implicit inputs")
		}
	})
	callGood(inFrame, []interface{}{"s", vals.MakeList("foo", "bar")}, theOptions)
	if !called {
		t.Errorf("not called")
	}

	// Outputting of return values.
	outFrame := &Frame{ports: make([]*Port, 3)}
	ch = make(chan interface{}, 10)
	outFrame.ports[1] = &Port{Chan: ch}
	f = NewBuiltinFn("f90", func(s string) string { return s + "-ret" })
	callGood(outFrame, []interface{}{"arg"}, theOptions)
	select {
	case ret := <-ch:
		if ret != "arg-ret" {
			t.Errorf("Output is not the same as return value")
		}
	default:
		t.Errorf("Return value is not outputted")
	}

	// Conversion of return values.
	f = NewBuiltinFn("f100", func() int { return 314 })
	callGood(outFrame, nil, theOptions)
	select {
	case ret := <-ch:
		if ret != "314" {
			t.Errorf("Return value is not converted to string")
		}
	default:
		t.Errorf("Return value is not outputted")
	}

	// Passing of error return value.
	f = NewBuiltinFn("f110", func() error {
		return theError
	})
	if f.Call(outFrame, nil, theOptions) != theError {
		t.Errorf("Returned error is not passed")
	}
	select {
	case <-ch:
		t.Errorf("Return value is outputted when error is not nil")
	default:
	}

	// Too many arguments.
	f = NewBuiltinFn("f120", func() {
		t.Errorf("Function called when there are too many arguments")
	})
	callBad(theFrame, []interface{}{"x"}, theOptions)

	// Too few arguments.
	f = NewBuiltinFn("f130", func(x string) {
		t.Errorf("Function called when there are too few arguments")
	})
	callBad(theFrame, nil, theOptions)
	f = NewBuiltinFn("f140", func(x string, y ...string) {
		t.Errorf("Function called when there are too few arguments")
	})
	callBad(theFrame, nil, theOptions)

	// Options when the function does not accept options.
	f = NewBuiltinFn("f150", func() {
		t.Errorf("Function called when there are extra options")
	})
	callBad(theFrame, nil, RawOptions{"foo": "bar"})

	// Wrong argument type.
	f = NewBuiltinFn("f160", func(x string) {
		t.Errorf("Function called when arguments have wrong type")
	})
	callBad(theFrame, []interface{}{1}, theOptions)

	// Wrong argument type: cannot convert to int.
	f = NewBuiltinFn("f170", func(x int) {
		t.Errorf("Function called when arguments have wrong type")
	})
	callBad(theFrame, []interface{}{"x"}, theOptions)

	// Wrong argument type: cannot convert to float64.
	f = NewBuiltinFn("f180", func(x float64) {
		t.Errorf("Function called when arguments have wrong type")
	})
	callBad(theFrame, []interface{}{"x"}, theOptions)
}
