package dhcpv6

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRelayMsgParseOptRelayMsg(t *testing.T) {
	opt, err := parseOptRelayMsg([]byte{
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0, 0, // option value
	})
	if err != nil {
		t.Fatal(err)
	}
	if code := opt.Code(); code != OptionRelayMsg {
		t.Fatalf("Invalid option code. Expected OptionRelayMsg (%v), got %v",
			OptionRelayMsg, code,
		)
	}
}

func TestRelayMsgOptionsFromBytes(t *testing.T) {
	var opts Options
	err := opts.FromBytes([]byte{
		0, 9, // option: relay message
		0, 10, // relayed message length
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0, 0, // option value
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts) != 1 {
		t.Fatalf("Invalid number of options. Expected 1, got %v", len(opts))
	}
	opt := opts[0]
	if code := opt.Code(); code != OptionRelayMsg {
		t.Fatalf("Invalid option code. Expected OptionRelayMsg (%v), got %v",
			OptionRelayMsg, code,
		)
	}
}

func TestRelayMsgParseOptRelayMsgSingleEncapsulation(t *testing.T) {
	d, err := FromBytes([]byte{
		12,                                             // RELAY-FORW
		0,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, // peerAddr
		0, 9, // option: relay message
		0, 10, // relayed message length
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0x00, 0x01, // option value
	})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := d.(*RelayMessage)
	if !ok {
		t.Fatalf("Invalid DHCPv6 type. Expected RelayMessage, got %v",
			reflect.TypeOf(d),
		)
	}
	if mType := r.Type(); mType != MessageTypeRelayForward {
		t.Fatalf("Invalid messge type for relay. Expected %v, got %v", MessageTypeRelayForward, mType)
	}
	if len(r.Options.Options) != 1 {
		t.Fatalf("Invalid number of options. Expected 1, got %v", len(r.Options.Options))
	}
	ro := r.Options.RelayMessage()
	if ro == nil {
		t.Fatalf("No relay message")
	}
	innerDHCP, ok := ro.(*Message)
	if !ok {
		t.Fatalf("Invalid relay message type. Expected Message, got %v",
			reflect.TypeOf(innerDHCP),
		)
	}
	if dType := innerDHCP.Type(); dType != MessageTypeSolicit {
		t.Fatalf("Invalid inner DHCP type. Expected MessageTypeSolicit (%v), got %v",
			MessageTypeSolicit, dType,
		)
	}
	xid := TransactionID{0xaa, 0xbb, 0xcc}
	if tID := innerDHCP.TransactionID; tID != xid {
		t.Fatalf("Invalid inner DHCP transaction ID. Expected 0xaabbcc, got %v", tID)
	}
	if len(innerDHCP.Options.Options) != 1 {
		t.Fatalf("Invalid inner DHCP options length. Expected 1, got %v", len(innerDHCP.Options.Options))
	}
	eTime := innerDHCP.Options.ElapsedTime()
	if eTime != 10*time.Millisecond {
		t.Fatalf("Invalid elapsed time. Expected 0x1122, got 0x%04x", eTime)
	}
}

func TestSample(t *testing.T) {
	// Nested relay message. This test only checks if it parses correctly, but
	// could/should be extended to check all the fields like done in other tests
	buf := []byte{
		12,                                             // relay
		1,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // peerAddr
		// relay msg
		0, 9, // opt relay msg
		0, 66, // opt len
		// relay fwd
		12,                                             // relay
		0,                                              // hop count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // linkAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // peerAddr
		// opt interface ID
		0, 18, // interface id
		0, 6, // opt len
		0xba, 0xbe, 0xb1, 0xb0, 0xbe, 0xbe, // opt value
		// relay msg
		0, 9, // relay msg
		0, 18, // msg len
		// dhcpv6 msg
		1,                // solicit
		0xaa, 0xbb, 0xcc, // transaction ID
		// client ID
		0, 1, // opt client id
		0, 10, // opt len
		0, 3, // duid type
		0, 1, // hw type
		5, 6, 7, 8, 9, 10, // duid value
	}
	_, err := FromBytes(buf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRelayMsgParseOptRelayMsgTooShort(t *testing.T) {
	_, err := parseOptRelayMsg([]byte{
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		// no length/value for elapsed time option
	})
	require.Error(t, err, "ParseOptRelayMsg() should return an error if the encapsulated message is invalid")
}

func TestRelayMsgString(t *testing.T) {
	opt, err := parseOptRelayMsg([]byte{
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		0, 8, // option: elapsed time
		0, 2, // option length
		0, 0, // option value
	})
	require.NoError(t, err)
	require.Contains(
		t,
		opt.String(),
		"RelayMsg: Message",
		"String() should contain the relaymsg contents",
	)
}
