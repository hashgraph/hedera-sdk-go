package hedera

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type NftID struct {
	TokenID      TokenID
	SerialNumber int64
}

func NftIDFromString(s string) (NftID, error) {
	split := strings.Split(s, "@")
	shard, realm, num, checksum, err := idFromString(split[1])
	if err != nil {
		return NftID{}, err
	}

	serial, err := strconv.Atoi(split[0])
	if err != nil {
		return NftID{}, err
	}

	return NftID{
		TokenID: TokenID{
			Shard:    uint64(shard),
			Realm:    uint64(realm),
			Token:    uint64(num),
			checksum: checksum,
		},
		SerialNumber: int64(serial),
	}, nil
}

func (id *NftID) Validate(client *Client) error {
	if !id.isZero() && client != nil && client.network.networkName != nil {
		if err := id.TokenID.Validate(client); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (id NftID) String() string {
	return fmt.Sprintf("%d@%s", id.SerialNumber, id.TokenID.String())
}

func (id NftID) ToStringWithChecksum(client Client) (string, error) {
	token, err := id.TokenID.ToStringWithChecksum(client)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d@%s", id.SerialNumber, token), nil
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

	tokenID := TokenID{}
	if pb.TokenID != nil {
		tokenID = *tokenIDFromProtobuf(pb.TokenID)
	}

	return NftID{
		TokenID:      tokenID,
		SerialNumber: pb.SerialNumber,
	}
}

func (id NftID) isZero() bool {
	return id.TokenID.isZero() && id.SerialNumber == 0
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
