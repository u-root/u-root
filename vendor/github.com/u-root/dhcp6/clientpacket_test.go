package dhcp6

import (
	"net"
	"reflect"
	"testing"
)

var mac = net.HardwareAddr([]byte{0xb8, 0xae, 0xed, 0x7a, 0x10, 0x66})

func TestNewSolicitOptions(t *testing.T) {
	options, err := newSolicitOptions(mac)
	if err != nil {
		t.Fatalf("error in newSolicitOptions: %v\n", err)
	}
	expected := Options(map[OptionCode][][]byte{
		OptionIANA:        [][]byte{[]byte{0x72, 0x6f, 0x6f, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		OptionRapidCommit: [][]byte{nil},
		OptionElapsedTime: [][]byte{[]byte{0x00, 0x00}},
		OptionORO:         [][]byte{[]byte{0x00, byte(OptionDNSServers), 0x00, byte(OptionDomainList), 0x00, byte(OptionBootFileURL), 0x00, byte(OptionBootFileParam)}},
		OptionClientID:    [][]byte{[]byte{0x00, 0x03, 0x00, 0x06, 0xb8, 0xae, 0xed, 0x7a, 0x10, 0x66}},
	})

	optionsIANA, optionsSuccess, err := options.IANA()
	if err != nil {
		t.Fatalf("getting IANA from options: got %v, want nil", err)
	}
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: IANA does not exist\n")
	}

	expectedIANA, expectedSuccess, err := expected.IANA()
	if err != nil {
		t.Fatalf("getting IANA from expected: got %v, want nil", err)
	}
	if expectedSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: IANA does not exist\n")
	}
	if !reflect.DeepEqual(optionsIANA, expectedIANA) {
		t.Fatalf(
			"incorrect newSolicitOptions: IANAs do not match, get %v, but should be %v instead\n",
			optionsIANA, expectedIANA,
		)
	}

	optionsRapidCommit, err := options.RapidCommit()
	if err != nil {
		t.Fatalf("getting RapidCommit from options: got %v, want nil", err)
	}
	expectedRapidCommit, err := expected.RapidCommit()
	if err != nil {
		t.Fatalf("getting RapidCommit from expected: got %v, want nil", err)
	}
	if !reflect.DeepEqual(expectedIANA, optionsIANA) {
		t.Fatalf(
			"incorrect newSolicitOptions: Rapid Commits do not match, get %v, but should be %v instead\n",
			optionsRapidCommit, expectedRapidCommit,
		)
	}

	optionsOR, optionsSuccess, err := options.OptionRequest()
	if err != nil {
		t.Fatalf("getting OptionRequest from options: got %v, want nil", err)
	}
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: Option request does not exist\n")
	}
	expectedOR, expectedSuccess, err := expected.OptionRequest()
	if err != nil {
		t.Fatalf("getting OptionRequest from expected: got %v, want nil", err)
	}
	if expectedSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: OptionRequest does not exist\n")
	}
	if !reflect.DeepEqual(optionsOR, expectedOR) {
		t.Fatalf(
			"incorrect newSolicitOptions: Option request do not match, get %v, but should be %v instead\n",
			optionsOR, expectedOR,
		)
	}

	optionsElapsedTime, optionsSuccess, err := options.ElapsedTime()
	if err != nil {
		t.Fatalf("getting ElapsedTime from options: got %v, want nil", err)
	}
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: ElapsedTime does not exist\n")
	}
	expectedElapsedTime, expectedSuccess, err := expected.ElapsedTime()
	if err != nil {
		t.Fatalf("getting ElapsedTime from expected: got %v, want nil", err)
	}
	if expectedSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: ElapsedTime does not exist\n")
	}
	if !reflect.DeepEqual(optionsElapsedTime, expectedElapsedTime) {
		t.Fatalf(
			"incorrect newSolicitOptions: Elapsed time do not match, get %v, but should be %v instead\n",
			optionsElapsedTime, expectedElapsedTime,
		)
	}

	optionsClientID, optionsSuccess, err := options.ClientID()
	if err != nil {
		t.Fatalf("getting ClientID from options: got %v, want nil", err)
	}
	if optionsSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: Client ID does not exist\n")
	}
	expectedClientID, expectedSuccess, err := expected.ClientID()
	if err != nil {
		t.Fatalf("getting ClientID from expected: got %v, want nil", err)
	}
	if expectedSuccess != true {
		t.Fatalf("incorrect newSolicitOptions: ClientID does not exist\n")
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
	p, err := newSolicitPacket(mac)
	if err != nil {
		t.Fatalf("error in newSolicitPacket: %v\n", err)
	}

	options, err := newSolicitOptions(mac)
	expected := &Packet{
		MessageType:   MessageTypeSolicit,
		TransactionID: [3]byte{0x00, 0x01, 0x02},
		Options:       options,
	}
	if !reflect.DeepEqual(p, expected) {
		t.Fatalf("incorrect newSolicitPacket: get %v but should be %v\n", p, expected)
	}
}
