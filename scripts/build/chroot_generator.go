// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	archiveGenerators["chroot"] = chrootGenerator{}
}

type chrootGenerator struct {
}

// The chroot generator dumps the rootfs tree into a directory appropriate for
// running under a chroot and namespaces.
func (g chrootGenerator) generate(config Config, files []file) error {
	// Since files is sorted, we can guarantee that the directories are created
	// before their children.
	for _, f := range files {
		f.path = filepath.Join(config.OutputPath, f.path)
		if err := createFile(f); err != nil {
			return err
		}
	}
	return nil
}

// Create an actual file from the file struct. The file struct should already
// be prepended with the location of the chroot.
func createFile(f file) error {
	// TODO: set uid and gid
	mode := fmt.Sprintf("%03o", f.mode&os.ModePerm)
	major := fmt.Sprint(major(f.rdev))
	minor := fmt.Sprint(minor(f.rdev))

	// Special file types
	// I'm sure many of these commands can be done through syscalls rather than
	// forking, however we need to perform the priviledge escalation via sudo
	// (asking for their password).
	switch f.mode & (os.ModeType | os.ModeCharDevice) {
	case os.ModeDir:
		return exec.Command("sudo", "mkdir", "--mode="+mode, f.path).Run()
	case os.ModeDevice:
		return exec.Command("sudo", "mknod", "--mode="+mode, f.path, "b", major, minor).Run()
	case os.ModeDevice | os.ModeCharDevice:
		return exec.Command("sudo", "mknod", "--mode="+mode, f.path, "c", major, minor).Run()
	case os.ModeNamedPipe:
		return exec.Command("sudo", "mknod", "--mode="+mode, f.path, "p").Run()
	case os.ModeSymlink:
		data, err := ioutil.ReadAll(f.data)
		if err != nil {
			return err
		}
		if err := exec.Command("sudo", "ln", "-s", string(data), f.path).Run(); err != nil {
			return err
		}
		return exec.Command("sudo", "chmod", mode, f.path).Run()
	case os.ModeSocket:
		return errors.New("making sockets not supported yet") // TODO
	}

	// Regular files
	cmd := exec.Command("sudo", "tee", f.path)
	cmd.Stdin = f.data
	if err := cmd.Run(); err != nil {
		return err
	}
	return exec.Command("sudo", "chmod", mode, f.path).Run()
}

// Run the rootfs under a chroot jail.
func (g chrootGenerator) run(config Config) error {
	// TODO: Until https://github.com/golang/go/issues/19661 is fixed,
	// exec.Command is insufficient for making mount namespaces, so instead we
	// rely on a moderately up-to-date unshare command.
	args := []string{
		"unshare",

		// mounts disappear after the process ends
		"--mount", "--propagation=private",

		// new root filessytem
		"chroot", config.OutputPath, "/init",
	}
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
