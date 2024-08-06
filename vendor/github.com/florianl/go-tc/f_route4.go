package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaRoute4Unspec = iota
	tcaRoute4ClassID
	tcaRoute4To
	tcaRoute4From
	tcaRoute4IIf
	tcaRoute4Police
	tcaRoute4Act
)

// Route4 contains attributes of the route discipline
type Route4 struct {
	ClassID *uint32
	To      *uint32
	From    *uint32
	IIf     *uint32
	Actions *[]*Action
}

// unmarshalRoute4 parses the Route4-encoded data and stores the result in the value pointed to by info.
func unmarshalRoute4(data []byte, info *Route4) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaRoute4ClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaRoute4To:
			info.To = uint32Ptr(ad.Uint32())
		case tcaRoute4From:
			info.From = uint32Ptr(ad.Uint32())
		case tcaRoute4IIf:
			info.IIf = uint32Ptr(ad.Uint32())
		case tcaRoute4Act:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Actions = actions
		default:
			return fmt.Errorf("unmarshalRoute4()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalRoute4 returns the binary encoding of Route4
func marshalRoute4(info *Route4) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Route4: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations

	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRoute4ClassID, Data: uint32Value(info.ClassID)})
	}
	if info.To != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRoute4To, Data: uint32Value(info.To)})
	}
	if info.From != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRoute4From, Data: uint32Value(info.From)})
	}
	if info.IIf != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRoute4IIf, Data: uint32Value(info.IIf)})
	}
	if info.Actions != nil {
		data, err := marshalActions(0, *info.Actions)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRoute4Act, Data: data})
	}

	return marshalAttributes(options)
}
