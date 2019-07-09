// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func DummyCheck(args CheckArgs) error {
	return nil
}

var countCheckCallCount = 0

func CountCheck(args CheckArgs) error {
	countCheckCallCount++
	return fmt.Errorf("TEST ERROR")
}

func TestFuncName(t *testing.T) {
	require.Equal(t, "DummyCheck", funcName(DummyCheck))
}

func TestCall(t *testing.T) {
	registerCheckFun(CountCheck)

	err := Call("CountCheck", nil)

	require.Equal(t, "TEST ERROR", err.Error())
	require.Equal(t, 1, countCheckCallCount)
}

func TestCallUnregistered(t *testing.T) {
	err := Call("UnregisteredFunction", nil)

	require.Regexp(t, "Invalid CheckFun name: .+", err.Error())
}
