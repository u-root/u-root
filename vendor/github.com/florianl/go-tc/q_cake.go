package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCakeUnspec = iota
	tcaCakePad
	tcaCakeBaseRate64
	tcaCakeDiffServMode
	tcaCakeAtm
	tcaCakeFlowMode
	tcaCakeOverhead
	tcaCakeRtt
	tcaCakeTarget
	tcaCakeAutorate
	tcaCakeMemory
	tcaCakeNat
	tcaCakeRaw
	tcaCakeWash
	tcaCakeMpu
	tcaCakeIngress
	tcaCakeAckFilter
	tcaCakeSplitGso
	tcaCakeFwMark
)

// Cake contains attributes of the cake discipline.
// http://man7.org/linux/man-pages/man8/tc-cake.8.html
type Cake struct {
	BaseRate     *uint64
	DiffServMode *uint32
	Atm          *uint32
	FlowMode     *uint32
	Overhead     *uint32
	Rtt          *uint32
	Target       *uint32
	Autorate     *uint32
	Memory       *uint32
	Nat          *uint32
	Raw          *uint32
	Wash         *uint32
	Mpu          *uint32
	Ingress      *uint32
	AckFilter    *uint32
	SplitGso     *uint32
	FwMark       *uint32
}

// unmarshalCake parses the Cake-encoded data and stores the result in the value pointed to by info.
func unmarshalCake(data []byte, info *Cake) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaCakeBaseRate64:
			tmp := ad.Uint64()
			info.BaseRate = &tmp
		case tcaCakeDiffServMode:
			tmp := ad.Uint32()
			info.DiffServMode = &tmp
		case tcaCakeAtm:
			tmp := ad.Uint32()
			info.Atm = &tmp
		case tcaCakeFlowMode:
			tmp := ad.Uint32()
			info.FlowMode = &tmp
		case tcaCakeOverhead:
			tmp := ad.Uint32()
			info.Overhead = &tmp
		case tcaCakeRtt:
			tmp := ad.Uint32()
			info.Rtt = &tmp
		case tcaCakeTarget:
			tmp := ad.Uint32()
			info.Target = &tmp
		case tcaCakeAutorate:
			tmp := ad.Uint32()
			info.Autorate = &tmp
		case tcaCakeMemory:
			tmp := ad.Uint32()
			info.Memory = &tmp
		case tcaCakeNat:
			tmp := ad.Uint32()
			info.Nat = &tmp
		case tcaCakeRaw:
			tmp := ad.Uint32()
			info.Raw = &tmp
		case tcaCakeWash:
			tmp := ad.Uint32()
			info.Wash = &tmp
		case tcaCakeMpu:
			tmp := ad.Uint32()
			info.Mpu = &tmp
		case tcaCakeIngress:
			tmp := ad.Uint32()
			info.Ingress = &tmp
		case tcaCakeAckFilter:
			tmp := ad.Uint32()
			info.AckFilter = &tmp
		case tcaCakeSplitGso:
			tmp := ad.Uint32()
			info.SplitGso = &tmp
		case tcaCakeFwMark:
			tmp := ad.Uint32()
			info.FwMark = &tmp
		case tcaCakePad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalCake()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalCake returns the binary encoding of Red
func marshalCake(info *Cake) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Cake: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.BaseRate != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeBaseRate64, Data: *info.BaseRate})
	}
	if info.DiffServMode != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeDiffServMode, Data: *info.DiffServMode})
	}
	if info.Atm != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeAtm, Data: *info.Atm})
	}
	if info.FlowMode != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeFlowMode, Data: *info.FlowMode})
	}
	if info.Overhead != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeOverhead, Data: *info.Overhead})
	}
	if info.Rtt != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeRtt, Data: *info.Rtt})
	}
	if info.Target != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTarget, Data: *info.Target})
	}
	if info.Autorate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeAutorate, Data: *info.Autorate})
	}
	if info.Memory != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeMemory, Data: *info.Memory})
	}
	if info.Nat != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeNat, Data: *info.Nat})
	}
	if info.Raw != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeRaw, Data: *info.Raw})
	}
	if info.Wash != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeWash, Data: *info.Wash})
	}
	if info.Mpu != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeMpu, Data: *info.Mpu})
	}
	if info.Ingress != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeIngress, Data: *info.Ingress})
	}
	if info.AckFilter != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeAckFilter, Data: *info.AckFilter})
	}
	if info.SplitGso != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeSplitGso, Data: *info.SplitGso})
	}
	if info.FwMark != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeFwMark, Data: *info.FwMark})
	}
	return marshalAttributes(options)
}
