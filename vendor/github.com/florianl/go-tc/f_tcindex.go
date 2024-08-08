package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaTcIndexUnspec = iota
	tcaTcIndexHash
	tcaTcIndexMask
	tcaTcIndexShift
	tcaTcIndexFallThrough
	tcaTcIndexClassID
	tcaTcIndexPolice
	tcaTcIndexAct
)

// TcIndex contains attributes of the tcIndex discipline
type TcIndex struct {
	Hash        *uint32
	Mask        *uint16
	Shift       *uint32
	FallThrough *uint32
	ClassID     *uint32
	Actions     *[]*Action
}

// marshalTcIndex returns the binary encoding of TcIndex
func marshalTcIndex(info *TcIndex) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("TcIndex: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Hash != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTcIndexHash, Data: uint32Value(info.Hash)})
	}
	if info.Mask != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaTcIndexMask, Data: uint16Value(info.Mask)})
	}
	if info.Shift != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTcIndexShift, Data: uint32Value(info.Shift)})
	}
	if info.FallThrough != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTcIndexFallThrough, Data: uint32Value(info.FallThrough)})
	}
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTcIndexClassID, Data: uint32Value(info.ClassID)})
	}
	if info.Actions != nil {
		data, err := marshalActions(0, *info.Actions)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTcIndexAct, Data: data})
	}
	return marshalAttributes(options)
}

// unmarshalTcIndex parses the TcIndex-encoded data and stores the result in the value pointed to by info.
func unmarshalTcIndex(data []byte, info *TcIndex) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaTcIndexHash:
			info.Hash = uint32Ptr(ad.Uint32())
		case tcaTcIndexMask:
			info.Mask = uint16Ptr(ad.Uint16())
		case tcaTcIndexShift:
			info.Shift = uint32Ptr(ad.Uint32())
		case tcaTcIndexFallThrough:
			info.FallThrough = uint32Ptr(ad.Uint32())
		case tcaTcIndexClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaTcIndexAct:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Actions = actions
		default:
			return fmt.Errorf("unmarshalTcIndex()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
