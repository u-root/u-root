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
#include "Load_fp.h"

#if CC_Load  // Conditional expansion of this file

#include "Object_spt_fp.h"

/*(See part 3 specification)
// Load an ordinary or temporary object
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       'inPulblic' attributes are not allowed with selected
//                              parent
//      TPM_RC_BINDING          'inPrivate' and 'inPublic' are not
//                              cryptographically bound
//      TPM_RC_HASH             incorrect hash selection for signing key or
//                              the 'nameAlg' for 'inPubic is not valid
//      TPM_RC_INTEGRITY        HMAC on 'inPrivate' was not valid
//      TPM_RC_KDF              KDF selection not allowed
//      TPM_RC_KEY              the size of the object's 'unique' field is not
//                              consistent with the indicated size in the object's
//                              parameters
//      TPM_RC_OBJECT_MEMORY    no available object slot
//      TPM_RC_SCHEME           the signing scheme is not valid for the key
//      TPM_RC_SENSITIVE        the 'inPrivate' did not unmarshal correctly
//      TPM_RC_SIZE             'inPrivate' missing, or 'authPolicy' size for
//                              'inPublic' or is not valid
//      TPM_RC_SYMMETRIC        symmetric algorithm not provided when required
//      TPM_RC_TYPE             'parentHandle' is not a storage key, or the object
//                              to load is a storage key but its parameters do not
//                              match the parameters of the parent.
//      TPM_RC_VALUE            decryption failure
TPM_RC
TPM2_Load(
    Load_In         *in,            // IN: input parameter list
    Load_Out        *out            // OUT: output parameter list
    )
{
    TPM_RC                   result = TPM_RC_SUCCESS;
    TPMT_SENSITIVE           sensitive;
    OBJECT                  *parentObject;
    OBJECT                  *newObject;

// Input Validation
    // Don't get invested in loading if there is no place to put it.
    newObject = FindEmptyObjectSlot(&out->objectHandle);
    if(newObject == NULL)
        return TPM_RC_OBJECT_MEMORY;

    if(in->inPrivate.t.size == 0)
        return TPM_RCS_SIZE + RC_Load_inPrivate;

    parentObject = HandleToObject(in->parentHandle);
    pAssert(parentObject != NULL);
    // Is the object that is being used as the parent actually a parent.
    if(!ObjectIsParent(parentObject))
        return TPM_RCS_TYPE + RC_Load_parentHandle;

    // Compute the name of object. If there isn't one, it is because the nameAlg is
    // not valid.
    PublicMarshalAndComputeName(&in->inPublic.publicArea, &out->name);
    if(out->name.t.size == 0)
        return TPM_RCS_HASH + RC_Load_inPublic;

    // Retrieve sensitive data. 
    result = PrivateToSensitive(&in->inPrivate.b, &out->name.b, parentObject,
                                in->inPublic.publicArea.nameAlg,
                                &sensitive);
    if(result != TPM_RC_SUCCESS)
        return RcSafeAddToResult(result, RC_Load_inPrivate);

// Internal Data Update
    // Load and validate object
    result = ObjectLoad(newObject, parentObject, 
                        &in->inPublic.publicArea, &sensitive, 
                        RC_Load_inPublic, RC_Load_inPrivate,
                        &out->name);
    if(result == TPM_RC_SUCCESS)
    {
        // Set the common OBJECT attributes for a loaded object.
        ObjectSetLoadedAttributes(newObject, in->parentHandle);
    }
    return result;

}

#endif // CC_Load