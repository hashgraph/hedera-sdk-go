package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// NetworkVersionInfo is the version info for the Hiero network protobuf and services
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

// ToBytes returns the byte representation of the NetworkVersionInfo
func (version *NetworkVersionInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(version._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// NetworkVersionInfoFromBytes returns the NetworkVersionInfo from a raw byte array
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
