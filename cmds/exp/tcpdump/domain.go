// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"strings"

	"github.com/gopacket/gopacket/layers"
)

// dnsData returns a string representation of the DNS layer.
func dnsData(layer *layers.DNS) string {
	applicationData := fmt.Sprintf("%d", layer.ID)

	if layer.ResponseCode != layers.DNSResponseCodeNoErr {
		applicationData += fmt.Sprintf(" %s", layer.ResponseCode.String())
	}

	if layer.AA {
		applicationData += "*"
	} else if layer.RD {
		applicationData += "+"
	}

	if len(layer.Answers)+len(layer.Authorities)+len(layer.Additionals) > 1 {
		applicationData += fmt.Sprintf(" %d/%d/%d ", len(layer.Answers), len(layer.Authorities), len(layer.Additionals))
	}

	if layer.ResponseCode == layers.DNSResponseCodeNoErr {

		if len(layer.Answers) == 0 {
			for _, question := range layer.Questions {
				applicationData += fmt.Sprintf(" %s? %s, ", question.Type.String(), question.Name)
			}
		}

		for _, answer := range layer.Answers {
			applicationData += answer.String() + ", "
		}

		applicationData = strings.TrimRight(applicationData, ", ")
	}

	applicationData += fmt.Sprintf((" (%d)"), len(layer.Contents))

	return applicationData
}
