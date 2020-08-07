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

dnl ADD_COMPILER_FLAG:
dnl   A macro to add a CFLAG to the EXTRA_CFLAGS variable. This macro will
dnl   check to be sure the compiler supprts the flag. Flags can be made
dnl   mandatory (configure will fail).
dnl $1: C compiler flag to add to EXTRA_CFLAGS.
dnl $2: Set to "required" to cause configure failure if flag not supported..
AC_DEFUN([ADD_COMPILER_FLAG],[
    AX_CHECK_COMPILE_FLAG([$1],[
        EXTRA_CFLAGS="$EXTRA_CFLAGS $1"
        AC_SUBST([EXTRA_CFLAGS])],[
        AS_IF([test x$2 != xrequired],[
            AC_MSG_WARN([Optional CFLAG "$1" not supported by your compiler, continuing.])],[
            AC_MSG_ERROR([Required CFLAG "$1" not supported by your compiler, aborting.])]
        )],[
        -Wall -Werror]
    )]
)
dnl ADD_PREPROC_FLAG:
dnl   Add the provided preprocessor flag to the EXTRA_CFLAGS variable. This
dnl   macro will check to be sure the preprocessor supports the flag.
dnl   The flag can be made mandatory by provideing the string 'required' as
dnl   the second parameter.
dnl $1: Preprocessor flag to add to EXTRA_CFLAGS.
dnl $2: Set to "required" t ocause configure failure if preprocesor flag
dnl     is not supported.
AC_DEFUN([ADD_PREPROC_FLAG],[
    AX_CHECK_PREPROC_FLAG([$1],[
        EXTRA_CFLAGS="$EXTRA_CFLAGS $1"
        AC_SUBST([EXTRA_CFLAGS])],[
        AS_IF([test x$2 != xrequired],[
            AC_MSG_WARN([Optional preprocessor flag "$1" not supported by your compiler, continuing.])],[
            AC_MSG_ERROR([Required preprocessor flag "$1" not supported by your compiler, aborting.])]
        )],[
        -Wall -Werror]
    )]
)
dnl ADD_LINK_FLAG:
dnl   A macro to add a LDLAG to the EXTRA_LDFLAGS variable. This macro will
dnl   check to be sure the linker supprts the flag. Flags can be made
dnl   mandatory (configure will fail).
dnl $1: linker flag to add to EXTRA_LDFLAGS.
dnl $2: Set to "required" to cause configure failure if flag not supported.
AC_DEFUN([ADD_LINK_FLAG],[
    AX_CHECK_LINK_FLAG([$1],[
        EXTRA_LDFLAGS="$EXTRA_LDFLAGS $1"
        AC_SUBST([EXTRA_LDFLAGS])],[
        AS_IF([test x$2 != xrequired],[
            AC_MSG_WARN([Optional LDFLAG "$1" not supported by your linker, continuing.])],[
            AC_MSG_ERROR([Required LDFLAG "$1" not supported by your linker, aborting.])]
        )]
    )]
)
