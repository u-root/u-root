// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"math"
	"time"

	"github.com/google/netstack/sleep"
	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/tcpip/header"
	"github.com/google/netstack/tcpip/seqnum"
)

const (
	// minRTO is the minimium allowed value for the retransmit timeout.
	minRTO = 200 * time.Millisecond

	// initalCwnd is the initial congestion window.
	initialCwnd = 10
)

// sender holds the state necessary to send TCP segments.
type sender struct {
	ep *endpoint

	// lastSendTime is the timestamp when the last packet was sent.
	lastSendTime time.Time

	// dupAckCount is the number of duplicated acks received. It is used for
	// fast retransmit.
	dupAckCount int

	// fr holds state related to fast recovery.
	fr fastRecovery

	// sndCwnd is the congestion window, in packets.
	sndCwnd int

	// sndSsthresh is the threshold between slow start and congestion
	// avoidance.
	sndSsthresh int

	// sndCAAckCount is the number of packets acknowledged during congestion
	// avoidance. When enough packets have been ack'd (typically cwnd
	// packets), the congestion window is incremented by one.
	sndCAAckCount int

	// outstanding is the number of outstanding packets, that is, packets
	// that have been sent but not yet acknowledged.
	outstanding int

	// sndWnd is the send window size.
	sndWnd seqnum.Size

	// sndUna is the next unacknowledged sequence number.
	sndUna seqnum.Value

	// sndNxt is the sequence number of the next segment to be sent.
	sndNxt seqnum.Value

	// sndNxtList is the sequence number of the next segment to be added to
	// the send list.
	sndNxtList seqnum.Value

	// rttMeasureSeqNum is the sequence number being used for the latest RTT
	// measurement.
	rttMeasureSeqNum seqnum.Value

	// rttMeasureTime is the time when the rttMeasureSeqNum was sent.
	rttMeasureTime time.Time

	closed      bool
	writeNext   *segment
	writeList   segmentList
	resendTimer timer
	resendWaker sleep.Waker

	// srtt, rttvar & rto are the "smoothed round-trip time", "round-trip
	// time variation" and "retransmit timeout", as defined in section 2 of
	// RFC 6298.
	srtt       time.Duration
	rttvar     time.Duration
	rto        time.Duration
	srttInited bool

	// maxPayloadSize is the maximum size of the payload of a given segment.
	// It is initialized on demand.
	maxPayloadSize int

	// sndWndScale is the number of bits to shift left when reading the send
	// window size from a segment.
	sndWndScale uint8

	// maxSentAck is the maxium acknowledgement actually sent.
	maxSentAck seqnum.Value
}

// fastRecovery holds information related to fast recovery from a packet loss.
type fastRecovery struct {
	// active whether the endpoint is in fast recovery. The following fields
	// are only meaningful when active is true.
	active bool

	// first and last represent the inclusive sequence number range being
	// recovered.
	first seqnum.Value
	last  seqnum.Value

	// maxCwnd is the maximum value the congestion window may be inflated to
	// due to duplicate acks. This exists to avoid attacks where the
	// receiver intentionally sends duplicate acks to artificially inflate
	// the sender's cwnd.
	maxCwnd int
}

func newSender(ep *endpoint, iss, irs seqnum.Value, sndWnd seqnum.Size, mss uint16, sndWndScale int) *sender {
	s := &sender{
		ep:               ep,
		sndCwnd:          initialCwnd,
		sndSsthresh:      math.MaxInt64,
		sndWnd:           sndWnd,
		sndUna:           iss + 1,
		sndNxt:           iss + 1,
		sndNxtList:       iss + 1,
		rto:              1 * time.Second,
		rttMeasureSeqNum: iss + 1,
		lastSendTime:     time.Now(),
		maxPayloadSize:   int(mss),
		maxSentAck:       irs + 1,
	}

	// A negative sndWndScale means that no scaling is in use, otherwise we
	// store the scaling value.
	if sndWndScale > 0 {
		s.sndWndScale = uint8(sndWndScale)
	}

	m := int(ep.route.MTU()) - header.TCPMinimumSize
	// Adjust the maxPayloadsize to account for the timestamp option.
	if ep.sendTSOk {
		m -= header.TCPTimeStampOptionSize
	}
	if m < s.maxPayloadSize {
		s.maxPayloadSize = m
	}

	s.resendTimer.init(&s.resendWaker)

	return s
}

// sendAck sends an ACK segment.
func (s *sender) sendAck() {
	s.sendSegment(nil, flagAck, s.sndNxt)
}

// updateRTO updates the retransmit timeout when a new roud-trip time is
// available. This is done in accordance with section 2 of RFC 6298.
func (s *sender) updateRTO(rtt time.Duration) {
	if !s.srttInited {
		s.rttvar = rtt / 2
		s.srtt = rtt
		s.srttInited = true
	} else {
		diff := s.srtt - rtt
		if diff < 0 {
			diff = -diff
		}
		s.rttvar = (3*s.rttvar + diff) / 4
		s.srtt = (7*s.srtt + rtt) / 8
	}

	s.rto = s.srtt + 4*s.rttvar
	if s.rto < minRTO {
		s.rto = minRTO
	}
}

// resendSegment resends the first unacknowledged segment.
func (s *sender) resendSegment() {
	// Don't use any segments we already sent to measure RTT as they may
	// have been affected by packets being lost.
	s.rttMeasureSeqNum = s.sndNxt

	// Resend the segment.
	if seg := s.writeList.Front(); seg != nil {
		s.sendSegment(&seg.data, seg.flags, seg.sequenceNumber)
	}
}

// reduceSlowStartThreshold reduces the slow-start threshold per RFC 5681,
// page 6, eq. 4. It is called when we detect congestion in the network.
func (s *sender) reduceSlowStartThreshold() {
	s.sndSsthresh = s.outstanding / 2
	if s.sndSsthresh < 2 {
		s.sndSsthresh = 2
	}
}

// retransmitTimerExpired is called when the retransmit timer expires, and
// unacknowledged segments are assumed lost, and thus need to be resent.
// Returns true if the connection is still usable, or false if the connection
// is deemed lost.
func (s *sender) retransmitTimerExpired() bool {
	// Check if the timer actually expired or if it's a spurious wake due
	// to a previously orphaned runtime timer.
	if !s.resendTimer.checkExpiration() {
		return true
	}

	// Give up if we've waited more than a minute since the last resend.
	if s.rto >= 60*time.Second {
		return false
	}

	// Set new timeout. The timer will be restarted by the call to sendData
	// below.
	s.rto *= 2

	if s.fr.active {
		// We were attempting fast recovery but were not successfull.
		// Leave the state. We don't need to update ssthresh because it
		// has already been updated when entered fast-recovery.
		s.leaveFastRecovery()
	} else {
		// We lost a packet, so reduce ssthresh.
		s.reduceSlowStartThreshold()
	}

	// Reduce the congestion window to 1, i.e., enter slow-start.
	s.sndCwnd = 1

	// Mark the next segment to be sent as the first unacknowledged one and
	// start sending again. Set the number of outstanding packets to 0 so
	// that we'll be able to retransmit.
	//
	// We'll keep on transmitting (or retransmitting) as we get acks for
	// the data we transmit.
	s.outstanding = 0
	s.writeNext = s.writeList.Front()
	s.sendData()

	return true
}

// sendData sends new data segments. It is called when data becomes available or
// when the send window opens up.
func (s *sender) sendData() {
	limit := s.maxPayloadSize

	// Reduce the congestion window to min(IW, cwnd) per RFC 5681, page 10.
	// "A TCP SHOULD set cwnd to no more than RW before beginning
	// transmission if the TCP has not sent data in the interval exceeding
	// the retrasmission timeout."
	if !s.fr.active && time.Now().Sub(s.lastSendTime) > s.rto {
		if s.sndCwnd > initialCwnd {
			s.sndCwnd = initialCwnd
		}
	}

	// TODO: We currently don't merge multiple send buffers
	// into one segment if they happen to fit. We should do that
	// eventually.
	var seg *segment
	end := s.sndUna.Add(s.sndWnd)
	for seg = s.writeNext; seg != nil && s.outstanding < s.sndCwnd; seg = seg.Next() {
		// We abuse the flags field to determine if we have already
		// assigned a sequence number to this segment.
		if seg.flags == 0 {
			seg.sequenceNumber = s.sndNxt
			seg.flags = flagAck | flagPsh
		}

		var segEnd seqnum.Value
		if seg.data.Size() == 0 {
			// We're sending a FIN.
			seg.flags = flagAck | flagFin
			segEnd = seg.sequenceNumber.Add(1)
		} else {
			// We're sending a non-FIN segment.
			if !seg.sequenceNumber.LessThan(end) {
				break
			}

			available := int(seg.sequenceNumber.Size(end))
			if available > limit {
				available = limit
			}

			if seg.data.Size() > available {
				// Split this segment up.
				nSeg := seg.clone()
				nSeg.data.TrimFront(available)
				nSeg.sequenceNumber.UpdateForward(seqnum.Size(available))
				s.writeList.InsertAfter(seg, nSeg)
				seg.data.CapLength(available)
			}

			s.outstanding++
			segEnd = seg.sequenceNumber.Add(seqnum.Size(seg.data.Size()))
		}

		s.sendSegment(&seg.data, seg.flags, seg.sequenceNumber)

		// Update sndNxt if we actually sent new data (as opposed to
		// retransmitting some previously sent data).
		if s.sndNxt.LessThan(segEnd) {
			s.sndNxt = segEnd
		}
	}

	// Remember the next segment we'll write.
	s.writeNext = seg

	// Enable the timer if we have pending data and it's not enabled yet.
	if !s.resendTimer.enabled() && s.sndUna != s.sndNxt {
		s.resendTimer.enable(s.rto)
	}
}

func (s *sender) enterFastRecovery() {
	// Save state to reflect we're now in fast recovery.
	s.reduceSlowStartThreshold()
	s.sndCwnd = s.sndSsthresh
	s.fr.first = s.sndUna
	s.fr.last = s.sndNxt - 1
	s.fr.maxCwnd = s.sndCwnd + s.outstanding
	s.fr.active = true
}

func (s *sender) leaveFastRecovery() {
	s.fr.active = false

	// Deflate cwnd. It had been artifically inflated when new dups arrived.
	s.sndCwnd = s.sndSsthresh
}

// checkDuplicateAck is called when an ack is received. It manages the state
// related to duplicate acks and determines if a retransmit is needed according
// to the rules in RFC 6582 (NewReno).
func (s *sender) checkDuplicateAck(seg *segment) bool {
	ack := seg.ackNumber
	if s.fr.active {
		// We are in fast recovery mode. Ignore the ack if it's out of
		// range.
		if !ack.InRange(s.sndUna, s.sndNxt+1) {
			return false
		}

		// Leave fast recovery if it acknowleges all the data covered by
		// this fast recovery session.
		if s.fr.last.LessThan(ack) {
			s.leaveFastRecovery()
			return false
		}

		// Don't count this as a duplicate if it is carrying data or
		// updating the window.
		if seg.logicalLen() != 0 || s.sndWnd != seg.window {
			return false
		}

		// Inflate the congestion window if we're getting duplicate acks
		// for the packet we retransmitted.
		if ack == s.fr.first {
			// We received a dup, inflate the congestion window by 1
			// packet if we're not at the max yet.
			if s.sndCwnd < s.fr.maxCwnd {
				s.sndCwnd++
			}
			return false
		}

		// A partial ack was received. Retransmit this packet and
		// remember it so that we don't retransmit it again. We don't
		// inflate the window because we're putting the same packet back
		// onto the wire.
		//
		// N.B. The retransmit timer will be reset by the caller.
		s.fr.first = ack
		return true
	}

	// We're not in fast recovery yet. A segment is considered a duplicate
	// only if it doesn't carry any data and doesn't update the send window,
	// because if it does, it wasn't sent in response to an out-of-order
	// segment.
	if ack != s.sndUna || seg.logicalLen() != 0 || s.sndWnd != seg.window || ack == s.sndNxt {
		s.dupAckCount = 0
		return false
	}

	// Enter fast recovery when we reach 3 dups.
	s.dupAckCount++
	if s.dupAckCount != 3 {
		return false
	}

	s.enterFastRecovery()
	s.dupAckCount = 0
	return true
}

// updateCwnd updates the congestion window based on the number of packets that
// were acknowledged.
func (s *sender) updateCwnd(packetsAcked int) {
	if s.sndCwnd < s.sndSsthresh {
		// Don't let the congestion window cross into the congestion
		// avoidance range.
		newcwnd := s.sndCwnd + packetsAcked
		if newcwnd >= s.sndSsthresh {
			newcwnd = s.sndSsthresh
			s.sndCAAckCount = 0
		}

		packetsAcked -= newcwnd - s.sndCwnd
		s.sndCwnd = newcwnd
		if packetsAcked == 0 {
			// We've consumed all ack'd packets.
			return
		}
	}

	// Consume the packets in congestion avoidance mode.
	s.sndCAAckCount += packetsAcked
	if s.sndCAAckCount >= s.sndCwnd {
		s.sndCwnd += s.sndCAAckCount / s.sndCwnd
		s.sndCAAckCount = s.sndCAAckCount % s.sndCwnd
	}
}

// handleRcvdSegment is called when a segment is received; it is responsible for
// updating the send-related state.
func (s *sender) handleRcvdSegment(seg *segment) {
	// Check if we can extract an RTT measurement from this ack.
	if s.rttMeasureSeqNum.LessThan(seg.ackNumber) {
		s.updateRTO(time.Now().Sub(s.rttMeasureTime))
		s.rttMeasureSeqNum = s.sndNxt
	}

	// Update Timestamp if required. See RFC7323, section-4.3.
	s.ep.updateRecentTimestamp(seg.parsedOptions.TSVal, s.maxSentAck, seg.sequenceNumber)

	// Count the duplicates and do the fast retransmit if needed.
	rtx := s.checkDuplicateAck(seg)

	// Stash away the current window size.
	s.sndWnd = seg.window

	// Ignore ack if it doesn't acknowledge any new data.
	ack := seg.ackNumber
	if (ack - 1).InRange(s.sndUna, s.sndNxt) {
		// When an ack is received we must reset the timer. We stop it
		// here and it will be restarted later if needed.
		s.resendTimer.disable()

		// Remove all acknowledged data from the write list.
		acked := s.sndUna.Size(ack)
		s.sndUna = ack

		ackLeft := acked
		originalOutsanding := s.outstanding
		for ackLeft > 0 {
			// We use logicalLen here because we can have FIN
			// segments (which are always at the end of list) that
			// have no data, but do consume a sequence number.
			seg := s.writeList.Front()
			datalen := seg.logicalLen()

			if datalen > ackLeft {
				seg.data.TrimFront(int(ackLeft))
				break
			}

			if s.writeNext == seg {
				s.writeNext = seg.Next()
			}
			s.writeList.Remove(seg)
			s.outstanding--
			seg.decRef()
			ackLeft -= datalen
		}

		// Update the send buffer usage and notify potential waiters.
		s.ep.updateSndBufferUsage(int(acked))

		// Update the congestion window based on the number of
		// acknowledged packets.
		s.updateCwnd(originalOutsanding - s.outstanding)

		// It is possible for s.outstanding to drop below zero if we get
		// a retransmit timeout, reset outstanding to zero but later
		// get an ack that cover previously sent data.
		if s.outstanding < 0 {
			s.outstanding = 0
		}
	}

	// Now that we've popped all acknowledged data from the retransmit
	// queue, retransmit if needed.
	if rtx {
		s.resendSegment()
	}

	// Send more data now that some of the pending data has been ack'd, or
	// that the window opened up, or the congestion window was inflated due
	// to a duplicate ack during fast recovery. This will also re-enable
	// the retransmit timer if needed.
	s.sendData()
}

// sendSegment sends a new segment containing the given payload, flags and
// sequence number.
func (s *sender) sendSegment(data *buffer.VectorisedView, flags byte, seq seqnum.Value) *tcpip.Error {
	s.lastSendTime = time.Now()
	if seq == s.rttMeasureSeqNum {
		s.rttMeasureTime = s.lastSendTime
	}

	rcvNxt, rcvWnd := s.ep.rcv.getSendParams()

	// Remember the max sent ack.
	s.maxSentAck = rcvNxt

	if data == nil {
		return s.ep.sendRaw(nil, flags, seq, rcvNxt, rcvWnd)
	}

	if len(data.Views()) > 1 {
		panic("send path does not support views with multiple buffers")
	}

	return s.ep.sendRaw(data.First(), flags, seq, rcvNxt, rcvWnd)
}
