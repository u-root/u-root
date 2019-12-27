package eval

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/u-root/u-root/cmds/core/elvish/eval/vals"
	"github.com/u-root/u-root/cmds/core/elvish/hash"
)

var ErrArgs = errors.New("args error")

// BuiltinFn uses reflection to wrap Go functions into Elvish functions.
// Functions with simple signatures are handled by commonCallable,
// more elaborate cases are handled by CustomCallable interface, e.g.
// when importing the necessary types would have created a dependency loop.
//
// Parameters are converted using ScanToGo.
//
// Return values go to the channel part of the stdout port, after being
// converted using goToElv. If the last return value has type error and is not
// nil, it is turned into an exception and no ouputting happens. If the last
// return value is a nil error, it is ignored.
//
// Note: reflect.Call is deliberately not used as it disables DCE
// (see https://github.com/u-root/u-root/issues/1477).
type BuiltinFn struct {
	name string
	impl CustomCallable

	// Type information of impl.

	frame   bool
	options bool
	inputs  bool
	// Type of "normal" (non-Frame, non-Options, non-variadic) arguments.
	normalArgs []reflect.Type
	// Type of variadic arguments, nil if function is non-variadic
	variadicArg reflect.Type
}

var _ Callable = &BuiltinFn{}

// CustomCallable is an interface for functions that have complex signatures
// not covered by commonCallable.
type CustomCallable interface {
	// Target returns the callable's target function, needed to examine arguments.
	Target() interface{}
	// Call invokes the target function.
	Call(f *Frame, args []interface{}, opts RawOptions, inputs Inputs) ([]interface{}, error)
}

type (
	Inputs func(func(interface{}))
)

var (
	frameType      = reflect.TypeOf((*Frame)(nil))
	rawOptionsType = reflect.TypeOf(RawOptions(nil))
	inputsType     = reflect.TypeOf(Inputs(nil))
)

// NewBuiltinFnCustom creates a new ReflectBuiltinFn instance.
func NewBuiltinFnCustom(name string, impl CustomCallable) *BuiltinFn {
	implType := reflect.TypeOf(impl.Target())
	b := &BuiltinFn{name: name, impl: impl}

	i := 0
	if i < implType.NumIn() && implType.In(i) == frameType {
		b.frame = true
		i++
	}
	if i < implType.NumIn() && implType.In(i) == rawOptionsType {
		b.options = true
		i++
	}
	for ; i < implType.NumIn(); i++ {
		paramType := implType.In(i)
		if i == implType.NumIn()-1 {
			if implType.IsVariadic() {
				b.variadicArg = paramType.Elem()
				break
			} else if paramType == inputsType {
				b.inputs = true
				break
			}
		}
		b.normalArgs = append(b.normalArgs, paramType)
	}
	return b
}

// NewBuiltinFn creates a new ReflectBuiltinFn instance.
func NewBuiltinFn(name string, impl interface{}) *BuiltinFn {
	if _, err := callFunc(impl, nil, nil, nil, nil, true); err != nil {
		panic(fmt.Sprintf("%s: %v", name, err))
	}
	return NewBuiltinFnCustom(name, &commonCallable{target: impl})
}

// Kind returns "fn".
func (*BuiltinFn) Kind() string {
	return "fn"
}

// Equal compares identity.
func (b *BuiltinFn) Equal(rhs interface{}) bool {
	return b == rhs
}

// Hash hashes the address.
func (b *BuiltinFn) Hash() uint32 {
	return hash.Hash(b)
}

// Repr returns an opaque representation "<builtin $name>".
func (b *BuiltinFn) Repr(int) string {
	return "<builtin " + b.name + ">"
}

// error(nil) is treated as nil by reflect.TypeOf, so we first get the type of
// *error and use Elem to obtain type of error.
var errorType = reflect.TypeOf((*error)(nil)).Elem()

var errNoOptions = errors.New("function does not accept any options")

// Call calls the implementation using reflection.
func (b *BuiltinFn) Call(f *Frame, args []interface{}, opts map[string]interface{}) error {
	if b.variadicArg != nil {
		if len(args) < len(b.normalArgs) {
			return fmt.Errorf("%s: want %d or more arguments, got %d",
				b.name, len(b.normalArgs), len(args))
		}
	} else if b.inputs {
		if len(args) != len(b.normalArgs) && len(args) != len(b.normalArgs)+1 {
			return fmt.Errorf("%s: want %d or %d arguments, got %d",
				b.name, len(b.normalArgs), len(b.normalArgs)+1, len(args))
		}
	} else if len(args) != len(b.normalArgs) {
		return fmt.Errorf("%s: want %d arguments, got %d", b.name, len(b.normalArgs), len(args))
	}
	if !b.options && len(opts) > 0 {
		return errNoOptions
	}

	var goArgs []interface{}
	for i, arg := range args {
		var typ reflect.Type
		if i < len(b.normalArgs) {
			typ = b.normalArgs[i]
		} else if b.variadicArg != nil {
			typ = b.variadicArg
		} else if b.inputs {
			break // Handled after the loop
		} else {
			panic("impossible")
		}
		ptr := reflect.New(typ)
		err := vals.ScanToGo(arg, ptr.Interface())
		if err != nil {
			return fmt.Errorf("%s: wrong type of %d'th argument: %v", b.name, i+1, err)
		}
		goArgs = append(goArgs, ptr.Elem().Interface())
	}

	var inputs Inputs
	if b.inputs {
		if len(args) == len(b.normalArgs) {
			inputs = Inputs(f.IterateInputs)
		} else {
			// Wrap an iterable argument in Inputs.
			iterable := args[len(args)-1]
			inputs = Inputs(func(f func(interface{})) {
				err := vals.Iterate(iterable, func(v interface{}) bool {
					f(v)
					return true
				})
				maybeThrow(err)
			})
		}
	}

	outs, err := b.impl.Call(f, goArgs, opts, inputs)
	if err != nil {
		return err
	}
	for _, out := range outs {
		f.OutputChan() <- vals.FromGo(out)
	}
	return nil
}

type commonCallable struct {
	target interface{}
}

func (c *commonCallable) Target() interface{} {
	return c.target
}

func (c *commonCallable) Call(
	f *Frame, args []interface{}, opts RawOptions, inputs Inputs) ([]interface{}, error) {
	return callFunc(c.target, f, args, opts, inputs, false)
}

func callFunc(fnp interface{}, f *Frame, args []interface{}, opts RawOptions, inputs Inputs, checkOnly bool) ([]interface{}, error) {
	switch fn := fnp.(type) {
	case func():
		if checkOnly {
			return nil, nil
		}
		fn()
		return nil, nil

	case func() error:
		if checkOnly {
			return nil, nil
		}
		err := fn()
		return nil, err

	case func() int:
		if checkOnly {
			return nil, nil
		}
		out := fn()
		return []interface{}{out}, nil

	case func(*Frame):
		if checkOnly {
			return nil, nil
		}
		fn(f)
		return nil, nil

	case func(*Frame, RawOptions):
		if checkOnly {
			return nil, nil
		}
		fn(f, opts)
		return nil, nil

	case func(*Frame, RawOptions, Callable, Callable):
		if checkOnly {
			return nil, nil
		}
		fn(f, opts, args[0].(Callable), args[1].(Callable))
		return nil, nil

	case func(*Frame, RawOptions, string, Inputs):
		if checkOnly {
			return nil, nil
		}
		fn(f, opts, args[0].(string), inputs)
		return nil, nil

	case func(*Frame, ...int):
		if checkOnly {
			return nil, nil
		}
		var vargs []int
		for _, arg := range args {
			vargs = append(vargs, arg.(int))
		}
		fn(f, vargs...)
		return nil, nil

	case func(*Frame, ...interface{}) error:
		if checkOnly {
			return nil, nil
		}
		err := fn(f, args...)
		return nil, err

	case func(*Frame, interface{}, interface{}, interface{}):
		if checkOnly {
			return nil, nil
		}
		fn(f, args[0], args[1], args[2])
		return nil, nil

	case func(*Frame, ...int) error:
		if checkOnly {
			return nil, nil
		}
		var vargs []int
		for _, arg := range args {
			vargs = append(vargs, arg.(int))
		}
		err := fn(f, vargs...)
		return nil, err

	case func(*Frame, ...string) error:
		if checkOnly {
			return nil, nil
		}
		var vargs []string
		for _, arg := range args {
			vargs = append(vargs, arg.(string))
		}
		err := fn(f, vargs...)
		return nil, err

	case func(*Frame, string):
		if checkOnly {
			return nil, nil
		}
		fn(f, args[0].(string))
		return nil, nil

	case func(*Frame, string) error:
		if checkOnly {
			return nil, nil
		}
		err := fn(f, args[0].(string))
		return nil, err

	case func(Inputs):
		if checkOnly {
			return nil, nil
		}
		fn(inputs)
		return nil, nil

	case func(RawOptions):
		if checkOnly {
			return nil, nil
		}
		fn(opts)
		return nil, nil

	case func(RawOptions, ...interface{}):
		if checkOnly {
			return nil, nil
		}
		fn(opts, args...)
		return nil, nil

	case func(int, float64):
		if checkOnly {
			return nil, nil
		}
		fn(args[0].(int), args[1].(float64))
		return nil, nil

	case func(...int) error:
		if checkOnly {
			return nil, nil
		}
		var vargs []int
		for _, arg := range args {
			vargs = append(vargs, arg.(int))
		}
		err := fn(vargs...)
		return nil, err

	case func(float64):
		if checkOnly {
			return nil, nil
		}
		fn(args[0].(float64))
		return nil, nil

	case func(int):
		if checkOnly {
			return nil, nil
		}
		fn(args[0].(int))
		return nil, nil

	case func(string):
		if checkOnly {
			return nil, nil
		}
		fn(args[0].(string))
		return nil, nil

	case func(string) string:
		if checkOnly {
			return nil, nil
		}
		out := fn(args[0].(string))
		return []interface{}{out}, nil

	case func(string, ...string):
		if checkOnly {
			return nil, nil
		}
		var vargs []string
		for _, arg := range args[1:] {
			vargs = append(vargs, arg.(string))
		}
		fn(args[0].(string), vargs...)
		return nil, nil
	}

	sig := fmt.Sprintf("%#v", fnp)[1:]
	sig = sig[:strings.LastIndex(sig, "(")-1]
	return nil, fmt.Errorf("unsupported function signature: %s", sig)
}
