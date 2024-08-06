package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaDsmarkUnspec = iota
	tcaDsmarkIndices
	tcaDsmarkDefaultIndex
	tcaDsmarkSetTCIndex
	tcaDsmarkMask
	tcaDsmarkValue
)

// Dsmark contains attributes of the dsmark discipline
type Dsmark struct {
	Indices      *uint16
	DefaultIndex *uint16
	SetTCIndex   *bool
	Mask         *uint8
	Value        *uint8
}

// unmarshalDsmark parses the Dsmark-encoded data and stores the result in the value pointed to by info.
func unmarshalDsmark(data []byte, info *Dsmark) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaDsmarkIndices:
			info.Indices = uint16Ptr(ad.Uint16())
		case tcaDsmarkDefaultIndex:
			info.DefaultIndex = uint16Ptr(ad.Uint16())
		case tcaDsmarkSetTCIndex:
			info.SetTCIndex = boolPtr(ad.Flag())
		case tcaDsmarkMask:
			info.Mask = uint8Ptr(ad.Uint8())
		case tcaDsmarkValue:
			info.Value = uint8Ptr(ad.Uint8())
		default:
			return fmt.Errorf("UnmarshalDsmark()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalDsmark returns the binary encoding of Qfq
func marshalDsmark(info *Dsmark) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Dsmark: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Indices != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaDsmarkIndices, Data: uint16Value(info.Indices)})
	}
	if info.DefaultIndex != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaDsmarkDefaultIndex, Data: uint16Value(info.DefaultIndex)})
	}
	if info.Mask != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaDsmarkMask, Data: uint8Value(info.Mask)})
	}
	if info.Value != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaDsmarkValue, Data: uint8Value(info.Value)})
	}
	if info.SetTCIndex != nil {
		options = append(options, tcOption{Interpretation: vtFlag, Type: tcaDsmarkSetTCIndex, Data: boolValue(info.SetTCIndex)})
	}
	return marshalAttributes(options)
}
