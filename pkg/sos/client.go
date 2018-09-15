// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
)

// RegistersNecessaryPatterns registers all the necessary patterns needed
// to make a service becomes a SoS client.
func RegistersNecessaryPatterns(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")
}

// RegisterServiceWithSos tries to register a service with SoS.
// If an non-nil error is returned, the service needs to exit immediately.
func RegisterServiceWithSos(service string, port uint) error {
	return registerServiceWithSos(service, port, "http://localhost:"+PortNum)
}

func registerServiceWithSos(service string, port uint, sosServerURL string) error {
	m := RegisterReqJson{service, port}
	return makeRequestToServer("POST", sosServerURL+"/register", m)
}

// UnregisterServiceWithSos makes a request to SoS Server to unregister the service.
// This function should be called before a service exit.
func UnregisterServiceWithSos(service string) error {
	return unregisterServiceWithSos(service, "http://localhost:"+PortNum)
}

func unregisterServiceWithSos(service string, sosServerURL string) error {
	m := UnRegisterReqJson{service}
	return makeRequestToServer("POST", sosServerURL+"/unregister", m)
}

func makeRequestToServer(reqType, url string, reqJSON interface{}) error {
	b, err := json.Marshal(reqJSON)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		decoder := json.NewDecoder(res.Body)
		defer res.Body.Close()
		var retMsg struct{ Error string }
		if err := decoder.Decode(&retMsg); err != nil {
			return err
		}
		if retMsg.Error != "" {
			return fmt.Errorf(retMsg.Error)
		}
	}

	return nil
}

// GetListener starts listener on a random port in localhost
// and returns the listener and the port that the listener is on.
// Remember to close the listener with:
//
// defer listener.close()
//
// or if it's used in a server, remember to shutdown the server.
func GetListener() (net.Listener, uint, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, 0, err
	}

	addrSplit := strings.Split(listener.Addr().String(), ":")
	if len(addrSplit) != 2 {
		listener.Close()
		return nil, 0, fmt.Errorf("Address format not recognized: %v", listener.Addr().String())
	}

	port, err := strconv.ParseUint(addrSplit[1], 10, 32)
	if err != nil {
		listener.Close()
		return nil, 0, err
	}
	return listener, uint(port), nil
}

// StartServiceServer establishes registers all necessary patterns to the router passed in,
// registers the service with SoS using the port passed in, and starts serving the service on the
// listener passed in. If any of the above step fails, this function will return an error.
// This function wraps around RegistersNecessaryPatterns, RegisterServiceWithSos, and UnregisterServiceWithSos.
// If no extenral settings are required, instead of calling each of the above separately, one can call
// the GetListener function and pass the result into this function to start and serve their HTTP server right away.
func StartServiceServer(router *mux.Router, serviceName string, listener net.Listener, port uint) error {
	RegistersNecessaryPatterns(router)
	if err := RegisterServiceWithSos(serviceName, port); err != nil {
		return err
	}
	defer UnregisterServiceWithSos(serviceName)

	shutdownChan := make(chan bool, 2)
	server := http.Server{Handler: router}
	defer func() {
		shutdownChan <- true
	}() // Use to collect any other failure besides signals

	// Signals Collector
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		fmt.Printf("Received: %v\n", sig)
		shutdownChan <- true
	}()

	// Server Shutdown code
	go func() {
		<-shutdownChan
		fmt.Println("Shutting down...")
		server.Shutdown(context.Background())
	}()

	return server.Serve(listener)
}
