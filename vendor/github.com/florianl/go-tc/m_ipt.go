package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaIptUnspec = iota
	tcaIptTable
	tcaIptHook
	tcaIptIndex
	tcaIptCnt
	tcaIptTm
	tcaIptTarg
	tcaIptPad
)

// Ipt contains attribute of the ipt discipline
type Ipt struct {
	Table *string
	Hook  *uint32
	Index *uint32
	Cnt   *IptCnt
	Tm    *Tcft
}

// IptCnt as tc_cnt from include/uapi/linux/pkt_cls.h
type IptCnt struct {
	RefCnt  uint32
	BindCnt uint32
}

// unmarshalIpt parses the ipt-encoded data and stores the result in the value pointed to by info.
func unmarshalIpt(data []byte, info *Ipt) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaIptTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaIptTable:
			info.Table = stringPtr(ad.String())
		case tcaIptHook:
			info.Hook = uint32Ptr(ad.Uint32())
		case tcaIptIndex:
			info.Index = uint32Ptr(ad.Uint32())
		case tcaIptCnt:
			tmp := &IptCnt{}
			err = unmarshalStruct(ad.Bytes(), tmp)
			multiError = concatError(multiError, err)
			info.Cnt = tmp
		case tcaIptPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalIpt()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalIpt returns the binary encoding of Ipt
func marshalIpt(info *Ipt) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ipt: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Table != nil {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaIptTable, Data: stringValue(info.Table)})
	}
	if info.Hook != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIptHook, Data: uint32Value(info.Hook)})
	}
	if info.Index != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIptIndex, Data: uint32Value(info.Index)})
	}
	if info.Cnt != nil {
		data, err := marshalStruct(info.Cnt)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIptCnt, Data: data})
	}

	return marshalAttributes(options)
}
