package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type SemanticVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Pre   string
	Build string
}

func _SemanticVersionFromProtobuf(version *services.SemanticVersion) SemanticVersion {
	if version == nil {
		return SemanticVersion{}
	}
	return SemanticVersion{
		Major: uint32(version.GetMajor()),
		Minor: uint32(version.GetMinor()),
		Patch: uint32(version.GetPatch()),
		Pre:   version.GetPre(),
		Build: version.GetBuild(),
	}
}

func (version *SemanticVersion) _ToProtobuf() *services.SemanticVersion {
	return &services.SemanticVersion{
		Major: int32(version.Major),
		Minor: int32(version.Minor),
		Patch: int32(version.Patch),
		Pre:   version.Pre,
		Build: version.Build,
	}
}
