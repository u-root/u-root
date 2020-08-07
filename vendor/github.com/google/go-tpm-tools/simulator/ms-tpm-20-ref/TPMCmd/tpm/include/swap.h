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
#ifndef _SWAP_H
#define _SWAP_H

#if LITTLE_ENDIAN_TPM
#define TO_BIG_ENDIAN_UINT16(i)     REVERSE_ENDIAN_16(i)
#define FROM_BIG_ENDIAN_UINT16(i)   REVERSE_ENDIAN_16(i)
#define TO_BIG_ENDIAN_UINT32(i)     REVERSE_ENDIAN_32(i)
#define FROM_BIG_ENDIAN_UINT32(i)   REVERSE_ENDIAN_32(i)
#define TO_BIG_ENDIAN_UINT64(i)     REVERSE_ENDIAN_64(i)
#define FROM_BIG_ENDIAN_UINT64(i)   REVERSE_ENDIAN_64(i)
#else
#define TO_BIG_ENDIAN_UINT16(i)     (i)
#define FROM_BIG_ENDIAN_UINT16(i)   (i)
#define TO_BIG_ENDIAN_UINT32(i)     (i)
#define FROM_BIG_ENDIAN_UINT32(i)   (i)
#define TO_BIG_ENDIAN_UINT64(i)     (i)
#define FROM_BIG_ENDIAN_UINT64(i)   (i)
#endif

#if   AUTO_ALIGN == NO 

// The aggregation macros for machines that do not allow unaligned access or for
// little-endian machines.

// Aggregate bytes into an UINT

#define BYTE_ARRAY_TO_UINT8(b)  (uint8_t)((b)[0])
#define BYTE_ARRAY_TO_UINT16(b) ByteArrayToUint16((BYTE *)(b))
#define BYTE_ARRAY_TO_UINT32(b) ByteArrayToUint32((BYTE *)(b))
#define BYTE_ARRAY_TO_UINT64(b) ByteArrayToUint64((BYTE *)(b))
#define UINT8_TO_BYTE_ARRAY(i, b) ((b)[0] = (uint8_t)(i))
#define UINT16_TO_BYTE_ARRAY(i, b)  Uint16ToByteArray((i), (BYTE *)(b))
#define UINT32_TO_BYTE_ARRAY(i, b)  Uint32ToByteArray((i), (BYTE *)(b))
#define UINT64_TO_BYTE_ARRAY(i, b)  Uint64ToByteArray((i), (BYTE *)(b))


#else // AUTO_ALIGN

#if BIG_ENDIAN_TPM
// the big-endian macros for machines that allow unaligned memory access
// Aggregate a byte array into a UINT
#define BYTE_ARRAY_TO_UINT8(b)        *((uint8_t  *)(b))
#define BYTE_ARRAY_TO_UINT16(b)       *((uint16_t *)(b))
#define BYTE_ARRAY_TO_UINT32(b)       *((uint32_t *)(b))
#define BYTE_ARRAY_TO_UINT64(b)       *((uint64_t *)(b))

// Disaggregate a UINT into a byte array

#define UINT8_TO_BYTE_ARRAY(i, b)   {*((uint8_t  *)(b)) = (i);}
#define UINT16_TO_BYTE_ARRAY(i, b)  {*((uint16_t *)(b)) = (i);}
#define UINT32_TO_BYTE_ARRAY(i, b)  {*((uint32_t *)(b)) = (i);}
#define UINT64_TO_BYTE_ARRAY(i, b)  {*((uint64_t *)(b)) = (i);}
#else
// the little endian macros for machines that allow unaligned memory access
// the big-endian macros for machines that allow unaligned memory access
// Aggregate a byte array into a UINT
#define BYTE_ARRAY_TO_UINT8(b)        *((uint8_t  *)(b))
#define BYTE_ARRAY_TO_UINT16(b)       REVERSE_ENDIAN_16(*((uint16_t *)(b)))
#define BYTE_ARRAY_TO_UINT32(b)       REVERSE_ENDIAN_32(*((uint32_t *)(b)))
#define BYTE_ARRAY_TO_UINT64(b)       REVERSE_ENDIAN_64(*((uint64_t *)(b)))

// Disaggregate a UINT into a byte array

#define UINT8_TO_BYTE_ARRAY(i, b)   {*((uint8_t  *)(b)) = (i);}
#define UINT16_TO_BYTE_ARRAY(i, b)  {*((uint16_t *)(b)) = REVERSE_ENDIAN_16(i);}
#define UINT32_TO_BYTE_ARRAY(i, b)  {*((uint32_t *)(b)) = REVERSE_ENDIAN_32(i);}
#define UINT64_TO_BYTE_ARRAY(i, b)  {*((uint64_t *)(b)) = REVERSE_ENDIAN_64(i);}
#endif   // BIG_ENDIAN_TPM

#endif  // AUTO_ALIGN == NO

#endif  // _SWAP_H
