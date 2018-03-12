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

type RegisterJson struct {
	ServiceName string
	PortNumber  uint
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
	var msg RegisterJson
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	if port, exists := read(msg.ServiceName); exists {
		err := fmt.Sprintf("error: %v already exists at %v", msg.ServiceName, port)
		fmt.Println(err)
		json.NewEncoder(w).Encode(struct{ Error string }{err})
		return
	}
	register(msg.ServiceName, msg.PortNumber)
	json.NewEncoder(w).Encode(nil)
}

type UnregisterJson struct {
	ServiceName string
}

func unregisterHandle(w http.ResponseWriter, r *http.Request) {
	var msg UnregisterJson
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

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buildHtmlPage(w)
	})
	http.HandleFunc("/register", registerHandle)
	http.HandleFunc("/unregister", unregisterHandle)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", PortNum), nil))
}

func buildHtmlPage(wr io.Writer) error {
	s := snapshotRegistry()
	tmpl := template.Must(template.New("SoS").Parse(HtmlPage))
	return tmpl.Execute(wr, s)
}
