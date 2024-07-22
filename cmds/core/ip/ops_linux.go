// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/vishvananda/netlink"
)

func (cmd cmd) showAllLinks(withAddresses bool, filterByType ...string) error {
	links, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("can't enumerate interfaces: %v", err)
	}
	return cmd.showLinks(withAddresses, links, filterByType...)
}

func (cmd cmd) showLink(link netlink.Link, withAddresses bool, filterByType ...string) error {
	return cmd.showLinks(withAddresses, []netlink.Link{link}, filterByType...)
}

type Link struct {
	IfIndex   int        `json:"ifindex,omitempty"`
	IfName    string     `json:"ifname"`
	Flags     []string   `json:"flags"`
	MTU       int        `json:"mtu,omitempty"`
	Operstate string     `json:"operstate"`
	Group     string     `json:"group,omitempty"`
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

func (cmd cmd) showLinks(withAddresses bool, links []netlink.Link, filterByType ...string) error {
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
				link.Group = fmt.Sprintf("%v", v.Attrs().Group)

				if !f.numeric && v.Attrs().Group == 0 {
					link.Group = "default"
				}

				link.Txqlen = v.Attrs().TxQLen
			}

			if withAddresses {

				addrs, err := cmd.handle.AddrList(v, family)
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

		return printJSON(cmd.out, linkObs)
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

				fmt.Fprintf(cmd.out, "%-20s %-10s", l.Name, l.OperState.String())

				for _, addr := range addrs {
					fmt.Fprintf(cmd.out, " %s", addr.IP)
				}

				fmt.Fprintf(cmd.out, "\n")

				continue
			}

			addr := " "
			if l.HardwareAddr != nil {
				addr = fmt.Sprintf(" %s ", l.HardwareAddr.String())
			}

			fmt.Fprintf(cmd.out, "%-25s %-10s%-20s <%s>\n", l.Name,
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

		group := fmt.Sprintf("%v", l.Group)

		if !f.numeric && l.Group == 0 {
			group = "default"
		}

		fmt.Fprintf(cmd.out, "%d: %s: <%s> mtu %d %sstate %s group %s\n", l.Index, l.Name,
			strings.Replace(strings.ToUpper(l.Flags.String()), "|", ",", -1),
			l.MTU, master, strings.ToUpper(l.OperState.String()), group)

		fmt.Fprintf(cmd.out, "    link/%s %s\n", l.EncapType, l.HardwareAddr)

		if f.details {
			switch v := v.(type) {
			case *netlink.Bridge:
				fmt.Fprintf(cmd.out, "    bridge hello_time %d ageing_time %d vlan_filtering %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.HelloTime, v.AgeingTime, v.VlanFiltering, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Vlan:
				fmt.Fprintf(cmd.out, "    vlan %s vlan-id %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.VlanProtocol, v.VlanId, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Macvlan:
				fmt.Fprintf(cmd.out, "    macvlan mode %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Macvtap:
				fmt.Fprintf(cmd.out, "    macvtap mode %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Tuntap:
				fmt.Fprintf(cmd.out, "    %s owner %d group %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.Owner, v.Group, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Veth:
				fmt.Fprintf(cmd.out, "    peer %s peer-address %s numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.PeerName, v.PeerHardwareAddr, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Vxlan:
				fmt.Fprintf(cmd.out, "    vxlan id %d src %s group %s ttl %d tos %d learning %t proxy %t rsc %t age %d limit %d port %d port-low %d port-high %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.VxlanId, v.SrcAddr, v.Group, v.TTL, v.TOS, v.Learning, v.Proxy, v.RSC, v.Age, v.Limit, v.Port, v.PortLow, v.PortHigh, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.IPVlan:
				fmt.Fprintf(cmd.out, "    ipvlan mode %d flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.Flags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.IPVtap:
				fmt.Fprintf(cmd.out, "    ipvtap mode %d flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.Flags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Bond:
				fmt.Fprintf(cmd.out, "    bond mode active slave %d %d miimon %d updelay %d downdelay %d use_carrier %d arp_interval %d arp_validate %s arp_all_targets %s primary %d primary_reselect %s fail_over_mac %s %s resend_igmp %d num_peer_notif %d all_slaves_active %d min_links %d lp_interval %d packets_per_slave %d lacp_rate %d ad_select %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Mode, v.ActiveSlave, v.Miimon, v.UpDelay, v.DownDelay, v.UseCarrier, v.ArpInterval, v.ArpValidate, v.ArpAllTargets, v.Primary, v.PrimaryReselect, v.FailOverMac, v.XmitHashPolicy, v.ResendIgmp, v.NumPeerNotif, v.AllSlavesActive, v.MinLinks, v.LpInterval, v.PacketsPerSlave, v.LacpRate, v.AdSelect, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Geneve:
				fmt.Fprintf(cmd.out, "    geneve id %d remote %s ttl %d tos %d dport %d udpcsum %d udp_zero_csum_6TX %d udp_zero_csum_6RX %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.ID, v.Remote, v.Ttl, v.Tos, v.Dport, v.UdpCsum, v.UdpZeroCsum6Tx, v.UdpZeroCsum6Rx, v.Link, v.FlowBased, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Gretap:
				fmt.Fprintf(cmd.out, "    gretap i_key %d o_key %d encap_src_port %d encap_dst_port %d local %s remote %s iflags %d oflags %d pmtudisc %d ttl %d tos %d encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.IKey, v.OKey, v.EncapSport, v.EncapDport, v.Local, v.Remote, v.IFlags, v.OFlags, v.PMtuDisc, v.Ttl, v.Tos, v.EncapType, v.EncapFlags, v.Link, v.FlowBased, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Iptun:
				fmt.Fprintf(cmd.out, "    iptun local %s remote %s encap_type %d encap_flags %d link %d flow_based %t numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.EncapType, v.EncapFlags, v.Link, v.FlowBased, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Ip6tnl:
				fmt.Fprintf(cmd.out, "    ip6tnl local %s remote %s ttl %d tos %d proto %d flow_info %d encap_limit %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.Ttl, v.Tos, v.Proto, v.FlowInfo, v.EncapLimit, v.EncapType, v.EncapSport, v.EncapDport, v.EncapFlags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Sittun:
				fmt.Fprintf(cmd.out, "    sittun local %s remote %s ttl %d tos %d proto %d encap_limit %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.Ttl, v.Tos, v.Proto, v.EncapLimit, v.EncapType, v.EncapSport, v.EncapDport, v.EncapFlags, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Vti:
				fmt.Fprintf(cmd.out, "    vti local %s remote %s ikey %d okey %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.IKey, v.OKey, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Gretun:
				fmt.Fprintf(cmd.out, "    gretun local %s remote %s ttl %d tos %d ptmudisc %d encap_type %d encap_src_port %d encap_dst_port %d encap_flags %d ikey %d okey %d numtxqueues %d numrxqueues %d gso_max_size %d gso_max_segs %d\n",
					v.Local, v.Remote, v.Ttl, v.Tos, v.PMtuDisc, v.EncapType, v.EncapSport, v.EncapDport, v.EncapFlags, v.IKey, v.OKey, v.NumTxQueues, v.NumRxQueues, v.GSOMaxSize, v.GSOMaxSegs)
			case *netlink.Xfrmi:
				fmt.Fprintf(cmd.out, "    xfrmi if_id %d", v.Ifid)
			case *netlink.Can:
				fmt.Fprintf(cmd.out, "    can state %d bitrate %d sample-point %d tq %d prop-seg %d phase-seg1 %d phase-seg2 %d\n",
					v.State, v.BitRate, v.SamplePoint, v.TimeQuanta, v.PropagationSegment, v.PhaseSegment1, v.PhaseSegment2)
			case *netlink.IPoIB:
				fmt.Fprintf(cmd.out, "    ipoib pkey %d mode %d umcast %d\n", v.Pkey, v.Mode, v.Umcast)
			case *netlink.BareUDP:
				fmt.Fprintf(cmd.out, "    port %d ethertype %d srcport %d min multi_proto %t\n", v.Port, v.EtherType, v.SrcPortMin, v.MultiProto)

			}
		}

		if f.stats {
			stats := l.Statistics
			if stats != nil {
				fmt.Fprintf(cmd.out, "    RX: bytes %d packets %d errors %d dropped %d missed %d mcast %d\n",
					stats.RxBytes, stats.RxPackets, stats.RxErrors, stats.RxDropped, stats.RxMissedErrors, stats.Multicast)
				fmt.Fprintf(cmd.out, "    TX: bytes %d packets %d errors %d dropped %d carrier %d collsns %d\n",
					stats.TxBytes, stats.TxPackets, stats.TxErrors, stats.TxDropped, stats.TxCarrierErrors, stats.Collisions)
			}
		}

		if withAddresses {
			cmd.showLinkAddresses(v)
		}
	}
	return nil
}

func (cmd cmd) showLinkAddresses(link netlink.Link) error {
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

		fmt.Fprintf(cmd.out, "    %s %s", inet, addr.IP)

		if addr.Broadcast != nil {
			fmt.Fprintf(cmd.out, " brd %s", addr.Broadcast)
		}

		fmt.Fprintf(cmd.out, " scope %s %s\n", addrScopes[netlink.Scope(addr.Scope)], addr.Label)

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

		fmt.Fprintf(cmd.out, "       valid_lft %s preferred_lft %s\n", validLft, preferredLft)
	}
	return nil
}
