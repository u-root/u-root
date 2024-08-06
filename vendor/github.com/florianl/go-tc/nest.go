package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

type valueType int

const (
	vtUint8 valueType = iota
	vtUint16
	vtUint32
	vtUint64
	vtString
	vtBytes
	vtFlag
	vtInt8
	vtInt16
	vtInt32
	vtInt64
	vtUint16Be
	vtUint32Be
	vtInt16Be
)

type tcOption struct {
	Interpretation valueType
	Type           uint16
	Data           interface{}
}

// NLA_F_NESTED from include/uapi/linux/netlink.h
const nlaFNnested = (1 << 15)

func marshalAttributes(options []tcOption) ([]byte, error) {
	if len(options) == 0 {
		return []byte{}, nil
	}
	var multiError error
	ad := netlink.NewAttributeEncoder()

	for _, option := range options {
		switch option.Interpretation {
		case vtUint8:
			ad.Uint8(option.Type, (option.Data).(uint8))
		case vtUint16:
			ad.Uint16(option.Type, (option.Data).(uint16))
		case vtUint32:
			ad.Uint32(option.Type, (option.Data).(uint32))
		case vtUint64:
			ad.Uint64(option.Type, (option.Data).(uint64))
		case vtString:
			ad.String(option.Type, (option.Data).(string))
		case vtBytes:
			ad.Bytes(option.Type, (option.Data).([]byte))
		case vtFlag:
			ad.Flag(option.Type, true)
		case vtInt8:
			data := bytes.NewBuffer(make([]byte, 0, 1))
			err := binary.Write(data, nativeEndian, (option.Data).(int8))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		case vtInt16:
			data := bytes.NewBuffer(make([]byte, 0, 2))
			err := binary.Write(data, nativeEndian, (option.Data).(int16))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		case vtInt32:
			data := bytes.NewBuffer(make([]byte, 0, 4))
			err := binary.Write(data, nativeEndian, (option.Data).(int32))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		case vtInt64:
			data := bytes.NewBuffer(make([]byte, 0, 8))
			err := binary.Write(data, nativeEndian, (option.Data).(int64))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		case vtUint16Be:
			data := bytes.NewBuffer(make([]byte, 0, 2))
			err := binary.Write(data, binary.BigEndian, (option.Data).(uint16))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		case vtUint32Be:
			data := bytes.NewBuffer(make([]byte, 0, 4))
			err := binary.Write(data, binary.BigEndian, (option.Data).(uint32))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		case vtInt16Be:
			data := bytes.NewBuffer(make([]byte, 0, 2))
			err := binary.Write(data, binary.BigEndian, (option.Data).(int16))
			multiError = concatError(multiError, err)
			ad.Bytes(option.Type, data.Bytes())
		default:
			multiError = fmt.Errorf("unknown interpretation (%d)", option.Interpretation)
		}
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return ad.Encode()
}

func unmarshalNetlinkAttribute(data []byte, val interface{}) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, nativeEndian, val)
	return err
}
