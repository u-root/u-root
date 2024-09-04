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
#include "field_values_cb.h"
extern int go_dcgmFieldValueEntityEnumeration(dcgm_field_entity_group_t entityGroupId,
            dcgm_field_eid_t entityId,
            dcgmFieldValue_v1 *values,
            int numValues,
            void *userData);
*/
import "C"
import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

type callback struct {
	mu     sync.Mutex
	Values []FieldValue_v2
}

func (cb *callback) processValues(entityGroup Field_Entity_Group, entityID uint, cvalues []C.dcgmFieldValue_v1) {
	values := dcgmFieldValue_v1ToFieldValue_v2(entityGroup, entityID, cvalues)
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Values = append(cb.Values, values...)
}

//export go_dcgmFieldValueEntityEnumeration
func go_dcgmFieldValueEntityEnumeration(
	entityGroup C.dcgm_field_entity_group_t,
	entityID C.dcgm_field_eid_t,
	values *C.dcgmFieldValue_v1,
	numValues C.int,
	userData unsafe.Pointer) C.int {
	ptrValues := unsafe.Pointer(values)
	if ptrValues != nil {
		valuesSlice := (*[1 << 30]C.dcgmFieldValue_v1)(ptrValues)[0:numValues]
		if userData != nil {
			processor := (*callback)(userData)
			processor.processValues(Field_Entity_Group(entityGroup), uint(entityID), valuesSlice)
		}
	}
	return 0
}

// GetValuesSince reads and returns field values for a specified group of entities, such as GPUs,
// that have been updated since a given timestamp. It allows for targeted data retrieval based on time criteria.
//
// GPUGroup is a GroupHandle that identifies the group of entities to operate on. It can be obtained from CreateGroup
// for a specific group of GPUs or use GroupAllGPUs() to target all GPUs.
//
// fieldGroup is a FieldHandle representing the group of fields for which data is requested.
//
// sinceTime is a time.Time value representing the timestamp from which to request updated values.
// A zero value (time.Time{}) requests all available data.
//
// Returns []FieldValue_v2 slice containing the requested field values, a time.Time indicating the time
// of the latest data retrieval, and an error if there is any issue during the operation.
func GetValuesSince(GPUGroup GroupHandle, fieldGroup FieldHandle, sinceTime time.Time) ([]FieldValue_v2, time.Time, error) {
	var (
		nextSinceTimestamp C.longlong
	)

	cbResult := &callback{}

	result := C.dcgmGetValuesSince_v2(handle.handle,
		GPUGroup.handle,
		C.dcgmFieldGrp_t(fieldGroup.handle),
		C.longlong(sinceTime.UnixMicro()),
		&nextSinceTimestamp,
		(C.dcgmFieldValueEnumeration_f)(unsafe.Pointer(C.fieldValueEntityCallback)),
		unsafe.Pointer(cbResult))
	if result != C.DCGM_ST_OK {
		return nil, time.Time{}, fmt.Errorf("dcgmGetValuesSince_v2 failed with error code %d", int(result))
	}

	return cbResult.Values, timestampUSECToTime(int64(nextSinceTimestamp)), nil
}

func timestampUSECToTime(timestampUSEC int64) time.Time {
	// Convert microseconds to seconds and nanoseconds
	sec := timestampUSEC / 1000000           // Convert microseconds to seconds
	nsec := (timestampUSEC % 1000000) * 1000 // Convert the remaining microseconds to nanoseconds
	// Use time.Unix to get a time.Time object
	return time.Unix(sec, nsec)
}
