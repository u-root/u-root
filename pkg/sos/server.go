// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sos

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	PortNum     = "8000"
	DefHtmlPage = `
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
<h1>Current Services (html embedded)</h1>
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
	<td><a href="http://localhost:{{$value}}" target="_blank">Go there!</td>
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

// default path
var htmlRoot = "/etc/sos/html"

type SosServer struct {
	service *SosService
}

type RegisterReqJson struct {
	Service string
	Port    uint
}

// set htmlRoot var to array of dirs p
func SetHTMLRoot(p ...string) {
	if len(p) > 0 {
		htmlRoot = filepath.Join(p...)
	}
}

// HTMLPath returns the HTMLPath formed by joining the arguments together.
// If there are no arguments, it simply returns the HTML root directory.
func HTMLPath(n ...string) string {
	return filepath.Join(htmlRoot, filepath.Join(n...))
}

func (s SosServer) registerHandle(w http.ResponseWriter, r *http.Request) {
	var msg RegisterReqJson
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}

	if err := s.service.Register(msg.Service, msg.Port); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

type UnRegisterReqJson struct {
	ServiceName string
}

func (s SosServer) unregisterHandle(w http.ResponseWriter, r *http.Request) {
	var msg UnRegisterReqJson
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	s.service.Unregister(msg.ServiceName)
	json.NewEncoder(w).Encode(nil)
}

type GetServiceResJson struct {
	Port uint
}

func (s SosServer) getServiceHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	port, err := s.service.Read(vars["service"])
	if err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(GetServiceResJson{port})
}

func (s SosServer) redirectToResourceHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	port, err := s.service.Read(vars["service"])
	if err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("http://localhost:%v/", port), http.StatusTemporaryRedirect)
}

func (s SosServer) displaySosHandle(w http.ResponseWriter, r *http.Request) {
	snap := s.service.SnapshotRegistry()
	var tmpl *template.Template
	file, err := ioutil.ReadFile(HTMLPath("sos.html"))
	if err == nil {
		html := string(file)
		tmpl = template.Must(template.New("SoS").Parse(html))
	} else {
		tmpl = template.Must(template.New("SoS").Parse(DefHtmlPage))
	}
	tmpl.Execute(w, snap)
}

func (s SosServer) buildRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", s.displaySosHandle).Methods("GET")
	r.HandleFunc("/register", s.registerHandle).Methods("POST")
	r.HandleFunc("/unregister", s.unregisterHandle).Methods("POST")
	r.HandleFunc("/service/{service}", s.getServiceHandle).Methods("GET")
	r.HandleFunc("/go/{service}", s.redirectToResourceHandle).Methods("GET")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(HTMLPath("css")))))
	return r
}

func StartServer(service *SosService) {
	server := SosServer{service}
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", PortNum), server.buildRouter()))
}
