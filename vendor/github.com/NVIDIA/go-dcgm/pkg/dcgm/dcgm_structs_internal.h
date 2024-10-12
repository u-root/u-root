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
 * File:   dcgm_structs_internal.h
 */

#ifndef DCGM_STRUCTS_INTERNAL_H
#define DCGM_STRUCTS_INTERNAL_H

/* Make sure that dcgm_structs.h is loaded first. This file depends on it */
#include "dcgm_agent.h"
#include "dcgm_structs.h"
#include "dcgm_test_structs.h"
#include <dcgm_nvml.h>

#ifdef INJECTION_LIBRARY_AVAILABLE
#include <nvml_injection.h>
#endif

#ifdef __cplusplus
extern "C" {
#endif


/*
 * The following is a compile time assertion.  It makes use of the
 * restriction that you cannot have an array with a negative size.
 * If the expression resolves to 0, then the index to the array is
 * defined as -1, and a compile time error is generated.  Note that
 * all three macros are needed because of the way the preprocessor
 * evaluates the directives.  Also note that the line number is
 * embedded in the name of the array so that the array name is unique
 * and we can have multiple calls to the assert with the same msg.
 *
 * Usage would be like this:
 * DCGM_CASSERT(DCGM_VGPU_NAME_BUFFER_SIZE == NVML_VGPU_NAME_BUFFER_SIZE, DCGM_VGPU_NAME_BUFFER_SIZE);
 *
 */
#define _DCGM_CASSERT_SYMBOL_INNER(line, msg) COMPILE_TIME_ASSERT_DETECTED_AT_LINE_##line##__##msg
#define _DCGM_CASSERT_SYMBOL(line, msg)       _DCGM_CASSERT_SYMBOL_INNER(line, msg)
#define DCGM_CASSERT(expression, msg) \
    __attribute__((unused)) typedef char _DCGM_CASSERT_SYMBOL(__LINE__, msg)[((expression) ? 1 : -1)]

/**
 * Max length of the DCGM string field
 */
#define DCGM_MAX_STR_LENGTH 256

typedef struct
{
    unsigned int gpuId;             /* DCGM GPU ID */
    char uuid[DCGM_MAX_STR_LENGTH]; /* UUID String */
} dcgmGpuInfo_t;

/* Below is a test API simply to make sure versioning is working correctly
 */

typedef struct
{
    // version must always be first
    unsigned int version;

    unsigned int a;
} dcgmVersionTest_v1;

typedef struct
{
    // version must always be first
    unsigned int version;

    unsigned int a;
    unsigned int b;
} dcgmVersionTest_v2;

typedef dcgmVersionTest_v2 dcgmVersionTest_t;
#define dcgmVersionTest_version1 MAKE_DCGM_VERSION(dcgmVersionTest_v1, 1)
#define dcgmVersionTest_version2 MAKE_DCGM_VERSION(dcgmVersionTest_v2, 2)
#define dcgmVersionTest_version3 MAKE_DCGM_VERSION(dcgmVersionTest_v2, 3)
#define dcgmVersionTest_version  dcgmVersionTest_version2

/**
 * Represents a command to save or load a JSON file to/from the DcgmCacheManager
 */

typedef enum dcgmStatsFileType_enum
{
    DCGM_STATS_FILE_TYPE_JSON = 0 /* JSON */
} dcgmStatsFileType_t;

typedef struct
{
    // version must always be first
    unsigned int version;

    dcgmStatsFileType_t fileType; /* File type to save to/load from */
    char filename[256];           /* Filename to save to/load from */
} dcgmCacheManagerSave_v1_t;

#define dcgmCacheManagerSave_version1 MAKE_DCGM_VERSION(dcgmCacheManagerSave_v1_t, 1)
#define dcgmCacheManagerSave_version  dcgmCacheManagerSave_version1

typedef dcgmCacheManagerSave_v1_t dcgmCacheManagerSave_t;

/* Same message contents for now */
typedef dcgmCacheManagerSave_v1_t dcgmCacheManagerLoad_v1_t;

typedef dcgmCacheManagerLoad_v1_t dcgmCacheManagerLoad_t;

#define dcgmCacheManagerLoad_version1 MAKE_DCGM_VERSION(dcgmCacheManagerLoad_v1_t, 1)
#define dcgmCacheManagerLoad_version  dcgmCacheManagerLoad_version1

#define dcgmWatchFieldValue_version1 1
#define dcgmWatchFieldValue_version  dcgmWatchFieldValue_version1

#define dcgmUpdateAllFields_version1 1
#define dcgmUpdateAllFields_version  dcgmUpdateAllFields_version1

#define dcgmGetMultipleValuesForField_version1 1
#define dcgmGetMultipleValuesForField_version  dcgmGetMultipleValuesForField_version1

#define dcgmUnwatchFieldValue_version1 1
#define dcgmUnwatchFieldValue_version  dcgmUnwatchFieldValue_version1

/**
 * This structure is used to represent a field value to be injected into
 * the cache manager
 */
typedef dcgmFieldValue_v1 dcgmInjectFieldValue_v1;
typedef dcgmInjectFieldValue_v1 dcgmInjectFieldValue_t;
#define dcgmInjectFieldValue_version1 MAKE_DCGM_VERSION(dcgmInjectFieldValue_v1, 1)
#define dcgmInjectFieldValue_version  dcgmInjectFieldValue_version1

#define dcgmGetMultipleValuesForFieldResponse_version1 1
#define dcgmGetMultipleValuesForFieldResponse_version  dcgmGetMultipleValuesForFieldResponse_version1

/* Underlying structure for the GET_MULTIPLE_LATEST_VALUES request */
typedef struct
{
    unsigned int version;                                    /* Set this to dcgmGetMultipleLatestValues_version1 */
    dcgmGpuGrp_t groupId;                                    /* Entity group to retrieve values for. This is only
                                            looked at if entitiesCount is 0 */
    unsigned int entitiesCount;                              /* Number of entities provided in entities[]. This
                                            should only be provided if you aren't also setting
                                            entityGroupId */
    dcgmGroupEntityPair_t entities[DCGM_GROUP_MAX_ENTITIES]; /* Entities to retrieve values for.
                                            Only looked at if entitiesCount > 0 */
    dcgmFieldGrp_t fieldGroupId;                             /* Field group to retrive values for. This is onlu looked
                                                                        at if fieldIdCount is 0 */
    unsigned int fieldIdCount;                               /* Number of field IDs in fieldIds[] that are valid. This
                                                                        should only be set if fieldGroupId is not set */
    unsigned short fieldIds[DCGM_MAX_FIELD_IDS_PER_FIELD_GROUP]; /* Field IDs for which values should
                                            be retrieved. only looked at if fieldIdCount is > 0 */
    unsigned int flags;                                          /* Mask of DCGM_FV_FLAG_? #defines that affect this
                                            request */

} dcgmGetMultipleLatestValues_v1, dcgmGetMultipleLatestValues_t;

#define dcgmGetMultipleLatestValues_version1 MAKE_DCGM_VERSION(dcgmGetMultipleLatestValues_v1, 1)
#define dcgmGetMultipleLatestValues_version  dcgmGetMultipleLatestValues_version1

/* Represents cached record metadata */

/* Represents a unique watcher of an entity in DCGM */

/* Watcher types. Each watcher type's watches are tracked separately within subsystems */
typedef enum
{
    DcgmWatcherTypeClient          = 0, /* Embedded or remote client via external APIs */
    DcgmWatcherTypeHostEngine      = 1, /* Watcher is DcgmHostEngineHandler */
    DcgmWatcherTypeHealthWatch     = 2, /* Watcher is DcgmHealthWatch */
    DcgmWatcherTypePolicyManager   = 3, /* Watcher is DcgmPolicyMgr */
    DcgmWatcherTypeCacheManager    = 4, /* Watcher is DcgmCacheManager */
    DcgmWatcherTypeConfigManager   = 5, /* Watcher is DcgmConfigMgr */
    DcgmWatcherTypeNvSwitchManager = 6, /* Watcher is NvSwitchManager */

    DcgmWatcherTypeCount /* Should always be last */
} DcgmWatcherType_t;


/* ID of a remote client connection within the host engine */
typedef unsigned int dcgm_connection_id_t;

/* Special constant for not connected */
#define DCGM_CONNECTION_ID_NONE ((dcgm_connection_id_t)0)

/* Cache Manager Info flags */
#define DCGM_CMI_F_WATCHED 0x00000001 /* Is this field being watched? */

/* This structure mirrors the DcgmWatcher object */
typedef struct dcgm_cm_field_info_watcher_t
{
    DcgmWatcherType_t watcherType;     /* Type of watcher. See DcgmWatcherType_t */
    dcgm_connection_id_t connectionId; /* Connection ID of the watcher */
    long long monitorIntervalUsec;     /* How often this field should be sampled */
    long long maxAgeUsec;              /* Maximum time to cache samples of this
                                       field. If 0, the class default is used */
} dcgm_cm_field_info_watcher_t, *dcgm_cm_field_info_watcher_p;

/**
 * Number of watchers to show for each field
 */
#define DCGM_CM_FIELD_INFO_NUM_WATCHERS 10

typedef struct dcgmCacheManagerFieldInfo_v4_t
{
    unsigned int version;          /* Version. Check against dcgmCacheManagerInfo_version */
    unsigned int flags;            /* Bitmask of DCGM_CMI_F_? #defines that apply to this field */
    unsigned int entityId;         /* ordinal id for this entity */
    unsigned int entityGroupId;    /* the type of entity, see dcgm_field_entity_group_t */
    unsigned short fieldId;        /* Field ID of this field */
    short lastStatus;              /* Last nvml status returned for this field when taking a sample */
    long long oldestTimestamp;     /* Timestamp of the oldest record. 0=no records or single
                             non-time series record */
    long long newestTimestamp;     /* Timestamp of the newest record. 0=no records or
                             single non-time series record */
    long long monitorIntervalUsec; /* How often is this field updated in usec */
    long long maxAgeUsec;          /* How often is this field updated */
    long long execTimeUsec;        /* Cumulative time spent updating this
                             field since the cache manager started */
    long long fetchCount;          /* Number of times that this field has been
                             fetched from the driver */
    int numSamples;                /* Number of samples currently cached for this field */
    int numWatchers;               /* Number of watchers that are valid in watchers[] */
    dcgm_cm_field_info_watcher_t watchers[DCGM_CM_FIELD_INFO_NUM_WATCHERS]; /* Who are the first 10
                                                                           watchers of this field? */
} dcgmCacheManagerFieldInfo_v4_t, *dcgmCacheManagerFieldInfo_v4_p;

#define dcgmCacheManagerFieldInfo_version4 MAKE_DCGM_VERSION(dcgmCacheManagerFieldInfo_v4_t, 4)

/**
 * The maximum number of topology elements possible given DCGM_MAX_NUM_DEVICES
 * calculated using arithmetic sequence formula
 * (DCGM_MAX_NUM_DEVICES - 1) * (1 + (DCGM_MAX_NUM_DEVICES-2)/2)
 */
#define DCGM_TOPOLOGY_MAX_ELEMENTS 496

/**
 * Topology element structure
 */
typedef struct
{
    unsigned int dcgmGpuA;       //!< GPU A
    unsigned int dcgmGpuB;       //!< GPU B
    unsigned int AtoBNvLinkIds;  //!< bits representing the links connected from GPU A to GPU B
                                 //!< e.g. if this field == 3, links 0 and 1 are connected,
                                 //!< field is only valid if NVLINKS actually exist between GPUs
    unsigned int BtoANvLinkIds;  //!< bits representing the links connected from GPU B to GPU A
                                 //!< e.g. if this field == 3, links 0 and 1 are connected,
                                 //!< field is only valid if NVLINKS actually exist between GPUs
    dcgmGpuTopologyLevel_t path; //!< path between A and B
} dcgmTopologyElement_t;

/**
 * Topology results structure
 */
typedef struct
{
    unsigned int version;     //!< version number (dcgmTopology_version)
    unsigned int numElements; //!< number of valid dcgmTopologyElement_t elements

    dcgmTopologyElement_t element[DCGM_TOPOLOGY_MAX_ELEMENTS];
} dcgmTopology_v1;

/**
 * Typedef for \ref dcgmTopology_v1
 */
typedef dcgmTopology_v1 dcgmTopology_t;

/**
 * Version 1 for \ref dcgmTopology_v1
 */
#define dcgmTopology_version1 MAKE_DCGM_VERSION(dcgmTopology_v1, 1)

/**
 * Latest version for \ref dcgmTopology_t
 */
#define dcgmTopology_version dcgmTopology_version1

typedef struct
{
    unsigned int numGpus;
    struct
    {
        unsigned int dcgmGpuId;
        unsigned long bitmask[DCGM_AFFINITY_BITMASK_ARRAY_SIZE];
    } affinityMasks[DCGM_MAX_NUM_DEVICES];
} dcgmAffinity_t;


typedef struct
{
    unsigned int version;                                       //!< IN: Version number (dcgmCreateFakeEntities_version)
    unsigned int numToCreate;                                   //!< IN: Number of fake entities to create
    dcgmMigHierarchyInfo_t entityList[DCGM_MAX_HIERARCHY_INFO]; //!< IN: specifies who to create and the parent
} dcgmCreateFakeEntities_v2;

typedef dcgmCreateFakeEntities_v2 dcgmCreateFakeEntities_t;

/**
 * Version 2 for \ref dcgmCreateFakeEntities_t
 */
#define dcgmCreateFakeEntities_version2 MAKE_DCGM_VERSION(dcgmCreateFakeEntities_v2, 2)

/**
 * Latest version for \ref dcgmCreateFakeEntities_t
 */
#define dcgmCreateFakeEntities_version dcgmCreateFakeEntities_version2


/* Field watch predefined groups */
typedef enum
{
    DCGM_WATCH_PREDEF_INVALID = 0,
    DCGM_WATCH_PREDEF_PID, /*!< PID stats */
    DCGM_WATCH_PREDEF_JOB, /*!< Job stats */
} dcgmWatchPredefinedType_t;

typedef struct
{
    unsigned int version;
    dcgmWatchPredefinedType_t watchPredefType; /*!< Which type of predefined watch are we adding? */

    dcgmGpuGrp_t groupId; /*!< GPU group to watch fields for */
    long long updateFreq; /*!< How often to update the fields in usec */
    double maxKeepAge;    /*!< How long to keep values for the fields in seconds */
    int maxKeepSamples;   /*!< Maximum number of samples we should keep at a time */
} dcgmWatchPredefined_v1;

typedef dcgmWatchPredefined_v1 dcgmWatchPredefined_t;

/**
 * Version 1 for \ref dcgmWatchPredefined_t
 */
#define dcgmWatchPredefined_version1 MAKE_DCGM_VERSION(dcgmWatchPredefined_v1, 1)

/**
 * Latest version for \ref dcgmWatchPredefined_t
 */
#define dcgmWatchPredefined_version dcgmWatchPredefined_version1

/**
 * Request to set a NvLink link state for an entity
 */
typedef struct
{
    unsigned int version;                    /*!< Version. Should be dcgmSetNvLinkLinkState_version1 */
    dcgm_field_entity_group_t entityGroupId; /*!< Entity group of the entity to set the link state of */
    dcgm_field_eid_t entityId;               /*!< ID of the entity to set the link state of */
    unsigned int linkId;                     /*!< Link (or portId) of the link to set the state of */
    dcgmNvLinkLinkState_t linkState;         /*!< State to set the link to */
    unsigned int unused;                     /*!< Not used for now. Set to 0 */
} dcgmSetNvLinkLinkState_v1;

#define dcgmSetNvLinkLinkState_version1 MAKE_DCGM_VERSION(dcgmSetNvLinkLinkState_v1, 1)


/**
 * Request to add a module ID to the denylist
 */
typedef struct
{
    unsigned int version;    /*!< Version. Should be dcgmModuleDenylist_version */
    dcgmModuleId_t moduleId; /*!< Module to add to the denylist */
} dcgmModuleDenylist_v1;

#define dcgmModuleDenylist_version1 MAKE_DCGM_VERSION(dcgmModuleDenylist_v1, 1)


/**
 * Counter to use for NvLink
 */
#define DCGMCM_NVLINK_COUNTER_BYTES 0

/**
 * The Brand of the GPU. These are 1:1 with NVML_BRAND_*. There's a DCGM_CASSERT() below that tests that
 */
typedef enum dcgmGpuBrandType_enum
{
    DCGM_GPU_BRAND_UNKNOWN = 0,
    DCGM_GPU_BRAND_QUADRO  = 1,
    DCGM_GPU_BRAND_TESLA   = 2,
    DCGM_GPU_BRAND_NVS     = 3,
    DCGM_GPU_BRAND_GRID    = 4,
    DCGM_GPU_BRAND_GEFORCE = 5,
    DCGM_GPU_BRAND_TITAN   = 6,
    /* The following are new as of r460 TRD2's nvml.h */
    DCGM_BRAND_NVIDIA_VAPPS   = 7,  // NVIDIA Virtual Applications
    DCGM_BRAND_NVIDIA_VPC     = 8,  // NVIDIA Virtual PC
    DCGM_BRAND_NVIDIA_VCS     = 9,  // NVIDIA Virtual Compute Server
    DCGM_BRAND_NVIDIA_VWS     = 10, // NVIDIA RTX Virtual Workstation
    DCGM_BRAND_NVIDIA_VGAMING = 11, // NVIDIA vGaming
    DCGM_BRAND_QUADRO_RTX     = 12,
    DCGM_BRAND_NVIDIA_RTX     = 13,
    DCGM_BRAND_NVIDIA         = 14,
    DCGM_BRAND_GEFORCE_RTX    = 15,
    DCGM_BRAND_TITAN_RTX      = 16,
    // Keep this last
    DCGM_GPU_BRAND_COUNT
} dcgmGpuBrandType_t;

/*****************************************************************************/
typedef enum dcgmEntityStatusType_enum
{
    DcgmEntityStatusUnknown = 0,  /* Entity has not been referenced yet */
    DcgmEntityStatusOk,           /* Entity is known and OK */
    DcgmEntityStatusUnsupported,  /* Entity is unsupported by DCGM */
    DcgmEntityStatusInaccessible, /* Entity is inaccessible, usually due to cgroups */
    DcgmEntityStatusLost,         /* Entity has been lost. Usually set from NVML
                                   returning NVML_ERROR_GPU_IS_LOST */
    DcgmEntityStatusFake,         /* Entity is a fake, injection-only entity for testing */
    DcgmEntityStatusDisabled,     /* Don't collect values from this GPU */
    DcgmEntityStatusDetached      /* Entity is detached, not good for any uses */
} DcgmEntityStatus_t;

/**
 * Making these internal so that client apps must be explicit with struct versions.
 */

/**
 * Typedef for \ref dcgmRunDiag_t
 */
typedef dcgmRunDiag_v7 dcgmRunDiag_t;

/**
 * Latest version for \ref dcgmRunDiag_t
 */
#define dcgmRunDiag_version dcgmRunDiag_version7

/**
 * Version 1 of dcgmCreateGroup_t
 */

typedef struct
{
    dcgmGroupType_t groupType; //!< Type of group to create
    char groupName[1024];      //!< Name to give new group
    dcgmGpuGrp_t newGroupId;   //!< On success, the ID of the newly created group
    dcgmReturn_t cmdRet;       //!< Error code generated when creating new group
} dcgmCreateGroup_v1;

/**
 * Version 1 of dcgmRemoveEntity_t
 */

typedef struct
{
    unsigned int groupId;       //!< IN: Group id from which entity should be removed
    unsigned int entityGroupId; //!< IN: Entity group that entity belongs to
    unsigned int entityId;      //!< IN: Entity id to remove
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmAddRemoveEntity_v1;

/**
 * Version 1 of dcgmGroupDestroy_t
 */

typedef struct
{
    unsigned int groupId; //!< IN: Group to remove
    unsigned int cmdRet;  //!< OUT: Error code generated
} dcgmGroupDestroy_v1;

/**
 * Version 1 of dcgmGetEntityGroupEntities_t
 */

typedef struct
{
    unsigned int entityGroup;                       //!< IN: Entity of group to list entities
    unsigned int entities[DCGM_GROUP_MAX_ENTITIES]; //!< OUT: Array of entities for entityGroup
    unsigned int numEntities;                       //!< IN/OUT: Upon calling, this should be the number of
                                                    //           entities that entityList[] can hold. Upon
                                                    //           return, this will contain the number of
                                                    //           entities actually saved to entityList.
    unsigned int flags;                             //!< IN: Flags to modify the behavior of this request.
                                                    //       See DCGM_GEGE_FLAG_*
    unsigned int cmdRet;                            //!< OUT: Error code generated
} dcgmGetEntityGroupEntities_v1;

/**
 * Version 1 of dcgmGroupGetAllIds_t
 */

typedef struct
{
    unsigned int groupIds[DCGM_MAX_NUM_GROUPS]; //!< OUT: List of group ids
    unsigned int numGroups;                     //!< OUT: Number of group ids in the list
    unsigned int cmdRet;                        //!< OUT: Error code generated
} dcgmGroupGetAllIds_v1;

/**
 * Version 1 of dcgmGroupGetInfo_t
 */

typedef struct
{
    unsigned int groupId;      //!< IN: Group ID for which information to be fetched
    dcgmGroupInfo_t groupInfo; //!< OUT: Group Information
    long long timestamp;       //!< OUT: Timestamp of information
    unsigned int cmdRet;       //!< OUT: Error code generated
} dcgmGroupGetInfo_v1;

#define SAMPLES_BUFFER_SIZE_V1 16384

/**
 * Version 1 of dcgmEntitiesGetLatestValues_t
 */
typedef struct
{
    unsigned int groupId;                                    //!< IN: Optional group id for information to be fetched
    dcgmGroupEntityPair_t entities[DCGM_GROUP_MAX_ENTITIES]; //!< IN: List of entities to get values for
    unsigned int entitiesCount;                              //!< IN: Number of entries in entities[]
    unsigned int fieldGroupId; //!< IN: Optional fieldGroupId that will be resolved by the host engine.
                               //!<     This is ignored if fieldIdList[] is provided
    unsigned short fieldIdList[DCGM_MAX_FIELD_IDS_PER_FIELD_GROUP]; //!< IN: Field IDs to return data for
    unsigned int fieldIdCount;                                      //!< IN: Number of field IDs in fieldIdList[] array.
    unsigned int flags;                  //!< IN: Optional flags that affect how this request is processed.
    unsigned int cmdRet;                 //!< OUT: Error code generated
    unsigned int bufferSize;             //!< OUT: Length of populated buffer
    char buffer[SAMPLES_BUFFER_SIZE_V1]; //!< OUT: this field is last, and can be truncated for speed */
} dcgmEntitiesGetLatestValues_v1;

#define SAMPLES_BUFFER_SIZE_V2 4186112 // 4MB - 8k for header

/**
 * Version 2 of dcgmEntitiesGetLatestValues_t
 */
typedef struct
{
    unsigned int groupId;                                    //!< IN: Optional group id for information to be fetched
    dcgmGroupEntityPair_t entities[DCGM_GROUP_MAX_ENTITIES]; //!< IN: List of entities to get values for
    unsigned int entitiesCount;                              //!< IN: Number of entries in entities[]
    unsigned int fieldGroupId; //!< IN: Optional fieldGroupId that will be resolved by the host engine.
                               //!<     This is ignored if fieldIdList[] is provided
    unsigned short fieldIdList[DCGM_MAX_FIELD_IDS_PER_FIELD_GROUP]; //!< IN: Field IDs to return data for
    unsigned int fieldIdCount;                                      //!< IN: Number of field IDs in fieldIdList[] array.
    unsigned int flags;                  //!< IN: Optional flags that affect how this request is processed.
    unsigned int cmdRet;                 //!< OUT: Error code generated
    unsigned int bufferSize;             //!< OUT: Length of populated buffer
    char buffer[SAMPLES_BUFFER_SIZE_V2]; //!< OUT: this field is last, and can be truncated for speed */
} dcgmEntitiesGetLatestValues_v2;

/**
 * Version 1 of dcgmGetMultipleValuesForField
 */
typedef struct
{
    unsigned int entityGroupId;          //!< IN: Optional group id for information to be fetched
    unsigned int entityId;               //!< IN: Optional entity id for information to be fetched
    unsigned int fieldId;                //!< IN: Field id to fetch
    long long startTs;                   //!< IN: Starting timestamp
    long long endTs;                     //!< IN: End timestamp
    unsigned int order;                  //!< IN: Order for output data, see dcgmOrder_t
    unsigned int count;                  //!< IN: Number of values to retrieve (may be limited by size of buffer)
    unsigned int cmdRet;                 //!< OUT: Error code generated
    unsigned int bufferSize;             //!< OUT: Length of populated buffer
    char buffer[SAMPLES_BUFFER_SIZE_V1]; //!< OUT:: this field is last, and can be truncated for speed */
} dcgmGetMultipleValuesForField_v1;

/**
 * Version 2 of dcgmGetMultipleValuesForField
 */
typedef struct
{
    unsigned int entityGroupId;          //!< IN: Optional group id for information to be fetched
    unsigned int entityId;               //!< IN: Optional entity id for information to be fetched
    unsigned int fieldId;                //!< IN: Field id to fetch
    long long startTs;                   //!< IN: Starting timestamp
    long long endTs;                     //!< IN: End timestamp
    unsigned int order;                  //!< IN: Order for output data, see dcgmOrder_t
    unsigned int count;                  //!< IN: Number of values to retrieve (may be limited by size of buffer)
    unsigned int cmdRet;                 //!< OUT: Error code generated
    unsigned int bufferSize;             //!< OUT: Length of populated buffer
    char buffer[SAMPLES_BUFFER_SIZE_V2]; //!< OUT:: this field is last, and can be truncated for speed */
} dcgmGetMultipleValuesForField_v2;

/**
 * Version 1 of dcgmJobCmd_t
 */

typedef struct
{
    unsigned int groupId; //!< IN: optional group id
    char jobId[64];       //!< IN: job id
    unsigned int cmdRet;  //!< OUT: Error code generated
} dcgmJobCmd_v1;

/**
 * Version 1 of dcgmJobGetStats_t
 */

typedef struct
{
    char jobId[64];         //!< IN: job id
    dcgmJobInfo_t jobStats; //!< OUT: job stats
    unsigned int cmdRet;    //!< OUT: Error code generated
} dcgmJobGetStats_v1;

/**
 * Version 1 of dcgmWatchFieldValue_t (DCGM 2.x)
 */
typedef struct
{
    int gpuId;                  //!< IN: GPU ID to watch field on
    unsigned int entityGroupId; //!< IN: Optional entity group id
    unsigned short fieldId;     //!< IN: Field ID to watch
    long long updateFreq;       //!< IN: How often to update this field in usec
    double maxKeepAge;          //!< IN: How long to keep data for this field in seconds
    int maxKeepSamples;         //!< IN: Maximum number of samples to keep. 0=no limit
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmWatchFieldValue_v1;

/**
 * Version 2 of dcgmWatchFieldValue_t (DCGM 3.x+)
 */
typedef struct
{
    unsigned int entityId;      //!< IN: entityId (gpuId for GPUs) to watch field on
    unsigned int entityGroupId; //!< IN: Optional entity group id
    unsigned short fieldId;     //!< IN: Field ID to watch
    unsigned char unused[6];    //!< IN: Unused. Aligns next member to 8-byte boundary
    long long updateFreq;       //!< IN: How often to update this field in usec
    double maxKeepAge;          //!< IN: How long to keep data for this field in seconds
    int maxKeepSamples;         //!< IN: Maximum number of samples to keep. 0=no limit
    int updateOnFirstWatcher;   //!< IN: Should we do an UpdateAllFields() automatically if we are the first watcher?
                                //!< 1=yes. 0=no.
    int wereFirstWatcher;       //!< OUT: Returns 1 if we were the first watcher. 0 if not */
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmWatchFieldValue_v2;

/**
 * Version 1 of dcgmUpdateAllFields_v1
 */
typedef struct
{
    int waitForUpdate;   //!< IN: Whether or not to wait for the update loop to complete before returning to the
                         //       caller 1=wait. 0=do not wait.
    unsigned int cmdRet; //!< OUT: Error code generated
} dcgmUpdateAllFields_v1;

/**
 * Version 1 of dcgmUnwatchFieldValue_t
 */
typedef struct
{
    int gpuId;                  //!< IN: GPU ID to watch field on
    unsigned int entityGroupId; //!< IN: Optional entity group id
    unsigned short fieldId;     //!< IN: Field id to unwatch
    int clearCache;             //!< IN: Whether or not to clear all cached data for
                                //       the field after the watch is removed
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmUnwatchFieldValue_v1;

/**
 * Version 1 of dcgmInjectFieldValue_t
 */
typedef struct
{
    unsigned int entityGroupId;   //!< IN: entity group id
    unsigned int entityId;        //!< IN: entity id
    dcgmFieldValue_v1 fieldValue; //!< IN: field value to insert
    unsigned int cmdRet;          //!< OUT: Error code generated
} dcgmInjectFieldValueMsg_v1;

#define dcgmInjectFieldValueMsg_version1 MAKE_DCGM_VERSION(dcgmInjectFieldValueMsg_v1, 1)
#define dcgmInjectFieldValueMsg_version  dcgmInjectFieldValueMsg_version1
typedef dcgmInjectFieldValueMsg_v1 dcgmInjectFieldValueMsg_t;

/**
 * Version 2 of dcgmGetCacheManagerFieldInfo_t
 */
typedef struct
{
    dcgmCacheManagerFieldInfo_v4_t
        fieldInfo;       //!< IN/OUT: Structure to populate. fieldInfo->gpuId and fieldInfo->fieldId must
                         //           be populated on calling for this call to work
    unsigned int cmdRet; //!< OUT: Error code generated
} dcgmGetCacheManagerFieldInfo_v2;

typedef struct
{
    unsigned int groupId;      //!< IN: Group ID representing collection of one or more entities
    unsigned int fieldGroupId; //!< IN: Fields to watch.
    long long updateFreq;      //!< IN: How often to update this field in usec
    double maxKeepAge;         //!< IN: How long to keep data for this field in seconds
    int maxKeepSamples;        //!< IN: Maximum number of samples to keep. 0=no limit
    unsigned int cmdRet;       //!< OUT: Error code generated
} dcgmWatchFields_v1;

#define dcgmWatchFields_version1 1
#define dcgmWatchFields_version  dcgmWatchFields_version1

typedef struct
{
    unsigned int groupId;    //!< IN: Group ID representing collection of one or more entities
    dcgmTopology_t topology; //!< OUT: populated struct
    unsigned int cmdRet;     //!< OUT: Error code generated
} dcgmGetTopologyMsg_v1;

typedef struct
{
    unsigned int groupId;    //!< IN: Group ID representing collection of one or more entities
    dcgmAffinity_t affinity; //!< OUT: populated struct
    unsigned int cmdRet;     //!< OUT: Error code generated
} dcgmGetTopologyAffinityMsg_v1;

typedef struct
{
    uint64_t inputGpus;  //!< IN: bitmask of available gpus
    uint32_t numGpus;    //!< IN: number of gpus needed
    uint64_t flags;      //!< IN: Hints to ignore certain factors for the scheduling hint
    uint64_t outputGpus; //!< OUT: bitmask of selected gpus
    unsigned int cmdRet; //!< OUT: Error code generated
} dcgmSelectGpusByTopologyMsg_v1;

typedef struct
{
    int supported;                              //!< IN: boolean to ONLY include Ids of supported GPUs
    unsigned int devices[DCGM_MAX_NUM_DEVICES]; //!< OUT: GPU Ids present on the system.
    int count;                                  //!< OUT: Number of devices returned in "devices"
    unsigned int cmdRet;                        //!< OUT: Error code generated
} dcgmGetAllDevicesMsg_v1;

typedef struct
{
    int persistAfterDisconnect; //!< IN: boolean whether to persist groups, etc after client is disconnected
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmClientLogin_v1;

typedef struct
{
    dcgmFieldGroupInfo_t fg; //!< IN/OUT: field group info populated on success
    unsigned int cmdRet;     //!< OUT: Error code generated
} dcgmFieldGroupOp_v1;

typedef struct
{
    unsigned int groupId;  //!< IN: group id for query
    dcgmPidInfo_t pidInfo; //!< IN/OUT: pid info populated on success
    unsigned int cmdRet;   //!< OUT: Error code generated
} dcgmPidGetInfo_v1;

typedef struct
{
    dcgmFieldSummaryRequest_t fsr; //!< IN/OUT: field summary populated on success
    unsigned int cmdRet;           //!< OUT: Error code generated
} dcgmGetFieldSummary_v1;

typedef struct
{
    dcgmNvLinkStatus_v3 ls; //!< IN/OUT: nvlink status populated on success
    unsigned int cmdRet;    //!< OUT: Error code generated
} dcgmGetNvLinkStatus_v2;

typedef struct
{
    dcgmCreateFakeEntities_v2 fe; //!< IN/OUT: fake entity info, populated on success
    unsigned int cmdRet;          //!< OUT: Error code generated
} dcgmMsgCreateFakeEntities_v1;

typedef struct
{
    dcgmWatchPredefined_t wpf; //!< IN: watch info
    unsigned int cmdRet;       //!< OUT: Error code generated
} dcgmWatchPredefinedFields_v1;

typedef struct
{
    unsigned int moduleId; //!< IN: Module to add to the denylist
    unsigned int cmdRet;   //!< OUT: Error code generated
} dcgmMsgModuleDenylist_v1;

typedef struct
{
    dcgmModuleGetStatuses_t st; //!< IN/OUT: module status
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmMsgModuleStatus_v1;

typedef struct
{
    unsigned int overallHealth; //!< IN/OUT: hostengine health
    unsigned int cmdRet;        //!< OUT: Error code generated
} dcgmMsgHostEngineHealth_v1;

typedef struct
{
    dcgmAllFieldGroup_t fg; //!< IN/OUT: hostengine health
    unsigned int cmdRet;    //!< OUT: Error code generated
} dcgmGetAllFieldGroup_v1;

typedef struct
{
    dcgmMigHierarchy_v2 data; //!< OUT: populated on success

    unsigned int cmdRet; //!< OUT: Error code generated
} dcgmMsgGetGpuInstanceHierarchy_v1;

typedef struct
{
    unsigned int index;  //!< IN: the index of the GPU to create
    unsigned int cmdRet; //!< OUT: Error code generated
} dcgmMsgNvmlCreateInjectionGpu_v1;

#ifdef INJECTION_LIBRARY_AVAILABLE
#define DCGM_MAX_EXTRA_KEYS 4
typedef struct
{
    unsigned int gpuId;                             //!< IN: the DCGM gpu id of the device being injected
    char key[DCGM_MAX_STR_LENGTH];                  //!< IN: The key for the NVML injected value
    injectNvmlVal_t extraKeys[DCGM_MAX_EXTRA_KEYS]; //!< IN: extra keys, optional
    unsigned int extraKeyCount;                     //!< IN: the number of extra keys
    injectNvmlVal_t value;                          //!< IN: the NVML value being injected
    unsigned int cmdRet;                            //!< OUT: Error code generated
} dcgmMsgNvmlInjectDevice_v1;
#endif

/**
 * Verify that DCGM definitions that are copies of NVML ones match up with their NVML counterparts
 */
DCGM_CASSERT(DCGM_VGPU_NAME_BUFFER_SIZE == NVML_VGPU_NAME_BUFFER_SIZE, NVML_VGPU_NAME_BUFFER_SIZE);
DCGM_CASSERT(DCGM_GRID_LICENSE_BUFFER_SIZE == NVML_GRID_LICENSE_BUFFER_SIZE, NVML_GRID_LICENSE_BUFFER_SIZE);
DCGM_CASSERT(DCGM_DEVICE_UUID_BUFFER_SIZE == NVML_DEVICE_UUID_BUFFER_SIZE, NVML_DEVICE_UUID_BUFFER_SIZE);
DCGM_CASSERT(DCGM_NVLINK_MAX_LINKS_PER_GPU == NVML_NVLINK_MAX_LINKS, NVML_NVLINK_MAX_LINKS);
DCGM_CASSERT((int)DCGM_GPU_BRAND_COUNT == (int)NVML_BRAND_COUNT, NVML_BRAND_COUNT);

DCGM_CASSERT((int)DCGM_GPU_VIRTUALIZATION_MODE_NONE == (int)NVML_GPU_VIRTUALIZATION_MODE_NONE,
             NVML_GPU_VIRTUALIZATION_MODE_NONE);
DCGM_CASSERT((int)DCGM_GPU_VIRTUALIZATION_MODE_PASSTHROUGH == (int)NVML_GPU_VIRTUALIZATION_MODE_PASSTHROUGH,
             NVML_GPU_VIRTUALIZATION_MODE_PASSTHROUGH);
DCGM_CASSERT((int)DCGM_GPU_VIRTUALIZATION_MODE_VGPU == (int)NVML_GPU_VIRTUALIZATION_MODE_VGPU,
             NVML_GPU_VIRTUALIZATION_MODE_VGPU);
DCGM_CASSERT((int)DCGM_GPU_VIRTUALIZATION_MODE_HOST_VGPU == (int)NVML_GPU_VIRTUALIZATION_MODE_HOST_VGPU,
             NVML_GPU_VIRTUALIZATION_MODE_HOST_VGPU);
DCGM_CASSERT((int)DCGM_GPU_VIRTUALIZATION_MODE_HOST_VSGA == (int)NVML_GPU_VIRTUALIZATION_MODE_HOST_VSGA,
             NVML_GPU_VIRTUALIZATION_MODE_HOST_VSGA);

/**
 *  Verify correct version of APIs that use a versioned structure
 */

DCGM_CASSERT(dcgmPidInfo_version == (long)0x02004528, 1);
DCGM_CASSERT(dcgmConfig_version == (long)16777256, 1);
DCGM_CASSERT(dcgmConnectV2Params_version1 == (long)16777224, 1);
DCGM_CASSERT(dcgmConnectV2Params_version == (long)0x02000010, 1);
DCGM_CASSERT(dcgmCpuHierarchyOwnedCores_version1 == (long)0x1000088, 1);
DCGM_CASSERT(dcgmCpuHierarchy_version1 == (long)0x1000488, 1);
DCGM_CASSERT(dcgmFieldGroupInfo_version == (long)16777744, 1);
DCGM_CASSERT(dcgmAllFieldGroup_version == (long)16811016, 1);
DCGM_CASSERT(dcgmDeviceAttributes_version3 == (long)0x3001464, 1);
DCGM_CASSERT(dcgmDeviceAttributes_version == (long)0x3001464, 1);
DCGM_CASSERT(dcgmHealthResponse_version4 == (long)0x0401050C, 1);
DCGM_CASSERT(dcgmIntrospectMemory_version == (long)16777232, 1);
DCGM_CASSERT(dcgmIntrospectCpuUtil_version == (long)16777248, 1);
DCGM_CASSERT(dcgmJobInfo_version == (long)0x030098A8, 1);
DCGM_CASSERT(dcgmPolicy_version == (long)16777360, 1);
DCGM_CASSERT(dcgmPolicyCallbackResponse_version == (long)16777240, 1);
DCGM_CASSERT(dcgmDiagResponse_version7 == (long)0x07099290, 1);
DCGM_CASSERT(dcgmDiagResponse_version8 == (long)0x80d9690, 8);
DCGM_CASSERT(dcgmDiagResponse_version9 == (long)0x914f4dc, 9);
DCGM_CASSERT(dcgmDiagResponse_version == (long)0x914f4dc, 9);
DCGM_CASSERT(dcgmRunDiag_version7 == (long)0x70054D0, 1);
DCGM_CASSERT(dcgmVgpuDeviceAttributes_version6 == (long)16787744, 1);
DCGM_CASSERT(dcgmVgpuDeviceAttributes_version7 == (long)117451168, 1);
DCGM_CASSERT(dcgmVgpuDeviceAttributes_version == (long)117451168, 1);
DCGM_CASSERT(dcgmVgpuInstanceAttributes_version == (long)16777556, 1);
DCGM_CASSERT(dcgmVgpuConfig_version == (long)16777256, 1);
DCGM_CASSERT(dcgmModuleGetStatuses_version == (long)0x01000088, 1);
DCGM_CASSERT(dcgmModuleDenylist_version1 == (long)0x01000008, 1);
DCGM_CASSERT(dcgmSettingsSetLoggingSeverity_version1 == (long)0x01000008, 1);
DCGM_CASSERT(dcgmVersionInfo_version == (long)0x2000204, 1);
DCGM_CASSERT(dcgmStartEmbeddedV2Params_version1 == (long)0x01000048, 1);
DCGM_CASSERT(dcgmStartEmbeddedV2Params_version2 == (long)0x02000050, 2);
DCGM_CASSERT(dcgmInjectFieldValue_version1 == (long)0x1001018, 1);
DCGM_CASSERT(dcgmInjectFieldValue_version == (long)0x1001018, 1);
DCGM_CASSERT(dcgmNvLinkStatus_version3 == (long)0x30015bc, 3);

#ifndef DCGM_ARRAY_CAPACITY
#ifdef __cplusplus
#define DCGM_ARRAY_CAPACITY(a) std::extent<decltype(a)>::value
static_assert(NVML_COMPUTE_INSTANCE_PROFILE_COUNT == 0x08);
static_assert(NVML_GPU_INSTANCE_PROFILE_1_SLICE_REV2 == 0x09);
#endif
#endif

#ifndef DCGM_ARRAY_CAPACITY
#define DCGM_ARRAY_CAPACITY(a) (sizeof(a) / sizeof(a[0]))
#endif

#ifdef __cplusplus
}
#endif

#endif /* DCGM_STRUCTS_H */
