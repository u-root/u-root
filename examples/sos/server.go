// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/u-root/u-root/pkg/sos"
)

const (
	// it's ugly, but we have to define a basic HTML string to fall back on if our .html is
	// missing. If you are implementing a simple service, you may be able to get away with
	// only this. example_sos.html is exactly the same as this string, except we use divs
	// and the CSS stylesheet instead of tables, with more documentation.
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
      function sendString() {
        ex = document.getElementById("field").value
        fetch("http://localhost:{{.Port}}/your_url_here_1", {
          method: 'Post',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            Example: ex
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

      function sendEmpty() {
        fetch("http://localhost:{{.Port}}/your_url_here_2", {
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
    </script>
    </head>
    <body>
      {{$example := .Example}}
      <h1>Example</h1>
      <table style="width:100%">
        <tr>
          <td><input type="text" id="field"  class="text" value="{{$example}}"></td>
          <td><input type="submit" class="btn" value="Submit" onclick=sendString()></td>
          <td><input type="submit" class="btn" value="Clear" onclick=sendEmpty()></td>
        </tr>
      </table>
    </body>
  `
)

var Port uint

// our server contains an instance of our service.
type ExampleServer struct {
	service *ExampleService
}

// displayStateHandle renders our webpage with the data from our service
// * update your service (if required).
// * copy your service's fields to a data struct, with the port we got from SoS.
// * load the HTML file. If it doesn't exist, use the HTML string defined above.
// * Render the HTML with our data.
func (es *ExampleServer) displayStateHandle(w http.ResponseWriter, r *http.Request) {
	// if your service has an update function, call it here to refresh its fields.
	// es.service.Update()
	exampleData := struct {
		Example string
		Port    uint
	}{es.service.Example, Port}
	var tmpl *template.Template
	file, err := ioutil.ReadFile(sos.HTMLPath("example.html"))
	if err == nil {
		html := string(file)
		tmpl = template.Must(template.New("SoS").Parse(html))
	} else {
		tmpl = template.Must(template.New("SoS").Parse(DefHtmlPage))
	}
	tmpl.Execute(w, exampleData)
}

// define a JsonMsg struct for easily passing around our data. Though it isn't
// necessary in this example, once you have multiple values associated with your
// service it makes things much easier. This JSON must be defined exactly as it is
// in the HTML file.
type ExampleJsonMsg struct {
	Example string
}

// this is a basic JSON handler function, dealing with string passed back as JSON.
// * decode the returned JSON message and store the fields in our JsonMsg struct,
//   with proper error checking.
// * call the associated service function, passing in this JsonMsg, with proper error checking.
func (es *ExampleServer) yourHandleFunc1(w http.ResponseWriter, r *http.Request) {
	var msg ExampleJsonMsg
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	if err := es.service.ExampleServiceFunc1(msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

// if your service function doesn't require input, i.e. a button press, you can omit the
// JSON decoding step.
func (es *ExampleServer) yourHandleFunc2(w http.ResponseWriter, r *http.Request) {
	if err := es.service.ExampleServiceFunc2(); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

// define our router
// * the "/" handler is required, as it is what updates our service and renders the webpage.
// * if you plan on using the Material Design CSS in pkg/sos/html/css, you need to define
//   the pathPrefix for it.
// * All other handlers are user-defined.
func (es *ExampleServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", es.displayStateHandle).Methods("GET")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(sos.HTMLPath("css")))))
	r.HandleFunc("/your_url_here_1", es.yourHandleFunc1).Methods("POST")
	r.HandleFunc("/your_url_here_2", es.yourHandleFunc2).Methods("POST")
	return r
}

// start our server
// * get a listener and an open port from the SoS
// * build our router (see above) and add our service to the SoS table, with the label "example"
func (es *ExampleServer) Start() {
	listener, port, err := sos.GetListener()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	Port = port
	sos.StartServiceServer(es.buildRouter(), "example", listener, Port)
}

// build a new server with the service passed in.
func NewExampleServer(service *ExampleService) *ExampleServer {
	return &ExampleServer{
		service: service,
	}
}
