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
	checksum *string
}

func tokenIDFromProtobuf(tokenID *proto.TokenID) TokenID {
	if tokenID == nil {
		return TokenID{}
	}

	id := TokenID{
		Shard: uint64(tokenID.ShardNum),
		Realm: uint64(tokenID.RealmNum),
		Token: uint64(tokenID.TokenNum),
	}

	return id
}

func (id *TokenID) toProtobuf() *proto.TokenID {
	return &proto.TokenID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TokenNum: int64(id.Token),
	}
}

func (id TokenID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token)
}

func (id TokenID) ToStringWithChecksum(client Client) (string, error) {
	if client.networkName == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := checksumParseAddress(client.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Token, checksum.correctChecksum), nil
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

func (id *TokenID) Nft(serial int64) NftID {
	return NftID{
		TokenID:      *id,
		SerialNumber: serial,
	}
}

// TokenIDFromString constructs an TokenID from a string formatted as
// `Shard.Realm.TokenID` (for example "0.0.3")
func TokenIDFromString(data string) (TokenID, error) {
	shard, realm, num, checksum, err := idFromString(data)
	if err != nil {
		return TokenID{}, err
	}

	return TokenID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Token:    uint64(num),
		checksum: checksum,
	}, nil
}

func (id *TokenID) Validate(client *Client) error {
	if !id.isZero() && client != nil && client.networkName != nil {
		tempChecksum, err := checksumParseAddress(client.networkName.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
		if err != nil {
			return err
		}
		err = checksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}
		if id.checksum == nil {
			id.checksum = &tempChecksum.correctChecksum
			return nil
		}
		if tempChecksum.correctChecksum != *id.checksum {
			return errNetworkMismatch
		}
	}

	return nil
}

func (id *TokenID) setNetworkWithClient(client *Client) {
	if client.networkName != nil {
		id.setNetwork(*client.networkName)
	}
}

func (id *TokenID) setNetwork(name NetworkName) {
	checksum := checkChecksum(name.ledgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
	id.checksum = &checksum
}

func (id TokenID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Token == 0
}
