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

#ifndef DCGMFIELDS_H
#define DCGMFIELDS_H

#ifdef __cplusplus
extern "C" {
#endif

#define DCGM_PUBLIC_API

/***************************************************************************************************/
/** @defgroup dcgmFieldTypes Field Types
 *  Field Types are a single byte.
 *  @{
 */
/***************************************************************************************************/

/**
 * Blob of binary data representing a structure
 */
#define DCGM_FT_BINARY 'b'

/**
 * 8-byte double precision
 */
#define DCGM_FT_DOUBLE 'd'

/**
 * 8-byte signed integer
 */
#define DCGM_FT_INT64 'i'

/**
 * Null-terminated ASCII Character string
 */
#define DCGM_FT_STRING 's'

/**
 * 8-byte signed integer usec since 1970
 */
#define DCGM_FT_TIMESTAMP 't'

/** @} */


/***************************************************************************************************/
/** @defgroup dcgmFieldScope Field Scope
 *  Represents field association with entity scope or global scope.
 *  @{
 */
/***************************************************************************************************/

/**
 * Field is global (ex: driver version)
 */
#define DCGM_FS_GLOBAL 0

/**
 * Field is associated with an entity (GPU, VGPU...etc)
 */
#define DCGM_FS_ENTITY 1

/**
 * Field is associated with a device. Deprecated. Use DCGM_FS_ENTITY
 */
#define DCGM_FS_DEVICE DCGM_FS_ENTITY

/** @} */

/***************************************************************************************************/
/** @defgroup dcgmFieldConstants Field Constants
 *  Constants that represent contents of individual field values.
 *  @{
 */
/***************************************************************************************************/

/**
 * DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY is 16 bits of major version followed by
 * 16 bits of the minor version. These macros separate the two.
 */
#define DCGM_CUDA_COMPUTE_CAPABILITY_MAJOR(x) ((uint64_t)(x)&0xFFFF0000)
#define DCGM_CUDA_COMPUTE_CAPABILITY_MINOR(x) ((uint64_t)(x)&0x0000FFFF)

/**
 * DCGM_FI_DEV_CLOCK_THROTTLE_REASONS is a bitmap of why the clock is throttled.
 * These macros are masks for relevant throttling, and are a 1:1 map to the NVML
 * reasons documented in nvml.h. The notes for the header are copied blow:
 */
/** Nothing is running on the GPU and the clocks are dropping to Idle state
 * \note This limiter may be removed in a later release
 */
#define DCGM_CLOCKS_THROTTLE_REASON_GPU_IDLE 0x0000000000000001LL
/** GPU clocks are limited by current setting of applications clocks
 */
#define DCGM_CLOCKS_THROTTLE_REASON_CLOCKS_SETTING 0x0000000000000002LL
/** SW Power Scaling algorithm is reducing the clocks below requested clocks
 */
#define DCGM_CLOCKS_THROTTLE_REASON_SW_POWER_CAP 0x0000000000000004LL
/** HW Slowdown (reducing the core clocks by a factor of 2 or more) is engaged
 *
 * This is an indicator of:
 *  - temperature being too high
 *  - External Power Brake Assertion is triggered (e.g. by the system power supply)
 *  - Power draw is too high and Fast Trigger protection is reducing the clocks
 *  - May be also reported during PState or clock change
 *  - This behavior may be removed in a later release.
 */
#define DCGM_CLOCKS_THROTTLE_REASON_HW_SLOWDOWN 0x0000000000000008LL
/** Sync Boost
 *
 * This GPU has been added to a Sync boost group with nvidia-smi or DCGM in
 * order to maximize performance per watt. All GPUs in the sync boost group
 * will boost to the minimum possible clocks across the entire group. Look at
 * the throttle reasons for other GPUs in the system to see why those GPUs are
 * holding this one at lower clocks.
 */
#define DCGM_CLOCKS_THROTTLE_REASON_SYNC_BOOST 0x0000000000000010LL
/** SW Thermal Slowdown
 *
 * This is an indicator of one or more of the following:
 *  - Current GPU temperature above the GPU Max Operating Temperature
 *  - Current memory temperature above the Memory Max Operating Temperature
 */
#define DCGM_CLOCKS_THROTTLE_REASON_SW_THERMAL 0x0000000000000020LL
/** HW Thermal Slowdown (reducing the core clocks by a factor of 2 or more) is engaged
 *
 * This is an indicator of:
 *  - temperature being too high
 */
#define DCGM_CLOCKS_THROTTLE_REASON_HW_THERMAL 0x0000000000000040LL
/** HW Power Brake Slowdown (reducing the core clocks by a factor of 2 or more) is engaged
 *
 * This is an indicator of:
 *  - External Power Brake Assertion being triggered (e.g. by the system power supply)
 */
#define DCGM_CLOCKS_THROTTLE_REASON_HW_POWER_BRAKE 0x0000000000000080LL
/** GPU clocks are limited by current setting of Display clocks
 */
#define DCGM_CLOCKS_THROTTLE_REASON_DISPLAY_CLOCKS 0x0000000000000100LL

/**
 * GPU virtualization mode types for DCGM_FI_DEV_VIRTUAL_MODE
 */
typedef enum
{
    DCGM_GPU_VIRTUALIZATION_MODE_NONE        = 0, //!< Represents Bare Metal GPU
    DCGM_GPU_VIRTUALIZATION_MODE_PASSTHROUGH = 1, //!< Device is associated with GPU-Passthrough
    DCGM_GPU_VIRTUALIZATION_MODE_VGPU        = 2, //!< Device is associated with vGPU inside virtual machine.
    DCGM_GPU_VIRTUALIZATION_MODE_HOST_VGPU   = 3, //!< Device is associated with VGX hypervisor in vGPU mode
    DCGM_GPU_VIRTUALIZATION_MODE_HOST_VSGA   = 4, //!< Device is associated with VGX hypervisor in vSGA mode
} dcgmGpuVirtualizationMode_t;


/** @} */

/***************************************************************************************************/
/** @defgroup dcgmFieldEntity Field Entity
 *  Represents field association with a particular entity
 *  @{
 */
/***************************************************************************************************/

/**
 * Enum of possible field entity groups
 */
typedef enum dcgm_field_entity_group_t
{
    DCGM_FE_NONE = 0, /*!< Field is not associated with an entity. Field scope should be DCGM_FS_GLOBAL */
    DCGM_FE_GPU,      /*!< Field is associated with a GPU entity */
    DCGM_FE_VGPU,     /*!< Field is associated with a VGPU entity */
    DCGM_FE_SWITCH,   /*!< Field is associated with a Switch entity */
    DCGM_FE_GPU_I,    /*!< Field is associated with a GPU Instance entity */
    DCGM_FE_GPU_CI,   /*!< Field is associated with a GPU Compute Instance entity */
    DCGM_FE_LINK,     /*!< Field is associated with an NVLink */
    DCGM_FE_CPU,      /*!< Field is associated with a CPU node */
    DCGM_FE_CPU_CORE, /*!< Field is associated with a CPU */

    DCGM_FE_COUNT /*!< Number of elements in this enumeration. Keep this entry last */
} dcgm_field_entity_group_t;

/**
 * Represents an identifier for an entity within a field entity. For instance, this is the gpuId for DCGM_FE_GPU.
 */
typedef unsigned int dcgm_field_eid_t;


/** @} */

/***************************************************************************************************/
/** @defgroup dcgmFieldIdentifiers Field Identifiers
 *  Field Identifiers
 *  @{
 */
/***************************************************************************************************/

/**
 * NULL field
 */
#define DCGM_FI_UNKNOWN 0

/**
 * Driver Version
 */
#define DCGM_FI_DRIVER_VERSION 1

/* Underlying NVML version */
#define DCGM_FI_NVML_VERSION 2

/*
 * Process Name
 */
#define DCGM_FI_PROCESS_NAME 3

/**
 * Number of Devices on the node
 */
#define DCGM_FI_DEV_COUNT 4

/**
 * Cuda Driver Version
 * Retrieves a number with the major value in the thousands place and the minor value in the hundreds place.
 * CUDA 11.1 = 11100
 */
#define DCGM_FI_CUDA_DRIVER_VERSION 5


/**
 * Name of the GPU device
 */
#define DCGM_FI_DEV_NAME 50

/**
 * Device Brand
 */
#define DCGM_FI_DEV_BRAND 51

/**
 * NVML index of this GPU
 */
#define DCGM_FI_DEV_NVML_INDEX 52

/**
 * Device Serial Number
 */
#define DCGM_FI_DEV_SERIAL 53

/**
 * UUID corresponding to the device
 */
#define DCGM_FI_DEV_UUID 54

/**
 * Device node minor number /dev/nvidia#
 */
#define DCGM_FI_DEV_MINOR_NUMBER 55

/**
 * OEM inforom version
 */
#define DCGM_FI_DEV_OEM_INFOROM_VER 56

/**
 * PCI attributes for the device
 */
#define DCGM_FI_DEV_PCI_BUSID 57

/**
 * The combined 16-bit device id and 16-bit vendor id
 */
#define DCGM_FI_DEV_PCI_COMBINED_ID 58

/**
 * The 32-bit Sub System Device ID
 */
#define DCGM_FI_DEV_PCI_SUBSYS_ID 59

/**
 * Topology of all GPUs on the system via PCI (static)
 */
#define DCGM_FI_GPU_TOPOLOGY_PCI 60

/**
 * Topology of all GPUs on the system via NVLINK (static)
 */
#define DCGM_FI_GPU_TOPOLOGY_NVLINK 61

/**
 * Affinity of all GPUs on the system (static)
 */
#define DCGM_FI_GPU_TOPOLOGY_AFFINITY 62

/**
 * Cuda compute capability for the device.
 * The major version is the upper 32 bits and
 * the minor version is the lower 32 bits.
 */
#define DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY 63

/**
 * Compute mode for the device
 */
#define DCGM_FI_DEV_COMPUTE_MODE 65

/**
 * Persistence mode for the device
 * Boolean: 0 is disabled, 1 is enabled
 */
#define DCGM_FI_DEV_PERSISTENCE_MODE 66

/**
 * MIG mode for the device
 * Boolean: 0 is disabled, 1 is enabled
 */
#define DCGM_FI_DEV_MIG_MODE 67

/**
 * The string that CUDA_VISIBLE_DEVICES should
 * be set to for this entity (including MIG)
 */
#define DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR 68

/**
 * The maximum number of MIG slices supported by this GPU
 */
#define DCGM_FI_DEV_MIG_MAX_SLICES 69

/**
 * Device CPU affinity. part 1/8 = cpus 0 - 63
 */
#define DCGM_FI_DEV_CPU_AFFINITY_0 70

/**
 * Device CPU affinity. part 1/8 = cpus 64 - 127
 */
#define DCGM_FI_DEV_CPU_AFFINITY_1 71

/**
 * Device CPU affinity. part 2/8 = cpus 128 - 191
 */
#define DCGM_FI_DEV_CPU_AFFINITY_2 72

/**
 * Device CPU affinity. part 3/8 = cpus 192 - 255
 */
#define DCGM_FI_DEV_CPU_AFFINITY_3 73

/**
 * ConfidentialCompute/AmpereProtectedMemory status for this system
 * 0 = disabled
 * 1 = enabled
 */
#define DCGM_FI_DEV_CC_MODE 74

/**
 * Attributes for the given MIG device handles
 */
#define DCGM_FI_DEV_MIG_ATTRIBUTES 75

/**
 * GPU instance profile information
 */
#define DCGM_FI_DEV_MIG_GI_INFO 76

/**
 * Compute instance profile information
 */
#define DCGM_FI_DEV_MIG_CI_INFO 77

/**
 * ECC inforom version
 */
#define DCGM_FI_DEV_ECC_INFOROM_VER 80

/**
 * Power management object inforom version
 */
#define DCGM_FI_DEV_POWER_INFOROM_VER 81

/**
 * Inforom image version
 */
#define DCGM_FI_DEV_INFOROM_IMAGE_VER 82

/**
 * Inforom configuration checksum
 */
#define DCGM_FI_DEV_INFOROM_CONFIG_CHECK 83

/**
 * Reads the infoROM from the flash and verifies the checksums
 */
#define DCGM_FI_DEV_INFOROM_CONFIG_VALID 84

/**
 * VBIOS version of the device
 */
#define DCGM_FI_DEV_VBIOS_VERSION 85

/**
 * Device Memory node affinity, 0-63
 */
#define DCGM_FI_DEV_MEM_AFFINITY_0 86

/**
 * Device Memory node affinity, 64-127
 */
#define DCGM_FI_DEV_MEM_AFFINITY_1 87

/**
 * Device Memory node affinity, 128-191
 */
#define DCGM_FI_DEV_MEM_AFFINITY_2 88

/**
 * Device Memory node affinity, 192-255
 */
#define DCGM_FI_DEV_MEM_AFFINITY_3 89

/**
 * Total BAR1 of the GPU in MB
 */
#define DCGM_FI_DEV_BAR1_TOTAL 90

/**
 * Deprecated - Sync boost settings on the node
 */
#define DCGM_FI_SYNC_BOOST 91

/**
 * Used BAR1 of the GPU in MB
 */
#define DCGM_FI_DEV_BAR1_USED 92

/**
 * Free BAR1 of the GPU in MB
 */
#define DCGM_FI_DEV_BAR1_FREE 93

/**
 * SM clock for the device
 */
#define DCGM_FI_DEV_SM_CLOCK 100

/**
 * Memory clock for the device
 */
#define DCGM_FI_DEV_MEM_CLOCK 101

/**
 * Video encoder/decoder clock for the device
 */
#define DCGM_FI_DEV_VIDEO_CLOCK 102

/**
 * SM Application clocks
 */
#define DCGM_FI_DEV_APP_SM_CLOCK 110

/**
 * Memory Application clocks
 */
#define DCGM_FI_DEV_APP_MEM_CLOCK 111

/**
 * Current clock throttle reasons (bitmask of DCGM_CLOCKS_THROTTLE_REASON_*)
 */
#define DCGM_FI_DEV_CLOCK_THROTTLE_REASONS 112

/**
 * Maximum supported SM clock for the device
 */
#define DCGM_FI_DEV_MAX_SM_CLOCK 113

/**
 * Maximum supported Memory clock for the device
 */
#define DCGM_FI_DEV_MAX_MEM_CLOCK 114

/**
 * Maximum supported Video encoder/decoder clock for the device
 */
#define DCGM_FI_DEV_MAX_VIDEO_CLOCK 115

/**
 * Auto-boost for the device (1 = enabled. 0 = disabled)
 */
#define DCGM_FI_DEV_AUTOBOOST 120

/**
 * Supported clocks for the device
 */
#define DCGM_FI_DEV_SUPPORTED_CLOCKS 130

/**
 * Memory temperature for the device
 */
#define DCGM_FI_DEV_MEMORY_TEMP 140

/**
 * Current temperature readings for the device, in degrees C
 */
#define DCGM_FI_DEV_GPU_TEMP 150

/**
 * Maximum operating temperature for the memory of this GPU
 */
#define DCGM_FI_DEV_MEM_MAX_OP_TEMP 151

/**
 * Maximum operating temperature for this GPU
 */
#define DCGM_FI_DEV_GPU_MAX_OP_TEMP 152


/**
 * Power usage for the device in Watts
 */
#define DCGM_FI_DEV_POWER_USAGE 155

/**
 * Total energy consumption for the GPU in mJ since the driver was last reloaded
 */
#define DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION 156

/**
 * Current instantaneous power usage of the device in Watts
 */
#define DCGM_FI_DEV_POWER_USAGE_INSTANT 157

/**
 * Slowdown temperature for the device
 */
#define DCGM_FI_DEV_SLOWDOWN_TEMP 158

/**
 * Shutdown temperature for the device
 */
#define DCGM_FI_DEV_SHUTDOWN_TEMP 159

/**
 * Current Power limit for the device
 */
#define DCGM_FI_DEV_POWER_MGMT_LIMIT 160

/**
 * Minimum power management limit for the device
 */
#define DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN 161

/**
 * Maximum power management limit for the device
 */
#define DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX 162

/**
 * Default power management limit for the device
 */
#define DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF 163

/**
 * Effective power limit that the driver enforces after taking into account all limiters
 */
#define DCGM_FI_DEV_ENFORCED_POWER_LIMIT 164

/**
 * Performance state (P-State) 0-15. 0=highest
 */
#define DCGM_FI_DEV_PSTATE 190

/**
 * Fan speed for the device in percent 0-100
 */
#define DCGM_FI_DEV_FAN_SPEED 191

/**
 * PCIe Tx utilization information
 *
 * Deprecated: Use DCGM_FI_PROF_PCIE_TX_BYTES instead.
 */
#define DCGM_FI_DEV_PCIE_TX_THROUGHPUT 200

/**
 * PCIe Rx utilization information
 *
 * Deprecated: Use DCGM_FI_PROF_PCIE_RX_BYTES instead.
 */
#define DCGM_FI_DEV_PCIE_RX_THROUGHPUT 201

/**
 * PCIe replay counter
 */
#define DCGM_FI_DEV_PCIE_REPLAY_COUNTER 202

/**
 * GPU Utilization
 */
#define DCGM_FI_DEV_GPU_UTIL 203

/**
 * Memory Utilization
 */
#define DCGM_FI_DEV_MEM_COPY_UTIL 204

/**
 * Process accounting stats.
 *
 * This field is only supported when the host engine is running as root unless you
 * enable accounting ahead of time. Accounting mode can be enabled by
 * running "nvidia-smi -am 1" as root on the same node the host engine is running on.
 */
#define DCGM_FI_DEV_ACCOUNTING_DATA 205

/**
 * Encoder Utilization
 */
#define DCGM_FI_DEV_ENC_UTIL 206

/**
 * Decoder Utilization
 */
#define DCGM_FI_DEV_DEC_UTIL 207

/* Fields 210, 211, 220 and 221 are internal-only. See dcgm_fields_internal.hpp */

/**
 * XID errors. The value is the specific XID error
 */
#define DCGM_FI_DEV_XID_ERRORS 230

/**
 * PCIe Max Link Generation
 */
#define DCGM_FI_DEV_PCIE_MAX_LINK_GEN 235

/**
 * PCIe Max Link Width
 */
#define DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH 236

/**
 * PCIe Current Link Generation
 */
#define DCGM_FI_DEV_PCIE_LINK_GEN 237

/**
 * PCIe Current Link Width
 */
#define DCGM_FI_DEV_PCIE_LINK_WIDTH 238

/**
 * Power Violation time in usec
 */
#define DCGM_FI_DEV_POWER_VIOLATION 240

/**
 * Thermal Violation time in usec
 */
#define DCGM_FI_DEV_THERMAL_VIOLATION 241

/**
 * Sync Boost Violation time in usec
 */
#define DCGM_FI_DEV_SYNC_BOOST_VIOLATION 242

/**
 * Board violation limit.
 */
#define DCGM_FI_DEV_BOARD_LIMIT_VIOLATION 243

/**
 *Low utilisation violation limit.
 */
#define DCGM_FI_DEV_LOW_UTIL_VIOLATION 244

/**
 *Reliability violation limit.
 */
#define DCGM_FI_DEV_RELIABILITY_VIOLATION 245

/**
 * App clock violation limit.
 */
#define DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION 246

/**
 * Base clock violation limit.
 */
#define DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION 247

/**
 * Total Frame Buffer of the GPU in MB
 */
#define DCGM_FI_DEV_FB_TOTAL 250

/**
 * Free Frame Buffer in MB
 */
#define DCGM_FI_DEV_FB_FREE 251

/**
 * Used Frame Buffer in MB
 */
#define DCGM_FI_DEV_FB_USED 252

/**
 * Reserved Frame Buffer in MB
 */
#define DCGM_FI_DEV_FB_RESERVED 253

/**
 * Percentage used of Frame Buffer: 'Used/(Total - Reserved)'. Range 0.0-1.0
 */
#define DCGM_FI_DEV_FB_USED_PERCENT 254

/**
 * C2C Link Count
 */
#define DCGM_FI_DEV_C2C_LINK_COUNT 285

/**
 * C2C Link Status
 * The value of 0 the link is INACTIVE.
 * The value of 1 the link is ACTIVE.
 */
#define DCGM_FI_DEV_C2C_LINK_STATUS 286

/**
 * C2C Max Bandwidth
 * The value indicates the link speed in MB/s.
 */
#define DCGM_FI_DEV_C2C_MAX_BANDWIDTH 287

/**
 * Current ECC mode for the device
 */
#define DCGM_FI_DEV_ECC_CURRENT 300

/**
 * Pending ECC mode for the device
 */
#define DCGM_FI_DEV_ECC_PENDING 301

/**
 * Total single bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_SBE_VOL_TOTAL 310

/**
 * Total double bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_DBE_VOL_TOTAL 311

/**
 * Total single bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_SBE_AGG_TOTAL 312

/**
 * Total double bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_DBE_AGG_TOTAL 313

/**
 * L1 cache single bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_SBE_VOL_L1 314

/**
 * L1 cache double bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_DBE_VOL_L1 315

/**
 * L2 cache single bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_SBE_VOL_L2 316

/**
 * L2 cache double bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_DBE_VOL_L2 317

/**
 * Device memory single bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_SBE_VOL_DEV 318

/**
 * Device memory double bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_DBE_VOL_DEV 319

/**
 * Register file single bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_SBE_VOL_REG 320

/**
 * Register file double bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_DBE_VOL_REG 321

/**
 * Texture memory single bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_SBE_VOL_TEX 322

/**
 * Texture memory double bit volatile ECC errors
 */
#define DCGM_FI_DEV_ECC_DBE_VOL_TEX 323

/**
 * L1 cache single bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_SBE_AGG_L1 324

/**
 * L1 cache double bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_DBE_AGG_L1 325

/**
 * L2 cache single bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_SBE_AGG_L2 326

/**
 * L2 cache double bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_DBE_AGG_L2 327

/**
 * Device memory single bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_SBE_AGG_DEV 328

/**
 * Device memory double bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_DBE_AGG_DEV 329

/**
 * Register File single bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_SBE_AGG_REG 330

/**
 * Register File double bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_DBE_AGG_REG 331

/**
 * Texture memory single bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_SBE_AGG_TEX 332

/**
 * Texture memory double bit aggregate (persistent) ECC errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_ECC_DBE_AGG_TEX 333

/**
 * Historical max available spare memory rows per memory bank
 */
#define DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX 385

/**
 * Historical high mark of available spare memory rows per memory bank
 */
#define DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH 386

/**
 * Historical mark of partial available spare memory rows per memory bank
 */
#define DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL 387

/**
 * Historical low mark of available spare memory rows per memory bank
 */
#define DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW 388

/**
 * Historical marker of memory banks with no available spare memory rows
 */
#define DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE 389

/**
 * Number of retired pages because of single bit errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_RETIRED_SBE 390

/**
 * Number of retired pages because of double bit errors
 * Note: monotonically increasing
 */
#define DCGM_FI_DEV_RETIRED_DBE 391

/**
 * Number of pages pending retirement
 */
#define DCGM_FI_DEV_RETIRED_PENDING 392

/**
 * Number of remapped rows for uncorrectable errors
 */
#define DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS 393

/**
 * Number of remapped rows for correctable errors
 */
#define DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS 394

/**
 * Whether remapping of rows has failed
 */
#define DCGM_FI_DEV_ROW_REMAP_FAILURE 395

/**
 * Whether remapping of rows is pending
 */
#define DCGM_FI_DEV_ROW_REMAP_PENDING 396

/*
 * NV Link flow control CRC  Error Counter for Lane 0
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0 400

/*
 * NV Link flow control CRC  Error Counter for Lane 1
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1 401

/*
 * NV Link flow control CRC  Error Counter for Lane 2
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2 402

/*
 * NV Link flow control CRC  Error Counter for Lane 3
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3 403

/*
 * NV Link flow control CRC  Error Counter for Lane 4
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4 404

/*
 * NV Link flow control CRC  Error Counter for Lane 5
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5 405

/*
 * NV Link flow control CRC  Error Counter total for all Lanes
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL 409

/*
 * NV Link data CRC Error Counter for Lane 0
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0 410

/*
 * NV Link data CRC Error Counter for Lane 1
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1 411

/*
 * NV Link data CRC Error Counter for Lane 2
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2 412

/*
 * NV Link data CRC Error Counter for Lane 3
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3 413

/*
 * NV Link data CRC Error Counter for Lane 4
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4 414

/*
 * NV Link data CRC Error Counter for Lane 5
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5 415

/*
 * NV Link data CRC Error Counter total for all Lanes
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL 419

/*
 * NV Link Replay Error Counter for Lane 0
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0 420

/*
 * NV Link Replay Error Counter for Lane 1
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1 421

/*
 * NV Link Replay Error Counter for Lane 2
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2 422

/*
 * NV Link Replay Error Counter for Lane 3
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3 423

/*
 * NV Link Replay Error Counter for Lane 4
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4 424

/*
 * NV Link Replay Error Counter for Lane 5
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5 425

/*
 * NV Link Replay Error Counter total for all Lanes
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL 429

/*
 * NV Link Recovery Error Counter for Lane 0
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0 430

/*
 * NV Link Recovery Error Counter for Lane 1
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1 431

/*
 * NV Link Recovery Error Counter for Lane 2
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2 432

/*
 * NV Link Recovery Error Counter for Lane 3
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3 433

/*
 * NV Link Recovery Error Counter for Lane 4
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4 434

/*
 * NV Link Recovery Error Counter for Lane 5
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5 435

/*
 * NV Link Recovery Error Counter total for all Lanes
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL 439

/*
 * NV Link Bandwidth Counter for Lane 0
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L0 440

/*
 * NV Link Bandwidth Counter for Lane 1
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L1 441

/*
 * NV Link Bandwidth Counter for Lane 2
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L2 442

/*
 * NV Link Bandwidth Counter for Lane 3
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L3 443

/*
 * NV Link Bandwidth Counter for Lane 4
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L4 444

/*
 * NV Link Bandwidth Counter for Lane 5
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L5 445

/*
 * NV Link Bandwidth Counter total for all Lanes
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL 449

/*
 * GPU NVLink error information
 */
#define DCGM_FI_DEV_GPU_NVLINK_ERRORS 450

/*
 * NV Link flow control CRC  Error Counter for Lane 6
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6 451

/*
 * NV Link flow control CRC  Error Counter for Lane 7
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7 452

/*
 * NV Link flow control CRC  Error Counter for Lane 8
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8 453

/*
 * NV Link flow control CRC  Error Counter for Lane 9
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9 454

/*
 * NV Link flow control CRC  Error Counter for Lane 10
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10 455

/*
 * NV Link flow control CRC  Error Counter for Lane 11
 */
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11 456

/*
 * NV Link data CRC Error Counter for Lane 6
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6 457

/*
 * NV Link data CRC Error Counter for Lane 7
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7 458

/*
 * NV Link data CRC Error Counter for Lane 8
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8 459

/*
 * NV Link data CRC Error Counter for Lane 9
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9 460

/*
 * NV Link data CRC Error Counter for Lane 10
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10 461

/*
 * NV Link data CRC Error Counter for Lane 11
 */
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11 462

/*
 * NV Link Replay Error Counter for Lane 6
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6 463

/*
 * NV Link Replay Error Counter for Lane 7
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7 464

/*
 * NV Link Replay Error Counter for Lane 8
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8 465

/*
 * NV Link Replay Error Counter for Lane 9
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9 466

/*
 * NV Link Replay Error Counter for Lane 10
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10 467

/*
 * NV Link Replay Error Counter for Lane 11
 */
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11 468

/*
 * NV Link Recovery Error Counter for Lane 6
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6 469

/*
 * NV Link Recovery Error Counter for Lane 7
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7 470

/*
 * NV Link Recovery Error Counter for Lane 8
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8 471

/*
 * NV Link Recovery Error Counter for Lane 9
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9 472

/*
 * NV Link Recovery Error Counter for Lane 10
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10 473

/*
 * NV Link Recovery Error Counter for Lane 11
 */
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11 474

/*
 * NV Link Bandwidth Counter for Lane 6
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L6 475

/*
 * NV Link Bandwidth Counter for Lane 7
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L7 476

/*
 * NV Link Bandwidth Counter for Lane 8
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L8 477

/*
 * NV Link Bandwidth Counter for Lane 9
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L9 478

/*
 * NV Link Bandwidth Counter for Lane 10
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L10 479

/*
 * NV Link Bandwidth Counter for Lane 11
 */
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L11 480

#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12 406
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13 407
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14 408
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15 481
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16 482
#define DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17 483

#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12 416
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13 417
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14 418
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15 484
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16 485
#define DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17 486

#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12 426
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13 427
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14 428
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15 487
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16 488
#define DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17 489

#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12 436
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13 437
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14 438
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15 491
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16 492
#define DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17 493

#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L12 446
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L13 447
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L14 448
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L15 494
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L16 495
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_L17 496

/**
 * Virtualization Mode corresponding to the GPU.
 *
 * One of DCGM_GPU_VIRTUALIZATION_MODE_* constants.
 */
#define DCGM_FI_DEV_VIRTUAL_MODE 500

/**
 * Includes Count and Static info of vGPU types supported on a device
 */
#define DCGM_FI_DEV_SUPPORTED_TYPE_INFO 501

/**
 * Includes Count and currently Creatable vGPU types on a device
 */
#define DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS 502

/**
 * Includes Count and currently Active vGPU Instances on a device
 */
#define DCGM_FI_DEV_VGPU_INSTANCE_IDS 503

/**
 * Utilization values for vGPUs running on the device
 */
#define DCGM_FI_DEV_VGPU_UTILIZATIONS 504

/**
 * Utilization values for processes running within vGPU VMs using the device
 */
#define DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION 505

/**
 * Current encoder statistics for a given device
 */
#define DCGM_FI_DEV_ENC_STATS 506

/**
 * Statistics of current active frame buffer capture sessions on a given device
 */
#define DCGM_FI_DEV_FBC_STATS 507

/**
 * Information about active frame buffer capture sessions on a target device
 */
#define DCGM_FI_DEV_FBC_SESSIONS_INFO 508

/**
 * Includes Count and currently Supported vGPU types on a device
 */
#define DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS 509

/**
 * Includes Static info of vGPU types supported on a device
 */
#define DCGM_FI_DEV_VGPU_TYPE_INFO 510

/**
 * Includes the name of a vGPU type supported on a device
 */
#define DCGM_FI_DEV_VGPU_TYPE_NAME 511

/**
 * Includes the class of a vGPU type supported on a device
 */
#define DCGM_FI_DEV_VGPU_TYPE_CLASS 512

/**
 * Includes the license info for a vGPU type supported on a device
 */
#define DCGM_FI_DEV_VGPU_TYPE_LICENSE 513

/**
 * VM ID of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_VM_ID 520

/**
 * VM name of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_VM_NAME 521

/**
 * vGPU type of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_TYPE 522

/**
 * UUID of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_UUID 523

/**
 * Driver version of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_DRIVER_VERSION 524

/**
 * Memory usage of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_MEMORY_USAGE 525

/**
 * License status of the vGPU
 */
#define DCGM_FI_DEV_VGPU_LICENSE_STATUS 526

/**
 * Frame rate limit of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT 527

/**
 * Current encoder statistics of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_ENC_STATS 528

/**
 * Information about all active encoder sessions on the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO 529

/**
 * Statistics of current active frame buffer capture sessions on the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_FBC_STATS 530

/**
 * Information about active frame buffer capture sessions on the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO 531

/**
 * License state information of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE 532

/**
 * PCI Id of the vGPU instance
 */
#define DCGM_FI_DEV_VGPU_PCI_ID 533

/**
 * GPU Instance ID for the given vGPU Instance
 */
#define DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID 534

/**
 * Starting field ID of the vGPU instance
 */
#define DCGM_FI_FIRST_VGPU_FIELD_ID 520

/**
 * Last field ID of the vGPU instance
 */
#define DCGM_FI_LAST_VGPU_FIELD_ID 570

/**
 * For now max vGPU field Ids taken as difference of DCGM_FI_LAST_VGPU_FIELD_ID and DCGM_FI_LAST_VGPU_FIELD_ID i.e. 50
 */
#define DCGM_FI_MAX_VGPU_FIELDS DCGM_FI_LAST_VGPU_FIELD_ID - DCGM_FI_FIRST_VGPU_FIELD_ID

/**
 * Starting ID for all the internal fields
 */
#define DCGM_FI_INTERNAL_FIELDS_0_START 600

/**
 * Last ID for all the internal fields
 */

/**
 * <p>&nbsp;</p>
 * <p>&nbsp;</p>
 * <p>&nbsp;</p>
 * <p>NVSwitch entity field IDs start here.</p>
 * <p>&nbsp;</p>
 * <p>&nbsp;</p>
 * <p>NVSwitch latency bins for port 0</p>
 */

#define DCGM_FI_INTERNAL_FIELDS_0_END 699

/**
 * Starting field ID of the NVSwitch instance
 */
#define DCGM_FI_FIRST_NVSWITCH_FIELD_ID 700

/**
 * NvSwitch voltage
 */
#define DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT 701

/**
 * NvSwitch Current IDDQ
 */
#define DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ 702

/**
 * NvSwitch Current IDDQ Rev
 */
#define DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV 703

/**
 * NvSwitch Current IDDQ Rev DVDD
 */
#define DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD 704

/**
 * NvSwitch Power VDD in watts
 */
#define DCGM_FI_DEV_NVSWITCH_POWER_VDD 705

/**
 * NvSwitch Power DVDD in watts
 */
#define DCGM_FI_DEV_NVSWITCH_POWER_DVDD 706

/**
 * NvSwitch Power HVDD in watts
 */
#define DCGM_FI_DEV_NVSWITCH_POWER_HVDD 707

/**
 * <p>NVSwitch Tx Throughput Counter for ports 0-17</p>
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX 780
/**
 * NVSwitch Rx Throughput Counter for ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX 781

/**
 * NvSwitch fatal_errors for ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS 782

/**
 * NvSwitch non_fatal_errors for ports 0-17
 *
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS 783

/**
 * NvSwitch replay_count_errors for ports  0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS 784

/**
 * NvSwitch recovery_count_errors for ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS 785

/**
 * NvSwitch filt_err_count_errors for ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS 786

/**
 * NvLink lane_crs_err_count_aggregate_errors for ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS 787

/**
 * NvLink lane ecc_err_count_aggregate_errors for ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS 788

/**
 * Nvlink lane latency low lane0 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0 789

/**
 * Nvlink lane latency low lane1 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1 790

/**
 * Nvlink lane latency low lane2 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2 791

/**
 * Nvlink lane latency low lane3 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3 792

/**
 * Nvlink lane latency medium lane0 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0 793

/**
 * Nvlink lane latency medium lane1 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1 794

/**
 * Nvlink lane latency medium lane2 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2 795

/**
 * Nvlink lane latency medium lane3 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3 796

/**
 * Nvlink lane latency high lane0 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0 797

/**
 * Nvlink lane latency high lane1 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1 798

/**
 * Nvlink lane latency high lane2 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2 799

/**
 * Nvlink lane latency high lane3 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3 800

/**
 * Nvlink lane latency panic lane0 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0 801

/**
 * Nvlink lane latency panic lane1 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1 802

/**
 * Nvlink lane latency panic lane2 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2 803

/**
 * Nvlink lane latency panic lane2 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3 804

/**
 * Nvlink lane latency count lane0 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0 805

/**
 * Nvlink lane latency count lane1 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1 806

/**
 * Nvlink lane latency count lane2 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2 807

/**
 * Nvlink lane latency count lane3 counter.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3 808

/**
 * NvLink lane crc_err_count for lane 0 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0 809

/**
 * NvLink lane crc_err_count for lane 1 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1 810

/**
 * NvLink lane crc_err_count for lane 2 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2 811

/**
 * NvLink lane crc_err_count for lane 3 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3 812

/**
 * NvLink lane ecc_err_count for lane 0 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0 813

/**
 * NvLink lane ecc_err_count for lane 1 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1 814

/**
 * NvLink lane ecc_err_count for lane 2 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2 815

/**
 * NvLink lane ecc_err_count for lane 3 on ports 0-17
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3 816

/**
 * NVSwitch fatal error information.
 * Note: value field indicates the specific SXid reported
 */
#define DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS 856

/**
 * NVSwitch non fatal error information.
 * Note: value field indicates the specific SXid reported
 */
#define DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS 857

/**
 * NVSwitch current temperature.
 */
#define DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT 858

/**
 * NVSwitch limit slowdown temperature.
 */
#define DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN 859

/**
 * NVSwitch limit shutdown temperature.
 */
#define DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN 860

/**
 * NVSwitch throughput Tx.
 */
#define DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX 861

/**
 * NVSwitch throughput Rx.
 */
#define DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX 862

/*
 * NVSwitch Physical ID.
 */
#define DCGM_FI_DEV_NVSWITCH_PHYS_ID 863

/**
 * NVSwitch reset required.
 */
#define DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED 864

/**
 * NvSwitch NvLink ID
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_ID 865

/**
 * NvSwitch PCIE domain
 */
#define DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN 866

/**
 * NvSwitch PCIE bus
 */
#define DCGM_FI_DEV_NVSWITCH_PCIE_BUS 867

/**
 * NvSwitch PCIE device
 */
#define DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE 868

/**
 * NvSwitch PCIE function
 */
#define DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION 869

/**
 * NvLink status.  UNKNOWN:-1 OFF:0 SAFE:1 ACTIVE:2 ERROR:3
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_STATUS 870

/**
 * NvLink device type (GPU/Switch).
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_TYPE 871

/**
 * NvLink device pcie domain.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN 872

/**
 * NvLink device pcie bus.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS 873

/**
 * NvLink device pcie device.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE 874
/**
 * NvLink device pcie function.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION 875

/**
 * NvLink device link ID
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID 876

/**
 * NvLink device SID.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID 877

/**
 * NvLink device link uid.
 */
#define DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_UUID 878

/**
 * Last field ID of the NVSwitch instance
 */
#define DCGM_FI_LAST_NVSWITCH_FIELD_ID 899

/**
 * For now max NVSwitch field Ids taken as difference of DCGM_FI_LAST_NVSWITCH_FIELD_ID and
 * DCGM_FI_FIRST_NVSWITCH_FIELD_ID + 1 i.e. 200
 */
#define DCGM_FI_MAX_NVSWITCH_FIELDS DCGM_FI_LAST_NVSWITCH_FIELD_ID - DCGM_FI_FIRST_NVSWITCH_FIELD_ID + 1

/**
 * Profiling Fields. These all start with DCGM_FI_PROF_*
 */

/**
 * Ratio of time the graphics engine is active. The graphics engine is
 * active if a graphics/compute context is bound and the graphics pipe or
 * compute pipe is busy.
 */
#define DCGM_FI_PROF_GR_ENGINE_ACTIVE 1001

/**
 * The ratio of cycles an SM has at least 1 warp assigned
 * (computed from the number of cycles and elapsed cycles)
 */
#define DCGM_FI_PROF_SM_ACTIVE 1002

/**
 * The ratio of number of warps resident on an SM.
 * (number of resident as a ratio of the theoretical
 * maximum number of warps per elapsed cycle)
 */
#define DCGM_FI_PROF_SM_OCCUPANCY 1003

/**
 * The ratio of cycles the any tensor pipe is active
 * (off the peak sustained elapsed cycles)
 */
#define DCGM_FI_PROF_PIPE_TENSOR_ACTIVE 1004

/**
 * The ratio of cycles the device memory interface is
 * active sending or receiving data.
 */
#define DCGM_FI_PROF_DRAM_ACTIVE 1005

/**
 * Ratio of cycles the fp64 pipe is active.
 */
#define DCGM_FI_PROF_PIPE_FP64_ACTIVE 1006

/**
 * Ratio of cycles the fp32 pipe is active.
 */
#define DCGM_FI_PROF_PIPE_FP32_ACTIVE 1007

/**
 * Ratio of cycles the fp16 pipe is active. This does not include HMMA.
 */
#define DCGM_FI_PROF_PIPE_FP16_ACTIVE 1008

/**
 * The number of bytes of active PCIe tx (transmit) data including both header and payload.
 *
 * Note that this is from the perspective of the GPU, so copying data from device to host (DtoH)
 * would be reflected in this metric.
 */
#define DCGM_FI_PROF_PCIE_TX_BYTES 1009

/**
 * The number of bytes of active PCIe rx (read) data including both header and payload.
 *
 * Note that this is from the perspective of the GPU, so copying data from host to device (HtoD)
 * would be reflected in this metric.
 */
#define DCGM_FI_PROF_PCIE_RX_BYTES 1010

/**
 * The total number of bytes of active NvLink tx (transmit) data including both header and payload.
 * Per-link fields are available below
 */
#define DCGM_FI_PROF_NVLINK_TX_BYTES 1011

/**
 * The total number of bytes of active NvLink rx (read) data including both header and payload.
 * Per-link fields are available below
 */
#define DCGM_FI_PROF_NVLINK_RX_BYTES 1012

/**
 * The ratio of cycles the tensor (IMMA) pipe is active (off the peak sustained elapsed cycles)
 */
#define DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE 1013

/**
 * The ratio of cycles the tensor (HMMA) pipe is active (off the peak sustained elapsed cycles)
 */
#define DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE 1014

/**
 * The ratio of cycles the tensor (DFMA) pipe is active (off the peak sustained elapsed cycles)
 */
#define DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE 1015

/**
 * Ratio of cycles the integer pipe is active.
 */
#define DCGM_FI_PROF_PIPE_INT_ACTIVE 1016

/**
 * Ratio of cycles each of the NVDEC engines are active.
 */
#define DCGM_FI_PROF_NVDEC0_ACTIVE 1017
#define DCGM_FI_PROF_NVDEC1_ACTIVE 1018
#define DCGM_FI_PROF_NVDEC2_ACTIVE 1019
#define DCGM_FI_PROF_NVDEC3_ACTIVE 1020
#define DCGM_FI_PROF_NVDEC4_ACTIVE 1021
#define DCGM_FI_PROF_NVDEC5_ACTIVE 1022
#define DCGM_FI_PROF_NVDEC6_ACTIVE 1023
#define DCGM_FI_PROF_NVDEC7_ACTIVE 1024

/**
 * Ratio of cycles each of the NVJPG engines are active.
 */
#define DCGM_FI_PROF_NVJPG0_ACTIVE 1025
#define DCGM_FI_PROF_NVJPG1_ACTIVE 1026
#define DCGM_FI_PROF_NVJPG2_ACTIVE 1027
#define DCGM_FI_PROF_NVJPG3_ACTIVE 1028
#define DCGM_FI_PROF_NVJPG4_ACTIVE 1029
#define DCGM_FI_PROF_NVJPG5_ACTIVE 1030
#define DCGM_FI_PROF_NVJPG6_ACTIVE 1031
#define DCGM_FI_PROF_NVJPG7_ACTIVE 1032

/**
 * Ratio of cycles each of the NVOFA engines are active.
 */
#define DCGM_FI_PROF_NVOFA0_ACTIVE 1033

/**
 * The per-link number of bytes of active NvLink TX (transmit) or RX (transmit) data including both header and payload.
 * For example: DCGM_FI_PROF_NVLINK_L0_TX_BYTES -> L0 TX
 * To get the bandwidth for a link, add the RX and TX value together like
 * total = DCGM_FI_PROF_NVLINK_L0_TX_BYTES + DCGM_FI_PROF_NVLINK_L0_RX_BYTES
 */
#define DCGM_FI_PROF_NVLINK_L0_TX_BYTES  1040
#define DCGM_FI_PROF_NVLINK_L0_RX_BYTES  1041
#define DCGM_FI_PROF_NVLINK_L1_TX_BYTES  1042
#define DCGM_FI_PROF_NVLINK_L1_RX_BYTES  1043
#define DCGM_FI_PROF_NVLINK_L2_TX_BYTES  1044
#define DCGM_FI_PROF_NVLINK_L2_RX_BYTES  1045
#define DCGM_FI_PROF_NVLINK_L3_TX_BYTES  1046
#define DCGM_FI_PROF_NVLINK_L3_RX_BYTES  1047
#define DCGM_FI_PROF_NVLINK_L4_TX_BYTES  1048
#define DCGM_FI_PROF_NVLINK_L4_RX_BYTES  1049
#define DCGM_FI_PROF_NVLINK_L5_TX_BYTES  1050
#define DCGM_FI_PROF_NVLINK_L5_RX_BYTES  1051
#define DCGM_FI_PROF_NVLINK_L6_TX_BYTES  1052
#define DCGM_FI_PROF_NVLINK_L6_RX_BYTES  1053
#define DCGM_FI_PROF_NVLINK_L7_TX_BYTES  1054
#define DCGM_FI_PROF_NVLINK_L7_RX_BYTES  1055
#define DCGM_FI_PROF_NVLINK_L8_TX_BYTES  1056
#define DCGM_FI_PROF_NVLINK_L8_RX_BYTES  1057
#define DCGM_FI_PROF_NVLINK_L9_TX_BYTES  1058
#define DCGM_FI_PROF_NVLINK_L9_RX_BYTES  1059
#define DCGM_FI_PROF_NVLINK_L10_TX_BYTES 1060
#define DCGM_FI_PROF_NVLINK_L10_RX_BYTES 1061
#define DCGM_FI_PROF_NVLINK_L11_TX_BYTES 1062
#define DCGM_FI_PROF_NVLINK_L11_RX_BYTES 1063
#define DCGM_FI_PROF_NVLINK_L12_TX_BYTES 1064
#define DCGM_FI_PROF_NVLINK_L12_RX_BYTES 1065
#define DCGM_FI_PROF_NVLINK_L13_TX_BYTES 1066
#define DCGM_FI_PROF_NVLINK_L13_RX_BYTES 1067
#define DCGM_FI_PROF_NVLINK_L14_TX_BYTES 1068
#define DCGM_FI_PROF_NVLINK_L14_RX_BYTES 1069
#define DCGM_FI_PROF_NVLINK_L15_TX_BYTES 1070
#define DCGM_FI_PROF_NVLINK_L15_RX_BYTES 1071
#define DCGM_FI_PROF_NVLINK_L16_TX_BYTES 1072
#define DCGM_FI_PROF_NVLINK_L16_RX_BYTES 1073
#define DCGM_FI_PROF_NVLINK_L17_TX_BYTES 1074
#define DCGM_FI_PROF_NVLINK_L17_RX_BYTES 1075

/**
 * NVLink throughput First.
 */
#define DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST DCGM_FI_PROF_NVLINK_L0_TX_BYTES

/**
 * NVLink throughput Last.
 */
#define DCGM_FI_PROF_NVLINK_THROUGHPUT_LAST DCGM_FI_PROF_NVLINK_L17_RX_BYTES

/**
 * CPU Utilization, total
 */
#define DCGM_FI_DEV_CPU_UTIL_TOTAL 1100

/**
 * CPU Utilization, user
 */
#define DCGM_FI_DEV_CPU_UTIL_USER 1101

/**
 * CPU Utilization, nice
 */
#define DCGM_FI_DEV_CPU_UTIL_NICE 1102

/**
 * CPU Utilization, system time
 */
#define DCGM_FI_DEV_CPU_UTIL_SYS 1103

/**
 * CPU Utilization, interrupt servicing
 */
#define DCGM_FI_DEV_CPU_UTIL_IRQ 1104

/**
 * CPU temperature
 */
#define DCGM_FI_DEV_CPU_TEMP_CURRENT 1110

/**
 * CPU Warning Temperature
 */
#define DCGM_FI_DEV_CPU_TEMP_WARNING 1111

/**
 * CPU Critical Temperature
 */
#define DCGM_FI_DEV_CPU_TEMP_CRITICAL 1112

/**
 * CPU instantaneous clock speed
 */
#define DCGM_FI_DEV_CPU_CLOCK_CURRENT 1120

/**
 * CPU power utilization
 */
#define DCGM_FI_DEV_CPU_POWER_UTIL_CURRENT 1130

/**
 * CPU power limit
 */
#define DCGM_FI_DEV_CPU_POWER_LIMIT 1131

/**
 * CPU vendor name
 */
#define DCGM_FI_DEV_CPU_VENDOR 1140

/**
 * CPU model name
 */
#define DCGM_FI_DEV_CPU_MODEL 1141

/**
 * 1 greater than maximum fields above. This is the 1 greater than the maximum field id that could be allocated
 */
#define DCGM_FI_MAX_FIELDS 1142


/** @} */

/*****************************************************************************/

/**
 * Structure for formating the output for dmon.
 * Used as a member in dcgm_field_meta_p
 */
typedef struct
{
    char shortName[10]; /*!< Short name corresponding to field. This short name is used to identify columns in dmon
                             output.*/
    char unit[4];       /*!< The unit of value. Eg: C(elsius), W(att), MB/s*/
    short width;        /*!< Maximum width/number of digits that a value for field can have.*/
} dcgm_field_output_format_t, *dcgm_field_output_format_p;

/**
 * Structure to store meta data for the field
 */

typedef struct
{
    unsigned short fieldId; /*!< Field identifier. DCGM_FI_? #define */
    char fieldType;         /*!< Field type. DCGM_FT_? #define */
    unsigned char size;     /*!< field size in bytes (raw value size). 0=variable (like DCGM_FT_STRING) */
    char tag[48];           /*!< Tag for this field for serialization like 'device_temperature' */
    int scope;              /*!< Field scope. DCGM_FS_? #define of this field's association */
    int nvmlFieldId;        /*!< Optional NVML field this DCGM field maps to. 0 = no mapping.
                                 Otherwise, this should be a NVML_FI_? #define from nvml.h */
    dcgm_field_entity_group_t
        entityLevel; /*!< Field entity level. DCGM_FE_? specifying at what level the field is queryable */

    dcgm_field_output_format_p valueFormat; /*!< pointer to the structure that holds the formatting the
                                                 values for fields */
} dcgm_field_meta_t;

typedef const dcgm_field_meta_t *dcgm_field_meta_p;

/***************************************************************************************************/
/** @addtogroup dcgmFieldIdentifiers
 *  @{
 */
/***************************************************************************************************/

/**
 * Get a pointer to the metadata for a field by its field ID. See DCGM_FI_? for a list of field IDs.
 *
 * @param fieldId     IN: One of the field IDs (DCGM_FI_?)
 *
 * @return
 *        0     On Failure
 *       >0     Pointer to field metadata structure if found.
 *
 */
dcgm_field_meta_p DCGM_PUBLIC_API DcgmFieldGetById(unsigned short fieldId);

/**
 * Get a pointer to the metadata for a field by its field tag.
 *
 * @param tag       IN: Tag for the field of interest
 *
 * @return
 *        0     On failure or not found
 *       >0     Pointer to field metadata structure if found
 *
 */
dcgm_field_meta_p DCGM_PUBLIC_API DcgmFieldGetByTag(const char *tag);

/**
 * Initialize the DcgmFields module. Call this once from inside
 * your program
 *
 * @return
 *        0     On success
 *       <0     On error
 *
 */
int DCGM_PUBLIC_API DcgmFieldsInit(void);

/**
 * Terminates the DcgmFields module. Call this once from inside your program
 *
 * @return
 *        0     On success
 *       <0     On error
 *
 */
int DCGM_PUBLIC_API DcgmFieldsTerm(void);

/**
 * Get the string version of a entityGroupId
 *
 * @returns
 *         - Pointer to a string like GPU/NvSwitch..etc
 *         - Null on error
 *
 */
DCGM_PUBLIC_API const char *DcgmFieldsGetEntityGroupString(dcgm_field_entity_group_t entityGroupId);

/** @} */


#ifdef __cplusplus
}
#endif


#endif // DCGMFIELDS_H
