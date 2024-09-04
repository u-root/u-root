package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"unsafe"
)

type MetricGroup struct {
	Major    uint
	Minor    uint
	FieldIds []uint
}

func getSupportedMetricGroups(gpuId uint) (groups []MetricGroup, err error) {

	var groupInfo C.dcgmProfGetMetricGroups_t
	groupInfo.version = makeVersion3(unsafe.Sizeof(groupInfo))

	groupInfo.gpuId = C.uint(gpuId)

	result := C.dcgmProfGetSupportedMetricGroups(handle.handle, &groupInfo)

	if err = errorString(result); err != nil {
		return groups, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	var count = uint(groupInfo.numMetricGroups)

	for i := uint(0); i < count; i++ {
		var group MetricGroup
		group.Major = uint(groupInfo.metricGroups[i].majorId)
		group.Minor = uint(groupInfo.metricGroups[i].minorId)

		var fieldCount = uint(groupInfo.metricGroups[i].numFieldIds)

		for j := uint(0); j < fieldCount; j++ {
			group.FieldIds = append(group.FieldIds, uint(groupInfo.metricGroups[i].fieldIds[j]))
		}
		groups = append(groups, group)
	}

	return groups, nil
}
