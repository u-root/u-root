/*
 * Copyright (c) 2020, NVIDIA CORPORATION.  All rights reserved.
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

package dcgm

type MigProfile int

const (
	MigProfileNone                      MigProfile = 0  /*!< No profile (for GPUs) */
	MigProfileGPUInstanceSlice1         MigProfile = 1  /*!< GPU instance slice 1 */
	MigProfileGPUInstanceSlice2         MigProfile = 2  /*!< GPU instance slice 2 */
	MigProfileGPUInstanceSlice3         MigProfile = 3  /*!< GPU instance slice 3 */
	MigProfileGPUInstanceSlice4         MigProfile = 4  /*!< GPU instance slice 4 */
	MigProfileGPUInstanceSlice7         MigProfile = 5  /*!< GPU instance slice 7 */
	MigProfileGPUInstanceSlice8         MigProfile = 6  /*!< GPU instance slice 8 */
	MigProfileGPUInstanceSlice6         MigProfile = 7  /*!< GPU instance slice 6 */
	MigProfileGPUInstanceSlice1Rev1     MigProfile = 8  /*!< GPU instance slice 1 revision 1 */
	MigProfileGPUInstanceSlice2Rev1     MigProfile = 9  /*!< GPU instance slice 2 revision 1 */
	MigProfileGPUInstanceSlice1Rev2     MigProfile = 10 /*!< GPU instance slice 1 revision 2 */
	MigProfileComputeInstanceSlice1     MigProfile = 30 /*!< compute instance slice 1 */
	MigProfileComputeInstanceSlice2     MigProfile = 31 /*!< compute instance slice 2 */
	MigProfileComputeInstanceSlice3     MigProfile = 32 /*!< compute instance slice 3 */
	MigProfileComputeInstanceSlice4     MigProfile = 33 /*!< compute instance slice 4*/
	MigProfileComputeInstanceSlice7     MigProfile = 34 /*!< compute instance slice 7 */
	MigProfileComputeInstanceSlice8     MigProfile = 35 /*!< compute instance slice 8 */
	MigProfileComputeInstanceSlice6     MigProfile = 36 /*!< compute instance slice 6 */
	MigProfileComputeInstanceSlice1Rev1 MigProfile = 37 /*!< compute instance slice 1 revision 1 */
)
