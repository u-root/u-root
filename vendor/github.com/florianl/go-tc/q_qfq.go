package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaQfqUnspec = iota
	tcaQfqWeight
	tcaQfqLmax
)

// Qfq contains attributes of the qfq discipline
type Qfq struct {
	Weight *uint32
	Lmax   *uint32
}

// unmarshalQfq parses the Qfq-encoded data and stores the result in the value pointed to by info.
func unmarshalQfq(data []byte, info *Qfq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaQfqWeight:
			info.Weight = uint32Ptr(ad.Uint32())
		case tcaQfqLmax:
			info.Lmax = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("UnmarshalQfq()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalQfq returns the binary encoding of Qfq
func marshalQfq(info *Qfq) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Qfq: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Weight != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaQfqWeight, Data: uint32Value(info.Weight)})
	}
	if info.Lmax != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaQfqLmax, Data: uint32Value(info.Lmax)})
	}
	return marshalAttributes(options)
}
