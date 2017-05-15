package dhcp6client

import (
	"net"
	"reflect"
	"testing"

	"github.com/mdlayher/dhcp6"
)

const (
	v6addr = "fe80::baae:edff:fe79:6191"
)

func TestSolicitAndAdvertise(t *testing.T) {
	p, _ := newSolicitPacket(&mac)
	pb, _ := p.MarshalBinary()

	r := &testMessage{
		addr: &net.UDPAddr{
			IP: net.ParseIP(v6addr),
		},
	}
	r.b.Write(pb)

	reply, err := serve(r)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if reply.MessageType != dhcp6.MessageTypeAdvertise {
		t.Fatalf("Reply does not have a correct type\nShould be MessageTypeAdvertise but %v instead\n", dhcp6.MessageType(reply.MessageType))
	}
	if !reflect.DeepEqual(reply.TransactionID, [3]byte{0x00, 0x01, 0x02}) {
		t.Fatalf("Reply txID does not match\n")
	}
	iana, containsIANA, err := reply.Options.IANA()
	if !containsIANA {
		t.Fatalf("Reply does not contain IANA\n")
	}
	if err != nil {
		t.Fatalf("Reply does not contain valid IANA\n")
	}

	iaaddr, containsIAAddr, err := iana[0].Options.IAAddr()
	if !containsIAAddr {
		t.Fatalf("Reply does not contain IAAddr\n")
	}
	if err != nil {
		t.Fatalf("Reply does not contain valid IAAddr\n")
	}
	t.Logf("Get assigned ipv6 addr from server: %+v\n", iaaddr[0])
}
