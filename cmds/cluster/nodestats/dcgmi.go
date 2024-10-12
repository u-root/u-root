// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/u-root/u-root/pkg/cluster/health"
)

func dcgmi() ([]health.NvidiaHealth, error) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		return nil, err

	}

	var all []health.NvidiaHealth
	var errs error
	for _, gpu := range gpus {
		hc, err := dcgm.HealthCheckByGpuId(gpu)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		h := health.NvidiaHealth{DeviceHealth: dcgm.DeviceHealth{GPU: gpu, Status: hc.Status}}
		s, err := dcgm.GetDeviceStatus(gpu)
		if err != nil {
			errs = errors.Join(errs, err)
		}
		h.DeviceStatus = s

		all = append(all, h)
	}
	return all, errs
}
