package dhcpv4

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/uio"
)

func TestParseOption(t *testing.T) {
	for _, tt := range []struct {
		code  OptionCode
		value []byte
		want  string
	}{
		{
			code:  OptionNameServer,
			value: []byte{192, 168, 1, 254},
			want:  "[192 168 1 254]",
		},
		{
			code:  OptionSubnetMask,
			value: []byte{255, 255, 255, 0},
			want:  "ffffff00",
		},
		{
			code:  OptionRouter,
			value: []byte{192, 168, 1, 1, 192, 168, 2, 1},
			want:  "192.168.1.1, 192.168.2.1",
		},
		{
			code:  OptionDomainNameServer,
			value: []byte{192, 168, 1, 1, 192, 168, 2, 1},
			want:  "192.168.1.1, 192.168.2.1",
		},
		{
			code:  OptionNTPServers,
			value: []byte{192, 168, 1, 1, 192, 168, 2, 1},
			want:  "192.168.1.1, 192.168.2.1",
		},
		{
			code:  OptionServerIdentifier,
			value: []byte{192, 168, 1, 1, 192, 168, 2, 1},
			want:  "192.168.1.1, 192.168.2.1",
		},
		{
			code:  OptionHostName,
			value: []byte("test"),
			want:  "test",
		},
		{
			code:  OptionDomainName,
			value: []byte("test"),
			want:  "test",
		},
		{
			code:  OptionRootPath,
			value: []byte("test"),
			want:  "test",
		},
		{
			code:  OptionClassIdentifier,
			value: []byte("test"),
			want:  "test",
		},
		{
			code:  OptionTFTPServerName,
			value: []byte("test"),
			want:  "test",
		},
		{
			code:  OptionBootfileName,
			value: []byte("test"),
			want:  "test",
		},
		{
			code:  OptionBroadcastAddress,
			value: []byte{192, 168, 1, 1},
			want:  "192.168.1.1",
		},
		{
			code:  OptionRequestedIPAddress,
			value: []byte{192, 168, 1, 1},
			want:  "192.168.1.1",
		},
		{
			code:  OptionIPAddressLeaseTime,
			value: []byte{0, 0, 0, 12},
			want:  "12s",
		},
		{
			code:  OptionDHCPMessageType,
			value: []byte{1},
			want:  "DISCOVER",
		},
		{
			code:  OptionParameterRequestList,
			value: []byte{3, 4, 5},
			want:  "Router, Time Server, Name Server",
		},
		{
			code:  OptionMaximumDHCPMessageSize,
			value: []byte{1, 2},
			want:  "258",
		},
		{
			code:  OptionUserClassInformation,
			value: []byte{4, 't', 'e', 's', 't', 3, 'f', 'o', 'o'},
			want:  "test, foo",
		},
		{
			code:  OptionRelayAgentInformation,
			value: []byte{1, 4, 129, 168, 0, 1},
			want:  "    unknown (1): [129 168 0 1]\n",
		},
		{
			code:  OptionClientSystemArchitectureType,
			value: []byte{0, 0},
			want:  "Intel x86PC",
		},
	} {
		s := parseOption(tt.code, tt.value)
		if got := s.String(); got != tt.want {
			t.Errorf("parseOption(%s, %v) = %s, want %s", tt.code, tt.value, got, tt.want)
		}
	}
}

func TestOptionToBytes(t *testing.T) {
	o := Option{
		Code:  OptionDHCPMessageType,
		Value: &OptionGeneric{[]byte{byte(MessageTypeDiscover)}},
	}
	serialized := o.Value.ToBytes()
	expected := []byte{1}
	require.Equal(t, expected, serialized)
}

func TestOptionString(t *testing.T) {
	o := Option{
		Code:  OptionDHCPMessageType,
		Value: MessageTypeDiscover,
	}
	require.Equal(t, "DHCP Message Type: DISCOVER", o.String())
}

func TestOptionStringUnknown(t *testing.T) {
	o := Option{
		Code:  GenericOptionCode(102), // Returend option code.
		Value: &OptionGeneric{[]byte{byte(MessageTypeDiscover)}},
	}
	require.Equal(t, "unknown (102): [1]", o.String())
}

func TestOptionsMarshal(t *testing.T) {
	for i, tt := range []struct {
		opts Options
		want []byte
	}{
		{
			opts: nil,
			want: nil,
		},
		{
			opts: Options{
				5: []byte{1, 2, 3, 4},
			},
			want: []byte{
				5 /* key */, 4 /* length */, 1, 2, 3, 4,
			},
		},
		{
			// Test sorted key order.
			opts: Options{
				5:   []byte{1, 2, 3},
				100: []byte{101, 102, 103},
				255: []byte{},
			},
			want: []byte{
				5, 3, 1, 2, 3,
				100, 3, 101, 102, 103,
			},
		},
		{
			// Test RFC 3396.
			opts: Options{
				5: bytes.Repeat([]byte{10}, math.MaxUint8+1),
			},
			want: append(append(
				[]byte{5, math.MaxUint8}, bytes.Repeat([]byte{10}, math.MaxUint8)...),
				5, 1, 10,
			),
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			require.Equal(t, uio.ToBigEndian(tt.opts), tt.want)
		})
	}
}

func TestOptionsUnmarshal(t *testing.T) {
	for i, tt := range []struct {
		input     []byte
		want      Options
		wantError bool
	}{
		{
			// Buffer missing data.
			input: []byte{
				3 /* key */, 3 /* length */, 1,
			},
			wantError: true,
		},
		{
			input: []byte{
				// This may look too long, but 0 is padding.
				// The issue here is the missing OptionEnd.
				3, 3, 0, 0, 0, 0, 0, 0, 0,
			},
			wantError: true,
		},
		{
			// Only OptionPad and OptionEnd can stand on their own
			// without a length field. So this is too short.
			input: []byte{
				3,
			},
			wantError: true,
		},
		{
			// Option present after the End is a nono.
			input:     []byte{byte(OptionEnd), 3},
			wantError: true,
		},
		{
			input: []byte{byte(OptionEnd)},
			want:  Options{},
		},
		{
			input: []byte{
				3, 2, 5, 6,
				byte(OptionEnd),
			},
			want: Options{
				3: []byte{5, 6},
			},
		},
		{
			// Test RFC 3396.
			input: append(
				append([]byte{3, math.MaxUint8}, bytes.Repeat([]byte{10}, math.MaxUint8)...),
				3, 5, 10, 10, 10, 10, 10,
				byte(OptionEnd),
			),
			want: Options{
				3: bytes.Repeat([]byte{10}, math.MaxUint8+5),
			},
		},
		{
			input: []byte{
				10, 2, 255, 254,
				11, 3, 5, 5, 5,
				byte(OptionEnd),
			},
			want: Options{
				10: []byte{255, 254},
				11: []byte{5, 5, 5},
			},
		},
		{
			input: append(
				append([]byte{10, 2, 255, 254}, bytes.Repeat([]byte{byte(OptionPad)}, 255)...),
				byte(OptionEnd),
			),
			want: Options{
				10: []byte{255, 254},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			opt := make(Options)
			err := opt.fromBytesCheckEnd(tt.input, true)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, opt, tt.want)
			}
		})
	}
}
