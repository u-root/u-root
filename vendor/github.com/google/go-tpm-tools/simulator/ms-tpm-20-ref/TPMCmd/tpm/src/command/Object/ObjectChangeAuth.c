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
#include "ObjectChangeAuth_fp.h"

#if CC_ObjectChangeAuth  // Conditional expansion of this file

#include "Object_spt_fp.h"

/*(See part 3 specification)
// Create an object
*/
//  Return Type: TPM_RC
//      TPM_RC_SIZE             'newAuth' is larger than the size of the digest
//                              of the Name algorithm of 'objectHandle'
//      TPM_RC_TYPE             the key referenced by 'parentHandle' is not the
//                              parent of the object referenced by 'objectHandle';
//                              or 'objectHandle' is a sequence object.
TPM_RC
TPM2_ObjectChangeAuth(
    ObjectChangeAuth_In     *in,            // IN: input parameter list
    ObjectChangeAuth_Out    *out            // OUT: output parameter list
    )
{
    TPMT_SENSITIVE           sensitive;

    OBJECT                  *object = HandleToObject(in->objectHandle);
    TPM2B_NAME               QNCompare;

// Input Validation

    // Can not change authorization on sequence object
    if(ObjectIsSequence(object))
        return TPM_RCS_TYPE + RC_ObjectChangeAuth_objectHandle;

    // Make sure that the authorization value is consistent with the nameAlg
    if(!AdjustAuthSize(&in->newAuth, object->publicArea.nameAlg))
        return TPM_RCS_SIZE + RC_ObjectChangeAuth_newAuth;

    // Parent handle should be the parent of object handle.  In this
    // implementation we verify this by checking the QN of object.  Other
    // implementation may choose different method to verify this attribute.
    ComputeQualifiedName(in->parentHandle,
                         object->publicArea.nameAlg,
                         &object->name, &QNCompare);
    if(!MemoryEqual2B(&object->qualifiedName.b, &QNCompare.b))
        return TPM_RCS_TYPE + RC_ObjectChangeAuth_parentHandle;

// Command Output
    // Prepare the sensitive area with the new authorization value
    sensitive = object->sensitive;
    sensitive.authValue = in->newAuth;

    // Protect the sensitive area
    SensitiveToPrivate(&sensitive, &object->name, HandleToObject(in->parentHandle),
                       object->publicArea.nameAlg,
                       &out->outPrivate);
    return TPM_RC_SUCCESS;
}

#endif // CC_ObjectChangeAuth