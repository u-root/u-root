package tc

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

const (
	tcaTunnelUnspec = iota
	tcaTunnelKeyTm
	tcaTunnelKeyParms
	tcaTunnelKeyEncIPv4Src
	tcaTunnelKeyEncIPv4Dst
	tcaTunnelKeyEncIPv6Src
	tcaTunnelKeyEncIPv6Dst
	tcaTunnelKeyEncKeyID
	tcaTunnelKeyPad
	tcaTunnelKeyEncDstPort
	tcaTunnelKeyNoCSUM
	tcaTunnelKeyEncOpts
	tcaTunnelKeyEncTOS
	tcaTunnelKeyEncTTL
)

// TunnelKey contains attribute of the TunnelKey discipline
type TunnelKey struct {
	Parms         *TunnelParms
	Tm            *Tcft
	KeyEncSrc     *net.IP
	KeyEncDst     *net.IP
	KeyEncKeyID   *uint32
	KeyEncDstPort *uint16
	KeyNoCSUM     *uint8
	KeyEncTOS     *uint8
	KeyEncTTL     *uint8
}

// TunnelParms from include/uapi/linux/tc_act/tc_tunnel_key.h
type TunnelParms struct {
	Index           uint32
	Capab           uint32
	Action          uint32
	RefCnt          uint32
	BindCnt         uint32
	TunnelKeyAction uint32
}

// marshalTunnelKey returns the binary encoding of TunnelKey
func marshalTunnelKey(info *TunnelKey) ([]byte, error) {
	options := []tcOption{}
	if info == nil {
		return []byte{}, fmt.Errorf("TunnelKey: %w", ErrNoArg)
	}

	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTunnelKeyParms, Data: data})
	}
	if info.KeyEncSrc != nil {
		if info.KeyEncSrc.To4() != nil {
			tmp, err := ipToUint32(*info.KeyEncSrc)
			if err != nil {
				return []byte{}, fmt.Errorf("TunnelKey - KeyEncIPv4Src: %w", err)
			}
			options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTunnelKeyEncIPv4Src, Data: tmp})
		} else {
			tmp := ipToBytes(*info.KeyEncSrc)
			options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTunnelKeyEncIPv6Src, Data: tmp})
		}
	}
	if info.KeyEncDst != nil {
		if info.KeyEncDst.To4() != nil {
			tmp, err := ipToUint32(*info.KeyEncDst)
			if err != nil {
				return []byte{}, fmt.Errorf("TunnelKey - KeyEncIPv4Src: %w", err)
			}
			options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTunnelKeyEncIPv4Dst, Data: tmp})
		} else {
			tmp := ipToBytes(*info.KeyEncSrc)
			options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTunnelKeyEncIPv6Dst, Data: tmp})
		}
	}
	if info.KeyEncKeyID != nil {
		options = append(options, tcOption{Interpretation: vtUint32Be, Type: tcaTunnelKeyEncKeyID, Data: *info.KeyEncKeyID})
	}
	if info.KeyEncDstPort != nil {
		options = append(options, tcOption{Interpretation: vtUint16Be, Type: tcaTunnelKeyEncDstPort, Data: *info.KeyEncDstPort})
	}
	if info.KeyNoCSUM != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaTunnelKeyNoCSUM, Data: *info.KeyNoCSUM})
	}
	if info.KeyEncTOS != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaTunnelKeyEncTOS, Data: *info.KeyEncTOS})
	}
	if info.KeyEncTTL != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaTunnelKeyEncTTL, Data: *info.KeyEncTTL})
	}

	return marshalAttributes(options)
}

// unmarshalTunnelKey parses the TunnelKey-encoded data and stores the result in the value pointed to by info.
func unmarshalTunnelKey(data []byte, info *TunnelKey) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaTunnelKeyTm:
			tm := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tm)
			multiError = concatError(multiError, err)
			info.Tm = tm
		case tcaTunnelKeyParms:
			parms := &TunnelParms{}
			err = unmarshalStruct(ad.Bytes(), parms)
			multiError = concatError(multiError, err)
			info.Parms = parms
		case tcaTunnelKeyEncIPv4Src:
			tmp := uint32ToIP(ad.Uint32())
			info.KeyEncSrc = &tmp
		case tcaTunnelKeyEncIPv4Dst:
			tmp := uint32ToIP(ad.Uint32())
			info.KeyEncDst = &tmp
		case tcaTunnelKeyEncIPv6Src:
			tmp, err := bytesToIP(ad.Bytes())
			multiError = concatError(multiError, err)
			info.KeyEncSrc = &tmp
		case tcaTunnelKeyEncIPv6Dst:
			tmp, err := bytesToIP(ad.Bytes())
			multiError = concatError(multiError, err)
			info.KeyEncDst = &tmp
		case tcaTunnelKeyEncKeyID:
			tmp := ad.Uint32()
			info.KeyEncKeyID = &tmp
		case tcaTunnelKeyEncDstPort:
			tmp := ad.Uint16()
			info.KeyEncDstPort = &tmp
		case tcaTunnelKeyNoCSUM:
			tmp := ad.Uint8()
			info.KeyNoCSUM = &tmp
		case tcaTunnelKeyEncTOS:
			tmp := ad.Uint8()
			info.KeyEncTOS = &tmp
		case tcaTunnelKeyEncTTL:
			tmp := ad.Uint8()
			info.KeyEncTTL = &tmp
		case tcaTunnelKeyPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalTunnelKey()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
