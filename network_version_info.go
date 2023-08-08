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

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// NetworkVersionInfo is the version info for the Hedera network protobuf and services
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
