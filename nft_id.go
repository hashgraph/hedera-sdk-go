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
	checksum     *string
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
	if !id.isZero() && client != nil && client.networkName != nil {
		return id.TokenID.Validate(client)
	}

	return nil
}

func (id *NftID) setNetworkWithClient(client *Client) {
	if client.networkName != nil {
		id.setNetwork(*client.networkName)
	}
}

func (id *NftID) setNetwork(name NetworkName) {
	checksum := checkChecksum(name.ledgerID(), fmt.Sprintf("%d.%d.%d", id.TokenID.Shard, id.TokenID.Realm, id.TokenID.Token))
	id.checksum = &checksum
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
	return NftID{
		TokenID:      tokenIDFromProtobuf(pb.TokenID),
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
