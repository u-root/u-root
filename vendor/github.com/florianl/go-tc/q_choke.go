package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaChokeUnspec = iota
	tcaChokeParms
	tcaChokeStab
	tcaChokeMaxP
)

// Choke contains attributes of the choke discipline
type Choke struct {
	Parms *RedQOpt
	MaxP  *uint32
}

// unmarshalChoke parses the Choke-encoded data and stores the result in the value pointed to by info.
func unmarshalChoke(data []byte, info *Choke) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaChokeParms:
			opt := &RedQOpt{}
			err = unmarshalStruct(ad.Bytes(), opt)
			multiError = concatError(multiError, err)
			info.Parms = opt
		case tcaChokeMaxP:
			info.MaxP = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalChoke()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalChoke returns the binary encoding of Choke
func marshalChoke(info *Choke) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Choke: %w", ErrNoArg)
	}

	var multiError error
	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaChokeParms, Data: data})
	}

	if info.MaxP != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaChokeMaxP, Data: uint32Value(info.MaxP)})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}
