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
 *  Date: Mar 28, 2019  Time: 08:25:18PM
 */

#ifndef    _CONTEXT_SPT_FP_H_
#define    _CONTEXT_SPT_FP_H_

//*** ComputeContextProtectionKey()
// This function retrieves the symmetric protection key for context encryption
// It is used by TPM2_ConextSave and TPM2_ContextLoad to create the symmetric
// encryption key and iv
//  Return Type: void
void
ComputeContextProtectionKey(
    TPMS_CONTEXT    *contextBlob,   // IN: context blob
    TPM2B_SYM_KEY   *symKey,        // OUT: the symmetric key
    TPM2B_IV        *iv             // OUT: the IV.
);

//*** ComputeContextIntegrity()
// Generate the integrity hash for a context
//       It is used by TPM2_ContextSave to create an integrity hash
//       and by TPM2_ContextLoad to compare an integrity hash
//  Return Type: void
void
ComputeContextIntegrity(
    TPMS_CONTEXT    *contextBlob,   // IN: context blob
    TPM2B_DIGEST    *integrity      // OUT: integrity
);

//*** SequenceDataExport();
// This function is used scan through the sequence object and
// either modify the hash state data for export (contextSave) or to
// import it into the internal format (contextLoad).
// This function should only be called after the sequence object has been copied
// to the context buffer (contextSave) or from the context buffer into the sequence
// object. The presumption is that the context buffer version of the data is the
// same size as the internal representation so nothing outsize of the hash context
// area gets modified.
void
SequenceDataExport(
    HASH_OBJECT         *object,        // IN: an internal hash object
    HASH_OBJECT_BUFFER  *exportObject   // OUT: a sequence context in a buffer
);

//*** SequenceDataImport();
// This function is used scan through the sequence object and
// either modify the hash state data for export (contextSave) or to
// import it into the internal format (contextLoad).
// This function should only be called after the sequence object has been copied
// to the context buffer (contextSave) or from the context buffer into the sequence
// object. The presumption is that the context buffer version of the data is the
// same size as the internal representation so nothing outsize of the hash context
// area gets modified.
void
SequenceDataImport(
    HASH_OBJECT         *object,        // IN/OUT: an internal hash object
    HASH_OBJECT_BUFFER  *exportObject   // IN/OUT: a sequence context in a buffer
);

#endif  // _CONTEXT_SPT_FP_H_
