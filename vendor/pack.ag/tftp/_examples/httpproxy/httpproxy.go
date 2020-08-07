// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package main

import (
	"io"
	"log"
	"net/http"

	"pack.ag/tftp"
)

const baseURL = "https://api.github.com/"

func main() {
	// Create a new server listening on port 6900, all interfaces
	server, err := tftp.NewServer(":6900")
	if err != nil {
		log.Fatal(err)
	}

	// Make proxyTFTP a ReadHandler with the ReadHandlerFunc adapter
	readHandler := tftp.ReadHandlerFunc(proxyTFTP)

	// Set the server's read handler, write requests will be rejected.
	server.ReadHandler(readHandler)

	// Start the server, if it fails error will be printed by log.Fatal
	log.Fatal(server.ListenAndServe())
}

func proxyTFTP(w tftp.ReadRequest) {
	// Append the requested path to the baseURL
	url := baseURL + w.Name()

	// Send the HTTP request
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		// This could send more specific errors, but here we'read
		// choosing to simply send "file not found"" with the error
		// message from the HTTP client back to the TFTP client.
		w.WriteError(tftp.ErrCodeFileNotFound, err.Error())
		return
	}
	defer resp.Body.Close()

	// Copy the body of the response to the TFTP client.
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}
