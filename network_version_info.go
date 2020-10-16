package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type NetworkVersionInfo struct {
	ProtobufVersion SemanticVersion
	ServicesVersion SemanticVersion
}

func newNetworkVersionInfo(hapi SemanticVersion, hedera SemanticVersion) NetworkVersionInfo {
	return NetworkVersionInfo{
		ProtobufVersion: hapi,
		ServicesVersion: hedera,
	}
}

func networkVersionInfoFromProtobuf(version *proto.NetworkGetVersionInfoResponse) NetworkVersionInfo {
	return NetworkVersionInfo{
		ProtobufVersion: semanticVersionFromProtobuf(version.HapiProtoVersion),
		ServicesVersion: semanticVersionFromProtobuf(version.HederaServicesVersion),
	}
}

func (version *NetworkVersionInfo) toProtobuf() *proto.NetworkGetVersionInfoResponse {
	return &proto.NetworkGetVersionInfoResponse{
		HapiProtoVersion:      version.ProtobufVersion.toProtobuf(),
		HederaServicesVersion: version.ServicesVersion.toProtobuf(),
	}
}
