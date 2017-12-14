package dhcp4

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/u-root/dhcp4/util"
)

func TestOptionsMarshal(t *testing.T) {
	for i, tt := range []struct {
		opts Options
		want []byte
	}{
		{
			opts: nil,
			want: []byte{255},
		},
		{
			opts: Options{
				5: []byte{1, 2, 3, 4},
			},
			want: []byte{
				5 /* key */, 4 /* length */, 1, 2, 3, 4,
				255, /* end key */
			},
		},
		{
			// Test sorted key order.
			opts: Options{
				5:   []byte{1, 2, 3},
				100: []byte{101, 102, 103},
			},
			want: []byte{
				5, 3, 1, 2, 3,
				100, 3, 101, 102, 103,
				255,
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
				255,
			),
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			b := util.NewBuffer(nil)
			tt.opts.Marshal(b)
			if !bytes.Equal(b.Data(), tt.want) {
				t.Errorf("got %v want %v", b.Data(), tt.want)
			}
		})
	}
}

func TestOptionsUnmarshal(t *testing.T) {
	for i, tt := range []struct {
		input []byte
		want  Options
		err   error
	}{
		{
			input: nil,
			err:   io.ErrUnexpectedEOF,
		},
		{
			input: []byte{},
			err:   io.ErrUnexpectedEOF,
		},
		{
			input: []byte{
				3 /* key */, 3 /* length */, 1,
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			input: []byte{
				// This may look too long, but 0 is padding.
				// The issue here is the missing EOF.
				3, 3, 0, 0, 0, 0, 0, 0, 0,
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			input: []byte{
				3,
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			input: []byte{byte(End), 3},
			err:   ErrInvalidOptions,
		},
		{
			input: []byte{byte(End)},
			want:  Options{},
		},
		{
			input: []byte{
				3, 2, 5, 6,
				byte(End),
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
				byte(End),
			),
			want: Options{
				3: bytes.Repeat([]byte{10}, math.MaxUint8+5),
			},
		},
		{
			input: []byte{
				10, 2, 255, 254,
				11, 3, 5, 5, 5,
				byte(End),
			},
			want: Options{
				10: []byte{255, 254},
				11: []byte{5, 5, 5},
			},
		},
		{
			input: append(
				append([]byte{10, 2, 255, 254}, bytes.Repeat([]byte{byte(Pad)}, 255)...),
				byte(End),
			),
			want: Options{
				10: []byte{255, 254},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			var got Options
			if err := (&got).Unmarshal(util.NewBuffer(tt.input)); err != tt.err {
				t.Fatalf("got %v want %v", err, tt.err)
			} else if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}
