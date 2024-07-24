package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCbsUnspec = iota
	tcaCbsParms
)

// CbsOpt contains attributes of the cbs discipline
type CbsOpt struct {
	Offload   uint8
	Pad       [3]uint8
	HiCredit  int32
	LoCredit  int32
	IdleSlope int32
	SendSlope int32
}

// Cbs contains attributes of the cbs discipline
type Cbs struct {
	Parms *CbsOpt
}

// unmarshalCbs parses the Cbs-encoded data and stores the result in the value pointed to by info.
func unmarshalCbs(data []byte, info *Cbs) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaCbsParms:
			opt := &CbsOpt{}
			err := unmarshalStruct(ad.Bytes(), opt)
			multiError = concatError(multiError, err)
			info.Parms = opt
		default:
			return fmt.Errorf("unmarshalCbs()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalCbs returns the binary encoding of Qbs
func marshalCbs(info *Cbs) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Cbs: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	data, err := marshalStruct(info.Parms)
	if err != nil {
		return []byte{}, err
	}
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCbsParms, Data: data})

	return marshalAttributes(options)

}
