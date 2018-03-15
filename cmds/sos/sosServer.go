// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PortNum  = "8000"
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
</style>
<script>
</script>
</head>
<body>
<h1>Current Services</h1> 
<table style="width:100%">
	<tr>
    	<th>Service</th>
    	<th>Port Number</th>
    	<th></th>
  	</tr>
	{{range $key, $value := .}}
	<tr>
	<td>{{$key}}</td>
	<td>{{$value}}</td>
	<td><a href="localhost:{{$value}}" target="_blank">Go there!</td>
	</tr>
	{{else}}
	<tr>
		<td colspan="3" style="text-align:center">No services</td>
	</tr>
	{{end}}
</table>
</body>
`
)

type RegisterReqJson struct {
	Service string
	Port    uint
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
	var msg RegisterReqJson
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}

	if err := register(msg.Service, msg.Port); err != nil {
		fmt.Printf("error: %v", err)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

type UnRegisterReqJson struct {
	ServiceName string
}

func unregisterHandle(w http.ResponseWriter, r *http.Request) {
	var msg UnRegisterReqJson
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	unregister(msg.ServiceName)
	json.NewEncoder(w).Encode(nil)
}

type GetServiceResJson struct {
	Port uint
}

func getServiceHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	port, err := read(vars["service"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(GetServiceResJson{port})
}

func redirectToResourceHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	port, err := read(vars["service"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("localhost:%v/", port), http.StatusPermanentRedirect)
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buildHtmlPage(w)
	}).Methods("GET")
	r.HandleFunc("/register", registerHandle).Methods("POST")
	r.HandleFunc("/unregister", unregisterHandle).Methods("POST")
	r.HandleFunc("/service/{service}", getServiceHandle).Methods("GET")
	r.HandleFunc("/go/{service}", redirectToResourceHandle).Methods("GET")
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", PortNum), r))
}

func buildHtmlPage(wr io.Writer) error {
	s := snapshotRegistry()
	tmpl := template.Must(template.New("SoS").Parse(HtmlPage))
	return tmpl.Execute(wr, s)
}
