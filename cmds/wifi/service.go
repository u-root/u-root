// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/u-root/u-root/pkg/wifi"
)

// Exported Constants

var (
	DefaultBufferSize = 4
)

// Exported Types

type State struct {
	CurEssid        string
	ConnectingEssid string
	Refreshing      bool
	NearbyWifis     []wifi.WifiOption
}

type ConnectReqMsg struct {
	// private channel that the requesting routine is listening on
	c    chan (error)
	args []string
}

type RefreshReqMsg chan error

type StateReqMsg chan State

type WifiService struct {
	// Share Resource between goroutines
	wifiWorker wifi.Wifi

	// Communicating Channels between internal goroutines
	connectArbitratorQuit chan bool
	refreshPoolerQuit     chan bool
	stateTrackerQuit      chan bool
	stateUpdateChan       chan stateUpdateMsg

	// Communicating Channels with Server goroutines
	connectReqChan chan ConnectReqMsg
	refreshReqChan chan RefreshReqMsg
	stateReqChan   chan StateReqMsg
}

// Internal type

type stateComponent int

const (
	curEssidComp stateComponent = iota
	connectingEssidComp
	refreshingComp
	nearbyWifisComp
)

type stateUpdateMsg struct {
	key        stateComponent
	val        interface{}
	doneUpdate chan bool // Used when need to ensure happening before relationship

}

func (ws WifiService) startConnectWifiArbitrator() {
	curEssid, connectingEssid := "", ""
	workDone := make(chan error, 1)
	var winningChan chan error
	for {
		select {
		case req := <-ws.connectReqChan:
			if connectingEssid == "" {
				// The requesting routine wins
				connectingEssid = req.args[0]
				ws.stateUpdateChan <- stateUpdateMsg{connectingEssidComp, connectingEssid, nil}
				winningChan = req.c
				// Starts connection
				go func(args ...string) {
					workDone <- ws.wifiWorker.Connect(args...)
				}(req.args...)
			} else {
				// The requesting routine loses
				req.c <- fmt.Errorf("Service is trying to connect to %s", connectingEssid)
			}
		case err := <-workDone:
			// Update states
			if err != nil {
				curEssid, _ = ws.wifiWorker.ScanCurrentWifi()
			} else {
				curEssid = connectingEssid
			}
			doneUpdate := make(chan bool, 2)
			ws.stateUpdateChan <- stateUpdateMsg{curEssidComp, curEssid, doneUpdate}
			connectingEssid = ""
			ws.stateUpdateChan <- stateUpdateMsg{connectingEssidComp, connectingEssid, doneUpdate}
			<-doneUpdate
			<-doneUpdate
			winningChan <- err
		case <-ws.connectArbitratorQuit:
			return
		}
	}
}

func (ws WifiService) startRefreshPooler() {
	workDone := make(chan bool, 1)
	pool := make(chan RefreshReqMsg, DefaultBufferSize)
	refreshing := false
	// Pooler
	for {
		select {
		case req := <-ws.refreshReqChan:
			if !refreshing {
				refreshing = true
				ws.stateUpdateChan <- stateUpdateMsg{refreshingComp, refreshing, nil}

				// Notifier
				go func(p chan RefreshReqMsg) {
					o, err := ws.wifiWorker.ScanWifi()
					doneUpdate := make(chan bool, 1)
					ws.stateUpdateChan <- stateUpdateMsg{nearbyWifisComp, o, doneUpdate}
					<-doneUpdate
					workDone <- true
					for ch := range p {
						ch <- err
					}
				}(pool)
			}
			pool <- req
		case <-workDone:
			close(pool)
			refreshing = false
			ws.stateUpdateChan <- stateUpdateMsg{refreshingComp, refreshing, nil}
			pool = make(chan RefreshReqMsg, DefaultBufferSize)
		case <-ws.refreshPoolerQuit:
			return
		}
	}
}

func (ws WifiService) startStateTracker() {
	state := State{
		CurEssid:        "",
		ConnectingEssid: "",
		Refreshing:      false,
		NearbyWifis:     nil,
	}
	for {
		select {
		case r := <-ws.stateReqChan:
			r <- state
		case updateMsg := <-ws.stateUpdateChan:
			updateState(&state, updateMsg)
			if updateMsg.doneUpdate != nil {
				updateMsg.doneUpdate <- true
			}
		case <-ws.stateTrackerQuit:
			return
		}
	}
}

func updateState(state *State, update stateUpdateMsg) {
	switch update.key {
	case curEssidComp:
		state.CurEssid = update.val.(string)
	case connectingEssidComp:
		state.ConnectingEssid = update.val.(string)
	case refreshingComp:
		state.Refreshing = update.val.(bool)
	case nearbyWifisComp:
		state.NearbyWifis = update.val.([]wifi.WifiOption)
	}
}

func NewWifiService(w wifi.Wifi) WifiService {
	return WifiService{
		wifiWorker:            w,
		connectArbitratorQuit: make(chan bool, 1),
		refreshPoolerQuit:     make(chan bool, 1),
		stateTrackerQuit:      make(chan bool, 1),
		stateUpdateChan:       make(chan stateUpdateMsg, 4),
		connectReqChan:        make(chan ConnectReqMsg, DefaultBufferSize),
		refreshReqChan:        make(chan RefreshReqMsg, DefaultBufferSize),
		stateReqChan:          make(chan StateReqMsg, DefaultBufferSize),
	}
}

func (ws WifiService) Start() {
	go ws.startConnectWifiArbitrator()
	go ws.startRefreshPooler()
	go ws.startStateTracker()
}

func (ws WifiService) Shutdown() {
	ws.connectArbitratorQuit <- true
	ws.refreshPoolerQuit <- true
	ws.stateTrackerQuit <- true
}

func (ws WifiService) GetState() State {
	c := make(chan State, 1)
	ws.stateReqChan <- (c)
	return <-c
}

func (ws WifiService) Connect(args []string) error {
	c := make(chan error, 1)
	ws.connectReqChan <- ConnectReqMsg{c, args}
	return <-c
}

func (ws WifiService) Refresh() error {
	c := make(chan error, 1)
	ws.refreshReqChan <- (c)
	return <-c
}
