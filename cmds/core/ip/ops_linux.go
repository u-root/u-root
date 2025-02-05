// Copyright 2012-2017 the u-root Authors. All rights reserved
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
			typeName:       link.Type(),
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

// printLinks prints the list of links. If addresses is not nil, it prints the
// the link's addresses as well.
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

		line = fmt.Sprintf("%d: %s: <%s> mtu %d master %s state %s group %s",
			link.attrs.Index, link.attrs.Name, p.flagsStr(link.attrs.Flags), link.attrs.MTU, link.masterName, p.operStateStr(link.attrs.OperState), p.groupStr(link.attrs.Group))
		fmt.Fprintln(p.out, line)

		line = fmt.Sprintf("    link/%s %s", link.attrs.EncapType, link.attrs.HardwareAddr)
		fmt.Fprintln(p.out, line)

		if p.withDetails {
			if line := p.deviceDetailsLine(link.specificDevice); line != "" {
				fmt.Fprintln(p.out, line)
			}
		}

		if p.withStats && link.attrs.Statistics != nil {
			stats := link.attrs.Statistics
			line = fmt.Sprintf("    RX: bytes %d packets %d errors %d dropped %d missed %d mcast %d",
				stats.RxBytes, stats.RxPackets, stats.RxErrors, stats.RxDropped, stats.RxMissedErrors, stats.Multicast)
			fmt.Fprintln(p.out, line)
			line = fmt.Sprintf("    TX: bytes %d packets %d errors %d dropped %d carrier %d collsns %d",
				stats.TxBytes, stats.TxPackets, stats.TxErrors, stats.TxDropped, stats.TxCarrierErrors, stats.Collisions)
			fmt.Fprintln(p.out, line)
		}

		if p.withAddresses {
			for _, addr := range link.addresses {
				line = fmt.Sprintf("    %s %s brd %s scope %s %s",
					p.ipNetStr(addr), addr.IP, addr.Broadcast, addrScopeStr(netlink.Scope(addr.Scope)), addr.Label)
				fmt.Fprintln(p.out, line)
				line = fmt.Sprintf("       valid_lft %s preferred_lft %s",
					p.lifetimeStr(addr.ValidLft), p.lifetimeStr(addr.PreferedLft))
				fmt.Fprintln(p.out, line)
			}
		}

	}
}

func (p *linkPrinter) printBrief() {
	for _, link := range p.data {
		if link.attrs == nil {
			continue
		}

		var line string

		line = fmt.Sprintf("%-25s %-10s", link.attrs.Name, p.operStateStr(link.attrs.OperState))
		if p.withAddresses {
			for _, addr := range link.addresses {
				line += fmt.Sprintf(" %s", addr.IP)
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
	if uint32(lft) == math.MaxUint32 {
		return "forever"
	}

	return fmt.Sprintf("%dsec", lft)
}

func (p *linkPrinter) deviceDetailsLine(t any) string {
	var line string

	switch dev := t.(type) {
	case *netlink.Bridge:
		line = fmt.Sprintf("    bridge hello_time %d ageing_time %d vlan_filtering %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.HelloTime, dev.AgeingTime, dev.VlanFiltering, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
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
		line = fmt.Sprintf("    bond mode active slave %d %d miimon %d updelay %d downdelay %d use_carrier %d arp_interval %d arp_validate %s arp_all_targets %s primary %d primary_reselect %s fail_over_mac %s %s resend_igmp %d num_peer_notif %d all_slaves_active %d min_links %d lp_interval %d packets_per_slave %d lacp_rate %d ad_select %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.Mode, dev.ActiveSlave, dev.Miimon, dev.UpDelay, dev.DownDelay, dev.UseCarrier, dev.ArpInterval, dev.ArpValidate, dev.ArpAllTargets, dev.Primary, dev.PrimaryReselect, dev.FailOverMac, dev.XmitHashPolicy, dev.ResendIgmp, dev.NumPeerNotif, dev.AllSlavesActive, dev.MinLinks, dev.LpInterval, dev.PacketsPerSlave, dev.LacpRate, dev.AdSelect, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Geneve:
		line = fmt.Sprintf("    geneve id %d remote %s ttl %d tos %d dport %d udpcsum %d udp_zero_csum_6TX %d udp_zero_csum_6RX %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.ID, dev.Remote, dev.Ttl, dev.Tos, dev.Dport, dev.UdpCsum, dev.UdpZeroCsum6Tx, dev.UdpZeroCsum6Rx, dev.Link, dev.FlowBased, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Gretap:
		line = fmt.Sprintf("    gretap i_key %d o_key %d encap_src_port %d encap_dst_port %d local %s remote %s iflags %d oflags %d pmtudisc %d ttl %d tos %d encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d",
			dev.IKey, dev.OKey, dev.EncapSport, dev.EncapDport, dev.Local, dev.Remote, dev.IFlags, dev.OFlags, dev.PMtuDisc, dev.Ttl, dev.Tos, dev.EncapType, dev.EncapFlags, dev.Link, dev.FlowBased, dev.NumTxQueues, dev.NumRxQueues, dev.GSOMaxSize, dev.GSOMaxSegs)
	case *netlink.Iptun:
		line = fmt.Sprintf("    iptun local %s remote %s encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
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

				addrInfo := AddrJSON{
					Local:     addr.IPNet.IP.String(),
					PrefixLen: addr.IPNet.Mask.String(),
				}

				if !cmd.Opts.Brief {
					if addr.Broadcast != nil {
						addrInfo.Family = family
						addrInfo.Scope = addrScopeStr(netlink.Scope(addr.Scope))
						addrInfo.Label = addr.Label
						addrInfo.ValidLifeTime = fmt.Sprintf("%dsec", addr.ValidLft)
						addrInfo.PreferredLifeTime = fmt.Sprintf("%dsec", addr.PreferedLft)
					}

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
