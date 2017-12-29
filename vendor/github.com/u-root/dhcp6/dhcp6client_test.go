package dhcp6

import (
	"net"
	"reflect"
	"testing"
)

const (
	v6addr = "fe80::baae:edff:fe79:6191"
)

func TestSolicitAndAdvertise(t *testing.T) {
	p, _ := newSolicitPacket(mac)
	pb, _ := p.MarshalBinary()

	r := &testMessage{
		addr: &net.UDPAddr{
			IP: net.ParseIP(v6addr),
		},
	}
	r.b.Write(pb)

	reply, err := serve(r)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("Got a reply: %v", reply)
	if reply == nil {
		t.Fatalf("The reply is nil?")
	}

	if reply.MessageType != MessageTypeAdvertise {
		t.Fatalf("Reply does not have a correct typeShould be MessageTypeAdvertise but %v instead\n", MessageType(reply.MessageType))
	}
	t.Logf("It is a MessageTypeAdvertise")
	if !reflect.DeepEqual(reply.TransactionID, [3]byte{0x00, 0x01, 0x02}) {
		t.Fatalf("Reply txID does not match")
	}
	t.Logf("txID matches")
	iana, containsIANA, err := reply.Options.IANA()
	if !containsIANA {
		t.Fatalf("Reply does not contain IANA")
	}
	t.Logf("It contains an IANA")
	if err != nil {
		t.Fatalf("Reply does not contain valid IANA")
	}

	t.Logf("The IANA is valid")
	iaaddr, containsIAAddr, err := iana[0].Options.IAAddr()
	if !containsIAAddr {
		t.Fatalf("Reply does not contain IAAddr")
	}
	t.Logf("It contains an IAAddr")
	if err != nil {
		t.Fatalf("Reply does not contain valid IAAddr")
	}
	t.Logf("Get assigned ipv6 addr from server: %+v", iaaddr[0])
}
