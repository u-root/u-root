// Copyright 2013-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Get the time the machine has been up
// Synopsis:
//     uptime
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// loadavg takes in the contents of proc/loadavg,it then extracts and returns the three load averages as a string
func loadavg(contents string) (loadaverage string, err error) {
	loadavg := strings.Fields(contents)
	if len(loadavg) < 3 {
		return "", fmt.Errorf("error:invalid contents:the contents of proc/loadavg we are trying to process contain less than the required 3 loadavgs")
	}
	return loadavg[0] + "," + loadavg[1] + "," + loadavg[2], nil
}

// uptime takes in the contents of proc/uptime it then extracts and returns the uptime in the format Days , Hours , Minutes ,Seconds
func uptime(contents string) (*time.Time, error) {
	uptimeArray := strings.Fields(contents)
	if len(uptimeArray) == 0 {
		return nil, errors.New("error:the contents of proc/uptime we are trying to read are empty")
	}
	uptimeDuration, err := time.ParseDuration(string(uptimeArray[0]) + "s")
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	uptime := time.Time{}.Add(uptimeDuration)

	return &uptime, nil
}

func main() {
	procUptimeOutput, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		log.Fatalf("error reading /proc/uptime: %v \n", err)
	}
	uptimeTime, err := uptime(string(procUptimeOutput))
	if err != nil {
		log.Fatal(err)
	}
	procLoadAvgOutput, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		log.Fatalf("error reading /proc/loadavg: %v \n", err)
	}
	loadAverage, err := loadavg(string(procLoadAvgOutput))
	if err != nil {
		log.Fatal(err)
	}
	//Subtracted one from time.Day() because time.Add(Duration) starts counting at 1 day instead of zero days.
	fmt.Printf(" %s up %d days, %d hours , %d min ,loadaverage: %s \n", time.Now().Format("15:04:05"), (uptimeTime.Day() - 1), uptimeTime.Hour(), uptimeTime.Minute(), loadAverage)
}
