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
	  function sendEdit() {
	    fetch("http://localhost:{{.Port}}/edit", {
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
	  function sendSubmit() {
	    username = document.getElementById("user").value
	    dirserver = document.getElementById("dir").value
	    storeserver = document.getElementById("store").value
	    secretseed = document.getElementById("seed").value
	    fetch("http://localhost:{{.Port}}/submit", {
	      method: 'Post',
	      headers: {
	  			'Accept': 'application/json',
	  			'Content-Type': 'application/json'
	  		},
	  		body: JSON.stringify({
	  			User:  username,
	  			Dir:   dirserver,
	        Store: storeserver,
	        Seed:  secretseed
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

	  function setFieldsOnLoad(config) {
	    if (config) {
	      disableFields();
	    }
	    else {
	      enableFields();
	    }
	  }

	  function enableFields() {
	    fields = document.getElementsByClassName("text");
	    for(let field of fields) {
	      field.removeAttribute("disabled")
	    }
	    document.getElementById("button").setAttribute("value", "Submit");
	    document.getElementById("button").setAttribute("onclick", "sendSubmit()");
	  }
	  function disableFields() {
	    fields = document.getElementsByClassName("text");
	    for(let field of fields) {
	      field.setAttribute("disabled", "true")
	    }
	    document.getElementById("button").setAttribute("value", "Edit");
	    document.getElementById("button").setAttribute("onclick", "sendEdit()");
	  }

	</script>
	</head>
	<body onload="setFieldsOnLoad({{.Configured}})">
	  <!-- Copy to local variables -->
	  {{$user := .User}}
	  {{$dir := .Dir}}
	  {{$store := .Store}}
	  {{$seed := .Seed}}
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
	      <td><input type="text" id="user"  class="text" value="{{$user}}"></td>
	      <td><input type="text" id="dir"   class="text" value="{{$dir}}"></td>
	      <td><input type="text" id="store" class="text" value="{{$store}}"></td>
	      <td><input type="text" id="seed"  class="text" value="{{$seed}}"></td>
	      <td><input type="submit" id="button" class="btn"></td>
	    </tr>
	  </table>
	</body>
`
)

type UpspinServer struct {
	service *UpspinService
}

var (
	Port uint
)

func (us *UpspinServer) editHandle(w http.ResponseWriter, r *http.Request) {
	us.service.ToggleFlag()
	json.NewEncoder(w).Encode(nil)
}

type UpspinAcctJsonMsg struct {
	User  string
	Dir   string
	Store string
	Seed  string
}

func (us *UpspinServer) submitHandle(w http.ResponseWriter, r *http.Request) {
	us.service.ToggleFlag()
	var msg UpspinAcctJsonMsg
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	if err := us.service.SetConfig(msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

func (us *UpspinServer) displayStateHandle(w http.ResponseWriter, r *http.Request) {
	us.service.Update()
	upspinData := struct {
		Configured bool
		User       string
		Dir        string
		Store      string
		Seed       string
		Port       uint
	}{us.service.Configured, us.service.User, us.service.Dir, us.service.Store, us.service.Seed, Port}
	var tmpl *template.Template
	file, err := ioutil.ReadFile(sos.HTMLPath("upspin.html"))
	if err == nil {
		html := string(file)
		tmpl = template.Must(template.New("SoS").Parse(html))
	} else {
		tmpl = template.Must(template.New("SoS").Parse(DefHtmlPage))
	}
	tmpl.Execute(w, upspinData)
}

func (us *UpspinServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", us.displayStateHandle).Methods("GET")
	r.HandleFunc("/edit", us.editHandle).Methods("POST")
	r.HandleFunc("/submit", us.submitHandle).Methods("POST")
	sos.SetHTMLPath([]string{"/etc/sos/html"})
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(sos.HTMLPath("css")))))
	return r
}

func (us *UpspinServer) Start() {
	listener, port, err := sos.GetListener()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	Port = port
	fmt.Println(sos.StartServiceServer(us.buildRouter(), "upspin", listener, Port))
}

func NewUpspinServer(service *UpspinService) *UpspinServer {
	return &UpspinServer{
		service: service,
	}
}
