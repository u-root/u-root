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

#ifndef    _OBJECT_FP_H_
#define    _OBJECT_FP_H_

//*** ObjectFlush()
// This function marks an object slot as available.
// Since there is no checking of the input parameters, it should be used
// judiciously.
// Note: This could be converted to a macro.
void
ObjectFlush(
    OBJECT          *object
);

//*** ObjectSetInUse()
// This access function sets the occupied attribute of an object slot.
void
ObjectSetInUse(
    OBJECT          *object
);

//*** ObjectStartup()
// This function is called at TPM2_Startup() to initialize the object subsystem.
BOOL
ObjectStartup(
    void
);

//*** ObjectCleanupEvict()
//
// In this implementation, a persistent object is moved from NV into an object slot
// for processing. It is flushed after command execution. This function is called
// from ExecuteCommand().
void
ObjectCleanupEvict(
    void
);

//*** IsObjectPresent()
// This function checks to see if a transient handle references a loaded
// object.  This routine should not be called if the handle is not a
// transient handle. The function validates that the handle is in the
// implementation-dependent allowed in range for loaded transient objects.
//  Return Type: BOOL
//      TRUE(1)         handle references a loaded object
//      FALSE(0)        handle is not an object handle, or it does not
//                      reference to a loaded object
BOOL
IsObjectPresent(
    TPMI_DH_OBJECT   handle         // IN: handle to be checked
);

//*** ObjectIsSequence()
// This function is used to check if the object is a sequence object. This function
// should not be called if the handle does not reference a loaded object.
//  Return Type: BOOL
//      TRUE(1)         object is an HMAC, hash, or event sequence object
//      FALSE(0)        object is not an HMAC, hash, or event sequence object
BOOL
ObjectIsSequence(
    OBJECT          *object         // IN: handle to be checked
);

//*** HandleToObject()
// This function is used to find the object structure associated with a handle.
//
// This function requires that 'handle' references a loaded object or a permanent
// handle.
OBJECT*
HandleToObject(
    TPMI_DH_OBJECT   handle         // IN: handle of the object
);

//*** GetQualifiedName()
// This function returns the Qualified Name of the object. In this implementation,
// the Qualified Name is computed when the object is loaded and is saved in the
// internal representation of the object. The alternative would be to retain the
// Name of the parent and compute the QN when needed. This would take the same
// amount of space so it is not recommended that the alternate be used.
//
// This function requires that 'handle' references a loaded object.
void
GetQualifiedName(
    TPMI_DH_OBJECT   handle,        // IN: handle of the object
    TPM2B_NAME      *qualifiedName  // OUT: qualified name of the object
);

//*** ObjectGetHierarchy()
// This function returns the handle for the hierarchy of an object.
TPMI_RH_HIERARCHY
ObjectGetHierarchy(
    OBJECT          *object         // IN :object
);

//*** GetHeriarchy()
// This function returns the handle of the hierarchy to which a handle belongs.
// This function is similar to ObjectGetHierarchy() but this routine takes
// a handle but ObjectGetHierarchy() takes an pointer to an object.
//
// This function requires that 'handle' references a loaded object.
TPMI_RH_HIERARCHY
GetHeriarchy(
    TPMI_DH_OBJECT   handle         // IN :object handle
);

//*** FindEmptyObjectSlot()
// This function finds an open object slot, if any. It will clear the attributes
// but will not set the occupied attribute. This is so that a slot may be used
// and discarded if everything does not go as planned.
//  Return Type: OBJECT *
//      NULL        no open slot found
//      != NULL     pointer to available slot
OBJECT *
FindEmptyObjectSlot(
    TPMI_DH_OBJECT  *handle         // OUT: (optional)
);

//*** ObjectAllocateSlot()
// This function is used to allocate a slot in internal object array.
OBJECT *
ObjectAllocateSlot(
    TPMI_DH_OBJECT  *handle        // OUT: handle of allocated object
);

//*** ObjectSetLoadedAttributes()
// This function sets the internal attributes for a loaded object. It is called to
// finalize the OBJECT attributes (not the TPMA_OBJECT attributes) for a loaded
// object.
void
ObjectSetLoadedAttributes(
    OBJECT          *object,        // IN: object attributes to finalize
    TPM_HANDLE       parentHandle   // IN: the parent handle
);

//*** ObjectLoad()
// Common function to load an object. A loaded object has its public area validated
// (unless its 'nameAlg' is TPM_ALG_NULL). If a sensitive part is loaded, it is
// verified to be correct and if both public and sensitive parts are loaded, then
// the cryptographic binding between the objects is validated. This function does
// not cause the allocated slot to be marked as in use.
TPM_RC
ObjectLoad(
    OBJECT          *object,        // IN: pointer to object slot
                                    //     object
    OBJECT          *parent,        // IN: (optional) the parent object
    TPMT_PUBLIC     *publicArea,    // IN: public area to be installed in the object
    TPMT_SENSITIVE  *sensitive,     // IN: (optional) sensitive area to be
                                    //      installed in the object
    TPM_RC           blamePublic,   // IN: parameter number to associate with the
                                    //     publicArea errors
    TPM_RC           blameSensitive,// IN: parameter number to associate with the
                                    //     sensitive area errors
    TPM2B_NAME      *name           // IN: (optional)
);

#if CC_HMAC_Start || CC_MAC_Start
//*** ObjectCreateHMACSequence()
// This function creates an internal HMAC sequence object.
//  Return Type: TPM_RC
//      TPM_RC_OBJECT_MEMORY        if there is no free slot for an object
TPM_RC
ObjectCreateHMACSequence(
    TPMI_ALG_HASH    hashAlg,       // IN: hash algorithm
    OBJECT          *keyObject,     // IN: the object containing the HMAC key
    TPM2B_AUTH      *auth,          // IN: authValue
    TPMI_DH_OBJECT  *newHandle      // OUT: HMAC sequence object handle
);
#endif

//*** ObjectCreateHashSequence()
// This function creates a hash sequence object.
//  Return Type: TPM_RC
//      TPM_RC_OBJECT_MEMORY        if there is no free slot for an object
TPM_RC
ObjectCreateHashSequence(
    TPMI_ALG_HASH    hashAlg,       // IN: hash algorithm
    TPM2B_AUTH      *auth,          // IN: authValue
    TPMI_DH_OBJECT  *newHandle      // OUT: sequence object handle
);

//*** ObjectCreateEventSequence()
// This function creates an event sequence object.
//  Return Type: TPM_RC
//      TPM_RC_OBJECT_MEMORY        if there is no free slot for an object
TPM_RC
ObjectCreateEventSequence(
    TPM2B_AUTH      *auth,          // IN: authValue
    TPMI_DH_OBJECT  *newHandle      // OUT: sequence object handle
);

//*** ObjectTerminateEvent()
// This function is called to close out the event sequence and clean up the hash
// context states.
void
ObjectTerminateEvent(
    void
);

//*** ObjectContextLoad()
// This function loads an object from a saved object context.
//  Return Type: OBJECT *
//      NULL        if there is no free slot for an object
//      != NULL     points to the loaded object
OBJECT *
ObjectContextLoad(
    ANY_OBJECT_BUFFER   *object,        // IN: pointer to object structure in saved
                                        //     context
    TPMI_DH_OBJECT      *handle         // OUT: object handle
);

//*** FlushObject()
// This function frees an object slot.
//
// This function requires that the object is loaded.
void
FlushObject(
    TPMI_DH_OBJECT   handle         // IN: handle to be freed
);

//*** ObjectFlushHierarchy()
// This function is called to flush all the loaded transient objects associated
// with a hierarchy when the hierarchy is disabled.
void
ObjectFlushHierarchy(
    TPMI_RH_HIERARCHY    hierarchy      // IN: hierarchy to be flush
);

//*** ObjectLoadEvict()
// This function loads a persistent object into a transient object slot.
//
// This function requires that 'handle' is associated with a persistent object.
//  Return Type: TPM_RC
//      TPM_RC_HANDLE               the persistent object does not exist
//                                  or the associated hierarchy is disabled.
//      TPM_RC_OBJECT_MEMORY        no object slot
TPM_RC
ObjectLoadEvict(
    TPM_HANDLE      *handle,        // IN:OUT: evict object handle.  If success, it
                                    // will be replace by the loaded object handle
    COMMAND_INDEX    commandIndex   // IN: the command being processed
);

//*** ObjectComputeName()
// This does the name computation from a public area (can be marshaled or not).
TPM2B_NAME *
ObjectComputeName(
    UINT32           size,          // IN: the size of the area to digest
    BYTE            *publicArea,    // IN: the public area to digest
    TPM_ALG_ID       nameAlg,       // IN: the hash algorithm to use
    TPM2B_NAME      *name           // OUT: Computed name
);

//*** PublicMarshalAndComputeName()
// This function computes the Name of an object from its public area.
TPM2B_NAME *
PublicMarshalAndComputeName(
    TPMT_PUBLIC     *publicArea,    // IN: public area of an object
    TPM2B_NAME      *name           // OUT: name of the object
);

//*** ComputeQualifiedName()
// This function computes the qualified name of an object.
void
ComputeQualifiedName(
    TPM_HANDLE       parentHandle,  // IN: parent's handle
    TPM_ALG_ID       nameAlg,       // IN: name hash
    TPM2B_NAME      *name,          // IN: name of the object
    TPM2B_NAME      *qualifiedName  // OUT: qualified name of the object
);

//*** ObjectIsStorage()
// This function determines if an object has the attributes associated
// with a parent. A parent is an asymmetric or symmetric block cipher key
// that has its 'restricted' and 'decrypt' attributes SET, and 'sign' CLEAR.
//  Return Type: BOOL
//      TRUE(1)         object is a storage key
//      FALSE(0)        object is not a storage key
BOOL
ObjectIsStorage(
    TPMI_DH_OBJECT   handle         // IN: object handle
);

//*** ObjectCapGetLoaded()
// This function returns a a list of handles of loaded object, starting from
// 'handle'. 'Handle' must be in the range of valid transient object handles,
// but does not have to be the handle of a loaded transient object.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
ObjectCapGetLoaded(
    TPMI_DH_OBJECT   handle,        // IN: start handle
    UINT32           count,         // IN: count of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
);

//*** ObjectCapGetTransientAvail()
// This function returns an estimate of the number of additional transient
// objects that could be loaded into the TPM.
UINT32
ObjectCapGetTransientAvail(
    void
);

//*** ObjectGetPublicAttributes()
// Returns the attributes associated with an object handles.
TPMA_OBJECT
ObjectGetPublicAttributes(
    TPM_HANDLE       handle
);

OBJECT_ATTRIBUTES
ObjectGetProperties(
    TPM_HANDLE       handle
);

#endif  // _OBJECT_FP_H_
