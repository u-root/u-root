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

// This function is called to process a _TPM_Hash_End indication.
LIB_EXPORT void
_TPM_Hash_End(
    void
    )
{
    UINT32          i;
    TPM2B_DIGEST    digest;
    HASH_OBJECT    *hashObject;
    TPMI_DH_PCR     pcrHandle;

    // If the DRTM handle is not being used, then either _TPM_Hash_Start has not
    // been called, _TPM_Hash_End was previously called, or some other command
    // was executed and the sequence was aborted.
    if(g_DRTMHandle == TPM_RH_UNASSIGNED)
        return;

    // Get DRTM sequence object
    hashObject = (HASH_OBJECT *)HandleToObject(g_DRTMHandle);

    // Is this _TPM_Hash_End after Startup or before
    if(TPMIsStarted())
    {
        // After

        // Reset the DRTM PCR
        PCRResetDynamics();

        // Extend the DRTM_PCR.
        pcrHandle = PCR_FIRST + DRTM_PCR;

        // DRTM sequence increments restartCount
        gr.restartCount++;
    }
    else
    {
        pcrHandle = PCR_FIRST + HCRTM_PCR;
        g_DrtmPreStartup = TRUE;
    }

    // Complete hash and extend PCR, or if this is an HCRTM, complete
    // the hash, reset the H-CRTM register (PCR[0]) to 0...04, and then
    // extend the H-CRTM data
    for(i = 0; i < HASH_COUNT; i++)
    {
        TPMI_ALG_HASH       hash = CryptHashGetAlgByIndex(i);
        // make sure that the PCR is implemented for this algorithm
        if(PcrIsAllocated(pcrHandle,
                          hashObject->state.hashState[i].hashAlg))
        {
            // Complete hash
            digest.t.size = CryptHashGetDigestSize(hash);
            CryptHashEnd2B(&hashObject->state.hashState[i], &digest.b);

            PcrDrtm(pcrHandle, hash, &digest);
        }
    }

    // Flush sequence object.
    FlushObject(g_DRTMHandle);

    g_DRTMHandle = TPM_RH_UNASSIGNED;


    return;
}