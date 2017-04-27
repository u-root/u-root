// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"fmt"
	gobuild "go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	archivers["cpio"] = cpioArchiver{}
}

type cpioArchiver struct {
}

// This archiver creates a cpio archive from the given list of files.
// TODO: Replace this implementation with one which does not require sudo.
// TODO: Preferably, it will share code with cmds/cpio.
func (a cpioArchiver) generate(config Config, files []file) error {
	// TODO: Delete this temporary directory which is too scary because of "sudo rm -rf" --
	// especially considering the directory could contain mount points.
	chrootConfig := config
	tmpDir, err := ioutil.TempDir("", "uroot")
	if err != nil {
		return err
	}

	// We cheat by calling the chroot archiver to create the directory
	// structure and running cpio over it.
	chrootConfig.OutputPath = filepath.Join(tmpDir, "chroot")
	if err := (chrootArchiver{}).generate(chrootConfig, files); err != nil {
		return err
	}

	// TODO: This is crud.
	outputCpio, err := filepath.Abs(config.OutputPath)
	if err != nil {
		return err
	}
	cmd := exec.Command("sh", "-c", "find * | sudo cpio -H newc -o > "+outputCpio)
	cmd.Dir = chrootConfig.OutputPath
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run the cpio file under Linux in QEMU. This requires a bit of setup.
func (a cpioArchiver) run(config Config) error {
	envName := fmt.Sprintf("UROOT_CPIO_RUN_%s", gobuild.Default.GOARCH)
	env := os.Getenv(envName)
	if env == "" {
		return fmt.Errorf(`%s is unset.
To run the cpio file, set UROOT_CPIO_RUN_$GOARCH to a command to be run under sh.
In the command, {} is replaced by the path of the cpio file.  For example:
	$ export UROOT_CPIO_RUN_amd64="qemu-system-x86_64 -kernel $YOUR_BZIMAGE -initrd {} -nographic -m 1G"
	$ u-root --format=cpio --run`, envName)
	}

	cmd := exec.Command("sh", "-c", strings.Replace(env, "{}", config.OutputPath, -1))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
