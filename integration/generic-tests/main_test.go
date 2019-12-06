// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"flag"
	"os"
	"testing"
)

var (
	kernelPath    = flag.String("kernel", "", "path to the Linux kernel binary to use for tests")
	qemuPath      = flag.String("qemu", "", "path to the QEMU binary to use for tests")
	testarch      = flag.String("testarch", "", "name of the architecture to use for tests")
	initramfsPath = flag.String("initramfs", "", "path to a custom initramfs to use for tests")
)

var tests = []struct {
	name    string
	runTest func(t *testing.T, initramfs string)
}{
	// uinit_test.go
	{
		name:    "HelloWorld",
		runTest: RunTestHelloWorld,
	},
	{
		name:    "HelloWorldNegative",
		runTest: RunTestHelloWorldNegative,
	},
	{
		name:    "Script",
		runTest: RunTestScript,
	},
	// dhclient_test.go
	{
		name:    "Dhclient",
		runTest: RunTestDhclient,
	},
	{
		name:    "Pxeboot",
		runTest: RunTestPxeboot,
	},
	{
		name:    "QEMUDHCPTimesOut",
		runTest: RunTestQEMUDHCPTimesOut,
	},
	// kexec_test.go
	{
		name:    "MountKexec",
		runTest: RunTestMountKexec,
	},
	// multiboot_test.go
	{
		name:    "Multiboot",
		runTest: RunTestMultiboot,
	},
	// tcz_test.go
	{
		name:    "Tcz",
		runTest: RunTestTczclient,
	},
	// io_test.go
	{
		name:    "IO",
		runTest: RunTestIO,
	},
}

func TestGeneric(t *testing.T) {
	// We use environment variables here for consistency with the automated CI
	// Dockerfiles.
	if len(*kernelPath) > 0 {
		os.Setenv("UROOT_KERNEL", *kernelPath)
	}
	if len(*qemuPath) > 0 {
		os.Setenv("UROOT_QEMU", *qemuPath)
	}
	if len(*testarch) > 0 {
		os.Setenv("UROOT_TESTARCH", *testarch)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runTest == nil {
				t.Fatalf("No runner function found for %s", tt.name)
			}

			tt.runTest(t, *initramfsPath)
		})
	}
}
