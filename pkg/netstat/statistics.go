// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

var (
	SNMP4file   = "/proc/net/snmp"
	netstatfile = "/proc/net/netstat"
)

type NetStat struct {
	TCPExt   tcpExt
	IPExt    ipExt
	MPTCPExt mptcpext
}

func (n *NetStat) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "%s\n", n.TCPExt.String())
	fmt.Fprintf(&s, "%s\n", n.IPExt.String())
	fmt.Fprintf(&s, "%s\n", n.MPTCPExt.String())

	return s.String()
}

type tcpExt struct {
	SyncookiesSent            string `opt:"" text:"%s SYN cookies sent"`
	SyncookiesRecv            string `opt:"" text:"%s SYN cookies received"`
	SyncookiesFailed          string `opt:"" text:"%s invalid SYN cookies received"`
	EmbryonicRsts             string `opt:"" text:"%s resets received for embryonic SYN_RECV sockets"`
	PruneCalled               string `opt:"" text:"%s packets pruned from receive queue because of socket buffer overrun"`
	RcvPruned                 string `opt:"" text:"%s packets pruned from receive queue"`
	OfoPruned                 string `opt:"" text:"%s packets dropped from out-of-order queue because of socket buffer overrun"`
	OutOfWindowIcmps          string `opt:"" text:"%s ICMP packets dropped because they were out-of-window"`
	LockDroppedIcmps          string `opt:"" text:"%s ICMP packets dropped because socket was locked"`
	ArpFilter                 string `opt:"" text:""`
	TW                        string `opt:"" text:"%s TCP sockets finished time wait in fast timer"`
	TWRecycled                string `opt:"" text:"%s time wait sockets recycled by time stamp"`
	TWKilled                  string `opt:"" text:"%s TCP sockets finished time wait in slow timer"`
	PAWSActive                string `opt:"" text:"%s active connections rejected because of time stamp"`
	PAWSEstab                 string `opt:"" text:"%s packets rejected in established connections because of timestamp"`
	DelayedACKs               string `opt:"" text:"%s delayed acks sent"`
	DelayedACKLocked          string `opt:"" text:"%s delayed acks further delayed because of locked socket"`
	DelayedACKLost            string `opt:"" text:"Quick ack mode was activated %s times"`
	ListenOverflows           string `opt:"" text:"%s times the listen queue of a socket overflowed"`
	ListenDrops               string `opt:"" text:"%s SYNs to LISTEN sockets dropped"`
	TCPHPHits                 string `opt:"" text:"%s packet headers predicted"`
	TCPPureAcks               string `opt:"" text:"%s acknowledgments not containing data payload received"`
	TCPHPAcks                 string `opt:"" text:"%s predicted acknowledgments"`
	TCPRenoRecovery           string `opt:"" text:"%s times recovered from packet loss due to fast retransmit"`
	TCPSackRecovery           string `opt:"" text:"%s times recovered from packet loss by selective acknowledgements"`
	TCPSACKReneging           string `opt:"" text:"%s bad SACK blocks received"`
	TCPSACKReorder            string `opt:"" text:"Detected reordering %s times using SACK"`
	TCPRenoReorder            string `opt:"" text:"Detected reordering %s times using reno fast retransmit"`
	TCPTSReorder              string `opt:"" text:"Detected reordering %s times using time stamp"`
	TCPPartialUndo            string `opt:"" text:"%s congestion windows partially recovered using Hoe heuristic"`
	TCPDSACKUndo              string `opt:"" text:"%s congestion windows recovered without slow start by DSACK"`
	TCPLossUndo               string `opt:"" text:"%s congestion windows recovered without slow start after partial ack"`
	TCPLostRetransmit         string `opt:"" text:"%s retransmits lost"`
	TCPRenoFailures           string `opt:"" text:"%s timeouts after reno fast retransmit"`
	TCPSackFailures           string `opt:"" text:"%s timeouts after SACK recovery"`
	TCPLossFailures           string `opt:"" text:"%s timeouts in loss state"`
	TCPFastRetrans            string `opt:"" text:"%s fast retransmits"`
	TCPSlowStartRetrans       string `opt:"" text:"%s retransmits in slow start"`
	TCPTimeouts               string `opt:"" text:"%s other TCP timeouts"`
	TCPLossProbes             string `opt:"" text:""`
	TCPLossProbeRecovery      string `opt:"" text:""`
	TCPRenoRecoveryFail       string `opt:"" text:"%s classic Reno fast retransmits failed"`
	TCPSackRecoveryFail       string `opt:"" text:"%s SACK retransmits failed"`
	TCPRcvCollapsed           string `opt:"" text:"%s packets collapsed in receive queue due to low socket buffer"`
	TCPBacklogCoalesce        string `opt:"" text:""`
	TCPDSACKOldSent           string `opt:"" text:"%s DSACKs sent for old packets"`
	TCPDSACKOfoSent           string `opt:"" text:"%s DSACKs sent for out of order packets"`
	TCPDSACKRecv              string `opt:"" text:"%s DSACKs received"`
	TCPDSACKOfoRecv           string `opt:"" text:"%s DSACKs for out of order packets received"`
	TCPAbortOnData            string `opt:"" text:"%s connections reset due to unexpected data"`
	TCPAbortOnClose           string `opt:"" text:"%s connections reset due to early user close"`
	TCPAbortOnMemory          string `opt:"" text:"%s connections aborted due to memory pressure"`
	TCPAbortOnTimeout         string `opt:"" text:"%s connections aborted due to timeout"`
	TCPAbortOnLinger          string `opt:"" text:"%s connections aborted after user close in linger timeout"`
	TCPAbortFailed            string `opt:"" text:"%s times unable to send RST due to no memory"`
	TCPMemoryPressures        string `opt:"" text:"TCP ran low on memory %s times"`
	TCPMemoryPressuresChrono  string `opt:"" text:""`
	TCPSACKDiscard            string `opt:"" text:""`
	TCPDSACKIgnoredOld        string `opt:"" text:""`
	TCPDSACKIgnoredNoUndo     string `opt:"" text:""`
	TCPSpuriousRTOs           string `opt:"" text:""`
	TCPMD5NotFound            string `opt:"" text:""`
	TCPMD5Unexpected          string `opt:"" text:""`
	TCPMD5Failure             string `opt:"" text:""`
	TCPSackShifted            string `opt:"" text:""`
	TCPSackMerged             string `opt:"" text:""`
	TCPSackShiftFallback      string `opt:"" text:""`
	TCPBacklogDrop            string `opt:"" text:""`
	PFMemallocDrop            string `opt:"" text:""`
	TCPMinTTLDrop             string `opt:"" text:""`
	TCPDeferAcceptDrop        string `opt:"" text:""`
	IPReversePathFilter       string `opt:"" text:""`
	TCPTimeWaitOverflow       string `opt:"" text:""`
	TCPReqQFullDoCookies      string `opt:"" text:""`
	TCPReqQFullDrop           string `opt:"" text:""`
	TCPRetransFail            string `opt:"" text:""`
	TCPRcvCoalesce            string `opt:"" text:""`
	TCPOFOQueue               string `opt:"" text:""`
	TCPOFODrop                string `opt:"" text:""`
	TCPOFOMerge               string `opt:"" text:""`
	TCPChallengeACK           string `opt:"" text:""`
	TCPSYNChallenge           string `opt:"" text:""`
	TCPFastOpenActive         string `opt:"" text:""`
	TCPFastOpenActiveFail     string `opt:"" text:""`
	TCPFastOpenPassive        string `opt:"" text:""`
	TCPFastOpenPassiveFail    string `opt:"" text:""`
	TCPFastOpenListenOverflow string `opt:"" text:""`
	TCPFastOpenCookieReqd     string `opt:"" text:""`
	TCPFastOpenBlackhole      string `opt:"" text:""`
	TCPSpuriousRtxHostQueues  string `opt:"" text:""`
	BusyPollRxPackets         string `opt:"" text:""`
	TCPAutoCorking            string `opt:"" text:""`
	TCPFromZeroWindowAdv      string `opt:"" text:""`
	TCPToZeroWindowAdv        string `opt:"" text:""`
	TCPWantZeroWindowAdv      string `opt:"" text:""`
	TCPSynRetrans             string `opt:"" text:""`
	TCPOrigDataSent           string `opt:"" text:""`
	TCPHystartTrainDetect     string `opt:"" text:""`
	TCPHystartTrainCwnd       string `opt:"" text:""`
	TCPHystartDelayDetect     string `opt:"" text:""`
	TCPHystartDelayCwnd       string `opt:"" text:""`
	TCPACKSkippedSynRecv      string `opt:"" text:""`
	TCPACKSkippedPAWS         string `opt:"" text:""`
	TCPACKSkippedSeq          string `opt:"" text:""`
	TCPACKSkippedFinWait2     string `opt:"" text:""`
	TCPACKSkippedTimeWait     string `opt:"" text:""`
	TCPACKSkippedChallenge    string `opt:"" text:""`
	TCPWinProbe               string `opt:"" text:""`
	TCPKeepAlive              string `opt:"" text:""`
	TCPMTUPFail               string `opt:"" text:""`
	TCPMTUPSuccess            string `opt:"" text:""`
	TCPDelivered              string `opt:"" text:""`
	TCPDeliveredCE            string `opt:"" text:""`
	TCPAckCompressed          string `opt:"" text:""`
	TCPZeroWindowDrop         string `opt:"" text:""`
	TCPRcvQDrop               string `opt:"" text:""`
	TCPWqueueTooBig           string `opt:"" text:""`
	TCPFastOpenPassiveAltKey  string `opt:"" text:""`
	TcpTimeoutRehash          string `opt:"" text:""` //nolint:revive
	TcpDuplicateDataRehash    string `opt:"" text:""` //nolint:revive
	TCPDSACKRecvSegs          string `opt:"" text:""`
	TCPDSACKIgnoredDubious    string `opt:"" text:""`
	TCPMigrateReqSuccess      string `opt:"" text:""`
	TCPMigrateReqFailure      string `opt:"" text:""`
	TCPPLBRehash              string `opt:"" text:""`
	TCPAORequired             string `opt:"" text:""`
	TCPAOBad                  string `opt:"" text:""`
	TCPAOKeyNotFound          string `opt:"" text:""`
	TCPAOGood                 string `opt:"" text:""`
	TCPAODroppedIcmps         string `opt:"" text:""`
}

func (t *tcpExt) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*t)
	refVal := reflect.ValueOf(*t)
	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

type ipExt struct {
	InNoRoutes      string `opt:"" text:""`
	InTruncatedPkts string `opt:"" text:""`
	InMcastPkts     string `opt:"" text:""`
	OutMcastPkts    string `opt:"" text:""`
	InBcastPkts     string `opt:"" text:""`
	OutBcastPkts    string `opt:"" text:""`
	InOctets        string `opt:"" text:""`
	OutOctets       string `opt:"" text:""`
	InMcastOctets   string `opt:"" text:""`
	OutMcastOctets  string `opt:"" text:""`
	InBcastOctets   string `opt:"" text:""`
	OutBcastOctets  string `opt:"" text:""`
	InCsumErrors    string `opt:"" text:""`
	InNoECTPkts     string `opt:"" text:""`
	InECT1Pkts      string `opt:"" text:""`
	InECT0Pkts      string `opt:"" text:""`
	InCEPkts        string `opt:"" text:""`
	ReasmOverlaps   string `opt:"" text:""`
}

func (i *ipExt) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*i)
	refVal := reflect.ValueOf(*i)
	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

type mptcpext struct {
	MPCapableSYNRX          string `opt:"" text:""`
	MPCapableSYNTX          string `opt:"" text:""`
	MPCapableSYNACKRX       string `opt:"" text:""`
	MPCapableACKRX          string `opt:"" text:""`
	MPCapableFallbackACK    string `opt:"" text:""`
	MPCapableFallbackSYNACK string `opt:"" text:""`
	MPFallbackTokenInit     string `opt:"" text:""`
	MPTCPRetrans            string `opt:"" text:""`
	MPJoinNoTokenFound      string `opt:"" text:""`
	MPJoinSynRx             string `opt:"" text:""`
	MPJoinSynAckRx          string `opt:"" text:""`
	MPJoinSynAckHMacFailure string `opt:"" text:""`
	MPJoinAckRx             string `opt:"" text:""`
	MPJoinAckHMacFailure    string `opt:"" text:""`
	DSSNotMatching          string `opt:"" text:""`
	InfiniteMapTx           string `opt:"" text:""`
	InfiniteMapRx           string `opt:"" text:""`
	DSSNoMatchTCP           string `opt:"" text:""`
	DataCsumErr             string `opt:"" text:""`
	OFOQueueTail            string `opt:"" text:""`
	OFOQueue                string `opt:"" text:""`
	OFOMerge                string `opt:"" text:""`
	NoDSSInWindow           string `opt:"" text:""`
	DuplicateData           string `opt:"" text:""`
	AddAddr                 string `opt:"" text:""`
	AddAddrTx               string `opt:"" text:""`
	AddAddrTxDrop           string `opt:"" text:""`
	EchoAdd                 string `opt:"" text:""`
	EchoAddTx               string `opt:"" text:""`
	EchoAddTxDrop           string `opt:"" text:""`
	PortAdd                 string `opt:"" text:""`
	AddAddrDrop             string `opt:"" text:""`
	MPJoinPortSynRx         string `opt:"" text:""`
	MPJoinPortSynAckRx      string `opt:"" text:""`
	MPJoinPortAckRx         string `opt:"" text:""`
	MismatchPortSynRx       string `opt:"" text:""`
	MismatchPortAckRx       string `opt:"" text:""`
	RmAddr                  string `opt:"" text:""`
	RmAddrDrop              string `opt:"" text:""`
	RmAddrTx                string `opt:"" text:""`
	RmAddrTxDrop            string `opt:"" text:""`
	RmSubflow               string `opt:"" text:""`
	MPPrioTx                string `opt:"" text:""`
	MPPrioRx                string `opt:"" text:""`
	MPFailTx                string `opt:"" text:""`
	MPFailRx                string `opt:"" text:""`
	MPFastcloseTx           string `opt:"" text:""`
	MPFastcloseRx           string `opt:"" text:""`
	MPRstTx                 string `opt:"" text:""`
	MPRstRx                 string `opt:"" text:""`
	RcvPruned               string `opt:"" text:""`
	SubflowStale            string `opt:"" text:""`
	SubflowRecover          string `opt:"" text:""`
	SndWndShared            string `opt:"" text:""`
	RcvWndShared            string `opt:"" text:""`
	RcvWndConflictUpdate    string `opt:"" text:""`
	RcvWndConflict          string `opt:"" text:""`
	MPCurrEstab             string `opt:"" text:""`
}

func (m *mptcpext) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*m)
	refVal := reflect.ValueOf(*m)
	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

type SNMP struct {
	IP      ip
	ICMP    icmp
	ICMPMsg icmpmsg
	TCP     tcp
	UDP     udp
	UDPL    udp
}

func (s *SNMP) String() string {
	var str strings.Builder

	fmt.Fprintf(&str, "%s\n", s.IP.String())
	fmt.Fprintf(&str, "%s\n", s.ICMP.String())
	fmt.Fprintf(&str, "%s\n", s.ICMPMsg.String())
	fmt.Fprintf(&str, "%s\n", s.TCP.String())
	fmt.Fprintf(&str, "%s\n", s.UDP.String())

	return str.String()
}

type ip struct {
	Forwarding      string `req:"" text:"Forwarding is %s"`
	DefaultTTL      string `req:"" text:"Default TTL is %s"`
	InReceives      string `req:"" text:"%s total packets received"`
	InHdrErrors     string `opts:"" text:"%s with invalid headers"`
	InAddrErrors    string `opts:"" text:"%s with invalid addresses"`
	ForwDatagrams   string `req:"" text:"%s forwarded"`
	InUnknownProtos string `opts:"" text:"%s with unknown protocol"`
	InDiscards      string `req:"" text:"%s incoming packets discarded"`
	InDelivers      string `req:"" text:"%s incoming packets delivered"`
	OutRequests     string `req:"" text:"%s requests sent out"`
	OutDiscards     string `opts:"" text:"%s outgoing packets dropped"`
	OutNoRoutes     string `opts:"" text:"%s dropped because of missing route"`
	ReasmTimeout    string `opts:"" text:"%s fragments dropped after timeout"`
	ReasmReqds      string `opts:"" text:"%s reassemblies required"`
	ReasmOKs        string `opts:"" text:"%s packets reassembled ok"`
	ReasmFails      string `opts:"" text:"%s packet reassemblies failed"`
	FragOKs         string `opts:"" text:"%s outgoing packets fragmented ok"`
	FragFails       string `opts:"" text:"%s outgoing packets failed fragmentation"`
	FragCreates     string `opts:"" text:"%s fragments created"`
	OutTransmits    string `opts:"" text:"OutTransmits: %s"`
}

func (i *ip) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*i)
	refVal := reflect.ValueOf(*i)
	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

type icmp struct {
	InMsgs   string `req:"" text:"%s ICMP messages received"`
	InErrors string `req:"" text:"%s input ICMP message failed"`
	// Input histogram
	InCsumErrors    string `text:""`
	InDestUnreachs  string `text:"destination unreachable: %s" input:""`
	InTimeExcds     string `text:"timeout in transit: %s" input:""`
	InParmProbs     string `text:"wrong parameters: %s" input:""`
	InSrcQuenchs    string `text:"source quenches: %s" input:""`
	InRedirects     string `text:"redirects: %s" input:""`
	InEchos         string `text:"echo requests: %s" input:""`
	InEchoReps      string `text:"echo replies: %s" input:""`
	InTimestamps    string `text:"timestamp request: %s" input:""`
	InTimestampReps string `text:"timestamp reply: %s" input:""`
	InAddrMasks     string `text:"address mask request: %s" input:""`
	InAddrMaskReps  string `text:"address mask replies: %s" input:""`
	// ICMPv6
	InPktTooBigs             string `text:"packets too big: %s" input:""`
	InGroupMembQueries       string `text:"group member queries: %s" input:""`
	InGroupMembResponses     string `text:"group member responses: %s" input:""`
	InGroupMembReductions    string `text:"group member reductions: %s" input:""`
	InRouterSolicits         string `text:"router solicits: %s" input:""`
	InRouterAdvertisements   string `text:"router advertisement: %s" input:""`
	InNeighborSolicits       string `text:"neighbour solicits: %s" input:""`
	InNeighborAdvertisements string `text:"neighbour advertisement: %s" input:""`
	InMLDv2Reports           string `text:""`

	OutMsgs   string `req:"" text:"%s ICMP messages sent"`
	OutErrors string `req:"" text:"%s ICMP messages failed"`
	// Output histogram
	OutRateLimitGlobal string `text:""`
	OutRateLimitHost   string `text:""`
	OutDestUnreachs    string `text:"destination unreachable: %s" output:""`
	OutTimeExcds       string `text:"time exceeded: %s" output:""`
	OutParmProbs       string `text:"wrong parameters: %s" output:""`
	OutSrcQuenchs      string `text:"source quench: %s" output:""`
	OutRedirects       string `text:"redirect: %s" output:""`
	OutEchos           string `text:"echo requests: %s" output:""`
	OutEchoReps        string `text:"echo replies: %s" output:""`
	OutTimestamps      string `text:"timestamp requests: %s" output:""`
	OutTimestampReps   string `text:"timestamp replies: %s" output:""`
	OutAddrMasks       string `text:"address mask requests: %s" output:""`
	OutAddrMaskReps    string `text:"address mask replies: %s" output:""`
	// ICMPv6
	OutPktTooBigs             string `text:"packets too big: %llu" output:""`
	OutEchoReplies            string `text:"echo replies: %s" output:""`
	OutGroupMembQueries       string `text:"group member queries: %s" output:""`
	OutGroupMembResponses     string `text:"group member responses: %s" output:""`
	OutGroupMembReductions    string `text:"group member reductions: %s" output:""`
	OutRouterSolicits         string `text:"router solicits: %s" output:""`
	OutRouterAdvertisements   string `text:"router advertisement: %s" output:""`
	OutNeighborSolicits       string `text:"neighbor solicits: %s" output:""`
	OutNeighborAdvertisements string `text:"neighbor advertisements: %s" output:""`
	OutMLDv2Reports           string `text:""`
	// ICMPv6 In/Out Types
	InType136  string
	InType143  string
	InType133  string
	OutType135 string
	OutType143 string
}

func (i *icmp) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*i)
	refVal := reflect.ValueOf(*i)

	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

type icmpmsg struct {
	InType3  string
	OutType3 string
}

func (i *icmpmsg) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "IcmpMsg:\n")
	fmt.Fprintf(&s, "\tInType3: %s\n", i.InType3)
	fmt.Fprintf(&s, "\tOutType3: %s\n", i.OutType3)

	return s.String()
}

type tcp struct {
	RtoAlgorithm string
	RtoMin       string
	RtoMax       string
	MaxConn      string
	ActiveOpens  string `req:"" text:"%s active connection openings"`
	PassiveOpens string `req:"" text:"%s passive connection openings"`
	AttemptFails string `req:"" text:"%s failed connection attempts"`
	EstabResets  string `req:"" text:"%s connection resets received"`
	CurrEstab    string `req:"" text:"%s connections established"`
	InSegs       string `req:"" text:"%s segments received"`
	OutSegs      string `req:"" text:"%s resets sent"`
	RetransSegs  string `req:"" text:"%s segments retransmitted"`
	InErrs       string `req:"" text:"%s bad segments received"`
	OutRsts      string `req:"" text:"%s segments sent out"`
	InCsumErrors string
}

func (t *tcp) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*t)
	refVal := reflect.ValueOf(*t)
	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

type udp struct {
	InDatagrams  string `req:"" text:"%s packets received"`
	NoPorts      string `req:"" text:"%s packets to unknown port received"`
	InErrors     string `req:"" text:"%s packet receive errors"`
	OutDatagrams string `req:"" text:"%s packets sent"`
	RcvbufErrors string `req:"" text:"%s receive buffer errors"`
	SndbufErrors string `req:"" text:"%s send buffer errors"`
	InCsumErrors string
	IngoredMulti string
	MemErrors    string
}

func (u *udp) String() string {
	var s strings.Builder

	typeVal := reflect.TypeOf(*u)
	refVal := reflect.ValueOf(*u)
	parseFromTags(typeVal, refVal, &s)

	return s.String()
}

func parseFromTags(t reflect.Type, v reflect.Value, w io.Writer) {
	var in, out strings.Builder
	numFields := t.NumField()

	fmt.Fprintf(w, "%s:\n", t.Name())
	for j := 0; j < numFields; j++ {
		field := t.Field(j)
		outputfmt := field.Tag.Get("text")
		refFieldVal := v.Field(j).String()
		_, req := field.Tag.Lookup("req")

		var inb, outb, hasVal bool
		_, inb = field.Tag.Lookup("input")
		_, outb = field.Tag.Lookup("output")

		hasVal = (refFieldVal != "0" && refFieldVal != "")

		if (hasVal || (req && refFieldVal != "")) && outputfmt != "" && !inb && !outb {
			fmt.Fprintf(w, "\t"+outputfmt+"\n", refFieldVal)
		} else if hasVal && outputfmt == "" && !inb && !outb {
			fmt.Fprintf(w, "\t%s: %s\n", field.Name, refFieldVal)
		} else if inb && hasVal {
			// Add to input histogram
			fmt.Fprintf(&in, "\t\t"+outputfmt+"\n", refFieldVal)
		} else if outb && hasVal {
			// Add to output histogram
			fmt.Fprintf(&out, "\t\t"+outputfmt+"\n", refFieldVal)
		}
	}

	if in.Len() > 0 {
		fmt.Fprintf(w, "\t%s\n%s\n", "Input historam:", in.String())
	}

	if out.Len() > 0 {
		fmt.Fprintf(w, "\t%s\n%s\n", "Output historam:", out.String())
	}
}
