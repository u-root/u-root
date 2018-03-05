// Copyright 2017 the u-root Authors. All rights reserved
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
	PortNum  = "8080"
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

.essid {
	border-width: 0;
}
</style>
<script>
function replaceWithConnecting(elem) {
    connectingTxt = document.createTextNode("Connecting...");
    elem.style.display = "none";
    elem.parentNode.appendChild(connectingTxt);
	disableOtherButtons(elem);
}

function sendRefresh(elem) {
	elem.setAttribute("disabled", "true");
	elem.setAttribute("value","Refreshing");
	disableOtherButtons(elem);
	fetch("http://localhost:8080/refresh")
	.then(r => r.json())
	.then( s => {
		if (s !== null) {
			alert(s.Error);
		}
		else {
			window.location.reload();
		}
	})
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
<h1 style="float:left">Please choose your Wifi</h1> 
<table style="width:100%">
	<tr>
    	<th>Essid</th>
    	<th>Identity</th>
    	<th>Password / Passphrase</th>
    	<th><input type="submit" class="btn" onclick=sendRefresh(this) value="Refresh"></button></th>
  	</tr>
	{{range $idx, $opt := .WifiOpts}}
		<form id="f{{$idx}}" method="post"></form>
		{{if eq $opt.AuthSuite $NoEnc}}
			<tr>
    			<td><input type="text" name="essid" class="essid" form="f{{$idx}}" readonly value={{$opt.Essid}}></td>
    			<td></td>
    			<td></td>
    			{{if and (eq $connectedEssid $opt.Essid) (ne $connectedEssid "")}}
    				<td>Connected</td>
				{{else}}
    				<td><input type="submit" class="btn" form="f{{$idx}}" onclick=replaceWithConnecting(this) value="Connect"></td>
    			{{end}}
  			</tr>
		{{else if eq $opt.AuthSuite $WpaPsk}}
			<tr>
    			<td><input type="text" name="essid" class="essid" form="f{{$idx}}" readonly value={{$opt.Essid}}></td>
    			<td></td>
    			<td><input type="password" name="pass" form="f{{$idx}}"></td>
    			{{if and (eq $connectedEssid $opt.Essid) (ne $connectedEssid "")}}
    				<td>Connected</td>
				{{else}}
    				<td><input type="submit" class="btn" form="f{{$idx}}" onclick=replaceWithConnecting(this) value="Connect"></td>
    			{{end}}
  			</tr>
		{{else if eq $opt.AuthSuite $WpaEap}}
			<tr>
    			<td><input type="text" name="essid" class="essid" form="f{{$idx}}" readonly value={{$opt.Essid}}></td>
    			<td><input type="text" name="identity" form="f{{$idx}}"></td>
				<td><input type="password" name="pass" form="f{{$idx}}"></td>
    			{{if and (eq $connectedEssid $opt.Essid) (ne $connectedEssid "")}}
    				<td>Connected</td>
				{{else}}
    				<td><input type="submit" class="btn" form="f{{$idx}}" onclick=replaceWithConnecting(this) value="Connect"></td>
    			{{end}}
  			</tr>
		{{else}}
			<tr>
    			<td><input type="text" name="essid" class="essid" form="f{{$idx}}" readonly value={{$opt.Essid}}></td>
    			<td colspan="3">Not a supported protocol</td>
  			</tr>
		{{end}}
	{{else}}
		<td colspan="4">No essids found</td>
	{{end}}
</table>
</body>
`
)

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

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := getState()

		if r.Method != http.MethodPost {
			err := displayWifi(w, s.nearbyWifis, s.curEssid)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			return
		}

		a, err := userInputValidation(r.FormValue("essid"), r.FormValue("pass"), r.FormValue("identity"))
		if err != nil {
			// TODO: Need proper error handling
			log.Printf("error: %v", err)
			return
		}

		UserInputChannel <- UserInputMessage{args: a}
		sMsg := <-StateChannel
		displayWifi(w, sMsg.nearbyWifis, sMsg.curEssid)
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if err := scanWifi(); err != nil {
			json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
			return
		}
		json.NewEncoder(w).Encode(nil)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", PortNum), nil)
}

func displayWifi(wr io.Writer, wifiOpts []WifiOption, connectedEssid string) error {
	wifiData := struct {
		WifiOpts       []WifiOption
		ConnectedEssid string
	}{wifiOpts, connectedEssid}

	tmpl := template.Must(template.New("name").Parse(HtmlPage))

	return tmpl.Execute(wr, wifiData)
}
