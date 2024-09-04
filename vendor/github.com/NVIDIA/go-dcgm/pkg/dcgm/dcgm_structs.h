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
 * File: dcgm_structs.h
 */

#ifndef DCGM_STRUCTS_H
#define DCGM_STRUCTS_H

#include "dcgm_fields.h"
#include <limits.h>
#include <stdint.h>


/***************************************************************************************************/
/** @defgroup dcgmReturnEnums Enums and Macros
 *  @{
 */
/***************************************************************************************************/

/**
 * Creates a unique version number for each struct
 */
#define MAKE_DCGM_VERSION(typeName, ver) (unsigned int)(sizeof(typeName) | ((unsigned long)(ver) << 24U))

/**
 * Represents value of the field which can be returned by Host Engine in case the
 * operation is not successful
 */
#ifndef DCGM_BLANK_VALUES
#define DCGM_BLANK_VALUES

/**
 * Base value for 32 bits integer blank. can be used as an unspecified blank
 */
#define DCGM_INT32_BLANK 0x7ffffff0

/**
 * Base value for 64 bits integer blank. can be used as an unspecified blank
 */
#define DCGM_INT64_BLANK 0x7ffffffffffffff0

/**
 * Base value for double blank. 2 ** 47. FP 64 has 52 bits of mantissa,
 * so 47 bits can still increment by 1 and represent each value from 0-15
 */
#define DCGM_FP64_BLANK 140737488355328.0

/**
 * Base value for string blank.
 */
#define DCGM_STR_BLANK "<<<NULL>>>"

/**
 * Represents an error where INT32 data was not found
 */
#define DCGM_INT32_NOT_FOUND (DCGM_INT32_BLANK + 1)

/**
 * Represents an error where INT64 data was not found
 */
#define DCGM_INT64_NOT_FOUND (DCGM_INT64_BLANK + 1)

/**
 * Represents an error where FP64 data was not found
 */
#define DCGM_FP64_NOT_FOUND (DCGM_FP64_BLANK + 1.0)

/**
 * Represents an error where STR data was not found
 */
#define DCGM_STR_NOT_FOUND "<<<NOT_FOUND>>>"

/**
 * Represents an error where fetching the INT32 value is not supported
 */
#define DCGM_INT32_NOT_SUPPORTED (DCGM_INT32_BLANK + 2)

/**
 * Represents an error where fetching the INT64 value is not supported
 */
#define DCGM_INT64_NOT_SUPPORTED (DCGM_INT64_BLANK + 2)

/**
 * Represents an error where fetching the FP64 value is not supported
 */
#define DCGM_FP64_NOT_SUPPORTED (DCGM_FP64_BLANK + 2.0)

/**
 * Represents an error where fetching the STR value is not supported
 */
#define DCGM_STR_NOT_SUPPORTED "<<<NOT_SUPPORTED>>>"

/**
 *  Represents and error where fetching the INT32 value is not allowed with our current credentials
 */
#define DCGM_INT32_NOT_PERMISSIONED (DCGM_INT32_BLANK + 3)

/**
 *  Represents and error where fetching the INT64 value is not allowed with our current credentials
 */
#define DCGM_INT64_NOT_PERMISSIONED (DCGM_INT64_BLANK + 3)

/**
 *  Represents and error where fetching the FP64 value is not allowed with our current credentials
 */
#define DCGM_FP64_NOT_PERMISSIONED (DCGM_FP64_BLANK + 3.0)

/**
 *  Represents and error where fetching the STR value is not allowed with our current credentials
 */
#define DCGM_STR_NOT_PERMISSIONED "<<<NOT_PERM>>>"

/**
 * Macro to check if a INT32 value is blank or not
 */
#define DCGM_INT32_IS_BLANK(val) (((val) >= DCGM_INT32_BLANK) ? 1 : 0)

/**
 * Macro to check if a INT64 value is blank or not
 */
#define DCGM_INT64_IS_BLANK(val) (((val) >= DCGM_INT64_BLANK) ? 1 : 0)

/**
 * Macro to check if a FP64 value is blank or not
 */
#define DCGM_FP64_IS_BLANK(val) (((val) >= DCGM_FP64_BLANK ? 1 : 0))

/**
 * Macro to check if a STR value is blank or not
 * Works on (char *). Looks for <<< at first position and >>> inside string
 */
#define DCGM_STR_IS_BLANK(val) (val == strstr(val, "<<<") && strstr(val, ">>>"))

#endif // DCGM_BLANK_VALUES

/**
 * Max number of GPUs supported by DCGM
 */
#define DCGM_MAX_NUM_DEVICES 32 /* DCGM 2.0 and newer = 32. DCGM 1.8 and older = 16. */

/**
 * Number of NvLink links per GPU supported by DCGM
 * 18 for Hopper, 12 for Ampere, 6 for Volta, and 4 for Pascal
 */
#define DCGM_NVLINK_MAX_LINKS_PER_GPU 18

/**
 * Number of nvlink errors supported by DCGM
 * @see NVML_NVLINK_ERROR_COUNT
 *
 * NVML_NVLINK_ERROR_DL_ECC_DATA not currently supported
 */
#define DCGM_NVLINK_ERROR_COUNT 4

/**
 * Number of nvlink error types: @see NVML_NVLINK_ERROR_COUNT
 * TODO: update with refactor of ampere-next nvlink APIs (JIRA DCGM-2628)
 */
#define DCGM_HEALTH_WATCH_NVLINK_ERROR_NUM_FIELDS 4

/**
 * Maximum NvLink links pre-Ampere
 */
#define DCGM_NVLINK_MAX_LINKS_PER_GPU_LEGACY1 6

/**
 * Maximum NvLink links pre-Hopper
 */
#define DCGM_NVLINK_MAX_LINKS_PER_GPU_LEGACY2 12

/**
 * Max number of NvSwitches supported by DCGM
 **/
#define DCGM_MAX_NUM_SWITCHES 12

/**
 * Number of NvLink links per NvSwitch supported by DCGM
 */
#define DCGM_NVLINK_MAX_LINKS_PER_NVSWITCH 64

/**
 * Number of Lines per NvSwitch NvLink supported by DCGM
 */
#define DCGM_LANE_MAX_LANES_PER_NVSWICH_LINK 4

/**
 * Maximum number of vGPU instances per physical GPU
 */
#define DCGM_MAX_VGPU_INSTANCES_PER_PGPU 32

/**
 * Max number of CPU nodes
 **/
#define DCGM_MAX_NUM_CPUS 8

/**
 * Max number of CPUs
 **/
#define DCGM_MAX_NUM_CPU_CORES 1024

/**
 * Max length of the DCGM string field
 */
#define DCGM_MAX_STR_LENGTH 256

/**
 * Default maximum age of samples kept (usec)
 */
#define DCGM_MAX_AGE_USEC_DEFAULT 30000000

/**
 * Max number of clocks supported for a device
 */
#define DCGM_MAX_CLOCKS 256

/**
 * Max limit on the number of groups supported by DCGM
 */
#define DCGM_MAX_NUM_GROUPS 64

/**
 * Max number of active FBC sessions
 */
#define DCGM_MAX_FBC_SESSIONS 256

/**
 * Represents the size of a buffer that holds a vGPU type Name or vGPU class type or name of process running on vGPU
 * instance.
 */
#define DCGM_VGPU_NAME_BUFFER_SIZE 64

/**
 * Represents the size of a buffer that holds a vGPU license string
 */
#define DCGM_GRID_LICENSE_BUFFER_SIZE 128

/**
 * Default compute mode -- multiple contexts per device
 */
#define DCGM_CONFIG_COMPUTEMODE_DEFAULT 0

/**
 * Compute-prohibited mode -- no contexts per device
 */
#define DCGM_CONFIG_COMPUTEMODE_PROHIBITED 1

/**
 * Compute-exclusive-process mode -- only one context per device, usable from multiple threads at a time
 */
#define DCGM_CONFIG_COMPUTEMODE_EXCLUSIVE_PROCESS 2

/**
 * Default Port Number for DCGM Host Engine
 */
#define DCGM_HE_PORT_NUMBER 5555

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Operation mode for DCGM
 *
 * DCGM can run in auto-mode where it runs additional threads in the background to collect
 * any metrics of interest and auto manages any operations needed for policy management.
 *
 * DCGM can also operate in manual-mode where it's execution is controlled by the user. In
 * this mode, the user has to periodically call APIs such as \ref dcgmPolicyTrigger and
 * \ref dcgmUpdateAllFields which tells DCGM to wake up and perform data collection and
 * operations needed for policy management.
 */
typedef enum dcgmOperationMode_enum
{
    DCGM_OPERATION_MODE_AUTO   = 1,
    DCGM_OPERATION_MODE_MANUAL = 2
} dcgmOperationMode_t;

/**
 * When more than one value is returned from a query, which order should it be returned in?
 */
typedef enum dcgmOrder_enum
{
    DCGM_ORDER_ASCENDING  = 1, //!< Data with earliest (lowest) timestamps returned first
    DCGM_ORDER_DESCENDING = 2  //!< Data with latest (highest) timestamps returned first
} dcgmOrder_t;

/**
 * Return values for DCGM API calls.
 */
typedef enum dcgmReturn_enum
{
    DCGM_ST_OK                   = 0,   //!< Success
    DCGM_ST_BADPARAM             = -1,  //!< A bad parameter was passed to a function
    DCGM_ST_GENERIC_ERROR        = -3,  //!< A generic, unspecified error
    DCGM_ST_MEMORY               = -4,  //!< An out of memory error occurred
    DCGM_ST_NOT_CONFIGURED       = -5,  //!< Setting not configured
    DCGM_ST_NOT_SUPPORTED        = -6,  //!< Feature not supported
    DCGM_ST_INIT_ERROR           = -7,  //!< DCGM Init error
    DCGM_ST_NVML_ERROR           = -8,  //!< When NVML returns error
    DCGM_ST_PENDING              = -9,  //!< Object is in pending state of something else
    DCGM_ST_UNINITIALIZED        = -10, //!< Object is in undefined state
    DCGM_ST_TIMEOUT              = -11, //!< Requested operation timed out
    DCGM_ST_VER_MISMATCH         = -12, //!< Version mismatch between received and understood API
    DCGM_ST_UNKNOWN_FIELD        = -13, //!< Unknown field id
    DCGM_ST_NO_DATA              = -14, //!< No data is available
    DCGM_ST_STALE_DATA           = -15, //!< Data is considered stale
    DCGM_ST_NOT_WATCHED          = -16, //!< The given field id is not being updated by the cache manager
    DCGM_ST_NO_PERMISSION        = -17, //!< Do not have permission to perform the desired action
    DCGM_ST_GPU_IS_LOST          = -18, //!< GPU is no longer reachable
    DCGM_ST_RESET_REQUIRED       = -19, //!< GPU requires a reset
    DCGM_ST_FUNCTION_NOT_FOUND   = -20, //!< The function that was requested was not found (bindings only error)
    DCGM_ST_CONNECTION_NOT_VALID = -21, //!< The connection to the host engine is not valid any longer
    DCGM_ST_GPU_NOT_SUPPORTED    = -22, //!< This GPU is not supported by DCGM
    DCGM_ST_GROUP_INCOMPATIBLE   = -23, //!< The GPUs of the provided group are not compatible with each other for the
                                        //!< requested operation
    DCGM_ST_MAX_LIMIT                   = -24, //!< Max limit reached for the object
    DCGM_ST_LIBRARY_NOT_FOUND           = -25, //!< DCGM library could not be found
    DCGM_ST_DUPLICATE_KEY               = -26, //!< Duplicate key passed to a function
    DCGM_ST_GPU_IN_SYNC_BOOST_GROUP     = -27, //!< GPU is already a part of a sync boost group
    DCGM_ST_GPU_NOT_IN_SYNC_BOOST_GROUP = -28, //!< GPU is not a part of a sync boost group
    DCGM_ST_REQUIRES_ROOT     = -29, //!< This operation cannot be performed when the host engine is running as non-root
    DCGM_ST_NVVS_ERROR        = -30, //!< DCGM GPU Diagnostic was successfully executed, but reported an error.
    DCGM_ST_INSUFFICIENT_SIZE = -31, //!< An input argument is not large enough
    DCGM_ST_FIELD_UNSUPPORTED_BY_API = -32, //!< The given field ID is not supported by the API being called
    DCGM_ST_MODULE_NOT_LOADED = -33, //!< This request is serviced by a module of DCGM that is not currently loaded
    DCGM_ST_IN_USE            = -34, //!< The requested operation could not be completed because the affected
                                     //!< resource is in use
    DCGM_ST_GROUP_IS_EMPTY = -35,    //!< This group is empty and the requested operation is not valid on an empty group
    DCGM_ST_PROFILING_NOT_SUPPORTED = -36,     //!< Profiling is not supported for this group of GPUs or GPU.
    DCGM_ST_PROFILING_LIBRARY_ERROR = -37,     //!< The third-party Profiling module returned an unrecoverable error.
    DCGM_ST_PROFILING_MULTI_PASS    = -38,     //!< The requested profiling metrics cannot be collected in a single pass
    DCGM_ST_DIAG_ALREADY_RUNNING    = -39,     //!< A diag instance is already running, cannot run a new diag until
                                               //!< the current one finishes.
    DCGM_ST_DIAG_BAD_JSON               = -40, //!< The DCGM GPU Diagnostic returned JSON that cannot be parsed
    DCGM_ST_DIAG_BAD_LAUNCH             = -41, //!< Error while launching the DCGM GPU Diagnostic
    DCGM_ST_DIAG_UNUSED                 = -42, //!< Unused
    DCGM_ST_DIAG_THRESHOLD_EXCEEDED     = -43, //!< A field value met or exceeded the error threshold.
    DCGM_ST_INSUFFICIENT_DRIVER_VERSION = -44, //!< The installed driver version is insufficient for this API
    DCGM_ST_INSTANCE_NOT_FOUND          = -45, //!< The specified GPU instance does not exist
    DCGM_ST_COMPUTE_INSTANCE_NOT_FOUND  = -46, //!< The specified GPU compute instance does not exist
    DCGM_ST_CHILD_NOT_KILLED            = -47, //!< Couldn't kill a child process within the retries
    DCGM_ST_3RD_PARTY_LIBRARY_ERROR     = -48, //!< Detected an error in a 3rd-party library
    DCGM_ST_INSUFFICIENT_RESOURCES      = -49, //!< Not enough resources available
    DCGM_ST_PLUGIN_EXCEPTION            = -50, //!< Exception thrown from a diagnostic plugin
    DCGM_ST_NVVS_ISOLATE_ERROR    = -51, //!< The diagnostic returned an error that indicates the need for isolation
    DCGM_ST_NVVS_BINARY_NOT_FOUND = -52, //!< The NVVS binary was not found in the specified location
    DCGM_ST_NVVS_KILLED           = -53, //!< The NVVS process was killed by a signal
    DCGM_ST_PAUSED                = -54, //!< The hostengine and all modules are paused
    DCGM_ST_ALREADY_INITIALIZED   = -55, //!< The object is already initialized
} dcgmReturn_t;

const char *errorString(dcgmReturn_t result);

/**
 * Type of GPU groups
 */
typedef enum dcgmGroupType_enum
{
    DCGM_GROUP_DEFAULT                   = 0, //!< All the GPUs on the node are added to the group
    DCGM_GROUP_EMPTY                     = 1, //!< Creates an empty group
    DCGM_GROUP_DEFAULT_NVSWITCHES        = 2, //!< All NvSwitches of the node are added to the group
    DCGM_GROUP_DEFAULT_INSTANCES         = 3, //!< All GPU instances of the node are added to the group
    DCGM_GROUP_DEFAULT_COMPUTE_INSTANCES = 4, //!< All compute instances of the node are added to the group
    DCGM_GROUP_DEFAULT_EVERYTHING        = 5, //!< All entities are added to this default group
} dcgmGroupType_t;

/**
 * Identifies for special DCGM groups
 */
#define DCGM_GROUP_ALL_GPUS              0x7fffffff
#define DCGM_GROUP_ALL_NVSWITCHES        0x7ffffffe
#define DCGM_GROUP_ALL_INSTANCES         0x7ffffffd
#define DCGM_GROUP_ALL_COMPUTE_INSTANCES 0x7ffffffc
#define DCGM_GROUP_ALL_ENTITIES          0x7ffffffb

/**
 * Maximum number of entities per entity group
 */
#define DCGM_GROUP_MAX_ENTITIES 64

/**
 * Simplified chip architecture. Note that these are made to match nvmlChipArchitecture_t and thus
 * do not start at 0.
 */
typedef enum dcgmChipArchitecture_enum
{
    DCGM_CHIP_ARCH_OLDER   = 1, //!< All GPUs older than Kepler
    DCGM_CHIP_ARCH_KEPLER  = 2, //!< All Kepler-architecture parts
    DCGM_CHIP_ARCH_MAXWELL = 3, //!< All Maxwell-architecture parts
    DCGM_CHIP_ARCH_PASCAL  = 4, //!< All Pascal-architecture parts
    DCGM_CHIP_ARCH_VOLTA   = 5, //!< All Volta-architecture parts
    DCGM_CHIP_ARCH_TURING  = 6, //!< All Turing-architecture parts
    DCGM_CHIP_ARCH_AMPERE  = 7, //!< All Ampere-architecture parts
    DCGM_CHIP_ARCH_ADA     = 8, //!< All Ada-architecture parts
    DCGM_CHIP_ARCH_HOPPER  = 9, //!< All Hopper-architecture parts

    DCGM_CHIP_ARCH_COUNT, //!< Keep this second to last, exclude unknown

    DCGM_CHIP_ARCH_UNKNOWN = 0xffffffff //!< Anything else, presumably something newer
} dcgmChipArchitecture_t;

/**
 * Represents the type of configuration to be fetched from the GPUs
 */
typedef enum dcgmConfigType_enum
{
    DCGM_CONFIG_TARGET_STATE  = 0, //!< The target configuration values to be applied
    DCGM_CONFIG_CURRENT_STATE = 1, //!< The current configuration state
} dcgmConfigType_t;

/**
 * Represents the power cap for each member of the group.
 */
typedef enum dcgmConfigPowerLimitType_enum
{
    DCGM_CONFIG_POWER_CAP_INDIVIDUAL = 0, //!< Represents the power cap to be applied for each member of the group
    DCGM_CONFIG_POWER_BUDGET_GROUP   = 1, //!< Represents the power budget for the entire group
} dcgmConfigPowerLimitType_t;

/** @} */


/***************************************************************************************************/
/** @defgroup dcgmStructs Structure definitions
 *  @{
 */
/***************************************************************************************************/
typedef uintptr_t dcgmHandle_t;   //!< Identifier for DCGM Handle
typedef uintptr_t dcgmGpuGrp_t;   //!< Identifier for a group of GPUs. A group can have one or more GPUs
typedef uintptr_t dcgmFieldGrp_t; //!< Identifier for a group of fields.
typedef uintptr_t dcgmStatus_t;   //!< Identifier for list of status codes

/**
 * DCGM Logging Severities. These match up with plog severities defined in Severity.h
 * Each level includes all of the levels above it. For instance, level 4 includes 3,2, and 1 as well
 */
typedef enum
{
    DcgmLoggingSeverityUnspecified = -1, /*!< Don't care/inherit from the environment */
    DcgmLoggingSeverityNone        = 0,  /*!< No logging */
    DcgmLoggingSeverityFatal       = 1,  /*!< Fatal Errors */
    DcgmLoggingSeverityError       = 2,  /*!< Errors */
    DcgmLoggingSeverityWarning     = 3,  /*!< Warnings */
    DcgmLoggingSeverityInfo        = 4,  /*!< Informative */
    DcgmLoggingSeverityDebug       = 5,  /*!< Debug information (will generate large logs) */
    DcgmLoggingSeverityVerbose     = 6   /*!< Verbose debugging information */
} DcgmLoggingSeverity_t;

/**
 * Represents a link object. type should be one of DCGM_FE_GPU or
 * DCGM_FE_SWITCH; gpuId or switchID is the associated gpu or switch; and index
 * is the link index, 0-based, with TX (even) coming before RX (odd).
 */
#pragma pack(push, 1)
typedef struct dcgm_link_s
{
    union
    {
        struct
        {
            dcgm_field_entity_group_t type : 8; /*!< Entity Group */
            uint8_t index                  : 8; /*!< Link Index Tx before Rx */
            union
            {
                dcgm_field_eid_t gpuId    : 16; /*!< Physical GPU ID */
                dcgm_field_eid_t switchId : 16; /*!< Physical Switch ID */
            };
        } parsed;             /*!< Broken out Link identifier GPU/SW:[GPU|SW]:Index */
        dcgm_field_eid_t raw; /*!< Raw Link ID */
    };
} dcgm_link_t;
#pragma pack(pop)

/**
 * Connection options for dcgmConnect_v2 (v1)
 *
 * NOTE: This version is deprecated. use dcgmConnectV2Params_v2
 */
typedef struct
{
    unsigned int version;                /*!< Version number. Use dcgmConnectV2Params_version */
    unsigned int persistAfterDisconnect; /*!< Whether to persist DCGM state modified by this connection
                                              once the connection is terminated. Normally, all field
                                              watches created by a connection are removed once a
                                              connection goes away.
                                              1 = do not clean up after this connection.
                                              0 = clean up after this connection */
} dcgmConnectV2Params_v1;

/**
 * Version 1 for \ref dcgmConnectV2Params_v1
 */
#define dcgmConnectV2Params_version1 MAKE_DCGM_VERSION(dcgmConnectV2Params_v1, 1)

/**
 * Connection options for dcgmConnect_v2 (v2)
 */
typedef struct
{
    unsigned int version;                /*!< Version number. Use dcgmConnectV2Params_version */
    unsigned int persistAfterDisconnect; /*!< Whether to persist DCGM state modified by this connection once the
                                              connection is terminated. Normally, all field watches created by a
                                              connection are removed once a connection goes away. 1 = do not clean up
                                              after this connection. 0 = clean up after this connection */
    unsigned int timeoutMs;              /*!< When attempting to connect to the specified host engine, how long should
                                              we wait in milliseconds before giving up */
    unsigned int addressIsUnixSocket;    /*!< Whether or not the passed-in address is a unix socket filename (1) or a
                                              TCP/IP address (0) */
} dcgmConnectV2Params_v2;

/**
 * Typedef for \ref dcgmConnectV2Params_v2
 */
typedef dcgmConnectV2Params_v2 dcgmConnectV2Params_t;

/**
 * Version 2 for \ref dcgmConnectV2Params_v2
 */
#define dcgmConnectV2Params_version2 MAKE_DCGM_VERSION(dcgmConnectV2Params_v2, 2)

/**
 * Latest version for \ref dcgmConnectV2Params_t
 */
#define dcgmConnectV2Params_version dcgmConnectV2Params_version2

/**
 * Typedef for \ref dcgmHostengineHealth_v1
 */
typedef struct
{
    unsigned int version;       //!< The version of this request
    unsigned int overallHealth; //!< 0 to indicate healthy, or a code to indicate the error
                                //   For now, this will always be populated with 0 if the
                                //   hostengine can respond. In the future this will be
                                //   updated to have other options like NVML unresponsive,
                                //   no GPUs on system, etc.
} dcgmHostengineHealth_v1;

/**
 * Typedef for \ref dcgmHostengineHealth_t
 */
typedef dcgmHostengineHealth_v1 dcgmHostengineHealth_t;

#define dcgmHostengineHealth_version1 MAKE_DCGM_VERSION(dcgmHostengineHealth_v1, 1)

/**
 * Latest version for \ref dcgmHostengineHealth_t
 */
#define dcgmHostengineHealth_version dcgmHostengineHealth_version1

/**
 * Represents a entityGroupId + entityId pair to uniquely identify a given entityId inside a group of entities
 *
 * Added in DCGM 1.5.0
 */
typedef struct
{
    dcgm_field_entity_group_t entityGroupId; //!< Entity Group ID entity belongs to
    dcgm_field_eid_t entityId;               //!< Entity ID of the entity
} dcgmGroupEntityPair_t;

/**
 * Structure to store information for DCGM group
 *
 * Added in DCGM 1.5.0
 */
typedef struct
{
    unsigned int version;                                      //!< Version Number (use dcgmGroupInfo_version2)
    unsigned int count;                                        //!< count of entityIds returned in \a entityList
    char groupName[DCGM_MAX_STR_LENGTH];                       //!< Group Name
    dcgmGroupEntityPair_t entityList[DCGM_GROUP_MAX_ENTITIES]; //!< List of the entities that are in this group
} dcgmGroupInfo_v2;

/**
 * Typedef for \ref dcgmGroupInfo_v2
 */
typedef dcgmGroupInfo_v2 dcgmGroupInfo_t;

/**
 * Version 2 for \ref dcgmGroupInfo_v2
 */
#define dcgmGroupInfo_version2 MAKE_DCGM_VERSION(dcgmGroupInfo_v2, 2)

/**
 * Latest version for \ref dcgmGroupInfo_t
 */
#define dcgmGroupInfo_version dcgmGroupInfo_version2

/**
 * Enum for the different kinds of MIG profiles
 */
typedef enum
{
    DcgmMigProfileNone                      = 0,  /*!< No profile (for GPUs) */
    DcgmMigProfileGpuInstanceSlice1         = 1,  /*!< GPU instance slice 1 */
    DcgmMigProfileGpuInstanceSlice2         = 2,  /*!< GPU instance slice 2 */
    DcgmMigProfileGpuInstanceSlice3         = 3,  /*!< GPU instance slice 3 */
    DcgmMigProfileGpuInstanceSlice4         = 4,  /*!< GPU instance slice 4 */
    DcgmMigProfileGpuInstanceSlice7         = 5,  /*!< GPU instance slice 7 */
    DcgmMigProfileGpuInstanceSlice8         = 6,  /*!< GPU instance slice 8 */
    DcgmMigProfileGpuInstanceSlice6         = 7,  /*!< GPU instance slice 6 */
    DcgmMigProfileGpuInstanceSlice1Rev1     = 8,  /*!< GPU instance slice 1 revision 1 */
    DcgmMigProfileGpuInstanceSlice2Rev1     = 9,  /*!< GPU instance slice 2 revision 1 */
    DcgmMigProfileGpuInstanceSlice1Rev2     = 10, /*!< GPU instance slice 1 revision 2 */
    DcgmMigProfileComputeInstanceSlice1     = 30, /*!< compute instance slice 1 */
    DcgmMigProfileComputeInstanceSlice2     = 31, /*!< compute instance slice 2 */
    DcgmMigProfileComputeInstanceSlice3     = 32, /*!< compute instance slice 3 */
    DcgmMigProfileComputeInstanceSlice4     = 33, /*!< compute instance slice 4*/
    DcgmMigProfileComputeInstanceSlice7     = 34, /*!< compute instance slice 7 */
    DcgmMigProfileComputeInstanceSlice8     = 35, /*!< compute instance slice 8 */
    DcgmMigProfileComputeInstanceSlice6     = 36, /*!< compute instance slice 6 */
    DcgmMigProfileComputeInstanceSlice1Rev1 = 37, /*!< compute instance slice 1 revision 1 */
} dcgmMigProfile_t;

/**
 * Represents a pair of entity pairings to uniquely identify an entity and its place in the hierarchy.
 */
typedef struct
{
    dcgmGroupEntityPair_t entity;  //!< Entity id and type for the entity in question
    dcgmGroupEntityPair_t parent;  //!< Entity id and type for the parent of the entity in question
    dcgmMigProfile_t sliceProfile; //!< Entity MIG profile identifier
} dcgmMigHierarchyInfo_t;

/**
 * Provides additional information about location of MIG entities.
 */
typedef struct
{
    char gpuUuid[128];                  /*!< GPU UUID */
    unsigned int nvmlGpuIndex;          /*!< GPU index from NVML */
    unsigned int nvmlInstanceId;        /*!< GPU instance index within GPU. 0 to N. -1 for GPU entities */
    unsigned int nvmlComputeInstanceId; /*!< GPU Compute instance index within GPU instance. 0 to N. -1 for GPU
                                         *      Instance and GPU entities */
    unsigned int nvmlMigProfileId;      /*!< Unique profile ID for GPU or Compute instances. -1 GPU entities
                                         *      \see nvmlComputeInstanceProfileInfo_st
                                         *      \see nvmlGpuInstanceProfileInfo_st */
    unsigned int nvmlProfileSlices;     /*!< Number of slices in the MIG profile */
} dcgmMigEntityInfo_t;

typedef struct
{
    dcgmGroupEntityPair_t entity;
    dcgmGroupEntityPair_t parent;
    dcgmMigEntityInfo_t info;
} dcgmMigHierarchyInfo_v2;

#define DCGM_MAX_INSTANCES_PER_GPU 8
// There can never be more compute instances per GPU than instances per GPU because a compute instance is part
// of an instance
#define DCGM_MAX_COMPUTE_INSTANCES_PER_GPU DCGM_MAX_INSTANCES_PER_GPU
// Currently, there cannot be more than 14 instances + compute instances. There are always 7 compute instances
// and never more than 7 instances
#define DCGM_MAX_TOTAL_INSTANCES_PER_GPU 14
#define DCGM_MAX_HIERARCHY_INFO          DCGM_MAX_NUM_DEVICES *DCGM_MAX_TOTAL_INSTANCES_PER_GPU
#define DCGM_MAX_INSTANCES               DCGM_MAX_NUM_DEVICES *DCGM_MAX_INSTANCES_PER_GPU
// The maximum compute instances are always the same as the maximum instances because each compute instance is
// part of an instance.
#define DCGM_MAX_COMPUTE_INSTANCES DCGM_MAX_INSTANCES

typedef struct
{
    unsigned int version;
    unsigned int count;
    dcgmMigHierarchyInfo_v2 entityList[DCGM_MAX_HIERARCHY_INFO];
} dcgmMigHierarchy_v2;

#define dcgmMigHierarchy_version2 MAKE_DCGM_VERSION(dcgmMigHierarchy_v2, 2)

#define dcgmMigHierarchy_version dcgmMigHierarchy_version2

/**
 * Bitmask indicating which cores are owned by this CPUs
 */
#define DCGM_CPU_CORE_BITMASK_COUNT_V1 (DCGM_MAX_NUM_CPU_CORES / sizeof(uint64_t) / CHAR_BIT)
typedef struct
{
    unsigned int version;
    uint64_t bitmask[DCGM_CPU_CORE_BITMASK_COUNT_V1];
} dcgmCpuHierarchyOwnedCores_v1;

typedef dcgmCpuHierarchyOwnedCores_v1 dcgmCpuHierarchyOwnedCores_t;

#define dcgmCpuHierarchyOwnedCores_version1 MAKE_DCGM_VERSION(dcgmCpuHierarchyOwnedCores_v1, 1)

/**
 * Hierarchy of CPUs and their cores
 */
typedef struct
{
    unsigned int version;
    unsigned int numCpus;
    struct dcgmCpuHierarchyCpu_v1
    {
        unsigned int cpuId;
        dcgmCpuHierarchyOwnedCores_v1 ownedCores;
    } cpus[DCGM_MAX_NUM_CPUS];
} dcgmCpuHierarchy_v1;

typedef dcgmCpuHierarchy_v1 dcgmCpuHierarchy_t;

/**
 * Version 1 for dcgmCpuHierarchy_t
 */
#define dcgmCpuHierarchy_version1 MAKE_DCGM_VERSION(dcgmCpuHierarchy_v1, 1)

/**
 * Maximum number of field groups that can exist
 */
#define DCGM_MAX_NUM_FIELD_GROUPS 64

/**
 * Maximum number of field IDs that can be in a single field group
 */
#define DCGM_MAX_FIELD_IDS_PER_FIELD_GROUP 128

/**
 * Structure to represent information about a field group
 */
typedef struct
{
    unsigned int version;                                        //!< Version number (dcgmFieldGroupInfo_version)
    unsigned int numFieldIds;                                    //!< Number of entries in fieldIds[] that are valid
    dcgmFieldGrp_t fieldGroupId;                                 //!< ID of this field group
    char fieldGroupName[DCGM_MAX_STR_LENGTH];                    //!< Field Group Name
    unsigned short fieldIds[DCGM_MAX_FIELD_IDS_PER_FIELD_GROUP]; //!< Field ids that belong to this group
} dcgmFieldGroupInfo_v1;

typedef dcgmFieldGroupInfo_v1 dcgmFieldGroupInfo_t;

/**
 * Version 1 for dcgmFieldGroupInfo_v1
 */
#define dcgmFieldGroupInfo_version1 MAKE_DCGM_VERSION(dcgmFieldGroupInfo_v1, 1)

/**
 * Latest version for dcgmFieldGroupInfo_t
 */
#define dcgmFieldGroupInfo_version dcgmFieldGroupInfo_version1

typedef struct
{
    unsigned int version;        //!< Version number (dcgmAllFieldGroupInfo_version)
    unsigned int numFieldGroups; //!< Number of entries in fieldGroups[] that are populated
    dcgmFieldGroupInfo_t fieldGroups[DCGM_MAX_NUM_FIELD_GROUPS]; //!< Info about each field group
} dcgmAllFieldGroup_v1;

typedef dcgmAllFieldGroup_v1 dcgmAllFieldGroup_t;

/**
 * Version 1 for dcgmAllFieldGroup_v1
 */
#define dcgmAllFieldGroup_version1 MAKE_DCGM_VERSION(dcgmAllFieldGroup_v1, 1)

/**
 * Latest version for dcgmAllFieldGroup_t
 */
#define dcgmAllFieldGroup_version dcgmAllFieldGroup_version1

/**
 * Structure to represent error attributes
 */
typedef struct
{
    unsigned int gpuId; //!< Represents GPU ID
    short fieldId;      //!< One of DCGM_FI_?
    int status;         //!< One of DCGM_ST_?
} dcgmErrorInfo_t;

/**
 * Represents a set of memory, SM, and video clocks for a device. This can be current values or a target values
 * based on context
 */
typedef struct
{
    int version;           //!< Version Number (dcgmClockSet_version)
    unsigned int memClock; //!< Memory Clock (Memory Clock value OR DCGM_INT32_BLANK to Ignore/Use compatible
                           //!< value with smClk)
    unsigned int smClock;  //!< SM Clock (SM Clock value OR DCGM_INT32_BLANK to Ignore/Use compatible value with memClk)
} dcgmClockSet_v1;

/**
 * Typedef for \ref dcgmClockSet_v1
 */
typedef dcgmClockSet_v1 dcgmClockSet_t;

/**
 * Version 1 for \ref dcgmClockSet_v1
 */
#define dcgmClockSet_version1 MAKE_DCGM_VERSION(dcgmClockSet_v1, 1)

/**
 * Latest version for \ref dcgmClockSet_t
 */
#define dcgmClockSet_version dcgmClockSet_version1

/**
 * Represents list of supported clock sets for a device
 */
typedef struct
{
    unsigned int version;                     //!< Version Number (dcgmDeviceSupportedClockSets_version)
    unsigned int count;                       //!< Number of supported clocks
    dcgmClockSet_t clockSet[DCGM_MAX_CLOCKS]; //!< Valid clock sets for the device. Upto \ref count entries are filled
} dcgmDeviceSupportedClockSets_v1;

/**
 * Typedef for \ref dcgmDeviceSupportedClockSets_v1
 */
typedef dcgmDeviceSupportedClockSets_v1 dcgmDeviceSupportedClockSets_t;

/**
 * Version 1 for \ref dcgmDeviceSupportedClockSets_v1
 */
#define dcgmDeviceSupportedClockSets_version1 MAKE_DCGM_VERSION(dcgmDeviceSupportedClockSets_v1, 1)

/**
 * Latest version for \ref dcgmDeviceSupportedClockSets_t
 */
#define dcgmDeviceSupportedClockSets_version dcgmDeviceSupportedClockSets_version1

/**
 * Represents accounting data for one process
 */
typedef struct
{
    unsigned int version;              //!< Version Number. Should match dcgmDevicePidAccountingStats_version
    unsigned int pid;                  //!< Process id of the process these stats are for
    unsigned int gpuUtilization;       //!< Percent of time over the process's lifetime during which one or more kernels
                                       //!< was executing on the GPU.
                                       //!< Set to DCGM_INT32_NOT_SUPPORTED if is not supported
    unsigned int memoryUtilization;    //!< Percent of time over the process's lifetime during which global (device)
                                       //!< memory was being read or written.
                                       //!< Set to DCGM_INT32_NOT_SUPPORTED if is not supported
    unsigned long long maxMemoryUsage; //!< Maximum total memory in bytes that was ever allocated by the process.
                                       //!< Set to DCGM_INT64_NOT_SUPPORTED if is not supported
    unsigned long long startTimestamp; //!< CPU Timestamp in usec representing start time for the process
    unsigned long long activeTimeUsec; //!< Amount of time in usec during which the compute context was active.
                                       //!< Note that this does not mean the context was being used. endTimestamp
                                       //!< can be computed as startTimestamp + activeTime
} dcgmDevicePidAccountingStats_v1;

/**
 * Typedef for \ref dcgmDevicePidAccountingStats_v1
 */
typedef dcgmDevicePidAccountingStats_v1 dcgmDevicePidAccountingStats_t;

/**
 * Version 1 for \ref dcgmDevicePidAccountingStats_v1
 */
#define dcgmDevicePidAccountingStats_version1 MAKE_DCGM_VERSION(dcgmDevicePidAccountingStats_v1, 1)

/**
 * Latest version for \ref dcgmDevicePidAccountingStats_t
 */
#define dcgmDevicePidAccountingStats_version dcgmDevicePidAccountingStats_version1

/**
 * Represents thermal information
 */
typedef struct
{
    unsigned int version;      //!< Version Number
    unsigned int slowdownTemp; //!< Slowdown temperature
    unsigned int shutdownTemp; //!< Shutdown temperature
} dcgmDeviceThermals_v1;

/**
 * Typedef for \ref dcgmDeviceThermals_v1
 */
typedef dcgmDeviceThermals_v1 dcgmDeviceThermals_t;

/**
 * Version 1 for \ref dcgmDeviceThermals_v1
 */
#define dcgmDeviceThermals_version1 MAKE_DCGM_VERSION(dcgmDeviceThermals_v1, 1)

/**
 * Latest version for \ref dcgmDeviceThermals_t
 */
#define dcgmDeviceThermals_version dcgmDeviceThermals_version1

/**
 * Represents various power limits
 */
typedef struct
{
    unsigned int version;            //!< Version Number
    unsigned int curPowerLimit;      //!< Power management limit associated with this device (in W)
    unsigned int defaultPowerLimit;  //!< Power management limit effective at device boot (in W)
    unsigned int enforcedPowerLimit; //!< Effective power limit that the driver enforces after taking into account
                                     //!< all limiters (in W)
    unsigned int minPowerLimit;      //!< Minimum power management limit (in W)
    unsigned int maxPowerLimit;      //!< Maximum power management limit (in W)
} dcgmDevicePowerLimits_v1;

/**
 * Typedef for \ref dcgmDevicePowerLimits_v1
 */
typedef dcgmDevicePowerLimits_v1 dcgmDevicePowerLimits_t;

/**
 * Version 1 for \ref dcgmDevicePowerLimits_v1
 */
#define dcgmDevicePowerLimits_version1 MAKE_DCGM_VERSION(dcgmDevicePowerLimits_v1, 1)

/**
 * Latest version for \ref dcgmDevicePowerLimits_t
 */
#define dcgmDevicePowerLimits_version dcgmDevicePowerLimits_version1

/**
 * Represents device identifiers
 */
typedef struct
{
    unsigned int version;                          //!< Version Number (dcgmDeviceIdentifiers_version)
    char brandName[DCGM_MAX_STR_LENGTH];           //!< Brand Name
    char deviceName[DCGM_MAX_STR_LENGTH];          //!< Name of the device
    char pciBusId[DCGM_MAX_STR_LENGTH];            //!< PCI Bus ID
    char serial[DCGM_MAX_STR_LENGTH];              //!< Serial for the device
    char uuid[DCGM_MAX_STR_LENGTH];                //!< UUID for the device
    char vbios[DCGM_MAX_STR_LENGTH];               //!< VBIOS version
    char inforomImageVersion[DCGM_MAX_STR_LENGTH]; //!< Inforom Image version
    unsigned int pciDeviceId;                      //!< The combined 16-bit device id and 16-bit vendor id
    unsigned int pciSubSystemId;                   //!< The 32-bit Sub System Device ID
    char driverVersion[DCGM_MAX_STR_LENGTH];       //!< Driver Version
    unsigned int virtualizationMode;               //!< Virtualization Mode
} dcgmDeviceIdentifiers_v1;

/**
 * Typedef for \ref dcgmDeviceIdentifiers_v1
 */
typedef dcgmDeviceIdentifiers_v1 dcgmDeviceIdentifiers_t;

/**
 * Version 1 for \ref dcgmDeviceIdentifiers_v1
 */
#define dcgmDeviceIdentifiers_version1 MAKE_DCGM_VERSION(dcgmDeviceIdentifiers_v1, 1)

/**
 * Latest version for \ref dcgmDeviceIdentifiers_t
 */
#define dcgmDeviceIdentifiers_version dcgmDeviceIdentifiers_version1

/**
 * Represents device memory and usage
 */
typedef struct
{
    unsigned int version;   //!< Version Number (dcgmDeviceMemoryUsage_version)
    unsigned int bar1Total; //!< Total BAR1 size in megabytes
    unsigned int fbTotal;   //!< Total framebuffer memory in megabytes
    unsigned int fbUsed;    //!< Used framebuffer memory in megabytes
    unsigned int fbFree;    //!< Free framebuffer memory in megabytes
} dcgmDeviceMemoryUsage_v1;

/**
 * Typedef for \ref dcgmDeviceMemoryUsage_v1
 */
typedef dcgmDeviceMemoryUsage_v1 dcgmDeviceMemoryUsage_t;

/**
 * Version 1 for \ref dcgmDeviceMemoryUsage_v1
 */
#define dcgmDeviceMemoryUsage_version1 MAKE_DCGM_VERSION(dcgmDeviceMemoryUsage_v1, 1)

/**
 * Latest version for \ref dcgmDeviceMemoryUsage_t
 */
#define dcgmDeviceMemoryUsage_version dcgmDeviceMemoryUsage_version1

/**
 * Represents utilization values for vGPUs running on the device
 */
typedef struct
{
    unsigned int version; //!< Version Number (dcgmDeviceVgpuUtilInfo_version)
    unsigned int vgpuId;  //!< vGPU instance ID
    unsigned int smUtil;  //!< GPU utilization for vGPU
    unsigned int memUtil; //!< Memory utilization for vGPU
    unsigned int encUtil; //!< Encoder utilization for vGPU
    unsigned int decUtil; //!< Decoder utilization for vGPU
} dcgmDeviceVgpuUtilInfo_v1;

/**
 * Typedef for \ref dcgmDeviceVgpuUtilInfo_v1
 */
typedef dcgmDeviceVgpuUtilInfo_v1 dcgmDeviceVgpuUtilInfo_t;

/**
 * Version 1 for \ref dcgmDeviceVgpuUtilInfo_v1
 */
#define dcgmDeviceVgpuUtilInfo_version1 MAKE_DCGM_VERSION(dcgmDeviceVgpuUtilInfo_v1, 1)

/**
 * Latest version for \ref dcgmDeviceVgpuUtilInfo_t
 */
#define dcgmDeviceVgpuUtilInfo_version dcgmDeviceVgpuUtilInfo_version1

/**
 * Represents current encoder statistics for the given device/vGPU instance
 */
typedef struct
{
    unsigned int version;        //!< Version Number (dcgmDeviceEncStats_version)
    unsigned int sessionCount;   //!< Count of active encoder sessions
    unsigned int averageFps;     //!< Trailing average FPS of all active sessions
    unsigned int averageLatency; //!< Encode latency in milliseconds
} dcgmDeviceEncStats_v1;

/**
 * Typedef for \ref dcgmDeviceEncStats_v1
 */
typedef dcgmDeviceEncStats_v1 dcgmDeviceEncStats_t;

/**
 * Version 1 for \ref dcgmDeviceEncStats_v1
 */
#define dcgmDeviceEncStats_version1 MAKE_DCGM_VERSION(dcgmDeviceEncStats_v1, 1)

/**
 * Latest version for \ref dcgmDeviceEncStats_t
 */
#define dcgmDeviceEncStats_version dcgmDeviceEncStats_version1

/**
 * Represents current frame buffer capture sessions statistics for the given device/vGPU instance
 */
typedef struct
{
    unsigned int version;        //!< Version Number (dcgmDeviceFbcStats_version)
    unsigned int sessionCount;   //!< Count of active FBC sessions
    unsigned int averageFps;     //!< Moving average new frames captured per second
    unsigned int averageLatency; //!< Moving average new frame capture latency in microseconds
} dcgmDeviceFbcStats_v1;

/**
 * Typedef for \ref dcgmDeviceFbcStats_v1
 */
typedef dcgmDeviceFbcStats_v1 dcgmDeviceFbcStats_t;

/**
 * Version 1 for \ref dcgmDeviceFbcStats_v1
 */
#define dcgmDeviceFbcStats_version1 MAKE_DCGM_VERSION(dcgmDeviceFbcStats_v1, 1)

/**
 * Latest version for \ref dcgmDeviceEncStats_t
 */
#define dcgmDeviceFbcStats_version dcgmDeviceFbcStats_version1

/*
 * Represents frame buffer capture session type
 */
typedef enum dcgmFBCSessionType_enum
{
    DCGM_FBC_SESSION_TYPE_UNKNOWN = 0, //!< Unknown
    DCGM_FBC_SESSION_TYPE_TOSYS,       //!< FB capture for a system buffer
    DCGM_FBC_SESSION_TYPE_CUDA,        //!< FB capture for a cuda buffer
    DCGM_FBC_SESSION_TYPE_VID,         //!< FB capture for a Vid buffer
    DCGM_FBC_SESSION_TYPE_HWENC,       //!< FB capture for a NVENC HW buffer
} dcgmFBCSessionType_t;

/**
 * Represents information about active FBC session on the given device/vGPU instance
 */
typedef struct
{
    unsigned int version;             //!< Version Number (dcgmDeviceFbcSessionInfo_version)
    unsigned int sessionId;           //!< Unique session ID
    unsigned int pid;                 //!< Owning process ID
    unsigned int vgpuId;              //!< vGPU instance ID (only valid on vGPU hosts, otherwise zero)
    unsigned int displayOrdinal;      //!< Display identifier
    dcgmFBCSessionType_t sessionType; //!< Type of frame buffer capture session
    unsigned int sessionFlags;        //!< Session flags
    unsigned int hMaxResolution;      //!< Max horizontal resolution supported by the capture session
    unsigned int vMaxResolution;      //!< Max vertical resolution supported by the capture session
    unsigned int hResolution;         //!< Horizontal resolution requested by caller in capture call
    unsigned int vResolution;         //!< Vertical resolution requested by caller in capture call
    unsigned int averageFps;          //!< Moving average new frames captured per second
    unsigned int averageLatency;      //!< Moving average new frame capture latency in microseconds
} dcgmDeviceFbcSessionInfo_v1;

/**
 * Typedef for \ref dcgmDeviceFbcSessionInfo_v1
 */
typedef dcgmDeviceFbcSessionInfo_v1 dcgmDeviceFbcSessionInfo_t;

/**
 * Version 1 for \ref dcgmDeviceFbcSessionInfo_v1
 */
#define dcgmDeviceFbcSessionInfo_version1 MAKE_DCGM_VERSION(dcgmDeviceFbcSessionInfo_v1, 1)

/**
 * Latest version for \ref dcgmDeviceFbcSessionInfo_t
 */
#define dcgmDeviceFbcSessionInfo_version dcgmDeviceFbcSessionInfo_version1

/**
 * Represents all the active FBC sessions on the given device/vGPU instance
 */
typedef struct
{
    unsigned int version;                                          //!< Version Number (dcgmDeviceFbcSessions_version)
    unsigned int sessionCount;                                     //!< Count of active FBC sessions
    dcgmDeviceFbcSessionInfo_t sessionInfo[DCGM_MAX_FBC_SESSIONS]; //!< Info about the active FBC session
} dcgmDeviceFbcSessions_v1;

/**
 * Typedef for \ref dcgmDeviceFbcSessions_v1
 */
typedef dcgmDeviceFbcSessions_v1 dcgmDeviceFbcSessions_t;

/**
 * Version 1 for \ref dcgmDeviceFbcSessions_v1
 */
#define dcgmDeviceFbcSessions_version1 MAKE_DCGM_VERSION(dcgmDeviceFbcSessions_v1, 1)

/**
 * Latest version for \ref dcgmDeviceFbcSessions_t
 */
#define dcgmDeviceFbcSessions_version dcgmDeviceFbcSessions_version1

/*
 * Represents type of encoder for capacity can be queried
 */
typedef enum dcgmEncoderQueryType_enum
{
    DCGM_ENCODER_QUERY_H264 = 0,
    DCGM_ENCODER_QUERY_HEVC = 1
} dcgmEncoderType_t;

/**
 * Represents information about active encoder sessions on the given vGPU instance
 */
typedef struct
{
    unsigned int version; //!< Version Number (dcgmDeviceVgpuEncSessions_version)
    union
    {
        unsigned int vgpuId; //!< vGPU instance ID
        unsigned int sessionCount;
    } encoderSessionInfo;
    unsigned int sessionId;      //!< Unique session ID
    unsigned int pid;            //!< Process ID
    dcgmEncoderType_t codecType; //!< Video encoder type
    unsigned int hResolution;    //!< Current encode horizontal resolution
    unsigned int vResolution;    //!< Current encode vertical resolution
    unsigned int averageFps;     //!< Moving average encode frames per second
    unsigned int averageLatency; //!< Moving average encode latency in milliseconds
} dcgmDeviceVgpuEncSessions_v1;

/**
 * Typedef for \ref dcgmDeviceVgpuEncSessions_v1
 */
typedef dcgmDeviceVgpuEncSessions_v1 dcgmDeviceVgpuEncSessions_t;

/**
 * Version 1 for \ref dcgmDeviceVgpuEncSessions_v1
 */
#define dcgmDeviceVgpuEncSessions_version1 MAKE_DCGM_VERSION(dcgmDeviceVgpuEncSessions_v1, 1)

/**
 * Latest version for \ref dcgmDeviceVgpuEncSessions_t
 */
#define dcgmDeviceVgpuEncSessions_version dcgmDeviceVgpuEncSessions_version1

/**
 * Represents utilization values for processes running in vGPU VMs using the device
 */
typedef struct
{
    unsigned int version; //!< Version Number (dcgmDeviceVgpuProcessUtilInfo_version)
    union
    {
        unsigned int vgpuId;                  //!< vGPU instance ID
        unsigned int vgpuProcessSamplesCount; //!< Count of processes running in the vGPU VM,for which utilization
                                              //!< rates are being reported in this cycle.
    } vgpuProcessUtilInfo;
    unsigned int pid;                             //!< Process ID of the process running in the vGPU VM.
    char processName[DCGM_VGPU_NAME_BUFFER_SIZE]; //!< Process Name of process running in the vGPU VM.
    unsigned int smUtil;                          //!< GPU utilization of process running in the vGPU VM.
    unsigned int memUtil;                         //!< Memory utilization of process running in the vGPU VM.
    unsigned int encUtil;                         //!< Encoder utilization of process running in the vGPU VM.
    unsigned int decUtil;                         //!< Decoder utilization of process running in the vGPU VM.
} dcgmDeviceVgpuProcessUtilInfo_v1;

/**
 * Typedef for \ref dcgmDeviceVgpuProcessUtilInfo_v1
 */
typedef dcgmDeviceVgpuProcessUtilInfo_v1 dcgmDeviceVgpuProcessUtilInfo_t;

/**
 * Version 1 for \ref dcgmDeviceVgpuProcessUtilInfo_v1
 */
#define dcgmDeviceVgpuProcessUtilInfo_version1 MAKE_DCGM_VERSION(dcgmDeviceVgpuProcessUtilInfo_v1, 1)

/**
 * Latest version for \ref dcgmDeviceVgpuProcessUtilInfo_t
 */
#define dcgmDeviceVgpuProcessUtilInfo_version dcgmDeviceVgpuProcessUtilInfo_version1

/**
 * Represents static info related to vGPUs supported on the device.
 */
typedef struct
{
    unsigned int version; //!< Version number (dcgmDeviceVgpuTypeInfo_version)
    union
    {
        unsigned int vgpuTypeId;
        unsigned int supportedVgpuTypeCount;
    } vgpuTypeInfo;                                      //!< vGPU type ID and Supported vGPU type count
    char vgpuTypeName[DCGM_VGPU_NAME_BUFFER_SIZE];       //!< vGPU type Name
    char vgpuTypeClass[DCGM_VGPU_NAME_BUFFER_SIZE];      //!< Class of vGPU type
    char vgpuTypeLicense[DCGM_GRID_LICENSE_BUFFER_SIZE]; //!< license of vGPU type
    int deviceId;                                        //!< device ID of vGPU type
    int subsystemId;                                     //!< Subsystem ID of vGPU type
    int numDisplayHeads;                                 //!< Count of vGPU's supported display heads
    int maxInstances;   //!< maximum number of vGPU instances creatable on a device for given vGPU type
    int frameRateLimit; //!< Frame rate limit value of the vGPU type
    int maxResolutionX; //!< vGPU display head's maximum supported resolution in X dimension
    int maxResolutionY; //!< vGPU display head's maximum supported resolution in Y dimension
    int fbTotal;        //!< vGPU Total framebuffer size in megabytes
} dcgmDeviceVgpuTypeInfo_v1;

/**
 * Version 1 for \ref dcgmDeviceVgpuTypeInfo_v1
 */
#define dcgmDeviceVgpuTypeInfo_version1 MAKE_DCGM_VERSION(dcgmDeviceVgpuTypeInfo_v1, 1)

typedef struct
{
    unsigned int version; //!< Version number (dcgmDeviceVgpuTypeInfo_version2)
    union
    {
        unsigned int vgpuTypeId;
        unsigned int supportedVgpuTypeCount;
    } vgpuTypeInfo;                                      //!< vGPU type ID and Supported vGPU type count
    char vgpuTypeName[DCGM_VGPU_NAME_BUFFER_SIZE];       //!< vGPU type Name
    char vgpuTypeClass[DCGM_VGPU_NAME_BUFFER_SIZE];      //!< Class of vGPU type
    char vgpuTypeLicense[DCGM_GRID_LICENSE_BUFFER_SIZE]; //!< license of vGPU type
    int deviceId;                                        //!< device ID of vGPU type
    int subsystemId;                                     //!< Subsystem ID of vGPU type
    int numDisplayHeads;                                 //!< Count of vGPU's supported display heads
    int maxInstances;         //!< maximum number of vGPU instances creatable on a device for given vGPU type
    int frameRateLimit;       //!< Frame rate limit value of the vGPU type
    int maxResolutionX;       //!< vGPU display head's maximum supported resolution in X dimension
    int maxResolutionY;       //!< vGPU display head's maximum supported resolution in Y dimension
    int fbTotal;              //!< vGPU Total framebuffer size in megabytes
    int gpuInstanceProfileId; //!< GPU Instance Profile ID for the given vGPU type
} dcgmDeviceVgpuTypeInfo_v2;

/**
 * Typedef for \ref dcgmDeviceVgpuTypeInfo_v2
 */
typedef dcgmDeviceVgpuTypeInfo_v2 dcgmDeviceVgpuTypeInfo_t;

/**
 * Version 2 for \ref dcgmDeviceVgpuTypeInfo_v2
 */
#define dcgmDeviceVgpuTypeInfo_version2 MAKE_DCGM_VERSION(dcgmDeviceVgpuTypeInfo_v2, 2)

/**
 * Latest version for \ref dcgmDeviceVgpuTypeInfo_t
 */
#define dcgmDeviceVgpuTypeInfo_version dcgmDeviceVgpuTypeInfo_version2

/**
 * Represents the info related to vGPUs supported on the device.
 */
typedef struct
{
    unsigned int version;              //!< Version number (dcgmDeviceSupportedVgpuTypeInfo_version)
    unsigned long long deviceId;       //!< device ID of vGPU type
    unsigned long long subsystemId;    //!< Subsystem ID of vGPU type
    unsigned int numDisplayHeads;      //!< Count of vGPU's supported display heads
    unsigned int maxInstances;         //!< maximum number of vGPU instances creatable on a device for given vGPU type
    unsigned int frameRateLimit;       //!< Frame rate limit value of the vGPU type
    unsigned int maxResolutionX;       //!< vGPU display head's maximum supported resolution in X dimension
    unsigned int maxResolutionY;       //!< vGPU display head's maximum supported resolution in Y dimension
    unsigned long long fbTotal;        //!< vGPU Total framebuffer size in megabytes
    unsigned int gpuInstanceProfileId; //!< GPU Instance Profile ID for the given vGPU type
} dcgmDeviceSupportedVgpuTypeInfo_v1;

/**
 * Typedef for \ref dcgmDeviceSupportedVgpuTypeInfo_v1
 */
typedef dcgmDeviceSupportedVgpuTypeInfo_v1 dcgmDeviceSupportedVgpuTypeInfo_t;

/**
 * Version 1 for \ref dcgmDeviceSupportedVgpuTypeInfo_v1
 */
#define dcgmDeviceSupportedVgpuTypeInfo_version1 MAKE_DCGM_VERSION(dcgmDeviceSupportedVgpuTypeInfo_v1, 1)

/**
 * Latest version for \ref dcgmDeviceSupportedVgpuTypeInfo_t
 */
#define dcgmDeviceSupportedVgpuTypeInfo_version dcgmDeviceSupportedVgpuTypeInfo_version1

typedef struct
{
    unsigned int version;
    unsigned int persistenceModeEnabled;
    unsigned int migModeEnabled;
    unsigned int confidentialComputeMode;
} dcgmDeviceSettings_v2;

typedef dcgmDeviceSettings_v2 dcgmDeviceSettings_t;

#define dcgmDeviceSettings_version2 MAKE_DCGM_VERSION(dcgmDeviceSettings_v2, 2)

#define dcgmDeviceSettings_version dcgmDeviceSettings_version2

typedef struct
{
    unsigned int version;                     //!< Version number (dcgmDeviceAttributes_version)
    dcgmDeviceSupportedClockSets_t clockSets; //!< Supported clocks for the device
    dcgmDeviceThermals_t thermalSettings;     //!< Thermal settings for the device
    dcgmDevicePowerLimits_t powerLimits;      //!< Various power limits for the device
    dcgmDeviceIdentifiers_t identifiers;      //!< Identifiers for the device
    dcgmDeviceMemoryUsage_t memoryUsage;      //!< Memory usage info for the device
    dcgmDeviceSettings_v2 settings;           //!< Basic device settings
} dcgmDeviceAttributes_v3;

/**
 * Typedef for \ref dcgmDeviceAttributes_v3
 */
typedef dcgmDeviceAttributes_v3 dcgmDeviceAttributes_t;

/**
 * Version 3 for \ref dcgmDeviceAttributes_v3
 */
#define dcgmDeviceAttributes_version3 MAKE_DCGM_VERSION(dcgmDeviceAttributes_v3, 3)

/**
 * Latest version for \ref dcgmDeviceAttributes_t
 */
#define dcgmDeviceAttributes_version dcgmDeviceAttributes_version3

/**
 * Structure to represent attributes info for a MIG device
 */
typedef struct
{
    unsigned int version;                   //!< Version Number (dcgmDeviceMigAttributesInfo_version)
    unsigned int gpuInstanceId;             //!< GPU instance ID
    unsigned int computeInstanceId;         //!< Compute instance ID
    unsigned int multiprocessorCount;       //!< Streaming Multiprocessor count
    unsigned int sharedCopyEngineCount;     //!< Shared Copy Engine count
    unsigned int sharedDecoderCount;        //!< Shared Decoder Engine count
    unsigned int sharedEncoderCount;        //!< Shared Encoder Engine count
    unsigned int sharedJpegCount;           //!< Shared JPEG Engine count
    unsigned int sharedOfaCount;            //!< Shared OFA Engine count
    unsigned int gpuInstanceSliceCount;     //!< GPU instance slice count
    unsigned int computeInstanceSliceCount; //!< Compute instance slice count
    unsigned long long memorySizeMB;        //!< Device memory size (in MiB)
} dcgmDeviceMigAttributesInfo_v1;

/**
 * Typedef for \ref dcgmDeviceMigAttributesInfo_v1
 */
typedef dcgmDeviceMigAttributesInfo_v1 dcgmDeviceMigAttributesInfo_t;

/**
 * Version 1 for \ref dcgmDeviceMigAttributesInfo_v1
 */
#define dcgmDeviceMigAttributesInfo_version1 MAKE_DCGM_VERSION(dcgmDeviceMigAttributesInfo_v1, 1)

/**
 * Latest version for \ref dcgmDeviceMigAttributesInfo_t
 */
#define dcgmDeviceMigAttributesInfo_version dcgmDeviceMigAttributesInfo_version1

/**
 * Structure to represent attributes for a MIG device
 */
typedef struct
{
    unsigned int version;                             //!< Version Number (dcgmDeviceMigAttributes_version)
    unsigned int migDevicesCount;                     //!< Count of MIG devices
    dcgmDeviceMigAttributesInfo_v1 migAttributesInfo; //!< MIG attributes information
} dcgmDeviceMigAttributes_v1;

/**
 * Typedef for \ref dcgmDeviceMigAttributes_v1
 */
typedef dcgmDeviceMigAttributes_v1 dcgmDeviceMigAttributes_t;

/**
 * Version 1 for \ref dcgmDeviceMigAttributes_v1
 */
#define dcgmDeviceMigAttributes_version1 MAKE_DCGM_VERSION(dcgmDeviceMigAttributes_v1, 1)

/**
 * Latest version for \ref dcgmDeviceMigAttributes_t
 */
#define dcgmDeviceMigAttributes_version dcgmDeviceMigAttributes_version1

/**
 * Structure to represent GPU instance profile information
 */
typedef struct
{
    unsigned int version;             //!< Version Number (dcgmGpuInstanceProfileInfo_version)
    unsigned int id;                  //!< Unique profile ID within the device
    unsigned int isP2pSupported;      //!< Peer-to-Peer support
    unsigned int sliceCount;          //!< GPU Slice count
    unsigned int instanceCount;       //!< GPU instance count
    unsigned int multiprocessorCount; //!< Streaming Multiprocessor count
    unsigned int copyEngineCount;     //!< Copy Engine count
    unsigned int decoderCount;        //!< Decoder Engine count
    unsigned int encoderCount;        //!< Encoder Engine count
    unsigned int jpegCount;           //!< JPEG Engine count
    unsigned int ofaCount;            //!< OFA Engine count
    unsigned long long memorySizeMB;  //!< Memory size in MBytes
} dcgmGpuInstanceProfileInfo_v1;

/**
 * Typedef for \ref dcgmGpuInstanceProfileInfo_v1
 */
typedef dcgmGpuInstanceProfileInfo_v1 dcgmGpuInstanceProfileInfo_t;

/**
 * Version 1 for \ref dcgmGpuInstanceProfileInfo_v1
 */
#define dcgmGpuInstanceProfileInfo_version1 MAKE_DCGM_VERSION(dcgmGpuInstanceProfileInfo_v1, 1)

/**
 * Latest version for \ref dcgmGpuInstanceProfileInfo_t
 */
#define dcgmGpuInstanceProfileInfo_version dcgmGpuInstanceProfileInfo_version1

/**
 * Structure to represent GPU instance profiles
 */
typedef struct
{
    unsigned int version;                      //!< Version Number (dcgmGpuInstanceProfiles_version)
    unsigned int profileCount;                 //!< Profile count
    dcgmGpuInstanceProfileInfo_v1 profileInfo; //!< GPU instance profile information
} dcgmGpuInstanceProfiles_v1;

/**
 * Typedef for \ref dcgmGpuInstanceProfiles_v1
 */
typedef dcgmGpuInstanceProfiles_v1 dcgmGpuInstanceProfiles_t;

/**
 * Version 1 for \ref dcgmGpuInstanceProfiles_v1
 */
#define dcgmGpuInstanceProfiles_version1 MAKE_DCGM_VERSION(dcgmGpuInstanceProfiles_v1, 1)

/**
 * Latest version for \ref dcgmGpuInstanceProfiles_t
 */
#define dcgmGpuInstanceProfiles_version dcgmGpuInstanceProfiles_version1

/**
 * Structure to represent Compute instance profile information
 */
typedef struct
{
    unsigned int version;               //!< Version Number (dcgmComputeInstanceProfileInfo_version)
    unsigned int gpuInstanceId;         //!< GPU instance ID
    unsigned int id;                    //!< Unique profile ID within the GPU instance
    unsigned int sliceCount;            //!< GPU Slice count
    unsigned int instanceCount;         //!< Compute instance count
    unsigned int multiprocessorCount;   //!< Streaming Multiprocessor count
    unsigned int sharedCopyEngineCount; //!< Shared Copy Engine count
    unsigned int sharedDecoderCount;    //!< Shared Decoder Engine count
    unsigned int sharedEncoderCount;    //!< Shared Encoder Engine count
    unsigned int sharedJpegCount;       //!< Shared JPEG Engine count
    unsigned int sharedOfaCount;        //!< Shared OFA Engine count
} dcgmComputeInstanceProfileInfo_v1;

/**
 * Typedef for \ref dcgmComputeInstanceProfileInfo_v1
 */
typedef dcgmComputeInstanceProfileInfo_v1 dcgmComputeInstanceProfileInfo_t;

/**
 * Version 1 for \ref dcgmComputeInstanceProfileInfo_v1
 */
#define dcgmComputeInstanceProfileInfo_version1 MAKE_DCGM_VERSION(dcgmComputeInstanceProfileInfo_v1, 1)

/**
 * Latest version for \ref dcgmComputeInstanceProfileInfo_t
 */
#define dcgmComputeInstanceProfileInfo_version dcgmComputeInstanceProfileInfo_version1

/**
 * Structure to represent Compute instance profiles
 */
typedef struct
{
    unsigned int version;                          //!< Version Number (dcgmComputeInstanceProfiles_version)
    unsigned int profileCount;                     //!< Profile count
    dcgmComputeInstanceProfileInfo_v1 profileInfo; //!< Compute instance profile information
} dcgmComputeInstanceProfiles_v1;

/**
 * Typedef for \ref dcgmComputeInstanceProfiles_v1
 */
typedef dcgmComputeInstanceProfiles_v1 dcgmComputeInstanceProfiles_t;

/**
 * Version 1 for \ref dcgmComputeInstanceProfiles_v1
 */
#define dcgmComputeInstanceProfiles_version1 MAKE_DCGM_VERSION(dcgmComputeInstanceProfiles_v1, 1)

/**
 * Latest version for \ref dcgmComputeInstanceProfiles_t
 */
#define dcgmComputeInstanceProfiles_version dcgmComputeInstanceProfiles_version1

/**
 * Maximum number of vGPU types per physical GPU
 */
#define DCGM_MAX_VGPU_TYPES_PER_PGPU 32

/**
 * Represents the size of a buffer that holds string related to attributes specific to vGPU instance
 */
#define DCGM_DEVICE_UUID_BUFFER_SIZE 80

/**
 * Used to represent Performance state settings
 */
typedef struct
{
    unsigned int syncBoost;      //!< Sync Boost Mode (0: Disabled, 1 : Enabled, DCGM_INT32_BLANK : Ignored). Note that
                                 //!< using this setting may result in lower clocks than targetClocks
    dcgmClockSet_t targetClocks; //!< Target clocks. Set smClock and memClock to DCGM_INT32_BLANK to ignore/use
                                 //!< compatible values. For GPUs > Maxwell, setting this implies autoBoost=0
} dcgmConfigPerfStateSettings_t;

/**
 * Used to represents the power capping limit for each GPU in the group or to represent the power
 * budget for the entire group
 */
typedef struct
{
    dcgmConfigPowerLimitType_t type; //!< Flag to represent power cap for each GPU or power budget for the group of GPUs
    unsigned int val;                //!< Power Limit in Watts (Set a value OR DCGM_INT32_BLANK to Ignore)
} dcgmConfigPowerLimit_t;

/**
 * Structure to represent default and target configuration for a device
 */
typedef struct
{
    unsigned int version;     //!< Version number (dcgmConfig_version)
    unsigned int gpuId;       //!< GPU ID
    unsigned int eccMode;     //!< ECC Mode  (0: Disabled, 1 : Enabled, DCGM_INT32_BLANK : Ignored)
    unsigned int computeMode; //!< Compute Mode (One of DCGM_CONFIG_COMPUTEMODE_? OR DCGM_INT32_BLANK to Ignore)
    dcgmConfigPerfStateSettings_t perfState; //!< Performance State Settings (clocks / boost mode)
    dcgmConfigPowerLimit_t powerLimit;       //!< Power Limits
} dcgmConfig_v1;

/**
 * Typedef for \ref dcgmConfig_v1
 */
typedef dcgmConfig_v1 dcgmConfig_t;

/**
 * Version 1 for \ref dcgmConfig_v1
 */
#define dcgmConfig_version1 MAKE_DCGM_VERSION(dcgmConfig_v1, 1)

/**
 * Latest version for \ref dcgmConfig_t
 */
#define dcgmConfig_version dcgmConfig_version1

/**
 * Represents a callback to receive updates from asynchronous functions.
 * Currently the only implemented callback function is dcgmPolicyRegister
 * and the void * data will be a pointer to dcgmPolicyCallbackResponse_t.
 * Ex.
 * dcgmPolicyCallbackResponse_t *callbackResponse = (dcgmPolicyCallbackResponse_t *) userData;
 *
 */
typedef int (*fpRecvUpdates)(void *userData);

/*Remove from doxygen documentation
 *
 * Define the structure that contains specific policy information
 */
typedef struct
{
    // version must always be first
    unsigned int version; //!< Version number (dcgmPolicyViolation_version)

    unsigned int notifyOnEccDbe;          //!< true/false notification on ECC Double Bit Errors
    unsigned int notifyOnPciEvent;        //!< true/false notification on PCI Events
    unsigned int notifyOnMaxRetiredPages; //!< number of retired pages to occur before notification
} dcgmPolicyViolation_v1;

/*Remove from doxygen documentation
 *
 * Represents the versioning for the dcgmPolicyViolation_v1 structure
 */

/*
 * Typedef for \ref dcgmPolicyViolation_v1
 */
typedef dcgmPolicyViolation_v1 dcgmPolicyViolation_t;

/*
 * Version 1 for \ref dcgmPolicyViolation_v1
 */
#define dcgmPolicyViolation_version1 MAKE_DCGM_VERSION(dcgmPolicyViolation_v1, 1)

/*
 * Latest version for \ref dcgmPolicyViolation_t
 */
#define dcgmPolicyViolation_version dcgmPolicyViolation_version1

/**
 * Enumeration for policy conditions.
 * When used as part of dcgmPolicy_t these have corresponding parameters to
 * allow them to be switched on/off or set specific violation thresholds
 */
typedef enum dcgmPolicyConditionIdx_enum
{
    // These are sequential rather than bitwise.
    DCGM_POLICY_COND_IDX_DBE = 0,           //!< Double bit errors -- boolean in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_IDX_PCI,               //!< PCI events/errors -- boolean in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_IDX_MAX_PAGES_RETIRED, //!< Maximum number of retired pages -- number
                                            //!< required in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_IDX_THERMAL,           //!< Thermal violation -- number required in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_IDX_POWER,             //!< Power violation -- number required in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_IDX_NVLINK,            //!< NVLINK errors -- boolean in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_IDX_XID                //!< XID errors -- number required in dcgmPolicyConditionParams_t
} dcgmPolicyConditionIdx_t;

#define DCGM_POLICY_COND_IDX_MAX 7
#define DCGM_POLICY_COND_MAX     DCGM_POLICY_COND_IDX_MAX

/**
 * Bitmask enumeration for policy conditions.
 * When used as part of dcgmPolicy_t these have corresponding parameters to
 * allow them to be switched on/off or set specific violation thresholds
 */
typedef enum dcgmPolicyCondition_enum
{
    // These are bitwise rather than sequential.
    DCGM_POLICY_COND_DBE               = 0x1, //!< Double bit errors -- boolean in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_PCI               = 0x2, //!< PCI events/errors -- boolean in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_MAX_PAGES_RETIRED = 0x4, //!< Maximum number of retired pages -- number
                                              //!< required in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_THERMAL = 0x8,           //!< Thermal violation -- number required in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_POWER   = 0x10,          //!< Power violation -- number required in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_NVLINK  = 0x20,          //!< NVLINK errors -- boolean in dcgmPolicyConditionParams_t
    DCGM_POLICY_COND_XID     = 0x40,          //!< XID errors -- number required in dcgmPolicyConditionParams_t
} dcgmPolicyCondition_t;

/**
 * Structure for policy condition parameters.
 * This structure contains a tag that represents the type of the value being passed
 * as well as a "val" which is a union of the possible value types.  For example,
 * to pass a true boolean: tag = BOOL, val.boolean = 1.
 */
typedef struct dcgmPolicyConditionParams_st
{
    enum
    {
        BOOL,
        LLONG
    } tag;
    union
    {
        unsigned int boolean;
        unsigned long long llval;
    } val;
} dcgmPolicyConditionParams_t;

/**
 * Enumeration for policy modes
 */
typedef enum dcgmPolicyMode_enum
{
    DCGM_POLICY_MODE_AUTOMATED = 0, //!< automatic mode
    DCGM_POLICY_MODE_MANUAL    = 1, //!< manual mode
} dcgmPolicyMode_t;

/**
 * Enumeration for policy isolation modes
 */
typedef enum dcgmPolicyIsolation_enum
{
    DCGM_POLICY_ISOLATION_NONE = 0, //!< no isolation of GPUs on error
} dcgmPolicyIsolation_t;

/**
 * Enumeration for policy actions
 */
typedef enum dcgmPolicyAction_enum
{
    DCGM_POLICY_ACTION_NONE     = 0, //!< no action
    DCGM_POLICY_ACTION_GPURESET = 1, //!< Deprecated - perform a GPU reset on violation
} dcgmPolicyAction_t;

/**
 * Enumeration for policy validation actions
 */
typedef enum dcgmPolicyValidation_enum
{
    DCGM_POLICY_VALID_NONE     = 0, //!< no validation after an action is performed
    DCGM_POLICY_VALID_SV_SHORT = 1, //!< run a short System Validation on the system after failure
    DCGM_POLICY_VALID_SV_MED   = 2, //!< run a medium System Validation test after failure
    DCGM_POLICY_VALID_SV_LONG  = 3, //!< run a extensive System Validation test after failure
    DCGM_POLICY_VALID_SV_XLONG = 4, //!< run a more extensive System Validation test after failure
} dcgmPolicyValidation_t;

/**
 * Enumeration for policy failure responses
 */
typedef enum dcgmPolicyFailureResp_enum
{
    DCGM_POLICY_FAILURE_NONE = 0, //!< on failure of validation perform no action
} dcgmPolicyFailureResp_t;

/**
 * Structure to fill when a user queries for policy violations
 */
typedef struct
{
    unsigned int gpuId;             //!< gpu ID
    unsigned int violationOccurred; //!< a violation based on the bit values in \ref dcgmPolicyCondition_t
} dcgmPolicyViolationNotify_t;

/**
 * Define the structure that specifies a policy to be enforced for a GPU
 */
typedef struct
{
    // version must always be first
    unsigned int version; //!< version number (dcgmPolicy_version)

    dcgmPolicyCondition_t condition;   //!< Condition(s) to access \ref dcgmPolicyCondition_t
    dcgmPolicyMode_t mode;             //!< Mode of operation \ref dcgmPolicyMode_t
    dcgmPolicyIsolation_t isolation;   //!< Isolation level after a policy violation \ref dcgmPolicyIsolation_t
    dcgmPolicyAction_t action;         //!< Action to perform after a policy violation \ref dcgmPolicyAction_t action
    dcgmPolicyValidation_t validation; //!< Validation to perform after action is taken \ref dcgmPolicyValidation_t
    dcgmPolicyFailureResp_t response;  //!< Failure to validation response \ref dcgmPolicyFailureResp_t
    dcgmPolicyConditionParams_t parms[DCGM_POLICY_COND_MAX]; //!< Parameters for the \a condition fields
} dcgmPolicy_v1;

/**
 * Typedef for \ref dcgmPolicy_v1
 */
typedef dcgmPolicy_v1 dcgmPolicy_t;

/**
 * Version 1 for \ref dcgmPolicy_v1
 */
#define dcgmPolicy_version1 MAKE_DCGM_VERSION(dcgmPolicy_v1, 1)

/**
 * Latest version for \ref dcgmPolicy_t
 */
#define dcgmPolicy_version dcgmPolicy_version1


/**
 * Define the ECC DBE return structure
 */
typedef struct
{
    long long timestamp; //!< timestamp of the error
    enum
    {
        L1,
        L2,
        DEVICE,
        REGISTER,
        TEXTURE
    } location;             //!< location of the error
    unsigned int numerrors; //!< number of errors
} dcgmPolicyConditionDbe_t;

/**
 * Define the PCI replay error return structure
 */
typedef struct
{
    long long timestamp;  //!< timestamp of the error
    unsigned int counter; //!< value of the PCIe replay counter
} dcgmPolicyConditionPci_t;

/**
 * Define the maximum pending retired pages limit return structure
 */
typedef struct
{
    long long timestamp;   //!< timestamp of the error
    unsigned int sbepages; //!< number of pending pages due to SBE
    unsigned int dbepages; //!< number of pending pages due to DBE
} dcgmPolicyConditionMpr_t;

/**
 * Define the thermal policy violations return structure
 */
typedef struct
{
    long long timestamp;           //!< timestamp of the error
    unsigned int thermalViolation; //!< Temperature reached that violated policy
} dcgmPolicyConditionThermal_t;

/**
 * Define the power policy violations return structure
 */
typedef struct
{
    long long timestamp;         //!< timestamp of the error
    unsigned int powerViolation; //!< Power value reached that violated policy
} dcgmPolicyConditionPower_t;

/**
 * Define the nvlink policy violations return structure
 */
typedef struct
{
    long long timestamp;    //!< timestamp of the error
    unsigned short fieldId; //!< Nvlink counter field ID that violated policy
    unsigned int counter;   //!< Nvlink counter value that violated policy
} dcgmPolicyConditionNvlink_t;

/**
 * Define the xid policy violations return structure
 */
typedef struct
{
    long long timestamp; //!< Timestamp of the error
    unsigned int errnum; //!< The XID error number
} dcgmPolicyConditionXID_t;


/**
 * Define the structure that is given to the callback function
 */
typedef struct
{
    // version must always be first
    unsigned int version; //!< version number (dcgmPolicyCallbackResponse_version)

    dcgmPolicyCondition_t condition; //!< Condition that was violated
    union
    {
        dcgmPolicyConditionDbe_t dbe;         //!< ECC DBE return structure
        dcgmPolicyConditionPci_t pci;         //!< PCI replay error return structure
        dcgmPolicyConditionMpr_t mpr;         //!< Max retired pages limit return structure
        dcgmPolicyConditionThermal_t thermal; //!< Thermal policy violations return structure
        dcgmPolicyConditionPower_t power;     //!< Power policy violations return structure
        dcgmPolicyConditionNvlink_t nvlink;   //!< Nvlink policy violations return structure
        dcgmPolicyConditionXID_t xid;         //!< XID policy violations return structure
    } val;
} dcgmPolicyCallbackResponse_v1;


/**
 * Typedef for \ref dcgmPolicyCallbackResponse_v1
 */
typedef dcgmPolicyCallbackResponse_v1 dcgmPolicyCallbackResponse_t;

/**
 * Version 1 for \ref dcgmPolicyCallbackResponse_v1
 */
#define dcgmPolicyCallbackResponse_version1 MAKE_DCGM_VERSION(dcgmPolicyCallbackResponse_v1, 1)

/**
 * Latest version for \ref dcgmPolicyCallbackResponse_t
 */
#define dcgmPolicyCallbackResponse_version dcgmPolicyCallbackResponse_version1

/**
 * Set above size of largest blob entry. Currently this is dcgmDeviceVgpuTypeInfo_v1
 */
#define DCGM_MAX_BLOB_LENGTH 4096

/**
 * This structure is used to represent value for the field to be queried.
 */
typedef struct
{
    // version must always be first
    unsigned int version; //!< version number (dcgmFieldValue_version1)

    unsigned short fieldId;   //!< One of DCGM_FI_?
    unsigned short fieldType; //!< One of DCGM_FT_?
    int status;               //!< Status for the querying the field. DCGM_ST_OK or one of DCGM_ST_?
    int64_t ts;               //!< Timestamp in usec since 1970
    union
    {
        int64_t i64;                     //!< Int64 value
        double dbl;                      //!< Double value
        char str[DCGM_MAX_STR_LENGTH];   //!< NULL terminated string
        char blob[DCGM_MAX_BLOB_LENGTH]; //!< Binary blob
    } value;                             //!< Value
} dcgmFieldValue_v1;

/**
 * Version 1 for \ref dcgmFieldValue_v1
 */
#define dcgmFieldValue_version1 MAKE_DCGM_VERSION(dcgmFieldValue_v1, 1)

/**
 * This structure is used to represent value for the field to be queried.
 */
typedef struct
{
    // version must always be first
    unsigned int version;                    //!< version number (dcgmFieldValue_version2)
    dcgm_field_entity_group_t entityGroupId; //!< Entity group this field value's entity belongs to
    dcgm_field_eid_t entityId;               //!< Entity this field value belongs to
    unsigned short fieldId;                  //!< One of DCGM_FI_?
    unsigned short fieldType;                //!< One of DCGM_FT_?
    int status;                              //!< Status for the querying the field. DCGM_ST_OK or one of DCGM_ST_?
    unsigned int unused;                     //!< Unused for now to align ts to an 8-byte boundary.
    int64_t ts;                              //!< Timestamp in usec since 1970
    union
    {
        int64_t i64;                     //!< Int64 value
        double dbl;                      //!< Double value
        char str[DCGM_MAX_STR_LENGTH];   //!< NULL terminated string
        char blob[DCGM_MAX_BLOB_LENGTH]; //!< Binary blob
    } value;                             //!< Value
} dcgmFieldValue_v2;

/**
 * Version 2 for \ref dcgmFieldValue_v2
 */
#define dcgmFieldValue_version2 MAKE_DCGM_VERSION(dcgmFieldValue_v2, 2)

/**
 * Field value flags used by \ref dcgmEntitiesGetLatestValues
 *
 * Retrieve live data from the driver rather than cached data.
 * Warning: Setting this flag will result in multiple calls to the NVIDIA driver that will be much slower than
 *          retrieving a cached value.
 */
#define DCGM_FV_FLAG_LIVE_DATA 0x00000001

/**
 * User callback function for processing one or more field updates. This callback will
 * be invoked one or more times per field until all of the expected field values have been
 * enumerated. It is up to the callee to detect when the field id changes
 *
 * @param gpuId                IN: GPU ID of the GPU this field value set belongs to
 * @param values               IN: Field values. These values must be copied as they will be destroyed as soon as this
 *                                 call returns.
 * @param numValues            IN: Number of entries that are valid in values[]
 * @param userData             IN: User data pointer passed to the update function that generated this callback
 *
 * @returns
 *          0 if OK
 *         <0 if enumeration should stop. This allows to callee to abort field value enumeration.
 *
 */
typedef int (*dcgmFieldValueEnumeration_f)(unsigned int gpuId,
                                           dcgmFieldValue_v1 *values,
                                           int numValues,
                                           void *userData);

/**
 * User callback function for processing one or more field updates. This callback will
 * be invoked one or more times per field until all of the expected field values have been
 * enumerated. It is up to the callee to detect when the field id changes
 *
 * @param entityGroupId        IN: entityGroup of the entity this field value set belongs to
 * @param entityId             IN: Entity this field value set belongs to
 * @param values               IN: Field values. These values must be copied as they will be destroyed as soon as this
 *                                 call returns.
 * @param numValues            IN: Number of entries that are valid in values[]
 * @param userData             IN: User data pointer passed to the update function that generated this callback
 *
 * @returns
 *          0 if OK
 *         <0 if enumeration should stop. This allows to callee to abort field value enumeration.
 *
 */
typedef int (*dcgmFieldValueEntityEnumeration_f)(dcgm_field_entity_group_t entityGroupId,
                                                 dcgm_field_eid_t entityId,
                                                 dcgmFieldValue_v1 *values,
                                                 int numValues,
                                                 void *userData);


/**
 * Summary of time series data in int64 format.
 *
 * Each value will either be set or be a BLANK value.
 * Check for blank with the DCGM_INT64_IS_BLANK() macro.
 * \sa See dcgmvalue.h for the actual values of BLANK values
 */
typedef struct
{
    long long minValue; //!< Minimum value of the samples looked at
    long long maxValue; //!< Maximum value of the samples looked at
    long long average;  //!< Simple average of the samples looked at. Blank values are ignored for this calculation
} dcgmStatSummaryInt64_t;

/**
 * Same as dcgmStatSummaryInt64_t, but with 32-bit integer values
 */
typedef struct
{
    int minValue; //!< Minimum value of the samples looked at
    int maxValue; //!< Maximum value of the samples looked at
    int average;  //!< Simple average of the samples looked at. Blank values are ignored for this calculation
} dcgmStatSummaryInt32_t;

/**
 * Summary of time series data in double-precision format.
 * Each value will either be set or be a BLANK value.
 * Check for blank with the DCGM_FP64_IS_BLANK() macro.
 * \sa See dcgmvalue.h for the actual values of BLANK values
 */
typedef struct
{
    double minValue; //!< Minimum value of the samples looked at
    double maxValue; //!< Maximum value of the samples looked at
    double average;  //!< Simple average of the samples looked at. Blank values are ignored for this calculation
} dcgmStatSummaryFp64_t;

/**
 * Systems structure used to enable or disable health watch systems
 */
typedef enum dcgmHealthSystems_enum
{
    DCGM_HEALTH_WATCH_PCIE              = 0x1,   //!< PCIe system watches (must have 1m of data before query)
    DCGM_HEALTH_WATCH_NVLINK            = 0x2,   //!< NVLINK system watches
    DCGM_HEALTH_WATCH_PMU               = 0x4,   //!< Power management unit watches
    DCGM_HEALTH_WATCH_MCU               = 0x8,   //!< Micro-controller unit watches
    DCGM_HEALTH_WATCH_MEM               = 0x10,  //!< Memory watches
    DCGM_HEALTH_WATCH_SM                = 0x20,  //!< Streaming multiprocessor watches
    DCGM_HEALTH_WATCH_INFOROM           = 0x40,  //!< Inforom watches
    DCGM_HEALTH_WATCH_THERMAL           = 0x80,  //!< Temperature watches (must have 1m of data before query)
    DCGM_HEALTH_WATCH_POWER             = 0x100, //!< Power watches (must have 1m of data before query)
    DCGM_HEALTH_WATCH_DRIVER            = 0x200, //!< Driver-related watches
    DCGM_HEALTH_WATCH_NVSWITCH_NONFATAL = 0x400, //!< Non-fatal errors in NvSwitch
    DCGM_HEALTH_WATCH_NVSWITCH_FATAL    = 0x800, //!< Fatal errors in NvSwitch

    // ...
    DCGM_HEALTH_WATCH_ALL = 0xFFFFFFFF //!< All watches enabled
} dcgmHealthSystems_t;

#define DCGM_HEALTH_WATCH_COUNT_V1 10 /*!< For iterating through the dcgmHealthSystems_v1 enum */
#define DCGM_HEALTH_WATCH_COUNT_V2 12 /*!< For iterating through the dcgmHealthSystems_v2 enum */

/**
 * Health Watch test results
 */
typedef enum dcgmHealthWatchResult_enum
{
    DCGM_HEALTH_RESULT_PASS = 0,  //!< All results within this system are reporting normal
    DCGM_HEALTH_RESULT_WARN = 10, //!< A warning has been issued, refer to the response for more information
    DCGM_HEALTH_RESULT_FAIL = 20, //!< A failure has been issued, refer to the response for more information
} dcgmHealthWatchResults_t;

typedef struct
{
    char msg[1024];
    unsigned int code;
} dcgmDiagErrorDetail_t;

#define DCGM_ERR_MSG_LENGTH 512
/**
 * Error details
 *
 * Since DCGM 3.3
 */
typedef struct
{
    char msg[DCGM_ERR_MSG_LENGTH];
    int gpuId;
    unsigned int code;
    unsigned int category; //!< See dcgmErrorCategory_t
    unsigned int severity; //!< See dcgmErrorSeverity_t
} dcgmDiagErrorDetail_v2;

#define DCGM_HEALTH_WATCH_MAX_INCIDENTS DCGM_GROUP_MAX_ENTITIES

typedef struct
{
    dcgmHealthSystems_t system;       //!< system to which this information belongs
    dcgmHealthWatchResults_t health;  //!< health diagnosis of this incident
    dcgmDiagErrorDetail_t error;      //!< Information about the error(s) and their error codes
    dcgmGroupEntityPair_t entityInfo; //!< identify which entity has this error
} dcgmIncidentInfo_t;

/**
 * Health response structure version 4 - Simply list the incidents instead of reporting by entity
 *
 * Since DCGM 2.0
 */
typedef struct
{
    unsigned int version;                   //!< The version number of this struct
    dcgmHealthWatchResults_t overallHealth; //!< The overall health of this entire host
    unsigned int incidentCount;             //!< The number of health incidents reported in this struct
    dcgmIncidentInfo_t incidents[DCGM_HEALTH_WATCH_MAX_INCIDENTS]; //!< Report of the errors detected
} dcgmHealthResponse_v4;

/**
 * Version 4 for \ref dcgmHealthResponse_v4
 */
#define dcgmHealthResponse_version4 MAKE_DCGM_VERSION(dcgmHealthResponse_v4, 4)

/**
 * Latest version for \ref dcgmHealthResponse_t
 */
#define dcgmHealthResponse_version dcgmHealthResponse_version4

/**
 * Typedef for \ref dcgmHealthResponse_v4
 */
typedef dcgmHealthResponse_v4 dcgmHealthResponse_t;

/**
 * Structure used to set health watches via the dcgmHealthSet_v2 API
 */
typedef struct
{
    unsigned int version;        /*!< Version of this struct. Should be dcgmHealthSet_version2 */
    dcgmGpuGrp_t groupId;        /*!< Group ID representing collection of one or more entities. Look
                                      at \ref dcgmGroupCreate for details on creating the group.
                                      Alternatively, pass in the group id as \a DCGM_GROUP_ALL_GPUS
                                      to perform operation on all the GPUs or \a DCGM_GROUP_ALL_NVSWITCHES
                                      to perform operation on all the NvSwitches. */
    dcgmHealthSystems_t systems; /*!< An enum representing systems that should be enabled for health
                                      checks logically OR'd together. Refer to \ref dcgmHealthSystems_t
                                      for details. */
    long long updateInterval;    /*!< How often to query the underlying health information from the
                                      NVIDIA driver in usec. This should be the same as how often you call
                                      dcgmHealthCheck */
    double maxKeepAge;           /*!< How long to keep data cached for this field in seconds. This should
                                      be at least your maximum time between calling dcgmHealthCheck */
} dcgmHealthSetParams_v2;

/**
 * Version 2 for \ref dcgmHealthSet_v2
 */
#define dcgmHealthSetParams_version2 MAKE_DCGM_VERSION(dcgmHealthSetParams_v2, 2)


#define DCGM_MAX_PID_INFO_NUM 16
/**
 * per process utilization rates
 */
typedef struct
{
    unsigned int pid;
    double smUtil;
    double memUtil;
} dcgmProcessUtilInfo_t;

/**
 *Internal structure used to get the PID and the corresponding utilization rate
 */
typedef struct
{
    double util;
    unsigned int pid;
} dcgmProcessUtilSample_t;

/**
 * Info corresponding to single PID
 */
typedef struct
{
    unsigned int gpuId; //!< ID of the GPU this pertains to. GPU_ID_INVALID = summary information for multiple GPUs

    /* All of the following are during the process's lifetime */

    long long energyConsumed;               //!< Energy consumed by the gpu in milli-watt/seconds
    dcgmStatSummaryInt64_t pcieRxBandwidth; //!< PCI-E bytes read from the GPU
    dcgmStatSummaryInt64_t pcieTxBandwidth; //!< PCI-E bytes written to the GPU
    long long pcieReplays;                  //!< Count of PCI-E replays that occurred
    long long startTime;                    //!< Process start time in microseconds since 1970
    long long endTime; //!< Process end time in microseconds since 1970 or reported as 0 if the process is not completed
    dcgmProcessUtilInfo_t processUtilization; //!< Process SM and Memory Utilization (in percent)
    dcgmStatSummaryInt32_t smUtilization;     //!< GPU SM Utilization in percent
    dcgmStatSummaryInt32_t memoryUtilization; //!< GPU Memory Utilization in percent
    unsigned int eccSingleBit;                //!< Deprecated - Count of ECC single bit errors that occurred
    unsigned int eccDoubleBit;                //!< Count of ECC double bit errors that occurred
    dcgmStatSummaryInt32_t memoryClock;       //!< Memory clock in MHz
    dcgmStatSummaryInt32_t smClock;           //!< SM clock in MHz

    int numXidCriticalErrors;          //!< Number of valid entries in xidCriticalErrorsTs
    long long xidCriticalErrorsTs[10]; //!< Timestamps of the critical XID errors that occurred

    int numOtherComputePids;                              //!< Count of otherComputePids entries that are valid
    unsigned int otherComputePids[DCGM_MAX_PID_INFO_NUM]; //!< Other compute processes that ran. 0=no process

    int numOtherGraphicsPids;                              //!< Count of otherGraphicsPids entries that are valid
    unsigned int otherGraphicsPids[DCGM_MAX_PID_INFO_NUM]; //!< Other graphics processes that ran. 0=no process

    long long maxGpuMemoryUsed; //!< Maximum amount of GPU memory that was used in bytes

    long long powerViolationTime;       //!< Number of microseconds we were at reduced clocks due to power violation
    long long thermalViolationTime;     //!< Number of microseconds we were at reduced clocks due to thermal violation
    long long reliabilityViolationTime; //!< Amount of microseconds we were at reduced clocks
                                        //!< due to the reliability limit
    long long boardLimitViolationTime;  //!< Amount of microseconds we were at reduced clocks due to being at the
                                        //!< board's max voltage
    long long lowUtilizationTime;       //!< Amount of microseconds we were at reduced clocks due to low utilization
    long long syncBoostTime;            //!< Amount of microseconds we were at reduced clocks due to sync boost
    dcgmHealthWatchResults_t overallHealth; //!< The overall health of the system. \ref dcgmHealthWatchResults_t
    unsigned int incidentCount;
    struct
    {
        dcgmHealthSystems_t system;      //!< system to which this information belongs
        dcgmHealthWatchResults_t health; //!< health of the specified system on this GPU
    } systems[DCGM_HEALTH_WATCH_COUNT_V1];
} dcgmPidSingleInfo_t;

/**
 * To store process statistics
 */
typedef struct
{
    unsigned int version; //!< Version of this message  (dcgmPidInfo_version)
    unsigned int pid;     //!< PID of the process
    unsigned int unused;
    int numGpus;                                    //!< Number of GPUs that are valid in GPUs
    dcgmPidSingleInfo_t summary;                    //!< Summary information for all GPUs listed in gpus[]
    dcgmPidSingleInfo_t gpus[DCGM_MAX_NUM_DEVICES]; //!< Per-GPU information for this PID
} dcgmPidInfo_v2;

/**
 * Typedef for \ref dcgmPidInfo_v2
 */
typedef dcgmPidInfo_v2 dcgmPidInfo_t;

/**
 * Version 2 for \ref dcgmPidInfo_v2
 */
#define dcgmPidInfo_version2 MAKE_DCGM_VERSION(dcgmPidInfo_v2, 2)

/**
 * Latest version for \ref dcgmPidInfo_t
 */
#define dcgmPidInfo_version dcgmPidInfo_version2

/**
 * Info corresponding to the job on a GPU
 */
typedef struct
{
    unsigned int gpuId; //!< ID of the GPU this pertains to. GPU_ID_INVALID = summary information for multiple GPUs

    /* All of the following are during the job's lifetime */

    long long energyConsumed;                 //!< Energy consumed in milli-watt/seconds
    dcgmStatSummaryFp64_t powerUsage;         //!< Power usage Min/Max/Avg in watts
    dcgmStatSummaryInt64_t pcieRxBandwidth;   //!< PCI-E bytes read from the GPU
    dcgmStatSummaryInt64_t pcieTxBandwidth;   //!< PCI-E bytes written to the GPU
    long long pcieReplays;                    //!< Count of PCI-E replays that occurred
    long long startTime;                      //!< User provided job start time in microseconds since 1970
    long long endTime;                        //!< User provided job end time in microseconds since 1970
    dcgmStatSummaryInt32_t smUtilization;     //!< GPU SM Utilization in percent
    dcgmStatSummaryInt32_t memoryUtilization; //!< GPU Memory Utilization in percent
    unsigned int eccSingleBit;                //!< Deprecated - Count of ECC single bit errors that occurred
    unsigned int eccDoubleBit;                //!< Count of ECC double bit errors that occurred
    dcgmStatSummaryInt32_t memoryClock;       //!< Memory clock in MHz
    dcgmStatSummaryInt32_t smClock;           //!< SM clock in MHz

    int numXidCriticalErrors;          //!< Number of valid entries in xidCriticalErrorsTs
    long long xidCriticalErrorsTs[10]; //!< Timestamps of the critical XID errors that occurred

    int numComputePids;                                          //!< Count of computePids entries that are valid
    dcgmProcessUtilInfo_t computePidInfo[DCGM_MAX_PID_INFO_NUM]; //!< List of compute processes that ran during the job
                                                                 //!< 0=no process

    int numGraphicsPids;                                          //!< Count of graphicsPids entries that are valid
    dcgmProcessUtilInfo_t graphicsPidInfo[DCGM_MAX_PID_INFO_NUM]; //!< List of compute processes that ran during the job
                                                                  //!< 0=no process

    long long maxGpuMemoryUsed; //!< Maximum amount of GPU memory that was used in bytes

    long long powerViolationTime;       //!< Number of microseconds we were at reduced clocks due to power violation
    long long thermalViolationTime;     //!< Number of microseconds we were at reduced clocks due to thermal violation
    long long reliabilityViolationTime; //!< Amount of microseconds we were at reduced clocks
                                        //!< due to the reliability limit
    long long boardLimitViolationTime;  //!< Amount of microseconds we were at reduced clocks
                                        //!< due to being at the board's max voltage
    long long lowUtilizationTime;       //!< Amount of microseconds we were at reduced clocks due to low utilization
    long long syncBoostTime;            //!< Amount of microseconds we were at reduced clocks due to sync boost
    dcgmHealthWatchResults_t overallHealth; //!< The overall health of the system. \ref dcgmHealthWatchResults_t
    unsigned int incidentCount;
    struct
    {
        dcgmHealthSystems_t system;      //!< system to which this information belongs
        dcgmHealthWatchResults_t health; //!< health of the specified system on this GPU
    } systems[DCGM_HEALTH_WATCH_COUNT_V1];
} dcgmGpuUsageInfo_t;


/**
 * To store job statistics
 * The following fields are not applicable in the summary info:
 * - pcieRxBandwidth (Min/Max)
 * - pcieTxBandwidth (Min/Max)
 * - smUtilization (Min/Max)
 * - memoryUtilization (Min/Max)
 * - memoryClock (Min/Max)
 * - smClock (Min/Max)
 * - processSamples
 *
 * The average value in the above fields (in the summary) is the
 * average of the averages of respective fields from all GPUs
 */
typedef struct
{
    unsigned int version;                          //!< Version of this message  (dcgmPidInfo_version)
    int numGpus;                                   //!< Number of GPUs that are valid in gpus[]
    dcgmGpuUsageInfo_t summary;                    //!< Summary information for all GPUs listed in gpus[]
    dcgmGpuUsageInfo_t gpus[DCGM_MAX_NUM_DEVICES]; //!< Per-GPU information for this PID
} dcgmJobInfo_v3;

/**
 * Typedef for \ref dcgmJobInfo_v3
 */
typedef dcgmJobInfo_v3 dcgmJobInfo_t;

/**
 * Version 3 for \ref dcgmJobInfo_v3
 */
#define dcgmJobInfo_version3 MAKE_DCGM_VERSION(dcgmJobInfo_v3, 3)

/**
 * Latest version for \ref dcgmJobInfo_t
 */
#define dcgmJobInfo_version dcgmJobInfo_version3


/**
 * Running process information for a compute or graphics process
 */
typedef struct
{
    unsigned int version;          //!< Version of this message (dcgmRunningProcess_version)
    unsigned int pid;              //!< PID of the process
    unsigned long long memoryUsed; //!< GPU memory used by this process in bytes.
} dcgmRunningProcess_v1;

/**
 * Typedef for \ref dcgmRunningProcess_v1
 */
typedef dcgmRunningProcess_v1 dcgmRunningProcess_t;

/**
 * Version 1 for \ref dcgmRunningProcess_v1
 */
#define dcgmRunningProcess_version1 MAKE_DCGM_VERSION(dcgmRunningProcess_v1, 1)

/**
 * Latest version for \ref dcgmRunningProcess_t
 */
#define dcgmRunningProcess_version dcgmRunningProcess_version1

/**
 * Enumeration for diagnostic levels
 */
typedef enum
{
    DCGM_DIAG_LVL_INVALID = 0,  //!< Uninitialized
    DCGM_DIAG_LVL_SHORT   = 10, //!< run a very basic health check on the system
    DCGM_DIAG_LVL_MED     = 20, //!< run a medium-length diagnostic (a few minutes)
    DCGM_DIAG_LVL_LONG    = 30, //!< run a extensive diagnostic (several minutes)
    DCGM_DIAG_LVL_XLONG   = 40, //!< run a very extensive diagnostic (many minutes)
} dcgmDiagnosticLevel_t;

/**
 * Diagnostic test results
 */
typedef enum dcgmDiagResult_enum
{
    DCGM_DIAG_RESULT_PASS    = 0, //!< This test passed as diagnostics
    DCGM_DIAG_RESULT_SKIP    = 1, //!< This test was skipped
    DCGM_DIAG_RESULT_WARN    = 2, //!< This test passed with warnings
    DCGM_DIAG_RESULT_FAIL    = 3, //!< This test failed the diagnostics
    DCGM_DIAG_RESULT_NOT_RUN = 4, //!< This test wasn't executed
} dcgmDiagResult_t;

typedef struct
{
    dcgmDiagResult_t status;     //!< The result of the test
    dcgmDiagErrorDetail_t error; //!< The error message and error code, if any
    char info[1024];             //!< Information details returned from the test, if any
} dcgmDiagTestResult_v2;

#define DCGM_MAX_ERRORS 5
typedef struct
{
    dcgmDiagResult_t status;                       //!< The result of the test
    dcgmDiagErrorDetail_v2 error[DCGM_MAX_ERRORS]; //!< The error message and error code, if any
    char info[DCGM_ERR_MSG_LENGTH];                //!< Information details returned from the test, if any
} dcgmDiagTestResult_v3;

/**
 * Diagnostic per gpu tests - fixed indices for dcgmDiagResponsePerGpu_t.results[]
 */
typedef enum dcgmPerGpuTestIndices_enum
{
    DCGM_MEMORY_INDEX           = 0, //!< Memory test index
    DCGM_DIAGNOSTIC_INDEX       = 1, //!< Diagnostic test index
    DCGM_PCI_INDEX              = 2, //!< PCIe test index
    DCGM_SM_STRESS_INDEX        = 3, //!< SM Stress test index
    DCGM_TARGETED_STRESS_INDEX  = 4, //!< Targeted Stress test index
    DCGM_TARGETED_POWER_INDEX   = 5, //!< Targeted Power test index
    DCGM_MEMORY_BANDWIDTH_INDEX = 6, //!< Memory bandwidth test index
    DCGM_MEMTEST_INDEX          = 7, //!< Memtest test index
    DCGM_PULSE_TEST_INDEX       = 8, //!< Pulse test index
    DCGM_EUD_TEST_INDEX         = 9, //!< EUD test index
    // Remaining tests are included for convenience but have different execution rules
    // See DCGM_PER_GPU_TEST_COUNT
    DCGM_UNUSED2_TEST_INDEX   = 10,
    DCGM_UNUSED3_TEST_INDEX   = 11,
    DCGM_UNUSED4_TEST_INDEX   = 12,
    DCGM_UNUSED5_TEST_INDEX   = 13,
    DCGM_SOFTWARE_INDEX       = 14, //!< Software test index
    DCGM_CONTEXT_CREATE_INDEX = 15, //!< Context create test index
    DCGM_UNKNOWN_INDEX        = 16  //!< Unknown test
} dcgmPerGpuTestIndices_t;

// TODO: transition these to dcgm_deprecated.h
#define DCGM_SM_PERF_INDEX       DCGM_SM_STRESS_INDEX
#define DCGM_TARGETED_PERF_INDEX DCGM_TARGETED_PERF_INDEX

// Number of diag tests
// NOTE: does not include software and context_create which have different execution rules
#define DCGM_PER_GPU_TEST_COUNT_V8 13
#define DCGM_PER_GPU_TEST_COUNT_V7 9

/**
 * Per GPU diagnostics result structure
 */
typedef struct
{
    unsigned int gpuId;                                        //!< ID for the GPU this information pertains
    unsigned int hwDiagnosticReturn;                           //!< Per GPU hardware diagnostic test return code
    dcgmDiagTestResult_v2 results[DCGM_PER_GPU_TEST_COUNT_V8]; //!< Array with a result for each per-gpu test
} dcgmDiagResponsePerGpu_v4;

typedef struct
{
    unsigned int gpuId;                                        //!< ID for the GPU this information pertains
    unsigned int hwDiagnosticReturn;                           //!< Per GPU hardware diagnostic test return code
    dcgmDiagTestResult_v3 results[DCGM_PER_GPU_TEST_COUNT_V8]; //!< Array with a result for each per-gpu test
} dcgmDiagResponsePerGpu_v5;

/**
 * Per gpu response structure v3
 *
 * Since DCGM 2.4
 */
typedef struct
{
    unsigned int gpuId;                                        //!< ID for the GPU this information pertains
    unsigned int hwDiagnosticReturn;                           //!< Per GPU hardware diagnostic test return code
    dcgmDiagTestResult_v2 results[DCGM_PER_GPU_TEST_COUNT_V7]; //!< Array with a result for each per-gpu test
} dcgmDiagResponsePerGpu_v3;


#define DCGM_SWTEST_COUNT     10
#define LEVEL_ONE_MAX_RESULTS 16

typedef enum dcgmSoftwareTest_enum
{
    DCGM_SWTEST_DENYLIST             = 0, //!< test for presence of drivers on the denylist (e.g. nouveau)
    DCGM_SWTEST_NVML_LIBRARY         = 1, //!< test for presence (and version) of NVML lib
    DCGM_SWTEST_CUDA_MAIN_LIBRARY    = 2, //!< test for presence (and version) of CUDA lib
    DCGM_SWTEST_CUDA_RUNTIME_LIBRARY = 3, //!< test for presence (and version) of CUDA RT lib
    DCGM_SWTEST_PERMISSIONS          = 4, //!< test for character device permissions
    DCGM_SWTEST_PERSISTENCE_MODE     = 5, //!< test for persistence mode enabled
    DCGM_SWTEST_ENVIRONMENT          = 6, //!< test for CUDA environment vars that may slow tests
    DCGM_SWTEST_PAGE_RETIREMENT      = 7, //!< test for pending frame buffer page retirement
    DCGM_SWTEST_GRAPHICS_PROCESSES   = 8, //!< test for graphics processes running
    DCGM_SWTEST_INFOROM              = 9, //!< test for inforom corruption
} dcgmSoftwareTest_t;

#define DCGM_DEVICE_ID_LEN 5
#define DCGM_VERSION_LEN   12

/**
 * Global diagnostics result structure v9
 *
 * Since DCGM 3.3
 */
typedef struct
{
    unsigned int version;           //!< version number (dcgmDiagResult_version)
    unsigned int gpuCount;          //!< number of valid per GPU results
    unsigned int levelOneTestCount; //!< number of valid levelOne results

    dcgmDiagTestResult_v3 levelOneResults[LEVEL_ONE_MAX_RESULTS];    //!< Basic, system-wide test results.
    dcgmDiagResponsePerGpu_v5 perGpuResponses[DCGM_MAX_NUM_DEVICES]; //!< per GPU test results
    dcgmDiagErrorDetail_v2 systemError;                              //!< System-wide error reported from NVVS
    char devIds[DCGM_MAX_NUM_DEVICES][DCGM_DEVICE_ID_LEN];           //!< The SKU device id for each GPU
    char devSerials[DCGM_MAX_NUM_DEVICES][DCGM_MAX_STR_LENGTH];      //!< Serial for the device
    char dcgmVersion[DCGM_VERSION_LEN];                              //!< A string representing DCGM's version
    char driverVersion[DCGM_MAX_STR_LENGTH];                         //!< A string representing the driver version
    char _unused[596];                                               //!< No longer used
} dcgmDiagResponse_v9;

/**
 * Global diagnostics result structure v8
 *
 * Since DCGM 3.0
 */
typedef struct
{
    unsigned int version;           //!< version number (dcgmDiagResult_version)
    unsigned int gpuCount;          //!< number of valid per GPU results
    unsigned int levelOneTestCount; //!< number of valid levelOne results

    dcgmDiagTestResult_v2 levelOneResults[LEVEL_ONE_MAX_RESULTS];    //!< Basic, system-wide test results.
    dcgmDiagResponsePerGpu_v4 perGpuResponses[DCGM_MAX_NUM_DEVICES]; //!< per GPU test results
    dcgmDiagErrorDetail_t systemError;                               //!< System-wide error reported from NVVS
    char devIds[DCGM_MAX_NUM_DEVICES][DCGM_DEVICE_ID_LEN];           //!< The SKU device id for each GPU
    char dcgmVersion[DCGM_VERSION_LEN];                              //!< A string representing DCGM's version
    char driverVersion[DCGM_MAX_STR_LENGTH];                         //!< A string representing the driver version
    char _unused[596];                                               //!< No longer used
} dcgmDiagResponse_v8;

/**
 * Global diagnostics result structure v7
 *
 * Since DCGM 2.4
 */
typedef struct
{
    unsigned int version;           //!< version number (dcgmDiagResult_version)
    unsigned int gpuCount;          //!< number of valid per GPU results
    unsigned int levelOneTestCount; //!< number of valid levelOne results

    dcgmDiagTestResult_v2 levelOneResults[LEVEL_ONE_MAX_RESULTS];    //!< Basic, system-wide test results.
    dcgmDiagResponsePerGpu_v3 perGpuResponses[DCGM_MAX_NUM_DEVICES]; //!< per GPU test results
    dcgmDiagErrorDetail_t systemError;                               //!< System-wide error reported from NVVS
    char _unused[1024];                                              //!< No longer used
} dcgmDiagResponse_v7;

/**
 * Typedef for \ref dcgmDiagResponse_v9
 */
typedef dcgmDiagResponse_v9 dcgmDiagResponse_t;

/**
 * Version 9 for \ref dcgmDiagResponse_v9
 */
#define dcgmDiagResponse_version9 MAKE_DCGM_VERSION(dcgmDiagResponse_v9, 9)

/**
 * Version 8 for \ref dcgmDiagResponse_v8
 */
#define dcgmDiagResponse_version8 MAKE_DCGM_VERSION(dcgmDiagResponse_v8, 8)

/**
 * Version 7 for \ref dcgmDiagResponse_v7
 */
#define dcgmDiagResponse_version7 MAKE_DCGM_VERSION(dcgmDiagResponse_v7, 7)

/**
 * Latest version for \ref dcgmDiagResponse_t
 */
#define dcgmDiagResponse_version dcgmDiagResponse_version9

/**
 * Represents level relationships within a system between two GPUs
 * The enums are spaced to allow for future relationships.
 * These match the definitions in nvml.h
 */
typedef enum dcgmGpuLevel_enum
{
    DCGM_TOPOLOGY_UNINITIALIZED = 0x0,

    /** \name PCI connectivity states */
    /**@{*/
    DCGM_TOPOLOGY_BOARD      = 0x1, //!< multi-GPU board
    DCGM_TOPOLOGY_SINGLE     = 0x2, //!< all devices that only need traverse a single PCIe switch
    DCGM_TOPOLOGY_MULTIPLE   = 0x4, //!< all devices that need not traverse a host bridge
    DCGM_TOPOLOGY_HOSTBRIDGE = 0x8, //!< all devices that are connected to the same host bridge
    DCGM_TOPOLOGY_CPU    = 0x10, //!< all devices that are connected to the same CPU but possibly multiple host bridges
    DCGM_TOPOLOGY_SYSTEM = 0x20, //!< all devices in the system
    /**@}*/

    /** \name NVLINK connectivity states */
    /**@{*/
    DCGM_TOPOLOGY_NVLINK1  = 0x0100,    //!< GPUs connected via a single NVLINK link
    DCGM_TOPOLOGY_NVLINK2  = 0x0200,    //!< GPUs connected via two NVLINK links
    DCGM_TOPOLOGY_NVLINK3  = 0x0400,    //!< GPUs connected via three NVLINK links
    DCGM_TOPOLOGY_NVLINK4  = 0x0800,    //!< GPUs connected via four NVLINK links
    DCGM_TOPOLOGY_NVLINK5  = 0x1000,    //!< GPUs connected via five NVLINK links
    DCGM_TOPOLOGY_NVLINK6  = 0x2000,    //!< GPUs connected via six NVLINK links
    DCGM_TOPOLOGY_NVLINK7  = 0x4000,    //!< GPUs connected via seven NVLINK links
    DCGM_TOPOLOGY_NVLINK8  = 0x8000,    //!< GPUs connected via eight NVLINK links
    DCGM_TOPOLOGY_NVLINK9  = 0x10000,   //!< GPUs connected via nine NVLINK links
    DCGM_TOPOLOGY_NVLINK10 = 0x20000,   //!< GPUs connected via ten NVLINK links
    DCGM_TOPOLOGY_NVLINK11 = 0x40000,   //!< GPUs connected via eleven NVLINK links
    DCGM_TOPOLOGY_NVLINK12 = 0x80000,   //!< GPUs connected via twelve NVLINK links
    DCGM_TOPOLOGY_NVLINK13 = 0x100000,  //!< GPUs connected via twelve NVLINK links
    DCGM_TOPOLOGY_NVLINK14 = 0x200000,  //!< GPUs connected via twelve NVLINK links
    DCGM_TOPOLOGY_NVLINK15 = 0x400000,  //!< GPUs connected via twelve NVLINK links
    DCGM_TOPOLOGY_NVLINK16 = 0x800000,  //!< GPUs connected via twelve NVLINK links
    DCGM_TOPOLOGY_NVLINK17 = 0x1000000, //!< GPUs connected via twelve NVLINK links
    DCGM_TOPOLOGY_NVLINK18 = 0x2000000, //!< GPUs connected via twelve NVLINK links
    /**@}*/
} dcgmGpuTopologyLevel_t;

// the PCI paths are the lower 8 bits of the path information
#define DCGM_TOPOLOGY_PATH_PCI(x) (dcgmGpuTopologyLevel_t)((unsigned int)(x)&0xFF)

// the NVLINK paths are the upper 24 bits of the path information
#define DCGM_TOPOLOGY_PATH_NVLINK(x) (dcgmGpuTopologyLevel_t)((unsigned int)(x)&0xFFFFFF00)

#define DCGM_AFFINITY_BITMASK_ARRAY_SIZE 8

/**
 * Device topology information
 */
typedef struct
{
    unsigned int version; //!< version number (dcgmDeviceTopology_version)

    unsigned long cpuAffinityMask[DCGM_AFFINITY_BITMASK_ARRAY_SIZE]; //!< affinity mask for the specified GPU
                                                                     //!< a 1 represents affinity to the CPU in that
                                                                     //!< bit position supports up to 256 cores
    unsigned int numGpus;                                            //!< number of valid entries in gpuPaths

    struct
    {
        unsigned int gpuId;          //!< gpuId to which the path represents
        dcgmGpuTopologyLevel_t path; //!< path to the gpuId from this GPU. Note that this is a bit-mask
                                     //!< of DCGM_TOPOLOGY_* values and can contain both PCIe topology
                                     //!< and NvLink topology where applicable. For instance:
                                     //!< 0x210 = DCGM_TOPOLOGY_CPU | DCGM_TOPOLOGY_NVLINK2
                                     //!< Use the macros DCGM_TOPOLOGY_PATH_NVLINK and
                                     //!< DCGM_TOPOLOGY_PATH_PCI to mask the NvLink and PCI paths, respectively.
        unsigned int localNvLinkIds; //!< bits representing the local links connected to gpuId
                                     //!< e.g. if this field == 3, links 0 and 1 are connected,
                                     //!< field is only valid if NVLINKS actually exist between GPUs
    } gpuPaths[DCGM_MAX_NUM_DEVICES - 1];
} dcgmDeviceTopology_v1;

/**
 * Typedef for \ref dcgmDeviceTopology_v1
 */
typedef dcgmDeviceTopology_v1 dcgmDeviceTopology_t;

/**
 * Version 1 for \ref dcgmDeviceTopology_v1
 */
#define dcgmDeviceTopology_version1 MAKE_DCGM_VERSION(dcgmDeviceTopology_v1, 1)

/**
 * Latest version for \ref dcgmDeviceTopology_t
 */
#define dcgmDeviceTopology_version dcgmDeviceTopology_version1

/**
 * Group topology information
 */
typedef struct
{
    unsigned int version; //!< version number (dcgmGroupTopology_version)

    unsigned long
        groupCpuAffinityMask[DCGM_AFFINITY_BITMASK_ARRAY_SIZE]; //!< the CPU affinity mask for all GPUs in the group
                                                                //!< a 1 represents affinity to the CPU in that bit
                                                                //!< position supports up to 256 cores
    unsigned int numaOptimalFlag;                               //!< a zero value indicates that 1 or more GPUs
                                                                //!< in the group have a different CPU affinity and thus
                                                                //!< may not be optimal for certain algorithms
    dcgmGpuTopologyLevel_t slowestPath;                         //!< the slowest path amongst GPUs in the group
} dcgmGroupTopology_v1;

/**
 * Typedef for \ref dcgmGroupTopology_v1
 */
typedef dcgmGroupTopology_v1 dcgmGroupTopology_t;

/**
 * Version 1 for \ref dcgmGroupTopology_v1
 */
#define dcgmGroupTopology_version1 MAKE_DCGM_VERSION(dcgmGroupTopology_v1, 1)

/**
 * Latest version for \ref dcgmGroupTopology_t
 */
#define dcgmGroupTopology_version dcgmGroupTopology_version1

/**
 * DCGM Memory usage information
 */
typedef struct
{
    unsigned int version; //!< version number (dcgmIntrospectMemory_version)
    long long bytesUsed;  //!< number of bytes
} dcgmIntrospectMemory_v1;

/**
 * Typedef for \ref dcgmIntrospectMemory_t
 */
typedef dcgmIntrospectMemory_v1 dcgmIntrospectMemory_t;

/**
 * Version 1 for \ref dcgmIntrospectMemory_t
 */
#define dcgmIntrospectMemory_version1 MAKE_DCGM_VERSION(dcgmIntrospectMemory_v1, 1)

/**
 * Latest version for \ref dcgmIntrospectMemory_t
 */
#define dcgmIntrospectMemory_version dcgmIntrospectMemory_version1

/**
 * DCGM CPU Utilization information.  Multiply values by 100 to get them in %.
 */
typedef struct
{
    unsigned int version; //!< version number (dcgmMetadataCpuUtil_version)
    double total;         //!< fraction of device's CPU resources that were used
    double kernel;        //!< fraction of device's CPU resources that were used in kernel mode
    double user;          //!< fraction of device's CPU resources that were used in user mode
} dcgmIntrospectCpuUtil_v1;

/**
 * Typedef for \ref dcgmIntrospectCpuUtil_t
 */
typedef dcgmIntrospectCpuUtil_v1 dcgmIntrospectCpuUtil_t;

/**
 * Version 1 for \ref dcgmIntrospectCpuUtil_t
 */
#define dcgmIntrospectCpuUtil_version1 MAKE_DCGM_VERSION(dcgmIntrospectCpuUtil_v1, 1)

/**
 * Latest version for \ref dcgmIntrospectCpuUtil_t
 */
#define dcgmIntrospectCpuUtil_version dcgmIntrospectCpuUtil_version1

#define DCGM_MAX_CONFIG_FILE_LEN 10000
#define DCGM_MAX_TEST_NAMES      20
#define DCGM_MAX_TEST_NAMES_LEN  50
#define DCGM_MAX_TEST_PARMS      100
#define DCGM_MAX_TEST_PARMS_LEN  100
#define DCGM_GPU_LIST_LEN        50
#define DCGM_FILE_LEN            30
#define DCGM_PATH_LEN            128
#define DCGM_THROTTLE_MASK_LEN   50

/**
 * Flags options for running the GPU diagnostic
 * @{
 *
 */

#define DCGM_HOME_DIR_VAR_NAME "DCGM_HOME_DIR"

/**
 * Output in verbose mode; include information as well as warnings
 */
#define DCGM_RUN_FLAGS_VERBOSE 0x0001

/**
 * Output stats only on failure
 */
#define DCGM_RUN_FLAGS_STATSONFAIL 0x0002

/**
 * UNUSED Train DCGM diagnostic and output a configuration file with golden values
 */
#define DCGM_RUN_FLAGS_TRAIN 0x0004

/**
 * UNUSED Ignore warnings against training the diagnostic and train anyway
 */
#define DCGM_RUN_FLAGS_FORCE_TRAIN 0x0008

/**
 * Enable fail early checks for the Targeted Stress, Targeted Power, SM Stress, and Diagnostic tests
 */
#define DCGM_RUN_FLAGS_FAIL_EARLY 0x0010

/**
 * @}
 */

/*
 * Run diagnostic structure v7
 */
typedef struct
{
    unsigned int version;            //!< version of this message
    unsigned int flags;              //!< flags specifying binary options for running it. See DCGM_RUN_FLAGS_*
    unsigned int debugLevel;         //!< 0-5 for the debug level the GPU diagnostic will use for logging.
    dcgmGpuGrp_t groupId;            //!< group of GPUs to verify. Cannot be specified together with gpuList.
    dcgmPolicyValidation_t validate; //!< 0-3 for which tests to run. Optional.
    char testNames[DCGM_MAX_TEST_NAMES][DCGM_MAX_TEST_NAMES_LEN]; //!< Specified list of test names. Optional.
    char testParms[DCGM_MAX_TEST_PARMS][DCGM_MAX_TEST_PARMS_LEN]; //!< Parameters to set for specified tests
                                                                  //!< in the format:
                                                                  //!< testName.parameterName=parameterValue. Optional.
    char fakeGpuList[DCGM_GPU_LIST_LEN]; //!< Comma-separated list of GPUs. Cannot be specified with the groupId.
    char gpuList[DCGM_GPU_LIST_LEN];     //!< Comma-separated list of GPUs. Cannot be specified with the groupId.
    char debugLogFile[DCGM_PATH_LEN];    //!< Alternate name for the debug log file that should be used
    char statsPath[DCGM_PATH_LEN];       //!< Path that the plugin's statistics files should be written to
    char configFileContents[DCGM_MAX_CONFIG_FILE_LEN]; //!< Contents of nvvs config file (likely yaml)
    char throttleMask[DCGM_THROTTLE_MASK_LEN]; //!< Throttle reasons to ignore as either integer mask or csv list of
                                               //!< reasons
    char pluginPath[DCGM_PATH_LEN]; //!< Custom path to the diagnostic plugins - No longer supported as of 2.2.9

    unsigned int currentIteration;  //!< The current iteration that will be executed
    unsigned int totalIterations;   //!< The total iterations that will be executed
    unsigned int _unusedInt1;       //!< No longer used
    char _unusedBuf[DCGM_PATH_LEN]; //!< No longer used
    unsigned int failCheckInterval; //!< How often the fail early checks should occur when enabled.
} dcgmRunDiag_v7;

/**
 * Version 7 for \ref dcgmRunDiag_t
 */
#define dcgmRunDiag_version7 MAKE_DCGM_VERSION(dcgmRunDiag_v7, 7)

/**
 * Flags for dcgmGetEntityGroupEntities's flags parameter
 *
 * Only return entities that are supported by DCGM.
 * This mimics the behavior of dcgmGetAllSupportedDevices().
 */
#define DCGM_GEGE_FLAG_ONLY_SUPPORTED 0x00000001

/**
 * Identifies a GPU NVLink error type returned by DCGM_FI_DEV_GPU_NVLINK_ERRORS
 */
typedef enum dcgmGpuNVLinkErrorType_enum
{
    DCGM_GPU_NVLINK_ERROR_RECOVERY_REQUIRED = 1, //!< NVLink link recovery error occurred
    DCGM_GPU_NVLINK_ERROR_FATAL,                 //!< NVLink link fatal error occurred
} dcgmGpuNVLinkErrorType_t;

/** Topology hints for dcgmSelectGpusByTopology()
 * @{
 */

/** No hints specified */
#define DCGM_TOPO_HINT_F_NONE 0x00000000

/** Ignore the health of the GPUs when picking GPUs for job
 * execution. By default, only healthy GPUs are considered.
 */
#define DCGM_TOPO_HINT_F_IGNOREHEALTH 0x00000001

/**
 * @}
 */


typedef struct
{
    unsigned int version; //!< version of this message
    uint64_t inputGpuIds; //!< bit-mask of the GPU ids to choose from
    uint32_t numGpus;     //!< the number of GPUs that DCGM should choose
    uint64_t hintFlags;   //!< Hints to ignore certain factors for the scheduling hint
} dcgmTopoSchedHint_v1;

typedef dcgmTopoSchedHint_v1 dcgmTopoSchedHint_t;

#define dcgmTopoSchedHint_version1 MAKE_DCGM_VERSION(dcgmTopoSchedHint_v1, 1)

/**
 * NvLink link states
 */
typedef enum dcgmNvLinkLinkState_enum
{
    DcgmNvLinkLinkStateNotSupported = 0, //!< NvLink is unsupported by this GPU (Default for GPUs)
    DcgmNvLinkLinkStateDisabled     = 1, //!< NvLink is supported for this link but this link is disabled
                                         //!< (Default for NvSwitches)
    DcgmNvLinkLinkStateDown = 2,         //!< This NvLink link is down (inactive)
    DcgmNvLinkLinkStateUp   = 3          //!< This NvLink link is up (active)
} dcgmNvLinkLinkState_t;

/**
 * State of NvLink links for a GPU
 */
typedef struct
{
    dcgm_field_eid_t entityId;                                              //!< Entity ID of the GPU (gpuId)
    dcgmNvLinkLinkState_t linkState[DCGM_NVLINK_MAX_LINKS_PER_GPU_LEGACY1]; //!< Per-GPU link states
} dcgmNvLinkGpuLinkStatus_v1;

typedef struct
{
    dcgm_field_eid_t entityId;                                              //!< Entity ID of the GPU (gpuId)
    dcgmNvLinkLinkState_t linkState[DCGM_NVLINK_MAX_LINKS_PER_GPU_LEGACY2]; //!< Per-GPU link states
} dcgmNvLinkGpuLinkStatus_v2;


typedef struct
{
    dcgm_field_eid_t entityId;                                      //!< Entity ID of the GPU (gpuId)
    dcgmNvLinkLinkState_t linkState[DCGM_NVLINK_MAX_LINKS_PER_GPU]; //!< Per-GPU link states
} dcgmNvLinkGpuLinkStatus_v3;

/**
 * State of NvLink links for a NvSwitch
 */
typedef struct
{
    dcgm_field_eid_t entityId;                                           //!< Entity ID of the NvSwitch (physicalId)
    dcgmNvLinkLinkState_t linkState[DCGM_NVLINK_MAX_LINKS_PER_NVSWITCH]; //!< Per-NvSwitch link states
} dcgmNvLinkNvSwitchLinkStatus_t;

/**
 * Status of all of the NvLinks in a given system
 */
typedef struct
{
    unsigned int version; //!< Version of this request. Should be dcgmNvLinkStatus_version1
    unsigned int numGpus; //!< Number of entries in gpus[] that are populated
    dcgmNvLinkGpuLinkStatus_v3 gpus[DCGM_MAX_NUM_DEVICES]; //!< Per-GPU NvLink link statuses
    unsigned int numNvSwitches;                            //!< Number of entries in nvSwitches[] that are populated
    dcgmNvLinkNvSwitchLinkStatus_t nvSwitches[DCGM_MAX_NUM_SWITCHES]; //!< Per-NvSwitch link statuses
} dcgmNvLinkStatus_v3;

typedef dcgmNvLinkStatus_v3 dcgmNvLinkStatus_t;

/**
 * Version 3 of dcgmNvLinkStatus
 */
#define dcgmNvLinkStatus_version3 MAKE_DCGM_VERSION(dcgmNvLinkStatus_v3, 3)

/* Bitmask values for dcgmGetFieldIdSummary - Sync with DcgmcmSummaryType_t */
#define DCGM_SUMMARY_MIN      0x00000001
#define DCGM_SUMMARY_MAX      0x00000002
#define DCGM_SUMMARY_AVG      0x00000004
#define DCGM_SUMMARY_SUM      0x00000008
#define DCGM_SUMMARY_COUNT    0x00000010
#define DCGM_SUMMARY_INTEGRAL 0x00000020
#define DCGM_SUMMARY_DIFF     0x00000040
#define DCGM_SUMMARY_SIZE     7

/* dcgmSummaryResponse_t is part of dcgmFieldSummaryRequest, so it uses dcgmFieldSummaryRequest's version. */

typedef struct
{
    unsigned int fieldType;    //!< type of field that is summarized (int64 or fp64)
    unsigned int summaryCount; //!< the number of populated summaries in \ref values
    union
    {
        int64_t i64;
        double fp64;
    } values[DCGM_SUMMARY_SIZE]; //!< array for storing the values of each summary. The summaries are stored
                                 //!< in order. For example, if MIN AND MAX are requested, then 0 will be MIN
                                 //!< and 1 will be MAX. If AVG and DIFF were requested, then AVG would be 0
                                 //!< and 1 would be DIFF
} dcgmSummaryResponse_t;

typedef struct
{
    unsigned int version;                    //!< version of this message - dcgmFieldSummaryRequest_v1
    unsigned short fieldId;                  //!< field id to be summarized
    dcgm_field_entity_group_t entityGroupId; //!< the type of entity whose field we're getting
    dcgm_field_eid_t entityId;               //!< ordinal id for this entity
    uint32_t summaryTypeMask;                //!< bit-mask of DCGM_SUMMARY_*, the requested summaries
    uint64_t startTime;                      //!< start time for the interval being summarized. 0 means to use
                                             //!< any data before.
    uint64_t endTime;                        //!< end time for the interval being summarized. 0 means to use
                                             //!< any data after.
    dcgmSummaryResponse_t response;          //!< response data for this request
} dcgmFieldSummaryRequest_v1;

typedef dcgmFieldSummaryRequest_v1 dcgmFieldSummaryRequest_t;

#define dcgmFieldSummaryRequest_version1 MAKE_DCGM_VERSION(dcgmFieldSummaryRequest_v1, 1)

/**
 * Module IDs
 */
typedef enum
{
    DcgmModuleIdCore       = 0, //!< Core DCGM - always loaded
    DcgmModuleIdNvSwitch   = 1, //!< NvSwitch Module
    DcgmModuleIdVGPU       = 2, //!< VGPU Module
    DcgmModuleIdIntrospect = 3, //!< Introspection Module
    DcgmModuleIdHealth     = 4, //!< Health Module
    DcgmModuleIdPolicy     = 5, //!< Policy Module
    DcgmModuleIdConfig     = 6, //!< Config Module
    DcgmModuleIdDiag       = 7, //!< GPU Diagnostic Module
    DcgmModuleIdProfiling  = 8, //!< Profiling Module
    DcgmModuleIdSysmon     = 9, //!< System Monitoring Module

    DcgmModuleIdCount //!< Always last. 1 greater than largest value above
} dcgmModuleId_t;

/**
 * Module Status. Modules are lazy loaded, so they will be in status DcgmModuleStatusNotLoaded
 * until they are used. One modules are used, they will move to another status.
 */
typedef enum
{
    DcgmModuleStatusNotLoaded  = 0, //!< Module has not been loaded yet
    DcgmModuleStatusDenylisted = 1, //!< Module is on the denylist; can't be loaded
    DcgmModuleStatusFailed     = 2, //!< Loading the module failed
    DcgmModuleStatusLoaded     = 3, //!< Module has been loaded
    DcgmModuleStatusUnloaded   = 4, //!< Module has been unloaded, happens during shutdown
    DcgmModuleStatusPaused     = 5, /*!< Module has been paused. This is a temporary state that will
                                         move to DcgmModuleStatusLoaded once the module is resumed.
                                         This status implies that the module is loaded. */
} dcgmModuleStatus_t;

/**
 * Status of all of the modules of the host engine
 */
typedef struct
{
    dcgmModuleId_t id;         //!< ID of this module
    dcgmModuleStatus_t status; //!< Status of this module
} dcgmModuleGetStatusesModule_t;

/* This is larger than DcgmModuleIdCount so we can add modules without versioning this request */
#define DCGM_MODULE_STATUSES_CAPACITY 16

typedef struct
{
    unsigned int version;     //!< Version of this request. Should be dcgmModuleGetStatuses_version1
    unsigned int numStatuses; //!< Number of entries in statuses[] that are populated
    dcgmModuleGetStatusesModule_t statuses[DCGM_MODULE_STATUSES_CAPACITY]; //!< Per-module status information
} dcgmModuleGetStatuses_v1;

/**
 * Version 1 of dcgmModuleGetStatuses
 */
#define dcgmModuleGetStatuses_version1 MAKE_DCGM_VERSION(dcgmModuleGetStatuses_v1, 1)
#define dcgmModuleGetStatuses_version  dcgmModuleGetStatuses_version1
typedef dcgmModuleGetStatuses_v1 dcgmModuleGetStatuses_t;

/**
 * Options for dcgmStartEmbedded_v2
 *
 * Added in DCGM 2.0.0
 */
typedef struct
{
    unsigned int version;                     /*!< Version number. Use dcgmStartEmbeddedV2Params_version1 */
    dcgmOperationMode_t opMode;               /*!< IN: Collect data automatically or manually when asked by the user. */
    dcgmHandle_t dcgmHandle;                  /*!< OUT: DCGM Handle to use for API calls */
    const char *logFile;                      /*!< IN: File that DCGM should log to. NULL = do not log. '-' = stdout */
    DcgmLoggingSeverity_t severity;           /*!< IN: Severity at which DCGM should log to logFile */
    unsigned int denyListCount;               /*!< IN: Number of modules in denyList[] */
    unsigned int denyList[DcgmModuleIdCount]; /* IN: IDs of modules to add to the denylist */
} dcgmStartEmbeddedV2Params_v1;

/**
 * Version 1 for \ref dcgmStartEmbeddedV2Params_v1
 */
#define dcgmStartEmbeddedV2Params_version1 MAKE_DCGM_VERSION(dcgmStartEmbeddedV2Params_v1, 1)

/**
 * Options for dcgmStartEmbeddedV2Params_v2
 *
 * Added in DCGM 2.4.0, renamed members in 3.0.0
 */
typedef struct
{
    unsigned int version;                     /*!< Version number. Use dcgmStartEmbeddedV2Params_version2 */
    dcgmOperationMode_t opMode;               /*!< IN: Collect data automatically or manually when asked by the user. */
    dcgmHandle_t dcgmHandle;                  /*!< OUT: DCGM Handle to use for API calls */
    const char *logFile;                      /*!< IN: File that DCGM should log to. NULL = do not log. '-' = stdout */
    DcgmLoggingSeverity_t severity;           /*!< IN: Severity at which DCGM should log to logFile */
    unsigned int denyListCount;               /*!< IN: Number of modules to be added to the denylist in denyList[] */
    const char *serviceAccount;               /*!< IN: Service account for unprivileged processes */
    unsigned int denyList[DcgmModuleIdCount]; /*!< IN: IDs of modules to be added to the denylist */
} dcgmStartEmbeddedV2Params_v2;

/**
 * Version 2 for \ref dcgmStartEmbeddedV2Params
 */
#define dcgmStartEmbeddedV2Params_version2 MAKE_DCGM_VERSION(dcgmStartEmbeddedV2Params_v2, 2)

/**
 * Maximum number of metric ID groups that can exist in DCGM
 */
#define DCGM_PROF_MAX_NUM_GROUPS_V2 10

/**
 * Maximum number of field IDs that can be in a single DCGM profiling metric group
 */
#define DCGM_PROF_MAX_FIELD_IDS_PER_GROUP_V2 64

/**
 * Structure to return all of the profiling metric groups that are available for the given groupId.
 */
typedef struct
{
    unsigned short majorId;   //!< Major ID of this metric group. Metric groups with the same majorId cannot be
                              //!< watched concurrently with other metric groups with the same majorId
    unsigned short minorId;   //!< Minor ID of this metric group. This distinguishes metric groups within the same
                              //!< major metric group from each other
    unsigned int numFieldIds; //!< Number of field IDs that are populated in fieldIds[]
    unsigned short fieldIds[DCGM_PROF_MAX_FIELD_IDS_PER_GROUP_V2]; //!< DCGM Field IDs that are part of this profiling
                                                                   //!< group. See DCGM_FI_PROF_* definitions in
                                                                   //!< dcgm_fields.h for details.
} dcgmProfMetricGroupInfo_v2;

typedef struct
{
    /** \name Input parameters
     * @{
     */
    unsigned int version; //!< Version of this request. Should be dcgmProfGetMetricGroups_version
    unsigned int unused;  //!< Not used for now. Set to 0
    unsigned int gpuId;   //!< GPU ID we should get the metric groups for.
    /**
     * @}
     */

    /** \name Output
     * @{
     */
    unsigned int numMetricGroups; //!< Number of entries in metricGroups[] that are populated
    dcgmProfMetricGroupInfo_v2 metricGroups[DCGM_PROF_MAX_NUM_GROUPS_V2]; //!< Info for each metric group
    /**
     * @}
     */
} dcgmProfGetMetricGroups_v3;

/**
 * Version 3 of dcgmProfGetMetricGroups_t. See dcgm_structs_24.h for v2
 */
#define dcgmProfGetMetricGroups_version3 MAKE_DCGM_VERSION(dcgmProfGetMetricGroups_v3, 3)
#define dcgmProfGetMetricGroups_version  dcgmProfGetMetricGroups_version3
typedef dcgmProfGetMetricGroups_v3 dcgmProfGetMetricGroups_t;

/**
 * Structure to pass to dcgmProfWatchFields() when watching profiling metrics
 */
typedef struct
{
    unsigned int version;        //!< Version of this request. Should be dcgmProfWatchFields_version
    dcgmGpuGrp_t groupId;        //!< Group ID representing collection of one or more GPUs. Look at \ref dcgmGroupCreate
                                 //!< for details on creating the group. Alternatively, pass in the group id as \a
                                 //!< DCGM_GROUP_ALL_GPUS to perform operation on all the GPUs. The GPUs of the group
                                 //!< must all be identical or DCGM_ST_GROUP_INCOMPATIBLE will be returned by this API.
    unsigned int numFieldIds;    //!< Number of field IDs that are being passed in fieldIds[]
    unsigned short fieldIds[64]; //!< DCGM_FI_PROF_? field IDs to watch
    long long updateFreq;        //!< How often to update this field in usec. Note that profiling metrics may need to be
                                 //!< sampled more frequently than this value. See
                                 //!< dcgmProfMetricGroupInfo_t.minUpdateFreqUsec of the metric group matching
                                 //!< metricGroupTag to see what this minimum is. If minUpdateFreqUsec < updateFreq
                                 //!< then samples will be aggregated to updateFreq intervals in DCGM's internal cache.
    double maxKeepAge;           //!< How long to keep data for every fieldId in seconds
    int maxKeepSamples;          //!< Maximum number of samples to keep for each fieldId. 0=no limit
    unsigned int flags;          //!< For future use. Set to 0 for now.
} dcgmProfWatchFields_v2;

/**
 * Version 2 of dcgmProfWatchFields_v2
 */
#define dcgmProfWatchFields_version2 MAKE_DCGM_VERSION(dcgmProfWatchFields_v2, 2)
#define dcgmProfWatchFields_version  dcgmProfWatchFields_version2
typedef dcgmProfWatchFields_v2 dcgmProfWatchFields_t;

/**
 * Structure to pass to dcgmProfUnwatchFields when unwatching profiling metrics
 */
typedef struct
{
    unsigned int version; //!< Version of this request. Should be dcgmProfUnwatchFields_version
    dcgmGpuGrp_t groupId; //!< Group ID representing collection of one or more GPUs. Look at
                          //!< \ref dcgmGroupCreate for details on creating the group.
                          //!< Alternatively, pass in the group id as \a DCGM_GROUP_ALL_GPUS
                          //!< to perform operation on all the GPUs. The GPUs of the group must all be
                          //!< identical or DCGM_ST_GROUP_INCOMPATIBLE will be returned by this API.
    unsigned int flags;   //!< For future use. Set to 0 for now.
} dcgmProfUnwatchFields_v1;

/**
 * Version 1 of dcgmProfUnwatchFields_v1
 */
#define dcgmProfUnwatchFields_version1 MAKE_DCGM_VERSION(dcgmProfUnwatchFields_v1, 1)
#define dcgmProfUnwatchFields_version  dcgmProfUnwatchFields_version1
typedef dcgmProfUnwatchFields_v1 dcgmProfUnwatchFields_t;

/**
 * Version 1 of dcgmSettingsSetLoggingSeverity_t
 */
typedef struct
{
    int targetLogger;
    DcgmLoggingSeverity_t targetSeverity;
} dcgmSettingsSetLoggingSeverity_v1;


#define dcgmSettingsSetLoggingSeverity_version1 MAKE_DCGM_VERSION(dcgmSettingsSetLoggingSeverity_v1, 1)
#define dcgmSettingsSetLoggingSeverity_version  dcgmSettingsSetLoggingSeverity_version1
typedef dcgmSettingsSetLoggingSeverity_v1 dcgmSettingsSetLoggingSeverity_t;

/**
 * Structure to describe the DCGM build environment ver 2.0
 */
typedef struct
{
    unsigned int version; //<! Structure version
    /**
     * Raw form of the DCGM build info. There may be multiple kv-pairs separated by semicolon (;).<br>
     * Every pair is separated by a colon char (:). Only the very first colon is considered as a separation.<br>
     * Values can contain colon chars. Values and Keys cannot contain semicolon chars.<br>
     * Usually defined keys are:
     *      <p style="margin-left:20px">
     *      <i>version</i> : DCGM Version.<br>
     *      <i>arch</i>    : Target DCGM Architecture.<br>
     *      <i>buildid</i> : Build ID. Usually a sequential number.<br>
     *      <i>commit</i>  : Commit ID (Usually a git commit hash).<br>
     *      <i>author</i>  : Author of the commit above.<br>
     *      <i>branch</i>  : Branch (Usually a git branch that was used for the build).<br>
     *      <i>buildtype</i> : Build Type.<br>
     *      <i>builddate</i> : Date of the build.<br>
     *      <i>buildplatform</i>   : Platform where the build was made.<br>
     *      </p>
     * Any or all keys may be absent.<br>
     * This values are for reference only are not supposed to participate in some complicated logic.<br>
     */
    char rawBuildInfoString[DCGM_MAX_STR_LENGTH * 2];
} dcgmVersionInfo_v2;

/**
 * Version 2 of the dcgmVersionInfo_v2
 */
#define dcgmVersionInfo_version2 MAKE_DCGM_VERSION(dcgmVersionInfo_v2, 2)

#define dcgmVersionInfo_version dcgmVersionInfo_version2
typedef dcgmVersionInfo_v2 dcgmVersionInfo_t;

/** @} */

#ifdef __cplusplus
}
#endif

#endif /* DCGM_STRUCTS_H */
