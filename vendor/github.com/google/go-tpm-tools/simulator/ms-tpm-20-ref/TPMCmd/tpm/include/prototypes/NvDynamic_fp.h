/* Microsoft Reference Implementation for TPM 2.0
 *
 *  The copyright in this software is being made available under the BSD License,
 *  included below. This software may be subject to other third party and
 *  contributor rights, including patent rights, and no such rights are granted
 *  under this license.
 *
 *  Copyright (c) Microsoft Corporation
 *
 *  All rights reserved.
 *
 *  BSD License
 *
 *  Redistribution and use in source and binary forms, with or without modification,
 *  are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list
 *  of conditions and the following disclaimer.
 *
 *  Redistributions in binary form must reproduce the above copyright notice, this
 *  list of conditions and the following disclaimer in the documentation and/or
 *  other materials provided with the distribution.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS ""AS IS""
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 *  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
 *  ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 *  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 *  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 *  ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 *  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 *  SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
/*(Auto-generated)
 *  Created by TpmPrototypes; Version 3.0 July 18, 2017
 *  Date: Apr  7, 2019  Time: 06:58:58PM
 */

#ifndef    _NV_DYNAMIC_FP_H_
#define    _NV_DYNAMIC_FP_H_

//*** NvWriteNvListEnd()
// Function to write the list terminator.
NV_REF
NvWriteNvListEnd(
    NV_REF           end
);

//*** NvUpdateIndexOrderlyData()
// This function is used to cause an update of the orderly data to the NV backing
// store.
void
NvUpdateIndexOrderlyData(
    void
);

//*** NvReadIndex()
// This function is used to read the NV Index NV_INDEX. This is used so that the
// index information can be compressed and only this function would be needed
// to decompress it. Mostly, compression would only be able to save the space
// needed by the policy.
void
NvReadNvIndexInfo(
    NV_REF           ref,           // IN: points to NV where index is located
    NV_INDEX        *nvIndex        // OUT: place to receive index data
);

//*** NvReadObject()
// This function is used to read a persistent object. This is used so that the
// object information can be compressed and only this function would be needed
// to uncompress it.
void
NvReadObject(
    NV_REF           ref,           // IN: points to NV where index is located
    OBJECT          *object         // OUT: place to receive the object data
);

//*** NvIndexIsDefined()
// See if an index is already defined
BOOL
NvIndexIsDefined(
    TPM_HANDLE       nvHandle       // IN: Index to look for
);

//*** NvIsPlatformPersistentHandle()
// This function indicates if a handle references a persistent object in the
// range belonging to the platform.
//  Return Type: BOOL
//      TRUE(1)         handle references a platform persistent object
//                      and may reference an owner persistent object either
//      FALSE(0)        handle does not reference platform persistent object
BOOL
NvIsPlatformPersistentHandle(
    TPM_HANDLE       handle         // IN: handle
);

//*** NvIsOwnerPersistentHandle()
// This function indicates if a handle references a persistent object in the
// range belonging to the owner.
//  Return Type: BOOL
//      TRUE(1)         handle is owner persistent handle
//      FALSE(0)        handle is not owner persistent handle and may not be
//                      a persistent handle at all
BOOL
NvIsOwnerPersistentHandle(
    TPM_HANDLE       handle         // IN: handle
);

//*** NvIndexIsAccessible()
//
// This function validates that a handle references a defined NV Index and
// that the Index is currently accessible.
//  Return Type: TPM_RC
//      TPM_RC_HANDLE           the handle points to an undefined NV Index
//                              If shEnable is CLEAR, this would include an index
//                              created using ownerAuth. If phEnableNV is CLEAR,
//                              this would include and index created using
//                              platformAuth
//      TPM_RC_NV_READLOCKED    Index is present but locked for reading and command
//                              does not write to the index
//      TPM_RC_NV_WRITELOCKED   Index is present but locked for writing and command
//                              writes to the index
TPM_RC
NvIndexIsAccessible(
    TPMI_RH_NV_INDEX     handle        // IN: handle
);

//*** NvGetEvictObject()
// This function is used to dereference an evict object handle and get a pointer
// to the object.
//  Return Type: TPM_RC
//      TPM_RC_HANDLE           the handle does not point to an existing
//                              persistent object
TPM_RC
NvGetEvictObject(
    TPM_HANDLE       handle,        // IN: handle
    OBJECT          *object         // OUT: object data
);

//*** NvIndexCacheInit()
// Function to initialize the Index cache
void
NvIndexCacheInit(
    void
);

//*** NvGetIndexData()
// This function is used to access the data in an NV Index. The data is returned
// as a byte sequence.
//
// This function requires that the NV Index be defined, and that the
// required data is within the data range.  It also requires that TPMA_NV_WRITTEN
// of the Index is SET.
void
NvGetIndexData(
    NV_INDEX            *nvIndex,       // IN: the in RAM index descriptor
    NV_REF               locator,       // IN: where the data is located
    UINT32               offset,        // IN: offset of NV data
    UINT16               size,          // IN: number of octets of NV data to read
    void                *data           // OUT: data buffer
);

//*** NvHashIndexData()
// This function adds Index data to a hash. It does this in parts to avoid large stack
// buffers.
void
NvHashIndexData(
    HASH_STATE          *hashState,     // IN: Initialized hash state
    NV_INDEX            *nvIndex,       // IN: Index
    NV_REF               locator,       // IN: where the data is located
    UINT32               offset,        // IN: starting offset
    UINT16               size           // IN: amount to hash
);

//*** NvGetUINT64Data()
// Get data in integer format of a bit or counter NV Index.
//
// This function requires that the NV Index is defined and that the NV Index
// previously has been written.
UINT64
NvGetUINT64Data(
    NV_INDEX            *nvIndex,       // IN: the in RAM index descriptor
    NV_REF               locator        // IN: where index exists in NV
);

//*** NvWriteIndexAttributes()
// This function is used to write just the attributes of an index.
//  Return type: TPM_RC
//      TPM_RC_NV_RATE          NV is rate limiting so retry
//      TPM_RC_NV_UNAVAILABLE   NV is not available
TPM_RC
NvWriteIndexAttributes(
    TPM_HANDLE       handle,
    NV_REF           locator,       // IN: location of the index
    TPMA_NV          attributes     // IN: attributes to write
);

//*** NvWriteIndexAuth()
// This function is used to write the authValue of an index. It is used by
// TPM2_NV_ChangeAuth()
//  Return type: TPM_RC
//      TPM_RC_NV_RATE          NV is rate limiting so retry
//      TPM_RC_NV_UNAVAILABLE   NV is not available
TPM_RC
NvWriteIndexAuth(
    NV_REF           locator,       // IN: location of the index
    TPM2B_AUTH      *authValue      // IN: the authValue to write
);

//*** NvGetIndexInfo()
// This function loads the nvIndex Info into the NV cache and returns a pointer
// to the NV_INDEX. If the returned value is zero, the index was not found.
// The 'locator' parameter, if not NULL, will be set to the offset in NV of the
// Index (the location of the handle of the Index).
//
// This function will set the index cache. If the index is orderly, the attributes
// from RAM are substituted for the attributes in the cached index
NV_INDEX *
NvGetIndexInfo(
    TPM_HANDLE       nvHandle,      // IN: the index handle
    NV_REF          *locator        // OUT: location of the index
);

//*** NvWriteIndexData()
// This function is used to write NV index data. It is intended to be used to
// update the data associated with the default index.
//
// This function requires that the NV Index is defined, and the data is
// within the defined data range for the index.
//
// Index data is only written due to a command that modifies the data in a single
// index. There is no case where changes are made to multiple indexes data at the
// same time. Multiple attributes may be change but not multiple index data. This
// is important because we will normally be handling the index for which we have
// the cached pointer values.
//  Return type: TPM_RC
//      TPM_RC_NV_RATE          NV is rate limiting so retry
//      TPM_RC_NV_UNAVAILABLE   NV is not available
TPM_RC
NvWriteIndexData(
    NV_INDEX        *nvIndex,       // IN: the description of the index
    UINT32           offset,        // IN: offset of NV data
    UINT32           size,          // IN: size of NV data
    void            *data           // IN: data buffer
);

//*** NvWriteUINT64Data()
// This function to write back a UINT64 value. The various UINT64 values (bits,
// counters, and PINs) are kept in canonical format but manipulate in native
// format. This takes a native format value converts it and saves it back as
// in canonical format.
//
// This function will return the value from NV or RAM depending on the type of the
// index (orderly or not)
//
TPM_RC
NvWriteUINT64Data(
    NV_INDEX        *nvIndex,       // IN: the description of the index
    UINT64           intValue       // IN: the value to write
);

//*** NvGetIndexName()
// This function computes the Name of an index
// The 'name' buffer receives the bytes of the Name and the return value
// is the number of octets in the Name.
//
// This function requires that the NV Index is defined.
TPM2B_NAME *
NvGetIndexName(
    NV_INDEX        *nvIndex,       // IN: the index over which the name is to be
                                    //     computed
    TPM2B_NAME      *name           // OUT: name of the index
);

//*** NvGetNameByIndexHandle()
// This function is used to compute the Name of an NV Index referenced by handle.
//
// The 'name' buffer receives the bytes of the Name and the return value
// is the number of octets in the Name.
//
// This function requires that the NV Index is defined.
TPM2B_NAME *
NvGetNameByIndexHandle(
    TPMI_RH_NV_INDEX     handle,        // IN: handle of the index
    TPM2B_NAME          *name           // OUT: name of the index
);

//*** NvDefineIndex()
// This function is used to assign NV memory to an NV Index.
//
//  Return Type: TPM_RC
//      TPM_RC_NV_SPACE         insufficient NV space
TPM_RC
NvDefineIndex(
    TPMS_NV_PUBLIC  *publicArea,    // IN: A template for an area to create.
    TPM2B_AUTH      *authValue      // IN: The initial authorization value
);

//*** NvAddEvictObject()
// This function is used to assign NV memory to a persistent object.
//  Return Type: TPM_RC
//      TPM_RC_NV_HANDLE        the requested handle is already in use
//      TPM_RC_NV_SPACE         insufficient NV space
TPM_RC
NvAddEvictObject(
    TPMI_DH_OBJECT   evictHandle,   // IN: new evict handle
    OBJECT          *object         // IN: object to be added
);

//*** NvDeleteIndex()
// This function is used to delete an NV Index.
//  Return Type: TPM_RC
//      TPM_RC_NV_UNAVAILABLE   NV is not accessible
//      TPM_RC_NV_RATE          NV is rate limiting
TPM_RC
NvDeleteIndex(
    NV_INDEX        *nvIndex,       // IN: an in RAM index descriptor
    NV_REF           entityAddr     // IN: location in NV
);

TPM_RC
NvDeleteEvict(
    TPM_HANDLE       handle         // IN: handle of entity to be deleted
);

//*** NvFlushHierarchy()
// This function will delete persistent objects belonging to the indicated hierarchy.
// If the storage hierarchy is selected, the function will also delete any
// NV Index defined using ownerAuth.
//  Return Type: TPM_RC
//      TPM_RC_NV_RATE           NV is unavailable because of rate limit
//      TPM_RC_NV_UNAVAILABLE    NV is inaccessible
TPM_RC
NvFlushHierarchy(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy to be flushed.
);

//*** NvSetGlobalLock()
// This function is used to SET the TPMA_NV_WRITELOCKED attribute for all
// NV indexes that have TPMA_NV_GLOBALLOCK SET. This function is use by
// TPM2_NV_GlobalWriteLock().
//  Return Type: TPM_RC
//      TPM_RC_NV_RATE           NV is unavailable because of rate limit
//      TPM_RC_NV_UNAVAILABLE    NV is inaccessible
TPM_RC
NvSetGlobalLock(
    void
);

//*** NvCapGetPersistent()
// This function is used to get a list of handles of the persistent objects,
// starting at 'handle'.
//
// 'Handle' must be in valid persistent object handle range, but does not
// have to reference an existing persistent object.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
NvCapGetPersistent(
    TPMI_DH_OBJECT   handle,        // IN: start handle
    UINT32           count,         // IN: maximum number of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
);

//*** NvCapGetIndex()
// This function returns a list of handles of NV indexes, starting from 'handle'.
// 'Handle' must be in the range of NV indexes, but does not have to reference
// an existing NV Index.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles to report
//      NO          all the available handles has been reported
TPMI_YES_NO
NvCapGetIndex(
    TPMI_DH_OBJECT   handle,        // IN: start handle
    UINT32           count,         // IN: max number of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
);

//*** NvCapGetIndexNumber()
// This function returns the count of NV Indexes currently defined.
UINT32
NvCapGetIndexNumber(
    void
);

//*** NvCapGetPersistentNumber()
// Function returns the count of persistent objects currently in NV memory.
UINT32
NvCapGetPersistentNumber(
    void
);

//*** NvCapGetPersistentAvail()
// This function returns an estimate of the number of additional persistent
// objects that could be loaded into NV memory.
UINT32
NvCapGetPersistentAvail(
    void
);

//*** NvCapGetCounterNumber()
// Get the number of defined NV Indexes that are counter indexes.
UINT32
NvCapGetCounterNumber(
    void
);

//*** NvEntityStartup()
//  This function is called at TPM_Startup(). If the startup completes
//  a TPM Resume cycle, no action is taken. If the startup is a TPM Reset
//  or a TPM Restart, then this function will:
//  1. clear read/write lock;
//  2. reset NV Index data that has TPMA_NV_CLEAR_STCLEAR SET; and
//  3. set the lower bits in orderly counters to 1 for a non-orderly startup
//
//  It is a prerequisite that NV be available for writing before this
//  function is called.
BOOL
NvEntityStartup(
    STARTUP_TYPE     type           // IN: start up type
);

//*** NvCapGetCounterAvail()
// This function returns an estimate of the number of additional counter type
// NV indexes that can be defined.
UINT32
NvCapGetCounterAvail(
    void
);

//*** NvFindHandle()
// this function returns the offset in NV memory of the entity associated
// with the input handle.  A value of zero indicates that handle does not
//  exist reference an existing persistent object or defined NV Index.
NV_REF
NvFindHandle(
    TPM_HANDLE       handle
);

//*** NvReadMaxCount()
// This function returns the max NV counter value.
//
UINT64
NvReadMaxCount(
    void
);

//*** NvUpdateMaxCount()
// This function updates the max counter value to NV memory. This is just staging
// for the actual write that will occur when the NV index memory is modified.
//
void
NvUpdateMaxCount(
    UINT64           count
);

//*** NvSetMaxCount()
// This function is used at NV initialization time to set the initial value of
// the maximum counter.
void
NvSetMaxCount(
    UINT64          value
);

//*** NvGetMaxCount()
// Function to get the NV max counter value from the end-of-list marker
UINT64
NvGetMaxCount(
    void
);

#endif  // _NV_DYNAMIC_FP_H_
