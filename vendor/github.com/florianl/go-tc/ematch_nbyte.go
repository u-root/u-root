package tc

import (
	"fmt"
)

// NByteMatch contains attributes of the Nbyte match discipline
type NByteMatch struct {
	Offset uint16
	Layer  uint8
	Needle []byte
}

type tcfEmNByte struct {
	off   uint16
	len   uint16
	layer uint8
}

func unmarshalNByteMatch(data []byte, info *NByteMatch) error {
	if len(data) < 8 {
		return fmt.Errorf("unmarshalNByteMatch: incomplete data")
	}

	nbyte := tcfEmNByte{}
	if err := unmarshalStruct(data[:5], nbyte); err != nil {
		return err
	}
	info.Offset = nbyte.off
	info.Layer = nbyte.layer
	info.Needle = data[8:]

	return nil
}

func marshalNByteMatch(info *NByteMatch) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("marshalNByteMatch: %w", ErrNoArg)
	}
	nbyte := tcfEmNByte{
		off:   info.Offset,
		len:   uint16(len(info.Needle)),
		layer: info.Layer,
	}

	tmp, err := marshalAndAlignStruct(nbyte)
	if err != nil {
		return []byte{}, err
	}
	tmp = append(tmp, info.Needle...)
	return tmp, nil
}
