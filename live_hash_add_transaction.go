package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	transaction
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
	this := LiveHashAddTransaction{
		transaction: _NewTransaction(),
	}
	this._SetDefaultMaxTransactionFee(NewHbar(2))
	this.e = &this
	return &this
}

func _LiveHashAddTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *LiveHashAddTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.GetKeys())
	duration := _DurationFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.Duration)

	resultTx := &LiveHashAddTransaction{
		transaction: this,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoAddLiveHash().GetLiveHash().GetAccountId()),
		hash:        pb.GetCryptoAddLiveHash().LiveHash.Hash,
		keys:        &keys,
		duration:    &duration,
	}
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *LiveHashAddTransaction) SetGrpcDeadline(deadline *time.Duration) *LiveHashAddTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetHash Sets the SHA-384 hash of a credential or certificate
func (this *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.hash = hash
	return this
}

func (this *LiveHashAddTransaction) GetHash() []byte {
	return this.hash
}

// SetKeys Sets a list of keys (primitive or threshold), all of which must sign to attach the livehash to an account.
// Any one of which can later delete it.
func (this *LiveHashAddTransaction) SetKeys(keys ...Key) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	if this.keys == nil {
		this.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	this.keys = keyList

	return this
}

func (this *LiveHashAddTransaction) GetKeys() KeyList {
	if this.keys != nil {
		return *this.keys
	}

	return KeyList{}
}

// SetDuration Set the duration for which the livehash will remain valid
func (this *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.duration = &duration
	return this
}

// GetDuration returns the duration for which the livehash will remain valid
func (this *LiveHashAddTransaction) GetDuration() time.Duration {
	if this.duration != nil {
		return *this.duration
	}

	return time.Duration(0)
}

// SetAccountID Sets the account to which the livehash is attached
func (this *LiveHashAddTransaction) SetAccountID(accountID AccountID) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.accountID = &accountID
	return this
}

// GetAccountID returns the account to which the livehash is attached
func (this *LiveHashAddTransaction) GetAccountID() AccountID {
	if this.accountID == nil {
		return AccountID{}
	}

	return *this.accountID
}

// Sign uses the provided privateKey to sign the transaction.
func (this *LiveHashAddTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashAddTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *LiveHashAddTransaction) SignWithOperator(
	client *Client,
) (*LiveHashAddTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *LiveHashAddTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashAddTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *LiveHashAddTransaction) Freeze() (*LiveHashAddTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *LiveHashAddTransaction) FreezeWith(client *Client) (*LiveHashAddTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *LiveHashAddTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *LiveHashAddTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *LiveHashAddTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *LiveHashAddTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

func (this *LiveHashAddTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashAddTransaction.
func (this *LiveHashAddTransaction) SetTransactionMemo(memo string) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration sets the duration that this transaction is valid for.
// This is defaulted by the SDK to 120 seconds (or two minutes).
func (this *LiveHashAddTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashAddTransaction.
func (this *LiveHashAddTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	 LiveHashAddTransaction.
func (this *LiveHashAddTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashAddTransaction.
func (this *LiveHashAddTransaction) SetTransactionID(transactionID TransactionID) *LiveHashAddTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this LiveHashAddTransaction.
func (this *LiveHashAddTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashAddTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *LiveHashAddTransaction) SetMaxRetry(count int) *LiveHashAddTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *LiveHashAddTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashAddTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *LiveHashAddTransaction) SetMaxBackoff(max time.Duration) *LiveHashAddTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *LiveHashAddTransaction) SetMinBackoff(min time.Duration) *LiveHashAddTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *LiveHashAddTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("LiveHashAddTransaction:%d", timestamp.UnixNano())
}

func (this *LiveHashAddTransaction) SetLogLevel(level LogLevel) *LiveHashAddTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *LiveHashAddTransaction) getName() string {
	return "LiveHashAddTransaction"
}
func (this *LiveHashAddTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.accountID != nil {
		if err := this.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *LiveHashAddTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoAddLiveHash{
			CryptoAddLiveHash: this.buildProtoBody(),
		},
	}
}

func (this *LiveHashAddTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

func (this *LiveHashAddTransaction) buildProtoBody() *services.CryptoAddLiveHashTransactionBody {
	body := &services.CryptoAddLiveHashTransactionBody{
		LiveHash: &services.LiveHash{},
	}

	if this.accountID != nil {
		body.LiveHash.AccountId = this.accountID._ToProtobuf()
	}

	if this.duration != nil {
		body.LiveHash.Duration = _DurationToProtobuf(*this.duration)
	}

	if this.keys != nil {
		body.LiveHash.Keys = this.keys._ToProtoKeyList()
	}

	if this.hash != nil {
		body.LiveHash.Hash = this.hash
	}

	return body
}

func (this *LiveHashAddTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().AddLiveHash,
	}
}
