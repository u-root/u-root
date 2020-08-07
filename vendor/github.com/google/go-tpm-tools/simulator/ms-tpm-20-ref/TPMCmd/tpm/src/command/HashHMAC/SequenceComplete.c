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
#include "SequenceComplete_fp.h"

#if CC_SequenceComplete  // Conditional expansion of this file

/*(See part 3 specification)
// Complete a sequence and flush the object.
*/
//  Return Type: TPM_RC
//      TPM_RC_MODE             'sequenceHandle' does not reference a hash or HMAC
//                              sequence object
TPM_RC
TPM2_SequenceComplete(
    SequenceComplete_In     *in,            // IN: input parameter list
    SequenceComplete_Out    *out            // OUT: output parameter list
    )
{
    HASH_OBJECT                      *hashObject;
// Input validation
    // Get hash object pointer
    hashObject = (HASH_OBJECT *)HandleToObject(in->sequenceHandle);

    // input handle must be a hash or HMAC sequence object.
    if(hashObject->attributes.hashSeq == CLEAR
       && hashObject->attributes.hmacSeq == CLEAR)
        return TPM_RCS_MODE + RC_SequenceComplete_sequenceHandle;
// Command Output
    if(hashObject->attributes.hashSeq == SET)           // sequence object for hash
    {
       // Get the hash algorithm before the algorithm is lost in CryptHashEnd
        TPM_ALG_ID       hashAlg = hashObject->state.hashState[0].hashAlg;

        // Update last piece of the data
        CryptDigestUpdate2B(&hashObject->state.hashState[0], &in->buffer.b);

        // Complete hash
        out->result.t.size = CryptHashEnd(&hashObject->state.hashState[0], 
                                         sizeof(out->result.t.buffer),
                                         out->result.t.buffer);
        // Check if the first block of the sequence has been received
        if(hashObject->attributes.firstBlock == CLEAR)
        {
            // If not, then this is the first block so see if it is 'safe'
            // to sign.
            if(TicketIsSafe(&in->buffer.b))
                hashObject->attributes.ticketSafe = SET;
        }
        // Output ticket
        out->validation.tag = TPM_ST_HASHCHECK;
        out->validation.hierarchy = in->hierarchy;

        if(in->hierarchy == TPM_RH_NULL)
        {
            // Ticket is not required
            out->validation.digest.t.size = 0;
        }
        else if(hashObject->attributes.ticketSafe == CLEAR)
        {
            // Ticket is not safe to generate
            out->validation.hierarchy = TPM_RH_NULL;
            out->validation.digest.t.size = 0;
        }
        else
        {
            // Compute ticket
            TicketComputeHashCheck(out->validation.hierarchy, hashAlg,
                                   &out->result, &out->validation);
        }
    }
    else
    {
        //   Update last piece of data
        CryptDigestUpdate2B(&hashObject->state.hmacState.hashState, &in->buffer.b);
#if !SMAC_IMPLEMENTED
        // Complete HMAC
        out->result.t.size = CryptHmacEnd(&(hashObject->state.hmacState),
                                          sizeof(out->result.t.buffer),
                                          out->result.t.buffer);
#else
        // Complete the MAC
        out->result.t.size = CryptMacEnd(&hashObject->state.hmacState, 
                                         sizeof(out->result.t.buffer),
                                         out->result.t.buffer);
#endif
        // No ticket is generated for HMAC sequence
        out->validation.tag = TPM_ST_HASHCHECK;
        out->validation.hierarchy = TPM_RH_NULL;
        out->validation.digest.t.size = 0;
    }
// Internal Data Update
    // mark sequence object as evict so it will be flushed on the way out
    hashObject->attributes.evict = SET;

    return TPM_RC_SUCCESS;
}

#endif // CC_SequenceComplete