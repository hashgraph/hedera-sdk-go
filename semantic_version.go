package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func semanticVersionFromProtobuf(version *proto.SemanticVersion) SemanticVersion {
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

func (version *SemanticVersion) toProtobuf() *proto.SemanticVersion {
	return &proto.SemanticVersion{
		Major: int32(version.Major),
		Minor: int32(version.Minor),
		Patch: int32(version.Patch),
		Pre:   version.Pre,
		Build: version.Build,
	}
}
