// Copyright 2012-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// testramfs acts just like a u-root command, with one difference: it will boot
// a VM with the generated initramfs.
// This code calls the u-root command in an effort to replicate, as much as possible,
// what a user would do. It therefore has no switches.
// All these variations work
// GOARCH=amd64 GOOS=linux ./testramfs  ../..
// GOARCH=arm64 GOOS=linux ./testramfs  ../..
// GOARCH=riscv64  GOOS=linux ./testramfs  ../..
// For extra convenience, we leave in the vm/ package cpud so you can,
// in addition to everything else, cpu in and mess around.
// Note: it boots fine, but cpud is not working on mac right now; qemu issue?
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/u-root/cpu/vm"
)

func Execute(stdin io.Reader, stdout io.Writer, stderr io.Writer, goos, goarch string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: testramfs <uroot-dir> [args...]:%w", os.ErrInvalid)
	}

	if len(goos) == 0 {
		goos = runtime.GOOS
	}

	if len(goarch) == 0 {
		goarch = runtime.GOARCH
	}
	dir := args[1]
	args = args[2:]

	// Figure out where the user wants the initramfs to be.
	// The default is the current working directory.
	// TODO: process -o for u-root command
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	ramfs := filepath.Join(cwd, fmt.Sprintf("initramfs.%s_%s.cpio", goos, goarch))

	start := time.Now()
	c := exec.Command("u-root", append([]string{"-o", ramfs}, args...)...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	c.Dir = dir
	if err := c.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "************************************************\n")
	fmt.Fprintf(os.Stderr, "u-root build took %v seconds\n", time.Since(start).Seconds())
	fmt.Fprintf(os.Stderr, "************************************************\n")

	// This will write files for the kernel into cwd.
	// as a side effect, it provides users with a kernel they can use.
	vmi, err := vm.New(goos, goarch)
	if err != nil {
		return err
	}

	c, err = vmi.CommandContext(context.Background(), cwd, ramfs)
	if err != nil {
		return err
	}
	c.Stdin, c.Stdout, c.Stderr = stdin, stdout, stderr
	if err := vmi.StartVM(c); err != nil {
		return err
	}
	return c.Wait()
}

func main() {
	if err := Execute(os.Stdin, os.Stdout, os.Stderr, os.Getenv("GOOS"), os.Getenv("GOARCH"), os.Args); err != nil {
		log.Fatal(err)
	}

}
