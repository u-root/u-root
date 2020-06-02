// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pxe implements the PXE config file parsing.
//
// See http://www.pix.net/software/pxeboot/archive/pxespec.pdf
package pxe

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/curl"
)

// ParseConfig probes for config files based on the Mac and IP given
// and uses s to fetch files.
func ParseConfig(ctx context.Context, workingDir *url.URL, mac net.HardwareAddr, ip net.IP, s curl.Schemes) ([]boot.OSImage, error) {
	rootDir := *workingDir
	rootDir.Path = ""

	for _, relname := range probeFiles(mac, ip) {
		// "When booting, the initial working directory for PXELINUX
		// will be the parent directory of pxelinux.0 unless overridden
		// with DHCP option 210."
		//
		// https://wiki.syslinux.org/wiki/index.php?title=Config#Working_directory
		imgs, err := syslinux.ParseConfigFile(ctx, s, path.Join("pxelinux.cfg", relname), &rootDir, workingDir.Path)
		if curl.IsURLError(err) {
			// We didn't find the file.
			// TODO(hugelgupf): log this.
			continue
		}
		return imgs, err
	}
	return nil, fmt.Errorf("no valid pxelinux config found")
}

func probeFiles(ethernetMac net.HardwareAddr, ip net.IP) []string {
	files := make([]string, 0, 10)
	// Skipping client UUID. Figure that out later.

	// MAC address.
	files = append(files, fmt.Sprintf("01-%s", strings.ToLower(strings.Replace(ethernetMac.String(), ":", "-", -1))))

	// IP address in upper case hex, chopping one letter off at a time.
	if ip != nil {
		ipf := strings.ToUpper(hex.EncodeToString(ip))
		for n := len(ipf); n >= 1; n-- {
			files = append(files, ipf[:n])
		}
	}
	files = append(files, "default")
	return files
}
