// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
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

func linkSet() error {
	iface, err := parseDeviceName(true)
	if err != nil {
		return err
	}

	for cursor < len(arg)-1 {
		cursor++
		whatIWant = []string{"address", "up", "down", "arp", "promisc", "multicast", "allmulticast", "mtu", "name", "alias", "vf", "master", "nomaster", "netns", "txqueuelen", "txqlen", "group"}
		switch arg[cursor] {
		case "address":
			return setLinkHardwareAddress(iface)
		case "up":
			if err := netlink.LinkSetUp(iface); err != nil {
				return fmt.Errorf("%v can't make it up: %v", iface.Attrs().Name, err)
			}
		case "down":
			if err := netlink.LinkSetDown(iface); err != nil {
				return fmt.Errorf("%v can't make it down: %v", iface.Attrs().Name, err)
			}
		case "arp":
			cursor++
			whatIWant = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return netlink.LinkSetARPOn(iface)
			case "off":
				return netlink.LinkSetARPOff(iface)
			}
		case "promisc":
			cursor++
			whatIWant = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return netlink.SetPromiscOn(iface)
			case "off":
				return netlink.SetPromiscOff(iface)
			}
		case "multicast":
			cursor++
			whatIWant = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return netlink.LinkSetMulticastOn(iface)
			case "off":
				return netlink.LinkSetMulticastOff(iface)
			}
		case "allmulticast":
			cursor++
			whatIWant = []string{"on", "off"}
			switch arg[cursor] {
			case "on":
				return netlink.LinkSetAllmulticastOn(iface)
			case "off":
				return netlink.LinkSetAllmulticastOff(iface)
			}
		case "mtu":
			return setLinkMTU(iface)
		case "name":
			return setLinkName(iface)
		case "alias":
			return setLinkAlias(iface)
		case "vf":
			return setLinkVf(iface)
		case "master":
			cursor++
			whatIWant = []string{"device name"}
			master, err := netlink.LinkByName(arg[cursor])
			if err != nil {
				return err
			}
			return netlink.LinkSetMaster(iface, master)
		case "nomaster":
			return netlink.LinkSetNoMaster(iface)
		case "netns":
			return setLinkNetns(iface)
		case "txqueuelen", "txqlen":
			return setLinkTxQLen(iface)
		case "group":

		}
	}

	return nil
}

func setLinkHardwareAddress(iface netlink.Link) error {
	hwAddr, err := net.ParseMAC(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid mac address %v: %v", arg[cursor], err)
	}

	err = netlink.LinkSetHardwareAddr(iface, hwAddr)
	if err != nil {
		return fmt.Errorf("%v cant set mac addr %v: %v", iface.Attrs().Name, hwAddr, err)
	}

	return nil
}

func setLinkMTU(iface netlink.Link) error {
	cursor++

	mtu, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid mtu %v: %v", arg[cursor], err)
	}

	return netlink.LinkSetMTU(iface, mtu)
}

func setLinkGroup(iface netlink.Link) error {
	cursor++

	group, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid group %v: %v", arg[cursor], err)
	}

	return netlink.LinkSetMTU(iface, group)
}

func setLinkName(iface netlink.Link) error {
	cursor++
	whatIWant = []string{"<name>"}
	name := arg[cursor]

	return netlink.LinkSetName(iface, name)
}

func setLinkAlias(iface netlink.Link) error {
	cursor++
	whatIWant = []string{"<alias name>"}
	alias := arg[cursor]

	return netlink.LinkSetAlias(iface, alias)
}

func setLinkTxQLen(iface netlink.Link) error {
	cursor++
	whatIWant = []string{"<qlen>"}
	qlen, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid queuelen %v: %v", arg[cursor], err)
	}

	return netlink.LinkSetTxQLen(iface, qlen)
}

func setLinkNetns(iface netlink.Link) error {
	cursor++
	whatIWant = []string{"<netns pid>, <netns path>"}
	ns, err := strconv.Atoi(arg[cursor])
	if err != nil {
		return fmt.Errorf("invalid int %v: %v", arg[cursor], err)
	}

	if err := netlink.LinkSetNsPid(iface, ns); err != nil {
		if err := netlink.LinkSetNsFd(iface, ns); err != nil {
			return fmt.Errorf("failed to set netns: %v", err)
		}
	}

	return nil
}

func setLinkVf(iface netlink.Link) error {
	vf, err := parseInt()
	if err != nil {
		return err
	}

	cursor++

	whatIWant = []string{"vlan", "mac", "qos", "rate", "max_tx_rate", "min_tx_rate", "state", "spoofchk", "trust", "node_guid", "port_guid"}
	for cursor < len(arg)-1 {
		switch arg[cursor] {
		case "mac":
			addr, err := parseHardwareAddress()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfHardwareAddr(iface, vf, addr)
		case "vlan":
			vlan, err := parseInt()
			if err != nil {
				return err
			}

			if cursor == len(arg)-1 {
				return netlink.LinkSetVfVlan(iface, vf, vlan)
			}

			cursor++
			whatIWant = []string{"qos"}
			switch arg[cursor] {
			case "qos":
				qos, err := parseInt()
				if err != nil {
					return err
				}

				return netlink.LinkSetVfVlanQos(iface, vf, vlan, qos)
			default:
				return usage()
			}
		case "rate":
			rate, err := parseInt()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfTxRate(iface, vf, rate)
		case "max_tx_rate":
			rate, err := parseInt()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfRate(iface, vf, int(iface.Attrs().Vfs[0].MinTxRate), rate)
		case "min_tx_rate":
			rate, err := parseInt()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfRate(iface, vf, rate, int(iface.Attrs().Vfs[0].MaxTxRate))
		case "state":
			state, err := parseUint32()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfState(iface, vf, state)
		case "spoofchk":
			check, err := parseBool()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfSpoofchk(iface, vf, check)
		case "trust":
			trust, err := parseBool()
			if err != nil {
				return err
			}

			return netlink.LinkSetVfTrust(iface, vf, trust)
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

func linkAdd() error {
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
		whatIWant = []string{"type", "txqueuelen", "txqlen", "address", "mtu", "index", "numtxqueues", "numrxqueues"}
		switch arg[cursor] {
		case "txqueuelen", "txqlen":
			qlen, err := parseInt()
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
			mtu, err := parseInt()
			if err != nil {
				return err
			}
			attrs.MTU = mtu
		case "index":
			index, err := parseInt()
			if err != nil {
				return err
			}
			attrs.Index = index
		case "numtxqueues":
			numtxqueues, err := parseInt()
			if err != nil {
				return err
			}

			attrs.NumTxQueues = numtxqueues
		case "numrxqueues":
			numrxqueues, err := parseInt()
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
		return netlink.LinkAdd(&netlink.Dummy{LinkAttrs: attrs})
	case "ifb":
		return netlink.LinkAdd(&netlink.Ifb{LinkAttrs: attrs})
	case "vlan":
		return netlink.LinkAdd(&netlink.Vlan{LinkAttrs: attrs})
	case "macvlan":
		return netlink.LinkAdd(&netlink.Macvlan{LinkAttrs: attrs})
	case "veth":
		return netlink.LinkAdd(&netlink.Veth{LinkAttrs: attrs})
	case "vxlan":
		return netlink.LinkAdd(&netlink.Vxlan{LinkAttrs: attrs})
	case "ipvlan":
		return netlink.LinkAdd(&netlink.IPVlan{LinkAttrs: attrs})
	case "ipvtap":
		return netlink.LinkAdd(&netlink.IPVtap{IPVlan: netlink.IPVlan{LinkAttrs: attrs}})
	case "bond":
		return netlink.LinkAdd(netlink.NewLinkBond(attrs))
	case "geneve":
		return netlink.LinkAdd(&netlink.Geneve{LinkAttrs: attrs})
	case "gretap":
		return netlink.LinkAdd(&netlink.Gretap{LinkAttrs: attrs})
	case "ipip":
		return netlink.LinkAdd(&netlink.Iptun{LinkAttrs: attrs})
	case "ip6tln":
		return netlink.LinkAdd(&netlink.Ip6tnl{LinkAttrs: attrs})
	case "sit":
		return netlink.LinkAdd(&netlink.Sittun{LinkAttrs: attrs})
	case "vti":
		return netlink.LinkAdd(&netlink.Vti{LinkAttrs: attrs})
	case "gre":
		return netlink.LinkAdd(&netlink.Gretun{LinkAttrs: attrs})
	case "vrf":
		cursor++
		whatIWant = []string{"table"}
		if arg[cursor] != "table" {
			return usage()
		}
		tableID, err := parseUint32("TABLE")
		if err != nil {
			return err
		}

		return netlink.LinkAdd(&netlink.Vrf{LinkAttrs: attrs, Table: tableID})
	case "bridge":
		return netlink.LinkAdd(&netlink.Bridge{LinkAttrs: attrs})
	case "xfrm":
		return netlink.LinkAdd(&netlink.Xfrmi{LinkAttrs: attrs})
	case "ipoib":
		return netlink.LinkAdd(&netlink.IPoIB{LinkAttrs: attrs})
	case "bareudp":
		return netlink.LinkAdd(&netlink.BareUDP{LinkAttrs: attrs})
	default:
		return fmt.Errorf("unsupported link type %s", typeName)
	}
}

func linkDel() error {
	link, err := parseDeviceName(true)
	if err != nil {
		return err
	}

	return netlink.LinkDel(link)
}

func linkShow(w io.Writer) error {
	dev, err := parseDeviceName(false)
	if errors.Is(err, ErrNotFound) {
		return showAllLinks(w, false)
	}

	typeName, err := parseType()
	if errors.Is(err, ErrNotFound) {
		return showLink(w, dev, false)
	}

	return showLink(w, dev, false, typeName)
}

func link(w io.Writer) error {
	if len(arg) == 1 {
		return linkShow(w)
	}

	cursor++
	whatIWant = []string{"show", "set", "add", "delete", "help"}
	cmd := arg[cursor]

	switch findPrefix(cmd, whatIWant) {
	case "show":
		return linkShow(w)
	case "set":
		return linkSet()
	case "add":
		return linkAdd()
	case "delete":
		return linkDel()
	case "help":
		fmt.Fprint(w, linkHelp)
		return nil
	default:
		return usage()
	}
}
