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
//** Includes, Defines, and Types
#define     TPM_FAIL_C
#include    "Tpm.h"
#include    <assert.h>

// On MS C compiler, can save the alignment state and set the alignment to 1 for
// the duration of the TpmTypes.h include.  This will avoid a lot of alignment
// warnings from the compiler for the unaligned structures. The alignment of the
// structures is not important as this function does not use any of the structures
// in TpmTypes.h and only include it for the #defines of the capabilities,
// properties, and command code values.
#include "TpmTypes.h"

//** Typedefs
// These defines are used primarily for sizing of the local response buffer.
typedef struct 
{
    TPM_ST          tag;
    UINT32          size;
    TPM_RC          code;
} HEADER;

typedef struct 
{
    BYTE            tag[sizeof(TPM_ST)];
    BYTE            size[sizeof(UINT32)];
    BYTE            code[sizeof(TPM_RC)];
} PACKED_HEADER;

typedef struct 
{
    BYTE             size[sizeof(UINT16)];
    struct
    {
        BYTE         function[sizeof(UINT32)];
        BYTE         line[sizeof(UINT32)];
        BYTE         code[sizeof(UINT32)];
    } values;
    BYTE             returnCode[sizeof(TPM_RC)];
} GET_TEST_RESULT_PARAMETERS;

typedef struct
{
    BYTE         moreData[sizeof(TPMI_YES_NO)];
    BYTE         capability[sizeof(TPM_CAP)]; // Always TPM_CAP_TPM_PROPERTIES
    BYTE         tpmProperty[sizeof(TPML_TAGGED_TPM_PROPERTY)]; 
} GET_CAPABILITY_PARAMETERS;

typedef struct
{
    BYTE         header[sizeof(PACKED_HEADER)];
    BYTE         getTestResult[sizeof(GET_TEST_RESULT_PARAMETERS)];
} TEST_RESPONSE;

typedef struct
{
    BYTE         header[sizeof(PACKED_HEADER)];
    BYTE         getCap[sizeof(GET_CAPABILITY_PARAMETERS)];
} CAPABILITY_RESPONSE;

typedef union
{
    BYTE         test[sizeof(TEST_RESPONSE)];
    BYTE         cap[sizeof(CAPABILITY_RESPONSE)];
} RESPONSES;

// Buffer to hold the responses. This may be a little larger than
// required due to padding that a compiler might add.
// Note: This is not in Global.c because of the specialized data definitions above.
// Since the data contained in this structure is not relevant outside of the
// execution of a single command (when the TPM is in failure mode. There is no
// compelling reason to move all the typedefs to Global.h and this structure
// to Global.c.
#ifndef __IGNORE_STATE__ // Don't define this value
static BYTE response[sizeof(RESPONSES)];
#endif

//** Local Functions

//*** MarshalUint16()
// Function to marshal a 16 bit value to the output buffer.
static INT32
MarshalUint16(
    UINT16          integer,
    BYTE            **buffer
    )
{
    UINT16_TO_BYTE_ARRAY(integer, *buffer);
    *buffer += 2;
    return 2;
}

//*** MarshalUint32()
// Function to marshal a 32 bit value to the output buffer.
static INT32
MarshalUint32(
    UINT32           integer,
    BYTE            **buffer
    )
{
    UINT32_TO_BYTE_ARRAY(integer, *buffer);
    *buffer += 4;
    return 4;
}

//***Unmarshal32()
static BOOL Unmarshal32(
    UINT32          *target,
    BYTE            **buffer,
    INT32           *size
    )
{
    if((*size -= 4) < 0)
        return FALSE;
    *target = BYTE_ARRAY_TO_UINT32(*buffer);
    *buffer += 4;
    return TRUE;
}

//***Unmarshal16()
static BOOL Unmarshal16(
    UINT16          *target,
    BYTE           **buffer,
    INT32           *size
)
{
    if((*size -= 2) < 0)
        return FALSE;
    *target = BYTE_ARRAY_TO_UINT16(*buffer);
    *buffer += 2;
    return TRUE;
}

//** Public Functions

//*** SetForceFailureMode()
// This function is called by the simulator to enable failure mode testing.
#if SIMULATION
LIB_EXPORT void
SetForceFailureMode(
    void
    )
{
    g_forceFailureMode = TRUE;
    return;
}
#endif

//*** TpmLogFailure()
// This function saves the failure values when the code will continue to operate. It
// if similar to TpmFail() but returns to the caller. The assumption is that the
// caller will propagate a failure back up the stack.
void
TpmLogFailure(
#if FAIL_TRACE
    const char      *function,
    int              line,
#endif
    int              code
)
{
    // Save the values that indicate where the error occurred.
    // On a 64-bit machine, this may truncate the address of the string
    // of the function name where the error occurred.
#if FAIL_TRACE
    s_failFunction = (UINT32)(ptrdiff_t)function;
    s_failLine = line;
#else
    s_failFunction = 0;
    s_failLine = 0;
#endif
    s_failCode = code;

    // We are in failure mode
    g_inFailureMode = TRUE;

    return;
}

//*** TpmFail()
// This function is called by TPM.lib when a failure occurs. It will set up the
// failure values to be returned on TPM2_GetTestResult().
NORETURN void
TpmFail(
#if FAIL_TRACE
    const char      *function,
    int              line,
#endif
    int              code
    )
{
    // Save the values that indicate where the error occurred.
    // On a 64-bit machine, this may truncate the address of the string
    // of the function name where the error occurred.
#if FAIL_TRACE
    s_failFunction = (UINT32)(ptrdiff_t)function;
    s_failLine = line;
#else
    s_failFunction = (UINT32)(ptrdiff_t)NULL;
    s_failLine = 0;
#endif
    s_failCode = code;

    // We are in failure mode
    g_inFailureMode = TRUE;

    // if asserts are enabled, then do an assert unless the failure mode code
    // is being tested.
#if SIMULATION
#   ifndef NDEBUG
    assert(g_forceFailureMode);
#   endif
    // Clear this flag
    g_forceFailureMode = FALSE;
#endif
    // Jump to the failure mode code.
    // Note: only get here if asserts are off or if we are testing failure mode
    _plat__Fail();
}

//*** TpmFailureMode(
// This function is called by the interface code when the platform is in failure
// mode.
void
TpmFailureMode(
    unsigned int     inRequestSize,     // IN: command buffer size
    unsigned char   *inRequest,         // IN: command buffer
    unsigned int    *outResponseSize,   // OUT: response buffer size
    unsigned char   **outResponse       // OUT: response buffer
    )
{
    UINT32           marshalSize;
    UINT32           capability;
    HEADER           header;    // unmarshaled command header
    UINT32           pt;    // unmarshaled property type
    UINT32           count; // unmarshaled property count
    UINT8           *buffer = inRequest;
    INT32            size = inRequestSize;

    // If there is no command buffer, then just return TPM_RC_FAILURE
    if(inRequestSize == 0 || inRequest == NULL)
        goto FailureModeReturn;
    // If the header is not correct for TPM2_GetCapability() or
    // TPM2_GetTestResult() then just return the in failure mode response;
    if(! (Unmarshal16(&header.tag,  &buffer, &size)
       && Unmarshal32(&header.size, &buffer, &size)
       && Unmarshal32(&header.code, &buffer, &size)))
        goto FailureModeReturn;
    if(header.tag != TPM_ST_NO_SESSIONS
       || header.size < 10)
        goto FailureModeReturn;
    switch(header.code)
    {
        case TPM_CC_GetTestResult:
            // make sure that the command size is correct
            if(header.size != 10)
                goto FailureModeReturn;
            buffer = &response[10];
            marshalSize = MarshalUint16(3 * sizeof(UINT32), &buffer);
            marshalSize += MarshalUint32(s_failFunction, &buffer);
            marshalSize += MarshalUint32(s_failLine, &buffer);
            marshalSize += MarshalUint32(s_failCode, &buffer);
            if(s_failCode == FATAL_ERROR_NV_UNRECOVERABLE)
                marshalSize += MarshalUint32(TPM_RC_NV_UNINITIALIZED, &buffer);
            else
                marshalSize += MarshalUint32(TPM_RC_FAILURE, &buffer);
            break;
        case TPM_CC_GetCapability:
            // make sure that the size of the command is exactly the size
            // returned for the capability, property, and count
            if(header.size != (10 + (3 * sizeof(UINT32)))
                    // also verify that this is requesting TPM properties
               || !Unmarshal32(&capability, &buffer, &size)
               || capability != TPM_CAP_TPM_PROPERTIES
               || !Unmarshal32(&pt, &buffer, &size)
               || !Unmarshal32(&count, &buffer, &size))
                goto FailureModeReturn;
            // If in failure mode because of an unrecoverable read error, and the
            // property is 0 and the count is 0, then this is an indication to
            // re-manufacture the TPM. Do the re-manufacture but stay in failure
            // mode until the TPM is reset.
            // Note: this behavior is not required by the specification and it is
            // OK to leave the TPM permanently bricked due to an unrecoverable NV
            // error.
            if(count == 0 && pt == 0 && s_failCode == FATAL_ERROR_NV_UNRECOVERABLE)
            {
                g_manufactured = FALSE;
                TPM_Manufacture(0);
            }
            if(count > 0)
                count = 1;
            else if(pt > TPM_PT_FIRMWARE_VERSION_2)
                count = 0;
            if(pt < TPM_PT_MANUFACTURER)
                pt = TPM_PT_MANUFACTURER;
            // set up for return
            buffer = &response[10];
            // if the request was for a PT less than the last one
            // then we indicate more, otherwise, not.
            if(pt < TPM_PT_FIRMWARE_VERSION_2)
                *buffer++ = YES;
            else
                *buffer++ = NO;
            marshalSize = 1;

            // indicate the capability type
            marshalSize += MarshalUint32(capability, &buffer);
            // indicate the number of values that are being returned (0 or 1)
            marshalSize += MarshalUint32(count, &buffer);
            // indicate the property
            marshalSize += MarshalUint32(pt, &buffer);

            if(count > 0)
                switch(pt)
                {
                    case TPM_PT_MANUFACTURER:
                    // the vendor ID unique to each TPM manufacturer
#ifdef  MANUFACTURER
                        pt = *(UINT32*)MANUFACTURER;
#else
                        pt = 0;
#endif
                        break;
                    case TPM_PT_VENDOR_STRING_1:
                        // the first four characters of the vendor ID string
#ifdef  VENDOR_STRING_1
                        pt = *(UINT32*)VENDOR_STRING_1;
#else
                        pt = 0;
#endif
                        break;
                    case TPM_PT_VENDOR_STRING_2:
                        // the second four characters of the vendor ID string
#ifdef  VENDOR_STRING_2
                        pt = *(UINT32*)VENDOR_STRING_2;
#else
                        pt = 0;
#endif
                        break;
                    case TPM_PT_VENDOR_STRING_3:
                        // the third four characters of the vendor ID string
#ifdef  VENDOR_STRING_3
                        pt = *(UINT32*)VENDOR_STRING_3;
#else
                        pt = 0;
#endif
                        break;
                    case TPM_PT_VENDOR_STRING_4:
                        // the fourth four characters of the vendor ID string
#ifdef  VENDOR_STRING_4
                        pt = *(UINT32*)VENDOR_STRING_4;
#else
                        pt = 0;
#endif
                        break;
                    case TPM_PT_VENDOR_TPM_TYPE:
                        // vendor-defined value indicating the TPM model
                        // We just make up a number here
                        pt = 1;
                        break;
                    case TPM_PT_FIRMWARE_VERSION_1:
                        // the more significant 32-bits of a vendor-specific value
                        // indicating the version of the firmware
#ifdef  FIRMWARE_V1
                        pt = FIRMWARE_V1;
#else
                        pt = 0;
#endif
                        break;
                    default: // TPM_PT_FIRMWARE_VERSION_2:
                        // the less significant 32-bits of a vendor-specific value
                        // indicating the version of the firmware
#ifdef  FIRMWARE_V2
                        pt = FIRMWARE_V2;
#else
                        pt = 0;
#endif
                        break;
                }
            marshalSize += MarshalUint32(pt, &buffer);
            break;
        default: // default for switch (cc)
            goto FailureModeReturn;
    }
    // Now do the header
    buffer = response;
    marshalSize = marshalSize + 10; // Add the header size to the
                                    // stuff already marshaled
    MarshalUint16(TPM_ST_NO_SESSIONS, &buffer); // structure tag
    MarshalUint32(marshalSize, &buffer);  // responseSize
    MarshalUint32(TPM_RC_SUCCESS, &buffer); // response code

    *outResponseSize = marshalSize;
    *outResponse = (unsigned char *)&response;
    return;
FailureModeReturn:
    buffer = response;
    marshalSize = MarshalUint16(TPM_ST_NO_SESSIONS, &buffer);
    marshalSize += MarshalUint32(10, &buffer);
    marshalSize += MarshalUint32(TPM_RC_FAILURE, &buffer);
    *outResponseSize = marshalSize;
    *outResponse = (unsigned char *)response;
    return;
}

//*** UnmarshalFail()
// This is a stub that is used to catch an attempt to unmarshal an entry
// that is not defined. Don't ever expect this to be called but...
void
UnmarshalFail(
    void            *type,
    BYTE            **buffer,
    INT32           *size
    )
{
    NOT_REFERENCED(type);
    NOT_REFERENCED(buffer);
    NOT_REFERENCED(size);
    FAIL(FATAL_ERROR_INTERNAL);
}