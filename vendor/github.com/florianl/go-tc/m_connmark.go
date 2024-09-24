package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaConnmarkUnspec = iota
	tcaConnmarkParms
	tcaConnmarkTm
	tcaConnmarkPad
)

// Connmark represents policing attributes of various filters and classes
type Connmark struct {
	Parms *ConnmarkParam
	Tm    *Tcft
}

// ConnmarkParam from include/uapi/linux/tc_act/tc_connmark.h
type ConnmarkParam struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
	Zone    uint16
}

func unmarshalConnmark(data []byte, info *Connmark) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaConnmarkParms:
			param := &ConnmarkParam{}
			err = unmarshalStruct(ad.Bytes(), param)
			multiError = concatError(multiError, err)
			info.Parms = param
		case tcaConnmarkTm:
			tm := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tm)
			multiError = concatError(multiError, err)
			info.Tm = tm
		case tcaConnmarkPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalConnmark()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalConnmark returns the binary encoding of ActBpf
func marshalConnmark(info *Connmark) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Connmark: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Parms != nil {
		data, err := marshalAndAlignStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaConnmarkParms, Data: data})
	}
	return marshalAttributes(options)
}
