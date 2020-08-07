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
#include "ContextSave_fp.h"

#if CC_ContextSave  // Conditional expansion of this file

#include "Context_spt_fp.h"

/*(See part 3 specification)
 Save context
*/
//  Return Type: TPM_RC
//      TPM_RC_CONTEXT_GAP          a contextID could not be assigned for a session
//                                  context save
//      TPM_RC_TOO_MANY_CONTEXTS    no more contexts can be saved as the counter has
//                                  maxed out
TPM_RC
TPM2_ContextSave(
    ContextSave_In      *in,            // IN: input parameter list
    ContextSave_Out     *out            // OUT: output parameter list
    )
{
    TPM_RC          result = TPM_RC_SUCCESS;
    UINT16          fingerprintSize;    // The size of fingerprint in context
    // blob.
    UINT64          contextID = 0;      // session context ID
    TPM2B_SYM_KEY   symKey;
    TPM2B_IV        iv;

    TPM2B_DIGEST    integrity;
    UINT16          integritySize;
    BYTE            *buffer;

    // This command may cause the orderlyState to be cleared due to
    // the update of state reset data. If the state is orderly and
    // cannot be changed, exit early.
    RETURN_IF_ORDERLY;
    
// Internal Data Update

// This implementation does not do things in quite the same way as described in
// Part 2 of the specification. In Part 2, it indicates that the 
// TPMS_CONTEXT_DATA contains two TPM2B values. That is not how this is 
// implemented. Rather, the size field of the TPM2B_CONTEXT_DATA is used to 
// determine the amount of data in the encrypted data. That part is not 
// independently sized. This makes the actual size 2 bytes smaller than 
// calculated using Part 2. Since this is opaque to the caller, it is not 
// necessary to fix. The actual size is returned by TPM2_GetCapabilties().

    // Initialize output handle.  At the end of command action, the output
    // handle of an object will be replaced, while the output handle
    // for a session will be the same as input
    out->context.savedHandle = in->saveHandle;

    // Get the size of fingerprint in context blob.  The sequence value in
    // TPMS_CONTEXT structure is used as the fingerprint
    fingerprintSize = sizeof(out->context.sequence);

    // Compute the integrity size at the beginning of context blob
    integritySize = sizeof(integrity.t.size)
        + CryptHashGetDigestSize(CONTEXT_INTEGRITY_HASH_ALG);

// Perform object or session specific context save
    switch(HandleGetType(in->saveHandle))
    {
        case TPM_HT_TRANSIENT:
        {
            OBJECT              *object = HandleToObject(in->saveHandle);
            ANY_OBJECT_BUFFER   *outObject;
            UINT16               objectSize = ObjectIsSequence(object)
                ? sizeof(HASH_OBJECT) : sizeof(OBJECT);

            outObject = (ANY_OBJECT_BUFFER *)(out->context.contextBlob.t.buffer
                                              + integritySize + fingerprintSize);

            // Set size of the context data.  The contents of context blob is vendor
            // defined.  In this implementation, the size is size of integrity
            // plus fingerprint plus the whole internal OBJECT structure
            out->context.contextBlob.t.size = integritySize +
                fingerprintSize + objectSize;
#if ALG_RSA
            // For an RSA key, make sure that the key has had the private exponent
            // computed before saving.
            if(object->publicArea.type == TPM_ALG_RSA &&
               !(object->attributes.publicOnly))
                CryptRsaLoadPrivateExponent(&object->publicArea, &object->sensitive);
#endif
            // Make sure things fit
            pAssert(out->context.contextBlob.t.size
                    <= sizeof(out->context.contextBlob.t.buffer));
            // Copy the whole internal OBJECT structure to context blob
            MemoryCopy(outObject, object, objectSize);

            // Increment object context ID
            gr.objectContextID++;
            // If object context ID overflows, TPM should be put in failure mode
            if(gr.objectContextID == 0)
                FAIL(FATAL_ERROR_INTERNAL);

            // Fill in other return values for an object.
            out->context.sequence = gr.objectContextID;
            // For regular object, savedHandle is 0x80000000.  For sequence object,
            // savedHandle is 0x80000001.  For object with stClear, savedHandle
            // is 0x80000002
            if(ObjectIsSequence(object))
            {
                out->context.savedHandle = 0x80000001;
                SequenceDataExport((HASH_OBJECT *)object,
                                   (HASH_OBJECT_BUFFER *)outObject);
            }
            else
                out->context.savedHandle = (object->attributes.stClear == SET)
                ? 0x80000002 : 0x80000000;
// Get object hierarchy
            out->context.hierarchy = ObjectGetHierarchy(object);

            break;
        }
        case TPM_HT_HMAC_SESSION:
        case TPM_HT_POLICY_SESSION:
        {
            SESSION         *session = SessionGet(in->saveHandle);

            // Set size of the context data.  The contents of context blob is vendor
            // defined.  In this implementation, the size of context blob is the
            // size of a internal session structure plus the size of
            // fingerprint plus the size of integrity
            out->context.contextBlob.t.size = integritySize +
                fingerprintSize + sizeof(*session);

            // Make sure things fit
            pAssert(out->context.contextBlob.t.size
                    < sizeof(out->context.contextBlob.t.buffer));

            // Copy the whole internal SESSION structure to context blob.
            // Save space for fingerprint at the beginning of the buffer
            // This is done before anything else so that the actual context
            // can be reclaimed after this call
            pAssert(sizeof(*session) <= sizeof(out->context.contextBlob.t.buffer)
                    - integritySize - fingerprintSize);
            MemoryCopy(out->context.contextBlob.t.buffer + integritySize 
                       + fingerprintSize, session, sizeof(*session));
           // Fill in the other return parameters for a session
           // Get a context ID and set the session tracking values appropriately
           // TPM_RC_CONTEXT_GAP is a possible error.
           // SessionContextSave() will flush the in-memory context
           // so no additional errors may occur after this call.
            result = SessionContextSave(out->context.savedHandle, &contextID);
            if(result != TPM_RC_SUCCESS)
                return result;
            // sequence number is the current session contextID
            out->context.sequence = contextID;

            // use TPM_RH_NULL as hierarchy for session context
            out->context.hierarchy = TPM_RH_NULL;

            break;
        }
        default:
            // SaveContext may only take an object handle or a session handle.
            // All the other handle type should be filtered out at unmarshal
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }

    // Save fingerprint at the beginning of encrypted area of context blob.
    // Reserve the integrity space
    pAssert(sizeof(out->context.sequence) <=
            sizeof(out->context.contextBlob.t.buffer) - integritySize);
    MemoryCopy(out->context.contextBlob.t.buffer + integritySize,
               &out->context.sequence, sizeof(out->context.sequence));

    // Compute context encryption key
    ComputeContextProtectionKey(&out->context, &symKey, &iv);

    // Encrypt context blob
    CryptSymmetricEncrypt(out->context.contextBlob.t.buffer + integritySize,
                          CONTEXT_ENCRYPT_ALG, CONTEXT_ENCRYPT_KEY_BITS,
                          symKey.t.buffer, &iv, ALG_CFB_VALUE,
                          out->context.contextBlob.t.size - integritySize,
                          out->context.contextBlob.t.buffer + integritySize);

    // Compute integrity hash for the object
    // In this implementation, the same routine is used for both sessions
    // and objects.
    ComputeContextIntegrity(&out->context, &integrity);

    // add integrity at the beginning of context blob
    buffer = out->context.contextBlob.t.buffer;
    TPM2B_DIGEST_Marshal(&integrity, &buffer, NULL);

    // orderly state should be cleared because of the update of state reset and
    // state clear data
    g_clearOrderly = TRUE;

    return result;
}

#endif // CC_ContextSave