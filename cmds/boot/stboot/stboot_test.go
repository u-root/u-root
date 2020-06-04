// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForceHTTPS(t *testing.T) {
	var urls = []string{"stboot.dev", "http://stboot.dev", "stboot.dev/file", "http://stboot.dev/file"}
	err := forceHTTPS(urls)
	require.NoError(t, err)
	for _, raw := range urls {
		url, err := url.Parse(raw)
		require.NoError(t, err)
		require.Equal(t, url.Scheme, "https")
	}
}
