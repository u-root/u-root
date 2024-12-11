package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaActBpfUnspec = iota
	tcaActBpfTm
	tcaActBpfParms
	tcaActBpfOpsLen
	tcaActBpfOps
	tcaActBpfFD
	tcaActBpfName
	tcaActBpfPad
	tcaActBpfTag
	tcaActBpfID
)

// ActBpf represents policing attributes of various filters and classes
type ActBpf struct {
	Tm     *Tcft
	Parms  *ActBpfParms
	Ops    *[]byte
	OpsLen *uint16
	FD     *uint32
	Name   *string
	Tag    *[]byte
	ID     *uint32
}

// unmarshalActBpf parses the ActBpf-encoded data and stores the result in the value pointed to by info.
func unmarshalActBpf(data []byte, info *ActBpf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaActBpfTm:
			tm := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tm)
			multiError = concatError(multiError, err)
			info.Tm = tm
		case tcaActBpfParms:
			parms := &ActBpfParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaActBpfOpsLen:
			info.OpsLen = uint16Ptr(ad.Uint16())
		case tcaActBpfOps:
			info.Ops = bytesPtr(ad.Bytes())
		case tcaActBpfFD:
			info.FD = uint32Ptr(ad.Uint32())
		case tcaActBpfName:
			info.Name = stringPtr(ad.String())
		case tcaActBpfTag:
			info.Tag = bytesPtr(ad.Bytes())
		case tcaActBpfID:
			info.ID = uint32Ptr(ad.Uint32())
		case tcaActBpfPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalActBpf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalActBpf returns the binary encoding of ActBpf
func marshalActBpf(info *ActBpf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("ActBpf: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Name != nil {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaActBpfName, Data: stringValue(info.Name)})
	}
	if info.Tag != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActBpfTag, Data: bytesValue(info.Tag)})
	}
	if info.FD != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaActBpfFD, Data: uint32Value(info.FD)})
	}
	if info.ID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaActBpfID, Data: uint32Value(info.ID)})
	}
	if info.Ops != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActBpfOps, Data: bytesValue(info.Ops)})
	}
	if info.OpsLen != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaActBpfOpsLen, Data: uint16Value(info.OpsLen)})
	}
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActBpfParms, Data: data})
	}
	return marshalAttributes(options)
}

// ActBpfParms from include/uapi/linux/tc_act/tc_bpf.h
type ActBpfParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	Refcnt  uint32
	Bindcnt uint32
}
