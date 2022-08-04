// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config manages configuratino settings for secure launch.
package config

// Config contains the configuration for secure launch u-root.
type Config struct {
	Collectors     bool `json:"collectors"`
	EventLog       bool `json:"eventlog"`
	Measurements   bool `json:"measurements"`
	MeasurementPCR int  `json:"measurement pcr"`
}

// New returns a new default config.
func New() Config {
	var config Config

	config.Collectors = true
	config.EventLog = true
	config.Measurements = true
	config.MeasurementPCR = 22

	return config
}

// Conf holds the global configuration settings. It is initialized with the
// default settings but can be overridden by explicitly setting some options in
// the policy file or by leaving sections of the policy file out.
var Conf = New()
