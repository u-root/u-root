// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"fmt"
	"reflect"

	"github.com/insomniacslk/makefunc"
)

func makeFunctionThatReturnsError(f interface{}, args ...interface{}) (func() error, error) {
	if err := makefunc.ValidateFunction(f, args...); err != nil {
		return nil, fmt.Errorf("invalid function: %v", err)
	}
	argValues := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}
	return func() error {
		ret := reflect.ValueOf(f).Call(argValues)
		if ret[0].IsNil() {
			return nil
		}
		// the return value being an `error` object is guaranteed by
		// `validateFunction` called above.
		return ret[0].Interface().(error)
	}, nil
}

func makeCheckRunner(f interface{}, args ...interface{}) (CheckRunner, error) {
	return makeFunctionThatReturnsError(f, args...)
}

func makeRemediationRunner(f interface{}, args ...interface{}) (RemediationRunner, error) {
	return makeFunctionThatReturnsError(f, args...)
}
