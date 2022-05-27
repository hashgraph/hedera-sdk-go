package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

// AccountInfo is info about the account returned from an AccountInfoQuery
type AccountInfo struct {
	AccountID         AccountID
	ContractAccountID string
	IsDeleted         bool
	// Deprecated
	ProxyAccountID                 AccountID
	ProxyReceived                  Hbar
	Key                            Key
	Balance                        Hbar
	GenerateSendRecordThreshold    Hbar
	GenerateReceiveRecordThreshold Hbar
	ReceiverSigRequired            bool
	ExpirationTime                 time.Time
	AutoRenewPeriod                time.Duration
	LiveHashes                     []*LiveHash
	TokenRelationships             []*TokenRelationship
	AccountMemo                    string
	OwnedNfts                      int64
	MaxAutomaticTokenAssociations  uint32
	AliasKey                       *PublicKey
	LedgerID                       LedgerID
	// Deprecated
	HbarAllowances []HbarAllowance
	// Deprecated
	NftAllowances []TokenNftAllowance
	// Deprecated
	TokenAllowances []TokenAllowance
	EthereumNonce   int64
	StakingInfo     *StakingInfo
}

func _AccountInfoFromProtobuf(pb *services.CryptoGetInfoResponse_AccountInfo) (AccountInfo, error) {
	if pb == nil {
		return AccountInfo{}, errParameterNull
	}

	pubKey, err := _KeyFromProtobuf(pb.Key)
	if err != nil {
		return AccountInfo{}, err
	}

	tokenRelationship := make([]*TokenRelationship, len(pb.TokenRelationships))

	if pb.TokenRelationships != nil {
		for i, relationship := range pb.TokenRelationships {
			singleRelationship := _TokenRelationshipFromProtobuf(relationship)
			tokenRelationship[i] = &singleRelationship
		}
	}

	liveHashes := make([]*LiveHash, len(pb.LiveHashes))

	if pb.LiveHashes != nil {
		for i, liveHash := range pb.LiveHashes {
			singleRelationship, err := _LiveHashFromProtobuf(liveHash)
			if err != nil {
				return AccountInfo{}, err
			}
			liveHashes[i] = &singleRelationship
		}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	var alias *PublicKey
	if len(pb.Alias) != 0 {
		pbKey := services.Key{}
		_ = protobuf.Unmarshal(pb.Alias, &pbKey)
		initialKey, _ := _KeyFromProtobuf(&pbKey)
		switch t2 := initialKey.(type) { //nolint
		case PublicKey:
			alias = &t2
		}
	}

	var stakingInfo StakingInfo
	if pb.StakingInfo != nil {
		stakingInfo = _StakingInfoFromProtobuf(pb.StakingInfo)
	}

	return AccountInfo{
		AccountID:                      accountID,
		ContractAccountID:              pb.ContractAccountID,
		IsDeleted:                      pb.Deleted,
		ProxyReceived:                  HbarFromTinybar(pb.ProxyReceived),
		Key:                            pubKey,
		Balance:                        HbarFromTinybar(int64(pb.Balance)),
		GenerateSendRecordThreshold:    HbarFromTinybar(int64(pb.GenerateSendRecordThreshold)),    // nolint
		GenerateReceiveRecordThreshold: HbarFromTinybar(int64(pb.GenerateReceiveRecordThreshold)), // nolint
		ReceiverSigRequired:            pb.ReceiverSigRequired,
		TokenRelationships:             tokenRelationship,
		ExpirationTime:                 _TimeFromProtobuf(pb.ExpirationTime),
		AccountMemo:                    pb.Memo,
		AutoRenewPeriod:                _DurationFromProtobuf(pb.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		OwnedNfts:                      pb.OwnedNfts,
		MaxAutomaticTokenAssociations:  uint32(pb.MaxAutomaticTokenAssociations),
		AliasKey:                       alias,
		LedgerID:                       LedgerID{pb.LedgerId},
		EthereumNonce:                  pb.EthereumNonce,
		StakingInfo:                    &stakingInfo,
	}, nil
}

func (info AccountInfo) _ToProtobuf() *services.CryptoGetInfoResponse_AccountInfo {
	tokenRelationship := make([]*services.TokenRelationship, len(info.TokenRelationships))

	for i, relationship := range info.TokenRelationships {
		singleRelationship := relationship._ToProtobuf()
		tokenRelationship[i] = singleRelationship
	}

	liveHashes := make([]*services.LiveHash, len(info.LiveHashes))

	for i, liveHash := range info.LiveHashes {
		singleRelationship := liveHash._ToProtobuf()
		liveHashes[i] = singleRelationship
	}

	var alias []byte
	if info.AliasKey != nil {
		alias, _ = protobuf.Marshal(info.AliasKey._ToProtoKey())
	}

	body := &services.CryptoGetInfoResponse_AccountInfo{
		AccountID:                      info.AccountID._ToProtobuf(),
		ContractAccountID:              info.ContractAccountID,
		Deleted:                        info.IsDeleted,
		ProxyReceived:                  info.ProxyReceived.tinybar,
		Key:                            info.Key._ToProtoKey(),
		Balance:                        uint64(info.Balance.tinybar),
		GenerateSendRecordThreshold:    uint64(info.GenerateSendRecordThreshold.tinybar),
		GenerateReceiveRecordThreshold: uint64(info.GenerateReceiveRecordThreshold.tinybar),
		ReceiverSigRequired:            info.ReceiverSigRequired,
		ExpirationTime:                 _TimeToProtobuf(info.ExpirationTime),
		AutoRenewPeriod:                _DurationToProtobuf(info.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		TokenRelationships:             tokenRelationship,
		Memo:                           info.AccountMemo,
		OwnedNfts:                      info.OwnedNfts,
		MaxAutomaticTokenAssociations:  int32(info.MaxAutomaticTokenAssociations),
		Alias:                          alias,
		LedgerId:                       info.LedgerID.ToBytes(),
		EthereumNonce:                  info.EthereumNonce,
	}

	if info.StakingInfo != nil {
		body.StakingInfo = info.StakingInfo._ToProtobuf()
	}

	return body
}

func (info AccountInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(info._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func AccountInfoFromBytes(data []byte) (AccountInfo, error) {
	if data == nil {
		return AccountInfo{}, errByteArrayNull
	}
	pb := services.CryptoGetInfoResponse_AccountInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AccountInfo{}, err
	}

	info, err := _AccountInfoFromProtobuf(&pb)
	if err != nil {
		return AccountInfo{}, err
	}

	return info, nil
}
