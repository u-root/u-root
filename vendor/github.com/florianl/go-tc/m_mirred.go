package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaMirredUnspec = iota
	tcaMirredTm
	tcaMirredParms
	tcaMirredPad
)

// Mirred represents policing attributes of various filters and classes
type Mirred struct {
	Parms *MirredParam
	Tm    *Tcft
}

// MirredParam from include/uapi/linux/tc_act/tc_mirred.h
type MirredParam struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
	Eaction uint32
	IfIndex uint32
}

func unmarshalMirred(data []byte, info *Mirred) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaMirredParms:
			param := &MirredParam{}
			err = unmarshalStruct(ad.Bytes(), param)
			multiError = concatError(multiError, err)
			info.Parms = param
		case tcaMirredTm:
			tm := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tm)
			multiError = concatError(multiError, err)
			info.Tm = tm
		case tcaMirredPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalMirred()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalMirred returns the binary encoding of Mirred
func marshalMirred(info *Mirred) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Mirred: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaMirredParms, Data: data})
	}
	return marshalAttributes(options)
}
