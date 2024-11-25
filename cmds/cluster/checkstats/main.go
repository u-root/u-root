// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// checkstats looks for issues in a JSON slice that could be
// important.
// Example usage: clusterstats | checkstats
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cluster/health"
)

func versions(stats []health.Stat) error {
	var err error
	v := stats[0].Kernel.Version
	b := stats[0].Info.BIOS.Version
	for _, s := range stats[1:] {
		if s.Kernel.Version != v {
			err = errors.Join(err, fmt.Errorf("%s:Version %q differs from %q", s.Hostname, s.Kernel.Version, v))
		}
		if s.Info.BIOS.Version != b {
			err = errors.Join(err, fmt.Errorf("%s:BIOS Version %q differs from %q", s.Hostname, s.Info.BIOS.Version, b))
		}
	}
	return err
}

func hardware(stats []health.Stat) error {
	var err error
	c := stats[0].Info.CPU
	cores := c.TotalCores
	threads := c.TotalThreads
	for _, s := range stats[1:] {
		c := s.Info.CPU
		if c.TotalCores != cores {
			err = errors.Join(err, fmt.Errorf("%s:TotalCores %d differs from %d", s.Hostname, c.TotalCores, cores))
		}
		if c.TotalThreads != threads {
			err = errors.Join(err, fmt.Errorf("%s:TotalThreads %d differs from %d", s.Hostname, c.TotalThreads, threads))
		}
	}
	return err
}

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var s []health.Stat
	if err := json.Unmarshal(b, &s); err != nil {
		log.Fatal(err)
	}

	good := make([]health.Stat, 0, len(s))
	for i, node := range s {
		if len(node.Err) > 0 {
			log.Printf("dropping %d:%s as it has error:%s", i, node.Hostname, node.Err)
			continue
		}
		good = append(good, node)
	}

	if len(good) == 0 {
		log.Fatalf("Removed all nodes with errors: none left.")
	}

	if err := versions(good); err != nil {
		log.Printf("versions check:%v", err)
	}

}
