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
#include "ContextLoad_fp.h"

#if CC_ContextLoad  // Conditional expansion of this file

#include "Context_spt_fp.h"

/*(See part 3 specification)
// Load context
*/

//  Return Type: TPM_RC
//      TPM_RC_CONTEXT_GAP          there is only one available slot and this is not
//                                  the oldest saved session context
//      TPM_RC_HANDLE               context.savedHandle' does not reference a saved
//                                  session
//      TPM_RC_HIERARCHY            'context.hierarchy' is disabled
//      TPM_RC_INTEGRITY            'context' integrity check fail
//      TPM_RC_OBJECT_MEMORY        no free slot for an object
//      TPM_RC_SESSION_MEMORY       no free session slots
//      TPM_RC_SIZE                 incorrect context blob size
TPM_RC
TPM2_ContextLoad(
    ContextLoad_In      *in,            // IN: input parameter list
    ContextLoad_Out     *out            // OUT: output parameter list
    )
{
    TPM_RC              result;
    TPM2B_DIGEST        integrityToCompare;
    TPM2B_DIGEST        integrity;
    BYTE                *buffer;    // defined to save some typing
    INT32               size;       // defined to save some typing
    TPM_HT              handleType;
    TPM2B_SYM_KEY       symKey;
    TPM2B_IV            iv;

// Input Validation

// See discussion about the context format in TPM2_ContextSave Detailed Actions

    // IF this is a session context, make sure that the sequence number is
    // consistent with the version in the slot

    // Check context blob size
    handleType = HandleGetType(in->context.savedHandle);

    // Get integrity from context blob
    buffer = in->context.contextBlob.t.buffer;
    size = (INT32)in->context.contextBlob.t.size;
    result = TPM2B_DIGEST_Unmarshal(&integrity, &buffer, &size);
    if(result != TPM_RC_SUCCESS)
        return result;

    // the size of the integrity value has to match the size of digest produced
    // by the integrity hash
    if(integrity.t.size != CryptHashGetDigestSize(CONTEXT_INTEGRITY_HASH_ALG))
        return TPM_RCS_SIZE + RC_ContextLoad_context;

    // Make sure that the context blob has enough space for the fingerprint. This
    // is elastic pants to go with the belt and suspenders we already have to make
    // sure that the context is complete and untampered.
    if((unsigned)size < sizeof(in->context.sequence))
        return TPM_RCS_SIZE + RC_ContextLoad_context;

    // After unmarshaling the integrity value, 'buffer' is pointing at the first
    // byte of the integrity protected and encrypted buffer and 'size' is the number
    // of integrity protected and encrypted bytes.

    // Compute context integrity
    ComputeContextIntegrity(&in->context, &integrityToCompare);

    // Compare integrity
    if(!MemoryEqual2B(&integrity.b, &integrityToCompare.b))
        return TPM_RCS_INTEGRITY + RC_ContextLoad_context;
    // Compute context encryption key
    ComputeContextProtectionKey(&in->context, &symKey, &iv);

    // Decrypt context data in place
    CryptSymmetricDecrypt(buffer, CONTEXT_ENCRYPT_ALG, CONTEXT_ENCRYPT_KEY_BITS,
                          symKey.t.buffer, &iv, ALG_CFB_VALUE, size, buffer);
    // See if the fingerprint value matches. If not, it is symptomatic of either
    // a broken TPM or that the TPM is under attack so go into failure mode.
    if(!MemoryEqual(buffer, &in->context.sequence, sizeof(in->context.sequence)))
        FAIL(FATAL_ERROR_INTERNAL);

    // step over fingerprint
    buffer += sizeof(in->context.sequence);

    // set the remaining size of the context
    size -= sizeof(in->context.sequence);

    // Perform object or session specific input check
    switch(handleType)
    {
        case TPM_HT_TRANSIENT:
        {
            OBJECT      *outObject;

            if(size > (INT32)sizeof(OBJECT))
                FAIL(FATAL_ERROR_INTERNAL);

            // Discard any changes to the handle that the TRM might have made
            in->context.savedHandle = TRANSIENT_FIRST;

            // If hierarchy is disabled, no object context can be loaded in this
            // hierarchy
            if(!HierarchyIsEnabled(in->context.hierarchy))
                return TPM_RCS_HIERARCHY + RC_ContextLoad_context;

            // Restore object. If there is no empty space, indicate as much
            outObject = ObjectContextLoad((ANY_OBJECT_BUFFER *)buffer,
                                          &out->loadedHandle);
            if(outObject == NULL)
                return TPM_RC_OBJECT_MEMORY;

            break;
        }
        case TPM_HT_POLICY_SESSION:
        case TPM_HT_HMAC_SESSION:
        {
            if(size != sizeof(SESSION))
                FAIL(FATAL_ERROR_INTERNAL);

            // This command may cause the orderlyState to be cleared due to
            // the update of state reset data.  If this is the case, check if NV is
            // available first
            RETURN_IF_ORDERLY;

            // Check if input handle points to a valid saved session and that the
            // sequence number makes sense
            if(!SequenceNumberForSavedContextIsValid(&in->context))
                return TPM_RCS_HANDLE + RC_ContextLoad_context;

            // Restore session.  A TPM_RC_SESSION_MEMORY, TPM_RC_CONTEXT_GAP error
            // may be returned at this point
            result = SessionContextLoad((SESSION_BUF *)buffer,
                                        &in->context.savedHandle);
            if(result != TPM_RC_SUCCESS)
                return result;

            out->loadedHandle = in->context.savedHandle;

            // orderly state should be cleared because of the update of state
            // reset and state clear data
            g_clearOrderly = TRUE;

            break;
        }
        default:
            // Context blob may only have an object handle or a session handle.
            // All the other handle type should be filtered out at unmarshal
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_ContextLoad