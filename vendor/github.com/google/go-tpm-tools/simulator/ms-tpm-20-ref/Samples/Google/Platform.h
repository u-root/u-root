/* Microsoft Reference Implementation for TPM 2.0
 *
 *  The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 *  Copyright (c) Microsoft Corporation
 *
 *  All rights reserved.
 *
 *  BSD License
 *
 *  Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 *  Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS ""AS
 * IS"" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
 * THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
 * PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
 * OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
 * OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
 * ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
// External interface to the vTPM

#ifndef _PLATFORM_H_
#define _PLATFORM_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>

//***_plat__RunCommand()
// This version of RunCommand will set up a jum_buf and call ExecuteCommand().
// If the command executes without failing, it will return and RunCommand will
// return. If there is a failure in the command, then _plat__Fail() is called
// and it will longjump back to RunCommand which will call ExecuteCommand again.
// However, this time, the TPM will be in failure mode so ExecuteCommand will
// simply build a failure response and return.
void _plat__RunCommand(uint32_t requestSize,     // IN: command buffer size
                       unsigned char *request,   // IN: command buffer
                       uint32_t *responseSize,   // IN/OUT: response buffer size
                       unsigned char **response  // IN/OUT: response buffer
);

//*** _plat_Reset()
// Reset the TPM. This should always be called before _plat__RunCommand. The
// first time this function is called, the TPM will be manufactured. Pass true
// for forceManufacture to perfrom a manufacturer reset.
void _plat__Reset(bool forceManufacture);

#ifdef __cplusplus
}
#endif

#endif  // _PLATFORM_H_
