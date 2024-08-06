// The userspace part transforms the logic expressions into an array
// consisting of multiple sequences of interconnected ematches separated
// by markers. Precedence is implemented by a special ematch kind
// referencing a sequence beyond the marker of the current sequence
// causing the current position in the sequence to be pushed onto a stack
// to allow the current position to be overwritten by the position referenced
// in the special ematch. Matching continues in the new sequence until a
// marker is reached causing the position to be restored from the stack.
//
// Example:
//          A AND (B1 OR B2) AND C AND D
//
//              ------->-PUSH-------
//    -->--    /         -->--      \   -->--
//   /     \  /         /     \      \ /     \
// +-------+-------+-------+-------+-------+--------+
// | A AND | B AND | C AND | D END | B1 OR | B2 END |
// +-------+-------+-------+-------+-------+--------+
//                    \                      /
//                     --------<-POP---------
//
// where B is a virtual ematch referencing to sequence starting with B1. and B
// implemented with Container.
//
// When the skb input ematch module is used, the ematch match logic operates on
// the entire array. If the kernel finds the kind to be a container, it goes back
// to B until it reaches the end. This is implemented using recursive functions.
//
// The above is kernel logic.
//
// In userspace, the ematch array needs to be encapsulated. Logical combinations
// need to update flags and use containers. The updated ematch array would look like this:
// -------------------------------------------------------------------------------------------------
// index |      0       |        1        |      2       |     3        |     4       |     5      |
// ematch|      A       |        B        |      C       |     D        |     B1      |     B2     |
// kind  | EmatchIPSet  | EmatchContainer | EmatchIPT    | EmatchNByte  | EmatchCmp   | EmatchU32  |
// flags | EmatchRelAnd | EmatchRelAnd    | EmatchRelAnd | EmatchRelEND | EmatchRelOr |EmatchRelEnd|
// extend|      ...     |    Pos=4        |      ...     |    ...       |     ...     |     ...    |
// --------------------------------------------------------------------------------------------------
//
// last match order is：
// index： 0 --> 1 --> 4 --> 5 --> 2 --> 3

package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

// EmatchLayer defines the layer the match will be applied upon.
type EmatchLayer uint8

// Various Ematch network layers.
const (
	EmatchLayerLink      = EmatchLayer(0)
	EmatchLayerNetwork   = EmatchLayer(1)
	EmatchLayerTransport = EmatchLayer(2)
)

// EmatchOpnd defines how matches are concatenated.
type EmatchOpnd uint8

// Various Ematch operands
const (
	EmatchOpndEq = EmatchOpnd(0)
	EmatchOpndGt = EmatchOpnd(1)
	EmatchOpndLt = EmatchOpnd(2)
)

const (
	tcaEmatchTreeUnspec = iota
	tcaEmatchTreeHdr
	tcaEmatchTreeList
)

// EmatchKind defines the matching module.
type EmatchKind uint16

// Various Ematch kinds
const (
	EmatchContainer = EmatchKind(iota)
	EmatchCmp
	EmatchNByte
	EmatchU32
	EmatchMeta
	EmatchText
	EmatchVLan
	EmatchCanID
	EmatchIPSet
	EmatchIPT
	ematchInvalid
)

// Various Ematch flags
const (
	EmatchRelEnd uint16 = 0
	EmatchRelAnd uint16 = 1 << (iota - 1)
	EmatchRelOr
	EmatchInvert
	EmatchSimple
)

// Ematch contains attributes of the ematch discipline
// https://man7.org/linux/man-pages/man8/tc-ematch.8.html
type Ematch struct {
	Hdr     *EmatchTreeHdr
	Matches *[]EmatchMatch
}

// EmatchTreeHdr from tcf_ematch_tree_hdr in include/uapi/linux/pkt_cls.h
type EmatchTreeHdr struct {
	NMatches uint16
	ProgID   uint16
}

// EmatchHdr from tcf_ematch_hdr in include/uapi/linux/pkt_cls.h
type EmatchHdr struct {
	MatchID uint16
	Kind    EmatchKind
	Flags   uint16
	Pad     uint16
}

// EmatchMatch contains attributes of the ematch discipline
type EmatchMatch struct {
	Hdr            EmatchHdr
	U32Match       *U32Match
	CmpMatch       *CmpMatch
	IPSetMatch     *IPSetMatch
	IptMatch       *IptMatch
	ContainerMatch *ContainerMatch
	NByteMatch     *NByteMatch
}

// unmarshalEmatch parses the Ematch-encoded data and stores the result in the value pointed to by info.
func unmarshalEmatch(data []byte, info *Ematch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaEmatchTreeHdr:
			hdr := &EmatchTreeHdr{}
			err := unmarshalStruct(ad.Bytes(), hdr)
			multiError = concatError(multiError, err)
			info.Hdr = hdr
		case tcaEmatchTreeList:
			list := []EmatchMatch{}
			err := unmarshalEmatchTreeList(ad.Bytes(), &list)
			multiError = concatError(multiError, err)
			info.Matches = &list
		default:
			return fmt.Errorf("UnmarshalEmatch()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalEmatch returns the binary encoding of Ematch
func marshalEmatch(info *Ematch) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ematch: %w", ErrNoArg)
	}
	var multiError error

	if info.Hdr != nil {
		data, err := marshalStruct(info.Hdr)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmatchTreeHdr, Data: data})
	}
	if info.Matches != nil {
		data, err := marshalEmatchTreeList(info.Matches)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmatchTreeList | nlaFNnested, Data: data})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}

func unmarshalEmatchTreeList(data []byte, info *[]EmatchMatch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		match := EmatchMatch{}
		tmp := ad.Bytes()
		if err := unmarshalStruct(tmp[:8], &match.Hdr); err != nil {
			return err
		}
		switch match.Hdr.Kind {
		case EmatchU32:
			expr := &U32Match{}
			err := unmarshalU32Match(tmp[8:], expr)
			multiError = concatError(multiError, err)
			match.U32Match = expr
		case EmatchCmp:
			expr := &CmpMatch{}
			err := unmarshalCmpMatch(tmp[8:], expr)
			multiError = concatError(multiError, err)
			match.CmpMatch = expr
		case EmatchIPSet:
			expr := &IPSetMatch{}
			err := unmarshalIPSetMatch(tmp[8:], expr)
			multiError = concatError(multiError, err)
			match.IPSetMatch = expr
		case EmatchIPT:
			expr := &IptMatch{}
			err := unmarshalIptMatch(tmp[8:], expr)
			multiError = concatError(multiError, err)
			match.IptMatch = expr
		case EmatchContainer:
			expr := &ContainerMatch{}
			err := unmarshalContainerMatch(tmp[8:], expr)
			multiError = concatError(multiError, err)
			match.ContainerMatch = expr
		case EmatchNByte:
			expr := &NByteMatch{}
			err := unmarshalNByteMatch(tmp[8:], expr)
			multiError = concatError(multiError, err)
			match.NByteMatch = expr
		default:
			return fmt.Errorf("unmarshalEmatchTreeList() kind %d is not yet implemented", match.Hdr.Kind)
		}
		*info = append(*info, match)
	}
	return concatError(multiError, ad.Err())
}

func marshalEmatchTreeList(info *[]EmatchMatch) ([]byte, error) {
	options := []tcOption{}

	for i, m := range *info {
		payload, err := marshalStruct(m.Hdr)
		if err != nil {
			return []byte{}, err
		}
		var expr []byte
		switch m.Hdr.Kind {
		case EmatchU32:
			expr, err = marshalU32Match(m.U32Match)
		case EmatchCmp:
			expr, err = marshalCmpMatch(m.CmpMatch)
		case EmatchIPSet:
			expr, err = marshalIPSetMatch(m.IPSetMatch)
		case EmatchIPT:
			expr, err = marshalIptMatch(m.IptMatch)
		case EmatchContainer:
			expr, err = marshalContainerMatch(m.ContainerMatch)
		case EmatchNByte:
			expr, err = marshalNByteMatch(m.NByteMatch)
		default:
			return []byte{}, fmt.Errorf("marshalEmatchTreeList() kind %d is not yet implemented", m.Hdr.Kind)
		}
		if err != nil {
			return []byte{}, fmt.Errorf("marshalEmatchTreeList(): %v", err)
		}
		payload = append(payload, expr...)
		options = append(options, tcOption{Interpretation: vtBytes, Type: uint16(i + 1), Data: payload})
	}
	return marshalAttributes(options)
}
