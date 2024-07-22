// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/vishvananda/netlink"
)

const linkHelp = `Usage: ip link add  [ name ] NAME
		    [ txqueuelen PACKETS ]
		    [ address LLADDR ]
		    [ broadcast LLADDR ]
		    [ mtu MTU ] [index IDX ]
		    [ numtxqueues QUEUE_COUNT ]
		    [ numrxqueues QUEUE_COUNT ]
		    type TYPE 

	ip link delete { DEVICE | dev DEVICE  } 

	ip link set { DEVICE | dev DEVICE | group DEVGROUP }
			[ { up | down } ]
			[ type TYPE ARGS ]
		[ arp { on | off } ]
		[ multicast { on | off } ]
		[ allmulticast { on | off } ]
		[ promisc { on | off } ]
		[ txqueuelen PACKETS ]
		[ name NEWNAME ]
		[ address LLADDR ]
		[ mtu MTU ]
		[ netns { PID | NAME } ]
		[ alias NAME ]
		[ vf NUM [ mac LLADDR ]
			 [ vlan VLANID [ qos VLAN-QOS ] [ proto VLAN-PROTO ] ]
			 [ rate TXRATE ]
			 [ max_tx_rate TXRATE ]
			 [ min_tx_rate TXRATE ]
			 [ spoofchk { on | off} ]
			 [ state { auto | enable | disable} ]
			 [ trust { on | off} ]
			 [ node_guid EUI64 ]
			 [ port_guid EUI64 ] ]

	ip link show [ DEVICE | group GROUP ] [type TYPE]

	ip link help

TYPE := { bareudp | bond |bridge | dummy |
          geneve | gre | gretap | ifb |
          ip6gre | ip6gretap | ip6tnl | ipip |
          ipoib | ipvlan | ipvtap | macvlan |
          macvlan | sit | vlan | vrf |
          vti | vxlan | xfrm }

`

func (cmd cmd) linkSet() error {
	iface, err := parseDeviceName(true)
	if err != nil {
		return err
	}

	for cursor < len(arg)-1 {
		cursor++
		expectedValues = []string{"address", "up", "down", "arp", "promisc", "multicast", "allmulticast", "mtu", "name", "alias", "vf", "master", "nomaster", "netns", "txqueuelen", "txqlen", "group"}
		switch arg[cursor] {
		case "address":
			return cmd.setLinkHardwareAddress(iface)
		case "up":
			if err := cmd.handle.LinkSetUp(iface); err != nil {
				return fmt.Errorf("%v can't make it up: %v", iface.Attrs().Name, err)
			}
		case "down":
			if err := cmd.handle.LinkSetDown(iface); err != nil {
				return fmt.Errorf("%v can't make it down: %v", iface.Attrs().Name, err)
			}
		case "arp":
			cursor++
			expectedValues = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return cmd.handle.LinkSetARPOn(iface)
			case "off":
				return cmd.handle.LinkSetARPOff(iface)
			}
		case "promisc":
			cursor++
			expectedValues = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return cmd.handle.SetPromiscOn(iface)
			case "off":
				return cmd.handle.SetPromiscOff(iface)
			}
		case "multicast":
			cursor++
			expectedValues = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return cmd.handle.LinkSetMulticastOn(iface)
			case "off":
				return cmd.handle.LinkSetMulticastOff(iface)
			}
		case "allmulticast":
			cursor++
			expectedValues = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return cmd.handle.LinkSetAllmulticastOn(iface)
			case "off":
				return cmd.handle.LinkSetAllmulticastOff(iface)
			}
		case "mtu":
			return cmd.setLinkMTU(iface)
		case "name":
			return cmd.setLinkName(iface)
		case "alias":
			return cmd.setLinkAlias(iface)
		case "vf":
			return cmd.setLinkVf(iface)
		case "master":
			cursor++
			expectedValues = []string{"device name"}
			master, err := cmd.handle.LinkByName(arg[cursor])
			if err != nil {
				return err
			}
			return cmd.handle.LinkSetMaster(iface, master)
		case "nomaster":
			return cmd.handle.LinkSetNoMaster(iface)
		case "netns":
			return cmd.setLinkNetns(iface)
		case "txqueuelen", "txqlen":
			return cmd.setLinkTxQLen(iface)
		case "group":

		}
	}

	return nil
}

func (cmd cmd) setLinkHardwareAddress(iface netlink.Link) error {
	hwAddr, err := net.ParseMAC(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid mac address %v: %v", arg[cursor], err)
	}

	err = cmd.handle.LinkSetHardwareAddr(iface, hwAddr)
	if err != nil {
		return fmt.Errorf("%v cant set mac addr %v: %v", iface.Attrs().Name, hwAddr, err)
	}

	return nil
}

func (cmd cmd) setLinkMTU(iface netlink.Link) error {
	cursor++

	mtu, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid mtu %v: %v", arg[cursor], err)
	}

	return cmd.handle.LinkSetMTU(iface, mtu)
}

func (cmd cmd) setLinkGroup(iface netlink.Link) error {
	cursor++

	group, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid group %v: %v", arg[cursor], err)
	}

	return cmd.handle.LinkSetMTU(iface, group)
}

func (cmd cmd) setLinkName(iface netlink.Link) error {
	cursor++
	expectedValues = []string{"<name>"}
	name := arg[cursor]

	return cmd.handle.LinkSetName(iface, name)
}

func (cmd cmd) setLinkAlias(iface netlink.Link) error {
	cursor++
	expectedValues = []string{"<alias name>"}
	alias := arg[cursor]

	return cmd.handle.LinkSetAlias(iface, alias)
}

func (cmd cmd) setLinkTxQLen(iface netlink.Link) error {
	cursor++
	expectedValues = []string{"<qlen>"}
	qlen, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid queuelen %v: %v", arg[cursor], err)
	}

	return cmd.handle.LinkSetTxQLen(iface, qlen)
}

func (cmd cmd) setLinkNetns(iface netlink.Link) error {
	cursor++
	expectedValues = []string{"<netns pid>, <netns path>"}
	ns, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid int %v: %v", arg[cursor], err)
	}

	if err := cmd.handle.LinkSetNsPid(iface, ns); err != nil {
		if err := cmd.handle.LinkSetNsFd(iface, ns); err != nil {
			return fmt.Errorf("failed to set netns: %v", err)
		}
	}

	return nil
}

func (cmd cmd) setLinkVf(iface netlink.Link) error {
	vf, err := parseInt("VF")
	if err != nil {
		return err
	}

	cursor++

	expectedValues = []string{"vlan", "mac", "qos", "rate", "max_tx_rate", "min_tx_rate", "state", "spoofchk", "trust", "node_guid", "port_guid"}
	for cursor < len(arg)-1 {
		switch arg[cursor] {
		case "mac":
			addr, err := parseHardwareAddress()
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfHardwareAddr(iface, vf, addr)
		case "vlan":
			vlan, err := parseInt("VLANID")
			if err != nil {
				return err
			}

			if cursor == len(arg)-1 {
				return cmd.handle.LinkSetVfVlan(iface, vf, vlan)
			}

			cursor++
			expectedValues = []string{"qos"}
			switch arg[cursor] {
			case "qos":
				qos, err := parseInt("VLAN-QOS")
				if err != nil {
					return err
				}

				return cmd.handle.LinkSetVfVlanQos(iface, vf, vlan, qos)
			default:
				return usage()
			}
		case "rate":
			rate, err := parseInt("TXRATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfTxRate(iface, vf, rate)
		case "max_tx_rate":
			rate, err := parseInt("TXRATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfRate(iface, vf, int(iface.Attrs().Vfs[0].MinTxRate), rate)
		case "min_tx_rate":
			rate, err := parseInt("TXRATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfRate(iface, vf, rate, int(iface.Attrs().Vfs[0].MaxTxRate))
		case "state":
			state, err := parseUint32("STATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfState(iface, vf, state)
		case "spoofchk":
			check, err := parseBool()
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfSpoofchk(iface, vf, check)
		case "trust":
			trust, err := parseBool()
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfTrust(iface, vf, trust)
		case "node_guid":
			nodeguid, err := parseHardwareAddress()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfNodeGUID(iface, vf, nodeguid)

		case "port_guid":
			portguid, err := parseHardwareAddress()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfPortGUID(iface, vf, portguid)
		}
	}
	return usage()
}

func (cmd cmd) linkAdd() error {
	name, err := parseName()
	if err != nil {
		return err
	}
	attrs := netlink.LinkAttrs{Name: name}

	// Parse link attributes
	optionsDone := false

	for {
		if cursor == len(arg)-1 || optionsDone {
			break
		}

		cursor++
		expectedValues = []string{"type", "txqueuelen", "txqlen", "address", "mtu", "index", "numtxqueues", "numrxqueues"}
		switch arg[cursor] {
		case "txqueuelen", "txqlen":
			qlen, err := parseInt("PACKETS")
			if err != nil {
				return err
			}
			attrs.TxQLen = qlen
		case "address":
			hwAddr, err := parseHardwareAddress()
			if err != nil {
				return err
			}
			attrs.HardwareAddr = hwAddr
		case "mtu":
			mtu, err := parseInt("MTU")
			if err != nil {
				return err
			}
			attrs.MTU = mtu
		case "index":
			index, err := parseInt("IDX")
			if err != nil {
				return err
			}
			attrs.Index = index
		case "numtxqueues":
			numtxqueues, err := parseInt("QUEUE_COUNT")
			if err != nil {
				return err
			}

			attrs.NumTxQueues = numtxqueues
		case "numrxqueues":
			numrxqueues, err := parseInt("QUEUE_COUNT")
			if err != nil {
				return err
			}

			attrs.NumRxQueues = numrxqueues
		default:
			optionsDone = true
		}
	}

	cursor--
	typeName, err := parseType()
	if err != nil {
		return err
	}

	switch typeName {
	case "dummy":
		return cmd.handle.LinkAdd(&netlink.Dummy{LinkAttrs: attrs})
	case "ifb":
		return cmd.handle.LinkAdd(&netlink.Ifb{LinkAttrs: attrs})
	case "vlan":
		return cmd.handle.LinkAdd(&netlink.Vlan{LinkAttrs: attrs})
	case "macvlan":
		return cmd.handle.LinkAdd(&netlink.Macvlan{LinkAttrs: attrs})
	case "veth":
		return cmd.handle.LinkAdd(&netlink.Veth{LinkAttrs: attrs})
	case "vxlan":
		return cmd.handle.LinkAdd(&netlink.Vxlan{LinkAttrs: attrs})
	case "ipvlan":
		return cmd.handle.LinkAdd(&netlink.IPVlan{LinkAttrs: attrs})
	case "ipvtap":
		return cmd.handle.LinkAdd(&netlink.IPVtap{IPVlan: netlink.IPVlan{LinkAttrs: attrs}})
	case "bond":
		return cmd.handle.LinkAdd(netlink.NewLinkBond(attrs))
	case "geneve":
		return cmd.handle.LinkAdd(&netlink.Geneve{LinkAttrs: attrs})
	case "gretap":
		return cmd.handle.LinkAdd(&netlink.Gretap{LinkAttrs: attrs})
	case "ipip":
		return cmd.handle.LinkAdd(&netlink.Iptun{LinkAttrs: attrs})
	case "ip6tln":
		return cmd.handle.LinkAdd(&netlink.Ip6tnl{LinkAttrs: attrs})
	case "sit":
		return cmd.handle.LinkAdd(&netlink.Sittun{LinkAttrs: attrs})
	case "vti":
		return cmd.handle.LinkAdd(&netlink.Vti{LinkAttrs: attrs})
	case "gre":
		return cmd.handle.LinkAdd(&netlink.Gretun{LinkAttrs: attrs})
	case "vrf":
		cursor++
		expectedValues = []string{"table"}
		if arg[cursor] != "table" {
			return usage()
		}
		tableID, err := parseUint32("TABLE")
		if err != nil {
			return err
		}

		return cmd.handle.LinkAdd(&netlink.Vrf{LinkAttrs: attrs, Table: tableID})
	case "bridge":
		return cmd.handle.LinkAdd(&netlink.Bridge{LinkAttrs: attrs})
	case "xfrm":
		return cmd.handle.LinkAdd(&netlink.Xfrmi{LinkAttrs: attrs})
	case "ipoib":
		return cmd.handle.LinkAdd(&netlink.IPoIB{LinkAttrs: attrs})
	case "bareudp":
		return cmd.handle.LinkAdd(&netlink.BareUDP{LinkAttrs: attrs})
	default:
		return fmt.Errorf("unsupported link type %s", typeName)
	}
}

func (cmd cmd) linkDel() error {
	link, err := parseDeviceName(true)
	if err != nil {
		return err
	}

	return cmd.handle.LinkDel(link)
}

func (cmd cmd) linkShow() error {
	dev, err := parseDeviceName(false)
	if errors.Is(err, ErrNotFound) {
		return cmd.showAllLinks(false)
	}

	typeName, err := parseType()
	if errors.Is(err, ErrNotFound) {
		return cmd.showLink(dev, false)
	}

	return cmd.showLink(dev, false, typeName)
}

func (cmd cmd) link() error {
	if len(arg) == 1 {
		return cmd.linkShow()
	}

	cursor++
	expectedValues = []string{"show", "set", "add", "delete", "help"}
	argument := arg[cursor]

	switch findPrefix(argument, expectedValues) {
	case "show":
		return cmd.linkShow()
	case "set":
		return cmd.linkSet()
	case "add":
		return cmd.linkAdd()
	case "delete":
		return cmd.linkDel()
	case "help":
		fmt.Fprint(cmd.out, linkHelp)
		return nil
	default:
		return usage()
	}
}
