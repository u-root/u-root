package dhcpv6

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRelayMessage(t *testing.T) {
	ll := net.IP{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xbb, 0xcc, 0xff, 0xfe, 0xdd, 0xee, 0xff}
	ma := net.IP{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	r := RelayMessage{
		MessageType: MessageTypeRelayForward,
		HopCount:    10,
		LinkAddr:    ll,
		PeerAddr:    ma,
		// options is left empty here for testing purposes, even if it's
		// mandatory to have at least a relay message option
	}
	if mt := r.Type(); mt != MessageTypeRelayForward {
		t.Fatalf("Invalid message type. Expected %v, got %v", MessageTypeRelayForward, mt)
	}
	if hc := r.HopCount; hc != 10 {
		t.Fatalf("Invalid hop count. Expected 10, got %v", hc)
	}
	if la := r.LinkAddr; !la.Equal(ll) {
		t.Fatalf("Invalid link address. Expected %v, got %v", ll, la)
	}
	if pa := r.PeerAddr; !pa.Equal(ma) {
		t.Fatalf("Invalid peer address. Expected %v, got %v", ma, pa)
	}
	if opts := r.Options; len(opts) != 0 {
		t.Fatalf("Invalid options. Expected none, got %v", opts)
	}
}

func TestRelayMessageToBytes(t *testing.T) {
	expected := []byte{
		12,                                                      // MessageTypeRelayForward
		1,                                                       // hop count
		0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, // link addr
		0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, // peer addr
		// option relay message
		0, 9, // relay msg
		0, 10, // option length
		// inner dhcp solicit
		1,                // MessageTypeSolicit
		0xaa, 0xbb, 0xcc, // transaction ID
		// inner option - elapsed time
		0, 8, // elapsed time
		0, 2, // length
		0, 0,
	}
	r := RelayMessage{
		MessageType: MessageTypeRelayForward,
		HopCount:    1,
		LinkAddr:    net.IPv6interfacelocalallnodes,
		PeerAddr:    net.IPv6linklocalallrouters,
	}
	opt := OptRelayMsg{
		relayMessage: &Message{
			MessageType:   MessageTypeSolicit,
			TransactionID: TransactionID{0xaa, 0xbb, 0xcc},
			Options: []Option{
				&OptElapsedTime{
					ElapsedTime: 0,
				},
			},
		},
	}
	r.AddOption(&opt)
	relayBytes := r.ToBytes()
	if !bytes.Equal(expected, relayBytes) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, relayBytes)
	}
}

func TestNewRelayRepFromRelayForw(t *testing.T) {
	// create a new relay forward
	rf := RelayMessage{}
	rf.MessageType = MessageTypeRelayForward
	rf.PeerAddr = net.IPv6linklocalallrouters
	rf.LinkAddr = net.IPv6interfacelocalallnodes
	rf.AddOption(&OptInterfaceId{})
	rf.AddOption(&OptRemoteId{})

	// create the inner message
	s, err := NewMessage()
	require.NoError(t, err)
	s.AddOption(&OptClientId{})
	orm := OptRelayMsg{}
	orm.SetRelayMessage(s)
	rf.AddOption(&orm)

	a, err := NewAdvertiseFromSolicit(s)
	require.NoError(t, err)
	rr, err := NewRelayReplFromRelayForw(&rf, a)
	require.NoError(t, err)
	relay := rr.(*RelayMessage)
	require.Equal(t, rr.Type(), MessageTypeRelayReply)
	require.Equal(t, relay.HopCount, rf.HopCount)
	require.Equal(t, relay.PeerAddr, rf.PeerAddr)
	require.Equal(t, relay.LinkAddr, rf.LinkAddr)
	require.NotNil(t, rr.GetOneOption(OptionInterfaceID))
	require.NotNil(t, rr.GetOneOption(OptionRemoteID))
	m, err := relay.GetInnerMessage()
	require.NoError(t, err)
	require.Equal(t, m, a)

	_, err = NewRelayReplFromRelayForw(nil, a)
	require.Error(t, err)
	_, err = NewRelayReplFromRelayForw(&rf, nil)
	require.Error(t, err)
}
