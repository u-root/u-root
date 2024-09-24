package tc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/josharian/native"
	"github.com/mdlayher/netlink"
)

// tcConn defines a subset of netlink.Conn.
type tcConn interface {
	Close() error
	JoinGroup(group uint32) error
	LeaveGroup(group uint32) error
	Receive() ([]netlink.Message, error)
	Send(m netlink.Message) (netlink.Message, error)
	SetOption(option netlink.ConnOption, enable bool) error
	SetReadDeadline(t time.Time) error
}

var _ tcConn = &netlink.Conn{}

// Tc represents a RTNETLINK wrapper
type Tc struct {
	con tcConn

	logger *log.Logger
}

var nativeEndian = native.Endian

// Open establishes a RTNETLINK socket for traffic control
func Open(config *Config) (*Tc, error) {
	var tc Tc

	if config == nil {
		config = &Config{}
	}

	con, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{NetNS: config.NetNS})
	if err != nil {
		return nil, err
	}
	tc.con = con

	if config.Logger == nil {
		tc.logger = setDummyLogger()
	} else {
		tc.logger = config.Logger
	}

	return &tc, nil
}

// SetOption allows to enable or disable netlink socket options.
func (tc *Tc) SetOption(o netlink.ConnOption, enable bool) error {
	return tc.con.SetOption(o, enable)
}

// Close the connection
func (tc *Tc) Close() error {
	return tc.con.Close()
}

func (tc *Tc) query(req netlink.Message) ([]netlink.Message, error) {
	verify, err := tc.con.Send(req)
	if err != nil {
		return nil, err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return nil, err
	}

	return tc.con.Receive()
}

func (tc *Tc) action(action int, flags netlink.HeaderFlags, msg interface{}, opts []tcOption) error {
	tcminfo, err := marshalStruct(msg)
	if err != nil {
		return err
	}

	var data []byte
	data = append(data, tcminfo...)

	attrs, err := marshalAttributes(opts)
	if err != nil {
		return err
	}
	data = append(data, attrs...)
	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Acknowledge | flags,
		},
		Data: data,
	}

	msgs, err := tc.query(req)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		if msg.Header.Type == netlink.Error {
			// see https://www.infradead.org/~tgr/libnl/doc/core.html#core_errmsg
			tc.logger.Printf("received netlink.Error in action()\n")
		}
	}

	return nil
}

func (tc *Tc) get(action int, i *Msg) ([]Object, error) {
	var results []Object

	tcminfo, err := marshalStruct(i)
	if err != nil {
		return results, err
	}

	var data []byte
	data = append(data, tcminfo...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := tc.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		var result Object
		if err := unmarshalStruct(msg.Data[:20], &result.Msg); err != nil {
			return results, err
		}
		if err := extractTcmsgAttributes(action, msg.Data[20:], &result.Attribute); err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

// Object represents a generic traffic control object
type Object struct {
	Msg
	Attribute
}

// Msg represents a Traffic Control Message
type Msg struct {
	Family  uint32
	Ifindex uint32
	Handle  uint32
	Parent  uint32
	Info    uint32
}

// Attribute contains various elements for traffic control
type Attribute struct {
	Kind         string
	EgressBlock  *uint32
	IngressBlock *uint32
	HwOffload    *uint8
	Chain        *uint32
	Stats        *Stats
	XStats       *XStats
	Stats2       *Stats2
	Stab         *Stab
	ExtWarnMsg   string

	// Filters
	Basic    *Basic
	BPF      *Bpf
	Cgroup   *Cgroup
	U32      *U32
	Rsvp     *Rsvp
	Route4   *Route4
	Fw       *Fw
	Flow     *Flow
	Flower   *Flower
	Matchall *Matchall
	TcIndex  *TcIndex

	// Classless qdiscs
	Cake    *Cake
	FqCodel *FqCodel
	Codel   *Codel
	Fq      *Fq
	Pie     *Pie
	Hhf     *Hhf
	Tbf     *Tbf
	Sfb     *Sfb
	Sfq     *Sfq
	Red     *Red
	MqPrio  *MqPrio
	Pfifo   *FifoOpt
	Bfifo   *FifoOpt
	Choke   *Choke
	Netem   *Netem
	Plug    *Plug

	// Classful qdiscs
	Cbs      *Cbs
	Htb      *Htb
	Hfsc     *Hfsc
	HfscQOpt *HfscQOpt
	Dsmark   *Dsmark
	Drr      *Drr
	Cbq      *Cbq
	Atm      *Atm
	Qfq      *Qfq
	Prio     *Prio
	TaPrio   *TaPrio
}

// XStats contains further statistics to the TCA_KIND
type XStats struct {
	Sfb     *SfbXStats
	Sfq     *SfqXStats
	Red     *RedXStats
	Choke   *ChokeXStats
	Htb     *HtbXStats
	Cbq     *CbqXStats
	Codel   *CodelXStats
	Hhf     *HhfXStats
	Pie     *PieXStats
	FqCodel *FqCodelXStats
	Fq      *FqQdStats
	Hfsc    *HfscXStats
}

func marshalXStats(v XStats) ([]byte, error) {
	if v.Sfb != nil {
		return marshalStruct(v.Sfb)
	} else if v.Sfq != nil {
		return marshalStruct(v.Sfq)
	} else if v.Red != nil {
		return marshalStruct(v.Red)
	} else if v.Choke != nil {
		return marshalStruct(v.Choke)
	} else if v.Htb != nil {
		return marshalStruct(v.Htb)
	} else if v.Cbq != nil {
		return marshalStruct(v.Cbq)
	} else if v.Codel != nil {
		return marshalStruct(v.Codel)
	} else if v.Hhf != nil {
		return marshalStruct(v.Hhf)
	} else if v.Pie != nil {
		return marshalStruct(v.Pie)
	} else if v.FqCodel != nil {
		return marshalFqCodelXStats(v.FqCodel)
	}
	return []byte{}, fmt.Errorf("could not marshal XStat")
}

// HookFunc is a function, which is called for each altered RTNETLINK Object.
// Return something different than 0, to stop receiving messages.
// action will have the value of unix.RTM_[NEW|GET|DEL][QDISC|TCLASS|FILTER].
type HookFunc func(action uint16, m Object) int

// ErrorFunc is a function that receives all errors that happen while reading
// from a Netlinkgroup. To stop receiving messages return something different than 0.
type ErrorFunc func(e error) int

// MonitorWithErrorFunc handles NETLINK_ROUTE messages and calls for each HookFunc.
// Received errors tigger the given ErrorFunc.
func (tc *Tc) MonitorWithErrorFunc(ctx context.Context, deadline time.Duration,
	fn HookFunc, errfn ErrorFunc) error {
	return tc.monitor(ctx, deadline, fn, errfn)
}

// Monitor NETLINK_ROUTE messages
//
// Deprecated: Use MonitorWithErrorFunc() instead.
func (tc *Tc) Monitor(ctx context.Context, deadline time.Duration, fn HookFunc) error {
	return tc.monitor(ctx, deadline, fn, func(err error) int {
		if opError, ok := err.(*netlink.OpError); ok {
			if opError.Timeout() || opError.Temporary() {
				return 0
			}
		}
		tc.logger.Printf("Could not receive message: %v\n", err)
		return 1
	})
}

func (tc *Tc) monitor(ctx context.Context, deadline time.Duration,
	fn HookFunc, errfn ErrorFunc) error {
	ifinfomsg, err := marshalStruct(unix.IfInfomsg{
		Family: unix.AF_UNSPEC,
	})
	if err != nil {
		return err
	}

	rtattr, err := marshalAttributes([]tcOption{
		{Interpretation: vtUint32, Type: unix.IFLA_EXT_MASK, Data: uint32(1)},
	})
	if err != nil {
		return err
	}

	data := ifinfomsg
	data = append(data, rtattr...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(unix.RTM_GETLINK),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	if err := tc.con.JoinGroup(unix.RTNLGRP_TC); err != nil {
		return err
	}

	verify, err := tc.con.Send(req)
	if err != nil {
		tc.con.LeaveGroup(unix.RTNLGRP_TC)
		return err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		tc.con.LeaveGroup(unix.RTNLGRP_TC)
		return err
	}

	go func() {
		go func() {
			<-ctx.Done()
			stop := time.Now().Add(deadline)
			tc.con.SetReadDeadline(stop)
			tc.con.LeaveGroup(unix.RTNLGRP_TC)
		}()
		for {
			msgs, err := tc.con.Receive()
			if err != nil {
				if ret := errfn(err); ret != 0 {
					return
				}
				if ctx.Err() != nil {
					return
				}
				continue
			}
			for _, msg := range msgs {
				var monitored Object
				if err := unmarshalStruct(msg.Data[:20], &monitored.Msg); err != nil {
					tc.logger.Printf("could not extract tc.Msg from %v\n", msg.Data[:20])
					continue
				}
				if err := extractTcmsgAttributes(int(msg.Header.Type), msg.Data[20:],
					&monitored.Attribute); err != nil {
					tc.logger.Printf("could not extract attributes from %v\n", msg.Data[20:36])
					continue
				}
				if fn(uint16(msg.Header.Type), monitored) != 0 {
					return
				}
			}
		}
	}()
	return nil
}
