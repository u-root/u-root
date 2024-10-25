// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bufio"
	"fmt"
	"os"
)

// Holds allowed and denied hosts for the connection.
// The map is a string representation of the host is allowed or denied.
// The key is the host and the value is true for allowed and false for denied.
// If the host is not in the map, it is allowed by default.
// The map is populated by the access cli flags, with access flags taking precedence over deny flags.
type AccessControlOptions struct {
	SetAllowed     bool
	ConnectionList map[string]bool
}

// Check if the hostName is allowed to connect.
func (ac *AccessControlOptions) IsAllowed(hostNames []string) bool {
	// If atleast one item is in the allowed list, the hostName is only allowed if also mentioned.
	if ac.SetAllowed {
		for _, hostName := range hostNames {
			allowed, ok := ac.ConnectionList[hostName]
			if !ok {
				// If the hostname is not in the list, check the next hostname.
				continue
			}

			// If the host name is allowed, return true.
			return allowed
		}

		// If the host name is not in the list, it is denied by default.
		return false
	}

	// If the host is part of the list and not one entry is set to allowed, it is denied by default.
	for _, hostName := range hostNames {
		// Otherwise check if the host name is denied.
		if _, ok := ac.ConnectionList[hostName]; ok {
			return false
		}
	}

	// If the host is not in the list, it is allowed by default.
	return true
}

func ParseAccessControl(connectionAllowFile string, connectionAllowList []string, connectionDenyFile string, connectionDenyList []string) (AccessControlOptions, error) {
	accessControl := AccessControlOptions{}

	accessControl.ConnectionList = make(map[string]bool)

	if connectionDenyFile != "" {
		denyFile, err := os.Open(connectionDenyFile)
		if err != nil {
			return accessControl, fmt.Errorf("failed to open deny file: %w", err)
		}
		defer denyFile.Close()

		scanner := bufio.NewScanner(denyFile)
		for scanner.Scan() {
			line := scanner.Text()
			accessControl.ConnectionList[line] = false
		}

		if err := scanner.Err(); err != nil {
			return accessControl, fmt.Errorf("failed to read deny file: %w", err)
		}
	}

	for _, deniedLine := range connectionDenyList {
		accessControl.ConnectionList[deniedLine] = false
	}

	// allowed hosts are parsed secondly as they take precedence
	if len(connectionAllowList) > 0 || connectionAllowFile != "" {
		accessControl.SetAllowed = true
	}

	if connectionAllowFile != "" {
		allowFile, err := os.Open(connectionAllowFile)
		if err != nil {
			return accessControl, fmt.Errorf("failed to open allow file: %w", err)
		}
		defer allowFile.Close()

		scanner := bufio.NewScanner(allowFile)
		for scanner.Scan() {
			line := scanner.Text()
			accessControl.ConnectionList[line] = true
		}

		if err := scanner.Err(); err != nil {
			return accessControl, fmt.Errorf("failed to read allow file: %w", err)
		}
	}

	for _, allowedLine := range connectionAllowList {
		accessControl.ConnectionList[allowedLine] = true
	}

	return accessControl, nil
}
