package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaSkbEditUnspec = iota
	tcaSkbEditTm
	tcaSkbEditParms
	tcaSkbEditPriority
	tcaSkbEditQueueMapping
	tcaSkbEditMark
	tcaSkbEditPad
	tcaSkbEditPtype
	tcaSkbEditMask
	tcaSkbEditFlags
	tcaSkbEditQueueMappingMax
)

// SkbEdit  contains attribute of the SkbEdit discipline
type SkbEdit struct {
	Tm              *Tcft
	Parms           *SkbEditParms
	Priority        *uint32
	QueueMapping    *uint16
	Mark            *uint32
	Ptype           *uint16
	Mask            *uint32
	Flags           *uint64
	QueueMappingMax *uint16
}

// SkbEditParms from include/uapi/linux/tc_act/tc_skbedit.h
type SkbEditParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
}

// unmarshalSkbEdit parses the skbedit-encoded data and stores the result in the value pointed to by info.
func unmarshalSkbEdit(data []byte, info *SkbEdit) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaSkbEditTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaSkbEditParms:
			parms := &SkbEditParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaSkbEditPriority:
			info.Priority = uint32Ptr(ad.Uint32())
		case tcaSkbEditQueueMapping:
			info.QueueMapping = uint16Ptr(ad.Uint16())
		case tcaSkbEditMark:
			info.Mark = uint32Ptr(ad.Uint32())
		case tcaSkbEditPtype:
			info.Ptype = uint16Ptr(ad.Uint16())
		case tcaSkbEditMask:
			info.Mask = uint32Ptr(ad.Uint32())
		case tcaSkbEditFlags:
			info.Flags = uint64Ptr(ad.Uint64())
		case tcaSkbEditQueueMappingMax:
			info.QueueMappingMax = uint16Ptr(ad.Uint16())
		case tcaSkbEditPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalSkbEdit()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalSkbEdit returns the binary encoding of SkbEdit
func marshalSkbEdit(info *SkbEdit) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("SkbEdit: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaSkbEditParms, Data: data})
	}
	if info.Priority != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaSkbEditPriority, Data: uint32Value(info.Priority)})
	}
	if info.QueueMapping != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaSkbEditQueueMapping, Data: uint16Value(info.QueueMapping)})
	}
	if info.Mark != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaSkbEditMark, Data: uint32Value(info.Mark)})
	}
	if info.Ptype != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaSkbEditPtype, Data: uint16Value(info.Ptype)})
	}
	if info.Mask != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaSkbEditMask, Data: uint32Value(info.Mask)})
	}
	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaSkbEditFlags, Data: uint64Value(info.Flags)})
	}
	if info.QueueMappingMax != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaSkbEditQueueMappingMax, Data: uint16Value(info.QueueMappingMax)})
	}
	return marshalAttributes(options)
}
