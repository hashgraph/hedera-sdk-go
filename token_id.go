package hedera

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenID struct {
	Shard    uint64
	Realm    uint64
	Token    uint64
	checksum *string
}

type _TokenIDs struct {
	tokenIDs []TokenID
}

func _TokenIDFromProtobuf(tokenID *services.TokenID) *TokenID {
	if tokenID == nil {
		return nil
	}

	return &TokenID{
		Shard: uint64(tokenID.ShardNum),
		Realm: uint64(tokenID.RealmNum),
		Token: uint64(tokenID.TokenNum),
	}
}

func (id *TokenID) _ToProtobuf() *services.TokenID {
	return &services.TokenID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
		TokenNum: int64(id.Token),
	}
}

func (id TokenID) String() string {
	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token)
}

func (id TokenID) ToStringWithChecksum(client Client) (string, error) {
	if client.GetNetworkName() == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Token, checksum.correctChecksum), nil
}

func (id TokenID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenIDFromBytes(data []byte) (TokenID, error) {
	if data == nil {
		return TokenID{}, errByteArrayNull
	}
	pb := services.TokenID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenID{}, err
	}

	return *_TokenIDFromProtobuf(&pb), nil
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
	shard, realm, num, checksum, err := _IdFromString(data)
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

func (id *TokenID) ValidateChecksum(client *Client) error {
	if !id._IsZero() && client != nil && client.network.networkName != nil {
		tempChecksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
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

// TokenIDFromSolidityAddress constructs a TokenID from a string
// representation of a _Solidity address
func TokenIDFromSolidityAddress(s string) (TokenID, error) {
	shard, realm, account, err := _IdFromSolidityAddress(s)
	if err != nil {
		return AccountID{}, err
	}

	return TokenID{
		Shard:    shard,
		Realm:    realm,
		Account:  account,
		checksum: nil,
	}, nil
}

// ToSolidityAddress returns the string representation of the TokenID as a
// _Solidity address.
func (id TokenID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Token)
}

// Deprecated
func (id *TokenID) Validate(client *Client) error {
	if !id._IsZero() && client != nil && client.GetNetworkName() != nil {
		tempChecksum, err := _ChecksumParseAddress(client.GetNetworkName()._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
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

func (id TokenID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Token == 0
}

func (id TokenID) Compare(given TokenID) int {
	if id.Shard > given.Shard { //nolint
		return 1
	} else if id.Shard < given.Shard {
		return -1
	}

	if id.Realm > given.Realm { //nolint
		return 1
	} else if id.Realm < given.Realm {
		return -1
	}

	if id.Token > given.Token { //nolint
		return 1
	} else if id.Token < given.Token {
		return -1
	} else { //nolint
		return 0
	}
}

func (tokenIDs _TokenIDs) Len() int {
	return len(tokenIDs.tokenIDs)
}
func (tokenIDs _TokenIDs) Swap(i, j int) {
	tokenIDs.tokenIDs[i], tokenIDs.tokenIDs[j] = tokenIDs.tokenIDs[j], tokenIDs.tokenIDs[i]
}

func (tokenIDs _TokenIDs) Less(i, j int) bool {
	if tokenIDs.tokenIDs[i].Shard < tokenIDs.tokenIDs[j].Shard { //nolint
		return true
	} else if tokenIDs.tokenIDs[i].Shard > tokenIDs.tokenIDs[j].Shard {
		return false
	}

	if tokenIDs.tokenIDs[i].Realm < tokenIDs.tokenIDs[j].Realm { //nolint
		return true
	} else if tokenIDs.tokenIDs[i].Realm > tokenIDs.tokenIDs[j].Realm {
		return false
	}

	if tokenIDs.tokenIDs[i].Token < tokenIDs.tokenIDs[j].Token { //nolint
		return true
	} else { //nolint
		return false
	}
}
