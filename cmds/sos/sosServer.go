// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"io"
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

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buildHtmlPage(w)
	})

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", PortNum), nil))
}

func buildHtmlPage(wr io.Writer) error {
	m := make(map[string]int)
	m["a"] = 8080
	tmpl := template.Must(template.New("SoS").Parse(HtmlPage))
	return tmpl.Execute(wr, m)
}
