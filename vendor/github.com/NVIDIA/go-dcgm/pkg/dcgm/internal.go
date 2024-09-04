/*
 * Copyright (c) 2023, NVIDIA CORPORATION.  All rights reserved.
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
#cgo linux LDFLAGS: -ldl -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
#cgo darwin LDFLAGS: -ldl -Wl,--export-dynamic -Wl,-undefined,dynamic_lookup

#include "dcgm_test_apis.h"
#include "dcgm_test_structs.h"
#include "dcgm_structs_internal.h"
*/
import "C"

import (
	"unsafe"
)

type MigHierarchyInfo struct {
	Entity       GroupEntityPair
	Parent       GroupEntityPair
	SliceProfile MigProfile
}

func CreateFakeEntities(entities []MigHierarchyInfo) ([]uint, error) {
	ccfe := C.dcgmCreateFakeEntities_v2{
		version:     C.dcgmCreateFakeEntities_version2,
		numToCreate: C.uint(len(entities)),
		entityList:  [C.DCGM_MAX_HIERARCHY_INFO]C.dcgmMigHierarchyInfo_t{},
	}
	for i := range entities {
		if i >= C.DCGM_MAX_HIERARCHY_INFO {
			break
		}
		entity := entities[i]
		ccfe.entityList[i] = C.dcgmMigHierarchyInfo_t{
			entity: C.dcgmGroupEntityPair_t{
				entityGroupId: C.dcgm_field_entity_group_t(entity.Entity.EntityGroupId),
				entityId:      C.uint(entity.Entity.EntityId),
			},
			parent: C.dcgmGroupEntityPair_t{
				entityGroupId: C.dcgm_field_entity_group_t(entity.Parent.EntityGroupId),
				entityId:      C.uint(entity.Parent.EntityId),
			},
			sliceProfile: C.dcgmMigProfile_t(entity.SliceProfile),
		}
	}
	result := C.dcgmCreateFakeEntities(handle.handle, &ccfe)

	if err := errorString(result); err != nil {
		return nil, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}
	gpuIDs := make([]uint, ccfe.numToCreate)
	for i := 0; i < int(ccfe.numToCreate); i++ {
		gpuIDs[i] = uint(ccfe.entityList[i].entity.entityId)
	}

	return gpuIDs, nil
}

func InjectFieldValue(gpu uint, fieldID uint, fieldType uint, status int, ts int64, value interface{}) error {
	field := C.dcgmInjectFieldValue_t{
		version:   C.dcgmInjectFieldValue_version1,
		fieldId:   C.ushort(fieldID),
		fieldType: C.ushort(fieldType),
		status:    C.int(status),
		ts:        C.long(ts),
	}

	switch fieldType {
	case DCGM_FT_INT64:
		i64Val := value.(int64)
		ptr := (*C.int64_t)(unsafe.Pointer(&field.value[0]))
		*ptr = C.int64_t(i64Val)
	case DCGM_FT_DOUBLE:
		dbVal := value.(float64)
		ptr := (*C.double)(unsafe.Pointer(&field.value[0]))
		*ptr = C.double(dbVal)
	}

	result := C.dcgmInjectFieldValue(handle.handle, C.uint(gpu), &field)

	if err := errorString(result); err != nil {
		return &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	return nil
}
