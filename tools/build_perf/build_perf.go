// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Measure the performance of building all the Go commands under various GOGC
// values. The output is four csv files:
//
// - build_perf_real.csv
// - build_perf_user.csv
// - build_perf_sys.csv
// - build_perf_max_rss.csv
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	cmdsPath  = "$GOPATH/src/github.com/u-root/u-root/cmds"
	gogcBegin = 50
	gogcEnd   = 2000
	gogcStep  = 50
)

// The fields profiled by a single `go build`.
type measurement struct {
	realTime float64 // seconds
	userTime float64 // seconds
	sysTime  float64 // seconds
	maxRss   int64   // KiB
}

// Each CSV file only stores one field of the measurement.
// This struct describes a single CSV file.
type csvDesc struct {
	filename string
	field    func(measurement) string
	// Measurements are sent from the `measureBuilds` to the `writeCsv`
	// functions via this channel.
	c chan []*measurement
}

var descs = []csvDesc{
	{
		"build_perf_real.csv",
		func(m measurement) string { return fmt.Sprint(m.realTime) },
		make(chan []*measurement),
	}, {
		"build_perf_user.csv",
		func(m measurement) string { return fmt.Sprint(m.userTime) },
		make(chan []*measurement),
	}, {
		"build_perf_sys.csv",
		func(m measurement) string { return fmt.Sprint(m.sysTime) },
		make(chan []*measurement),
	}, {
		"build_perf_max_rss.csv",
		func(m measurement) string { return fmt.Sprint(m.maxRss) },
		make(chan []*measurement),
	},
}

var wg sync.WaitGroup

// Return a list of command names.
func getCmdNames() ([]string, error) {
	files, err := os.ReadDir(os.ExpandEnv(cmdsPath))
	if err != nil {
		return nil, err
	}
	cmds := []string{}
	for _, file := range files {
		if file.IsDir() {
			cmds = append(cmds, file.Name())
		}
	}
	return cmds, nil
}

func buildCmd(cmd string, gogc int) (*measurement, error) {
	start := time.Now()
	c := exec.Command("go", "build")
	c.Dir = filepath.Join(os.ExpandEnv(cmdsPath), cmd)
	c.Env = append(os.Environ(), fmt.Sprintf("GOGC=%d", gogc))
	if err := c.Run(); err != nil {
		return nil, err
	}
	return &measurement{
		realTime: time.Since(start).Seconds(),
		userTime: c.ProcessState.UserTime().Seconds(),
		sysTime:  c.ProcessState.SystemTime().Seconds(),
		maxRss:   c.ProcessState.SysUsage().(*syscall.Rusage).Maxrss,
	}, nil
}

func measureBuilds(cmds []string, gogcs []int) {
	for _, cmd := range cmds {
		measurements := make([]*measurement, len(gogcs))
		for i, gogc := range gogcs {
			m, err := buildCmd(cmd, gogc)
			if err != nil {
				log.Printf("%v: %v", cmd, err)
			}
			measurements[i] = m
		}
		// Write to all csv files.
		for _, d := range descs {
			d.c <- measurements
		}
		fmt.Print(".")
	}
}

func writeCsv(cmds []string, gogcs []int, d csvDesc) {
	// Create the csv writer.
	f, err := os.Create(d.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	// Write header.
	header := make([]string, len(gogcs)+1)
	header[0] = "cmd \\ gogc"
	for i, gogc := range gogcs {
		header[i+1] = fmt.Sprint(gogc)
	}
	w.Write(header)

	// Iterator over all the measurements.
	row := 0
	for measurements := range d.c {
		record := make([]string, len(measurements)+1)
		record[0] = cmds[row]
		// Iterate over measurements for a single command.
		for i, m := range measurements {
			if m == nil {
				record[i+1] = "err"
			} else {
				record[i+1] = d.field(*m)
			}
		}

		// Write to CSV.
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatalln("error flushing csv:", err)
		}
		row++
	}
	wg.Done()
}

func main() {
	// Get list of commands.
	cmds, err := getCmdNames()
	if err != nil {
		log.Fatal("Cannot get list of commands:", err)
	}

	// Create range of GOGC values.
	gogcs := []int{}
	for i := gogcBegin; i <= gogcEnd; i += gogcStep {
		gogcs = append(gogcs, i)
	}

	wg.Add(len(descs))
	for _, d := range descs {
		go writeCsv(cmds, gogcs, d)
	}

	measureBuilds(cmds, gogcs)

	for _, d := range descs {
		close(d.c)
	}
	wg.Wait()
	fmt.Println("\nDone!")
}
