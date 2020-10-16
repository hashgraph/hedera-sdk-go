package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type SemanticVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
}

func newSemanticVersion(major uint32, minor uint32, patch uint32) SemanticVersion {
	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func semanticVersionFromProtobuf(version *proto.SemanticVersion) SemanticVersion {
	return SemanticVersion{
		Major: uint32(version.GetMajor()),
		Minor: uint32(version.GetMinor()),
		Patch: uint32(version.GetPatch()),
	}
}

func (version *SemanticVersion) toProtobuf() *proto.SemanticVersion {
	return &proto.SemanticVersion{
		Major: int32(version.Major),
		Minor: int32(version.Minor),
		Patch: int32(version.Patch),
	}
}
