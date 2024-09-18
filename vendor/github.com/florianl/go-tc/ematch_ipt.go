package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaEmIptUnspec = iota
	tcaEmIptHook
	tcaEmIptMatchName
	tcaEmIptMatchRevision
	tcaEmIptNFProto
	tcaEmIptMatchData
)

// IptMatch contains attributes of the ipt match discipline
type IptMatch struct {
	Hook      *uint32
	MatchName *string
	Revision  *uint8
	NFProto   *uint8
	MatchData *[]byte
}

func unmarshalIptMatch(data []byte, info *IptMatch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaEmIptHook:
			info.Hook = uint32Ptr(ad.Uint32())
		case tcaEmIptMatchName:
			info.MatchName = stringPtr(ad.String())
		case tcaEmIptMatchRevision:
			info.Revision = uint8Ptr(ad.Uint8())
		case tcaEmIptNFProto:
			info.NFProto = uint8Ptr(ad.Uint8())
		case tcaEmIptMatchData:
			info.MatchData = bytesPtr(ad.Bytes())
		default:
			return fmt.Errorf("unmarshalIptMatch()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

func marshalIptMatch(info *IptMatch) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("IptMatch: %w", ErrNoArg)
	}
	var multiError error

	// TODO: improve logic and check combinations
	if info.Hook != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaEmIptHook, Data: uint32Value(info.Hook)})
	}
	if info.MatchName != nil {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaEmIptMatchName, Data: stringValue(info.MatchName)})
	}
	if info.Revision != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaEmIptMatchRevision, Data: uint8Value(info.Revision)})
	}
	if info.NFProto != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaEmIptNFProto, Data: uint8Value(info.NFProto)})
	}
	if info.MatchData != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmIptMatchData, Data: bytesValue(info.MatchData)})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}
