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

//** Includes and Data Definitions

// This definition allows this module to "see" the values that are private
// to this module but kept in Global.c for ease of state migration.
#define IO_BUFFER_C
#include "Tpm.h"
#include "IoBuffers_fp.h"

//** Buffers and Functions

// These buffers are set aside to hold command and response values. In this
// implementation, it is not guaranteed that the code will stop accessing
// the s_actionInputBuffer before starting to put values in the
// s_actionOutputBuffer so different buffers are required.
//

//*** MemoryIoBufferAllocationReset()
// This function is used to reset the allocation of buffers.
void
MemoryIoBufferAllocationReset(
    void
)
{
    s_actionIoAllocation = 0;
}

//*** MemoryIoBufferZero()
// Function zeros the action I/O buffer at the end of a command. Calling this is
// not mandatory for proper functionality.
void
MemoryIoBufferZero(
    void
)
{
    memset(s_actionIoBuffer, 0, s_actionIoAllocation);
}

//*** MemoryGetInBuffer()
// This function returns the address of the buffer into which the
// command parameters will be unmarshaled in preparation for calling
// the command actions.
BYTE *
MemoryGetInBuffer(
    UINT32           size           // Size, in bytes, required for the input
                                    // unmarshaling
    )
{
    pAssert(size <= sizeof(s_actionIoBuffer));
    // In this implementation, a static buffer is set aside for the command action
    // buffers. The buffer is shared between input and output. This is because
    // there is no need to allocate for the worst case input and worst case output
    // at the same time.
    // Round size up
    #define UoM  (sizeof(s_actionIoBuffer[0]))
    size = (size + (UoM - 1)) & (UINT32_MAX - (UoM - 1));
    memset(s_actionIoBuffer, 0, size);
    s_actionIoAllocation = size;
    return (BYTE *)&s_actionIoBuffer[0];
}

//*** MemoryGetOutBuffer()
// This function returns the address of the buffer into which the command
// action code places its output values.
BYTE *
MemoryGetOutBuffer(
    UINT32           size           // required size of the buffer
    )
{
    BYTE        *retVal = (BYTE *)(&s_actionIoBuffer[s_actionIoAllocation / UoM]);
    pAssert((size + s_actionIoAllocation) < (sizeof(s_actionIoBuffer)));
    // In this implementation, a static buffer is set aside for the command action
    // output buffer.
    memset(retVal, 0, size);
    s_actionIoAllocation += size;
    return retVal;
}

//*** IsLabelProperlyFormatted()
// This function checks that a label is a null-terminated string.
// NOTE: this function is here because there was no better place for it.
//  Return Type: BOOL
//      TRUE(1)         string is null terminated
//      FALSE(0)        string is not null terminated
BOOL
IsLabelProperlyFormatted(
    TPM2B           *x
    )
{
    return (((x)->size == 0) || ((x)->buffer[(x)->size - 1] == 0));
}
