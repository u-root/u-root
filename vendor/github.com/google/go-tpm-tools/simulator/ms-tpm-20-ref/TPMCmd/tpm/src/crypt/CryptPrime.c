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
// This file contains the code for prime validation.

#include "Tpm.h"
#include "CryptPrime_fp.h"

//#define CPRI_PRIME
//#include "PrimeTable.h"

#include "CryptPrimeSieve_fp.h"

extern const uint32_t      s_LastPrimeInTable;
extern const uint32_t      s_PrimeTableSize;
extern const uint32_t      s_PrimesInTable;
extern const unsigned char s_PrimeTable[];
extern bigConst            s_CompositeOfSmallPrimes;

//** Functions

//*** Root2()
// This finds ceil(sqrt(n)) to use as a stopping point for searching the prime
// table.
static uint32_t
Root2(
    uint32_t             n
    )
{
    int32_t              last = (int32_t)(n >> 2);
    int32_t              next = (int32_t)(n >> 1);
    int32_t              diff;
    int32_t              stop = 10;
//
    // get a starting point
    for(; next != 0; last >>= 1, next >>= 2);
    last++;
    do
    {
        next = (last + (n / last)) >> 1;
        diff = next - last;
        last = next;
        if(stop-- == 0)
            FAIL(FATAL_ERROR_INTERNAL);
    } while(diff < -1 || diff > 1);
    if((n / next) > (unsigned)next)
        next++;
    pAssert(next != 0);
    pAssert(((n / next) <= (unsigned)next) && (n / (next + 1) < (unsigned)next));
    return next;
}

//*** IsPrimeInt()
// This will do a test of a word of up to 32-bits in size.
BOOL
IsPrimeInt(
    uint32_t            n
    )
{
    uint32_t            i;
    uint32_t            stop;
    if(n < 3 || ((n & 1) == 0))
        return (n == 2);
    if(n <= s_LastPrimeInTable)
    {
        n >>= 1;
        return ((s_PrimeTable[n >> 3] >> (n & 7)) & 1);
    }
    // Need to search
    stop = Root2(n) >> 1;
    // starting at 1 is equivalent to staring at  (1 << 1) + 1 = 3
    for(i = 1; i < stop; i++)
    {
        if((s_PrimeTable[i >> 3] >> (i & 7)) & 1)
            // see if this prime evenly divides the number
            if((n % ((i << 1) + 1)) == 0)
                return FALSE;
    }
    return TRUE;
}

//*** BnIsProbablyPrime()
// This function is used when the key sieve is not implemented. This function
// Will try to eliminate some of the obvious things before going on
// to perform MillerRabin as a final verification of primeness.
BOOL
BnIsProbablyPrime(
    bigNum          prime,           // IN:
    RAND_STATE      *rand            // IN: the random state just
                                     //     in case Miller-Rabin is required
    )
{
#if RADIX_BITS > 32
    if(BnUnsignedCmpWord(prime, UINT32_MAX) <= 0)
#else
    if(BnGetSize(prime) == 1)
#endif
        return IsPrimeInt((uint32_t)prime->d[0]);

    if(BnIsEven(prime))
        return FALSE;
    if(BnUnsignedCmpWord(prime, s_LastPrimeInTable) <= 0)
    {
        crypt_uword_t       temp = prime->d[0] >> 1;
        return ((s_PrimeTable[temp >> 3] >> (temp & 7)) & 1);
    }
    {
        BN_VAR(n, LARGEST_NUMBER_BITS);
        BnGcd(n, prime, s_CompositeOfSmallPrimes);
        if(!BnEqualWord(n, 1))
            return FALSE;
    }
    return MillerRabin(prime, rand);
}

//*** MillerRabinRounds()
// Function returns the number of Miller-Rabin rounds necessary to give an
// error probability equal to the security strength of the prime. These values
// are from FIPS 186-3.
UINT32
MillerRabinRounds(
    UINT32           bits           // IN: Number of bits in the RSA prime
    )
{
    if(bits < 511) return 8;    // don't really expect this
    if(bits < 1536) return 5;   // for 512 and 1K primes
    return 4;                   // for 3K public modulus and greater
}

//*** MillerRabin()
// This function performs a Miller-Rabin test from FIPS 186-3. It does
// 'iterations' trials on the number. In all likelihood, if the number
// is not prime, the first test fails.
//  Return Type: BOOL
//      TRUE(1)         probably prime
//      FALSE(0)        composite
BOOL
MillerRabin(
    bigNum           bnW,
    RAND_STATE      *rand
    )
{
    BN_MAX(bnWm1);
    BN_PRIME(bnM);
    BN_PRIME(bnB);
    BN_PRIME(bnZ);
    BOOL             ret = FALSE;   // Assumed composite for easy exit
    unsigned int     a;
    unsigned int     j;
    int              wLen;
    int              i;
    int              iterations = MillerRabinRounds(BnSizeInBits(bnW));
//
   INSTRUMENT_INC(MillerRabinTrials[PrimeIndex]);

    pAssert(bnW->size > 1);
    // Let a be the largest integer such that 2^a divides w1.
    BnSubWord(bnWm1, bnW, 1);
    pAssert(bnWm1->size != 0);

    // Since w is odd (w-1) is even so start at bit number 1 rather than 0
    // Get the number of bits in bnWm1 so that it doesn't have to be recomputed
    // on each iteration.
    i = (int)(bnWm1->size * RADIX_BITS);
  // Now find the largest power of 2 that divides w1
    for(a = 1;
    (a < (bnWm1->size * RADIX_BITS)) &&
        (BnTestBit(bnWm1, a) == 0);
        a++);
 // 2. m = (w1) / 2^a
        BnShiftRight(bnM, bnWm1, a);
      // 3. wlen = len (w).
    wLen = BnSizeInBits(bnW);
  // 4. For i = 1 to iterations do
    for(i = 0; i < iterations; i++)
    {
        // 4.1 Obtain a string b of wlen bits from an RBG.
        // Ensure that 1 < b < w1.
        // 4.2 If ((b <= 1) or (b >= w1)), then go to step 4.1.
        while(BnGetRandomBits(bnB, wLen, rand) && ((BnUnsignedCmpWord(bnB, 1) <= 0)
            || (BnUnsignedCmp(bnB, bnWm1) >= 0)));
        if(g_inFailureMode)
            return FALSE;

       // 4.3 z = b^m mod w.
       // if ModExp fails, then say this is not
       // prime and bail out.
        BnModExp(bnZ, bnB, bnM, bnW);

        // 4.4 If ((z == 1) or (z = w == 1)), then go to step 4.7.
        if((BnUnsignedCmpWord(bnZ, 1) == 0)
           || (BnUnsignedCmp(bnZ, bnWm1) == 0))
            goto step4point7;
        // 4.5 For j = 1 to a  1 do.
        for(j = 1; j < a; j++)
        {
          // 4.5.1 z = z^2 mod w.
            BnModMult(bnZ, bnZ, bnZ, bnW);
          // 4.5.2 If (z = w1), then go to step 4.7.
            if(BnUnsignedCmp(bnZ, bnWm1) == 0)
                goto step4point7;
          // 4.5.3 If (z = 1), then go to step 4.6.
            if(BnEqualWord(bnZ, 1))
                goto step4point6;
        }
        // 4.6 Return COMPOSITE.
step4point6:
        INSTRUMENT_INC(failedAtIteration[i]);
        goto end;
      // 4.7 Continue. Comment: Increment i for the do-loop in step 4.
step4point7:
        continue;
    }
  // 5. Return PROBABLY PRIME
    ret = TRUE;
end:
    return ret;
}

#if ALG_RSA

//*** RsaCheckPrime()
// This will check to see if a number is prime and appropriate for an
// RSA prime.
//
// This has different functionality based on whether we are using key
// sieving or not. If not, the number checked to see if it is divisible by
// the public exponent, then the number is adjusted either up or down
// in order to make it a better candidate. It is then checked for being
// probably prime.
//
// If sieving is used, the number is used to root a sieving process.
//
TPM_RC
RsaCheckPrime(
    bigNum           prime,
    UINT32           exponent,
    RAND_STATE      *rand
    )
{
#if !RSA_KEY_SIEVE
    TPM_RC          retVal = TPM_RC_SUCCESS;
    UINT32          modE = BnModWord(prime, exponent);

    NOT_REFERENCED(rand);

    if(modE == 0)
        // evenly divisible so add two keeping the number odd
        BnAddWord(prime, prime, 2);
    // want 0 != (p - 1) mod e
    // which is 1 != p mod e
    else if(modE == 1)
        // subtract 2 keeping number odd and insuring that
        // 0 != (p - 1) mod e
        BnSubWord(prime, prime, 2);

    if(BnIsProbablyPrime(prime, rand) == 0)
        ERROR_RETURN(g_inFailureMode ? TPM_RC_FAILURE : TPM_RC_VALUE);
Exit:
    return retVal;
#else
    return PrimeSelectWithSieve(prime, exponent, rand);
#endif
}

//*** AdjustPrimeCandiate()
// For this math, we assume that the RSA numbers are fixed-point numbers with
// the decimal point to the "left" of the most significant bit. This approach helps
// make it clear what is happening with the MSb of the values.
// The two RSA primes have to be large enough so that their product will be a number 
// with the necessary number of significant bits. For example, we want to be able 
// to multiply two 1024-bit numbers to produce a number with 2028 significant bits. If
// we accept any 1024-bit prime that has its MSb set, then it is possible to produce a
// product that does not have the MSb SET. For example, if we use tiny keys of 16 bits
// and have two 8-bit 'primes' of 0x80, then the public key would be 0x4000 which is
// only 15-bits. So, what we need to do is made sure that each of the primes is large
// enough so that the product of the primes is twice as large as each prime. A little
// arithmetic will show that the only way to do this is to make sure that each of the
// primes is no less than root(2)/2. That's what this functions does.
// This function adjusts the candidate prime so that it is odd and >= root(2)/2.
// This allows the product of these two numbers to be .5, which, in fixed point
// notation means that the most significant bit is 1.
// For this routine, the root(2)/2 (0.7071067811865475) approximated with 0xB505 
// which is, in fixed point, 0.7071075439453125 or an error of 0.000108%. Just setting
// the upper two bits would give a value > 0.75 which is an error of > 6%. Given the 
// amount of time all the other computations take, reducing the error is not much of
// a cost, but it isn't totally required either.
//
// This function can be replaced with a function that just sets the two most 
// significant bits of each prime candidate without introducing any computational 
// issues. 
//
//
LIB_EXPORT void
RsaAdjustPrimeCandidate(
    bigNum          prime
    )
{
    UINT32          msw;
    UINT32          adjusted;

    // If the radix is 32, the compiler should turn this into a simple assignment
    msw = prime->d[prime->size - 1] >> ((RADIX_BITS == 64) ? 32 : 0);
    // Multiplying 0xff...f by 0x4AFB gives 0xff..f - 0xB5050...0
    adjusted = (msw >> 16) * 0x4AFB;
    adjusted += ((msw & 0xFFFF) * 0x4AFB) >> 16;
    adjusted += 0xB5050000UL;
#if RADIX_BITS == 64
    // Save the low-order 32 bits
    prime->d[prime->size - 1] &= 0xFFFFFFFFUL;
    // replace the upper 32-bits
    prime->d[prime->size -1] |= ((crypt_uword_t)adjusted << 32);
#else
    prime->d[prime->size - 1] = (crypt_uword_t)adjusted;
#endif
    // make sure the number is odd
    prime->d[0] |= 1;
}

//***BnGeneratePrimeForRSA()
// Function to generate a prime of the desired size with the proper attributes
// for an RSA prime.
TPM_RC
BnGeneratePrimeForRSA(
    bigNum          prime,          // IN/OUT: points to the BN that will get the 
                                    //  random value
    UINT32          bits,           // IN: number of bits to get
    UINT32          exponent,       // IN: the exponent 
    RAND_STATE      *rand           // IN: the random state
    )
{
    BOOL            found = FALSE;
//
    // Make sure that the prime is large enough
    pAssert(prime->allocated >= BITS_TO_CRYPT_WORDS(bits));
    // Only try to handle specific sizes of keys in order to save overhead
    pAssert((bits % 32) == 0);
    prime->size = BITS_TO_CRYPT_WORDS(bits);
    while(!found)
    {
// The change below is to make sure that all keys that are generated from the same
// seed value will be the same regardless of the endianess or word size of the CPU.
//       DRBG_Generate(rand, (BYTE *)prime->d, (UINT16)BITS_TO_BYTES(bits));// old
//       if(g_inFailureMode)                                                // old
       if(!BnGetRandomBits(prime, bits, rand))                              // new
            return TPM_RC_FAILURE;
        RsaAdjustPrimeCandidate(prime);
        found = RsaCheckPrime(prime, exponent, rand) == TPM_RC_SUCCESS;
    }
    return TPM_RC_SUCCESS;
}

#endif // ALG_RSA