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
 * File:   dcgm_test_apis.h
 */

#ifndef DCGM_AGENT_INTERNAL_H
#define DCGM_AGENT_INTERNAL_H

#include "dcgm_api_export.h"
#include "dcgm_structs.h"
#include "dcgm_structs_internal.h"
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/*****************************************************************************
 *****************************************************************************/
/*****************************************************************************
 * DCGM Test Methods, only used for testing, not officially supported
 *****************************************************************************/
/*****************************************************************************
 *****************************************************************************/

#define DCGM_EMBEDDED_HANDLE 0x7fffffff

/**
 * This method starts the Host Engine Server
 *
 * @param portNumber      IN: TCP port to listen on. This is only used if isTcp == 1.
 * @param socketPath      IN: This is the path passed to bind() when creating the socket
 *                            For isConnectionTCP == 1, this is the bind address. "" or NULL = All interfaces
 *                            For isConnectionTCP == 0, this is the path to the domain socket to use
 * @param isConnectionTCP IN: Whether to listen on a TCP/IP socket (1) or a unix domain socket (0)
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmEngineRun(unsigned short portNumber,
                                           char const *socketPath,
                                           unsigned int isConnectionTCP);

/**
 * This method is used to get values corresponding to the fields.
 * @return
 *      - \ref DCGM_ST_SUCCESS  On Success. Even when the API returns success, check for
 *                              individual status inside each field.
 *                              Look at values[index].status. The field values will be
 *                              populated only when status in each field is DCGM_ST_SUCCESS
 *      - DCGM_ST_?             In case of error
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetLatestValuesForFields(dcgmHandle_t pDcgmHandle,
                                                          int gpuId,
                                                          unsigned short fieldIds[],
                                                          unsigned int count,
                                                          dcgmFieldValue_v1 values[]);

/**
 * This method is used to get multiple values for a single field
 *
 * @return
 *      - \ref DCGM_ST_SUCCESS      on success.
 *      - DCGM_ST_?                 error code on failure
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetMultipleValuesForField(dcgmHandle_t pDcgmHandle,
                                                           int gpuId,
                                                           unsigned short fieldId,
                                                           int *count,
                                                           long long startTs,
                                                           long long endTs,
                                                           dcgmOrder_t order,
                                                           dcgmFieldValue_v1 values[]);

/**
 * Request updates for all field values that have updated since a given timestamp
 *
 * @param groupId             IN: Group ID representing a collection of one or more GPUs
 *                                Refer to \ref dcgmEngineGroupCreate for details on creating a group
 * @param sinceTimestamp      IN: Timestamp to request values since in usec since 1970. This will
 *                                be returned in nextSinceTimestamp for subsequent calls
 *                                0 = request all data
 * @param fieldIds            IN: Fields to return data for
 * @param numFieldIds         IN: Number of entries in fieldIds
 * @param nextSinceTimestamp OUT: Timestamp to use for sinceTimestamp on next call to this function
 * @param enumCB              IN: Callback to invoke for every field value update. Note that
 *                                multiple updates can be returned in each invocation
 * @param userData            IN: User data pointer to pass to the userData field of enumCB.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetFieldValuesSince(dcgmHandle_t pDcgmHandle,
                                                     dcgmGpuGrp_t groupId,
                                                     long long sinceTimestamp,
                                                     unsigned short *fieldIds,
                                                     int numFieldIds,
                                                     long long *nextSinceTimestamp,
                                                     dcgmFieldValueEnumeration_f enumCB,
                                                     void *userData);

/**
 * This method is used to tell the cache manager to watch a field value
 *
 * @param gpuId                               GPU ID to watch field on
 * @param fieldId                             Field ID to watch
 * @param updateFreq                          How often to update this field in usec
 * @param maxKeepAge                          How long to keep data for this field in seconds
 * @param maxKeepSamples                      Maximum number of samples to keep. 0=no limit
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a gpuId, \a fieldId, \a updateFreq, \a maxKeepAge,
 *                                            or \a maxKeepSamples are invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmWatchFieldValue(dcgmHandle_t pDcgmHandle,
                                                 int gpuId,
                                                 unsigned short fieldId,
                                                 long long updateFreq,
                                                 double maxKeepAge,
                                                 int maxKeepSamples);

/**
 * This method is used to tell the cache manager to unwatch a field value
 *
 * @param gpuId                               GPU ID to watch field on
 * @param fieldId                             Field ID to watch
 * @param clearCache                          Whether or not to clear all cached data for
 *                                            the field after the watch is removed
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a gpuId, \a fieldId, or \a clearCache is invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmUnwatchFieldValue(dcgmHandle_t pDcgmHandle,
                                                   int gpuId,
                                                   unsigned short fieldId,
                                                   int clearCache);

/*************************************************************************/
/**
 * Used to set vGPU configuration for the group of one or more GPUs identified by \a groupId.
 *
 * The configuration settings specified in \a pDeviceConfig are applied to all the GPUs in the
 * group. Since DCGM groups are a logical grouping of GPUs, the configuration settings Set for a group
 * stay intact for the individual GPUs even after the group is destroyed.
 *
 * If the user wishes to ignore the configuration of one or more properties in the input
 * \a pDeviceConfig then the property should be specified as one of \a DCGM_INT32_BLANK,
 * \a DCGM_INT64_BLANK, \a DCGM_FP64_BLANK or \a DCGM_STR_BLANK based on the data type of the
 * property to be ignored.
 *
 * If any of the properties fail to be configured for any of the GPUs in the group then the API
 * returns an error. The status handle \a statusHandle should be further evaluated to access error
 * attributes for the failed operations. Please refer to status management APIs at \ref DCGMAPI_ST
 * to access the error attributes.
 *
 * @param pDcgmHandle           IN: DCGM Handle
 *
 * @param groupId               IN: Group ID representing collection of one or more GPUs. Look
 *                                  at \ref dcgmGroupCreate for details on creating the group.
 *                                  applied for all the GPU in the group represented by
 *                                  \a groupId. The caller must populate the version field of
 *                                  \a pDeviceConfig.
 * @param statusHandle       IN/OUT: Resulting error status for multiple operations. Pass it as
 *                                   NULL if the detailed error information is not needed.
 *                                   Look at \ref dcgmStatusCreate for details on creating
 *                                   status handle.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the configuration has been successfully set.
 *        - \ref DCGM_ST_BADPARAM             if any of \a groupId or \a pDeviceConfig is invalid.
 *        - \ref DCGM_ST_VER_MISMATCH         if \a pDeviceConfig has the incorrect version.
 *        - \ref DCGM_ST_GENERIC_ERROR        if an unknown error has occurred.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmVgpuConfigSet(dcgmHandle_t pDcgmHandle,
                                               dcgmGpuGrp_t groupId,
                                               dcgmVgpuConfig_t *pDeviceConfig,
                                               dcgmStatus_t statusHandle);

/*************************************************************************/
/**
 * Used to get vGPU configuration for all the GPUs present in the group.
 *
 * This API can get the most recent target or desired configuration set by \ref dcgmConfigSet.
 * Set type as \a DCGM_CONFIG_TARGET_STATE to get target configuration. The target configuration
 * properties are maintained by DCGM and are automatically enforced after a GPU reset or
 * reinitialization is completed.
 *
 * The method can also be used to get the actual configuration state for the GPUs in the group.
 * Set type as \a DCGM_CONFIG_CURRENT_STATE to get the actually configuration state. Ideally, the
 * actual configuration state will be exact same as the target configuration state.
 *
 * If any of the property in the target configuration is unknown then the property value in the
 * output is populated as  one of DCGM_INT32_BLANK, DCGM_INT64_BLANK, DCGM_FP64_BLANK or
 * DCGM_STR_BLANK based on the data type of the property.
 *
 * If any of the property in the current configuration state is not supported then the property
 * value in the output is populated as one of DCGM_INT32_NOT_SUPPORTED, DCGM_INT64_NOT_SUPPORTED,
 * DCGM_FP64_NOT_SUPPORTED or DCGM_STR_NOT_SUPPORTED based on the data type of the property.
 *
 * If any of the properties can't be fetched for any of the GPUs in the group then the API returns
 * an error. The status handle \a statusHandle should be further evaluated to access error
 * attributes for the failed operations. Please refer to status management APIs at \ref DCGMAPI_ST
 * to access the error attributes.
 *
 * @param pDcgmHandle           IN: DCGM Handle
 * @param groupId               IN: Group ID representing collection of one or more GPUs. Look
 *                                  at \ref dcgmGroupCreate for details on creating the group.
 * @param type                  IN: Type of configuration values to be fetched.
 * @param count                 IN: The number of entries that \a deviceConfigList array can store.
 * @param deviceConfigList      OUT: Pointer to memory to hold requested configuration
 *                                   corresponding to all the GPUs in the group (\a groupId). The
 *                                   size of the memory must be greater than or equal to hold
 *                                   output information for the number of GPUs present in the
 *                                   group (\a groupId).
 * @param statusHandle       IN/OUT: Resulting error status for multiple operations. Pass it as
 *                                   NULL if the detailed error information is not needed.
 *                                   Look at \ref dcgmStatusCreate for details on creating
 *                                   status handle.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the configuration has been successfully fetched.
 *        - \ref DCGM_ST_BADPARAM             if any of \a groupId, \a type, \a count,
 *                                            or \a deviceConfigList is invalid.
 *        - \ref DCGM_ST_NOT_CONFIGURED       if the target configuration is not already set.
 *        - \ref DCGM_ST_VER_MISMATCH         if \a deviceConfigList has the incorrect version.
 *        - \ref DCGM_ST_GENERIC_ERROR        if an unknown error has occurred.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmVgpuConfigGet(dcgmHandle_t pDcgmHandle,
                                               dcgmGpuGrp_t groupId,
                                               dcgmConfigType_t type,
                                               int count,
                                               dcgmVgpuConfig_t deviceConfigList[],
                                               dcgmStatus_t statusHandle);

/*************************************************************************/
/**
 * Used to enforce previously set vGPU configuration for all the GPUs present in the group.
 *
 * This API provides a mechanism to the users to manually enforce the configuration at any point of
 * time. The configuration can only be enforced if it's already configured using the API \ref
 * dcgmConfigSet.
 *
 * If any of the properties can't be enforced for any of the GPUs in the group then the API returns
 * an error. The status handle \a statusHandle should be further evaluated to access error
 * attributes for the failed operations. Please refer to status management APIs at \ref DCGMAPI_ST
 * to access the error attributes.
 *
 * @param pDcgmHandle           IN: DCGM Handle
 *
 * @param groupId               IN: Group ID representing collection of one or more GPUs. Look at
 *                                  \ref dcgmGroupCreate for details on creating the group.
 *                                  Alternatively, pass in the group id as \a DCGM_GROUP_ALL_GPUS
 *                                  to perform operation on all the GPUs.
 * @param statusHandle       IN/OUT: Resulting error status for multiple operations. Pass it as
 *                                   NULL if the detailed error information is not needed.
 *                                   Look at \ref dcgmStatusCreate for details on creating
 *                                   status handle.
 * @return
 *        - \ref DCGM_ST_OK                   if the configuration has been successfully enforced.
 *        - \ref DCGM_ST_BADPARAM             if \a groupId is invalid.
 *        - \ref DCGM_ST_NOT_CONFIGURED       if the target configuration is not already set.
 *        - \ref DCGM_ST_GENERIC_ERROR        if an unknown error has occurred.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmVgpuConfigEnforce(dcgmHandle_t pDcgmHandle,
                                                   dcgmGpuGrp_t groupId,
                                                   dcgmStatus_t statusHandle);

/*************************************************************************/
/**
 * Gets vGPU device attributes corresponding to the \a gpuId. If operation is not successful for any of
 * the requested fields then the field is populated with one of DCGM_BLANK_VALUES defined in
 * dcgm_structs.h.
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param gpuId         IN: GPU Id corresponding to which the attributes
 *                          should be fetched
 * @param pDcgmVgpuAttr IN/OUT: vGPU Device attributes corresponding to \a gpuId.<br>
 *                              .version should be set to \ref dcgmVgpuDeviceAttributes_version before this call.
 *
 * @return
 *        - \ref DCGM_ST_OK            if the call was successful.
 *        - \ref DCGM_ST_VER_MISMATCH  if version is not set or is invalid.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetVgpuDeviceAttributes(dcgmHandle_t pDcgmHandle,
                                                         unsigned int gpuId,
                                                         dcgmVgpuDeviceAttributes_t *pDcgmVgpuAttr);

/*************************************************************************/
/**
 * Gets vGPU attributes corresponding to the \a vgpuId. If operation is not successful for any of
 * the requested fields then the field is populated with one of DCGM_BLANK_VALUES defined in
 * dcgm_structs.h.
 *
 * @param pDcgmHandle       IN: DCGM Handle
 * @param vgpuId            IN: vGPU Id corresponding to which the attributes should be fetched
 * @param pDcgmVgpuInstanceAttr IN/OUT: vGPU attributes corresponding to \a vgpuId.<br> .version should be set to
 *                                      \ref dcgmVgpuInstanceAttributes_version before this call.
 *
 * @return
 *        - \ref DCGM_ST_OK            if the call was successful.
 *        - \ref DCGM_ST_VER_MISMATCH  if .version is not set or is invalid.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetVgpuInstanceAttributes(dcgmHandle_t pDcgmHandle,
                                                           unsigned int vgpuId,
                                                           dcgmVgpuInstanceAttributes_t *pDcgmVgpuInstanceAttr);

/**
 * Stop a diagnostic if there is one currently running.
 *
 * @param pDcgmHandle                   IN: DCGM Handle
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a provided parameter is invalid or missing
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStopDiagnostic(dcgmHandle_t pDcgmHandle);

/**
 * This method injects a sample into the cache manager
 *
 * @param gpuId
 * @param dcgmInjectFieldValue
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmInjectFieldValue(dcgmHandle_t pDcgmHandle,
                                                  unsigned int gpuId,
                                                  dcgmInjectFieldValue_t *dcgmInjectFieldValue);

/**
 * This method retries the state of a field within the cache manager
 *
 * @param fieldInfo Structure to populate. fieldInfo->gpuId and fieldInfo->fieldId must
 *                  be populated on calling for this call to work
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetCacheManagerFieldInfo(dcgmHandle_t pDcgmHandle,
                                                          dcgmCacheManagerFieldInfo_v4_t *fieldInfo);

/**
 * This method returns the status of the gpu
 *
 * @param gpuId
 * @param DcgmEntityStatus_t
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetGpuStatus(dcgmHandle_t pDcgmHandle, unsigned int gpuId, DcgmEntityStatus_t *status);

/**
 * Create fake entities for injection testing
 *
 * @param createFakeEntities Details about the number and type of entities to create
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmCreateFakeEntities(dcgmHandle_t pDcgmHandle,
                                                    dcgmCreateFakeEntities_t *createFakeEntities);

/**
 * This method injects a sample into the cache manager
 *
 * @param entityGroupId
 * @param entityId
 * @param dcgmInjectFieldValue
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmEntityInjectFieldValue(dcgmHandle_t pDcgmHandle,
                                                        dcgm_field_entity_group_t entityGroupId,
                                                        dcgm_field_eid_t entityId,
                                                        dcgmInjectFieldValue_t *dcgmInjectFieldValue);

/**
 * This method sets the link state of an entity's NvLink
 *
 * dcgmHandle_t dcgmHandle
 * linkState    contains details about the link state to set
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmSetEntityNvLinkLinkState(dcgmHandle_t dcgmHandle,
                                                          dcgmSetNvLinkLinkState_v1 *linkState);

/**
 * Creates a MIG entity with according to the specifications in the struct
 *
 * @param dcgmHandle       IN: DCGM Handle
 * @param cme              IN: struct stating which kind of entity is being created, who the parent entity is, flags
 *                             for processing it, and the profile to specify what size of that entity to create.
 * @return
 *        - \ref DCGM_ST_OK                if the call was successful.
 *        - \ref DCGM_ST_BADPARAM          if any parameter is invalid
 *        - \ref DCGM_ST_REQUIRES_ROOT     if the hostengine is not running as root
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmCreateMigEntity(dcgmHandle_t dcgmHandle, dcgmCreateMigEntity_t *cme);

/**
 * Delete the specified MIG entity
 *
 * @param dcgmHandle     IN: DCGM Handle
 * @param dme            IN: struct specifying which entity to delete with flags to suggest how to process it.
 * @return
 *        - \ref DCGM_ST_OK                if the call was successful.
 *        - \ref DCGM_ST_BADPARAM          if any parameter is invalid
 *        - \ref DCGM_ST_REQUIRES_ROOT     if the hostengine is not running as root
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmDeleteMigEntity(dcgmHandle_t dcgmHandle, dcgmDeleteMigEntity_t *dme);

/**
 * @brief Pauses all DCGM modules from updating field values
 *
 * This method sends a pause message to each loaded module.
 * It's up to the module to decide whether to handle or ignore the message.
 *
 * @param[in] pDcgmHandle DCGM Handle of an active connection
 *
 * @return
 *      - \ref DCGM_ST_OK if successful
 *      - \ref DCGM_ST_* on error
 *
 * @note If this function fails, the modules may be in an inconsistent state.
 * @note You may call \ref dcgmModuleGetStatuses to see which modules are paused.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmPauseTelemetryForDiag(dcgmHandle_t pDcgmHandle);

/**
 * @brief Resumes all DCGM modules to updating field values
 *
 * This method sends a resume message to each loaded module.
 * It's up to the module to decide whether to handle or ignore the message.
 *
 * @param[in] pDcgmHandle DCGM Handle of an active connection
 *
 * @return
 *      - \ref DCGM_ST_OK if successful
 *      - \ref DCGM_ST_* on error
 *
 * @note If this function fails, the modules may be in an inconsistent state.
 * @note You may call \ref dcgmModuleGetStatuses to see which modules are resumed. The satus of the resumed modules
 *       should be \ref dcgmModuleStatus_t::DcgmModuleStatusLoaded.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmResumeTelemetryForDiag(dcgmHandle_t pDcgmHandle);

#ifdef __cplusplus
}
#endif

#endif /* DCGM_AGENT_INTERNAL_H */
