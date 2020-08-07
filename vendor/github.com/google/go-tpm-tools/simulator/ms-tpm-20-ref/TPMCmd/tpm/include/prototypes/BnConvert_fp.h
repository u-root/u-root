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

#ifndef    _BN_CONVERT_FP_H_
#define    _BN_CONVERT_FP_H_

//*** BnFromBytes()
// This function will convert a big-endian byte array to the internal number
// format. If bn is NULL, then the output is NULL. If bytes is null or the
// required size is 0, then the output is set to zero
LIB_EXPORT bigNum
BnFromBytes(
    bigNum           bn,
    const BYTE      *bytes,
    NUMBYTES         nBytes
);

//*** BnFrom2B()
// Convert an TPM2B to a BIG_NUM.
// If the input value does not exist, or the output does not exist, or the input
// will not fit into the output the function returns NULL
LIB_EXPORT bigNum
BnFrom2B(
    bigNum           bn,         // OUT:
    const TPM2B     *a2B         // IN: number to convert
);

//*** BnFromHex()
// Convert a hex string into a bigNum. This is primarily used in debugging.
LIB_EXPORT bigNum
BnFromHex(
    bigNum          bn,         // OUT:
    const char      *hex        // IN:
);

//*** BnToBytes()
// This function converts a BIG_NUM to a byte array. It converts the bigNum to a
// big-endian byte string and sets 'size' to the normalized value. If  'size' is an
// input 0, then the receiving buffer is guaranteed to be large enough for the result
// and the size will be set to the size required for bigNum (leading zeros
// suppressed).
//
// The conversion for a little-endian machine simply requires that all significant
// bytes of the bigNum be reversed. For a big-endian machine, rather than
// unpack each word individually, the bigNum is converted to little-endian words,
// copied, and then converted back to big-endian.
LIB_EXPORT BOOL
BnToBytes(
    bigConst             bn,
    BYTE                *buffer,
    NUMBYTES            *size           // This the number of bytes that are
                                        // available in the buffer. The result
                                        // should be this big.
);

//*** BnTo2B()
// Function to convert a BIG_NUM to TPM2B.
// The TPM2B size is set to the requested 'size' which may require padding.
// If 'size' is non-zero and less than required by the value in 'bn' then an error
// is returned. If 'size' is zero, then the TPM2B is assumed to be large enough
// for the data and a2b->size will be adjusted accordingly.
LIB_EXPORT BOOL
BnTo2B(
    bigConst         bn,                // IN:
    TPM2B           *a2B,               // OUT:
    NUMBYTES         size               // IN: the desired size
);
#if ALG_ECC

//*** BnPointFrom2B()
// Function to create a BIG_POINT structure from a 2B point.
// A point is going to be two ECC values in the same buffer. The values are going
// to be the size of the modulus.  They are in modular form.
LIB_EXPORT bn_point_t   *
BnPointFrom2B(
    bigPoint             ecP,         // OUT: the preallocated point structure
    TPMS_ECC_POINT      *p            // IN: the number to convert
);

//*** BnPointTo2B()
// This function converts a BIG_POINT into a TPMS_ECC_POINT. A TPMS_ECC_POINT
// contains two TPM2B_ECC_PARAMETER values. The maximum size of the parameters
// is dependent on the maximum EC key size used in an implementation.
// The presumption is that the TPMS_ECC_POINT is large enough to hold 2 TPM2B
// values, each as large as a MAX_ECC_PARAMETER_BYTES
LIB_EXPORT BOOL
BnPointTo2B(
    TPMS_ECC_POINT  *p,             // OUT: the converted 2B structure
    bigPoint         ecP,           // IN: the values to be converted
    bigCurve         E              // IN: curve descriptor for the point
);
#endif // ALG_ECC

#endif  // _BN_CONVERT_FP_H_
