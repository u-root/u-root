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
//** Includes and Typedefs
#include "Tpm.h"

#if TABLE_DRIVEN_DISPATCH

typedef TPM_RC(NoFlagFunction)(void *target, BYTE **buffer, INT32 *size);
typedef TPM_RC(FlagFunction)(void *target, BYTE **buffer, INT32 *size, BOOL flag);

typedef FlagFunction *UNMARSHAL_t;

typedef INT16(MarshalFunction)(void *source, BYTE **buffer, INT32 *size);
typedef MarshalFunction *MARSHAL_t;

typedef TPM_RC(COMMAND_NO_ARGS)(void);
typedef TPM_RC(COMMAND_IN_ARG)(void *in);
typedef TPM_RC(COMMAND_OUT_ARG)(void *out);
typedef TPM_RC(COMMAND_INOUT_ARG)(void *in, void *out);

typedef union COMMAND_t
{
    COMMAND_NO_ARGS         *noArgs;
    COMMAND_IN_ARG          *inArg;
    COMMAND_OUT_ARG         *outArg;
    COMMAND_INOUT_ARG       *inOutArg;
} COMMAND_t;

// This structure is used by ParseHandleBuffer() and CommandDispatcher(). The
// parameters in this structure are unique for each command. The parameters are:
// command      holds the address of the command processing function that is called
//              by Command Dispatcher.
// inSize       this is the size of the command-dependent input structure. The
//              input structure holds the unmarshaled handles and command
//              parameters. If the command takes no arguments (handles or
//              parameters) then inSize will have a value of 0.
// outSize      this is the size of the command-dependent output structure. The
//              output structure holds the results of the command in an unmarshaled
//              form. When command processing is completed, these values are
//              marshaled into the output buffer. It is always the case that the
//              unmarshaled version of an output structure is larger then the
//              marshaled version. This is because the marshaled version contains
//              the exact same number of significant bytes but with padding removed.
// typesOffsets    this parameter points to the list of data types that are to be
//              marshaled or unmarshaled. The list of types follows the 'offsets'
//              array. The offsets array is variable sized so the typesOffset filed
//              is necessary for the handle and command processing to be able to
//              find the types that are being handled. The 'offsets' array may be
//              empty. The types structure is described below.
// offsets      this is an array of offsets of each of the parameters in the
//              command or response. When processing the command parameters (not
//              handles) the list contains the offset of the next parameter. For
//              example, if the first command parameter has a size of 4 and there is
//              a second command parameter, then the offset would be 4, indicating
//              that the second parameter starts at 4. If the second parameter has
//              a size of 8, and there is a third parameter, then the second entry
//              in offsets is 12 (4 for the first parameter and 8 for the second).
//              An offset value of 0 in the list indicates the start of the response
//              parameter list. When CommandDispatcher hits this value, it will stop
//              unmarshaling the parameters and call 'command'. If a command has no
//              response parameters and only one command parameter, then offsets can
//              be an empty list.

typedef struct COMMAND_DESCRIPTOR_t
{
    COMMAND_t       command;        // Address of the command
    UINT16          inSize;         // Maximum size of the input structure
    UINT16          outSize;        // Maximum size of the output structure
    UINT16          typesOffset;    // address of the types field
    UINT16          offsets[1];
} COMMAND_DESCRIPTOR_t;

// The 'types' list is an encoded byte array. The byte value has two parts. The most
// significant bit is used when a parameter takes a flag and indicates if the flag
// should be SET or not. The remaining 7 bits are an index into an array of
// addresses of marshaling and unmarshaling functions.
// The array of functions is divided into 6 sections with a value assigned
// to denote the start of that section (and the end of the previous section). The
// defined offset values for each section are:
// 0                                unmarshaling for handles that do not take flags
// HANDLE_FIRST_FLAG_TYPE           unmarshaling for handles that take flags
// PARAMETER_FIRST_TYPE             unmarshaling for parameters that do not take flags
// PARAMETER_FIRST_FLAG_TYPE        unmarshaling for parameters that take flags
// PARAMETER_LAST_TYPE + 1          marshaling for handles
// RESPONSE_PARAMETER_FIRST_TYPE    marshaling for parameters
// RESPONSE_PARAMETER_LAST_TYPE     is the last value in the list of marshaling and
//                                  unmarshaling functions.
//
// The types list is constructed with a byte of 0xff at the end of the command
// parameters and with an 0xff at the end of the response parameters.

#if COMPRESSED_LISTS
#   define PAD_LIST 0
#else
#   define PAD_LIST 1
#endif
#define _COMMAND_TABLE_DISPATCH_
#include "CommandDispatchData.h"

#define TEST_COMMAND    TPM_CC_Startup

#define NEW_CC

#else

#include "Commands.h"

#endif

//** Marshal/Unmarshal Functions

//*** ParseHandleBuffer()
// This is the table-driven version of the handle buffer unmarshaling code
TPM_RC
ParseHandleBuffer(
    COMMAND                 *command
    )
{
    TPM_RC                   result;
#if TABLE_DRIVEN_DISPATCH
    COMMAND_DESCRIPTOR_t    *desc;
    BYTE                    *types;
    BYTE                     type;
    BYTE                     dType;

    // Make sure that nothing strange has happened
    pAssert(command->index
            < sizeof(s_CommandDataArray) / sizeof(COMMAND_DESCRIPTOR_t *));
    // Get the address of the descriptor for this command
    desc = s_CommandDataArray[command->index];

    pAssert(desc != NULL);
    // Get the associated list of unmarshaling data types.
    types = &((BYTE *)desc)[desc->typesOffset];

//    if(s_ccAttr[commandIndex].commandIndex == TEST_COMMAND)
//        commandIndex = commandIndex;
    // No handles yet
    command->handleNum = 0;

    // Get the first type value
    for(type = *types++;
        // check each byte to make sure that we have not hit the start
        // of the parameters
    (dType = (type & 0x7F)) < PARAMETER_FIRST_TYPE;
    // get the next type
        type = *types++)
    {
        // See if unmarshaling of this handle type requires a flag
        if(dType < HANDLE_FIRST_FLAG_TYPE)
        {
            // Look up the function to do the unmarshaling
            NoFlagFunction  *f = (NoFlagFunction *)UnmarshalArray[dType];
            // call it
            result = f(&(command->handles[command->handleNum]),
                       &command->parameterBuffer,
                       &command->parameterSize);
        }
        else
        {
            //  Look up the function
            FlagFunction    *f = UnmarshalArray[dType];

            // Call it setting the flag to the appropriate value
            result = f(&(command->handles[command->handleNum]),
                       &command->parameterBuffer,
                       &command->parameterSize, (type & 0x80) != 0);
        }
        // Got a handle
        // We do this first so that the match for the handle offset of the
        // response code works correctly.
        command->handleNum += 1;
        if(result != TPM_RC_SUCCESS)
            // if the unmarshaling failed, return the response code with the
            // handle indication set
            return result + TPM_RC_H + (command->handleNum * TPM_RC_1);
    }
#else
    BYTE            **handleBufferStart = &command->parameterBuffer;
    INT32           *bufferRemainingSize = &command->parameterSize;
    TPM_HANDLE      *handles = &command->handles[0];
    UINT32          *handleCount = &command->handleNum;
    *handleCount = 0;
    switch(command->code)
    {
#include "HandleProcess.h"
#undef handles
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
#endif
    return TPM_RC_SUCCESS;
}

//*** CommandDispatcher()
// Function to unmarshal the command parameters, call the selected action code, and
// marshal the response parameters.
TPM_RC
CommandDispatcher(
    COMMAND                 *command
    )
{
#if !TABLE_DRIVEN_DISPATCH
    TPM_RC       result;
    BYTE        **paramBuffer = &command->parameterBuffer;
    INT32       *paramBufferSize = &command->parameterSize;
    BYTE        **responseBuffer = &command->responseBuffer;
    INT32       *respParmSize = &command->parameterSize;
    INT32        rSize;
    TPM_HANDLE  *handles = &command->handles[0];
//
    command->handleNum = 0;                 // The command-specific code knows how
                                            // many handles there are. This is for
                                            // cataloging the number of response
                                            // handles
    MemoryIoBufferAllocationReset();        // Initialize so that allocation will
                                            // work properly
    switch(GetCommandCode(command->index))
    {
#include "CommandDispatcher.h"

        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
Exit:
    MemoryIoBufferZero();
    return result;
#else
    COMMAND_DESCRIPTOR_t    *desc;
    BYTE                    *types;
    BYTE                     type;
    UINT16                  *offsets;
    UINT16                   offset = 0;
    UINT32                   maxInSize;
    BYTE                    *commandIn;
    INT32                    maxOutSize;
    BYTE                    *commandOut;
    COMMAND_t                cmd;
    TPM_HANDLE              *handles;
    UINT32                   hasInParameters = 0;
    BOOL                     hasOutParameters = FALSE;
    UINT32                   pNum = 0;
    BYTE                     dType;     // dispatch type
    TPM_RC                   result;
//
    // Get the address of the descriptor for this command
    pAssert(command->index
            < sizeof(s_CommandDataArray) / sizeof(COMMAND_DESCRIPTOR_t *));
    desc = s_CommandDataArray[command->index];

    // Get the list of parameter types for this command
    pAssert(desc != NULL);
    types = &((BYTE *)desc)[desc->typesOffset];

    // Get a pointer to the list of parameter offsets
    offsets = &desc->offsets[0];
    // pointer to handles
    handles = command->handles;

    // Get the size required to hold all the unmarshaled parameters for this command
    maxInSize = desc->inSize;
    // and the size of the output parameter structure returned by this command
    maxOutSize = desc->outSize;

    MemoryIoBufferAllocationReset();
    // Get a buffer for the input parameters
    commandIn = MemoryGetInBuffer(maxInSize);
    // And the output parameters
    commandOut = (BYTE *)MemoryGetOutBuffer((UINT32)maxOutSize);

    // Get the address of the action code dispatch
    cmd = desc->command;

    // Copy any handles into the input buffer
    for(type = *types++; (type & 0x7F) < PARAMETER_FIRST_TYPE; type = *types++)
    {
        // 'offset' was initialized to zero so the first unmarshaling will always
        // be to the start of the data structure
        *(TPM_HANDLE *)&(commandIn[offset]) = *handles++;
        // This check is used so that we don't have to add an additional offset
        // value to the offsets list to correspond to the stop value in the
        // command parameter list.
        if(*types != 0xFF)
            offset = *offsets++;
//        maxInSize -= sizeof(TPM_HANDLE);
        hasInParameters++;
    }
    // Exit loop with type containing the last value read from types
    // maxInSize has the amount of space remaining in the command action input
    // buffer. Make sure that we don't have more data to unmarshal than is going to
    // fit.

    // type contains the last value read from types so it is not necessary to
    // reload it, which is good because *types now points to the next value
    for(; (dType = (type & 0x7F)) <= PARAMETER_LAST_TYPE; type = *types++)
    {
        pNum++;
        if(dType < PARAMETER_FIRST_FLAG_TYPE)
        {
            NoFlagFunction      *f = (NoFlagFunction *)UnmarshalArray[dType];
            result = f(&commandIn[offset], &command->parameterBuffer,
                       &command->parameterSize);
        }
        else
        {
            FlagFunction        *f = UnmarshalArray[dType];
            result = f(&commandIn[offset], &command->parameterBuffer,
                       &command->parameterSize,
                       (type & 0x80) != 0);
        }
        if(result != TPM_RC_SUCCESS)
        {
            result += TPM_RC_P + (TPM_RC_1 * pNum);
            goto Exit;
        }

        // This check is used so that we don't have to add an additional offset
        // value to the offsets list to correspond to the stop value in the
        // command parameter list.
        if(*types != 0xFF)
            offset = *offsets++;
        hasInParameters++;
    }
    // Should have used all the bytes in the input
    if(command->parameterSize != 0)
    {
        result = TPM_RC_SIZE;
        goto Exit;
    }

    // The command parameter unmarshaling stopped when it hit a value that was out
    // of range for unmarshaling values and left *types pointing to the first
    // marshaling type. If that type happens to be the STOP value, then there
    // are no response parameters. So, set the flag to indicate if there are
    // output parameters.
    hasOutParameters = *types != 0xFF;

    // There are four cases for calling, with and without input parameters and with
    // and without output parameters.
    if(hasInParameters > 0)
    {
        if(hasOutParameters)
            result = cmd.inOutArg(commandIn, commandOut);
        else
            result = cmd.inArg(commandIn);
    }
    else
    {
        if(hasOutParameters)
            result = cmd.outArg(commandOut);
        else
            result = cmd.noArgs();
    }
    if(result != TPM_RC_SUCCESS)
       goto Exit;

    // Offset in the marshaled output structure
    offset = 0;

    // Process the return handles, if any
    command->handleNum = 0;

    // Could make this a loop to process output handles but there is only ever
    // one handle in the outputs (for now).
    type = *types++;
    if((dType = (type & 0x7F)) < RESPONSE_PARAMETER_FIRST_TYPE)
    {
        // The out->handle value was referenced as TPM_HANDLE in the
        // action code so it has to be properly aligned.
        command->handles[command->handleNum++] =
            *((TPM_HANDLE *)&(commandOut[offset]));
        maxOutSize -= sizeof(UINT32);
        type = *types++;
        offset = *offsets++;
    }
    // Use the size of the command action output buffer as the maximum for the
    // number of bytes that can get marshaled. Since the marshaling code has
    // no pointers to data, all of the data being returned has to be in the
    // command action output buffer. If we try to marshal more bytes than
    // could fit into the output buffer, we need to fail.
    for(;(dType = (type & 0x7F)) <= RESPONSE_PARAMETER_LAST_TYPE 
        && !g_inFailureMode; type = *types++)
    {
        const MARSHAL_t     f = MarshalArray[dType];

        command->parameterSize += f(&commandOut[offset], 
                                    &command->responseBuffer,
                                    &maxOutSize);
        offset = *offsets++;
    }
    result = (maxOutSize < 0) ? TPM_RC_FAILURE : TPM_RC_SUCCESS;
Exit:
    MemoryIoBufferZero();
    return result;
#endif
}