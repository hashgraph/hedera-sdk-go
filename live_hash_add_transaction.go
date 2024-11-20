package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"

	"time"
)

// LiveHashAddTransaction At consensus, attaches the given livehash to the given account.  The hash can be deleted by the
// key controlling the account, or by any of the keys associated to the livehash.  Hence livehashes
// provide a revocation service for their implied credentials; for example, when an authority grants
// a credential to the account, the account owner will cosign with the authority (or authorities) to
// attach a hash of the credential to the account---hence proving the grant. If the credential is
// revoked, then any of the authorities may delete it (or the account owner). In this way, the
// livehash mechanism acts as a revocation service.  An account cannot have two identical livehashes
// associated. To modify the list of keys in a livehash, the livehash should first be deleted, then
// recreated with a new list of keys.
type LiveHashAddTransaction struct {
	*Transaction[*LiveHashAddTransaction]
	accountID *AccountID
	hash      []byte
	keys      *KeyList
	duration  *time.Duration
}

// NewLiveHashAddTransaction creates LiveHashAddTransaction which at consensus, attaches the given livehash to the given account.
// The hash can be deleted by the key controlling the account, or by any of the keys associated to the livehash.  Hence livehashes
// provide a revocation service for their implied credentials; for example, when an authority grants
// a credential to the account, the account owner will cosign with the authority (or authorities) to
// attach a hash of the credential to the account---hence proving the grant. If the credential is
// revoked, then any of the authorities may delete it (or the account owner). In this way, the
// livehash mechanism acts as a revocation service.  An account cannot have two identical livehashes
// associated. To modify the list of keys in a livehash, the livehash should first be deleted, then
// recreated with a new list of keys.
func NewLiveHashAddTransaction() *LiveHashAddTransaction {
	tx := &LiveHashAddTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))
	return tx
}

func _LiveHashAddTransactionFromProtobuf(tx Transaction[*LiveHashAddTransaction], pb *services.TransactionBody) LiveHashAddTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.GetKeys())
	duration := _DurationFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.Duration)

	liveHashAddTransaction := LiveHashAddTransaction{
		accountID: _AccountIDFromProtobuf(pb.GetCryptoAddLiveHash().GetLiveHash().GetAccountId()),
		hash:      pb.GetCryptoAddLiveHash().LiveHash.Hash,
		keys:      &keys,
		duration:  &duration,
	}
	tx.childTransaction = &liveHashAddTransaction
	liveHashAddTransaction.Transaction = &tx
	return liveHashAddTransaction
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *LiveHashAddTransaction) SetGrpcDeadline(deadline *time.Duration) *LiveHashAddTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

// SetHash Sets the SHA-384 hash of a credential or certificate
func (tx *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.hash = hash
	return tx
}

func (tx *LiveHashAddTransaction) GetHash() []byte {
	return tx.hash
}

// SetKeys Sets a list of keys (primitive or threshold), all of which must sign to attach the livehash to an account.
// Any one of which can later delete it.
func (tx *LiveHashAddTransaction) SetKeys(keys ...Key) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	if tx.keys == nil {
		tx.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	tx.keys = keyList

	return tx
}

func (tx *LiveHashAddTransaction) GetKeys() KeyList {
	if tx.keys != nil {
		return *tx.keys
	}

	return KeyList{}
}

// SetDuration Set the duration for which the livehash will remain valid
func (tx *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.duration = &duration
	return tx
}

// GetDuration returns the duration for which the livehash will remain valid
func (tx *LiveHashAddTransaction) GetDuration() time.Duration {
	if tx.duration != nil {
		return *tx.duration
	}

	return time.Duration(0)
}

// SetAccountID Sets the account to which the livehash is attached
func (tx *LiveHashAddTransaction) SetAccountID(accountID AccountID) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the account to which the livehash is attached
func (tx *LiveHashAddTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ----------- Overridden functions ----------------

func (tx LiveHashAddTransaction) getName() string {
	return "LiveHashAddTransaction"
}
func (tx LiveHashAddTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx LiveHashAddTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoAddLiveHash{
			CryptoAddLiveHash: tx.buildProtoBody(),
		},
	}
}

func (tx LiveHashAddTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

func (tx LiveHashAddTransaction) buildProtoBody() *services.CryptoAddLiveHashTransactionBody {
	body := &services.CryptoAddLiveHashTransactionBody{
		LiveHash: &services.LiveHash{},
	}

	if tx.accountID != nil {
		body.LiveHash.AccountId = tx.accountID._ToProtobuf()
	}

	if tx.duration != nil {
		body.LiveHash.Duration = _DurationToProtobuf(*tx.duration)
	}

	if tx.keys != nil {
		body.LiveHash.Keys = tx.keys._ToProtoKeyList()
	}

	if tx.hash != nil {
		body.LiveHash.Hash = tx.hash
	}

	return body
}

func (tx LiveHashAddTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().AddLiveHash,
	}
}

func (tx LiveHashAddTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx LiveHashAddTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
