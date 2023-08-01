// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)
// +build linux,arm64 linux,amd64 linux,riscv64

package strace

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uio/uiotest"
)

func prepareTestCmd(t *testing.T, cmd string) {
	// VM environment doesn't have makefiles or whatever. Just skip it.
	testutil.SkipIfInVMTest(t)

	if _, err := os.Stat(cmd); !os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Failed to find test program %q: %v", cmd, err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	w := uiotest.TestLineWriter(t, "make all")

	c := exec.CommandContext(ctx, "make", "all")
	c.Stdout = w
	c.Stderr = w
	c.Dir = "./test"
	if err := c.Run(); err != nil {
		t.Fatalf("make failed: %v", err)
	}
}

func runAndCollectTrace(t *testing.T, cmd *exec.Cmd) []*TraceRecord {
	// Write strace logs to t.Logf.
	w := uiotest.TestLineWriter(t, "")
	traceChan := make(chan *TraceRecord)
	done := make(chan error, 1)

	go func() {
		done <- Trace(cmd, PrintTraces(w), RecordTraces(traceChan))
		close(traceChan)
	}()

	var events []*TraceRecord
	for r := range traceChan {
		events = append(events, r)
	}

	if err := <-done; err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Trace exited with error -- did you compile the test programs? (cd ./test && make all): %v", err)
		} else {
			t.Errorf("Trace exited with error: %v", err)
		}
	}
	return events
}

func TestSingleThreaded(t *testing.T) {
	prepareTestCmd(t, "./test/hello")

	var b bytes.Buffer
	cmd := exec.Command("./test/hello")
	cmd.Stdout = &b

	runAndCollectTrace(t, cmd)
}

func TestMultiProcess(t *testing.T) {
	prepareTestCmd(t, "./test/fork")

	var b bytes.Buffer
	cmd := exec.Command("./test/fork")
	cmd.Stdout = &b

	runAndCollectTrace(t, cmd)
}
