package hedera

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenID struct {
	Shard    uint64
	Realm    uint64
	Token    uint64
	checksum *string
}

func _TokenIDFromProtobuf(tokenID *proto.TokenID) *TokenID {
	if tokenID == nil {
		return nil
	}

	return &TokenID{
		Shard: uint64(tokenID.ShardNum),
		Realm: uint64(tokenID.RealmNum),
		Token: uint64(tokenID.TokenNum),
	}
}

func (id *TokenID) _ToProtobuf() *proto.TokenID {
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
	if client.network.networkName == nil {
		return "", errNetworkNameMissing
	}
	checksum, err := _ChecksumParseAddress(client.network.networkName._LedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Token))
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
	pb := proto.TokenID{}
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

// Deprecated
func (id *TokenID) Validate(client *Client) error {
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

func (id TokenID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Token == 0
}
