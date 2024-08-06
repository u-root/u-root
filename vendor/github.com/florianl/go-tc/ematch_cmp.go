package tc

// CmpMatchAlign defines byte alignments.
type CmpMatchAlign uint8

// Various byte alignments.
const (
	CmpMatchU8  = CmpMatchAlign(1)
	CmpMatchU16 = CmpMatchAlign(2)
	CmpMatchU32 = CmpMatchAlign(4)
)

// CmpMatchFlag defines available flags.
type CmpMatchFlag uint8

// Various flags.
const (
	CmpMatchTrans = CmpMatchFlag(1)
)

// CmpMatch contains attributes of the cmp match discipline
type CmpMatch struct {
	Val   uint32
	Mask  uint32
	Off   uint16
	Align CmpMatchAlign
	Flags CmpMatchFlag
	Layer EmatchLayer
	Opnd  EmatchOpnd
}

type cmpMatch struct {
	Val  uint32
	Mask uint32
	Off  uint16
	Opts uint16
}

func unmarshalCmpMatch(data []byte, info *CmpMatch) error {
	tmp := cmpMatch{}
	if err := unmarshalStruct(data, &tmp); err != nil {
		return err
	}
	info.Val = tmp.Val
	info.Mask = tmp.Mask
	info.Off = tmp.Off
	info.Align = CmpMatchAlign((tmp.Opts) & 0xf)
	info.Flags = CmpMatchFlag((tmp.Opts >> 4) & 0xf)
	info.Layer = EmatchLayer((tmp.Opts >> 8) & 0xf)
	info.Opnd = EmatchOpnd((tmp.Opts >> 12) & 0xf)
	return nil
}

func marshalCmpMatch(info *CmpMatch) ([]byte, error) {
	if info == nil {
		return []byte{}, ErrNoArg
	}
	var opts uint16
	opts |= uint16(info.Align & 0xf)
	opts |= uint16(info.Flags&0xf) << 4
	opts |= uint16(info.Layer&0xf) << 8
	opts |= uint16(info.Opnd&0xf) << 12

	tmp := cmpMatch{
		Val:  info.Val,
		Mask: info.Mask,
		Off:  info.Off,
		Opts: opts,
	}
	return marshalStruct(&tmp)
}
