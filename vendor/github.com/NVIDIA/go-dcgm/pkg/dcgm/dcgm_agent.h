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

#ifndef DCGM_AGENT_H
#define DCGM_AGENT_H

#define DCGM_PUBLIC_API
#include "dcgm_structs.h"

#ifdef __cplusplus
extern "C" {
#endif


/***************************************************************************************************/
/** @defgroup DCGMAPI_Admin Administrative
 *
 *  This chapter describes the administration interfaces for DCGM.
 *  It is the user's responsibility to call \ref dcgmInit() before calling any other methods,
 *  and \ref dcgmShutdown() once DCGM is no longer being used. The APIs in Administrative module
 *  can be broken down into following categories:
 *  @{
 */
/***************************************************************************************************/

/***************************************************************************************************/
/** @defgroup DCGMAPI_Admin_InitShut Init and Shutdown
 *
 *  Describes APIs to Initialize and Shutdown the DCGM Engine.
 *  @{
 */
/***************************************************************************************************/

/**
 * This method is used to initialize DCGM within this process. This must be called before
 * dcgmStartEmbedded() or dcgmConnect()
 *
 *  * @return
 *        - \ref DCGM_ST_OK                   if DCGM has been properly initialized
 *        - \ref DCGM_ST_INIT_ERROR           if there was an error initializing the library
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmInit(void);

/**
 * This method is used to shut down DCGM. Any embedded host engines or remote connections will automatically
 * be shut down as well.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if DCGM has been properly shut down
 *        - \ref DCGM_ST_UNINITIALIZED        if the library was not shut down properly
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmShutdown(void);

/**
 * Start an embedded host engine agent within this process.
 *
 * The agent is loaded as a shared library. This mode is provided to avoid any
 * extra jitter associated with an additional autonomous agent needs to be managed. In
 * this mode, the user has to periodically call APIs such as \ref dcgmPolicyTrigger and
 * \ref dcgmUpdateAllFields which tells DCGM to wake up and perform data collection and
 * operations needed for policy management.
 *
 * @param opMode       IN: Collect data automatically or manually when asked by the user.
 * @param pDcgmHandle OUT: DCGM Handle to use for API calls
 *
 * @return
 *         - \ref DCGM_ST_OK                if DCGM was started successfully within our process
 *         - \ref DCGM_ST_UNINITIALIZED     if DCGM has not been initialized with \ref dcgmInit yet
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStartEmbedded(dcgmOperationMode_t opMode, dcgmHandle_t *pDcgmHandle);

/**
 * Start an embedded host engine agent within this process.
 *
 * The agent is loaded as a shared library. This mode is provided to avoid any
 * extra jitter associated with an additional autonomous agent needs to be managed. In
 * this mode, the user has to periodically call APIs such as \c dcgmPolicyTrigger and
 * \c dcgmUpdateAllFields which tells DCGM to wake up and perform data collection and
 * operations needed for policy management.
 *
 * @param[in,out] params    A pointer to either \c dcgmStartEmbeddedV2Params_v1 or \c dcgmStartEmbeddedV2Params_v2.
 *
 * @return \c DCGM_ST_OK                if DCGM was started successfully within our process
 * @return \c DCGM_ST_UNINITIALIZED     if DCGM has not been initialized with \c dcgmInit yet
 * @note This function has a versioned argument that can be actually called with two different types. The behavior will
 *       depend on the params->version value.
 * @see dcgmStartEmbeddedV2Params_v1
 * @see dcgmStartEmbeddedV2Params_v2
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStartEmbedded_v2(dcgmStartEmbeddedV2Params_v1 *params);

/**
 * Stop the embedded host engine within this process that was started with dcgmStartEmbedded
 *
 * @param pDcgmHandle IN : DCGM Handle of the embedded host engine that came from dcgmStartEmbedded
 *
 * @return
 *         - \ref DCGM_ST_OK                if DCGM was stopped successfully within our process
 *         - \ref DCGM_ST_UNINITIALIZED     if DCGM has not been initialized with \ref dcgmInit or
 *                                          the embedded host engine was not running.
 *         - \ref DCGM_ST_BADPARAM          if an invalid parameter was provided
 *         - \ref DCGM_ST_INIT_ERROR        if an error occurred while trying to start the host engine.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStopEmbedded(dcgmHandle_t pDcgmHandle);

/**
 * This method is used to connect to a stand-alone host engine process. Remote host engines are started
 * by running the nv-hostengine command.
 *
 * NOTE: dcgmConnect_v2 provides additional connection options.
 *
 * @param ipAddress    IN: Valid IP address for the remote host engine to connect to.
 *                         If ipAddress is specified as x.x.x.x it will attempt to connect to the default
 *                         port specified by DCGM_HE_PORT_NUMBER
 *                         If ipAddress is specified as x.x.x.x:yyyy it will attempt to connect to the
 *                         port specified by yyyy
 * @param pDcgmHandle OUT: DCGM Handle of the remote host engine
 *
 * @return
 *         - \ref DCGM_ST_OK                   if we successfully connected to the remote host engine
 *         - \ref DCGM_ST_CONNECTION_NOT_VALID if the remote host engine could not be reached
 *         - \ref DCGM_ST_UNINITIALIZED        if DCGM has not been initialized with \ref dcgmInit.
 *         - \ref DCGM_ST_BADPARAM             if pDcgmHandle is NULL or ipAddress is invalid
 *         - \ref DCGM_ST_INIT_ERROR           if DCGM encountered an error while initializing the remote client library
 *         - \ref DCGM_ST_UNINITIALIZED        if DCGM has not been initialized with \ref dcgmInit
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmConnect(const char *ipAddress, dcgmHandle_t *pDcgmHandle);

/**
 * This method is used to connect to a stand-alone host engine process. Remote host engines are started
 * by running the nv-hostengine command.
 *
 * @param ipAddress     IN: Valid IP address for the remote host engine to connect to.
 *                          If ipAddress is specified as x.x.x.x it will attempt to connect to the default port
 *                          specified by DCGM_HE_PORT_NUMBER.
 *                          If ipAddress is specified as x.x.x.x:yyyy it will attempt to connect to the port
 *                          specified by yyyy
 * @param connectParams IN: Additional connection parameters. See \ref dcgmConnectV2Params_t for details.
 * @param pDcgmHandle  OUT: DCGM Handle of the remote host engine
 *
 * @return
 *         - \ref DCGM_ST_OK                   if we successfully connected to the remote host engine
 *         - \ref DCGM_ST_CONNECTION_NOT_VALID if the remote host engine could not be reached
 *         - \ref DCGM_ST_UNINITIALIZED        if DCGM has not been initialized with \ref dcgmInit.
 *         - \ref DCGM_ST_BADPARAM             if pDcgmHandle is NULL or ipAddress is invalid
 *         - \ref DCGM_ST_INIT_ERROR           if DCGM encountered an error while initializing the remote client library
 *         - \ref DCGM_ST_UNINITIALIZED        if DCGM has not been initialized with \ref dcgmInit
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmConnect_v2(const char *ipAddress,
                                            dcgmConnectV2Params_t *connectParams,
                                            dcgmHandle_t *pDcgmHandle);

/**
 * This method is used to disconnect from a stand-alone host engine process.
 *
 * @param pDcgmHandle IN: DCGM Handle that came from dcgmConnect
 *
 * @return
 *         - \ref DCGM_ST_OK                if we successfully disconnected from the host engine
 *         - \ref DCGM_ST_UNINITIALIZED     if DCGM has not been initialized with \ref dcgmInit
 *         - \ref DCGM_ST_BADPARAM          if pDcgmHandle is not a valid DCGM handle
 *         - \ref DCGM_ST_GENERIC_ERROR     if an unspecified internal error occurred
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmDisconnect(dcgmHandle_t pDcgmHandle);


/** @} */ // Closing for DCGMAPI_Admin_InitShut

/***************************************************************************************************/
/** @defgroup DCGMAPI_Admin_Info Auxilary information about DCGM engine.
 *
 *  Describes APIs to get generic information about the DCGM Engine.
 *  @{
 */
/***************************************************************************************************/

/**
 * This method is used to return information about the build environment where DCGM was built.
 *
 * @param pVersionInfo OUT: Build environment information
 *
 * @return
 *          - \ref DCGM_ST_OK           if build information is sucessfully obtained
 *          - \ref DCGM_ST_BADPARAM     if pVersionInfo is null
 *          - \ref DCGM_ST_VER_MISMATCH if the expected and provided versions of dcgmVersionInfo_t do not match
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmVersionInfo(dcgmVersionInfo_t *pVersionInfo);

/**
 * This method is used to return information about the build environment of the hostengine.
 *
 * @param pDcgmHandle  IN:  DCGM Handle that came from dcgmConnect
 * @param pVersionInfo OUT: Build environment information
 *
 * @return
 *          - \ref DCGM_ST_OK           if build information is sucessfully obtained
 *          - \ref DCGM_ST_BADPARAM     if pVersionInfo is null
 *          - \ref DCGM_ST_VER_MISMATCH if the expected and provided versions of dcgmVersionInfo_t do not match
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmHostengineVersionInfo(dcgmHandle_t pDcgmHandle, dcgmVersionInfo_t *pVersionInfo);


/**
 * This method is used to set the logging severity on HostEngine for the specified logger
 *
 * @param pDcgmHandle  IN: DCGM Handle
 * @param logging      IN: dcgmSettingsSetLoggingSeverity_t struct containing the target logger and severity
 *
 * @return
 *          - \ref DCGM_ST_OK           Severity successfuly set
 *          - \ref DCGM_ST_BADPARAM     Bad logger/severity string
 *          - \ref DCGM_ST_VER_MISMATCH if the expected and provided versions of dcgmSettingsSetLoggingSeverity_t
 *                                      do not match
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmHostengineSetLoggingSeverity(dcgmHandle_t pDcgmHandle,
                                                              dcgmSettingsSetLoggingSeverity_t *logging);

/**
 * This function is used to return whether or not the host engine considers itself healthy
 *
 * @param[in]  pDcgmHandle - the handle to DCGM
 * @param[out] heHealth - struct describing the health of the hostengine. if heHealth.hostengineHealth is 0,
 *                        then the hostengine is healthy. Non-zero indicates not healthy with error codes
 *                        determining the cause.
 *
 * @return
 *          - \ref DCGM_ST_OK         Able to gauge health
 *          - \ref DCGM_ST_BADPARAM   isHealthy is not a valid pointer
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmHostengineIsHealthy(dcgmHandle_t pDcgmHandle, dcgmHostengineHealth_t *heHealth);


/**
 * This function describes DCGM error codes in human readable form
 *
 * @param[in] result    - DCGM return code to describe
 *
 * @return
 *          - Human readable string with the DCGM error code description if the code is valid.
 *          - nullptr if there is not such error code
 */
DCGM_PUBLIC_API const char *errorString(dcgmReturn_t result);

/**
 * This function describes DCGM Module by given Module ID
 *
 * @param id[in]        - Module ID to name.
 * @param name[out]     - Module name will be provided via this argument.
 * @return
 *          - \ref DCGM_ST_OK           Module name has valid value
 *          - \ref DCGM_ST_BADPARAM     There is no module with specified ID. Name value is not changed.
 */
DCGM_PUBLIC_API dcgmReturn_t dcgmModuleIdToName(dcgmModuleId_t id, char const **name);

/** @} */ // Closing DCGMAPI_Admin_Info

/** @} */ // Closing for DCGMAPI_Admin


/***************************************************************************************************/
/** @defgroup DCGMAPI_SYS System
 *  @{
 *  This chapter describes the APIs used to identify entities on the node, grouping functions to
 *  provide mechanism to operate on a group of entities, and status management APIs in
 *  order to get individual statuses for each operation. The APIs in System module can be
 *  broken down into following categories:
 */
/***************************************************************************************************/

/***************************************************************************************************/
/** @defgroup DCGM_DISCOVERY Discovery
 *  The following APIs are used to discover GPUs and their attributes on a Node.
 *  @{
 */
/***************************************************************************************************/

/**
 * This method is used to get identifiers corresponding to all the devices on the system. The
 * identifier represents DCGM GPU Id corresponding to each GPU on the system and is immutable during
 * the lifespan of the engine. The list should be queried again if the engine is restarted.
 *
 * The GPUs returned from this function include gpuIds of GPUs that are not supported by DCGM.
 * To only get gpuIds of GPUs that are supported by DCGM, use dcgmGetAllSupportedDevices().
 *
 * @param pDcgmHandle                    IN: DCGM Handle
 * @param gpuIdList                     OUT: Array reference to fill GPU Ids present on the system.
 * @param count                         OUT: Number of GPUs returned in \a gpuIdList.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful.
 *        - \ref DCGM_ST_BADPARAM             if \a gpuIdList or \a count were not valid.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetAllDevices(dcgmHandle_t pDcgmHandle,
                                               unsigned int gpuIdList[DCGM_MAX_NUM_DEVICES],
                                               int *count);

/**
 * This method is used to get identifiers corresponding to all the DCGM-supported devices on the system. The
 * identifier represents DCGM GPU Id corresponding to each GPU on the system and is immutable during
 * the lifespan of the engine. The list should be queried again if the engine is restarted.
 *
 * The GPUs returned from this function ONLY includes gpuIds of GPUs that are supported by DCGM.
 * To get gpuIds of all GPUs in the system, use dcgmGetAllDevices().
 *
 *
 * @param pDcgmHandle                    IN: DCGM Handle
 * @param gpuIdList                     OUT: Array reference to fill GPU Ids present on the system.
 * @param count                         OUT: Number of GPUs returned in \a gpuIdList.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful.
 *        - \ref DCGM_ST_BADPARAM             if \a gpuIdList or \a count were not valid.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetAllSupportedDevices(dcgmHandle_t pDcgmHandle,
                                                        unsigned int gpuIdList[DCGM_MAX_NUM_DEVICES],
                                                        int *count);

/**
 * Gets device attributes corresponding to the \a gpuId. If operation is not successful for any of
 * the requested fields then the field is populated with one of DCGM_BLANK_VALUES defined in
 * dcgm_structs.h.
 *
 * @param pDcgmHandle    IN: DCGM Handle
 * @param gpuId          IN: GPU Id corresponding to which the attributes should be fetched
 * @param pDcgmAttr  IN/OUT: Device attributes corresponding to \a gpuId.<br> pDcgmAttr->version should be set to
 *                           \ref dcgmDeviceAttributes_version before this call.
 *
 * @return
 *        - \ref DCGM_ST_OK            if the call was successful.
 *        - \ref DCGM_ST_VER_MISMATCH  if pDcgmAttr->version is not set or is invalid.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetDeviceAttributes(dcgmHandle_t pDcgmHandle,
                                                     unsigned int gpuId,
                                                     dcgmDeviceAttributes_t *pDcgmAttr);

/**
 * Gets the list of entities that exist for a given entity group. This API can be used in place of
 * \ref dcgmGetAllDevices.
 *
 * @param dcgmHandle      IN: DCGM Handle
 * @param entityGroup     IN: Entity group to list entities of
 * @param entities       OUT: Array of entities for entityGroup
 * @param numEntities IN/OUT: Upon calling, this should be the number of entities that entityList[] can hold. Upon
 *                            return, this will contain the number of entities actually saved to entityList.
 * @param flags           IN: Flags to modify the behavior of this request.
 *                            See DCGM_GEGE_FLAG_* #defines in dcgm_structs.h
 *
 * @return
 *        - \ref DCGM_ST_OK                if the call was successful.
 *        - \ref DCGM_ST_INSUFFICIENT_SIZE if numEntities was not large enough to hold the number of entities in the
 *                                         entityGroup. numEntities will contain the capacity needed to complete this
 *                                         request successfully.
 *        - \ref DCGM_ST_NOT_SUPPORTED     if the given entityGroup does not support enumeration.
 *        - \ref DCGM_ST_BADPARAM          if any parameter is invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetEntityGroupEntities(dcgmHandle_t dcgmHandle,
                                                        dcgm_field_entity_group_t entityGroup,
                                                        dcgm_field_eid_t *entities,
                                                        int *numEntities,
                                                        unsigned int flags);

/**
 * Gets the hierarchy of GPUs, GPU Instances, and Compute Instances by populating a list of each entity with
 * a reference to their parent
 *
 * @param dcgmHandle       IN: DCGM Handle
 * @param entities        OUT: array of entities in the hierarchy
 * @param numEntities  IN/OUT: Upon calling, this should be the capacity of entities.
 *                             Upon return, this will contain the number of entities actually saved to entities.
 *
 * @return
 *        - \ref DCGM_ST_OK                if the call was successful.
 *        - \ref DCGM_ST_VER_MISMATCH      if the struct version is incorrect
 *        - \ref DCGM_ST_BADPARAM          if any parameter is invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetGpuInstanceHierarchy(dcgmHandle_t dcgmHandle, dcgmMigHierarchy_v2 *hierarchy);

/**
 * Get the NvLink link status for every NvLink in this system. This includes the NvLinks of both GPUs and
 * NvSwitches. Note that only NvSwitches and GPUs that are visible to the current environment will be
 * returned in this structure.
 *
 * @param dcgmHandle  IN: DCGM Handle
 * @param linkStatus OUT: Structure in which to store NvLink link statuses. .version should be set to
 *                        dcgmNvLinkStatus_version1 before calling this.
 *
 * @return
 *        - \ref DCGM_ST_OK                if the call was successful.
 *        - \ref DCGM_ST_NOT_SUPPORTED     if the given entityGroup does not support enumeration.
 *        - \ref DCGM_ST_BADPARAM          if any parameter is invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetNvLinkLinkStatus(dcgmHandle_t dcgmHandle, dcgmNvLinkStatus_v3 *linkStatus);


/**
 * List supported CPUs and their cores present on the system
 *
 * This and other CPU APIs only support datacenter NVIDIA CPUs
 *
 * @param dcgmHandle   IN: DCGM Handle
 * @param cpuHierarchy OUT: Structure where the CPUs and their associated cores will be enumerated
 *
 * @return
 *        - \ref DCGM_ST_OK                if the call was successful.
 *        - \ref DCGM_ST_NOT_SUPPORTED     if the device is unsupported
 *        - \ref DCGM_ST_MODULE_NOT_LOADED if the sysmon module could not be loaded
 *        - \ref DCGM_ST_BADPARAM          if any parameter is invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetCpuHierarchy(dcgmHandle_t dcgmHandle, dcgmCpuHierarchy_v1 *cpuHierarchy);

/** @} */

/***************************************************************************************************/
/** @defgroup DCGM_GROUPING Grouping
 *  The following APIs are used for group management. The user can create a group of entities and
 *  perform an operation on a group of entities. If grouping is not needed and the user wishes
 *  to run commands on all GPUs seen by DCGM then the user can use DCGM_GROUP_ALL_GPUS or
 *  DCGM_GROUP_ALL_NVSWITCHES in place of group IDs when needed.
 *  @{
 */
/***************************************************************************************************/

/**
 * Used to create a entity group handle which can store one or more entity Ids as an opaque handle
 * returned in \a pDcgmGrpId. Instead of executing an operation separately for each entity, the
 * DCGM group enables the user to execute same operation on all the entities present in the group as a
 * single API call.
 *
 * To create the group with all the entities present on the system, the \a type field should be
 * specified as \a DCGM_GROUP_DEFAULT or \a DCGM_GROUP_ALL_NVSWITCHES. To create an empty group,
 * the \a type field should be specified as \a DCGM_GROUP_EMPTY. The empty group can be updated
 * with the desired set of entities using the APIs \ref dcgmGroupAddDevice, \ref dcgmGroupAddEntity,
 * \ref dcgmGroupRemoveDevice, and \ref dcgmGroupRemoveEntity.
 *
 * @param pDcgmHandle    IN: DCGM Handle
 * @param type           IN: Type of Entity Group to be formed
 * @param groupName      IN: Desired name of the GPU group specified as NULL terminated C string
 * @param pDcgmGrpId    OUT: Reference to group ID
 *
 * @return
 *  - \ref DCGM_ST_OK                if the group has been created
 *  - \ref DCGM_ST_BADPARAM          if any of \a type, \a groupName, \a length or \a pDcgmGrpId is invalid
 *  - \ref DCGM_ST_MAX_LIMIT         if number of groups on the system has reached the max limit \a DCGM_MAX_NUM_GROUPS
 *  - \ref DCGM_ST_INIT_ERROR        if the library has not been successfully initialized
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupCreate(dcgmHandle_t pDcgmHandle,
                                             dcgmGroupType_t type,
                                             const char *groupName,
                                             dcgmGpuGrp_t *pDcgmGrpId);

/**
 * Used to destroy a group represented by \a groupId.
 * Since DCGM group is a logical grouping of entities, the properties applied on the group stay intact
 * for the individual entities even after the group is destroyed.
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param groupId       IN: Group ID
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the group has been destroyed
 *  - \ref DCGM_ST_BADPARAM             if \a groupId is invalid
 *  - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized
 *  - \ref DCGM_ST_NOT_CONFIGURED       if entry corresponding to the group does not exists
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupDestroy(dcgmHandle_t pDcgmHandle, dcgmGpuGrp_t groupId);

/**
 * Used to add specified GPU Id to the group represented by \a groupId.
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param groupId       IN: Group Id to which device should be added
 * @param gpuId         IN: DCGM GPU Id
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the GPU Id has been successfully added to the group
 *  - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized
 *  - \ref DCGM_ST_NOT_CONFIGURED       if entry corresponding to the group (\a groupId) does not exists
 *  - \ref DCGM_ST_BADPARAM             if \a gpuId is invalid or already part of the specified group
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupAddDevice(dcgmHandle_t pDcgmHandle, dcgmGpuGrp_t groupId, unsigned int gpuId);

/**
 * Used to add specified entity to the group represented by \a groupId.
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param groupId       IN: Group Id to which device should be added
 * @param entityGroupId IN: Entity group that entityId belongs to
 * @param entityId      IN: DCGM entityId
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the entity has been successfully added to the group
 *  - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized
 *  - \ref DCGM_ST_NOT_CONFIGURED       if entry corresponding to the group (\a groupId) does not exists
 *  - \ref DCGM_ST_BADPARAM             if \a entityId is invalid or already part of the specified group
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupAddEntity(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupId,
                                                dcgm_field_entity_group_t entityGroupId,
                                                dcgm_field_eid_t entityId);

/**
 * Used to remove specified GPU Id from the group represented by \a groupId.
 * @param pDcgmHandle   IN: DCGM Handle
 * @param groupId       IN: Group ID from which device should be removed
 * @param gpuId         IN: DCGM GPU Id
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the GPU Id has been successfully removed from the group
 *  - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized
 *  - \ref DCGM_ST_NOT_CONFIGURED       if entry corresponding to the group (\a groupId) does not exists
 *  - \ref DCGM_ST_BADPARAM             if \a gpuId is invalid or not part of the specified group
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupRemoveDevice(dcgmHandle_t pDcgmHandle, dcgmGpuGrp_t groupId, unsigned int gpuId);

/**
 * Used to remove specified entity from the group represented by \a groupId.
 * @param pDcgmHandle   IN: DCGM Handle
 * @param groupId       IN: Group ID from which device should be removed
 * @param entityGroupId IN: Entity group that entityId belongs to
 * @param entityId      IN: DCGM entityId
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the entity has been successfully removed from the group
 *  - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized
 *  - \ref DCGM_ST_NOT_CONFIGURED       if entry corresponding to the group (\a groupId) does not exists
 *  - \ref DCGM_ST_BADPARAM             if \a entityId is invalid or not part of the specified group
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupRemoveEntity(dcgmHandle_t pDcgmHandle,
                                                   dcgmGpuGrp_t groupId,
                                                   dcgm_field_entity_group_t entityGroupId,
                                                   dcgm_field_eid_t entityId);

/**
 * Used to get information corresponding to the group represented by \a groupId. The information
 * returned in \a pDcgmGroupInfo consists of group name, and the list of entities present in the
 * group.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID for which information to be fetched
 * @param pDcgmGroupInfo    OUT: Group Information
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the group info is successfully received.
 *  - \ref DCGM_ST_BADPARAM             if any of \a groupId or \a pDcgmGroupInfo is invalid.
 *  - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized.
 *  - \ref DCGM_ST_MAX_LIMIT            if the group does not contain the GPU
 *  - \ref DCGM_ST_NOT_CONFIGURED       if entry corresponding to the group (\a groupId) does not exists
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupGetInfo(dcgmHandle_t pDcgmHandle,
                                              dcgmGpuGrp_t groupId,
                                              dcgmGroupInfo_t *pDcgmGroupInfo);

/**
 * Used to get the Ids of all groups of entities. The information returned is a list of group ids
 * in \a groupIdList as well as a count of how many ids there are in \a count. Please allocate enough
 * memory for \a groupIdList. Memory of size MAX_NUM_GROUPS should be allocated for \a groupIdList.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupIdList       OUT: List of Group Ids
 * @param count             OUT: The number of Group ids in the list
 *
 * @return
 *  - \ref DCGM_ST_OK               if the ids of the groups were successfully retrieved
 *  - \ref DCGM_ST_BADPARAM         if either of the \a groupIdList or \a count is null
 *  - \ref DCGM_ST_GENERIC_ERROR    if an unknown error has occurred
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGroupGetAllIds(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupIdList[],
                                                unsigned int *count);

/** @} */

/***************************************************************************************************/
/** @defgroup DCGM_FIELD_GROUPING Field Grouping
 *  The following APIs are used for field group management. The user can create a group of fields and
 *  perform an operation on a group of fields at once.
 *  @{
 */

/**
 * Used to create a group of fields and return the handle in dcgmFieldGroupId
 *
 * @param dcgmHandle         IN: DCGM handle
 * @param numFieldIds        IN: Number of field IDs that are being provided in fieldIds[]. Must be between 1 and
 *                               DCGM_MAX_FIELD_IDS_PER_FIELD_GROUP.
 * @param fieldIds           IN: Field IDs to be added to the newly-created field group
 * @param fieldGroupName     IN: Unique name for this group of fields. This must not be the same as any existing field
 *                               groups.
 * @param dcgmFieldGroupId  OUT: Handle to the newly-created field group
 *
 * @return
 * - \ref DCGM_ST_OK                   if the field group was successfully created.
 * - \ref DCGM_ST_BADPARAM             if any parameters were bad
 * - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized.
 * - \ref DCGM_ST_MAX_LIMIT            if too many field groups already exist
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmFieldGroupCreate(dcgmHandle_t dcgmHandle,
                                                  int numFieldIds,
                                                  unsigned short *fieldIds,
                                                  const char *fieldGroupName,
                                                  dcgmFieldGrp_t *dcgmFieldGroupId);

/**
 * Used to remove a field group that was created with \ref dcgmFieldGroupCreate
 *
 * @param dcgmHandle         IN: DCGM handle
 * @param dcgmFieldGroupId   IN: Field group to remove
 *
 * @return
 * - \ref DCGM_ST_OK                   if the field group was successfully removed
 * - \ref DCGM_ST_BADPARAM             if any parameters were bad
 * - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmFieldGroupDestroy(dcgmHandle_t dcgmHandle, dcgmFieldGrp_t dcgmFieldGroupId);


/**
 * Used to get information about a field group that was created with \ref dcgmFieldGroupCreate.
 *
 * @param dcgmHandle         IN: DCGM handle
 * @param fieldGroupInfo IN/OUT: Info about all of the field groups that exist.<br>
 *                               .version should be set to \ref dcgmFieldGroupInfo_version before this call<br>
 *                               .fieldGroupId should contain the fieldGroupId you are interested in querying
 *                               information for.
 *
 * @return
 * - \ref DCGM_ST_OK                   if the field group info was returned successfully
 * - \ref DCGM_ST_BADPARAM             if any parameters were bad
 * - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized.
 * - \ref DCGM_ST_VER_MISMATCH         if .version is not set or is invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmFieldGroupGetInfo(dcgmHandle_t dcgmHandle, dcgmFieldGroupInfo_t *fieldGroupInfo);

/**
 * Used to get information about all field groups in the system.
 *
 * @param dcgmHandle         IN: DCGM handle
 * @param allGroupInfo   IN/OUT: Info about all of the field groups that exist.<br>
 *                               .version should be set to \ref dcgmAllFieldGroup_version before this call.
 *
 * @return
 * - \ref DCGM_ST_OK                   if the field group info was successfully returned
 * - \ref DCGM_ST_BADPARAM             if any parameters were bad
 * - \ref DCGM_ST_INIT_ERROR           if the library has not been successfully initialized.
 * - \ref DCGM_ST_VER_MISMATCH         if .version is not set or is invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmFieldGroupGetAll(dcgmHandle_t dcgmHandle, dcgmAllFieldGroup_t *allGroupInfo);

/** @} */


/***************************************************************************************************/
/** @defgroup DCGMAPI_ST Status handling
 * The following APIs are used to manage statuses for multiple operations on one or more GPUs.
 *  @{
 */
/***************************************************************************************************/

/**
 * Creates reference to DCGM status handler which can be used to get the statuses for multiple
 * operations on one or more devices.
 *
 * The multiple statuses are useful when the operations are performed at group level. The status
 * handle provides a mechanism to access error attributes for the failed operations.
 *
 * The number of errors stored behind the opaque handle can be accessed using the the API
 * \ref dcgmStatusGetCount. The errors are accessed from the opaque handle \a statusHandle
 * using the API \ref dcgmStatusPopError. The user can invoke \ref dcgmStatusPopError
 * for the number of errors or until all the errors are fetched.
 *
 * When the status handle is not required any further then it should be deleted using the API
 * \ref dcgmStatusDestroy.
 * @param statusHandle   OUT: Reference to handle for list of statuses
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the status handle is successfully created
 *  - \ref DCGM_ST_BADPARAM             if \a statusHandle is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStatusCreate(dcgmStatus_t *statusHandle);

/**
 * Used to destroy status handle created using \ref dcgmStatusCreate.
 * @param statusHandle   IN: Handle to list of statuses
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the status handle is successfully created
 *  - \ref DCGM_ST_BADPARAM             if \a statusHandle is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStatusDestroy(dcgmStatus_t statusHandle);

/**
 * Used to get count of error entries stored inside the opaque handle \a statusHandle.
 * @param statusHandle   IN: Handle to list of statuses
 * @param count         OUT: Number of error entries present in the list of statuses
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the error count is successfully received
 *  - \ref DCGM_ST_BADPARAM             if any of \a statusHandle or \a count is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStatusGetCount(dcgmStatus_t statusHandle, unsigned int *count);

/**
 * Used to iterate through the list of errors maintained behind \a statusHandle. The method pops the
 * first error from the list of DCGM statuses. In order to iterate through all the errors, the user
 * can invoke this API for the number of errors or until all the errors are fetched.
 * @param statusHandle       IN: Handle to list of statuses
 * @param pDcgmErrorInfo    OUT: First error from the list of statuses
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the error entry is successfully fetched
 *  - \ref DCGM_ST_BADPARAM             if any of \a statusHandle or \a pDcgmErrorInfo is invalid
 *  - \ref DCGM_ST_NO_DATA              if the status handle list is empty
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStatusPopError(dcgmStatus_t statusHandle, dcgmErrorInfo_t *pDcgmErrorInfo);

/**
 * Used to clear all the errors in the status handle created by the API
 * \ref dcgmStatusCreate. After one set of operation, the \a statusHandle
 * can be cleared and reused for the next set of operation.
 * @param statusHandle   IN: Handle to list of statuses
 *
 * @return
 *  - \ref DCGM_ST_OK                   if the errors are successfully cleared
 *  - \ref DCGM_ST_BADPARAM             if \a statusHandle is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmStatusClear(dcgmStatus_t statusHandle);

/** @} */ // Closing for DCGMAPI_ST


/** @} */ // Closing for DCGMAPI_SYS

/***************************************************************************************************/
/** @defgroup DCGMAPI_DC Configuration
 *  This chapter describes the methods that handle device configuration retrieval and
 *  default settings. The APIs in Configuration module can be broken down into following
 *  categories:
 *  @{
 */
/***************************************************************************************************/

/***************************************************************************************************/
/** @defgroup DCGMAPI_DC_Setup Setup and management
 *  Describes APIs to Get/Set configuration on the group of GPUs.
 *  @{
 */
/***************************************************************************************************/

/**
* Used to set configuration for the group of one or more GPUs identified by \a groupId.
*
* The configuration settings specified in \a pDeviceConfig are applied to all the GPUs in the
* group. Since DCGM group is a logical grouping of GPUs, the configuration settings stays intact
* for the individual GPUs even after the group is destroyed.
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
* To find out valid supported clock values that can be passed to dcgmConfigSet, look at the device
* attributes of a GPU in the group using the API dcgmGetDeviceAttributes.

* @param pDcgmHandle            IN: DCGM Handle
* @param groupId                IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate
*                                   for details on creating the group.
* @param pDeviceConfig          IN: Pointer to memory to hold desired configuration to be applied for all the GPU in the
*                                   group represented by \a groupId.
*                                   The caller must populate the version field of \a pDeviceConfig.
* @param statusHandle       IN/OUT: Resulting error status for multiple operations. Pass it as NULL if the detailed
*                                   error information is not needed.
*                                   Look at \ref dcgmStatusCreate for details on creating status handle.

* @return
*        - \ref DCGM_ST_OK                   if the configuration has been successfully set.
*        - \ref DCGM_ST_BADPARAM             if any of \a groupId or \a pDeviceConfig is invalid.
*        - \ref DCGM_ST_VER_MISMATCH         if \a pDeviceConfig has the incorrect version.
*        - \ref DCGM_ST_GENERIC_ERROR        if an unknown error has occurred.
*
*/
dcgmReturn_t DCGM_PUBLIC_API dcgmConfigSet(dcgmHandle_t pDcgmHandle,
                                           dcgmGpuGrp_t groupId,
                                           dcgmConfig_t *pDeviceConfig,
                                           dcgmStatus_t statusHandle);

/**
* Used to get configuration for all the GPUs present in the group.
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
* @param pDcgmHandle            IN: DCGM Handle
* @param groupId                IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate
*                                   for details on creating the group.
* @param type                   IN: Type of configuration values to be fetched.
* @param count                  IN: The number of entries that \a deviceConfigList array can store.
* @param deviceConfigList      OUT: Pointer to memory to hold requested configuration corresponding to all the GPUs in
*                                   the group (\a groupId). The size of the memory must be greater than or equal to hold
*                                   output information for the number of GPUs present in the group (\a groupId).
* @param statusHandle       IN/OUT: Resulting error status for multiple operations. Pass it as NULL if the detailed
*                                   error information is not needed.
*                                   Look at \ref dcgmStatusCreate for details on creating status handle.

* @return
*        - \ref DCGM_ST_OK                   if the configuration has been successfully fetched.
*        - \ref DCGM_ST_BADPARAM             if any of \a groupId, \a type, \a count, or \a deviceConfigList is invalid.
*        - \ref DCGM_ST_NOT_CONFIGURED       if the target configuration is not already set.
*        - \ref DCGM_ST_VER_MISMATCH         if \a deviceConfigList has the incorrect version.
*        - \ref DCGM_ST_GENERIC_ERROR        if an unknown error has occurred.
*
*/
dcgmReturn_t DCGM_PUBLIC_API dcgmConfigGet(dcgmHandle_t pDcgmHandle,
                                           dcgmGpuGrp_t groupId,
                                           dcgmConfigType_t type,
                                           int count,
                                           dcgmConfig_t deviceConfigList[],
                                           dcgmStatus_t statusHandle);

/** @} */ // Closing for DCGMAPI_DC_Setup


/***************************************************************************************************/
/** @defgroup DCGMAPI_DC_MI Manual Invocation
 *  Describes APIs used to manually enforce the desired configuration on a group of GPUs.
 *  @{
 */
/***************************************************************************************************/

/**
 * Used to enforce previously set configuration for all the GPUs present in the group.
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
 * @param pDcgmHandle            IN: DCGM Handle
 * @param groupId                IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate
 *                                   for details on creating the group. Alternatively, pass in the group id as
 *                                   \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param statusHandle       IN/OUT: Resulting error status for multiple operations. Pass it as NULL if the detailed
 *                                   error information is not needed. Look at \ref dcgmStatusCreate for details on
 *                                   creating status handle.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the configuration has been successfully enforced.
 *        - \ref DCGM_ST_BADPARAM             if \a groupId is invalid.
 *        - \ref DCGM_ST_NOT_CONFIGURED       if the target configuration is not already set.
 *        - \ref DCGM_ST_GENERIC_ERROR        if an unknown error has occurred.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmConfigEnforce(dcgmHandle_t pDcgmHandle,
                                               dcgmGpuGrp_t groupId,
                                               dcgmStatus_t statusHandle);

/** @} */ // Closing for DCGMAPI_DC_MI

/** @} */ // Closing for DCGMAPI_DC

/***************************************************************************************************/
/** @defgroup DCGMAPI_FI Field APIs
 *
 *   These APIs are responsible for watching, unwatching, and updating specific fields as defined
 *   by DCGM_FI_*
 *
 *  @{
 */
/***************************************************************************************************/

/**
 * Request that DCGM start recording updates for a given field collection.
 *
 * Note that the first update of the field will not occur until the next field update cycle.
 * To force a field update cycle, call dcgmUpdateAllFields(1).
 *
 * @param pDcgmHandle         IN: DCGM Handle
 * @param groupId             IN: Group ID representing collection of one or more entities. Look at \ref dcgmGroupCreate
 *                                for details on creating the group. Alternatively, pass in the group id as
 *                                \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs or
 *                                \a DCGM_GROUP_ALL_NVSWITCHES to to perform the operation on all NvSwitches.
 * @param fieldGroupId        IN: Fields to watch.
 * @param updateFreq          IN: How often to update this field in usec
 * @param maxKeepAge          IN: How long to keep data for this field in seconds
 * @param maxKeepSamples      IN: Maximum number of samples to keep. 0=no limit
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */

dcgmReturn_t DCGM_PUBLIC_API dcgmWatchFields(dcgmHandle_t pDcgmHandle,
                                             dcgmGpuGrp_t groupId,
                                             dcgmFieldGrp_t fieldGroupId,
                                             long long updateFreq,
                                             double maxKeepAge,
                                             int maxKeepSamples);

/**
 * Request that DCGM stop recording updates for a given field collection.
 *
 * @param pDcgmHandle         IN: DCGM Handle
 * @param groupId             IN: Group ID representing collection of one or more entities. Look at \ref dcgmGroupCreate
 *                                for details on creating the group. Alternatively, pass in the group id as
 *                                \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs or
 *                                \a DCGM_GROUP_ALL_NVSWITCHES to to perform the operation on all NvSwitches.
 * @param fieldGroupId        IN: Fields to unwatch.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmUnwatchFields(dcgmHandle_t pDcgmHandle,
                                               dcgmGpuGrp_t groupId,
                                               dcgmFieldGrp_t fieldGroupId);

/**
 * Request updates for all field values that have updated since a given timestamp
 *
 * This version only works with GPU entities. Use \ref dcgmGetValuesSince_v2 for entity groups
 * containing NvSwitches.
 *
 * @param pDcgmHandle         IN: DCGM Handle
 * @param groupId             IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                                details on creating the group. Alternatively, pass in the group id as
 *                                \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param fieldGroupId        IN: Fields to return data for
 * @param sinceTimestamp      IN: Timestamp to request values since in usec since 1970. This will be returned in
 *                                nextSinceTimestamp for subsequent calls 0 = request all data
 * @param nextSinceTimestamp OUT: Timestamp to use for sinceTimestamp on next call to this function
 * @param enumCB              IN: Callback to invoke for every field value update. Note that multiple updates can be
 *                                returned in each invocation
 * @param userData            IN: User data pointer to pass to the userData field of enumCB.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_NOT_SUPPORTED        if one of the entities was from a non-GPU type
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetValuesSince(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupId,
                                                dcgmFieldGrp_t fieldGroupId,
                                                long long sinceTimestamp,
                                                long long *nextSinceTimestamp,
                                                dcgmFieldValueEnumeration_f enumCB,
                                                void *userData);

/**
 * Request updates for all field values that have updated since a given timestamp
 *
 * This version works with non-GPU entities like NvSwitches
 *
 * @param pDcgmHandle         IN: DCGM Handle
 * @param groupId             IN: Group ID representing collection of one or more entities. Look at \ref dcgmGroupCreate
 *                                for details on creating the group. Alternatively, pass in the group id as
 *                                \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs or
 *                                \a DCGM_GROUP_ALL_NVSWITCHES to perform the operation on all NvSwitches.
 * @param fieldGroupId        IN: Fields to return data for
 * @param sinceTimestamp      IN: Timestamp to request values since in usec since 1970. This will be returned in
 *                                nextSinceTimestamp for subsequent calls 0 = request all data
 * @param nextSinceTimestamp OUT: Timestamp to use for sinceTimestamp on next call to this function
 * @param enumCB              IN: Callback to invoke for every field value update. Note that multiple updates can be
 *                                returned in each invocation
 * @param userData            IN: User data pointer to pass to the userData field of enumCB.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetValuesSince_v2(dcgmHandle_t pDcgmHandle,
                                                   dcgmGpuGrp_t groupId,
                                                   dcgmFieldGrp_t fieldGroupId,
                                                   long long sinceTimestamp,
                                                   long long *nextSinceTimestamp,
                                                   dcgmFieldValueEntityEnumeration_f enumCB,
                                                   void *userData);

/**
 * Request latest cached field value for a field value collection
 *
 * This version only works with GPU entities. Use \ref dcgmGetLatestValues_v2 for entity groups
 * containing NvSwitches.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                               details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param fieldGroupId       IN: Fields to return data for.
 * @param enumCB             IN: Callback to invoke for every field value update. Note that multiple updates can be
 *                               returned in each invocation
 * @param userData           IN: User data pointer to pass to the userData field of enumCB.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_NOT_SUPPORTED        if one of the entities was from a non-GPU type
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetLatestValues(dcgmHandle_t pDcgmHandle,
                                                 dcgmGpuGrp_t groupId,
                                                 dcgmFieldGrp_t fieldGroupId,
                                                 dcgmFieldValueEnumeration_f enumCB,
                                                 void *userData);

/**
 * Request latest cached field value for a field value collection
 *
 * This version works with non-GPU entities like NvSwitches
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more entities. Look at \ref dcgmGroupCreate
 *                               for details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs or
 *                               \a DCGM_GROUP_ALL_NVSWITCHES to perform the operation on all NvSwitches.
 * @param fieldGroupId       IN: Fields to return data for.
 * @param enumCB             IN: Callback to invoke for every field value update. Note that multiple updates can be
 *                               returned in each invocation
 * @param userData           IN: User data pointer to pass to the userData field of enumCB.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_NOT_SUPPORTED        if one of the entities was from a non-GPU type
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetLatestValues_v2(dcgmHandle_t pDcgmHandle,
                                                    dcgmGpuGrp_t groupId,
                                                    dcgmFieldGrp_t fieldGroupId,
                                                    dcgmFieldValueEntityEnumeration_f enumCB,
                                                    void *userData);

/**
 * Request latest cached field value for a GPU
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param gpuId         IN: Gpu ID representing the GPU for which the fields are being requested.
 * @param fields        IN: Field IDs to return data for. See the definitions in dcgm_fields.h that start with DCGM_FI_.
 * @param count         IN: Number of field IDs in fields[] array.
 * @param values       OUT: Latest field values for the fields in fields[].
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetLatestValuesForFields(dcgmHandle_t pDcgmHandle,
                                                          int gpuId,
                                                          unsigned short fields[],
                                                          unsigned int count,
                                                          dcgmFieldValue_v1 values[]);
/**
 * Request latest cached field value for a group of fields for a specific entity
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param entityGroup   IN: entity_group_t (e.g. switch)
 * @param entityId      IN: entity ID representing the rntity for which the fields are being requested.
 * @param fields        IN: Field IDs to return data for. See the definitions in dcgm_fields.h that start with DCGM_FI_.
 * @param count         IN: Number of field IDs in fields[] array.
 * @param values       OUT: Latest field values for the fields in fields[].
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmEntityGetLatestValues(dcgmHandle_t pDcgmHandle,
                                                       dcgm_field_entity_group_t entityGroup,
                                                       int entityId,
                                                       unsigned short fields[],
                                                       unsigned int count,
                                                       dcgmFieldValue_v1 values[]);

/**
 * Request the latest cached or live field value for a list of fields for a group of entities
 *
 * Note: The returned entities are not guaranteed to be in any order. Reordering can occur internally
 *       in order to optimize calls to the NVIDIA driver.
 *
 * @param pDcgmHandle   IN: DCGM Handle
 * @param entities      IN: List of entities to get values for
 * @param entityCount   IN: Number of entries in entities[]
 * @param fields        IN: Field IDs to return data for. See the definitions in dcgm_fields.h that start with DCGM_FI_.
 * @param fieldCount    IN: Number of field IDs in fields[] array.
 * @param flags         IN: Optional flags that affect how this request is processed. Pass \ref DCGM_FV_FLAG_LIVE_DATA
 *                          here to retrieve a live driver value rather than a cached value. See that flag's
 *                          documentation for caveats.
 * @param values       OUT: Latest field values for the fields requested. This must be able to hold entityCount *
 *                          fieldCount field value records.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmEntitiesGetLatestValues(dcgmHandle_t pDcgmHandle,
                                                         dcgmGroupEntityPair_t entities[],
                                                         unsigned int entityCount,
                                                         unsigned short fields[],
                                                         unsigned int fieldCount,
                                                         unsigned int flags,
                                                         dcgmFieldValue_v2 values[]);

/*************************************************************************/
/**
 * Get a summary of the values for a field id over a period of time.
 *
 * @param pDcgmHandle       IN: DCGM Handle
 * @param request       IN/OUT: a pointer to the struct detailing the request and containing the response
 *
 * @return
 *       - \ref DCGM_ST_OK                if the call was successful
 *       - \ref DCGM_ST_FIELD_UNSUPPORTED_BY_API if the field is not int64 or double type
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetFieldSummary(dcgmHandle_t pDcgmHandle, dcgmFieldSummaryRequest_t *request);

/** @} */

/***************************************************************************************************/
/** @addtogroup DCGMAPI_Admin_ExecCtrl
 *  @{
 */
/***************************************************************************************************/

/**
 * This method is used to tell the DCGM module to update all the fields being watched.
 *
 * Note: If the if the operation mode was set to manual mode (DCGM_OPERATION_MODE_MANUAL) during
 * initialization (\ref dcgmInit), this method must be caused periodically to allow field value watches
 * the opportunity to gather samples.
 *
 * @param pDcgmHandle           IN: DCGM Handle
 * @param waitForUpdate         IN: Whether or not to wait for the update loop to complete before returning to the
 *                                  caller 1=wait. 0=do not wait.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a waitForUpdate is invalid
 *        - \ref DCGM_ST_GENERIC_ERROR        if an unspecified DCGM error occurs
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmUpdateAllFields(dcgmHandle_t pDcgmHandle, int waitForUpdate);

/** @} */ // Closing for DCGMAPI_Admin_ExecCtrl


/***************************************************************************************************/
/** @defgroup DCGMAPI_PROCESS_STATS Process Statistics
 *  Describes APIs to investigate statistics such as accounting, performance and errors during the
 *  lifetime of a GPU process
 *  @{
 */
/***************************************************************************************************/

/**
 * Request that DCGM start recording stats for fields that can be queried with dcgmGetPidInfo().
 *
 * Note that the first update of the field will not occur until the next field update cycle.
 * To force a field update cycle, call dcgmUpdateAllFields(1).
 *
 * @param pDcgmHandle         IN: DCGM Handle
 * @param groupId             IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                                details on creating the group. Alternatively, pass in the group id as
 *                                \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param updateFreq          IN: How often to update this field in usec
 * @param maxKeepAge          IN: How long to keep data for this field in seconds
 * @param maxKeepSamples      IN: Maximum number of samples to keep. 0=no limit
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *        - \ref DCGM_ST_REQUIRES_ROOT        if the host engine is being run as non-root, and accounting mode could not
 *                                            be enabled (requires root). Run "nvidia-smi -am 1" as root on the node
 *                                            before starting DCGM to fix this.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmWatchPidFields(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupId,
                                                long long updateFreq,
                                                double maxKeepAge,
                                                int maxKeepSamples);

/**
 *
 * Get information about all GPUs while the provided pid was running
 *
 * In order for this request to work, you must first call dcgmWatchPidFields() to
 * make sure that DCGM is watching the appropriate field IDs that will be
 * populated in pidInfo
 *
 * @param pDcgmHandle IN: DCGM Handle
 * @param groupId     IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate
 *                        for details on creating the group. Alternatively, pass in the group id as
 *                        \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param pidInfo IN/OUT: Structure to return information about pid in. pidInfo->pid must be set to the pid in question.
 *                        pidInfo->version should be set to dcgmPidInfo_version.
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_NO_DATA             if the PID did not run on any GPU
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetPidInfo(dcgmHandle_t pDcgmHandle, dcgmGpuGrp_t groupId, dcgmPidInfo_t *pidInfo);

/** @} */ // Closing for DCGMAPI_PROCESS_STATS

/***************************************************************************************************/
/** @defgroup DCGMAPI_JOB_STATS Job Statistics
 * The client can invoke DCGM APIs to start and stop collecting the stats at the process boundaries
 * (during prologue and epilogue). This will enable DCGM to monitor all the PIDs while the job is
 * in progress, and provide a summary of active processes and resource usage during the window of
 * interest.
 *  @{
 */
/***************************************************************************************************/

/**
 * Request that DCGM start recording stats for fields that are queried with dcgmJobGetStats()
 *
 * Note that the first update of the field will not occur until the next field update cycle.
 * To force a field update cycle, call dcgmUpdateAllFields(1).
 *
 * @param pDcgmHandle         IN: DCGM Handle
 * @param groupId             IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                                details on creating the group. Alternatively, pass in the group id as
 *                                \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param updateFreq          IN: How often to update this field in usec
 * @param maxKeepAge          IN: How long to keep data for this field in seconds
 * @param maxKeepSamples      IN: Maximum number of samples to keep. 0=no limit
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *        - \ref DCGM_ST_REQUIRES_ROOT        if the host engine is being run as non-root, and
 *                                            accounting mode could not be enabled (requires root).
 *                                            Run "nvidia-smi -am 1" as root on the node before starting
 *                                            DCGM to fix this.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmWatchJobFields(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupId,
                                                long long updateFreq,
                                                double maxKeepAge,
                                                int maxKeepSamples);

/**
 * This API is used by the client to notify DCGM about the job to be started. Should be invoked as
 * part of job prologue
 *
 * @param pDcgmHandle       IN: DCGM Handle
 * @param groupId           IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                              details on creating the group. Alternatively, pass in the group id as
 *                              \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param jobId             IN: User provided string to represent the job
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *       - \ref DCGM_ST_DUPLICATE_KEY       if the specified \a jobId is already in use
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmJobStartStats(dcgmHandle_t pDcgmHandle, dcgmGpuGrp_t groupId, char jobId[64]);

/**
 * This API is used by the clients to notify DCGM to stop collecting stats for the job represented
 * by job id. Should be invoked as part of job epilogue.
 * The job Id remains available to view the stats at any point but cannot be used to start a new job.
 * You must call dcgmWatchJobFields() before this call to enable watching of job
 *
 * @param pDcgmHandle       IN: DCGM Handle
 * @param jobId             IN: User provided string to represent the job
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *       - \ref DCGM_ST_NO_DATA             if \a jobId is not a valid job identifier.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmJobStopStats(dcgmHandle_t pDcgmHandle, char jobId[64]);

/**
 * Get stats for the job identified by DCGM generated job id. The stats can be retrieved at any
 * point when the job is in process.
 * If you want to reuse this jobId, call \ref dcgmJobRemove after this call.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param jobId              IN: User provided string to represent the job
 * @param pJobInfo       IN/OUT: Structure to return information about the job.<br> .version should be set to
 *                               \ref dcgmJobInfo_version before this call.
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *       - \ref DCGM_ST_NO_DATA             if \a jobId is not a valid job identifier.
 *       - \ref DCGM_ST_VER_MISMATCH        if .version is not set or is invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmJobGetStats(dcgmHandle_t pDcgmHandle, char jobId[64], dcgmJobInfo_t *pJobInfo);

/**
 * This API tells DCGM to stop tracking the job given by jobId. After this call, you will no longer
 * be able to call dcgmJobGetStats() on this jobId. However, you will be able to reuse jobId after
 * this call.
 *
 * @param pDcgmHandle       IN: DCGM Handle
 * @param jobId             IN: User provided string to represent the job
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *       - \ref DCGM_ST_NO_DATA             if \a jobId is not a valid job identifier.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmJobRemove(dcgmHandle_t pDcgmHandle, char jobId[64]);

/**
 * This API tells DCGM to stop tracking all jobs. After this call, you will no longer
 * be able to call dcgmJobGetStats() any jobs until you call dcgmJobStartStats again.
 * You will be able to reuse any previously-used jobIds after this call.
 *
 * @param pDcgmHandle       IN: DCGM Handle
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmJobRemoveAll(dcgmHandle_t pDcgmHandle);

/** @} */ // Closing for DCGMAPI_JOB_STATS

/***************************************************************************************************/
/** @defgroup DCGMAPI_HM Health Monitor
 *
 *  This chapter describes the methods that handle the GPU health monitor.
 *
 *  @{
 */
/***************************************************************************************************/

/**
 * Enable the DCGM health check system for the given systems defined in \ref dcgmHealthSystems_t
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more entities. Look at \ref dcgmGroupCreate
 *                               for details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs or
 *                               \a DCGM_GROUP_ALL_NVSWITCHES to perform operation on all the NvSwitches.
 * @param systems            IN: An enum representing systems that should be enabled for health checks logically OR'd
 *                               together. Refer to \ref dcgmHealthSystems_t for details.
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmHealthSet(dcgmHandle_t pDcgmHandle, dcgmGpuGrp_t groupId, dcgmHealthSystems_t systems);

/**
 * Enable the DCGM health check system for the given systems defined in \ref dcgmHealthSystems_t
 *
 * Since DCGM 2.0
 *
 * @param pDcgmHandle                   IN: DCGM Handle
 * @param healthSet                     IN: Parameters to use when setting health watches. See
 *                                          \ref dcgmHealthSetParams_v2 for the description of each parameter.
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 */

dcgmReturn_t DCGM_PUBLIC_API dcgmHealthSet_v2(dcgmHandle_t pDcgmHandle, dcgmHealthSetParams_v2 *params);

/**
 * Retrieve the current state of the DCGM health check system
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more entities. Look at \ref dcgmGroupCreate
 *                               for details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs or
 *                               \a DCGM_GROUP_ALL_NVSWITCHES to perform operation on all the NvSwitches.
 * @param systems           OUT: An integer representing the enabled systems for the given group Refer to
 *                               \ref dcgmHealthSystems_t for details.
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmHealthGet(dcgmHandle_t pDcgmHandle,
                                           dcgmGpuGrp_t groupId,
                                           dcgmHealthSystems_t *systems);


/**
 * Check the configured watches for any errors/failures/warnings that have occurred
 * since the last time this check was invoked.  On the first call, stateful information
 * about all of the enabled watches within a group is created but no error results are
 * provided.  On subsequent calls, any error information will be returned.
 *
 *
 * @param pDcgmHandle                   IN: DCGM Handle
 * @param groupId                       IN: Group ID representing a collection of one or more entities.
 *                                          Refer to \ref dcgmGroupCreate for details on creating a group
 * @param results                      OUT: A reference to the dcgmHealthResponse_t structure to populate.
 *                                          results->version must be set to dcgmHealthResponse_version.
 *
 * @return
 *       - \ref DCGM_ST_OK                  if the call was successful
 *       - \ref DCGM_ST_BADPARAM            if a parameter is invalid
 *       - \ref DCGM_ST_VER_MISMATCH        if results->version is not dcgmHealthResponse_version
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmHealthCheck(dcgmHandle_t pDcgmHandle,
                                             dcgmGpuGrp_t groupId,
                                             dcgmHealthResponse_t *results);

/** @} */

/***************************************************************************************************/
/** @defgroup DCGMAPI_PO Policies
 *
 *  This chapter describes the methods that handle system policy management and violation settings.
 *  The APIs in Policies module can be broken down into following categories:
 *
 *  @{
 */
/***************************************************************************************************/

/***************************************************************************************************/
/** @defgroup DCGMAPI_PO_Setup Setup and Management
 *  Describes APIs for setting up policies and registering callbacks to receive notification in
 *  case specific policy condition has been violated.
 *  @{
 */
/***************************************************************************************************/

/**
 * Set the current violation policy inside the policy manager.  Given the conditions within the
 * \ref dcgmPolicy_t structure, if a violation has occurred, subsequent action(s) may be performed to
 * either report or contain the failure.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                               details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param policy             IN: A reference to \ref dcgmPolicy_t that will be applied to all GPUs in the group.
 * @param statusHandle   IN/OUT: Resulting status for the operation.  Pass it as NULL if the detailed error information
 *                               is not needed. Refer to \ref dcgmStatusCreate for details on creating a status handle.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a groupId or \a policy is invalid
 *        - \ref DCGM_ST_NOT_SUPPORTED        if any unsupported GPUs are part of the GPU group specified in groupId
 *        - DCGM_ST_*                         a different error has occurred and is stored in \a statusHandle.
 *                                            Refer to \ref dcgmReturn_t
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmPolicySet(dcgmHandle_t pDcgmHandle,
                                           dcgmGpuGrp_t groupId,
                                           dcgmPolicy_t *policy,
                                           dcgmStatus_t statusHandle);

/**
 * Get the current violation policy inside the policy manager. Given a groupId, a number of
 * policy structures are retrieved.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                               details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param count              IN: The size of the policy array.  This is the maximum number of policies that will be
 *                               retrieved and ultimately should correspond to the number of GPUs specified in the
 *                               group.
 * @param policy             OUT: A reference to \ref dcgmPolicy_t that will used as storage for the current policies
 *                                applied to each GPU in the group.
 * @param statusHandle    IN/OUT: Resulting status for the operation. Pass it as NULL if the detailed error information
 *                                for the operation is not needed. Refer to \ref dcgmStatusCreate for details on
 *                                creating a status handle.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a groupId or \a policy is invalid
 *        - DCGM_ST_*                         a different error has occurred and is stored in \a statusHandle.
 *                                            Refer to \ref dcgmReturn_t
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmPolicyGet(dcgmHandle_t pDcgmHandle,
                                           dcgmGpuGrp_t groupId,
                                           int count,
                                           dcgmPolicy_t *policy,
                                           dcgmStatus_t statusHandle);

/**
 * Register a function to be called when a specific policy condition (see \ref dcgmPolicyCondition_t) has been
 * violated.  This callback(s) will be called automatically when in DCGM_OPERATION_MODE_AUTO mode and only after
 * dcgmPolicyTrigger when in DCGM_OPERATION_MODE_MANUAL mode.  All callbacks are made within a separate thread.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                               details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param condition          IN: The set of conditions specified as an OR'd list (see \ref dcgmPolicyCondition_t) for
 *                               which to register a callback function
 * @param beginCallback      IN: A reference to a function that should be called should a violation occur.
 *                               This function will be called prior to any actions specified by the policy are taken.
 * @param finishCallback     IN: A reference to a function that should be called should a violation occur.
 *                           This function will be called after any action specified by the policy are completed.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a groupId, \a condition, is invalid, \a beginCallback, or
 *                                            \a finishCallback is NULL
 *        - \ref DCGM_ST_NOT_SUPPORTED        if any unsupported GPUs are part of the GPU group specified in groupId
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmPolicyRegister(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupId,
                                                dcgmPolicyCondition_t condition,
                                                fpRecvUpdates beginCallback,
                                                fpRecvUpdates finishCallback);

/**
 * Unregister a function to be called for a specific policy condition (see \ref dcgmPolicyCondition_t).
 * This function will unregister all callbacks for a given condition and handle.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                               details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param condition          IN: The set of conditions specified as an OR'd list (see \ref dcgmPolicyCondition_t) for
 *                               which to unregister a callback function
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if \a groupId, \a condition, is invalid or \a callback is NULL
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmPolicyUnregister(dcgmHandle_t pDcgmHandle,
                                                  dcgmGpuGrp_t groupId,
                                                  dcgmPolicyCondition_t condition);

/** @} */ // Closing for DCGMAPI_PO_Setup

/***************************************************************************************************/
/** @defgroup DCGMAPI_PO_MI Manual Invocation
 *  Describes APIs which can be used to perform direct actions (e.g. Perform GPU Reset, Run Health
 *  Diagnostics) on a group of GPUs.
 *  @{
 */
/***************************************************************************************************/

/**
 * Inform the action manager to perform a manual validation of a group of GPUs on the system
 *
 * *************************************** DEPRECATED ***************************************
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate for
 *                               details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param validate           IN: The validation to perform after the action.
 * @param response          OUT: Result of the validation process. Refer to \ref dcgmDiagResponse_t for details.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_NOT_SUPPORTED        if running the specified \a validate is not supported. This is usually due
 *                                            to the Tesla recommended driver not being installed on the system.
 *        - \ref DCGM_ST_BADPARAM             if \a groupId, \a validate, or \a statusHandle is invalid
 *        - \ref DCGM_ST_GENERIC_ERROR        an internal error has occurred
 *        - \ref DCGM_ST_GROUP_INCOMPATIBLE   if \a groupId refers to a group of non-homogeneous GPUs. This is currently
 *                                            not allowed.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmActionValidate(dcgmHandle_t pDcgmHandle,
                                                dcgmGpuGrp_t groupId,
                                                dcgmPolicyValidation_t validate,
                                                dcgmDiagResponse_t *response);

/**
 * Inform the action manager to perform a manual validation of a group of GPUs on the system
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param drd                IN: Contains the group id, test names, test parameters, struct version, and the validation
 *                               that should be performed. Look at \ref dcgmGroupCreate for details on creating the
 *                               group. Alternatively, pass in the group id as \a DCGM_GROUP_ALL_GPUS to perform
 *                               operation on all the GPUs.
 * @param response          OUT: Result of the validation process. Refer to \ref dcgmDiagResponse_t for details.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_NOT_SUPPORTED        if running the specified \a validate is not supported. This is usually
 *                                            due to the Tesla recommended driver not being installed on the system.
 *        - \ref DCGM_ST_BADPARAM             if \a groupId, \a validate, or \a statusHandle is invalid
 *        - \ref DCGM_ST_GENERIC_ERROR        an internal error has occurred
 *        - \ref DCGM_ST_GROUP_INCOMPATIBLE   if \a groupId refers to a group of non-homogeneous GPUs. This is
 *                                            currently not allowed.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmActionValidate_v2(dcgmHandle_t pDcgmHandle,
                                                   dcgmRunDiag_v7 *drd,
                                                   dcgmDiagResponse_t *response);

/**
 * Run a diagnostic on a group of GPUs
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param groupId            IN: Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate
 *                               for details on creating the group. Alternatively, pass in the group id as
 *                               \a DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs.
 * @param diagLevel          IN: Diagnostic level to run
 * @param diagResponse   IN/OUT: Result of running the DCGM diagnostic.<br>
 *                               .version should be set to \ref dcgmDiagResponse_version before this call.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_NOT_SUPPORTED        if running the diagnostic is not supported. This is usually due to the
 *                                            Tesla recommended driver not being installed on the system.
 *        - \ref DCGM_ST_BADPARAM             if a provided parameter is invalid or missing
 *        - \ref DCGM_ST_GENERIC_ERROR        an internal error has occurred
 *        - \ref DCGM_ST_GROUP_INCOMPATIBLE   if \a groupId refers to a group of non-homogeneous GPUs. This is
 *                                            currently not allowed.
 *        - \ref DCGM_ST_VER_MISMATCH         if .version is not set or is invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmRunDiagnostic(dcgmHandle_t pDcgmHandle,
                                               dcgmGpuGrp_t groupId,
                                               dcgmDiagnosticLevel_t diagLevel,
                                               dcgmDiagResponse_t *diagResponse);

/** @} */ // Closing for DCGMAPI_PO_MI

/** @} */ // Closing for DCGMAPI_PO

/***************************************************************************************************/
/** @addtogroup DCGMAPI_Admin_ExecCtrl
 *  @{
 */
/***************************************************************************************************/

/**
 * Inform the policy manager loop to perform an iteration and trigger the callbacks of any
 * registered functions. Callback functions will be called from a separate thread as the calling function.
 *
 * Note: The GPU monitoring and management agent must call this method periodically if the operation
 * mode is set to manual mode (DCGM_OPERATION_MODE_MANUAL) during initialization
 * (\ref dcgmInit).
 *
 * @param pDcgmHandle                   IN: DCGM Handle
 *
 * @return
 *        - \ref DCGM_ST_OK                   If the call was successful
 *        - DCGM_ST_GENERIC_ERROR             The policy manager was unable to perform another iteration.
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmPolicyTrigger(dcgmHandle_t pDcgmHandle);

/** @} */ // Closing for DCGMAPI_Admin_ExecCtrl

/***************************************************************************************************/
/** @defgroup DCGMAPI_Topo Topology
 *  @{
 */
/***************************************************************************************************/

/**
 * Gets device topology corresponding to the \a gpuId.
 *
 * @param pDcgmHandle             IN: DCGM Handle
 * @param gpuId                   IN: GPU Id corresponding to which topology information should be fetched
 * @param pDcgmDeviceTopology IN/OUT: Topology information corresponding to \a gpuId. pDcgmDeviceTopology->version must
 *                                    be set to dcgmDeviceTopology_version before this call.
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful.
 *        - \ref DCGM_ST_BADPARAM             if \a gpuId or \a pDcgmDeviceTopology were not valid.
 *        - \ref DCGM_ST_VER_MISMATCH         if pDcgmDeviceTopology->version was not set to dcgmDeviceTopology_version.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetDeviceTopology(dcgmHandle_t pDcgmHandle,
                                                   unsigned int gpuId,
                                                   dcgmDeviceTopology_t *pDcgmDeviceTopology);

/**
 * Gets group topology corresponding to the \a groupId.
 *
 * @param pDcgmHandle            IN: DCGM Handle
 * @param groupId                IN: GroupId corresponding to which topology information should be fetched
 * @param pDcgmGroupTopology IN/OUT: Topology information corresponding to \a groupId. pDcgmgroupTopology->version must
 *                                   be set to dcgmGroupTopology_version.
 * @return
 *        - \ref DCGM_ST_OK             if the call was successful.
 *        - \ref DCGM_ST_BADPARAM       if \a groupId or \a pDcgmGroupTopology were not valid.
 *        - \ref DCGM_ST_VER_MISMATCH   if pDcgmgroupTopology->version was not set to dcgmGroupTopology_version.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmGetGroupTopology(dcgmHandle_t pDcgmHandle,
                                                  dcgmGpuGrp_t groupId,
                                                  dcgmGroupTopology_t *pDcgmGroupTopology);

/** @} */ // Closing for DCGMAPI_Topo

/***************************************************************************************************/
/** @defgroup DCGMAPI_METADATA Metadata
 * @{
 *  This chapter describes the methods that query for DCGM metadata.
 */
/***************************************************************************************************/

/*************************************************************************/
/**
 * Retrieve the total amount of memory that the hostengine process is currently using.
 * This measurement represents both the resident set size (what is currently in RAM) and
 * the swapped memory that belongs to the process.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param memoryInfo     IN/OUT: see \ref dcgmIntrospectMemory_t. memoryInfo->version must be set to
 *                               dcgmIntrospectMemory_version prior to this call.
 * @param waitIfNoData       IN: if no metadata is gathered wait till this occurs (!0) or return DCGM_ST_NO_DATA (0)
 *
 * @return
 *       - \ref DCGM_ST_OK                   if the call was successful
 *       - \ref DCGM_ST_NOT_CONFIGURED       if metadata gathering state is \a DCGM_INTROSPECT_STATE_DISABLED
 *       - \ref DCGM_ST_NO_DATA              if \a waitIfNoData is false and metadata has not been gathered yet
 *       - \ref DCGM_ST_VER_MISMATCH         if memoryInfo->version is 0 or invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmIntrospectGetHostengineMemoryUsage(dcgmHandle_t pDcgmHandle,
                                                                    dcgmIntrospectMemory_t *memoryInfo,
                                                                    int waitIfNoData);

/*************************************************************************/
/**
 * Retrieve the CPU utilization of the DCGM hostengine process.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param cpuUtil        IN/OUT: see \ref dcgmIntrospectCpuUtil_t. cpuUtil->version must be set to
 *                               dcgmIntrospectCpuUtil_version prior to this call.
 * @param waitIfNoData       IN: if no metadata is gathered wait till this occurs (!0) or return DCGM_ST_NO_DATA (0)
 *
 * @return
 *       - \ref DCGM_ST_OK                   if the call was successful
 *       - \ref DCGM_ST_NOT_CONFIGURED       if metadata gathering state is \a DCGM_INTROSPECT_STATE_DISABLED
 *       - \ref DCGM_ST_NO_DATA              if \a waitIfNoData is false and metadata has not been gathered yet
 *       - \ref DCGM_ST_VER_MISMATCH         if cpuUtil->version or execTime->version is 0 or invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmIntrospectGetHostengineCpuUtilization(dcgmHandle_t pDcgmHandle,
                                                                       dcgmIntrospectCpuUtil_t *cpuUtil,
                                                                       int waitIfNoData);

/** @} */ // Closing for DCGMAPI_METADATA

/***************************************************************************************************/
/** @defgroup DCGMAPI_TOPOLOGY Topology
 * @{
 *  This chapter describes the methods that query for DCGM topology information.
 */
/***************************************************************************************************/

/*************************************************************************/
/**
 * Get the best group of gpus from the specified bitmask according to topological proximity: cpuAffinity, NUMA
 * node, and NVLink.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param inputGpuIds        IN: a bitmask of which GPUs DCGM should consider. If some of the GPUs on the system are
 *                               already in use, they shouldn't be included in the bitmask. 0 means that all of the GPUs
 *                               in the system should be considered.
 * @param numGpus            IN: the number of GPUs that are desired from inputGpuIds. If this number is greater than
 *                               the number of healthy GPUs in inputGpuIds, then less than numGpus gpus will be
 *                               specified in outputGpuIds.
 * @param outputGpuIds      OUT: a bitmask of numGpus or fewer GPUs from inputGpuIds that represent the best placement
 *                               available from inputGpuIds.
 * @param hintFlags          IN: a bitmask of DCGM_TOPO_HINT_F_ #defines of hints that should be taken into account when
 *                               assigning outputGpuIds.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmSelectGpusByTopology(dcgmHandle_t pDcgmHandle,
                                                      uint64_t inputGpuIds,
                                                      uint32_t numGpus,
                                                      uint64_t *outputGpuIds,
                                                      uint64_t hintFlags);

/** @} */ // Closing for DCGMAPI_TOPOLOGY

/***************************************************************************************************/
/** @defgroup DCGMAPI_MODULES Modules
 * @{
 *  This chapter describes the methods that query and configure DCGM modules.
 */
/***************************************************************************************************/

/*************************************************************************/
/**
 * Add a module to the denylist. This module will be prevented from being loaded
 * if it hasn't been loaded already. Modules are lazy-loaded as they are used by
 * DCGM APIs, so it's important to call this API soon after the host engine has been started.
 * You can also pass --denylist-modules to the nv-hostengine binary to make sure modules
 * get add to the denylist immediately after the host engine starts up.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param moduleId           IN: ID of the module to denylist. Use \ref dcgmModuleGetStatuses to get a list of valid
 *                               module IDs.
 *
 * @return
 *        - \ref DCGM_ST_OK         if the module has been add to the denylist.
 *        - \ref DCGM_ST_IN_USE     if the module has already been loaded and cannot add to the denylist.
 *        - \ref DCGM_ST_BADPARAM   if a parameter is missing or bad.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmModuleDenylist(dcgmHandle_t pDcgmHandle, dcgmModuleId_t moduleId);

/*************************************************************************/
/**
 * Get the status of all of the DCGM modules.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param moduleStatuses    OUT: Module statuses.<br>
 *                               .version should be set to dcgmModuleStatuses_version upon calling.
 *
 * @return
 *        - \ref DCGM_ST_OK         if the request succeeds.
 *        - \ref DCGM_ST_BADPARAM   if a parameter is missing or bad.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmModuleGetStatuses(dcgmHandle_t pDcgmHandle, dcgmModuleGetStatuses_t *moduleStatuses);

/** @} */ // Closing for DCGMAPI_MODULES

/*************************************************************************/
/** @defgroup DCGMAPI_PROFILING Profiling
 * @{
 *  This chapter describes the methods that watch profiling fields from within DCGM.
 */
/*************************************************************************/

/*************************************************************************/
/**
 * Get all of the profiling metric groups for a given GPU group.
 *
 * Profiling metrics are watched in groups of fields that are all watched together. For instance, if you want
 * to watch DCGM_FI_PROF_GR_ENGINE_ACTIVITY, this might also be in the same group as DCGM_FI_PROF_SM_EFFICIENCY.
 * Watching this group would result in DCGM storing values for both of these metrics.
 *
 * Some groups cannot be watched concurrently as others as they utilize the same hardware resource. For instance,
 * you may not be able to watch DCGM_FI_PROF_TENSOR_OP_UTIL at the same time as DCGM_FI_PROF_GR_ENGINE_ACTIVITY
 * on your hardware. At the same time, you may be able to watch DCGM_FI_PROF_TENSOR_OP_UTIL at the same time as
 * DCGM_FI_PROF_NVLINK_TX_DATA.
 *
 * Metrics that can be watched concurrently will have different .majorId fields in their dcgmProfMetricGroupInfo_t
 *
 * See \ref dcgmGroupCreate for details on creating a GPU group
 * See \ref dcgmProfWatchFields to actually watch a metric group
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param metricGroups   IN/OUT: Metric groups supported for metricGroups->groupId.<br>
 *                               metricGroups->version should be set to dcgmProfGetMetricGroups_version upon calling.
 *
 * @return
 *        - \ref DCGM_ST_OK                     if the request succeeds.
 *        - \ref DCGM_ST_BADPARAM               if a parameter is missing or bad.
 *        - \ref DCGM_ST_GROUP_INCOMPATIBLE     if metricGroups->groupId's GPUs are not identical GPUs.
 *        - \ref DCGM_ST_NOT_SUPPORTED          if profiling metrics are not supported for the given GPU group.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmProfGetSupportedMetricGroups(dcgmHandle_t pDcgmHandle,
                                                              dcgmProfGetMetricGroups_t *metricGroups);

/**
 * Request that DCGM start recording updates for a given list of profiling field IDs.
 *
 * Once metrics have been watched by this API, any of the normal DCGM field-value retrieval APIs can be used on
 * the underlying fieldIds of this metric group. See \ref dcgmGetLatestValues_v2, \ref dcgmGetLatestValuesForFields,
 * \ref dcgmEntityGetLatestValues, and \ref dcgmEntitiesGetLatestValues.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param watchFields        IN: Details of which metric groups to watch for which GPUs. See \ref dcgmProfWatchFields_v1
 *                               for details of what should be put in each struct member. watchFields->version should be
 *                               set to dcgmProfWatchFields_version upon calling.
 *
 * @return
 *        - \ref DCGM_ST_OK                     if the call was successful
 *        - \ref DCGM_ST_BADPARAM               if a parameter is invalid
 *        - \ref DCGM_ST_NOT_SUPPORTED          if profiling metric group metricGroupTag is not supported for the given
 *                                              GPU group.
 *        - \ref DCGM_ST_GROUP_INCOMPATIBLE     if groupId's GPUs are not identical GPUs. Profiling metrics are only
 *                                              support for homogenous groups of GPUs.
 *        - \ref DCGM_ST_PROFILING_MULTI_PASS   if any of the metric groups could not be watched concurrently due to
 *                                              requiring the hardware to gather them with multiple passes
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmProfWatchFields(dcgmHandle_t pDcgmHandle, dcgmProfWatchFields_t *watchFields);

/**
 * Request that DCGM stop recording updates for all profiling field IDs for all GPUs
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param unwatchFields      IN: Details of which metric groups to unwatch for which GPUs. See \ref
 *                               dcgmProfUnwatchFields_v1 for details of what should be put in each struct member.
 *                               unwatchFields->version should be set to dcgmProfUnwatchFields_version upon calling.
 *
 * @return
 *        - \ref DCGM_ST_OK                   if the call was successful
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmProfUnwatchFields(dcgmHandle_t pDcgmHandle, dcgmProfUnwatchFields_t *unwatchFields);

/**
 * Pause profiling activities in DCGM. This should be used when you are monitoring profiling fields
 * from DCGM but want to be able to still run developer tools like nvprof, nsight systems, and nsight compute.
 * Profiling fields start with DCGM_PROF_ and are in the field ID range 1001-1012.
 *
 * Call this API before you launch one of those tools and dcgmProfResume() after the tool has completed.
 *
 * DCGM will save BLANK values while profiling is paused.
 *
 * Calling this while profiling activities are already paused is fine and will be treated as a no-op.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 *
 * @return
 *        - \ref DCGM_ST_OK                   If the call was successful.
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmProfPause(dcgmHandle_t pDcgmHandle);

/**
 * Resume profiling activities in DCGM that were previously paused with dcgmProfPause().
 *
 * Call this API after you have completed running other NVIDIA developer tools to reenable DCGM
 * profiling metrics.
 *
 * DCGM will save BLANK values while profiling is paused.
 *
 * Calling this while profiling activities have already been resumed is fine and will be treated as a no-op.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 *
 * @return
 *        - \ref DCGM_ST_OK                   If the call was successful.
 *        - \ref DCGM_ST_BADPARAM             if a parameter is invalid.
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmProfResume(dcgmHandle_t pDcgmHandle);

/** @} */ // Closing for DCGMAPI_PROFILING

/**
 * Adds fake GPU instances and or compute instances for testing purposes. The entity IDs specified for
 * the GPU instances and compute instances are only guaranteed to be used by DCGM if MIG mode is not active.
 *
 * NOTE: this API will not work on a real system reading actual values from NVML, and it may even cause
 * the real instances to malfunction. This API is for testing purposes only.
 *
 * @param pDcgmHandle        IN: DCGM Handle
 * @param hierarchy
 *
 * @return
 *        - \ref DCGM_ST_OK
 *
 */
dcgmReturn_t DCGM_PUBLIC_API dcgmAddFakeInstances(dcgmHandle_t pDcgmHandle, dcgmMigHierarchy_v2 *hierarchy);

#ifdef __cplusplus
}
#endif

#endif /* DCGM_AGENT_H */
