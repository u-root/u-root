// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	data, err := os.ReadFile("/etc/services")
	if err != nil {
		log.Fatalf("Error opening /etc/services: %v", err)
	}

	re := regexp.MustCompile(`^(\S+)\s+(\d+)/\S+`)

	ports := make(map[string]string)

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			service := matches[1]
			port := matches[2]
			ports[port] = service
		}
	}

	// Generate the content for the Go file
	content := `// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var wellKnownPortsMap = map[string]string{
`

	// Append the ports to the content
	for port, service := range ports {
		content += `"` + port + `":"` + service + `",` + "\n"
	}
	content += "}"

	// Write the content to the Go file
	if err := os.WriteFile("well_known_ports.go", []byte(content), 0o644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	if err := exec.Command("gofmt", "-w", "well_known_ports.go").Run(); err != nil {
		log.Fatalf("Error running gofmt: %v", err)
	}
}
