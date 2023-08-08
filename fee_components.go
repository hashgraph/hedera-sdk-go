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

// nolint
type FeeComponents struct {
	Min                        int64
	Max                        int64
	Constant                   int64
	TransactionBandwidthByte   int64
	TransactionVerification    int64
	TransactionRamByteHour     int64
	TransactionStorageByteHour int64
	ContractTransactionGas     int64
	TransferVolumeHbar         int64
	ResponseMemoryByte         int64
	ResponseDiscByte           int64
}

func _FeeComponentsFromProtobuf(feeComponents *services.FeeComponents) (FeeComponents, error) {
	if feeComponents == nil {
		return FeeComponents{}, errParameterNull
	}

	return FeeComponents{
		Min:                        feeComponents.GetMin(),
		Max:                        feeComponents.GetMax(),
		Constant:                   feeComponents.GetConstant(),
		TransactionBandwidthByte:   feeComponents.GetBpt(),
		TransactionVerification:    feeComponents.GetVpt(),
		TransactionRamByteHour:     feeComponents.GetRbh(),
		TransactionStorageByteHour: feeComponents.GetSbh(),
		ContractTransactionGas:     feeComponents.GetGas(),
		TransferVolumeHbar:         feeComponents.GetTv(),
		ResponseMemoryByte:         feeComponents.GetBpr(),
		ResponseDiscByte:           feeComponents.GetSbpr(),
	}, nil
}

func (feeComponents FeeComponents) _ToProtobuf() *services.FeeComponents {
	return &services.FeeComponents{
		Min:      feeComponents.Min,
		Max:      feeComponents.Max,
		Constant: feeComponents.Constant,
		Bpt:      feeComponents.TransactionBandwidthByte,
		Vpt:      feeComponents.TransactionVerification,
		Rbh:      feeComponents.TransactionRamByteHour,
		Sbh:      feeComponents.TransactionStorageByteHour,
		Gas:      feeComponents.ContractTransactionGas,
		Tv:       feeComponents.TransferVolumeHbar,
		Bpr:      feeComponents.ResponseMemoryByte,
		Sbpr:     feeComponents.ResponseDiscByte,
	}
}

// ToBytes returns the byte representation of the FeeComponents
func (feeComponents FeeComponents) ToBytes() []byte {
	data, err := protobuf.Marshal(feeComponents._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// FeeComponentsFromBytes returns the FeeComponents from a byte array representation
func FeeComponentsFromBytes(data []byte) (FeeComponents, error) {
	if data == nil {
		return FeeComponents{}, errByteArrayNull
	}
	pb := services.FeeComponents{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeComponents{}, err
	}

	info, err := _FeeComponentsFromProtobuf(&pb)
	if err != nil {
		return FeeComponents{}, err
	}

	return info, nil
}

// String returns a string representation of the FeeComponents
func (feeComponents FeeComponents) String() string {
	return fmt.Sprintf("Min: %d, Max: %d, Constant: %d,TransactionBandwithByte: %d,TransactionVerification: %d,TransactionRamByteHour: %d,TransactionStorageByteHour: %d, ContractTransactionGas: %d,TransferVolumeHbar: %d, ResponseMemoryByte: %d, ResponseDiscByte: %d", feeComponents.Min, feeComponents.Max, feeComponents.Constant, feeComponents.TransactionBandwidthByte, feeComponents.TransactionVerification, feeComponents.TransactionRamByteHour, feeComponents.TransactionStorageByteHour, feeComponents.ContractTransactionGas, feeComponents.TransferVolumeHbar, feeComponents.ResponseMemoryByte, feeComponents.ResponseDiscByte)
}
