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
#include "NV_DefineSpace_fp.h"

#if CC_NV_DefineSpace  // Conditional expansion of this file

/*(See part 3 specification)
// Define a NV index space
*/
//  Return Type: TPM_RC
//      TPM_RC_HIERARCHY            for authorizations using TPM_RH_PLATFORM
//                                  phEnable_NV is clear preventing access to NV
//                                  data in the platform hierarchy.
//      TPM_RC_ATTRIBUTES           attributes of the index are not consistent
//      TPM_RC_NV_DEFINED           index already exists
//      TPM_RC_NV_SPACE             insufficient space for the index
//      TPM_RC_SIZE                 'auth->size' or 'publicInfo->authPolicy.size' is
//                                  larger than the digest size of
//                                  'publicInfo->nameAlg'; or 'publicInfo->dataSize'
//                                  is not consistent with 'publicInfo->attributes'
//                                  (this includes the case when the index is
//                                   larger than a MAX_NV_BUFFER_SIZE but the
//                                   TPMA_NV_WRITEALL attribute is SET)
TPM_RC
TPM2_NV_DefineSpace(
    NV_DefineSpace_In   *in             // IN: input parameter list
    )
{
    TPMA_NV         attributes = in->publicInfo.nvPublic.attributes;
    UINT16          nameSize;

    nameSize = CryptHashGetDigestSize(in->publicInfo.nvPublic.nameAlg);

// Input Validation

    // Checks not specific to type

    // If the UndefineSpaceSpecial command is not implemented, then can't have
    // an index that can only be deleted with policy
#if CC_NV_UndefineSpaceSpecial == NO
    if(IS_ATTRIBUTE(attributes, TPMA_NV, POLICY_DELETE))
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;
#endif

    // check that the authPolicy consistent with hash algorithm

    if(in->publicInfo.nvPublic.authPolicy.t.size != 0
       && in->publicInfo.nvPublic.authPolicy.t.size != nameSize)
        return TPM_RCS_SIZE + RC_NV_DefineSpace_publicInfo;

   // make sure that the authValue is not too large
    if(MemoryRemoveTrailingZeros(&in->auth)
       > CryptHashGetDigestSize(in->publicInfo.nvPublic.nameAlg))
        return TPM_RCS_SIZE + RC_NV_DefineSpace_auth;

    // If an index is being created by the owner and shEnable is
    // clear, then we would not reach this point because ownerAuth
    // can't be given when shEnable is CLEAR. However, if phEnable
    // is SET but phEnableNV is CLEAR, we have to check here
    if(in->authHandle == TPM_RH_PLATFORM && gc.phEnableNV == CLEAR)
        return TPM_RCS_HIERARCHY + RC_NV_DefineSpace_authHandle;

    // Attribute checks
    // Eliminate the unsupported types
    switch(GET_TPM_NT(attributes))
    {
#if CC_NV_Increment == YES
        case TPM_NT_COUNTER:
#endif
#if CC_NV_SetBits == YES
        case TPM_NT_BITS:
#endif
#if CC_NV_Extend == YES
        case TPM_NT_EXTEND:
#endif
#if CC_PolicySecret == YES && defined TPM_NT_PIN_PASS
        case TPM_NT_PIN_PASS:
        case TPM_NT_PIN_FAIL:
#endif
        case TPM_NT_ORDINARY:
            break;
        default:
            return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;
            break;
    }
    // Check that the sizes are OK based on the type
    switch(GET_TPM_NT(attributes))
    {
        case TPM_NT_ORDINARY:
            // Can't exceed the allowed size for the implementation
            if(in->publicInfo.nvPublic.dataSize > MAX_NV_INDEX_SIZE)
                return TPM_RCS_SIZE + RC_NV_DefineSpace_publicInfo;
            break;
        case TPM_NT_EXTEND:
            if(in->publicInfo.nvPublic.dataSize != nameSize)
                return TPM_RCS_SIZE + RC_NV_DefineSpace_publicInfo;
            break;
        default:
            // Everything else needs a size of 8
            if(in->publicInfo.nvPublic.dataSize != 8)
                return TPM_RCS_SIZE + RC_NV_DefineSpace_publicInfo;
            break;
    }
    // Handle other specifics
    switch(GET_TPM_NT(attributes))
    {
        case TPM_NT_COUNTER:
            // Counter can't have TPMA_NV_CLEAR_STCLEAR SET (don't clear counters)
            if(IS_ATTRIBUTE(attributes, TPMA_NV, CLEAR_STCLEAR))
                return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;
            break;
#ifdef TPM_NT_PIN_FAIL
        case TPM_NT_PIN_FAIL:
            // NV_NO_DA must be SET and AUTHWRITE must be CLEAR
            // NOTE: As with a PIN_PASS index, the authValue of the index is not
            // available until the index is written. If AUTHWRITE is the only way to
            // write then index, it could never be written. Rather than go through
            // all of the other possible ways to write the Index, it is simply
            // prohibited to write the index with the authValue. Other checks
            // below will insure that there seems to be a way to write the index
            // (i.e., with platform authorization , owner authorization,
            // or with policyAuth.)
            // It is not allowed to create a PIN Index that can't be modified.
            if(!IS_ATTRIBUTE(attributes, TPMA_NV, NO_DA))
                return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;
#endif
#ifdef TPM_NT_PIN_PASS
        case TPM_NT_PIN_PASS:
            // AUTHWRITE must be CLEAR (see note above to TPM_NT_PIN_FAIL)
            if(IS_ATTRIBUTE(attributes, TPMA_NV, AUTHWRITE)
               || IS_ATTRIBUTE(attributes, TPMA_NV, GLOBALLOCK)
               || IS_ATTRIBUTE(attributes, TPMA_NV, WRITEDEFINE))
                return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;
#endif  // this comes before break because PIN_FAIL falls through
            break;
        default:
            break;
    }

    // Locks may not be SET and written cannot be SET
    if(IS_ATTRIBUTE(attributes, TPMA_NV, WRITTEN)
       || IS_ATTRIBUTE(attributes, TPMA_NV, WRITELOCKED)
       || IS_ATTRIBUTE(attributes, TPMA_NV, READLOCKED))
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;

    // There must be a way to read the index.
    if(!IS_ATTRIBUTE(attributes, TPMA_NV, OWNERREAD)
       && !IS_ATTRIBUTE(attributes, TPMA_NV, PPREAD)
       && !IS_ATTRIBUTE(attributes, TPMA_NV, AUTHREAD)
       && !IS_ATTRIBUTE(attributes, TPMA_NV, POLICYREAD))
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;

    // There must be a way to write the index
    if(!IS_ATTRIBUTE(attributes, TPMA_NV, OWNERWRITE)
       && !IS_ATTRIBUTE(attributes, TPMA_NV, PPWRITE)
       && !IS_ATTRIBUTE(attributes, TPMA_NV, AUTHWRITE)
       && !IS_ATTRIBUTE(attributes, TPMA_NV, POLICYWRITE))
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;

    // An index with TPMA_NV_CLEAR_STCLEAR can't have TPMA_NV_WRITEDEFINE SET
    if(IS_ATTRIBUTE(attributes, TPMA_NV, CLEAR_STCLEAR)
       &&  IS_ATTRIBUTE(attributes, TPMA_NV, WRITEDEFINE))
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;

    // Make sure that the creator of the index can delete the index
    if((IS_ATTRIBUTE(attributes, TPMA_NV, PLATFORMCREATE)
        && in->authHandle == TPM_RH_OWNER)
       || (!IS_ATTRIBUTE(attributes, TPMA_NV, PLATFORMCREATE)
           && in->authHandle == TPM_RH_PLATFORM))
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_authHandle;

    // If TPMA_NV_POLICY_DELETE is SET, then the index must be defined by
    // the platform
    if(IS_ATTRIBUTE(attributes, TPMA_NV, POLICY_DELETE)
       &&  TPM_RH_PLATFORM != in->authHandle)
        return TPM_RCS_ATTRIBUTES + RC_NV_DefineSpace_publicInfo;

    // Make sure that the TPMA_NV_WRITEALL is not set if the index size is larger
    // than the allowed NV buffer size.
    if(in->publicInfo.nvPublic.dataSize > MAX_NV_BUFFER_SIZE
       &&  IS_ATTRIBUTE(attributes, TPMA_NV, WRITEALL))
        return TPM_RCS_SIZE + RC_NV_DefineSpace_publicInfo;

    // And finally, see if the index is already defined.
    if(NvIndexIsDefined(in->publicInfo.nvPublic.nvIndex))
        return TPM_RC_NV_DEFINED;

// Internal Data Update
    // define the space.  A TPM_RC_NV_SPACE error may be returned at this point
    return NvDefineIndex(&in->publicInfo.nvPublic, &in->auth);
}

#endif // CC_NV_DefineSpace