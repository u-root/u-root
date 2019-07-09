// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
)

var funcNameRegex = regexp.MustCompile("[^.]+$")

var globalCheckRepo = make(map[string]CheckFun)

// registerCheckFun registers a function so that it can be referenced by name in Call()
func registerCheckFun(checkFun CheckFun) {
	name := funcName(checkFun)
	globalCheckRepo[name] = checkFun
}

// Call a registered CheckFun of given name
func Call(name string, args CheckArgs) error {
	checkFun := globalCheckRepo[name]
	if checkFun == nil {
		return fmt.Errorf("Invalid CheckFun name: %v. Please ensure the name is correct and that the check function was properly registered", name)
	}
	return checkFun(args)
}

// ListRegistered returns all registered CheckFuns
func ListRegistered() []string {
	registered := make([]string, 0)
	for name := range globalCheckRepo {
		registered = append(registered, name)
	}
	return registered
}

func funcName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name = funcNameRegex.FindString(name)
	return name
}
