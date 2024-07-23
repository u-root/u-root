// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

type Printable interface {
	Link | []Link | Vrf | []Vrf | Neigh | []Neigh | Route | []Route | Tunnel | []Tunnel | Tuntap | []Tuntap
}

func printJSON[T Printable](cmd cmd, data T) error {
	var jsonData []byte
	var err error

	if cmd.opts.prettify {
		jsonData, err = json.MarshalIndent(data, "", "    ") // Use 4 spaces for indentation
	} else {
		jsonData, err = json.Marshal(data)
	}
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %v", err)
	}

	_, err = cmd.out.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing JSON data to writer: %v", err)
	}

	return nil
}
