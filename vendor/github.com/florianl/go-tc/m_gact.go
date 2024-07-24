package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaGactUnspec = iota
	tcaGactTm
	tcaGactParm
	tcaGactProb
	tcaGactPad
)

// Gact contains attributes of the gact discipline
type Gact struct {
	Tm    *Tcft
	Parms *GactParms
	Prob  *GactProb
}

// GactProb from include/uapi/linux/tc_act/tc_gact.h
type GactProb struct {
	PType   uint16
	PVal    uint16
	PAction uint32
}

// GactParms from include/uapi/linux/tc_act/tc_gact.h
type GactParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
}

// marshalGact returns the binary encoding of Gact
func marshalGact(info *Gact) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Gact: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	var multiError error

	if info.Prob != nil {
		data, err := marshalStruct(info.Prob)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaGactProb, Data: data})
	}
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaGactParm, Data: data})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}

// unmarshalGact parses the gact-encoded data and stores the result in the value pointed to by info.
func unmarshalGact(data []byte, info *Gact) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaGactTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaGactParm:
			parms := &GactParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaGactProb:
			prob := &GactProb{}
			err = unmarshalStruct(ad.Bytes(), prob)
			multiError = concatError(multiError, err)
			info.Prob = prob
		case tcaGactPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalGact()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
