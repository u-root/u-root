package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaVLanUnspec = iota
	tcaVLanTm
	tcaVLanParms
	tcaVLanPushVLanID
	tcaVLanPushVLanProtocol
	tcaVLanPad
	tcaVLanPushVLanPriority
)

// VLan contains attribute of the VLan discipline
type VLan struct {
	Parms        *VLanParms
	Tm           *Tcft
	PushID       *uint16
	PushProtocol *uint16
	PushPriority *uint32
}

// VLanParms from include/uapi/linux/tc_act/tc_vlan.h
type VLanParms struct {
	Index      uint32
	Capab      uint32
	Action     uint32
	RefCnt     uint32
	BindCnt    uint32
	VLanAction uint32
}

// marshalVLan returns the binary encoding of Vlan
func marshalVlan(info *VLan) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("VLan: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaVLanParms, Data: data})
	}
	if info.PushID != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaVLanPushVLanID, Data: *info.PushID})
	}
	if info.PushProtocol != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaVLanPushVLanProtocol, Data: *info.PushProtocol})
	}
	if info.PushPriority != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaVLanPushVLanPriority, Data: *info.PushPriority})
	}
	return marshalAttributes(options)
}

// unmarshalVLan parses the VLan-encoded data and stores the result in the value pointed to by info.
func unmarshalVLan(data []byte, info *VLan) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaVLanParms:
			parms := &VLanParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaVLanTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaVLanPushVLanID:
			tmp := ad.Uint16()
			info.PushID = &tmp
		case tcaVLanPushVLanProtocol:
			tmp := ad.Uint16()
			info.PushProtocol = &tmp
		case tcaVLanPushVLanPriority:
			tmp := ad.Uint32()
			info.PushPriority = &tmp
		case tcaVLanPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalVLan()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
