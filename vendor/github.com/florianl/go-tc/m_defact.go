package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaDefUnspec = iota
	tcaDefTm
	tcaDefParms
	tcaDefData
	tcaDefPad
)

// Defact contains attributes of the defact discipline
type Defact struct {
	Parms *DefactParms
	Tm    *Tcft
	Data  *string
}

// DefactParms from include/uapi/linux/tc_act/tc_defact.h
type DefactParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
}

// marshalDefact returns the binary encoding of Defact
func marshalDefact(info *Defact) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Defact: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaDefParms, Data: data})
	}
	if info.Data != nil {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaDefData, Data: *info.Data})
	}
	return marshalAttributes(options)
}

// unmarshalDefact parses the defact-encoded data and stores the result in the value pointed to by info.
func unmarshalDefact(data []byte, info *Defact) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaDefParms:
			parms := &DefactParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaDefTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaDefData:
			tmp := ad.String()
			info.Data = &tmp
		case tcaDefPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalDefact()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
