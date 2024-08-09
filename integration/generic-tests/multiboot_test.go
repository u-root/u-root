// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
	"github.com/u-root/u-root/pkg/boot/multiboot"
)

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

func testMultiboot(t *testing.T, kernel string) {
	src := filepath.Join(os.Getenv("UROOT_MULTIBOOT_TEST_KERNEL_DIR"), kernel)
	if _, err := os.Stat(src); err != nil && os.IsNotExist(err) {
		t.Skip("multiboot kernel is not present")
	} else if err != nil {
		t.Error(err)
	}

	script := `
		kexec -l kernel -d --module="/kernel foo=bar" --module="/bbin/bb"
		sync
		kexec -e
	`
	var b bytes.Buffer
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			// Build kexec as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands(
				"github.com/u-root/u-root/cmds/core/kexec",
				"github.com/u-root/u-root/cmds/core/sync",
			),
			uimage.WithFiles(
				src+":kernel",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithSerialOutput(nopCloser{&b}),
			qemu.WithVMTimeout(time.Minute),
		),
	)

	if _, err := vm.Console.ExpectString(`"status": "ok"`); err != nil {
		t.Errorf(`expected '"status": "ok"', got error: %v`, err)
	}
	if _, err := vm.Console.ExpectString(`}`); err != nil {
		t.Errorf(`expected '}' = end of JSON, got error: %v`, err)
	}
	if err := vm.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = vm.Wait()

	output := b.Bytes()

	i := bytes.Index(output, []byte(multiboot.DebugPrefix))
	if i == -1 {
		t.Fatalf("%q prefix not found in output", multiboot.DebugPrefix)
	}
	output = output[i+len(multiboot.DebugPrefix):]
	if i = bytes.Index(output, []byte{'\n'}); i == -1 {
		t.Fatalf("Cannot find newline character")
	}
	var want multiboot.Description
	if err := json.Unmarshal(output[:i], &want); err != nil {
		t.Fatalf("Cannot unmarshal multiboot debug information: %v", err)
	}

	const starting = "Starting multiboot kernel"
	if i = bytes.Index(output, []byte(starting)); i == -1 {
		t.Fatalf("Multiboot kernel was not executed")
	}
	output = output[i+len(starting):]

	var got multiboot.Description
	if err := json.Unmarshal(output, &got); err != nil {
		t.Fatalf("Cannot unmarshal multiboot information from executed kernel: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("kexec failed: got\n%#v, want\n%#v", got, want)
	}
}

func TestMultiboot(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64, qemu.ArchArm64)

	for _, kernel := range []string{"/kernel", "/kernel.gz"} {
		t.Run(kernel, func(t *testing.T) {
			testMultiboot(t, kernel)
		})
	}
}
