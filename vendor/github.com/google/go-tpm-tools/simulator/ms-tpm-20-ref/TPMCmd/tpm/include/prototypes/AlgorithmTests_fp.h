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
/*(Auto-generated)
 *  Created by TpmPrototypes; Version 3.0 July 18, 2017
 *  Date: Mar 28, 2019  Time: 08:25:18PM
 */

#ifndef    _ALGORITHM_TESTS_FP_H_
#define    _ALGORITHM_TESTS_FP_H_

#if SELF_TEST

//*** TestAlgorithm()
// Dispatches to the correct test function for the algorithm or gets a list of
// testable algorithms.
//
// If 'toTest' is not NULL, then the test decisions are based on the algorithm
// selections in 'toTest'. Otherwise, 'g_toTest' is used. When bits are clear in
// 'g_toTest' they will also be cleared 'toTest'.
//
// If there doesn't happen to be a test for the algorithm, its associated bit is
// quietly cleared.
//
// If 'alg' is zero (TPM_ALG_ERROR), then the toTest vector is cleared of any bits
// for which there is no test (i.e. no tests are actually run but the vector is
// cleared).
//
// Note: 'toTest' will only ever have bits set for implemented algorithms but 'alg'
// can be anything.
//  Return Type: TPM_RC
//      TPM_RC_CANCELED     test was canceled
LIB_EXPORT
TPM_RC
TestAlgorithm(
    TPM_ALG_ID               alg,
    ALGORITHM_VECTOR        *toTest
);
#endif // SELF_TESTS

#endif  // _ALGORITHM_TESTS_FP_H_
