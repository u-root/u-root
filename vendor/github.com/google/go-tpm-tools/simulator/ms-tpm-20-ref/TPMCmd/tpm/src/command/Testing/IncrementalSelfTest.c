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
#include "IncrementalSelfTest_fp.h"

#if CC_IncrementalSelfTest  // Conditional expansion of this file

/*(See part 3 specification)
// perform a test of selected algorithms
*/
//  Return Type: TPM_RC
//      TPM_RC_CANCELED         the command was canceled (some tests may have
//                              completed)
//      TPM_RC_VALUE            an algorithm in the toTest list is not implemented
TPM_RC
TPM2_IncrementalSelfTest(
    IncrementalSelfTest_In      *in,            // IN: input parameter list
    IncrementalSelfTest_Out     *out            // OUT: output parameter list
    )
{
    TPM_RC                       result;
// Command Output

    // Call incremental self test function in crypt module. If this function
    // returns TPM_RC_VALUE, it means that an algorithm on the 'toTest' list is
    // not implemented.
    result = CryptIncrementalSelfTest(&in->toTest, &out->toDoList);
    if(result == TPM_RC_VALUE)
        return TPM_RCS_VALUE + RC_IncrementalSelfTest_toTest;
    return result;
}

#endif // CC_IncrementalSelfTest