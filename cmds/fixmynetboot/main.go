package main

import (
	"github.com/systemboot/systemboot/pkg/checker"
)

// fixmynetboot is a troubleshooting tool that can help you identify issues that
// won't let your system boot over the network.

func main() {
	checklist := []checker.Check{
		checker.Check{"wlp2s0 exists", checker.InterfaceExists("wlp2s0"), nil, false},
		checker.Check{"eth0 exists", checker.InterfaceExists("eth0"), checker.InterfaceRemediate("eth0"), false},
		checker.Check{"eth0 link speed", checker.LinkSpeed("eth0", 100), nil, false},
		checker.Check{"eth0 link autoneg", checker.LinkAutoneg("eth0", true), nil, false},
		checker.Check{"eth0 has link-local", checker.InterfaceHasLinkLocalAddress("wlp2s0"), nil, false},
		checker.Check{"eth0 has global addresses", checker.InterfaceHasGlobalAddresses("wlp2s0"), nil, false},
	}

	checker.Run(checklist)

}
