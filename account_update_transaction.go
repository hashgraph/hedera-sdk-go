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
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountUpdateTransaction
// Change properties for the given account. Any null field is ignored (left unchanged). This
// transaction must be signed by the existing key for this account. If the transaction is changing
// the key field, then the transaction must be signed by both the old key (from before the change)
// and the new key. The old key must sign for security. The new key must sign as a safeguard to
// avoid accidentally changing to an invalid key, and then having no way to recover.
type AccountUpdateTransaction struct {
	Transaction
	accountID                     *AccountID
	proxyAccountID                *AccountID
	key                           Key
	autoRenewPeriod               *time.Duration
	memo                          string
	receiverSignatureRequired     bool
	expirationTime                *time.Time
	maxAutomaticTokenAssociations uint32
	aliasKey                      *PublicKey
	stakedAccountID               *AccountID
	stakedNodeID                  *int64
	declineReward                 bool
}

// NewAccountUpdateTransaction
// Creates AccoutnUppdateTransaction which changes properties for the given account.
// Any null field is ignored (left unchanged).
// This transaction must be signed by the existing key for this account. If the transaction is changing
// the key field, then the transaction must be signed by both the old key (from before the change)
// and the new key. The old key must sign for security. The new key must sign as a safeguard to
// avoid accidentally changing to an invalid key, and then having no way to recover.
func NewAccountUpdateTransaction() *AccountUpdateTransaction {
	transaction := AccountUpdateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _AccountUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountUpdateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetCryptoUpdateAccount().GetKey())
	var receiverSignatureRequired bool

	switch s := pb.GetCryptoUpdateAccount().GetReceiverSigRequiredField().(type) {
	case *services.CryptoUpdateTransactionBody_ReceiverSigRequired:
		receiverSignatureRequired = s.ReceiverSigRequired // nolint
	case *services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper:
		receiverSignatureRequired = s.ReceiverSigRequiredWrapper.Value // nolint
	}

	autoRenew := _DurationFromProtobuf(pb.GetCryptoUpdateAccount().AutoRenewPeriod)
	expiration := _TimeFromProtobuf(pb.GetCryptoUpdateAccount().ExpirationTime)

	stakedNodeID := pb.GetCryptoUpdateAccount().GetStakedNodeId()

	var stakeNodeAccountID *AccountID
	if pb.GetCryptoUpdateAccount().GetStakedAccountId() != nil {
		stakeNodeAccountID = _AccountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetStakedAccountId())
	}

	return &AccountUpdateTransaction{
		Transaction:                   transaction,
		accountID:                     _AccountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetAccountIDToUpdate()),
		key:                           key,
		autoRenewPeriod:               &autoRenew,
		memo:                          pb.GetCryptoUpdateAccount().GetMemo().Value,
		receiverSignatureRequired:     receiverSignatureRequired,
		expirationTime:                &expiration,
		maxAutomaticTokenAssociations: uint32(pb.GetCryptoUpdateAccount().MaxAutomaticTokenAssociations.GetValue()),
		stakedAccountID:               stakeNodeAccountID,
		stakedNodeID:                  &stakedNodeID,
		declineReward:                 pb.GetCryptoUpdateAccount().GetDeclineReward().GetValue(),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *AccountUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountUpdateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetKey Sets the new key for the Account
func (transaction *AccountUpdateTransaction) SetKey(key Key) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.key = key
	return transaction
}

func (transaction *AccountUpdateTransaction) GetKey() (Key, error) {
	return transaction.key, nil
}

// SetAccountID Sets the account ID which is being updated in this transaction.
func (transaction *AccountUpdateTransaction) SetAccountID(accountID AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

// Deprecated
func (transaction *AccountUpdateTransaction) SetAliasKey(alias PublicKey) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.aliasKey = &alias
	return transaction
}

// Deprecated
func (transaction *AccountUpdateTransaction) GetAliasKey() PublicKey {
	if transaction.aliasKey == nil {
		return PublicKey{}
	}

	return *transaction.aliasKey
}

func (transaction *AccountUpdateTransaction) SetStakedAccountID(id AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedAccountID = &id
	return transaction
}

func (transaction *AccountUpdateTransaction) GetStakedAccountID() AccountID {
	if transaction.stakedAccountID != nil {
		return *transaction.stakedAccountID
	}

	return AccountID{}
}

func (transaction *AccountUpdateTransaction) SetStakedNodeID(id int64) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedNodeID = &id
	return transaction
}

func (transaction *AccountUpdateTransaction) GetStakedNodeID() int64 {
	if transaction.stakedNodeID != nil {
		return *transaction.stakedNodeID
	}

	return 0
}

func (transaction *AccountUpdateTransaction) SetDeclineStakingReward(decline bool) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.declineReward = decline
	return transaction
}

func (transaction *AccountUpdateTransaction) ClearStakedAccountID() *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedAccountID = &AccountID{Account: 0}
	return transaction
}

func (transaction *AccountUpdateTransaction) ClearStakedNodeID() *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	*transaction.stakedNodeID = -1
	return transaction
}

func (transaction *AccountUpdateTransaction) GetDeclineStakingReward() bool {
	return transaction.declineReward
}

// SetMaxAutomaticTokenAssociations
// Sets the maximum number of tokens that an Account can be implicitly associated with. Up to a 1000
// including implicit and explicit associations.
func (transaction *AccountUpdateTransaction) SetMaxAutomaticTokenAssociations(max uint32) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

func (transaction *AccountUpdateTransaction) GetMaxAutomaticTokenAssociations() uint32 {
	return transaction.maxAutomaticTokenAssociations
}

// SetReceiverSignatureRequired
// If true, this account's key must sign any transaction depositing into this account (in
// addition to all withdrawals)
func (transaction *AccountUpdateTransaction) SetReceiverSignatureRequired(receiverSignatureRequired bool) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.receiverSignatureRequired = receiverSignatureRequired
	return transaction
}

func (transaction *AccountUpdateTransaction) GetReceiverSignatureRequired() bool {
	return transaction.receiverSignatureRequired
}

// Deprecated
// SetProxyAccountID Sets the ID of the account to which this account is proxy staked.
func (transaction *AccountUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.proxyAccountID = &proxyAccountID
	return transaction
}

// Deprecated
func (transaction *AccountUpdateTransaction) GetProxyAccountID() AccountID {
	if transaction.proxyAccountID == nil {
		return AccountID{}
	}

	return *transaction.proxyAccountID
}

// SetAutoRenewPeriod Sets the duration in which it will automatically extend the expiration period.
func (transaction *AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *AccountUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime sets the new expiration time to extend to (ignored if equal to or before the current one)
func (transaction *AccountUpdateTransaction) SetExpirationTime(expirationTime time.Time) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expirationTime
	return transaction
}

func (transaction *AccountUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}
	return time.Time{}
}

// SetAccountMemo sets the new memo to be associated with the account (UTF-8 encoding max 100 bytes)
func (transaction *AccountUpdateTransaction) SetAccountMemo(memo string) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *AccountUpdateTransaction) GetAccountMemo() string {
	return transaction.memo
}

func (transaction *AccountUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.proxyAccountID != nil {
		if err := transaction.proxyAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *AccountUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoUpdateTransactionBody{
		ReceiverSigRequiredField: &services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: transaction.receiverSignatureRequired},
		},
		Memo:                          &wrapperspb.StringValue{Value: transaction.memo},
		MaxAutomaticTokenAssociations: &wrapperspb.Int32Value{Value: int32(transaction.maxAutomaticTokenAssociations)},
		DeclineReward:                 &wrapperspb.BoolValue{Value: transaction.declineReward},
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.accountID != nil {
		body.AccountIDToUpdate = transaction.accountID._ToProtobuf()
	}

	if transaction.key != nil {
		body.Key = transaction.key._ToProtoKey()
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	pb := services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}

	body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: int32(transaction.maxAutomaticTokenAssociations)}

	return &pb
}

// Schedule Prepares a ScheduleCreateTransaction containing this transaction.
func (transaction *AccountUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoUpdateTransactionBody{
		ReceiverSigRequiredField: &services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: transaction.receiverSignatureRequired},
		},
		Memo:                          &wrapperspb.StringValue{Value: transaction.memo},
		DeclineReward:                 &wrapperspb.BoolValue{Value: transaction.declineReward},
		MaxAutomaticTokenAssociations: &wrapperspb.Int32Value{Value: int32(transaction.maxAutomaticTokenAssociations)},
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.accountID != nil {
		body.AccountIDToUpdate = transaction.accountID._ToProtobuf()
	}

	if transaction.key != nil {
		body.Key = transaction.key._ToProtoKey()
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}, nil
}

func _AccountUpdateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().UpdateAccount,
	}
}

func (transaction *AccountUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountUpdateTransaction) Sign(
	privateKey PrivateKey,
) *AccountUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *AccountUpdateTransaction) SignWithOperator(
	client *Client,
) (*AccountUpdateTransaction, error) {
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
func (transaction *AccountUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountUpdateTransaction) Execute(
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

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountUpdateTransactionGetMethod,
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

func (transaction *AccountUpdateTransaction) Freeze() (*AccountUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountUpdateTransaction) FreezeWith(client *Client) (*AccountUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &AccountUpdateTransaction{}, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountUpdateTransaction) SetMaxTransactionFee(fee Hbar) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountUpdateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionMemo(memo string) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *AccountUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetTransactionID(transactionID TransactionID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountUpdateTransaction.
func (transaction *AccountUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *AccountUpdateTransaction) SetMaxRetry(count int) *AccountUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *AccountUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountUpdateTransaction {
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
func (transaction *AccountUpdateTransaction) SetMaxBackoff(max time.Duration) *AccountUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *AccountUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *AccountUpdateTransaction) SetMinBackoff(min time.Duration) *AccountUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *AccountUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountUpdateTransaction:%d", timestamp.UnixNano())
}

func (transaction *AccountUpdateTransaction) SetLogLevel(level LogLevel) *AccountUpdateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
