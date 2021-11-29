package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type NetworkVersionInfo struct {
	ProtobufVersion SemanticVersion
	ServicesVersion SemanticVersion
}

func _NetworkVersionInfoFromProtobuf(version *services.NetworkGetVersionInfoResponse) NetworkVersionInfo {
	if version == nil {
		return NetworkVersionInfo{}
	}
	return NetworkVersionInfo{
		ProtobufVersion: _SemanticVersionFromProtobuf(version.HapiProtoVersion),
		ServicesVersion: _SemanticVersionFromProtobuf(version.HederaServicesVersion),
	}
}

func (version *NetworkVersionInfo) _ToProtobuf() *services.NetworkGetVersionInfoResponse {
	return &services.NetworkGetVersionInfoResponse{
		HapiProtoVersion:      version.ProtobufVersion._ToProtobuf(),
		HederaServicesVersion: version.ServicesVersion._ToProtobuf(),
	}
}

func (version *NetworkVersionInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(version._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func NetworkVersionInfoFromBytes(data []byte) (NetworkVersionInfo, error) {
	if data == nil {
		return NetworkVersionInfo{}, errByteArrayNull
	}
	pb := services.NetworkGetVersionInfoResponse{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return NetworkVersionInfo{}, err
	}

	info := _NetworkVersionInfoFromProtobuf(&pb)

	return info, nil
}
