package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
	"strings"
)

// AccountID is the ID for a Hedera account
type AccountID struct {
	Shard    uint64
	Realm    uint64
	Account  uint64
	Checksum *string
	Network  *NetworkName
}

// AccountIDFromString constructs an AccountID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func AccountIDFromString(data string) (AccountID, error) {
	var checksum parseAddressResult
	var err error

	var networkNames = []NetworkName{
		NetworkNameMainnet,
		NetworkNameTestnet,
		NetworkNamePreviewnet,
	}

	var network NetworkName
	for _, name := range networkNames {
		checksum, err = checksumParseAddress(name.Network(), data)
		if err != nil {
			return AccountID{}, err
		}
		if checksum.status != 1 {
			network = name
			break
		}
	}

	err = checksumVerify(checksum.status)
	if err != nil {
		return AccountID{}, err
	}

	tempChecksum := checksum.correctChecksum

	return AccountID{
		Shard:    uint64(checksum.num1),
		Realm:    uint64(checksum.num2),
		Account:  uint64(checksum.num3),
		Checksum: &tempChecksum,
		Network:  &network,
	}, nil
}

func (id *AccountID) SetNetworkName(network NetworkName) {
	id.Network = &network
	checksum := checkChecksum(id.Network.Network(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
	id.Checksum = &checksum
}

// AccountIDFromSolidityAddress constructs an AccountID from a string
// representation of a solidity address
func AccountIDFromSolidityAddress(s string) (AccountID, error) {
	shard, realm, account, err := idFromSolidityAddress(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:    shard,
		Realm:    realm,
		Account:  account,
		Checksum: nil,
		Network:  nil,
	}, nil
}

func (id *AccountID) GetNetworkFromChecksum() error {
	if id.Checksum == nil {
		return errors.New("Checksum is missing.")
	}
	var checksum parseAddressResult
	var err error

	var networkNames = []NetworkName{
		NetworkNameMainnet,
		NetworkNameTestnet,
		NetworkNamePreviewnet,
	}

	var network NetworkName
	for _, name := range networkNames {
		checksum, err = checksumParseAddress(name.Network(), fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Account, *id.Checksum))
		if err != nil {
			return err
		}
		if checksum.status != 1 {
			network = name
			break
		}
	}

	id.Network = &network
	return nil
}

func AccountIDValidateNetworkOnIDs(id AccountID, other *Client) error {
	if !id.isZero() && other != nil && id.Network != nil && other.networkName != nil && *id.Network != *other.networkName {
		return errNetworkMismatch
	}

	return nil
}

// String returns the string representation of an AccountID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id AccountID) String() string {
	if id.Network == nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Account, *id.Checksum)
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

	*id = accountID

	return nil
}

func accountIDFromProtobuf(pb *proto.AccountID) AccountID {
	if pb == nil {
		return AccountID{}
	}
	return AccountID{
		Shard:    uint64(pb.ShardNum),
		Realm:    uint64(pb.RealmNum),
		Account:  uint64(pb.AccountNum),
		Checksum: nil,
		Network:  nil,
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
	if data == nil {
		return AccountID{}, errByteArrayNull
	}
	pb := proto.AccountID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AccountID{}, err
	}

	return accountIDFromProtobuf(&pb), nil
}
