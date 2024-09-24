package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaNatUnspec = iota
	tcaNatParms
	tcaNatTm
	tcaNatPad
)

// Nat contains attribute of the nat discipline
type Nat struct {
	Parms *NatParms
	Tm    *Tcft
}

// NatParms from include/uapi/linux/tc_act/tc_nat.h
type NatParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
	OldAddr uint32
	NewAddr uint32
	Mask    uint32
	Flags   uint32
}

// marshalNat returns the binary encoding of Ife
func marshalNat(info *Nat) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Nat: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaNatParms, Data: data})
	}
	return marshalAttributes(options)
}

// unmarshalNat parses the nat-encoded data and stores the result in the value pointed to by info.
func unmarshalNat(data []byte, info *Nat) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaNatParms:
			parms := &NatParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaNatTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaNatPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalNat()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
