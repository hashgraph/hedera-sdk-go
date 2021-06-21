package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenID struct {
	Shard    uint64
	Realm    uint64
	Token    uint64
	Checksum *string
	Network  *NetworkName
}

func tokenIDFromProtobuf(tokenID *proto.TokenID) TokenID {
	if tokenID == nil {
		return TokenID{}
	}
	return TokenID{
		Shard: uint64(tokenID.ShardNum),
		Realm: uint64(tokenID.RealmNum),
		Token: uint64(tokenID.TokenNum),
	}
}

func (id *TokenID) toProtobuf() *proto.TokenID {
	return &proto.TokenID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TokenNum: int64(id.Token),
	}
}

func (id TokenID) String() string {
	checksum, err := checksumParseAddress("", fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Token, checksum.correctChecksum)
}

func (id TokenID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenIDFromBytes(data []byte) (TokenID, error) {
	if data == nil {
		return TokenID{}, errByteArrayNull
	}
	pb := proto.TokenID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenID{}, err
	}

	return tokenIDFromProtobuf(&pb), nil
}

// TokenIDFromString constructs an TokenID from a string formatted as
// `Shard.Realm.TokenID` (for example "0.0.3")
func TokenIDFromString(data string) (TokenID, error) {
	var checksum parseAddressResult
	var err error

	var networkNames = []NetworkName{
		Mainnet,
		Testnet,
		Previewnet,
	}

	var network NetworkName
	for _, name := range networkNames {
		checksum, err = checksumParseAddress(name.Network(), data)
		if err != nil {
			return TokenID{}, err
		}
		if checksum.status != 1 {
			network = name
			break
		}
	}

	err = checksumVerify(checksum.status)
	if err != nil {
		return TokenID{}, err
	}

	tempChecksum := checksum.correctChecksum

	return TokenID{
		Shard:    uint64(checksum.num1),
		Realm:    uint64(checksum.num2),
		Token:    uint64(checksum.num3),
		Checksum: &tempChecksum,
		Network:  &network,
	}, nil
}

func TokenIDValidateNetworkOnIDs(id TokenID, other AccountID) error {
	if !id.isZero() && !other.isZero() && id.Network != nil && other.Network != nil && *id.Network != *other.Network {
		return errNetworkMismatch
	}

	return nil
}

func (id TokenID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Token == 0
}
