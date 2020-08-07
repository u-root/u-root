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

#ifndef    _TPM_TO_OSSL_SUPPORT_FP_H_
#define    _TPM_TO_OSSL_SUPPORT_FP_H_

#ifdef MATH_LIB_OSSL

//*** SupportLibInit()
// This does any initialization required by the support library.
LIB_EXPORT int
SupportLibInit(
    void
);

//*** OsslContextEnter()
// This function is used to initialize an OpenSSL context at the start of a function
// that will call to an OpenSSL math function.
BN_CTX *
OsslContextEnter(
    void
);

//*** OsslContextLeave()
// This is the companion function to OsslContextEnter().
void
OsslContextLeave(
    BN_CTX          *CTX
);

//*** OsslPushContext()
// This function is used to create a frame in a context. All values allocated within
// this context after the frame is started will be automatically freed when the
// context (OsslPopContext()
BN_CTX *
OsslPushContext(
    BN_CTX          *CTX
);

//*** OsslPopContext()
// This is the companion function to OsslPushContext().
void
OsslPopContext(
    BN_CTX          *CTX
);
#endif // MATH_LIB_OSSL

#endif  // _TPM_TO_OSSL_SUPPORT_FP_H_
