package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalln("usage: %s <cpio-path>", os.Args[0])
	}

	if err := Test(flag.Args()[0]); err != nil {
		log.Fatalln(err)
	}
}

func Test(name string) error {
	// So, what's the plan here?
	//
	// - new mount namespace
	//   - root mount is a tmpfs mount filled with the archive.
	//
	// - new PID namespace
	//   - archive/init actually runs as PID 1.

	// Whatever setup we do shouldn't affect the caller.
	if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
		return fmt.Errorf("unshare: %v", err)
	}

	tempDir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	if err := syscall.Mount("", tempDir, "tmpfs", 0, ""); err != nil {
		return fmt.Errorf("mount: %v", err)
	}
	defer syscall.Unmount(tempDir, syscall.MNT_DETACH)

	f, err := os.Open(name)
	if err != nil {
		return err
	}
	archiver, err := cpio.Format("newc")
	if err != nil {
		return err
	}

	r := archiver.Reader(f)
	for {
		rec, err := r.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		cpio.CreateFileInRoot(rec, tempDir)
	}

	if err := syscall.Chdir(tempDir); err != nil {
		return fmt.Errorf("chdir: %v", err)
	}
	pivotDir, err := ioutil.TempDir(tempDir, ".pivot-root")
	if err != nil {
		return err
	}
	if err := syscall.PivotRoot(tempDir, pivotDir); err != nil {
		return fmt.Errorf("pivot_root: %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir: %v", err)
	}
	if err := syscall.Unmount(filepath.Base(pivotDir), syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount: %v", err)
	}

	cmd := exec.Command("/init")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Setctty:    true,
		//Setsid:     true,
		Cloneflags: syscall.CLONE_NEWPID,
	}
	return cmd.Run()
}
