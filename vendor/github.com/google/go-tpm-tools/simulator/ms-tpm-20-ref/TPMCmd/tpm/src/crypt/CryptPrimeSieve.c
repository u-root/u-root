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
//** Includes and defines

#include "Tpm.h"

#if RSA_KEY_SIEVE

#include "CryptPrimeSieve_fp.h"

// This determines the number of bits in the largest sieve field.
#define MAX_FIELD_SIZE  2048

extern const uint32_t      s_LastPrimeInTable;
extern const uint32_t      s_PrimeTableSize;
extern const uint32_t      s_PrimesInTable;
extern const unsigned char s_PrimeTable[];

// This table is set of prime markers. Each entry is the prime value
// for the ((n + 1) * 1024) prime. That is, the entry in s_PrimeMarkers[1]
// is the value for the 2,048th prime. This is used in the PrimeSieve
// to adjust the limit for the prime search. When processing smaller
// prime candidates, fewer primes are checked directly before going to
// Miller-Rabin. As the prime grows, it is worth spending more time eliminating
// primes as, a) the density is lower, and b) the cost of Miller-Rabin is
// higher.
const uint32_t      s_PrimeMarkersCount = 6;
const uint32_t      s_PrimeMarkers[] = {
    8167, 17881, 28183, 38891, 49871, 60961 };
uint32_t   primeLimit;

//** Functions

//*** RsaAdjustPrimeLimit()
// This used during the sieve process. The iterator for getting the
// next prime (RsaNextPrime()) will return primes until it hits the
// limit (primeLimit) set up by this function. This causes the sieve
// process to stop when an appropriate number of primes have been
// sieved.
LIB_EXPORT void
RsaAdjustPrimeLimit(
    uint32_t        requestedPrimes
    )
{
    if(requestedPrimes == 0 || requestedPrimes > s_PrimesInTable)
        requestedPrimes = s_PrimesInTable;
    requestedPrimes = (requestedPrimes - 1) / 1024;
    if(requestedPrimes < s_PrimeMarkersCount)
        primeLimit = s_PrimeMarkers[requestedPrimes];
    else
        primeLimit = s_LastPrimeInTable;
    primeLimit >>= 1;

}

//*** RsaNextPrime()
// This the iterator used during the sieve process. The input is the
// last prime returned (or any starting point) and the output is the
// next higher prime. The function returns 0 when the primeLimit is
// reached.
LIB_EXPORT uint32_t
RsaNextPrime(
    uint32_t    lastPrime
    )
{
    if(lastPrime == 0)
        return 0;
    lastPrime >>= 1;
    for(lastPrime += 1; lastPrime <= primeLimit; lastPrime++)
    {
        if(((s_PrimeTable[lastPrime >> 3] >> (lastPrime & 0x7)) & 1) == 1)
            return ((lastPrime << 1) + 1);
    }
    return 0;
}

// This table contains a previously sieved table. It has
// the bits for 3, 5, and 7 removed. Because of the
// factors, it needs to be aligned to 105 and has
// a repeat of 105.
const BYTE   seedValues[] = {
    0x16, 0x29, 0xcb, 0xa4, 0x65, 0xda, 0x30, 0x6c,
    0x99, 0x96, 0x4c, 0x53, 0xa2, 0x2d, 0x52, 0x96,
    0x49, 0xcb, 0xb4, 0x61, 0xd8, 0x32, 0x2d, 0x99,
    0xa6, 0x44, 0x5b, 0xa4, 0x2c, 0x93, 0x96, 0x69,
    0xc3, 0xb0, 0x65, 0x5a, 0x32, 0x4d, 0x89, 0xb6,
    0x48, 0x59, 0x26, 0x2d, 0xd3, 0x86, 0x61, 0xcb,
    0xb4, 0x64, 0x9a, 0x12, 0x6d, 0x91, 0xb2, 0x4c,
    0x5a, 0xa6, 0x0d, 0xc3, 0x96, 0x69, 0xc9, 0x34,
    0x25, 0xda, 0x22, 0x65, 0x99, 0xb4, 0x4c, 0x1b,
    0x86, 0x2d, 0xd3, 0x92, 0x69, 0x4a, 0xb4, 0x45,
    0xca, 0x32, 0x69, 0x99, 0x36, 0x0c, 0x5b, 0xa6,
    0x25, 0xd3, 0x94, 0x68, 0x8b, 0x94, 0x65, 0xd2,
    0x32, 0x6d, 0x18, 0xb6, 0x4c, 0x4b, 0xa6, 0x29,
    0xd1};

#define USE_NIBBLE

#ifndef USE_NIBBLE
static const BYTE bitsInByte[256] = {
    0x00, 0x01, 0x01, 0x02, 0x01, 0x02, 0x02, 0x03,
    0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04,
    0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
    0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
    0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
    0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
    0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
    0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
    0x05, 0x06, 0x06, 0x07, 0x06, 0x07, 0x07, 0x08
};
#define BitsInByte(x)   bitsInByte[(unsigned char)x]
#else
const BYTE bitsInNibble[16] = {
    0x00, 0x01, 0x01, 0x02, 0x01, 0x02, 0x02, 0x03,
    0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04};
#define BitsInByte(x)                                       \
            (bitsInNibble[(unsigned char)(x) & 0xf]         \
        +   bitsInNibble[((unsigned char)(x) >> 4) & 0xf])
#endif

//*** BitsInArry()
// This function counts the number of bits set in an array of bytes.
static int
BitsInArray(
    const unsigned char     *a,             // IN: A pointer to an array of bytes
    unsigned int             aSize          // IN: the number of bytes to sum
    )
{
    int     j = 0;
    for(; aSize; a++, aSize--)
        j += BitsInByte(*a);
    return j;
}

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
    )
{
    UINT16       i;
    int          retValue;
    UINT32       sum = 0;
    BYTE         sel;

    //find the bit
    for(i = 0; (i < (int)aSize) && (sum < n); i++)
        sum += BitsInByte(a[i]);
    i--;
    // The chosen bit is in the byte that was just accessed
    // Compute the offset to the start of that byte
    retValue = i * 8 - 1;
    sel = a[i];
    // Subtract the bits in the last byte added.
    sum -= BitsInByte(sel);
    // Now process the byte, one bit at a time.
    for(; (sel != 0) && (sum != n); retValue++, sel = sel >> 1)
        sum += (sel & 1) != 0;
    return (sum == n) ? retValue : -1;
}

typedef struct
{
    UINT16      prime;
    UINT16      count;
} SIEVE_MARKS;

const SIEVE_MARKS sieveMarks[5] = {
    {31, 7}, {73, 5}, {241, 4}, {1621, 3}, {UINT16_MAX, 2}};


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
    )
{
    UINT32           i;
    UINT32           j;
    UINT32           fieldBits = fieldSize * 8;
    UINT32           r;
    BYTE            *pField;
    INT32            iter;
    UINT32           adjust;
    UINT32           mark = 0;
    UINT32           count = sieveMarks[0].count;
    UINT32           stop = sieveMarks[0].prime;
    UINT32           composite;
    UINT32           pList[8];
    UINT32           next;

    pAssert(field != NULL && bnN != NULL);

    // If the remainder is odd, then subtracting the value will give an even number,
    // but we want an odd number, so subtract the 105+rem. Otherwise, just subtract
    // the even remainder.
    adjust = (UINT32)BnModWord(bnN, 105);
    if(adjust & 1)
        adjust += 105;

    // Adjust the input number so that it points to the first number in a
    // aligned field.
    BnSubWord(bnN, bnN, adjust);
//    pAssert(BnModWord(bnN, 105) == 0);
    pField = field;
    for(i = fieldSize; i >= sizeof(seedValues);
        pField += sizeof(seedValues), i -= sizeof(seedValues))
    {
        memcpy(pField, seedValues, sizeof(seedValues));
    }
    if(i != 0)
        memcpy(pField, seedValues, i);

    // Cycle through the primes, clearing bits
    // Have already done 3, 5, and 7
    iter = 7;

#define NEXT_PRIME(iter)    (iter = RsaNextPrime(iter))
    // Get the next N primes where N is determined by the mark in the sieveMarks
    while((composite = NEXT_PRIME(iter)) != 0)
    {
        next = 0;
        i = count;
        pList[i--] = composite;
        for(; i > 0; i--)
        {
            next = NEXT_PRIME(iter);
            pList[i] = next;
            if(next != 0)
                composite *= next;
        }
        // Get the remainder when dividing the base field address
        // by the composite
        composite = (UINT32)BnModWord(bnN, composite);
        // 'composite' is divisible by the composite components. for each of the
        // composite components, divide 'composite'. That remainder (r) is used to
        // pick a starting point for clearing the array. The stride is equal to the
        // composite component. Note, the field only contains odd numbers. If the
        // field were expanded to contain all numbers, then half of the bits would
        // have already been cleared. We can save the trouble of clearing them a
        // second time by having a stride of 2*next. Or we can take all of the even
        // numbers out of the field and use a stride of 'next'
        for(i = count; i > 0; i--)
        {
            next = pList[i];
            if(next == 0)
                goto done;
            r = composite % next;
        // these computations deal with the fact that we have picked a field-sized
        // range that is aligned to a 105 count boundary. The problem is, this field
        // only contains odd numbers. If we take our prime guess and walk through all 
        // the numbers using that prime as the 'stride', then every other 'stride' is
        // going to be an even number. So, we are actually counting by 2 * the stride
        // We want the count to start on an odd number at the start of our field. That
        // is, we want to assume that we have counted up to the edge of the field by
        // the 'stride' and now we are going to start flipping bits in the field as we
        // continue to count up by 'stride'. If we take the base of our field and
        // divide by the stride, we find out how much we find out how short the last
        // count was from reaching the edge of the bit field. Say we get a quotient of
        // 3 and remainder of 1. This means that after 3 strides, we are 1 short of
        // the start of the field and the next stride will either land within the
        // field or step completely over it. The confounding factor is that our field 
        // only contains odd numbers and our stride is actually 2 * stride. If the
        // quoitent is even, then that means that when we add 2 * stride, we are going
        // to hit another even number. So, we have to know if we need to back off
        // by 1 stride before we start couting by 2 * stride. 
        // We can tell from the remainder whether we are on an even or odd
        // stride when we hit the beginning of the table. If we are on an odd stride
        // (r & 1), we would start half a stride in (next - r)/2. If we are on an
        // even stride, we need 0.5 strides (next - r/2) because the table only has
        // odd numbers. If the remainder happens to be zero, then the start of the
        // table is on stride so no adjustment is necessary.
            if(r & 1)           j = (next - r) / 2;
            else if(r == 0)     j = 0;
            else                 j = next - (r / 2); 
            for(; j < fieldBits; j += next)
                ClearBit(j, field, fieldSize);
        }
        if(next >= stop)
        {
            mark++;
            count = sieveMarks[mark].count;
            stop = sieveMarks[mark].prime;
        }
    }
done:
    INSTRUMENT_INC(totalFieldsSieved[PrimeIndex]);
    i = BitsInArray(field, fieldSize);
    INSTRUMENT_ADD(bitsInFieldAfterSieve[PrimeIndex], i);
    INSTRUMENT_ADD(emptyFieldsSieved[PrimeIndex], (i == 0));
    return i;
}



#ifdef SIEVE_DEBUG
static uint32_t fieldSize = 210;

//***SetFieldSize()
// Function to set the field size used for prime generation. Used for tuning.
LIB_EXPORT uint32_t
SetFieldSize(
    uint32_t         newFieldSize
    )
{
    if(newFieldSize == 0 || newFieldSize > MAX_FIELD_SIZE)
        fieldSize = MAX_FIELD_SIZE;
    else
        fieldSize = newFieldSize;
    return fieldSize;
}
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
    )
{
    BYTE             field[MAX_FIELD_SIZE];
    UINT32           first;
    UINT32           ones;
    INT32            chosen;
    BN_PRIME(test);
    UINT32           modE;
#ifndef SIEVE_DEBUG
    UINT32           fieldSize = MAX_FIELD_SIZE;
#endif
    UINT32           primeSize;
//
    // Adjust the field size and prime table list to fit the size of the prime
    // being tested. This is done to try to optimize the trade-off between the 
    // dividing done for sieving and the time for Miller-Rabin. When the size
    // of the prime is large, the cost of Miller-Rabin is fairly high, as is the
    // cost of the sieving. However, the time for Miller-Rabin goes up considerably
    // faster than the cost of dividing by a number of primes.
    primeSize = BnSizeInBits(candidate);

    if(primeSize <= 512)
    {
        RsaAdjustPrimeLimit(1024); // Use just the first 1024 primes
    }
    else if(primeSize <= 1024)
    {
        RsaAdjustPrimeLimit(4096); // Use just the first 4K primes
    }
    else
    {
        RsaAdjustPrimeLimit(0);     // Use all available
    }

    // Save the low-order word to use as a search generator and make sure that
    // it has some interesting range to it
    first = (UINT32)(candidate->d[0] | 0x80000000);

    // Sieve the field
    ones = PrimeSieve(candidate, fieldSize, field);
    pAssert(ones > 0 && ones < (fieldSize * 8));
    for(; ones > 0; ones--)
    {
        // Decide which bit to look at and find its offset
        chosen = FindNthSetBit((UINT16)fieldSize, field, ((first % ones) + 1));

        if((chosen < 0) || (chosen >= (INT32)(fieldSize * 8)))
            FAIL(FATAL_ERROR_INTERNAL);

        // Set this as the trial prime
        BnAddWord(test, candidate, (crypt_uword_t)(chosen * 2));

        // The exponent might not have been one of the tested primes so
        // make sure that it isn't divisible and make sure that 0 != (p-1) mod e
        // Note: This is the same as 1 != p mod e 
        modE = (UINT32)BnModWord(test, e);
        if((modE != 0) && (modE != 1) && MillerRabin(test, rand))
        {
            BnCopy(candidate, test);
            return TPM_RC_SUCCESS;
        }
        // Clear the bit just tested
        ClearBit(chosen, field, fieldSize);
    }
    // Ran out of bits and couldn't find a prime in this field
    INSTRUMENT_INC(noPrimeFields[PrimeIndex]);
    return (g_inFailureMode ? TPM_RC_FAILURE : TPM_RC_NO_RESULT);
}

#if RSA_INSTRUMENT
static char            a[256];

//*** PrintTuple()
char *
PrintTuple(
    UINT32      *i
    )
{
    sprintf(a, "{%d, %d, %d}", i[0], i[1], i[2]);
    return a;
}

#define CLEAR_VALUE(x)    memset(x, 0, sizeof(x))

//*** RsaSimulationEnd()
void
RsaSimulationEnd(
    void
    )
{
    int         i;
    UINT32      averages[3];
    UINT32      nonFirst = 0;
    if((PrimeCounts[0] + PrimeCounts[1] + PrimeCounts[2]) != 0)
    {
        printf("Primes generated = %s\n", PrintTuple(PrimeCounts));
        printf("Fields sieved = %s\n", PrintTuple(totalFieldsSieved));
        printf("Fields with no primes = %s\n", PrintTuple(noPrimeFields));
        printf("Primes checked with Miller-Rabin = %s\n",
               PrintTuple(MillerRabinTrials));
        for(i = 0; i < 3; i++)
            averages[i] = (totalFieldsSieved[i]
                           != 0 ? bitsInFieldAfterSieve[i] / totalFieldsSieved[i]
                           : 0);
        printf("Average candidates in field %s\n", PrintTuple(averages));
        for(i = 1; i < (sizeof(failedAtIteration) / sizeof(failedAtIteration[0])); 
        i++)
            nonFirst += failedAtIteration[i];
        printf("Miller-Rabin failures not in first round = %d\n", nonFirst);
            
    }
    CLEAR_VALUE(PrimeCounts);
    CLEAR_VALUE(totalFieldsSieved);
    CLEAR_VALUE(noPrimeFields);
    CLEAR_VALUE(MillerRabinTrials);
    CLEAR_VALUE(bitsInFieldAfterSieve);
}

//*** GetSieveStats()
LIB_EXPORT void
GetSieveStats(
    uint32_t        *trials,
    uint32_t        *emptyFields,
    uint32_t        *averageBits
    )
{
    uint32_t        totalBits;
    uint32_t        fields;
    *trials = MillerRabinTrials[0] + MillerRabinTrials[1] + MillerRabinTrials[2];
    *emptyFields = noPrimeFields[0] + noPrimeFields[1] + noPrimeFields[2];
    fields = totalFieldsSieved[0] + totalFieldsSieved[1] 
        + totalFieldsSieved[2];
    totalBits = bitsInFieldAfterSieve[0] + bitsInFieldAfterSieve[1] 
        + bitsInFieldAfterSieve[2];
    if(fields != 0)
        *averageBits = totalBits / fields;
    else
        *averageBits = 0;
    CLEAR_VALUE(PrimeCounts);
    CLEAR_VALUE(totalFieldsSieved);
    CLEAR_VALUE(noPrimeFields);
    CLEAR_VALUE(MillerRabinTrials);
    CLEAR_VALUE(bitsInFieldAfterSieve);

}
#endif

#endif // RSA_KEY_SIEVE

#if !RSA_INSTRUMENT

//*** RsaSimulationEnd()
// Stub for call when not doing instrumentation. 
void
RsaSimulationEnd(
    void
    )
{
    return;
}
#endif