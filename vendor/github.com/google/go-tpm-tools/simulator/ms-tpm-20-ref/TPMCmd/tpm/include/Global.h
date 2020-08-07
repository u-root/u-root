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

// This file contains internal global type definitions and data declarations that
// are need between subsystems. The instantiation of global data is in Global.c.
// The initialization of global data is in the subsystem that is the primary owner
// of the data.
//
// The first part of this file has the typedefs for structures and other defines
// used in many portions of the code. After the typedef section, is a section that
// defines global values that are only present in RAM. The next three sections
// define the structures for the NV data areas: persistent, orderly, and state
// save. Additional sections define the data that is used in specific modules. That
// data is private to the module but is collected here to simplify the management
// of the instance data.
// All the data is instanced in Global.c.
#if !defined _TPM_H_
#error "Should only be instanced in TPM.h"
#endif


//** Includes

#ifndef         GLOBAL_H
#define         GLOBAL_H

#ifdef GLOBAL_C
#define EXTERN
#define INITIALIZER(_value_)  = _value_
#else
#define EXTERN  extern
#define INITIALIZER(_value_)
#endif

_REDUCE_WARNING_LEVEL_(2)
#include <string.h>
#include <stddef.h>
_NORMAL_WARNING_LEVEL_

#if SIMULATION
#undef CONTEXT_SLOT
#  define CONTEXT_SLOT    UINT8
#endif
#include "Capabilities.h"
#include "TpmTypes.h"
#include "CommandAttributes.h"
#include "CryptTest.h"
#include "BnValues.h"
#include "CryptHash.h"
#include "CryptSym.h"
#include "CryptRand.h"
#include "CryptEcc.h"
#include "CryptRsa.h"
#include "CryptTest.h"
#include "TpmError.h"
#include "NV.h"

//** Defines and Types

//*** Size Types
// These types are used to differentiate the two different size values used.
//
// NUMBYTES is used when a size is a number of bytes (usually a TPM2B)
typedef UINT16  NUMBYTES;

//*** Other Types
// An AUTH_VALUE is a BYTE array containing a digest (TPMU_HA)
typedef BYTE    AUTH_VALUE[sizeof(TPMU_HA)];

// A TIME_INFO is a BYTE array that can contain a TPMS_TIME_INFO
typedef BYTE    TIME_INFO[sizeof(TPMS_TIME_INFO)];

// A NAME is a BYTE array that can contain a TPMU_NAME
typedef BYTE    NAME[sizeof(TPMU_NAME)];

// Definition for a PROOF value
TPM2B_TYPE(PROOF, PROOF_SIZE);

// Definition for a Primary Seed value
TPM2B_TYPE(SEED, PRIMARY_SEED_SIZE);


// A CLOCK_NONCE is used to tag the time value in the authorization session and
// in the ticket computation so that the ticket expires when there is a time
// discontinuity. When the clock stops during normal operation, the nonce is
// 64-bit value kept in RAM but it is a 32-bit counter when the clock only stops
// during power events.
#if CLOCK_STOPS
typedef UINT64          CLOCK_NONCE;
#else
typedef UINT32          CLOCK_NONCE;
#endif

//** Loaded Object Structures
//*** Description
// The structures in this section define the object layout as it exists in TPM
// memory.
//
// Two types of objects are defined: an ordinary object such as a key, and a
// sequence object that may be a hash, HMAC, or event.
//
//*** OBJECT_ATTRIBUTES
// An OBJECT_ATTRIBUTES structure contains the variable attributes of an object.
// These properties are not part of the public properties but are used by the
// TPM in managing the object. An OBJECT_ATTRIBUTES is used in the definition of
// the OBJECT data type.

typedef struct
{
    unsigned            publicOnly : 1;     //0) SET if only the public portion of
                                            //   an object is loaded
    unsigned            epsHierarchy : 1;   //1) SET if the object belongs to EPS
                                            //   Hierarchy
    unsigned            ppsHierarchy : 1;   //2) SET if the object belongs to PPS
                                            //   Hierarchy
    unsigned            spsHierarchy : 1;   //3) SET f the object belongs to SPS
                                            //   Hierarchy
    unsigned            evict : 1;          //4) SET if the object is a platform or
                                            //   owner evict object.  Platform-
                                            //   evict object belongs to PPS
                                            //   hierarchy, owner-evict object
                                            //   belongs to SPS or EPS hierarchy.
                                            //   This bit is also used to mark a
                                            //   completed sequence object so it
                                            //   will be flush when the
                                            //   SequenceComplete command succeeds.
    unsigned            primary : 1;        //5) SET for a primary object
    unsigned            temporary : 1;      //6) SET for a temporary object
    unsigned            stClear : 1;        //7) SET for an stClear object
    unsigned            hmacSeq : 1;        //8) SET for an HMAC or MAC sequence 
                                            //   object
    unsigned            hashSeq : 1;        //9) SET for a hash sequence object
    unsigned            eventSeq : 1;       //10) SET for an event sequence object
    unsigned            ticketSafe : 1;     //11) SET if a ticket is safe to create
                                            //    for hash sequence object
    unsigned            firstBlock : 1;     //12) SET if the first block of hash
                                            //    data has been received.  It
                                            //    works with ticketSafe bit
    unsigned            isParent : 1;       //13) SET if the key has the proper
                                            //    attributes to be a parent key
//   unsigned            privateExp : 1;    //14) SET when the private exponent
//                                          //    of an RSA key has been validated.
    unsigned            not_used_14 : 1;
    unsigned            occupied : 1;       //15) SET when the slot is occupied.
    unsigned            derivation : 1;     //16) SET when the key is a derivation
                                            //        parent
    unsigned            external : 1;       //17) SET when the object is loaded with
                                            //    TPM2_LoadExternal();
} OBJECT_ATTRIBUTES;

#if ALG_RSA
// There is an overload of the sensitive.rsa.t.size field of a TPMT_SENSITIVE when an 
// RSA key is loaded. When the sensitive->sensitive contains an RSA key with all of 
// the CRT values, then the MSB of the size field will be set to indicate that the
// buffer contains all 5 of the CRT private key values.
#define     RSA_prime_flag      0x8000
#endif


//*** OBJECT Structure
// An OBJECT structure holds the object public, sensitive, and meta-data
// associated. This structure is implementation dependent. For this
// implementation, the structure is not optimized for space but rather
// for clarity of the reference implementation. Other implementations
// may choose to overlap portions of the structure that are not used
// simultaneously. These changes would necessitate changes to the source
// code but those changes would be compatible with the reference
// implementation.

typedef struct OBJECT
{
    // The attributes field is required to be first followed by the publicArea.
    // This allows the overlay of the object structure and a sequence structure
    OBJECT_ATTRIBUTES   attributes;         // object attributes
    TPMT_PUBLIC         publicArea;         // public area of an object
    TPMT_SENSITIVE      sensitive;          // sensitive area of an object
    TPM2B_NAME          qualifiedName;      // object qualified name
    TPMI_DH_OBJECT      evictHandle;        // if the object is an evict object,
                                            // the original handle is kept here.
                                            // The 'working' handle will be the
                                            // handle of an object slot.
    TPM2B_NAME          name;               // Name of the object name. Kept here
                                            // to avoid repeatedly computing it.
} OBJECT;

//*** HASH_OBJECT Structure
// This structure holds a hash sequence object or an event sequence object.
//
// The first four components of this structure are manually set to be the same as
// the first four components of the object structure. This prevents the object
// from being inadvertently misused as sequence objects occupy the same memory as
// a regular object. A debug check is present to make sure that the offsets are
// what they are supposed to be.
// NOTE: In a future version, this will probably be renamed as SEQUENCE_OBJECT
typedef struct HASH_OBJECT
{
    OBJECT_ATTRIBUTES   attributes;         // The attributes of the HASH object
    TPMI_ALG_PUBLIC     type;               // algorithm
    TPMI_ALG_HASH       nameAlg;            // name algorithm
    TPMA_OBJECT         objectAttributes;   // object attributes

    // The data below is unique to a sequence object
    TPM2B_AUTH          auth;               // authorization for use of sequence
    union
    {
        HASH_STATE      hashState[HASH_COUNT];
        HMAC_STATE      hmacState;
    }                   state;
} HASH_OBJECT;

typedef BYTE  HASH_OBJECT_BUFFER[sizeof(HASH_OBJECT)];

//*** ANY_OBJECT
// This is the union for holding either a sequence object or a regular object.
// for ContextSave and ContextLoad
typedef union ANY_OBJECT
{
    OBJECT              entity;
    HASH_OBJECT         hash;
} ANY_OBJECT;

typedef BYTE    ANY_OBJECT_BUFFER[sizeof(ANY_OBJECT)];

//**AUTH_DUP Types
// These values are used in the authorization processing.

typedef UINT32          AUTH_ROLE;
#define AUTH_NONE       ((AUTH_ROLE)(0))
#define AUTH_USER       ((AUTH_ROLE)(1))
#define AUTH_ADMIN      ((AUTH_ROLE)(2))
#define AUTH_DUP        ((AUTH_ROLE)(3))

//** Active Session Context
//*** Description
// The structures in this section define the internal structure of a session
// context.
//
//*** SESSION_ATTRIBUTES
// The attributes in the SESSION_ATTRIBUTES structure track the various properties
// of the session. It maintains most of the tracking state information for the
// policy session. It is used within the SESSION structure.

typedef struct SESSION_ATTRIBUTES
{
    unsigned    isPolicy : 1;           //1) SET if the session may only be used 
                                        //   for policy
    unsigned    isAudit : 1;            //2) SET if the session is used for audit
    unsigned    isBound : 1;            //3) SET if the session is bound to with an
                                        //   entity. This attribute will be CLEAR 
                                        //   if either isPolicy or isAudit is SET.
    unsigned    isCpHashDefined : 1;    //3) SET if the cpHash has been defined 
                                        //   This attribute is not SET unless
                                        //   'isPolicy' is SET.
    unsigned    isAuthValueNeeded : 1;  //5) SET if the authValue is required for 
                                        //   computing the session HMAC. This 
                                        //   attribute is not SET unless 'isPolicy'
                                        //   is SET.
    unsigned    isPasswordNeeded : 1;   //6) SET if a password authValue is required
                                        //   for authorization This attribute is not
                                        //   SET unless 'isPolicy' is SET.
    unsigned    isPPRequired : 1;       //7) SET if physical presence is required to
                                        //   be asserted when the authorization is 
                                        //   checked. This attribute is not SET 
                                        //   unless 'isPolicy' is SET.
    unsigned    isTrialPolicy : 1;      //8) SET if the policy session is created 
                                        //   for trial of the policy's policyHash 
                                        //   generation. This attribute is not SET 
                                        //   unless 'isPolicy' is SET.
    unsigned    isDaBound : 1;          //9) SET if the bind entity had noDA CLEAR. 
                                        //   If this is SET, then an authorization 
                                        //   failure using this session will count 
                                        //   against lockout even if the object 
                                        //   being authorized is exempt from DA.
    unsigned    isLockoutBound : 1;     //10) SET if the session is bound to 
                                        //    lockoutAuth.
    unsigned    includeAuth : 1;        //11) This attribute is SET when the
                                        //    authValue of an object is to be
                                        //    included in the computation of the
                                        //    HMAC key for the command and response 
                                        //    computations. (was 'requestWasBound')
    unsigned    checkNvWritten : 1;     //12) SET if the TPMA_NV_WRITTEN attribute 
                                        //    needs to be checked when the policy is
                                        //    used for authorization for NV access.
                                        //    If this is SET for any other type, the
                                        //    policy will fail.
    unsigned    nvWrittenState : 1;     //13) SET if TPMA_NV_WRITTEN is required to 
                                        //    be SET. Used when 'checkNvWritten' is
                                        //    SET
    unsigned    isTemplateSet : 1;      //14) SET if the templateHash needs to be 
                                        //    checked for Create, CreatePrimary, or
                                        //    CreateLoaded.
} SESSION_ATTRIBUTES;

//*** SESSION Structure
// The SESSION structure contains all the context of a session except for the
// associated contextID.
//
// Note: The contextID of a session is only relevant when the session context
// is stored off the TPM.

typedef struct SESSION
{
    SESSION_ATTRIBUTES  attributes;         // session attributes
    UINT32              pcrCounter;         // PCR counter value when PCR is
                                            // included (policy session)
                                            // If no PCR is included, this
                                            // value is 0.
    UINT64              startTime;          // The value in g_time when the session 
                                            // was started (policy session)
    UINT64              timeout;            // The timeout relative to g_time
                                            // There is no timeout if this value
                                            // is 0.
    CLOCK_NONCE         epoch;              // The g_clockEpoch value when the
                                            // session was started. If g_clockEpoch
                                            // does not match this value when the
                                            // timeout is used, then 
                                            // then the command will fail.
    TPM_CC              commandCode;        // command code (policy session)
    TPM_ALG_ID          authHashAlg;        // session hash algorithm
    TPMA_LOCALITY       commandLocality;    // command locality (policy session)
    TPMT_SYM_DEF        symmetric;          // session symmetric algorithm (if any)
    TPM2B_AUTH          sessionKey;         // session secret value used for
                                            // this session
    TPM2B_NONCE         nonceTPM;           // last TPM-generated nonce for
                                            // generating HMAC and encryption keys
   union
    {
        TPM2B_NAME      boundEntity;        // value used to track the entity to
                                            // which the session is bound

        TPM2B_DIGEST    cpHash;             // the required cpHash value for the
                                            // command being authorized
        TPM2B_DIGEST    nameHash;           // the required nameHash
        TPM2B_DIGEST    templateHash;       // the required template for creation
    } u1;

    union
    {
        TPM2B_DIGEST    auditDigest;        // audit session digest
        TPM2B_DIGEST    policyDigest;       // policyHash
    } u2;                                   // audit log and policyHash may
                                            // share space to save memory
} SESSION;

#define     EXPIRES_ON_RESET    INT32_MIN
#define     TIMEOUT_ON_RESET    UINT64_MAX
#define     EXPIRES_ON_RESTART  (INT32_MIN + 1)
#define     TIMEOUT_ON_RESTART  (UINT64_MAX - 1)

typedef BYTE        SESSION_BUF[sizeof(SESSION)];

//*********************************************************************************
//** PCR
//*********************************************************************************
//***PCR_SAVE Structure
// The PCR_SAVE structure type contains the PCR data that are saved across power
// cycles. Only the static PCR are required to be saved across power cycles. The
// DRTM and resettable PCR are not saved. The number of static and resettable PCR
// is determined by the platform-specific specification to which the TPM is built.

typedef struct PCR_SAVE
{
#if     ALG_SHA1
    BYTE                sha1[NUM_STATIC_PCR][SHA1_DIGEST_SIZE];
#endif
#if     ALG_SHA256
    BYTE                sha256[NUM_STATIC_PCR][SHA256_DIGEST_SIZE];
#endif
#if     ALG_SHA384
    BYTE                sha384[NUM_STATIC_PCR][SHA384_DIGEST_SIZE];
#endif
#if     ALG_SHA512
    BYTE                sha512[NUM_STATIC_PCR][SHA512_DIGEST_SIZE];
#endif
#if     ALG_SM3_256
    BYTE                sm3_256[NUM_STATIC_PCR][SM3_256_DIGEST_SIZE];
#endif

    // This counter increments whenever the PCR are updated.
    // NOTE: A platform-specific specification may designate
    //       certain PCR changes as not causing this counter
    //       to increment.
    UINT32              pcrCounter;
} PCR_SAVE;

//***PCR_POLICY
#if defined NUM_POLICY_PCR_GROUP && NUM_POLICY_PCR_GROUP > 0
// This structure holds the PCR policies, one for each group of PCR controlled
// by policy.
typedef struct PCR_POLICY
{
    TPMI_ALG_HASH       hashAlg[NUM_POLICY_PCR_GROUP];
    TPM2B_DIGEST        a;
    TPM2B_DIGEST        policy[NUM_POLICY_PCR_GROUP];
} PCR_POLICY;
#endif

//***PCR_AUTHVALUE
// This structure holds the PCR policies, one for each group of PCR controlled
// by policy.
typedef struct PCR_AUTH_VALUE
{
    TPM2B_DIGEST        auth[NUM_AUTHVALUE_PCR_GROUP];
} PCR_AUTHVALUE;



//**STARTUP_TYPE
// This enumeration is the possible startup types. The type is determined
// by the combination of TPM2_ShutDown and TPM2_Startup.
typedef enum
{
    SU_RESET,
    SU_RESTART,
    SU_RESUME
} STARTUP_TYPE;

//**NV

//***NV_INDEX
// The NV_INDEX structure defines the internal format for an NV index.
// The 'indexData' size varies according to the type of the index.
// In this implementation, all of the index is manipulated as a unit.
typedef struct NV_INDEX
{
    TPMS_NV_PUBLIC      publicArea;
    TPM2B_AUTH          authValue;
} NV_INDEX;

//*** NV_REF
// An NV_REF is an opaque value returned by the NV subsystem. It is used to
// reference and NV Index in a relatively efficient way. Rather than having to
// continually search for an Index, its reference value may be used. In this
// implementation, an NV_REF is a byte pointer that points to the copy of the
// NV memory that is kept in RAM.
typedef UINT32           NV_REF;

typedef BYTE            *NV_RAM_REF;
//***NV_PIN
// This structure deals with the possible endianess differences between the
// canonical form of the TPMS_NV_PIN_COUNTER_PARAMETERS structure and the internal
// value. The structures allow the data in a PIN index to be read as an 8-octet
// value using NvReadUINT64Data(). That function will byte swap all the values on a
// little endian system. This will put the bytes with the 4-octet values in the
// correct order but will swap the pinLimit and pinCount values. When written, the
// PIN index is simply handled as a normal index with the octets in canonical order.
#if BIG_ENDIAN_TPM
typedef struct
{
    UINT32      pinCount;
    UINT32      pinLimit;
} PIN_DATA;
#else
typedef struct
{
    UINT32      pinLimit;
    UINT32      pinCount;
} PIN_DATA;
#endif

typedef union
{
    UINT64     intVal;
    PIN_DATA   pin;
} NV_PIN;

//**COMMIT_INDEX_MASK
// This is the define for the mask value that is used when manipulating
// the bits in the commit bit array. The commit counter is a 64-bit
// value and the low order bits are used to index the commitArray.
// This mask value is applied to the commit counter to extract the
// bit number in the array.
#if     ALG_ECC

#define COMMIT_INDEX_MASK ((UINT16)((sizeof(gr.commitArray)*8)-1))

#endif

//*****************************************************************************
//*****************************************************************************
//** RAM Global Values
//*****************************************************************************
//*****************************************************************************
//*** Description
// The values in this section are only extant in RAM or ROM as constant values.

//*** Crypto Self-Test Values
EXTERN ALGORITHM_VECTOR     g_implementedAlgorithms;
EXTERN ALGORITHM_VECTOR     g_toTest;

//*** g_rcIndex[]
// This array is used to contain the array of values that are added to a return
// code when it is a parameter-, handle-, or session-related error.
// This is an implementation choice and the same result can be achieved by using
// a macro.
#define g_rcIndexInitializer {  TPM_RC_1, TPM_RC_2, TPM_RC_3, TPM_RC_4,             \
                                TPM_RC_5, TPM_RC_6, TPM_RC_7, TPM_RC_8,             \
                                TPM_RC_9, TPM_RC_A, TPM_RC_B, TPM_RC_C,             \
                                TPM_RC_D, TPM_RC_E, TPM_RC_F }
EXTERN const UINT16     g_rcIndex[15] INITIALIZER(g_rcIndexInitializer);

//*** g_exclusiveAuditSession
// This location holds the session handle for the current exclusive audit
// session. If there is no exclusive audit session, the location is set to
// TPM_RH_UNASSIGNED.
EXTERN TPM_HANDLE       g_exclusiveAuditSession;

//*** g_time
// This is the value in which we keep the current command time. This is initialized
// at the start of each command. The time is the accumulated time since the last
// time that the TPM's timer was last powered up. Clock is the accumulated time
// since the last time that the TPM was cleared. g_time is in mS.
EXTERN  UINT64          g_time;

//*** g_timeEpoch
// This value contains the current clock Epoch. It changes when there is a clock 
// discontinuity. It may be necessary to place this in NV should the timer be able
// to run across a power down of the TPM but not in all cases (e.g. dead battery).
// If the nonce is placed in NV, it should go in gp because it should be changing
// slowly.
#if CLOCK_STOPS
EXTERN CLOCK_NONCE       g_timeEpoch;
#else
#define g_timeEpoch      gp.timeEpoch
#endif
 

//*** g_phEnable
// This is the platform hierarchy control and determines if the platform hierarchy
// is available. This value is SET on each TPM2_Startup(). The default value is
// SET.
EXTERN BOOL             g_phEnable;

//*** g_pcrReConfig
// This value is SET if a TPM2_PCR_Allocate command successfully executed since
// the last TPM2_Startup(). If so, then the next shutdown is required to be
// Shutdown(CLEAR).
EXTERN BOOL             g_pcrReConfig;

//*** g_DRTMHandle
// This location indicates the sequence object handle that holds the DRTM
// sequence data. When not used, it is set to TPM_RH_UNASSIGNED. A sequence
// DRTM sequence is started on either _TPM_Init or _TPM_Hash_Start.
EXTERN TPMI_DH_OBJECT   g_DRTMHandle;

//*** g_DrtmPreStartup
// This value indicates that an H-CRTM occurred after _TPM_Init but before
// TPM2_Startup(). The define for PRE_STARTUP_FLAG is used to add the
// g_DrtmPreStartup value to gp_orderlyState at shutdown. This hack is to avoid
// adding another NV variable.
EXTERN  BOOL            g_DrtmPreStartup;

//*** g_StartupLocality3
// This value indicates that a TPM2_Startup() occurred at locality 3. Otherwise, it
// at locality 0. The define for STARTUP_LOCALITY_3 is to
// indicate that the startup was not at locality 0. This hack is to avoid
// adding another NV variable.
EXTERN  BOOL            g_StartupLocality3;

//***TPM_SU_NONE
// Part 2 defines the two shutdown/startup types that may be used in
// TPM2_Shutdown() and TPM2_Starup(). This additional define is
// used by the TPM to indicate that no shutdown was received.
// NOTE: This is a reserved value.
#define SU_NONE_VALUE           (0xFFFF)
#define TPM_SU_NONE             (TPM_SU)(SU_NONE_VALUE)

//*** TPM_SU_DA_USED
// As with TPM_SU_NONE, this value is added to allow indication that the shutdown
// was not orderly and that a DA=protected object was reference during the previous
// cycle.
#define SU_DA_USED_VALUE    (SU_NONE_VALUE - 1)
#define TPM_SU_DA_USED      (TPM_SU)(SU_DA_USED_VALUE)



//*** Startup Flags
// These flags are included in gp.orderlyState. These are hacks and are being
// used to avoid having to change the layout of gp. The PRE_STARTUP_FLAG indicates
// that a _TPM_Hash_Start/_Data/_End sequence was received after _TPM_Init but
// before TPM2_StartUp(). STARTUP_LOCALITY_3 indicates that the last TPM2_Startup()
// was received at locality 3. These flags are only  relevant if after a 
// TPM2_Shutdown(STATE).
#define PRE_STARTUP_FLAG     0x8000
#define STARTUP_LOCALITY_3   0x4000

#if USE_DA_USED
//*** g_daUsed
// This location indicates if a DA-protected value is accessed during a boot
// cycle. If none has, then there is no need to increment 'failedTries' on the
// next non-orderly startup. This bit is merged with gp.orderlyState when that
// gp.orderly is set to SU_NONE_VALUE
EXTERN  BOOL                 g_daUsed;
#endif

//*** g_updateNV
// This flag indicates if NV should be updated at the end of a command.
// This flag is set to UT_NONE at the beginning of each command in ExecuteCommand().
// This flag is checked in ExecuteCommand() after the detailed actions of a command
// complete. If the command execution was successful and this flag is not UT_NONE,
// any pending NV writes will be committed to NV.
// UT_ORDERLY causes any RAM data to be written to the orderly space for staging
// the write to NV.
typedef BYTE        UPDATE_TYPE; 
#define UT_NONE     (UPDATE_TYPE)0
#define UT_NV       (UPDATE_TYPE)1
#define UT_ORDERLY  (UPDATE_TYPE)(UT_NV + 2)
EXTERN UPDATE_TYPE          g_updateNV;

//*** g_powerWasLost
// This flag is used to indicate if the power was lost. It is SET in _TPM__Init.
// This flag is cleared by TPM2_Startup() after all power-lost activities are
// completed.
// Note: When power is applied, this value can come up as anything. However, 
// _plat__WasPowerLost() will provide the proper indication in that case. So, when
// power is actually lost, we get the correct answer. When power was not lost, but
// the power-lost processing has not been completed before the next _TPM_Init(), 
// then the TPM still does the correct thing.
EXTERN BOOL             g_powerWasLost;

//*** g_clearOrderly
// This flag indicates if the execution of a command should cause the orderly
// state to be cleared.  This flag is set to FALSE at the beginning of each
// command in ExecuteCommand() and is checked in ExecuteCommand() after the
// detailed actions of a command complete but before the check of
// 'g_updateNV'. If this flag is TRUE, and the orderly state is not
// SU_NONE_VALUE, then the orderly state in NV memory will be changed to
// SU_NONE_VALUE or SU_DA_USED_VALUE.
EXTERN BOOL             g_clearOrderly;

//*** g_prevOrderlyState
// This location indicates how the TPM was shut down before the most recent
// TPM2_Startup(). This value, along with the startup type, determines if
// the TPM should do a TPM Reset, TPM Restart, or TPM Resume.
EXTERN TPM_SU           g_prevOrderlyState;

//*** g_nvOk
// This value indicates if the NV integrity check was successful or not. If not and
// the failure was severe, then the TPM would have been put into failure mode after
// it had been re-manufactured. If the NV failure was in the area where the state-save
// data is kept, then this variable will have a value of FALSE indicating that
// a TPM2_Startup(CLEAR) is required.
EXTERN BOOL             g_nvOk;
// NV availability is sampled as the start of each command and stored here
// so that its value remains consistent during the command execution
EXTERN TPM_RC           g_NvStatus;

#ifdef  VENDOR_PERMANENT
//*** g_platformUnique
// This location contains the unique value(s) used to identify the TPM. It is
// loaded on every _TPM2_Startup()
// The first value is used to seed the RNG. The second value is used as a vendor
// authValue. The value used by the RNG would be the value derived from the
// chip unique value (such as fused) with a dependency on the authorities of the
// code in the TPM boot path. The second would be derived from the chip unique value
// with a dependency on the details of the code in the boot path. That is, the
// first value depends on the various signers of the code and the second depends on
// what was signed. The TPM vendor should not be able to know the first value but
// they are expected to know the second.
EXTERN TPM2B_AUTH       g_platformUniqueAuthorities; // Reserved for RNG

EXTERN TPM2B_AUTH       g_platformUniqueDetails;   // referenced by VENDOR_PERMANENT
#endif

//*********************************************************************************
//*********************************************************************************
//** Persistent Global Values
//*********************************************************************************
//*********************************************************************************
//*** Description
// The values in this section are global values that are persistent across power
// events. The lifetime of the values determines the structure in which the value
// is placed.

//*********************************************************************************
//*** PERSISTENT_DATA
//*********************************************************************************
// This structure holds the persistent values that only change as a consequence
// of a specific Protected Capability and are not affected by TPM power events
// (TPM2_Startup() or TPM2_Shutdown().
typedef struct
{
//*********************************************************************************
//          Hierarchy
//*********************************************************************************
// The values in this section are related to the hierarchies.

    BOOL                disableClear;       // TRUE if TPM2_Clear() using
                                            // lockoutAuth is disabled

    // Hierarchy authPolicies
    TPMI_ALG_HASH       ownerAlg;
    TPMI_ALG_HASH       endorsementAlg;
    TPMI_ALG_HASH       lockoutAlg;
    TPM2B_DIGEST        ownerPolicy;
    TPM2B_DIGEST        endorsementPolicy;
    TPM2B_DIGEST        lockoutPolicy;

    // Hierarchy authValues
    TPM2B_AUTH          ownerAuth;
    TPM2B_AUTH          endorsementAuth;
    TPM2B_AUTH          lockoutAuth;

    // Primary Seeds
    TPM2B_SEED          EPSeed;
    TPM2B_SEED          SPSeed;
    TPM2B_SEED          PPSeed;
    // Note there is a nullSeed in the state_reset memory.

    // Hierarchy proofs
    TPM2B_PROOF          phProof;
    TPM2B_PROOF          shProof;
    TPM2B_PROOF          ehProof;
    // Note there is a nullProof in the state_reset memory.

//*********************************************************************************
//          Reset Events
//*********************************************************************************
// A count that increments at each TPM reset and never get reset during the life
// time of TPM.  The value of this counter is initialized to 1 during TPM
// manufacture process. It is used to invalidate all saved contexts after a TPM
// Reset.
    UINT64              totalResetCount;

// This counter increments on each TPM Reset. The counter is reset by
// TPM2_Clear().
    UINT32              resetCount;

//*********************************************************************************
//          PCR
//*********************************************************************************
// This structure hold the policies for those PCR that have an update policy.
// This implementation only supports a single group of PCR controlled by
// policy. If more are required, then this structure would be changed to
// an array.
#if defined NUM_POLICY_PCR_GROUP && NUM_POLICY_PCR_GROUP > 0
    PCR_POLICY          pcrPolicies;
#endif

// This structure indicates the allocation of PCR. The structure contains a
// list of PCR allocations for each implemented algorithm. If no PCR are
// allocated for an algorithm, a list entry still exists but the bit map
// will contain no SET bits.
    TPML_PCR_SELECTION  pcrAllocated;

//*********************************************************************************
//          Physical Presence
//*********************************************************************************
// The PP_LIST type contains a bit map of the commands that require physical
// to be asserted when the authorization is evaluated. Physical presence will be
// checked if the corresponding bit in the array is SET and if the authorization
// handle is TPM_RH_PLATFORM.
//
// These bits may be changed with TPM2_PP_Commands().
    BYTE                ppList[(COMMAND_COUNT + 7) / 8];

//*********************************************************************************
//          Dictionary attack values
//*********************************************************************************
// These values are used for dictionary attack tracking and control.
    UINT32              failedTries;        // the current count of unexpired
                                            // authorization failures

    UINT32              maxTries;           // number of unexpired authorization
                                            // failures before the TPM is in
                                            // lockout

    UINT32              recoveryTime;       // time between authorization failures
                                            // before failedTries is decremented

    UINT32              lockoutRecovery;    // time that must expire between
                                            // authorization failures associated
                                            // with lockoutAuth

    BOOL                lockOutAuthEnabled; // TRUE if use of lockoutAuth is
                                            // allowed

//*****************************************************************************
//            Orderly State
//*****************************************************************************
// The orderly state for current cycle
    TPM_SU              orderlyState;

//*****************************************************************************
//           Command audit values.
//*****************************************************************************
    BYTE                auditCommands[((COMMAND_COUNT + 1) + 7) / 8];
    TPMI_ALG_HASH       auditHashAlg;
    UINT64              auditCounter;

//*****************************************************************************
//           Algorithm selection
//*****************************************************************************
//
// The 'algorithmSet' value indicates the collection of algorithms that are
// currently in used on the TPM.  The interpretation of value is vendor dependent.
    UINT32              algorithmSet;

//*****************************************************************************
//           Firmware version
//*****************************************************************************
// The firmwareV1 and firmwareV2 values are instanced in TimeStamp.c. This is
// a scheme used in development to allow determination of the linker build time
// of the TPM. An actual implementation would implement these values in a way that
// is consistent with vendor needs. The values are maintained in RAM for simplified
// access with a master version in NV.  These values are modified in a
// vendor-specific way.

// g_firmwareV1 contains the more significant 32-bits of the vendor version number.
// In the reference implementation, if this value is printed as a hex
// value, it will have the format of YYYYMMDD
    UINT32              firmwareV1;

// g_firmwareV1 contains the less significant 32-bits of the vendor version number.
// In the reference implementation, if this value is printed as a hex
// value, it will have the format of 00 HH MM SS
    UINT32              firmwareV2;
//*****************************************************************************
//           Timer Epoch
//*****************************************************************************
// timeEpoch contains a nonce that has a vendor=specific size (should not be
// less than 8 bytes. This nonce changes when the clock epoch changes. The clock
// epoch changes when there is a discontinuity in the timing of the TPM.
#if !CLOCK_STOPS
    CLOCK_NONCE         timeEpoch;
#endif

} PERSISTENT_DATA;

EXTERN PERSISTENT_DATA  gp;

//*********************************************************************************
//*********************************************************************************
//*** ORDERLY_DATA
//*********************************************************************************
//*********************************************************************************
// The data in this structure is saved to NV on each TPM2_Shutdown().
typedef struct orderly_data
{
//*****************************************************************************
//           TIME
//*****************************************************************************

// Clock has two parts. One is the state save part and one is the NV part. The
// state save version is updated on each command. When the clock rolls over, the
// NV version is updated. When the TPM starts up, if the TPM was shutdown in and
// orderly way, then the sClock value is used to initialize the clock. If the
// TPM shutdown was not orderly, then the persistent value is used and the safe
// attribute is clear.

    UINT64              clock;              // The orderly version of clock
    TPMI_YES_NO         clockSafe;          // Indicates if the clock value is
                                            // safe.

    // In many implementations, the quality of the entropy available is not that
    // high. To compensate, the current value of the drbgState can be saved and
    // restored on each power cycle. This prevents the internal state from reverting
    // to the initial state on each power cycle and starting with a limited amount
    // of entropy. By keeping the old state and adding entropy, the entropy will
    // accumulate.
    DRBG_STATE          drbgState;

// These values allow the accumulation of self-healing time across orderly shutdown
// of the TPM.
#if ACCUMULATE_SELF_HEAL_TIMER 
    UINT64              selfHealTimer;  // current value of s_selfHealTimer
    UINT64              lockoutTimer;   // current value of s_lockoutTimer
    UINT64              time;           // current value of g_time at shutdown
#endif // ACCUMULATE_SELF_HEAL_TIMER

} ORDERLY_DATA;

#if ACCUMULATE_SELF_HEAL_TIMER
#define     s_selfHealTimer     go.selfHealTimer
#define     s_lockoutTimer      go.lockoutTimer
#endif  // ACCUMULATE_SELF_HEAL_TIMER

#  define drbgDefault go.drbgState

EXTERN ORDERLY_DATA     go;

//*********************************************************************************
//*********************************************************************************
//*** STATE_CLEAR_DATA
//*********************************************************************************
//*********************************************************************************
// This structure contains the data that is saved on Shutdown(STATE)
// and restored on Startup(STATE).  The values are set to their default
// settings on any Startup(Clear). In other words, the data is only persistent
// across TPM Resume.
//
// If the comments associated with a parameter indicate a default reset value, the
// value is applied on each Startup(CLEAR).

typedef struct state_clear_data
{
//*****************************************************************************
//           Hierarchy Control
//*****************************************************************************
    BOOL                shEnable;           // default reset is SET
    BOOL                ehEnable;           // default reset is SET
    BOOL                phEnableNV;         // default reset is SET
    TPMI_ALG_HASH       platformAlg;        // default reset is TPM_ALG_NULL
    TPM2B_DIGEST        platformPolicy;     // default reset is an Empty Buffer
    TPM2B_AUTH          platformAuth;       // default reset is an Empty Buffer

//*****************************************************************************
//           PCR
//*****************************************************************************
// The set of PCR to be saved on Shutdown(STATE)
    PCR_SAVE            pcrSave;            // default reset is 0...0

// This structure hold the authorization values for those PCR that have an
// update authorization.
// This implementation only supports a single group of PCR controlled by
// authorization. If more are required, then this structure would be changed to
// an array.
    PCR_AUTHVALUE       pcrAuthValues;
} STATE_CLEAR_DATA;

EXTERN STATE_CLEAR_DATA gc;

//*********************************************************************************
//*********************************************************************************
//***  State Reset Data
//*********************************************************************************
//*********************************************************************************
// This structure contains data is that is saved on Shutdown(STATE) and restored on
// the subsequent Startup(ANY). That is, the data is preserved across TPM Resume
// and TPM Restart.
//
// If a default value is specified in the comments this value is applied on
// TPM Reset.

typedef struct state_reset_data
{
//*****************************************************************************
//          Hierarchy Control
//*****************************************************************************
    TPM2B_PROOF         nullProof;          // The proof value associated with
                                            // the TPM_RH_NULL hierarchy. The
                                            // default reset value is from the RNG.

    TPM2B_SEED          nullSeed;           // The seed value for the TPM_RN_NULL
                                            // hierarchy. The default reset value
                                            // is from the RNG.

//*****************************************************************************
//           Context
//*****************************************************************************
// The 'clearCount' counter is incremented each time the TPM successfully executes
// a TPM Resume. The counter is included in each saved context that has 'stClear'
// SET (including descendants of keys that have 'stClear' SET). This prevents these
// objects from being loaded after a TPM Resume.
// If 'clearCount' is at its maximum value when the TPM receives a Shutdown(STATE),
// the TPM will return TPM_RC_RANGE and the TPM will only accept Shutdown(CLEAR).
    UINT32              clearCount;         // The default reset value is 0.

    UINT64              objectContextID;    // This is the context ID for a saved
                                            //  object context. The default reset
                                            //  value is 0.
#ifndef NDEBUG
#undef  CONTEXT_SLOT
#define CONTEXT_SLOT     BYTE
#endif

    CONTEXT_SLOT        contextArray[MAX_ACTIVE_SESSIONS];    // This array contains
                                            // contains the values used to track
                                            // the version numbers of saved
                                            // contexts (see
                                            // Session.c in for details). The
                                            // default reset value is {0}.

    CONTEXT_COUNTER     contextCounter;     // This is the value from which the
                                            // 'contextID' is derived. The
                                            // default reset value is {0}.

//*****************************************************************************
//           Command Audit
//*****************************************************************************
// When an audited command completes, ExecuteCommand() checks the return
// value.  If it is TPM_RC_SUCCESS, and the command is an audited command, the
// TPM will extend the cpHash and rpHash for the command to this value. If this
// digest was the Zero Digest before the cpHash was extended, the audit counter
// is incremented.

    TPM2B_DIGEST        commandAuditDigest; // This value is set to an Empty Digest
                                            // by TPM2_GetCommandAuditDigest() or a
                                            // TPM Reset.

//*****************************************************************************
//           Boot counter
//*****************************************************************************

    UINT32              restartCount;       // This counter counts TPM Restarts.
                                            // The default reset value is 0.

//*********************************************************************************
//            PCR
//*********************************************************************************
// This counter increments whenever the PCR are updated. This counter is preserved
// across TPM Resume even though the PCR are not preserved. This is because
// sessions remain active across TPM Restart and the count value in the session
// is compared to this counter so this counter must have values that are unique
// as long as the sessions are active.
// NOTE: A platform-specific specification may designate that certain PCR changes
//       do not increment this counter to increment.
    UINT32              pcrCounter;         // The default reset value is 0.

#if     ALG_ECC

//*****************************************************************************
//         ECDAA
//*****************************************************************************
    UINT64              commitCounter;      // This counter increments each time
                                            // TPM2_Commit() returns
                                            // TPM_RC_SUCCESS. The default reset
                                            // value is 0.

    TPM2B_NONCE         commitNonce;        // This random value is used to compute
                                            // the commit values. The default reset
                                            // value is from the RNG.

// This implementation relies on the number of bits in g_commitArray being a
// power of 2 (8, 16, 32, 64, etc.) and no greater than 64K.
    BYTE                 commitArray[16];   // The default reset value is {0}.

#endif // ALG_ECC
} STATE_RESET_DATA;

EXTERN STATE_RESET_DATA gr;

//** NV Layout
// The NV data organization is
// 1) a PERSISTENT_DATA structure
// 2) a STATE_RESET_DATA structure
// 3) a STATE_CLEAR_DATA structure
// 4) an ORDERLY_DATA structure
// 5) the user defined NV index space
#define NV_PERSISTENT_DATA  (0)
#define NV_STATE_RESET_DATA (NV_PERSISTENT_DATA + sizeof(PERSISTENT_DATA))
#define NV_STATE_CLEAR_DATA (NV_STATE_RESET_DATA + sizeof(STATE_RESET_DATA))
#define NV_ORDERLY_DATA     (NV_STATE_CLEAR_DATA + sizeof(STATE_CLEAR_DATA))
#define NV_INDEX_RAM_DATA   (NV_ORDERLY_DATA + sizeof(ORDERLY_DATA))
#define NV_USER_DYNAMIC     (NV_INDEX_RAM_DATA + sizeof(s_indexOrderlyRam))
#define NV_USER_DYNAMIC_END     NV_MEMORY_SIZE

//** Global Macro Definitions
// The NV_READ_PERSISTENT and NV_WRITE_PERSISTENT macros are used to access members
// of the PERSISTENT_DATA structure in NV.
#define NV_READ_PERSISTENT(to, from)                \
            NvRead(&to, offsetof(PERSISTENT_DATA, from), sizeof(to))

#define NV_WRITE_PERSISTENT(to, from)               \
            NvWrite(offsetof(PERSISTENT_DATA, to), sizeof(gp.to), &from)

#define CLEAR_PERSISTENT(item)                      \
            NvClearPersistent(offsetof(PERSISTENT_DATA, item), sizeof(gp.item))

#define NV_SYNC_PERSISTENT(item) NV_WRITE_PERSISTENT(item, gp.item)

// At the start of command processing, the index of the command is determined. This
// index value is used to access the various data tables that contain per-command
// information. There are multiple options for how the per-command tables can be
// implemented. This is resolved in GetClosestCommandIndex().
typedef UINT16      COMMAND_INDEX;
#define UNIMPLEMENTED_COMMAND_INDEX     ((COMMAND_INDEX)(~0))

typedef struct _COMMAND_FLAGS_ 
{
    unsigned    trialPolicy : 1;    //1) If SET, one of the handles references a
                                    //   trial policy and authorization may be
                                    //   skipped. This is only allowed for a policy
                                    //   command.
} COMMAND_FLAGS;
        
// This structure is used to avoid having to manage a large number of
// parameters being passed through various levels of the command input processing.
//
typedef struct _COMMAND_
{
    TPM_ST           tag;               // the parsed command tag
    TPM_CC           code;              // the parsed command code
    COMMAND_INDEX    index;             // the computed command index
    UINT32           handleNum;         // the number of entity handles in the 
                                        //   handle area of the command
    TPM_HANDLE       handles[MAX_HANDLE_NUM]; // the parsed handle values
    UINT32           sessionNum;        // the number of sessions found
    INT32            parameterSize;     // starts out with the parsed command size
                                        // and is reduced and values are
                                        // unmarshaled. Just before calling the 
                                        // command actions, this should be zero. 
                                        // After the command actions, this number
                                        // should grow as values are marshaled
                                        // in to the response buffer.
    INT32            authSize;          // this is initialized with the parsed size
                                        // of authorizationSize field and should
                                        // be zero when the authorizations are
                                        // parsed.
    BYTE            *parameterBuffer;   // input to ExecuteCommand
    BYTE            *responseBuffer;    // input to ExecuteCommand
#if ALG_SHA1
    TPM2B_SHA1_DIGEST   sha1CpHash;
    TPM2B_SHA1_DIGEST   sha1RpHash;
#endif
#if ALG_SHA256
    TPM2B_SHA256_DIGEST sha256CpHash;
    TPM2B_SHA256_DIGEST sha256RpHash;
#endif
#if ALG_SHA384
    TPM2B_SHA384_DIGEST sha384CpHash;
    TPM2B_SHA384_DIGEST sha384RpHash;
#endif
#if ALG_SHA512
    TPM2B_SHA512_DIGEST sha512CpHash;
    TPM2B_SHA512_DIGEST sha512RpHash;
#endif
#if ALG_SM3_256
    TPM2B_SM3_256_DIGEST sm3_256CpHash;
    TPM2B_SM3_256_DIGEST sm3_256RpHash;
#endif
} COMMAND;

// Global sting constants for consistency in KDF function calls.
// These string constants are shared across functions to make sure that they 
// are all using consistent sting values.

#define STRING_INITIALIZER(value)   {{sizeof(value), {value}}}
#define TPM2B_STRING(name, value)                                                   \
typedef union name##_ {                                                             \
        struct  {                                                                   \
            UINT16  size;                                                           \
            BYTE    buffer[sizeof(value)];                                          \
        } t;                                                                        \
        TPM2B   b;                                                                  \
    } TPM2B_##name##_;                                                              \
EXTERN  const TPM2B_##name##_      name##_ INITIALIZER(STRING_INITIALIZER(value));  \
EXTERN  const TPM2B               *name INITIALIZER(&name##_.b)

TPM2B_STRING(PRIMARY_OBJECT_CREATION, "Primary Object Creation");
TPM2B_STRING(CFB_KEY, "CFB");
TPM2B_STRING(CONTEXT_KEY, "CONTEXT");
TPM2B_STRING(INTEGRITY_KEY, "INTEGRITY");
TPM2B_STRING(SECRET_KEY, "SECRET");
TPM2B_STRING(SESSION_KEY, "ATH");
TPM2B_STRING(STORAGE_KEY, "STORAGE");
TPM2B_STRING(XOR_KEY, "XOR");
TPM2B_STRING(COMMIT_STRING, "ECDAA Commit");
TPM2B_STRING(DUPLICATE_STRING, "DUPLICATE");
TPM2B_STRING(IDENTITY_STRING, "IDENTITY");
TPM2B_STRING(OBFUSCATE_STRING, "OBFUSCATE");
#if SELF_TEST
TPM2B_STRING(OAEP_TEST_STRING, "OAEP Test Value");
#endif // SELF_TEST

//*****************************************************************************
//** From CryptTest.c
//*****************************************************************************
// This structure contains the self-test state values for the cryptographic modules.
EXTERN CRYPTO_SELF_TEST_STATE   g_cryptoSelfTestState;

//*****************************************************************************
//** From Manufacture.c
//*****************************************************************************
EXTERN BOOL              g_manufactured INITIALIZER(FALSE);

// This value indicates if a TPM2_Startup commands has been
// receive since the power on event.  This flag is maintained in power
// simulation module because this is the only place that may reliably set this
// flag to FALSE.
EXTERN BOOL              g_initialized;

//** Private data

//*****************************************************************************
//*** From SessionProcess.c
//*****************************************************************************
#if defined SESSION_PROCESS_C || defined GLOBAL_C || defined MANUFACTURE_C
// The following arrays are used to save command sessions information so that the
// command handle/session buffer does not have to be preserved for the duration of
// the command. These arrays are indexed by the session index in accordance with
// the order of sessions in the session area of the command.
//
// Array of the authorization session handles
EXTERN TPM_HANDLE       s_sessionHandles[MAX_SESSION_NUM];

// Array of authorization session attributes
EXTERN TPMA_SESSION     s_attributes[MAX_SESSION_NUM];

// Array of handles authorized by the corresponding authorization sessions;
// and if none, then TPM_RH_UNASSIGNED value is used
EXTERN TPM_HANDLE       s_associatedHandles[MAX_SESSION_NUM];

// Array of nonces provided by the caller for the corresponding sessions
EXTERN TPM2B_NONCE      s_nonceCaller[MAX_SESSION_NUM];

// Array of authorization values (HMAC's or passwords) for the corresponding
// sessions
EXTERN TPM2B_AUTH       s_inputAuthValues[MAX_SESSION_NUM];

// Array of pointers to the SESSION structures for the sessions in a command
EXTERN SESSION          *s_usedSessions[MAX_SESSION_NUM];

// Special value to indicate an undefined session index
#define             UNDEFINED_INDEX     (0xFFFF)

// Index of the session used for encryption of a response parameter
EXTERN UINT32           s_encryptSessionIndex;

// Index of the session used for decryption of a command parameter
EXTERN UINT32           s_decryptSessionIndex;

// Index of a session used for audit
EXTERN UINT32           s_auditSessionIndex;

// The cpHash for command audit
#ifdef  TPM_CC_GetCommandAuditDigest
EXTERN TPM2B_DIGEST    s_cpHashForCommandAudit;
#endif

// Flag indicating if NV update is pending for the lockOutAuthEnabled or
// failedTries DA parameter
EXTERN BOOL             s_DAPendingOnNV;

#endif // SESSION_PROCESS_C

//*****************************************************************************
//*** From DA.c
//*****************************************************************************
#if defined DA_C || defined GLOBAL_C || defined MANUFACTURE_C
// This variable holds the accumulated time since the last time
// that 'failedTries' was decremented. This value is in millisecond.
#if !ACCUMULATE_SELF_HEAL_TIMER
EXTERN UINT64       s_selfHealTimer;

// This variable holds the accumulated time that the lockoutAuth has been
// blocked.
EXTERN UINT64       s_lockoutTimer;
#endif // ACCUMULATE_SELF_HEAL_TIMER

#endif // DA_C

//*****************************************************************************
//*** From NV.c
//*****************************************************************************
#if defined NV_C || defined GLOBAL_C
// This marks the end of the NV area. This is a run-time variable as it might
// not be compile-time constant.
EXTERN NV_REF   s_evictNvEnd;

// This space is used to hold the index data for an orderly Index. It also contains
// the attributes for the index.
EXTERN BYTE      s_indexOrderlyRam[RAM_INDEX_SPACE];   // The orderly NV Index data

// This value contains the current max counter value. It is written to the end of
// allocatable NV space each time an index is deleted or added. This value is
// initialized on Startup. The indices are searched and the maximum of all the
// current counter indices and this value is the initial value for this.
EXTERN UINT64    s_maxCounter;

// This is space used for the NV Index cache. As with a persistent object, the
// contents of a referenced index are copied into the cache so that the
// NV Index memory scanning and data copying can be reduced.
// Only code that operates on NV Index data should use this cache directly. When
// that action code runs, s_lastNvIndex will contain the index header information.
// It will have been loaded when the handles were verified.
// NOTE: An NV index handle can appear in many commands that do not operate on the
// NV data (e.g. TPM2_StartAuthSession). However, only one NV Index at a time is
// ever directly referenced by any command. If that changes, then the NV Index
// caching needs to be changed to accommodate that. Currently, the code will verify
// that only one NV Index is referenced by the handles of the command.
EXTERN      NV_INDEX         s_cachedNvIndex;
EXTERN      NV_REF           s_cachedNvRef;
EXTERN      BYTE            *s_cachedNvRamRef;

// Initial NV Index/evict object iterator value
#define     NV_REF_INIT     (NV_REF)0xFFFFFFFF

#endif

//*****************************************************************************
//*** From Object.c
//*****************************************************************************
#if defined OBJECT_C || defined GLOBAL_C
// This type is the container for an object.

EXTERN OBJECT           s_objects[MAX_LOADED_OBJECTS];

#endif // OBJECT_C

//*****************************************************************************
//*** From PCR.c
//*****************************************************************************
#if defined PCR_C || defined GLOBAL_C
typedef struct
{
#if     ALG_SHA1
    // SHA1 PCR
    BYTE    sha1Pcr[SHA1_DIGEST_SIZE];
#endif
#if     ALG_SHA256
    // SHA256 PCR
    BYTE    sha256Pcr[SHA256_DIGEST_SIZE];
#endif
#if     ALG_SHA384
    // SHA384 PCR
    BYTE    sha384Pcr[SHA384_DIGEST_SIZE];
#endif
#if     ALG_SHA512
    // SHA512 PCR
    BYTE    sha512Pcr[SHA512_DIGEST_SIZE];
#endif
#if     ALG_SM3_256
    // SHA256 PCR
    BYTE    sm3_256Pcr[SM3_256_DIGEST_SIZE];
#endif
} PCR;

typedef struct
{
    unsigned int    stateSave : 1;              // if the PCR value should be
                                                // saved in state save
    unsigned int    resetLocality : 5;          // The locality that the PCR
                                                // can be reset
    unsigned int    extendLocality : 5;         // The locality that the PCR
                                                // can be extend
} PCR_Attributes;

EXTERN PCR          s_pcrs[IMPLEMENTATION_PCR];

#endif // PCR_C

//*****************************************************************************
//*** From Session.c
//*****************************************************************************
#if defined SESSION_C || defined GLOBAL_C
// Container for HMAC or policy session tracking information
typedef struct
{
    BOOL                occupied;
    SESSION             session;        // session structure
} SESSION_SLOT;

EXTERN SESSION_SLOT     s_sessions[MAX_LOADED_SESSIONS];

//  The index in contextArray that has the value of the oldest saved session
//  context. When no context is saved, this will have a value that is greater
//  than or equal to MAX_ACTIVE_SESSIONS.
EXTERN UINT32            s_oldestSavedSession;

// The number of available session slot openings.  When this is 1,
// a session can't be created or loaded if the GAP is maxed out.
// The exception is that the oldest saved session context can always
// be loaded (assuming that there is a space in memory to put it)
EXTERN int               s_freeSessionSlots;

#endif // SESSION_C

//*****************************************************************************
//*** From IoBuffers.c
//*****************************************************************************
#if defined IO_BUFFER_C || defined GLOBAL_C
// Each command function is allowed a structure for the inputs to the function and
// a structure for the outputs. The command dispatch code unmarshals the input butter 
// to the command action input structure starting at the first byte of
// s_actionIoBuffer. The value of s_actionIoAllocation is the number of UINT64 values 
// allocated. It is used to set the pointer for the response structure. The command 
// dispatch code will marshal the response values into the final output buffer.
EXTERN UINT64   s_actionIoBuffer[768];      // action I/O buffer
EXTERN UINT32   s_actionIoAllocation;       // number of UIN64 allocated for the
                                            // action input structure
#endif // IO_BUFFER_C

//*****************************************************************************
//*** From TPMFail.c
//*****************************************************************************
// This value holds the address of the string containing the name of the function
// in which the failure occurred. This address value isn't useful for anything
// other than helping the vendor to know in which file the failure  occurred.
EXTERN BOOL      g_inFailureMode;       // Indicates that the TPM is in failure mode
#if SIMULATION
EXTERN BOOL      g_forceFailureMode;    // flag to force failure mode during test
#endif

typedef void(FailFunction)(const char *function, int line, int code);

#if defined TPM_FAIL_C || defined GLOBAL_C
EXTERN UINT32    s_failFunction;
EXTERN UINT32    s_failLine;            // the line in the file at which
                                        // the error was signaled
EXTERN UINT32    s_failCode;            // the error code used

EXTERN FailFunction    *LibFailCallback;

#endif // TPM_FAIL_C

//*****************************************************************************
//*** From CommandCodeAttributes.c
//*****************************************************************************
// This array is instanced in CommandCodeAttributes.c when it includes
// CommandCodeAttributes.h. Don't change the extern to EXTERN.
extern  const  TPMA_CC               s_ccAttr[];
extern  const  COMMAND_ATTRIBUTES    s_commandAttributes[];

#endif // GLOBAL_H
