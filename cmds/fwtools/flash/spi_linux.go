// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/u-root/u-root/pkg/flash"
	"github.com/u-root/u-root/pkg/spidev"
)

type spidevProgrammer struct {
	*flash.Flash
	spi *spidev.SPI
}

func (p *spidevProgrammer) Close() error {
	return p.spi.Close()
}

func init() {
	supportedProgrammers["linux_spi"] = func(params programmerParams) (programmer, error) {
		dev, ok := params["dev"]
		if !ok {
			return nil, fmt.Errorf("dev is a required parameter for linux_spi")
		}
		delete(params, "dev")

		spi, err := spidev.Open(dev, spidev.WithLogger(log.Printf))
		if err != nil {
			return nil, err
		}

		if spiSpeed, ok := params["spispeed"]; ok {
			spiSpeedKHz, err := strconv.Atoi(spiSpeed)
			if err != nil {
				return nil, err
			}
			if spiSpeedKHz > math.MaxUint32/1000 {
				return nil, fmt.Errorf("spispeed is larger than max uint32")
			}
			if err := spi.SetSpeedHz(uint32(spiSpeedKHz) * 1000); err != nil {
				return nil, err
			}
			delete(params, "spispeed")
		}

		if len(params) != 0 {
			return nil, fmt.Errorf("unrecognized parameters: %v", params)
		}
		f, err := flash.New(spi)
		if err != nil {
			return nil, err
		}
		return &spidevProgrammer{
			Flash: f,
			spi:   spi,
		}, nil
	}
}
