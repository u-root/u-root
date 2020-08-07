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
    
#ifndef     _VENDOR_STRING_H
#define     _VENDOR_STRING_H

// Define up to 4-byte values for MANUFACTURER.  This value defines the response  
// for TPM_PT_MANUFACTURER in TPM2_GetCapability.
// The following line should be un-commented and a vendor specific string 
// should be provided here.  
#define    MANUFACTURER    "MSFT"

// The following #if macro may be deleted after a proper MANUFACTURER is provided.
#ifndef MANUFACTURER
#error MANUFACTURER is not provided. \
Please modify include/VendorString.h to provide a specific \
manufacturer name.
#endif

// Define up to 4, 4-byte values. The values must each be 4 bytes long and the last
// value used may contain trailing zeros.
// These values define the response for TPM_PT_VENDOR_STRING_(1-4)
// in TPM2_GetCapability.
// The following line should be un-commented and a vendor specific string 
// should be provided here.  
// The vendor strings 2-4 may also be defined as appropriate.
#define       VENDOR_STRING_1       "xCG "
#define       VENDOR_STRING_2       "fTPM"
// #define       VENDOR_STRING_3 
// #define       VENDOR_STRING_4

// The following #if macro may be deleted after a proper VENDOR_STRING_1 
// is provided.
#ifndef VENDOR_STRING_1
#error VENDOR_STRING_1 is not provided. \
Please modify include/VendorString.h to provide a vendor-specific string.
#endif

// the more significant 32-bits of a vendor-specific value 
// indicating the version of the firmware
// The following line should be un-commented and a vendor specific firmware V1 
// should be provided here. 
// The FIRMWARE_V2 may also be defined as appropriate.
#define   FIRMWARE_V1         (0x20170619)
// the less significant 32-bits of a vendor-specific value 
// indicating the version of the firmware
#define   FIRMWARE_V2         (0x00163636)

// The following #if macro may be deleted after a proper FIRMWARE_V1 is provided.
#ifndef FIRMWARE_V1
#error  FIRMWARE_V1 is not provided. \
Please modify include/VendorString.h to provide a vendor-specific firmware \
version
#endif

#endif
