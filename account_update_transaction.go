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
	transaction
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
	this := AccountUpdateTransaction{
		transaction: _NewTransaction(),
	}
	this.e = &this
	this.SetAutoRenewPeriod(7890000 * time.Second)
	this._SetDefaultMaxTransactionFee(NewHbar(2))

	return &this
}

func _AccountUpdateTransactionFromProtobuf(transact transaction, pb *services.TransactionBody) *AccountUpdateTransaction {
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
		transaction:                   transact,
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
func (this *AccountUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountUpdateTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetKey Sets the new key for the Account
func (this *AccountUpdateTransaction) SetKey(key Key) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.key = key
	return this
}

func (this *AccountUpdateTransaction) GetKey() (Key, error) {
	return this.key, nil
}

// SetAccountID Sets the account ID which is being updated in this transaction.
func (this *AccountUpdateTransaction) SetAccountID(accountID AccountID) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.accountID = &accountID
	return this
}

func (this *AccountUpdateTransaction) GetAccountID() AccountID {
	if this.accountID == nil {
		return AccountID{}
	}

	return *this.accountID
}

// Deprecated
func (this *AccountUpdateTransaction) SetAliasKey(alias PublicKey) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.aliasKey = &alias
	return this
}

// Deprecated
func (this *AccountUpdateTransaction) GetAliasKey() PublicKey {
	if this.aliasKey == nil {
		return PublicKey{}
	}

	return *this.aliasKey
}

func (this *AccountUpdateTransaction) SetStakedAccountID(id AccountID) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.stakedAccountID = &id
	return this
}

func (this *AccountUpdateTransaction) GetStakedAccountID() AccountID {
	if this.stakedAccountID != nil {
		return *this.stakedAccountID
	}

	return AccountID{}
}

func (this *AccountUpdateTransaction) SetStakedNodeID(id int64) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.stakedNodeID = &id
	return this
}

func (this *AccountUpdateTransaction) GetStakedNodeID() int64 {
	if this.stakedNodeID != nil {
		return *this.stakedNodeID
	}

	return 0
}

func (this *AccountUpdateTransaction) SetDeclineStakingReward(decline bool) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.declineReward = decline
	return this
}

func (this *AccountUpdateTransaction) ClearStakedAccountID() *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.stakedAccountID = &AccountID{Account: 0}
	return this
}

func (this *AccountUpdateTransaction) ClearStakedNodeID() *AccountUpdateTransaction {
	this._RequireNotFrozen()
	*this.stakedNodeID = -1
	return this
}

func (this *AccountUpdateTransaction) GetDeclineStakingReward() bool {
	return this.declineReward
}

// SetMaxAutomaticTokenAssociations
// Sets the maximum number of tokens that an Account can be implicitly associated with. Up to a 1000
// including implicit and explicit associations.
func (this *AccountUpdateTransaction) SetMaxAutomaticTokenAssociations(max uint32) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.maxAutomaticTokenAssociations = max
	return this
}

func (this *AccountUpdateTransaction) GetMaxAutomaticTokenAssociations() uint32 {
	return this.maxAutomaticTokenAssociations
}

// SetReceiverSignatureRequired
// If true, this account's key must sign any transaction depositing into this account (in
// addition to all withdrawals)
func (this *AccountUpdateTransaction) SetReceiverSignatureRequired(receiverSignatureRequired bool) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.receiverSignatureRequired = receiverSignatureRequired
	return this
}

func (this *AccountUpdateTransaction) GetReceiverSignatureRequired() bool {
	return this.receiverSignatureRequired
}

// Deprecated
// SetProxyAccountID Sets the ID of the account to which this account is proxy staked.
func (this *AccountUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.proxyAccountID = &proxyAccountID
	return this
}

// Deprecated
func (this *AccountUpdateTransaction) GetProxyAccountID() AccountID {
	if this.proxyAccountID == nil {
		return AccountID{}
	}

	return *this.proxyAccountID
}

// SetAutoRenewPeriod Sets the duration in which it will automatically extend the expiration period.
func (this *AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.autoRenewPeriod = &autoRenewPeriod
	return this
}

func (this *AccountUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if this.autoRenewPeriod != nil {
		return *this.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime sets the new expiration time to extend to (ignored if equal to or before the current one)
func (this *AccountUpdateTransaction) SetExpirationTime(expirationTime time.Time) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.expirationTime = &expirationTime
	return this
}

func (this *AccountUpdateTransaction) GetExpirationTime() time.Time {
	if this.expirationTime != nil {
		return *this.expirationTime
	}
	return time.Time{}
}

// SetAccountMemo sets the new memo to be associated with the account (UTF-8 encoding max 100 bytes)
func (this *AccountUpdateTransaction) SetAccountMemo(memo string) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.memo = memo

	return this
}

func (this *AccountUpdateTransaction) GetAccountMemo() string {
	return this.memo
}

// Schedule Prepares a ScheduleCreateTransaction containing this transaction.
func (this *AccountUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (this *AccountUpdateTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *AccountUpdateTransaction) Sign(
	privateKey PrivateKey,
) *AccountUpdateTransaction {
	this.transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *AccountUpdateTransaction) SignWithOperator(
	client *Client,
) (*AccountUpdateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *AccountUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountUpdateTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *AccountUpdateTransaction) Freeze() (*AccountUpdateTransaction, error) {
	_,err := this.transaction.Freeze()
	return this, err
}

func (this *AccountUpdateTransaction) FreezeWith(client *Client) (*AccountUpdateTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountUpdateTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountUpdateTransaction) SetMaxTransactionFee(fee Hbar) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *AccountUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *AccountUpdateTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountUpdateTransaction.
func (this *AccountUpdateTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountUpdateTransaction.
func (this *AccountUpdateTransaction) SetTransactionMemo(memo string) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *AccountUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountUpdateTransaction.
func (this *AccountUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this AccountUpdateTransaction.
func (this *AccountUpdateTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountUpdateTransaction.
func (this *AccountUpdateTransaction) SetTransactionID(transactionID TransactionID) *AccountUpdateTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountUpdateTransaction.
func (this *AccountUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountUpdateTransaction) SetMaxRetry(count int) *AccountUpdateTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *AccountUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountUpdateTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountUpdateTransaction) SetMaxBackoff(max time.Duration) *AccountUpdateTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *AccountUpdateTransaction) SetMinBackoff(min time.Duration) *AccountUpdateTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (transaction *AccountUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountUpdateTransaction:%d", timestamp.UnixNano())
}

func (transaction *AccountUpdateTransaction) SetLogLevel(level LogLevel) *AccountUpdateTransaction {
	transaction.transaction.SetLogLevel(level)
	return transaction
}

// ----------- overriden functions ----------------

func (this *AccountUpdateTransaction) getName() string {
	return "AccountUpdateTransaction"
}

func (this *AccountUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.accountID != nil {
		if err := this.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if this.proxyAccountID != nil {
		if err := this.proxyAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *AccountUpdateTransaction) build() *services.TransactionBody {
	body := &services.CryptoUpdateTransactionBody{
		ReceiverSigRequiredField: &services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: this.receiverSignatureRequired},
		},
		Memo:                          &wrapperspb.StringValue{Value: this.memo},
		MaxAutomaticTokenAssociations: &wrapperspb.Int32Value{Value: int32(this.maxAutomaticTokenAssociations)},
		DeclineReward:                 &wrapperspb.BoolValue{Value: this.declineReward},
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*this.expirationTime)
	}

	if this.accountID != nil {
		body.AccountIDToUpdate = this.accountID._ToProtobuf()
	}

	if this.key != nil {
		body.Key = this.key._ToProtoKey()
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	pb := services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}

	body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: int32(this.maxAutomaticTokenAssociations)}

	return &pb
}
func (this *AccountUpdateTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoUpdateTransactionBody{
		ReceiverSigRequiredField: &services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: this.receiverSignatureRequired},
		},
		Memo:                          &wrapperspb.StringValue{Value: this.memo},
		DeclineReward:                 &wrapperspb.BoolValue{Value: this.declineReward},
		MaxAutomaticTokenAssociations: &wrapperspb.Int32Value{Value: int32(this.maxAutomaticTokenAssociations)},
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*this.expirationTime)
	}

	if this.accountID != nil {
		body.AccountIDToUpdate = this.accountID._ToProtobuf()
	}

	if this.key != nil {
		body.Key = this.key._ToProtoKey()
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}, nil
}

func (this *AccountUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().UpdateAccount,
	}
}