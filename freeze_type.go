package hiero

// SPDX-License-Identifier: Apache-2.0

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
