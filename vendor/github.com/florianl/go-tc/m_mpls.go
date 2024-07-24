package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaMPLSUnspec = iota
	tcaMPLSTm
	tcaMPLSParms
	tcaMPLSPad
	tcaMPLSProto /* be16; eth_type of pushed or next (for pop) header. */
	tcaMPLSLabel
	tcaMPLSTC
	tcaMPLSTTL
	tcaMPLSBOS
)

// MPLSAction defines MPS actions.
type MPLSAction int32

// Various MPLS actions
const (
	MPLSActPop     = MPLSAction(1)
	MPLSActPush    = MPLSAction(2)
	MPLSActModify  = MPLSAction(3)
	MPLSActDecTTL  = MPLSAction(4)
	MPLSActMACPush = MPLSAction(5)
)

// MPLS contains attributes of the mpls discipline
// https://man7.org/linux/man-pages/man8/tc-mpls.8.html
type MPLS struct {
	Parms *MPLSParam
	Tm    *Tcft
	Proto *int16
	Label *uint32
	TC    *uint8
	TTL   *uint8
	BOS   *uint8
}

// MPLSParam contains further MPLS attributes.
type MPLSParam struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
	MAction MPLSAction
}

func unmarshalMPLS(data []byte, info *MPLS) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaMPLSTm:
			tm := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tm)
			multiError = concatError(multiError, err)
			info.Tm = tm
		case tcaMPLSParms:
			param := &MPLSParam{}
			err = unmarshalStruct(ad.Bytes(), param)
			multiError = concatError(multiError, err)
			info.Parms = param
		case tcaMPLSPad:
			// padding does not contain data, we just skip it
		case tcaMPLSProto: /* be16; eth_type of pushed or next (for pop) header. */
			info.Proto = int16Ptr(ad.Int16())
		case tcaMPLSLabel:
			info.Label = uint32Ptr(ad.Uint32())
		case tcaMPLSTC:
			info.TC = uint8Ptr(ad.Uint8())
		case tcaMPLSTTL:
			info.TTL = uint8Ptr(ad.Uint8())
		case tcaMPLSBOS:
			info.BOS = uint8Ptr(ad.Uint8())
		default:
			return fmt.Errorf("unmarshalMPLS()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

func marshalMPLS(info *MPLS) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("MPLS: %w", ErrNoArg)
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
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaMPLSParms, Data: data})
	}
	if info.Proto != nil {
		options = append(options, tcOption{Interpretation: vtInt16Be, Type: tcaMPLSProto, Data: *info.Proto})
	}
	if info.Label != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaMPLSLabel, Data: *info.Label})
	}
	if info.TC != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaMPLSTC, Data: *info.TC})
	}
	if info.TTL != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaMPLSTTL, Data: *info.TTL})
	}
	if info.BOS != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaMPLSBOS, Data: *info.BOS})
	}

	return marshalAttributes(options)
}
