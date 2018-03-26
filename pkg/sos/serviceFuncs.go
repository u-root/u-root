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

func RegistersNeccesaryPatterns(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")
}

func RegisterServiceWithSos(service string, port uint) error {
	return registerServiceWithSos(service, port, "http://localhost:"+PortNum)
}

func registerServiceWithSos(service string, port uint, sosServerURL string) error {
	m := RegisterReqJson{service, port}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", sosServerURL+"/register", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var retMsg struct{ Error string }
	if err := decoder.Decode(&retMsg); err != nil {
		return err
	}

	if retMsg != struct{ Error string }{} {
		return fmt.Errorf(retMsg.Error)
	}

	return nil
}

func UnregisterServiceWithSos(service string) error {
	return unregisterServiceWithSos(service, "http://localhost:"+PortNum)
}

func unregisterServiceWithSos(service string, sosServerURL string) error {
	m := UnRegisterReqJson{service}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", sosServerURL+"/unregister", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var retMsg struct{ Error string }
	if err := decoder.Decode(&retMsg); err != nil {
		return err
	}

	if retMsg != struct{ Error string }{} {
		return fmt.Errorf(retMsg.Error)
	}

	return nil
}

func StartServiceServer(router *mux.Router, serviceName string, portNumReq *uint) error {
	listener, err := net.Listen("tcp", "localhost:0")
	defer listener.Close()
	if err != nil {
		return err
	}
	port, err := strconv.ParseUint(strings.Split(listener.Addr().String(), ":")[1], 10, 32)
	if err != nil {
		return err
	}
	if portNumReq != nil {
		*portNumReq = uint(port)
	}
	RegistersNeccesaryPatterns(router)
	if err := RegisterServiceWithSos(serviceName, uint(port)); err != nil {
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
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP)
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
