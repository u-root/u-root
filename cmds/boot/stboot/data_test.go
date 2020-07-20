// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchBootballFiles(t *testing.T) {
	dir := "testdata/datapartition/bootballs"
	ret, err := searchBootballFiles(dir)
	t.Log(ret)
	require.NotEmpty(t, ret)
	require.NoError(t, err)
}
