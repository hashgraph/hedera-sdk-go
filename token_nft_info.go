package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

type TokenNftInfo struct {
	NftID        NftID
	AccountID    AccountID
	CreationTime time.Time
	Metadata     []byte
}

func tokenNftInfoFromProtobuf(pb *services.TokenNftInfo, networkName *NetworkName) TokenNftInfo {
	if pb == nil {
		return TokenNftInfo{}
	}

	return TokenNftInfo{
		NftID:        nftIDFromProtobuf(pb.NftID, networkName),
		AccountID:    accountIDFromProtobuf(pb.AccountID, networkName),
		CreationTime: timeFromProtobuf(pb.CreationTime),
		Metadata:     pb.Metadata,
	}
}

func (tokenNftInfo *TokenNftInfo) toProtobuf() *services.TokenNftInfo {
	return &services.TokenNftInfo{
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
	pb := services.TokenNftInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenNftInfo{}, err
	}

	return tokenNftInfoFromProtobuf(&pb, nil), nil
}
