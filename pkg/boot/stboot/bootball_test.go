// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBootBallFromConfig(t *testing.T) {
	file := "testdata/testConfigDir/stconfig.json"
	ball, err := BootBallFromConfig(file)
	t.Logf("tmp config dir: %s", ball.dir)
	require.NoError(t, err)
	_, err = os.Stat(ball.dir)
	require.NoError(t, err)
	// todo: test files, too
}
