package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaBasicUnspec = iota
	tcaBasicClassID
	tcaBasicEmatches
	tcaBasicAct
	tcaBasicPolice
	tcaBasicPCNT
)

// Basic contains attributes of the basic discipline
type Basic struct {
	ClassID *uint32
	Police  *Police
	Ematch  *Ematch
	Actions *[]*Action
}

// unmarshalBasic parses the Basic-encoded data and stores the result in the value pointed to by info.
func unmarshalBasic(data []byte, info *Basic) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaBasicPolice:
			pol := &Police{}
			err := unmarshalPolice(ad.Bytes(), pol)
			multiError = concatError(multiError, err)
			info.Police = pol
		case tcaBasicClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaBasicEmatches:
			ematch := &Ematch{}
			err := unmarshalEmatch(ad.Bytes(), ematch)
			multiError = concatError(multiError, err)
			info.Ematch = ematch
		case tcaBasicAct:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Actions = actions
		case tcaBasicPCNT:
			continue
		default:
			return fmt.Errorf("unmarshalBasic()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalBasic returns the binary encoding of Basic
func marshalBasic(info *Basic) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Basic: %w", ErrNoArg)
	}
	var multiError error

	// TODO: improve logic and check combinations
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBasicClassID, Data: uint32Value(info.ClassID)})
	}
	if info.Ematch != nil {
		data, err := marshalEmatch(info.Ematch)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBasicEmatches, Data: data})
	}
	if info.Police != nil {
		data, err := marshalPolice(info.Police)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBasicPolice, Data: data})
	}
	if info.Actions != nil {
		data, err := marshalActions(0, *info.Actions)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBasicAct, Data: data})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}
