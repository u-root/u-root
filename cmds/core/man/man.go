// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// man - print manual entry for command.
//
// Synopsis:
//
//	man COMMAND
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/cmds/core/man/data"
)

//go:generate go run gen/gen.go ../../../cmds data/data.go

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: man COMMAND")
	}
	mans := make(map[string]string)
	if err := json.Unmarshal([]byte(data.Data), &mans); err != nil {
		log.Fatal(err)
	}
	cmd := os.Args[1]
	man, ok := mans[cmd]
	if !ok {
		log.Fatalf("No manual entry for %q", cmd)
	}
	fmt.Println(man)
}
