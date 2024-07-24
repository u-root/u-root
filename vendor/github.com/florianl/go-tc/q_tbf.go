package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaTbfUnspec = iota
	tcaTbfParms
	tcaTbfRtab
	tcaTbfPtab
	tcaTbfRate64
	tcaTbfPrate64
	tcaTbfBurst
	tcaTbfPburst
	tcaTbfPad
)

// Tbf contains attributes of the TBF discipline
type Tbf struct {
	Parms  *TbfQopt
	Burst  *uint32
	Pburst *uint32
}

// unmarshalTbf parses the FqCodel-encoded data and stores the result in the value pointed to by info.
func unmarshalTbf(data []byte, info *Tbf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaTbfParms:
			qopt := &TbfQopt{}
			err := unmarshalStruct(ad.Bytes(), qopt)
			multiError = concatError(multiError, err)
			info.Parms = qopt
		case tcaTbfBurst:
			info.Burst = uint32Ptr(ad.Uint32())
		case tcaTbfPburst:
			info.Pburst = uint32Ptr(ad.Uint32())
		case tcaTbfPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalTbf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalTbf returns the binary encoding of Tbf
func marshalTbf(info *Tbf) ([]byte, error) {
	options := []tcOption{}

	if info == nil || info.Parms == nil {
		return []byte{}, fmt.Errorf("Tbf: %w", ErrNoArg)
	}
	var multiError error
	// TODO: improve logic and check combinations
	if info.Parms.Rate.Rate != 0 {
		ratePolicy := Policy{}
		ratePolicy.Burst = uint32Value(info.Burst)
		ratePolicy.Action = PolicyOk
		ratePolicy.Limit = info.Parms.Limit
		ratePolicy.Rate.Rate = info.Parms.Rate.Rate

		rtab, err := generateRateTable(&ratePolicy)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTbfRtab, Data: rtab})
	}
	if info.Parms.PeakRate.Rate != 0 {
		ratePolicy := Policy{}
		ratePolicy.Burst = uint32Value(info.Pburst)
		ratePolicy.Action = PolicyOk
		ratePolicy.Limit = info.Parms.Limit
		ratePolicy.PeakRate.Rate = info.Parms.PeakRate.Rate

		ptab, err := generateRateTable(&ratePolicy)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTbfPtab, Data: ptab})
	}
	data, err := marshalStruct(info.Parms)
	multiError = concatError(multiError, err)
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTbfParms, Data: data})

	if info.Burst != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTbfBurst, Data: uint32Value(info.Burst)})
	}
	if info.Pburst != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTbfPburst, Data: uint32Value(info.Pburst)})
	}

	if multiError != nil {
		return []byte{}, fmt.Errorf("Tbf: %w", multiError)
	}

	return marshalAttributes(options)
}

// TbfQopt from include/uapi/linux/pkt_sched.h
type TbfQopt struct {
	Rate     RateSpec
	PeakRate RateSpec
	Limit    uint32
	Buffer   uint32
	Mtu      uint32
}
