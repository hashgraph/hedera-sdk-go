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
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type FeeData struct {
	NodeData    *FeeComponents
	NetworkData *FeeComponents
	ServiceData *FeeComponents
}

func _FeeDataFromProtobuf(feeData *services.FeeData) (FeeData, error) {
	if feeData == nil {
		return FeeData{}, errParameterNull
	}

	nodeData, err := _FeeComponentsFromProtobuf(feeData.Nodedata)
	if err != nil {
		return FeeData{}, err
	}

	networkData, err := _FeeComponentsFromProtobuf(feeData.Networkdata)
	if err != nil {
		return FeeData{}, err
	}

	serviceData, err := _FeeComponentsFromProtobuf(feeData.Servicedata)
	if err != nil {
		return FeeData{}, err
	}

	return FeeData{
		NodeData:    &nodeData,
		NetworkData: &networkData,
		ServiceData: &serviceData,
	}, nil
}

func (feeData FeeData) _ToProtobuf() *services.FeeData {
	var nodeData *services.FeeComponents
	if feeData.NodeData != nil {
		nodeData = feeData.NodeData._ToProtobuf()
	}

	var networkData *services.FeeComponents
	if feeData.NetworkData != nil {
		networkData = feeData.NetworkData._ToProtobuf()
	}

	var serviceData *services.FeeComponents
	if feeData.ServiceData != nil {
		serviceData = feeData.ServiceData._ToProtobuf()
	}

	return &services.FeeData{
		Nodedata:    nodeData,
		Networkdata: networkData,
		Servicedata: serviceData,
	}
}

// ToBytes returns the byte representation of the FeeData
func (feeData FeeData) ToBytes() []byte {
	data, err := protobuf.Marshal(feeData._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// FeeDataFromBytes returns a FeeData struct from a raw byte array
func FeeDataFromBytes(data []byte) (FeeData, error) {
	if data == nil {
		return FeeData{}, errByteArrayNull
	}
	pb := services.FeeData{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeData{}, err
	}

	info, err := _FeeDataFromProtobuf(&pb)
	if err != nil {
		return FeeData{}, err
	}

	return info, nil
}

// String returns a string representation of the FeeData
func (feeData FeeData) String() string {
	return fmt.Sprintf("\nNodedata: %s\nNetworkdata: %s\nServicedata: %s\n", feeData.NodeData.String(), feeData.NetworkData.String(), feeData.ServiceData.String())
}
