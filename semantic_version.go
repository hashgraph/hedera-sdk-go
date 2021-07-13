package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type SemanticVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Pre   string
	Build string
}

func newSemanticVersion(major uint32, minor uint32, patch uint32, pre string, build string) SemanticVersion {
	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
		Pre:   pre,
		Build: build,
	}
}

func semanticVersionFromProtobuf(version *services.SemanticVersion) SemanticVersion {
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

func (version *SemanticVersion) toProtobuf() *services.SemanticVersion {
	return &services.SemanticVersion{
		Major: int32(version.Major),
		Minor: int32(version.Minor),
		Patch: int32(version.Patch),
		Pre:   version.Pre,
		Build: version.Build,
	}
}
