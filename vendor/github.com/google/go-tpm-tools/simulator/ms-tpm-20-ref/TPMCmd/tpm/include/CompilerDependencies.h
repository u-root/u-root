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
// This file contains the build switches. This contains switches for multiple
// versions of the crypto-library so some may not apply to your environment.
//

#ifndef _COMPILER_DEPENDENCIES_H_
#define _COMPILER_DEPENDENCIES_H_

#ifdef GCC
#   undef _MSC_VER
#   undef WIN32
#endif

#ifdef _MSC_VER
// These definitions are for the Microsoft compiler

// Endian conversion for aligned structures
#   define REVERSE_ENDIAN_16(_Number) _byteswap_ushort(_Number)
#   define REVERSE_ENDIAN_32(_Number) _byteswap_ulong(_Number)
#   define REVERSE_ENDIAN_64(_Number) _byteswap_uint64(_Number)

// Avoid compiler warning for in line of stdio (or not)
//#define _NO_CRT_STDIO_INLINE

// This macro is used to handle LIB_EXPORT of function and variable names in lieu
// of a .def file. Visual Studio requires that functions be explicitly exported and
// imported.
#   define LIB_EXPORT __declspec(dllexport) // VS compatible version
#   define LIB_IMPORT __declspec(dllimport)

// This is defined to indicate a function that does not return. Microsoft compilers
// do not support the _Noretrun function parameter.
#   define NORETURN  __declspec(noreturn)
#   if _MSC_VER >= 1400     // SAL processing when needed
#       include <sal.h>
#   endif

#   ifdef _WIN64
#       define _INTPTR 2
#    else
#       define _INTPTR 1
#    endif


#define NOT_REFERENCED(x)   (x)

// Lower the compiler error warning for system include
// files. They tend not to be that clean and there is no
// reason to sort through all the spurious errors that they
// generate when the normal error level is set to /Wall
#   define _REDUCE_WARNING_LEVEL_(n)                    \
__pragma(warning(push, n))
// Restore the compiler warning level
#   define _NORMAL_WARNING_LEVEL_                       \
__pragma(warning(pop))
#   include <stdint.h>
#endif

#ifndef _MSC_VER
#ifndef WINAPI
#   define WINAPI
#endif
#   define __pragma(x)
#   define REVERSE_ENDIAN_16(_Number) __builtin_bswap16(_Number)
#   define REVERSE_ENDIAN_32(_Number) __builtin_bswap32(_Number)
#   define REVERSE_ENDIAN_64(_Number) __builtin_bswap64(_Number)
#endif

#if defined(__GNUC__)
#   define NORETURN                     __attribute__((noreturn))
#   include <stdint.h>
#endif

// Things that are not defined should be defined as NULL
#ifndef NORETURN
#   define NORETURN
#endif
#ifndef LIB_EXPORT
#   define LIB_EXPORT
#endif
#ifndef LIB_IMPORT
#   define LIB_IMPORT
#endif
#ifndef _REDUCE_WARNING_LEVEL_
#   define _REDUCE_WARNING_LEVEL_(n)
#endif
#ifndef _NORMAL_WARNING_LEVEL_
#   define _NORMAL_WARNING_LEVEL_
#endif
#ifndef NOT_REFERENCED
#   define  NOT_REFERENCED(x) (x = x)
#endif

#ifdef _POSIX_
typedef int SOCKET;
#endif


#endif // _COMPILER_DEPENDENCIES_H_