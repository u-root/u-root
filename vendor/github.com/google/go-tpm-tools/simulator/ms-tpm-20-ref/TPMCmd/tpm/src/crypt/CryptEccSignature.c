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
//** Includes and Defines
#include "Tpm.h"
#include "CryptEccSignature_fp.h"

#if ALG_ECC

//** Utility Functions

//*** EcdsaDigest()
// Function to adjust the digest so that it is no larger than the order of the
// curve. This is used for ECDSA sign and verification.
static bigNum
EcdsaDigest(
    bigNum               bnD,           // OUT: the adjusted digest
    const TPM2B_DIGEST  *digest,        // IN: digest to adjust
    bigConst             max            // IN: value that indicates the maximum
                                        //     number of bits in the results
    )
{
    int              bitsInMax = BnSizeInBits(max);
    int              shift;
//
    if(digest == NULL)
        BnSetWord(bnD, 0);
    else
    {
        BnFromBytes(bnD, digest->t.buffer, 
                    (NUMBYTES)MIN(digest->t.size, BITS_TO_BYTES(bitsInMax)));
        shift = BnSizeInBits(bnD) - bitsInMax;
        if(shift > 0)
            BnShiftRight(bnD, bnD, shift);
    }
    return bnD;
}

//*** BnSchnorrSign()
// This contains the Schnorr signature computation. It is used by both ECDSA and
// Schnorr signing. The result is computed as: ['s' = 'k' + 'r' * 'd' (mod 'n')] 
// where
// 1) 's' is the signature
// 2) 'k' is a random value
// 3) 'r' is the value to sign
// 4) 'd' is the private EC key
// 5) 'n' is the order of the curve
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        the result of the operation was zero or 'r' (mod 'n') 
//                              is zero
static TPM_RC
BnSchnorrSign(
    bigNum                   bnS,           // OUT: 's' component of the signature
    bigConst                 bnK,           // IN: a random value
    bigNum                   bnR,           // IN: the signature 'r' value
    bigConst                 bnD,           // IN: the private key
    bigConst                 bnN            // IN: the order of the curve
    )
{
    // Need a local temp value to store the intermediate computation because product
    // size can be larger than will fit in bnS.
    BN_VAR(bnT1, MAX_ECC_PARAMETER_BYTES * 2 * 8);
//
    // Reduce bnR without changing the input value
    BnDiv(NULL, bnT1, bnR, bnN);
    if(BnEqualZero(bnT1))
        return TPM_RC_NO_RESULT;
    // compute s = (k + r * d)(mod n)
        // r * d
    BnMult(bnT1, bnT1, bnD);
        // k * r * d
    BnAdd(bnT1, bnT1, bnK);
        // k + r * d (mod n)
    BnDiv(NULL, bnS, bnT1, bnN);
    return (BnEqualZero(bnS)) ? TPM_RC_NO_RESULT : TPM_RC_SUCCESS;
}

//** Signing Functions

//*** BnSignEcdsa()
// This function implements the ECDSA signing algorithm. The method is described
// in the comments below.
TPM_RC
BnSignEcdsa(
    bigNum                   bnR,           // OUT: 'r' component of the signature
    bigNum                   bnS,           // OUT: 's' component of the signature
    bigCurve                 E,             // IN: the curve used in the signature
                                            //     process
    bigNum                   bnD,           // IN: private signing key
    const TPM2B_DIGEST      *digest,        // IN: the digest to sign
    RAND_STATE              *rand           // IN: used in debug of signing
    )
{
    ECC_NUM(bnK);
    ECC_NUM(bnIk);
    BN_VAR(bnE, MAX(MAX_ECC_KEY_BYTES, MAX_DIGEST_SIZE) * 8);
    POINT(ecR);
    bigConst                order = CurveGetOrder(AccessCurveData(E));
    TPM_RC                  retVal = TPM_RC_SUCCESS;
    INT32                   tries = 10;
    BOOL                    OK = FALSE;
//
    pAssert(digest != NULL);
    // The algorithm as described in "Suite B Implementer's Guide to FIPS
    // 186-3(ECDSA)"
    // 1. Use one of the routines in Appendix A.2 to generate (k, k^-1), a
    //    per-message secret number and its inverse modulo n. Since n is prime,
    //    the output will be invalid only if there is a failure in the RBG.
    // 2. Compute the elliptic curve point R = [k]G = (xR, yR) using EC scalar
    //    multiplication (see [Routines]), where G is the base point included in
    //    the set of domain parameters.
    // 3. Compute r = xR mod n. If r = 0, then return to Step 1. 1.
    // 4. Use the selected hash function to compute H = Hash(M).
    // 5. Convert the bit string H to an integer e as described in Appendix B.2.
    // 6. Compute s = (k^-1 *  (e + d *  r)) mod q. If s = 0, return to Step 1.2.
    // 7. Return (r, s).
    // In the code below, q is n (that it, the order of the curve is p)

    do // This implements the loop at step 6. If s is zero, start over.
    {
        for(; tries > 0; tries--)
        {
            // Step 1 and 2 -- generate an ephemeral key and the modular inverse
            // of the private key.
            if(!BnEccGenerateKeyPair(bnK, ecR, E, rand))
                continue;
         // x coordinate is mod p.  Make it mod q
            BnMod(ecR->x, order);
            // Make sure that it is not zero;
            if(BnEqualZero(ecR->x))
                continue;
            // write the modular reduced version of r as part of the signature
            BnCopy(bnR, ecR->x);
            // Make sure that a modular inverse exists and try again if not
            OK = (BnModInverse(bnIk, bnK, order));
            if(OK)
                break;
        }
        if(!OK)
            goto Exit;

        EcdsaDigest(bnE, digest, order);

        // now have inverse of K (bnIk), e (bnE), r (bnR),  d (bnD) and
        // CurveGetOrder(E)
        // Compute s = k^-1 (e + r*d)(mod q)
        //  first do s = r*d mod q
        BnModMult(bnS, bnR, bnD, order);
        // s = e + s = e + r * d
        BnAdd(bnS, bnE, bnS);
        // s = k^(-1)s (mod n) = k^(-1)(e + r * d)(mod n)
        BnModMult(bnS, bnIk, bnS, order);

        // If S is zero, try again
    } while(BnEqualZero(bnS));
Exit:
    return retVal;
}

#if ALG_ECDAA

//*** BnSignEcdaa()
//
// This function performs 's' = 'r' + 'T' * 'd' mod 'q' where
// 1) 'r is a random, or pseudo-random value created in the commit phase
// 2) 'nonceK' is a TPM-generated, random value 0 < 'nonceK' < 'n'
// 3) 'T' is mod 'q' of "Hash"('nonceK' || 'digest'), and
// 4) 'd' is a private key.
//
// The signature is the tuple ('nonceK', 's')
//
// Regrettably, the parameters in this function kind of collide with the parameter
// names used in ECSCHNORR making for a lot of confusion.
//  Return Type: TPM_RC
//      TPM_RC_SCHEME       unsupported hash algorithm
//      TPM_RC_NO_RESULT    cannot get values from random number generator
static TPM_RC
BnSignEcdaa(
    TPM2B_ECC_PARAMETER     *nonceK,        // OUT: 'nonce' component of the signature
    bigNum                   bnS,           // OUT: 's' component of the signature
    bigCurve                 E,             // IN: the curve used in signing
    bigNum                   bnD,           // IN: the private key
    const TPM2B_DIGEST      *digest,        // IN: the value to sign (mod 'q')
    TPMT_ECC_SCHEME         *scheme,        // IN: signing scheme (contains the
                                            //      commit count value).
    OBJECT                  *eccKey,        // IN: The signing key
    RAND_STATE              *rand           // IN: a random number state
    )
{
    TPM_RC                   retVal;
    TPM2B_ECC_PARAMETER      r;
    HASH_STATE               state;
    TPM2B_DIGEST             T;
    BN_MAX(bnT);
//
    NOT_REFERENCED(rand);
    if(!CryptGenerateR(&r, &scheme->details.ecdaa.count,
                       eccKey->publicArea.parameters.eccDetail.curveID,
                       &eccKey->name))
        retVal = TPM_RC_VALUE;
    else
    {
        // This allocation is here because 'r' doesn't have a value until 
        // CrypGenerateR() is done.
        ECC_INITIALIZED(bnR, &r);
        do
        {
            // generate nonceK such that 0 < nonceK < n
            // use bnT as a temp.
            if(!BnEccGetPrivate(bnT, AccessCurveData(E), rand))
            {
                retVal = TPM_RC_NO_RESULT;
                break;
            }
            BnTo2B(bnT, &nonceK->b, 0);

            T.t.size = CryptHashStart(&state, scheme->details.ecdaa.hashAlg);
            if(T.t.size == 0)
            {
                retVal = TPM_RC_SCHEME;
            }
            else
            {
                CryptDigestUpdate2B(&state, &nonceK->b);
                CryptDigestUpdate2B(&state, &digest->b);
                CryptHashEnd2B(&state, &T.b);
                BnFrom2B(bnT, &T.b);
                // Watch out for the name collisions in this call!!
                retVal = BnSchnorrSign(bnS, bnR, bnT, bnD,
                                       AccessCurveData(E)->order);
            }
        } while(retVal == TPM_RC_NO_RESULT);
        // Because the rule is that internal state is not modified if the command 
        // fails, only end the commit if the command succeeds.
        // NOTE that if the result of the Schnorr computation was zero 
        // it will probably not be worthwhile to run the same command again because
        // the result will still be zero. This means that the Commit command will 
        // need to be run again to get a new commit value for the signature.
        if(retVal == TPM_RC_SUCCESS)
            CryptEndCommit(scheme->details.ecdaa.count);
    }
    return retVal;
}
#endif // ALG_ECDAA

#if ALG_ECSCHNORR

//*** SchnorrReduce()
// Function to reduce a hash result if it's magnitude is too large. The size of
// 'number' is set so that it has no more bytes of significance than 'reference'
// value. If the resulting number can have more bits of significance than
// 'reference'.
static void
SchnorrReduce(
    TPM2B       *number,        // IN/OUT: Value to reduce
    bigConst     reference      // IN: the reference value
    )
{
    UINT16      maxBytes = (UINT16)BITS_TO_BYTES(BnSizeInBits(reference));
    if(number->size > maxBytes)
        number->size = maxBytes;
}

//*** SchnorrEcc()
// This function is used to perform a modified Schnorr signature.
//
// This function will generate a random value 'k' and compute
// a) ('xR', 'yR') = ['k']'G'
// b) 'r' = "Hash"('xR' || 'P')(mod 'q')
// c) 'rT' = truncated 'r'
// d) 's'= 'k' + 'rT' * 'ds' (mod 'q')
// e) return the tuple 'rT', 's'
//
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        failure in the Schnorr sign process
//      TPM_RC_SCHEME           hashAlg can't produce zero-length digest
static TPM_RC
BnSignEcSchnorr(
    bigNum                   bnR,           // OUT: 'r' component of the signature
    bigNum                   bnS,           // OUT: 's' component of the signature
    bigCurve                 E,             // IN: the curve used in signing
    bigNum                   bnD,           // IN: the signing key
    const TPM2B_DIGEST      *digest,        // IN: the digest to sign
    TPM_ALG_ID               hashAlg,       // IN: signing scheme (contains a hash)
    RAND_STATE              *rand           // IN: non-NULL when testing
    )
{
    HASH_STATE               hashState;
    UINT16                   digestSize = CryptHashGetDigestSize(hashAlg);
    TPM2B_TYPE(T, MAX(MAX_DIGEST_SIZE, MAX_ECC_KEY_BYTES));
    TPM2B_T                  T2b;
    TPM2B                   *e = &T2b.b;
    TPM_RC                   retVal = TPM_RC_NO_RESULT;
    const ECC_CURVE_DATA    *C;
    bigConst                 order;
    bigConst                 prime;
    ECC_NUM(bnK);
    POINT(ecR);
//
    // Parameter checks
    if(E == NULL)
        ERROR_RETURN(TPM_RC_VALUE);
    C = AccessCurveData(E);
    order = CurveGetOrder(C);
    prime = CurveGetOrder(C);

    // If the digest does not produce a hash, then null the signature and return
    // a failure.
    if(digestSize == 0)
    {
        BnSetWord(bnR, 0);
        BnSetWord(bnS, 0);
        ERROR_RETURN(TPM_RC_SCHEME);
    }
    do
    {
        // Generate a random key pair
        if(!BnEccGenerateKeyPair(bnK, ecR, E, rand))
            break;
        // Convert R.x to a string
        BnTo2B(ecR->x, e, (NUMBYTES)BITS_TO_BYTES(BnSizeInBits(prime)));

        // f) compute r = Hash(e || P) (mod n)
        CryptHashStart(&hashState, hashAlg);
        CryptDigestUpdate2B(&hashState, e);
        CryptDigestUpdate2B(&hashState, &digest->b);
        e->size = CryptHashEnd(&hashState, digestSize, e->buffer);
        // Reduce the hash size if it is larger than the curve order
        SchnorrReduce(e, order);
        // Convert hash to number
        BnFrom2B(bnR, e);
        // Do the Schnorr computation
        retVal = BnSchnorrSign(bnS, bnK, bnR, bnD, CurveGetOrder(C));
    } while(retVal == TPM_RC_NO_RESULT);
Exit:
    return retVal;
}

#endif // ALG_ECSCHNORR

#if ALG_SM2
#ifdef  _SM2_SIGN_DEBUG

//*** BnHexEqual()
// This function compares a bignum value to a hex string.
//  Return Type: BOOL
//      TRUE(1)         values equal
//      FALSE(0)        values not equal
static BOOL
BnHexEqual(
    bigNum           bn,        //IN: big number value
    const char      *c          //IN: character string number
    )
{
    ECC_NUM(bnC);
    BnFromHex(bnC, c);
    return (BnUnsignedCmp(bn, bnC) == 0);
}
#endif // _SM2_SIGN_DEBUG

//*** BnSignEcSm2()
// This function signs a digest using the method defined in SM2 Part 2. The method
// in the standard will add a header to the message to be signed that is a hash of
// the values that define the key. This then hashed with the message to produce a
// digest ('e'). This function signs 'e'.
//  Return Type: TPM_RC
//      TPM_RC_VALUE         bad curve
static TPM_RC
BnSignEcSm2(
    bigNum                   bnR,       // OUT: 'r' component of the signature
    bigNum                   bnS,       // OUT: 's' component of the signature
    bigCurve                 E,         // IN: the curve used in signing
    bigNum                   bnD,       // IN: the private key
    const TPM2B_DIGEST      *digest,    // IN: the digest to sign
    RAND_STATE              *rand       // IN: random number generator (mostly for
                                        //     debug)
    )
{
    BN_MAX_INITIALIZED(bnE, digest);    // Don't know how big digest might be 
    ECC_NUM(bnN);
    ECC_NUM(bnK);
    ECC_NUM(bnT);                       // temp
    POINT(Q1);
    bigConst                  order = (E != NULL)
        ? CurveGetOrder(AccessCurveData(E)) : NULL;
//
#ifdef _SM2_SIGN_DEBUG
    BnFromHex(bnE, "B524F552CD82B8B028476E005C377FB1"
                   "9A87E6FC682D48BB5D42E3D9B9EFFE76");
    BnFromHex(bnD, "128B2FA8BD433C6C068C8D803DFF7979"
                   "2A519A55171B1B650C23661D15897263");
#endif
    // A3: Use random number generator to generate random number 1 <= k <= n-1;
    // NOTE: Ax: numbers are from the SM2 standard
loop:
    {
        // Get a random number 0 < k < n
        BnGenerateRandomInRange(bnK, order, rand);
#ifdef _SM2_SIGN_DEBUG
        BnFromHex(bnK, "6CB28D99385C175C94F94E934817663F"
                       "C176D925DD72B727260DBAAE1FB2F96F");
#endif
        // A4: Figure out the point of elliptic curve (x1, y1)=[k]G, and according
        // to details specified in 4.2.7 in Part 1 of this document, transform the
        // data type of x1 into an integer;
        if(!BnEccModMult(Q1, NULL, bnK, E))
            goto loop;
        // A5: Figure out 'r' = ('e' + 'x1') mod 'n',
        BnAdd(bnR, bnE, Q1->x);
        BnMod(bnR, order);
#ifdef _SM2_SIGN_DEBUG
        pAssert(BnHexEqual(bnR, "40F1EC59F793D9F49E09DCEF49130D41"
                                "94F79FB1EED2CAA55BACDB49C4E755D1"));
#endif
        // if r=0 or r+k=n, return to A3;
        if(BnEqualZero(bnR))
            goto loop;
        BnAdd(bnT, bnK, bnR);
        if(BnUnsignedCmp(bnT, bnN) == 0)
            goto loop;
        // A6: Figure out s = ((1 + dA)^-1  (k - r  dA)) mod n,
        // if s=0, return to A3;
        // compute t = (1+dA)^-1
        BnAddWord(bnT, bnD, 1);
        BnModInverse(bnT, bnT, order);
#ifdef _SM2_SIGN_DEBUG
        pAssert(BnHexEqual(bnT, "79BFCF3052C80DA7B939E0C6914A18CB"
                                "B2D96D8555256E83122743A7D4F5F956"));
#endif
        // compute s = t * (k - r * dA) mod n
        BnModMult(bnS, bnR, bnD, order);
        // k - r * dA mod n = k + n - ((r * dA) mod n)
        BnSub(bnS, order, bnS);
        BnAdd(bnS, bnK, bnS);
        BnModMult(bnS, bnS, bnT, order);
#ifdef _SM2_SIGN_DEBUG
        pAssert(BnHexEqual(bnS, "6FC6DAC32C5D5CF10C77DFB20F7C2EB6"
                                "67A457872FB09EC56327A67EC7DEEBE7"));
#endif
        if(BnEqualZero(bnS))
            goto loop;
    }
    // A7: According to details specified in 4.2.1 in Part 1 of this document,
    // transform the data type of r, s into bit strings, signature of message M
    // is (r, s).
    // This is handled by the common return code
#ifdef _SM2_SIGN_DEBUG
    pAssert(BnHexEqual(bnR, "40F1EC59F793D9F49E09DCEF49130D41"
                            "94F79FB1EED2CAA55BACDB49C4E755D1"));
    pAssert(BnHexEqual(bnS, "6FC6DAC32C5D5CF10C77DFB20F7C2EB6"
                            "67A457872FB09EC56327A67EC7DEEBE7"));
#endif
    return TPM_RC_SUCCESS;
}
#endif // ALG_SM2

//*** CryptEccSign()
// This function is the dispatch function for the various ECC-based
// signing schemes.
// There is a bit of ugliness to the parameter passing. In order to test this,
// we sometime would like to use a deterministic RNG so that we can get the same
// signatures during testing. The easiest way to do this for most schemes is to
// pass in a deterministic RNG and let it return canned values during testing.
// There is a competing need for a canned parameter to use in ECDAA. To accommodate
// both needs with minimal fuss, a special type of RAND_STATE is defined to carry
// the address of the commit value. The setup and handling of this is not very
// different for the caller than what was in previous versions of the code.
//  Return Type: TPM_RC
//      TPM_RC_SCHEME            'scheme' is not supported
LIB_EXPORT TPM_RC
CryptEccSign(
    TPMT_SIGNATURE          *signature,     // OUT: signature
    OBJECT                  *signKey,       // IN: ECC key to sign the hash
    const TPM2B_DIGEST      *digest,        // IN: digest to sign
    TPMT_ECC_SCHEME         *scheme,        // IN: signing scheme
    RAND_STATE              *rand
    )
{
    CURVE_INITIALIZED(E, signKey->publicArea.parameters.eccDetail.curveID);
    ECC_INITIALIZED(bnD, &signKey->sensitive.sensitive.ecc.b);
    ECC_NUM(bnR);
    ECC_NUM(bnS);
    const ECC_CURVE_DATA   *C;
    TPM_RC                  retVal = TPM_RC_SCHEME;
//
    NOT_REFERENCED(scheme);
    if(E == NULL)
        ERROR_RETURN(TPM_RC_VALUE);
     C = AccessCurveData(E);
    signature->signature.ecdaa.signatureR.t.size
        = sizeof(signature->signature.ecdaa.signatureR.t.buffer);
    signature->signature.ecdaa.signatureS.t.size
        = sizeof(signature->signature.ecdaa.signatureS.t.buffer);
    TEST(signature->sigAlg);
    switch(signature->sigAlg)
    {
        case ALG_ECDSA_VALUE:
            retVal = BnSignEcdsa(bnR, bnS, E, bnD, digest, rand);
            break;
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            retVal = BnSignEcdaa(&signature->signature.ecdaa.signatureR, bnS, E, 
                                 bnD, digest, scheme, signKey, rand);
            bnR = NULL;
            break;
#endif
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            retVal = BnSignEcSchnorr(bnR, bnS, E, bnD, digest,
                                     signature->signature.ecschnorr.hash,
                                     rand);
            break;
#endif
#if ALG_SM2
        case ALG_SM2_VALUE:
            retVal = BnSignEcSm2(bnR, bnS, E, bnD, digest, rand);
            break;
#endif
        default:
            break;
    }
    // If signature generation worked, convert the results.
    if(retVal == TPM_RC_SUCCESS)
    {
        NUMBYTES     orderBytes =
            (NUMBYTES)BITS_TO_BYTES(BnSizeInBits(CurveGetOrder(C)));
        if(bnR != NULL)
            BnTo2B(bnR, &signature->signature.ecdaa.signatureR.b, orderBytes);
        if(bnS != NULL)
            BnTo2B(bnS, &signature->signature.ecdaa.signatureS.b, orderBytes);
    }
Exit:
    CURVE_FREE(E);
    return retVal;
}

//********************* Signature Validation   ********************

#if ALG_ECDSA 

//*** BnValidateSignatureEcdsa()
// This function validates an ECDSA signature. rIn and sIn should have been checked
// to make sure that they are in the range 0 < 'v' < 'n'
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE           signature not valid
TPM_RC
BnValidateSignatureEcdsa(
    bigNum                   bnR,           // IN: 'r' component of the signature
    bigNum                   bnS,           // IN: 's' component of the signature
    bigCurve                 E,             // IN: the curve used in the signature
                                            //     process
    bn_point_t              *ecQ,           // IN: the public point of the key
    const TPM2B_DIGEST      *digest         // IN: the digest that was signed
    )
{
    // Make sure that the allocation for the digest is big enough for a maximum
    // digest
    BN_VAR(bnE, MAX(MAX_ECC_KEY_BYTES, MAX_DIGEST_SIZE) * 8);
    POINT(ecR);
    ECC_NUM(bnU1);
    ECC_NUM(bnU2);
    ECC_NUM(bnW);
    bigConst                 order = CurveGetOrder(AccessCurveData(E));
    TPM_RC                   retVal = TPM_RC_SIGNATURE;
//
    // Get adjusted digest
    EcdsaDigest(bnE, digest, order);
    // 1. If r and s are not both integers in the interval [1, n - 1], output
    //    INVALID.
    //  bnR  and bnS were validated by the caller
    // 2. Use the selected hash function to compute H0 = Hash(M0).
    // This is an input parameter
    // 3. Convert the bit string H0 to an integer e as described in Appendix B.2.
    // Done at entry
    // 4. Compute w = (s')^-1 mod n, using the routine in Appendix B.1.
    if(!BnModInverse(bnW, bnS, order))
        goto Exit;
    // 5. Compute u1 = (e' *   w) mod n, and compute u2 = (r' *  w) mod n.
    BnModMult(bnU1, bnE, bnW, order);
    BnModMult(bnU2, bnR, bnW, order);
    // 6. Compute the elliptic curve point R = (xR, yR) = u1G+u2Q, using EC
    //    scalar multiplication and EC addition (see [Routines]). If R is equal to
    //    the point at infinity O, output INVALID.
    if(BnPointMult(ecR, CurveGetG(AccessCurveData(E)), bnU1, ecQ, bnU2, E)
       != TPM_RC_SUCCESS)
        goto Exit;
    // 7. Compute v = Rx mod n.
    BnMod(ecR->x, order);
    // 8. Compare v and r0. If v = r0, output VALID; otherwise, output INVALID
    if(BnUnsignedCmp(ecR->x, bnR) != 0)
        goto Exit;

    retVal = TPM_RC_SUCCESS;
Exit:
    return retVal;
}

#endif      // ALG_ECDSA

#if ALG_SM2

//*** BnValidateSignatureEcSm2()
// This function is used to validate an SM2 signature.
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE            signature not valid
static TPM_RC
BnValidateSignatureEcSm2(
    bigNum                   bnR,       // IN: 'r' component of the signature
    bigNum                   bnS,       // IN: 's' component of the signature
    bigCurve                 E,         // IN: the curve used in the signature
                                        //     process
    bigPoint                 ecQ,       // IN: the public point of the key
    const TPM2B_DIGEST      *digest     // IN: the digest that was signed
    )
{
    POINT(P);
    ECC_NUM(bnRp);
    ECC_NUM(bnT);
    BN_MAX_INITIALIZED(bnE, digest);
    BOOL                     OK;
    bigConst                 order = CurveGetOrder(AccessCurveData(E));

#ifdef _SM2_SIGN_DEBUG
    // Make sure that the input signature is the test signature
    pAssert(BnHexEqual(bnR,
                       "40F1EC59F793D9F49E09DCEF49130D41"
                       "94F79FB1EED2CAA55BACDB49C4E755D1"));
    pAssert(BnHexEqual(bnS,
                       "6FC6DAC32C5D5CF10C77DFB20F7C2EB6"
                       "67A457872FB09EC56327A67EC7DEEBE7"));
#endif
    // b)   compute t  := (r + s) mod n
    BnAdd(bnT, bnR, bnS);
    BnMod(bnT, order);
#ifdef _SM2_SIGN_DEBUG
    pAssert(BnHexEqual(bnT,
                       "2B75F07ED7ECE7CCC1C8986B991F441A"
                       "D324D6D619FE06DD63ED32E0C997C801"));
#endif
    // c)   verify that t > 0
    OK = !BnEqualZero(bnT);
    if(!OK)
        // set T to a value that should allow rest of the computations to run
        // without trouble
        BnCopy(bnT, bnS);
    // d)   compute (x, y) := [s]G + [t]Q
    OK = BnEccModMult2(P, NULL, bnS, ecQ, bnT, E);
#ifdef  _SM2_SIGN_DEBUG
    pAssert(OK && BnHexEqual(P->x,
                             "110FCDA57615705D5E7B9324AC4B856D"
                             "23E6D9188B2AE47759514657CE25D112"));
#endif
    // e)   compute r' := (e + x) mod n (the x coordinate is in bnT)
    OK = OK && BnAdd(bnRp, bnE, P->x);
    OK = OK && BnMod(bnRp, order);

    // f)   verify that r' = r
    OK = OK && (BnUnsignedCmp(bnR, bnRp) == 0);

    if(!OK)
        return TPM_RC_SIGNATURE;
    else
        return TPM_RC_SUCCESS;
}

#endif  // ALG_SM2

#if ALG_ECSCHNORR

//*** BnValidateSignatureEcSchnorr()
// This function is used to validate an EC Schnorr signature.
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE        signature not valid
static TPM_RC
BnValidateSignatureEcSchnorr(
    bigNum               bnR,       // IN: 'r' component of the signature
    bigNum               bnS,       // IN: 's' component of the signature
    TPM_ALG_ID           hashAlg,   // IN: hash algorithm of the signature
    bigCurve             E,         // IN: the curve used in the signature
                                    //     process
    bigPoint             ecQ,       // IN: the public point of the key
    const TPM2B_DIGEST  *digest     // IN: the digest that was signed
    )
{
    BN_MAX(bnRn);
    POINT(ecE);
    BN_MAX(bnEx);
    const ECC_CURVE_DATA    *C = AccessCurveData(E);
    bigConst                 order = CurveGetOrder(C);
    UINT16                   digestSize = CryptHashGetDigestSize(hashAlg);
    HASH_STATE               hashState;
    TPM2B_TYPE(BUFFER, MAX(MAX_ECC_PARAMETER_BYTES, MAX_DIGEST_SIZE));
    TPM2B_BUFFER             Ex2 = {{sizeof(Ex2.t.buffer),{ 0 }}};
    BOOL                     OK;
//
    // E = [s]G - [r]Q
    BnMod(bnR, order);
    // Make -r = n - r
    BnSub(bnRn, order, bnR);
    // E = [s]G + [-r]Q
    OK = BnPointMult(ecE, CurveGetG(C), bnS, ecQ, bnRn, E) == TPM_RC_SUCCESS;
//   // reduce the x portion of E mod q
//    OK = OK && BnMod(ecE->x, order);
    // Convert to byte string
    OK = OK && BnTo2B(ecE->x, &Ex2.b,
                      (NUMBYTES)(BITS_TO_BYTES(BnSizeInBits(order))));
    if(OK)
    {
// Ex = h(pE.x || digest)
        CryptHashStart(&hashState, hashAlg);
        CryptDigestUpdate(&hashState, Ex2.t.size, Ex2.t.buffer);
        CryptDigestUpdate(&hashState, digest->t.size, digest->t.buffer);
        Ex2.t.size = CryptHashEnd(&hashState, digestSize, Ex2.t.buffer);
        SchnorrReduce(&Ex2.b, order);
        BnFrom2B(bnEx, &Ex2.b);
        // see if Ex matches R
        OK = BnUnsignedCmp(bnEx, bnR) == 0;
    }
    return (OK) ? TPM_RC_SUCCESS : TPM_RC_SIGNATURE;
}
#endif  // ALG_ECSCHNORR

//*** CryptEccValidateSignature()
// This function validates an EcDsa or EcSchnorr signature.
// The point 'Qin' needs to have been validated to be on the curve of 'curveId'.
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE            not a valid signature
LIB_EXPORT TPM_RC
CryptEccValidateSignature(
    TPMT_SIGNATURE          *signature,     // IN: signature to be verified
    OBJECT                  *signKey,       // IN: ECC key signed the hash
    const TPM2B_DIGEST      *digest         // IN: digest that was signed
    )
{
    CURVE_INITIALIZED(E, signKey->publicArea.parameters.eccDetail.curveID);
    ECC_NUM(bnR);
    ECC_NUM(bnS);
    POINT_INITIALIZED(ecQ, &signKey->publicArea.unique.ecc);
    bigConst                 order;
    TPM_RC                   retVal;
    
    if(E == NULL)
        ERROR_RETURN(TPM_RC_VALUE);

    order = CurveGetOrder(AccessCurveData(E));

//    // Make sure that the scheme is valid
    switch(signature->sigAlg)
    {
        case ALG_ECDSA_VALUE:
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
#endif
#if ALG_SM2
        case ALG_SM2_VALUE:
#endif
            break;
        default:
            ERROR_RETURN(TPM_RC_SCHEME);
            break;
    }
    // Can convert r and s after determining that the scheme is an ECC scheme. If
    // this conversion doesn't work, it means that the unmarshaling code for 
    // an ECC signature is broken.
    BnFrom2B(bnR, &signature->signature.ecdsa.signatureR.b);
    BnFrom2B(bnS, &signature->signature.ecdsa.signatureS.b);

    // r and s have to be greater than 0 but less than the curve order
    if(BnEqualZero(bnR) || BnEqualZero(bnS))
        ERROR_RETURN(TPM_RC_SIGNATURE);
    if((BnUnsignedCmp(bnS, order) >= 0)
       || (BnUnsignedCmp(bnR, order) >= 0))
        ERROR_RETURN(TPM_RC_SIGNATURE);

    switch(signature->sigAlg)
    {
        case ALG_ECDSA_VALUE:
            retVal = BnValidateSignatureEcdsa(bnR, bnS, E, ecQ, digest);
            break;

#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            retVal = BnValidateSignatureEcSchnorr(bnR, bnS,
                                                  signature->signature.any.hashAlg, 
                                                  E, ecQ, digest);
            break;
#endif
#if ALG_SM2
        case ALG_SM2_VALUE:
            retVal = BnValidateSignatureEcSm2(bnR, bnS, E, ecQ, digest);
            break;
#endif
        default:
            FAIL(FATAL_ERROR_INTERNAL);
    }
Exit:
    CURVE_FREE(E);
    return retVal;
}

//***CryptEccCommitCompute()
// This function performs the point multiply operations required by TPM2_Commit.
//
// If 'B' or 'M' is provided, they must be on the curve defined by 'curveId'. This
// routine does not check that they are on the curve and results are unpredictable
// if they are not.
//
// It is a fatal error if 'r' is NULL. If 'B' is not NULL, then it is a
// fatal error if 'd' is NULL or if 'K' and 'L' are both NULL. 
// If 'M' is not NULL, then it is a fatal error if 'E' is NULL.
//
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        if 'K', 'L' or 'E' was computed to be the point
//                              at infinity
//      TPM_RC_CANCELED         a cancel indication was asserted during this
//                              function
LIB_EXPORT TPM_RC
CryptEccCommitCompute(
    TPMS_ECC_POINT          *K,             // OUT: [d]B or [r]Q
    TPMS_ECC_POINT          *L,             // OUT: [r]B
    TPMS_ECC_POINT          *E,             // OUT: [r]M
    TPM_ECC_CURVE            curveId,       // IN: the curve for the computations
    TPMS_ECC_POINT          *M,             // IN: M (optional)
    TPMS_ECC_POINT          *B,             // IN: B (optional)
    TPM2B_ECC_PARAMETER     *d,             // IN: d (optional)
    TPM2B_ECC_PARAMETER     *r              // IN: the computed r value (required)
    )
{
    CURVE_INITIALIZED(curve, curveId);  // Normally initialize E as the curve, but 
                                        // E means something else in this function
    ECC_INITIALIZED(bnR, r);
    TPM_RC               retVal = TPM_RC_SUCCESS;
//
    // Validate that the required parameters are provided.
    // Note: E has to be provided if computing E := [r]Q or E := [r]M. Will do
    // E := [r]Q if both M and B are NULL.
    pAssert(r != NULL && E != NULL);

    // Initialize the output points in case they are not computed
    ClearPoint2B(K);
    ClearPoint2B(L);
    ClearPoint2B(E);

    // Sizes of the r parameter may not be zero
    pAssert(r->t.size > 0);

    // If B is provided, compute K=[d]B and L=[r]B
    if(B != NULL)
    {
        ECC_INITIALIZED(bnD, d);
        POINT_INITIALIZED(pB, B);
        POINT(pK);
        POINT(pL);
//
        pAssert(d != NULL && K != NULL && L != NULL);

        if(!BnIsOnCurve(pB, AccessCurveData(curve)))
            ERROR_RETURN(TPM_RC_VALUE);
         // do the math for K = [d]B
        if((retVal = BnPointMult(pK, pB, bnD, NULL, NULL, curve)) != TPM_RC_SUCCESS)
            goto Exit;
        // Convert BN K to TPM2B K
        BnPointTo2B(K, pK, curve);
        //  compute L= [r]B after checking for cancel
        if(_plat__IsCanceled())
            ERROR_RETURN(TPM_RC_CANCELED);
        // compute L = [r]B
        if(!BnIsValidPrivateEcc(bnR, curve))
            ERROR_RETURN(TPM_RC_VALUE);
        if((retVal = BnPointMult(pL, pB, bnR, NULL, NULL, curve)) != TPM_RC_SUCCESS)
            goto Exit;
        // Convert BN L to TPM2B L
        BnPointTo2B(L, pL, curve);
    }
    if((M != NULL) || (B == NULL))
    {
        POINT_INITIALIZED(pM, M);
        POINT(pE);
//
        // Make sure that a place was provided for the result
        pAssert(E != NULL);

        // if this is the third point multiply, check for cancel first
        if((B != NULL) && _plat__IsCanceled())
            ERROR_RETURN(TPM_RC_CANCELED);

        // If M provided, then pM will not be NULL and will compute E = [r]M. 
        // However, if M was not provided, then pM will be NULL and E = [r]G 
        // will be computed
        if((retVal = BnPointMult(pE, pM, bnR, NULL, NULL, curve)) != TPM_RC_SUCCESS)
            goto Exit;
        // Convert E to 2B format
        BnPointTo2B(E, pE, curve);
    }
Exit:
    CURVE_FREE(curve);
    return retVal;
}

#endif  // ALG_ECC