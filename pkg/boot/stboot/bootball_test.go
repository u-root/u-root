// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestBootBallFromConfig(t *testing.T) {
// 	file := "testdata/testConfigDir/stconfig.json"
// 	ball, err := BootBallFromConfig(file)
// 	t.Logf("tmp config dir: %s", ball.dir)
// 	require.NoError(t, err)
// 	_, err = os.Stat(ball.dir)
// 	require.NoError(t, err)
// 	// todo: test files, too
// }

func TestCreateFileTree(t *testing.T) {

	kernel := "testdata/files/kernel"
	initramfs := "testdata/files/initramfs"
	tboot := "testdata/files/tboot"
	cert := "testdata/files/cert"
	acm1 := "testdata/files/acm1"
	acm2 := "testdata/files/acm2"
	acm3 := "testdata/files/acm3"

	kernelRelPath := filepath.Join(bootfilesDir, filepath.Base(kernel))
	initramfsRelPath := filepath.Join(bootfilesDir, filepath.Base(initramfs))
	tbootRelPath := filepath.Join(bootfilesDir, filepath.Base(tboot))
	certRelPath := rootCertPath
	acm1RelPath := filepath.Join(acmDir, filepath.Base(acm1))
	acm2RelPath := filepath.Join(acmDir, filepath.Base(acm2))
	acm3RelPath := filepath.Join(acmDir, filepath.Base(acm3))

	var goldenCfg = Stconfig{
		Kernel:    kernelRelPath,
		Initramfs: initramfsRelPath,
		Tboot:     tbootRelPath,
		ACMs:      []string{acm1RelPath, acm2RelPath, acm3RelPath},
	}

	dir, cfg, err := createFileTree(kernel, initramfs, tboot, cert, []string{acm1, acm2, acm3})
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(dir, kernelRelPath))
	require.FileExists(t, filepath.Join(dir, initramfsRelPath))
	require.FileExists(t, filepath.Join(dir, tbootRelPath))
	require.FileExists(t, filepath.Join(dir, certRelPath))
	require.FileExists(t, filepath.Join(dir, acm1RelPath))
	require.FileExists(t, filepath.Join(dir, acm2RelPath))
	require.FileExists(t, filepath.Join(dir, acm3RelPath))
	require.Equal(t, goldenCfg, cfg)

	err = os.RemoveAll(dir)
	require.NoError(t, err)
}

func TestCreateFileTreeFail(t *testing.T) {
	kernel := "wrong/path/kernel"
	initramfs := "testdata/files/initramfs"
	preexec := "testdata/files/preexec"
	cert := "testdata/files/cert"
	acm1 := "testdata/files/acm1"
	acm2 := "testdata/files/acm2"
	acm3 := "testdata/files/acm3"
	dir, _, err := createFileTree(kernel, initramfs, preexec, cert, []string{acm1, acm2, acm3})
	require.Error(t, err)
	require.NoDirExists(t, dir)
}
