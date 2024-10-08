package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"unsafe"
)

const DIAG_RESULT_STRING_SIZE = 1024

type DiagType int

const (
	DiagQuick    DiagType = 1
	DiagMedium            = 2
	DiagLong              = 3
	DiagExtended          = 4
)

type DiagResult struct {
	Status       string
	TestName     string
	TestOutput   string
	ErrorCode    uint
	ErrorMessage string
}

type GpuResult struct {
	GPU         uint
	RC          uint
	DiagResults []DiagResult
}

type DiagResults struct {
	Software []DiagResult
	PerGpu   []GpuResult
}

func diagResultString(r int) string {
	switch r {
	case C.DCGM_DIAG_RESULT_PASS:
		return "pass"
	case C.DCGM_DIAG_RESULT_SKIP:
		return "skipped"
	case C.DCGM_DIAG_RESULT_WARN:
		return "warn"
	case C.DCGM_DIAG_RESULT_FAIL:
		return "fail"
	case C.DCGM_DIAG_RESULT_NOT_RUN:
		return "notrun"
	}
	return ""
}

func swTestName(t int) string {
	switch t {
	case C.DCGM_SWTEST_DENYLIST:
		return "presence of drivers on the denylist (e.g. nouveau)"
	case C.DCGM_SWTEST_NVML_LIBRARY:
		return "presence (and version) of NVML lib"
	case C.DCGM_SWTEST_CUDA_MAIN_LIBRARY:
		return "presence (and version) of CUDA lib"
	case C.DCGM_SWTEST_CUDA_RUNTIME_LIBRARY:
		return "presence (and version) of CUDA RT lib"
	case C.DCGM_SWTEST_PERMISSIONS:
		return "character device permissions"
	case C.DCGM_SWTEST_PERSISTENCE_MODE:
		return "persistence mode enabled"
	case C.DCGM_SWTEST_ENVIRONMENT:
		return "CUDA environment vars that may slow tests"
	case C.DCGM_SWTEST_PAGE_RETIREMENT:
		return "pending frame buffer page retirement"
	case C.DCGM_SWTEST_GRAPHICS_PROCESSES:
		return "graphics processes running"
	case C.DCGM_SWTEST_INFOROM:
		return "inforom corruption"
	}

	return ""
}

func gpuTestName(t int) string {

	switch t {
	case C.DCGM_MEMORY_INDEX:
		return "Memory"
	case C.DCGM_DIAGNOSTIC_INDEX:
		return "Diagnostic"
	case C.DCGM_PCI_INDEX:
		return "PCIe"
	case C.DCGM_SM_STRESS_INDEX:
		return "SM Stress"
	case C.DCGM_TARGETED_STRESS_INDEX:
		return "Targeted Stress"
	case C.DCGM_TARGETED_POWER_INDEX:
		return "Targeted Power"
	case C.DCGM_MEMORY_BANDWIDTH_INDEX:
		return "Memory bandwidth"
	case C.DCGM_MEMTEST_INDEX:
		return "Memtest"
	case C.DCGM_PULSE_TEST_INDEX:
		return "Pulse"
	case C.DCGM_EUD_TEST_INDEX:
		return "EUD"
	case C.DCGM_SOFTWARE_INDEX:
		return "Software"
	case C.DCGM_CONTEXT_CREATE_INDEX:
		return "Context create"
	}
	return ""
}

func newDiagResult(testResult C.dcgmDiagTestResult_v3, testName string) DiagResult {
	msg := C.GoString((*C.char)(unsafe.Pointer(&testResult.error[0].msg)))
	info := C.GoString((*C.char)(unsafe.Pointer(&testResult.info)))

	return DiagResult{
		Status:       diagResultString(int(testResult.status)),
		TestName:     testName,
		TestOutput:   info,
		ErrorCode:    uint(testResult.error[0].code),
		ErrorMessage: msg,
	}
}

func diagLevel(diagType DiagType) C.dcgmDiagnosticLevel_t {
	switch diagType {
	case DiagQuick:
		return C.DCGM_DIAG_LVL_SHORT
	case DiagMedium:
		return C.DCGM_DIAG_LVL_MED
	case DiagLong:
		return C.DCGM_DIAG_LVL_LONG
	case DiagExtended:
		return C.DCGM_DIAG_LVL_XLONG
	}
	return C.DCGM_DIAG_LVL_INVALID
}

func RunDiag(diagType DiagType, groupId GroupHandle) (DiagResults, error) {
	var diagResults C.dcgmDiagResponse_v9
	diagResults.version = makeVersion9(unsafe.Sizeof(diagResults))

	result := C.dcgmRunDiagnostic(handle.handle, groupId.handle, diagLevel(diagType), (*C.dcgmDiagResponse_v9)(unsafe.Pointer(&diagResults)))
	if err := errorString(result); err != nil {
		return DiagResults{}, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	var diagRun DiagResults
	for i := 0; i < int(diagResults.levelOneTestCount); i++ {
		dr := newDiagResult(diagResults.levelOneResults[i], swTestName(i))
		diagRun.Software = append(diagRun.Software, dr)
	}

	for i := uint(0); i < uint(diagResults.gpuCount); i++ {
		r := diagResults.perGpuResponses[i]
		gr := GpuResult{GPU: uint(r.gpuId), RC: uint(r.hwDiagnosticReturn)}
		for j := 0; j < int(C.DCGM_PER_GPU_TEST_COUNT_V8); j++ {
			dr := newDiagResult(r.results[j], gpuTestName(j))
			gr.DiagResults = append(gr.DiagResults, dr)
		}
		diagRun.PerGpu = append(diagRun.PerGpu, gr)
	}

	return diagRun, nil
}
