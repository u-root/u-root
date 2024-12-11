package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaStabUnspec = iota
	tcaStabBase
	tcaStabData
)

// SizeSpec implements tc_sizespec
type SizeSpec struct {
	CellLog   uint8
	SizeLog   uint8
	CellAlign int16
	Overhead  int32
	LinkLayer uint32
	MPU       uint32
	MTU       uint32
	TSize     uint32
}

// Stab contains attributes of a stab
// http://man7.org/linux/man-pages/man8/tc-stab.8.html
type Stab struct {
	Base *SizeSpec
	Data *[]byte
}

func unmarshalStab(data []byte, stab *Stab) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaStabBase:
			base := &SizeSpec{}
			err := unmarshalStruct(ad.Bytes(), base)
			multiError = concatError(multiError, err)
			stab.Base = base
		case tcaStabData:
			tmp := ad.Bytes()
			stab.Data = &tmp
		default:
			return fmt.Errorf("unmarshalStab()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

func marshalStab(info *Stab) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Stab: %w", ErrNoArg)
	}

	var multiError error

	// TODO: improve logic and check combination
	if info.Base != nil {
		data, err := marshalStruct(info.Base)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStabBase, Data: data})
	}
	if info.Data != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStabData, Data: *info.Data})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}
