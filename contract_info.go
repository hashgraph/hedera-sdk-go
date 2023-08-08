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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// Current information on the smart contract instance, including its balance.
type ContractInfo struct {
	AccountID                     AccountID
	ContractID                    ContractID
	ContractAccountID             string
	AdminKey                      Key
	ExpirationTime                time.Time
	AutoRenewPeriod               time.Duration
	Storage                       uint64
	ContractMemo                  string
	Balance                       uint64
	LedgerID                      LedgerID
	AutoRenewAccountID            *AccountID
	MaxAutomaticTokenAssociations int32
	StakingInfo                   *StakingInfo
}

func _ContractInfoFromProtobuf(contractInfo *services.ContractGetInfoResponse_ContractInfo) (ContractInfo, error) {
	if contractInfo == nil {
		return ContractInfo{}, errParameterNull
	}

	var adminKey Key
	var err error
	if contractInfo.GetAdminKey() != nil {
		adminKey, err = _KeyFromProtobuf(contractInfo.GetAdminKey())
		if err != nil {
			return ContractInfo{}, err
		}
	}

	accountID := AccountID{}
	if contractInfo.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(contractInfo.AccountID)
	}

	contractID := ContractID{}
	if contractInfo.ContractID != nil {
		contractID = *_ContractIDFromProtobuf(contractInfo.ContractID)
	}

	var autoRenewAccountID *AccountID
	if contractInfo.AutoRenewAccountId != nil {
		autoRenewAccountID = _AccountIDFromProtobuf(contractInfo.AutoRenewAccountId)
	}

	var stakingInfo StakingInfo
	if contractInfo.StakingInfo != nil {
		stakingInfo = _StakingInfoFromProtobuf(contractInfo.StakingInfo)
	}

	return ContractInfo{
		AccountID:                     accountID,
		ContractID:                    contractID,
		ContractAccountID:             contractInfo.ContractAccountID,
		AdminKey:                      adminKey,
		ExpirationTime:                _TimeFromProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:               _DurationFromProtobuf(contractInfo.AutoRenewPeriod),
		Storage:                       uint64(contractInfo.Storage),
		ContractMemo:                  contractInfo.Memo,
		Balance:                       contractInfo.Balance,
		LedgerID:                      LedgerID{contractInfo.LedgerId},
		AutoRenewAccountID:            autoRenewAccountID,
		MaxAutomaticTokenAssociations: contractInfo.MaxAutomaticTokenAssociations,
		StakingInfo:                   &stakingInfo,
	}, nil
}

func (contractInfo *ContractInfo) _ToProtobuf() *services.ContractGetInfoResponse_ContractInfo {
	body := &services.ContractGetInfoResponse_ContractInfo{
		ContractID:        contractInfo.ContractID._ToProtobuf(),
		AccountID:         contractInfo.AccountID._ToProtobuf(),
		ContractAccountID: contractInfo.ContractAccountID,
		AdminKey:          contractInfo.AdminKey._ToProtoKey(),
		ExpirationTime:    _TimeToProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:   _DurationToProtobuf(contractInfo.AutoRenewPeriod),
		Storage:           int64(contractInfo.Storage),
		Memo:              contractInfo.ContractMemo,
		Balance:           contractInfo.Balance,
		LedgerId:          contractInfo.LedgerID.ToBytes(),
	}

	if contractInfo.AutoRenewAccountID != nil {
		body.AutoRenewAccountId = contractInfo.AutoRenewAccountID._ToProtobuf()
	}

	if contractInfo.StakingInfo != nil {
		body.StakingInfo = contractInfo.StakingInfo._ToProtobuf()
	}

	return body
}

// ToBytes returns a serialized version of the ContractInfo object
func (contractInfo ContractInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(contractInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// ContractInfoFromBytes returns a ContractInfo object deserialized from bytes
func ContractInfoFromBytes(data []byte) (ContractInfo, error) {
	if data == nil {
		return ContractInfo{}, errByteArrayNull
	}
	pb := services.ContractGetInfoResponse_ContractInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractInfo{}, err
	}

	info, err := _ContractInfoFromProtobuf(&pb)
	if err != nil {
		return ContractInfo{}, err
	}

	return info, nil
}
