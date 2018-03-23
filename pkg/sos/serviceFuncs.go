// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Mux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func RegistersNeccesaryPatterns(m Mux) {
	m.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}

func RegisterServiceWithSos(service string, port uint) error {
	return registerServiceWithSos(service, port, "localhost:"+PortNum)
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
