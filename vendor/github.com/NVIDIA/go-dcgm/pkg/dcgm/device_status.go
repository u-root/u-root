package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
)

type PerfState uint

const (
	PerfStateMax     = 0
	PerfStateMin     = 15
	PerfStateUnknown = 32
)

func (p PerfState) String() string {
	if p >= PerfStateMax && p <= PerfStateMin {
		return fmt.Sprintf("P%d", p)
	}
	return "Unknown"
}

type UtilizationInfo struct {
	GPU     int64 // %
	Memory  int64 // %
	Encoder int64 // %
	Decoder int64 // %
}

type ECCErrorsInfo struct {
	SingleBit int64
	DoubleBit int64
}

type MemoryInfo struct {
	GlobalUsed int64
	ECCErrors  ECCErrorsInfo
}

type ClockInfo struct {
	Cores  int64 // MHz
	Memory int64 // MHz
}

type PCIThroughputInfo struct {
	Rx      int64 // MB
	Tx      int64 // MB
	Replays int64
}

type PCIStatusInfo struct {
	BAR1Used   int64 // MB
	Throughput PCIThroughputInfo
	FBUsed     int64
}

type DeviceStatus struct {
	Power       float64 // W
	Temperature int64   // °C
	Utilization UtilizationInfo
	Memory      MemoryInfo
	Clocks      ClockInfo
	PCI         PCIStatusInfo
	Performance PerfState
	FanSpeed    int64 // %
}

func latestValuesForDevice(gpuId uint) (status DeviceStatus, err error) {
	const (
		pwr int = iota
		temp
		sm
		mem
		enc
		dec
		smClock
		memClock
		bar1Used
		pcieRxThroughput
		pcieTxThroughput
		pcieReplay
		fbUsed
		sbe
		dbe
		pstate
		fanSpeed
		fieldsCount
	)

	deviceFields := make([]Short, fieldsCount)
	deviceFields[pwr] = C.DCGM_FI_DEV_POWER_USAGE
	deviceFields[temp] = C.DCGM_FI_DEV_GPU_TEMP
	deviceFields[sm] = C.DCGM_FI_DEV_GPU_UTIL
	deviceFields[mem] = C.DCGM_FI_DEV_MEM_COPY_UTIL
	deviceFields[enc] = C.DCGM_FI_DEV_ENC_UTIL
	deviceFields[dec] = C.DCGM_FI_DEV_DEC_UTIL
	deviceFields[smClock] = C.DCGM_FI_DEV_SM_CLOCK
	deviceFields[memClock] = C.DCGM_FI_DEV_MEM_CLOCK
	deviceFields[bar1Used] = C.DCGM_FI_DEV_BAR1_USED
	deviceFields[pcieRxThroughput] = C.DCGM_FI_DEV_PCIE_RX_THROUGHPUT
	deviceFields[pcieTxThroughput] = C.DCGM_FI_DEV_PCIE_TX_THROUGHPUT
	deviceFields[pcieReplay] = C.DCGM_FI_DEV_PCIE_REPLAY_COUNTER
	deviceFields[fbUsed] = C.DCGM_FI_DEV_FB_USED
	deviceFields[sbe] = C.DCGM_FI_DEV_ECC_SBE_AGG_TOTAL
	deviceFields[dbe] = C.DCGM_FI_DEV_ECC_DBE_AGG_TOTAL
	deviceFields[pstate] = C.DCGM_FI_DEV_PSTATE
	deviceFields[fanSpeed] = C.DCGM_FI_DEV_FAN_SPEED

	fieldsName := fmt.Sprintf("devStatusFields%d", rand.Uint64())
	fieldsId, err := FieldGroupCreate(fieldsName, deviceFields)
	if err != nil {
		return
	}

	groupName := fmt.Sprintf("devStatus%d", rand.Uint64())
	groupId, err := WatchFields(gpuId, fieldsId, groupName)
	if err != nil {
		_ = FieldGroupDestroy(fieldsId)
		return
	}

	values, err := GetLatestValuesForFields(gpuId, deviceFields)
	if err != nil {
		_ = FieldGroupDestroy(fieldsId)
		_ = DestroyGroup(groupId)
		return status, err
	}

	power := values[pwr].Float64()

	gpuUtil := UtilizationInfo{
		GPU:     values[sm].Int64(),
		Memory:  values[mem].Int64(),
		Encoder: values[enc].Int64(),
		Decoder: values[dec].Int64(),
	}

	memory := MemoryInfo{
		ECCErrors: ECCErrorsInfo{
			SingleBit: values[sbe].Int64(),
			DoubleBit: values[dbe].Int64(),
		},
	}

	clocks := ClockInfo{
		Cores:  values[smClock].Int64(),
		Memory: values[memClock].Int64(),
	}

	pci := PCIStatusInfo{
		BAR1Used: values[bar1Used].Int64(),
		Throughput: PCIThroughputInfo{
			Rx:      values[pcieRxThroughput].Int64(),
			Tx:      values[pcieTxThroughput].Int64(),
			Replays: values[pcieReplay].Int64(),
		},
		FBUsed: values[fbUsed].Int64(),
	}

	status = DeviceStatus{
		Power:       power,
		Temperature: values[temp].Int64(),
		Utilization: gpuUtil,
		Memory:      memory,
		Clocks:      clocks,
		PCI:         pci,
		Performance: PerfState(values[pstate].Int64()),
		FanSpeed:    values[fanSpeed].Int64(),
	}

	_ = FieldGroupDestroy(fieldsId)
	_ = DestroyGroup(groupId)
	return
}
