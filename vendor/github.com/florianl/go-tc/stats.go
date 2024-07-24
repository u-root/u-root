package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaStatsUnspec = iota
	tcaStatsBasic
	tcaStatsRateEst
	tcaStatsQueue
	tcaStatsApp
	tcaStatsRateEst64
	tcaStatsPad
	tcaStatsBasicHw
	tcaStatsPkt64
)

// GenStats from include/uapi/linux/gen_stats.h
type GenStats struct {
	Basic     *GenBasic
	RateEst   *GenRateEst
	Queue     *GenQueue
	RateEst64 *GenRateEst64
	BasicHw   *GenBasic
}

// GenBasic from include/uapi/linux/gen_stats.h
type GenBasic struct {
	Bytes   uint64
	Packets uint32
}

// GenRateEst from include/uapi/linux/gen_stats.h
type GenRateEst struct {
	BytePerSecond   uint32
	PacketPerSecond uint32
}

// GenRateEst64 from include/uapi/linux/gen_stats.h
type GenRateEst64 struct {
	BytePerSecond   uint64
	PacketPerSecond uint64
}

// GenQueue from include/uapi/linux/gen_stats.h
type GenQueue struct {
	QueueLen   uint32
	Backlog    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

// unmarshalGenStats parses the Pie-encoded data and stores the result in the value pointed to by info.
func unmarshalGenStats(data []byte, info *GenStats) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaStatsBasic:
			stat := &GenBasic{}
			err = unmarshalStruct(ad.Bytes(), stat)
			multiError = concatError(multiError, err)
			info.Basic = stat
		case tcaStatsRateEst:
			stat := &GenRateEst{}
			err = unmarshalStruct(ad.Bytes(), stat)
			multiError = concatError(multiError, err)
			info.RateEst = stat
		case tcaStatsQueue:
			stat := &GenQueue{}
			err = unmarshalStruct(ad.Bytes(), stat)
			multiError = concatError(multiError, err)
			info.Queue = stat
		case tcaStatsRateEst64:
			stat := &GenRateEst64{}
			err = unmarshalStruct(ad.Bytes(), stat)
			multiError = concatError(multiError, err)
			info.RateEst64 = stat
		case tcaStatsBasicHw:
			stat := &GenBasic{}
			err = unmarshalStruct(ad.Bytes(), stat)
			multiError = concatError(multiError, err)
			info.BasicHw = stat
		case tcaStatsPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalGenStats()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalGenStats returns the binary encoding of GenStats
func marshalGenStats(info *GenStats) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("GenStats: %w", ErrNoArg)
	}

	var multiError error

	if info.Basic != nil {
		data, err := marshalStruct(info.Basic)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStatsBasic, Data: data})
	}
	if info.RateEst != nil {
		data, err := marshalStruct(info.RateEst)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStatsRateEst, Data: data})
	}
	if info.Queue != nil {
		data, err := marshalStruct(info.Queue)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStatsQueue, Data: data})
	}
	if info.RateEst64 != nil {
		data, err := marshalStruct(info.RateEst64)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStatsRateEst64, Data: data})
	}
	if info.BasicHw != nil {
		data, err := marshalStruct(info.BasicHw)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStatsBasicHw, Data: data})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}
