package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// TokenNftInfo is the information about a NFT
type TokenNftInfo struct {
	NftID        NftID
	AccountID    AccountID
	CreationTime time.Time
	Metadata     []byte
	LedgerID     LedgerID
	SpenderID    AccountID
}

func _TokenNftInfoFromProtobuf(pb *services.TokenNftInfo) TokenNftInfo {
	if pb == nil {
		return TokenNftInfo{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	spenderID := AccountID{}
	if pb.SpenderId != nil {
		spenderID = *_AccountIDFromProtobuf(pb.SpenderId)
	}

	return TokenNftInfo{
		NftID:        _NftIDFromProtobuf(pb.NftID),
		AccountID:    accountID,
		CreationTime: _TimeFromProtobuf(pb.CreationTime),
		Metadata:     pb.Metadata,
		LedgerID:     LedgerID{pb.LedgerId},
		SpenderID:    spenderID,
	}
}

func (tokenNftInfo *TokenNftInfo) _ToProtobuf() *services.TokenNftInfo {
	return &services.TokenNftInfo{
		NftID:        tokenNftInfo.NftID._ToProtobuf(),
		AccountID:    tokenNftInfo.AccountID._ToProtobuf(),
		CreationTime: _TimeToProtobuf(tokenNftInfo.CreationTime),
		Metadata:     tokenNftInfo.Metadata,
		LedgerId:     tokenNftInfo.LedgerID.ToBytes(),
		SpenderId:    tokenNftInfo.SpenderID._ToProtobuf(),
	}
}

// ToBytes returns the byte representation of the TokenNftInfo
func (tokenNftInfo *TokenNftInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(tokenNftInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenNftInfoFromBytes returns the TokenNftInfo from a byte array representation
func TokenNftInfoFromBytes(data []byte) (TokenNftInfo, error) {
	if data == nil {
		return TokenNftInfo{}, errByteArrayNull
	}
	pb := services.TokenNftInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenNftInfo{}, err
	}

	return _TokenNftInfoFromProtobuf(&pb), nil
}
