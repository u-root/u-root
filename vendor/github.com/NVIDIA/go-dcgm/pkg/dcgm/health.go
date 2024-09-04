/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
)

type SystemWatch struct {
	Type   string
	Status string
	Error  string
}

type DeviceHealth struct {
	GPU     uint
	Status  string
	Watches []SystemWatch
}

// HealthSet enable the DCGM health check system for the given systems
func HealthSet(groupId GroupHandle, systems HealthSystem) (err error) {
	result := C.dcgmHealthSet(handle.handle, groupId.handle, C.dcgmHealthSystems_t(systems))
	if err = errorString(result); err != nil {
		return fmt.Errorf("error setting health watches: %w", err)
	}
	return
}

// HealthGet retrieve the current state of the DCGM health check system
func HealthGet(groupId GroupHandle) (HealthSystem, error) {
	var systems C.dcgmHealthSystems_t

	result := C.dcgmHealthGet(handle.handle, groupId.handle, (*C.dcgmHealthSystems_t)(unsafe.Pointer(&systems)))
	if err := errorString(result); err != nil {
		return HealthSystem(0), err
	}
	return HealthSystem(systems), nil
}

type DiagErrorDetail struct {
	Message string
	Code    HealthCheckErrorCode
}

type Incident struct {
	System     HealthSystem
	Health     HealthResult
	Error      DiagErrorDetail
	EntityInfo GroupEntityPair
}

type HealthResponse struct {
	OverallHealth HealthResult
	Incidents     []Incident
}

// HealthCheck check the configured watches for any errors/failures/warnings that have occurred
// since the last time this check was invoked.  On the first call, stateful information
// about all of the enabled watches within a group is created but no error results are
// provided. On subsequent calls, any error information will be returned.
func HealthCheck(groupId GroupHandle) (HealthResponse, error) {
	var healthResults C.dcgmHealthResponse_v4
	healthResults.version = makeVersion4(unsafe.Sizeof(healthResults))

	result := C.dcgmHealthCheck(handle.handle, groupId.handle, (*C.dcgmHealthResponse_t)(unsafe.Pointer(&healthResults)))

	if err := errorString(result); err != nil {
		return HealthResponse{}, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	response := HealthResponse{
		OverallHealth: HealthResult(healthResults.overallHealth),
	}

	// number of watches that encountred error/warning
	incidents := uint(healthResults.incidentCount)

	response.Incidents = make([]Incident, incidents)

	for i := uint(0); i < incidents; i++ {
		response.Incidents[i] = Incident{
			System: HealthSystem(healthResults.incidents[i].system),
			Health: HealthResult(healthResults.incidents[i].health),
			Error: DiagErrorDetail{
				Message: *stringPtr(&healthResults.incidents[i].error.msg[0]),
				Code:    HealthCheckErrorCode(healthResults.incidents[i].error.code),
			},
			EntityInfo: GroupEntityPair{
				EntityGroupId: Field_Entity_Group(healthResults.incidents[i].entityInfo.entityGroupId),
				EntityId:      uint(healthResults.incidents[i].entityInfo.entityId),
			},
		}
	}

	return response, nil
}

func healthCheckByGpuId(gpuId uint) (deviceHealth DeviceHealth, err error) {
	name := fmt.Sprintf("health%d", rand.Uint64())
	groupId, err := CreateGroup(name)
	if err != nil {
		return
	}

	err = AddToGroup(groupId, gpuId)
	if err != nil {
		return
	}

	err = HealthSet(groupId, DCGM_HEALTH_WATCH_ALL)
	if err != nil {
		return
	}

	result, err := HealthCheck(groupId)
	if err != nil {
		return
	}

	status := healthStatus(result.OverallHealth)
	watches := []SystemWatch{}

	// number of watches that encountred error/warning
	incidents := len(result.Incidents)

	for j := 0; j < incidents; j++ {
		watch := SystemWatch{
			Type:   systemWatch(result.Incidents[j].System),
			Status: healthStatus(result.Incidents[j].Health),

			Error: result.Incidents[j].Error.Message,
		}
		watches = append(watches, watch)
	}

	deviceHealth = DeviceHealth{
		GPU:     gpuId,
		Status:  status,
		Watches: watches,
	}
	_ = DestroyGroup(groupId)
	return
}

func healthStatus(status HealthResult) string {
	switch status {
	case 0:
		return "Healthy"
	case 10:
		return "Warning"
	case 20:
		return "Failure"
	}
	return "N/A"
}

func systemWatch(watch HealthSystem) string {
	switch watch {
	case 1:
		return "PCIe watches"
	case 2:
		return "NVLINK watches"
	case 4:
		return "Power Managemnt unit watches"
	case 8:
		return "Microcontroller unit watches"
	case 16:
		return "Memory watches"
	case 32:
		return "Streaming Multiprocessor watches"
	case 64:
		return "Inforom watches"
	case 128:
		return "Temperature watches"
	case 256:
		return "Power watches"
	case 512:
		return "Driver-related watches"
	}
	return "N/A"
}
