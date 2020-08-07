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
#include "Tpm.h"
#include "DictionaryAttackParameters_fp.h"

#if CC_DictionaryAttackParameters  // Conditional expansion of this file

/*(See part 3 specification)
// change the lockout parameters
*/
TPM_RC
TPM2_DictionaryAttackParameters(
    DictionaryAttackParameters_In   *in             // IN: input parameter list
    )
{
    // The command needs NV update.
    RETURN_IF_NV_IS_NOT_AVAILABLE;

// Internal Data Update

    // Set dictionary attack parameters
    gp.maxTries = in->newMaxTries;
    gp.recoveryTime = in->newRecoveryTime;
    gp.lockoutRecovery = in->lockoutRecovery;

#if 0   // Errata eliminates this code
    // This functionality has been disabled. The preferred implementation is now
    // to leave failedTries unchanged when the parameters are changed. This could
    // have the effect of putting the TPM into DA lockout if in->newMaxTries is
    // not greater than the current value of gp.failedTries.
    // Set failed tries to 0
    gp.failedTries = 0;
#endif

    // Record the changes to NV
    NV_SYNC_PERSISTENT(failedTries);
    NV_SYNC_PERSISTENT(maxTries);
    NV_SYNC_PERSISTENT(recoveryTime);
    NV_SYNC_PERSISTENT(lockoutRecovery);

    return TPM_RC_SUCCESS;
}

#endif // CC_DictionaryAttackParameters