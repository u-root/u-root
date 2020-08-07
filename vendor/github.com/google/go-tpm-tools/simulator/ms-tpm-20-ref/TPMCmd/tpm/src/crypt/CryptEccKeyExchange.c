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
// This file contains the functions that are used for the two-phase, ECC, 
// key-exchange protocols


#include "Tpm.h"

#if CC_ZGen_2Phase == YES

//** Functions

#if ALG_ECMQV

//*** avf1()
// This function does the associated value computation required by MQV key
// exchange.
// Process:
// 1. Convert 'xQ' to an integer 'xqi' using the convention specified in Appendix C.3.
// 2. Calculate
//        xqm = xqi mod 2^ceil(f/2) (where f = ceil(log2(n)).
// 3. Calculate the associate value function
//        avf(Q) = xqm + 2ceil(f / 2)
// Always returns TRUE(1).
static BOOL
avf1(
    bigNum               bnX,           // IN/OUT: the reduced value
    bigNum               bnN            // IN: the order of the curve
    )
{
// compute f = 2^(ceil(ceil(log2(n)) / 2))
    int                      f = (BnSizeInBits(bnN) + 1) / 2;
// x' = 2^f + (x mod 2^f)
    BnMaskBits(bnX, f);   // This is mod 2*2^f but it doesn't matter because
                            // the next operation will SET the extra bit anyway
    BnSetBit(bnX, f);
    return TRUE;
}

//*** C_2_2_MQV()
// This function performs the key exchange defined in SP800-56A
// 6.1.1.4 Full MQV, C(2, 2, ECC MQV).
//
// CAUTION: Implementation of this function may require use of essential claims in
// patents not owned by TCG members.
//
// Points 'QsB' and 'QeB' are required to be on the curve of 'inQsA'. The function 
// will fail, possibly catastrophically, if this is not the case.
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        the value for dsA does not give a valid point on the
//                              curve
static TPM_RC
C_2_2_MQV(
    TPMS_ECC_POINT          *outZ,         // OUT: the computed point
    TPM_ECC_CURVE            curveId,      // IN: the curve for the computations
    TPM2B_ECC_PARAMETER     *dsA,          // IN: static private TPM key
    TPM2B_ECC_PARAMETER     *deA,          // IN: ephemeral private TPM key
    TPMS_ECC_POINT          *QsB,          // IN: static public party B key
    TPMS_ECC_POINT          *QeB           // IN: ephemeral public party B key
    )
{
    CURVE_INITIALIZED(E, curveId);
    const ECC_CURVE_DATA    *C;
    POINT(pQeA);
    POINT_INITIALIZED(pQeB, QeB);
    POINT_INITIALIZED(pQsB, QsB);
    ECC_NUM(bnTa);
    ECC_INITIALIZED(bnDeA, deA);
    ECC_INITIALIZED(bnDsA, dsA);
    ECC_NUM(bnN);
    ECC_NUM(bnXeB);
    TPM_RC                 retVal;
//
    // Parameter checks
    if(E == NULL)
        ERROR_RETURN(TPM_RC_VALUE);
    pAssert(outZ != NULL && pQeB != NULL && pQsB != NULL && deA != NULL 
            && dsA != NULL);
    C = AccessCurveData(E);
// Process:
//  1. implicitsigA = (de,A + avf(Qe,A)ds,A ) mod n.
//  2. P = h(implicitsigA)(Qe,B + avf(Qe,B)Qs,B).
//  3. If P = O, output an error indicator.
//  4. Z=xP, where xP is the x-coordinate of P.

    // Compute the public ephemeral key pQeA = [de,A]G
    if((retVal = BnPointMult(pQeA, CurveGetG(C), bnDeA, NULL, NULL, E))
       != TPM_RC_SUCCESS)
        goto Exit;

//  1. implicitsigA = (de,A + avf(Qe,A)ds,A ) mod n.
//  tA := (ds,A + de,A  avf(Xe,A)) mod n    (3)
//  Compute 'tA' = ('deA' +  'dsA'  avf('XeA')) mod n
    // Ta = avf(XeA);
    BnCopy(bnTa, pQeA->x);
    avf1(bnTa, bnN);
    // do Ta = ds,A * Ta mod n = dsA * avf(XeA) mod n
    BnModMult(bnTa, bnDsA, bnTa, bnN);
    // now Ta = deA + Ta mod n =  deA + dsA * avf(XeA) mod n
    BnAdd(bnTa, bnTa, bnDeA);
    BnMod(bnTa, bnN);

//  2. P = h(implicitsigA)(Qe,B + avf(Qe,B)Qs,B).
// Put this in because almost every case of h is == 1 so skip the call when
    // not necessary.
    if(!BnEqualWord(CurveGetCofactor(C), 1))
        // Cofactor is not 1 so compute Ta := Ta * h mod n
        BnModMult(bnTa, bnTa, CurveGetCofactor(C), CurveGetOrder(C));

    // Now that 'tA' is (h * 'tA' mod n)
    // 'outZ' = (tA)(Qe,B + avf(Qe,B)Qs,B).

    // first, compute XeB = avf(XeB)
    avf1(bnXeB, bnN);

    // QsB := [XeB]QsB
    BnPointMult(pQsB, pQsB, bnXeB, NULL, NULL, E);
    BnEccAdd(pQeB, pQeB, pQsB, E);

    // QeB := [tA]QeB = [tA](QsB + [Xe,B]QeB) and check for at infinity
    // If the result is not the point at infinity, return QeB
    BnPointMult(pQeB, pQeB, bnTa, NULL, NULL, E);
    if(BnEqualZero(pQeB->z))
        ERROR_RETURN(TPM_RC_NO_RESULT);
    // Convert BIGNUM E to TPM2B E
    BnPointTo2B(outZ, pQeB, E);

Exit:
    CURVE_FREE(E);
    return retVal;
}

#endif // ALG_ECMQV

//*** C_2_2_ECDH()
// This function performs the two phase key exchange defined in SP800-56A,
// 6.1.1.2 Full Unified Model, C(2, 2, ECC CDH).
//
static TPM_RC
C_2_2_ECDH(
    TPMS_ECC_POINT          *outZs,         // OUT: Zs
    TPMS_ECC_POINT          *outZe,         // OUT: Ze
    TPM_ECC_CURVE            curveId,       // IN: the curve for the computations
    TPM2B_ECC_PARAMETER     *dsA,           // IN: static private TPM key
    TPM2B_ECC_PARAMETER     *deA,           // IN: ephemeral private TPM key
    TPMS_ECC_POINT          *QsB,           // IN: static public party B key
    TPMS_ECC_POINT          *QeB            // IN: ephemeral public party B key
    )
{
    CURVE_INITIALIZED(E, curveId);
    ECC_INITIALIZED(bnAs, dsA);
    ECC_INITIALIZED(bnAe, deA);
    POINT_INITIALIZED(ecBs, QsB);
    POINT_INITIALIZED(ecBe, QeB);
    POINT(ecZ);
    TPM_RC            retVal;
//
    // Parameter checks
    if(E == NULL)
        ERROR_RETURN(TPM_RC_CURVE);
    pAssert(outZs != NULL && dsA != NULL && deA != NULL && QsB != NULL 
            && QeB != NULL);

    // Do the point multiply for the Zs value ([dsA]QsB)
    retVal = BnPointMult(ecZ, ecBs, bnAs, NULL, NULL, E);
    if(retVal == TPM_RC_SUCCESS)
    {
        // Convert the Zs value.
        BnPointTo2B(outZs, ecZ, E);
        // Do the point multiply for the Ze value ([deA]QeB)
        retVal = BnPointMult(ecZ, ecBe, bnAe, NULL, NULL, E);
        if(retVal == TPM_RC_SUCCESS)
            BnPointTo2B(outZe, ecZ, E);
    }
Exit:
    CURVE_FREE(E);
    return retVal;
}

//*** CryptEcc2PhaseKeyExchange()
// This function is the dispatch routine for the EC key exchange functions that use
// two ephemeral and two static keys.
//  Return Type: TPM_RC
//      TPM_RC_SCHEME             scheme is not defined
LIB_EXPORT TPM_RC
CryptEcc2PhaseKeyExchange(
    TPMS_ECC_POINT          *outZ1,         // OUT: a computed point
    TPMS_ECC_POINT          *outZ2,         // OUT: and optional second point
    TPM_ECC_CURVE            curveId,       // IN: the curve for the computations
    TPM_ALG_ID               scheme,        // IN: the key exchange scheme
    TPM2B_ECC_PARAMETER     *dsA,           // IN: static private TPM key
    TPM2B_ECC_PARAMETER     *deA,           // IN: ephemeral private TPM key
    TPMS_ECC_POINT          *QsB,           // IN: static public party B key
    TPMS_ECC_POINT          *QeB            // IN: ephemeral public party B key
    )
{
    pAssert(outZ1 != NULL
            && dsA != NULL && deA != NULL
            && QsB != NULL && QeB != NULL);

    // Initialize the output points so that they are empty until one of the
    // functions decides otherwise
    outZ1->x.b.size = 0;
    outZ1->y.b.size = 0;
    if(outZ2 != NULL)
    {
        outZ2->x.b.size = 0;
        outZ2->y.b.size = 0;
    }
    switch(scheme)
    {
        case ALG_ECDH_VALUE:
            return C_2_2_ECDH(outZ1, outZ2, curveId, dsA, deA, QsB, QeB);
            break;
#if ALG_ECMQV
        case ALG_ECMQV_VALUE:
            return C_2_2_MQV(outZ1, curveId, dsA, deA, QsB, QeB);
            break;
#endif
#if ALG_SM2
        case ALG_SM2_VALUE:
            return SM2KeyExchange(outZ1, curveId, dsA, deA, QsB, QeB);
            break;
#endif
        default:
            return TPM_RC_SCHEME;
    }
}

#if ALG_SM2

//*** ComputeWForSM2()
// Compute the value for w used by SM2
static UINT32
ComputeWForSM2(
    bigCurve        E
    )
{
    //  w := ceil(ceil(log2(n)) / 2) - 1
    return (BnMsb(CurveGetOrder(AccessCurveData(E))) / 2 - 1);
}

//*** avfSm2()
// This function does the associated value computation required by SM2 key
// exchange. This is different from the avf() in the international standards
// because it returns a value that is half the size of the value returned by the
// standard avf(). For example, if 'n' is 15, 'Ws' ('w' in the standard) is 2 but 
// the 'W' here is 1. This means that an input value of 14 (1110b) would return a 
// value of 110b with the standard but 10b with the scheme in SM2.
static bigNum
avfSm2(
    bigNum              bn,           // IN/OUT: the reduced value
    UINT32              w              // IN: the value of w
    )
{
    // a)   set w := ceil(ceil(log2(n)) / 2) - 1
    // b)   set x' := 2^w + ( x & (2^w - 1))
    // This is just like the avf for MQV where x' = 2^w + (x mod 2^w)

    BnMaskBits(bn, w);   // as with avf1, this is too big by a factor of 2 but
                         // it doesn't matter because we SET the extra bit
                         // anyway
    BnSetBit(bn, w);
    return bn;
}

//*** SM2KeyExchange()
// This function performs the key exchange defined in SM2.
// The first step is to compute
//  'tA' = ('dsA' + 'deA'  avf(Xe,A)) mod 'n'
// Then, compute the 'Z' value from
// 'outZ' = ('h'  'tA' mod 'n') ('QsA' + [avf('QeB.x')]('QeB')).
// The function will compute the ephemeral public key from the ephemeral
// private key.
// All points are required to be on the curve of 'inQsA'. The function will fail
// catastrophically if this is not the case
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        the value for dsA does not give a valid point on the
//                              curve
LIB_EXPORT TPM_RC
SM2KeyExchange(
    TPMS_ECC_POINT        *outZ,         // OUT: the computed point
    TPM_ECC_CURVE          curveId,      // IN: the curve for the computations
    TPM2B_ECC_PARAMETER   *dsAIn,        // IN: static private TPM key
    TPM2B_ECC_PARAMETER   *deAIn,        // IN: ephemeral private TPM key
    TPMS_ECC_POINT        *QsBIn,        // IN: static public party B key
    TPMS_ECC_POINT        *QeBIn         // IN: ephemeral public party B key
    )
{
    CURVE_INITIALIZED(E, curveId);
    const ECC_CURVE_DATA      *C;
    ECC_INITIALIZED(dsA, dsAIn);
    ECC_INITIALIZED(deA, deAIn);
    POINT_INITIALIZED(QsB, QsBIn);
    POINT_INITIALIZED(QeB, QeBIn);
    BN_WORD_INITIALIZED(One, 1);
    POINT(QeA);
    ECC_NUM(XeB);
    POINT(Z);
    ECC_NUM(Ta);
    UINT32                   w;
    TPM_RC                 retVal = TPM_RC_NO_RESULT;
//
    // Parameter checks
    if(E == NULL)
        ERROR_RETURN(TPM_RC_CURVE);
    C = AccessCurveData(E);
    pAssert(outZ != NULL && dsA != NULL && deA != NULL &&  QsB != NULL 
            && QeB != NULL);

    // Compute the value for w
    w = ComputeWForSM2(E);

    // Compute the public ephemeral key pQeA = [de,A]G
    if(!BnEccModMult(QeA, CurveGetG(C), deA, E))
        goto Exit;

    //  tA := (ds,A + de,A  avf(Xe,A)) mod n    (3)
    //  Compute 'tA' = ('dsA' +  'deA'  avf('XeA')) mod n
    // Ta = avf(XeA);
    // do Ta = de,A * Ta = deA * avf(XeA)
    BnMult(Ta, deA, avfSm2(QeA->x, w));
    // now Ta = dsA + Ta =  dsA + deA * avf(XeA)
    BnAdd(Ta, dsA, Ta);
    BnMod(Ta, CurveGetOrder(C));

    //  outZ = [h  tA mod n] (Qs,B + [avf(Xe,B)](Qe,B)) (4)
    // Put this in because almost every case of h is == 1 so skip the call when
    // not necessary.
    if(!BnEqualWord(CurveGetCofactor(C), 1))
        // Cofactor is not 1 so compute Ta := Ta * h mod n
        BnModMult(Ta, Ta, CurveGetCofactor(C), CurveGetOrder(C));
    // Now that 'tA' is (h * 'tA' mod n)
    // 'outZ' = ['tA'](QsB + [avf(QeB.x)](QeB)).
    BnCopy(XeB, QeB->x);
    if(!BnEccModMult2(Z, QsB, One, QeB, avfSm2(XeB, w), E))
        goto Exit;
    // QeB := [tA]QeB = [tA](QsB + [Xe,B]QeB) and check for at infinity
    if(!BnEccModMult(Z, Z, Ta, E))
        goto Exit;
    // Convert BIGNUM E to TPM2B E
    BnPointTo2B(outZ, Z, E);
    retVal = TPM_RC_SUCCESS;
Exit:
    CURVE_FREE(E);
    return retVal;
}
#endif

#endif // CC_ZGen_2Phase