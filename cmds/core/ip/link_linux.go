// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"

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
		[ group GROUP ]
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

	ip link show [ DEVICE ] [type TYPE]

	ip link help

TYPE := { bareudp | bond |bridge | dummy | ifb | vxlan }

`

// link is the entry point for 'ip link' command.
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

// linkSet performs 'ip link set' command.
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
		case "group":
			if err := cmd.setLinkGroup(iface); err != nil {
				return fmt.Errorf("%v can't set group: %w", iface.Attrs().Name, err)
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

	return cmd.handle.LinkSetGroup(iface, group)
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

// linkAdd performs 'ip link add' command.
func (cmd *cmd) linkAdd() error {
	typeName, attrs, err := cmd.parseLinkAdd()
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

// parseLinkAdd returns arguments to 'ip link add' from the cmdline.
func (cmd *cmd) parseLinkAdd() (string, netlink.LinkAttrs, error) {
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

// linkDel performs 'ip link delete' command.
func (cmd *cmd) linkDel() error {
	link, err := cmd.parseDeviceName(true)
	if err != nil {
		return err
	}

	return cmd.handle.LinkDel(link)
}

// linkShow performs 'ip link show' command.
func (cmd *cmd) linkShow() error {
	name, types := cmd.parseLinkShow()

	links, err := cmd.getLinkDevices(false, linkNameFilter([]string{name}), linkTypeFilter(types))
	if err != nil {
		return fmt.Errorf("link show: %w", err)
	}

	err = cmd.printLinks(false, links)
	if err != nil {
		return fmt.Errorf("link show: %w", err)
	}

	return nil
}

// parseLinkShow returns arguments to 'ip link show' from the cmdline.
func (cmd *cmd) parseLinkShow() (name string, types []string) {
	devSeen := false

	for cmd.tokenRemains() {
		switch c := cmd.nextToken("DEVICE", "dev", "type"); c {
		case "dev":
			if !devSeen {
				name = cmd.nextToken("DEVICE")
				devSeen = true
			}
		case "type":
			for cmd.tokenRemains() {
				if cmd.peekToken("dev") == "dev" {
					break
				}
				types = append(types, cmd.nextToken("TYPE"))
			}
		default:
			if name != "" {
				continue // ignore multiple link device names, taking the first one only
			}
			name = c
		}
	}

	return
}

// linkData contains information about a link device.
type linkData struct {
	attrs      *netlink.LinkAttrs
	typeName   string
	masterName string
	addresses  []netlink.Addr
	// concreteDevice can be any of the netlink.Link types. It can be casted to access the type-specific fields.
	specificDevice any
}

// getLinkDevices performs system I/O  to enumerate a list of link devices that match the given filters.
// If no error occurs, at least one linke device was found and the returned linkData objects have a non-nil attrs field.
// The addresses field is only populated if withAddresses is true.
func (cmd *cmd) getLinkDevices(withAddresses bool, filter ...linkfilter) ([]linkData, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("get link device list: %w", err)
	}

	links = filterLinks(links, filter)

	list := make([]linkData, 0, len(links))

	for _, link := range links {
		attrs := link.Attrs()
		if attrs == nil {
			return nil, errors.New("link device found, but does not provide any attributes")
		}

		// Get the master device name.
		var masterName string
		if attrs.MasterIndex != 0 {
			master, err := netlink.LinkByIndex(attrs.MasterIndex)
			if err != nil {
				return nil, fmt.Errorf("get master device for link device %q: %w", attrs.Name, err)
			}
			masterName = master.Attrs().Name
		}

		// Get the addresses for the link.
		var addresses []netlink.Addr
		if withAddresses {
			addresses, err = netlink.AddrList(link, cmd.Family)
			if err != nil {
				return nil, fmt.Errorf("get addresses for link device %q: %w", attrs.Name, err)
			}
		}

		list = append(list, linkData{
			attrs:          attrs,
			masterName:     masterName,
			addresses:      addresses,
			specificDevice: link,
		})
	}

	return list, nil
}

type linkfilter func(link netlink.Link) bool

// linkTypeFilter returns a linkfilter that filters links by type.
func linkTypeFilter(linkTypes []string) linkfilter {
	return func(link netlink.Link) bool {
		if len(linkTypes) == 0 {
			return true
		}
		for _, linkType := range linkTypes {
			if linkType == "" || link.Type() == linkType {
				return true
			}
		}
		return false
	}
}

// linkNameFilter returns a linkfilter that filters links by name.
func linkNameFilter(linkNames []string) linkfilter {
	return func(link netlink.Link) bool {
		if len(linkNames) == 0 {
			return true
		}
		for _, linkName := range linkNames {
			if linkName == "" || link.Attrs().Name == linkName {
				return true
			}
		}
		return false
	}
}

// filterLinks applies the given filters to the list of links.
func filterLinks(links []netlink.Link, lf []linkfilter) []netlink.Link {
	var (
		unfilteredLinks []netlink.Link
		filteredLinks   []netlink.Link
	)

	filteredLinks = links

	if len(lf) > 0 {
		for _, filter := range lf {
			unfilteredLinks = filteredLinks
			filteredLinks = make([]netlink.Link, 0)
			for _, link := range unfilteredLinks {
				if filter(link) {
					filteredLinks = append(filteredLinks, link)
				}
			}
		}
	}

	return filteredLinks
}

// printLinks prints the list of links. If withAddresses is true, the link's IP
// addresses are printed as well.
func (cmd *cmd) printLinks(withAddresses bool, links []linkData) error {
	if cmd.Opts.JSON {
		return cmd.printLinkJSON(links)
	}

	p := linkPrinter{
		out:           cmd.Out,
		data:          links,
		withStats:     cmd.Opts.Stats,
		withDetails:   cmd.Opts.Details,
		withAddresses: withAddresses,
		numeric:       cmd.Opts.Numeric,
	}

	switch {
	case cmd.Opts.Brief:
		p.printBrief()
	case cmd.Opts.Oneline:
		p.printOneline()
	default:
		p.printDefault()
	}

	return nil
}

type linkPrinter struct {
	out           io.Writer
	data          []linkData
	withStats     bool
	withDetails   bool
	withAddresses bool
	numeric       bool
}

func (p *linkPrinter) printDefault() {
	for _, link := range p.data {
		if link.attrs == nil {
			continue
		}

		var line string

		line = fmt.Sprintf("%d: %s: <%s> mtu %d", link.attrs.Index, link.attrs.Name, p.flagsStr(link.attrs.Flags), link.attrs.MTU)
		if link.masterName != "" {
			line += fmt.Sprintf(" master %s", link.masterName)
		}
		line += fmt.Sprintf(" state %s group %s", p.operStateStr(link.attrs.OperState), p.groupStr(link.attrs.Group))
		fmt.Fprintln(p.out, line)

		line = fmt.Sprintf("    link/%s %s", link.attrs.EncapType, link.attrs.HardwareAddr)
		fmt.Fprintln(p.out, line)

		if p.withDetails {
			if line := p.deviceDetailsLine(link.specificDevice); line != "" {
				fmt.Fprintln(p.out, line)
			}
		}

		if p.withAddresses {
			for _, addr := range link.addresses {
				line = fmt.Sprintf("    %s %s", p.ipNetStr(addr), addr.IPNet)
				if addr.IP.To4() != nil {
					line += fmt.Sprintf(" brd %s", addr.Broadcast)
				}
				line += fmt.Sprintf(" scope %s %s", addrScopeStr(netlink.Scope(addr.Scope)), addr.Label)
				fmt.Fprintln(p.out, line)

				line = fmt.Sprintf("       valid_lft %s preferred_lft %s", p.lifetimeStr(addr.ValidLft), p.lifetimeStr(addr.PreferedLft))
				fmt.Fprintln(p.out, line)
			}
		}

		if p.withStats && link.attrs.Statistics != nil {
			stats := link.attrs.Statistics
			fmt.Fprintln(p.out, "    RX:  bytes  packets errors dropped  missed   mcast")
			line = fmt.Sprintf("%14d %7d %6d %7d %7d %7d",
				stats.RxBytes, stats.RxPackets, stats.RxErrors, stats.RxDropped, stats.RxMissedErrors, stats.Multicast)
			fmt.Fprintln(p.out, line)

			fmt.Fprintln(p.out, "    TX:  bytes  packets errors dropped carrier collsns")
			line = fmt.Sprintf("%14d %7d %6d %7d %7d %7d",
				stats.TxBytes, stats.TxPackets, stats.TxErrors, stats.TxDropped, stats.TxCarrierErrors, stats.Collisions)
			fmt.Fprintln(p.out, line)
		}
	}
}

func (p *linkPrinter) printOneline() {
	for _, link := range p.data {
		if link.attrs == nil {
			continue
		}

		var line string

		if p.withAddresses {
			for _, addr := range link.addresses {
				line = fmt.Sprintf("%d: %s    %s %s", link.attrs.Index, link.attrs.Name, p.ipNetStr(addr), addr.IPNet)
				if addr.IP.To4() != nil {
					line += fmt.Sprintf(" brd %s", addr.Broadcast)
				}
				line += fmt.Sprintf(" scope %s %s", addrScopeStr(netlink.Scope(addr.Scope)), addr.Label)
				line += fmt.Sprintf("\\       valid_lft %s preferred_lft %s", p.lifetimeStr(addr.ValidLft), p.lifetimeStr(addr.PreferedLft))
				fmt.Fprintln(p.out, line)
			}
		} else {
			line = fmt.Sprintf("%d: %s: <%s> mtu %d", link.attrs.Index, link.attrs.Name, p.flagsStr(link.attrs.Flags), link.attrs.MTU)
			if link.masterName != "" {
				line += fmt.Sprintf(" master %s", link.masterName)
			}
			line += fmt.Sprintf(" state %s group %s", p.operStateStr(link.attrs.OperState), p.groupStr(link.attrs.Group))
			line += fmt.Sprintf("\\    link/%s %s", link.attrs.EncapType, link.attrs.HardwareAddr)

			if p.withDetails {
				if details := p.deviceDetailsLine(link.specificDevice); details != "" {
					line += fmt.Sprintf("\\ %s", details)
				}
			}

			if p.withStats && link.attrs.Statistics != nil {
				stats := link.attrs.Statistics
				line += "\\    RX:  bytes  packets errors dropped  missed   mcast"
				line += fmt.Sprintf("\\%14d %7d %6d %7d %7d %7d",
					stats.RxBytes, stats.RxPackets, stats.RxErrors, stats.RxDropped, stats.RxMissedErrors, stats.Multicast)

				line += "\\    TX:  bytes  packets errors dropped carrier collsns"
				line += fmt.Sprintf("\\%14d %7d %6d %7d %7d %7d",
					stats.TxBytes, stats.TxPackets, stats.TxErrors, stats.TxDropped, stats.TxCarrierErrors, stats.Collisions)
			}
			fmt.Fprintln(p.out, line)
		}
	}
}

func (p *linkPrinter) printBrief() {
	for _, link := range p.data {
		if link.attrs == nil {
			continue
		}

		var line string

		line = fmt.Sprintf("%-20s %-10s", link.attrs.Name, p.operStateStr(link.attrs.OperState))
		if p.withAddresses {
			for _, addr := range link.addresses {
				line += fmt.Sprintf(" %s", addr.IPNet)
			}
		} else {
			line += fmt.Sprintf(" %-20s <%s>", link.attrs.HardwareAddr, p.flagsStr(link.attrs.Flags))
		}
		fmt.Fprintln(p.out, line)
	}
}

func (p *linkPrinter) groupStr(group uint32) string {
	if group == 0 && !p.numeric {
		return "default"
	}

	return strconv.Itoa(int(group))
}

func (p *linkPrinter) flagsStr(flags net.Flags) string {
	return strings.Replace(strings.ToUpper(flags.String()), "|", ",", -1)
}

func (p *linkPrinter) operStateStr(state netlink.LinkOperState) string {
	return strings.ToUpper(state.String())
}

func (p *linkPrinter) ipNetStr(addr netlink.Addr) string {
	if addr.IPNet != nil {
		if addr.IPNet.IP.To4() != nil {
			return "inet"
		}
		if addr.IPNet.IP.To16() != nil {
			return "inet6"
		}
	}

	return ""
}

func (p *linkPrinter) lifetimeStr(lft int) string {
	// fix vishnavanda/netlink. *Lft should be uint32, not int.
	if uint32(lft) == math.MaxUint32 {
		return "forever"
	}

	return fmt.Sprintf("%dsec", lft)
}

func (p *linkPrinter) deviceDetailsLine(t any) string {
	var line string

	b2int := func(b bool) int8 {
		if b {
			return 1
		}
		return 0
	}

	switch dev := t.(type) {
	case *netlink.Bridge:
		line = fmt.Sprintf("    bridge hello_time %d ageing_time %d vlan_filtering %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			*dev.HelloTime, *dev.AgeingTime, b2int(*dev.VlanFiltering), dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Vlan:
		line = fmt.Sprintf("    vlan %s vlan-id %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.VlanProtocol, dev.VlanId, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Macvlan:
		line = fmt.Sprintf("    macvlan mode %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Macvtap:
		line = fmt.Sprintf("    macvtap mode %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Tuntap:
		line = fmt.Sprintf("    tuntap mode %s owner %d group %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.Owner, dev.Group, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Veth:
		line = fmt.Sprintf("    peer %s peer-address %s numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.PeerName, dev.PeerHardwareAddr, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Vxlan:
		line = fmt.Sprintf("    vxlan id %d src %s group %s ttl %d tos %d learning %t proxy %t rsc %t age %d limit %d port %d port-low %d port-high %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.VxlanId, dev.SrcAddr, dev.Group, dev.TTL, dev.TOS, dev.Learning, dev.Proxy, dev.RSC, dev.Age, dev.Limit, dev.Port, dev.PortLow, dev.PortHigh, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.IPVlan:
		line = fmt.Sprintf("    ipvlan mode %d flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.Flags, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.IPVtap:
		line = fmt.Sprintf("    ipvtap mode %d flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.Flags, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Bond:
		line = fmt.Sprintf("    bond mode active slave %d %d miimon %d updelay %d downdelay %d use_carrier %d arp_interval %d arp_validate %s arp_all_targets %s primary %d primary_reselect %s fail_over_mac %s %s resend_igmp %d num_peer_notif %d all_slaves_active %d min_links %d lp_interval %d packets_per_slave %d lacp_rate %s ad_select %s numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.ActiveSlave, dev.Miimon, dev.UpDelay, dev.DownDelay, dev.UseCarrier, dev.ArpInterval, dev.ArpValidate, dev.ArpAllTargets, dev.Primary, dev.PrimaryReselect, dev.FailOverMac, dev.XmitHashPolicy, dev.ResendIgmp, dev.NumPeerNotif, dev.AllSlavesActive, dev.MinLinks, dev.LpInterval, dev.PacketsPerSlave, dev.LacpRate, dev.AdSelect, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Geneve:
		line = fmt.Sprintf("    geneve id %d remote %s ttl %d tos %d dport %d udpcsum %d udp_zero_csum_6TX %d udp_zero_csum_6RX %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.ID, dev.Remote, dev.Ttl, dev.Tos, dev.Dport, dev.UdpCsum, dev.UdpZeroCsum6Tx, dev.UdpZeroCsum6Rx, dev.Link, dev.FlowBased, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Gretap:
		line = fmt.Sprintf("    gretap i_key %d o_key %d encap_src_port %d encap_dst_port %d local %s remote %s iflags %d oflags %d pmtudisc %d ttl %d tos %d encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.IKey, dev.OKey, dev.EncapSport, dev.EncapDport, dev.Local, dev.Remote, dev.IFlags, dev.OFlags, dev.PMtuDisc, dev.Ttl, dev.Tos, dev.EncapType, dev.EncapFlags, dev.Link, dev.FlowBased, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Iptun:
		line = fmt.Sprintf("    iptun local %s remote %s encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Local, dev.Remote, dev.EncapType, dev.EncapFlags, dev.Link, dev.FlowBased, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Ip6tnl:
		line = fmt.Sprintf("    ip6tnl local %s remote %s ttl %d tos %d proto %d flow_info %d encap_limit %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Local, dev.Remote, dev.Ttl, dev.Tos, dev.Proto, dev.FlowInfo, dev.EncapLimit, dev.EncapType, dev.EncapSport, dev.EncapDport, dev.EncapFlags, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Sittun:
		line = fmt.Sprintf("    sittun local %s remote %s ttl %d tos %d proto %d encap_limit %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Local, dev.Remote, dev.Ttl, dev.Tos, dev.Proto, dev.EncapLimit, dev.EncapType, dev.EncapSport, dev.EncapDport, dev.EncapFlags, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Vti:
		line = fmt.Sprintf("    vti local %s remote %s ikey %d okey %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Local, dev.Remote, dev.IKey, dev.OKey, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Gretun:
		line = fmt.Sprintf("    gretun local %s remote %s ttl %d tos %d ptmudisc %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d ikey %d okey %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Local, dev.Remote, dev.Ttl, dev.Tos, dev.PMtuDisc, dev.EncapType, dev.EncapSport, dev.EncapDport, dev.EncapFlags, dev.IKey, dev.OKey, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Xfrmi:
		line = fmt.Sprintf("    xfrmi if_id %d", dev.Ifid)
	case *netlink.Can:
		line = fmt.Sprintf("    can state %d bitrate %d sample-point %d tq %d prop-seg %d phase-seg1 %d phase-seg2 %d",
			dev.State, dev.BitRate, dev.SamplePoint, dev.TimeQuanta, dev.PropagationSegment, dev.PhaseSegment1, dev.PhaseSegment2)
	case *netlink.IPoIB:
		line = fmt.Sprintf("    ipoib pkey %d mode %d umcast %d", dev.Pkey, dev.Mode, dev.Umcast)
	case *netlink.BareUDP:
		line = fmt.Sprintf("    port %d ethertype %d srcport %d min multi_proto %t", dev.Port, dev.EtherType, dev.SrcPortMin, dev.MultiProto)
	}

	return line
}

// LinkJSON represents a link device for JSON output format.
type LinkJSON struct {
	IfIndex   int        `json:"ifindex,omitempty"`
	IfName    string     `json:"ifname"`
	Flags     []string   `json:"flags"`
	MTU       int        `json:"mtu,omitempty"`
	Operstate string     `json:"operstate"`
	Group     string     `json:"group,omitempty"`
	Txqlen    int        `json:"txqlen,omitempty"`
	LinkType  string     `json:"link_type,omitempty"`
	Address   string     `json:"address"`
	AddrInfo  []AddrJSON `json:"addr_info,omitempty"`
}

// AddrInfo represents an address for JSON output format.
type AddrJSON struct {
	Family            string `json:"ip,omitempty"`
	Local             string `json:"local"`
	PrefixLen         string `json:"prefixlen"`
	Broadcast         string `json:"broadcast,omitempty"`
	Scope             string `json:"scope,omitempty"`
	Label             string `json:"label,omitempty"`
	ValidLifeTime     string `json:"valid_life_time,omitempty"`
	PreferredLifeTime string `json:"preferred_life_time,omitempty"`
}

func (cmd *cmd) printLinkJSON(links []linkData) error {
	linkObs := make([]LinkJSON, 0)

	for _, v := range links {
		link := LinkJSON{
			IfName:    v.attrs.Name,
			Flags:     strings.Split(v.attrs.Flags.String(), "|"),
			Operstate: v.attrs.OperState.String(),
			Address:   v.attrs.HardwareAddr.String(),
		}

		if !cmd.Opts.Brief {
			link.IfIndex = v.attrs.Index
			link.MTU = v.attrs.MTU
			link.LinkType = v.typeName
			link.Group = fmt.Sprintf("%v", v.attrs.Group)

			if !cmd.Opts.Numeric && v.attrs.Group == 0 {
				link.Group = "default"
			}

			link.Txqlen = v.attrs.TxQLen
		}

		if v.addresses != nil {
			link.AddrInfo = make([]AddrJSON, 0)

			for _, addr := range v.addresses {

				family := "inet"
				if addr.IP.To4() == nil {
					family = "inet6"
				}

				ip := strings.Split(addr.IPNet.String(), "/")[0]
				prefixlen := strings.Split(addr.IPNet.String(), "/")[1]

				addrInfo := AddrJSON{
					Local:     ip,
					PrefixLen: prefixlen,
				}

				if !cmd.Opts.Brief {
					addrInfo.Family = family
					addrInfo.Scope = addrScopeStr(netlink.Scope(addr.Scope))
					addrInfo.Label = addr.Label
					addrInfo.ValidLifeTime = fmt.Sprintf("%dsec", addr.ValidLft)
					addrInfo.PreferredLifeTime = fmt.Sprintf("%dsec", addr.PreferedLft)

					if addr.Broadcast != nil {
						addrInfo.Broadcast = addr.Broadcast.String()
					}
				}

				link.AddrInfo = append(link.AddrInfo, addrInfo)

			}
		}
		linkObs = append(linkObs, link)
	}

	return printJSON(*cmd, linkObs)
}
