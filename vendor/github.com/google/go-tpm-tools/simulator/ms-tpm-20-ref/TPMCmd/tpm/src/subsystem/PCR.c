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
//
// This function contains the functions needed for PCR access and manipulation.
//
// This implementation uses a static allocation for the PCR. The amount of
// memory is allocated based on the number of PCR in the implementation and
// the number of implemented hash algorithms. This is not the expected
// implementation. PCR SPACE DEFINITIONS.
//
// In the definitions below, the g_hashPcrMap is a bit array that indicates
// which of the PCR are implemented. The g_hashPcr array is an array of digests.
// In this implementation, the space is allocated whether the PCR is implemented
// or not.

//** Includes, Defines, and Data Definitions
#define PCR_C
#include "Tpm.h"

// The initial value of PCR attributes.  The value of these fields should be
// consistent with PC Client specification
// In this implementation, we assume the total number of implemented PCR is 24.
static const PCR_Attributes s_initAttributes[] =
{
    // PCR 0 - 15, static RTM
    {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F},
    {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F},
    {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F},
    {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F}, {1, 0, 0x1F},

    {0, 0x0F, 0x1F},        // PCR 16, Debug
    {0, 0x10, 0x1C},        // PCR 17, Locality 4
    {0, 0x10, 0x1C},        // PCR 18, Locality 3
    {0, 0x10, 0x0C},        // PCR 19, Locality 2
    {0, 0x14, 0x0E},        // PCR 20, Locality 1
    {0, 0x14, 0x04},        // PCR 21, Dynamic OS
    {0, 0x14, 0x04},        // PCR 22, Dynamic OS
    {0, 0x0F, 0x1F},        // PCR 23, Application specific
    {0, 0x0F, 0x1F}         // PCR 24, testing policy
};

//** Functions

//*** PCRBelongsAuthGroup()
// This function indicates if a PCR belongs to a group that requires an authValue
// in order to modify the PCR.  If it does, 'groupIndex' is set to value of
// the group index.  This feature of PCR is decided by the platform specification.
//  Return Type: BOOL
//      TRUE(1)         PCR belongs an authorization group
//      FALSE(0)        PCR does not belong an authorization group
BOOL
PCRBelongsAuthGroup(
    TPMI_DH_PCR      handle,        // IN: handle of PCR
    UINT32          *groupIndex     // OUT: group index if PCR belongs a
                                    //      group that allows authValue.  If PCR
                                    //      does not belong to an authorization
                                    //      group, the value in this parameter is
                                    //      invalid
    )
{
#if defined NUM_AUTHVALUE_PCR_GROUP && NUM_AUTHVALUE_PCR_GROUP > 0
    // Platform specification determines to which authorization group a PCR belongs
    // (if any). In this implementation, we assume there is only
    // one authorization group which contains PCR[20-22].  If the platform
    // specification requires differently, the implementation should be changed
    // accordingly
    if(handle >= 20 && handle <= 22)
    {
        *groupIndex = 0;
        return TRUE;
    }

#endif
    return FALSE;
}

//*** PCRBelongsPolicyGroup()
// This function indicates if a PCR belongs to a group that requires a policy
// authorization in order to modify the PCR.  If it does, 'groupIndex' is set
// to value of the group index.  This feature of PCR is decided by the platform
// specification.
//  Return Type: BOOL
//      TRUE(1)         PCR belongs a policy group
//      FALSE(0)        PCR does not belong a policy group
BOOL
PCRBelongsPolicyGroup(
    TPMI_DH_PCR      handle,        // IN: handle of PCR
    UINT32          *groupIndex     // OUT: group index if PCR belongs a group that
                                    //     allows policy.  If PCR does not belong to
                                    //     a policy group, the value in this
                                    //     parameter is invalid
    )
{
#if defined NUM_POLICY_PCR_GROUP && NUM_POLICY_PCR_GROUP > 0
    // Platform specification decides if a PCR belongs to a policy group and
    // belongs to which group.  In this implementation, we assume there is only
    // one policy group which contains PCR20-22.  If the platform specification
    // requires differently, the implementation should be changed accordingly
    if(handle >= 20 && handle <= 22)
    {
        *groupIndex = 0;
        return TRUE;
    }
#endif
    return FALSE;
}

//*** PCRBelongsTCBGroup()
// This function indicates if a PCR belongs to the TCB group.
//  Return Type: BOOL
//      TRUE(1)         PCR belongs to TCB group
//      FALSE(0)        PCR does not belong to TCB group
static BOOL
PCRBelongsTCBGroup(
    TPMI_DH_PCR      handle         // IN: handle of PCR
    )
{
#if ENABLE_PCR_NO_INCREMENT == YES
    // Platform specification decides if a PCR belongs to a TCB group.  In this
    // implementation, we assume PCR[20-22] belong to TCB group.  If the platform
    // specification requires differently, the implementation should be
    // changed accordingly
    if(handle >= 20 && handle <= 22)
        return TRUE;

#endif
    return FALSE;
}

//*** PCRPolicyIsAvailable()
// This function indicates if a policy is available for a PCR.
//  Return Type: BOOL
//      TRUE(1)         the PCR should be authorized by policy
//      FALSE(0)        the PCR does not allow policy
BOOL
PCRPolicyIsAvailable(
    TPMI_DH_PCR      handle         // IN: PCR handle
    )
{
    UINT32          groupIndex;

    return PCRBelongsPolicyGroup(handle, &groupIndex);
}

//*** PCRGetAuthValue()
// This function is used to access the authValue of a PCR.  If PCR does not
// belong to an authValue group, an EmptyAuth will be returned.
TPM2B_AUTH *
PCRGetAuthValue(
    TPMI_DH_PCR      handle         // IN: PCR handle
    )
{
    UINT32      groupIndex;

    if(PCRBelongsAuthGroup(handle, &groupIndex))
    {
        return &gc.pcrAuthValues.auth[groupIndex];
    }
    else
    {
        return NULL;
    }
}

//*** PCRGetAuthPolicy()
// This function is used to access the authorization policy of a PCR. It sets
// 'policy' to the authorization policy and returns the hash algorithm for policy
//  If the PCR does not allow a policy, TPM_ALG_NULL is returned.
TPMI_ALG_HASH
PCRGetAuthPolicy(
    TPMI_DH_PCR      handle,        // IN: PCR handle
    TPM2B_DIGEST    *policy         // OUT: policy of PCR
    )
{
    UINT32           groupIndex;

    if(PCRBelongsPolicyGroup(handle, &groupIndex))
    {
        *policy = gp.pcrPolicies.policy[groupIndex];
        return gp.pcrPolicies.hashAlg[groupIndex];
    }
    else
    {
        policy->t.size = 0;
        return TPM_ALG_NULL;
    }
}

//*** PCRSimStart()
// This function is used to initialize the policies when a TPM is manufactured.
// This function would only be called in a manufacturing environment or in
// a TPM simulator.
void
PCRSimStart(
    void
    )
{
    UINT32  i;
#if defined NUM_POLICY_PCR_GROUP && NUM_POLICY_PCR_GROUP > 0
    for(i = 0; i < NUM_POLICY_PCR_GROUP; i++)
    {
        gp.pcrPolicies.hashAlg[i] = TPM_ALG_NULL;
        gp.pcrPolicies.policy[i].t.size = 0;
    }
#endif
#if defined NUM_AUTHVALUE_PCR_GROUP && NUM_AUTHVALUE_PCR_GROUP > 0
    for(i = 0; i < NUM_AUTHVALUE_PCR_GROUP; i++)
    {
        gc.pcrAuthValues.auth[i].t.size = 0;
    }
#endif
    // We need to give an initial configuration on allocated PCR before
    // receiving any TPM2_PCR_Allocate command to change this configuration
    // When the simulation environment starts, we allocate all the PCRs
    for(gp.pcrAllocated.count = 0; gp.pcrAllocated.count < HASH_COUNT;
    gp.pcrAllocated.count++)
    {
        gp.pcrAllocated.pcrSelections[gp.pcrAllocated.count].hash
            = CryptHashGetAlgByIndex(gp.pcrAllocated.count);

        gp.pcrAllocated.pcrSelections[gp.pcrAllocated.count].sizeofSelect
            = PCR_SELECT_MAX;
        for(i = 0; i < PCR_SELECT_MAX; i++)
            gp.pcrAllocated.pcrSelections[gp.pcrAllocated.count].pcrSelect[i]
            = 0xFF;
    }

    // Store the initial configuration to NV
    NV_SYNC_PERSISTENT(pcrPolicies);
    NV_SYNC_PERSISTENT(pcrAllocated);

    return;
}

//*** GetSavedPcrPointer()
// This function returns the address of an array of state saved PCR based
// on the hash algorithm.
//  Return Type: BYTE *
//      NULL            no such algorithm
//      != NULL         pointer to the 0th byte of the 0th PCR
static BYTE *
GetSavedPcrPointer(
    TPM_ALG_ID       alg,           // IN: algorithm for bank
    UINT32           pcrIndex       // IN: PCR index in PCR_SAVE
    )
{
    BYTE            *retVal;
    switch(alg)
    {
#if ALG_SHA1
        case ALG_SHA1_VALUE:
            retVal = gc.pcrSave.sha1[pcrIndex];
            break;
#endif
#if ALG_SHA256
        case ALG_SHA256_VALUE:
            retVal = gc.pcrSave.sha256[pcrIndex];
            break;
#endif
#if ALG_SHA384
        case ALG_SHA384_VALUE:
            retVal = gc.pcrSave.sha384[pcrIndex];
            break;
#endif

#if ALG_SHA512
        case ALG_SHA512_VALUE:
            retVal = gc.pcrSave.sha512[pcrIndex];
            break;
#endif
#if ALG_SM3_256
        case ALG_SM3_256_VALUE:
            retVal = gc.pcrSave.sm3_256[pcrIndex];
            break;
#endif
        default:
            FAIL(FATAL_ERROR_INTERNAL);
    }
    return retVal;
}

//*** PcrIsAllocated()
// This function indicates if a PCR number for the particular hash algorithm
// is allocated.
//  Return Type: BOOL
//      TRUE(1)         PCR is allocated
//      FALSE(0)        PCR is not allocated
BOOL
PcrIsAllocated(
    UINT32           pcr,           // IN: The number of the PCR
    TPMI_ALG_HASH    hashAlg        // IN: The PCR algorithm
    )
{
    UINT32               i;
    BOOL                 allocated = FALSE;

    if(pcr < IMPLEMENTATION_PCR)
    {
        for(i = 0; i < gp.pcrAllocated.count; i++)
        {
            if(gp.pcrAllocated.pcrSelections[i].hash == hashAlg)
            {
                if(((gp.pcrAllocated.pcrSelections[i].pcrSelect[pcr / 8])
                    & (1 << (pcr % 8))) != 0)
                    allocated = TRUE;
                else
                    allocated = FALSE;
                break;
            }
        }
    }
    return allocated;
}

//*** GetPcrPointer()
// This function returns the address of an array of PCR based on the
// hash algorithm.
//  Return Type: BYTE *
//      NULL            no such algorithm
//      != NULL         pointer to the 0th byte of the 0th PCR
static BYTE *
GetPcrPointer(
    TPM_ALG_ID       alg,           // IN: algorithm for bank
    UINT32           pcrNumber      // IN: PCR number
    )
{
    static BYTE     *pcr = NULL;

    if(!PcrIsAllocated(pcrNumber, alg))
        return NULL;

    switch(alg)
    {
#if ALG_SHA1
        case ALG_SHA1_VALUE:
            pcr = s_pcrs[pcrNumber].sha1Pcr;
            break;
#endif
#if ALG_SHA256
        case ALG_SHA256_VALUE:
            pcr = s_pcrs[pcrNumber].sha256Pcr;
            break;
#endif
#if ALG_SHA384
        case ALG_SHA384_VALUE:
            pcr = s_pcrs[pcrNumber].sha384Pcr;
            break;
#endif
#if ALG_SHA512
        case ALG_SHA512_VALUE:
            pcr = s_pcrs[pcrNumber].sha512Pcr;
            break;
#endif
#if ALG_SM3_256
        case ALG_SM3_256_VALUE:
            pcr = s_pcrs[pcrNumber].sm3_256Pcr;
            break;
#endif
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
    return pcr;
}

//*** IsPcrSelected()
// This function indicates if an indicated PCR number is selected by the bit map in
// 'selection'.
//  Return Type: BOOL
//      TRUE(1)         PCR is selected
//      FALSE(0)        PCR is not selected
static BOOL
IsPcrSelected(
    UINT32               pcr,           // IN: The number of the PCR
    TPMS_PCR_SELECTION  *selection      // IN: The selection structure
    )
{
    BOOL                 selected;
    selected = (pcr < IMPLEMENTATION_PCR
                && ((selection->pcrSelect[pcr / 8]) & (1 << (pcr % 8))) != 0);
    return selected;
}

//*** FilterPcr()
// This function modifies a PCR selection array based on the implemented
// PCR.
static void
FilterPcr(
    TPMS_PCR_SELECTION  *selection      // IN: input PCR selection
    )
{
    UINT32     i;
    TPMS_PCR_SELECTION      *allocated = NULL;

    // If size of select is less than PCR_SELECT_MAX, zero the unspecified PCR
    for(i = selection->sizeofSelect; i < PCR_SELECT_MAX; i++)
        selection->pcrSelect[i] = 0;

    // Find the internal configuration for the bank
    for(i = 0; i < gp.pcrAllocated.count; i++)
    {
        if(gp.pcrAllocated.pcrSelections[i].hash == selection->hash)
        {
            allocated = &gp.pcrAllocated.pcrSelections[i];
            break;
        }
    }

    for(i = 0; i < selection->sizeofSelect; i++)
    {
        if(allocated == NULL)
        {
            // If the required bank does not exist, clear input selection
            selection->pcrSelect[i] = 0;
        }
        else
            selection->pcrSelect[i] &= allocated->pcrSelect[i];
    }

    return;
}

//*** PcrDrtm()
// This function does the DRTM and H-CRTM processing it is called from
// _TPM_Hash_End.
void
PcrDrtm(
    const TPMI_DH_PCR        pcrHandle,     // IN: the index of the PCR to be
                                            //     modified
    const TPMI_ALG_HASH      hash,          // IN: the bank identifier
    const TPM2B_DIGEST      *digest         // IN: the digest to modify the PCR
    )
{
    BYTE        *pcrData = GetPcrPointer(hash, pcrHandle);

    if(pcrData != NULL)
    {
        // Rest the PCR to zeros
        MemorySet(pcrData, 0, digest->t.size);

        // if the TPM has not started, then set the PCR to 0...04 and then extend
        if(!TPMIsStarted())
        {
            pcrData[digest->t.size - 1] = 4;
        }
        // Now, extend the value
        PCRExtend(pcrHandle, hash, digest->t.size, (BYTE *)digest->t.buffer);
    }
}

//*** PCR_ClearAuth()
// This function is used to reset the PCR authorization values. It is called
// on TPM2_Startup(CLEAR) and TPM2_Clear().
void
PCR_ClearAuth(
    void
    )
{
#if defined NUM_AUTHVALUE_PCR_GROUP && NUM_AUTHVALUE_PCR_GROUP > 0
    int         j;
    for(j = 0; j < NUM_AUTHVALUE_PCR_GROUP; j++)
    {
        gc.pcrAuthValues.auth[j].t.size = 0;
    }
#endif
}

//*** PCRStartup()
// This function initializes the PCR subsystem at TPM2_Startup().
BOOL
PCRStartup(
    STARTUP_TYPE     type,          // IN: startup type
    BYTE             locality       // IN: startup locality
    )
{
    UINT32              pcr, j;
    UINT32              saveIndex = 0;

    g_pcrReConfig = FALSE;

    // Don't test for SU_RESET because that should be the default when nothing 
    // else is selected
    if(type != SU_RESUME && type != SU_RESTART)
    {
        // PCR generation counter is cleared at TPM_RESET
        gr.pcrCounter = 0;
    }

    // Initialize/Restore PCR values
    for(pcr = 0; pcr < IMPLEMENTATION_PCR; pcr++)
    {
        // On resume, need to know if this PCR had its state saved or not
        UINT32      stateSaved;

        if(type == SU_RESUME
           &&  s_initAttributes[pcr].stateSave == SET)
        {
            stateSaved = 1;
        }
        else
        {
            stateSaved = 0;
            PCRChanged(pcr);
        }

        // If this is the H-CRTM PCR and we are not doing a resume and we
        // had an H-CRTM event, then we don't change this PCR
        if(pcr == HCRTM_PCR && type != SU_RESUME && g_DrtmPreStartup == TRUE)
            continue;

        // Iterate each hash algorithm bank
        for(j = 0; j < gp.pcrAllocated.count; j++)
        {
            TPMI_ALG_HASH    hash = gp.pcrAllocated.pcrSelections[j].hash;
            BYTE            *pcrData = GetPcrPointer(hash, pcr);
            UINT16           pcrSize = CryptHashGetDigestSize(hash);

            if(pcrData != NULL)
            {
                // if state was saved
                if(stateSaved == 1)
                {
                    // Restore saved PCR value
                    BYTE    *pcrSavedData;
                    pcrSavedData = GetSavedPcrPointer(
                        gp.pcrAllocated.pcrSelections[j].hash,
                        saveIndex);
                    if(pcrSavedData == NULL)
                        return FALSE;
                    MemoryCopy(pcrData, pcrSavedData, pcrSize);
                }
                else
                    // PCR was not restored by state save
                {
                    // If the reset locality of the PCR is 4, then
                    // the reset value is all one's, otherwise it is
                    // all zero.
                    if((s_initAttributes[pcr].resetLocality & 0x10) != 0)
                        MemorySet(pcrData, 0xFF, pcrSize);
                    else
                    {
                        MemorySet(pcrData, 0, pcrSize);
                        if(pcr == HCRTM_PCR)
                            pcrData[pcrSize - 1] = locality;
                    }
                }
            }
        }
        saveIndex += stateSaved;
    }
    // Reset authValues on TPM2_Startup(CLEAR)
    if(type != SU_RESUME)
        PCR_ClearAuth();
    return TRUE;
}

//*** PCRStateSave()
// This function is used to save the PCR values that will be restored on TPM Resume.
void
PCRStateSave(
    TPM_SU           type           // IN: startup type
    )
{
    UINT32           pcr, j;
    UINT32           saveIndex = 0;

    // if state save CLEAR, nothing to be done.  Return here
    if(type == TPM_SU_CLEAR) 
        return;

    // Copy PCR values to the structure that should be saved to NV
    for(pcr = 0; pcr < IMPLEMENTATION_PCR; pcr++)
    {
        UINT32  stateSaved = (s_initAttributes[pcr].stateSave == SET) ? 1 : 0;

        // Iterate each hash algorithm bank
        for(j = 0; j < gp.pcrAllocated.count; j++)
        {
            BYTE    *pcrData;
            UINT32  pcrSize;

            pcrData = GetPcrPointer(gp.pcrAllocated.pcrSelections[j].hash, pcr);

            if(pcrData != NULL)
            {
                pcrSize
                    = CryptHashGetDigestSize(gp.pcrAllocated.pcrSelections[j].hash);

                if(stateSaved == 1)
                {
                    // Restore saved PCR value
                    BYTE    *pcrSavedData;
                    pcrSavedData
                        = GetSavedPcrPointer(gp.pcrAllocated.pcrSelections[j].hash,
                                             saveIndex);
                    MemoryCopy(pcrSavedData, pcrData, pcrSize);
                }
            }
        }
        saveIndex += stateSaved;
    }

    return;
}

//*** PCRIsStateSaved()
// This function indicates if the selected PCR is a PCR that is state saved
// on TPM2_Shutdown(STATE). The return value is based on PCR attributes.
//  Return Type: BOOL
//      TRUE(1)         PCR is state saved
//      FALSE(0)        PCR is not state saved
BOOL
PCRIsStateSaved(
    TPMI_DH_PCR      handle         // IN: PCR handle to be extended
    )
{
    UINT32              pcr = handle - PCR_FIRST;

    if(s_initAttributes[pcr].stateSave == SET)
        return TRUE;
    else
        return FALSE;
}

//*** PCRIsResetAllowed()
// This function indicates if a PCR may be reset by the current command locality.
// The return value is based on PCR attributes, and not the PCR allocation.
//  Return Type: BOOL
//      TRUE(1)         TPM2_PCR_Reset is allowed
//      FALSE(0)        TPM2_PCR_Reset is not allowed
BOOL
PCRIsResetAllowed(
    TPMI_DH_PCR      handle         // IN: PCR handle to be extended
    )
{
    UINT8               commandLocality;
    UINT8               localityBits = 1;
    UINT32              pcr = handle - PCR_FIRST;

    // Check for the locality
    commandLocality = _plat__LocalityGet();

#ifdef DRTM_PCR
    // For a TPM that does DRTM, Reset is not allowed at locality 4
    if(commandLocality == 4)
        return FALSE;
#endif

    localityBits = localityBits << commandLocality;
    if((localityBits & s_initAttributes[pcr].resetLocality) == 0)
        return FALSE;
    else
        return TRUE;
}

//*** PCRChanged()
// This function checks a PCR handle to see if the attributes for the PCR are set
// so that any change to the PCR causes an increment of the pcrCounter. If it does,
// then the function increments the counter. Will also bump the counter if the 
// handle is zero which means that PCR 0 can not be in the TCB group. Bump on zero
// is used by TPM2_Clear().
void
PCRChanged(
    TPM_HANDLE       pcrHandle      // IN: the handle of the PCR that changed.
    )
{
    // For the reference implementation, the only change that does not cause
    // increment is a change to a PCR in the TCB group.
    if((pcrHandle == 0) || !PCRBelongsTCBGroup(pcrHandle))
    {
        gr.pcrCounter++;
        if(gr.pcrCounter == 0)
            FAIL(FATAL_ERROR_COUNTER_OVERFLOW);
    }
}

//*** PCRIsExtendAllowed()
// This function indicates a PCR may be extended at the current command locality.
// The return value is based on PCR attributes, and not the PCR allocation.
//  Return Type: BOOL
//      TRUE(1)         extend is allowed
//      FALSE(0)        extend is not allowed
BOOL
PCRIsExtendAllowed(
    TPMI_DH_PCR      handle         // IN: PCR handle to be extended
    )
{
    UINT8               commandLocality;
    UINT8               localityBits = 1;
    UINT32              pcr = handle - PCR_FIRST;

    // Check for the locality
    commandLocality = _plat__LocalityGet();
    localityBits = localityBits << commandLocality;
    if((localityBits & s_initAttributes[pcr].extendLocality) == 0)
        return FALSE;
    else
        return TRUE;
}

//*** PCRExtend()
// This function is used to extend a PCR in a specific bank.
void
PCRExtend(
    TPMI_DH_PCR      handle,        // IN: PCR handle to be extended
    TPMI_ALG_HASH    hash,          // IN: hash algorithm of PCR
    UINT32           size,          // IN: size of data to be extended
    BYTE            *data           // IN: data to be extended
    )
{
    BYTE                *pcrData;
    HASH_STATE           hashState;
    UINT16               pcrSize;

    pcrData = GetPcrPointer(hash, handle - PCR_FIRST);

    // Extend PCR if it is allocated
    if(pcrData != NULL)
    {
        pcrSize = CryptHashGetDigestSize(hash);
        CryptHashStart(&hashState, hash);
        CryptDigestUpdate(&hashState, pcrSize, pcrData);
        CryptDigestUpdate(&hashState, size, data);
        CryptHashEnd(&hashState, pcrSize, pcrData);

        // PCR has changed so update the pcrCounter if necessary
        PCRChanged(handle);
    }

    return;
}

//*** PCRComputeCurrentDigest()
// This function computes the digest of the selected PCR.
//
// As a side-effect, 'selection' is modified so that only the implemented PCR
// will have their bits still set.
void
PCRComputeCurrentDigest(
    TPMI_ALG_HASH        hashAlg,       // IN: hash algorithm to compute digest
    TPML_PCR_SELECTION  *selection,     // IN/OUT: PCR selection (filtered on
                                        //     output)
    TPM2B_DIGEST        *digest         // OUT: digest
    )
{
    HASH_STATE               hashState;
    TPMS_PCR_SELECTION      *select;
    BYTE                    *pcrData;   // will point to a digest
    UINT32                   pcrSize;
    UINT32                   pcr;
    UINT32                   i;

    // Initialize the hash
    digest->t.size = CryptHashStart(&hashState, hashAlg);
    pAssert(digest->t.size > 0 && digest->t.size < UINT16_MAX);

    // Iterate through the list of PCR selection structures
    for(i = 0; i < selection->count; i++)
    {
        // Point to the current selection
        select = &selection->pcrSelections[i]; // Point to the current selection
        FilterPcr(select);      // Clear out the bits for unimplemented PCR

        // Need the size of each digest
        pcrSize = CryptHashGetDigestSize(selection->pcrSelections[i].hash);

        // Iterate through the selection
        for(pcr = 0; pcr < IMPLEMENTATION_PCR; pcr++)
        {
            if(IsPcrSelected(pcr, select))         // Is this PCR selected
            {
                // Get pointer to the digest data for the bank
                pcrData = GetPcrPointer(selection->pcrSelections[i].hash, pcr);
                pAssert(pcrData != NULL);
                CryptDigestUpdate(&hashState, pcrSize, pcrData);  // add to digest
            }
        }
    }
    // Complete hash stack
    CryptHashEnd2B(&hashState, &digest->b);

    return;
}

//*** PCRRead()
// This function is used to read a list of selected PCR.  If the requested PCR
// number exceeds the maximum number that can be output, the 'selection' is
// adjusted to reflect the actual output PCR.
void
PCRRead(
    TPML_PCR_SELECTION  *selection,     // IN/OUT: PCR selection (filtered on
                                        //     output)
    TPML_DIGEST         *digest,        // OUT: digest
    UINT32              *pcrCounter     // OUT: the current value of PCR generation
                                        //     number
    )
{
    TPMS_PCR_SELECTION      *select;
    BYTE                    *pcrData;       // will point to a digest
    UINT32                   pcr;
    UINT32                   i;

    digest->count = 0;

    // Iterate through the list of PCR selection structures
    for(i = 0; i < selection->count; i++)
    {
        // Point to the current selection
        select = &selection->pcrSelections[i]; // Point to the current selection
        FilterPcr(select);      // Clear out the bits for unimplemented PCR

        // Iterate through the selection
        for(pcr = 0; pcr < IMPLEMENTATION_PCR; pcr++)
        {
            if(IsPcrSelected(pcr, select))         // Is this PCR selected
            {
                // Check if number of digest exceed upper bound
                if(digest->count > 7)
                {
                    // Clear rest of the current select bitmap
                    while(pcr < IMPLEMENTATION_PCR
                             // do not round up!
                          && (pcr / 8) < select->sizeofSelect)
                    {
                        // do not round up!
                        select->pcrSelect[pcr / 8] &= (BYTE)~(1 << (pcr % 8));
                        pcr++;
                    }
                    // Exit inner loop
                    break;
                }
                // Need the size of each digest
                digest->digests[digest->count].t.size =
                    CryptHashGetDigestSize(selection->pcrSelections[i].hash);

                // Get pointer to the digest data for the bank
                pcrData = GetPcrPointer(selection->pcrSelections[i].hash, pcr);
                pAssert(pcrData != NULL);
                // Add to the data to digest
                MemoryCopy(digest->digests[digest->count].t.buffer,
                           pcrData,
                           digest->digests[digest->count].t.size);
                digest->count++;
            }
        }
        // If we exit inner loop because we have exceed the output upper bound
        if(digest->count > 7 && pcr < IMPLEMENTATION_PCR)
        {
            // Clear rest of the selection
            while(i < selection->count)
            {
                MemorySet(selection->pcrSelections[i].pcrSelect, 0,
                          selection->pcrSelections[i].sizeofSelect);
                i++;
            }
            // exit outer loop
            break;
        }
    }

    *pcrCounter = gr.pcrCounter;

    return;
}

//*** PCRAllocate()
// This function is used to change the PCR allocation.
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        allocate failed
//      TPM_RC_PCR              improper allocation
TPM_RC
PCRAllocate(
    TPML_PCR_SELECTION  *allocate,      // IN: required allocation
    UINT32              *maxPCR,        // OUT: Maximum number of PCR
    UINT32              *sizeNeeded,    // OUT: required space
    UINT32              *sizeAvailable  // OUT: available space
    )
{
    UINT32                  i, j, k;
    TPML_PCR_SELECTION      newAllocate;
    // Initialize the flags to indicate if HCRTM PCR and DRTM PCR are allocated.
    BOOL                    pcrHcrtm = FALSE;
    BOOL                    pcrDrtm = FALSE;

    // Create the expected new PCR allocation based on the existing allocation
    // and the new input:
    //  1. if a PCR bank does not appear in the new allocation, the existing
    //     allocation of this PCR bank will be preserved.
    //  2. if a PCR bank appears multiple times in the new allocation, only the
    //     last one will be in effect.
    newAllocate = gp.pcrAllocated;
    for(i = 0; i < allocate->count; i++)
    {
        for(j = 0; j < newAllocate.count; j++)
        {
            // If hash matches, the new allocation covers the old allocation
            // for this particular bank.
            // The assumption is the initial PCR allocation (from manufacture)
            // has all the supported hash algorithms with an assigned bank
            // (possibly empty).  So there must be a match for any new bank
            // allocation from the input.
            if(newAllocate.pcrSelections[j].hash ==
               allocate->pcrSelections[i].hash)
            {
                newAllocate.pcrSelections[j] = allocate->pcrSelections[i];
                break;
            }
        }
        // The j loop must exit with a match.
        pAssert(j < newAllocate.count);
    }

    // Max PCR in a bank is MIN(implemented PCR, PCR with attributes defined)
    *maxPCR = sizeof(s_initAttributes) / sizeof(PCR_Attributes);
    if(*maxPCR > IMPLEMENTATION_PCR)
        *maxPCR = IMPLEMENTATION_PCR;

    // Compute required size for allocation
    *sizeNeeded = 0;
    for(i = 0; i < newAllocate.count; i++)
    {
        UINT32      digestSize
            = CryptHashGetDigestSize(newAllocate.pcrSelections[i].hash);
#if defined(DRTM_PCR)
        // Make sure that we end up with at least one DRTM PCR
        pcrDrtm = pcrDrtm || TestBit(DRTM_PCR,
                                     newAllocate.pcrSelections[i].pcrSelect,
                                     newAllocate.pcrSelections[i].sizeofSelect);

#else   // if DRTM PCR is not required, indicate that the allocation is OK
        pcrDrtm = TRUE;
#endif

#if defined(HCRTM_PCR)
        // and one HCRTM PCR (since this is usually PCR 0...)
        pcrHcrtm = pcrHcrtm || TestBit(HCRTM_PCR,
                                       newAllocate.pcrSelections[i].pcrSelect,
                                       newAllocate.pcrSelections[i].sizeofSelect);
#else
        pcrHcrtm = TRUE;
#endif
        for(j = 0; j < newAllocate.pcrSelections[i].sizeofSelect; j++)
        {
            BYTE        mask = 1;
            for(k = 0; k < 8; k++)
            {
                if((newAllocate.pcrSelections[i].pcrSelect[j] & mask) != 0)
                    *sizeNeeded += digestSize;
                mask = mask << 1;
            }
        }
    }

    if(!pcrDrtm || !pcrHcrtm)
        return TPM_RC_PCR;

    // In this particular implementation, we always have enough space to
    // allocate PCR.  Different implementation may return a sizeAvailable less
    // than the sizeNeed.
    *sizeAvailable = sizeof(s_pcrs);

    // Save the required allocation to NV.  Note that after NV is written, the
    // PCR allocation in NV is no longer consistent with the RAM data
    // gp.pcrAllocated.  The NV version reflect the allocate after next
    // TPM_RESET, while the RAM version reflects the current allocation
    NV_WRITE_PERSISTENT(pcrAllocated, newAllocate);

    return TPM_RC_SUCCESS;
}

//*** PCRSetValue()
// This function is used to set the designated PCR in all banks to an initial value.
// The initial value is signed and will be sign extended into the entire PCR.
//
void
PCRSetValue(
    TPM_HANDLE       handle,        // IN: the handle of the PCR to set
    INT8             initialValue   // IN: the value to set
    )
{
    int              i;
    UINT32           pcr = handle - PCR_FIRST;
    TPMI_ALG_HASH    hash;
    UINT16           digestSize;
    BYTE            *pcrData;

    // Iterate supported PCR bank algorithms to reset
    for(i = 0; i < HASH_COUNT; i++)
    {
        hash = CryptHashGetAlgByIndex(i);
        // Prevent runaway
        if(hash == TPM_ALG_NULL)
            break;

        // Get a pointer to the data
        pcrData = GetPcrPointer(gp.pcrAllocated.pcrSelections[i].hash, pcr);

        // If the PCR is allocated
        if(pcrData != NULL)
        {
            // And the size of the digest
            digestSize = CryptHashGetDigestSize(hash);

            // Set the LSO to the input value
            pcrData[digestSize - 1] = initialValue;

            // Sign extend
            if(initialValue >= 0)
                MemorySet(pcrData, 0, digestSize - 1);
            else
                MemorySet(pcrData, -1, digestSize - 1);
        }
    }
}

//*** PCRResetDynamics
// This function is used to reset a dynamic PCR to 0.  This function is used in
// DRTM sequence.
void
PCRResetDynamics(
    void
    )
{
    UINT32              pcr, i;

    // Initialize PCR values
    for(pcr = 0; pcr < IMPLEMENTATION_PCR; pcr++)
    {
        // Iterate each hash algorithm bank
        for(i = 0; i < gp.pcrAllocated.count; i++)
        {
            BYTE    *pcrData;
            UINT32  pcrSize;

            pcrData = GetPcrPointer(gp.pcrAllocated.pcrSelections[i].hash, pcr);

            if(pcrData != NULL)
            {
                pcrSize =
                    CryptHashGetDigestSize(gp.pcrAllocated.pcrSelections[i].hash);

                // Reset PCR
                // Any PCR can be reset by locality 4 should be reset to 0
                if((s_initAttributes[pcr].resetLocality & 0x10) != 0)
                    MemorySet(pcrData, 0, pcrSize);
            }
        }
    }
    return;
}

//*** PCRCapGetAllocation()
// This function is used to get the current allocation of PCR banks.
//  Return Type: TPMI_YES_NO
//      YES         if the return count is 0
//      NO          if the return count is not 0
TPMI_YES_NO
PCRCapGetAllocation(
    UINT32               count,         // IN: count of return
    TPML_PCR_SELECTION  *pcrSelection   // OUT: PCR allocation list
    )
{
    if(count == 0)
    {
        pcrSelection->count = 0;
        return YES;
    }
    else
    {
        *pcrSelection = gp.pcrAllocated;
        return NO;
    }
}

//*** PCRSetSelectBit()
// This function sets a bit in a bitmap array.
static void
PCRSetSelectBit(
    UINT32           pcr,           // IN: PCR number
    BYTE            *bitmap         // OUT: bit map to be set
    )
{
    bitmap[pcr / 8] |= (1 << (pcr % 8));
    return;
}

//*** PCRGetProperty()
// This function returns the selected PCR property.
//  Return Type: BOOL
//      TRUE(1)         the property type is implemented
//      FALSE(0)        the property type is not implemented
static BOOL
PCRGetProperty(
    TPM_PT_PCR               property,
    TPMS_TAGGED_PCR_SELECT  *select
    )
{
    UINT32              pcr;
    UINT32              groupIndex;

    select->tag = property;
    // Always set the bitmap to be the size of all PCR
    select->sizeofSelect = (IMPLEMENTATION_PCR + 7) / 8;

    // Initialize bitmap
    MemorySet(select->pcrSelect, 0, select->sizeofSelect);

    // Collecting properties
    for(pcr = 0; pcr < IMPLEMENTATION_PCR; pcr++)
    {
        switch(property)
        {
            case TPM_PT_PCR_SAVE:
                if(s_initAttributes[pcr].stateSave == SET)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_EXTEND_L0:
                if((s_initAttributes[pcr].extendLocality & 0x01) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_RESET_L0:
                if((s_initAttributes[pcr].resetLocality & 0x01) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_EXTEND_L1:
                if((s_initAttributes[pcr].extendLocality & 0x02) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_RESET_L1:
                if((s_initAttributes[pcr].resetLocality & 0x02) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_EXTEND_L2:
                if((s_initAttributes[pcr].extendLocality & 0x04) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_RESET_L2:
                if((s_initAttributes[pcr].resetLocality & 0x04) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_EXTEND_L3:
                if((s_initAttributes[pcr].extendLocality & 0x08) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_RESET_L3:
                if((s_initAttributes[pcr].resetLocality & 0x08) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_EXTEND_L4:
                if((s_initAttributes[pcr].extendLocality & 0x10) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_RESET_L4:
                if((s_initAttributes[pcr].resetLocality & 0x10) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
            case TPM_PT_PCR_DRTM_RESET:
                // DRTM reset PCRs are the PCR reset by locality 4
                if((s_initAttributes[pcr].resetLocality & 0x10) != 0)
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
#if defined NUM_POLICY_PCR_GROUP && NUM_POLICY_PCR_GROUP > 0
            case TPM_PT_PCR_POLICY:
                if(PCRBelongsPolicyGroup(pcr + PCR_FIRST, &groupIndex))
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
#endif
#if defined NUM_AUTHVALUE_PCR_GROUP && NUM_AUTHVALUE_PCR_GROUP > 0
            case TPM_PT_PCR_AUTH:
                if(PCRBelongsAuthGroup(pcr + PCR_FIRST, &groupIndex))
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
#endif
#if ENABLE_PCR_NO_INCREMENT == YES
            case TPM_PT_PCR_NO_INCREMENT:
                if(PCRBelongsTCBGroup(pcr + PCR_FIRST))
                    PCRSetSelectBit(pcr, select->pcrSelect);
                break;
#endif
            default:
                // If property is not supported, stop scanning PCR attributes
                // and return.
                return FALSE;
                break;
        }
    }
    return TRUE;
}

//*** PCRCapGetProperties()
// This function returns a list of PCR properties starting at 'property'.
//  Return Type: TPMI_YES_NO
//      YES         if no more property is available
//      NO          if there are more properties not reported
TPMI_YES_NO
PCRCapGetProperties(
    TPM_PT_PCR                   property,      // IN: the starting PCR property
    UINT32                       count,         // IN: count of returned properties
    TPML_TAGGED_PCR_PROPERTY    *select         // OUT: PCR select
    )
{
    TPMI_YES_NO     more = NO;
    UINT32          i;

    // Initialize output property list
    select->count = 0;

    // The maximum count of properties we may return is MAX_PCR_PROPERTIES
    if(count > MAX_PCR_PROPERTIES) count = MAX_PCR_PROPERTIES;

    // TPM_PT_PCR_FIRST is defined as 0 in spec.  It ensures that property
    // value would never be less than TPM_PT_PCR_FIRST
    cAssert(TPM_PT_PCR_FIRST == 0);

    // Iterate PCR properties. TPM_PT_PCR_LAST is the index of the last property
    // implemented on the TPM.
    for(i = property; i <= TPM_PT_PCR_LAST; i++)
    {
        if(select->count < count)
        {
            // If we have not filled up the return list, add more properties to it
            if(PCRGetProperty(i, &select->pcrProperty[select->count]))
                // only increment if the property is implemented
                select->count++;
        }
        else
        {
            // If the return list is full but we still have properties
            // available, report this and stop iterating.
            more = YES;
            break;
        }
    }
    return more;
}

//*** PCRCapGetHandles()
// This function is used to get a list of handles of PCR, started from 'handle'.
// If 'handle' exceeds the maximum PCR handle range, an empty list will be
// returned and the return value will be NO.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
PCRCapGetHandles(
    TPMI_DH_PCR      handle,        // IN: start handle
    UINT32           count,         // IN: count of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
    )
{
    TPMI_YES_NO     more = NO;
    UINT32          i;

    pAssert(HandleGetType(handle) == TPM_HT_PCR);

    // Initialize output handle list
    handleList->count = 0;

    // The maximum count of handles we may return is MAX_CAP_HANDLES
    if(count > MAX_CAP_HANDLES) count = MAX_CAP_HANDLES;

    // Iterate PCR handle range
    for(i = handle & HR_HANDLE_MASK; i <= PCR_LAST; i++)
    {
        if(handleList->count < count)
        {
            // If we have not filled up the return list, add this PCR
            // handle to it
            handleList->handle[handleList->count] = i + PCR_FIRST;
            handleList->count++;
        }
        else
        {
            // If the return list is full but we still have PCR handle
            // available, report this and stop iterating
            more = YES;
            break;
        }
    }
    return more;
}