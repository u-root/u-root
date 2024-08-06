package tc

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

const (
	tcaIfeUnspec = iota
	tcaIfeParms
	tcaIfeTm
	tcaIfeDMac
	tcaIfeSMac
	tcaIfeType
	tcaIfeMetaList
	tcaIfePad
)

// Ife contains attribute of the ife discipline
type Ife struct {
	Parms *IfeParms
	SMac  *net.HardwareAddr
	DMac  *net.HardwareAddr
	Type  *uint16
	Tm    *Tcft
}

// IfeParms from include/uapi/linux/tc_act/tc_ife.h
type IfeParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
	Flags   uint16
}

// marshalIfe returns the binary encoding of Ife
func marshalIfe(info *Ife) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ife: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIfeParms, Data: data})
	}
	if info.SMac != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIfeSMac, Data: []byte(*info.SMac)})
	}
	if info.DMac != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIfeDMac, Data: []byte(*info.DMac)})
	}
	if info.Type != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaIfeType, Data: *info.Type})
	}
	return marshalAttributes(options)
}

// unmarshalIfe parses the ife-encoded data and stores the result in the value pointed to by info.
func unmarshalIfe(data []byte, info *Ife) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaIfeParms:
			parms := &IfeParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaIfeSMac:
			tmp := net.HardwareAddr(ad.Bytes())
			info.SMac = &tmp
		case tcaIfeDMac:
			tmp := net.HardwareAddr(ad.Bytes())
			info.DMac = &tmp
		case tcaIfeTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaIfeType:
			tmp := ad.Uint16()
			info.Type = &tmp
		case tcaIfePad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalIfe()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
