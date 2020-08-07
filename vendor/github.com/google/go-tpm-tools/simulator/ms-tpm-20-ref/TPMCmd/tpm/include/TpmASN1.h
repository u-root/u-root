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
//** Introduction
// This file contains the macro and structure definitions for the X509 commands and
// functions.

#ifndef _TPMASN1_H_
#define _TPMASN1_H_

//** Includes

#include "Tpm.h"
#include "OIDs.h"

//** Defined Constants
//*** ASN.1 Universal Types (Class 00b)
#define ASN1_EOC                0x00
#define ASN1_BOOLEAN            0x01
#define ASN1_INTEGER            0x02
#define ASN1_BITSTRING          0x03
#define ASN1_OCTET_STRING       0x04
#define ASN1_NULL               0x05
#define ASN1_OBJECT_IDENTIFIER  0x06
#define ASN1_OBJECT_DESCRIPTOR  0x07
#define ASN1_EXTERNAL           0x08
#define ASN1_REAL               0x09
#define ASN1_ENUMERATED         0x0A
#define ASN1_EMBEDDED           0x0B
#define ASN1_UTF8String         0x0C
#define ASN1_RELATIVE_OID       0x0D
#define ASN1_SEQUENCE           0x10        // Primitive + Constructed + 0x10
#define ASN1_SET                0x11        // Primitive + Constructed + 0x11
#define ASN1_NumericString      0x12
#define ASN1_PrintableString    0x13
#define ASN1_T61String          0x14
#define ASN1_VideoString        0x15
#define ASN1_IA5String          0x16
#define ASN1_UTCTime            0x17
#define ASN1_GeneralizeTime     0x18
#define ASN1_VisibleString	    0x1A
#define ASN1_GeneralString	    0x1B
#define ASN1_UniversalString	0x1C
#define ASN1_CHARACTER STRING	0x1D
#define ASN1_BMPString	        0x1E
#define ASN1_CONSTRUCTED        0x20

#define ASN1_APPLICAIION_SPECIFIC   0xA0

#define ASN1_CONSTRUCTED_SEQUENCE (ASN1_SEQUENCE + ASN1_CONSTRUCTED)

#define MAX_DEPTH              10  // maximum push depth for marshaling context.

//** Macros

//*** Unmarshaling Macros
#ifndef VERIFY
#define VERIFY(_X_) {if(!(_X_)) goto Error; }
#endif
// Checks the validity of the size making sure that there is no wrap around
#define CHECK_SIZE(context, length)                                                 \
            VERIFY(   (((length) + (context)->offset) >= (context)->offset)         \
                   && (((length) + (context)->offset) <= (context)->size))
#define NEXT_OCTET(context) ((context)->buffer[(context)->offset++])
#define PEEK_NEXT(context) ((context)->buffer[(context)->offset])

//*** Marshaling Macros

// Marshaling works in reverse order. The offset is set to the top of the buffer and,
// as the buffer is filled, 'offset' counts down to zero. When the full thing is 
// encoded it can be moved to the top of the buffer. This happens when the last 
// context is closed. 

#define CHECK_SPACE(context, length)    VERIFY(context->offset > length)

//** Structures

typedef struct ASN1UnmarshalContext {
    BYTE            *buffer;    // pointer to the buffer
    INT16            size;      // size of the buffer (a negative number indicates 
                                // a parsing failure).
    INT16            offset;    // current offset into the buffer (a negative number 
                                // indicates a parsing failure). Not used
    BYTE             tag;       // The last unmarshaled tag
} ASN1UnmarshalContext;

typedef struct ASN1MarshalContext {
    BYTE            *buffer;    // pointer to the start of the buffer
    INT16            offset;    // place on the top where the last entry was added
                                // items are added from the bottom up.
    INT16            end;       // the end offset of the current value
    INT16            depth;     // how many pushed end values.
    INT16            ends[MAX_DEPTH];
} ASN1MarshalContext;

#endif // _TPMASN1_H_
