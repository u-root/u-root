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
#ifndef     _MSBN_DRBG_TEST_VECTORS_H
#define     _MSBN_DRBG_TEST_VECTORS_H

//#if DRBG_ALGORITHM == TPM_ALG_AES && DRBG_KEY_BITS == 256
#if DRBG_KEY_SIZE_BITS == 256

/*(NIST test vector)
[AES-256 no df]
[PredictionResistance = False]
[EntropyInputLen = 384]
[NonceLen = 128]
[PersonalizationStringLen = 0]
[AdditionalInputLen = 0]

COUNT = 0
EntropyInput = 0d15aa80 b16c3a10 906cfedb 795dae0b 5b81041c 5c5bfacb
               373d4440 d9120f7e 3d6cf909 86cf52d8 5d3e947d 8c061f91
Nonce = 06caef5f b538e08e 1f3b0452 03f8f4b2
PersonalizationString = 
AdditionalInput = 
    INTERMEDIATE Key = be5df629 34cc1230 166a6773 345bbd6b 
                       4c8869cf 8aec1c3b 1aa98bca 37cacf61
    INTERMEDIATE V = 3182dd1e 7638ec70 014e93bd 813e524c
    INTERMEDIATE ReturnedBits = 28e0ebb8 21016650 8c8f65f2 207bd0a3
EntropyInputReseed = 6ee793a3 3955d72a d12fd80a 8a3fcf95 ed3b4dac 5795fe25 
                     cf869f7c 27573bbc 56f1acae 13a65042 b340093c 464a7a22
AdditionalInputReseed = 
AdditionalInput = 
ReturnedBits = 946f5182 d54510b9 461248f5 71ca06c9
*/


// Entropy is the size of the state. The state is the size of the key
// plus the IV. The IV is a block. If Key = 256 and Block = 128 then State = 384
#   define DRBG_TEST_INITIATE_ENTROPY                   \
        0x0d, 0x15, 0xaa, 0x80, 0xb1, 0x6c, 0x3a, 0x10, \
        0x90, 0x6c, 0xfe, 0xdb, 0x79, 0x5d, 0xae, 0x0b, \
        0x5b, 0x81, 0x04, 0x1c, 0x5c, 0x5b, 0xfa, 0xcb, \
        0x37, 0x3d, 0x44, 0x40, 0xd9, 0x12, 0x0f, 0x7e, \
        0x3d, 0x6c, 0xf9, 0x09, 0x86, 0xcf, 0x52, 0xd8, \
        0x5d, 0x3e, 0x94, 0x7d, 0x8c, 0x06, 0x1f, 0x91

#   define DRBG_TEST_RESEED_ENTROPY                     \
        0x6e, 0xe7, 0x93, 0xa3, 0x39, 0x55, 0xd7, 0x2a, \
        0xd1, 0x2f, 0xd8, 0x0a, 0x8a, 0x3f, 0xcf, 0x95, \
        0xed, 0x3b, 0x4d, 0xac, 0x57, 0x95, 0xfe, 0x25, \
        0xcf, 0x86, 0x9f, 0x7c, 0x27, 0x57, 0x3b, 0xbc, \
        0x56, 0xf1, 0xac, 0xae, 0x13, 0xa6, 0x50, 0x42, \
        0xb3, 0x40, 0x09, 0x3c, 0x46, 0x4a, 0x7a, 0x22

#   define DRBG_TEST_GENERATED_INTERM                   \
        0x28, 0xe0, 0xeb, 0xb8, 0x21, 0x01, 0x66, 0x50, \
        0x8c, 0x8f, 0x65, 0xf2, 0x20, 0x7b, 0xd0, 0xa3


#   define DRBG_TEST_GENERATED                          \
        0x94, 0x6f, 0x51, 0x82, 0xd5, 0x45, 0x10, 0xb9, \
        0x46, 0x12, 0x48, 0xf5, 0x71, 0xca, 0x06, 0xc9
#elif DRBG_KEY_SIZE_BITS == 128
/*(NIST test vector)
[AES-128 no df]
[PredictionResistance = False]
[EntropyInputLen = 256]
[NonceLen = 64]
[PersonalizationStringLen = 0]
[AdditionalInputLen = 0]

COUNT = 0
EntropyInput = 8fc11bdb5aabb7e093b61428e0907303cb459f3b600dad870955f22da80a44f8
Nonce = be1f73885ddd15aa
PersonalizationString = 
AdditionalInput = 
    INTERMEDIATE Key = b134ecc836df6dbd624900af118dd7e6
    INTERMEDIATE V = 01bb09e86dabd75c9f26dbf6f9531368
    INTERMEDIATE ReturnedBits = dc3cf6bf5bd341135f2c6811a1071c87
EntropyInputReseed = 
                 0cd53cd5eccd5a10d7ea266111259b05574fc6ddd8bed8bd72378cf82f1dba2a
AdditionalInputReseed = 
AdditionalInput = 
ReturnedBits = b61850decfd7106d44769a8e6e8c1ad4
*/

#   define DRBG_TEST_INITIATE_ENTROPY                   \
        0x8f, 0xc1, 0x1b, 0xdb, 0x5a, 0xab, 0xb7, 0xe0, \
        0x93, 0xb6, 0x14, 0x28, 0xe0, 0x90, 0x73, 0x03, \
        0xcb, 0x45, 0x9f, 0x3b, 0x60, 0x0d, 0xad, 0x87, \
        0x09, 0x55, 0xf2, 0x2d, 0xa8, 0x0a, 0x44, 0xf8
        
#   define DRBG_TEST_RESEED_ENTROPY                     \
        0x0c, 0xd5, 0x3c, 0xd5, 0xec, 0xcd, 0x5a, 0x10, \
        0xd7, 0xea, 0x26, 0x61, 0x11, 0x25, 0x9b, 0x05, \
        0x57, 0x4f, 0xc6, 0xdd, 0xd8, 0xbe, 0xd8, 0xbd, \
        0x72, 0x37, 0x8c, 0xf8, 0x2f, 0x1d, 0xba, 0x2a  
        
#define DRBG_TEST_GENERATED_INTERM                      \
        0xdc, 0x3c, 0xf6, 0xbf, 0x5b, 0xd3, 0x41, 0x13, \
        0x5f, 0x2c, 0x68, 0x11, 0xa1, 0x07, 0x1c, 0x87  

#   define DRBG_TEST_GENERATED                          \
        0xb6, 0x18, 0x50, 0xde, 0xcf, 0xd7, 0x10, 0x6d, \
        0x44, 0x76, 0x9a, 0x8e, 0x6e, 0x8c, 0x1a, 0xd4 

#endif


#endif      //     _MSBN_DRBG_TEST_VECTORS_H