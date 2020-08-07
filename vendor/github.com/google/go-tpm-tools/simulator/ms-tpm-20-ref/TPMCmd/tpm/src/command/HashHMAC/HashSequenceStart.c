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
#include "HashSequenceStart_fp.h"

#if CC_HashSequenceStart  // Conditional expansion of this file

/*(See part 3 specification)
// Start a hash or an event sequence
*/
//  Return Type: TPM_RC
//      TPM_RC_OBJECT_MEMORY        no space to create an internal object
TPM_RC
TPM2_HashSequenceStart(
    HashSequenceStart_In    *in,            // IN: input parameter list
    HashSequenceStart_Out   *out            // OUT: output parameter list
    )
{
// Internal Data Update

    if(in->hashAlg == TPM_ALG_NULL)
        // Start a event sequence.  A TPM_RC_OBJECT_MEMORY error may be
        // returned at this point
        return ObjectCreateEventSequence(&in->auth, &out->sequenceHandle);

    // Start a hash sequence.  A TPM_RC_OBJECT_MEMORY error may be
    // returned at this point
    return ObjectCreateHashSequence(in->hashAlg, &in->auth, &out->sequenceHandle);
}

#endif // CC_HashSequenceStart