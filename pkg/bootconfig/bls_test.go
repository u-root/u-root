// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

var testconfigs = map[string]BootConfig{
	"bls_fromspec.conf": BootConfig{
		Name:       "Fedora 19 (Rawhide)",
		Kernel:     "/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/linux",
		KernelArgs: "root=UUID=6d3376e4-fc93-4509-95ec-a21d68011da2",
		Initramfs:  "/6a9857a393724b7a981ebb5b8495b9ea/3.8.0-2.fc19.x86_64/initrd",
	},
	"bls_rhel8.conf": BootConfig{
		Name:       "Red Hat Enterprise Linux (4.16.18-114_internal1_3521_gf2e37788fa4a) 8.0 (Ootpa)",
		Kernel:     "/vmlinuz-4.16.18-114_internal1_3521_gf2e37788fa4a",
		KernelArgs: "$kernelopts",
		Initramfs:  "/initramfs-4.16.18-114_internal1_3521_gf2e37788fa4a.img",
	},
}

func TestParseBLSConfigs(t *testing.T) {
	for f, want := range testconfigs {
		data, err := ioutil.ReadFile(path.Join("testdata/loader/entries", f))
		if err != nil {
			panic(err)
		}
		got, err := parseBLSEntry(string(data))
		require.NoError(t, err)
		require.Equal(t, want, *got)
	}
}
