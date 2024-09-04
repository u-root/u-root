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
#ifndef DCGM_ERRORS_H
#define DCGM_ERRORS_H

#define DCGM_PUBLIC_API
#include "dcgm_structs.h"

/***************************************************************************************************/
/** @defgroup dcgmErrorEnums Error Codes
 *  @{
 */
/***************************************************************************************************/
/*
 * Error codes for passive and active health checks.
 * New error codes must be added to end of enum to maintain backwards compatibility.
 */
typedef enum dcgmError_enum
{
    DCGM_FR_OK                         = 0,  //!< 0 No error
    DCGM_FR_UNKNOWN                    = 1,  //!< 1 Unknown error code
    DCGM_FR_UNRECOGNIZED               = 2,  //!< 2 Unrecognized error code
    DCGM_FR_PCI_REPLAY_RATE            = 3,  //!< 3 Unacceptable rate of PCI errors
    DCGM_FR_VOLATILE_DBE_DETECTED      = 4,  //!< 4 Uncorrectable volatile double bit error
    DCGM_FR_VOLATILE_SBE_DETECTED      = 5,  //!< 5 Unacceptable rate of volatile single bit errors
    DCGM_FR_PENDING_PAGE_RETIREMENTS   = 6,  //!< 6 Pending page retirements detected
    DCGM_FR_RETIRED_PAGES_LIMIT        = 7,  //!< 7 Unacceptable total page retirements detected
    DCGM_FR_RETIRED_PAGES_DBE_LIMIT    = 8,  //!< 8 Unacceptable total page retirements due to uncorrectable errors
    DCGM_FR_CORRUPT_INFOROM            = 9,  //!< 9 Corrupt inforom found
    DCGM_FR_CLOCK_THROTTLE_THERMAL     = 10, //!< 10 Clocks being throttled due to overheating
    DCGM_FR_POWER_UNREADABLE           = 11, //!< 11 Cannot get a reading for power from NVML
    DCGM_FR_CLOCK_THROTTLE_POWER       = 12, //!< 12 Clock being throttled due to power restrictions
    DCGM_FR_NVLINK_ERROR_THRESHOLD     = 13, //!< 13 Unacceptable rate of NVLink errors
    DCGM_FR_NVLINK_DOWN                = 14, //!< 14 NVLink is down
    DCGM_FR_NVSWITCH_FATAL_ERROR       = 15, //!< 15 Fatal errors on the NVSwitch
    DCGM_FR_NVSWITCH_NON_FATAL_ERROR   = 16, //!< 16 Non-fatal errors on the NVSwitch
    DCGM_FR_NVSWITCH_DOWN              = 17, //!< 17 NVSwitch is down - NOT USED: DEPRECATED
    DCGM_FR_NO_ACCESS_TO_FILE          = 18, //!< 18 Cannot access a file
    DCGM_FR_NVML_API                   = 19, //!< 19 Error occurred on an NVML API - NOT USED: DEPRECATED
    DCGM_FR_DEVICE_COUNT_MISMATCH      = 20, //!< 20 Disagreement in GPU count between /dev and NVML
    DCGM_FR_BAD_PARAMETER              = 21, //!< 21 Bad parameter passed to API
    DCGM_FR_CANNOT_OPEN_LIB            = 22, //!< 22 Cannot open a library that must be accessed
    DCGM_FR_DENYLISTED_DRIVER          = 23, //!< 23 A driver on the denylist (nouveau) is active
    DCGM_FR_NVML_LIB_BAD               = 24, //!< 24 NVML library is missing expected functions - NOT USED: DEPRECATED
    DCGM_FR_GRAPHICS_PROCESSES         = 25, //!< 25 Graphics processes are active on this GPU
    DCGM_FR_HOSTENGINE_CONN            = 26, //!< 26 Bad connection to nv-hostengine - NOT USED: DEPRECATED
    DCGM_FR_FIELD_QUERY                = 27, //!< 27 Error querying a field from DCGM
    DCGM_FR_BAD_CUDA_ENV               = 28, //!< 28 The environment has variables that hurt CUDA
    DCGM_FR_PERSISTENCE_MODE           = 29, //!< 29 Persistence mode is disabled
    DCGM_FR_LOW_BANDWIDTH              = 30, //!< 30 The bandwidth is unacceptably low
    DCGM_FR_HIGH_LATENCY               = 31, //!< 31 Latency is too high
    DCGM_FR_CANNOT_GET_FIELD_TAG       = 32, //!< 32 Cannot find a tag for a field
    DCGM_FR_FIELD_VIOLATION            = 33, //!< 33 The value for the specified error field is above 0
    DCGM_FR_FIELD_THRESHOLD            = 34, //!< 34 The value for the specified field is above the threshold
    DCGM_FR_FIELD_VIOLATION_DBL        = 35, //!< 35 The value for the specified error field is above 0
    DCGM_FR_FIELD_THRESHOLD_DBL        = 36, //!< 36 The value for the specified field is above the threshold
    DCGM_FR_UNSUPPORTED_FIELD_TYPE     = 37, //!< 37 Field type cannot be supported
    DCGM_FR_FIELD_THRESHOLD_TS         = 38, //!< 38 The value for the specified field is above the threshold
    DCGM_FR_FIELD_THRESHOLD_TS_DBL     = 39, //!< 39 The value for the specified field is above the threshold
    DCGM_FR_THERMAL_VIOLATIONS         = 40, //!< 40 Thermal violations detected
    DCGM_FR_THERMAL_VIOLATIONS_TS      = 41, //!< 41 Thermal violations detected with a timestamp
    DCGM_FR_TEMP_VIOLATION             = 42, //!< 42 Temperature is too high
    DCGM_FR_THROTTLING_VIOLATION       = 43, //!< 43 Non-benign clock throttling is occurring
    DCGM_FR_INTERNAL                   = 44, //!< 44 An internal error was detected
    DCGM_FR_PCIE_GENERATION            = 45, //!< 45 PCIe generation is too low
    DCGM_FR_PCIE_WIDTH                 = 46, //!< 46 PCIe width is too low
    DCGM_FR_ABORTED                    = 47, //!< 47 Test was aborted by a user signal
    DCGM_FR_TEST_DISABLED              = 48, //!< 48 This test is disabled for this GPU
    DCGM_FR_CANNOT_GET_STAT            = 49, //!< 49 Cannot get telemetry for a needed value
    DCGM_FR_STRESS_LEVEL               = 50, //!< 50 Stress level is too low (bad performance)
    DCGM_FR_CUDA_API                   = 51, //!< 51 Error calling the specified CUDA API
    DCGM_FR_FAULTY_MEMORY              = 52, //!< 52 Faulty memory detected on this GPU
    DCGM_FR_CANNOT_SET_WATCHES         = 53, //!< 53 Unable to set field watches in DCGM - NOT USED: DEPRECATED
    DCGM_FR_CUDA_UNBOUND               = 54, //!< 54 CUDA context is no longer bound
    DCGM_FR_ECC_DISABLED               = 55, //!< 55 ECC memory is disabled right now
    DCGM_FR_MEMORY_ALLOC               = 56, //!< 56 Cannot allocate memory on the GPU
    DCGM_FR_CUDA_DBE                   = 57, //!< 57 CUDA detected unrecovable double-bit error
    DCGM_FR_MEMORY_MISMATCH            = 58, //!< 58 Memory error detected
    DCGM_FR_CUDA_DEVICE                = 59, //!< 59 No CUDA device discoverable for existing GPU
    DCGM_FR_ECC_UNSUPPORTED            = 60, //!< 60 ECC memory is unsupported by this SKU
    DCGM_FR_ECC_PENDING                = 61, //!< 61 ECC memory is in a pending state - NOT USED: DEPRECATED
    DCGM_FR_MEMORY_BANDWIDTH           = 62, //!< 62 Memory bandwidth is too low
    DCGM_FR_TARGET_POWER               = 63, //!< 63 Cannot hit the target power draw
    DCGM_FR_API_FAIL                   = 64, //!< 64 The specified API call failed
    DCGM_FR_API_FAIL_GPU               = 65, //!< 65 The specified API call failed for the specified GPU
    DCGM_FR_CUDA_CONTEXT               = 66, //!< 66 Cannot create a CUDA context on this GPU
    DCGM_FR_DCGM_API                   = 67, //!< 67 DCGM API failure
    DCGM_FR_CONCURRENT_GPUS            = 68, //!< 68 Need multiple GPUs to run this test
    DCGM_FR_TOO_MANY_ERRORS            = 69, //!< 69 More errors than fit in the return struct - NOT USED: DEPRECATED
    DCGM_FR_NVLINK_CRC_ERROR_THRESHOLD = 70, //!< 70 More than 100 CRC errors are happening per second
    DCGM_FR_NVLINK_ERROR_CRITICAL      = 71, //!< 71 NVLink error for a field that should always be 0
    DCGM_FR_ENFORCED_POWER_LIMIT       = 72, //!< 72 The enforced power limit is too low to hit the target
    DCGM_FR_MEMORY_ALLOC_HOST          = 73, //!< 73 Cannot allocate memory on the host
    DCGM_FR_GPU_OP_MODE                = 74, //!< 74 Bad GPU operating mode for running plugin - NOT USED: DEPRECATED
    DCGM_FR_NO_MEMORY_CLOCKS           = 75, //!< 75 No memory clocks with the needed MHz found - NOT USED: DEPRECATED
    DCGM_FR_NO_GRAPHICS_CLOCKS         = 76, //!< 76 No graphics clocks with the needed MHz found - NOT USED: DEPRECATED
    DCGM_FR_HAD_TO_RESTORE_STATE       = 77, //!< 77 Note that we had to restore a GPU's state
    DCGM_FR_L1TAG_UNSUPPORTED          = 78, //!< 78 L1TAG test is unsupported by this SKU
    DCGM_FR_L1TAG_MISCOMPARE           = 79, //!< 79 L1TAG test failed on a miscompare
    DCGM_FR_ROW_REMAP_FAILURE          = 80, //!< 80 Row remapping failed (Ampere or newer GPUs)
    DCGM_FR_UNCONTAINED_ERROR          = 81, //!< 81 Uncontained error - XID 95
    DCGM_FR_EMPTY_GPU_LIST             = 82, //!< 82 No GPU information given to plugin
    DCGM_FR_DBE_PENDING_PAGE_RETIREMENTS = 83, //!< 83 Pending page retirements due to a DBE
    DCGM_FR_UNCORRECTABLE_ROW_REMAP      = 84, //!< 84 Uncorrectable row remapping
    DCGM_FR_PENDING_ROW_REMAP            = 85, //!< 85 Row remapping is pending
    DCGM_FR_BROKEN_P2P_MEMORY_DEVICE     = 86, //!< 86 P2P copy test detected an error writing to this GPU
    DCGM_FR_BROKEN_P2P_WRITER_DEVICE     = 87, //!< 87 P2P copy test detected an error writing from this GPU
    DCGM_FR_NVSWITCH_NVLINK_DOWN    = 88, //!< 88 An NvLink is down for the specified NVSwitch - NOT USED: DEPRECATED
    DCGM_FR_EUD_BINARY_PERMISSIONS  = 89, //!< 89 EUD binary permissions are incorrect
    DCGM_FR_EUD_NON_ROOT_USER       = 90, //!< 90 EUD plugin is not running as root
    DCGM_FR_EUD_SPAWN_FAILURE       = 91, //!< 91 EUD plugin failed to spawn the EUD binary
    DCGM_FR_EUD_TIMEOUT             = 92, //!< 92 EUD plugin timed out
    DCGM_FR_EUD_ZOMBIE              = 93, //!< 93 EUD process remains running after the plugin considers it finished
    DCGM_FR_EUD_NON_ZERO_EXIT_CODE  = 94, //!< 94 EUD process exited with a non-zero exit code
    DCGM_FR_EUD_TEST_FAILED         = 95, //!< 95 EUD test failed
    DCGM_FR_FILE_CREATE_PERMISSIONS = 96, //!< 96 We cannot create a file in this directory.
    DCGM_FR_PAUSE_RESUME_FAILED     = 97, //!< 97 Pause/Resume failed
    DCGM_FR_PCIE_H_REPLAY_VIOLATION = 98, //!< 98 PCIe test caught correctable errors
    DCGM_FR_GPU_EXPECTED_NVLINKS_UP = 99, //!< 99 Expected nvlinks up per gpu
    DCGM_FR_NVSWITCH_EXPECTED_NVLINKS_UP    = 100, //!< 100 Expected nvlinks up per nvswitch
    DCGM_FR_XID_ERROR                       = 101, //!< 101 XID error detected
    DCGM_FR_SBE_VIOLATION                   = 102, //!< 102 Single bit error detected
    DCGM_FR_DBE_VIOLATION                   = 103, //!< 103 Double bit error detected
    DCGM_FR_PCIE_REPLAY_VIOLATION           = 104, //!< 104 PCIe replay errors detected
    DCGM_FR_SBE_THRESHOLD_VIOLATION         = 105, //!< 105 SBE threshold violated
    DCGM_FR_DBE_THRESHOLD_VIOLATION         = 106, //!< 106 DBE threshold violated
    DCGM_FR_PCIE_REPLAY_THRESHOLD_VIOLATION = 107, //!< 107 PCIE replay count violated
    DCGM_FR_CUDA_FM_NOT_INITIALIZED         = 108, //!< 108 The fabricmanager is not initialized
    DCGM_FR_SXID_ERROR                      = 109, //!< 109 NvSwitch fatal error detected
    DCGM_FR_ERROR_SENTINEL                  = 110, //!< 110 MUST BE THE LAST ERROR CODE
} dcgmError_t;

typedef enum dcgmErrorSeverity_enum
{
    DCGM_ERROR_NONE    = 0, //!< 0 NONE
    DCGM_ERROR_MONITOR = 1, //!< 1 Can perform workload, but needs to be monitored.
    DCGM_ERROR_ISOLATE = 2, //!< 2 Cannot perform workload. GPU should be isolated.
    DCGM_ERROR_UNKNOWN = 3, //!< 3 This error code is not recognized
    DCGM_ERROR_TRIAGE  = 4, //!< 4 This error should be triaged
    DCGM_ERROR_CONFIG  = 5, //!< 5 This error can be configured
    DCGM_ERROR_RESET   = 6, //!< 6 Drain and reset GPU
} dcgmErrorSeverity_t;

typedef enum dcgmErrorCategory_enum
{
    DCGM_FR_EC_NONE              = 0,  //!< 0 NONE
    DCGM_FR_EC_PERF_THRESHOLD    = 1,  //!< 1 Performance Threshold
    DCGM_FR_EC_PERF_VIOLATION    = 2,  //!< 2 Performance Violation
    DCGM_FR_EC_SOFTWARE_CONFIG   = 3,  //!< 3 Software Configuration
    DCGM_FR_EC_SOFTWARE_LIBRARY  = 4,  //!< 4 Software Library
    DCGM_FR_EC_SOFTWARE_XID      = 5,  //!< 5 Software XID
    DCGM_FR_EC_SOFTWARE_CUDA     = 6,  //!< 6 Software Cuda
    DCGM_FR_EC_SOFTWARE_EUD      = 7,  //!< 7 Software EUD
    DCGM_FR_EC_SOFTWARE_OTHER    = 8,  //!< 8 Software Other
    DCGM_FR_EC_HARDWARE_THERMAL  = 9,  //!< 9 Hardware Thermal
    DCGM_FR_EC_HARDWARE_MEMORY   = 10, //!< 10 Hardware Memory
    DCGM_FR_EC_HARDWARE_NVLINK   = 11, //!< 11 Hardware NvLink
    DCGM_FR_EC_HARDWARE_NVSWITCH = 12, //!< 12 Hardware NvSwitch
    DCGM_FR_EC_HARDWARE_PCIE     = 13, //!< 13 Hardware PCIe
    DCGM_FR_EC_HARDWARE_POWER    = 14, //!< 14 Hardware Power
    DCGM_FR_EC_HARDWARE_OTHER    = 15, //!< 15 Hardware Other
    DCGM_FR_EC_INTERNAL_OTHER    = 16, //!< 16 Internal Other
} dcgmErrorCategory_t;

typedef struct
{
    dcgmError_t errorId;
    const char *msgFormat;
    const char *suggestion;
    int severity;
    int category;
} dcgm_error_meta_t;

extern dcgm_error_meta_t dcgmErrorMeta[];


/* Standard message for running a field diagnostic */
#define TRIAGE_RUN_FIELD_DIAG_MSG "Run a field diagnostic on the GPU."
#define DEBUG_COOLING_MSG                                                         \
    "Verify that the cooling on this machine is functional, including external, " \
    "thermal material interface, fans, and any other components."
#define BUG_REPORT_MSG    "Please capture an nvidia-bug-report and send it to NVIDIA."
#define SYSTEM_TRIAGE_MSG "Check DCGM and system logs for errors. Reset GPU. Restart DCGM. Rerun diagnostics."
#define CONFIG_MSG        "Check DCGM and system configuration. This error may be eliminated with an updated configuration."

/*
 * Messages for the error codes. All messages must be defined in the ERROR_CODE_MSG <msg> format
 * where <msg> is the actual message.
 */
#define DCGM_FR_OK_MSG           "The operation completed successfully."
#define DCGM_FR_UNKNOWN_MSG      "Unknown error."
#define DCGM_FR_UNRECOGNIZED_MSG "Unrecognized error code."
// replay limit, gpu id, replay errors detected
#define DCGM_FR_PCI_REPLAY_RATE_MSG "Detected more than %u PCIe replays per minute for GPU %u : %d"
// dbes deteced, gpu id
#define DCGM_FR_VOLATILE_DBE_DETECTED_MSG "Detected %d volatile double-bit ECC error(s) in GPU %u."
// sbe limit, gpu id, sbes detected
#define DCGM_FR_VOLATILE_SBE_DETECTED_MSG "More than %u single-bit ECC error(s) detected in GPU %u Volatile SBEs: %lld"
// gpu id
#define DCGM_FR_PENDING_PAGE_RETIREMENTS_MSG "A pending retired page has been detected in GPU %u."
// retired pages detected, gpud id
#define DCGM_FR_RETIRED_PAGES_LIMIT_MSG "%u or more retired pages have been detected in GPU %u. "
// retired pages due to dbes detected, gpu id
#define DCGM_FR_RETIRED_PAGES_DBE_LIMIT_MSG                            \
    "An excess of %u retired pages due to DBEs have been detected and" \
    " more than one page has been retired due to DBEs in the past"     \
    " week in GPU %u."
// gpu id
#define DCGM_FR_CORRUPT_INFOROM_MSG "A corrupt InfoROM has been detected in GPU %u."
// gpu id
#define DCGM_FR_CLOCK_THROTTLE_THERMAL_MSG "Detected clock throttling due to thermal violation in GPU %u."
// gpu id
#define DCGM_FR_POWER_UNREADABLE_MSG "Cannot reliably read the power usage for GPU %u."
// gpu id
#define DCGM_FR_CLOCK_THROTTLE_POWER_MSG "Detected clock throttling due to power violation in GPU %u."
// nvlink errors detected, nvlink id, error threshold
#define DCGM_FR_NVLINK_ERROR_THRESHOLD_MSG                            \
    "Detected %ld %s NvLink errors on GPU %u's NVLink which exceeds " \
    "threshold of %u"
// gpu id, nvlink id
#define DCGM_FR_NVLINK_DOWN_MSG "GPU %u's NvLink link %d is currently down"
// nvlinks up, expected nvlinks up
#define DCGM_FR_GPU_EXPECTED_NVLINKS_UP_MSG "Only %u NvLinks are up out of the expected %u"
// switch id, nvlinks up, expected nvlinks up
#define DCGM_FR_NVSWITCH_EXPECTED_NVLINKS_UP_MSG "NvSwitch %u - Only %u NvLinks are up out of the expected %u"
// nvswitch id, nvlink id
#define DCGM_FR_NVSWITCH_FATAL_ERROR_MSG "Detected fatal errors on NvSwitch %u link %u"
// nvswitch id, nvlink id
#define DCGM_FR_NVSWITCH_NON_FATAL_ERROR_MSG "Detected nonfatal errors on NvSwitch %u link %u"
// nvswitch id, nvlink port
#define DCGM_FR_NVSWITCH_DOWN_MSG "NvSwitch physical ID %u's NvLink port %d is currently down."
// file path, error detail
#define DCGM_FR_NO_ACCESS_TO_FILE_MSG "File %s could not be accessed directly: %s"
// purpose for communicating with NVML, NVML error as string, NVML error
#define DCGM_FR_NVML_API_MSG "Error calling NVML API %s: %s"
#define DCGM_FR_DEVICE_COUNT_MISMATCH_MSG                              \
    "The number of devices NVML returns is different than the number " \
    "of devices in /dev."
// function name
#define DCGM_FR_BAD_PARAMETER_MSG "Bad parameter to function %s cannot be processed"
// library name, error returned from dlopen
#define DCGM_FR_CANNOT_OPEN_LIB_MSG "Cannot open library %s: '%s'"
// the name of the denylisted driver
#define DCGM_FR_DENYLISTED_DRIVER_MSG "Found driver on the denylist: %s"
// the name of the function that wasn't found
#define DCGM_FR_NVML_LIB_BAD_MSG "Cannot get pointer to %s from libnvidia-ml.so"
#define DCGM_FR_GRAPHICS_PROCESSES_MSG                                                 \
    "NVVS has detected processes with graphics contexts open running on at least one " \
    "GPU. This may cause some tests to fail."
// error message from the API call
#define DCGM_FR_HOSTENGINE_CONN_MSG "Could not connect to the host engine: '%s'"
// field name, gpu id
#define DCGM_FR_FIELD_QUERY_MSG "Could not query field %s for GPU %u"
// environment variable name
#define DCGM_FR_BAD_CUDA_ENV_MSG "Found CUDA performance-limiting environment variable '%s'."
// gpu id
#define DCGM_FR_PERSISTENCE_MODE_MSG "Persistence mode for GPU %u is disabled."
// gpu id, direction (d2h, e.g.), measured bandwidth, expected bandwidth
#define DCGM_FR_LOW_BANDWIDTH_MSG                                 \
    "Bandwidth of GPU %u in direction %s of %.2f did not exceed " \
    "minimum required bandwidth of %.2f."
// gpu id, direction (d2h, e.g.), measured latency, expected latency
#define DCGM_FR_HIGH_LATENCY_MSG                                     \
    "Latency type %s of GPU %u value %.2f exceeded maximum allowed " \
    "latency of %.2f."
// field id
#define DCGM_FR_CANNOT_GET_FIELD_TAG_MSG "Unable to get field information for field id %hu"
// field value, field name, gpu id (this message is for fields that should always have a 0 value)
#define DCGM_FR_FIELD_VIOLATION_MSG "Detected %ld %s for GPU %u"
// field value, field name, gpu id, allowable threshold
#define DCGM_FR_FIELD_THRESHOLD_MSG "Detected %ld %s for GPU %u which is above the threshold %ld"
// field value, field name, gpu id (same as DCGM_FR_FIELD_VIOLATION, but it's a double)
#define DCGM_FR_FIELD_VIOLATION_DBL_MSG "Detected %.1f %s for GPU %u"
// field value, field name, gpu id, allowable threshold (same as DCGM_FR_FIELD_THRESHOLD, but it's a double)
#define DCGM_FR_FIELD_THRESHOLD_DBL_MSG "Detected %.1f %s for GPU %u which is above the threshold %.1f"
// field name
#define DCGM_FR_UNSUPPORTED_FIELD_TYPE_MSG                            \
    "Field %s is not supported by this API because it is neither an " \
    "int64 nor a double type."
// field name, allowable threshold, observed value, seconds
#define DCGM_FR_FIELD_THRESHOLD_TS_MSG                            \
    "%s met or exceeded the threshold of %lu per second: %lu at " \
    "%.1f seconds into the test."
// field name, allowable threshold, observed value, seconds (same as DCGM_FR_FIELD_THRESHOLD, but it's a double)
#define DCGM_FR_FIELD_THRESHOLD_TS_DBL_MSG                          \
    "%s met or exceeded the threshold of %.1f per second: %.1f at " \
    "%.1f seconds into the test."
// total seconds of violation, gpu id
#define DCGM_FR_THERMAL_VIOLATIONS_MSG "There were thermal violations totaling %.1f seconds for GPU %u"
// total seconds of violations, first instance, gpu id
#define DCGM_FR_THERMAL_VIOLATIONS_TS_MSG                               \
    "Thermal violations totaling %.1f seconds started at %.1f seconds " \
    "into the test for GPU %u"
// observed temperature, gpu id, max allowed temperature
#define DCGM_FR_TEMP_VIOLATION_MSG                                \
    "Temperature %lld of GPU %u exceeded user-specified maximum " \
    "allowed temperature %lld"
// gpu id, seconds into test, details about throttling
#define DCGM_FR_THROTTLING_VIOLATION_MSG                      \
    "Clocks are being throttled for GPU %u because of clock " \
    "throttling starting %.1f seconds into the test. %s"
// details about error
#define DCGM_FR_INTERNAL_MSG "There was an internal error during the test: '%s'"
// gpu id, PCIe generation, minimum allowed, parameter to control
#define DCGM_FR_PCIE_GENERATION_MSG                                \
    "GPU %u is running at PCI link generation %d, which is below " \
    "the minimum allowed link generation of %d (parameter '%s')"
// gpu id, PCIe width, minimum allowed, parameter to control
#define DCGM_FR_PCIE_WIDTH_MSG                                     \
    "GPU %u is running at PCI link width %dX, which is below the " \
    "minimum allowed link generation of %d (parameter '%s')"
#define DCGM_FR_ABORTED_MSG "Test was aborted early due to user signal"
// Test name
#define DCGM_FR_TEST_DISABLED_MSG "The %s test is skipped for this GPU."
// stat name, gpu id
#define DCGM_FR_CANNOT_GET_STAT_MSG "Unable to generate / collect stat %s for GPU %u"
// observed value, minimum allowed, gpu id
#define DCGM_FR_STRESS_LEVEL_MSG                                      \
    "Max stress level of %.1f did not reach desired stress level of " \
    "%.1f for GPU %u"
// CUDA API name
#define DCGM_FR_CUDA_API_MSG "Error using CUDA API %s"
// count, gpu id
#define DCGM_FR_FAULTY_MEMORY_MSG "Found %d faulty memory elements on GPU %u"
// error detail
#define DCGM_FR_CANNOT_SET_WATCHES_MSG "Unable to add field watches to DCGM: %s"
// gpu id
#define DCGM_FR_CUDA_UNBOUND_MSG "Cuda GPU %d is no longer bound to a CUDA context...Aborting"
// Test name, gpu id
#define DCGM_FR_ECC_DISABLED_MSG "Skipping test %s because ECC is not enabled on GPU %u"
// percentage of memory we tried to allocate, gpu id
#define DCGM_FR_MEMORY_ALLOC_MSG "Couldn't allocate at least %.1f%% of GPU memory on GPU %u"
// gpu id
#define DCGM_FR_CUDA_DBE_MSG                                    \
    "CUDA APIs have indicated that a double-bit ECC error has " \
    "occured on GPU %u."
// gpu id
#define DCGM_FR_MEMORY_MISMATCH_MSG                               \
    "A memory mismatch was detected on GPU %u, but no error was " \
    "reported by CUDA or NVML."
// gpu id, error detail
#define DCGM_FR_CUDA_DEVICE_MSG     "Unable to find a corresponding CUDA device for GPU %u: '%s'"
#define DCGM_FR_ECC_UNSUPPORTED_MSG "ECC Memory is not turned on or is unsupported. Skipping test."
// gpu id
#define DCGM_FR_ECC_PENDING_MSG "ECC memory for GPU %u is in a pending state."
// gpu id, observed bandwidth, required, test name
#define DCGM_FR_MEMORY_BANDWIDTH_MSG                                 \
    "GPU %u only achieved a memory bandwidth of %.2f GB/s, failing " \
    "to meet %.2f GB/s for test %d"
// power draw observed, field tag, minimum power draw required, gpu id
#define DCGM_FR_TARGET_POWER_MSG                                   \
    "Max power of %.1f did not reach desired power minimum %s of " \
    "%.1f for GPU %u"
// API name, error detail
#define DCGM_FR_API_FAIL_MSG "API call %s failed: '%s'"
// API name, gpu id, error detail
#define DCGM_FR_API_FAIL_GPU_MSG "API call %s failed for GPU %u: '%s'"
// gpu id, error detail
#define DCGM_FR_CUDA_CONTEXT_MSG "GPU %u failed to create a CUDA context: %s"
// DCGM API name
#define DCGM_FR_DCGM_API_MSG "Error using DCGM API %s"
#define DCGM_FR_CONCURRENT_GPUS_MSG                                   \
    "Unable to run concurrent pair bandwidth test without 2 or more " \
    "gpus. Skipping"
#define DCGM_FR_TOO_MANY_ERRORS_MSG                                  \
    "This API can only return up to four errors per system. "        \
    "Additional errors were found for this system that couldn't be " \
    "communicated."
// error count, gpu id
#define DCGM_FR_NVLINK_CRC_ERROR_THRESHOLD_MSG                    \
    "%.1f %s NvLink errors found occuring per second on GPU %u, " \
    "exceeding the limit of 100 per second."
// error count, field name, gpu id
#define DCGM_FR_NVLINK_ERROR_CRITICAL_MSG "Detected %ld %s NvLink errors on GPU %u's NVLink (should be 0)"
// gpu id, power limit, power reached
#define DCGM_FR_ENFORCED_POWER_LIMIT_MSG                               \
    "Enforced power limit on GPU %u set to %.1f, which is too low to " \
    "attempt to achieve target power %.1f"
// memory
#define DCGM_FR_MEMORY_ALLOC_HOST_MSG "Cannot allocate %zu bytes on the host"
#define DCGM_FR_GPU_OP_MODE_MSG       "Skipping plugin due to a GPU being in GPU Operating Mode: LOW_DP."
// clock, count
#define DCGM_FR_NO_MEMORY_CLOCKS_MSG "No memory clocks <= %u MHZ were found in %u supported memory clocks."
// clock, count, clock
#define DCGM_FR_NO_GRAPHICS_CLOCKS_MSG \
    "No graphics clocks <= %u MHZ were found in %u supported graphics clocks for memory clock %u MHZ."
// error detail
#define DCGM_FR_HAD_TO_RESTORE_STATE_MSG "Had to restore GPU state on NVML GPU(s): %s"
#define DCGM_FR_L1TAG_UNSUPPORTED_MSG    "This card does not support the L1 cache test. Skipping test."
#define DCGM_FR_L1TAG_MISCOMPARE_MSG     "Detected a miscompare failure in the L1 cache."
// gpu id
#define DCGM_FR_ROW_REMAP_FAILURE_MSG            "GPU %u had uncorrectable memory errors and row remapping failed."
#define DCGM_FR_UNCONTAINED_ERROR_MSG            "GPU had an uncontained error (XID 95)"
#define DCGM_FR_EMPTY_GPU_LIST_MSG               "No valid GPUs passed to plugin"
#define DCGM_FR_DBE_PENDING_PAGE_RETIREMENTS_MSG "Pending page retirements together with a DBE were detected on GPU %u."
// gpu id, rows remapped
#define DCGM_FR_UNCORRECTABLE_ROW_REMAP_MSG "GPU %u had uncorrectable memory errors and %u rows were remapped"
// gpu id
#define DCGM_FR_PENDING_ROW_REMAP_MSG "GPU %u had memory errors and row remappings are pending"
// gpu id, test name
#define DCGM_FR_BROKEN_P2P_MEMORY_DEVICE_MSG "GPU %u was unsuccessfully written to in a peer-to-peer test: %s"
// gpu id, test name
#define DCGM_FR_BROKEN_P2P_WRITER_DEVICE_MSG "GPU %u unsuccessfully wrote data in a peer-to-peer test: %s"
// nvswitch id, nvlink id
#define DCGM_FR_NVSWITCH_NVLINK_DOWN_MSG   "NVSwitch %u's NvLink %u is down."
#define DCGM_FR_EUD_BINARY_PERMISSIONS_MSG "" /* See message inplace */
#define DCGM_FR_EUD_NON_ROOT_USER_MSG      "" /* See message inplace */
#define DCGM_FR_EUD_SPAWN_FAILURE_MSG      "" /* See message inplace */
#define DCGM_FR_EUD_TIMEOUT_MSG            "" /* See message inplace */
#define DCGM_FR_EUD_ZOMBIE_MSG             "" /* See message inplace */
#define DCGM_FR_EUD_NON_ZERO_EXIT_CODE_MSG "" /* See message inplace */
#define DCGM_FR_EUD_TEST_FAILED_MSG        "" /* See message inplace */
#define DCGM_FR_FILE_CREATE_PERMISSIONS_MSG \
    "The DCGM Diagnostic does not have permissions to create a file in directory '%s'"
#define DCGM_FR_PAUSE_RESUME_FAILED_MSG "" /* See message inplace */
// gpu id
#define DCGM_FR_PCIE_H_REPLAY_VIOLATION_MSG "GPU %u host-side PCIe replay violation, see dmesg for more information"
// xid error, gpu id
#define DCGM_FR_XID_ERROR_MSG "Detected XID %u for GPU %u"
// count, field, gpu id
#define DCGM_FR_SBE_VIOLATION_MSG "Detected %ld %s for GPU %u"
// count, field, gpu id
#define DCGM_FR_DBE_VIOLATION_MSG "Detected %ld %s for GPU %u"
// count, field, gpu id
#define DCGM_FR_PCIE_REPLAY_VIOLATION_MSG "Detected %ld %s for GPU %u"
// count, field, gpu id, threshold
#define DCGM_FR_SBE_THRESHOLD_VIOLATION_MSG         "Detected %ld %s for GPU %u which is above the threshold %ld"
#define DCGM_FR_DBE_THRESHOLD_VIOLATION_MSG         "Detected %ld %s for GPU %u which is above the threshold %ld"
#define DCGM_FR_PCIE_REPLAY_THRESHOLD_VIOLATION_MSG "Detected %ld %s for GPU %u which is above the threshold %ld"
#define DCGM_FR_CUDA_FM_NOT_INITIALIZED_MSG         ""
#define DCGM_FR_SXID_ERROR_MSG                      "Detected fatal NvSwitch SXID %u"
#define DCGM_FR_ERROR_SENTINEL_MSG                  "" /* See message inplace */

/*
 * Suggestions for next steps for the corresponding error message
 */
#define DCGM_FR_OK_NEXT           "N/A"
#define DCGM_FR_UNKNOWN_NEXT      ""
#define DCGM_FR_UNRECOGNIZED_NEXT ""
#define DCGM_FR_PCI_REPLAY_RATE_NEXT                                   \
    "Reconnect PCIe card. Run system side PCIE diagnostic utilities "  \
    "to verify hops off the GPU board. If issue is on the board, run " \
    "the field diagnostic."
#define DCGM_FR_VOLATILE_DBE_DETECTED_NEXT    "Drain the GPU and reset it or reboot the node."
#define DCGM_FR_VOLATILE_SBE_DETECTED_NEXT    "Monitor - this GPU can still perform workload."
#define DCGM_FR_PENDING_PAGE_RETIREMENTS_NEXT "Monitor - this GPU can still perform workload"
#define DCGM_FR_RETIRED_PAGES_LIMIT_NEXT      TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_RETIRED_PAGES_DBE_LIMIT_NEXT  TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_CORRUPT_INFOROM_NEXT          "Flash the InfoROM to clear this corruption."
#define DCGM_FR_CLOCK_THROTTLE_THERMAL_NEXT   DEBUG_COOLING_MSG
#define DCGM_FR_POWER_UNREADABLE_NEXT         SYSTEM_TRIAGE_MSG
#define DCGM_FR_CLOCK_THROTTLE_POWER_NEXT     "Monitor the power conditions. This GPU can still perform workload."
#define DCGM_FR_NVLINK_ERROR_THRESHOLD_NEXT   "Monitor the NVLink. It can still perform workload."
#define DCGM_FR_NVLINK_DOWN_NEXT              SYSTEM_TRIAGE_MSG
#define DCGM_FR_NVSWITCH_FATAL_ERROR_NEXT     TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_NVSWITCH_NON_FATAL_ERROR_NEXT "Monitor the NVSwitch. It can still perform workload."
#define DCGM_FR_NVSWITCH_DOWN_NEXT            SYSTEM_TRIAGE_MSG
#define DCGM_FR_NO_ACCESS_TO_FILE_NEXT        "Check relevant permissions, access, and existence of the file."
#define DCGM_FR_GPU_EXPECTED_NVLINKS_UP_NEXT \
    "Ensure Fabric Manager is running. Check system logs, dmesg, and fabric-manager logs for more info."

#define DCGM_FR_NVSWITCH_EXPECTED_NVLINKS_UP_NEXT \
    "Ensure Fabric Manager is running. Check system logs, dmesg, and fabric-manager logs for more info."

#define DCGM_FR_NVML_API_NEXT                                          \
    "Check the error condition and ensure that appropriate libraries " \
    "are present and accessible."
#define DCGM_FR_DEVICE_COUNT_MISMATCH_NEXT                             \
    "Check for the presence of cgroups, operating system blocks, and " \
    "or unsupported / older cards"
#define DCGM_FR_BAD_PARAMETER_NEXT BUG_REPORT_MSG
#define DCGM_FR_CANNOT_OPEN_LIB_NEXT                                  \
    "Check for the existence of the library and set LD_LIBRARY_PATH " \
    "if needed."
#define DCGM_FR_DENYLISTED_DRIVER_NEXT "Please load the appropriate driver."
#define DCGM_FR_NVML_LIB_BAD_NEXT                             \
    "Make sure that the required version of libnvidia-ml.so " \
    "is present and accessible on the system."
#define DCGM_FR_GRAPHICS_PROCESSES_NEXT                               \
    "Stop the graphics processes or run this diagnostic on a server " \
    "that is not being used for display purposes."
#define DCGM_FR_HOSTENGINE_CONN_NEXT                                \
    "If hostengine is run separately, please ensure that it is up " \
    "and responsive."
#define DCGM_FR_FIELD_QUERY_NEXT  SYSTEM_TRIAGE_MSG
#define DCGM_FR_BAD_CUDA_ENV_NEXT "Please unset this environment variable to address test failures."
#define DCGM_FR_PERSISTENCE_MODE_NEXT                                 \
    "Enable persistence mode by running \"nvidia-smi -i <gpuId> -pm " \
    "1 \" as root."
#define DCGM_FR_LOW_BANDWIDTH_NEXT                                   \
    "Verify that your minimum bandwidth setting is appropriate for " \
    "the topology of each GPU. If so, and errors are consistent, "   \
    "please run a field diagnostic."
#define DCGM_FR_HIGH_LATENCY_NEXT                                  \
    "Verify that your maximum latency setting is appropriate for " \
    "the topology of each GPU. If so, and errors are consistent, " \
    "please run a field diagnostic."
#define DCGM_FR_CANNOT_GET_FIELD_TAG_NEXT   ""
#define DCGM_FR_FIELD_VIOLATION_NEXT        SYSTEM_TRIAGE_MSG
#define DCGM_FR_FIELD_THRESHOLD_NEXT        SYSTEM_TRIAGE_MSG
#define DCGM_FR_FIELD_VIOLATION_DBL_NEXT    SYSTEM_TRIAGE_MSG
#define DCGM_FR_FIELD_THRESHOLD_DBL_NEXT    SYSTEM_TRIAGE_MSG
#define DCGM_FR_UNSUPPORTED_FIELD_TYPE_NEXT SYSTEM_TRIAGE_MSG
#define DCGM_FR_FIELD_THRESHOLD_TS_NEXT     SYSTEM_TRIAGE_MSG
#define DCGM_FR_FIELD_THRESHOLD_TS_DBL_NEXT SYSTEM_TRIAGE_MSG
#define DCGM_FR_THERMAL_VIOLATIONS_NEXT     DEBUG_COOLING_MSG
#define DCGM_FR_THERMAL_VIOLATIONS_TS_NEXT  DEBUG_COOLING_MSG
#define DCGM_FR_TEMP_VIOLATION_NEXT                              \
    "Verify that the user-specified temperature maximum is set " \
    "correctly. If it is, check the cooling for this GPU and node: " DEBUG_COOLING_MSG
#define DCGM_FR_THROTTLING_VIOLATION_NEXT SYSTEM_TRIAGE_MSG
#define DCGM_FR_INTERNAL_NEXT             SYSTEM_TRIAGE_MSG
#define DCGM_FR_PCIE_GENERATION_NEXT      CONFIG_MSG
#define DCGM_FR_PCIE_WIDTH_NEXT           CONFIG_MSG
#define DCGM_FR_ABORTED_NEXT              ""
#define DCGM_FR_TEST_DISABLED_NEXT        CONFIG_MSG
#define DCGM_FR_CANNOT_GET_STAT_NEXT                               \
    "If running a standalone nv-hostengine, verify that it is up " \
    "and responsive."
#define DCGM_FR_STRESS_LEVEL_NEXT       SYSTEM_TRIAGE_MSG
#define DCGM_FR_CUDA_API_NEXT           SYSTEM_TRIAGE_MSG
#define DCGM_FR_FAULTY_MEMORY_NEXT      TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_CANNOT_SET_WATCHES_NEXT SYSTEM_TRIAGE_MSG
#define DCGM_FR_CUDA_UNBOUND_NEXT       SYSTEM_TRIAGE_MSG
#define DCGM_FR_ECC_DISABLED_NEXT                                  \
    "Enable ECC memory by running \"nvidia-smi -i <gpuId> -e 1\" " \
    "to enable. This may require a GPU reset or reboot to take effect."
#define DCGM_FR_MEMORY_ALLOC_NEXT    SYSTEM_TRIAGE_MSG
#define DCGM_FR_CUDA_DBE_NEXT        TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_MEMORY_MISMATCH_NEXT TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_CUDA_DEVICE_NEXT                                      \
    "Make sure CUDA_VISIBLE_DEVICES is not preventing visibility of " \
    "this GPU. Also check if CUDA libraries are compatible and "      \
    "correctly installed."
#define DCGM_FR_ECC_UNSUPPORTED_NEXT  CONFIG_MSG
#define DCGM_FR_ECC_PENDING_NEXT      "Reboot to complete activation of the ECC memory."
#define DCGM_FR_MEMORY_BANDWIDTH_NEXT SYSTEM_TRIAGE_MSG
#define DCGM_FR_TARGET_POWER_NEXT     "Verify that the clock speeds and GPU utilization are high."
#define DCGM_FR_API_FAIL_NEXT         SYSTEM_TRIAGE_MSG
#define DCGM_FR_API_FAIL_GPU_NEXT     SYSTEM_TRIAGE_MSG
#define DCGM_FR_CUDA_CONTEXT_NEXT                                   \
    "Please make sure the correct driver version is installed and " \
    "verify that no conflicting libraries are present."
#define DCGM_FR_DCGM_API_NEXT                   SYSTEM_TRIAGE_MSG
#define DCGM_FR_CONCURRENT_GPUS_NEXT            CONFIG_MSG
#define DCGM_FR_TOO_MANY_ERRORS_NEXT            ""
#define DCGM_FR_NVLINK_CRC_ERROR_THRESHOLD_NEXT TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_NVLINK_ERROR_CRITICAL_NEXT      TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_ENFORCED_POWER_LIMIT_NEXT                           \
    "If this enforced power limit is necessary, then this test "    \
    "cannot be run. If it is unnecessary, then raise the enforced " \
    "power limit setting to be able to run this test."
#define DCGM_FR_MEMORY_ALLOC_HOST_NEXT "Manually kill processes or restart your machine."
#define DCGM_FR_GPU_OP_MODE_NEXT                                     \
    "Fix by running nvidia-smi as root with: nvidia-smi --gom=0 -i " \
    "<gpu index>"
#define DCGM_FR_NO_MEMORY_CLOCKS_NEXT             ""
#define DCGM_FR_NO_GRAPHICS_CLOCKS_NEXT           ""
#define DCGM_FR_HAD_TO_RESTORE_STATE_NEXT         SYSTEM_TRIAGE_MSG
#define DCGM_FR_L1TAG_UNSUPPORTED_NEXT            CONFIG_MSG
#define DCGM_FR_L1TAG_MISCOMPARE_NEXT             TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_ROW_REMAP_FAILURE_NEXT            TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_UNCONTAINED_ERROR_NEXT            DCGM_FR_VOLATILE_DBE_DETECTED_NEXT
#define DCGM_FR_DBE_PENDING_PAGE_RETIREMENTS_NEXT "Drain the GPU and reset it or reboot the node to resolve this issue."
#define DCGM_FR_EMPTY_GPU_LIST_NEXT               CONFIG_MSG
#define DCGM_FR_UNCORRECTABLE_ROW_REMAP_NEXT      ""
#define DCGM_FR_PENDING_ROW_REMAP_NEXT            SYSTEM_TRIAGE_MSG
#define DCGM_FR_BROKEN_P2P_MEMORY_DEVICE_NEXT     BUG_REPORT_MSG
#define DCGM_FR_BROKEN_P2P_WRITER_DEVICE_NEXT     BUG_REPORT_MSG
#define DCGM_FR_NVSWITCH_NVLINK_DOWN_NEXT                                                      \
    "Please check fabric manager and initialization logs to figure out why the link is down. " \
    "You may also need to run a field diagnostic."
#define DCGM_FR_EUD_BINARY_PERMISSIONS_NEXT "" /* See message inplace */
#define DCGM_FR_EUD_NON_ROOT_USER_NEXT      "" /* See message inplace */
#define DCGM_FR_EUD_SPAWN_FAILURE_NEXT      "" /* See message inplace */
#define DCGM_FR_EUD_TIMEOUT_NEXT            "" /* See message inplace */
#define DCGM_FR_EUD_ZOMBIE_NEXT             "" /* See message inplace */
#define DCGM_FR_EUD_NON_ZERO_EXIT_CODE_NEXT "" /* See message inplace */
#define DCGM_FR_EUD_TEST_FAILED_NEXT        "" /* See message inplace */
#define DCGM_FR_FILE_CREATE_PERMISSIONS_NEXT                                                                 \
    "Please restart the hostengine with parameter --home-dir to specify a different home directory for the " \
    "diagnostic or change permissions in the current directory to allow the user to write files there."
#define DCGM_FR_PAUSE_RESUME_FAILED_NEXT             "" /* See message inplace */
#define DCGM_FR_PCIE_H_REPLAY_VIOLATION_NEXT         "" /* See message inplace */
#define DCGM_FR_XID_ERROR_NEXT                       "Please consult the documentation for details of this XID."
#define DCGM_FR_SBE_VIOLATION_NEXT                   TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_DBE_VIOLATION_NEXT                   TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_PCIE_REPLAY_VIOLATION_NEXT           TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_SBE_THRESHOLD_VIOLATION_NEXT         TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_DBE_THRESHOLD_VIOLATION_NEXT         TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_PCIE_REPLAY_THRESHOLD_VIOLATION_NEXT TRIAGE_RUN_FIELD_DIAG_MSG
#define DCGM_FR_CUDA_FM_NOT_INITIALIZED_NEXT         "Ensure that the FabricManager is running without errors."
#define DCGM_FR_SXID_ERROR_NEXT                      SYSTEM_TRIAGE_MSG
#define DCGM_FR_ERROR_SENTINEL_NEXT                  "" /* See message inplace */

#ifdef __cplusplus
extern "C" {
#endif

DCGM_PUBLIC_API dcgmErrorSeverity_t dcgmErrorGetPriorityByCode(unsigned int code);
DCGM_PUBLIC_API dcgmErrorCategory_t dcgmErrorGetCategoryByCode(unsigned int code);
DCGM_PUBLIC_API const char *dcgmErrorGetFormatMsgByCode(unsigned int code);

DCGM_PUBLIC_API const dcgm_error_meta_t *dcgmGetErrorMeta(dcgmError_t error);
DCGM_PUBLIC_API const char *errorString(dcgmReturn_t result);

/** @} */

#ifdef __cplusplus
} // extern "C"
#endif

#endif // DCGM_ERRORS_H
