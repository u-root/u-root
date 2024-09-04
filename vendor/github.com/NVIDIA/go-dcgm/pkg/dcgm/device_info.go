package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
	"unsafe"

	"github.com/bits-and-blooms/bitset"
)

type PCIInfo struct {
	BusID     string
	BAR1      uint  // MB
	FBTotal   uint  // MB
	Bandwidth int64 // MB/s
}

type DeviceIdentifiers struct {
	Brand               string
	Model               string
	Serial              string
	Vbios               string
	InforomImageVersion string
	DriverVersion       string
}

type Device struct {
	GPU           uint
	DCGMSupported string
	UUID          string
	Power         uint // W
	PCI           PCIInfo
	Identifiers   DeviceIdentifiers
	Topology      []P2PLink
	CPUAffinity   string
}

// getAllDeviceCount counts all GPUs on the system
func getAllDeviceCount() (gpuCount uint, err error) {
	var gpuIdList [C.DCGM_MAX_NUM_DEVICES]C.uint
	var count C.int

	result := C.dcgmGetAllDevices(handle.handle, &gpuIdList[0], &count)
	if err = errorString(result); err != nil {
		return gpuCount, fmt.Errorf("Error getting devices count: %s", err)
	}
	gpuCount = uint(count)
	return
}

// getAllDeviceCount counts all GPUs on the system
func getEntityGroupEntities(entityGroup Field_Entity_Group) (entities []uint, err error) {
	var pEntities [C.DCGM_MAX_NUM_DEVICES]C.uint
	var count C.int = C.DCGM_MAX_NUM_DEVICES

	result := C.dcgmGetEntityGroupEntities(handle.handle, C.dcgm_field_entity_group_t(entityGroup), &pEntities[0], &count, 0)
	if err = errorString(result); err != nil {
		return nil, fmt.Errorf("Error getting entity count: %s", err)
	}

	for i := 0; i < int(count); i++ {
		entities = append(entities, uint(pEntities[i]))
	}
	return entities, nil
}

// getSupportedDevices returns DCGM supported GPUs
func getSupportedDevices() (gpus []uint, err error) {
	var gpuIdList [C.DCGM_MAX_NUM_DEVICES]C.uint
	var count C.int

	result := C.dcgmGetAllSupportedDevices(handle.handle, &gpuIdList[0], &count)
	if err = errorString(result); err != nil {
		return gpus, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	numGpus := int(count)
	gpus = make([]uint, numGpus)
	for i := 0; i < numGpus; i++ {
		gpus[i] = uint(gpuIdList[i])
	}
	return
}

func getPciBandwidth(gpuId uint) (int64, error) {
	const (
		maxLinkGen int = iota
		maxLinkWidth
		fieldsCount
	)

	pciFields := make([]Short, fieldsCount)
	pciFields[maxLinkGen] = C.DCGM_FI_DEV_PCIE_MAX_LINK_GEN
	pciFields[maxLinkWidth] = C.DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH

	fieldsName := fmt.Sprintf("pciBandwidthFields%d", rand.Uint64())

	fieldsId, err := FieldGroupCreate(fieldsName, pciFields)
	if err != nil {
		return 0, err
	}

	groupName := fmt.Sprintf("pciBandwidth%d", rand.Uint64())
	groupId, err := WatchFields(gpuId, fieldsId, groupName)
	if err != nil {
		_ = FieldGroupDestroy(fieldsId)
		return 0, err
	}

	values, err := GetLatestValuesForFields(gpuId, pciFields)
	if err != nil {
		_ = FieldGroupDestroy(fieldsId)
		_ = DestroyGroup(groupId)
		return 0, fmt.Errorf("Error getting Pcie bandwidth: %s", err)
	}

	gen := values[maxLinkGen].Int64()
	width := values[maxLinkWidth].Int64()

	_ = FieldGroupDestroy(fieldsId)
	_ = DestroyGroup(groupId)

	genMap := map[int64]int64{
		1: 250, // MB/s
		2: 500,
		3: 985,
		4: 1969,
	}

	bandwidth := genMap[gen] * width
	return bandwidth, nil
}

func getCPUAffinity(gpuId uint) (string, error) {
	const (
		affinity0 int = iota
		affinity1
		affinity2
		affinity3
		fieldsCount
	)

	affFields := make([]Short, fieldsCount)
	affFields[affinity0] = C.DCGM_FI_DEV_CPU_AFFINITY_0
	affFields[affinity1] = C.DCGM_FI_DEV_CPU_AFFINITY_1
	affFields[affinity2] = C.DCGM_FI_DEV_CPU_AFFINITY_2
	affFields[affinity3] = C.DCGM_FI_DEV_CPU_AFFINITY_3

	fieldsName := fmt.Sprintf("cpuAffFields%d", rand.Uint64())

	fieldsId, err := FieldGroupCreate(fieldsName, affFields)
	if err != nil {
		return "N/A", err
	}
	defer FieldGroupDestroy(fieldsId)

	groupName := fmt.Sprintf("cpuAff%d", rand.Uint64())
	groupId, err := WatchFields(gpuId, fieldsId, groupName)
	if err != nil {
		return "N/A", err
	}
	defer DestroyGroup(groupId)

	values, err := GetLatestValuesForFields(gpuId, affFields)
	if err != nil {
		return "N/A", fmt.Errorf("Error getting cpu affinity: %s", err)
	}

	bits := make([]uint64, 4)
	bits[0] = uint64(values[affinity0].Int64())
	bits[1] = uint64(values[affinity1].Int64())
	bits[2] = uint64(values[affinity2].Int64())
	bits[3] = uint64(values[affinity3].Int64())

	b := bitset.From(bits)

	return b.String(), nil
}

func getDeviceInfo(gpuid uint) (deviceInfo Device, err error) {
	var device C.dcgmDeviceAttributes_t
	device.version = makeVersion3(unsafe.Sizeof(device))

	result := C.dcgmGetDeviceAttributes(handle.handle, C.uint(gpuid), &device)
	if err = errorString(result); err != nil {
		return deviceInfo, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	// check if the given GPU is DCGM supported
	gpus, err := getSupportedDevices()
	if err != nil {
		return
	}

	supported := "No"

	for _, gpu := range gpus {
		if gpuid == gpu {
			supported = "Yes"
			break
		}
	}

	busid := *stringPtr(&device.identifiers.pciBusId[0])

	cpuAffinity, err := getCPUAffinity(gpuid)
	if err != nil {
		return
	}

	var topology []P2PLink
	var bandwidth int64
	// get device topology and bandwidth only if its a DCGM supported device
	if supported == "Yes" {
		topology, err = getDeviceTopology(gpuid)
		if err != nil {
			return
		}
		bandwidth, err = getPciBandwidth(gpuid)
		if err != nil {
			return
		}
	}

	uuid := *stringPtr(&device.identifiers.uuid[0])
	power := *uintPtr(device.powerLimits.defaultPowerLimit)

	pci := PCIInfo{
		BusID:     busid,
		BAR1:      *uintPtr(device.memoryUsage.bar1Total),
		FBTotal:   *uintPtr(device.memoryUsage.fbTotal),
		Bandwidth: bandwidth,
	}

	identifiers := DeviceIdentifiers{
		Brand:               *stringPtr(&device.identifiers.brandName[0]),
		Model:               *stringPtr(&device.identifiers.deviceName[0]),
		Serial:              *stringPtr(&device.identifiers.serial[0]),
		Vbios:               *stringPtr(&device.identifiers.vbios[0]),
		InforomImageVersion: *stringPtr(&device.identifiers.inforomImageVersion[0]),
		DriverVersion:       *stringPtr(&device.identifiers.driverVersion[0]),
	}

	deviceInfo = Device{
		GPU:           gpuid,
		DCGMSupported: supported,
		UUID:          uuid,
		Power:         power,
		PCI:           pci,
		Identifiers:   identifiers,
		Topology:      topology,
		CPUAffinity:   cpuAffinity,
	}
	return
}
