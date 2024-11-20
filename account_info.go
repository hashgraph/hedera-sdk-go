package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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

	var tokenRelationships []*TokenRelationship
	if pb.TokenRelationships != nil { // nolint
		tokenRelationships = _TokenRelationshipsFromProtobuf(pb.TokenRelationships) // nolint
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
		ExpirationTime:                 _TimeFromProtobuf(pb.ExpirationTime),
		AccountMemo:                    pb.Memo,
		AutoRenewPeriod:                _DurationFromProtobuf(pb.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		TokenRelationships:             tokenRelationships,
		OwnedNfts:                      pb.OwnedNfts,
		MaxAutomaticTokenAssociations:  uint32(pb.MaxAutomaticTokenAssociations),
		AliasKey:                       alias,
		LedgerID:                       LedgerID{pb.LedgerId},
		EthereumNonce:                  pb.EthereumNonce,
		StakingInfo:                    &stakingInfo,
	}, nil
}

func (info AccountInfo) _ToProtobuf() *services.CryptoGetInfoResponse_AccountInfo {
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

// ToBytes returns the serialized bytes of an AccountInfo
func (info AccountInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(info._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// AccountInfoFromBytes returns an AccountInfo from byte array
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
