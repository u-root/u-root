package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCgroupUnspec = iota
	tcaCgroupAct
	tcaCgroupPolice
	tcaCgroupEmatches
)

// Cgroup contains attributes of the cgroup discipline
type Cgroup struct {
	Action *Action
	Ematch *Ematch
}

// marshalCgroup returns the binary encoding of Cgroup
func marshalCgroup(info *Cgroup) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Cgroup: %w", ErrNoArg)
	}
	var multiError error
	// TODO: improve logic and check combinations
	if info.Action != nil {
		data, err := marshalAction(0, info.Action, tcaActOptions)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCgroupAct, Data: data})

	}
	if info.Ematch != nil {
		data, err := marshalEmatch(info.Ematch)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCgroupEmatches, Data: data})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}

// unmarshalCgroup parses the Cgroup-encoded data and stores the result in the value pointed to by info.
func unmarshalCgroup(data []byte, info *Cgroup) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaCgroupAct:
			act := &Action{}
			err := unmarshalAction(ad.Bytes(), act)
			multiError = concatError(multiError, err)
			info.Action = act
		case tcaCgroupEmatches:
			ematch := &Ematch{}
			err := unmarshalEmatch(ad.Bytes(), ematch)
			multiError = concatError(multiError, err)
			info.Ematch = ematch
		default:
			return fmt.Errorf("unmarshalCgroup()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
