package tc

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

const (
	tcaCtUnspec = iota
	tcaCtParms
	tcaCtTm
	tcaCtAction     /* u16 */
	tcaCtZone       /* u16 */
	tcaCtMark       /* u32 */
	tcaCtMarkMask   /* u32 */
	tcaCtLabels     /* u128 */
	tcaCtLabelsMask /* u128 */
	tcaCtNatIPv4Min /* be32 */
	tcaCtNatIPv4Max /* be32 */
	tcaCtNatIPv6Min /* struct in6_addr */
	tcaCtNatIPv6Max /* struct in6_addr */
	tcaCtNatPortMin /* be16 */
	tcaCtNatPortMax /* be16 */
	tcaCtPad
)

// Ct contains attributes of the ct discipline
type Ct struct {
	Parms      *CtParms
	Tm         *Tcft
	Action     *uint16
	Zone       *uint16
	Mark       *uint32
	MarkMask   *uint32
	NatIPv4Min *net.IP
	NatIPv4Max *net.IP
	NatPortMin *uint16
	NatPortMax *uint16
}

// CtParms contains further ct attributes.
type CtParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
}

// unmarshalCt parses the ct-encoded data and stores the result in the value pointed to by info.
func unmarshalCt(data []byte, info *Ct) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaCtParms:
			parms := &CtParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaCtTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaCtAction:
			info.Action = uint16Ptr(ad.Uint16())
		case tcaCtZone:
			info.Zone = uint16Ptr(ad.Uint16())
		case tcaCtMark:
			info.Mark = uint32Ptr(ad.Uint32())
		case tcaCtMarkMask:
			info.MarkMask = uint32Ptr(ad.Uint32())
		case tcaCtNatIPv4Min:
			tmp := uint32ToIP(ad.Uint32())
			info.NatIPv4Min = &tmp
		case tcaCtNatIPv4Max:
			tmp := uint32ToIP(ad.Uint32())
			info.NatIPv4Max = &tmp
		case tcaCtNatPortMin:
			info.NatPortMin = uint16Ptr(ad.Uint16())
		case tcaCtNatPortMax:
			info.NatPortMax = uint16Ptr(ad.Uint16())
		case tcaCtPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalCt()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalCt returns the binary encoding of Ct
func marshalCt(info *Ct) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ct: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	var multiError error
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCtParms, Data: data})
	}
	if info.Action != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaCtAction, Data: uint16Value(info.Action)})
	}
	if info.Zone != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaCtZone, Data: uint16Value(info.Zone)})
	}
	if info.Mark != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCtMark, Data: uint32Value(info.Mark)})
	}
	if info.MarkMask != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCtMarkMask, Data: uint32Value(info.MarkMask)})
	}
	if info.NatIPv4Min != nil {
		tmp, err := ipToUint32(*info.NatIPv4Min)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCtNatIPv4Min, Data: tmp})
	}
	if info.NatIPv4Max != nil {
		tmp, err := ipToUint32(*info.NatIPv4Max)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCtNatIPv4Max, Data: tmp})
	}
	if info.NatPortMin != nil {
		options = append(options, tcOption{Interpretation: vtUint16Be, Type: tcaCtNatPortMin, Data: uint16Value(info.NatPortMin)})
	}
	if info.NatPortMax != nil {
		options = append(options, tcOption{Interpretation: vtUint16Be, Type: tcaCtNatPortMax, Data: uint16Value(info.NatPortMax)})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}
