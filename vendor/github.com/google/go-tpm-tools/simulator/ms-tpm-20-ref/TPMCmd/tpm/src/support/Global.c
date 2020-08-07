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
// This file will instance the TPM variables that are not stack allocated. 

// Descriptions of global variables are in Global.h. There macro macro definitions
// that allows a variable to be instanced or simply defined as an external variable.
// When global.h is included from this .c file, GLOBAL_C is defined and values are 
// instanced (and possibly initialized), but when global.h is included by any other 
// file, they are simply defined as external values. DO NOT DEFINE GLOBAL_C IN ANY 
// OTHER FILE.
//
// NOTE: This is a change from previous implementations where Global.h just contained
// the extern declaration and values were instanced in this file. This change keeps 
// the definition and instance in one file making maintenance easier. The instanced
// data will still be in the global.obj file.
//
// The OIDs.h file works in a way that is similar to the Global.h with the definition
// of the values in OIDs.h such that they are instanced in global.obj. The macros 
// that are defined in Global.h are used in OIDs.h in the same way as they are in 
// Global.h.

//** Defines and Includes
#define GLOBAL_C
#include "Tpm.h"
#include "OIDs.h"

