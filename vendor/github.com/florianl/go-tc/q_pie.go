package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaPieUnspec = iota
	tcaPieTarget
	tcaPieLimit
	tcaPieTUpdate
	tcaPieAlpha
	tcaPieBeta
	tcaPieECN
	tcaPieBytemode
)

// Pie contains attributes of the pie discipline
type Pie struct {
	Target   *uint32
	Limit    *uint32
	TUpdate  *uint32
	Alpha    *uint32
	Beta     *uint32
	ECN      *uint32
	Bytemode *uint32
}

// unmarshalPie parses the Pie-encoded data and stores the result in the value pointed to by info.
func unmarshalPie(data []byte, info *Pie) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaPieTarget:
			info.Target = uint32Ptr(ad.Uint32())
		case tcaPieLimit:
			info.Limit = uint32Ptr(ad.Uint32())
		case tcaPieTUpdate:
			info.TUpdate = uint32Ptr(ad.Uint32())
		case tcaPieAlpha:
			info.Alpha = uint32Ptr(ad.Uint32())
		case tcaPieBeta:
			info.Beta = uint32Ptr(ad.Uint32())
		case tcaPieECN:
			info.ECN = uint32Ptr(ad.Uint32())
		case tcaPieBytemode:
			info.Bytemode = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("extractPieOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalPie returns the binary encoding of Qfq
func marshalPie(info *Pie) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Pie: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Target != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieTarget, Data: uint32Value(info.Target)})
	}
	if info.Limit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieLimit, Data: uint32Value(info.Limit)})
	}
	if info.TUpdate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieTUpdate, Data: uint32Value(info.TUpdate)})
	}
	if info.Alpha != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieAlpha, Data: uint32Value(info.Alpha)})
	}
	if info.Beta != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieBeta, Data: uint32Value(info.Beta)})
	}
	if info.ECN != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieECN, Data: uint32Value(info.ECN)})
	}
	if info.Bytemode != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieBytemode, Data: uint32Value(info.Bytemode)})
	}
	return marshalAttributes(options)
}
