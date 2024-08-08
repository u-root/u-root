package tc

// ContainerMatch contains attributes of the container match discipline
type ContainerMatch struct {
	Pos uint32
}

func unmarshalContainerMatch(data []byte, info *ContainerMatch) error {
	if info == nil {
		return ErrNoArg
	}

	tmp := ContainerMatch{}
	if err := unmarshalStruct(data, &tmp); err != nil {
		return err
	}

	info.Pos = tmp.Pos
	return nil
}

func marshalContainerMatch(info *ContainerMatch) ([]byte, error) {
	if info == nil {
		return []byte{}, ErrNoArg
	}

	if info.Pos == 0 {
		return nil, ErrInvalidArg
	}

	return marshalStruct(info)
}
