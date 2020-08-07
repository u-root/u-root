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
//** Includes
#include "Tpm.h"

#define _OIDS_
#include "OIDs.h"

#include "TpmASN1.h"
#include "TpmASN1_fp.h"

//** Unmarshaling Functions

//*** ASN1UnmarshalContextInitialize()
// Function does standard initialization of a context.
//  Return Type: BOOL
//      TRUE(1)     success
//      FALSE(0)    failure
BOOL
ASN1UnmarshalContextInitialize(
    ASN1UnmarshalContext    *ctx,
    INT16                    size,
    BYTE                    *buffer
)
{
    VERIFY(buffer != NULL);
    VERIFY(size > 0);
    ctx->buffer = buffer;
    ctx->size = size;
    ctx->offset = 0;
    ctx->tag = 0xFF;
    return TRUE;
Error:
    return FALSE;
}

//***ASN1DecodeLength()
// This function extracts the length of an element from 'buffer' starting at 'offset'.
// Return Type: UINT16
//      >=0         the extracted length
//      <0          an error
INT16
ASN1DecodeLength(
    ASN1UnmarshalContext        *ctx
)
{
    BYTE                first;                  // Next octet in buffer
    INT16               value;
//
    VERIFY(ctx->offset < ctx->size);
    first = NEXT_OCTET(ctx);
    // If the number of octets of the entity is larger than 127, then the first octet
    // is the number of octets in the length specifier. 
    if(first >= 0x80)
    {
        // Make sure that this length field is contained with the structure being 
        // parsed
        CHECK_SIZE(ctx, (first & 0x7F));
        if(first == 0x82)
        {
            // Two octets of size
            // get the next value
            value = (INT16)NEXT_OCTET(ctx);
            // Make sure that the result will fit in an INT16
            VERIFY(value < 0x0080);
            // Shift up and add next octet
            value = (value << 8) + NEXT_OCTET(ctx);
        }
        else if(first == 0x81)
            value = NEXT_OCTET(ctx);
        // Sizes larger than will fit in a INT16 are an error 
        else
            goto Error;
    }
    else
        value = first;
    // Make sure that the size defined something within the current context
    CHECK_SIZE(ctx, value);
    return value;
Error:
    ctx->size = -1;             // Makes everything fail from now on.
    return -1;
}

//***ASN1NextTag()
// This function extracts the next type from 'buffer' starting at 'offset'. 
// It advances 'offset' as it parses the type and the length of the type. It returns
// the length of the type. On return, the 'length' octets starting at 'offset' are the
// octets of the type.
// Return Type: UINT
//     >=0          the number of octets in 'type'
//     <0           an error
INT16
ASN1NextTag(
    ASN1UnmarshalContext    *ctx 
)
{
    // A tag to get?
    VERIFY(ctx->offset < ctx->size);
    // Get it
    ctx->tag = NEXT_OCTET(ctx);
    // Make sure that it is not an extended tag
    VERIFY((ctx->tag & 0x1F) != 0x1F);
    // Get the length field and return that
    return ASN1DecodeLength(ctx);
    
Error:
    // Attempt to read beyond the end of the context or an illegal tag
    ctx->size = -1;         // Persistent failure
    ctx->tag = 0xFF;
    return -1;
}


//*** ASN1GetBitStringValue()
// Try to parse a bit string of up to 32 bits from a value that is expected to be
// a bit string.
// If there is a general parsing error, the context->size is set to -1.
//  Return Type: BOOL
//      TRUE(1)     success
//      FALSE(0)    failure
BOOL
ASN1GetBitStringValue(
    ASN1UnmarshalContext        *ctx,
    UINT32                      *val
)
{
    int                  shift;
    INT16                length;
    UINT32               value = 0;
//

    VERIFY((length = ASN1NextTag(ctx)) >= 1);
    VERIFY(ctx->tag == ASN1_BITSTRING);
    // Get the shift value for the bit field (how many bits to loop off of the end)
    shift = NEXT_OCTET(ctx);
    length--;
    // the shift count has to make sense
    VERIFY((shift < 8) && ((length > 0) || (shift == 0)));
    // if there are any bytes left
    for(; length > 0; length--)
    {
        if(length > 1)
        {
            // for all but the last octet, just shift and add the new octet
            VERIFY((value & 0xFF000000) == 0); // can't loose significant bits
            value = (value << 8) + NEXT_OCTET(ctx);
        }
        else
        {
            // for the last octet, just shift the accumulated value enough to 
            // accept the significant bits in the last octet and shift the last 
            // octet down
            VERIFY(((value & (0xFF000000 << (8 - shift)))) == 0);
            value = (value << (8 - shift)) + (NEXT_OCTET(ctx) >> shift);
        }
    }
    *val = value;
    return TRUE;
Error:
    ctx->size = -1;
    return FALSE;
}

//*******************************************************************
//** Marshaling Functions
//*******************************************************************

//*** Introduction
// Marshaling of an ASN.1 structure is accomplished from the bottom up. That is, 
// the things that will be at the end of the structure are added last. To manage the
// collecting of the relative sizes, start a context for the outermost container, if
// there is one, and then placing items in from the bottom up. If the bottom-most 
// item is also within a structure, create a nested context by calling 
// ASN1StartMarshalingContext().
//
// The context control structure contains a 'buffer' pointer, an 'offset', an 'end'
// and a stack. 'offset' is the offset from the start of the buffer of the last added
// byte. When 'offset' reaches 0, the buffer is full. 'offset' is a signed value so
// that, when it becomes negative, there is an overflow. Only two functions are 
// allowed to move bytes into the buffer: ASN1PushByte() and ASN1PushBytes(). These
// functions make sure that no data is written beyond the end of the buffer.
//
// When a new context is started, the current value of 'end' is pushed
// on the stack and 'end' is set to 'offset. As bytes are added, offset gets smaller.
// At any time, the count of bytes in the current context is simply 'end' - 'offset'.
//
// Since starting a new context involves setting 'end' = 'offset', the number of bytes
// in the context starts at 0. The nominal way of ending a context is to use
// 'end' - 'offset' to set the length value, and then a tag is added to the buffer. 
// Then the previous 'end' value is popped meaning that the context just ended 
// becomes a member of the now current context.
//
// The nominal strategy for building a completed ASN.1 structure is to push everything
// into the buffer and then move everything to the start of the buffer. The move is 
// simple as the size of the move is the initial 'end' value minus the final 'offset'
// value. The destination is 'buffer' and the source is 'buffer' + 'offset'. As Skippy
// would say "Easy peasy, Joe."
//
// It is not necessary to provide a buffer into which the data is placed. If no buffer
// is provided, then the marshaling process will return values needed for marshaling.
// On strategy for filling the buffer would be to execute the process for building
// the structure without using a buffer. This would return the overall size of the
// structure. Then that amount of data could be allocated for the buffer and the fill
// process executed again with the data going into the buffer. At the end, the data
// would be in its final resting place. 

//*** ASN1InitialializeMarshalContext()
// This creates a structure for handling marshaling of an ASN.1 formatted data 
// structure.
void
ASN1InitialializeMarshalContext(
    ASN1MarshalContext      *ctx,
    INT16                    length,
    BYTE                    *buffer
)
{
    ctx->buffer = buffer;
    if(buffer)
        ctx->offset = length;
    else
        ctx->offset = INT16_MAX;
    ctx->end = ctx->offset;
    ctx->depth = -1;
}

//*** ASN1StartMarshalContext()
// This starts a new constructed element. It is constructed on 'top' of the value
// that was previously placed in the structure.
void
ASN1StartMarshalContext(
    ASN1MarshalContext      *ctx
)
{
    pAssert((ctx->depth + 1) < MAX_DEPTH);
    ctx->depth++;
    ctx->ends[ctx->depth] = ctx->end;
    ctx->end = ctx->offset;
}

//*** ASN1EndMarshalContext()
// This function restores the end pointer for an encapsulating structure.
//  Return Type: INT16
//      > 0             the size of the encapsulated structure that was just ended
//      <= 0            an error
INT16
ASN1EndMarshalContext(
    ASN1MarshalContext      *ctx
)
{
    INT16                   length;
    pAssert(ctx->depth >= 0);
    length = ctx->end - ctx->offset;
    ctx->end = ctx->ends[ctx->depth--];
    if((ctx->depth == -1) && (ctx->buffer))
    {
        MemoryCopy(ctx->buffer, ctx->buffer + ctx->offset, ctx->end - ctx->offset);
    }
    return length;
}


//***ASN1EndEncapsulation()
// This function puts a tag and length in the buffer. In this function, an embedded
// BIT_STRING is assumed to be a collection of octets. To indicate that all bits
// are used, a byte of zero is prepended. If a raw bit-string is needed, a new
// function like ASN1PushInteger() would be needed.
//  Return Type: INT16
//      > 0         number of octets in the encapsulation
//      == 0        failure
UINT16
ASN1EndEncapsulation(
    ASN1MarshalContext          *ctx,
    BYTE                         tag
)
{
    // only add a leading zero for an encapsulated BIT STRING
    if (tag == ASN1_BITSTRING)
        ASN1PushByte(ctx, 0);
    ASN1PushTagAndLength(ctx, tag, ctx->end - ctx->offset);
    return ASN1EndMarshalContext(ctx);
}

//*** ASN1PushByte()
BOOL
ASN1PushByte(
    ASN1MarshalContext          *ctx,
    BYTE                         b
)
{
    if(ctx->offset > 0)
    {
        ctx->offset -= 1;
        if(ctx->buffer)
            ctx->buffer[ctx->offset] = b;
        return TRUE;
    }
    ctx->offset = -1;
    return FALSE;
}

//*** ASN1PushBytes()
// Push some raw bytes onto the buffer. 'count' cannot be zero.
//  Return Type: IN16
//      > 0             count bytes
//      == 0            failure unless count was zero
INT16
ASN1PushBytes(
    ASN1MarshalContext          *ctx,
    INT16                        count,
    const BYTE                  *buffer
)
{
    // make sure that count is not negative which would mess up the math; and that 
    // if there is a count, there is a buffer
    VERIFY((count >= 0) && ((buffer != NULL) || (count == 0)));
    // back up the offset to determine where the new octets will get pushed
    ctx->offset -= count;
    // can't go negative
    VERIFY(ctx->offset >= 0);
    // if there are buffers, move the data, otherwise, assume that this is just a
    // test. 
    if(count && buffer && ctx->buffer)
        MemoryCopy(&ctx->buffer[ctx->offset], buffer, count);
    return count;
Error:
    ctx->offset = -1;
    return 0;
}

//*** ASN1PushNull()
//  Return Type: IN16
//      > 0             count bytes
//      == 0            failure unless count was zero
INT16
ASN1PushNull(
    ASN1MarshalContext      *ctx
)
{
    ASN1PushByte(ctx, 0);
    ASN1PushByte(ctx, ASN1_NULL);
    return (ctx->offset >= 0) ? 2 : 0;
}

//*** ASN1PushLength()
// Push a length value. This will only handle length values that fit in an INT16.
//  Return Type: UINT16
//      > 0         number of bytes added
//      == 0        failure
INT16
ASN1PushLength(
    ASN1MarshalContext          *ctx,
    INT16                        len
)
{
    UINT16                       start = ctx->offset;
    VERIFY(len >= 0);
    if(len <= 127)
        ASN1PushByte(ctx, (BYTE)len);
    else
    {
        ASN1PushByte(ctx, (BYTE)(len & 0xFF));
        len >>= 8;
        if(len == 0)
            ASN1PushByte(ctx, 0x81);
        else
        {
            ASN1PushByte(ctx, (BYTE)(len));
            ASN1PushByte(ctx, 0x82);
        }
    }
    goto Exit;
Error:
    ctx->offset = -1;
Exit:
    return (ctx->offset > 0) ? start - ctx->offset : 0;
}

//*** ASN1PushTagAndLength()
//  Return Type: INT16
//      > 0         number of bytes added
//      == 0        failure
INT16
ASN1PushTagAndLength(
    ASN1MarshalContext          *ctx,
    BYTE                         tag,
    INT16                        length
)
{
    INT16       bytes;
    bytes = ASN1PushLength(ctx, length);
    bytes += (INT16)ASN1PushByte(ctx, tag);
    return (ctx->offset < 0) ? 0 : bytes;
}


//*** ASN1PushTaggedOctetString()
// This function will push a random octet string. 
//  Return Type: INT16
//      > 0         number of bytes added
//      == 0        failure
INT16
ASN1PushTaggedOctetString(
    ASN1MarshalContext          *ctx,
    INT16                        size,
    const BYTE                  *string,
    BYTE                         tag
)
{
    ASN1PushBytes(ctx, size, string);
    // PushTagAndLenght just tells how many octets it added so the total size of this
    // element is the sum of those octets and input size.
    size += ASN1PushTagAndLength(ctx, tag, size);
    return size;
}

//*** ASN1PushUINT()
// This function pushes an native-endian integer value. This just changes a
// native-endian integer into a big-endian byte string and calls ASN1PushInteger().
// That function will remove leading zeros and make sure that the number is positive.
//  Return Type: IN16
//      > 0             count bytes
//      == 0            failure unless count was zero
INT16
ASN1PushUINT(
    ASN1MarshalContext      *ctx,
    UINT32                   integer
)
{
    BYTE                    marshaled[4];
    UINT32_TO_BYTE_ARRAY(integer, marshaled);
    return ASN1PushInteger(ctx, 4, marshaled);
}

//*** ASN1PushInteger
// Push a big-endian integer on the end of the buffer 
//  Return Type: UINT16
//      > 0         the number of bytes marshaled for the integer
//      == 0        failure
INT16 
ASN1PushInteger(
    ASN1MarshalContext  *ctx,           // IN/OUT: buffer context
    INT16                iLen,          // IN: octets of the integer
    BYTE                *integer        // IN: big-endian integer
)
{
    // no leading 0's
    while((*integer == 0) && (--iLen > 0))
        integer++;
    // Move the bytes to the buffer
    ASN1PushBytes(ctx, iLen, integer);
    // if needed, add a leading byte of 0 to make the number positive
    if(*integer & 0x80)
        iLen += (INT16)ASN1PushByte(ctx, 0);
    // PushTagAndLenght just tells how many octets it added so the total size of this
    // element is the sum of those octets and the adjusted input size.
    iLen +=  ASN1PushTagAndLength(ctx, ASN1_INTEGER, iLen);
    return iLen;
}

//*** ASN1PushOID()
// This function is used to add an OID. An OID is 0x06 followed by a byte of size 
// followed by size bytes. This is used to avoid having to do anything special in the
// definition of an OID.
//  Return Type: UINT16
//      > 0         the number of bytes marshaled for the integer
//      == 0        failure
INT16
ASN1PushOID(
    ASN1MarshalContext          *ctx,
    const BYTE                  *OID
)
{
    if((*OID == ASN1_OBJECT_IDENTIFIER) && ((OID[1] & 0x80) == 0))
    {
        return ASN1PushBytes(ctx, OID[1] + 2, OID);
    }
    ctx->offset = -1;
    return 0;
}


