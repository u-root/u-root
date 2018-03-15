// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type WifiErrorTestCase struct {
	name   string
	args   []string
	expect string
}

var (
	errorTestcases = []WifiErrorTestCase{
		{
			name:   "More elements than needed",
			args:   []string{"a", "a", "a", "a"},
			expect: "Usage",
		},
		{
			name:   "Flags, More elements than needed",
			args:   []string{"-i=123", "a", "a", "a", "a"},
			expect: "Usage",
		},
	}
)

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestWifiErrors(t *testing.T) {
	// Set up
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Tests
	for _, test := range errorTestcases {
		c := exec.Command(execPath, test.args...)
		_, e, _ := run(c)
		if !strings.Contains(e, test.expect) {
			t.Logf("TEST %v", test.name)
			execStatement := fmt.Sprintf("exec(wifi %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
			t.Errorf("%s\ngot:%s\nwant:%s", execStatement, e, test.expect)
		}
	}
}

func connectWifiArbitratorSetup(curEssid, connEssid string, bufferSize int) {
	CurEssid = curEssid
	ConnectingEssid = connEssid
	ConnectReqChan = make(chan ConnectReqChanMsg, bufferSize)
	go connectWifiArbitrator()
}

func TestConnectWifiArbitrator(t *testing.T) {
	done := make(chan error)

	// Accept Req, Connect success
	func() {
		connectWifiArbitratorSetup("", "", 2)
		defer close(ConnectReqChan)
		c := make(chan error, 1)
		ConnectReqChan <- ConnectReqChanMsg{c, "stub", []byte("stub"), false}
		err := <-c
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
			done <- nil
			return
		}
		ConnectReqChan <- ConnectReqChanMsg{done, "stub", []byte("stub"), true}
	}()

	<-done

	if ConnectingEssid != "" {
		t.Errorf("\ngot: %v\nwant: %v", ConnectingEssid, "")
	}
	if CurEssid != "stub" {
		t.Errorf("\ngot: %v\nwant: %v", CurEssid, "stub")
	}

	// Accept Req, Connect fails
	func() {
		connectWifiArbitratorSetup("", "", 2)
		defer close(ConnectReqChan)
		c := make(chan error, 1)
		ConnectReqChan <- ConnectReqChanMsg{c, "stub", []byte("stub"), false}
		err := <-c
		if err != nil {
			t.Errorf("\ngot: %v\nwant: %v", err, nil)
			done <- nil
			return
		}
		ConnectReqChan <- ConnectReqChanMsg{done, "stub", []byte("stub"), false}
	}()

	<-done

	if ConnectingEssid != "" {
		t.Errorf("\ngot: %v\nwant: %v", ConnectingEssid, "")
	}

	// Reject Req
	func() {
		connectWifiArbitratorSetup("", "stub", 2)
		defer close(ConnectReqChan)
		c := make(chan error, 1)
		ConnectReqChan <- ConnectReqChanMsg{c, "stub2", []byte("stub"), false}
		err := <-c
		if !reflect.DeepEqual(err, fmt.Errorf("Service is trying to connect to %s", "stub")) {
			t.Errorf("\ngot: %v\nwant: %v", err, fmt.Errorf("Service is trying to connect to %s", "stub"))
		}
	}()

	if ConnectingEssid != "stub" {
		t.Errorf("\ngot: %v\nwant: %v", ConnectingEssid, "stub")
	}
	if CurEssid != "" {
		t.Errorf("\ngot: %v\nwant: %v", CurEssid, "")
	}

	// Two competing Go Routines
	func() {
		connectWifiArbitratorSetup("", "", 2)
		defer close(ConnectReqChan)
		c1 := make(chan error, 1)
		ConnectReqChan <- ConnectReqChanMsg{c1, "stub1", []byte("stub1"), false}
		<-c1 // Now the channel has accepted me
		go func() {
			c2 := make(chan error, 1)
			ConnectReqChan <- ConnectReqChanMsg{c2, "stub2", []byte("stub2"), false}
			err := <-c2
			if !reflect.DeepEqual(err, fmt.Errorf("Service is trying to connect to %s", "stub1")) {
				t.Errorf("\ngot: %v\nwant: %v", err, fmt.Errorf("Service is trying to connect to %s", "stub1"))
			}
			done <- nil
		}()
		<-done
	}()
}

func TestRaceCondConnectWifiArbitrator(t *testing.T) {
	//Set Up
	numGoRoutines := 100
	connectWifiArbitratorSetup("", "", numGoRoutines)
	defer close(ConnectReqChan)

	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			c := make(chan error, 1)
			routineIdStub := fmt.Sprintf("stub%v", idx)
			ConnectReqChan <- ConnectReqChanMsg{c, routineIdStub, []byte(routineIdStub), false}
		}(i)
	}
	wg.Wait()
}
