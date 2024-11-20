package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountUpdateTransaction
// Change properties for the given account. Any null field is ignored (left unchanged). This
// transaction must be signed by the existing key for this account. If the transaction is changing
// the key field, then the transaction must be signed by both the old key (from before the change)
// and the new key. The old key must sign for security. The new key must sign as a safeguard to
// avoid accidentally changing to an invalid key, and then having no way to recover.
type AccountUpdateTransaction struct {
	*Transaction[*AccountUpdateTransaction]
	accountID                     *AccountID
	proxyAccountID                *AccountID
	key                           Key
	autoRenewPeriod               *time.Duration
	memo                          string
	receiverSignatureRequired     bool
	expirationTime                *time.Time
	maxAutomaticTokenAssociations int32
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
	tx := &AccountUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _AccountUpdateTransactionFromProtobuf(tx Transaction[*AccountUpdateTransaction], pb *services.TransactionBody) AccountUpdateTransaction {
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

	accountUpdateTransaction := AccountUpdateTransaction{
		accountID:                     _AccountIDFromProtobuf(pb.GetCryptoUpdateAccount().GetAccountIDToUpdate()),
		key:                           key,
		autoRenewPeriod:               &autoRenew,
		memo:                          pb.GetCryptoUpdateAccount().GetMemo().Value,
		receiverSignatureRequired:     receiverSignatureRequired,
		expirationTime:                &expiration,
		maxAutomaticTokenAssociations: pb.GetCryptoUpdateAccount().MaxAutomaticTokenAssociations.GetValue(),
		stakedAccountID:               stakeNodeAccountID,
		stakedNodeID:                  &stakedNodeID,
		declineReward:                 pb.GetCryptoUpdateAccount().GetDeclineReward().GetValue(),
	}

	tx.childTransaction = &accountUpdateTransaction
	accountUpdateTransaction.Transaction = &tx
	return accountUpdateTransaction
}

// SetKey Sets the new key for the Account
func (tx *AccountUpdateTransaction) SetKey(key Key) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.key = key
	return tx
}

func (tx *AccountUpdateTransaction) GetKey() (Key, error) {
	return tx.key, nil
}

// SetAccountID Sets the account ID which is being updated in tx transaction.
func (tx *AccountUpdateTransaction) SetAccountID(accountID AccountID) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

func (tx *AccountUpdateTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// Deprecated
func (tx *AccountUpdateTransaction) SetAliasKey(alias PublicKey) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.aliasKey = &alias
	return tx
}

// Deprecated
func (tx *AccountUpdateTransaction) GetAliasKey() PublicKey {
	if tx.aliasKey == nil {
		return PublicKey{}
	}

	return *tx.aliasKey
}

func (tx *AccountUpdateTransaction) SetStakedAccountID(id AccountID) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.stakedAccountID = &id
	return tx
}

func (tx *AccountUpdateTransaction) GetStakedAccountID() AccountID {
	if tx.stakedAccountID != nil {
		return *tx.stakedAccountID
	}

	return AccountID{}
}

func (tx *AccountUpdateTransaction) SetStakedNodeID(id int64) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.stakedNodeID = &id
	return tx
}

func (tx *AccountUpdateTransaction) GetStakedNodeID() int64 {
	if tx.stakedNodeID != nil {
		return *tx.stakedNodeID
	}

	return 0
}

func (tx *AccountUpdateTransaction) SetDeclineStakingReward(decline bool) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.declineReward = decline
	return tx
}

func (tx *AccountUpdateTransaction) ClearStakedAccountID() *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.stakedAccountID = &AccountID{Account: 0}
	return tx
}

func (tx *AccountUpdateTransaction) ClearStakedNodeID() *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	*tx.stakedNodeID = -1
	return tx
}

func (tx *AccountUpdateTransaction) GetDeclineStakingReward() bool {
	return tx.declineReward
}

// SetMaxAutomaticTokenAssociations
// Sets the maximum number of tokens that an Account can be implicitly associated with. Up to a 1000
// including implicit and explicit associations.
func (tx *AccountUpdateTransaction) SetMaxAutomaticTokenAssociations(max int32) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.maxAutomaticTokenAssociations = max
	return tx
}

func (tx *AccountUpdateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return tx.maxAutomaticTokenAssociations
}

// SetReceiverSignatureRequired
// If true, this account's key must sign any transaction depositing into this account (in
// addition to all withdrawals)
func (tx *AccountUpdateTransaction) SetReceiverSignatureRequired(receiverSignatureRequired bool) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.receiverSignatureRequired = receiverSignatureRequired
	return tx
}

func (tx *AccountUpdateTransaction) GetReceiverSignatureRequired() bool {
	return tx.receiverSignatureRequired
}

// Deprecated
// SetProxyAccountID Sets the ID of the account to which this account is proxy staked.
func (tx *AccountUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.proxyAccountID = &proxyAccountID
	return tx
}

// Deprecated
func (tx *AccountUpdateTransaction) GetProxyAccountID() AccountID {
	if tx.proxyAccountID == nil {
		return AccountID{}
	}

	return *tx.proxyAccountID
}

// SetAutoRenewPeriod Sets the duration in which it will automatically extend the expiration period.
func (tx *AccountUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &autoRenewPeriod
	return tx
}

func (tx *AccountUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime sets the new expiration time to extend to (ignored if equal to or before the current one)
func (tx *AccountUpdateTransaction) SetExpirationTime(expirationTime time.Time) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expirationTime
	return tx
}

func (tx *AccountUpdateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}
	return time.Time{}
}

// SetAccountMemo sets the new memo to be associated with the account (UTF-8 encoding max 100 bytes)
func (tx *AccountUpdateTransaction) SetAccountMemo(memo string) *AccountUpdateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo

	return tx
}

func (tx *AccountUpdateTransaction) GetAccountMemo() string {
	return tx.memo
}

// ----------- Overridden functions ----------------

func (tx AccountUpdateTransaction) getName() string {
	return "AccountUpdateTransaction"
}

func (tx AccountUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.proxyAccountID != nil {
		if err := tx.proxyAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx AccountUpdateTransaction) build() *services.TransactionBody {
	body := tx.buildProtoBody()

	pb := services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: body,
		},
	}

	body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: tx.maxAutomaticTokenAssociations}

	return &pb
}
func (tx AccountUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: tx.buildProtoBody(),
		},
	}, nil
}
func (tx AccountUpdateTransaction) buildProtoBody() *services.CryptoUpdateTransactionBody {
	body := &services.CryptoUpdateTransactionBody{
		ReceiverSigRequiredField: &services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper{
			ReceiverSigRequiredWrapper: &wrapperspb.BoolValue{Value: tx.receiverSignatureRequired},
		},
		Memo:                          &wrapperspb.StringValue{Value: tx.memo},
		DeclineReward:                 &wrapperspb.BoolValue{Value: tx.declineReward},
		MaxAutomaticTokenAssociations: &wrapperspb.Int32Value{Value: tx.maxAutomaticTokenAssociations},
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*tx.expirationTime)
	}

	if tx.accountID != nil {
		body.AccountIDToUpdate = tx.accountID._ToProtobuf()
	}

	if tx.key != nil {
		body.Key = tx.key._ToProtoKey()
	}

	if tx.stakedAccountID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedAccountId{StakedAccountId: tx.stakedAccountID._ToProtobuf()}
	} else if tx.stakedNodeID != nil {
		body.StakedId = &services.CryptoUpdateTransactionBody_StakedNodeId{StakedNodeId: *tx.stakedNodeID}
	}

	return body
}

func (tx AccountUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().UpdateAccount,
	}
}

func (tx AccountUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx AccountUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
