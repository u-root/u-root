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
// This header file is used to select the library code that gets included in the
// TPM build.

#ifndef _LIB_SUPPORT_H_
#define _LIB_SUPPORT_H_

//*********************
#ifndef RADIX_BITS
#   if defined(__x86_64__) || defined(__x86_64)                                         \
        || defined(__amd64__) || defined(__amd64) || defined(_WIN64) || defined(_M_X64) \
        || defined(_M_ARM64) || defined(__aarch64__)
#       define RADIX_BITS                      64
#   elif defined(__i386__) || defined(__i386) || defined(i386)                          \
        || defined(_WIN32) || defined(_M_IX86)                                          \
        || defined(_M_ARM) || defined(__arm__) || defined(__thumb__)
#       define RADIX_BITS                      32
#   else
#       error Unable to determine RADIX_BITS from compiler environment
#   endif
#endif // RADIX_BITS

// These macros use the selected libraries to the proper include files. 
#define LIB_QUOTE(_STRING_) #_STRING_
#define LIB_INCLUDE2(_LIB_, _TYPE_) LIB_QUOTE(_LIB_/TpmTo##_LIB_##_TYPE_.h)
#define LIB_INCLUDE(_LIB_, _TYPE_) LIB_INCLUDE2(_LIB_, _TYPE_)

// Include the options for hashing and symmetric. Defer the load of the math package
// Until the bignum parameters are defined.
#include LIB_INCLUDE(SYM_LIB, Sym)
#include LIB_INCLUDE(HASH_LIB, Hash)

#undef MIN
#undef MAX

#endif // _LIB_SUPPORT_H_
