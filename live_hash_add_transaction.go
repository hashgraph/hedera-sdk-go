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
	transaction := LiveHashAddTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _LiveHashAddTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *LiveHashAddTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.GetKeys())
	duration := _DurationFromProtobuf(pb.GetCryptoAddLiveHash().LiveHash.Duration)

	return &LiveHashAddTransaction{
		Transaction: transaction,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoAddLiveHash().GetLiveHash().GetAccountId()),
		hash:        pb.GetCryptoAddLiveHash().LiveHash.Hash,
		keys:        &keys,
		duration:    &duration,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *LiveHashAddTransaction) SetGrpcDeadline(deadline *time.Duration) *LiveHashAddTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetHash Sets the SHA-384 hash of a credential or certificate
func (transaction *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.hash = hash
	return transaction
}

func (transaction *LiveHashAddTransaction) GetHash() []byte {
	return transaction.hash
}

// SetKeys Sets a list of keys (primitive or threshold), all of which must sign to attach the livehash to an account.
// Any one of which can later delete it.
func (transaction *LiveHashAddTransaction) SetKeys(keys ...Key) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	if transaction.keys == nil {
		transaction.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.keys = keyList

	return transaction
}

func (transaction *LiveHashAddTransaction) GetKeys() KeyList {
	if transaction.keys != nil {
		return *transaction.keys
	}

	return KeyList{}
}

// SetDuration Set the duration for which the livehash will remain valid
func (transaction *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.duration = &duration
	return transaction
}

// GetDuration returns the duration for which the livehash will remain valid
func (transaction *LiveHashAddTransaction) GetDuration() time.Duration {
	if transaction.duration != nil {
		return *transaction.duration
	}

	return time.Duration(0)
}

// SetAccountID Sets the account to which the livehash is attached
func (transaction *LiveHashAddTransaction) SetAccountID(accountID AccountID) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

// GetAccountID returns the account to which the livehash is attached
func (transaction *LiveHashAddTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

func (transaction *LiveHashAddTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *LiveHashAddTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoAddLiveHashTransactionBody{
		LiveHash: &services.LiveHash{},
	}

	if transaction.accountID != nil {
		body.LiveHash.AccountId = transaction.accountID._ToProtobuf()
	}

	if transaction.duration != nil {
		body.LiveHash.Duration = _DurationToProtobuf(*transaction.duration)
	}

	if transaction.keys != nil {
		body.LiveHash.Keys = transaction.keys._ToProtoKeyList()
	}

	if transaction.hash != nil {
		body.LiveHash.Hash = transaction.hash
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoAddLiveHash{
			CryptoAddLiveHash: body,
		},
	}
}

func (transaction *LiveHashAddTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

func _LiveHashAddTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().AddLiveHash,
	}
}

func (transaction *LiveHashAddTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *LiveHashAddTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashAddTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *LiveHashAddTransaction) SignWithOperator(
	client *Client,
) (*LiveHashAddTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *LiveHashAddTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashAddTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashAddTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_LiveHashAddTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
		transaction.maxBackoff,
		transaction.minBackoff,
		transaction.maxRetry,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID:  transaction.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *LiveHashAddTransaction) Freeze() (*LiveHashAddTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashAddTransaction) FreezeWith(client *Client) (*LiveHashAddTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &LiveHashAddTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *LiveHashAddTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *LiveHashAddTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *LiveHashAddTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *LiveHashAddTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *LiveHashAddTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionMemo(memo string) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration sets the duration that this transaction is valid for.
// This is defaulted by the SDK to 120 seconds (or two minutes).
func (transaction *LiveHashAddTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	 LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionID(transactionID TransactionID) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashAddTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *LiveHashAddTransaction) SetMaxRetry(count int) *LiveHashAddTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *LiveHashAddTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashAddTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *LiveHashAddTransaction) SetMaxBackoff(max time.Duration) *LiveHashAddTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *LiveHashAddTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *LiveHashAddTransaction) SetMinBackoff(min time.Duration) *LiveHashAddTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *LiveHashAddTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *LiveHashAddTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("LiveHashAddTransaction:%d", timestamp.UnixNano())
}

func (transaction *LiveHashAddTransaction) SetLogLevel(level LogLevel) *LiveHashAddTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
