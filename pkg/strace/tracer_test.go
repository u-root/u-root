// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package strace

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/testutil"
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

	r, w := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			t.Logf("Output: %s", s.Text())
		}
		wg.Done()
	}()

	c := exec.CommandContext(ctx, "make", "all")
	c.Stdout = w
	c.Stderr = w
	c.Dir = "./test"
	if err := c.Run(); err != nil {
		t.Fatalf("make failed: %v", err)
	}
	w.Close()
	wg.Wait()
}

func runAndCollectTrace(t *testing.T, cmd *exec.Cmd) []*TraceRecord {
	// Write strace logs to t.Logf.
	r, w := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			t.Logf("Output: %s", s.Text())
		}
		wg.Done()
	}()

	traceChan := make(chan *TraceRecord)
	done := make(chan error, 1)

	go func() {
		done <- Trace(cmd, PrintTraces(w), RecordTraces(traceChan))
		w.Close()
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
	wg.Wait()
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

func TestAddrProcess(t *testing.T) {
	prepareTestCmd(t, "./test/addr")

	var b bytes.Buffer
	cmd := exec.Command("./test/addr")
	cmd.Stdout = &b

	runAndCollectTrace(t, cmd)
}
