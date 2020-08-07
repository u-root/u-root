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
 *  Date: Apr  2, 2019  Time: 04:06:42PM
 */

#ifndef    _SESSION_FP_H_
#define    _SESSION_FP_H_

//** Startup Function -- SessionStartup()
// This function initializes the session subsystem on TPM2_Startup().
BOOL
SessionStartup(
    STARTUP_TYPE     type
);

//*** SessionIsLoaded()
// This function test a session handle references a loaded session.  The handle
// must have previously been checked to make sure that it is a valid handle for
// an authorization session.
// NOTE:    A PWAP authorization does not have a session.
//
//  Return Type: BOOL
//      TRUE(1)         session is loaded
//      FALSE(0)        session is not loaded
//
BOOL
SessionIsLoaded(
    TPM_HANDLE       handle         // IN: session handle
);

//*** SessionIsSaved()
// This function test a session handle references a saved session.  The handle
// must have previously been checked to make sure that it is a valid handle for
// an authorization session.
// NOTE:    An password authorization does not have a session.
//
// This function requires that the handle be a valid session handle.
//
//  Return Type: BOOL
//      TRUE(1)         session is saved
//      FALSE(0)        session is not saved
//
BOOL
SessionIsSaved(
    TPM_HANDLE       handle         // IN: session handle
);

//*** SequenceNumberForSavedContextIsValid()
// This function validates that the sequence number and handle value within a
// saved context are valid.
BOOL
SequenceNumberForSavedContextIsValid(
    TPMS_CONTEXT    *context        // IN: pointer to a context structure to be
                                    //     validated
);

//*** SessionPCRValueIsCurrent()
//
// This function is used to check if PCR values have been updated since the
// last time they were checked in a policy session.
//
// This function requires the session is loaded.
//  Return Type: BOOL
//      TRUE(1)         PCR value is current
//      FALSE(0)        PCR value is not current
BOOL
SessionPCRValueIsCurrent(
    SESSION         *session        // IN: session structure
);

//*** SessionGet()
// This function returns a pointer to the session object associated with a
// session handle.
//
// The function requires that the session is loaded.
SESSION *
SessionGet(
    TPM_HANDLE       handle         // IN: session handle
);

//*** SessionCreate()
//
//  This function does the detailed work for starting an authorization session.
//  This is done in a support routine rather than in the action code because
//  the session management may differ in implementations.  This implementation
//  uses a fixed memory allocation to hold sessions and a fixed allocation
//  to hold the contextID for the saved contexts.
//
//  Return Type: TPM_RC
//      TPM_RC_CONTEXT_GAP          need to recycle sessions
//      TPM_RC_SESSION_HANDLE       active session space is full
//      TPM_RC_SESSION_MEMORY       loaded session space is full
TPM_RC
SessionCreate(
    TPM_SE           sessionType,   // IN: the session type
    TPMI_ALG_HASH    authHash,      // IN: the hash algorithm
    TPM2B_NONCE     *nonceCaller,   // IN: initial nonceCaller
    TPMT_SYM_DEF    *symmetric,     // IN: the symmetric algorithm
    TPMI_DH_ENTITY   bind,          // IN: the bind object
    TPM2B_DATA      *seed,          // IN: seed data
    TPM_HANDLE      *sessionHandle, // OUT: the session handle
    TPM2B_NONCE     *nonceTpm       // OUT: the session nonce
);

//*** SessionContextSave()
// This function is called when a session context is to be saved.  The
// contextID of the saved session is returned.  If no contextID can be
// assigned, then the routine returns TPM_RC_CONTEXT_GAP.
// If the function completes normally, the session slot will be freed.
//
// This function requires that 'handle' references a loaded session.
// Otherwise, it should not be called at the first place.
//
//  Return Type: TPM_RC
//      TPM_RC_CONTEXT_GAP      a contextID could not be assigned.
//      TPM_RC_TOO_MANY_CONTEXTSthe counter maxed out
//
TPM_RC
SessionContextSave(
    TPM_HANDLE           handle,        // IN: session handle
    CONTEXT_COUNTER     *contextID      // OUT: assigned contextID
);

//*** SessionContextLoad()
// This function is used to load a session from saved context.  The session
// handle must be for a saved context.
//
// If the gap is at a maximum, then the only session that can be loaded is
// the oldest session, otherwise TPM_RC_CONTEXT_GAP is returned.
///
// This function requires that 'handle' references a valid saved session.
//
//  Return Type: TPM_RC
//      TPM_RC_SESSION_MEMORY       no free session slots
//      TPM_RC_CONTEXT_GAP          the gap count is maximum and this
//                                  is not the oldest saved context
//
TPM_RC
SessionContextLoad(
    SESSION_BUF     *session,       // IN: session structure from saved context
    TPM_HANDLE      *handle         // IN/OUT: session handle
);

//*** SessionFlush()
// This function is used to flush a session referenced by its handle.  If the
// session associated with 'handle' is loaded, the session array entry is
// marked as available.
//
// This function requires that 'handle' be a valid active session.
//
void
SessionFlush(
    TPM_HANDLE       handle         // IN: loaded or saved session handle
);

//*** SessionComputeBoundEntity()
// This function computes the binding value for a session.  The binding value
// for a reserved handle is the handle itself.  For all the other entities,
// the authValue at the time of binding is included to prevent squatting.
// For those values, the Name and the authValue are concatenated
// into the bind buffer.  If they will not both fit, the will be overlapped
// by XORing bytes.  If XOR is required, the bind value will be full.
void
SessionComputeBoundEntity(
    TPMI_DH_ENTITY       entityHandle,  // IN: handle of entity
    TPM2B_NAME          *bind           // OUT: binding value
);

//*** SessionSetStartTime()
// This function is used to initialize the session timing
void
SessionSetStartTime(
    SESSION         *session        // IN: the session to update
);

//*** SessionResetPolicyData()
// This function is used to reset the policy data without changing the nonce
// or the start time of the session.
void
SessionResetPolicyData(
    SESSION         *session        // IN: the session to reset
);

//*** SessionCapGetLoaded()
// This function returns a list of handles of loaded session, started
// from input 'handle'
//
// 'Handle' must be in valid loaded session handle range, but does not
// have to point to a loaded session.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
SessionCapGetLoaded(
    TPMI_SH_POLICY   handle,        // IN: start handle
    UINT32           count,         // IN: count of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
);

//*** SessionCapGetSaved()
// This function returns a list of handles for saved session, starting at
// 'handle'.
//
// 'Handle' must be in a valid handle range, but does not have to point to a
// saved session
//
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
SessionCapGetSaved(
    TPMI_SH_HMAC     handle,        // IN: start handle
    UINT32           count,         // IN: count of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
);

//*** SessionCapGetLoadedNumber()
// This function return the number of authorization sessions currently
// loaded into TPM RAM.
UINT32
SessionCapGetLoadedNumber(
    void
);

//*** SessionCapGetLoadedAvail()
// This function returns the number of additional authorization sessions, of
// any type, that could be loaded into TPM RAM.
// NOTE: In other implementations, this number may just be an estimate. The only
//       requirement for the estimate is, if it is one or more, then at least one
//       session must be loadable.
UINT32
SessionCapGetLoadedAvail(
    void
);

//*** SessionCapGetActiveNumber()
// This function returns the number of active authorization sessions currently
// being tracked by the TPM.
UINT32
SessionCapGetActiveNumber(
    void
);

//*** SessionCapGetActiveAvail()
// This function returns the number of additional authorization sessions, of any
// type, that could be created. This not the number of slots for sessions, but
// the number of additional sessions that the TPM is capable of tracking.
UINT32
SessionCapGetActiveAvail(
    void
);

#endif  // _SESSION_FP_H_
