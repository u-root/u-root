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
//** Description
// This file contains the function that performs the "manufacturing" of the TPM
// in a simulated environment. These functions should not be used outside of
// a manufacturing or simulation environment.

//** Includes and Data Definitions
#define MANUFACTURE_C
#include "Tpm.h"
#include "TpmSizeChecks_fp.h"

//** Functions

//*** TPM_Manufacture()
// This function initializes the TPM values in preparation for the TPM's first
// use. This function will fail if previously called. The TPM can be re-manufactured
// by calling TPM_Teardown() first and then calling this function again.
//  Return Type: int
//      0           success
//      1           manufacturing process previously performed
LIB_EXPORT int
TPM_Manufacture(
    int             firstTime       // IN: indicates if this is the first call from
                                    //     main()
    )
{
    TPM_SU          orderlyShutdown;

#if RUNTIME_SIZE_CHECKS 
    // Call the function to verify the sizes of values that result from different
    // compile options.
    TpmSizeChecks();
#endif

    // If TPM has been manufactured, return indication.
    if(!firstTime && g_manufactured)
        return 1;

    // Do power on initializations of the cryptographic libraries.
    CryptInit();

    s_DAPendingOnNV = FALSE;

    // initialize NV
    NvManufacture();

    // Clear the magic value in the DRBG state
    go.drbgState.magic = 0;

    CryptStartup(SU_RESET);

    // default configuration for PCR
    PCRSimStart();

    // initialize pre-installed hierarchy data
    // This should happen after NV is initialized because hierarchy data is
    // stored in NV.
    HierarchyPreInstall_Init();

    // initialize dictionary attack parameters
    DAPreInstall_Init();

    // initialize PP list
    PhysicalPresencePreInstall_Init();

    // initialize command audit list
    CommandAuditPreInstall_Init();

    // first start up is required to be Startup(CLEAR)
    orderlyShutdown = TPM_SU_CLEAR;
    NV_WRITE_PERSISTENT(orderlyState, orderlyShutdown);

    // initialize the firmware version
    gp.firmwareV1 = FIRMWARE_V1;
#ifdef FIRMWARE_V2
    gp.firmwareV2 = FIRMWARE_V2;
#else
    gp.firmwareV2 = 0;
#endif
    NV_SYNC_PERSISTENT(firmwareV1);
    NV_SYNC_PERSISTENT(firmwareV2);

    // initialize the total reset counter to 0
    gp.totalResetCount = 0;
    NV_SYNC_PERSISTENT(totalResetCount);

    // initialize the clock stuff
    go.clock = 0;
    go.clockSafe = YES;

    NvWrite(NV_ORDERLY_DATA, sizeof(ORDERLY_DATA), &go);

    // Commit NV writes.  Manufacture process is an artificial process existing
    // only in simulator environment and it is not defined in the specification
    // that what should be the expected behavior if the NV write fails at this
    // point.  Therefore, it is assumed the NV write here is always success and
    // no return code of this function is checked.
    NvCommit();

    g_manufactured = TRUE;

    return 0;
}

//*** TPM_TearDown()
// This function prepares the TPM for re-manufacture. It should not be implemented
// in anything other than a simulated TPM.
//
// In this implementation, all that is needs is to stop the cryptographic units
// and set a flag to indicate that the TPM can be re-manufactured. This should
// be all that is necessary to start the manufacturing process again.
//  Return Type: int
//      0        success
//      1        TPM not previously manufactured
LIB_EXPORT int
TPM_TearDown(
    void
    )
{
    g_manufactured = FALSE;
    return 0;
}


//*** TpmEndSimulation()
// This function is called at the end of the simulation run. It is used to provoke
// printing of any statistics that might be needed.
LIB_EXPORT void
TpmEndSimulation(
    void
    )
{
#if SIMULATION
    HashLibSimulationEnd();
    SymLibSimulationEnd();
    MathLibSimulationEnd();
#if ALG_RSA
    RsaSimulationEnd();
#endif
#if ALG_ECC
    EccSimulationEnd();
#endif
#endif // SIMULATION
}