package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFqUnspec = iota
	tcaFqPLimit
	tcaFqFlowPLimit
	tcaFqQuantum
	tcaFqInitQuantum
	tcaFqRateEnable
	tcaFqFlowDefaultRate
	tcaFqFlowMaxRate
	tcaFqBucketsLog
	tcaFqFlowRefillDelay
	tcaFqOrphanMask
	tcaFqLowRateThreshold
	tcaFqCEThreshold
	tcaFqTimerSlack
	tcaFqHorizon
	tcaFqHorizonDrop
	tcaFqPrioMap
	tcaFqWeights
)

// FqPrioQopt according to tc_prio_qopt in /include/uapi/linux/pkt_sched.h
type FqPrioQopt struct {
	Bands   int32
	PrioMap [16]uint8 // TC_PRIO_MAX + 1 = 16
}

// Fq contains attributes of the fq discipline
type Fq struct {
	PLimit           *uint32
	FlowPLimit       *uint32
	Quantum          *uint32
	InitQuantum      *uint32
	RateEnable       *uint32
	FlowDefaultRate  *uint32
	FlowMaxRate      *uint32
	BucketsLog       *uint32
	FlowRefillDelay  *uint32
	OrphanMask       *uint32
	LowRateThreshold *uint32
	CEThreshold      *uint32
	TimerSlack       *uint32
	Horizon          *uint32
	HorizonDrop      *uint8
	PrioMap          *FqPrioQopt
	Weights          *[]int32
}

// unmarshalFq parses the Fq-encoded data and stores the result in the value pointed to by info.
func unmarshalFq(data []byte, info *Fq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaFqPLimit:
			info.PLimit = uint32Ptr(ad.Uint32())
		case tcaFqFlowPLimit:
			info.FlowPLimit = uint32Ptr(ad.Uint32())
		case tcaFqQuantum:
			info.Quantum = uint32Ptr(ad.Uint32())
		case tcaFqInitQuantum:
			info.InitQuantum = uint32Ptr(ad.Uint32())
		case tcaFqRateEnable:
			info.RateEnable = uint32Ptr(ad.Uint32())
		case tcaFqFlowDefaultRate:
			info.FlowDefaultRate = uint32Ptr(ad.Uint32())
		case tcaFqFlowMaxRate:
			info.FlowMaxRate = uint32Ptr(ad.Uint32())
		case tcaFqBucketsLog:
			info.BucketsLog = uint32Ptr(ad.Uint32())
		case tcaFqFlowRefillDelay:
			info.FlowRefillDelay = uint32Ptr(ad.Uint32())
		case tcaFqOrphanMask:
			info.OrphanMask = uint32Ptr(ad.Uint32())
		case tcaFqLowRateThreshold:
			info.LowRateThreshold = uint32Ptr(ad.Uint32())
		case tcaFqCEThreshold:
			info.CEThreshold = uint32Ptr(ad.Uint32())
		case tcaFqTimerSlack:
			info.TimerSlack = uint32Ptr(ad.Uint32())
		case tcaFqHorizon:
			info.Horizon = uint32Ptr(ad.Uint32())
		case tcaFqHorizonDrop:
			info.HorizonDrop = uint8Ptr(ad.Uint8())
		case tcaFqPrioMap:
			priomap := &FqPrioQopt{}
			err := unmarshalStruct(ad.Bytes(), priomap)
			multiError = concatError(multiError, err)
			info.PrioMap = priomap
		case tcaFqWeights:
			size := len(ad.Bytes()) / 4
			weights := make([]int32, size)
			reader := bytes.NewReader(ad.Bytes())
			err := binary.Read(reader, nativeEndian, weights)
			multiError = concatError(multiError, err)
			info.Weights = &weights
		default:
			return fmt.Errorf("unmarshalFq()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalFq returns the binary encoding of Fq
func marshalFq(info *Fq) ([]byte, error) {
	options := []tcOption{}
	var multiError error

	if info == nil {
		return []byte{}, fmt.Errorf("Fq: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.PLimit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqPLimit, Data: uint32Value(info.PLimit)})
	}
	if info.FlowPLimit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowPLimit, Data: uint32Value(info.FlowPLimit)})
	}
	if info.Quantum != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqQuantum, Data: uint32Value(info.Quantum)})
	}
	if info.InitQuantum != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqInitQuantum, Data: uint32Value(info.InitQuantum)})
	}
	if info.RateEnable != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqRateEnable, Data: uint32Value(info.RateEnable)})
	}
	if info.FlowDefaultRate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowDefaultRate, Data: uint32Value(info.FlowDefaultRate)})
	}
	if info.FlowMaxRate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowMaxRate, Data: uint32Value(info.FlowMaxRate)})
	}
	if info.BucketsLog != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqBucketsLog, Data: uint32Value(info.BucketsLog)})
	}
	if info.FlowRefillDelay != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowRefillDelay, Data: uint32Value(info.FlowRefillDelay)})
	}
	if info.OrphanMask != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqOrphanMask, Data: uint32Value(info.OrphanMask)})
	}
	if info.LowRateThreshold != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqLowRateThreshold, Data: uint32Value(info.LowRateThreshold)})
	}
	if info.CEThreshold != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCEThreshold, Data: uint32Value(info.CEThreshold)})
	}
	if info.TimerSlack != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqTimerSlack, Data: uint32Value(info.TimerSlack)})
	}
	if info.Horizon != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqHorizon, Data: uint32Value(info.Horizon)})
	}
	if info.HorizonDrop != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaFqHorizonDrop, Data: uint8Value(info.HorizonDrop)})
	}
	if info.PrioMap != nil {
		data, err := marshalStruct(info.PrioMap)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFqPrioMap, Data: data})
	}
	if info.Weights != nil {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, nativeEndian, *info.Weights)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFqWeights, Data: buf.Bytes()})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}
