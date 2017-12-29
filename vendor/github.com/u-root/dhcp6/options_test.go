package dhcp6

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"io"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

// TestOptionsAddBinaryMarshaler verifies that Options.Add correctly creates or
// appends OptionCode keys with BinaryMarshaler bytes values to an Options map.
func TestOptionsAddBinaryMarshaler(t *testing.T) {
	var tests = []struct {
		desc    string
		code    OptionCode
		bin     encoding.BinaryMarshaler
		options Options
	}{
		{
			desc: "DUID-LLT",
			code: OptionClientID,
			bin: &DUIDLLT{
				Type:         DUIDTypeLLT,
				HardwareType: 1,
				Time:         duidLLTTime.Add(1 * time.Minute).Sub(duidLLTTime),
				HardwareAddr: net.HardwareAddr([]byte{0, 1, 0, 1, 0, 1}),
			},
			options: Options{
				OptionClientID: [][]byte{{
					0, 1,
					0, 1,
					0, 0, 0, 60,
					0, 1, 0, 1, 0, 1,
				}},
			},
		},
		{
			desc: "DUID-EN",
			code: OptionClientID,
			bin: &DUIDEN{
				Type:             DUIDTypeEN,
				EnterpriseNumber: 100,
				Identifier:       []byte{0, 1, 2, 3, 4},
			},
			options: Options{
				OptionClientID: [][]byte{{
					0, 2,
					0, 0, 0, 100,
					0, 1, 2, 3, 4,
				}},
			},
		},
		{
			desc: "DUID-LL",
			code: OptionClientID,
			bin: &DUIDLL{
				Type:         DUIDTypeLL,
				HardwareType: 1,
				HardwareAddr: net.HardwareAddr([]byte{0, 1, 0, 1, 0, 1}),
			},
			options: Options{
				OptionClientID: [][]byte{{
					0, 3,
					0, 1,
					0, 1, 0, 1, 0, 1,
				}},
			},
		},
		{
			desc: "DUID-UUID",
			code: OptionClientID,
			bin: &DUIDUUID{
				Type: DUIDTypeUUID,
				UUID: [16]byte{
					1, 1, 1, 1,
					2, 2, 2, 2,
					3, 3, 3, 3,
					4, 4, 4, 4,
				},
			},
			options: Options{
				OptionClientID: [][]byte{{
					0, 4,
					1, 1, 1, 1,
					2, 2, 2, 2,
					3, 3, 3, 3,
					4, 4, 4, 4,
				}},
			},
		},
		{
			desc: "IA_NA",
			code: OptionIANA,
			bin: &IANA{
				IAID: [4]byte{0, 1, 2, 3},
				T1:   30 * time.Second,
				T2:   60 * time.Second,
			},
			options: Options{
				OptionIANA: [][]byte{{
					0, 1, 2, 3,
					0, 0, 0, 30,
					0, 0, 0, 60,
				}},
			},
		},
		{
			desc: "IA_TA",
			code: OptionIATA,
			bin: &IATA{
				IAID: [4]byte{0, 1, 2, 3},
			},
			options: Options{
				OptionIATA: [][]byte{{
					0, 1, 2, 3,
				}},
			},
		},
		{
			desc: "IAAddr",
			code: OptionIAAddr,
			bin: &IAAddr{
				IP:                net.IPv6loopback,
				PreferredLifetime: 30 * time.Second,
				ValidLifetime:     60 * time.Second,
			},
			options: Options{
				OptionIAAddr: [][]byte{{
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
					0, 0, 0, 30,
					0, 0, 0, 60,
				}},
			},
		},
		{
			desc: "Preference",
			code: OptionPreference,
			bin:  Preference(255),
			options: Options{
				OptionPreference: [][]byte{{255}},
			},
		},
		{
			desc: "ElapsedTime",
			code: OptionElapsedTime,
			bin:  ElapsedTime(60 * time.Second),
			options: Options{
				OptionElapsedTime: [][]byte{{23, 112}},
			},
		},
		{
			desc: "Unicast IP",
			code: OptionUnicast,
			bin:  IP(net.IPv6loopback),
			options: Options{
				OptionUnicast: [][]byte{{
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 1,
				}},
			},
		},
		{
			desc: "StatusCode",
			code: OptionStatusCode,
			bin: &StatusCode{
				Code:    StatusSuccess,
				Message: "hello world",
			},
			options: Options{
				OptionStatusCode: [][]byte{{
					0, 0,
					'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
				}},
			},
		},
		{
			desc: "RapidCommit",
			code: OptionRapidCommit,
			bin:  nil,
			options: Options{
				OptionRapidCommit: [][]byte{nil},
			},
		},
		{
			desc: "Data (UserClass, VendorClass, BootFileParam)",
			code: OptionUserClass,
			bin: Data{
				[]byte{0},
				[]byte{0, 1},
				[]byte{0, 1, 2},
			},
			options: Options{
				OptionUserClass: [][]byte{{
					0, 1, 0,
					0, 2, 0, 1,
					0, 3, 0, 1, 2,
				}},
			},
		},
		{
			desc: "IA_PD",
			code: OptionIAPD,
			bin: &IAPD{
				IAID: [4]byte{0, 1, 2, 3},
				T1:   30 * time.Second,
				T2:   60 * time.Second,
			},
			options: Options{
				OptionIAPD: [][]byte{{
					0, 1, 2, 3,
					0, 0, 0, 30,
					0, 0, 0, 60,
				}},
			},
		},
		{
			desc: "IAPrefix",
			code: OptionIAPrefix,
			bin: &IAPrefix{
				PreferredLifetime: 30 * time.Second,
				ValidLifetime:     60 * time.Second,
				PrefixLength:      64,
				Prefix: net.IP{
					1, 1, 1, 1, 1, 1, 1, 1,
					0, 0, 0, 0, 0, 0, 0, 0,
				},
			},
			options: Options{
				OptionIAPrefix: [][]byte{{
					0, 0, 0, 30,
					0, 0, 0, 60,
					64,
					1, 1, 1, 1, 1, 1, 1, 1,
					0, 0, 0, 0, 0, 0, 0, 0,
				}},
			},
		},
		{
			desc: "URL",
			code: OptionBootFileURL,
			bin: &URL{
				Scheme: "tftp",
				Host:   "192.168.1.1:69",
			},
			options: Options{
				OptionBootFileURL: [][]byte{[]byte("tftp://192.168.1.1:69")},
			},
		},
		{
			desc: "ArchTypes",
			code: OptionClientArchType,
			bin: ArchTypes{
				ArchTypeEFIx8664,
				ArchTypeIntelx86PC,
				ArchTypeIntelLeanClient,
			},
			options: Options{
				OptionClientArchType: [][]byte{[]byte{0, 9, 0, 0, 0, 5}},
			},
		},
		{
			desc: "NII",
			code: OptionNII,
			bin: &NII{
				Type:  1,
				Major: 2,
				Minor: 3,
			},
			options: Options{
				OptionNII: [][]byte{[]byte{1, 2, 3}},
			},
		},
	}

	for i, tt := range tests {
		o := make(Options)
		if err := o.Add(tt.code, tt.bin); err != nil {
			t.Fatal(err)
		}

		if want, got := tt.options, o; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Options map:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptions_addRaw verifies that Options.addRaw correctly creates or appends
// key/value Option pairs to an Options map.
func TestOptions_addRaw(t *testing.T) {
	var tests = []struct {
		desc    string
		kv      []option
		options Options
	}{
		{
			desc: "one key/value pair",
			kv: []option{
				{
					Code: 1,
					Data: []byte("foo"),
				},
			},
			options: Options{
				1: [][]byte{[]byte("foo")},
			},
		},
		{
			desc: "two key/value pairs",
			kv: []option{
				{
					Code: 1,
					Data: []byte("foo"),
				},
				{
					Code: 2,
					Data: []byte("bar"),
				},
			},
			options: Options{
				1: [][]byte{[]byte("foo")},
				2: [][]byte{[]byte("bar")},
			},
		},
		{
			desc: "three key/value pairs, two with same key",
			kv: []option{
				{
					Code: 1,
					Data: []byte("foo"),
				},
				{
					Code: 1,
					Data: []byte("baz"),
				},
				{
					Code: 2,
					Data: []byte("bar"),
				},
			},
			options: Options{
				1: [][]byte{[]byte("foo"), []byte("baz")},
				2: [][]byte{[]byte("bar")},
			},
		},
	}

	for i, tt := range tests {
		o := make(Options)
		for _, p := range tt.kv {
			o.addRaw(p.Code, p.Data)
		}

		if want, got := tt.options, o; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Options map:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsGet verifies that Options.Get correctly selects the first value
// for a given key, if the value is not empty in an Options map.
func TestOptionsGet(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		key     OptionCode
		value   []byte
		ok      bool
	}{
		{
			desc: "nil Options map",
		},
		{
			desc:    "empty Options map",
			options: Options{},
		},
		{
			desc: "value not present in Options map",
			options: Options{
				2: [][]byte{[]byte("foo")},
			},
			key: 1,
		},
		{
			desc: "value present in Options map, but zero length value for key",
			options: Options{
				1: [][]byte{},
			},
			key: 1,
			ok:  true,
		},
		{
			desc: "value present in Options map",
			options: Options{
				1: [][]byte{[]byte("foo")},
			},
			key:   1,
			value: []byte("foo"),
			ok:    true,
		},
		{
			desc: "value present in Options map, with multiple values",
			options: Options{
				1: [][]byte{[]byte("foo"), []byte("bar")},
			},
			key:   1,
			value: []byte("foo"),
			ok:    true,
		},
	}

	for i, tt := range tests {
		value, ok := tt.options.Get(tt.key)

		if want, got := tt.value, value; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.Get(%v):\n- want: %v\n-  got: %v",
				i, tt.desc, tt.key, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.Get(%v): %v != %v",
				i, tt.desc, tt.key, want, got)
		}
	}
}

// TestOptionsClientID verifies that Options.ClientID properly parses and returns
// a DUID value, if one is available with OptionClientID.
func TestOptionsClientID(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		duid    DUID
		ok      bool
	}{
		{
			desc: "OptionClientID not present in Options map",
		},
		{
			desc: "OptionClientID present in Options map",
			options: Options{
				OptionClientID: [][]byte{{
					0, 3,
					0, 1,
					0, 1, 0, 1, 0, 1,
				}},
			},
			duid: &DUIDLL{
				Type:         DUIDTypeLL,
				HardwareType: 1,
				HardwareAddr: []byte{0, 1, 0, 1, 0, 1},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		// DUID parsing is tested elsewhere, so errors should automatically fail
		// test here
		duid, ok, err := tt.options.ClientID()
		if err != nil {
			t.Fatal(err)
		}

		if want, got := tt.duid, duid; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.ClientID():\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.ClientID(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsServerID verifies that Options.ServerID properly parses and returns
// a DUID value, if one is available with OptionServerID.
func TestOptionsServerID(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		duid    DUID
		ok      bool
	}{
		{
			desc: "OptionServerID not present in Options map",
		},
		{
			desc: "OptionServerID present in Options map",
			options: Options{
				OptionServerID: [][]byte{{
					0, 3,
					0, 1,
					0, 1, 0, 1, 0, 1,
				}},
			},
			duid: &DUIDLL{
				Type:         DUIDTypeLL,
				HardwareType: 1,
				HardwareAddr: []byte{0, 1, 0, 1, 0, 1},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		// DUID parsing is tested elsewhere, so errors should automatically fail
		// test here
		duid, ok, err := tt.options.ServerID()
		if err != nil {
			t.Fatal(err)
		}

		if want, got := tt.duid, duid; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.ServerID():\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.ServerID(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsIANA verifies that Options.IANA properly parses and
// returns multiple IANA values, if one or more are available with OptionIANA.
func TestOptionsIANA(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		iana    []*IANA
		ok      bool
		err     error
	}{
		{
			desc: "OptionIANA not present in Options map",
		},
		{
			desc: "OptionIANA present in Options map, but too short",
			options: Options{
				OptionIANA: [][]byte{bytes.Repeat([]byte{0}, 11)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "one OptionIANA present in Options map",
			options: Options{
				OptionIANA: [][]byte{{
					1, 2, 3, 4,
					0, 0, 0, 30,
					0, 0, 0, 60,
				}},
			},
			iana: []*IANA{
				{
					IAID: [4]byte{1, 2, 3, 4},
					T1:   30 * time.Second,
					T2:   60 * time.Second,
				},
			},
			ok: true,
		},
		{
			desc: "two OptionIANA present in Options map",
			options: Options{
				OptionIANA: [][]byte{
					append(bytes.Repeat([]byte{0}, 12), []byte{0, 1, 0, 1, 1}...),
					append(bytes.Repeat([]byte{0}, 12), []byte{0, 2, 0, 1, 2}...),
				},
			},
			iana: []*IANA{
				{
					Options: Options{
						OptionClientID: [][]byte{{1}},
					},
				},
				{
					Options: Options{
						OptionServerID: [][]byte{{2}},
					},
				},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		iana, ok, err := tt.options.IANA()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.IANA: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		for j := range tt.iana {
			want, err := tt.iana[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			got, err := iana[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.IANA():\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.IANA(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsIATA verifies that Options.IATA properly parses and
// returns multiple IATA values, if one or more are available with OptionIATA.
func TestOptionsIATA(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		iata    []*IATA
		ok      bool
		err     error
	}{
		{
			desc: "OptionIATA not present in Options map",
		},
		{
			desc: "OptionIATA present in Options map, but too short",
			options: Options{
				OptionIATA: [][]byte{{0, 0, 0}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "one OptionIATA present in Options map",
			options: Options{
				OptionIATA: [][]byte{{
					1, 2, 3, 4,
				}},
			},
			iata: []*IATA{
				{
					IAID: [4]byte{1, 2, 3, 4},
				},
			},
			ok: true,
		},
		{
			desc: "two OptionIATA present in Options map",
			options: Options{
				OptionIATA: [][]byte{
					[]byte{0, 1, 2, 3, 0, 1, 0, 1, 1},
					[]byte{4, 5, 6, 7, 0, 2, 0, 1, 2},
				},
			},
			iata: []*IATA{
				{
					IAID: [4]byte{0, 1, 2, 3},
					Options: Options{
						OptionClientID: [][]byte{{1}},
					},
				},
				{
					IAID: [4]byte{4, 5, 6, 7},
					Options: Options{
						OptionServerID: [][]byte{{2}},
					},
				},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		iata, ok, err := tt.options.IATA()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.IATA: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		for j := range tt.iata {
			want, err := tt.iata[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			got, err := iata[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.IATA():\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.IATA(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsIAAddr verifies that Options.IAAddr properly parses and
// returns multiple IAAddr values, if one or more are available with
// OptionIAAddr.
func TestOptionsIAAddr(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		iaaddr  []*IAAddr
		ok      bool
		err     error
	}{
		{
			desc: "OptionIAAddr not present in Options map",
		},
		{
			desc: "OptionIAAddr present in Options map, but too short",
			options: Options{
				OptionIAAddr: [][]byte{bytes.Repeat([]byte{0}, 23)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "one OptionIAAddr present in Options map",
			options: Options{
				OptionIAAddr: [][]byte{{
					0, 0, 0, 0,
					1, 1, 1, 1,
					2, 2, 2, 2,
					3, 3, 3, 3,
					0, 0, 0, 30,
					0, 0, 0, 60,
				}},
			},
			iaaddr: []*IAAddr{
				{
					IP: net.IP{
						0, 0, 0, 0,
						1, 1, 1, 1,
						2, 2, 2, 2,
						3, 3, 3, 3,
					},
					PreferredLifetime: 30 * time.Second,
					ValidLifetime:     60 * time.Second,
				},
			},
			ok: true,
		},
		{
			desc: "two OptionIAAddr present in Options map",
			options: Options{
				OptionIAAddr: [][]byte{
					bytes.Repeat([]byte{0}, 24),
					bytes.Repeat([]byte{0}, 24),
				},
			},
			iaaddr: []*IAAddr{
				{
					IP: net.IPv6zero,
				},
				{
					IP: net.IPv6zero,
				},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		iaaddr, ok, err := tt.options.IAAddr()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.IAAddr: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		for j := range tt.iaaddr {
			want, err := tt.iaaddr[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			got, err := iaaddr[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.IAAddr():\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.IAAddr(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsOptionRequest verifies that Options.OptionRequest properly parses
// and returns a slice of OptionCode values, if they are available with
// OptionORO.
func TestOptionsOptionRequest(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		codes   OptionRequestOption
		ok      bool
		err     error
	}{
		{
			desc: "OptionORO not present in Options map",
		},
		{
			desc: "OptionORO present in Options map, but not even length",
			options: Options{
				OptionORO: [][]byte{{0}},
			},
			err: errInvalidOptionRequest,
		},
		{
			desc: "OptionORO present in Options map",
			options: Options{
				OptionORO: [][]byte{{0, 1}},
			},
			codes: []OptionCode{1},
			ok:    true,
		},
		{
			desc: "OptionORO present in Options map, with multiple values",
			options: Options{
				OptionORO: [][]byte{{0, 1, 0, 2, 0, 3}},
			},
			codes: []OptionCode{1, 2, 3},
			ok:    true,
		},
	}

	for i, tt := range tests {
		codes, ok, err := tt.options.OptionRequest()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.OptionRequest(): %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.codes, codes; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.OptionRequest():\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.OptionRequest(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsPreference verifies that Options.Preference properly parses
// and returns an integer value, if it is available with OptionPreference.
func TestOptionsPreference(t *testing.T) {
	var tests = []struct {
		desc       string
		options    Options
		preference Preference
		ok         bool
		err        error
	}{
		{
			desc: "OptionPreference not present in Options map",
		},
		{
			desc: "OptionPreference present in Options map, but too short length",
			options: Options{
				OptionPreference: [][]byte{{}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionPreference present in Options map, but too long length",
			options: Options{
				OptionPreference: [][]byte{{0, 1}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionPreference present in Options map",
			options: Options{
				OptionPreference: [][]byte{{255}},
			},
			preference: 255,
			ok:         true,
		},
	}

	for i, tt := range tests {
		preference, ok, err := tt.options.Preference()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.Preference(): %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.preference, preference; want != got {
			t.Fatalf("[%02d] test %q, unexpected value for Options.Preference(): %v != %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.Preference(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsUnicast verifies that Options.Unicast properly parses
// and returns an IPv6 address or an error, if available with OptionUnicast.
func TestOptionsUnicast(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		ip      IP
		ok      bool
		err     error
	}{
		{
			desc: "OptionUnicast not present in Options map",
		},
		{
			desc: "OptionUnicast present in Options map, but too short length",
			options: Options{
				OptionUnicast: [][]byte{bytes.Repeat([]byte{0}, 15)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionUnicast present in Options map, but too long length",
			options: Options{
				OptionUnicast: [][]byte{bytes.Repeat([]byte{0}, 17)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionUnicast present in Options map with IPv4 address",
			options: Options{
				OptionUnicast: [][]byte{net.IPv4(192, 168, 1, 1)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionUnicast present in Options map with IPv6 address",
			options: Options{
				OptionUnicast: [][]byte{net.IPv6loopback},
			},
			ip: IP(net.IPv6loopback),
			ok: true,
		},
	}

	for i, tt := range tests {
		ip, ok, err := tt.options.Unicast()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.Unicast(): %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.ip, ip; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.Unicast():\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.Unicast(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsStatusCode verifies that Options.StatusCode properly parses
// and returns a StatusCode value, if it is available with OptionStatusCode.
func TestOptionsStatusCode(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		sc      *StatusCode
		ok      bool
		err     error
	}{
		{
			desc: "OptionStatusCode not present in Options map",
		},
		{
			desc: "OptionStatusCode present in Options map, but too short length",
			options: Options{
				OptionStatusCode: [][]byte{{}},
			},
			err: errInvalidStatusCode,
		},
		{
			desc: "OptionStatusCode present in Options map, no message",
			options: Options{
				OptionStatusCode: [][]byte{{0, 0}},
			},
			sc: &StatusCode{
				Code: StatusSuccess,
			},
			ok: true,
		},
		{
			desc: "OptionStatusCode present in Options map, with message",
			options: Options{
				OptionStatusCode: [][]byte{append([]byte{0, 0}, []byte("deadbeef")...)},
			},
			sc: &StatusCode{
				Code:    StatusSuccess,
				Message: "deadbeef",
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		sc, ok, err := tt.options.StatusCode()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.StatusCode(): %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.sc, sc; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.StatusCode():\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.StatusCode(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsElapsedTime verifies that Options.ElapsedTime properly parses and
// returns a time.Duration value, if one is available with OptionElapsedTime.
func TestOptionsElapsedTime(t *testing.T) {
	var tests = []struct {
		desc     string
		options  Options
		duration ElapsedTime
		ok       bool
		err      error
	}{
		{
			desc: "OptionElapsedTime not present in Options map",
		},
		{
			desc: "OptionElapsedTime present in Options map, but too short",
			options: Options{
				OptionElapsedTime: [][]byte{{1}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionElapsedTime present in Options map, but too long",
			options: Options{
				OptionElapsedTime: [][]byte{{1, 2, 3}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionElapsedTime present in Options map",
			options: Options{
				OptionElapsedTime: [][]byte{{1, 1}},
			},
			duration: ElapsedTime(2570 * time.Millisecond),
			ok:       true,
		},
	}

	for i, tt := range tests {
		duration, ok, err := tt.options.ElapsedTime()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.ElapsedTime: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.duration, duration; want != got {
			t.Fatalf("[%02d] test %q, unexpected value for Options.ElapsedTime(): %v != %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.ElapsedTime(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestElapsedTimeMarshalBinary verifies that Options.ElapsedTime properly
// marsharls into bytes array.
func TestElapsedTimeMarshalBinary(t *testing.T) {
	var tests = []struct {
		desc        string
		elapsedTime ElapsedTime
		buf         []byte
		err         error
	}{
		{
			desc: "OptionElapsedTime elapsed-time = 0",
			buf:  []byte{0, 0},
		},
		{
			desc:        "OptionElapsedTime elapsed-time = 65534 hundredths of a second",
			elapsedTime: ElapsedTime(655340 * time.Millisecond),
			buf:         []byte{0xff, 0xfe},
		},
		{
			desc:        "OptionElapsedTime elapsed-time = 65535 hundredths of a second",
			elapsedTime: ElapsedTime(655350 * time.Millisecond),
			buf:         []byte{0xff, 0xff},
		},
		{
			desc:        "OptionElapsedTime elapsed-time = 65537 hundredths of a second",
			elapsedTime: ElapsedTime(655370 * time.Millisecond),
			buf:         []byte{0xff, 0xff},
		},
	}

	for i, tt := range tests {
		buf, err := tt.elapsedTime.MarshalBinary()
		if want, got := tt.err, err; want != got {
			t.Fatalf("[%02d] test %q, unexpected error for Options.ElapsedTime\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if tt.err != nil {
			continue
		}

		if want, got := tt.buf, buf; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected error for Options.ElapsedTime\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsRelayMessage verifies that Options.RelayMessageOption properly parses and
// returns an relay message option value, if one is available with RelayMessageOption.
func TestOptionsRelayMessage(t *testing.T) {
	var tests = []struct {
		desc           string
		options        Options
		authentication RelayMessageOption
		ok             bool
		err            error
	}{
		{
			desc: "RelayMessageOption not present in Options map",
		},
		{
			desc: "RelayMessageOption present in Options map",
			options: Options{
				OptionRelayMsg: [][]byte{{1, 1, 2, 3}},
			},
			authentication: []byte{1, 1, 2, 3},
			ok:             true,
		},
	}

	for i, tt := range tests {
		relayMsg, ok, err := tt.options.RelayMessageOption()
		if want, got := tt.err, err; want != got {
			t.Fatalf("[%02d] test %q, unexpected error for Options.RelayMessageOption\n- want: %v\n-  got: %v", i, tt.desc, want, got)
		}

		if tt.err != nil {
			continue
		}

		if want, got := tt.authentication, relayMsg; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.RelayMessageOption()\n- want: %v\n-  got: %v", i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.RelayMessageOption(): %v != %v", i, tt.desc, want, got)
		}
	}
}

// TestAuthentication verifies that Options.Authentication properly parses and
// returns an authentication value, if one is available with Authentication.
func TestAuthentication(t *testing.T) {
	var tests = []struct {
		desc           string
		options        Options
		authentication *Authentication
		ok             bool
		err            error
	}{
		{
			desc: "Authentication not present in Options map",
		},
		{
			desc: "Authentication present in Options map, but too short",
			options: Options{
				OptionAuth: [][]byte{bytes.Repeat([]byte{0}, 10)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "Authentication present in Options map",
			options: Options{
				OptionAuth: [][]byte{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}},
			},
			authentication: &Authentication{
				Protocol:                  0,
				Algorithm:                 1,
				RDM:                       2,
				ReplayDetection:           binary.BigEndian.Uint64([]byte{3, 4, 5, 6, 7, 8, 9, 0xa}),
				AuthenticationInformation: []byte{0xb, 0xc, 0xd, 0xe, 0xf},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		authentication, ok, err := tt.options.Authentication()
		if want, got := tt.err, err; want != got {
			t.Fatalf("[%02d] test %q, unexpected error for Options.Authentication\n- want: %v\n-  got: %v", i, tt.desc, want, got)
		}

		if tt.err != nil {
			continue
		}

		if want, got := tt.authentication, authentication; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.Authentication()\n- want: %v\n-  got: %v", i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.Authentication(): %v != %v", i, tt.desc, want, got)
		}
	}
}

// TestOptionsRapidCommit verifies that Options.RapidCommit properly indicates
// if OptionRapidCommit was present in Options.
func TestOptionsRapidCommit(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		ok      bool
		err     error
	}{
		{
			desc: "OptionRapidCommit not present in Options map",
		},
		{
			desc: "OptionRapidCommit present in Options map, but non-empty",
			options: Options{
				OptionRapidCommit: [][]byte{{1}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionRapidCommit present in Options map, empty",
			options: Options{
				OptionRapidCommit: [][]byte{},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		ok, err := tt.options.RapidCommit()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.RapidCommit: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.RapidCommit(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsUserClass verifies that Options.UserClass properly parses
// and returns raw user class data, if it is available with OptionUserClass.
func TestOptionsUserClass(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		classes [][]byte
		ok      bool
		err     error
	}{
		{
			desc: "OptionUserClass not present in Options map",
		},
		{
			desc: "OptionUserClass present in Options map, but empty",
			options: Options{
				OptionUserClass: [][]byte{{}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionUserClass present in Options map, one item, zero length",
			options: Options{
				OptionUserClass: [][]byte{{
					0, 0,
				}},
			},
			classes: [][]byte{{}},
			ok:      true,
		},
		{
			desc: "OptionUserClass present in Options map, one item, extra byte",
			options: Options{
				OptionUserClass: [][]byte{{
					0, 1, 1, 255,
				}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionUserClass present in Options map, one item",
			options: Options{
				OptionUserClass: [][]byte{{
					0, 1, 1,
				}},
			},
			classes: [][]byte{{1}},
			ok:      true,
		},
		{
			desc: "OptionUserClass present in Options map, three items",
			options: Options{
				OptionUserClass: [][]byte{{
					0, 1, 1,
					0, 2, 2, 2,
					0, 3, 3, 3, 3,
				}},
			},
			classes: [][]byte{{1}, {2, 2}, {3, 3, 3}},
			ok:      true,
		},
	}

	for i, tt := range tests {
		classes, ok, err := tt.options.UserClass()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.UserClass: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := len(tt.classes), len(classes); want != got {
			t.Fatalf("[%02d] test %q, unexpected classes slice length: %v != %v",
				i, tt.desc, want, got)

		}

		for j := range classes {
			if want, got := tt.classes[j], classes[j]; !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.UserClass()\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.UserClass(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsVendorClass verifies that Options.VendorClass properly parses
// and returns raw vendor class data, if it is available with OptionVendorClass.
func TestOptionsVendorClass(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		classes [][]byte
		ok      bool
		err     error
	}{
		{
			desc: "OptionVendorClass not present in Options map",
		},
		{
			desc: "OptionVendorClass present in Options map, but empty",
			options: Options{
				OptionVendorClass: [][]byte{{}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionVendorClass present in Options map, zero item",
			options: Options{
				OptionVendorClass: [][]byte{{
					0, 0, 5, 0x58,
				}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionVendorClass present in Options map, one item, zero length",
			options: Options{
				OptionVendorClass: [][]byte{{
					0, 0, 5, 0x58,
					0, 0,
				}},
			},
			classes: [][]byte{{}},
			ok:      true,
		},
		{
			desc: "OptionVendorClass present in Options map, one item, extra byte",
			options: Options{
				OptionVendorClass: [][]byte{{
					0, 0, 5, 0x58,
					0, 1, 1, 255,
				}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionVendorClass present in Options map, one item",
			options: Options{
				OptionVendorClass: [][]byte{{
					0, 0, 5, 0x58,
					0, 1, 1,
				}},
			},
			classes: [][]byte{{1}},
			ok:      true,
		},
		{
			desc: "OptionVendorClass present in Options map, three items",
			options: Options{
				OptionVendorClass: [][]byte{{
					0, 0, 5, 0x58,
					0, 1, 1,
					0, 2, 2, 2,
					0, 3, 3, 3, 3,
				}},
			},
			classes: [][]byte{{1}, {2, 2}, {3, 3, 3}},
			ok:      true,
		},
	}

	for i, tt := range tests {
		classes, ok, err := tt.options.VendorClass()

		if want, got := tt.err, err; want != got {
			t.Fatalf("[%02d] test %q, unexpected error for Options.VendorClass: %v != %v",
				i, tt.desc, want, got)
		}

		if err != nil {
			continue
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.VendorClass(): %v != %v",
				i, tt.desc, want, got)
		}

		if !ok {
			continue
		}

		if want, got := len(tt.classes), len(classes.VendorClassData); want != got {
			t.Fatalf("[%02d] test %q, unexpected classes slice length: %v != %v",
				i, tt.desc, want, got)

		}

		for j := range classes.VendorClassData {
			if want, got := tt.classes[j], classes.VendorClassData[j]; !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.VendorClass()\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}
	}
}

// TestInterfaceID verifies that Options.InterfaceID properly parses
// and returns raw interface-id data, if it is available with InterfaceID.
func TestInterfaceID(t *testing.T) {
	var tests = []struct {
		desc        string
		options     Options
		interfaceID InterfaceID
		ok          bool
		err         error
	}{
		{
			desc: "InterfaceID not present in Options map",
		},
		{
			desc: "InterfaceID present in Options map, one item",
			options: Options{
				OptionInterfaceID: [][]byte{{
					0, 1, 1,
				}},
			},
			interfaceID: []byte{0, 1, 1},
			ok:          true,
		},
		{
			desc: "InterfaceID present in Options map with no interface-id data",
			options: Options{
				OptionInterfaceID: [][]byte{{}},
			},
			interfaceID: []byte{},
			ok:          true,
		},
	}

	for i, tt := range tests {
		interfaceID, ok, err := tt.options.InterfaceID()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.InterfaceID: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.interfaceID, interfaceID; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.InterfaceID()\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.InterfaceID(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsIAPD verifies that Options.IAPD properly parses and
// returns multiple IAPD values, if one or more are available with OptionIAPD.
func TestOptionsIAPD(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		iapd    []*IAPD
		ok      bool
		err     error
	}{
		{
			desc: "OptionIAPD not present in Options map",
		},
		{
			desc: "OptionIAPD present in Options map, but too short",
			options: Options{
				OptionIAPD: [][]byte{bytes.Repeat([]byte{0}, 11)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "one OptionIAPD present in Options map",
			options: Options{
				OptionIAPD: [][]byte{{
					1, 2, 3, 4,
					0, 0, 0, 30,
					0, 0, 0, 60,
				}},
			},
			iapd: []*IAPD{
				{
					IAID: [4]byte{1, 2, 3, 4},
					T1:   30 * time.Second,
					T2:   60 * time.Second,
				},
			},
			ok: true,
		},
		{
			desc: "two OptionIAPD present in Options map",
			options: Options{
				OptionIAPD: [][]byte{
					append(bytes.Repeat([]byte{0}, 12), []byte{0, 1, 0, 1, 1}...),
					append(bytes.Repeat([]byte{0}, 12), []byte{0, 2, 0, 1, 2}...),
				},
			},
			iapd: []*IAPD{
				{
					Options: Options{
						OptionClientID: [][]byte{{1}},
					},
				},
				{
					Options: Options{
						OptionServerID: [][]byte{{2}},
					},
				},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		iapd, ok, err := tt.options.IAPD()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.IAPD: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		for j := range tt.iapd {
			want, err := tt.iapd[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			got, err := iapd[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.IAPD():\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.IAPD(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsIAPrefix verifies that Options.IAPrefix properly parses and
// returns multiple IAPrefix values, if one or more are available with
// OptionIAPrefix.
func TestOptionsIAPrefix(t *testing.T) {
	var tests = []struct {
		desc     string
		options  Options
		iaprefix []*IAPrefix
		ok       bool
		err      error
	}{
		{
			desc: "OptionIAPrefix not present in Options map",
		},
		{
			desc: "OptionIAPrefix present in Options map, but too short",
			options: Options{
				OptionIAPrefix: [][]byte{bytes.Repeat([]byte{0}, 24)},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "one OptionIAPrefix present in Options map",
			options: Options{
				OptionIAPrefix: [][]byte{{
					0, 0, 0, 30,
					0, 0, 0, 60,
					32,
					32, 1, 13, 184, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				}},
			},
			iaprefix: []*IAPrefix{
				{
					PreferredLifetime: 30 * time.Second,
					ValidLifetime:     60 * time.Second,
					PrefixLength:      32,
					Prefix: net.IP{
						32, 1, 13, 184, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0,
					},
				},
			},
			ok: true,
		},
		{
			desc: "two OptionIAPrefix present in Options map",
			options: Options{
				OptionIAPrefix: [][]byte{
					bytes.Repeat([]byte{0}, 25),
					bytes.Repeat([]byte{0}, 25),
				},
			},
			iaprefix: []*IAPrefix{
				{
					Prefix: net.IPv6zero,
				},
				{
					Prefix: net.IPv6zero,
				},
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		iaprefix, ok, err := tt.options.IAPrefix()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.IAPrefix: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		for j := range tt.iaprefix {
			want, err := tt.iaprefix[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			got, err := iaprefix[j].MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.IAPrefix():\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.IAPrefix(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestRemoteIdentifier verifies that Options.RemoteIdentifier properly parses
// and returns a RemoteIdentifier, if it is available with OptionsRemoteIdentifier.
func TestOptionsRemoteIdentifier(t *testing.T) {
	var tests = []struct {
		desc             string
		options          Options
		remoteIdentifier *RemoteIdentifier
		ok               bool
		err              error
	}{
		{
			desc: "OptionsRemoteIdentifier not present in Options map",
		},
		{
			desc: "OptionsRemoteIdentifier present in Options map, but too short",
			options: Options{
				OptionRemoteIdentifier: [][]byte{{
					0, 0, 5, 0x58,
				}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionsRemoteIdentifier present in Options map",
			options: Options{
				OptionRemoteIdentifier: [][]byte{{
					0, 0, 5, 0x58,
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xe, 0xf,
				}},
			},
			remoteIdentifier: &RemoteIdentifier{
				EnterpriseNumber: 1368,
				RemoteId:         []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xe, 0xf},
			},
			ok: true,
		},
	}
	for i, tt := range tests {
		remoteIdentifier, ok, err := tt.options.RemoteIdentifier()
		if want, got := tt.err, err; want != got {
			t.Fatalf("[%02d] test %q, unexpected error for Options.RemoteIdentifier\n- want: %v\n-  got: %v", i, tt.desc, want, got)
		}

		if tt.err != nil {
			continue
		}

		if want, got := tt.remoteIdentifier, remoteIdentifier; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.RemoteIdentifier()\n- want: %v\n-  got: %v", i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.RemoteIdentifier(): %v != %v", i, tt.desc, want, got)
		}
	}
}

// TestOptionsBootFileURL verifies that Options.BootFileURL properly parses
// and returns a URL, if it is available with OptionBootFileURL.
func TestOptionsBootFileURL(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		u       *url.URL
		ok      bool
	}{
		{
			desc: "OptionBootFileURL not present in Options map",
		},
		{
			desc: "OptionBootFileURL present in Options map",
			options: Options{
				OptionBootFileURL: [][]byte{[]byte("tftp://192.168.1.1:69")},
			},
			u: &url.URL{
				Scheme: "tftp",
				Host:   "192.168.1.1:69",
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		u, ok, err := tt.options.BootFileURL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.BootFileURL(): %v != %v",
				i, tt.desc, want, got)
		}
		if !ok {
			continue
		}

		ttuu := url.URL(*tt.u)
		uu := url.URL(*u)
		if want, got := ttuu.String(), uu.String(); want != got {
			t.Fatalf("[%02d] test %q, unexpected value for Options.BootFileURL(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsBootFileParam verifies that Options.BootFileParam properly parses
// and returns boot file parameter data, if it is available with
// OptionBootFileParam.
func TestOptionsBootFileParam(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		param   Data
		ok      bool
		err     error
	}{
		{
			desc: "OptionBootFileParam not present in Options map",
		},
		{
			desc: "OptionBootFileParam present in Options map, but empty",
			options: Options{
				OptionBootFileParam: [][]byte{{}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionBootFileParam present in Options map, one item, zero length",
			options: Options{
				OptionBootFileParam: [][]byte{{
					0, 0,
				}},
			},
			param: Data{{}},
			ok:    true,
		},
		{
			desc: "OptionBootFileParam present in Options map, one item, extra byte",
			options: Options{
				OptionBootFileParam: [][]byte{{
					0, 1, 1, 255,
				}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionBootFileParam present in Options map, one item",
			options: Options{
				OptionBootFileParam: [][]byte{{
					0, 3, 'f', 'o', 'o',
				}},
			},
			param: Data{[]byte("foo")},
			ok:    true,
		},
		{
			desc: "OptionBootFileParam present in Options map, three items",
			options: Options{
				OptionBootFileParam: [][]byte{{
					0, 1, 'a',
					0, 2, 'a', 'b',
					0, 3, 'a', 'b', 'c',
				}},
			},
			param: Data{[]byte("a"), []byte("ab"), []byte("abc")},
			ok:    true,
		},
	}

	for i, tt := range tests {
		param, ok, err := tt.options.BootFileParam()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.BootFileParam: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := len(tt.param), len(param); want != got {
			t.Fatalf("[%02d] test %q, unexpected param slice length: %v != %v",
				i, tt.desc, want, got)

		}

		for j := range param {
			if want, got := tt.param[j], param[j]; !bytes.Equal(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.BootFileParam()\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.BootFileParam(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsClientArchType verifies that Options.ClientArchType properly parses
// and returns client architecture type data, if it is available with
// OptionClientArchType.
func TestOptionsClientArchType(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		arch    ArchTypes
		ok      bool
		err     error
	}{
		{
			desc: "OptionClientArchType not present in Options map",
		},
		{
			desc: "OptionClientArchType present in Options map, but empty",
			options: Options{
				OptionClientArchType: [][]byte{{}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionClientArchType present in Options map, but not divisible by 2",
			options: Options{
				OptionClientArchType: [][]byte{{0, 0, 0}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionClientArchType present in Options map, one architecture",
			options: Options{
				OptionClientArchType: [][]byte{{0, 9}},
			},
			arch: ArchTypes{ArchTypeEFIx8664},
			ok:   true,
		},
		{
			desc: "OptionClientArchType present in Options map, three architectures",
			options: Options{
				OptionClientArchType: [][]byte{{0, 5, 0, 9, 0, 0}},
			},
			arch: ArchTypes{
				ArchTypeIntelLeanClient,
				ArchTypeEFIx8664,
				ArchTypeIntelx86PC,
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		arch, ok, err := tt.options.ClientArchType()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.ClientArchType: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := len(tt.arch), len(arch); want != got {
			t.Fatalf("[%02d] test %q, unexpected arch slice length: %v != %v",
				i, tt.desc, want, got)
		}

		for j := range arch {
			if want, got := tt.arch[j], arch[j]; !reflect.DeepEqual(want, got) {
				t.Fatalf("[%02d:%02d] test %q, unexpected value for Options.ClientArchType()\n- want: %v\n-  got: %v",
					i, j, tt.desc, want, got)
			}
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.ClientArchType(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptionsNII verifies that Options.NII properly parses and returns a
// Network Interface Identifier value, if it is available with OptionNII.
func TestOptionsNII(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		nii     *NII
		ok      bool
		err     error
	}{
		{
			desc: "OptionNII not present in Options map",
		},
		{
			desc: "OptionNII present in Options map, but too short length",
			options: Options{
				OptionNII: [][]byte{{1, 2}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionNII present in Options map, but too long length",
			options: Options{
				OptionNII: [][]byte{{1, 2, 3, 4}},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			desc: "OptionNII present in Options map",
			options: Options{
				OptionNII: [][]byte{{1, 2, 3}},
			},
			nii: &NII{
				Type:  1,
				Major: 2,
				Minor: 3,
			},
			ok: true,
		},
	}

	for i, tt := range tests {
		nii, ok, err := tt.options.NII()
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for Options.NII(): %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.nii, nii; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected value for Options.NII(): %v != %v",
				i, tt.desc, want, got)
		}

		if want, got := tt.ok, ok; want != got {
			t.Fatalf("[%02d] test %q, unexpected ok for Options.NII(): %v != %v",
				i, tt.desc, want, got)
		}
	}
}

// TestOptions_enumerate verifies that Options.enumerate correctly enumerates
// and sorts an Options map into key/value option pairs.
func TestOptions_enumerate(t *testing.T) {
	var tests = []struct {
		desc    string
		options Options
		kv      optslice
	}{
		{
			desc: "one key/value pair",
			options: Options{
				1: [][]byte{[]byte("foo")},
			},
			kv: optslice{
				option{
					Code: 1,
					Data: []byte("foo"),
				},
			},
		},
		{
			desc: "two key/value pairs",
			options: Options{
				1: [][]byte{[]byte("foo")},
				2: [][]byte{[]byte("bar")},
			},
			kv: optslice{
				option{
					Code: 1,
					Data: []byte("foo"),
				},
				option{
					Code: 2,
					Data: []byte("bar"),
				},
			},
		},
		{
			desc: "four key/value pairs, two with same key",
			options: Options{
				1: [][]byte{[]byte("foo"), []byte("baz")},
				3: [][]byte{[]byte("qux")},
				2: [][]byte{[]byte("bar")},
			},
			kv: optslice{
				option{
					Code: 1,
					Data: []byte("foo"),
				},
				option{
					Code: 1,
					Data: []byte("baz"),
				},
				option{
					Code: 2,
					Data: []byte("bar"),
				},
				option{
					Code: 3,
					Data: []byte("qux"),
				},
			},
		},
	}

	for i, tt := range tests {
		if want, got := tt.kv, tt.options.enumerate(); !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected key/value options:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

// Test_parseOptions verifies that parseOptions parses correct option values
// from a slice of bytes, and that it returns an empty Options map if the byte
// slice cannot contain options.
func Test_parseOptions(t *testing.T) {
	var tests = []struct {
		desc    string
		buf     []byte
		options Options
		err     error
	}{
		{
			desc:    "nil options bytes",
			options: Options{},
		},
		{
			desc:    "empty options bytes",
			buf:     []byte{},
			options: Options{},
		},
		{
			desc: "too short options bytes",
			buf:  []byte{0},
			err:  errInvalidOptions,
		},
		{
			desc:    "zero code, zero length option bytes",
			buf:     []byte{0, 0, 0, 0},
			options: Options{},
		},
		{
			desc: "zero code, zero length option bytes with trailing byte",
			buf:  []byte{0, 0, 0, 0, 1},
			err:  errInvalidOptions,
		},
		{
			desc: "zero code, length 3, incorrect length for data",
			buf:  []byte{0, 0, 0, 3, 1, 2},
			err:  errInvalidOptions,
		},
		{
			desc: "client ID, length 1, value [1]",
			buf:  []byte{0, 1, 0, 1, 1},
			options: Options{
				OptionClientID: [][]byte{{1}},
			},
		},
		{
			desc: "client ID, length 2, value [1 1] + server ID, length 3, value [1 2 3]",
			buf: []byte{
				0, 1, 0, 2, 1, 1,
				0, 2, 0, 3, 1, 2, 3,
			},
			options: Options{
				OptionClientID: [][]byte{{1, 1}},
				OptionServerID: [][]byte{{1, 2, 3}},
			},
		},
	}

	for i, tt := range tests {
		options, err := parseOptions(tt.buf)
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for parseOptions(%v): %v != %v",
					i, tt.desc, tt.buf, want, got)
			}

			continue
		}
		if want, got := tt.options, options; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Options map for parseOptions(%v):\n- want: %v\n-  got: %v",
				i, tt.desc, tt.buf, want, got)
		}

		for k, v := range tt.options {
			for ii := range v {
				if want, got := cap(v[ii]), cap(options[k][ii]); want != got {
					t.Fatalf("[%02d] test %q, unexpected capacity option data:\n- want: %v\n-  got: %v",
						i, tt.desc, want, got)
				}
			}
		}
	}
}
