// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"encoding/json"
	"fmt"
)

type Printable interface {
	LinkJSON | []LinkJSON | VrfJSON | []VrfJSON | NeighJSON | []NeighJSON | RouteJSON | []RouteJSON | Tunnel | []Tunnel | Tuntap | []Tuntap
}

func printJSON[T Printable](cmd cmd, data T) error {
	var jsonData []byte
	var err error

	if cmd.Opts.Prettify {
		jsonData, err = json.MarshalIndent(data, "", "    ") // Use 4 spaces for indentation
	} else {
		jsonData, err = json.Marshal(data)
	}
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %w", err)
	}

	_, err = cmd.Out.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing JSON data to writer: %w", err)
	}

	return nil
}
