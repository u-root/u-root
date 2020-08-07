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
//**Introduction
/*
    The code in this file is used to manage the session context counter.
    The scheme implemented here is a "truncated counter".
    This scheme allows the TPM to not need TPM_SU_CLEAR for a
    very long period of time and still not have the context
    count for a session repeated.

    The counter (contextCounter)in this implementation is a UINT64 but
    can be smaller.  The "tracking array" (contextArray) only
    has 16-bits per context.  The tracking array is the data
    that needs to be saved and restored across TPM_SU_STATE so that
    sessions are not lost when the system enters the sleep state.
    Also, when the TPM is active, the tracking array is kept in
    RAM making it important that the number of bytes for each
    entry be kept as small as possible.

    The TPM prevents "collisions" of these truncated values by
    not allowing a contextID to be assigned if it would be the
    same as an existing value.  Since the array holds 16 bits,
    after a context has been saved, an additional 2^16-1 contexts
    may be saved before the count would again match.  The normal
    expectation is that the context will be flushed before its count
    value is needed again but it is always possible to have long-lived
    sessions.

    The contextID is assigned when the context is saved (TPM2_ContextSave()).
    At that time, the TPM will compare the low-order 16 bits of
    contextCounter to the existing values in contextArray and if one
    matches, the TPM will return TPM_RC_CONTEXT_GAP (by construction,
    the entry that contains the matching value is the oldest
    context).

    The expected remediation by the TRM is to load the oldest saved
    session context (the one found by the TPM), and save it.  Since loading
    the oldest session also eliminates its contextID value from
    contextArray, there TPM will always be able to load and save the oldest
    existing context.

    In the worst case, software may have to load and save several contexts
    in order to save an additional one.  This should happen very infrequently.

    When the TPM searches contextArray and finds that none of the contextIDs
    match the low-order 16-bits of contextCount, the TPM can copy the low bits
    to the contextArray associated with the session, and increment contextCount.

    There is one entry in contextArray for each of the active sessions
    allowed by the TPM implementation.  This array contains either a
    context count, an index, or a value indicating the slot is available (0).

    The index into the contextArray is the handle for the session with the region
    selector byte of the session set to zero.  If an entry in contextArray contains
    0, then the corresponding handle may be assigned to a session.  If the entry
    contains a value that is less than or equal to the number of loaded sessions
    for the TPM, then the array entry is the slot in which the context is loaded.

    EXAMPLE:    If the TPM allows 8 loaded sessions, then the slot numbers would
    be 1-8 and a contextArrary value in that range would represent the loaded
    session.

    NOTE:   When the TPM firmware determines that the array entry is for a loaded
    session, it will subtract 1 to create the zero-based slot number.

    There is one significant corner case in this scheme.  When the contextCount
    is equal to a value in the contextArray, the oldest session needs to be
    recycled or flushed. In order to recycle the session, it must be loaded.
    To be loaded, there must be an available slot.  Rather than require that a
    spare slot be available all the time, the TPM will check to see if the
    contextCount is equal to some value in the contextArray when a session is
    created.  This prevents the last session slot from being used when it
    is likely that a session will need to be recycled.

    If a TPM with both 1.2 and 2.0 functionality uses this scheme for both
    1.2 and 2.0 sessions, and the list of active contexts is read with
    TPM_GetCapabiltiy(), the TPM will create 32-bit representations of the
    list that contains 16-bit values (the TPM2_GetCapability() returns a list
    of handles for active sessions rather than a list of contextID).  The full
    contextID has high-order bits that are either the same as the current
    contextCount or one less.  It is one less if the 16-bits
    of the contextArray has a value that is larger than the low-order 16 bits
    of contextCount.
*/

//** Includes, Defines, and Local Variables
#define SESSION_C
#include "Tpm.h"

//** File Scope Function -- ContextIdSetOldest()
/*
    This function is called when the oldest contextID is being loaded or deleted.
    Once a saved context becomes the oldest, it stays the oldest until it is
    deleted.

    Finding the oldest is a bit tricky.  It is not just the numeric comparison of
    values but is dependent on the value of contextCounter.

    Assume we have a small contextArray with 8, 4-bit values with values 1 and 2
    used to indicate the loaded context slot number.  Also assume that the array
    contains hex values of (0 0 1 0 3 0 9 F) and that the contextCounter is an
    8-bit counter with a value of 0x37. Since the low nibble is 7, that means
    that values above 7 are older than values below it and, in this example,
    9 is the oldest value.

    Note if we subtract the counter value, from each slot that contains a saved
    contextID we get (- - - - B - 2 - 8) and the oldest entry is now easy to find.
*/
static void
ContextIdSetOldest(
    void
    )
{
    CONTEXT_SLOT    lowBits;
    CONTEXT_SLOT    entry;
    CONTEXT_SLOT    smallest = ((CONTEXT_SLOT)~0);
    UINT32  i;

    // Set oldestSaveContext to a value indicating none assigned
    s_oldestSavedSession = MAX_ACTIVE_SESSIONS + 1;

    lowBits = (CONTEXT_SLOT)gr.contextCounter;
    for(i = 0; i < MAX_ACTIVE_SESSIONS; i++)
    {
        entry = gr.contextArray[i];

        // only look at entries that are saved contexts
        if(entry > MAX_LOADED_SESSIONS)
        {
            // Use a less than or equal in case the oldest
            // is brand new (= lowBits-1) and equal to our initial
            // value for smallest.
            if(((CONTEXT_SLOT)(entry - lowBits)) <= smallest)
            {
                smallest = (entry - lowBits);
                s_oldestSavedSession = i;
            }
        }
    }
    // When we finish, either the s_oldestSavedSession still has its initial
    // value, or it has the index of the oldest saved context.
}

//** Startup Function -- SessionStartup()
// This function initializes the session subsystem on TPM2_Startup().
BOOL
SessionStartup(
    STARTUP_TYPE     type
    )
{
    UINT32               i;

    // Initialize session slots.  At startup, all the in-memory session slots
    // are cleared and marked as not occupied
    for(i = 0; i < MAX_LOADED_SESSIONS; i++)
        s_sessions[i].occupied = FALSE;   // session slot is not occupied

    // The free session slots the number of maximum allowed loaded sessions
    s_freeSessionSlots = MAX_LOADED_SESSIONS;

    // Initialize context ID data.  On a ST_SAVE or hibernate sequence, it will
    // scan the saved array of session context counts, and clear any entry that
    // references a session that was in memory during the state save since that
    // memory was not preserved over the ST_SAVE.
    if(type == SU_RESUME || type == SU_RESTART)
    {
        // On ST_SAVE we preserve the contexts that were saved but not the ones
        // in memory
        for(i = 0; i < MAX_ACTIVE_SESSIONS; i++)
        {
            // If the array value is unused or references a loaded session then
            // that loaded session context is lost and the array entry is
            // reclaimed.
            if(gr.contextArray[i] <= MAX_LOADED_SESSIONS)
                gr.contextArray[i] = 0;
        }
        // Find the oldest session in context ID data and set it in
        // s_oldestSavedSession
        ContextIdSetOldest();
    }
    else
    {
        // For STARTUP_CLEAR, clear out the contextArray
        for(i = 0; i < MAX_ACTIVE_SESSIONS; i++)
            gr.contextArray[i] = 0;

        // reset the context counter
        gr.contextCounter = MAX_LOADED_SESSIONS + 1;

        // Initialize oldest saved session
        s_oldestSavedSession = MAX_ACTIVE_SESSIONS + 1;
    }
    return TRUE;
}

//************************************************
//** Access Functions
//************************************************

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
    )
{
    pAssert(HandleGetType(handle) == TPM_HT_POLICY_SESSION
            || HandleGetType(handle) == TPM_HT_HMAC_SESSION);

    handle = handle & HR_HANDLE_MASK;

    // if out of range of possible active session, or not assigned to a loaded
    // session return false
    if(handle >= MAX_ACTIVE_SESSIONS
       || gr.contextArray[handle] == 0
       || gr.contextArray[handle] > MAX_LOADED_SESSIONS)
        return FALSE;

    return TRUE;
}

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
    )
{
    pAssert(HandleGetType(handle) == TPM_HT_POLICY_SESSION
            || HandleGetType(handle) == TPM_HT_HMAC_SESSION);

    handle = handle & HR_HANDLE_MASK;
    // if out of range of possible active session, or not assigned, or
    // assigned to a loaded session, return false
    if(handle >= MAX_ACTIVE_SESSIONS
       || gr.contextArray[handle] == 0
       || gr.contextArray[handle] <= MAX_LOADED_SESSIONS
       )
        return FALSE;

    return TRUE;
}

//*** SequenceNumberForSavedContextIsValid()
// This function validates that the sequence number and handle value within a 
// saved context are valid.
BOOL
SequenceNumberForSavedContextIsValid(
    TPMS_CONTEXT    *context        // IN: pointer to a context structure to be
                                    //     validated
    )
{
#define MAX_CONTEXT_GAP ((UINT64)((CONTEXT_SLOT) ~0) + 1)

    TPM_HANDLE           handle = context->savedHandle & HR_HANDLE_MASK;

    if(// Handle must be with the range of active sessions
       handle >= MAX_ACTIVE_SESSIONS
       // the array entry must be for a saved context
       || gr.contextArray[handle] <= MAX_LOADED_SESSIONS
           // the array entry must agree with the sequence number
       || gr.contextArray[handle] != (CONTEXT_SLOT)context->sequence
           // the provided sequence number has to be less than the current counter
       || context->sequence > gr.contextCounter
           // but not so much that it could not be a valid sequence number
       || gr.contextCounter - context->sequence > MAX_CONTEXT_GAP)
        return FALSE;

    return TRUE;
}

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
    )
{
    if(session->pcrCounter != 0
       && session->pcrCounter != gr.pcrCounter
       )
        return FALSE;
    else
        return TRUE;
}

//*** SessionGet()
// This function returns a pointer to the session object associated with a
// session handle.
//
// The function requires that the session is loaded.
SESSION *
SessionGet(
    TPM_HANDLE       handle         // IN: session handle
    )
{
    size_t          slotIndex;
    CONTEXT_SLOT    sessionIndex;

    pAssert(HandleGetType(handle) == TPM_HT_POLICY_SESSION
            || HandleGetType(handle) == TPM_HT_HMAC_SESSION
            );

    slotIndex = handle & HR_HANDLE_MASK;

    pAssert(slotIndex < MAX_ACTIVE_SESSIONS);

    // get the contents of the session array.  Because session is loaded, we
    // should always get a valid sessionIndex
    sessionIndex = gr.contextArray[slotIndex] - 1;

    pAssert(sessionIndex < MAX_LOADED_SESSIONS);

    return &s_sessions[sessionIndex].session;
}

//************************************************
//** Utility Functions
//************************************************

//*** ContextIdSessionCreate()
//
//  This function is called when a session is created.  It will check
//  to see if the current gap would prevent a context from being saved.  If
//  so it will return TPM_RC_CONTEXT_GAP.  Otherwise, it will try to find
//  an open slot in contextArray, set contextArray to the slot.
//  
//  This routine requires that the caller has determined the session array
//  index for the session.
//
//  Return Type: TPM_RC
//      TPM_RC_CONTEXT_GAP      can't assign a new contextID until the oldest
//                              saved session context is recycled
//      TPM_RC_SESSION_HANDLE   there is no slot available in the context array
//                              for tracking of this session context
static TPM_RC
ContextIdSessionCreate(
    TPM_HANDLE      *handle,        // OUT: receives the assigned handle. This will
                                    //     be an index that must be adjusted by the
                                    //     caller according to the type of the
                                    //     session created
    UINT32           sessionIndex   // IN: The session context array entry that will
                                    //     be occupied by the created session
    )
{
    pAssert(sessionIndex < MAX_LOADED_SESSIONS);

    // check to see if creating the context is safe
    // Is this going to be an assignment for the last session context
    // array entry?  If so, then there will be no room to recycle the
    // oldest context if needed.  If the gap is not at maximum, then
    // it will be possible to save a context if it becomes necessary.
    if(s_oldestSavedSession < MAX_ACTIVE_SESSIONS
       && s_freeSessionSlots == 1)
    {
        // See if the gap is at maximum
        // The current value of the contextCounter will be assigned to the next
        // saved context. If the value to be assigned would make the same as an
        // existing context, then we can't use it because of the ambiguity it would
        // create.
        if((CONTEXT_SLOT)gr.contextCounter
           == gr.contextArray[s_oldestSavedSession])
            return TPM_RC_CONTEXT_GAP;
    }

    // Find an unoccupied entry in the contextArray
    for(*handle = 0; *handle < MAX_ACTIVE_SESSIONS; (*handle)++)
    {
        if(gr.contextArray[*handle] == 0)
        {
            // indicate that the session associated with this handle
            // references a loaded session
            gr.contextArray[*handle] = (CONTEXT_SLOT)(sessionIndex + 1);
            return TPM_RC_SUCCESS;
        }
    }
    return TPM_RC_SESSION_HANDLES;
}

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
    )
{
    TPM_RC               result = TPM_RC_SUCCESS;
    CONTEXT_SLOT         slotIndex;
    SESSION             *session = NULL;

    pAssert(sessionType == TPM_SE_HMAC
            || sessionType == TPM_SE_POLICY
            || sessionType == TPM_SE_TRIAL);

    // If there are no open spots in the session array, then no point in searching
    if(s_freeSessionSlots == 0)
        return TPM_RC_SESSION_MEMORY;

    // Find a space for loading a session
    for(slotIndex = 0; slotIndex < MAX_LOADED_SESSIONS; slotIndex++)
    {
        // Is this available?
        if(s_sessions[slotIndex].occupied == FALSE)
        {
            session = &s_sessions[slotIndex].session;
            break;
        }
    }
    // if no spot found, then this is an internal error
    if(slotIndex >= MAX_LOADED_SESSIONS)
        FAIL(FATAL_ERROR_INTERNAL);

    // Call context ID function to get a handle.  TPM_RC_SESSION_HANDLE may be
    // returned from ContextIdHandelAssign()
    result = ContextIdSessionCreate(sessionHandle, slotIndex);
    if(result != TPM_RC_SUCCESS)
        return result;

    //*** Only return from this point on is TPM_RC_SUCCESS

    // Can now indicate that the session array entry is occupied.
    s_freeSessionSlots--;
    s_sessions[slotIndex].occupied = TRUE;

    // Initialize the session data
    MemorySet(session, 0, sizeof(SESSION));

    // Initialize internal session data
    session->authHashAlg = authHash;
    // Initialize session type
    if(sessionType == TPM_SE_HMAC)
    {
        *sessionHandle += HMAC_SESSION_FIRST;
    }
    else
    {
        *sessionHandle += POLICY_SESSION_FIRST;

        // For TPM_SE_POLICY or TPM_SE_TRIAL
        session->attributes.isPolicy = SET;
        if(sessionType == TPM_SE_TRIAL)
            session->attributes.isTrialPolicy = SET;

        SessionSetStartTime(session);

        // Initialize policyDigest.  policyDigest is initialized with a string of 0 
        // of session algorithm digest size. Since the session is already clear.
        // Just need to set the size
        session->u2.policyDigest.t.size = 
            CryptHashGetDigestSize(session->authHashAlg);
    }
    // Create initial session nonce
    session->nonceTPM.t.size = nonceCaller->t.size;
    CryptRandomGenerate(session->nonceTPM.t.size, session->nonceTPM.t.buffer);
    MemoryCopy2B(&nonceTpm->b, &session->nonceTPM.b, 
                 sizeof(nonceTpm->t.buffer));

    // Set up session parameter encryption algorithm
    session->symmetric = *symmetric;

    // If there is a bind object or a session secret, then need to compute
    // a sessionKey.
    if(bind != TPM_RH_NULL || seed->t.size != 0)
    {
        // sessionKey = KDFa(hash, (authValue || seed), "ATH", nonceTPM,
        //                      nonceCaller, bits)
        // The HMAC key for generating the sessionSecret can be the concatenation
        // of an authorization value and a seed value
        TPM2B_TYPE(KEY, (sizeof(TPMT_HA) + sizeof(seed->t.buffer)));
        TPM2B_KEY            key;

        // Get hash size, which is also the length of sessionKey
        session->sessionKey.t.size = CryptHashGetDigestSize(session->authHashAlg);

        // Get authValue of associated entity
        EntityGetAuthValue(bind, (TPM2B_AUTH *)&key);
        pAssert(key.t.size + seed->t.size <= sizeof(key.t.buffer));

        // Concatenate authValue and seed
        MemoryConcat2B(&key.b, &seed->b, sizeof(key.t.buffer));

        // Compute the session key
        CryptKDFa(session->authHashAlg, &key.b, SESSION_KEY, &session->nonceTPM.b,
                  &nonceCaller->b, 
                  session->sessionKey.t.size * 8, session->sessionKey.t.buffer,
                  NULL, FALSE);
    }

    // Copy the name of the entity that the HMAC session is bound to
    // Policy session is not bound to an entity
    if(bind != TPM_RH_NULL && sessionType == TPM_SE_HMAC)
    {
        session->attributes.isBound = SET;
        SessionComputeBoundEntity(bind, &session->u1.boundEntity);
    }
    // If there is a bind object and it is subject to DA, then use of this session
    // is subject to DA regardless of how it is used.
    session->attributes.isDaBound = (bind != TPM_RH_NULL)
        && (IsDAExempted(bind) == FALSE);

// If the session is bound, then check to see if it is bound to lockoutAuth
    session->attributes.isLockoutBound = (session->attributes.isDaBound == SET)
        && (bind == TPM_RH_LOCKOUT);
    return TPM_RC_SUCCESS;
}

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
//      TPM_RC_CONTEXT_GAP              a contextID could not be assigned
//      TPM_RC_TOO_MANY_CONTEXTS        the counter maxed out
//
TPM_RC
SessionContextSave(
    TPM_HANDLE           handle,        // IN: session handle
    CONTEXT_COUNTER     *contextID      // OUT: assigned contextID
    )
{
    UINT32                      contextIndex;
    CONTEXT_SLOT                slotIndex;

    pAssert(SessionIsLoaded(handle));

    // check to see if the gap is already maxed out
    // Need to have a saved session
    if(s_oldestSavedSession < MAX_ACTIVE_SESSIONS
          // if the oldest saved session has the same value as the low bits
          // of the contextCounter, then the GAP is maxed out.
       && gr.contextArray[s_oldestSavedSession] == (CONTEXT_SLOT)gr.contextCounter)
        return TPM_RC_CONTEXT_GAP;

    // if the caller wants the context counter, set it
    if(contextID != NULL)
        *contextID = gr.contextCounter;

    contextIndex = handle & HR_HANDLE_MASK;
    pAssert(contextIndex < MAX_ACTIVE_SESSIONS);

    // Extract the session slot number referenced by the contextArray
    // because we are going to overwrite this with the low order
    // contextID value.
    slotIndex = gr.contextArray[contextIndex] - 1;

    // Set the contextID for the contextArray
    gr.contextArray[contextIndex] = (CONTEXT_SLOT)gr.contextCounter;

    // Increment the counter
    gr.contextCounter++;

    // In the unlikely event that the 64-bit context counter rolls over...
    if(gr.contextCounter == 0)
    {
        // back it up
        gr.contextCounter--;
        // return an error
        return TPM_RC_TOO_MANY_CONTEXTS;
    }
    // if the low-order bits wrapped, need to advance the value to skip over
    // the values used to indicate that a session is loaded
    if(((CONTEXT_SLOT)gr.contextCounter) == 0)
        gr.contextCounter += MAX_LOADED_SESSIONS + 1;

    // If no other sessions are saved, this is now the oldest.
    if(s_oldestSavedSession >= MAX_ACTIVE_SESSIONS)
        s_oldestSavedSession = contextIndex;

    // Mark the session slot as unoccupied
    s_sessions[slotIndex].occupied = FALSE;

    // and indicate that there is an additional open slot
    s_freeSessionSlots++;

    return TPM_RC_SUCCESS;
}

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
    )
{
    UINT32              contextIndex;
    CONTEXT_SLOT        slotIndex;

    pAssert(HandleGetType(*handle) == TPM_HT_POLICY_SESSION
            || HandleGetType(*handle) == TPM_HT_HMAC_SESSION);

    // Don't bother looking if no openings
    if(s_freeSessionSlots == 0)
        return TPM_RC_SESSION_MEMORY;

    // Find a free session slot to load the session
    for(slotIndex = 0; slotIndex < MAX_LOADED_SESSIONS; slotIndex++)
        if(s_sessions[slotIndex].occupied == FALSE) break;

    // if no spot found, then this is an internal error
    pAssert(slotIndex < MAX_LOADED_SESSIONS);

    contextIndex = *handle & HR_HANDLE_MASK;   // extract the index

    // If there is only one slot left, and the gap is at maximum, the only session
    // context that we can safely load is the oldest one.
    if(s_oldestSavedSession < MAX_ACTIVE_SESSIONS
       && s_freeSessionSlots == 1
       && (CONTEXT_SLOT)gr.contextCounter == gr.contextArray[s_oldestSavedSession]
       && contextIndex != s_oldestSavedSession)
        return TPM_RC_CONTEXT_GAP;

    pAssert(contextIndex < MAX_ACTIVE_SESSIONS);

    // set the contextArray value to point to the session slot where
    // the context is loaded
    gr.contextArray[contextIndex] = slotIndex + 1;

    // if this was the oldest context, find the new oldest
    if(contextIndex == s_oldestSavedSession)
        ContextIdSetOldest();

    // Copy session data to session slot
    MemoryCopy(&s_sessions[slotIndex].session, session, sizeof(SESSION));

           // Set session slot as occupied
    s_sessions[slotIndex].occupied = TRUE;

    // Reduce the number of open spots
    s_freeSessionSlots--;

    return TPM_RC_SUCCESS;
}

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
    )
{
    CONTEXT_SLOT         slotIndex;
    UINT32               contextIndex;   // Index into contextArray

    pAssert((HandleGetType(handle) == TPM_HT_POLICY_SESSION
             || HandleGetType(handle) == TPM_HT_HMAC_SESSION
             )
            && (SessionIsLoaded(handle) || SessionIsSaved(handle))
            );

         // Flush context ID of this session
         // Convert handle to an index into the contextArray
    contextIndex = handle & HR_HANDLE_MASK;

    pAssert(contextIndex < sizeof(gr.contextArray) / sizeof(gr.contextArray[0]));

    // Get the current contents of the array
    slotIndex = gr.contextArray[contextIndex];

    // Mark context array entry as available
    gr.contextArray[contextIndex] = 0;

    // Is this a saved session being flushed
    if(slotIndex > MAX_LOADED_SESSIONS)
    {
        // Flushing the oldest session?
        if(contextIndex == s_oldestSavedSession)
            // If so, find a new value for oldest.
            ContextIdSetOldest();
    }
    else
    {
        // Adjust slot index to point to session array index
        slotIndex -= 1;

        // Free session array index
        s_sessions[slotIndex].occupied = FALSE;
        s_freeSessionSlots++;
    }

    return;
}

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
    )
{
    TPM2B_AUTH           auth;
    BYTE                *pAuth = auth.t.buffer;
    UINT16               i;

    // Get name
    EntityGetName(entityHandle, bind);

//    // The bound value of a reserved handle is the handle itself
//    if(bind->t.size == sizeof(TPM_HANDLE)) return;

    // For all the other entities, concatenate the authorization value to the name.
    // Get a local copy of the authorization value because some overlapping
    // may be necessary.
    EntityGetAuthValue(entityHandle, &auth);

    // Make sure that the extra space is zeroed
    MemorySet(&bind->t.name[bind->t.size], 0, sizeof(bind->t.name) - bind->t.size);
    // XOR the authValue at the end of the name
    for(i = sizeof(bind->t.name) - auth.t.size; i < sizeof(bind->t.name); i++)
        bind->t.name[i] ^= *pAuth++;

    // Set the bind value to the maximum size
    bind->t.size = sizeof(bind->t.name);

    return;
}


//*** SessionSetStartTime()
// This function is used to initialize the session timing 
void
SessionSetStartTime(
    SESSION         *session        // IN: the session to update
    )
{
    session->startTime = g_time;
    session->epoch = g_timeEpoch;
    session->timeout = 0;
}

//*** SessionResetPolicyData()
// This function is used to reset the policy data without changing the nonce
// or the start time of the session.
void
SessionResetPolicyData(
    SESSION         *session        // IN: the session to reset
    )
{
    SESSION_ATTRIBUTES      oldAttributes;
    pAssert(session != NULL);

    // Will need later
    oldAttributes = session->attributes;

    // No command
    session->commandCode = 0;

    // No locality selected
    MemorySet(&session->commandLocality, 0, sizeof(session->commandLocality));

    // The cpHash size to zero
    session->u1.cpHash.b.size = 0;

    // No timeout
    session->timeout = 0;

    // Reset the pcrCounter
    session->pcrCounter = 0;

    // Reset the policy hash
    MemorySet(&session->u2.policyDigest.t.buffer, 0,
              session->u2.policyDigest.t.size);

    // Reset the session attributes
    MemorySet(&session->attributes, 0, sizeof(SESSION_ATTRIBUTES));

    // Restore the policy attributes
    session->attributes.isPolicy = SET;
    session->attributes.isTrialPolicy = oldAttributes.isTrialPolicy;

    // Restore the bind attributes
    session->attributes.isDaBound = oldAttributes.isDaBound;
    session->attributes.isLockoutBound = oldAttributes.isLockoutBound;
}

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
    )
{
    TPMI_YES_NO     more = NO;
    UINT32          i;

    pAssert(HandleGetType(handle) == TPM_HT_LOADED_SESSION);

    // Initialize output handle list
    handleList->count = 0;

    // The maximum count of handles we may return is MAX_CAP_HANDLES
    if(count > MAX_CAP_HANDLES) count = MAX_CAP_HANDLES;

    // Iterate session context ID slots to get loaded session handles
    for(i = handle & HR_HANDLE_MASK; i < MAX_ACTIVE_SESSIONS; i++)
    {
        // If session is active
        if(gr.contextArray[i] != 0)
        {
            // If session is loaded
            if(gr.contextArray[i] <= MAX_LOADED_SESSIONS)
            {
                if(handleList->count < count)
                {
                    SESSION         *session;

                    // If we have not filled up the return list, add this
                    // session handle to it
                    // assume that this is going to be an HMAC session
                    handle = i + HMAC_SESSION_FIRST;
                    session = SessionGet(handle);
                    if(session->attributes.isPolicy)
                        handle = i + POLICY_SESSION_FIRST;
                    handleList->handle[handleList->count] = handle;
                    handleList->count++;
                }
                else
                {
                    // If the return list is full but we still have loaded object
                    // available, report this and stop iterating
                    more = YES;
                    break;
                }
            }
        }
    }

    return more;
}

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
    )
{
    TPMI_YES_NO     more = NO;
    UINT32          i;

#ifdef  TPM_HT_SAVED_SESSION
    pAssert(HandleGetType(handle) == TPM_HT_SAVED_SESSION);
#else
    pAssert(HandleGetType(handle) == TPM_HT_ACTIVE_SESSION);
#endif

    // Initialize output handle list
    handleList->count = 0;

    // The maximum count of handles we may return is MAX_CAP_HANDLES
    if(count > MAX_CAP_HANDLES) count = MAX_CAP_HANDLES;

    // Iterate session context ID slots to get loaded session handles
    for(i = handle & HR_HANDLE_MASK; i < MAX_ACTIVE_SESSIONS; i++)
    {
        // If session is active
        if(gr.contextArray[i] != 0)
        {
            // If session is saved
            if(gr.contextArray[i] > MAX_LOADED_SESSIONS)
            {
                if(handleList->count < count)
                {
                    // If we have not filled up the return list, add this
                    // session handle to it
                    handleList->handle[handleList->count] = i + HMAC_SESSION_FIRST;
                    handleList->count++;
                }
                else
                {
                    // If the return list is full but we still have loaded object
                    // available, report this and stop iterating
                    more = YES;
                    break;
                }
            }
        }
    }

    return more;
}

//*** SessionCapGetLoadedNumber()
// This function return the number of authorization sessions currently
// loaded into TPM RAM.
UINT32
SessionCapGetLoadedNumber(
    void
    )
{
    return MAX_LOADED_SESSIONS - s_freeSessionSlots;
}

//*** SessionCapGetLoadedAvail()
// This function returns the number of additional authorization sessions, of
// any type, that could be loaded into TPM RAM.
// NOTE: In other implementations, this number may just be an estimate. The only
//       requirement for the estimate is, if it is one or more, then at least one
//       session must be loadable.
UINT32
SessionCapGetLoadedAvail(
    void
    )
{
    return s_freeSessionSlots;
}

//*** SessionCapGetActiveNumber()
// This function returns the number of active authorization sessions currently
// being tracked by the TPM.
UINT32
SessionCapGetActiveNumber(
    void
    )
{
    UINT32              i;
    UINT32              num = 0;

    // Iterate the context array to find the number of non-zero slots
    for(i = 0; i < MAX_ACTIVE_SESSIONS; i++)
    {
        if(gr.contextArray[i] != 0) num++;
    }

    return num;
}

//*** SessionCapGetActiveAvail()
// This function returns the number of additional authorization sessions, of any
// type, that could be created. This not the number of slots for sessions, but
// the number of additional sessions that the TPM is capable of tracking.
UINT32
SessionCapGetActiveAvail(
    void
    )
{
    UINT32              i;
    UINT32              num = 0;

    // Iterate the context array to find the number of zero slots
    for(i = 0; i < MAX_ACTIVE_SESSIONS; i++)
    {
        if(gr.contextArray[i] == 0) num++;
    }

    return num;
}