// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func ipVM(t *testing.T, name, script string, net *qnetwork.InterVM) *qemu.VM {
	return scriptvm.Start(t, name, script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/grep",
				"github.com/u-root/u-root/cmds/core/sleep",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(4*time.Minute),
			net.NewVM(),
		),
	)
}

func TestIP(t *testing.T) {
	net := qnetwork.NewInterVM()

	script := `
		# Function to convert dotted decimal IP to byte-swapped hex
		ip_to_route_hex() {
  			local ip=$1
  			# Split the IP into octets
  			IFS='.' read -r o1 o2 o3 o4 <<< "$ip"
  			# Convert to hex, pad with zeros, and swap the byte order
  			printf "%02x%02x%02x%02x" $o4 $o3 $o2 $o1
		}

		# Verify that eth0 and lo exist
		test -d /sys/class/net/eth0 || exit 1
		test -d /sys/class/net/lo || exit 1

		# Bring the eth0 interface up
		ip link set eth0 up || exit 1
		sleep 3
		state=$(cat /sys/class/net/eth0/operstate)
		test "$state" = "up" || exit 1

		# Assign an IP address to eth0
		ip addr add 192.168.1.1/24 dev eth0 || exit 1
		grep -q "192.168.1.1" /proc/net/fib_trie || exit 1

		# Add a route via eth0
		cat /proc/net/route
		ip route add 192.168.2.0/24 via 192.168.1.1 dev eth0 || exit 1
		cat /proc/net/route
		hex_destination=$(ip_to_route_hex "192.168.2.0")
		hex_gateway=$(ip_to_route_hex "192.168.1.1")
		grep -iq "$hex_destination" /proc/net/route || exit 1
		grep -iq "$hex_gateway" /proc/net/route || exit 1

		# Delete the route
		cat /proc/net/route
		ip route del 192.168.2.0/24 || exit 1
		cat /proc/net/route
		#! grep -iq "$hex_destination" /proc/net/route || exit 1

		# Bring the eth0 interface down
		ip link set eth0 down || exit 1
		sleep 2
		state=$(cat /sys/class/net/eth0/operstate)
		test "$state" = "down" || exit 1

		# Delete the IP address from eth0
		ip addr del 192.168.1.1/24 dev eth0 || exit 1
		! grep -q "192.168.1.1" /proc/net/fib_trie || exit 1


		echo "TESTS PASSED MARKER"
	`

	vm := ipVM(t, "ip_test", script, net)

	if _, err := vm.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("ip_test: %v", err)
	}

	vm.Wait()
}
