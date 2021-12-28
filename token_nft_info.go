package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenNftInfo struct {
	NftID        NftID
	AccountID    AccountID
	CreationTime time.Time
	Metadata     []byte
	LedgerID     []byte
}

func _TokenNftInfoFromProtobuf(pb *services.TokenNftInfo) TokenNftInfo {
	if pb == nil {
		return TokenNftInfo{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return TokenNftInfo{
		NftID:        _NftIDFromProtobuf(pb.NftID),
		AccountID:    accountID,
		CreationTime: _TimeFromProtobuf(pb.CreationTime),
		Metadata:     pb.Metadata,
		LedgerID:     pb.LedgerId,
	}
}

func (tokenNftInfo *TokenNftInfo) _ToProtobuf() *services.TokenNftInfo {
	return &services.TokenNftInfo{
		NftID:        tokenNftInfo.NftID._ToProtobuf(),
		AccountID:    tokenNftInfo.AccountID._ToProtobuf(),
		CreationTime: _TimeToProtobuf(tokenNftInfo.CreationTime),
		Metadata:     tokenNftInfo.Metadata,
		LedgerId:     tokenNftInfo.LedgerID,
	}
}

func (tokenNftInfo *TokenNftInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(tokenNftInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

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
