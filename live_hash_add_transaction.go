package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
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
	Transaction
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
	tx := LiveHashAddTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))
	return &tx
}

func _LiveHashAddTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *LiveHashAddTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.GetKeys())
	duration := _DurationFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.Duration)

	return &LiveHashAddTransaction{
		Transaction: tx,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoAddLiveHash().GetLiveHash().GetAccountId()),
		hash:        pb.GetCryptoAddLiveHash().LiveHash.Hash,
		keys:        &keys,
		duration:    &duration,
	}
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *LiveHashAddTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashAddTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *LiveHashAddTransaction) SignWithOperator(
	client *Client,
) (*LiveHashAddTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *LiveHashAddTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashAddTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *LiveHashAddTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashAddTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

func (tx *LiveHashAddTransaction) Freeze() (*LiveHashAddTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *LiveHashAddTransaction) FreezeWith(client *Client) (*LiveHashAddTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *LiveHashAddTransaction) GetMaxTransactionFee() Hbar {
	return tx.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *LiveHashAddTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *LiveHashAddTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this LiveHashAddTransaction.
func (tx *LiveHashAddTransaction) SetTransactionMemo(memo string) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this LiveHashAddTransaction.
func (tx *LiveHashAddTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashAddTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *LiveHashAddTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GetTransactionID gets the TransactionID for this	 LiveHashAddTransaction.
func (tx *LiveHashAddTransaction) GetTransactionID() TransactionID {
	return tx.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashAddTransaction.
func (tx *LiveHashAddTransaction) SetTransactionID(transactionID TransactionID) *LiveHashAddTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountID sets the _Node AccountID for this LiveHashAddTransaction.
func (tx *LiveHashAddTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashAddTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *LiveHashAddTransaction) SetMaxRetry(count int) *LiveHashAddTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *LiveHashAddTransaction) SetMaxBackoff(max time.Duration) *LiveHashAddTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *LiveHashAddTransaction) SetMinBackoff(min time.Duration) *LiveHashAddTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *LiveHashAddTransaction) SetLogLevel(level LogLevel) *LiveHashAddTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *LiveHashAddTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *LiveHashAddTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *LiveHashAddTransaction) getName() string {
	return "LiveHashAddTransaction"
}
func (tx *LiveHashAddTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *LiveHashAddTransaction) build() *services.TransactionBody {
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

func (tx *LiveHashAddTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

func (tx *LiveHashAddTransaction) buildProtoBody() *services.CryptoAddLiveHashTransactionBody {
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

func (tx *LiveHashAddTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().AddLiveHash,
	}
}
func (tx *LiveHashAddTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
