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
 *  Date: Apr  2, 2019  Time: 04:06:42PM
 */

#ifndef    _CRYPT_PRIME_SIEVE_FP_H_
#define    _CRYPT_PRIME_SIEVE_FP_H_

#if RSA_KEY_SIEVE

//*** RsaAdjustPrimeLimit()
// This used during the sieve process. The iterator for getting the
// next prime (RsaNextPrime()) will return primes until it hits the
// limit (primeLimit) set up by this function. This causes the sieve
// process to stop when an appropriate number of primes have been
// sieved.
LIB_EXPORT void
RsaAdjustPrimeLimit(
    uint32_t        requestedPrimes
);

//*** RsaNextPrime()
// This the iterator used during the sieve process. The input is the
// last prime returned (or any starting point) and the output is the
// next higher prime. The function returns 0 when the primeLimit is
// reached.
LIB_EXPORT uint32_t
RsaNextPrime(
    uint32_t    lastPrime
);

//*** FindNthSetBit()
// This function finds the nth SET bit in a bit array. The 'n' parameter is
// between 1 and the number of bits in the array (always a multiple of 8).
// If called when the array does not have n bits set, it will return -1
//  Return Type: unsigned int
//      <0      no bit is set or no bit with the requested number is set
//      >=0    the number of the bit in the array that is the nth set
LIB_EXPORT int
FindNthSetBit(
    const UINT16     aSize,         // IN: the size of the array to check
    const BYTE      *a,             // IN: the array to check
    const UINT32     n              // IN, the number of the SET bit
);

//*** PrimeSieve()
// This function does a prime sieve over the input 'field' which has as its
// starting address the value in bnN. Since this initializes the Sieve
// using a precomputed field with the bits associated with 3, 5 and 7 already
// turned off, the value of pnN may need to be adjusted by a few counts to allow
// the precomputed field to be used without modification.
//
// To get better performance, one could address the issue of developing the
// composite numbers. When the size of the prime gets large, the time for doing
// the divisions goes up, noticeably. It could be better to develop larger composite
// numbers even if they need to be bigNum's themselves. The object would be to
// reduce the number of times that the large prime is divided into a few large
// divides and then use smaller divides to get to the final 16 bit (or smaller)
// remainders.
LIB_EXPORT UINT32
PrimeSieve(
    bigNum           bnN,       // IN/OUT: number to sieve
    UINT32           fieldSize, // IN: size of the field area in bytes
    BYTE            *field      // IN: field
);
#ifdef SIEVE_DEBUG

//***SetFieldSize()
// Function to set the field size used for prime generation. Used for tuning.
LIB_EXPORT uint32_t
SetFieldSize(
    uint32_t         newFieldSize
);
#endif // SIEVE_DEBUG

//*** PrimeSelectWithSieve()
// This function will sieve the field around the input prime candidate. If the
// sieve field is not empty, one of the one bits in the field is chosen for testing
// with Miller-Rabin. If the value is prime, 'pnP' is updated with this value
// and the function returns success. If this value is not prime, another
// pseudo-random candidate is chosen and tested. This process repeats until
// all values in the field have been checked. If all bits in the field have
// been checked and none is prime, the function returns FALSE and a new random
// value needs to be chosen.
//  Return Type: TPM_RC
//      TPM_RC_FAILURE      TPM in failure mode, probably due to entropy source
//      TPM_RC_SUCCESS      candidate is probably prime
//      TPM_RC_NO_RESULT    candidate is not prime and couldn't find and alternative
//                          in the field
LIB_EXPORT TPM_RC
PrimeSelectWithSieve(
    bigNum           candidate,         // IN/OUT: The candidate to filter
    UINT32           e,                 // IN: the exponent
    RAND_STATE      *rand               // IN: the random number generator state
);
#if RSA_INSTRUMENT

char *
PrintTuple(
    UINT32      *i
);

void
RsaSimulationEnd(
    void
);

LIB_EXPORT void
GetSieveStats(
    uint32_t        *trials,
    uint32_t        *emptyFields,
    uint32_t        *averageBits
);

#endif
#endif // RSA_KEY_SIEVE
#if !RSA_INSTRUMENT
void
RsaSimulationEnd(
    void
);
#endif

#endif  // _CRYPT_PRIME_SIEVE_FP_H_
