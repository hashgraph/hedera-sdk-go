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
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountCreateTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountCreateTransaction struct {
	Transaction
	proxyAccountID                *AccountID
	key                           Key
	initialBalance                uint64
	autoRenewPeriod               *time.Duration
	memo                          string
	receiverSignatureRequired     bool
	maxAutomaticTokenAssociations uint32
	stakedAccountID               *AccountID
	stakedNodeID                  *int64
	declineReward                 bool
	alias                         []byte
}

// NewAccountCreateTransaction creates an AccountCreateTransaction transaction which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	transaction := AccountCreateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _AccountCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetCryptoCreateAccount().GetKey())
	renew := _DurationFromProtobuf(pb.GetCryptoCreateAccount().GetAutoRenewPeriod())
	stakedNodeID := pb.GetCryptoCreateAccount().GetStakedNodeId()

	var stakeNodeAccountID *AccountID
	if pb.GetCryptoCreateAccount().GetStakedAccountId() != nil {
		stakeNodeAccountID = _AccountIDFromProtobuf(pb.GetCryptoCreateAccount().GetStakedAccountId())
	}

	body := AccountCreateTransaction{
		Transaction:                   transaction,
		key:                           key,
		initialBalance:                pb.GetCryptoCreateAccount().InitialBalance,
		autoRenewPeriod:               &renew,
		memo:                          pb.GetCryptoCreateAccount().GetMemo(),
		receiverSignatureRequired:     pb.GetCryptoCreateAccount().ReceiverSigRequired,
		maxAutomaticTokenAssociations: uint32(pb.GetCryptoCreateAccount().MaxAutomaticTokenAssociations),
		stakedAccountID:               stakeNodeAccountID,
		stakedNodeID:                  &stakedNodeID,
		declineReward:                 pb.GetCryptoCreateAccount().GetDeclineReward(),
	}

	if pb.GetCryptoCreateAccount().GetAlias() != nil {
		body.alias = pb.GetCryptoCreateAccount().GetAlias()
	}

	return &body
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *AccountCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountCreateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (transaction *AccountCreateTransaction) SetKey(key Key) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.key = key
	return transaction
}

// GetKey returns the key that must sign each transfer out of the account.
func (transaction *AccountCreateTransaction) GetKey() (Key, error) {
	return transaction.key, nil
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.initialBalance = uint64(initialBalance.AsTinybar())
	return transaction
}

// GetInitialBalance returns the initial number of Hbar to put into the account
func (transaction *AccountCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(int64(transaction.initialBalance))
}

// SetMaxAutomaticTokenAssociations
// Set the maximum number of tokens that an Account can be implicitly associated with. Defaults to 0
// and up to a maximum value of 1000.
func (transaction *AccountCreateTransaction) SetMaxAutomaticTokenAssociations(max uint32) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that an Account can be implicitly associated with.
func (transaction *AccountCreateTransaction) GetMaxAutomaticTokenAssociations() uint32 {
	return transaction.maxAutomaticTokenAssociations
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

// GetAutoRenewPeriod returns the time duration for when account is charged to extend its expiration date.
func (transaction *AccountCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a _Node, then this account is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was not set.
func (transaction *AccountCreateTransaction) SetProxyAccountID(id AccountID) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.proxyAccountID = &id
	return transaction
}

// Deprecated
func (transaction *AccountCreateTransaction) GetProxyAccountID() AccountID {
	if transaction.proxyAccountID == nil {
		return AccountID{}
	}

	return *transaction.proxyAccountID
}

// SetAccountMemo Sets the memo associated with the account (UTF-8 encoding max 100 bytes)
func (transaction *AccountCreateTransaction) SetAccountMemo(memo string) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

// GetAccountMemo Gets the memo associated with the account (UTF-8 encoding max 100 bytes)
func (transaction *AccountCreateTransaction) GetAccountMemo() string {
	return transaction.memo
}

// SetStakedAccountID Set the account to which this account will stake.
func (transaction *AccountCreateTransaction) SetStakedAccountID(id AccountID) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedAccountID = &id
	return transaction
}

// GetStakedAccountID returns the account to which this account will stake.
func (transaction *AccountCreateTransaction) GetStakedAccountID() AccountID {
	if transaction.stakedAccountID != nil {
		return *transaction.stakedAccountID
	}

	return AccountID{}
}

// SetStakedNodeID Set the node to which this account will stake
func (transaction *AccountCreateTransaction) SetStakedNodeID(id int64) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedNodeID = &id
	return transaction
}

// GetStakedNodeID returns the node to which this account will stake
func (transaction *AccountCreateTransaction) GetStakedNodeID() int64 {
	if transaction.stakedNodeID != nil {
		return *transaction.stakedNodeID
	}

	return 0
}

// SetDeclineStakingReward If set to true, the account declines receiving a staking reward. The default value is false.
func (transaction *AccountCreateTransaction) SetDeclineStakingReward(decline bool) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.declineReward = decline
	return transaction
}

// GetDeclineStakingReward returns true if the account declines receiving a staking reward.
func (transaction *AccountCreateTransaction) GetDeclineStakingReward() bool {
	return transaction.declineReward
}

func (transaction *AccountCreateTransaction) SetAlias(evmAddress string) *AccountCreateTransaction {
	transaction._RequireNotFrozen()

	evmAddress = strings.TrimPrefix(evmAddress, "0x")
	evmAddressBytes, _ := hex.DecodeString(evmAddress)

	transaction.alias = evmAddressBytes
	return transaction
}

func (transaction *AccountCreateTransaction) GetAlias() []byte {
	return transaction.alias
}

func (transaction *AccountCreateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.proxyAccountID != nil {
		if transaction.proxyAccountID != nil {
			if err := transaction.proxyAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (transaction *AccountCreateTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoCreateTransactionBody{
		InitialBalance:                transaction.initialBalance,
		ReceiverSigRequired:           transaction.receiverSignatureRequired,
		Memo:                          transaction.memo,
		MaxAutomaticTokenAssociations: int32(transaction.maxAutomaticTokenAssociations),
		DeclineReward:                 transaction.declineReward,
		Alias:                         transaction.alias,
	}

	if transaction.key != nil {
		body.Key = transaction.key._ToProtoKey()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	return &services.TransactionBody{
		TransactionID:            transaction.transactionID._ToProtobuf(),
		TransactionFee:           transaction.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		Memo:                     transaction.Transaction.memo,
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: body,
		},
	}
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (transaction *AccountCreateTransaction) SetReceiverSignatureRequired(required bool) *AccountCreateTransaction {
	transaction.receiverSignatureRequired = required
	return transaction
}

// GetReceiverSignatureRequired returns the receiverSigRequired flag.
func (transaction *AccountCreateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.receiverSignatureRequired
}

// Schedule a Create Account transaction
func (transaction *AccountCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoCreateTransactionBody{
		InitialBalance:                transaction.initialBalance,
		ReceiverSigRequired:           transaction.receiverSignatureRequired,
		Memo:                          transaction.memo,
		MaxAutomaticTokenAssociations: int32(transaction.maxAutomaticTokenAssociations),
		DeclineReward:                 transaction.declineReward,
	}

	if transaction.key != nil {
		body.Key = transaction.key._ToProtoKey()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	if transaction.alias != nil {
		body.Alias = transaction.alias
	}
	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: body,
		},
	}, nil
}

func _AccountCreateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CreateAccount,
	}
}

func (transaction *AccountCreateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountCreateTransaction) Sign(
	privateKey PrivateKey,
) *AccountCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *AccountCreateTransaction) SignWithOperator(
	client *Client,
) (*AccountCreateTransaction, error) {
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
func (transaction *AccountCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountCreateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountCreateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	if transaction.grpcDeadline == nil {
		transaction.grpcDeadline = client.requestTimeout
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountCreateTransactionGetMethod,
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

func (transaction *AccountCreateTransaction) Freeze() (*AccountCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountCreateTransaction) FreezeWith(client *Client) (*AccountCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	body := transaction._Build()
	if err != nil {
		return &AccountCreateTransaction{}, err
	}

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountCreateTransaction) SetMaxTransactionFee(fee Hbar) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountCreateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionMemo(memo string) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *AccountCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID returns the TransactionID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionID(transactionID TransactionID) *AccountCreateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *AccountCreateTransaction) SetMaxRetry(count int) *AccountCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *AccountCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountCreateTransaction {
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
func (transaction *AccountCreateTransaction) SetMaxBackoff(max time.Duration) *AccountCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *AccountCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *AccountCreateTransaction) SetMinBackoff(min time.Duration) *AccountCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *AccountCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountCreateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountCreateTransaction:%d", timestamp.UnixNano())
}

func (transaction *AccountCreateTransaction) SetLogLevel(level LogLevel) *AccountCreateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
