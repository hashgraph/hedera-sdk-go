package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import "fmt"

type FreezeType int32

const (
	FreezeTypeUnknown          FreezeType = 0
	FreezeTypeFreezeOnly       FreezeType = 1
	FreezeTypePrepareUpgrade   FreezeType = 2
	FreezeTypeFreezeUpgrade    FreezeType = 3
	FreezeTypeFreezeAbort      FreezeType = 4
	FreezeTypeTelemetryUpgrade FreezeType = 5
)

func (freezeType FreezeType) String() string {
	switch freezeType {
	case FreezeTypeUnknown:
		return "UNKNOWN_FREEZE_TYPE"
	case FreezeTypeFreezeOnly:
		return "FREEZE_ONLY"
	case FreezeTypePrepareUpgrade:
		return "PREPARE_UPGRADE"
	case FreezeTypeFreezeUpgrade:
		return "FREEZE_UPGRADE"
	case FreezeTypeFreezeAbort:
		return "FREEZE_ABORT"
	case FreezeTypeTelemetryUpgrade:
		return "TELEMETRY_UPGRADE"
	}

	panic(fmt.Sprintf("unreacahble: FreezeType.String() switch statement is non-exhaustive. Status: %v", uint32(freezeType)))
}
