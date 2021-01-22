package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// AccountID is the ID for a Hedera account
type AccountID struct {
	Shard   uint64
	Realm   uint64
	Account uint64
}

// AccountIDFromString constructs an AccountID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func AccountIDFromString(s string) (AccountID, error) {
	shard, realm, num, err := idFromString(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:   uint64(shard),
		Realm:   uint64(realm),
		Account: uint64(num),
	}, nil
}

// AccountIDFromSolidityAddress constructs an AccountID from a string
// representation of a solidity address
func AccountIDFromSolidityAddress(s string) (AccountID, error) {
	shard, realm, account, err := idFromSolidityAddress(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:   shard,
		Realm:   realm,
		Account: account,
	}, nil
}

// String returns the string representation of an AccountID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id AccountID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
}

// ToSolidityAddress returns the string representation of the AccountID as a
// solidity address.
func (id AccountID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.Account)
}

func (id AccountID) toProtobuf() *proto.AccountID {
	return &proto.AccountID{
		ShardNum:   int64(id.Shard),
		RealmNum:   int64(id.Realm),
		AccountNum: int64(id.Account),
	}
}

// UnmarshalJSON implements the encoding.JSON interface.
func (id *AccountID) UnmarshalJSON(data []byte) error {
	accountID, err := AccountIDFromString(strings.Replace(string(data), "\"", "", 2))

	if err != nil {
		return err
	}

	id = &accountID

	return nil
}

func accountIDFromProtobuf(pb *proto.AccountID) AccountID {
	return AccountID{
		Shard:   uint64(pb.ShardNum),
		Realm:   uint64(pb.RealmNum),
		Account: uint64(pb.AccountNum),
	}
}

func (id AccountID) isZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Account == 0
}

func (id AccountID) equals(other AccountID) bool {
	return id.Shard == other.Shard && id.Realm == other.Realm && id.Account == other.Account
}

func (id AccountID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func AccountIDFromBytes(data []byte) (AccountID, error) {
	pb := proto.AccountID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AccountID{}, err
	}

	return accountIDFromProtobuf(&pb), nil
}
