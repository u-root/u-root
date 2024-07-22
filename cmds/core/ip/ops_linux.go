// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/vishvananda/netlink"
)

func showAllLinks(w io.Writer, withAddresses bool, filterByType ...string) error {
	links, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("can't enumerate interfaces: %v", err)
	}
	return showLinks(w, withAddresses, links, filterByType...)
}

func showLink(w io.Writer, link netlink.Link, withAddresses bool, filterByType ...string) error {
	return showLinks(w, withAddresses, []netlink.Link{link}, filterByType...)
}

type Link struct {
	IfIndex   int        `json:"ifindex,omitempty"`
	IfName    string     `json:"ifname"`
	Flags     []string   `json:"flags"`
	MTU       int        `json:"mtu,omitempty"`
	Operstate string     `json:"operstate"`
	Group     uint32     `json:"group,omitempty"`
	Txqlen    int        `json:"txqlen,omitempty"`
	LinkType  string     `json:"link_type,omitempty"`
	Address   string     `json:"address"`
	AddrInfo  []AddrInfo `json:"addr_info,omitempty"`
}

type AddrInfo struct {
	Family            string `json:"ip,omitempty"`
	Local             string `json:"local"`
	PrefixLen         string `json:"prefixlen"`
	Broadcast         string `json:"broadcast,omitempty"`
	Scope             string `json:"scope,omitempty"`
	Label             string `json:"label,omitempty"`
	ValidLifeTime     string `json:"valid_life_time,omitempty"`
	PreferredLifeTime string `json:"preferred_life_time,omitempty"`
}

func showLinks(w io.Writer, withAddresses bool, links []netlink.Link, filterByType ...string) error {
	if f.json {
		linkObs := make([]Link, 0)

		for _, v := range links {
			link := Link{
				IfName:    v.Attrs().Name,
				Flags:     strings.Split(v.Attrs().Flags.String(), "|"),
				Operstate: v.Attrs().OperState.String(),
				Address:   v.Attrs().HardwareAddr.String(),
			}

			if !f.brief {
				link.IfIndex = v.Attrs().Index
				link.MTU = v.Attrs().MTU
				link.LinkType = v.Type()
				link.Group = v.Attrs().Group
				link.Txqlen = v.Attrs().TxQLen
			}

			if withAddresses {

				addrs, err := netlink.AddrList(v, family)
				if err != nil {
					return fmt.Errorf("can't enumerate addresses: %v", err)
				}

				link.AddrInfo = make([]AddrInfo, 0)
				for _, addr := range addrs {
					family := "inet"
					if addr.IP.To4() == nil {
						family = "inet6"
					}

					addrInfo := AddrInfo{
						Local:     addr.IPNet.IP.String(),
						PrefixLen: addr.IPNet.Mask.String(),
					}

					if !f.brief {
						if addr.Broadcast != nil {
							addrInfo.Family = family
							addrInfo.Scope = addrScopes[netlink.Scope(addr.Scope)]
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

		return printJSON(w, linkObs)
	}

	for _, v := range links {
		if withAddresses {

			addrs, err := netlink.AddrList(v, family)
			if err != nil {
				return fmt.Errorf("can't enumerate addresses: %v", err)
			}

			// if there are no addresses and the link is not a vrf (only wihout -4 or -6), skip it
			if len(addrs) == 0 && (v.Type() != "vrf" || family != netlink.FAMILY_ALL) {
				continue
			}
		}

		found := true

		// check if the link type is in the filter list if the filter list is not empty
		if len(filterByType) > 0 {
			found = false
		}

		for _, t := range filterByType {
			if v.Type() == t {
				found = true
			}
		}

		if !found {
			continue
		}

		l := v.Attrs()

		if f.brief {
			if withAddresses {
				addrs, err := netlink.AddrList(v, family)
				if err != nil {
					return fmt.Errorf("can't enumerate addresses: %v", err)
				}

				fmt.Fprintf(w, "%-20s %-10s", l.Name, l.OperState.String())

				for _, addr := range addrs {
					fmt.Fprintf(w, " %s", addr.IP)
				}

				fmt.Fprintf(w, "\n")

				continue
			}

			addr := " "
			if l.HardwareAddr != nil {
				addr = fmt.Sprintf(" %s ", l.HardwareAddr.String())
			}

			fmt.Fprintf(w, "%-25s %-10s%-20s <%s>\n", l.Name,
				l.OperState.String(), addr, strings.Replace(strings.ToUpper(l.Flags.String()), "|", ",", -1))

			continue
		}

		master := ""
		if l.MasterIndex != 0 {
			link, err := netlink.LinkByIndex(l.MasterIndex)
			if err != nil {
				return fmt.Errorf("can't get link with index %d: %v", l.MasterIndex, err)
			}
			master = fmt.Sprintf("master %s ", link.Attrs().Name)
		}

		fmt.Fprintf(w, "%d: %s: <%s> mtu %d %sstate %s\n", l.Index, l.Name,
			strings.Replace(strings.ToUpper(l.Flags.String()), "|", ",", -1),
			l.MTU, master, strings.ToUpper(l.OperState.String()))

		fmt.Fprintf(w, "    link/%s %s\n", l.EncapType, l.HardwareAddr)

		if f.details {
			switch v := v.(type) {
			case *netlink.Bridge:
				fmt.Fprintf(w, "    bridge hello_time %d ageing_time %d vlan_filtering %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.HelloTime, v.AgeingTime, v.VlanFiltering, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Vlan:
				fmt.Fprintf(w, "    vlan %s vlan-id %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.VlanProtocol, v.VlanId, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Macvlan:
				fmt.Fprintf(w, "    macvlan mode %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Macvtap:
				fmt.Fprintf(w, "    macvtap mode %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Tuntap:
				fmt.Fprintf(w, "    %s owner %d group %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.Owner, v.Group, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Veth:
				fmt.Fprintf(w, "    peer %s peer-address %s numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.PeerName, v.PeerHardwareAddr, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Vxlan:
				fmt.Fprintf(w, "    vxlan id %d src %s group %s ttl %d tos %d learning %t proxy %t rsc %t age %d limit %d port %d port-low %d port-high %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.VxlanId, v.SrcAddr, v.Group, v.TTL, v.TOS, v.Learning, v.Proxy, v.RSC, v.Age, v.Limit, v.Port, v.PortLow, v.PortHigh, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.IPVlan:
				fmt.Fprintf(w, "    ipvlan mode %d flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.Flags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.IPVtap:
				fmt.Fprintf(w, "    ipvtap mode %d flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.Flags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Bond:
				fmt.Fprintf(w, "    bond mode active slave %d %d miimon %d updelay %d downdelay %d use_carrier %d arp_interval %d arp_validate %s arp_all_targets %s primary %d primary_reselect %s fail_over_mac %s %s resend_igmp %d num_peer_notif %d all_slaves_active %d min_links %d lp_interval %d packets_per_slave %d lacp_rate %d ad_select %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.ActiveSlave, v.Miimon, v.UpDelay, v.DownDelay, v.UseCarrier, v.ArpInterval, v.ArpValidate, v.ArpAllTargets, v.Primary, v.PrimaryReselect, v.FailOverMac, v.XmitHashPolicy, v.ResendIgmp, v.NumPeerNotif, v.AllSlavesActive, v.MinLinks, v.LpInterval, v.PacketsPerSlave, v.LacpRate, v.AdSelect, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Geneve:
				fmt.Fprintf(w, "    geneve id %d remote %s ttl %d tos %d dport %d udpcsum %d udp_zero_csum_6TX %d udp_zero_csum_6RX %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.ID, v.Remote, v.Ttl, v.Tos, v.Dport, v.UdpCsum, v.UdpZeroCsum6Tx, v.UdpZeroCsum6Rx, v.Link, v.FlowBased, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Gretap:
				fmt.Fprintf(w, "    gretap i_key %d o_key %d encap_src_port %d encap_dst_port %d local %s remote %s iflags %d oflags %d pmtudisc %d ttl %d tos %d encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.IKey, v.OKey, v.EncapSport, v.EncapDport, v.Local, v.Remote, v.IFlags, v.OFlags, v.PMtuDisc, v.Ttl, v.Tos, v.EncapType, v.EncapFlags, v.Link, v.FlowBased, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Iptun:
				fmt.Fprintf(w, "    iptun local %s remote %s encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.EncapType, v.EncapFlags, v.Link, v.FlowBased, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Ip6tnl:
				fmt.Fprintf(w, "    ip6tnl local %s remote %s ttl %d tos %d proto %d flow_info %d encap_limit %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.Ttl, v.Tos, v.Proto, v.FlowInfo, v.EncapLimit, v.EncapType, v.EncapSport, v.EncapDport, v.EncapFlags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Sittun:
				fmt.Fprintf(w, "    sittun local %s remote %s ttl %d tos %d proto %d encap_limit %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.Ttl, v.Tos, v.Proto, v.EncapLimit, v.EncapType, v.EncapSport, v.EncapDport, v.EncapFlags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Vti:
				fmt.Fprintf(w, "    vti local %s remote %s ikey %d okey %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.IKey, v.OKey, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Gretun:
				fmt.Fprintf(w, "    gretun local %s remote %s ttl %d tos %d ptmudisc %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d ikey %d okey %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.Ttl, v.Tos, v.PMtuDisc, v.EncapType, v.EncapSport, v.EncapDport, v.EncapFlags, v.IKey, v.OKey, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Xfrmi:
				fmt.Fprintf(w, "    xfrmi if_id %d", v.Ifid)
			case *netlink.Can:
				fmt.Fprintf(w, "    can state %d bitrate %d sample-point %d tq %d prop-seg %d phase-seg1 %d phase-seg2 %d\n",
					v.State, v.BitRate, v.SamplePoint, v.TimeQuanta, v.PropagationSegment, v.PhaseSegment1, v.PhaseSegment2)
			case *netlink.IPoIB:
				fmt.Fprintf(w, "    ipoib pkey %d mode %d umcast %d\n", v.Pkey, v.Mode, v.Umcast)
			case *netlink.BareUDP:
				fmt.Fprintf(w, "    port %d ethertype %d srcport %d min multi_proto %t\n", v.Port, v.EtherType, v.SrcPortMin, v.MultiProto)

			}
		}

		if f.stats {
			stats := l.Statistics
			if stats != nil {
				fmt.Fprintf(w, "    RX: bytes %d packets %d errors %d dropped %d missed %d mcast %d\n",
					stats.RxBytes, stats.RxPackets, stats.RxErrors, stats.RxDropped, stats.RxMissedErrors, stats.Multicast)
				fmt.Fprintf(w, "    TX: bytes %d packets %d errors %d dropped %d carrier %d collsns %d\n",
					stats.TxBytes, stats.TxPackets, stats.TxErrors, stats.TxDropped, stats.TxCarrierErrors, stats.Collisions)
			}
		}

		if withAddresses {
			showLinkAddresses(w, v)
		}
	}
	return nil
}

func showLinkAddresses(w io.Writer, link netlink.Link) error {
	addrs, err := netlink.AddrList(link, family)
	if err != nil {
		return fmt.Errorf("can't enumerate addresses: %v", err)
	}

	for _, addr := range addrs {

		var inet string
		switch len(addr.IPNet.IP) {
		case 4:
			inet = "inet"
		case 16:
			inet = "inet6"
		default:
			return fmt.Errorf("can't figure out IP protocol version: IP length is %d", len(addr.IPNet.IP))
		}

		fmt.Fprintf(w, "    %s %s", inet, addr.IP)

		if addr.Broadcast != nil {
			fmt.Fprintf(w, " brd %s", addr.Broadcast)
		}

		fmt.Fprintf(w, " scope %s %s\n", addrScopes[netlink.Scope(addr.Scope)], addr.Label)

		var validLft, preferredLft string
		// TODO: fix vishnavanda/netlink. *Lft should be uint32, not int.
		if uint32(addr.PreferedLft) == math.MaxUint32 {
			preferredLft = "forever"
		} else {
			preferredLft = fmt.Sprintf("%dsec", addr.PreferedLft)
		}

		if uint32(addr.ValidLft) == math.MaxUint32 {
			validLft = "forever"
		} else {
			validLft = fmt.Sprintf("%dsec", addr.ValidLft)
		}

		fmt.Fprintf(w, "       valid_lft %s preferred_lft %s\n", validLft, preferredLft)
	}
	return nil
}
