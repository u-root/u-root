// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/u-root/u-root/pkg/sos"
)

const (
	HtmlPage = `
<head>
</head>
<body>
  <h1>Upspin</h1>
</body>
`
)

var (
	Port uint
)

type UpspinServer struct {
	// TODO: Implement an Upspin client manager
}

func (us UpspinServer) displayUpspinHandle(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("upspin").Parse(HtmlPage))
	tmpl.Execute(w, "")
}

func (us UpspinServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", us.displayUpspinHandle).Methods("GET")
	return r
}

func (us UpspinServer) Start() {
	listener, port, err := sos.GetListener()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	Port = port
	fmt.Println(sos.StartServiceServer(us.buildRouter(), "upspin", listener, Port))
}

func NewUpspinServer() *UpspinServer {
	return &UpspinServer{}
}
