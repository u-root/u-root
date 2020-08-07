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
 *  Date: Apr  2, 2019  Time: 04:23:27PM
 */

#ifndef    _HIERARCHY_FP_H_
#define    _HIERARCHY_FP_H_

//*** HierarchyPreInstall()
// This function performs the initialization functions for the hierarchy
// when the TPM is simulated. This function should not be called if the
// TPM is not in a manufacturing mode at the manufacturer, or in a simulated
// environment.
void
HierarchyPreInstall_Init(
    void
);

//*** HierarchyStartup()
// This function is called at TPM2_Startup() to initialize the hierarchy
// related values.
BOOL
HierarchyStartup(
    STARTUP_TYPE     type           // IN: start up type
);

//*** HierarchyGetProof()
// This function finds the proof value associated with a hierarchy.It returns a
// pointer to the proof value.
TPM2B_PROOF *
HierarchyGetProof(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy constant
);

//*** HierarchyGetPrimarySeed()
// This function returns the primary seed of a hierarchy.
TPM2B_SEED *
HierarchyGetPrimarySeed(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy
);

//*** HierarchyIsEnabled()
// This function checks to see if a hierarchy is enabled.
// NOTE: The TPM_RH_NULL hierarchy is always enabled.
//  Return Type: BOOL
//      TRUE(1)         hierarchy is enabled
//      FALSE(0)        hierarchy is disabled
BOOL
HierarchyIsEnabled(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy
);

#endif  // _HIERARCHY_FP_H_
