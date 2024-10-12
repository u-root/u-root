/*
 * Copyright (c) 2023, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#ifndef DCGM_DCGM_API_EXPORT_H
#define DCGM_DCGM_API_EXPORT_H

#undef DCGM_PUBLIC_API
#undef DCGM_PRIVATE_API

#if defined(DCGM_API_EXPORT)
#define DCGM_PUBLIC_API __attribute((visibility("default")))
#else
#define DCGM_PUBLIC_API
#if defined(ERROR_IF_NOT_PUBLIC)
#error(Should be public)
#endif
#endif

#define DCGM_PRIVATE_API __attribute((visibility("hidden")))


#endif // DCGM_DCGM_API_EXPORT_H
