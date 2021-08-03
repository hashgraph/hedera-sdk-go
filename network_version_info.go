package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

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
	if version == nil {
		return NetworkVersionInfo{}
	}
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

func (version *NetworkVersionInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(version.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func NetworkVersionInfoFromBytes(data []byte) (NetworkVersionInfo, error) {
	if data == nil {
		return NetworkVersionInfo{}, errByteArrayNull
	}
	pb := proto.NetworkGetVersionInfoResponse{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return NetworkVersionInfo{}, err
	}

	info := networkVersionInfoFromProtobuf(&pb)

	return info, nil
}
