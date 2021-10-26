package hedera

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

// AccountID is the ID for a Hedera account
type AccountID struct {
	Shard    uint64
	Realm    uint64
	Account  uint64
	checksum *string
}

// AccountIDFromString constructs an AccountID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func AccountIDFromString(data string) (AccountID, error) {
	shard, realm, num, checksum, err := _IdFromString(data)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:    uint64(shard),
		Realm:    uint64(realm),
		Account:  uint64(num),
		checksum: checksum,
	}, nil
}

// AccountIDFromSolidityAddress constructs an AccountID from a string
// representation of a _Solidity address
func AccountIDFromSolidityAddress(s string) (AccountID, error) {
	shard, realm, account, err := _IdFromSolidityAddress(s)
	if err != nil {
		return AccountID{}, err
	}

	return AccountID{
		Shard:    shard,
		Realm:    realm,
		Account:  account,
		checksum: nil,
	}, nil
}

func (id *AccountID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
		if err != nil {
			return err
		}
		err = _ChecksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}
		if id.checksum == nil {
			return errChecksumMissing
		}
		if tempChecksum.correctChecksum != *id.checksum {
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				*client.network.networkName))
		}
	}

	return nil
}

// Deprecated
func (id *AccountID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
		if err != nil {
			return err
		}
		err = _ChecksumVerify(tempChecksum.status)
		if err != nil {
			return err
		}

		if tempChecksum.correctChecksum != *id.checksum {
			return errNetworkMismatch
		}
	}

	return nil
}

// String returns the string representation of an AccountID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id AccountID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
}

func (id AccountID) ToStringWithChecksum(client *Client) (string, error) {
	if client.network.networkName == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Account, checksum.correctChecksum), nil
}

func (id AccountID) GetChecksum() *string {
	return id.checksum
}

// ToSolidityAddress returns the string representation of the AccountID as a
// _Solidity address.
func (id AccountID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Account)
}

func (id AccountID) _ToProtobuf() *proto.AccountID {
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

func _AccountIDFromProtobuf(accountID *proto.AccountID) *AccountID {
	if accountID == nil {
		return nil
	}

	return &AccountID{
		Shard:   uint64(accountID.ShardNum),
		Realm:   uint64(accountID.RealmNum),
		Account: uint64(accountID.AccountNum),
	}
}

func (id AccountID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Account == 0
}

func (id AccountID) _Equals(other AccountID) bool {
	return id.Shard == other.Shard && id.Realm == other.Realm && id.Account == other.Account
}

func (id AccountID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
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

	return *_AccountIDFromProtobuf(&pb), nil
}
