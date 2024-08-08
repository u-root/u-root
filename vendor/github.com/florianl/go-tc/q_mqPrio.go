package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaMqPrioUnspec = iota
	tcaMqPrioMode
	tcaMqPrioShaper
	tcaMqPrioMinRate64
	tcaMqPrioMaxRate64
)

// MqPrio contains attributes of the mqprio discipline
type MqPrio struct {
	Opt       *MqPrioQopt
	Mode      *uint16
	Shaper    *uint16
	MinRate64 *uint64
	MaxRate64 *uint64
}

// MqPrioQopt according to tc_mqprio_qopt in /include/uapi/linux/pkt_sched.h
type MqPrioQopt struct {
	NumTc     uint8
	PrioTcMap [16]uint8 //  TC_QOPT_BITMASK + 1 = 16
	Hw        uint8
	Count     [16]uint16 // TC_QOPT_MAX_QUEUE = 16
	Offset    [16]uint16 // TC_QOPT_MAX_QUEUE = 16
}

// unmarshalMqPrio parses the MqPrio-encoded data and stores the result in the value pointed to by info.
func unmarshalMqPrio(data []byte, info *MqPrio) error {
	opt := &MqPrioQopt{}
	if err := unmarshalStruct(data, opt); err != nil {
		return err
	}
	info.Opt = opt

	// The size of MqPrioQopt is 82 bytes. To align it to 4 byte boundaries
	// we add two.
	data = data[84:]

	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaMqPrioMode:
			info.Mode = uint16Ptr(ad.Uint16())
		case tcaMqPrioShaper:
			info.Shaper = uint16Ptr(ad.Uint16())
		case tcaMqPrioMinRate64:
			info.MinRate64 = uint64Ptr(ad.Uint64())
		case tcaMqPrioMaxRate64:
			info.MaxRate64 = uint64Ptr(ad.Uint64())
		default:
			return fmt.Errorf("unmarshalMqPrio()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalMqPrio returns the binary encoding of MqPrio
func marshalMqPrio(info *MqPrio) ([]byte, error) {
	options := []tcOption{}

	if info == nil || info.Opt == nil {
		return []byte{}, fmt.Errorf("MqPrio: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Mode != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaMqPrioMode, Data: uint16Value(info.Mode)})
	}
	if info.Shaper != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaMqPrioShaper, Data: uint16Value(info.Shaper)})
	}
	if info.MinRate64 != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaMqPrioMinRate64, Data: uint64Value(info.MinRate64)})
	}
	if info.MaxRate64 != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaMqPrioMaxRate64, Data: uint64Value(info.MaxRate64)})
	}

	opt, err := marshalAndAlignStruct(info.Opt)
	if err != nil {
		return []byte{}, err
	}
	adds, err := marshalAttributes(options)
	if err != nil {
		return []byte{}, err
	}
	opt = append(opt, adds...)
	return opt, nil
}
