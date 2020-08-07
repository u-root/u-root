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
#include "Commit_fp.h"

#if CC_Commit  // Conditional expansion of this file

/*(See part 3 specification)
// This command performs the point multiply operations for anonymous signing
// scheme.
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       'keyHandle' references a restricted key that is not a
//                              signing key
//      TPM_RC_ECC_POINT        either 'P1' or the point derived from 's2' is not on
//                              the curve of 'keyHandle'
//      TPM_RC_HASH             invalid name algorithm in 'keyHandle'
//      TPM_RC_KEY              'keyHandle' does not reference an ECC key
//      TPM_RC_SCHEME           the scheme of 'keyHandle' is not an anonymous scheme
//      TPM_RC_NO_RESULT        'K', 'L' or 'E' was a point at infinity; or
//                              failed to generate "r" value
//      TPM_RC_SIZE             's2' is empty but 'y2' is not or 's2' provided but
//                              'y2' is not
TPM_RC
TPM2_Commit(
    Commit_In       *in,            // IN: input parameter list
    Commit_Out      *out            // OUT: output parameter list
    )
{
    OBJECT                  *eccKey;
    TPMS_ECC_POINT           P2;
    TPMS_ECC_POINT          *pP2 = NULL;
    TPMS_ECC_POINT          *pP1 = NULL;
    TPM2B_ECC_PARAMETER      r;
    TPM2B_ECC_PARAMETER      p;
    TPM_RC                   result;
    TPMS_ECC_PARMS          *parms;

// Input Validation

    eccKey = HandleToObject(in->signHandle);
    parms = &eccKey->publicArea.parameters.eccDetail;

    // Input key must be an ECC key
    if(eccKey->publicArea.type != TPM_ALG_ECC)
        return TPM_RCS_KEY + RC_Commit_signHandle;

    // This command may only be used with a sign-only key using an anonymous
    // scheme.
    // NOTE: a sign + decrypt key has no scheme so it will not be an anonymous one
    // and an unrestricted sign key might no have a signing scheme but it can't
    // be use in Commit()
    if(!CryptIsSchemeAnonymous(parms->scheme.scheme))
        return TPM_RCS_SCHEME + RC_Commit_signHandle;

// Make sure that both parts of P2 are present if either is present
    if((in->s2.t.size == 0) != (in->y2.t.size == 0))
        return TPM_RCS_SIZE + RC_Commit_y2;

    // Get prime modulus for the curve. This is needed later but getting this now
    // allows confirmation that the curve exists.
    if(!CryptEccGetParameter(&p, 'p', parms->curveID))
        return TPM_RCS_KEY + RC_Commit_signHandle;

    // Get the random value that will be used in the point multiplications
    // Note: this does not commit the count.
    if(!CryptGenerateR(&r, NULL, parms->curveID, &eccKey->name))
        return TPM_RC_NO_RESULT;

    // Set up P2 if s2 and Y2 are provided
    if(in->s2.t.size != 0)
    {
        TPM2B_DIGEST             x2;

        pP2 = &P2;

        // copy y2 for P2
        P2.y = in->y2;

        // Compute x2  HnameAlg(s2) mod p
        //      do the hash operation on s2 with the size of curve 'p'
        x2.t.size = CryptHashBlock(eccKey->publicArea.nameAlg,
                                     in->s2.t.size,
                                     in->s2.t.buffer,
                                     sizeof(x2.t.buffer),
                                     x2.t.buffer);

        // If there were error returns in the hash routine, indicate a problem
        // with the hash algorithm selection
        if(x2.t.size == 0)
            return TPM_RCS_HASH + RC_Commit_signHandle;
        // The size of the remainder will be same as the size of p. DivideB() will
        // pad the results (leading zeros) if necessary to make the size the same
        P2.x.t.size = p.t.size;
        //  set p2.x = hash(s2) mod p
        if(DivideB(&x2.b, &p.b, NULL, &P2.x.b) != TPM_RC_SUCCESS)
            return TPM_RC_NO_RESULT;

        if(!CryptEccIsPointOnCurve(parms->curveID, pP2))
            return TPM_RCS_ECC_POINT + RC_Commit_s2;

        if(eccKey->attributes.publicOnly == SET)
            return TPM_RCS_KEY + RC_Commit_signHandle;
    }
    // If there is a P1, make sure that it is on the curve
    // NOTE: an "empty" point has two UINT16 values which are the size values
    // for each of the coordinates.
    if(in->P1.size > 4)
    {
        pP1 = &in->P1.point;
        if(!CryptEccIsPointOnCurve(parms->curveID, pP1))
            return TPM_RCS_ECC_POINT + RC_Commit_P1;
    }

    // Pass the parameters to CryptCommit.
    // The work is not done in-line because it does several point multiplies
    // with the same curve.  It saves work by not having to reload the curve
    // parameters multiple times.
    result = CryptEccCommitCompute(&out->K.point,
                                   &out->L.point,
                                   &out->E.point,
                                   parms->curveID,
                                   pP1,
                                   pP2,
                                   &eccKey->sensitive.sensitive.ecc,
                                   &r);
    if(result != TPM_RC_SUCCESS)
        return result;

    // The commit computation was successful so complete the commit by setting
    // the bit
    out->counter = CryptCommit();

    return TPM_RC_SUCCESS;
}

#endif // CC_Commit