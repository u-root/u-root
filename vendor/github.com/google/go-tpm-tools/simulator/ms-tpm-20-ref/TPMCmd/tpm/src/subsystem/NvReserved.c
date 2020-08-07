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
//** Introduction

// The NV memory is divided into two areas: dynamic space for user defined NV
// Indices and evict objects, and reserved space for TPM persistent and state save
// data.
//
// The entries in dynamic space are a linked list of entries. Each entry has, as its
// first field, a size. If the size field is zero, it marks the end of the
// list.
//
// An allocation of an Index or evict object may use almost all of the remaining
// NV space such that the size field will not fit. The functions that search the
// list are aware of this and will terminate the search if they either find a zero
// size or recognize that there is insufficient space for the size field.
//
// An Index allocation will contain an NV_INDEX structure. If the Index does not
// have the orderly attribute, the NV_INDEX is followed immediately by the NV data.
//
// An evict object entry contains a handle followed by an OBJECT structure. This
// results in both the Index and Evict Object having an identifying handle as the
// first field following the size field.
//
// When an Index has the orderly attribute, the data is kept in RAM. This RAM is
// saved to backing store in NV memory on any orderly shutdown. The entries in
// orderly memory are also a linked list using a size field as the first entry. As
// with the NV memory, the list is terminated by a zero size field or when the last
// entry leaves insufficient space for the terminating size field.
//
// The attributes of an orderly index are maintained in RAM memory in order to
// reduce the number of NV writes needed for orderly data. When an orderly index
// is created, an entry is made in the dynamic NV memory space that holds the Index
// authorizations (authPolicy and authValue) and the size of the data. This entry is
// only modified if the authValue  of the index is changed. The more volatile data
// of the index is kept in RAM. When an orderly Index is created or deleted, the
// RAM data is copied to NV backing store so that the image in the backing store
// matches the layout of RAM. In normal operation. The RAM data is also copied on
// any orderly shutdown. In normal operation, the only other reason for writing
// to the backing store for RAM is when a counter is first written (TPMA_NV_WRITTEN
// changes from CLEAR to SET) or when a counter "rolls over."
//
// Static space contains items that are individually modifiable. The values are in
// the 'gp' PERSISTEND_DATA structure in RAM and mapped to locations in NV.
//

//** Includes, Defines
#define NV_C
#include    "Tpm.h"

//************************************************
//** Functions
//************************************************


//*** NvInitStatic()
// This function initializes the static variables used in the NV subsystem.
static void
NvInitStatic(
    void
    )
{
    // In some implementations, the end of NV is variable and is set at boot time.
    // This value will be the same for each boot, but is not necessarily known
    // at compile time.
    s_evictNvEnd = (NV_REF)NV_MEMORY_SIZE;
    return;
}

//*** NvCheckState()
// Function to check the NV state by accessing the platform-specific function
// to get the NV state.  The result state is registered in s_NvIsAvailable
// that will be reported by NvIsAvailable.
//
// This function is called at the beginning of ExecuteCommand before any potential
// check of g_NvStatus.
void
NvCheckState(
    void
    )
{
    int     func_return;
//
    func_return = _plat__IsNvAvailable();
    if(func_return == 0)
        g_NvStatus = TPM_RC_SUCCESS;
    else if(func_return == 1)
        g_NvStatus = TPM_RC_NV_UNAVAILABLE;
    else
        g_NvStatus = TPM_RC_NV_RATE;
    return;
}

//*** NvCommit
// This is a wrapper for the platform function to commit pending NV writes.
BOOL
NvCommit(
    void
    )
{
    return (_plat__NvCommit() == 0);
}

//*** NvPowerOn()
//  This function is called at _TPM_Init to initialize the NV environment.
//  Return Type: BOOL
//      TRUE(1)         all NV was initialized
//      FALSE(0)        the NV containing saved state had an error and 
//                      TPM2_Startup(CLEAR) is required
BOOL
NvPowerOn(
    void
    )
{
    int          nvError = 0;
    // If power was lost, need to re-establish the RAM data that is loaded from
    // NV and initialize the static variables
    if(g_powerWasLost)
    {
        if((nvError = _plat__NVEnable(0)) < 0)
            FAIL(FATAL_ERROR_NV_UNRECOVERABLE);
        NvInitStatic();
    }
    return nvError == 0;
}

//*** NvManufacture()
// This function initializes the NV system at pre-install time.
//
// This function should only be called in a manufacturing environment or in a
// simulation.
//
// The layout of NV memory space is an implementation choice.
void
NvManufacture(
    void
    )
{
#if SIMULATION
    // Simulate the NV memory being in the erased state.
    _plat__NvMemoryClear(0, NV_MEMORY_SIZE);
#endif
    // Initialize static variables
    NvInitStatic();
    // Clear the RAM used for Orderly Index data
    MemorySet(s_indexOrderlyRam, 0, RAM_INDEX_SPACE);
    // Write that Orderly Index data to NV
    NvUpdateIndexOrderlyData();
    // Initialize the next offset of the first entry in evict/index list to 0 (the
    // end of list marker) and the initial s_maxCounterValue;
    NvSetMaxCount(0);
    // Put the end of list marker at the end of memory. This contains the MaxCount
    // value as well as the end marker.
    NvWriteNvListEnd(NV_USER_DYNAMIC);
    return;
}

//*** NvRead()
// This function is used to move reserved data from NV memory to RAM.
void
NvRead(
    void            *outBuffer,     // OUT: buffer to receive data
    UINT32           nvOffset,      // IN: offset in NV of value
    UINT32           size           // IN: size of the value to read
    )
{
    // Input type should be valid
    pAssert(nvOffset + size < NV_MEMORY_SIZE);
    _plat__NvMemoryRead(nvOffset, size, outBuffer);
    return;
}

//*** NvWrite()
// This function is used to post reserved data for writing to NV memory. Before
// the TPM completes the operation, the value will be written.
BOOL
NvWrite(
    UINT32           nvOffset,      // IN: location in NV to receive data
    UINT32           size,          // IN: size of the data to move
    void            *inBuffer       // IN: location containing data to write
    )
{
    // Input type should be valid
    if(nvOffset + size <= NV_MEMORY_SIZE)
    {
    // Set the flag that a NV write happened
    SET_NV_UPDATE(UT_NV);
        return _plat__NvMemoryWrite(nvOffset, size, inBuffer);
    }
    return FALSE;
}

//*** NvUpdatePersistent()
// This function is used to update a value in the PERSISTENT_DATA structure and
// commits the value to NV.
void
NvUpdatePersistent(
    UINT32           offset,        // IN: location in PERMANENT_DATA to be updated
    UINT32           size,          // IN: size of the value
    void            *buffer         // IN: the new data
    )
{
    pAssert(offset + size <= sizeof(gp));
    MemoryCopy(&gp + offset, buffer, size);
    NvWrite(offset, size, buffer);
}

//*** NvClearPersistent()
// This function is used to clear a persistent data entry and commit it to NV
void
NvClearPersistent(
    UINT32           offset,        // IN: the offset in the PERMANENT_DATA
                                    //     structure to be cleared (zeroed)
    UINT32           size           // IN: number of bytes to clear
    )
{
    pAssert(offset + size <= sizeof(gp));
    MemorySet((&gp) + offset, 0, size);
    NvWrite(offset, size, (&gp) + offset);
}

//*** NvReadPersistent()
// This function reads persistent data to the RAM copy of the 'gp' structure.
void
NvReadPersistent(
    void
    )
{
    NvRead(&gp, NV_PERSISTENT_DATA, sizeof(gp));
    return;
}