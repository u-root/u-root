// Copyright 2021-2025 the u-root Authors. All rights reserved
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

func brctlVM(t *testing.T, name, script string, net *qnetwork.InterVM) *qemu.VM {
	return scriptvm.Start(t, name, script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/brctl",
				"github.com/u-root/u-root/cmds/core/grep",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/cat",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(4*time.Minute),
			net.NewVM(),
		),
	)
}

func TestBrctl(t *testing.T) {
	net := qnetwork.NewInterVM()

	script := `
		# Add a bridge
		brctl addbr br0 || exit 1
		test -d /sys/class/net/br0/bridge || exit 1

		# Show bridges
		brctl show > show.out || exit 1
		grep -q br0 show.out || exit 1

		# Show a non-existing bridge
		brctl show non_existing > non_existing.out 2>&1 && exit 1
		grep -q "does not exist" non_existing.out || exit 1

		# Add an interface to the bridge
		ip link set eth0 up || exit 1
		brctl addif br0 eth0 || exit 1
		test -d /sys/class/net/br0/brif/eth0 || exit 1

		# Show the bridge with the interface
		brctl show br0 > show_br0.out || exit 1
		grep -q eth0 show_br0.out || exit 1

		 # Set hairpin mode on
		brctl hairpin br0 eth0 on || exit 1
		hairpin_mode=$(cat /sys/class/net/br0/brif/eth0/hairpin_mode)
		test "$hairpin_mode" -eq 1 || exit 1

		# Set hairpin mode off
		brctl hairpin br0 eth0 off || exit 1
		hairpin_mode=$(cat /sys/class/net/br0/brif/eth0/hairpin_mode)
		test "$hairpin_mode" -eq 0 || exit 1

		 # Set ageing time
		brctl setageing br0 10 || exit 1
		ageing_time=$(cat /sys/class/net/br0/bridge/ageing_time)
		test "$ageing_time" -eq 1000 || exit 1

		# Show MAC addresses
		brctl showmacs br0 > showmacs.out || exit 1
				expected_macs=""
		for iface in $(ls /sys/class/net/br0/brif/); do
			mac=$(cat /sys/class/net/$iface/address)
			expected_macs="$expected_macs $mac"
        done
				for mac in $expected_macs; do
		grep -q "$mac" showmacs.out || (echo "Missing MAC address: $mac" && exit 1)
		done

		# STP-related cases
		# Enable STP
		brctl stp br0 on || exit 1
		stp_state=$(cat /sys/class/net/br0/bridge/stp_state)
		test "$stp_state" -eq 1 || exit 1

		 # Define input values
		bridge_prio=0
		forward_delay=2
		hello_time=1
		max_age=20
		path_cost=10
		port_prio=1

		# Set bridge priority
		brctl setbridgeprio br0 $bridge_prio || exit 1
		sysfs_bridge_prio=$(cat /sys/class/net/br0/bridge/priority)
		test "$sysfs_bridge_prio" -eq $bridge_prio || exit 1

		# Set forward delay
		brctl setfd br0 $forward_delay || exit 1
		sysfs_forward_delay=$(cat /sys/class/net/br0/bridge/forward_delay)
		test "$sysfs_forward_delay" -eq 200 || exit 1

		# Set hello time
		brctl sethello br0 $hello_time || exit 1
		sysfs_hello_time=$(cat /sys/class/net/br0/bridge/hello_time)
		test "$sysfs_hello_time" -eq 100 || exit 1

		# Set max age
		brctl setmaxage br0 $max_age || exit 1
		sysfs_max_age=$(cat /sys/class/net/br0/bridge/max_age)
		test "$sysfs_max_age" -eq 2000 || exit 1

		# Set path cost
		brctl setpathcost br0 eth0 $path_cost || exit 1
		sysfs_path_cost=$(cat /sys/class/net/eth0/brport/path_cost)
		test "$sysfs_path_cost" -eq $path_cost || exit 1

		# Set port priority
		brctl setportprio br0 eth0 $port_prio || exit 1
		sysfs_port_prio=$(cat /sys/class/net/eth0/brport/priority)
		test "$sysfs_port_prio" -eq $port_prio || exit 1

		# Show STP information
				brctl showstp br0 > showstp.out || exit 1
		cat showstp.out
		grep -q "bridge id.*$(cat /sys/class/net/br0/bridge/bridge_id)" showstp.out || exit 1
		grep -q "designated root.*$(cat /sys/class/net/br0/bridge/root_id)" showstp.out || exit 1
		grep -q "path cost.*$(cat /sys/class/net/br0/bridge/root_path_cost)" showstp.out || exit 1
		grep -q "root port.*$(cat /sys/class/net/br0/bridge/root_port)" showstp.out || exit 1
		grep -q "max age.*$max_age" showstp.out || exit 1
		grep -q "hello time.*$hello_time" showstp.out || exit 1
		grep -q "forward delay.*$forward_delay" showstp.out || exit 1

		 # Validate interface information in STP output
		for iface in $(ls /sys/class/net/br0/brif/); do
			expected_port_id=$(cat /sys/class/net/$iface/brport/port_id)
			expected_path_cost=$(cat /sys/class/net/$iface/brport/path_cost)
			expected_priority=$(cat /sys/class/net/$iface/brport/priority)
			grep -q "$iface" showstp.out || exit 1
			grep -q "port id.*$expected_port_id" showstp.out || exit 1
			grep -q "path cost.*$expected_path_cost" showstp.out || exit 1
			grep -q "priority.*$expected_priority" showstp.out || exit 1
		done

		# Remove the interface from the bridge
		brctl delif br0 eth0 || exit 1
		test ! -d /sys/class/net/br0/brif/eth0 || exit 1

		# Delete the bridge
		brctl delbr br0 || exit 1
		test ! -d /sys/class/net/br0 || exit 1

		# Delete a non-existing bridge
		brctl delbr non_existing > del_non_existing.out 2>&1 && exit 1
		grep -q "no such device" del_non_existing.out || exit 1

		# Add an interface to a non-existing bridge
		brctl addif non_existing eth0 > addif_non_existing.out 2>&1 && exit 1
		grep -q "no such device" addif_non_existing.out || exit 1

		# Delete an interface from a non-existing bridge
		brctl delif non_existing eth0 > delif_non_existing.out 2>&1 && exit 1
		grep -q "no such device" delif_non_existing.out || exit 1

		# Delete an interface not part of the bridge
		brctl addbr br1 || exit 1
		brctl delif br1 eth0 > delif_not_part.out 2>&1 && exit 1
		grep -q "not a member" delif_not_part.out || exit 1
		brctl delbr br1 || exit 1

		echo "TESTS PASSED MARKER"
	`

	vm := brctlVM(t, "brctl_test", script, net)

	if _, err := vm.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("brctl_test: %v", err)
	}

	vm.Wait()
}
