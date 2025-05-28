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
		ip route add 192.168.2.0/24 via 192.168.1.1 dev eth0 || exit 1
		hex_destination=$(ip_to_route_hex "192.168.2.0")
		hex_gateway=$(ip_to_route_hex "192.168.1.1")
		grep -iq "$hex_destination" /proc/net/route || exit 1
		grep -iq "$hex_gateway" /proc/net/route || exit 1

		# Delete the route
		ip route del 192.168.2.0/24 || exit 1
		#! grep -iq "$hex_destination" /proc/net/route || exit 1

		# Add a tunnel
		ip tunnel add my_test_tunnel mode sit remote 192.168.2.1 local 192.168.1.1 ttl 64

		# Verify tunnel exists in /proc/net/dev
		grep -q "my_test_tunnel" /proc/net/dev || exit 1

		# Verify the tunnel is created and has the right parameters
		ip tunnel show my_test_tunnel || exit 1
		tunnel_info=$(ip tunnel show my_test_tunnel)
		echo "$tunnel_info" | grep -q "my_test_tunnel:" || exit 1
		echo "$tunnel_info" | grep -q "remote 192.168.2.1" || exit 1
		echo "$tunnel_info" | grep -q "local 192.168.1.1" || exit 1
		echo "$tunnel_info" | grep -q "ttl 64" || exit 1

		# Delete the tunnel
		ip tunnel del my_test_tunnel

		# Verify tunnel no longer exists in /proc/net/dev
		! grep -q "my_test_tunnel" /proc/net/dev || exit 1

		# Add a GRE tunnel with key and tos options
		ip tunnel add gre_tunnel mode gre remote 192.168.2.2 local 192.168.1.1 ttl 128 key 1234 tos 10

		# Verify GRE tunnel exists in /proc/net/dev
		grep -q "gre_tunnel" /proc/net/dev || exit 1

		# Verify GRE tunnel parameters
		gre_info=$(ip tunnel show gre_tunnel)
		echo "$gre_info" | grep -q "gre_tunnel:" || exit 1
		echo "$gre_info" | grep -q "remote 192.168.2.2" || exit 1
		echo "$gre_info" | grep -q "local 192.168.1.1" || exit 1
		echo "$gre_info" | grep -q "ttl 128" || exit 1
		echo "$gre_info" | grep -q "key 1234" || exit 1
		echo "$gre_info" | grep -q "tos 0xa" || exit 1

		# Configure GRE tunnel
		ip link set gre_tunnel up || exit 1
		ip addr add 10.0.0.1/24 dev gre_tunnel || exit 1
		grep -q "10.0.0.1" /proc/net/fib_trie || exit 1

		# Delete GRE tunnel
		ip link set gre_tunnel down || exit 1
		ip tunnel del gre_tunnel || exit 1
		! grep -q "gre_tunnel" /proc/net/dev || exit 1

		# Add a VTI tunnel
		ip tunnel add vti_tunnel mode vti remote 192.168.2.3 local 192.168.1.1 key 5678 

		# Verify VTI tunnel exists in /proc/net/dev
		grep -q "vti_tunnel" /proc/net/dev || exit 1

		# Verify VTI tunnel parameters
		vti_info=$(ip tunnel show vti_tunnel)
		echo "$vti_info" | grep -q "vti_tunnel:" || exit 1
		echo "$vti_info" | grep -q "remote 192.168.2.3" || exit 1
		echo "$vti_info" | grep -q "local 192.168.1.1" || exit 1
		echo "$vti_info" | grep -q "key 5678" || exit 1

		# Configure VTI tunnel
		ip link set vti_tunnel up || exit 1
		ip addr add 172.16.0.1/30 dev vti_tunnel || exit 1
		grep -q "172.16.0.1" /proc/net/fib_trie || exit 1

		# Delete VTI tunnel
		ip link set vti_tunnel down || exit 1
		ip tunnel del vti_tunnel || exit 1
		! grep -q "vti_tunnel" /proc/net/dev || exit 1

        # Add an IPIP tunnel 
        ip tunnel add ipip_tunnel mode ipip remote 192.168.3.1 local 192.168.1.1 ttl 64

        # Verify IPIP tunnel exists in /proc/net/dev
        grep -q "ipip_tunnel" /proc/net/dev || exit 1

        # Verify IPIP tunnel parameters
        ipip_info=$(ip tunnel show ipip_tunnel)
        echo "$ipip_info" | grep -q "ipip_tunnel:" || exit 1
        echo "$ipip_info" | grep -q "remote 192.168.3.1" || exit 1
        echo "$ipip_info" | grep -q "local 192.168.1.1" || exit 1
        echo "$ipip_info" | grep -q "ttl 64" || exit 1

        # Configure IPIP tunnel
        ip link set ipip_tunnel up || exit 1
        ip addr add 172.17.0.1/30 dev ipip_tunnel || exit 1
        grep -q "172.17.0.1" /proc/net/fib_trie || exit 1

        # Delete IPIP tunnel
        ip link set ipip_tunnel down || exit 1
        ip tunnel del ipip_tunnel || exit 1
        ! grep -q "ipip_tunnel" /proc/net/dev || exit 1

        # Add a neighbor (ARP entry) on eth0
		ip neigh add 192.168.1.2 lladdr 00:11:22:33:44:55 dev eth0 || exit 1
		grep -q "192.168.1.2" /proc/net/arp || exit 1

		# Verify the neighbor entry
		ip neigh show dev eth0 || exit 1
		neigh_entry=$(ip neigh show dev eth0 | grep "192.168.1.2")
		test "$neigh_entry" = "192.168.1.2 dev eth0 lladdr 00:11:22:33:44:55 PERMANENT" || exit 1

		# Replace the entry with another hwaddress, nud state and router flag
		ip neigh replace 192.168.1.2 lladdr 11:22:33:44:55:66 dev eth0 nud stale router || exit 1

		# Verify the modified flags
		ip neigh show dev eth0 || exit 1
		modified_entry=$(ip neigh show dev eth0 | grep "192.168.1.2")
		test "$modified_entry" = "192.168.1.2 dev eth0 lladdr 11:22:33:44:55:66 router STALE" || exit 1
		echo "Modified neighbor entry verified with router flag and STALE state"

		# Delete the neighbor
		ip neigh del 192.168.1.2 dev eth0 || exit 1
		! grep -q "192.168.1.2" /proc/net/arp || exit 1


		# Test IP Neighbor flush capability
		# Add 3 neighbors
		ip neigh add 192.168.1.5 lladdr aa:bb:cc:dd:ee:ff nud stale dev eth0 || exit 1
		ip neigh add 192.168.1.6 lladdr aa:bb:cc:11:22:33 nud stale dev eth0 || exit 1
		ip neigh add 192.168.1.7 lladdr aa:bb:cc:44:55:66 dev eth0 || exit 1

		# Verify all entries exist
		grep -q "192.168.1.5" /proc/net/arp || exit 1
		grep -q "192.168.1.6" /proc/net/arp || exit 1
		grep -q "192.168.1.7" /proc/net/arp || exit 1

		# Flush the 2 stale neighbors from the table for eth0
		ip neigh flush dev eth0 || exit 1

		# Verify the 2 stale entries are gone, the permanent one remains
		! grep -q "192.168.1.5" /proc/net/arp || exit 1
		! grep -q "192.168.1.6" /proc/net/arp || exit 1
		grep -q "192.168.1.7" /proc/net/arp || exit 1

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
