// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"reflect"
	"sync"
	"testing"

	"github.com/u-root/u-root/pkg/wifi"
)

func setupStubService() (*WifiService, error) {
	wifiWorker, err := wifi.NewStubWorker("", NearbyWifisStub...)
	if err != nil {
		return nil, err
	}
	return NewWifiService(wifiWorker)
}

func TestGetState(t *testing.T) {
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}
	service.Start()
	defer service.Shutdown()

	expectedInitialState := State{
		CurEssid:        "",
		ConnectingEssid: "",
		Refreshing:      false,
		NearbyWifis:     nil,
	}
	state := service.GetState()
	if !reflect.DeepEqual(expectedInitialState, state) {
		t.Errorf("\ngot:%v\nwant:%v\n", state, expectedInitialState)
	}
}

func TestRaceGetState(t *testing.T) {
	// Set Up
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}
	service.Start()
	defer service.Shutdown()
	numGoRoutines := 100

	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.GetState()
		}()
	}
	wg.Wait()
}

func TestConnect(t *testing.T) {
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}
	service.Start()
	defer service.Shutdown()

	if err := service.Connect([]string{EssidStub}); err != nil {
		t.Errorf("error: %v", err)
		return
	}
	s := service.GetState()
	if s.CurEssid != EssidStub {
		t.Errorf("\ngot:%v\nwant:%v\n", s.CurEssid, EssidStub)
	}
}

func TestRaceConnect(t *testing.T) {
	//Set Up
	numGoRoutines := 100
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}

	service.Start()
	defer service.Shutdown()

	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Connect([]string{EssidStub})
		}()
	}
	wg.Wait()
}

func TestRefresh(t *testing.T) {
	//Set Up
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}
	service.Start()
	defer service.Shutdown()

	if err := service.Refresh(); err != nil {
		t.Fatalf("Refresh: error: %v", err)
	}
	s := service.GetState()
	if !reflect.DeepEqual(s.NearbyWifis, NearbyWifisStub) {
		t.Errorf("\ngot:%v\nwant:%v\n", s.NearbyWifis, NearbyWifisStub)
	}
}

func TestRaceRefreshWithinDefaultBufferSize(t *testing.T) {
	//Set Up
	numGoRoutines := DefaultBufferSize
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}

	service.Start()
	defer service.Shutdown()

	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Refresh()
		}()
	}
	wg.Wait()
}

func TestRaceRefreshOverDefaultBufferSize(t *testing.T) {
	// Set Up
	numGoRoutines := DefaultBufferSize * 2
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}

	service.Start()
	defer service.Shutdown()

	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Refresh()
		}()
	}
	wg.Wait()
}

func TestRaceCond(t *testing.T) {
	// Set Up
	numConnectRoutines, numRefreshGoRoutines, numReadGoRoutines := 10, 10, 100
	service, err := setupStubService()
	if err != nil {
		t.Fatal(err)
	}

	service.Start()
	defer service.Shutdown()

	essidChoices := []string{"stub1", "stub2", "stub3"}

	var wg sync.WaitGroup

	for i := 0; i < numConnectRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx := rand.Intn(len(essidChoices))
			service.Connect([]string{essidChoices[idx]})
		}()
	}

	for i := 0; i < numRefreshGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Refresh()
		}()
	}

	for i := 0; i < numReadGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.GetState()
		}()
	}

	wg.Wait()
}
