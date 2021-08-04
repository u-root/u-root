// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/vpd"
)

var getter *Getter

func TestNewgetter(t *testing.T) {
	g := NewGetter()
	require.NotNil(t, g)
	require.NotNil(t, g.R)
	require.NotNil(t, g.Out)
}

func TestGetOne(t *testing.T) {
	var buf bytes.Buffer
	getter := NewGetter()
	getter.Out = &buf
	getter.R = vpd.NewReader()
	getter.R.VpdDir = "./tests"
	err := getter.Print("firmware_version")
	require.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "firmware_version(RO) => 1.2.3\n\nfirmware_version(RW) => 3.2.1\n\n", out)
}

func TestGetAll(t *testing.T) {
	var buf bytes.Buffer
	getter := NewGetter()
	getter.Out = &buf
	getter.R = vpd.NewReader()
	getter.R.VpdDir = "./tests"
	err := getter.Print("")
	require.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "firmware_version(RO) => 1.2.3\n\nsomething(RO) => else\n\nfirmware_version(RW) => 3.2.1\n\n", out)
}
