package dhcp6

import (
	"net"
	"reflect"
	"testing"
)

// TestParseRequest verifies that ParseRequest returns a consistent
// Request struct for use in Handler types.
func TestParseRequest(t *testing.T) {
	opt := option{
		Code: OptionClientID,
		Data: []byte{0, 1},
	}

	p := &Packet{
		MessageType:   MessageTypeSolicit,
		TransactionID: [3]byte{1, 2, 3},
		Options: Options{
			opt.Code: [][]byte{opt.Data},
		},
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 546,
	}

	buf, err := p.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	r := &Request{
		MessageType:   p.MessageType,
		TransactionID: p.TransactionID,
		Options:       make(Options),
		Length:        int64(len(buf)),
		RemoteAddr:    "[::1]:546",
	}
	r.Options.addRaw(opt.Code, opt.Data)

	gotR, err := ParseRequest(buf, addr)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := r, gotR; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected Request for ParseRequest(%v, %v)\n- want: %v\n-  got: %v",
			p, addr, want, got)
	}
}
