dnl The copyright in this software is being made available under the BSD License,
dnl included below. This software may be subject to other third party and
dnl contributor rights, including patent rights, and no such rights are granted
dnl under this license.
dnl
dnl Copyright (c) Intel Corporation
dnl
dnl All rights reserved.
dnl
dnl BSD License
dnl
dnl Redistribution and use in source and binary forms, with or without modification,
dnl are permitted provided that the following conditions are met:
dnl
dnl Redistributions of source code must retain the above copyright notice, this list
dnl of conditions and the following disclaimer.
dnl
dnl Redistributions in binary form must reproduce the above copyright notice, this
dnl list of conditions and the following disclaimer in the documentation and/or
dnl other materials provided with the distribution.
dnl
dnl THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS ""AS IS""
dnl AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
dnl IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
dnl DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
dnl ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
dnl (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
dnl LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
dnl ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
dnl (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
dnl SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

AC_INIT([ms-tpm-20-ref],
        [0.1],
        [https://github.com/microsoft/ms-tpm-20-ref/issues],
        [],
        [https://github.com/microsoft/ms-tpm-20-ref])
AC_CONFIG_MACRO_DIR([.])
AC_PROG_CC
AC_PROG_LN_S
AC_PROG_RANLIB
AM_INIT_AUTOMAKE([foreign subdir-objects])
AC_CONFIG_FILES([Makefile])
AC_SUBST([DISTCHECK_CONFIGURE_FLAGS],[$ac_configure_args])

dnl By enabling this feature tpm simulator gets seeds derived from hardware parameters.
dnl It is enabled only for linux devices.
dnl Note that the seeds are not derived from secure hardware source.

AC_ARG_ENABLE(usedeviceid,
    AS_HELP_STRING([--enable-usedeviceid],
    [tpm simulator get seeds derived from hardware parameters. Seeds are not derived from secure hardware source.]))

PKG_CHECK_MODULES([LIBCRYPTO], [libcrypto])
AS_IF([test "x$enable_usedeviceid" = "xyes"], [
    PKG_CHECK_MODULES([LIBUDEV], [libudev])
    [ADDITIONAL_LIBS="-ludev"]
])
AX_PTHREAD([], [AC_MSG_ERROR([requires pthread])])

AC_DEFINE([HASH_LIB], [Ossl], [Crypto lib for hash algorithms])
AC_DEFINE([SYM_LIB], [Ossl], [Crypto lib for symmetric encryption algorithms])
AC_DEFINE([MATH_LIB], [Ossl], [Crypto lib for bignum operations])

ADD_COMPILER_FLAG([-std=gnu11])
ADD_COMPILER_FLAG([-Werror])
ADD_COMPILER_FLAG([-Wall])
ADD_COMPILER_FLAG([-Wformat-security])
ADD_COMPILER_FLAG([-fstack-protector-all])
ADD_COMPILER_FLAG([-fPIC])
ADD_COMPILER_FLAG([-Wno-error=empty-body])
ADD_COMPILER_FLAG([-Wno-error=expansion-to-defined])
ADD_COMPILER_FLAG([-Wno-error=parentheses])
ADD_COMPILER_FLAG([-Wno-error=pointer-to-int-cast])
ADD_COMPILER_FLAG([-Wno-error=missing-braces])
ADD_COMPILER_FLAG([-Wno-error=unused-result])

AS_IF([test "x$enable_usedeviceid" = "xyes"], [
    ADD_COMPILER_FLAG([-DNDEBUG])
    ADD_COMPILER_FLAG([-g])
    ADD_COMPILER_FLAG([-DUSE_PLATFORM_EPS])
    AC_SUBST(ADDITIONAL_LIBS)
])
ADD_LINK_FLAG([-Wl,--no-undefined])
ADD_LINK_FLAG([-Wl,-z,noexecstack])
ADD_LINK_FLAG([-Wl,-z,now])
ADD_LINK_FLAG([-Wl,-z,relro])

AC_OUTPUT
