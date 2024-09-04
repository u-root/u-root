package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

/*
 *See dcgm_structs.h
 *	DCGM_CPU_CORE_BITMASK_COUNT_V1 (DCGM_MAX_NUM_CPU_CORES / sizeof(uint64_t) / CHAR_BIT)
 *	or
 *	1024 / 8 / 8
 */

const (
	MAX_NUM_CPU_CORES          uint = C.DCGM_MAX_NUM_CPU_CORES
	MAX_NUM_CPUS               uint = C.DCGM_MAX_NUM_CPUS
	CHAR_BIT                   uint = C.CHAR_BIT
	MAX_CPU_CORE_BITMASK_COUNT uint = 1024 / 8 / 8
)

type CpuHierarchyCpu_v1 struct {
	CpuId      uint
	OwnedCores []uint64
}

type CpuHierarchy_v1 struct {
	Version uint
	NumCpus uint
	Cpus    [MAX_NUM_CPUS]CpuHierarchyCpu_v1
}

func GetCpuHierarchy() (hierarchy CpuHierarchy_v1, err error) {
	var c_hierarchy C.dcgmCpuHierarchy_v1
	c_hierarchy.version = C.dcgmCpuHierarchy_version1
	ptr_hierarchy := (*C.dcgmCpuHierarchy_v1)(unsafe.Pointer(&c_hierarchy))
	result := C.dcgmGetCpuHierarchy(handle.handle, ptr_hierarchy)

	if err = errorString(result); err != nil {
		return toCpuHierarchy(c_hierarchy), fmt.Errorf("Error retrieving DCGM MIG hierarchy: %s", err)
	}

	return toCpuHierarchy(c_hierarchy), nil
}

func toCpuHierarchy(c_hierarchy C.dcgmCpuHierarchy_v1) CpuHierarchy_v1 {
	var hierarchy CpuHierarchy_v1
	hierarchy.Version = uint(c_hierarchy.version)
	hierarchy.NumCpus = uint(c_hierarchy.numCpus)
	for i := uint(0); i < hierarchy.NumCpus; i++ {
		bits := make([]uint64, MAX_CPU_CORE_BITMASK_COUNT)

		for j := uint(0); j < MAX_CPU_CORE_BITMASK_COUNT; j++ {
			bits[j] = uint64(c_hierarchy.cpus[i].ownedCores.bitmask[j])
		}

		hierarchy.Cpus[i] = CpuHierarchyCpu_v1{
			CpuId:      uint(c_hierarchy.cpus[i].cpuId),
			OwnedCores: bits,
		}
	}

	return hierarchy
}
