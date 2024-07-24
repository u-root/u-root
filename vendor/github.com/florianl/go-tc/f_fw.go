package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFwUnspec = iota
	tcaFwClassID
	tcaFwPolice
	tcaFwInDev
	tcaFwAct
	tcaFwMask
)

// Fw contains attributes of the fw discipline
type Fw struct {
	ClassID *uint32
	Police  *Police
	InDev   *string
	Mask    *uint32
	Actions *[]*Action
}

// unmarshalFw parses the Fw-encoded data and stores the result in the value pointed to by info.
func unmarshalFw(data []byte, info *Fw) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaFwClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaFwInDev:
			info.InDev = stringPtr(ad.String())
		case tcaFwMask:
			info.Mask = uint32Ptr(ad.Uint32())
		case tcaFwPolice:
			pol := &Police{}
			err := unmarshalPolice(ad.Bytes(), pol)
			multiError = concatError(multiError, err)
			info.Police = pol
		case tcaFwAct:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Actions = actions
		default:
			return fmt.Errorf("unmarshalFw()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalFw returns the binary encoding of Fw
func marshalFw(info *Fw) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Fw: %w", ErrNoArg)
	}
	var multiError error

	// TODO: improve logic and check combinations
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFwClassID, Data: uint32Value(info.ClassID)})
	}
	if info.Mask != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFwMask, Data: uint32Value(info.Mask)})
	}
	if info.InDev != nil {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaFwInDev, Data: stringValue(info.InDev)})
	}
	if info.Police != nil {
		data, err := marshalPolice(info.Police)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFwPolice, Data: data})
	}
	if info.Actions != nil {
		data, err := marshalActions(0, *info.Actions)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFwAct, Data: data})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}
