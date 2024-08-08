package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaSfbUnspec = iota
	tcaSfbParms
)

// Sfb contains attributes of the SBF discipline
type Sfb struct {
	Parms *SfbQopt
}

// unmarshalSfb parses the Sfb-encoded data and stores the result in the value pointed to by info.
func unmarshalSfb(data []byte, info *Sfb) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaSfbParms:
			opt := &SfbQopt{}
			multiError = unmarshalStruct(ad.Bytes(), opt)
			info.Parms = opt
		default:
			return fmt.Errorf("extractSfbOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalSfb returns the binary encoding of Sfb
func marshalSfb(info *Sfb) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Sfb: %w", ErrNoArg)
	}

	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaSfbParms, Data: data})
	}

	// TODO: improve logic and check combinations
	return marshalAttributes(options)
}

// SfbQopt from include/uapi/linux/pkt_sched.h
type SfbQopt struct {
	RehashInterval uint32 // in ms
	WarmupTime     uint32 //  in ms
	Max            uint32
	BinSize        uint32
	Increment      uint32
	Decrement      uint32
	Limit          uint32
	PenaltyRate    uint32
	PenaltyBurst   uint32
}
