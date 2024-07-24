package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFlowUnspec = iota
	tcaFlowKeys
	tcaFlowMode
	tcaFlowBaseClass
	tcaFlowRShift
	tcaFlowAddend
	tcaFlowMask
	tcaFlowXOR
	tcaFlowDivisor
	tcaFlowAct
	tcaFlowPolice
	tcaFlowEMatches
	tcaFlowPerTurb
)

// Flow contains attributes of the flow discipline
type Flow struct {
	Keys      *uint32
	Mode      *uint32
	BaseClass *uint32
	RShift    *uint32
	Addend    *uint32
	Mask      *uint32
	XOR       *uint32
	Divisor   *uint32
	PerTurb   *uint32
	Ematch    *Ematch
	Actions   *[]*Action
}

// unmarshalFlow parses the Flow-encoded data and stores the result in the value pointed to by info.
func unmarshalFlow(data []byte, info *Flow) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaFlowKeys:
			info.Keys = uint32Ptr(ad.Uint32())
		case tcaFlowMode:
			info.Mode = uint32Ptr(ad.Uint32())
		case tcaFlowBaseClass:
			info.BaseClass = uint32Ptr(ad.Uint32())
		case tcaFlowRShift:
			info.RShift = uint32Ptr(ad.Uint32())
		case tcaFlowAddend:
			info.Addend = uint32Ptr(ad.Uint32())
		case tcaFlowMask:
			info.Mask = uint32Ptr(ad.Uint32())
		case tcaFlowXOR:
			info.XOR = uint32Ptr(ad.Uint32())
		case tcaFlowDivisor:
			info.Divisor = uint32Ptr(ad.Uint32())
		case tcaFlowPerTurb:
			info.PerTurb = uint32Ptr(ad.Uint32())
		case tcaFlowEMatches:
			ematch := &Ematch{}
			err := unmarshalEmatch(ad.Bytes(), ematch)
			multiError = concatError(multiError, err)
			info.Ematch = ematch
		case tcaFlowAct:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Actions = actions
		default:
			return fmt.Errorf("unmarshalFlow()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalFlow returns the binary encoding of Flow
func marshalFlow(info *Flow) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Flow: %w", ErrNoArg)
	}

	var multiError error
	// TODO: improve logic and check combinations
	if info.Keys != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowKeys, Data: uint32Value(info.Keys)})
	}
	if info.Mode != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowMode, Data: uint32Value(info.Mode)})
	}
	if info.BaseClass != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowBaseClass, Data: uint32Value(info.BaseClass)})
	}
	if info.RShift != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowRShift, Data: uint32Value(info.RShift)})
	}
	if info.Addend != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowAddend, Data: uint32Value(info.Addend)})
	}
	if info.Mask != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowMask, Data: uint32Value(info.Mask)})
	}
	if info.XOR != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowXOR, Data: uint32Value(info.XOR)})
	}
	if info.Divisor != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowDivisor, Data: uint32Value(info.Divisor)})
	}
	if info.PerTurb != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowPerTurb, Data: uint32Value(info.PerTurb)})
	}
	if info.Ematch != nil {
		data, err := marshalEmatch(info.Ematch)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFlowEMatches, Data: data})
	}
	if info.Actions != nil {
		data, err := marshalActions(0, *info.Actions)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFlowAct, Data: data})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}
