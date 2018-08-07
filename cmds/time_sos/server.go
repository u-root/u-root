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
    function sendSet() {
      fetch("http://localhost:{{.Port}}/auto", {
        method: 'Post'
      })
      .then(r => r.json())
      .then( s => {
        if (s !== null) {
          alert(s.Error);
        }
      })
      .catch(err => alert(err))
    }
  </script>
  </head>
  <body>
    <h1>Time</h1>
    <table style="width:100%">
      <tr>
        <td><input type="submit" id="button" class="btn" label="Set Time" onclick=sendSet()></td>
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

func (ts *TimeServer) buttonHandle(w http.ResponseWriter, r *http.Request) {
	if err := ts.service.AutoSetTime(); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{fmt.Sprintf("Unable to set time. Are you online?")})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

func (ts *TimeServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", ts.displayStateHandle).Methods("GET")
	r.HandleFunc("/auto", ts.buttonHandle).Methods("POST")
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
