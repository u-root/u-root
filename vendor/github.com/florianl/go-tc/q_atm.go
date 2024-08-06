package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaAtmUnspec = iota
	tcaAtmFD
	tcaAtmPtr
	tcaAtmHdr
	tcaAtmExcess
	tcaAtmAddr
	tcaAtmState
)

// Atm contains attributes of the atm discipline
type Atm struct {
	FD     *uint32
	Excess *uint32
	Addr   *AtmPvc
	State  *uint32
}

// unmarshalAtm parses the Atm-encoded data and stores the result in the value pointed to by info.
func unmarshalAtm(data []byte, info *Atm) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaAtmFD:
			info.FD = uint32Ptr(ad.Uint32())
		case tcaAtmExcess:
			info.Excess = uint32Ptr(ad.Uint32())
		case tcaAtmAddr:
			arg := &AtmPvc{}
			err := unmarshalStruct(ad.Bytes(), arg)
			multiError = concatError(multiError, err)
			info.Addr = arg
		case tcaAtmState:
			info.State = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalAtm()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	return concatError(multiError, ad.Err())
}

// marshalAtm returns the binary encoding of Atm
func marshalAtm(info *Atm) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Atm: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations

	if info.Addr != nil {
		data, err := marshalStruct(info.Addr)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaAtmAddr, Data: data})
	}
	if info.FD != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaAtmFD, Data: uint32Value(info.FD)})
	}
	if info.Excess != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaAtmExcess, Data: uint32Value(info.Excess)})
	}
	if info.State != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaAtmState, Data: uint32Value(info.State)})
	}

	return marshalAttributes(options)
}

// AtmPvc from include/uapi/linux/atm.h
type AtmPvc struct {
	SapFamily byte
	Itf       byte
	Vpi       byte
	Vci       byte
}
