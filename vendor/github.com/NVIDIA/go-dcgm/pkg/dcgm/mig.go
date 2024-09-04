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

type Field_Entity_Group uint

const (
	FE_NONE Field_Entity_Group = iota
	FE_GPU
	FE_VGPU
	FE_SWITCH
	FE_GPU_I
	FE_GPU_CI
	FE_LINK
	FE_CPU
	FE_CPU_CORE
	FE_COUNT
)

func (e Field_Entity_Group) String() string {
	switch e {
	case FE_GPU:
		return "GPU"
	case FE_VGPU:
		return "vGPU"
	case FE_SWITCH:
		return "NvSwitch"
	case FE_GPU_I:
		return "GPU Instance"
	case FE_GPU_CI:
		return "GPU Compute Instance"
	case FE_LINK:
		return "NvLink"
	case FE_CPU:
		return "CPU"
	case FE_CPU_CORE:
		return "CPU Core"
	}
	return "unknown"
}

type GroupEntityPair struct {
	EntityGroupId Field_Entity_Group
	EntityId      uint
}

type MigEntityInfo struct {
	GpuUuid               string
	NvmlGpuIndex          uint
	NvmlInstanceId        uint
	NvmlComputeInstanceId uint
	NvmlMigProfileId      uint
	NvmlProfileSlices     uint
}

type MigHierarchyInfo_v2 struct {
	Entity GroupEntityPair
	Parent GroupEntityPair
	Info   MigEntityInfo
}

const (
	MAX_NUM_DEVICES    uint = C.DCGM_MAX_NUM_DEVICES
	MAX_HIERARCHY_INFO uint = C.DCGM_MAX_HIERARCHY_INFO
)

type MigHierarchy_v2 struct {
	Version    uint
	Count      uint
	EntityList [C.DCGM_MAX_HIERARCHY_INFO]MigHierarchyInfo_v2
}

func GetGpuInstanceHierarchy() (hierarchy MigHierarchy_v2, err error) {
	var c_hierarchy C.dcgmMigHierarchy_v2
	c_hierarchy.version = C.dcgmMigHierarchy_version2
	ptr_hierarchy := (*C.dcgmMigHierarchy_v2)(unsafe.Pointer(&c_hierarchy))
	result := C.dcgmGetGpuInstanceHierarchy(handle.handle, ptr_hierarchy)

	if err = errorString(result); err != nil {
		return toMigHierarchy(c_hierarchy), fmt.Errorf("Error retrieving DCGM MIG hierarchy: %s", err)
	}

	return toMigHierarchy(c_hierarchy), nil
}

func toMigHierarchy(c_hierarchy C.dcgmMigHierarchy_v2) MigHierarchy_v2 {
	var hierarchy MigHierarchy_v2
	hierarchy.Version = uint(c_hierarchy.version)
	hierarchy.Count = uint(c_hierarchy.count)
	for i := uint(0); i < hierarchy.Count; i++ {
		hierarchy.EntityList[i] = MigHierarchyInfo_v2{
			Entity: GroupEntityPair{Field_Entity_Group(c_hierarchy.entityList[i].entity.entityGroupId), uint(c_hierarchy.entityList[i].entity.entityId)},
			Parent: GroupEntityPair{Field_Entity_Group(c_hierarchy.entityList[i].parent.entityGroupId), uint(c_hierarchy.entityList[i].parent.entityId)},
			Info: MigEntityInfo{
				GpuUuid:               *stringPtr(&c_hierarchy.entityList[i].info.gpuUuid[0]),
				NvmlGpuIndex:          uint(c_hierarchy.entityList[i].info.nvmlGpuIndex),
				NvmlInstanceId:        uint(c_hierarchy.entityList[i].info.nvmlInstanceId),
				NvmlComputeInstanceId: uint(c_hierarchy.entityList[i].info.nvmlComputeInstanceId),
				NvmlMigProfileId:      uint(c_hierarchy.entityList[i].info.nvmlMigProfileId),
				NvmlProfileSlices:     uint(c_hierarchy.entityList[i].info.nvmlProfileSlices),
			},
		}
	}

	return hierarchy
}
