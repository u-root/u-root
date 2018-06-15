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
	<style>
	  table {
	    font-family: arial, sans-serif;
	    border-collapse: collapse;
	    width: 100%;
	  }

	  td, th {
	    border: 1px solid #dddddd;
	    text-align: left;
	    padding: 8px;
	  }

	  input {
	    font-size: 120%;
	  }
	</style>
	</head>
	<body>
	  <h1>Upspin</h1>
	  <table style="width:100%">
	    <tr>
	      <th>Username</th>
	      <th>Dir Server</th>
	      <th>Store Server</th>
	      <th>Secret Seed</th>
	      <th></th>
	    </tr>
	    <tr>
	      <td><input type="text" id="user"  class="text" value="{{$Username}}"></td>
	      <td><input type="text" id="dir"   class="text" value="{{$DirServer}}"></td>
	      <td><input type="text" id="store" class="text" value="{{$StoreServer}}"></td>
	      <td><input type="text" id="seed"  class="text" value="{{$SecretSeed}}"></td>
	      <td><input type="submit" id="button" class="btn"></td>
	    </tr>
	  </table>
	</body>
`
)

var (
	Port uint
)

type UpspinServer struct {
	// TODO: Implement an Upspin client manager if needed
}

func (us UpspinServer) displayStateHandle(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("upspin").Parse(HtmlPage))
	tmpl.Execute(w, nil)
}

func (us UpspinServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", us.displayStateHandle).Methods("GET")
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
