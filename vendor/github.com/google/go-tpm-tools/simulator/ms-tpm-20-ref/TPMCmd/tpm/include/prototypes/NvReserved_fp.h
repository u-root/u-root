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

#ifndef    _NV_RESERVED_FP_H_
#define    _NV_RESERVED_FP_H_

//*** NvCheckState()
// Function to check the NV state by accessing the platform-specific function
// to get the NV state.  The result state is registered in s_NvIsAvailable
// that will be reported by NvIsAvailable.
//
// This function is called at the beginning of ExecuteCommand before any potential
// check of g_NvStatus.
void
NvCheckState(
    void
);

//*** NvCommit
// This is a wrapper for the platform function to commit pending NV writes.
BOOL
NvCommit(
    void
);

//*** NvPowerOn()
//  This function is called at _TPM_Init to initialize the NV environment.
//  Return Type: BOOL
//      TRUE(1)         all NV was initialized
//      FALSE(0)        the NV containing saved state had an error and
//                      TPM2_Startup(CLEAR) is required
BOOL
NvPowerOn(
    void
);

//*** NvManufacture()
// This function initializes the NV system at pre-install time.
//
// This function should only be called in a manufacturing environment or in a
// simulation.
//
// The layout of NV memory space is an implementation choice.
void
NvManufacture(
    void
);

//*** NvRead()
// This function is used to move reserved data from NV memory to RAM.
void
NvRead(
    void            *outBuffer,     // OUT: buffer to receive data
    UINT32           nvOffset,      // IN: offset in NV of value
    UINT32           size           // IN: size of the value to read
);

//*** NvWrite()
// This function is used to post reserved data for writing to NV memory. Before
// the TPM completes the operation, the value will be written.
BOOL
NvWrite(
    UINT32           nvOffset,      // IN: location in NV to receive data
    UINT32           size,          // IN: size of the data to move
    void            *inBuffer       // IN: location containing data to write
);

//*** NvUpdatePersistent()
// This function is used to update a value in the PERSISTENT_DATA structure and
// commits the value to NV.
void
NvUpdatePersistent(
    UINT32           offset,        // IN: location in PERMANENT_DATA to be updated
    UINT32           size,          // IN: size of the value
    void            *buffer         // IN: the new data
);

//*** NvClearPersistent()
// This function is used to clear a persistent data entry and commit it to NV
void
NvClearPersistent(
    UINT32           offset,        // IN: the offset in the PERMANENT_DATA
                                    //     structure to be cleared (zeroed)
    UINT32           size           // IN: number of bytes to clear
);

//*** NvReadPersistent()
// This function reads persistent data to the RAM copy of the 'gp' structure.
void
NvReadPersistent(
    void
);

#endif  // _NV_RESERVED_FP_H_
