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
#include "Tpm.h"
#include "EvictControl_fp.h"

#if CC_EvictControl  // Conditional expansion of this file

/*(See part 3 specification)
// Make a transient object persistent or evict a persistent object
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   an object with 'temporary', 'stClear' or 'publicOnly'
//                          attribute SET cannot be made persistent
//      TPM_RC_HIERARCHY    'auth' cannot authorize the operation in the hierarchy
//                          of 'evictObject'
//      TPM_RC_HANDLE       'evictHandle' of the persistent object to be evicted is
//                          not the same as the 'persistentHandle' argument
//      TPM_RC_NV_HANDLE    'persistentHandle' is unavailable
//      TPM_RC_NV_SPACE     no space in NV to make 'evictHandle' persistent
//      TPM_RC_RANGE        'persistentHandle' is not in the range corresponding to
//                          the hierarchy of 'evictObject'
TPM_RC
TPM2_EvictControl(
    EvictControl_In     *in             // IN: input parameter list
    )
{
    TPM_RC      result;
    OBJECT      *evictObject;

// Input Validation

    // Get internal object pointer
    evictObject = HandleToObject(in->objectHandle);

    // Temporary, stClear or public only objects can not be made persistent
    if(evictObject->attributes.temporary == SET
       || evictObject->attributes.stClear == SET
       || evictObject->attributes.publicOnly == SET)
        return TPM_RCS_ATTRIBUTES + RC_EvictControl_objectHandle;

    // If objectHandle refers to a persistent object, it should be the same as
    // input persistentHandle
    if(evictObject->attributes.evict == SET
       && evictObject->evictHandle != in->persistentHandle)
        return TPM_RCS_HANDLE + RC_EvictControl_objectHandle;

    // Additional authorization validation
    if(in->auth == TPM_RH_PLATFORM)
    {
        // To make persistent
        if(evictObject->attributes.evict == CLEAR)
        {
            // PlatformAuth can not set evict object in storage or endorsement
            // hierarchy
            if(evictObject->attributes.ppsHierarchy == CLEAR)
                return TPM_RCS_HIERARCHY + RC_EvictControl_objectHandle;
            // Platform cannot use a handle outside of platform persistent range.
            if(!NvIsPlatformPersistentHandle(in->persistentHandle))
                return TPM_RCS_RANGE + RC_EvictControl_persistentHandle;
        }
        // PlatformAuth can delete any persistent object
    }
    else if(in->auth == TPM_RH_OWNER)
    {
        // OwnerAuth can not set or clear evict object in platform hierarchy
        if(evictObject->attributes.ppsHierarchy == SET)
            return TPM_RCS_HIERARCHY + RC_EvictControl_objectHandle;

        // Owner cannot use a handle outside of owner persistent range.
        if(evictObject->attributes.evict == CLEAR
           && !NvIsOwnerPersistentHandle(in->persistentHandle))
            return TPM_RCS_RANGE + RC_EvictControl_persistentHandle;
    }
    else
    {
        // Other authorization is not allowed in this command and should have been
        // filtered out in unmarshal process
        FAIL(FATAL_ERROR_INTERNAL);
    }
// Internal Data Update
    // Change evict state
    if(evictObject->attributes.evict == CLEAR)
    {
        // Make object persistent
        if(NvFindHandle(in->persistentHandle) != 0)
            return TPM_RC_NV_DEFINED;
        // A TPM_RC_NV_HANDLE or TPM_RC_NV_SPACE error may be returned at this
        // point
        result = NvAddEvictObject(in->persistentHandle, evictObject);
    }
    else
    {
        // Delete the persistent object in NV
        result = NvDeleteEvict(evictObject->evictHandle);
    }
    return result;
}

#endif // CC_EvictControl