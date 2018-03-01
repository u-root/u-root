// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type SecProto int

const (
	NoEnc SecProto = iota
	WpaPsk
	WpaEap
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
    connectingTxt = document.createTextNode("Connecting...")
    elem.style.display = "none"
    elem.parentNode.appendChild(connectingTxt);
    btns = document.getElementsByClassName("btn");
    var i;
    for (i = 0; i < btns.length; i++) {
    	if (btns[i] === elem) {
    		continue;
    	}
    	btns[i].setAttribute("disabled", "true")
    }
}
</script>
</head>
<body>
{{$NoEnc := 0}}
{{$WpaPsk := 1}}
{{$WpaEap := 2}}
{{$connectedEssid := .ConnectedEssid}}
<h1>Please choose your Wifi</h1>
<table style="width:100%">
	<tr>
    	<th>Essid</th>
    	<th>Identity</th>
    	<th>Password / Passphrase</th>
    	<th></th>
  	</tr>
	{{range $idx, $opt := .WifiOpts}}
		<form id="f{{$idx}}" method="post"></form>
		{{if eq $opt.AuthSuite $NoEnc}}
			<tr>
    			<td><input type="text" name="essid" class="essid" form="f{{$idx}}" readonly value={{$opt.Essid}}></td>
    			<td></td>
    			<td></td>
    			{{if eq $connectedEssid $opt.Essid}}
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
    			{{if eq $connectedEssid $opt.Essid}}
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
				{{if eq $connectedEssid $opt.Essid}}
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

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		stubWifis := []WifiOptions{
			{"stub1", NoEnc},
			{"stub2", WpaPsk},
			{"stub3", WpaEap},
			{"stub4", 123},
		}

		if r.Method != http.MethodPost {
			err := displayWifi(w, stubWifis, "")
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			return
		}

		var a []string
		// user input validation
		switch {
		case r.FormValue("essid") != "" && r.FormValue("pass") != "" && r.FormValue("identity") != "":
			a = []string{r.FormValue("essid"), r.FormValue("pass"), r.FormValue("identity")}
		case r.FormValue("essid") != "" && r.FormValue("pass") != "" && r.FormValue("identity") == "":
			a = []string{r.FormValue("essid"), r.FormValue("pass")}
		case r.FormValue("essid") != "" && r.FormValue("pass") == "" && r.FormValue("identity") == "":
			a = []string{r.FormValue("essid")}
		default: // TODO: Error handling: for now, just exit
			os.Exit(1)
		}

		UserInputChannel <- UserInputMessage{args: a}
		statMsg := <-StatusChannel
		displayWifi(w, stubWifis, statMsg.essid)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", PortNum), nil)
}

type WifiOptions struct {
	Essid     string
	AuthSuite SecProto
}

func displayWifi(wr io.Writer, wifiOpts []WifiOptions, connectedEssid string) error {
	wifiData := struct {
		WifiOpts       []WifiOptions
		ConnectedEssid string
	}{wifiOpts, connectedEssid}

	tmpl := template.Must(template.New("name").Parse(HtmlPage))

	return tmpl.Execute(wr, wifiData)
}
