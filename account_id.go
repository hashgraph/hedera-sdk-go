package hiero

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// AccountID is the ID for a Hiero account
type AccountID struct {
	Shard           uint64
	Realm           uint64
	Account         uint64
	AliasKey        *PublicKey
	AliasEvmAddress *[]byte
	checksum        *string
}

type _AccountIDs struct { //nolint
	accountIDs []AccountID
}

// AccountIDFromString constructs an AccountID from a string formatted as
// `Shard.Realm.Account` (for example "0.0.3")
func AccountIDFromString(data string) (AccountID, error) {
	shard, realm, num, checksum, alias, aliasEvmAddress, err := _AccountIDFromString(data)
	if err != nil {
		return AccountID{}, err
	}

	if num == -1 {
		if alias != nil {
			return AccountID{
				Shard:           uint64(shard),
				Realm:           uint64(realm),
				Account:         0,
				AliasKey:        alias,
				AliasEvmAddress: nil,
				checksum:        checksum,
			}, nil
		}

		return AccountID{
			Shard:           uint64(shard),
			Realm:           uint64(realm),
			Account:         0,
			AliasKey:        nil,
			AliasEvmAddress: aliasEvmAddress,
			checksum:        checksum,
		}, nil
	}

	return AccountID{
		Shard:           uint64(shard),
		Realm:           uint64(realm),
		Account:         uint64(num),
		AliasKey:        nil,
		AliasEvmAddress: nil,
		checksum:        checksum,
	}, nil
}

// AccountIDFromEvmAddress constructs an AccountID from a string formatted as 0.0.<evm address>
func AccountIDFromEvmAddress(shard uint64, realm uint64, aliasEvmAddress string) (AccountID, error) {
	temp, err := hex.DecodeString(aliasEvmAddress)
	if err != nil {
		return AccountID{}, err
	}
	return AccountID{
		Shard:           shard,
		Realm:           realm,
		Account:         0,
		AliasEvmAddress: &temp,
		checksum:        nil,
	}, nil
}

// Returns an AccountID with EvmPublic address for the use of HIP-583
func AccountIDFromEvmPublicAddress(s string) (AccountID, error) {
	return AccountIDFromString(s)
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

// Verify that the client has a valid checksum.
func (id *AccountID) ValidateChecksum(client *Client) error {
	if id.AliasKey != nil {
		return errors.New("Account ID contains alias key, unable to validate")
	}
	if !id._IsZero() && client != nil {
		var tempChecksum _ParseAddressResult
		var err error
		tempChecksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
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
			networkName := NetworkNameOther
			if client.network.ledgerID != nil {
				networkName, _ = client.network.ledgerID.ToNetworkName()
			}
			return errors.New(fmt.Sprintf("network mismatch or wrong checksum given, given checksum: %s, correct checksum %s, network: %s",
				*id.checksum,
				tempChecksum.correctChecksum,
				networkName))
		}
	}

	return nil
}

// Deprecated - use ValidateChecksum instead
func (id *AccountID) Validate(client *Client) error {
	return id.ValidateChecksum(client)
}

// String returns the string representation of an AccountID in
// `Shard.Realm.Account` (for example "0.0.3")
func (id AccountID) String() string {
	if id.AliasKey != nil {
		return fmt.Sprintf("%d.%d.%s", id.Shard, id.Realm, id.AliasKey.String())
	} else if id.AliasEvmAddress != nil {
		return fmt.Sprintf("%d.%d.%s", id.Shard, id.Realm, hex.EncodeToString(*id.AliasEvmAddress))
	}

	return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account)
}

// ToStringWithChecksum returns the string representation of an AccountID in
// `Shard.Realm.Account-checksum` (for example "0.0.3-sdaf")
func (id AccountID) ToStringWithChecksum(client *Client) (string, error) {
	if id.AliasKey != nil {
		return "", errors.New("Account ID contains alias key, unable get checksum")
	}
	if client.GetNetworkName() == nil && client.GetLedgerID() == nil {
		return "", errNetworkNameMissing
	}
	var checksum _ParseAddressResult
	var err error
	if client.network.ledgerID != nil {
		checksum, err = _ChecksumParseAddress(client.GetLedgerID(), fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Account))
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Account, checksum.correctChecksum), nil
}

// GetChecksum Retrieve just the checksum
func (id AccountID) GetChecksum() *string {
	return id.checksum
}

// ToSolidityAddress returns the string representation of the AccountID as a
// _Solidity address.
func (id AccountID) ToSolidityAddress() string {
	return _IdToSolidityAddress(id.Shard, id.Realm, id.Account)
}

func (id AccountID) _ToProtobuf() *services.AccountID {
	resultID := &services.AccountID{
		ShardNum: int64(id.Shard),
		RealmNum: int64(id.Realm),
	}
	if id.AliasKey != nil {
		data, _ := protobuf.Marshal(id.AliasKey._ToProtoKey())
		resultID.Account = &services.AccountID_Alias{
			Alias: data,
		}

		return resultID
	} else if id.AliasEvmAddress != nil {
		resultID.Account = &services.AccountID_Alias{
			Alias: *id.AliasEvmAddress,
		}

		return resultID
	}

	resultID.Account = &services.AccountID_AccountNum{
		AccountNum: int64(id.Account),
	}

	return resultID
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

func _AccountIDFromProtobuf(accountID *services.AccountID) *AccountID {
	if accountID == nil {
		return nil
	}
	resultAccountID := &AccountID{
		Shard: uint64(accountID.ShardNum),
		Realm: uint64(accountID.RealmNum),
	}

	switch t := accountID.Account.(type) {
	case *services.AccountID_Alias:
		pb := services.Key{}
		_ = protobuf.Unmarshal(t.Alias, &pb)
		initialKey, err := _KeyFromProtobuf(&pb)
		if err != nil && t.Alias != nil {
			resultAccountID.Account = 0
			resultAccountID.AliasEvmAddress = &t.Alias
			return resultAccountID
		}
		if evm, ok := pb.Key.(*services.Key_ECDSASecp256K1); ok && len(evm.ECDSASecp256K1) == 20 {
			resultAccountID.Account = 0
			resultAccountID.AliasEvmAddress = &evm.ECDSASecp256K1
			return resultAccountID
		}
		switch t2 := initialKey.(type) {
		case PublicKey:
			resultAccountID.Account = 0
			resultAccountID.AliasKey = &t2
			return resultAccountID
		default:
			return &AccountID{}
		}
	case *services.AccountID_AccountNum:
		resultAccountID.Account = uint64(t.AccountNum)
		resultAccountID.AliasKey = nil
		return resultAccountID
	default:
		return &AccountID{}
	}
}

// IsZero returns true if this AccountID is the zero-value
func (id AccountID) IsZero() bool {
	return id._IsZero()
}

func (id AccountID) _IsZero() bool {
	return id.Shard == 0 && id.Realm == 0 && id.Account == 0 && id.AliasKey == nil
}

// Equals returns true if this AccountID and the given AccountID are identical
func (id AccountID) Equals(other AccountID) bool {
	return id._Equals(other)
}

func (id AccountID) _Equals(other AccountID) bool {
	initialAlias := ""
	otherAlias := ""
	if id.AliasKey != nil && other.AliasKey != nil {
		initialAlias = id.AliasKey.String()
		otherAlias = other.AliasKey.String()
	}

	return id.Shard == other.Shard && id.Realm == other.Realm && id.Account == other.Account && initialAlias == otherAlias
}

// ToBytes returns the wire-format encoding of AccountID
func (id AccountID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// AccountIDFromBytes converts wire-format encoding to Account ID
func AccountIDFromBytes(data []byte) (AccountID, error) {
	if data == nil {
		return AccountID{}, errByteArrayNull
	}
	pb := services.AccountID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AccountID{}, err
	}

	return *_AccountIDFromProtobuf(&pb), nil
}

type PopulateType int

const (
	Account PopulateType = iota
	EvmAddress
)

func (id *AccountID) _MirrorNodeRequest(client *Client, populateType string) (map[string]interface{}, error) {
	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return nil, errors.New("mirror node is not set")
	}

	mirrorUrl := client.GetMirrorNetwork()[0]
	index := strings.Index(mirrorUrl, ":")
	if index == -1 {
		return nil, errors.New("invalid mirrorUrl format")
	}
	mirrorUrl = mirrorUrl[:index]

	var url string
	protocol := "https"
	port := ""

	if client.GetLedgerID().String() == "" {
		protocol = "http"
		port = ":5551"
	}

	if populateType == "account" {
		url = fmt.Sprintf("%s://%s%s/api/v1/accounts/%s", protocol, mirrorUrl, port, hex.EncodeToString(*id.AliasEvmAddress))
	} else {
		url = fmt.Sprintf("%s://%s%s/api/v1/accounts/%s", protocol, mirrorUrl, port, id.String())
	}

	resp, err := http.Get(url) // #nosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// PopulateAccount gets the actual `Account` field of the `AccountId` from the Mirror Node.
// Should be used after generating `AccountId.FromEvmAddress()` because it sets the `Account` field to `0`
// automatically since there is no connection between the `Account` and the `evmAddress`
func (id *AccountID) PopulateAccount(client *Client) error {
	result, err := id._MirrorNodeRequest(client, "account")
	if err != nil {
		return err
	}

	mirrorAccountId, ok := result["account"].(string)
	if !ok {
		return errors.New("unexpected response format")
	}

	numStr := mirrorAccountId[strings.LastIndex(mirrorAccountId, ".")+1:]
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return err
	}
	id.Account = uint64(num)
	return nil
}

// PopulateEvmAddress gets the actual `AliasEvmAddress` field of the `AccountId` from the Mirror Node.
func (id *AccountID) PopulateEvmAddress(client *Client) error {
	result, err := id._MirrorNodeRequest(client, "evmAddress")
	if err != nil {
		return err
	}

	mirrorEvmAddress, ok := result["evm_address"].(string)
	if !ok {
		return errors.New("unexpected response format")
	}

	mirrorEvmAddress = strings.TrimPrefix(mirrorEvmAddress, "0x")
	asd, err := hex.DecodeString(mirrorEvmAddress)
	if err != nil {
		return err
	}
	id.AliasEvmAddress = &asd
	return nil
}

// Compare returns 0 if the two AccountID are identical, -1 if not.
func (id AccountID) Compare(given AccountID) int {
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

	if id.AliasKey != nil && given.AliasKey != nil {
		if id.AliasKey.String() > given.AliasKey.String() { //nolint
			return 1
		} else if id.AliasKey.String() < given.AliasKey.String() {
			return -1
		}
	}

	if id.AliasEvmAddress != nil && given.AliasEvmAddress != nil {
		originalEvmAddress := hex.EncodeToString(*id.AliasEvmAddress)
		givenEvmAddress := hex.EncodeToString(*given.AliasEvmAddress)
		if originalEvmAddress > givenEvmAddress { //nolint
			return 1
		} else if originalEvmAddress < givenEvmAddress {
			return -1
		}
	}

	if id.Account > given.Account { //nolint
		return 1
	} else if id.Account < given.Account {
		return -1
	} else {
		return 0
	}
}

// Len returns the number of elements in the collection.
func (accountIDs _AccountIDs) Len() int { //nolint
	return len(accountIDs.accountIDs)
}

func (accountIDs _AccountIDs) Swap(i, j int) { //nolint
	accountIDs.accountIDs[i], accountIDs.accountIDs[j] = accountIDs.accountIDs[j], accountIDs.accountIDs[i]
}

func (accountIDs _AccountIDs) Less(i, j int) bool { //nolint
	if accountIDs.accountIDs[i].Compare(accountIDs.accountIDs[j]) < 0 { //nolint
		return true
	}

	return false
}
