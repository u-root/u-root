// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"

	"github.com/u-root/u-root/pkg/flash"
	"github.com/u-root/u-root/pkg/flash/spimock"
)

type dummyProgrammer struct {
	*flash.Flash
	spi *spimock.MockSPI
}

func (p *dummyProgrammer) Close() error {
	return p.spi.Close()
}

func init() {
	supportedProgrammers["dummy"] = func(params programmerParams) (programmer, error) {
		spi := spimock.New()
		if image, ok := params["image"]; ok {
			var err error
			spi, err = spimock.NewFromFile(image)
			if err != nil {
				return nil, err
			}
			delete(params, "image")
		}
		if len(params) != 0 {
			return nil, fmt.Errorf("unrecognized parameters: %v", params)
		}
		flash, err := flash.New(spi)
		if err != nil {
			return nil, err
		}
		return &dummyProgrammer{
			Flash: flash,
			spi:   spi,
		}, nil
	}
}
