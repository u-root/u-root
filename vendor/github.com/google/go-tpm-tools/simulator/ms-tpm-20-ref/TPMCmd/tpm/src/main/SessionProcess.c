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
//**  Introduction
// This file contains the subsystem that process the authorization sessions
// including implementation of the Dictionary Attack logic. ExecCommand() uses
// ParseSessionBuffer() to process the authorization session area of a command and
// BuildResponseSession() to create the authorization session area of a response.

//**  Includes and Data Definitions

#define SESSION_PROCESS_C

#include "Tpm.h"

//
//**  Authorization Support Functions
//

//*** IsDAExempted()
// This function indicates if a handle is exempted from DA logic.
// A handle is exempted if it is
//  1. a primary seed handle,
//  2. an object with noDA bit SET,
//  3. an NV Index with TPMA_NV_NO_DA bit SET, or
//  4. a PCR handle.
//
//  Return Type: BOOL
//      TRUE(1)         handle is exempted from DA logic
//      FALSE(0)        handle is not exempted from DA logic
BOOL
IsDAExempted(
    TPM_HANDLE       handle         // IN: entity handle
    )
{
    BOOL        result = FALSE;
//
    switch(HandleGetType(handle))
    {
        case TPM_HT_PERMANENT:
            // All permanent handles, other than TPM_RH_LOCKOUT, are exempt from
            // DA protection.
            result = (handle != TPM_RH_LOCKOUT);
            break;
        // When this function is called, a persistent object will have been loaded
        // into an object slot and assigned a transient handle.
        case TPM_HT_TRANSIENT:
        {
            TPMA_OBJECT     attributes = ObjectGetPublicAttributes(handle);
            result = IS_ATTRIBUTE(attributes, TPMA_OBJECT, noDA);
            break;
        }
        case TPM_HT_NV_INDEX:
        {
            NV_INDEX            *nvIndex = NvGetIndexInfo(handle, NULL);
            result = IS_ATTRIBUTE(nvIndex->publicArea.attributes, TPMA_NV, NO_DA);
            break;
        }
        case TPM_HT_PCR:
            // PCRs are always exempted from DA.
            result = TRUE;
            break;
        default:
            break;
    }
    return result;
}

//*** IncrementLockout()
// This function is called after an authorization failure that involves use of
// an authValue. If the entity referenced by the handle is not exempt from DA
// protection, then the failedTries counter will be incremented.
//
//  Return Type: TPM_RC
//      TPM_RC_AUTH_FAIL    authorization failure that caused DA lockout to increment
//      TPM_RC_BAD_AUTH     authorization failure did not cause DA lockout to 
//                          increment
static TPM_RC
IncrementLockout(
    UINT32           sessionIndex
    )
{
    TPM_HANDLE       handle = s_associatedHandles[sessionIndex];
    TPM_HANDLE       sessionHandle = s_sessionHandles[sessionIndex];
    SESSION         *session = NULL;
//
    // Don't increment lockout unless the handle associated with the session
    // is DA protected or the session is bound to a DA protected entity.
    if(sessionHandle == TPM_RS_PW)
    {
        if(IsDAExempted(handle))
            return TPM_RC_BAD_AUTH;
    }
    else
    {
        session = SessionGet(sessionHandle);
        // If the session is bound to lockout, then use that as the relevant
        // handle. This means that an authorization failure with a bound session
        // bound to lockoutAuth will take precedence over any other
        // lockout check
        if(session->attributes.isLockoutBound == SET)
            handle = TPM_RH_LOCKOUT;
        if(session->attributes.isDaBound == CLEAR
           && (IsDAExempted(handle) || session->attributes.includeAuth == CLEAR))
           // If the handle was changed to TPM_RH_LOCKOUT, this will not return
           // TPM_RC_BAD_AUTH
            return TPM_RC_BAD_AUTH;
    }
    if(handle == TPM_RH_LOCKOUT)
    {
        pAssert(gp.lockOutAuthEnabled == TRUE);

        // lockout is no longer enabled
        gp.lockOutAuthEnabled = FALSE;

        // For TPM_RH_LOCKOUT, if lockoutRecovery is 0, no need to update NV since
        // the lockout authorization will be reset at startup.
        if(gp.lockoutRecovery != 0)
        {
            if(NV_IS_AVAILABLE)
                // Update NV.
                NV_SYNC_PERSISTENT(lockOutAuthEnabled);
            else
                // No NV access for now. Put the TPM in pending mode.
                s_DAPendingOnNV = TRUE;
        }
    }
    else
    {
        if(gp.recoveryTime != 0)
        {
            gp.failedTries++;
            if(NV_IS_AVAILABLE)
                // Record changes to NV. NvWrite will SET g_updateNV
                NV_SYNC_PERSISTENT(failedTries);
            else
                // No NV access for now.  Put the TPM in pending mode.
                s_DAPendingOnNV = TRUE;
        }
    }
    // Register a DA failure and reset the timers.
    DARegisterFailure(handle);

    return TPM_RC_AUTH_FAIL;
}

//*** IsSessionBindEntity()
// This function indicates if the entity associated with the handle is the entity,
// to which this session is bound. The binding would occur by making the "bind"
// parameter in TPM2_StartAuthSession() not equal to TPM_RH_NULL. The binding only
// occurs if the session is an HMAC session. The bind value is a combination of
// the Name and the authValue of the entity.
//
//  Return Type: BOOL
//      TRUE(1)         handle points to the session start entity
//      FALSE(0)        handle does not point to the session start entity
static BOOL
IsSessionBindEntity(
    TPM_HANDLE       associatedHandle,  // IN: handle to be authorized
    SESSION         *session            // IN: associated session
    )
{
    TPM2B_NAME     entity;             // The bind value for the entity
//
    // If the session is not bound, return FALSE.
    if(session->attributes.isBound)
    {
        // Compute the bind value for the entity.
        SessionComputeBoundEntity(associatedHandle, &entity);

        // Compare to the bind value in the session.
        return MemoryEqual2B(&entity.b, &session->u1.boundEntity.b);
    }
    return FALSE;
}

//*** IsPolicySessionRequired()
// Checks if a policy session is required for a command. If a command requires
// DUP or ADMIN role authorization, then the handle that requires that role is the
// first handle in the command. This simplifies this checking. If a new command
// is created that requires multiple ADMIN role authorizations, then it will
// have to be special-cased in this function.
// A policy session is required if:
//      1. the command requires the DUP role,
//      2. the command requires the ADMIN role and the authorized entity
//         is an object and its adminWithPolicy bit is SET, or
//      3. the command requires the ADMIN role and the authorized entity
//         is a permanent handle or an NV Index.
//      4. The authorized entity is a PCR belonging to a policy group, and
//         has its policy initialized
//  Return Type: BOOL
//      TRUE(1)         policy session is required
//      FALSE(0)        policy session is not required
static BOOL
IsPolicySessionRequired(
    COMMAND_INDEX    commandIndex,  // IN: command index
    UINT32           sessionIndex   // IN: session index
    )
{
    AUTH_ROLE       role = CommandAuthRole(commandIndex, sessionIndex);
    TPM_HT          type = HandleGetType(s_associatedHandles[sessionIndex]);
//
    if(role == AUTH_DUP)
        return TRUE;
    if(role == AUTH_ADMIN)
    {
        // We allow an exception for ADMIN role in a transient object. If the object
        // allows ADMIN role actions with authorization, then policy is not
        // required. For all other cases, there is no way to override the command
        // requirement that a policy be used
        if(type == TPM_HT_TRANSIENT)
        {
            OBJECT      *object = HandleToObject(s_associatedHandles[sessionIndex]);

            if(!IS_ATTRIBUTE(object->publicArea.objectAttributes, TPMA_OBJECT, 
                             adminWithPolicy))
                return FALSE;
        }
        return TRUE;
    }

    if(type == TPM_HT_PCR)
    {
        if(PCRPolicyIsAvailable(s_associatedHandles[sessionIndex]))
        {
            TPM2B_DIGEST        policy;
            TPMI_ALG_HASH       policyAlg;
            policyAlg = PCRGetAuthPolicy(s_associatedHandles[sessionIndex],
                                         &policy);
            if(policyAlg != TPM_ALG_NULL)
                return TRUE;
        }
    }
    return FALSE;
}

//*** IsAuthValueAvailable()
// This function indicates if authValue is available and allowed for USER role
// authorization of an entity.
//
// This function is similar to IsAuthPolicyAvailable() except that it does not
// check the size of the authValue as IsAuthPolicyAvailable() does (a null
// authValue is a valid authorization, but a null policy is not a valid policy).
//
// This function does not check that the handle reference is valid or if the entity
// is in an enabled hierarchy. Those checks are assumed to have been performed
// during the handle unmarshaling.
//
//  Return Type: BOOL
//      TRUE(1)         authValue is available
//      FALSE(0)        authValue is not available
static BOOL
IsAuthValueAvailable(
    TPM_HANDLE       handle,        // IN: handle of entity
    COMMAND_INDEX    commandIndex,  // IN: command index
    UINT32           sessionIndex   // IN: session index
    )
{
    BOOL             result = FALSE;
//
    switch(HandleGetType(handle))
    {
        case TPM_HT_PERMANENT:
            switch(handle)
            {
                    // At this point hierarchy availability has already been
                    // checked so primary seed handles are always available here
                case TPM_RH_OWNER:
                case TPM_RH_ENDORSEMENT:
                case TPM_RH_PLATFORM:
#ifdef VENDOR_PERMANENT
                    // This vendor defined handle associated with the
                    // manufacturer's shared secret
                case VENDOR_PERMANENT:
#endif
                    // The DA checking has been performed on LockoutAuth but we
                    // bypass the DA logic if we are using lockout policy. The
                    // policy would allow execution to continue an lockoutAuth
                    // could be used, even if direct use of lockoutAuth is disabled
                case TPM_RH_LOCKOUT:
                    // NullAuth is always available.
                case TPM_RH_NULL:
                    result = TRUE;
                    break;
                default:
                    // Otherwise authValue is not available.
                    break;
            }
            break;
        case TPM_HT_TRANSIENT:
            // A persistent object has already been loaded and the internal
            // handle changed.
        {
            OBJECT          *object;
            TPMA_OBJECT      attributes;
//
            object = HandleToObject(handle);
            attributes = object->publicArea.objectAttributes;

            // authValue is always available for a sequence object.
            // An alternative for this is to 
            // SET_ATTRIBUTE(object->publicArea, TPMA_OBJECT, userWithAuth) when the
            // sequence is started.
            if(ObjectIsSequence(object))
            {
                result = TRUE;
                break;
            }
            // authValue is available for an object if it has its sensitive
            // portion loaded and
            //  1. userWithAuth bit is SET, or
            //  2. ADMIN role is required
            if(object->attributes.publicOnly == CLEAR
               && (IS_ATTRIBUTE(attributes, TPMA_OBJECT, userWithAuth)
                   || (CommandAuthRole(commandIndex, sessionIndex) == AUTH_ADMIN
                       && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, adminWithPolicy))))
                result = TRUE;
        }
        break;
        case TPM_HT_NV_INDEX:
            // NV Index.
        {
            NV_REF           locator;
            NV_INDEX        *nvIndex = NvGetIndexInfo(handle, &locator);
            TPMA_NV          nvAttributes;
//
            pAssert(nvIndex != 0);
            
            nvAttributes = nvIndex->publicArea.attributes;

            if(IsWriteOperation(commandIndex))
            {
                // AuthWrite can't be set for a PIN index
                if(IS_ATTRIBUTE(nvAttributes, TPMA_NV, AUTHWRITE))
                    result = TRUE;
            }
            else
            {
                // A "read" operation
                // For a PIN Index, the authValue is available as long as the
                // Index has been written and the pinCount is less than pinLimit
                if(IsNvPinFailIndex(nvAttributes)
                   || IsNvPinPassIndex(nvAttributes))
                {
                    NV_PIN          pin;
                    if(!IS_ATTRIBUTE(nvAttributes, TPMA_NV, WRITTEN))
                        break; // return false
                    // get the index values
                    pin.intVal = NvGetUINT64Data(nvIndex, locator);
                    if(pin.pin.pinCount < pin.pin.pinLimit)
                        result = TRUE;
                }
                // For non-PIN Indexes, need to allow use of the authValue
                else if(IS_ATTRIBUTE(nvAttributes, TPMA_NV, AUTHREAD))
                    result = TRUE;
            }
        }
        break;
        case TPM_HT_PCR:
            // PCR handle.
            // authValue is always allowed for PCR
            result = TRUE;
            break;
        default:
            // Otherwise, authValue is not available
            break;
    }
    return result;
}

//*** IsAuthPolicyAvailable()
// This function indicates if an authPolicy is available and allowed.
//
// This function does not check that the handle reference is valid or if the entity
// is in an enabled hierarchy. Those checks are assumed to have been performed
// during the handle unmarshaling.
//
//  Return Type: BOOL
//      TRUE(1)         authPolicy is available
//      FALSE(0)        authPolicy is not available
static BOOL
IsAuthPolicyAvailable(
    TPM_HANDLE       handle,        // IN: handle of entity
    COMMAND_INDEX    commandIndex,  // IN: command index
    UINT32           sessionIndex   // IN: session index
    )
{
    BOOL            result = FALSE;
//
    switch(HandleGetType(handle))
    {
        case TPM_HT_PERMANENT:
            switch(handle)
            {
                // At this point hierarchy availability has already been checked.
                case TPM_RH_OWNER:
                    if(gp.ownerPolicy.t.size != 0)
                        result = TRUE;
                    break;
                case TPM_RH_ENDORSEMENT:
                    if(gp.endorsementPolicy.t.size != 0)
                        result = TRUE;
                    break;
                case TPM_RH_PLATFORM:
                    if(gc.platformPolicy.t.size != 0)
                        result = TRUE;
                    break;
                case TPM_RH_LOCKOUT:
                    if(gp.lockoutPolicy.t.size != 0)
                        result = TRUE;
                    break;
                default:
                    break;
            }
            break;
        case TPM_HT_TRANSIENT:
        {
            // Object handle.
            // An evict object would already have been loaded and given a
            // transient object handle by this point.
            OBJECT  *object = HandleToObject(handle);
            // Policy authorization is not available for an object with only
            // public portion loaded.
            if(object->attributes.publicOnly == CLEAR)
            {
                // Policy authorization is always available for an object but
                // is never available for a sequence.
                if(!ObjectIsSequence(object))
                    result = TRUE;
            }
            break;
        }
        case TPM_HT_NV_INDEX:
            // An NV Index.
        {
            NV_INDEX         *nvIndex = NvGetIndexInfo(handle, NULL);
            TPMA_NV           nvAttributes = nvIndex->publicArea.attributes;
//
            // If the policy size is not zero, check if policy can be used.
            if(nvIndex->publicArea.authPolicy.t.size != 0)
            {
                // If policy session is required for this handle, always
                // uses policy regardless of the attributes bit setting
                if(IsPolicySessionRequired(commandIndex, sessionIndex))
                    result = TRUE;
                // Otherwise, the presence of the policy depends on the NV
                // attributes.
                else if(IsWriteOperation(commandIndex))
                {
                    if(IS_ATTRIBUTE(nvAttributes, TPMA_NV, POLICYWRITE))
                        result = TRUE;
                }
                else
                {
                    if(IS_ATTRIBUTE(nvAttributes, TPMA_NV, POLICYREAD))
                        result = TRUE;
                }
            }
        }
        break;
        case TPM_HT_PCR:
            // PCR handle.
            if(PCRPolicyIsAvailable(handle))
                result = TRUE;
            break;
        default:
            break;
    }
    return result;
}

//**  Session Parsing Functions

//*** ClearCpRpHashes()
void
ClearCpRpHashes(
    COMMAND         *command
    )
{
#if ALG_SHA1
    command->sha1CpHash.t.size = 0;
    command->sha1RpHash.t.size = 0;
#endif
#if ALG_SHA256
    command->sha256CpHash.t.size = 0;
    command->sha256RpHash.t.size = 0;
#endif
#if ALG_SHA384
    command->sha384CpHash.t.size = 0;
    command->sha384RpHash.t.size = 0;
#endif
#if ALG_SHA512
    command->sha512CpHash.t.size = 0;
    command->sha512RpHash.t.size = 0;
#endif
#if ALG_SM3_256
    command->sm3_256CpHash.t.size = 0;
    command->sm3_256RpHash.t.size = 0;
#endif
}


//*** GetCpHashPointer()
// Function to get a pointer to the cpHash of the command
static TPM2B_DIGEST *
GetCpHashPointer(
    COMMAND         *command,
    TPMI_ALG_HASH    hashAlg
    )
{
    TPM2B_DIGEST     *retVal;
//
    switch(hashAlg)
    {
#if ALG_SHA1
        case ALG_SHA1_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha1CpHash;
			break;
#endif
#if ALG_SHA256
        case ALG_SHA256_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha256CpHash;
			break;
#endif
#if ALG_SHA384
        case ALG_SHA384_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha384CpHash;
			break;
#endif
#if ALG_SHA512
        case ALG_SHA512_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha512CpHash;
			break;
#endif
#if ALG_SM3_256
        case ALG_SM3_256_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sm3_256CpHash;
			break;
#endif
        default:
            retVal = NULL;
            break;
    }
    return retVal;
}

//*** GetRpHashPointer()
// Function to get a pointer to the RpHash of the command
static TPM2B_DIGEST *
GetRpHashPointer(
    COMMAND         *command,
    TPMI_ALG_HASH    hashAlg
    )
{
    TPM2B_DIGEST    *retVal;
//
    switch(hashAlg)
    {
#if ALG_SHA1
        case ALG_SHA1_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha1RpHash;
			break;
#endif
#if ALG_SHA256
        case ALG_SHA256_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha256RpHash;
			break;
#endif
#if ALG_SHA384
        case ALG_SHA384_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha384RpHash;
			break;
#endif
#if ALG_SHA512
        case ALG_SHA512_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sha512RpHash;
			break;
#endif
#if ALG_SM3_256
        case ALG_SM3_256_VALUE:
            retVal = (TPM2B_DIGEST *)&command->sm3_256RpHash;
			break;
#endif
        default:
            retVal = NULL;
            break;
    }
    return retVal;
}


//*** ComputeCpHash()
// This function computes the cpHash as defined in Part 2 and described in Part 1.
static TPM2B_DIGEST *
ComputeCpHash(
    COMMAND         *command,       // IN: command parsing structure
    TPMI_ALG_HASH    hashAlg        // IN: hash algorithm
    )
{
    UINT32               i;
    HASH_STATE           hashState;
    TPM2B_NAME           name;
    TPM2B_DIGEST        *cpHash;
//
    // cpHash = hash(commandCode [ || authName1
    //                           [ || authName2
    //                           [ || authName 3 ]]]
    //                           [ || parameters])
    // A cpHash can contain just a commandCode only if the lone session is
    // an audit session.
    // Get pointer to the hash value
    cpHash = GetCpHashPointer(command, hashAlg);
    if(cpHash->t.size == 0)
    {
        cpHash->t.size = CryptHashStart(&hashState, hashAlg);
            //  Add commandCode.
        CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), command->code);
            //  Add authNames for each of the handles.
        for(i = 0; i < command->handleNum; i++)
            CryptDigestUpdate2B(&hashState, &EntityGetName(command->handles[i],
                                                           &name)->b);
            //  Add the parameters.
        CryptDigestUpdate(&hashState, command->parameterSize, 
                          command->parameterBuffer);
            //  Complete the hash.
        CryptHashEnd2B(&hashState, &cpHash->b);
    }
    return cpHash;
}

//*** GetCpHash()
// This function is used to access a precomputed cpHash.
static TPM2B_DIGEST *
GetCpHash(
    COMMAND         *command,
    TPMI_ALG_HASH    hashAlg
    )
{
    TPM2B_DIGEST        *cpHash = GetCpHashPointer(command, hashAlg);
 //
    pAssert(cpHash->t.size != 0);
    return cpHash;
}

//*** CompareTemplateHash()
// This function computes the template hash and compares it to the session
// templateHash. It is the hash of the second parameter
// assuming that the command is TPM2_Create(), TPM2_CreatePrimary(), or
// TPM2_CreateLoaded()
//  Return Type: BOOL
//      TRUE(1)         template hash equal to session->templateHash
//      FALSE(0)        template hash not equal to session->templateHash
static BOOL
CompareTemplateHash(
    COMMAND         *command,       // IN: parsing structure
    SESSION         *session        // IN: session data
    )
{
    BYTE                *pBuffer = command->parameterBuffer;
    INT32                pSize = command->parameterSize;
    TPM2B_DIGEST         tHash;
    UINT16               size;
//
    // Only try this for the three commands for which it is intended
    if(command->code != TPM_CC_Create
       && command->code != TPM_CC_CreatePrimary
#if CC_CreateLoaded
       && command->code != TPM_CC_CreateLoaded
#endif
       )
        return FALSE;
    // Assume that the first parameter is a TPM2B and unmarshal the size field
    // Note: this will not affect the parameter buffer and size in the calling
    // function.
    if(UINT16_Unmarshal(&size, &pBuffer, &pSize) != TPM_RC_SUCCESS)
        return FALSE;
    // reduce the space in the buffer.
    // NOTE: this could make pSize go negative if the parameters are not correct but
    // the unmarshaling code does not try to unmarshal if the remaining size is
    // negative.
    pSize -= size;

    // Advance the pointer
    pBuffer += size;

    // Get the size of what should be the template
    if(UINT16_Unmarshal(&size, &pBuffer, &pSize) != TPM_RC_SUCCESS)
        return FALSE;
    // See if this is reasonable
    if(size > pSize)
        return FALSE;
    // Hash the template data
    tHash.t.size = CryptHashBlock(session->authHashAlg, size, pBuffer, 
                                  sizeof(tHash.t.buffer), tHash.t.buffer);
    return(MemoryEqual2B(&session->u1.templateHash.b, &tHash.b));
}

//*** CompareNameHash()
// This function computes the name hash and compares it to the nameHash in the
// session data.
BOOL
CompareNameHash(
    COMMAND         *command,       // IN: main parsing structure
    SESSION         *session        // IN: session structure with nameHash
    )
{
    HASH_STATE           hashState;
    TPM2B_DIGEST         nameHash;
    UINT32               i;
    TPM2B_NAME           name;
//
    nameHash.t.size = CryptHashStart(&hashState, session->authHashAlg);
    //  Add names.
    for(i = 0; i < command->handleNum; i++)
        CryptDigestUpdate2B(&hashState, &EntityGetName(command->handles[i],
                                                       &name)->b);
    //  Complete hash.
    CryptHashEnd2B(&hashState, &nameHash.b);
    // and compare
    return MemoryEqual(session->u1.nameHash.t.buffer, nameHash.t.buffer,
                       nameHash.t.size);
}

//*** CheckPWAuthSession()
// This function validates the authorization provided in a PWAP session. It
// compares the input value to authValue of the authorized entity. Argument
// sessionIndex is used to get handles handle of the referenced entities from
// s_inputAuthValues[] and s_associatedHandles[].
//
//  Return Type: TPM_RC
//        TPM_RC_AUTH_FAIL          authorization fails and increments DA failure
//                                  count
//        TPM_RC_BAD_AUTH           authorization fails but DA does not apply
//
static TPM_RC
CheckPWAuthSession(
    UINT32           sessionIndex   // IN: index of session to be processed
    )
{
    TPM2B_AUTH      authValue;
    TPM_HANDLE      associatedHandle = s_associatedHandles[sessionIndex];
//
    // Strip trailing zeros from the password.
    MemoryRemoveTrailingZeros(&s_inputAuthValues[sessionIndex]);

    // Get the authValue with trailing zeros removed
    EntityGetAuthValue(associatedHandle, &authValue);

    // Success if the values are identical.
    if(MemoryEqual2B(&s_inputAuthValues[sessionIndex].b, &authValue.b))
    {
        return TPM_RC_SUCCESS;
    }
    else                    // if the digests are not identical
    {
        // Invoke DA protection if applicable.
        return IncrementLockout(sessionIndex);
    }
}

//*** ComputeCommandHMAC()
// This function computes the HMAC for an authorization session in a command.
/*(See part 1 specification -- this tag keeps this comment from showing up in
// merged document which is probably good because this comment doesn't look right.
//      The sessionAuth value
//      authHMAC := HMACsHash((sessionKey | authValue),
//                  (pHash | nonceNewer | nonceOlder  | nonceTPMencrypt-only
//                   | nonceTPMaudit   | sessionAttributes))
// Where:
//      HMACsHash()     The HMAC algorithm using the hash algorithm specified
//                      when the session was started.
//
//      sessionKey      A value that is computed in a protocol-dependent way,
//                      using KDFa. When used in an HMAC or KDF, the size field
//                      for this value is not included.
//
//      authValue       A value that is found in the sensitive area of an entity.
//                      When used in an HMAC or KDF, the size field for this
//                      value is not included.
//
//      pHash           Hash of the command (cpHash) using the session hash.
//                      When using a pHash in an HMAC computation, only the
//                      digest is used.
//
//      nonceNewer      A value that is generated by the entity using the
//                      session. A new nonce is generated on each use of the
//                      session. For a command, this will be nonceCaller.
//                      When used in an HMAC or KDF, the size field is not used.
//
//      nonceOlder      A TPM2B_NONCE that was received the previous time the
//                      session was used. For a command, this is nonceTPM.
//                      When used in an HMAC or KDF, the size field is not used.
//
//      nonceTPMdecrypt     The nonceTPM of the decrypt session is included in
//                          the HMAC, but only in the command.
//
//      nonceTPMencrypt     The nonceTPM of the encrypt session is included in
//                          the HMAC but only in the command.
//
//      sessionAttributes   A byte indicating the attributes associated with the
//                          particular use of the session.
*/
static TPM2B_DIGEST *
ComputeCommandHMAC(
    COMMAND         *command,       // IN: primary control structure
    UINT32           sessionIndex,  // IN: index of session to be processed
    TPM2B_DIGEST    *hmac           // OUT: authorization HMAC
    )
{
    TPM2B_TYPE(KEY, (sizeof(AUTH_VALUE) * 2));
    TPM2B_KEY        key;
    BYTE             marshalBuffer[sizeof(TPMA_SESSION)];
    BYTE            *buffer;
    UINT32           marshalSize;
    HMAC_STATE       hmacState;
    TPM2B_NONCE     *nonceDecrypt;
    TPM2B_NONCE     *nonceEncrypt;
    SESSION         *session;
//
    nonceDecrypt = NULL;
    nonceEncrypt = NULL;

    // Determine if extra nonceTPM values are going to be required.
    // If this is the first session (sessionIndex = 0) and it is an authorization
    // session that uses an HMAC, then check if additional session nonces are to be
    // included.
    if(sessionIndex == 0
       && s_associatedHandles[sessionIndex] != TPM_RH_UNASSIGNED)
    {
        // If there is a decrypt session and if this is not the decrypt session,
        // then an extra nonce may be needed.
        if(s_decryptSessionIndex != UNDEFINED_INDEX
           && s_decryptSessionIndex != sessionIndex)
        {
            // Will add the nonce for the decrypt session.
            SESSION *decryptSession
                = SessionGet(s_sessionHandles[s_decryptSessionIndex]);
            nonceDecrypt = &decryptSession->nonceTPM;
        }
        // Now repeat for the encrypt session.
        if(s_encryptSessionIndex != UNDEFINED_INDEX
           && s_encryptSessionIndex != sessionIndex
           && s_encryptSessionIndex != s_decryptSessionIndex)
        {
            // Have to have the nonce for the encrypt session.
            SESSION *encryptSession
                = SessionGet(s_sessionHandles[s_encryptSessionIndex]);
            nonceEncrypt = &encryptSession->nonceTPM;
        }
    }

    // Continue with the HMAC processing.
    session = SessionGet(s_sessionHandles[sessionIndex]);

    // Generate HMAC key.
    MemoryCopy2B(&key.b, &session->sessionKey.b, sizeof(key.t.buffer));

    // Check if the session has an associated handle and if the associated entity
    // is the one to which the session is bound. If not, add the authValue of
    // this entity to the HMAC key.
    // If the session is bound to the object or the session is a policy session
    // with no authValue required, do not include the authValue in the HMAC key.
    // Note: For a policy session, its isBound attribute is CLEARED.
    //
    // Include the entity authValue if it is needed
    if(session->attributes.includeAuth == SET)
    {
        TPM2B_AUTH          authValue;
        // Get the entity authValue with trailing zeros removed
        EntityGetAuthValue(s_associatedHandles[sessionIndex], &authValue);
        // add the authValue to the HMAC key
        MemoryConcat2B(&key.b, &authValue.b, sizeof(key.t.buffer));
    }
     // if the HMAC key size is 0, a NULL string HMAC is allowed
    if(key.t.size == 0
       && s_inputAuthValues[sessionIndex].t.size == 0)
    {
        hmac->t.size = 0;
        return hmac;
    }
    // Start HMAC
    hmac->t.size = CryptHmacStart2B(&hmacState, session->authHashAlg, &key.b);

        //  Add cpHash
    CryptDigestUpdate2B(&hmacState.hashState,
                        &ComputeCpHash(command, session->authHashAlg)->b);
        //  Add nonces as required
    CryptDigestUpdate2B(&hmacState.hashState, &s_nonceCaller[sessionIndex].b);
    CryptDigestUpdate2B(&hmacState.hashState, &session->nonceTPM.b);
    if(nonceDecrypt != NULL)
        CryptDigestUpdate2B(&hmacState.hashState, &nonceDecrypt->b);
    if(nonceEncrypt != NULL)
        CryptDigestUpdate2B(&hmacState.hashState, &nonceEncrypt->b);
        //  Add sessionAttributes
    buffer = marshalBuffer;
    marshalSize = TPMA_SESSION_Marshal(&(s_attributes[sessionIndex]),
                                       &buffer, NULL);
    CryptDigestUpdate(&hmacState.hashState, marshalSize, marshalBuffer);
        // Complete the HMAC computation
    CryptHmacEnd2B(&hmacState, &hmac->b);

    return hmac;
}

//*** CheckSessionHMAC()
// This function checks the HMAC of in a session. It uses ComputeCommandHMAC()
// to compute the expected HMAC value and then compares the result with the
// HMAC in the authorization session. The authorization is successful if they
// are the same.
//
// If the authorizations are not the same, IncrementLockout() is called. It will
// return TPM_RC_AUTH_FAIL if the failure caused the failureCount to increment.
// Otherwise, it will return TPM_RC_BAD_AUTH.
//
//  Return Type: TPM_RC
//      TPM_RC_AUTH_FAIL        authorization failure caused failureCount increment
//      TPM_RC_BAD_AUTH         authorization failure did not cause failureCount
//                              increment
//
static TPM_RC
CheckSessionHMAC(
    COMMAND         *command,       // IN: primary control structure
    UINT32           sessionIndex   // IN: index of session to be processed
    )
{
    TPM2B_DIGEST        hmac;           // authHMAC for comparing
//
    // Compute authHMAC
    ComputeCommandHMAC(command, sessionIndex, &hmac);

    // Compare the input HMAC with the authHMAC computed above.
    if(!MemoryEqual2B(&s_inputAuthValues[sessionIndex].b, &hmac.b))
    {
        // If an HMAC session has a failure, invoke the anti-hammering
        // if it applies to the authorized entity or the session.
        // Otherwise, just indicate that the authorization is bad.
        return IncrementLockout(sessionIndex);
    }
    return TPM_RC_SUCCESS;
}

//*** CheckPolicyAuthSession()
//  This function is used to validate the authorization in a policy session.
//  This function performs the following comparisons to see if a policy
//  authorization is properly provided. The check are:
//  1. compare policyDigest in session with authPolicy associated with
//     the entity to be authorized;
//  2. compare timeout if applicable;
//  3. compare commandCode if applicable;
//  4. compare cpHash if applicable; and
//  5. see if PCR values have changed since computed.
//
// If all the above checks succeed, the handle is authorized.
// The order of these comparisons is not important because any failure will
// result in the same error code.
//
//  Return Type: TPM_RC
//      TPM_RC_PCR_CHANGED          PCR value is not current
//      TPM_RC_POLICY_FAIL          policy session fails
//      TPM_RC_LOCALITY             command locality is not allowed
//      TPM_RC_POLICY_CC            CC doesn't match
//      TPM_RC_EXPIRED              policy session has expired
//      TPM_RC_PP                   PP is required but not asserted
//      TPM_RC_NV_UNAVAILABLE       NV is not available for write
//      TPM_RC_NV_RATE              NV is rate limiting
static TPM_RC
CheckPolicyAuthSession(
    COMMAND         *command,       // IN: primary parsing structure
    UINT32           sessionIndex   // IN: index of session to be processed
    )
{
    SESSION             *session;
    TPM2B_DIGEST         authPolicy;
    TPMI_ALG_HASH        policyAlg;
    UINT8                locality;
//
    // Initialize pointer to the authorization session.
    session = SessionGet(s_sessionHandles[sessionIndex]);

    // If the command is TPM2_PolicySecret(), make sure that
    // either password or authValue is required
    if(command->code == TPM_CC_PolicySecret
       &&  session->attributes.isPasswordNeeded == CLEAR
       &&  session->attributes.isAuthValueNeeded == CLEAR)
        return TPM_RC_MODE;
    // See if the PCR counter for the session is still valid.
    if(!SessionPCRValueIsCurrent(session))
        return TPM_RC_PCR_CHANGED;
    // Get authPolicy.
    policyAlg = EntityGetAuthPolicy(s_associatedHandles[sessionIndex],
                                    &authPolicy);
    // Compare authPolicy.
    if(!MemoryEqual2B(&session->u2.policyDigest.b, &authPolicy.b))
        return TPM_RC_POLICY_FAIL;
    // Policy is OK so check if the other factors are correct

    // Compare policy hash algorithm.
    if(policyAlg != session->authHashAlg)
        return TPM_RC_POLICY_FAIL;

    // Compare timeout.
    if(session->timeout != 0)
    {
        // Cannot compare time if clock stop advancing.  An TPM_RC_NV_UNAVAILABLE
        // or TPM_RC_NV_RATE error may be returned here. This doesn't mean that
        // a new nonce will be created just that, because TPM time can't advance
        // we can't do time-based operations.
        RETURN_IF_NV_IS_NOT_AVAILABLE;

        if((session->timeout < g_time)
           || (session->epoch != g_timeEpoch))
            return TPM_RC_EXPIRED;
    }
    // If command code is provided it must match
    if(session->commandCode != 0)
    {
        if(session->commandCode != command->code)
            return TPM_RC_POLICY_CC;
    }
    else
    {
        // If command requires a DUP or ADMIN authorization, the session must have
        // command code set.
        AUTH_ROLE   role = CommandAuthRole(command->index, sessionIndex);
        if(role == AUTH_ADMIN || role == AUTH_DUP)
            return TPM_RC_POLICY_FAIL;
    }
    // Check command locality.
    {
        BYTE         sessionLocality[sizeof(TPMA_LOCALITY)];
        BYTE        *buffer = sessionLocality;

        // Get existing locality setting in canonical form
        sessionLocality[0] = 0;
        TPMA_LOCALITY_Marshal(&session->commandLocality, &buffer, NULL);

        // See if the locality has been set
        if(sessionLocality[0] != 0)
        {
            // If so, get the current locality
            locality = _plat__LocalityGet();
            if(locality < 5)
            {
                if(((sessionLocality[0] & (1 << locality)) == 0)
                   || sessionLocality[0] > 31)
                    return TPM_RC_LOCALITY;
            }
            else if(locality > 31)
            {
                if(sessionLocality[0] != locality)
                    return TPM_RC_LOCALITY;
            }
            else
            {
                // Could throw an assert here but a locality error is just
                // as good. It just means that, whatever the locality is, it isn't
                // the locality requested so...
                return TPM_RC_LOCALITY;
            }
        }
    } // end of locality check
    // Check physical presence.
    if(session->attributes.isPPRequired == SET
       && !_plat__PhysicalPresenceAsserted())
        return TPM_RC_PP;
    // Compare cpHash/nameHash if defined, or if the command requires an ADMIN or
    // DUP role for this handle.
    if(session->u1.cpHash.b.size != 0)
    {
        BOOL        OK;
        if(session->attributes.isCpHashDefined)
            // Compare cpHash.
            OK = MemoryEqual2B(&session->u1.cpHash.b,
                               &ComputeCpHash(command, session->authHashAlg)->b);
        else if(session->attributes.isTemplateSet)
            OK = CompareTemplateHash(command, session);
        else
            OK = CompareNameHash(command, session);
        if(!OK)
            return TPM_RCS_POLICY_FAIL;
    }
    if(session->attributes.checkNvWritten)
    {
        NV_REF           locator;
        NV_INDEX        *nvIndex;
//
        // If this is not an NV index, the policy makes no sense so fail it.
        if(HandleGetType(s_associatedHandles[sessionIndex]) != TPM_HT_NV_INDEX)
            return TPM_RC_POLICY_FAIL;
        // Get the index data
        nvIndex = NvGetIndexInfo(s_associatedHandles[sessionIndex], &locator);

        // Make sure that the TPMA_WRITTEN_ATTRIBUTE has the desired state
        if((IS_ATTRIBUTE(nvIndex->publicArea.attributes, TPMA_NV, WRITTEN))
           != (session->attributes.nvWrittenState == SET))
            return TPM_RC_POLICY_FAIL;
    }
    return TPM_RC_SUCCESS;
}

//*** RetrieveSessionData()
// This function will unmarshal the sessions in the session area of a command. The
// values are placed in the arrays that are defined at the beginning of this file.
// The normal unmarshaling errors are possible.
//
//  Return Type: TPM_RC
//      TPM_RC_SUCCSS       unmarshaled without error
//      TPM_RC_SIZE         the number of bytes unmarshaled is not the same
//                          as the value for authorizationSize in the command
//
static TPM_RC
RetrieveSessionData(
    COMMAND         *command        // IN: main parsing structure for command
    )
{
    int              i;
    TPM_RC           result;
    SESSION         *session;
    TPMA_SESSION     sessionAttributes;
    TPM_HT           sessionType;
    INT32            sessionIndex;
    TPM_RC           errorIndex;
//
    s_decryptSessionIndex = UNDEFINED_INDEX;
    s_encryptSessionIndex = UNDEFINED_INDEX;
    s_auditSessionIndex = UNDEFINED_INDEX;

    for(sessionIndex = 0; command->authSize > 0; sessionIndex++)
    {
        errorIndex = TPM_RC_S + g_rcIndex[sessionIndex];

        // If maximum allowed number of sessions has been parsed, return a size
        // error with a session number that is larger than the number of allowed
        // sessions
        if(sessionIndex == MAX_SESSION_NUM)
            return TPM_RCS_SIZE + errorIndex;
        // make sure that the associated handle for each session starts out
        // unassigned
        s_associatedHandles[sessionIndex] = TPM_RH_UNASSIGNED;

        // First parameter: Session handle.
        result = TPMI_SH_AUTH_SESSION_Unmarshal(
            &s_sessionHandles[sessionIndex],
            &command->parameterBuffer,
            &command->authSize, TRUE);
        if(result != TPM_RC_SUCCESS)
            return result + TPM_RC_S + g_rcIndex[sessionIndex];
        // Second parameter: Nonce.
        result = TPM2B_NONCE_Unmarshal(&s_nonceCaller[sessionIndex],
                                       &command->parameterBuffer,
                                       &command->authSize);
        if(result != TPM_RC_SUCCESS)
            return result + TPM_RC_S + g_rcIndex[sessionIndex];
        // Third parameter: sessionAttributes.
        result = TPMA_SESSION_Unmarshal(&s_attributes[sessionIndex],
                                        &command->parameterBuffer,
                                        &command->authSize);
        if(result != TPM_RC_SUCCESS)
            return result + TPM_RC_S + g_rcIndex[sessionIndex];
        // Fourth parameter: authValue (PW or HMAC).
        result = TPM2B_AUTH_Unmarshal(&s_inputAuthValues[sessionIndex],
                                      &command->parameterBuffer,
                                      &command->authSize);
        if(result != TPM_RC_SUCCESS)
            return result + errorIndex;

        sessionAttributes = s_attributes[sessionIndex];
        if(s_sessionHandles[sessionIndex] == TPM_RS_PW)
        {
            // A PWAP session needs additional processing.
            //     Can't have any attributes set other than continueSession bit
            if(IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, encrypt)
               || IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, decrypt)
               || IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, audit)
               || IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, auditExclusive)
               || IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, auditReset))
                return TPM_RCS_ATTRIBUTES + errorIndex;
            //     The nonce size must be zero.
            if(s_nonceCaller[sessionIndex].t.size != 0)
                return TPM_RCS_NONCE + errorIndex;
            continue;
        }
        // For not password sessions...
        // Find out if the session is loaded.
        if(!SessionIsLoaded(s_sessionHandles[sessionIndex]))
            return TPM_RC_REFERENCE_S0 + sessionIndex;
        sessionType = HandleGetType(s_sessionHandles[sessionIndex]);
        session = SessionGet(s_sessionHandles[sessionIndex]);

        // Check if the session is an HMAC/policy session.
        if((session->attributes.isPolicy == SET
            && sessionType == TPM_HT_HMAC_SESSION)
           || (session->attributes.isPolicy == CLEAR
               && sessionType == TPM_HT_POLICY_SESSION))
            return TPM_RCS_HANDLE + errorIndex;
        // Check that this handle has not previously been used.
        for(i = 0; i < sessionIndex; i++)
        {
            if(s_sessionHandles[i] == s_sessionHandles[sessionIndex])
                return TPM_RCS_HANDLE + errorIndex;
        }
        // If the session is used for parameter encryption or audit as well, set
        // the corresponding Indexes.

        // First process decrypt.
        if(IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, decrypt))
        {
            // Check if the commandCode allows command parameter encryption.
            if(DecryptSize(command->index) == 0)
                return TPM_RCS_ATTRIBUTES + errorIndex;
            // Encrypt attribute can only appear in one session
            if(s_decryptSessionIndex != UNDEFINED_INDEX)
                return TPM_RCS_ATTRIBUTES + errorIndex;
            // Can't decrypt if the session's symmetric algorithm is TPM_ALG_NULL
            if(session->symmetric.algorithm == TPM_ALG_NULL)
                return TPM_RCS_SYMMETRIC + errorIndex;
            // All checks passed, so set the index for the session used to decrypt
            // a command parameter.
            s_decryptSessionIndex = sessionIndex;
        }
        // Now process encrypt.
        if(IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, encrypt))
        {
            // Check if the commandCode allows response parameter encryption.
            if(EncryptSize(command->index) == 0)
                return TPM_RCS_ATTRIBUTES + errorIndex;
            // Encrypt attribute can only appear in one session.
            if(s_encryptSessionIndex != UNDEFINED_INDEX)
                return TPM_RCS_ATTRIBUTES + errorIndex;
            // Can't encrypt if the session's symmetric algorithm is TPM_ALG_NULL
            if(session->symmetric.algorithm == TPM_ALG_NULL)
                return TPM_RCS_SYMMETRIC + errorIndex;
            // All checks passed, so set the index for the session used to encrypt
            // a response parameter.
            s_encryptSessionIndex = sessionIndex;
        }
        // At last process audit.
        if(IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, audit))
        {
            // Audit attribute can only appear in one session.
            if(s_auditSessionIndex != UNDEFINED_INDEX)
                return TPM_RCS_ATTRIBUTES + errorIndex;
            // An audit session can not be policy session.
            if(HandleGetType(s_sessionHandles[sessionIndex])
               == TPM_HT_POLICY_SESSION)
                return TPM_RCS_ATTRIBUTES + errorIndex;
            // If this is a reset of the audit session, or the first use
            // of the session as an audit session, it doesn't matter what
            // the exclusive state is. The session will become exclusive.
            if(!IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, auditReset)
               && session->attributes.isAudit == SET)
            {
                // Not first use or reset. If auditExlusive is SET, then this
                // session must be the current exclusive session.
                if(IS_ATTRIBUTE(sessionAttributes, TPMA_SESSION, auditExclusive)
                   && g_exclusiveAuditSession != s_sessionHandles[sessionIndex])
                    return TPM_RC_EXCLUSIVE;
            }
            s_auditSessionIndex = sessionIndex;
        }
        // Initialize associated handle as undefined. This will be changed when
        // the handles are processed.
        s_associatedHandles[sessionIndex] = TPM_RH_UNASSIGNED;
    }
    command->sessionNum = sessionIndex;
    return TPM_RC_SUCCESS;
}

//*** CheckLockedOut()
// This function checks to see if the TPM is in lockout. This function should only
// be called if the entity being checked is subject to DA protection. The TPM
// is in lockout if the NV is not available and a DA write is pending. Otherwise
// the TPM is locked out if checking for lockoutAuth ('lockoutAuthCheck' == TRUE)
// and use of lockoutAuth is disabled, or 'failedTries' >= 'maxTries'
//  Return Type: TPM_RC
//      TPM_RC_NV_RATE          NV is rate limiting
//      TPM_RC_NV_UNAVAILABLE   NV is not available at this time
//      TPM_RC_LOCKOUT          TPM is in lockout
static TPM_RC
CheckLockedOut(
    BOOL             lockoutAuthCheck   // IN: TRUE if checking is for lockoutAuth
    )
{
    // If NV is unavailable, and current cycle state recorded in NV is not
    // SU_NONE_VALUE, refuse to check any authorization because we would
    // not be able to handle a DA failure.
    if(!NV_IS_AVAILABLE && NV_IS_ORDERLY)
        return g_NvStatus;
    // Check if DA info needs to be updated in NV.
    if(s_DAPendingOnNV)
    {
        // If NV is accessible,
        RETURN_IF_NV_IS_NOT_AVAILABLE;

        // ... write the pending DA data and proceed.
        NV_SYNC_PERSISTENT(lockOutAuthEnabled);
        NV_SYNC_PERSISTENT(failedTries);
        s_DAPendingOnNV = FALSE;
    }
    // Lockout is in effect if checking for lockoutAuth and use of lockoutAuth
    // is disabled...
    if(lockoutAuthCheck)
    {
        if(gp.lockOutAuthEnabled == FALSE)
            return TPM_RC_LOCKOUT;
    }
    else
    {
        // ... or if the number of failed tries has been maxed out.
        if(gp.failedTries >= gp.maxTries)
            return TPM_RC_LOCKOUT;
#if USE_DA_USED
        // If the daUsed flag is not SET, then no DA validation until the
        // daUsed state is written to NV
        if(!g_daUsed)
        {
            RETURN_IF_NV_IS_NOT_AVAILABLE;
            g_daUsed = TRUE;
            gp.orderlyState = SU_DA_USED_VALUE;
            NV_SYNC_PERSISTENT(orderlyState);
            return TPM_RC_RETRY;
        }
#endif
    }
    return TPM_RC_SUCCESS;
}

//*** CheckAuthSession()
// This function checks that the authorization session properly authorizes the
// use of the associated handle.
//
//  Return Type: TPM_RC
//      TPM_RC_LOCKOUT              entity is protected by DA and TPM is in
//                                  lockout, or TPM is locked out on NV update
//                                  pending on DA parameters
//
//      TPM_RC_PP                   Physical Presence is required but not provided
//      TPM_RC_AUTH_FAIL            HMAC or PW authorization failed
//                                  with DA side-effects (can be a policy session)
//
//      TPM_RC_BAD_AUTH             HMAC or PW authorization failed without DA
//                                  side-effects (can be a policy session)
//
//      TPM_RC_POLICY_FAIL          if policy session fails
//      TPM_RC_POLICY_CC            command code of policy was wrong
//      TPM_RC_EXPIRED              the policy session has expired
//      TPM_RC_PCR
//      TPM_RC_AUTH_UNAVAILABLE     authValue or authPolicy unavailable
static TPM_RC
CheckAuthSession(
    COMMAND         *command,       // IN: primary parsing structure
    UINT32           sessionIndex   // IN: index of session to be processed
    )
{
    TPM_RC           result = TPM_RC_SUCCESS;
    SESSION         *session = NULL;
    TPM_HANDLE       sessionHandle = s_sessionHandles[sessionIndex];
    TPM_HANDLE       associatedHandle = s_associatedHandles[sessionIndex];
    TPM_HT           sessionHandleType = HandleGetType(sessionHandle);
//
    pAssert(sessionHandle != TPM_RH_UNASSIGNED);

    // Take care of physical presence
    if(associatedHandle == TPM_RH_PLATFORM)
    {
        // If the physical presence is required for this command, check for PP
        // assertion. If it isn't asserted, no point going any further.
        if(PhysicalPresenceIsRequired(command->index)
           && !_plat__PhysicalPresenceAsserted())
            return TPM_RC_PP;
    }
    if(sessionHandle != TPM_RS_PW)
    {
        session = SessionGet(sessionHandle);

        // Set includeAuth to indicate if DA checking will be required and if the
        // authValue will be included in any HMAC.
        if(sessionHandleType == TPM_HT_POLICY_SESSION)
        {
            // For a policy session, will check the DA status of the entity if either
            // isAuthValueNeeded or isPasswordNeeded is SET.
            session->attributes.includeAuth =
                session->attributes.isAuthValueNeeded
                || session->attributes.isPasswordNeeded;
        }
        else
        {
            // For an HMAC session, need to check unless the session
            // is bound.
            session->attributes.includeAuth =
                !IsSessionBindEntity(s_associatedHandles[sessionIndex], session);
        }
    }
    // If the authorization session is going to use an authValue, then make sure
    // that access to that authValue isn't locked out.
    // Note: session == NULL for a PW session.
    if(session == NULL || session->attributes.includeAuth)
    {
        // See if entity is subject to lockout.
        if(!IsDAExempted(associatedHandle))
        {
            // See if in lockout 
            result = CheckLockedOut(associatedHandle == TPM_RH_LOCKOUT);
            if(result != TPM_RC_SUCCESS)
                return result;
        }
    }
    // Policy or HMAC+PW?
    if(sessionHandleType != TPM_HT_POLICY_SESSION)
    {
        // for non-policy session make sure that a policy session is not required
        if(IsPolicySessionRequired(command->index, sessionIndex))
            return TPM_RC_AUTH_TYPE;
        // The authValue must be available.
        // Note: The authValue is going to be "used" even if it is an EmptyAuth.
        // and the session is bound.
        if(!IsAuthValueAvailable(associatedHandle, command->index, sessionIndex))
            return TPM_RC_AUTH_UNAVAILABLE;
    }
    else
    {
        // ... see if the entity has a policy, ...
        // Note: IsAuthPolicyAvalable will return FALSE if the sensitive area of the
        // object is not loaded
        if(!IsAuthPolicyAvailable(associatedHandle, command->index, sessionIndex))
            return TPM_RC_AUTH_UNAVAILABLE;
        // ... and check the policy session.
        result = CheckPolicyAuthSession(command, sessionIndex);
        if(result != TPM_RC_SUCCESS)
            return result;
    }
    // Check authorization according to the type
    if(session == NULL || session->attributes.isPasswordNeeded == SET)
        result = CheckPWAuthSession(sessionIndex);
    else
        result = CheckSessionHMAC(command, sessionIndex);
    // Do processing for PIN Indexes are only three possibilities for 'result' at
    // this point: TPM_RC_SUCCESS, TPM_RC_AUTH_FAIL, and TPM_RC_BAD_AUTH.
    // For all these cases, we would have to process a PIN index if the
    // authValue of the index was used for authorization.
    // See if we need to do anything to a PIN index
    if(TPM_HT_NV_INDEX == HandleGetType(associatedHandle))
    {
        NV_REF           locator;
        NV_INDEX        *nvIndex = NvGetIndexInfo(associatedHandle, &locator);
        NV_PIN           pinData;
        TPMA_NV          nvAttributes;
//
        pAssert(nvIndex != NULL);
        nvAttributes = nvIndex->publicArea.attributes;
        // If this is a PIN FAIL index and the value has been written
        // then we can update the counter (increment or clear)
        if(IsNvPinFailIndex(nvAttributes) 
           && IS_ATTRIBUTE(nvAttributes, TPMA_NV, WRITTEN))
        {
            pinData.intVal = NvGetUINT64Data(nvIndex, locator);
            if(result != TPM_RC_SUCCESS)
                pinData.pin.pinCount++;
            else
                pinData.pin.pinCount = 0;
            NvWriteUINT64Data(nvIndex, pinData.intVal);
        }
        // If this is a PIN PASS Index, increment if we have used the
        // authorization value for anything other than NV_Read.
        // NOTE: If the counter has already hit the limit, then we
        // would not get here because the authorization value would not
        // be available and the TPM would have returned before it gets here
        else if(IsNvPinPassIndex(nvAttributes)
                && IS_ATTRIBUTE(nvAttributes, TPMA_NV, WRITTEN)
                && result == TPM_RC_SUCCESS)
        {
            // If the access is valid, then increment the use counter
            pinData.intVal = NvGetUINT64Data(nvIndex, locator);
            pinData.pin.pinCount++;
            NvWriteUINT64Data(nvIndex, pinData.intVal);
        }
    }
    return result;
}

#ifdef  TPM_CC_GetCommandAuditDigest
//*** CheckCommandAudit()
// This function is called before the command is processed if audit is enabled
// for the command. It will check to see if the audit can be performed and
// will ensure that the cpHash is available for the audit.
//  Return Type: TPM_RC
//      TPM_RC_NV_UNAVAILABLE       NV is not available for write
//      TPM_RC_NV_RATE              NV is rate limiting
static TPM_RC
CheckCommandAudit(
    COMMAND         *command
    )
{
    // If the audit digest is clear and command audit is required, NV must be
    // available so that TPM2_GetCommandAuditDigest() is able to increment
    // audit counter. If NV is not available, the function bails out to prevent
    // the TPM from attempting an operation that would fail anyway.
    if(gr.commandAuditDigest.t.size == 0
       || GetCommandCode(command->index) == TPM_CC_GetCommandAuditDigest)
    {
        RETURN_IF_NV_IS_NOT_AVAILABLE;
    }
    // Make sure that the cpHash is computed for the algorithm
    ComputeCpHash(command, gp.auditHashAlg);
    return TPM_RC_SUCCESS;
}
#endif

//*** ParseSessionBuffer()
// This function is the entry function for command session processing.
// It iterates sessions in session area and reports if the required authorization
// has been properly provided. It also processes audit session and passes the
// information of encryption sessions to parameter encryption module.
//
//  Return Type: TPM_RC
//        various           parsing failure or authorization failure
//
TPM_RC
ParseSessionBuffer(
    COMMAND         *command        // IN: the structure that contains
    )
{
    TPM_RC               result;
    UINT32               i;
    INT32                size = 0;
    TPM2B_AUTH           extraKey;
    UINT32               sessionIndex;
    TPM_RC               errorIndex;
    SESSION             *session = NULL;
//
    // Check if a command allows any session in its session area.
    if(!IsSessionAllowed(command->index))
        return TPM_RC_AUTH_CONTEXT;
    // Default-initialization.
    command->sessionNum = 0;

    result = RetrieveSessionData(command);
    if(result != TPM_RC_SUCCESS)
        return result;
    // There is no command in the TPM spec that has more handles than
    // MAX_SESSION_NUM.
    pAssert(command->handleNum <= MAX_SESSION_NUM);

    // Associate the session with an authorization handle.
    for(i = 0; i < command->handleNum; i++)
    {
        if(CommandAuthRole(command->index, i) != AUTH_NONE)
        {
            // If the received session number is less than the number of handles
            // that requires authorization, an error should be returned.
            // Note: for all the TPM 2.0 commands, handles requiring
            // authorization come first in a command input and there are only ever
            // two values requiring authorization
            if(i > (command->sessionNum - 1))
                return TPM_RC_AUTH_MISSING;
            // Record the handle associated with the authorization session
            s_associatedHandles[i] = command->handles[i];
        }
    }
    // Consistency checks are done first to avoid authorization failure when the
    // command will not be executed anyway.
    for(sessionIndex = 0; sessionIndex < command->sessionNum; sessionIndex++)
    {
        errorIndex = TPM_RC_S + g_rcIndex[sessionIndex];
        // PW session must be an authorization session
        if(s_sessionHandles[sessionIndex] == TPM_RS_PW)
        {
            if(s_associatedHandles[sessionIndex] == TPM_RH_UNASSIGNED)
                return TPM_RCS_HANDLE + errorIndex;
            // a password session can't be audit, encrypt or decrypt
            if(IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, audit)
               || IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, encrypt)
               || IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, decrypt))
                return TPM_RCS_ATTRIBUTES + errorIndex;
            session = NULL;
        }
        else
        {
            session = SessionGet(s_sessionHandles[sessionIndex]);

            // A trial session can not appear in session area, because it cannot
            // be used for authorization, audit or encrypt/decrypt.
            if(session->attributes.isTrialPolicy == SET)
                return TPM_RCS_ATTRIBUTES + errorIndex;

            // See if the session is bound to a DA protected entity
            // NOTE: Since a policy session is never bound, a policy is still
            // usable even if the object is DA protected and the TPM is in
            // lockout.
            if(session->attributes.isDaBound == SET)
            {
                result = CheckLockedOut(session->attributes.isLockoutBound == SET);
                if(result != TPM_RC_SUCCESS)
                    return result;
            }
            // If this session is for auditing, make sure the cpHash is computed.
            if(IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, audit))
                ComputeCpHash(command, session->authHashAlg);
        }

        // if the session has an associated handle, check the authorization
        if(s_associatedHandles[sessionIndex] != TPM_RH_UNASSIGNED)
        {
            result = CheckAuthSession(command, sessionIndex);
            if(result != TPM_RC_SUCCESS)
                return RcSafeAddToResult(result, errorIndex);
        }
        else
        {
            // a session that is not for authorization must either be encrypt,
            // decrypt, or audit
            if(!IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, audit)
               &&  !IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, encrypt)
               &&  !IS_ATTRIBUTE(s_attributes[sessionIndex], TPMA_SESSION, decrypt))
                return TPM_RCS_ATTRIBUTES + errorIndex;

            // no authValue included in any of the HMAC computations
            pAssert(session != NULL);
            session->attributes.includeAuth = CLEAR;

            // check HMAC for encrypt/decrypt/audit only sessions
            result = CheckSessionHMAC(command, sessionIndex);
            if(result != TPM_RC_SUCCESS)
                return RcSafeAddToResult(result, errorIndex);
        }
    }
#ifdef  TPM_CC_GetCommandAuditDigest
    // Check if the command should be audited. Need to do this before any parameter
    // encryption so that the cpHash for the audit is correct
    if(CommandAuditIsRequired(command->index))
    {
        result = CheckCommandAudit(command);
        if(result != TPM_RC_SUCCESS)
            return result;              // No session number to reference
    }
#endif
    // Decrypt the first parameter if applicable. This should be the last operation
    // in session processing.
    // If the encrypt session is associated with a handle and the handle's
    // authValue is available, then authValue is concatenated with sessionKey to
    // generate encryption key, no matter if the handle is the session bound entity
    // or not.
    if(s_decryptSessionIndex != UNDEFINED_INDEX)
    {
        // If this is an authorization session, include the authValue in the
        // generation of the decryption key
        if(s_associatedHandles[s_decryptSessionIndex] != TPM_RH_UNASSIGNED)
        {
            EntityGetAuthValue(s_associatedHandles[s_decryptSessionIndex], 
                               &extraKey);
        }
        else
        {
            extraKey.b.size = 0;
        }
        size = DecryptSize(command->index);
        result = CryptParameterDecryption(s_sessionHandles[s_decryptSessionIndex],
                                          &s_nonceCaller[s_decryptSessionIndex].b,
                                          command->parameterSize, (UINT16)size,
                                          &extraKey,
                                          command->parameterBuffer);
        if(result != TPM_RC_SUCCESS)
            return RcSafeAddToResult(result,
                                     TPM_RC_S + g_rcIndex[s_decryptSessionIndex]);
    }

    return TPM_RC_SUCCESS;
}

//*** CheckAuthNoSession()
// Function to process a command with no session associated.
// The function makes sure all the handles in the command require no authorization.
//
//  Return Type: TPM_RC
//      TPM_RC_AUTH_MISSING         failure - one or more handles require
//                                  authorization
TPM_RC
CheckAuthNoSession(
    COMMAND         *command        // IN: command parsing structure
    )
{
    UINT32 i;
    TPM_RC           result = TPM_RC_SUCCESS;
//
    // Check if the command requires authorization
    for(i = 0; i < command->handleNum; i++)
    {
        if(CommandAuthRole(command->index, i) != AUTH_NONE)
            return TPM_RC_AUTH_MISSING;
    }
#ifdef  TPM_CC_GetCommandAuditDigest
    // Check if the command should be audited.
    if(CommandAuditIsRequired(command->index))
    {
        result = CheckCommandAudit(command);
        if(result != TPM_RC_SUCCESS)
            return result;
    }
#endif
    // Initialize number of sessions to be 0
    command->sessionNum = 0;

    return TPM_RC_SUCCESS;
}

//** Response Session Processing
//*** Introduction
//
//  The following functions build the session area in a response and handle
//  the audit sessions (if present).
//

//*** ComputeRpHash()
// Function to compute rpHash (Response Parameter Hash). The rpHash is only
// computed if there is an HMAC authorization session and the return code is
// TPM_RC_SUCCESS.
static TPM2B_DIGEST *
ComputeRpHash(
    COMMAND         *command,       // IN: command structure
    TPM_ALG_ID       hashAlg        // IN: hash algorithm to compute rpHash
    )
{
    TPM2B_DIGEST    *rpHash = GetRpHashPointer(command, hashAlg);
    HASH_STATE       hashState;
//
    if(rpHash->t.size == 0)
    {
    //   rpHash := hash(responseCode || commandCode || parameters)

    // Initiate hash creation.
        rpHash->t.size = CryptHashStart(&hashState, hashAlg);

        // Add hash constituents.
        CryptDigestUpdateInt(&hashState, sizeof(TPM_RC), TPM_RC_SUCCESS);
        CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), command->code);
        CryptDigestUpdate(&hashState, command->parameterSize, 
                          command->parameterBuffer);
        // Complete hash computation.
        CryptHashEnd2B(&hashState, &rpHash->b);
    }
    return rpHash;
}

//*** InitAuditSession()
// This function initializes the audit data in an audit session.
static void
InitAuditSession(
    SESSION         *session        // session to be initialized
    )
{
    // Mark session as an audit session.
    session->attributes.isAudit = SET;

    // Audit session can not be bound.
    session->attributes.isBound = CLEAR;

    // Size of the audit log is the size of session hash algorithm digest.
    session->u2.auditDigest.t.size = CryptHashGetDigestSize(session->authHashAlg);

    // Set the original digest value to be 0.
    MemorySet(&session->u2.auditDigest.t.buffer,
              0,
              session->u2.auditDigest.t.size);
    return;
}

//*** UpdateAuditDigest
// Function to update an audit digest
static void
UpdateAuditDigest(
    COMMAND         *command,
    TPMI_ALG_HASH    hashAlg,
    TPM2B_DIGEST    *digest
    )
{
    HASH_STATE       hashState;
    TPM2B_DIGEST    *cpHash = GetCpHash(command, hashAlg);
    TPM2B_DIGEST    *rpHash = ComputeRpHash(command, hashAlg);
//
    pAssert(cpHash != NULL);

    // digestNew :=  hash (digestOld || cpHash || rpHash)
        // Start hash computation.
    digest->t.size = CryptHashStart(&hashState, hashAlg);
        // Add old digest.
    CryptDigestUpdate2B(&hashState, &digest->b);
        // Add cpHash 
    CryptDigestUpdate2B(&hashState, &cpHash->b);
        // Add rpHash
    CryptDigestUpdate2B(&hashState, &rpHash->b);
        // Finalize the hash.
    CryptHashEnd2B(&hashState, &digest->b);
}


//*** Audit()
//This function updates the audit digest in an audit session.
static void
Audit(
    COMMAND         *command,       // IN: primary control structure
    SESSION         *auditSession   // IN: loaded audit session
    )
{
    UpdateAuditDigest(command, auditSession->authHashAlg,
                      &auditSession->u2.auditDigest);
    return;
}

#ifdef  TPM_CC_GetCommandAuditDigest
//*** CommandAudit()
// This function updates the command audit digest.
static void
CommandAudit(
    COMMAND         *command        // IN:
    )
{
    // If the digest.size is one, it indicates the special case of changing
    // the audit hash algorithm. For this case, no audit is done on exit.
    // NOTE: When the hash algorithm is changed, g_updateNV is set in order to
    // force an update to the NV on exit so that the change in digest will
    // be recorded. So, it is safe to exit here without setting any flags
    // because the digest change will be written to NV when this code exits.
    if(gr.commandAuditDigest.t.size == 1)
    {
        gr.commandAuditDigest.t.size = 0;
        return;
    }
    // If the digest size is zero, need to start a new digest and increment
    // the audit counter.
    if(gr.commandAuditDigest.t.size == 0)
    {
        gr.commandAuditDigest.t.size = CryptHashGetDigestSize(gp.auditHashAlg);
        MemorySet(gr.commandAuditDigest.t.buffer,
                  0,
                  gr.commandAuditDigest.t.size);

        // Bump the counter and save its value to NV.
        gp.auditCounter++;
        NV_SYNC_PERSISTENT(auditCounter);
    }
    UpdateAuditDigest(command, gp.auditHashAlg, &gr.commandAuditDigest);
    return;
}
#endif

//*** UpdateAuditSessionStatus()
// Function to update the internal audit related states of a session. It
//  1. initializes the session as audit session and sets it to be exclusive if this
//     is the first time it is used for audit or audit reset was requested;
//  2. reports exclusive audit session;
//  3. extends audit log; and
//  4. clears exclusive audit session if no audit session found in the command.
static void
UpdateAuditSessionStatus(
    COMMAND         *command        // IN: primary control structure
    )
{
    UINT32           i;
    TPM_HANDLE       auditSession = TPM_RH_UNASSIGNED;
//
    // Iterate through sessions
    for(i = 0; i < command->sessionNum; i++)
    {
        SESSION     *session;
//
        // PW session do not have a loaded session and can not be an audit
        // session either.  Skip it.
        if(s_sessionHandles[i] == TPM_RS_PW)
            continue;
        session = SessionGet(s_sessionHandles[i]);

        // If a session is used for audit
        if(IS_ATTRIBUTE(s_attributes[i], TPMA_SESSION, audit))
        {
            // An audit session has been found
            auditSession = s_sessionHandles[i];

            // If the session has not been an audit session yet, or
            // the auditSetting bits indicate a reset, initialize it and set
            // it to be the exclusive session
            if(session->attributes.isAudit == CLEAR
               || IS_ATTRIBUTE(s_attributes[i], TPMA_SESSION, auditReset))
            {
                InitAuditSession(session);
                g_exclusiveAuditSession = auditSession;
            }
            else
            {
                // Check if the audit session is the current exclusive audit
                // session and, if not, clear previous exclusive audit session.
                if(g_exclusiveAuditSession != auditSession)
                    g_exclusiveAuditSession = TPM_RH_UNASSIGNED;
            }
            // Report audit session exclusivity.
            if(g_exclusiveAuditSession == auditSession)
            {
                SET_ATTRIBUTE(s_attributes[i], TPMA_SESSION, auditExclusive);
            }
            else
            {
                CLEAR_ATTRIBUTE(s_attributes[i], TPMA_SESSION, auditExclusive);
            }
            // Extend audit log.
            Audit(command, session);
        }
    }
    // If no audit session is found in the command, and the command allows
    // a session then, clear the current exclusive
    // audit session.
    if(auditSession == TPM_RH_UNASSIGNED && IsSessionAllowed(command->index))
    {
        g_exclusiveAuditSession = TPM_RH_UNASSIGNED;
    }
    return;
}

//*** ComputeResponseHMAC()
// Function to compute HMAC for authorization session in a response.
/*(See part 1 specification)
// Function: Compute HMAC for response sessions
//      The sessionAuth value
//          authHMAC := HMACsHASH((sessionAuth | authValue),
//                    (pHash | nonceTPM | nonceCaller | sessionAttributes))
//  Where:
//      HMACsHASH()     The HMAC algorithm using the hash algorithm specified when
//                      the session was started.
//
//      sessionAuth     A TPMB_MEDIUM computed in a protocol-dependent way, using
//                      KDFa. In an HMAC or KDF, only sessionAuth.buffer is used.
//
//      authValue       A TPM2B_AUTH that is found in the sensitive area of an
//                      object. In an HMAC or KDF, only authValue.buffer is used
//                      and all trailing zeros are removed.
//
//      pHash           Response parameters (rpHash) using the session hash. When
//                      using a pHash in an HMAC computation, both the algorithm ID
//                      and the digest are included.
//
//      nonceTPM        A TPM2B_NONCE that is generated by the entity using the
//                      session. In an HMAC or KDF, only nonceTPM.buffer is used.
//
//      nonceCaller     a TPM2B_NONCE that was received the previous time the
//                      session was used. In an HMAC or KDF, only
//                      nonceCaller.buffer is used.
//
//      sessionAttributes   A TPMA_SESSION that indicates the attributes associated
//                          with a particular use of the session.
*/
static void
ComputeResponseHMAC(
    COMMAND         *command,       // IN: command structure
    UINT32           sessionIndex,  // IN: session index to be processed
    SESSION         *session,       // IN: loaded session
    TPM2B_DIGEST    *hmac           // OUT: authHMAC
    )
{
    TPM2B_TYPE(KEY, (sizeof(AUTH_VALUE) * 2));
    TPM2B_KEY        key;       // HMAC key
    BYTE             marshalBuffer[sizeof(TPMA_SESSION)];
    BYTE            *buffer;
    UINT32           marshalSize;
    HMAC_STATE       hmacState;
    TPM2B_DIGEST    *rpHash = ComputeRpHash(command, session->authHashAlg);
//
    // Generate HMAC key
    MemoryCopy2B(&key.b, &session->sessionKey.b, sizeof(key.t.buffer));

    // Add the object authValue if required
    if(session->attributes.includeAuth == SET)
    {
        // Note: includeAuth may be SET for a policy that is used in
        // UndefineSpaceSpecial(). At this point, the Index has been deleted
        // so the includeAuth will have no meaning. However, the
        // s_associatedHandles[] value for the session is now set to TPM_RH_NULL so
        // this will return the authValue associated with TPM_RH_NULL and that is
        // and empty buffer.
        TPM2B_AUTH          authValue;
//
        // Get the authValue with trailing zeros removed
        EntityGetAuthValue(s_associatedHandles[sessionIndex], &authValue);

        // Add it to the key
        MemoryConcat2B(&key.b, &authValue.b, sizeof(key.t.buffer));
    }

    // if the HMAC key size is 0, the response HMAC is computed according to the
    // input HMAC
    if(key.t.size == 0
       && s_inputAuthValues[sessionIndex].t.size == 0)
    {
        hmac->t.size = 0;
        return;
    }
    // Start HMAC computation.
    hmac->t.size = CryptHmacStart2B(&hmacState, session->authHashAlg, &key.b);

    // Add hash components.
    CryptDigestUpdate2B(&hmacState.hashState, &rpHash->b);
    CryptDigestUpdate2B(&hmacState.hashState, &session->nonceTPM.b);
    CryptDigestUpdate2B(&hmacState.hashState, &s_nonceCaller[sessionIndex].b);

    // Add session attributes.
    buffer = marshalBuffer;
    marshalSize = TPMA_SESSION_Marshal(&s_attributes[sessionIndex], &buffer, NULL);
    CryptDigestUpdate(&hmacState.hashState, marshalSize, marshalBuffer);

    // Finalize HMAC.
    CryptHmacEnd2B(&hmacState, &hmac->b);

    return;
}

//*** UpdateInternalSession()
// Updates internal sessions:
//      1. Restarts session time.
//      2. Clears a policy session since nonce is rolling.
static void
UpdateInternalSession(
    SESSION         *session,       // IN: the session structure
    UINT32           i              // IN: session number
    )
{
    // If nonce is rolling in a policy session, the policy related data
    // will be re-initialized.
    if(HandleGetType(s_sessionHandles[i]) == TPM_HT_POLICY_SESSION
       && IS_ATTRIBUTE(s_attributes[i], TPMA_SESSION, continueSession))
    {
        // When the nonce rolls it starts a new timing interval for the
        // policy session.
        SessionResetPolicyData(session);
        SessionSetStartTime(session);
    }
    return;
}

//*** BuildSingleResponseAuth()
//   Function to compute response HMAC value for a policy or HMAC session.
static TPM2B_NONCE *
BuildSingleResponseAuth(
    COMMAND         *command,       // IN: command structure
    UINT32           sessionIndex,  // IN: session index to be processed
    TPM2B_AUTH      *auth           // OUT: authHMAC
    )
{
    // Fill in policy/HMAC based session response.
    SESSION     *session = SessionGet(s_sessionHandles[sessionIndex]);
//
    // If the session is a policy session with isPasswordNeeded SET, the
    // authorization field is empty.
    if(HandleGetType(s_sessionHandles[sessionIndex]) == TPM_HT_POLICY_SESSION
       && session->attributes.isPasswordNeeded == SET)
        auth->t.size = 0;
    else
        // Compute response HMAC.
        ComputeResponseHMAC(command, sessionIndex, session, auth);

    UpdateInternalSession(session, sessionIndex);
    return &session->nonceTPM;
}

//*** UpdateAllNonceTPM()
// Updates TPM nonce for all sessions in command.
static void
UpdateAllNonceTPM(
    COMMAND         *command        // IN: controlling structure
    )
{
    UINT32      i;
    SESSION     *session;
//
    for(i = 0; i < command->sessionNum; i++)
    {
        // If not a PW session, compute the new nonceTPM.
        if(s_sessionHandles[i] != TPM_RS_PW)
        {
            session = SessionGet(s_sessionHandles[i]);
            // Update nonceTPM in both internal session and response.
            CryptRandomGenerate(session->nonceTPM.t.size,
                                session->nonceTPM.t.buffer);
        }
    }
    return;
}



//*** BuildResponseSession()
// Function to build Session buffer in a response. The authorization data is added
// to the end of command->responseBuffer. The size of the authorization area is
// accumulated in command->authSize.
// When this is called, command->responseBuffer is pointing at the next location
// in the response buffer to be filled. This is where the authorization sessions 
// will go, if any. command->parameterSize is the number of bytes that have been
// marshaled as parameters in the output buffer.
void
BuildResponseSession(
    COMMAND         *command        // IN: structure that has relevant command
                                    //     information
    )
{
    pAssert(command->authSize == 0);

    // Reset the parameter buffer to point to the start of the parameters so that
    // there is a starting point for any rpHash that might be generated and so there
    // is a place where parameter encryption would start
    command->parameterBuffer = command->responseBuffer - command->parameterSize;

    // Session nonces should be updated before parameter encryption
    if(command->tag == TPM_ST_SESSIONS)
    {
        UpdateAllNonceTPM(command);

        // Encrypt first parameter if applicable. Parameter encryption should
        // happen after nonce update and before any rpHash is computed.
        // If the encrypt session is associated with a handle, the authValue of
        // this handle will be concatenated with sessionKey to generate
        // encryption key, no matter if the handle is the session bound entity
        // or not. The authValue is added to sessionKey only when the authValue
        // is available.
        if(s_encryptSessionIndex != UNDEFINED_INDEX)
        {
            UINT32          size;
            TPM2B_AUTH      extraKey;
//
            extraKey.b.size = 0;
            // If this is an authorization session, include the authValue in the
            // generation of the encryption key
            if(s_associatedHandles[s_encryptSessionIndex] != TPM_RH_UNASSIGNED)
            {
                EntityGetAuthValue(s_associatedHandles[s_encryptSessionIndex], 
                                   &extraKey);
            }
            size = EncryptSize(command->index);
            CryptParameterEncryption(s_sessionHandles[s_encryptSessionIndex],
                                     &s_nonceCaller[s_encryptSessionIndex].b,
                                     (UINT16)size,
                                     &extraKey,
                                     command->parameterBuffer);
        }
    }
    // Audit sessions should be processed regardless of the tag because
    // a command with no session may cause a change of the exclusivity state.
    UpdateAuditSessionStatus(command);
#if CC_GetCommandAuditDigest
    // Command Audit
    if(CommandAuditIsRequired(command->index))
        CommandAudit(command);
#endif
    // Process command with sessions.
    if(command->tag == TPM_ST_SESSIONS)
    {
        UINT32           i;
//
        pAssert(command->sessionNum > 0);

        // Iterate over each session in the command session area, and create
        // corresponding sessions for response.
        for(i = 0; i < command->sessionNum; i++)
        {
            TPM2B_NONCE     *nonceTPM;
            TPM2B_DIGEST     responseAuth;
            // Make sure that continueSession is SET on any Password session.
            // This makes it marginally easier for the management software
            // to keep track of the closed sessions.
            if(s_sessionHandles[i] == TPM_RS_PW)
            {
                SET_ATTRIBUTE(s_attributes[i], TPMA_SESSION, continueSession);
                responseAuth.t.size = 0;
                nonceTPM = (TPM2B_NONCE *)&responseAuth;                
            }
            else
            {
                // Compute the response HMAC and get a pointer to the nonce used.
                // This function will also update the values if needed. Note, the
                nonceTPM = BuildSingleResponseAuth(command, i, &responseAuth);
            }
            command->authSize += TPM2B_NONCE_Marshal(nonceTPM,
                                                     &command->responseBuffer,
                                                     NULL);
            command->authSize += TPMA_SESSION_Marshal(&s_attributes[i],
                                                      &command->responseBuffer,
                                                      NULL);
            command->authSize += TPM2B_DIGEST_Marshal(&responseAuth,
                                                      &command->responseBuffer,
                                                      NULL);
            if(!IS_ATTRIBUTE(s_attributes[i], TPMA_SESSION, continueSession))
                SessionFlush(s_sessionHandles[i]);
        }
    }
    return;
}

//*** SessionRemoveAssociationToHandle()
// This function deals with the case where an entity associated with an authorization
// is deleted during command processing. The primary use of this is to support
// UndefineSpaceSpecial().
void
SessionRemoveAssociationToHandle(
    TPM_HANDLE       handle
    )
{
    UINT32               i;
//
    for(i = 0; i < MAX_SESSION_NUM; i++)
    {
        if(s_associatedHandles[i] == handle)
        {
            s_associatedHandles[i] = TPM_RH_NULL;
        }
    }
}