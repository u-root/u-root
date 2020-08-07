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
#include "Sign_fp.h"

#if CC_Sign  // Conditional expansion of this file

#include "Attest_spt_fp.h"

/*(See part 3 specification)
// sign an externally provided hash using an asymmetric signing key
*/
//  Return Type: TPM_RC
//      TPM_RC_BINDING          The public and private portions of the key are not
//                              properly bound.
//      TPM_RC_KEY              'signHandle' does not reference a signing key;
//      TPM_RC_SCHEME           the scheme is not compatible with sign key type,
//                              or input scheme is not compatible with default
//                              scheme, or the chosen scheme is not a valid
//                              sign scheme
//      TPM_RC_TICKET           'validation' is not a valid ticket
//      TPM_RC_VALUE            the value to sign is larger than allowed for the
//                              type of 'keyHandle'

TPM_RC
TPM2_Sign(
    Sign_In         *in,            // IN: input parameter list
    Sign_Out        *out            // OUT: output parameter list
    )
{
    TPM_RC                   result;
    TPMT_TK_HASHCHECK        ticket;
    OBJECT                  *signObject = HandleToObject(in->keyHandle);
//
// Input Validation
    if(!IsSigningObject(signObject))
        return TPM_RCS_KEY + RC_Sign_keyHandle;

    // A key that will be used for x.509 signatures can't be used in TPM2_Sign().
    if(IS_ATTRIBUTE(signObject->publicArea.objectAttributes, TPMA_OBJECT, x509sign))
        return TPM_RCS_ATTRIBUTES + RC_Sign_keyHandle;

    // pick a scheme for sign.  If the input sign scheme is not compatible with
    // the default scheme, return an error.
    if(!CryptSelectSignScheme(signObject, &in->inScheme))
        return TPM_RCS_SCHEME + RC_Sign_inScheme;

    // If validation is provided, or the key is restricted, check the ticket
    if(in->validation.digest.t.size != 0
       || IS_ATTRIBUTE(signObject->publicArea.objectAttributes, 
                       TPMA_OBJECT, restricted))
    {
        // Compute and compare ticket
        TicketComputeHashCheck(in->validation.hierarchy,
                               in->inScheme.details.any.hashAlg,
                               &in->digest, &ticket);

        if(!MemoryEqual2B(&in->validation.digest.b, &ticket.digest.b))
            return TPM_RCS_TICKET + RC_Sign_validation;
    }
    else
    // If we don't have a ticket, at least verify that the provided 'digest'
    // is the size of the scheme hashAlg digest.
    // NOTE: this does not guarantee that the 'digest' is actually produced using
    // the indicated hash algorithm, but at least it might be.
    {
        if(in->digest.t.size
           != CryptHashGetDigestSize(in->inScheme.details.any.hashAlg))
            return TPM_RCS_SIZE + RC_Sign_digest;
    }

// Command Output
    // Sign the hash. A TPM_RC_VALUE or TPM_RC_SCHEME
    // error may be returned at this point
    result = CryptSign(signObject, &in->inScheme, &in->digest, &out->signature);

    return result;
}

#endif // CC_Sign