package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"strconv"
	"strings"
)

type NftID struct {
	TokenID      TokenID
	SerialNumber int64
}

func NftIDFromString(s string) (NftID, error) {
	split := strings.Split(s, "@")
	shard, realm, num, err := idFromString(split[1])
	if err != nil {
		return NftID{}, err
	}

	serial, err := strconv.Atoi(split[0])
	if err != nil {
		return NftID{}, err
	}

	return NftID{
		TokenID: TokenID{
			Shard: uint64(shard),
			Realm: uint64(realm),
			Token: uint64(num),
		},
		SerialNumber: int64(serial),
	}, nil
}

func (id NftID) String() string {
	return fmt.Sprintf("%d@%s", id.SerialNumber, id.TokenID.String())
}

func (id NftID) toProtobuf() *proto.NftID {
	return &proto.NftID{
		TokenID:      id.TokenID.toProtobuf(),
		SerialNumber: id.SerialNumber,
	}
}

func nftIDFromProtobuf(pb *proto.NftID) NftID {
	if pb == nil {
		return NftID{}
	}
	return NftID{
		TokenID:      tokenIDFromProtobuf(pb.TokenID),
		SerialNumber: pb.SerialNumber,
	}
}

func (id NftID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func NftIDFromBytes(data []byte) (NftID, error) {
	pb := proto.NftID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return NftID{}, err
	}

	return nftIDFromProtobuf(&pb), nil
}
