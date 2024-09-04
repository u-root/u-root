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

/*
 * File: dcgm_test_structs.h
 */

#ifndef DCGM_TEST_STRUCTS_H
#define DCGM_TEST_STRUCTS_H

#include "dcgm_fields.h"
#include "dcgm_structs.h"
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Structure to represent default and target vgpu configuration for a device
 */
typedef struct
{
    unsigned int version;     //!< Version number (dcgmConfig_version)
    unsigned int gpuId;       //!< GPU ID
    unsigned int eccMode;     //!< ECC Mode  (0: Disabled, 1 : Enabled, DCGM_INT32_BLANK : Ignored)
    unsigned int computeMode; //!< Compute Mode (One of DCGM_CONFIG_COMPUTEMODE_? OR DCGM_INT32_BLANK to Ignore)
    dcgmConfigPerfStateSettings_t perfState; //!< Performance State Settings (clocks / boost mode)
    dcgmConfigPowerLimit_t powerLimit;       //!< Power Limits
} dcgmVgpuConfig_v1;

/**
 * Typedef for \ref dcgmVgpuConfig_v1
 */
typedef dcgmVgpuConfig_v1 dcgmVgpuConfig_t;

/**
 * Version 1 for \ref dcgmVgpuConfig_v1
 */
#define dcgmVgpuConfig_version1 MAKE_DCGM_VERSION(dcgmVgpuConfig_v1, 1)

/**
 * Latest version for \ref dcgmVgpuConfig_t
 */
#define dcgmVgpuConfig_version dcgmVgpuConfig_version1

/**
 * Represents the vGPU attributes corresponding to a physical device
 */
typedef struct
{
    unsigned int version;                 //!< Version number (dcgmVgpuDeviceAttributes_version)
    unsigned int activeVgpuInstanceCount; //!< Count of active vGPU instances on the device
    unsigned int activeVgpuInstanceIds[DCGM_MAX_VGPU_INSTANCES_PER_PGPU]; //!< List of vGPU instances
    unsigned int creatableVgpuTypeCount;                                  //!< Creatable vGPU type count
    unsigned int creatableVgpuTypeIds[DCGM_MAX_VGPU_TYPES_PER_PGPU];      //!< List of Creatable vGPU types
    unsigned int supportedVgpuTypeCount;                                  //!< Supported vGPU type count
    dcgmDeviceVgpuTypeInfo_v1
        supportedVgpuTypeInfo[DCGM_MAX_VGPU_TYPES_PER_PGPU]; //!< Info related to vGPUs supported on the device
    dcgmDeviceVgpuUtilInfo_v1 vgpuUtilInfo[DCGM_MAX_VGPU_TYPES_PER_PGPU]; //!< Utilization specific to vGPU instance
    unsigned int gpuUtil;                                                 //!< GPU utilization
    unsigned int memCopyUtil;                                             //!< Memory utilization
    unsigned int encUtil;                                                 //!< Encoder utilization
    unsigned int decUtil;                                                 //!< Decoder utilization
} dcgmVgpuDeviceAttributes_v6;

/**
 * Version 6 for \ref dcgmVgpuDeviceAttributes_v6
 */
#define dcgmVgpuDeviceAttributes_version6 MAKE_DCGM_VERSION(dcgmVgpuDeviceAttributes_v6, 1)

typedef struct
{
    unsigned int version;                 //!< Version number (dcgmVgpuDeviceAttributes_version)
    unsigned int activeVgpuInstanceCount; //!< Count of active vGPU instances on the device
    unsigned int activeVgpuInstanceIds[DCGM_MAX_VGPU_INSTANCES_PER_PGPU]; //!< List of vGPU instances
    unsigned int creatableVgpuTypeCount;                                  //!< Creatable vGPU type count
    unsigned int creatableVgpuTypeIds[DCGM_MAX_VGPU_TYPES_PER_PGPU];      //!< List of Creatable vGPU types
    unsigned int supportedVgpuTypeCount;                                  //!< Supported vGPU type count
    dcgmDeviceVgpuTypeInfo_v2
        supportedVgpuTypeInfo[DCGM_MAX_VGPU_TYPES_PER_PGPU]; //!< Info related to vGPUs supported on the device
    dcgmDeviceVgpuUtilInfo_v1 vgpuUtilInfo[DCGM_MAX_VGPU_TYPES_PER_PGPU]; //!< Utilization specific to vGPU instance
    unsigned int gpuUtil;                                                 //!< GPU utilization
    unsigned int memCopyUtil;                                             //!< Memory utilization
    unsigned int encUtil;                                                 //!< Encoder utilization
    unsigned int decUtil;                                                 //!< Decoder utilization
} dcgmVgpuDeviceAttributes_v7;

/**
 *  * Typedef for \ref dcgmVgpuDeviceAttributes_v7
 *   */
typedef dcgmVgpuDeviceAttributes_v7 dcgmVgpuDeviceAttributes_t;

/**
 * Version 7 for \ref dcgmVgpuDeviceAttributes_v7
 */
#define dcgmVgpuDeviceAttributes_version7 MAKE_DCGM_VERSION(dcgmVgpuDeviceAttributes_v7, 7)

/**
 * Latest version for \ref dcgmVgpuDeviceAttributes_t
 */
#define dcgmVgpuDeviceAttributes_version dcgmVgpuDeviceAttributes_version7

/**
 * Represents attributes specific to vGPU instance
 */
typedef struct
{
    unsigned int version;                                 //!< Version number (dcgmVgpuInstanceAttributes_version)
    char vmId[DCGM_DEVICE_UUID_BUFFER_SIZE];              //!< VM ID of the vGPU instance
    char vmName[DCGM_DEVICE_UUID_BUFFER_SIZE];            //!< VM name of the vGPU instance
    unsigned int vgpuTypeId;                              //!< Type ID of the vGPU instance
    char vgpuUuid[DCGM_DEVICE_UUID_BUFFER_SIZE];          //!< UUID of the vGPU instance
    char vgpuDriverVersion[DCGM_DEVICE_UUID_BUFFER_SIZE]; //!< Driver version of the vGPU instance
    unsigned int fbUsage;                                 //!< Fb usage of the vGPU instance
    unsigned int licenseStatus;                           //!< License status of the vGPU instance
    unsigned int frameRateLimit;                          //!< Frame rate limit of the vGPU instance
} dcgmVgpuInstanceAttributes_v1;

/**
 * Typedef for \ref dcgmVgpuInstanceAttributes_v1
 */
typedef dcgmVgpuInstanceAttributes_v1 dcgmVgpuInstanceAttributes_t;

/**
 * Version 1 for \ref dcgmVgpuInstanceAttributes_v1
 */
#define dcgmVgpuInstanceAttributes_version1 MAKE_DCGM_VERSION(dcgmVgpuInstanceAttributes_v1, 1)

/**
 * Latest version for \ref dcgmVgpuInstanceAttributes_t
 */
#define dcgmVgpuInstanceAttributes_version dcgmVgpuInstanceAttributes_version1

/* Flags to ask nv-hostengine to process MIG events differently. */
/*
 * This flag is only meant to be used when running many commands that will trigger the
 * MIG configuration to get loaded again. The intended use is that if you are running
 * many commands that will cause the MIG configuration to change, then ask the hostengine
 * to only process the last one in order to prevent conflicts in how you are updating the
 * MIG information. For example, if you delete one compute instance on a GPU and
 * the hostengine processes the event from NVML before you delete the next one, the ID
 * of the compute instance will have changed in between. Using the flag asks the
 * hostengine to ignore those events temporarily while you are performing updates */
#define DCGM_MIG_RECONFIG_DELAY_PROCESSING 0x1

typedef struct
{
    unsigned int version;                    //!< Version number of this struct
    dcgm_field_entity_group_t entityGroupId; //!< entity group of the MIG entity being deleted
    dcgm_field_eid_t entityId;               //!< entity id of the MIG entity being deleted
    unsigned int flags;                      //!< flags to influence nv-hostengine's processing of the request
} dcgmDeleteMigEntity_v1;

#define dcgmDeleteMigEntity_version1 MAKE_DCGM_VERSION(dcgmDeleteMigEntity_v1, 1)

#define dcgmDeleteMigEntity_version dcgmDeleteMigEntity_version1

typedef dcgmDeleteMigEntity_v1 dcgmDeleteMigEntity_t;

/**
 * Enum for the kinds of MIG creations
 */
typedef enum
{
    DcgmMigCreateGpuInstance     = 0, /*!< Create a GPU instance */
    DcgmMigCreateComputeInstance = 1, /*!< Create a compute instance */
} dcgmMigCreate_t;

typedef struct
{
    unsigned int version;         //!< Version number of this request
    dcgm_field_eid_t parentId;    //!< The entity id of the parent (entity group is inferred from createOption
    dcgmMigProfile_t profile;     //!< Specify the MIG profile to create
    dcgmMigCreate_t createOption; //!< Specify if we're creating a GPU instance or compute instance
    unsigned int flags;           //!< flags to influence nv-hostengine's processing of the request
} dcgmCreateMigEntity_v1;

#define dcgmCreateMigEntity_version1 MAKE_DCGM_VERSION(dcgmCreateMigEntity_v1, 1)

#define dcgmCreateMigEntity_version dcgmCreateMigEntity_version1

typedef dcgmCreateMigEntity_v1 dcgmCreateMigEntity_t;

#ifdef __cplusplus
}
#endif

#endif /* DCGM_STRUCTS_H */
