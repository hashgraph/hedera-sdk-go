package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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
	*Transaction[*AccountCreateTransaction]
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
	tx := &AccountCreateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _AccountCreateTransactionFromProtobuf(tx Transaction[*AccountCreateTransaction], pb *services.TransactionBody) AccountCreateTransaction {
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

	accountCreateTransaction := AccountCreateTransaction{
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
		accountCreateTransaction.alias = pb.GetCryptoCreateAccount().GetAlias()
	}

	tx.childTransaction = &accountCreateTransaction
	accountCreateTransaction.Transaction = &tx
	return accountCreateTransaction
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
	tx.stakedNodeID = nil
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
	tx.stakedAccountID = nil
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

// ----------- Overridden functions ----------------

func (tx AccountCreateTransaction) getName() string {
	return "AccountCreateTransaction"
}

func (tx AccountCreateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx AccountCreateTransaction) build() *services.TransactionBody {
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
func (tx AccountCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: tx.buildProtoBody(),
		},
	}, nil
}

func (tx AccountCreateTransaction) buildProtoBody() *services.CryptoCreateTransactionBody {
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

func (tx AccountCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CreateAccount,
	}
}

func (tx AccountCreateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx AccountCreateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
