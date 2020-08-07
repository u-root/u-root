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
// This file contains the functions used for managing and accessing the
// hierarchy-related values.

//** Includes

#include "Tpm.h"

//** Functions

//*** HierarchyPreInstall()
// This function performs the initialization functions for the hierarchy
// when the TPM is simulated. This function should not be called if the
// TPM is not in a manufacturing mode at the manufacturer, or in a simulated
// environment.
void
HierarchyPreInstall_Init(
    void
    )
{
    // Allow lockout clear command
    gp.disableClear = FALSE;

    // Initialize Primary Seeds
    gp.EPSeed.t.size = sizeof(gp.EPSeed.t.buffer);
    gp.SPSeed.t.size = sizeof(gp.SPSeed.t.buffer);
    gp.PPSeed.t.size = sizeof(gp.PPSeed.t.buffer);
#if (defined USE_PLATFORM_EPS) && (USE_PLATFORM_EPS != NO)
    _plat__GetEPS(gp.EPSeed.t.size, gp.EPSeed.t.buffer);
#else
    CryptRandomGenerate(gp.EPSeed.t.size, gp.EPSeed.t.buffer);
#endif
    CryptRandomGenerate(gp.SPSeed.t.size, gp.SPSeed.t.buffer);
    CryptRandomGenerate(gp.PPSeed.t.size, gp.PPSeed.t.buffer);

    // Initialize owner, endorsement and lockout authorization
    gp.ownerAuth.t.size = 0;
    gp.endorsementAuth.t.size = 0;
    gp.lockoutAuth.t.size = 0;

    // Initialize owner, endorsement, and lockout policy
    gp.ownerAlg = TPM_ALG_NULL;
    gp.ownerPolicy.t.size = 0;
    gp.endorsementAlg = TPM_ALG_NULL;
    gp.endorsementPolicy.t.size = 0;
    gp.lockoutAlg = TPM_ALG_NULL;
    gp.lockoutPolicy.t.size = 0;

    // Initialize ehProof, shProof and phProof
    gp.phProof.t.size = sizeof(gp.phProof.t.buffer);
    gp.shProof.t.size = sizeof(gp.shProof.t.buffer);
    gp.ehProof.t.size = sizeof(gp.ehProof.t.buffer);
    CryptRandomGenerate(gp.phProof.t.size, gp.phProof.t.buffer);
    CryptRandomGenerate(gp.shProof.t.size, gp.shProof.t.buffer);
    CryptRandomGenerate(gp.ehProof.t.size, gp.ehProof.t.buffer);

    // Write hierarchy data to NV
    NV_SYNC_PERSISTENT(disableClear);
    NV_SYNC_PERSISTENT(EPSeed);
    NV_SYNC_PERSISTENT(SPSeed);
    NV_SYNC_PERSISTENT(PPSeed);
    NV_SYNC_PERSISTENT(ownerAuth);
    NV_SYNC_PERSISTENT(endorsementAuth);
    NV_SYNC_PERSISTENT(lockoutAuth);
    NV_SYNC_PERSISTENT(ownerAlg);
    NV_SYNC_PERSISTENT(ownerPolicy);
    NV_SYNC_PERSISTENT(endorsementAlg);
    NV_SYNC_PERSISTENT(endorsementPolicy);
    NV_SYNC_PERSISTENT(lockoutAlg);
    NV_SYNC_PERSISTENT(lockoutPolicy);
    NV_SYNC_PERSISTENT(phProof);
    NV_SYNC_PERSISTENT(shProof);
    NV_SYNC_PERSISTENT(ehProof);

    return;
}

//*** HierarchyStartup()
// This function is called at TPM2_Startup() to initialize the hierarchy
// related values.
BOOL
HierarchyStartup(
    STARTUP_TYPE     type           // IN: start up type
    )
{
    // phEnable is SET on any startup
    g_phEnable = TRUE;

    // Reset platformAuth, platformPolicy; enable SH and EH at TPM_RESET and
    // TPM_RESTART
    if(type != SU_RESUME)
    {
        gc.platformAuth.t.size = 0;
        gc.platformPolicy.t.size = 0;
        gc.platformAlg = TPM_ALG_NULL;

        // enable the storage and endorsement hierarchies and the platformNV
        gc.shEnable = gc.ehEnable = gc.phEnableNV = TRUE;
    }

    // nullProof and nullSeed are updated at every TPM_RESET
    if((type != SU_RESTART) && (type != SU_RESUME))
    {
        gr.nullProof.t.size = sizeof(gr.nullProof.t.buffer);
        CryptRandomGenerate(gr.nullProof.t.size, gr.nullProof.t.buffer);
        gr.nullSeed.t.size = sizeof(gr.nullSeed.t.buffer);
        CryptRandomGenerate(gr.nullSeed.t.size, gr.nullSeed.t.buffer);
    }

    return TRUE;
}

//*** HierarchyGetProof()
// This function finds the proof value associated with a hierarchy.It returns a
// pointer to the proof value.
TPM2B_PROOF *
HierarchyGetProof(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy constant
    )
{
    TPM2B_PROOF         *proof = NULL;

    switch(hierarchy)
    {
        case TPM_RH_PLATFORM:
            // phProof for TPM_RH_PLATFORM
            proof = &gp.phProof;
            break;
        case TPM_RH_ENDORSEMENT:
            // ehProof for TPM_RH_ENDORSEMENT
            proof = &gp.ehProof;
            break;
        case TPM_RH_OWNER:
            // shProof for TPM_RH_OWNER
            proof = &gp.shProof;
            break;
        default:
            // nullProof for TPM_RH_NULL or anything else
            proof = &gr.nullProof;
            break;
    }
    return proof;
}

//*** HierarchyGetPrimarySeed()
// This function returns the primary seed of a hierarchy.
TPM2B_SEED *
HierarchyGetPrimarySeed(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy
    )
{
    TPM2B_SEED          *seed = NULL;
    switch(hierarchy)
    {
        case TPM_RH_PLATFORM:
            seed = &gp.PPSeed;
            break;
        case TPM_RH_OWNER:
            seed = &gp.SPSeed;
            break;
        case TPM_RH_ENDORSEMENT:
            seed = &gp.EPSeed;
            break;
         default:
            seed = &gr.nullSeed;
            break;
    }
    return seed;
}

//*** HierarchyIsEnabled()
// This function checks to see if a hierarchy is enabled.
// NOTE: The TPM_RH_NULL hierarchy is always enabled.
//  Return Type: BOOL
//      TRUE(1)         hierarchy is enabled
//      FALSE(0)        hierarchy is disabled
BOOL
HierarchyIsEnabled(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy
    )
{
    BOOL            enabled = FALSE;

    switch(hierarchy)
    {
        case TPM_RH_PLATFORM:
            enabled = g_phEnable;
            break;
        case TPM_RH_OWNER:
            enabled = gc.shEnable;
            break;
        case TPM_RH_ENDORSEMENT:
            enabled = gc.ehEnable;
            break;
        case TPM_RH_NULL:
            enabled = TRUE;
            break;
        default:
            enabled = FALSE;
            break;
    }
    return enabled;
}