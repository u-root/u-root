// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race && amd64

package integration

import (
	"strings"
	"testing"
	"time"
)

func count(s string, t []string) int {
	var i int
	for _, n := range t {
		if strings.Contains(s, n) {
			i++
		}
	}
	return i
}

func all(s string, t []string) bool {
	return count(s, t) == len(t)
}

func some(s string, t []string) bool {
	return count(s, t) > 0
}

func none(s string, t []string) bool {
	return count(s, t) == 0
}

// tunnelModes returns the tunnel modes a group of tests needs.
//
// A group sets up and tears down together, so it is skipped whole rather than
// command by command, and the modes it needs are whatever its "ip tunnel add"
// commands name.
func tunnelModes(group []iptest) []string {
	var modes []string
	for _, tt := range group {
		var cmd string
		switch c := tt.cmd.(type) {
		case string:
			cmd = c
		case []string:
			cmd = strings.Join(c, " ")
		}
		f := strings.Fields(cmd)
		for n := range f {
			if f[n] == "mode" && n+1 < len(f) && strings.Contains(cmd, "tunnel add") {
				modes = append(modes, f[n+1])
			}
		}
	}
	return modes
}

// tunnelSupport answers, once per mode, whether this guest's kernel can
// create a tunnel of that mode.
//
// The embedded kernels differ by architecture: amd64 carries gre and vti, and
// the arm and arm64 ones carry neither, so those groups cannot pass there.
// Asking the guest is better than listing architectures, because it is the
// kernel that decides and these kernels are not ours.
type tunnelSupport struct {
	g    *guest
	seen map[string]bool
}

func (s *tunnelSupport) has(mode string) bool {
	if ok, done := s.seen[mode]; done {
		return ok
	}
	_, err := s.g.Run("/bbin/ip", "tunnel", "add", "probe0", "mode", mode,
		"remote", "192.168.242.9", "local", "192.168.241.9")
	ok := err == nil
	if ok {
		_, _ = s.g.Run("/bbin/ip", "tunnel", "del", "probe0")
	}
	s.seen[mode] = ok
	return ok
}

type iptest struct {
	cmd      any
	delay    time.Duration
	includes []string
	excludes []string
}

// TestIP tests creation and removal of addresses, tunnels, and
// ARP entries with the u-root ip command.
//
// Every architecture here now gets a guest of its own, which is new: until
// each VM had its own socket they all forwarded the same host port, the
// second qemu died with "Could not set up host forwarding rule", and every
// architecture's commands went to the first guest. So these tables have only
// ever run against amd64, and what the other architectures cannot do is only
// now visible.
func TestIP(t *testing.T) {
	t.Parallel()

	for _, arch := range []string{"amd64", "arm", "arm64", "riscv64"} {
		t.Run(arch, func(t *testing.T) {
			t.Parallel()

			if arch == "riscv64" {
				// The guest boots and takes a DHCPv4 lease, then
				// sits in DHCPv6 solicit and never starts cpud.
				// Measured at over ten minutes, alone and
				// unloaded. That is the embedded kernel and
				// init, not this test.
				t.Skip("the riscv64 guest does not reach cpud; see u-root/cpu")
			}

			g := bootGuest(t, arch)
			tun := &tunnelSupport{g: g, seen: map[string]bool{}}

			// This is a slice of iptest slices.
			// The intent is to create a grouping of
			// different tests. The groups can not run in parallel
			// as they would interfere with other.
			for _, group := range [][]iptest{
				{
					{cmd: "ip link set eth1 down", delay: 3 * time.Second},
					{cmd: "cat /sys/class/net/eth1/operstate", includes: []string{"down"}},
					{cmd: "ip link set eth1 up", delay: 3 * time.Second},
					{cmd: "cat /sys/class/net/eth1/operstate", includes: []string{"up"}},
					{cmd: []string{"ip", "addr", "add", "192.168.241.1/24", "dev", "eth1"}},
					{cmd: []string{"cat", "/proc/net/fib_trie"}, includes: []string{"192.168.241.1"}},
					{cmd: "ip route add 192.168.242.0/24 via 192.168.241.1 dev eth1"},
					{cmd: "cat /proc/net/route", includes: []string{"00F2A8C0", "01F1A8C0"}},
					{cmd: "ip route del 192.168.242.0/24"},
					{cmd: "cat /proc/net/route", excludes: []string{"00F2A8C0", "01F1A8C0"}},
					{cmd: "ip tunnel add my_test_tunnel mode sit remote 192.168.242.1 local 192.168.241.1 ttl 64"},
					{cmd: "cat /proc/net/dev", includes: []string{"my_test_tunnel"}},
					{cmd: "ip tunnel del my_test_tunnel"},
					{cmd: "cat /proc/net/dev", excludes: []string{"my_test_tunnel"}},
					{cmd: "ip tunnel add my_test_tunnel mode sit remote 192.168.242.1 local 192.168.241.1 ttl 64"},
					{cmd: "ip tunnel show my_test_tunnel", includes: []string{"my_test_tunnel", "remote 192.168.242.1", "local 192.168.241.1", "ttl 64"}},
					{cmd: "ip tunnel del my_test_tunnel"},
					{cmd: "cat /proc/net/dev", excludes: []string{"my_test_tunnel"}},
				},
				{
					// Various tunnel tests.
					// Add a GRE tunnel with key and tos options
					{cmd: "ip tunnel add gre_tunnel mode gre remote 192.168.242.2 local 192.168.241.1 ttl 128 key 1234 tos 10", delay: time.Second},
					{cmd: "cat /proc/net/dev", includes: []string{"gre_tunnel"}},
					// Verify GRE tunnel parameters
					{cmd: "ip tunnel show gre_tunnel", includes: []string{"gre_tunnel:", "remote 192.168.242.2", "local 192.168.241.1", "ttl 128", "key 1234", "tos 0xa"}},
					{cmd: "ip link set gre_tunnel up"},
					{cmd: "ip addr add 10.0.0.1/24 dev gre_tunnel"},
					{cmd: []string{"cat", "/proc/net/fib_trie"}, includes: []string{"10.0.0.1"}},
					{cmd: "ip link set gre_tunnel down"},
					{cmd: "ip tunnel del gre_tunnel"},
					{cmd: "cat /proc/net/dev", excludes: []string{"gre_tunnel"}},
				},
				{
					{cmd: "ip tunnel add vti_tunnel mode vti remote 192.168.242.3 local 192.168.241.1 key 5678"},

					// Verify VTI tunnel exists in /proc/net/dev
					{cmd: "cat /proc/net/dev", includes: []string{"vti_tunnel"}},

					// Verify VTI tunnel parameters
					{cmd: "ip tunnel show vti_tunnel", includes: []string{"vti_tunnel:", "remote 192.168.242.3", "local 192.168.241.1", "key 5678"}},
					{cmd: "ip link set vti_tunnel up"},
					{cmd: "ip addr add 172.16.0.1/30 dev vti_tunnel"},
					{cmd: []string{"cat", "/proc/net/fib_trie"}, includes: []string{"172.16.0.1"}},
					{cmd: "ip link set vti_tunnel down"},
					{cmd: "ip tunnel del vti_tunnel"},
					{cmd: "cat /proc/net/dev", excludes: []string{"vti_tunnel"}},
				},
				{
					{cmd: "ip tunnel add ipip_tunnel mode ipip remote 192.168.243.1 local 192.168.241.1 ttl 64"},

					// Verify IPIP tunnel exists in /proc/net/dev
					{cmd: "cat /proc/net/dev", includes: []string{"ipip_tunnel"}},

					// Verify IPIP tunnel parameters
					{cmd: "ip tunnel show ipip_tunnel", includes: []string{"ipip_tunnel:", "remote 192.168.243.1", "local 192.168.241.1", "ttl 64"}},
					{cmd: "ip link set ipip_tunnel up"},
					{cmd: "ip addr add 172.17.0.1/30 dev ipip_tunnel"},
					{cmd: []string{"cat", "/proc/net/fib_trie"}, includes: []string{"172.17.0.1"}},
					{cmd: "ip link set ipip_tunnel down"},
					{cmd: "ip tunnel del ipip_tunnel"},
					{cmd: "cat /proc/net/dev", excludes: []string{"ipip_tunnel"}},
				},
				{
					// ARP tests
					{cmd: "ip neigh add 192.168.241.2 lladdr 00:11:22:33:44:55 dev eth1"},
					{cmd: "cat /proc/net/arp", includes: []string{"192.168.241.2"}},

					// Verify the neighbor entry
					{cmd: "ip neigh show dev eth1", includes: []string{"192.168.241.2", "192.168.241.2 dev eth1 lladdr 00:11:22:33:44:55 PERMANENT"}},

					// Replace the entry with another hwaddress, nud state and router flag
					{cmd: "ip neigh replace 192.168.241.2 lladdr 11:22:33:44:55:66 dev eth1 nud stale router", includes: []string{}},

					// Verify the modified flags
					{cmd: "ip neigh show dev eth1", includes: []string{"192.168.241.2", "192.168.241.2 dev eth1 lladdr 11:22:33:44:55:66 router STALE"}},

					// Delete the neighbor
					{cmd: "ip neigh del 192.168.241.2 dev eth1", includes: []string{}},
					{cmd: "cat /proc/net/arp", excludes: []string{"192.168.241.2"}},

					// Test IP Neighbor flush capability
					// Add 3 neighbors
					{cmd: "ip neigh add 192.168.241.5 lladdr aa:bb:cc:dd:ee:ff nud stale dev eth1"},
					{cmd: "ip neigh add 192.168.241.6 lladdr aa:bb:cc:11:22:33 nud stale dev eth1"},
					{cmd: "ip neigh add 192.168.241.7 lladdr aa:bb:cc:44:55:66 dev eth1"},

					// Verify all entries exist
					{cmd: "cat /proc/net/arp", includes: []string{"192.168.241.5", "192.168.241.6", "192.168.241.7"}},

					// Flush the 2 stale neighbors from the table for eth1
					{cmd: "ip neigh flush dev eth1", includes: []string{}},

					// Verify the 2 stale entries are gone, the permanent one remains
					{cmd: "cat /proc/net/arp", includes: []string{"192.168.241.7"}, excludes: []string{"192.168.241.5", "192.168.241.6"}},

					// Delete the IP address from eth1
					{cmd: "ip addr del 192.168.241.1/24 dev eth1"},
					{cmd: []string{"cat", "/proc/net/fib_trie"}, excludes: []string{"192.168.241.1"}},
					// Bring the eth1 interface down
					{cmd: "ip link set eth1 down", delay: 2 * time.Second},
					{cmd: "cat /sys/class/net/eth1/operstate", includes: []string{"down"}},
				},
			} {

				var missing string
				for _, m := range tunnelModes(group) {
					if !tun.has(m) {
						missing = m
						break
					}
				}
				if missing != "" {
					t.Logf("skipping a group: this guest kernel cannot make a %q tunnel", missing)
					continue
				}
				for _, tt := range group {
					var cmd []string
					switch c := tt.cmd.(type) {
					case []string:
						cmd = c
					case string:
						cmd = strings.Fields(c)
					default:
						t.Fatalf("tt.cmd.Type(): type %T, want string or []string", tt.cmd)
					}
					t.Logf("Run %s", cmd)
					b, err := g.Run(cmd[0], cmd[1:]...)
					if err != nil {
						// Without b this reads only "Process
						// exited with status 1", which says
						// nothing about which operation the
						// kernel refused.
						t.Fatalf("%s: got %v (%s), want nil", cmd, err, b)
					}
					//t.Logf("%q, includes %s?", string(b), tt.includes)
					if !all(string(b), tt.includes) {
						t.Fatalf("%s: got %s, does not contain all of %s", cmd, string(b), tt.includes)
					}
					//t.Logf("%q, excludes %s?", string(b), tt.includes)
					if some(string(b), tt.excludes) {
						t.Fatalf("%s:got %s, contains some of %s and should not", cmd, string(b), tt.excludes)
					}
					if tt.delay > 0 {
						time.Sleep(tt.delay)
					}
				}
			}
		})
	}
}
