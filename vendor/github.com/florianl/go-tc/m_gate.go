package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaGateUnspec = iota
	tcaGateTm
	tcaGateParms
	tcaGatePad
	tcaGatePriority
	tcaGateEntryList
	tcaGateBaseTime
	tcaGateCycleTime
	tcaGateCycleTimeExt
	tcaGateFlags
	tcaGateClockID
)

// Gate contains attributes of the gate discipline
// https://man7.org/linux/man-pages/man8/tc-gate.8.html
type Gate struct {
	Tm           *Tcft
	Parms        *GateParms
	Priority     *int32
	BaseTime     *uint64
	CycleTime    *uint64
	CycleTimeExt *uint64
	Flags        *uint32
	ClockID      *int32
}

// marshalGate returns the binary encoding of Gate
func marshalGate(info *Gate) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Gate: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaGateParms, Data: data})
	}
	if info.Priority != nil {
		options = append(options, tcOption{Interpretation: vtInt32, Type: tcaGatePriority, Data: *info.Priority})
	}
	if info.BaseTime != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaGateBaseTime, Data: *info.BaseTime})
	}
	if info.CycleTime != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaGateCycleTime, Data: *info.CycleTime})
	}
	if info.CycleTimeExt != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaGateCycleTimeExt, Data: *info.CycleTimeExt})
	}
	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaGateFlags, Data: *info.Flags})
	}
	if info.ClockID != nil {
		options = append(options, tcOption{Interpretation: vtInt32, Type: tcaGateClockID, Data: *info.ClockID})
	}
	return marshalAttributes(options)
}

// unmarshalGate parses the gate-encoded data and stores the result in the value pointed to by info.
func unmarshalGate(data []byte, info *Gate) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaGateParms:
			parms := &GateParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaGateTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaGatePad:
			// padding does not contain data, we just skip it
		case tcaGatePriority:
			info.Priority = int32Ptr(ad.Int32())
		case tcaGateBaseTime:
			info.BaseTime = uint64Ptr(ad.Uint64())
		case tcaGateCycleTime:
			info.CycleTime = uint64Ptr(ad.Uint64())
		case tcaGateCycleTimeExt:
			info.CycleTimeExt = uint64Ptr(ad.Uint64())
		case tcaGateFlags:
			info.Flags = uint32Ptr(ad.Uint32())
		case tcaGateClockID:
			info.ClockID = int32Ptr(ad.Int32())
		default:
			return fmt.Errorf("UnmarshalGate()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// GateParms from include/uapi/linux/tc_act/tc_gate.h
type GateParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
}
