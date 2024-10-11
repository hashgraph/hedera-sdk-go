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
	"encoding/hex"
	"strings"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"
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
	maxAutomaticTokenAssociations int32
	stakedAccountID               *AccountID
	stakedNodeID                  *int64
	declineReward                 bool
	alias                         []byte
}

// NewAccountCreateTransaction creates an AccountCreateTransaction transaction which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	tx := AccountCreateTransaction{
		Transaction: _NewTransaction(),
	}

	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _AccountCreateTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *AccountCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetCryptoCreateAccount().GetKey())
	renew := _DurationFromProtobuf(pb.GetCryptoCreateAccount().GetAutoRenewPeriod())

	var stakedNodeID *int64
	if pb.GetCryptoCreateAccount().GetStakedNodeId() != 0 {
		nodeId := pb.GetCryptoCreateAccount().GetStakedNodeId()
		stakedNodeID = &nodeId
	}
	var stakeAccountID *AccountID
	if pb.GetCryptoCreateAccount().GetStakedAccountId() != nil {
		stakeAccountID = _AccountIDFromProtobuf(pb.GetCryptoCreateAccount().GetStakedAccountId())
	}

	body := AccountCreateTransaction{
		Transaction:                   tx,
		key:                           key,
		initialBalance:                pb.GetCryptoCreateAccount().InitialBalance,
		autoRenewPeriod:               &renew,
		memo:                          pb.GetCryptoCreateAccount().GetMemo(),
		receiverSignatureRequired:     pb.GetCryptoCreateAccount().ReceiverSigRequired,
		maxAutomaticTokenAssociations: pb.GetCryptoCreateAccount().MaxAutomaticTokenAssociations,
		stakedAccountID:               stakeAccountID,
		stakedNodeID:                  stakedNodeID,
		declineReward:                 pb.GetCryptoCreateAccount().GetDeclineReward(),
	}

	if pb.GetCryptoCreateAccount().GetAlias() != nil {
		body.alias = pb.GetCryptoCreateAccount().GetAlias()
	}

	return &body
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (tx *AccountCreateTransaction) SetKey(key Key) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.key = key
	return tx
}

// GetKey returns the key that must sign each transfer out of the account.
func (tx *AccountCreateTransaction) GetKey() (Key, error) {
	return tx.key, nil
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (tx *AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.initialBalance = uint64(initialBalance.AsTinybar())
	return tx
}

// GetInitialBalance returns the initial number of Hbar to put into the account
func (tx *AccountCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(int64(tx.initialBalance))
}

// SetMaxAutomaticTokenAssociations
// Set the maximum number of tokens that an Account can be implicitly associated with. Defaults to 0
// and up to a maximum value of 1000.
func (tx *AccountCreateTransaction) SetMaxAutomaticTokenAssociations(max int32) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.maxAutomaticTokenAssociations = max
	return tx
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that an Account can be implicitly associated with.
func (tx *AccountCreateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return tx.maxAutomaticTokenAssociations
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (tx *AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &autoRenewPeriod
	return tx
}

// GetAutoRenewPeriod returns the time duration for when account is charged to extend its expiration date.
func (tx *AccountCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a _Node, then this account is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was not set.
func (tx *AccountCreateTransaction) SetProxyAccountID(id AccountID) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.proxyAccountID = &id
	return tx
}

// Deprecated
func (tx *AccountCreateTransaction) GetProxyAccountID() AccountID {
	if tx.proxyAccountID == nil {
		return AccountID{}
	}

	return *tx.proxyAccountID
}

// SetAccountMemo Sets the memo associated with the account (UTF-8 encoding max 100 bytes)
func (tx *AccountCreateTransaction) SetAccountMemo(memo string) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetAccountMemo Gets the memo associated with the account (UTF-8 encoding max 100 bytes)
func (tx *AccountCreateTransaction) GetAccountMemo() string {
	return tx.memo
}

// SetStakedAccountID Set the account to which this account will stake.
func (tx *AccountCreateTransaction) SetStakedAccountID(id AccountID) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.stakedAccountID = &id
	return tx
}

// GetStakedAccountID returns the account to which this account will stake.
func (tx *AccountCreateTransaction) GetStakedAccountID() AccountID {
	if tx.stakedAccountID != nil {
		return *tx.stakedAccountID
	}

	return AccountID{}
}

// SetStakedNodeID Set the node to which this account will stake
func (tx *AccountCreateTransaction) SetStakedNodeID(id int64) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.stakedNodeID = &id
	return tx
}

// GetStakedNodeID returns the node to which this account will stake
func (tx *AccountCreateTransaction) GetStakedNodeID() int64 {
	if tx.stakedNodeID != nil {
		return *tx.stakedNodeID
	}

	return 0
}

// SetDeclineStakingReward If set to true, the account declines receiving a staking reward. The default value is false.
func (tx *AccountCreateTransaction) SetDeclineStakingReward(decline bool) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.declineReward = decline
	return tx
}

// GetDeclineStakingReward returns true if the account declines receiving a staking reward.
func (tx *AccountCreateTransaction) GetDeclineStakingReward() bool {
	return tx.declineReward
}

func (tx *AccountCreateTransaction) SetAlias(evmAddress string) *AccountCreateTransaction {
	tx._RequireNotFrozen()

	evmAddress = strings.TrimPrefix(evmAddress, "0x")
	evmAddressBytes, _ := hex.DecodeString(evmAddress)

	tx.alias = evmAddressBytes
	return tx
}

func (tx *AccountCreateTransaction) GetAlias() []byte {
	return tx.alias
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (tx *AccountCreateTransaction) SetReceiverSignatureRequired(required bool) *AccountCreateTransaction {
	tx.receiverSignatureRequired = required
	return tx
}

// GetReceiverSignatureRequired returns the receiverSigRequired flag.
func (tx *AccountCreateTransaction) GetReceiverSignatureRequired() bool {
	return tx.receiverSignatureRequired
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *AccountCreateTransaction) Sign(
	privateKey PrivateKey,
) *AccountCreateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *AccountCreateTransaction) SignWithOperator(
	client *Client,
) (*AccountCreateTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *AccountCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountCreateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *AccountCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountCreateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *AccountCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountCreateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *AccountCreateTransaction) Freeze() (*AccountCreateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *AccountCreateTransaction) FreezeWith(client *Client) (*AccountCreateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *AccountCreateTransaction) GetMaxTransactionFee() Hbar {
	return tx.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *AccountCreateTransaction) SetMaxTransactionFee(fee Hbar) *AccountCreateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *AccountCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountCreateTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this AccountCreateTransaction.
func (tx *AccountCreateTransaction) SetTransactionMemo(memo string) *AccountCreateTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this AccountCreateTransaction.
func (tx *AccountCreateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountCreateTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *AccountCreateTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this AccountCreateTransaction.
func (tx *AccountCreateTransaction) SetTransactionID(transactionID TransactionID) *AccountCreateTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountCreateTransaction.
func (tx *AccountCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountCreateTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *AccountCreateTransaction) SetMaxRetry(count int) *AccountCreateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *AccountCreateTransaction) SetMaxBackoff(max time.Duration) *AccountCreateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *AccountCreateTransaction) SetMinBackoff(min time.Duration) *AccountCreateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *AccountCreateTransaction) SetLogLevel(level LogLevel) *AccountCreateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *AccountCreateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *AccountCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *AccountCreateTransaction) getName() string {
	return "AccountCreateTransaction"
}

func (tx *AccountCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.proxyAccountID != nil {
		if tx.proxyAccountID != nil {
			if err := tx.proxyAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (tx *AccountCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionID:            tx.transactionID._ToProtobuf(),
		TransactionFee:           tx.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		Memo:                     tx.Transaction.memo,
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: tx.buildProtoBody(),
		},
	}
}
func (tx *AccountCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *AccountCreateTransaction) buildProtoBody() *services.CryptoCreateTransactionBody {
	body := &services.CryptoCreateTransactionBody{
		InitialBalance:                tx.initialBalance,
		ReceiverSigRequired:           tx.receiverSignatureRequired,
		Memo:                          tx.memo,
		MaxAutomaticTokenAssociations: tx.maxAutomaticTokenAssociations,
		DeclineReward:                 tx.declineReward,
		Alias:                         tx.alias,
	}

	if tx.key != nil {
		body.Key = tx.key._ToProtoKey()
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.stakedAccountID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedAccountId{StakedAccountId: tx.stakedAccountID._ToProtobuf()}
	} else if tx.stakedNodeID != nil {
		body.StakedId = &services.CryptoCreateTransactionBody_StakedNodeId{StakedNodeId: *tx.stakedNodeID}
	}

	return body
}

func (tx *AccountCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CreateAccount,
	}
}

func (tx *AccountCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
