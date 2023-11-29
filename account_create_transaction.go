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
	transaction
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
// execute a Crypto Create transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	this := AccountCreateTransaction{
		transaction: _NewTransaction(),
	}

	this.e = &this
	this.SetAutoRenewPeriod(7890000 * time.Second)
	this._SetDefaultMaxTransactionFee(NewHbar(5))

	return &this
}

func _AccountCreateTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *AccountCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetCryptoCreateAccount().GetKey())
	renew := _DurationFromProtobuf(pb.GetCryptoCreateAccount().GetAutoRenewPeriod())
	stakedNodeID := pb.GetCryptoCreateAccount().GetStakedNodeId()

	var stakeNodeAccountID *AccountID
	if pb.GetCryptoCreateAccount().GetStakedAccountId() != nil {
		stakeNodeAccountID = _AccountIDFromProtobuf(pb.GetCryptoCreateAccount().GetStakedAccountId())
	}

	body := AccountCreateTransaction{
		transaction:                   this,
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
func (this *AccountCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountCreateTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (this *AccountCreateTransaction) SetKey(key Key) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.key = key
	return this
}

// GetKey returns the key that must sign each transfer out of the account.
func (this *AccountCreateTransaction) GetKey() (Key, error) {
	return this.key, nil
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (this *AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.initialBalance = uint64(initialBalance.AsTinybar())
	return this
}

// GetInitialBalance returns the initial number of Hbar to put into the account
func (this *AccountCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(int64(this.initialBalance))
}

// SetMaxAutomaticTokenAssociations
// Set the maximum number of tokens that an Account can be implicitly associated with. Defaults to 0
// and up to a maximum value of 1000.
func (this *AccountCreateTransaction) SetMaxAutomaticTokenAssociations(max uint32) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.maxAutomaticTokenAssociations = max
	return this
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that an Account can be implicitly associated with.
func (this *AccountCreateTransaction) GetMaxAutomaticTokenAssociations() uint32 {
	return this.maxAutomaticTokenAssociations
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (this *AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.autoRenewPeriod = &autoRenewPeriod
	return this
}

// GetAutoRenewPeriod returns the time duration for when account is charged to extend its expiration date.
func (this *AccountCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if this.autoRenewPeriod != nil {
		return *this.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a _Node, then this account is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was not set.
func (this *AccountCreateTransaction) SetProxyAccountID(id AccountID) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.proxyAccountID = &id
	return this
}

// Deprecated
func (this *AccountCreateTransaction) GetProxyAccountID() AccountID {
	if this.proxyAccountID == nil {
		return AccountID{}
	}

	return *this.proxyAccountID
}

// SetAccountMemo Sets the memo associated with the account (UTF-8 encoding max 100 bytes)
func (this *AccountCreateTransaction) SetAccountMemo(memo string) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.memo = memo
	return this
}

// GetAccountMemo Gets the memo associated with the account (UTF-8 encoding max 100 bytes)
func (this *AccountCreateTransaction) GetAccountMemo() string {
	return this.memo
}

// SetStakedAccountID Set the account to which this account will stake.
func (this *AccountCreateTransaction) SetStakedAccountID(id AccountID) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.stakedAccountID = &id
	return this
}

// GetStakedAccountID returns the account to which this account will stake.
func (this *AccountCreateTransaction) GetStakedAccountID() AccountID {
	if this.stakedAccountID != nil {
		return *this.stakedAccountID
	}

	return AccountID{}
}

// SetStakedNodeID Set the node to which this account will stake
func (this *AccountCreateTransaction) SetStakedNodeID(id int64) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.stakedNodeID = &id
	return this
}

// GetStakedNodeID returns the node to which this account will stake
func (this *AccountCreateTransaction) GetStakedNodeID() int64 {
	if this.stakedNodeID != nil {
		return *this.stakedNodeID
	}

	return 0
}

// SetDeclineStakingReward If set to true, the account declines receiving a staking reward. The default value is false.
func (this *AccountCreateTransaction) SetDeclineStakingReward(decline bool) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.declineReward = decline
	return this
}

// GetDeclineStakingReward returns true if the account declines receiving a staking reward.
func (this *AccountCreateTransaction) GetDeclineStakingReward() bool {
	return this.declineReward
}

func (this *AccountCreateTransaction) SetAlias(evmAddress string) *AccountCreateTransaction {
	this._RequireNotFrozen()

	evmAddress = strings.TrimPrefix(evmAddress, "0x")
	evmAddressBytes, _ := hex.DecodeString(evmAddress)

	this.alias = evmAddressBytes
	return this
}

func (this *AccountCreateTransaction) GetAlias() []byte {
	return this.alias
}



// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (this *AccountCreateTransaction) SetReceiverSignatureRequired(required bool) *AccountCreateTransaction {
	this.receiverSignatureRequired = required
	return this
}

// GetReceiverSignatureRequired returns the receiverSigRequired flag.
func (this *AccountCreateTransaction) GetReceiverSignatureRequired() bool {
	return this.receiverSignatureRequired
}

// Schedule a Create Account transaction
func (this *AccountCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildScheduled()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (this *AccountCreateTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *AccountCreateTransaction) Sign(
	privateKey PrivateKey,
) *AccountCreateTransaction {
	this.transaction.Sign(privateKey);
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *AccountCreateTransaction) SignWithOperator(
	client *Client,
) (*AccountCreateTransaction, error) {
	_,err := this.transaction.SignWithOperator(client)
	return this,err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *AccountCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountCreateTransaction {
	this.transaction.SignWith(publicKey,signer);
	return this
}


func (this *AccountCreateTransaction) Freeze() (*AccountCreateTransaction, error) {
    _,err := this.transaction.Freeze()
	return this,err
}

func (this *AccountCreateTransaction) FreezeWith(client *Client) (*AccountCreateTransaction, error) {
	_,err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountCreateTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountCreateTransaction) SetMaxTransactionFee(fee Hbar) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *AccountCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *AccountCreateTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountCreateTransaction.
func (this *AccountCreateTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountCreateTransaction.
func (this *AccountCreateTransaction) SetTransactionMemo(memo string) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *AccountCreateTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountCreateTransaction.
func (this *AccountCreateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID returns the TransactionID for this AccountCreateTransaction.
func (this *AccountCreateTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountCreateTransaction.
func (this *AccountCreateTransaction) SetTransactionID(transactionID TransactionID) *AccountCreateTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountCreateTransaction.
func (this *AccountCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountCreateTransaction {
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountCreateTransaction) SetMaxRetry(count int) *AccountCreateTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *AccountCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountCreateTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountCreateTransaction) SetMaxBackoff(max time.Duration) *AccountCreateTransaction {
	this.transaction.SetMaxBackoff(max)
    return this
}

func (this *AccountCreateTransaction) GetMaxBackoff() time.Duration {
	return this.transaction.GetMaxBackoff();
}
	


// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *AccountCreateTransaction) SetMinBackoff(min time.Duration) *AccountCreateTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}


func (this *AccountCreateTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountCreateTransaction:%d", timestamp.UnixNano())
}

func (this *AccountCreateTransaction) SetLogLevel(level LogLevel) *AccountCreateTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (transaction *AccountCreateTransaction) getName() string {
	return "AccountCreateTransaction"
}

func (this *AccountCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.proxyAccountID != nil {
		if this.proxyAccountID != nil {
			if err := this.proxyAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *AccountCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionID:            this.transactionID._ToProtobuf(),
		TransactionFee:           this.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		Memo:                     this.transaction.memo,
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: this.buildProtoBody(),
		},
	}
}
func (this *AccountCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: this.buildProtoBody(),
		},
	}, nil
}

func (this *AccountCreateTransaction) buildProtoBody() *services.CryptoCreateTransactionBody {
	body := &services.CryptoCreateTransactionBody{
		InitialBalance:                this.initialBalance,
		ReceiverSigRequired:           this.receiverSignatureRequired,
		Memo:                          this.memo,
		MaxAutomaticTokenAssociations: int32(this.maxAutomaticTokenAssociations),
		DeclineReward:                 this.declineReward,
		Alias:                         this.alias,
	}

	if this.key != nil {
		body.Key = this.key._ToProtoKey()
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	return body
}

func (this *AccountCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CreateAccount,
	}
}
