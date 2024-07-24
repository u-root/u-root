package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaMatchallUnspec = iota
	tcaMatchallClassID
	tcaMatchallAct
	tcaMatchallFlags
	tcaMatchallPcnt
	tcaMatchallPad
)

// Matchall contains attributes of the matchall discipline
type Matchall struct {
	ClassID *uint32
	Actions *[]*Action
	Flags   *uint32
	Pcnt    *uint64
}

func unmarshalMatchall(data []byte, info *Matchall) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaMatchallClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaMatchallAct:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Actions = actions
		case tcaMatchallFlags:
			info.Flags = uint32Ptr(ad.Uint32())
		case tcaMatchallPcnt:
			info.Pcnt = uint64Ptr(ad.Uint64())
		case tcaMatchallPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalMatchall()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalMatchall returns the binary encoding of Matchall
func marshalMatchall(info *Matchall) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Matchall: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	var multiError error

	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaMatchallClassID, Data: uint32Value(info.ClassID)})
	}
	if info.Actions != nil {
		data, err := marshalActions(0, *info.Actions)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaMatchallAct, Data: data})
	}

	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaMatchallFlags, Data: uint32Value(info.Flags)})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}
