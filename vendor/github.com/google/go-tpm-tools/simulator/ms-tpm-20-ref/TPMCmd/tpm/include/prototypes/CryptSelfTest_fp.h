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
 *  Date: Mar 28, 2019  Time: 08:25:19PM
 */

#ifndef    _CRYPT_SELF_TEST_FP_H_
#define    _CRYPT_SELF_TEST_FP_H_

//*** CryptSelfTest()
// This function is called to start/complete a full self-test.
// If 'fullTest' is NO, then only the untested algorithms will be run. If
// 'fullTest' is YES, then 'g_untestedDecryptionAlgorithms' is reinitialized and then
// all tests are run.
// This implementation of the reference design does not support processing outside
// the framework of a TPM command. As a consequence, this command does not
// complete until all tests are done. Since this can take a long time, the TPM
// will check after each test to see if the command is canceled. If so, then the
// TPM will returned TPM_RC_CANCELLED. To continue with the self-tests, call
// TPM2_SelfTest(fullTest == No) and the TPM will complete the testing.
//  Return Type: TPM_RC
//      TPM_RC_CANCELED        if the command is canceled
LIB_EXPORT
TPM_RC
CryptSelfTest(
    TPMI_YES_NO      fullTest       // IN: if full test is required
);

//*** CryptIncrementalSelfTest()
// This function is used to perform an incremental self-test. This implementation
// will perform the toTest values before returning. That is, it assumes that the
// TPM cannot perform background tasks between commands.
//
// This command may be canceled. If it is, then there is no return result.
// However, this command can be run again and the incremental progress will not
// be lost.
//  Return Type: TPM_RC
//      TPM_RC_CANCELED         processing of this command was canceled
//      TPM_RC_TESTING          if toTest list is not empty
//      TPM_RC_VALUE            an algorithm in the toTest list is not implemented
TPM_RC
CryptIncrementalSelfTest(
    TPML_ALG            *toTest,        // IN: list of algorithms to be tested
    TPML_ALG            *toDoList       // OUT: list of algorithms needing test
);

//*** CryptInitializeToTest()
// This function will initialize the data structures for testing all the
// algorithms. This should not be called unless CryptAlgsSetImplemented() has
// been called
void
CryptInitializeToTest(
    void
);

//*** CryptTestAlgorithm()
// Only point of contact with the actual self tests. If a self-test fails, there
// is no return and the TPM goes into failure mode.
// The call to TestAlgorithm uses an algorithm selector and a bit vector. When the
// test is run, the corresponding bit in 'toTest' and in 'g_toTest' is CLEAR. If
// 'toTest' is NULL, then only the bit in 'g_toTest' is CLEAR.
// There is a special case for the call to TestAlgorithm(). When 'alg' is
// ALG_ERROR, TestAlgorithm() will CLEAR any bit in 'toTest' for which it has
// no test. This allows the knowledge about which algorithms have test to be
// accessed through the interface that provides the test.
//  Return Type: TPM_RC
//      TPM_RC_CANCELED     test was canceled
LIB_EXPORT
TPM_RC
CryptTestAlgorithm(
    TPM_ALG_ID           alg,
    ALGORITHM_VECTOR    *toTest
);

#endif  // _CRYPT_SELF_TEST_FP_H_
