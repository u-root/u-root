package tc

// IPSetDir defines the packet direction.
type IPSetDir uint8

// Various IP packet directions.
const (
	IPSetSrc = IPSetDir(1)
	IPSetDst = IPSetDir(2)
)

// IPSetMatch contains attributes of the ipset match discipline
type IPSetMatch struct {
	IPSetID uint16
	Dir     []IPSetDir
}

type ipsetMatch struct {
	ID    uint16
	Dim   uint8
	Flags uint8
}

func unmarshalIPSetMatch(data []byte, info *IPSetMatch) error {
	tmp := ipsetMatch{}
	if err := unmarshalStruct(data, &tmp); err != nil {
		return err
	}

	info.IPSetID = tmp.ID
	for i := uint8(1); i <= tmp.Dim; i++ {
		if (tmp.Flags & (1 << i)) == (1 << i) {
			info.Dir = append(info.Dir, IPSetSrc)
		} else {
			info.Dir = append(info.Dir, IPSetDst)
		}
	}
	return nil
}

func marshalIPSetMatch(info *IPSetMatch) ([]byte, error) {
	if info == nil {
		return []byte{}, ErrNoArg
	}

	tmp := ipsetMatch{
		ID: info.IPSetID,
	}

	if len(info.Dir) == 0 || len(info.Dir) > 3 {
		return nil, ErrInvalidArg
	}

	for _, dir := range info.Dir {
		tmp.Dim++
		if dir == IPSetSrc {
			tmp.Flags |= (1 << tmp.Dim)
		}
	}
	return marshalStruct(&tmp)
}
