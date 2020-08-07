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
// This file contains the macro and structure definitions for the X509 commands and
// functions.

#ifndef _X509_H_
#define _X509_H_

//** Includes

#include "Tpm.h"
#include "TpmASN1.h"

//** Defined Constants

//*** X509 Application-specific types 
#define X509_SELECTION          0xA0
#define X509_ISSUER_UNIQUE_ID   0xA1
#define X509_SUBJECT_UNIQUE_ID  0xA2
#define X509_EXTENSIONS         0xA3

// These defines give the order in which values appear in the TBScertificate
// of an x.509 certificate. These values are used to index into an array of
//
#define ENCODED_SIZE_REF        0
#define VERSION_REF             (ENCODED_SIZE_REF + 1)
#define SERIAL_NUMBER_REF       (VERSION_REF + 1)
#define SIGNATURE_REF           (SERIAL_NUMBER_REF + 1)
#define ISSUER_REF              (SIGNATURE_REF + 1)
#define VALIDITY_REF            (ISSUER_REF + 1)
#define SUBJECT_KEY_REF         (VALIDITY_REF + 1)
#define SUBJECT_PUBLIC_KEY_REF  (SUBJECT_KEY_REF + 1)
#define EXTENSIONS_REF          (SUBJECT_PUBLIC_KEY_REF + 1)
#define REF_COUNT               (EXTENSIONS_REF + 1)

#undef MAKE_OID
#ifdef _X509_SPT_
#   define MAKE_OID(NAME)                  \
        const BYTE      OID##NAME[] = {OID##NAME##_VALUE}
#else
#   define MAKE_OID(NAME)                   \
        extern const BYTE    OID##NAME[]
#endif


//** Structures


// Used to access the fields of a TBSsignature some of which are in the in_CertifyX509
// structure and some of which are in the out_CertifyX509 structure.
typedef struct stringRef
{
    BYTE        *buf;
    INT16        len;
} stringRef;


typedef union x509KeyUsageUnion {
    TPMA_X509_KEY_USAGE     x509;
    UINT32                  integer;
} x509KeyUsageUnion;


//** Global X509 Constants
// These values are instanced by X509_spt.c and referenced by other X509-related
// files.


// This is the DER-encoded value for the Key Usage OID  (2.5.29.15). This is the
// full OID, not just the numeric value
#define OID_KEY_USAGE_EXTENSTION_VALUE  0x06, 0x03, 0x55, 0x1D, 0x0F
MAKE_OID(_KEY_USAGE_EXTENSTION);
  
// This is the DER-encoded value for the TCG-defined TPMA_OBJECT OID
// (2.23.133.10.1.1.1)
#define OID_TCG_TPMA_OBJECT_VALUE       0x06, 0x07, 0x67, 0x81, 0x05, 0x0a, 0x01,   \
                                        0x01, 0x01
MAKE_OID(_TCG_TPMA_OBJECT);

#ifdef _X509_SPT_
const x509KeyUsageUnion keyUsageSign = { TPMA_X509_KEY_USAGE_INITIALIZER(
    /* digitalsignature */ 1, /* nonrepudiation   */ 0,
    /* keyencipherment  */ 0, /* dataencipherment */ 0,
    /* keyagreement     */ 0, /* keycertsign      */ 1,
    /* crlsign          */ 1, /* encipheronly     */ 0,
    /* decipheronly     */ 0, /* bits_at_9        */ 0) };

const x509KeyUsageUnion keyUsageDecrypt = { TPMA_X509_KEY_USAGE_INITIALIZER(
    /* digitalsignature */ 0, /* nonrepudiation   */ 0,
    /* keyencipherment  */ 1, /* dataencipherment */ 1,
    /* keyagreement     */ 1, /* keycertsign      */ 0,
    /* crlsign          */ 0, /* encipheronly     */ 1,
    /* decipheronly     */ 1, /* bits_at_9        */ 0) };
#else
extern x509KeyUsageUnion keyUsageSign;
extern x509KeyUsageUnion keyUsageDecrypt;
#endif

#undef MAKE_OID

#endif // _X509_H_
