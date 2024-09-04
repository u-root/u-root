package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"unsafe"
)

type DcgmStatus struct {
	Memory int64
	CPU    float64
}

func introspect() (engine DcgmStatus, err error) {
	var memory C.dcgmIntrospectMemory_t
	memory.version = makeVersion1(unsafe.Sizeof(memory))
	waitIfNoData := 1
	result := C.dcgmIntrospectGetHostengineMemoryUsage(handle.handle, &memory, C.int(waitIfNoData))

	if err = errorString(result); err != nil {
		return engine, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	var cpu C.dcgmIntrospectCpuUtil_t

	cpu.version = makeVersion1(unsafe.Sizeof(cpu))
	result = C.dcgmIntrospectGetHostengineCpuUtilization(handle.handle, &cpu, C.int(waitIfNoData))

	if err = errorString(result); err != nil {
		return engine, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	engine = DcgmStatus{
		Memory: toInt64(memory.bytesUsed) / 1024,
		CPU:    *dblToFloat(cpu.total) * 100,
	}
	return
}
