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
 *  Date: Apr  7, 2019  Time: 06:58:58PM
 */

#ifndef    _MEMORY_FP_H_
#define    _MEMORY_FP_H_

//*** MemoryCopy()
// This is an alias for memmove. This is used in place of memcpy because
// some of the moves may overlap and rather than try to make sure that
// memmove is used when necessary, it is always used.
void
MemoryCopy(
    void        *dest,
    const void  *src,
    int          sSize
);

//*** MemoryEqual()
// This function indicates if two buffers have the same values in the indicated
// number of bytes.
//  Return Type: BOOL
//      TRUE(1)         all octets are the same
//      FALSE(0)        all octets are not the same
BOOL
MemoryEqual(
    const void      *buffer1,       // IN: compare buffer1
    const void      *buffer2,       // IN: compare buffer2
    unsigned int     size           // IN: size of bytes being compared
);

//*** MemoryCopy2B()
// This function copies a TPM2B. This can be used when the TPM2B types are
// the same or different.
//
// This function returns the number of octets in the data buffer of the TPM2B.
LIB_EXPORT INT16
MemoryCopy2B(
    TPM2B           *dest,          // OUT: receiving TPM2B
    const TPM2B     *source,        // IN: source TPM2B
    unsigned int     dSize          // IN: size of the receiving buffer
);

//*** MemoryConcat2B()
// This function will concatenate the buffer contents of a TPM2B to an
// the buffer contents of another TPM2B and adjust the size accordingly
//      ('a' := ('a' | 'b')).
void
MemoryConcat2B(
    TPM2B           *aInOut,        // IN/OUT: destination 2B
    TPM2B           *bIn,           // IN: second 2B
    unsigned int     aMaxSize       // IN: The size of aInOut.buffer (max values for
                                    //     aInOut.size)
);

//*** MemoryEqual2B()
// This function will compare two TPM2B structures. To be equal, they
// need to be the same size and the buffer contexts need to be the same
// in all octets.
//  Return Type: BOOL
//      TRUE(1)         size and buffer contents are the same
//      FALSE(0)        size or buffer contents are not the same
BOOL
MemoryEqual2B(
    const TPM2B     *aIn,           // IN: compare value
    const TPM2B     *bIn            // IN: compare value
);

//*** MemorySet()
// This function will set all the octets in the specified memory range to
// the specified octet value.
// Note: A previous version had an additional parameter (dSize) that was
// intended to make sure that the destination would not be overrun. The
// problem is that, in use, all that was happening was that the value of
// size was used for dSize so there was no benefit in the extra parameter.
void
MemorySet(
    void            *dest,
    int              value,
    size_t           size
);

//*** MemoryPad2B()
// Function to pad a TPM2B with zeros and adjust the size.
void
MemoryPad2B(
    TPM2B           *b,
    UINT16           newSize
);

//*** Uint16ToByteArray()
// Function to write an integer to a byte array
void
Uint16ToByteArray(
    UINT16              i,
    BYTE                *a
);

//*** Uint32ToByteArray()
// Function to write an integer to a byte array
void
Uint32ToByteArray(
    UINT32              i,
    BYTE                *a
);

//*** Uint64ToByteArray()
// Function to write an integer to a byte array
void
Uint64ToByteArray(
    UINT64               i,
    BYTE                *a
);

//*** ByteArrayToUint8()
// Function to write a UINT8 to a byte array. This is included for completeness
// and to allow certain macro expansions
UINT8
ByteArrayToUint8(
    BYTE                *a
);

//*** ByteArrayToUint16()
// Function to write an integer to a byte array
UINT16
ByteArrayToUint16(
    BYTE                *a
);

//*** ByteArrayToUint32()
// Function to write an integer to a byte array
UINT32
ByteArrayToUint32(
    BYTE                *a
);

//*** ByteArrayToUint64()
// Function to write an integer to a byte array
UINT64
ByteArrayToUint64(
    BYTE                *a
);

#endif  // _MEMORY_FP_H_
