package tc

import (
	"encoding/binary"
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

	// We can not unmarshal elements of a non-public struct.
	// So we extract the elements by hand.
	info.Offset = binary.LittleEndian.Uint16(data[:2])
	needleLen := binary.LittleEndian.Uint16(data[2:4])
	info.Layer = uint8(data[4])
	if len(data) < (8 + int(needleLen)) {
		return fmt.Errorf("unmarshalNByteMatch: incomplete needle")
	}
	info.Needle = data[8 : 8+needleLen]

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
