package hedera

import (
	"fmt"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// TokenID is the ID for a Hedera token
type TokenID struct {
	Shard uint64
	Realm uint64
	Token uint64
}

// TokenIDFromString constructs an TokenID from a string formatted as
// `Shard.Realm.Token` (for example "0.0.3")
func TokenIDFromString(s string) (TokenID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return TokenID{}, err
	}

	return TokenID{
		Shard: uint64(shard),
		Realm: uint64(realm),
		Token: uint64(num),
	}, nil
}

// TokenIDFromSolidityAddress constructs an TokenID from a string
// representation of a solidity address
func TokenIDFromSolidityAddress(s string) (TokenID, error) {
	shard, realm, token, err := idFromSolidityAddress(s)
	if err != nil {
		return TokenID{}, err
	}

	return TokenID{
		Shard: shard,
		Realm: realm,
		Token: token,
	}, nil
}

// String returns the string representation of an TokenID in
// `Shard.Realm.Token` (for example "0.0.3")
func (id TokenID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token)
}

// ToSolidityAddress returns the string representation of the TokenID as a
// solidity address.
func (id TokenID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.Token)
}

func (id TokenID) toProto() *proto.TokenID {
	return &proto.TokenID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TokenNum: int64(id.Token),
	}
}

// UnmarshalJSON implements the encoding.JSON interface.
func (id *TokenID) UnmarshalJSON(data []byte) error {
	tokenID, err := TokenIDFromString(strings.Replace(string(data), "\"", "", 2))

	if err != nil {
		println("error was not nil")
		return err
	}

	id = &tokenID

	return nil
}

func tokenIDFromProto(pb *proto.TokenID) TokenID {
	return TokenID{
		Shard: uint64(pb.ShardNum),
		Realm: uint64(pb.RealmNum),
		Token: uint64(pb.TokenNum),
	}
}
