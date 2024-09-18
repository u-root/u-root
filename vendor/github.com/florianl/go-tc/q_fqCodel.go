package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFqCodelUnspec = iota
	tcaFqCodelTarget
	tcaFqCodelLimit
	tcaFqCodelInterval
	tcaFqCodelEcn
	tcaFqCodelFlows
	tcaFqCodelQuantum
	tcaFqCodelCeThreshold
	tcaFqCodelDropBatchSize
	tcaFqCodelMemoryLimit
)

const (
	tcaFqCodelXStatsQdisc = iota
	tcaFqCodelXStatsClass
)

// FqCodel contains attributes of the fq_codel discipline
type FqCodel struct {
	Target        *uint32
	Limit         *uint32
	Interval      *uint32
	ECN           *uint32
	Flows         *uint32
	Quantum       *uint32
	CEThreshold   *uint32
	DropBatchSize *uint32
	MemoryLimit   *uint32
}

// marshalFqCodel returns the binary encoding of FqCodel
func marshalFqCodel(info *FqCodel) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("FqCodel: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Target != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelTarget, Data: uint32Value(info.Target)})
	}

	if info.Limit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelLimit, Data: uint32Value(info.Limit)})
	}

	if info.Interval != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelInterval, Data: uint32Value(info.Interval)})
	}

	if info.ECN != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelEcn, Data: uint32Value(info.ECN)})
	}

	if info.Flows != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelFlows, Data: uint32Value(info.Flows)})
	}

	if info.Quantum != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelQuantum, Data: uint32Value(info.Quantum)})
	}

	if info.CEThreshold != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelCeThreshold, Data: uint32Value(info.CEThreshold)})
	}

	if info.DropBatchSize != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelDropBatchSize, Data: uint32Value(info.DropBatchSize)})
	}

	if info.MemoryLimit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelMemoryLimit, Data: uint32Value(info.MemoryLimit)})
	}

	return marshalAttributes(options)
}

// unmarshalFqCodel parses the FqCodel-encoded data and stores the result in the value pointed to by info.
func unmarshalFqCodel(data []byte, info *FqCodel) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaFqCodelTarget:
			info.Target = uint32Ptr(ad.Uint32())
		case tcaFqCodelLimit:
			info.Limit = uint32Ptr(ad.Uint32())
		case tcaFqCodelInterval:
			info.Interval = uint32Ptr(ad.Uint32())
		case tcaFqCodelEcn:
			info.ECN = uint32Ptr(ad.Uint32())
		case tcaFqCodelFlows:
			info.Flows = uint32Ptr(ad.Uint32())
		case tcaFqCodelQuantum:
			info.Quantum = uint32Ptr(ad.Uint32())
		case tcaFqCodelCeThreshold:
			info.CEThreshold = uint32Ptr(ad.Uint32())
		case tcaFqCodelDropBatchSize:
			info.DropBatchSize = uint32Ptr(ad.Uint32())
		case tcaFqCodelMemoryLimit:
			info.MemoryLimit = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalFqCodel()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}
