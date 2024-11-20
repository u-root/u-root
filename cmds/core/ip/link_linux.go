// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
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
		    type TYPE [ARGS]

	ip link delete { DEVICE | dev DEVICE  } 

	ip link set { DEVICE | dev DEVICE }
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

TYPE := { bareudp | bond |bridge | dummy | ifb | vxlan }

`

func (cmd *cmd) linkSet() error {
	iface, err := cmd.parseDeviceName(true)
	if err != nil {
		return err
	}

	for cmd.tokenRemains() {
		token := cmd.nextToken("address", "up", "down", "arp", "promisc", "multicast", "allmulticast", "mtu", "name", "alias", "vf", "master", "nomaster", "netns", "txqueuelen", "txqlen", "group")
		switch token {
		case "address":
			if err := cmd.setLinkHardwareAddress(iface); err != nil {
				return err
			}
		case "up":
			if err := cmd.handle.LinkSetUp(iface); err != nil {
				return fmt.Errorf("%v can't set it up: %w", iface.Attrs().Name, err)
			}
		case "down":
			if err := cmd.handle.LinkSetDown(iface); err != nil {
				return fmt.Errorf("%v can't set it down: %w", iface.Attrs().Name, err)
			}
		case "arp":
			switch cmd.nextToken("on", "off") {
			case "on":
				if err := cmd.handle.LinkSetARPOn(iface); err != nil {
					return fmt.Errorf("%v can't set arp on: %w", iface.Attrs().Name, err)
				}
			case "off":
				if err := cmd.handle.LinkSetARPOff(iface); err != nil {
					return fmt.Errorf("%v can't set arp off: %w", iface.Attrs().Name, err)
				}
			}
		case "promisc":
			switch cmd.nextToken("on", "off") {
			case "on":
				if err := cmd.handle.SetPromiscOn(iface); err != nil {
					return fmt.Errorf("%v can't set promisc on: %w", iface.Attrs().Name, err)
				}
			case "off":
				if err := cmd.handle.SetPromiscOff(iface); err != nil {
					return fmt.Errorf("%v can't set promisc off: %w", iface.Attrs().Name, err)
				}
			}
		case "multicast":
			switch cmd.nextToken("on", "off") {
			case "on":
				if err := cmd.handle.LinkSetMulticastOn(iface); err != nil {
					return fmt.Errorf("%v can't set multicast on: %w", iface.Attrs().Name, err)
				}
			case "off":
				if err := cmd.handle.LinkSetMulticastOff(iface); err != nil {
					return fmt.Errorf("%v can't set multicast off: %w", iface.Attrs().Name, err)
				}
			}
		case "allmulticast":
			switch cmd.nextToken("on", "off") {
			case "on":
				if err := cmd.handle.LinkSetAllmulticastOn(iface); err != nil {
					return fmt.Errorf("%v can't set allmulticast on: %w", iface.Attrs().Name, err)
				}
			case "off":
				if err := cmd.handle.LinkSetAllmulticastOff(iface); err != nil {
					return fmt.Errorf("%v can't set allmulticast off: %w", iface.Attrs().Name, err)
				}
			}
		case "mtu":
			if err := cmd.setLinkMTU(iface); err != nil {
				return fmt.Errorf("%v can't set mtu: %w", iface.Attrs().Name, err)
			}
		case "name":
			if err := cmd.setLinkName(iface); err != nil {
				return fmt.Errorf("%v can't set name: %w", iface.Attrs().Name, err)
			}
		case "alias":
			if err := cmd.setLinkAlias(iface); err != nil {
				return fmt.Errorf("%v can't set alias: %w", iface.Attrs().Name, err)
			}
		case "vf":
			if err := cmd.setLinkVf(iface); err != nil {
				return fmt.Errorf("%v can't set vf: %w", iface.Attrs().Name, err)
			}
		case "master":

			master, err := cmd.handle.LinkByName(cmd.nextToken("device name"))
			if err != nil {
				return err
			}

			if err := cmd.handle.LinkSetMaster(iface, master); err != nil {
				return fmt.Errorf("%v can't set master: %w", iface.Attrs().Name, err)
			}
		case "nomaster":
			if err := cmd.handle.LinkSetNoMaster(iface); err != nil {
				return fmt.Errorf("%v can't set no master: %w", iface.Attrs().Name, err)
			}
		case "netns":
			if err := cmd.setLinkNetns(iface); err != nil {
				return fmt.Errorf("%v can't set netns: %w", iface.Attrs().Name, err)
			}
		case "txqueuelen", "txqlen":
			if err := cmd.setLinkTxQLen(iface); err != nil {
				return fmt.Errorf("%v can't set txqueuelen: %w", iface.Attrs().Name, err)
			}
		}
	}

	return nil
}

func (cmd *cmd) setLinkHardwareAddress(iface netlink.Link) error {
	hwAddr, err := cmd.parseHardwareAddress()
	if err != nil {
		return err
	}

	err = cmd.handle.LinkSetHardwareAddr(iface, hwAddr)
	if err != nil {
		return fmt.Errorf("%v cant set mac addr %v: %w", iface.Attrs().Name, hwAddr, err)
	}

	return nil
}

func (cmd *cmd) setLinkMTU(iface netlink.Link) error {
	token := cmd.nextToken("MTU")

	mtu, err := strconv.Atoi(token)
	if err != nil {
		return fmt.Errorf("invalid mtu %v: %w", token, err)
	}

	return cmd.handle.LinkSetMTU(iface, mtu)
}

func (cmd *cmd) setLinkGroup(iface netlink.Link) error {
	token := cmd.nextToken("GROUP")

	group, err := strconv.Atoi(token)
	if err != nil {
		return fmt.Errorf("invalid group %v: %w", token, err)
	}

	return cmd.handle.LinkSetMTU(iface, group)
}

func (cmd *cmd) setLinkName(iface netlink.Link) error {
	return cmd.handle.LinkSetName(iface, cmd.nextToken("name"))
}

func (cmd *cmd) setLinkAlias(iface netlink.Link) error {
	return cmd.handle.LinkSetAlias(iface, cmd.nextToken("<alias name>"))
}

func (cmd *cmd) setLinkTxQLen(iface netlink.Link) error {
	token := cmd.nextToken("<qlen>")
	qlen, err := strconv.Atoi(token)
	if err != nil {
		return fmt.Errorf("invalid queuelen %v: %w", token, err)
	}

	return cmd.handle.LinkSetTxQLen(iface, qlen)
}

func (cmd *cmd) setLinkNetns(iface netlink.Link) error {
	token := cmd.nextToken("PID", "NAME")

	ns, err := strconv.Atoi(token)
	if err != nil {
		return fmt.Errorf("invalid int %v: %w", token, err)
	}

	if err := cmd.handle.LinkSetNsPid(iface, ns); err != nil {
		if err := cmd.handle.LinkSetNsFd(iface, ns); err != nil {
			return fmt.Errorf("failed to set netns: %w", err)
		}
	}

	return nil
}

func (cmd *cmd) setLinkVf(iface netlink.Link) error {
	vf, err := cmd.parseInt("VF")
	if err != nil {
		return err
	}

	for cmd.tokenRemains() {
		switch cmd.nextToken("vlan", "mac", "qos", "rate", "max_tx_rate", "min_tx_rate", "state", "spoofchk", "trust", "node_guid", "port_guid") {
		case "mac":
			addr, err := cmd.parseHardwareAddress()
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfHardwareAddr(iface, vf, addr)
		case "vlan":
			vlan, err := cmd.parseInt("VLANID")
			if err != nil {
				return err
			}

			if !cmd.tokenRemains() {
				return cmd.handle.LinkSetVfVlan(iface, vf, vlan)
			}

			switch cmd.nextToken("qos") {
			case "qos":
				qos, err := cmd.parseInt("VLAN-QOS")
				if err != nil {
					return err
				}

				return cmd.handle.LinkSetVfVlanQos(iface, vf, vlan, qos)
			default:
				return cmd.usage()
			}
		case "rate":
			rate, err := cmd.parseInt("TXRATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfTxRate(iface, vf, rate)
		case "max_tx_rate":
			rate, err := cmd.parseInt("TXRATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfRate(iface, vf, int(iface.Attrs().Vfs[0].MinTxRate), rate)
		case "min_tx_rate":
			rate, err := cmd.parseInt("TXRATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfRate(iface, vf, rate, int(iface.Attrs().Vfs[0].MaxTxRate))
		case "state":
			state, err := cmd.parseUint32("STATE")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfState(iface, vf, state)
		case "spoofchk":
			check, err := cmd.parseBool("on", "off")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfSpoofchk(iface, vf, check)
		case "trust":
			trust, err := cmd.parseBool("on", "off")
			if err != nil {
				return err
			}

			return cmd.handle.LinkSetVfTrust(iface, vf, trust)
		case "node_guid":
			nodeguid, err := cmd.parseHardwareAddress()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfNodeGUID(iface, vf, nodeguid)

		case "port_guid":
			portguid, err := cmd.parseHardwareAddress()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfPortGUID(iface, vf, portguid)
		}
	}
	return cmd.usage()
}

func (cmd *cmd) linkAdd() error {
	typeName, attrs, err := cmd.parseLinkAttrs()
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
		if cmd.nextToken("table") != "table" {
			return cmd.usage()
		}
		tableID, err := cmd.parseUint32("TABLE")
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

func (cmd *cmd) parseLinkAttrs() (string, netlink.LinkAttrs, error) {
	typeName := ""
	attrs := netlink.LinkAttrs{Name: cmd.parseName()}

	for cmd.tokenRemains() {
		switch cmd.nextToken("type", "txqueuelen", "txqlen", "address", "mtu", "index", "numtxqueues", "numrxqueues") {
		case "txqueuelen", "txqlen":
			qlen, err := cmd.parseInt("PACKETS")
			if err != nil {
				return "", netlink.LinkAttrs{}, err
			}
			attrs.TxQLen = qlen
		case "address":
			hwAddr, err := cmd.parseHardwareAddress()
			if err != nil {
				return "", netlink.LinkAttrs{}, err
			}
			attrs.HardwareAddr = hwAddr
		case "mtu":
			mtu, err := cmd.parseInt("MTU")
			if err != nil {
				return "", netlink.LinkAttrs{}, err
			}
			attrs.MTU = mtu
		case "index":
			index, err := cmd.parseInt("IDX")
			if err != nil {
				return "", netlink.LinkAttrs{}, err
			}
			attrs.Index = index
		case "numtxqueues":
			numtxqueues, err := cmd.parseInt("QUEUE_COUNT")
			if err != nil {
				return "", netlink.LinkAttrs{}, err
			}

			attrs.NumTxQueues = numtxqueues
		case "numrxqueues":
			numrxqueues, err := cmd.parseInt("QUEUE_COUNT")
			if err != nil {
				return "", netlink.LinkAttrs{}, err
			}

			attrs.NumRxQueues = numrxqueues
		case "type":
			typeName = cmd.nextToken("TYPE")
		default:
			return "", netlink.LinkAttrs{}, cmd.usage()
		}
	}

	if typeName == "" {
		return "", netlink.LinkAttrs{}, fmt.Errorf("type not specified")
	}

	return typeName, attrs, nil
}

func (cmd *cmd) linkDel() error {
	link, err := cmd.parseDeviceName(true)
	if err != nil {
		return err
	}

	return cmd.handle.LinkDel(link)
}

func (cmd *cmd) linkShow() error {
	dev, typeName, err := cmd.parseLinkShow()
	if err != nil {
		return err
	}

	if dev == nil {
		return cmd.showAllLinks(false, typeName...)
	}

	return cmd.showLink(dev, false, typeName...)
}

func (cmd *cmd) parseLinkShow() (netlink.Link, []string, error) {
	var (
		device netlink.Link
		err    error
	)

	typeNames := []string{}

	for cmd.tokenRemains() {
		switch c := cmd.nextToken("device name", "dev", "type"); c {
		default:
			device, err = netlink.LinkByName(c)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get link %v: %w", device, err)
			}
		case "dev":
			devName := cmd.nextToken("device name")
			device, err = netlink.LinkByName(devName)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get link %v: %w", device, err)
			}
		case "type":
			for cmd.tokenRemains() {
				if cmd.peekToken("dev") == "dev" {
					break
				}
				typeNames = append(typeNames, cmd.nextToken("type name"))
			}
		}
	}

	return device, typeNames, nil
}

func (cmd *cmd) link() error {
	if !cmd.tokenRemains() {
		return cmd.linkShow()
	}

	c := cmd.findPrefix("show", "set", "add", "delete", "help")
	switch c {
	case "show":
		return cmd.linkShow()
	case "set":
		return cmd.linkSet()
	case "add":
		return cmd.linkAdd()
	case "delete":
		return cmd.linkDel()
	case "help":
		fmt.Fprint(cmd.Out, linkHelp)
		return nil
	default:
		return cmd.usage()
	}
}
