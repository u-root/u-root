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

// This file contains the definitions needed for defining the internal BIGNUM
// structure.

// A BIGNUM is a pointer to a structure. The structure has three fields. The
// last field is and array (d) of crypt_uword_t. Each word is in machine format
// (big- or little-endian) with the words in ascending significance (i.e. words
// in little-endian order). This is the order that seems to be used in every
// big number library in the worlds, so...
//
// The first field in the structure (allocated) is the number of words in 'd'.
// This is the upper limit on the size of the number that can be held in the
// structure. This differs from libraries like OpenSSL as this is not intended
// to deal with numbers of arbitrary size; just numbers that are needed to deal
// with the algorithms that are defined in the TPM implementation.
//
// The second field in the structure (size) is the number of significant words
// in 'n'. When this number is zero, the number is zero. The word at used-1 should
// never be zero. All words between d[size] and d[allocated-1] should be zero.

//** Defines

#ifndef _BN_NUMBERS_H
#define _BN_NUMBERS_H

#if RADIX_BITS == 64
# define RADIX_LOG2         6
#elif RADIX_BITS == 32
#define RADIX_LOG2          5
#else
# error "Unsupported radix"
#endif

#define RADIX_MOD(x)        ((x) & ((1 << RADIX_LOG2) - 1))
#define RADIX_DIV(x)        ((x) >> RADIX_LOG2)
#define RADIX_MASK  ((((crypt_uword_t)1) << RADIX_LOG2) - 1)

#define BITS_TO_CRYPT_WORDS(bits)       RADIX_DIV((bits) + (RADIX_BITS - 1))
#define BYTES_TO_CRYPT_WORDS(bytes)     BITS_TO_CRYPT_WORDS(bytes * 8)
#define SIZE_IN_CRYPT_WORDS(thing)      BYTES_TO_CRYPT_WORDS(sizeof(thing))

#if RADIX_BITS == 64
#define SWAP_CRYPT_WORD(x)  REVERSE_ENDIAN_64(x)
    typedef uint64_t    crypt_uword_t;
    typedef int64_t     crypt_word_t;
#   define TO_CRYPT_WORD_64             BIG_ENDIAN_BYTES_TO_UINT64 
#   define TO_CRYPT_WORD_32(a, b, c, d) TO_CRYPT_WORD_64(0, 0, 0, 0, a, b, c, d)
#elif RADIX_BITS == 32
#   define SWAP_CRYPT_WORD(x)  REVERSE_ENDIAN_32((x))
    typedef uint32_t    crypt_uword_t;
    typedef int32_t     crypt_word_t;
#   define TO_CRYPT_WORD_64(a, b, c, d, e, f, g, h)                                \
        BIG_ENDIAN_BYTES_TO_UINT32(e, f, g, h),                                    \
        BIG_ENDIAN_BYTES_TO_UINT32(a, b, c, d)
#endif

#define MAX_CRYPT_UWORD (~((crypt_uword_t)0))
#define MAX_CRYPT_WORD  ((crypt_word_t)(MAX_CRYPT_UWORD >> 1))
#define MIN_CRYPT_WORD  (~MAX_CRYPT_WORD)

#define LARGEST_NUMBER (MAX((ALG_RSA * MAX_RSA_KEY_BYTES),                      \
                        MAX((ALG_ECC * MAX_ECC_KEY_BYTES), MAX_DIGEST_SIZE)))
#define LARGEST_NUMBER_BITS (LARGEST_NUMBER * 8)

#define MAX_ECC_PARAMETER_BYTES (MAX_ECC_KEY_BYTES * ALG_ECC)

// These are the basic big number formats. This is convertible to the library-
// specific format without to much difficulty. For the math performed using
// these numbers, the value is always positive.
#define BN_STRUCT_DEF(count) struct {       \
    crypt_uword_t       allocated;          \
    crypt_uword_t       size;               \
    crypt_uword_t       d[count];           \
    }

typedef BN_STRUCT_DEF(1) bignum_t;
#ifndef bigNum
typedef bignum_t       *bigNum;
typedef const bignum_t *bigConst;
#endif

extern const bignum_t   BnConstZero;

// The Functions to access the properties of a big number.
// Get number of allocated words
#define BnGetAllocated(x)   (unsigned)((x)->allocated)

// Get number of words used
#define BnGetSize(x)        ((x)->size)

// Get a pointer to the data array
#define BnGetArray(x)       ((crypt_uword_t *)&((x)->d[0]))

// Get the nth word of a BIGNUM (zero-based)
#define BnGetWord(x, i)     (crypt_uword_t)((x)->d[i])

// Some things that are done often.

// Test to see if a bignum_t is equal to zero
#define BnEqualZero(bn)   (BnGetSize(bn) == 0)

// Test to see if a bignum_t is equal to a word type
#define BnEqualWord(bn, word)                         \
            ((BnGetSize(bn) == 1) && (BnGetWord(bn, 0) == (crypt_uword_t)word))

// Determine if a BIGNUM is even. A zero is even. Although the
// indication that a number is zero is that it's size is zero,
// all words of the number are 0 so this test works on zero.
#define BnIsEven(n)     ((BnGetWord(n, 0) & 1) == 0)

// The macros below are used to define BIGNUM values of the required
// size. The values are allocated on the stack so they can be
// treated like simple local values.

// This will call the initialization function for a defined bignum_t.
// This sets the allocated and used fields and clears the words of 'n'.
#define BN_INIT(name)                                           \
    (bigNum)BnInit((bigNum)&(name),                             \
                BYTES_TO_CRYPT_WORDS(sizeof(name.d)))

// In some cases, a function will need the address of the structure
// associated with a variable. The structure for a BIGNUM variable
// of 'name' is 'name_'. Generally, when the structure is created, it
// is initialized and a parameter is created with a pointer to the
// structure. The pointer has the 'name' and the structure it points
// to is 'name_'
#define BN_ADDRESS(name) (bigNum)&name##_

#define BN_STRUCT_ALLOCATION(bits) (BITS_TO_CRYPT_WORDS(bits) + 1)

// Create a structure of the correct size.
#define BN_STRUCT(bits)                                         \
    BN_STRUCT_DEF(BN_STRUCT_ALLOCATION(bits))

// Define a BIGNUM type with a specific allocation
#define BN_TYPE(name, bits)                                     \
    typedef BN_STRUCT(bits) bn_##name##_t

// This creates a local BIGNUM variable of a specific size and
// initializes it from a TPM2B input parameter.
#define BN_INITIALIZED(name, bits, initializer)                 \
    BN_STRUCT(bits)  name##_;                                   \
    bigNum           name = BnFrom2B(BN_INIT(name##_),          \
                                    (const TPM2B *)initializer)

// Create a local variable that can hold a number with 'bits'
#define BN_VAR(name, bits)                                      \
    BN_STRUCT(bits)  _##name;                                    \
    bigNum           name = BN_INIT(_##name)

// Create a type that can hold the largest number defined by the
// implementation.
#define BN_MAX(name)   BN_VAR(name, LARGEST_NUMBER_BITS)
#define BN_MAX_INITIALIZED(name, initializer)                   \
    BN_INITIALIZED(name, LARGEST_NUMBER_BITS, initializer)

// A word size value is useful
#define BN_WORD(name)      BN_VAR(name, RADIX_BITS)

// This is used to created a word-size BIGNUM and initialize it with
// an input parameter to a function.
#define BN_WORD_INITIALIZED(name, initial)                              \
    BN_STRUCT(RADIX_BITS)  name##_;                                     \
    bigNum                 name = BnInitializeWord((bigNum)&name##_,    \
                                BN_STRUCT_ALLOCATION(RADIX_BITS), initial)

// ECC-Specific Values

// This is the format for a point. It is always in affine format. The Z value is
// carried as part of the point, primarily to simplify the interface to the support
// library. Rather than have the interface layer have to create space for the
// point each time it is used...
// The x, y, and z values are pointers to bigNum values and not in-line versions of
// the numbers. This is a relic of the days when there was no standard TPM format
// for the numbers
typedef struct _bn_point_t
{
    bigNum          x;
    bigNum          y;
    bigNum          z;
} bn_point_t;

typedef bn_point_t          *bigPoint;
typedef const bn_point_t    *pointConst;

typedef struct constant_point_t
{
    bigConst        x;
    bigConst        y;
    bigConst        z;
} constant_point_t;

#define ECC_BITS    (MAX_ECC_KEY_BYTES * 8)
BN_TYPE(ecc, ECC_BITS);
#define ECC_NUM(name)       BN_VAR(name, ECC_BITS)
#define ECC_INITIALIZED(name, initializer)                          \
    BN_INITIALIZED(name, ECC_BITS, initializer)

#define POINT_INSTANCE(name, bits)                                  \
    BN_STRUCT (bits)    name##_x =                                  \
                {BITS_TO_CRYPT_WORDS ( bits ), 0,{0}};              \
    BN_STRUCT ( bits )    name##_y =                                \
                {BITS_TO_CRYPT_WORDS ( bits ), 0,{0}};              \
    BN_STRUCT ( bits )    name##_z =                                \
                {BITS_TO_CRYPT_WORDS ( bits ), 0,{0}};              \
    bn_point_t name##_

#define POINT_INITIALIZER(name)                                      \
    BnInitializePoint(&name##_, (bigNum)&name##_x,                  \
                    (bigNum)&name##_y, (bigNum)&name##_z)

#define POINT_INITIALIZED(name, initValue)                         \
    POINT_INSTANCE(name, MAX_ECC_KEY_BITS);                         \
    bigPoint             name = BnPointFrom2B(                      \
                                    POINT_INITIALIZER(name),        \
                                    initValue)

#define POINT_VAR(name, bits)                                       \
    POINT_INSTANCE (name, bits);                                  \
    bigPoint            name = POINT_INITIALIZER(name)

#define POINT(name)      POINT_VAR(name, MAX_ECC_KEY_BITS)

// Structure for the curve parameters. This is an analog to the
// TPMS_ALGORITHM_DETAIL_ECC
typedef struct
{
    bigConst             prime;     // a prime number
    bigConst             order;     // the order of the curve
    bigConst             h;         // cofactor
    bigConst             a;         // linear coefficient
    bigConst             b;         // constant term
    constant_point_t     base;      // base point
} ECC_CURVE_DATA;

// Access macros for the ECC_CURVE structure. The parameter 'C' is a pointer
// to an ECC_CURVE_DATA structure. In some libraries, the curve structure contains
// a pointer to an ECC_CURVE_DATA structure as well as some other bits. For those
// cases, the AccessCurveData macro is used in the code to first get the pointer
// to the ECC_CURVE_DATA for access. In some cases, the macro does noting.
#define CurveGetPrime(C)    ((C)->prime)
#define CurveGetOrder(C)    ((C)->order)
#define CurveGetCofactor(C) ((C)->h)
#define CurveGet_a(C)       ((C)->a)
#define CurveGet_b(C)       ((C)->b)
#define CurveGetG(C)        ((pointConst)&((C)->base))
#define CurveGetGx(C)       ((C)->base.x)
#define CurveGetGy(C)       ((C)->base.y)


// Convert bytes in initializers according to the endianess of the system.
// This is used for CryptEccData.c.
#define     BIG_ENDIAN_BYTES_TO_UINT32(a, b, c, d)                      \
            (    ((UINT32)(a) << 24)                                    \
            +    ((UINT32)(b) << 16)                                    \
            +    ((UINT32)(c) << 8)                                     \
            +    ((UINT32)(d))                                          \
            )

#define     BIG_ENDIAN_BYTES_TO_UINT64(a, b, c, d, e, f, g, h)          \
            (    ((UINT64)(a) << 56)                                    \
            +    ((UINT64)(b) << 48)                                    \
            +    ((UINT64)(c) << 40)                                    \
            +    ((UINT64)(d) << 32)                                    \
            +    ((UINT64)(e) << 24)                                    \
            +    ((UINT64)(f) << 16)                                    \
            +    ((UINT64)(g) << 8)                                     \
            +    ((UINT64)(h))                                          \
            )

#ifndef RADIX_BYTES
#   if RADIX_BITS == 32
#       define RADIX_BYTES 4
#   elif RADIX_BITS == 64
#       define RADIX_BYTES 8
#   else
#       error "RADIX_BITS must either be 32 or 64"
#   endif
#endif

// Add implementation dependent definitions for other ECC Values and for linkages. 
#include LIB_INCLUDE(MATH_LIB, Math)


#endif // _BN_NUMBERS_H