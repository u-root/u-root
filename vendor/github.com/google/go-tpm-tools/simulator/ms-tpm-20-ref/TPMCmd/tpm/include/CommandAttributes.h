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
 *  Created by TpmStructures; Version 3.0 June 16, 2017
 *  Date: Aug 14, 2017  Time: 02:53:08PM
 */
// The attributes defined in this file are produced by the parser that
// creates the structure definitions from Part 3. The attributes are defined
// in that parser and should track the attributes being tested in
// CommandCodeAttributes.c. Generally, when an attribute is added to this list, 
// new code will be needed in CommandCodeAttributes.c to test it. 

#ifndef COMMAND_ATTRIBUTES_H
#define COMMAND_ATTRIBUTES_H

typedef UINT16              COMMAND_ATTRIBUTES;
#define NOT_IMPLEMENTED     (COMMAND_ATTRIBUTES)(0)
#define ENCRYPT_2           ((COMMAND_ATTRIBUTES)1 << 0)
#define ENCRYPT_4           ((COMMAND_ATTRIBUTES)1 << 1)
#define DECRYPT_2           ((COMMAND_ATTRIBUTES)1 << 2)
#define DECRYPT_4           ((COMMAND_ATTRIBUTES)1 << 3)
#define HANDLE_1_USER       ((COMMAND_ATTRIBUTES)1 << 4)
#define HANDLE_1_ADMIN      ((COMMAND_ATTRIBUTES)1 << 5)
#define HANDLE_1_DUP        ((COMMAND_ATTRIBUTES)1 << 6)
#define HANDLE_2_USER       ((COMMAND_ATTRIBUTES)1 << 7)
#define PP_COMMAND          ((COMMAND_ATTRIBUTES)1 << 8)
#define IS_IMPLEMENTED      ((COMMAND_ATTRIBUTES)1 << 9)
#define NO_SESSIONS         ((COMMAND_ATTRIBUTES)1 << 10)
#define NV_COMMAND          ((COMMAND_ATTRIBUTES)1 << 11)
#define PP_REQUIRED         ((COMMAND_ATTRIBUTES)1 << 12)
#define R_HANDLE            ((COMMAND_ATTRIBUTES)1 << 13)
#define ALLOW_TRIAL         ((COMMAND_ATTRIBUTES)1 << 14)

#endif // COMMAND_ATTRIBUTES_H
