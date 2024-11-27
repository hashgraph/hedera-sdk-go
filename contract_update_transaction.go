package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractUpdateTransaction is used to modify a smart contract instance to have the given parameter values. Any nil
// field is ignored (left unchanged). If only the contractInstanceExpirationTime is being modified, then no signature is
// needed on this transaction other than for the account paying for the transaction itself. But if any of the other
// fields are being modified, then it must be signed by the adminKey. The use of adminKey is not currently supported in
// this API, but in the future will be implemented to allow these fields to be modified, and also to make modifications
// to the state of the instance. If the contract is created with no admin key, then none of the fields can be changed
// that need an admin signature, and therefore no admin key can ever be added. So if there is no admin key, then things
// like the bytecode are immutable. But if there is an admin key, then they can be changed.
//
// For example, the admin key might be a threshold key, which requires 3 of 5 binding arbitration judges to agree before
// the bytecode can be changed. This can be used to add flexibility to the management of smart contract behavior. But
// this is optional. If the smart contract is created without an admin key, then such a key can never be added, and its
// bytecode will be immutable.
type ContractUpdateTransaction struct {
	*Transaction[*ContractUpdateTransaction]
	contractID                    *ContractID
	proxyAccountID                *AccountID
	bytecodeFileID                *FileID
	adminKey                      Key
	autoRenewPeriod               *time.Duration
	expirationTime                *time.Time
	memo                          string
	autoRenewAccountID            *AccountID
	maxAutomaticTokenAssociations int32
	stakedAccountID               *AccountID
	stakedNodeID                  *int64
	declineReward                 bool
}

// NewContractUpdateTransaction creates a ContractUpdateTransaction transaction which can be
// used to construct and execute a Contract Update Transaction.
// ContractUpdateTransaction is used to modify a smart contract instance to have the given parameter values. Any nil
// field is ignored (left unchanged). If only the contractInstanceExpirationTime is being modified, then no signature is
// needed on this transaction other than for the account paying for the transaction itself. But if any of the other
// fields are being modified, then it must be signed by the adminKey. The use of adminKey is not currently supported in
// this API, but in the future will be implemented to allow these fields to be modified, and also to make modifications
// to the state of the instance. If the contract is created with no admin key, then none of the fields can be changed
// that need an admin signature, and therefore no admin key can ever be added. So if there is no admin key, then things
// like the bytecode are immutable. But if there is an admin key, then they can be changed.
//
// For example, the admin key might be a threshold key, which requires 3 of 5 binding arbitration judges to agree before
// the bytecode can be changed. This can be used to add flexibility to the management of smart contract behavior. But
// this is optional. If the smart contract is created without an admin key, then such a key can never be added, and its
// bytecode will be immutable.
func NewContractUpdateTransaction() *ContractUpdateTransaction {
	tx := &ContractUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _ContractUpdateTransactionFromProtobuf(tx Transaction[*ContractUpdateTransaction], pb *services.TransactionBody) ContractUpdateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetContractUpdateInstance().AdminKey)
	autoRenew := _DurationFromProtobuf(pb.GetContractUpdateInstance().GetAutoRenewPeriod())
	expiration := _TimeFromProtobuf(pb.GetContractUpdateInstance().GetExpirationTime())
	var memo string

	switch m := pb.GetContractUpdateInstance().GetMemoField().(type) {
	case *services.ContractUpdateTransactionBody_Memo:
		memo = m.Memo // nolint
	case *services.ContractUpdateTransactionBody_MemoWrapper:
		memo = m.MemoWrapper.Value
	}

	stakedNodeID := pb.GetContractUpdateInstance().GetStakedNodeId()

	var stakeNodeAccountID *AccountID
	if pb.GetContractUpdateInstance().GetStakedAccountId() != nil {
		stakeNodeAccountID = _AccountIDFromProtobuf(pb.GetContractUpdateInstance().GetStakedAccountId())
	}

	var autoRenewAccountID *AccountID
	if pb.GetContractUpdateInstance().AutoRenewAccountId != nil {
		autoRenewAccountID = _AccountIDFromProtobuf(pb.GetContractUpdateInstance().GetAutoRenewAccountId())
	}

	contractUpdateTransaction := ContractUpdateTransaction{
		contractID:                    _ContractIDFromProtobuf(pb.GetContractUpdateInstance().GetContractID()),
		adminKey:                      key,
		autoRenewPeriod:               &autoRenew,
		expirationTime:                &expiration,
		memo:                          memo,
		autoRenewAccountID:            autoRenewAccountID,
		maxAutomaticTokenAssociations: pb.GetContractUpdateInstance().MaxAutomaticTokenAssociations.GetValue(),
		stakedAccountID:               stakeNodeAccountID,
		stakedNodeID:                  &stakedNodeID,
		declineReward:                 pb.GetContractUpdateInstance().GetDeclineReward().GetValue(),
	}
	tx.childTransaction = &contractUpdateTransaction
	contractUpdateTransaction.Transaction = &tx
	return contractUpdateTransaction
}

// SetContractID sets The Contract ID instance to update (this can't be changed on the contract)
func (tx *ContractUpdateTransaction) SetContractID(contractID ContractID) *ContractUpdateTransaction {
	tx.contractID = &contractID
	return tx
}

func (tx *ContractUpdateTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// Deprecated
func (tx *ContractUpdateTransaction) SetBytecodeFileID(bytecodeFileID FileID) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.bytecodeFileID = &bytecodeFileID
	return tx
}

// Deprecated
func (tx *ContractUpdateTransaction) GetBytecodeFileID() FileID {
	if tx.bytecodeFileID == nil {
		return FileID{}
	}

	return *tx.bytecodeFileID
}

// SetAdminKey sets the key which can be used to arbitrarily modify the state of the instance by signing a
// ContractUpdateTransaction to modify it. If the admin key was never set then such modifications are not possible,
// and there is no administrator that can overrIDe the normal operation of the smart contract instance.
func (tx *ContractUpdateTransaction) SetAdminKey(publicKey PublicKey) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = publicKey
	return tx
}

func (tx *ContractUpdateTransaction) GetAdminKey() (Key, error) {
	return tx.adminKey, nil
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this contract is proxy staked. If proxyAccountID is left unset,
// is an invalID account, or is an account that isn't a _Node, then this contract is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was never set.
func (tx *ContractUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.proxyAccountID = &proxyAccountID
	return tx
}

// Deprecated
func (tx *ContractUpdateTransaction) GetProxyAccountID() AccountID {
	if tx.proxyAccountID == nil {
		return AccountID{}
	}

	return *tx.proxyAccountID
}

// SetAutoRenewPeriod sets the duration for which the contract instance will automatically charge its account to
// renew for.
func (tx *ContractUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &autoRenewPeriod
	return tx
}

func (tx *ContractUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime extends the expiration of the instance and its account to the provIDed time. If the time provIDed
// is the current or past time, then there will be no effect.
func (tx *ContractUpdateTransaction) SetExpirationTime(expiration time.Time) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

func (tx *ContractUpdateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// SetContractMemo sets the memo associated with the contract (max 100 bytes)
func (tx *ContractUpdateTransaction) SetContractMemo(memo string) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	// if transaction.pb.GetMemoWrapper() != nil {
	//	transaction.pb.GetMemoWrapper().Value = memo
	// } else {
	//	transaction.pb.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
	//		MemoWrapper: &wrapperspb.StringValue{Value: memo},
	//	}
	// }

	return tx
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (tx *ContractUpdateTransaction) SetAutoRenewAccountID(id AccountID) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &id
	return tx
}

func (tx *ContractUpdateTransaction) GetAutoRenewAccountID() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (tx *ContractUpdateTransaction) SetMaxAutomaticTokenAssociations(max int32) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.maxAutomaticTokenAssociations = max
	return tx
}

func (tx *ContractUpdateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return tx.maxAutomaticTokenAssociations
}

func (tx *ContractUpdateTransaction) GetContractMemo() string {
	return tx.memo
}

func (tx *ContractUpdateTransaction) SetStakedAccountID(id AccountID) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.stakedAccountID = &id
	return tx
}

func (tx *ContractUpdateTransaction) GetStakedAccountID() AccountID {
	if tx.stakedAccountID != nil {
		return *tx.stakedAccountID
	}

	return AccountID{}
}

func (tx *ContractUpdateTransaction) SetStakedNodeID(id int64) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.stakedNodeID = &id
	return tx
}

func (tx *ContractUpdateTransaction) GetStakedNodeID() int64 {
	if tx.stakedNodeID != nil {
		return *tx.stakedNodeID
	}

	return 0
}

func (tx *ContractUpdateTransaction) SetDeclineStakingReward(decline bool) *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.declineReward = decline
	return tx
}

func (tx *ContractUpdateTransaction) GetDeclineStakingReward() bool {
	return tx.declineReward
}

func (tx *ContractUpdateTransaction) ClearStakedAccountID() *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	tx.stakedAccountID = &AccountID{Account: 0}
	return tx
}

func (tx *ContractUpdateTransaction) ClearStakedNodeID() *ContractUpdateTransaction {
	tx._RequireNotFrozen()
	*tx.stakedNodeID = -1
	return tx
}

// ----------- Overridden functions ----------------

func (tx ContractUpdateTransaction) getName() string {
	return "ContractUpdateTransaction"
}

func (tx ContractUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
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

func (tx ContractUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: tx.buildProtoBody(),
		},
	}
}

func (tx ContractUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: tx.buildProtoBody(),
		},
	}, nil
}

func (tx ContractUpdateTransaction) buildProtoBody() *services.ContractUpdateTransactionBody {
	body := &services.ContractUpdateTransactionBody{
		DeclineReward: &wrapperspb.BoolValue{Value: tx.declineReward},
	}

	if tx.maxAutomaticTokenAssociations != 0 {
		body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: tx.maxAutomaticTokenAssociations}
	}

	if tx.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*tx.expirationTime)
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.contractID != nil {
		body.ContractID = tx.contractID._ToProtobuf()
	}

	if tx.autoRenewAccountID != nil {
		body.AutoRenewAccountId = tx.autoRenewAccountID._ToProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = tx.memo
	} else {
		body.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrapperspb.StringValue{Value: tx.memo},
		}
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.contractID != nil {
		body.ContractID = tx.contractID._ToProtobuf()
	}

	if tx.stakedAccountID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedAccountId{StakedAccountId: tx.stakedAccountID._ToProtobuf()}
	} else if tx.stakedNodeID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedNodeId{StakedNodeId: *tx.stakedNodeID}
	}

	return body
}

func (tx ContractUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().UpdateContract,
	}
}
func (tx ContractUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx ContractUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
