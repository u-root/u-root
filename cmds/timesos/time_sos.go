// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

func main() {
	service, err := NewTimeService()
	if err != nil {
		log.Fatal(err)
	}
	NewTimeServer(service).Start()
}
