// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

func main() {
	// Build a new example service
	service, err := NewExampleService()
	if err != nil {
		log.Fatal(err)
	}
	// Build a new example server with this service, then start it
	NewExampleServer(service).Start()
}
