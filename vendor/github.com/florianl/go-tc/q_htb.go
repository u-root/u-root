package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaHtbUnspec = iota
	tcaHtbParms
	tcaHtbInit
	tcaHtbCtab
	tcaHtbRtab
	tcaHtbDirectQlen
	tcaHtbRate64
	tcaHtbCeil64
	tcaHtbPad
)

// Htb contains attributes of the HTB discipline
type Htb struct {
	Parms      *HtbOpt
	Init       *HtbGlob
	Ctab       *[]byte
	Rtab       *[]byte
	DirectQlen *uint32
	Rate64     *uint64
	Ceil64     *uint64
}

// unmarshalHtb parses the Htb-encoded data and stores the result in the value pointed to by info.
func unmarshalHtb(data []byte, info *Htb) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaHtbParms:
			opt := &HtbOpt{}
			err := unmarshalStruct(ad.Bytes(), opt)
			multiError = concatError(multiError, err)
			info.Parms = opt
		case tcaHtbInit:
			glob := &HtbGlob{}
			err := unmarshalStruct(ad.Bytes(), glob)
			multiError = concatError(multiError, err)
			info.Init = glob
		case tcaHtbCtab:
			info.Ctab = bytesPtr(ad.Bytes())
		case tcaHtbRtab:
			info.Rtab = bytesPtr(ad.Bytes())
		case tcaHtbDirectQlen:
			info.DirectQlen = uint32Ptr(ad.Uint32())
		case tcaHtbRate64:
			info.Rate64 = uint64Ptr(ad.Uint64())
		case tcaHtbCeil64:
			info.Ceil64 = uint64Ptr(ad.Uint64())
		case tcaHtbPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalHtb()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalHtb returns the binary encoding of Qfq
func marshalHtb(info *Htb) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Htb: %w", ErrNoArg)
	}
	var multiError error
	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHtbParms, Data: data})
	}
	if info.Init != nil {
		data, err := marshalStruct(info.Init)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHtbInit, Data: data})
	}
	if info.DirectQlen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHtbDirectQlen, Data: uint32Value(info.DirectQlen)})
	}
	if info.Rate64 != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaHtbRate64, Data: uint64Value(info.Rate64)})
	}
	if info.Ceil64 != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaHtbCeil64, Data: uint64Value(info.Ceil64)})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}

// HtbGlob from include/uapi/linux/pkt_sched.h
type HtbGlob struct {
	Version      uint32
	Rate2Quantum uint32
	Defcls       uint32
	Debug        uint32
	DirectPkts   uint32
}

// HtbOpt from include/uapi/linux/pkt_sched.h
type HtbOpt struct {
	Rate    RateSpec
	Ceil    RateSpec
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
}
