// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/quimage"
	"github.com/u-root/mkuimage/uimage"
)

// Regression test for #2938.
func TestGoshRegression2938(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64, qemu.ArchArm64)

	script := `#!/bin/sh
		echo "bar"
	`
	scriptPath := filepath.Join(t.TempDir(), "script.sh")
	_ = os.WriteFile(scriptPath, []byte(script), 0o777)

	vm := qemu.StartT(t, "vm", qemu.ArchUseEnvv,
		quimage.WithUimageT(t,
			uimage.WithInit("init"),
			uimage.WithUinit("/script.sh"),
			uimage.WithShell("gosh"),
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/sync",
				"github.com/u-root/u-root/cmds/core/echo",
			),
			uimage.WithBinaryCommands(
				"github.com/u-root/u-root/cmds/core/gosh",
			),
			uimage.WithFiles(
				fmt.Sprintf("%s:script.sh", scriptPath),
			),
		),
		qemu.WithVMTimeout(time.Minute),
		qemu.ArbitraryArgs("-m", "512"),
	)

	_ = vm.Wait()
}
