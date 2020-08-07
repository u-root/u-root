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
//** Includes

#include "Tpm.h"
#include "Context_spt_fp.h"

//** Functions

//*** ComputeContextProtectionKey()
// This function retrieves the symmetric protection key for context encryption
// It is used by TPM2_ConextSave and TPM2_ContextLoad to create the symmetric
// encryption key and iv
/*(See part 1 specification)
    KDFa is used to generate the symmetric encryption key and IV. The parameters
    of the call are:
        Symkey = KDFa(hashAlg, hProof, vendorString, sequence, handle, bits)
    where
    hashAlg         a vendor-defined hash algorithm
    hProof          the hierarchy proof as selected by the hierarchy parameter
                    of the TPMS_CONTEXT
    vendorString    a value used to differentiate the uses of the KDF
    sequence        the sequence parameter of the TPMS_CONTEXT
    handle          the handle parameter of the TPMS_CONTEXT
    bits            the number of bits needed for a symmetric key and IV for
                    the context encryption
*/
//  Return Type: void
void
ComputeContextProtectionKey(
    TPMS_CONTEXT    *contextBlob,   // IN: context blob
    TPM2B_SYM_KEY   *symKey,        // OUT: the symmetric key
    TPM2B_IV        *iv             // OUT: the IV.
    )
{
    UINT16           symKeyBits;    // number of bits in the parent's
                                    //   symmetric key
    TPM2B_PROOF     *proof = NULL;  // the proof value to use. Is null for
                                    //   everything but a primary object in
                                    //   the Endorsement Hierarchy

    BYTE             kdfResult[sizeof(TPMU_HA) * 2];// Value produced by the KDF

    TPM2B_DATA       sequence2B, handle2B;

    // Get proof value
    proof = HierarchyGetProof(contextBlob->hierarchy);

    // Get sequence value in 2B format
    sequence2B.t.size = sizeof(contextBlob->sequence);
    cAssert(sizeof(contextBlob->sequence) <= sizeof(sequence2B.t.buffer));
    MemoryCopy(sequence2B.t.buffer, &contextBlob->sequence,
               sizeof(contextBlob->sequence));

    // Get handle value in 2B format
    handle2B.t.size = sizeof(contextBlob->savedHandle);
    cAssert(sizeof(contextBlob->savedHandle) <= sizeof(handle2B.t.buffer));
    MemoryCopy(handle2B.t.buffer, &contextBlob->savedHandle,
               sizeof(contextBlob->savedHandle));

    // Get the symmetric encryption key size
    symKey->t.size = CONTEXT_ENCRYPT_KEY_BYTES;
    symKeyBits = CONTEXT_ENCRYPT_KEY_BITS;
    // Get the size of the IV for the algorithm
    iv->t.size = CryptGetSymmetricBlockSize(CONTEXT_ENCRYPT_ALG, symKeyBits);

    // KDFa to generate symmetric key and IV value
    CryptKDFa(CONTEXT_INTEGRITY_HASH_ALG, &proof->b, CONTEXT_KEY, &sequence2B.b,
              &handle2B.b, (symKey->t.size + iv->t.size) * 8, kdfResult, NULL, 
              FALSE);

         // Copy part of the returned value as the key
    pAssert(symKey->t.size <= sizeof(symKey->t.buffer));
    MemoryCopy(symKey->t.buffer, kdfResult, symKey->t.size);

    // Copy the rest as the IV
    pAssert(iv->t.size <= sizeof(iv->t.buffer));
    MemoryCopy(iv->t.buffer, &kdfResult[symKey->t.size], iv->t.size);

    return;
}

//*** ComputeContextIntegrity()
// Generate the integrity hash for a context
//       It is used by TPM2_ContextSave to create an integrity hash
//       and by TPM2_ContextLoad to compare an integrity hash
/*(See part 1 specification)
    The HMAC integrity computation for a saved context is:
    HMACvendorAlg(hProof, resetValue {|| clearCount} || sequence || handle ||
                encContext)
    where
    HMACvendorAlg       HMAC using a vendor-defined hash algorithm
    hProof              the hierarchy proof as selected by the hierarchy
                        parameter of the TPMS_CONTEXT
    resetValue          either a counter value that increments on each TPM Reset
                        and is not reset over the lifetime of the TPM or a random
                        value that changes on each TPM Reset and has the size of
                        the digest produced by vendorAlg
    clearCount          a counter value that is incremented on each TPM Reset
                        or TPM Restart. This value is only included if the handle
                        value is 0x80000002.
    sequence            the sequence parameter of the TPMS_CONTEXT
    handle              the handle parameter of the TPMS_CONTEXT
    encContext          the encrypted context blob
*/
//  Return Type: void
void
ComputeContextIntegrity(
    TPMS_CONTEXT    *contextBlob,   // IN: context blob
    TPM2B_DIGEST    *integrity      // OUT: integrity
    )
{
    HMAC_STATE          hmacState;
    TPM2B_PROOF        *proof;
    UINT16              integritySize;

    // Get proof value
    proof = HierarchyGetProof(contextBlob->hierarchy);

    // Start HMAC
    integrity->t.size = CryptHmacStart2B(&hmacState, CONTEXT_INTEGRITY_HASH_ALG,
                                         &proof->b);

    // Compute integrity size at the beginning of context blob
    integritySize = sizeof(integrity->t.size) + integrity->t.size;

    // Adding total reset counter so that the context cannot be
    // used after a TPM Reset
    CryptDigestUpdateInt(&hmacState.hashState, sizeof(gp.totalResetCount),
                         gp.totalResetCount);

    // If this is a ST_CLEAR object, add the clear count
    // so that this contest cannot be loaded after a TPM Restart
    if(contextBlob->savedHandle == 0x80000002)
        CryptDigestUpdateInt(&hmacState.hashState, sizeof(gr.clearCount),
                             gr.clearCount);

    // Adding sequence number to the HMAC to make sure that it doesn't
    // get changed
    CryptDigestUpdateInt(&hmacState.hashState, sizeof(contextBlob->sequence),
                         contextBlob->sequence);

    // Protect the handle
    CryptDigestUpdateInt(&hmacState.hashState, sizeof(contextBlob->savedHandle),
                         contextBlob->savedHandle);

    // Adding sensitive contextData, skip the leading integrity area
    CryptDigestUpdate(&hmacState.hashState,
                      contextBlob->contextBlob.t.size - integritySize,
                      contextBlob->contextBlob.t.buffer + integritySize);

    // Complete HMAC
    CryptHmacEnd2B(&hmacState, &integrity->b);

    return;
}

//*** SequenceDataExport();
// This function is used scan through the sequence object and
// either modify the hash state data for export (contextSave) or to
// import it into the internal format (contextLoad).
// This function should only be called after the sequence object has been copied
// to the context buffer (contextSave) or from the context buffer into the sequence
// object. The presumption is that the context buffer version of the data is the
// same size as the internal representation so nothing outsize of the hash context
// area gets modified.
void
SequenceDataExport(
    HASH_OBJECT         *object,        // IN: an internal hash object
    HASH_OBJECT_BUFFER  *exportObject   // OUT: a sequence context in a buffer
    )
{
    // If the hash object is not an event, then only one hash context is needed
    int                   count = (object->attributes.eventSeq) ? HASH_COUNT : 1;

    for(count--; count >= 0; count--)
    {
        HASH_STATE          *hash = &object->state.hashState[count];
        size_t               offset = (BYTE *)hash - (BYTE *)object;
        BYTE                *exportHash = &((BYTE *)exportObject)[offset];

        CryptHashExportState(hash, (EXPORT_HASH_STATE *)exportHash);
    }
}

//*** SequenceDataImport();
// This function is used scan through the sequence object and
// either modify the hash state data for export (contextSave) or to
// import it into the internal format (contextLoad).
// This function should only be called after the sequence object has been copied
// to the context buffer (contextSave) or from the context buffer into the sequence
// object. The presumption is that the context buffer version of the data is the
// same size as the internal representation so nothing outsize of the hash context
// area gets modified.
void
SequenceDataImport(
    HASH_OBJECT         *object,        // IN/OUT: an internal hash object
    HASH_OBJECT_BUFFER  *exportObject   // IN/OUT: a sequence context in a buffer
    )
{
    // If the hash object is not an event, then only one hash context is needed
    int                   count = (object->attributes.eventSeq) ? HASH_COUNT : 1;

    for(count--; count >= 0; count--)
    {
        HASH_STATE          *hash = &object->state.hashState[count];
        size_t               offset = (BYTE *)hash - (BYTE *)object;
        BYTE                *importHash = &((BYTE *)exportObject)[offset];
//
        CryptHashImportState(hash, (EXPORT_HASH_STATE *)importHash);
    }
}