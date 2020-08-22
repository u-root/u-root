// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

func testMultiboot(t *testing.T, kernel string) {
	var serial wc

	src := fmt.Sprintf("/home/circleci/%v", kernel)
	if tk := os.Getenv("UROOT_MULTIBOOT_TEST_KERNEL_DIR"); len(tk) > 0 {
		src = filepath.Join(tk, kernel)
	}
	if _, err := os.Stat(src); err != nil && os.IsNotExist(err) {
		t.Skip("multiboot kernel is not present")
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		BuildOpts: uroot.Opts{
			ExtraFiles: []string{
				src + ":kernel",
			},
		},
		TestCmds: []string{
			`kexec -l kernel -e -d --module="/kernel foo=bar" --module="/bbin/bb"`,
		},
		QEMUOpts: qemu.Options{
			SerialOutput: &serial,
		},
	})
	defer cleanup()

	if err := q.Expect(`"status": "ok"`); err != nil {
		t.Logf(serial.String())
		t.Fatalf(`expected '"status": "ok"', got error: %v`, err)
	}

	if err := q.Expect(`}`); err != nil {
		t.Logf(serial.String())
		t.Fatalf(`expected '}' = end of JSON, got error: %v`, err)
	}

	output := serial.Bytes()
	t.Logf(serial.String())

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
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	for _, kernel := range []string{"/kernel", "/kernel.gz"} {
		t.Run(kernel, func(t *testing.T) {
			testMultiboot(t, kernel)
		})
	}
}
