// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/u-root/u-root/pkg/sos"
)

const (
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
    input {
      font-size: 120%;
    }
  </style>
  <script>
		function sendAutoSet() {
			fetch("http://localhost:{{.Port}}/auto", {
				method: 'Post'
			})
			.then(r => r.json())
			.then( s => {
				if (s !== null) {
					alert(s.Error);
					window.location.reload();
				}
				else {
					window.location.reload();
				}
			})
			.catch(err => alert(err))
		}

		function sendManSet() {
			d = document.getElementById("date_field").value
			t = document.getElementById("time_field").value
			fetch("http://localhost:{{.Port}}/manual", {
				method: 'Post',
				headers: {
					'Accept': 'application/json',
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					Date: d,
					Time: t
				})
			})
			.then(r => r.json())
			.then( s => {
				if (s !== null) {
					alert(s.Error);
					window.location.reload();
				}
				else {
					window.location.reload();
				}
			})
			.catch(err => alert(err))
		}

		function setOnLoad(date, time) {
			document.getElementById("date_field").setAttribute("value", date)
			document.getElementById("time_field").setAttribute("value", time)
		}
  </script>
  </head>
	<body onload="setOnLoad({{.Date}}, {{.Time}})">
		<h1>System Time Settings</h2>
		<table style="width:100%">
			<tr>
				<td><input type="date" id="date_field"></td>
				<td><input type="time" id="time_field"></td>
				<td><input type="submit" id="button" value="Auto-Set" onclick=sendAutoSet()></td>
			</tr>
			<tr>
				<td colspan="3"><input type="submit" id="button" value="Save" onclick=sendManSet()></td>
			</tr>
		</table>
	</body>
  `
)

type TimeServer struct {
	service *TimeService
}

var (
	Port uint
)

func (ts *TimeServer) displayStateHandle(w http.ResponseWriter, r *http.Request) {
	ts.service.Update()
	timeData := struct {
		Date string
		Time string
		Port uint
	}{ts.service.Date, ts.service.Time, Port}
	var tmpl *template.Template
	file, err := ioutil.ReadFile(sos.HTMLPath("time.html"))
	if err == nil {
		html := string(file)
		tmpl = template.Must(template.New("SoS").Parse(html))
	} else {
		tmpl = template.Must(template.New("SoS").Parse(DefHtmlPage))
	}
	tmpl.Execute(w, timeData)
}

func (ts *TimeServer) autoHandle(w http.ResponseWriter, r *http.Request) {
	if err := ts.service.AutoSetTime(); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{fmt.Sprintf("Unable to set time. Are you online?")})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

type TimeJsonMsg struct {
	Date string
	Time string
}

func (ts *TimeServer) manHandle(w http.ResponseWriter, r *http.Request) {
	var msg TimeJsonMsg
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	if err := ts.service.ManSetTime(msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

func (ts *TimeServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", ts.displayStateHandle).Methods("GET")
	r.HandleFunc("/auto", ts.autoHandle).Methods("POST")
	r.HandleFunc("/manual", ts.manHandle).Methods("POST")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(sos.HTMLPath("css")))))
	return r
}

// Start opens the server at localhost:{port}, where port is provided automatically
// by the SoS.
func (ts *TimeServer) Start() {
	listener, port, err := sos.GetListener()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	Port = port
	fmt.Println(sos.StartServiceServer(ts.buildRouter(), "time", listener, Port))
}

// NewTimeServer creates a server with the given TimeService.
func NewTimeServer(service *TimeService) *TimeServer {
	return &TimeServer{
		service: service,
	}
}
