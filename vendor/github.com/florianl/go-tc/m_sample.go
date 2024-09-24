package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaSampleUnspec = iota
	tcaSampleTm
	tcaSampleParms
	tcaSampleRate
	tcaSampleTruncSize
	tcaSamplePSampleGroup
	tcaSamplePad
)

// Sample contains attribute of the Sample discipline
type Sample struct {
	Parms       *SampleParms
	Tm          *Tcft
	Rate        *uint32
	TruncSize   *uint32
	SampleGroup *uint32
}

// SampleParms from include/uapi/linux/tc_act/tc_sample.h
type SampleParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
}

// marshalSample returns the binary encoding of Sample
func marshalSample(info *Sample) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Sample: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaSampleParms, Data: data})
	}
	if info.Rate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaSampleRate, Data: *info.Rate})
	}
	if info.TruncSize != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaSampleTruncSize, Data: *info.TruncSize})
	}
	if info.SampleGroup != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaSamplePSampleGroup, Data: *info.SampleGroup})
	}
	return marshalAttributes(options)
}

// unmarshalSample parses the Sample-encoded data and stores the result in the value pointed to by info.
func unmarshalSample(data []byte, info *Sample) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaSampleParms:
			parms := &SampleParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaSampleTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaSampleRate:
			info.Rate = uint32Ptr(ad.Uint32())
		case tcaSampleTruncSize:
			info.TruncSize = uint32Ptr(ad.Uint32())
		case tcaSamplePSampleGroup:
			info.SampleGroup = uint32Ptr(ad.Uint32())
		case tcaSamplePad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalSample()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
