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
#ifndef     _CAPABILITIES_H
#define     _CAPABILITIES_H

#define    MAX_CAP_DATA         (MAX_CAP_BUFFER - sizeof(TPM_CAP)-sizeof(UINT32))
#define    MAX_CAP_ALGS         (MAX_CAP_DATA / sizeof(TPMS_ALG_PROPERTY))
#define    MAX_CAP_HANDLES      (MAX_CAP_DATA / sizeof(TPM_HANDLE))
#define    MAX_CAP_CC           (MAX_CAP_DATA / sizeof(TPM_CC))
#define    MAX_TPM_PROPERTIES   (MAX_CAP_DATA / sizeof(TPMS_TAGGED_PROPERTY))
#define    MAX_PCR_PROPERTIES   (MAX_CAP_DATA / sizeof(TPMS_TAGGED_PCR_SELECT))
#define    MAX_ECC_CURVES       (MAX_CAP_DATA / sizeof(TPM_ECC_CURVE))
#define    MAX_TAGGED_POLICIES  (MAX_CAP_DATA / sizeof(TPMS_TAGGED_POLICY))

#define    MAX_AC_CAPABILITIES  (MAX_CAP_DATA / sizeof(TPMS_AC_OUTPUT))

#endif
