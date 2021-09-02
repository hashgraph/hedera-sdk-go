package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenNftInfo struct {
	NftID        NftID
	AccountID    AccountID
	CreationTime time.Time
	Metadata     []byte
}

func tokenNftInfoFromProtobuf(pb *proto.TokenNftInfo) TokenNftInfo {
	if pb == nil {
		return TokenNftInfo{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *accountIDFromProtobuf(pb.AccountID)
	}

	return TokenNftInfo{
		NftID:        nftIDFromProtobuf(pb.NftID),
		AccountID:    accountID,
		CreationTime: timeFromProtobuf(pb.CreationTime),
		Metadata:     pb.Metadata,
	}
}

func (tokenNftInfo *TokenNftInfo) toProtobuf() *proto.TokenNftInfo {
	return &proto.TokenNftInfo{
		NftID:        tokenNftInfo.NftID.toProtobuf(),
		AccountID:    tokenNftInfo.AccountID.toProtobuf(),
		CreationTime: timeToProtobuf(tokenNftInfo.CreationTime),
		Metadata:     tokenNftInfo.Metadata,
	}
}

func (tokenNftInfo *TokenNftInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(tokenNftInfo.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenNftInfoFromBytes(data []byte) (TokenNftInfo, error) {
	if data == nil {
		return TokenNftInfo{}, errByteArrayNull
	}
	pb := proto.TokenNftInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenNftInfo{}, err
	}

	return tokenNftInfoFromProtobuf(&pb), nil
}
