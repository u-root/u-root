package dhcp6client

import (
	"math/rand"
	"net"
	"reflect"
	"testing"

	"github.com/mdlayher/dhcp6"
	"golang.org/x/net/ipv6"
	"golang.org/x/sys/unix"
)

var mac = net.HardwareAddr([]byte{0xb8, 0xae, 0xed, 0x7a, 0x10, 0x66})

func TestNewSolicitOptions(t *testing.T) {
	options, err := newSolicitOptions(&mac)
	if err != nil {
		t.Fatalf("error in newSolicitOptions: %v\n", err)
	}
	expected := dhcp6.Options(map[dhcp6.OptionCode][][]byte{
		dhcp6.OptionIANA:        [][]byte{[]byte{0x72, 0x6f, 0x6f, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		dhcp6.OptionRapidCommit: [][]byte{nil},
		dhcp6.OptionElapsedTime: [][]byte{[]byte{0x00, 0x00}},
		dhcp6.OptionORO:         [][]byte{[]byte{0x00, 0x17, 0x00, 0x18}},
		dhcp6.OptionClientID:    [][]byte{[]byte{0x00, 0x03, 0x00, 0x06, 0xb8, 0xae, 0xed, 0x7a, 0x10, 0x66}},
	})

	optionsIANA, optionsSuccess, _ := options.IANA()
	expectedIANA, _, _ := expected.IANA()
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: IANA does not exist\n")
	}
	if !reflect.DeepEqual(optionsIANA, expectedIANA) {
		t.Fatalf(
			"incorrect newSolicitOptions: IANAs do not match, get %v, but should be %v instead\n",
			optionsIANA, expectedIANA,
		)
	}

	optionsRapidCommit, _ := options.RapidCommit()
	expectedRapidCommit, _ := expected.RapidCommit()
	if !reflect.DeepEqual(expectedIANA, optionsIANA) {
		t.Fatalf(
			"incorrect newSolicitOptions: Rapid Commits do not match, get %v, but should be %v instead\n",
			optionsRapidCommit, expectedRapidCommit,
		)
	}

	optionsOR, optionsSuccess, _ := options.OptionRequest()
	expectedOR, _, _ := expected.OptionRequest()
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: Option request does not exist\n")
	}
	if !reflect.DeepEqual(optionsOR, expectedOR) {
		t.Fatalf(
			"incorrect newSolicitOptions: Option request do not match, get %v, but should be %v instead\n",
			optionsOR, expectedOR,
		)
	}

	optionsElapsedTime, optionsSuccess, _ := options.ElapsedTime()
	expectedElapsedTime, _, _ := expected.ElapsedTime()
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: ElapsedTime does not exist\n")
	}
	if !reflect.DeepEqual(optionsElapsedTime, expectedElapsedTime) {
		t.Fatalf(
			"incorrect newSolicitOptions: Elapsed time do not match, get %v, but should be %v instead\n",
			optionsElapsedTime, expectedElapsedTime,
		)
	}

	optionsClientID, optionsSuccess, _ := options.ClientID()
	expectedClientID, _, _ := expected.ClientID()
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: Client ID does not exist\n")
	}
	if !reflect.DeepEqual(optionsClientID, expectedClientID) {
		t.Fatalf(
			"incorrect newSolicitOptions: Client IDs do not match, get %v, but should be %v instead\n",
			optionsClientID, expectedClientID,
		)
	}

	if !reflect.DeepEqual(expected, options) {
		t.Fatalf("incorrect newSolicitOptions: extra unnecessary options\n%v\n%v\n", expected, options)
	}
}

func TestNewSolicitPacket(t *testing.T) {
	p, err := newSolicitPacket(&mac)
	if err != nil {
		t.Fatalf("error in newSolicitPacket: %v\n", err)
	}

	options, err := newSolicitOptions(&mac)
	expected := &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x00, 0x01, 0x02},
		Options:       options,
	}
	if !reflect.DeepEqual(p, expected) {
		t.Fatalf("incorrect newSolicitPacket: get %v but should be %v\n", p, expected)
	}
}

func TestMarshalPacket(t *testing.T) {
	p, err := newSolicitPacket(&mac)
	if err != nil {
		t.Fatalf("error in newSolicitPacket: %v\n", err)
	}

	pb, err := p.MarshalBinary()
	if err != nil {
		t.Fatalf("packet marshal to binary err: %v\n", err)
	}

	h1 := &ipv6.Header{
		Version:      ipv6.Version,
		TrafficClass: 0,
		FlowLabel:    rand.Int() & 0xfffff,
		PayloadLen:   udpHdrLen + len(pb),
		NextHeader:   unix.IPPROTO_UDP,
		HopLimit:     1,
		Src:          net.ParseIP("::0"),
		Dst:          net.ParseIP("FF02::1:2"),
	}

	h2 := &Udphdr{
		Src:    srcPort,
		Dst:    dstPort,
		Length: uint16(udpHdrLen + len(pb)),
	}
	pkt, err := marshalPacket(h1, h2, pb)
	if err != nil {
		t.Fatalf("failed to make a new packet: %v\n", err)
	}

	// Test if identical after marshal and unmarshal
	h1p, h2p, pp, err := unmarshalPacket(pkt)
	if err != nil {
		t.Fatalf("failed to unmarshal packet: %v\n", err)
	}
	if !reflect.DeepEqual(h1, h1p) {
		t.Fatalf("ip headers do not match:\n%v\n%v\n", h1, h1p)
	}
	if !reflect.DeepEqual(h2, h2p) {
		t.Fatalf("udp headers do not match:\n%v\n%v\n", h2, h2p)
	}

	// Packet unmarshal ignores rapid commit so it has to be added back
	rc, _ := p.Options.RapidCommit()
	if rc {
		pp.Options[dhcp6.OptionRapidCommit] = append(pp.Options[dhcp6.OptionRapidCommit], nil)
	}
	if !reflect.DeepEqual(p, pp) {
		t.Fatalf("dhcp packets do not match:\n%v\n%v\n", p, pp)
	}
}
