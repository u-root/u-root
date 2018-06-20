// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/u-root/u-root/pkg/sos"
	"github.com/u-root/u-root/pkg/wifi"
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
function sendConnect(elem, index) {
	replaceWithConnecting(elem);
	disableOtherButtons(elem);
	essid = document.getElementById("essid".concat(index)).innerHTML
	pass = document.getElementById("pass".concat(index)) ? 
		document.getElementById("pass".concat(index)).value : ""
	id = document.getElementById("id".concat(index)) ? 
		document.getElementById("id".concat(index)).value : ""
	fetch("http://localhost:{{.Port}}/connect", {
		method: 'Post',
		headers: {
			'Accept': 'application/json',
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			Essid: essid,
			Pass: pass,
			Id: id
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

function replaceWithConnecting(elem) {
    connectingTxt = document.createTextNode("Connecting...");
    elem.style.display = "none";
    elem.parentNode.appendChild(connectingTxt);
}

function sendRefresh(elem) {
	elem.setAttribute("disabled", "true");
	elem.setAttribute("value","Refreshing");
	disableOtherButtons(elem);
	fetch("http://localhost:{{.Port}}/refresh", {
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

function disableOtherButtons(elem) {
    btns = document.getElementsByClassName("btn");
    for (let btn of btns) {
    	if (btn === elem) {
    		continue;
    	}
    	btn.setAttribute("disabled", "true");
    }	
}
</script>
</head>
<body>
{{$NoEnc := 0}}
{{$WpaPsk := 1}}
{{$WpaEap := 2}}
{{$connectedEssid := .ConnectedEssid}}
{{$connectingEssid := .ConnectingEssid}}
<h1>Please choose your Wifi</h1> 
<table style="width:100%">
	<tr>
    	<th>Essid</th>
    	<th>Identity</th>
    	<th>Password / Passphrase</th>
    	<th><input type="submit" class="btn" onclick=sendRefresh(this) value="Refresh"></th>
  	</tr>
	{{range $idx, $opt := .WifiOpts}}
		{{if eq $opt.AuthSuite $NoEnc}}
			<tr>
    			<td id="essid{{$idx}}">{{$opt.Essid}}</td>
    			<td></td>
    			<td></td>
    			{{if and (eq $connectedEssid $opt.Essid) (ne $connectedEssid "")}}
    				<td>Connected</td>
				{{else if and (and (eq $connectingEssid $opt.Essid) (ne $connectingEssid "")) (ne $connectingEssid $connectedEssid) }}
    				<td>Connecting...</td>
    			{{else}}
    				<td><input type="submit" class="btn" onclick="sendConnect(this, {{$idx}})" value="Connect"></td>
    			{{end}}
  			</tr>
		{{else if eq $opt.AuthSuite $WpaPsk}}
			<tr>
    			<td id="essid{{$idx}}">{{$opt.Essid}}</td>
    			<td></td>
    			<td><input type="password" id="pass{{$idx}}"></td>
    			{{if and (eq $connectedEssid $opt.Essid) (ne $connectedEssid "")}}
    				<td>Connected</td>
				{{else if and (and (eq $connectingEssid $opt.Essid) (ne $connectingEssid "")) (ne $connectingEssid $connectedEssid) }}
    				<td>Connecting...</td>
    			{{else}}
    				<td><input type="submit" class="btn" onclick="sendConnect(this, {{$idx}})" value="Connect"></td>
    			{{end}}
       		</tr>
		{{else if eq $opt.AuthSuite $WpaEap}}
			<tr>
    			<td id="essid{{$idx}}">{{$opt.Essid}}</td>
    			<td><input type="text" id="id{{$idx}}"></td>
    			<td><input type="password" id="pass{{$idx}}"></td>
    			{{if and (eq $connectedEssid $opt.Essid) (ne $connectedEssid "")}}
    				<td>Connected</td>
				{{else if and (and (eq $connectingEssid $opt.Essid) (ne $connectingEssid "")) (ne $connectingEssid $connectedEssid) }}
    				<td>Connecting...</td>
    			{{else}}
    				<td><input type="submit" class="btn" onclick="sendConnect(this, {{$idx}})" value="Connect"></td>
    			{{end}}
  			</tr>
		{{else}}
			<tr>
    			<td id="essid{{$idx}}">{{$opt.Essid}}</td>
    			<td colspan="3">Not a supported protocol</td>
  			</tr>
		{{end}}
	{{else}}
		<td colspan="4">No essids found</td>
	{{end}}
</table>

{{if and (ne $connectingEssid "") (ne $connectingEssid $connectedEssid) }}
	<script>disableOtherButtons(null)</script>
{{end}}
</body>
`
)

var (
	Port uint
)

type WifiServer struct {
	service *WifiService
}

func userInputValidation(essid, pass, id string) ([]string, error) {
	switch {
	case essid != "" && pass != "" && id != "":
		return []string{essid, pass, id}, nil
	case essid != "" && pass != "" && id == "":
		return []string{essid, pass}, nil
	case essid != "" && pass == "" && id == "":
		return []string{essid}, nil
	default:
		return nil, fmt.Errorf("Invalid user input")
	}
}

func (ws WifiServer) refreshHandle(w http.ResponseWriter, r *http.Request) {
	if err := ws.service.Refresh(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(nil)
}

type ConnectJsonMsg struct {
	Essid string
	Pass  string
	Id    string
}

func (ws WifiServer) connectHandle(w http.ResponseWriter, r *http.Request) {
	var msg ConnectJsonMsg
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	a, err := userInputValidation(msg.Essid, msg.Pass, msg.Id)
	if err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}

	if err := ws.service.Connect(a); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		return
	}
	// Connect Successful
	json.NewEncoder(w).Encode(nil)
}

func (ws WifiServer) displayStateHandle(w http.ResponseWriter, r *http.Request) {
	s := ws.service.GetState()
	displayWifi(w, s.NearbyWifis, s.CurEssid, s.ConnectingEssid)
}

func (ws WifiServer) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", ws.displayStateHandle).Methods("GET")
	r.HandleFunc("/refresh", ws.refreshHandle).Methods("POST")
	r.HandleFunc("/connect", ws.connectHandle).Methods("POST")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("/etc/sos/html/css"))))
	return r
}

func (ws WifiServer) Start() {
	defer ws.service.Shutdown()
	listener, port, err := sos.GetListener()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	Port = port
	fmt.Println(sos.StartServiceServer(ws.buildRouter(), "wifi", listener, Port))
}

func displayWifi(wr io.Writer, wifiOpts []wifi.Option, connectedEssid, connectingEssid string) error {
	wifiData := struct {
		WifiOpts        []wifi.Option
		ConnectedEssid  string
		ConnectingEssid string
		Port            uint
	}{wifiOpts, connectedEssid, connectingEssid, Port}

	var tmpl *template.Template
	file, err := ioutil.ReadFile("/etc/sos/html/wifi.html")
	if err == nil {
		html := string(file)
		tmpl = template.Must(template.New("name").Parse(html))
	} else {
		tmpl = template.Must(template.New("name").Parse(DefHtmlPage))
	}
	return tmpl.Execute(wr, wifiData)
}

func NewWifiServer(service *WifiService) *WifiServer {
	return &WifiServer{
		service: service,
	}
}
