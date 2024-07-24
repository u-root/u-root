package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaBpfUnspec = iota
	tcaBpfAct
	tcaBpfPolice
	tcaBpfClassID
	tcaBpfOpsLen
	tcaBpfOps
	tcaBpfFd
	tcaBpfName
	tcaBpfFlags
	tcaBpfFlagsGen
	tcaBpfTag
	tcaBpfID
)

// Bpf contains attributes of the bpf discipline
type Bpf struct {
	Action   *Action
	Police   *Police
	ClassID  *uint32
	OpsLen   *uint16
	Ops      *[]byte
	FD       *uint32
	Name     *string
	Flags    *uint32
	FlagsGen *uint32
	Tag      *[]byte
	ID       *uint32
}

// Flags defined by the kernel for the BPF filter
const (
	BpfActDirect = 1
)

// unmarshalBpf parses the Bpf-encoded data and stores the result in the value pointed to by info.
func unmarshalBpf(data []byte, info *Bpf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaBpfAct:
			actions := &[]*Action{}
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
			info.Action = (*actions)[0]
		case tcaBpfPolice:
			pol := &Police{}
			err := unmarshalPolice(ad.Bytes(), pol)
			multiError = concatError(multiError, err)
			info.Police = pol
		case tcaBpfClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaBpfOpsLen:
			info.OpsLen = uint16Ptr(ad.Uint16())
		case tcaBpfOps:
			info.Ops = bytesPtr(ad.Bytes())
		case tcaBpfFd:
			info.FD = uint32Ptr(ad.Uint32())
		case tcaBpfName:
			info.Name = stringPtr(ad.String())
		case tcaBpfFlags:
			info.Flags = uint32Ptr(ad.Uint32())
		case tcaBpfFlagsGen:
			info.FlagsGen = uint32Ptr(ad.Uint32())
		case tcaBpfTag:
			info.Tag = bytesPtr(ad.Bytes())
		case tcaBpfID:
			info.ID = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalBpf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalBpf returns the binary encoding of Bpf
func marshalBpf(info *Bpf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Bpf: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Ops != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfOps, Data: bytesValue(info.Ops)})
	}
	if info.OpsLen != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaBpfOpsLen, Data: uint16Value(info.OpsLen)})
	}
	if info.FD != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFd, Data: uint32Value(info.FD)})
	}
	if info.Name != nil {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaBpfName, Data: stringValue(info.Name)})
	}
	if info.ID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfID, Data: uint32Value(info.ID)})
	}
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfClassID, Data: uint32Value(info.ClassID)})
	}
	if info.Tag != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfTag, Data: bytesValue(info.Tag)})
	}
	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlags, Data: uint32Value(info.Flags)})
	}
	if info.FlagsGen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlagsGen, Data: uint32Value(info.FlagsGen)})
	}
	if info.Action != nil {
		actions := []*Action{info.Action}
		data, err := marshalActions(0, actions)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfAct, Data: data})
	}
	return marshalAttributes(options)
}
