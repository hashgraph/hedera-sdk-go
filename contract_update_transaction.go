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
	transaction
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
// used to construct and execute a Contract Update transaction.
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
	this := ContractUpdateTransaction{
		transaction: _NewTransaction(),
	}
	this._SetDefaultMaxTransactionFee(NewHbar(2))
	this.e = &this

	return &this
}

func _ContractUpdateTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *ContractUpdateTransaction {
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

	resultTx := &ContractUpdateTransaction{
		transaction:                   this,
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
	resultTx.e = resultTx
	return resultTx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *ContractUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractUpdateTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetContractID sets The Contract ID instance to update (this can't be changed on the contract)
func (this *ContractUpdateTransaction) SetContractID(contractID ContractID) *ContractUpdateTransaction {
	this.contractID = &contractID
	return this
}

func (this *ContractUpdateTransaction) GetContractID() ContractID {
	if this.contractID == nil {
		return ContractID{}
	}

	return *this.contractID
}

// Deprecated
func (this *ContractUpdateTransaction) SetBytecodeFileID(bytecodeFileID FileID) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.bytecodeFileID = &bytecodeFileID
	return this
}

// Deprecated
func (this *ContractUpdateTransaction) GetBytecodeFileID() FileID {
	if this.bytecodeFileID == nil {
		return FileID{}
	}

	return *this.bytecodeFileID
}

// SetAdminKey sets the key which can be used to arbitrarily modify the state of the instance by signing a
// ContractUpdateTransaction to modify it. If the admin key was never set then such modifications are not possible,
// and there is no administrator that can overrIDe the normal operation of the smart contract instance.
func (this *ContractUpdateTransaction) SetAdminKey(publicKey PublicKey) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.adminKey = publicKey
	return this
}

func (this *ContractUpdateTransaction) GetAdminKey() (Key, error) {
	return this.adminKey, nil
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this contract is proxy staked. If proxyAccountID is left unset,
// is an invalID account, or is an account that isn't a _Node, then this contract is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was never set.
func (this *ContractUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.proxyAccountID = &proxyAccountID
	return this
}

// Deprecated
func (this *ContractUpdateTransaction) GetProxyAccountID() AccountID {
	if this.proxyAccountID == nil {
		return AccountID{}
	}

	return *this.proxyAccountID
}

// SetAutoRenewPeriod sets the duration for which the contract instance will automatically charge its account to
// renew for.
func (this *ContractUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.autoRenewPeriod = &autoRenewPeriod
	return this
}

func (this *ContractUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if this.autoRenewPeriod != nil {
		return *this.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime extends the expiration of the instance and its account to the provIDed time. If the time provIDed
// is the current or past time, then there will be no effect.
func (this *ContractUpdateTransaction) SetExpirationTime(expiration time.Time) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.expirationTime = &expiration
	return this
}

func (this *ContractUpdateTransaction) GetExpirationTime() time.Time {
	if this.expirationTime != nil {
		return *this.expirationTime
	}

	return time.Time{}
}

// SetContractMemo sets the memo associated with the contract (max 100 bytes)
func (this *ContractUpdateTransaction) SetContractMemo(memo string) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.memo = memo
	// if transaction.pb.GetMemoWrapper() != nil {
	//	transaction.pb.GetMemoWrapper().Value = memo
	// } else {
	//	transaction.pb.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
	//		MemoWrapper: &wrapperspb.StringValue{Value: memo},
	//	}
	// }

	return this
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (this *ContractUpdateTransaction) SetAutoRenewAccountID(id AccountID) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.autoRenewAccountID = &id
	return this
}

func (this *ContractUpdateTransaction) GetAutoRenewAccountID() AccountID {
	if this.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *this.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (this *ContractUpdateTransaction) SetMaxAutomaticTokenAssociations(max int32) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.maxAutomaticTokenAssociations = max
	return this
}

func (this *ContractUpdateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return this.maxAutomaticTokenAssociations
}

func (this *ContractUpdateTransaction) GetContractMemo() string {
	return this.memo
}

func (this *ContractUpdateTransaction) SetStakedAccountID(id AccountID) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.stakedAccountID = &id
	return this
}

func (this *ContractUpdateTransaction) GetStakedAccountID() AccountID {
	if this.stakedAccountID != nil {
		return *this.stakedAccountID
	}

	return AccountID{}
}

func (this *ContractUpdateTransaction) SetStakedNodeID(id int64) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.stakedNodeID = &id
	return this
}

func (this *ContractUpdateTransaction) GetStakedNodeID() int64 {
	if this.stakedNodeID != nil {
		return *this.stakedNodeID
	}

	return 0
}

func (this *ContractUpdateTransaction) SetDeclineStakingReward(decline bool) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.declineReward = decline
	return this
}

func (this *ContractUpdateTransaction) GetDeclineStakingReward() bool {
	return this.declineReward
}

func (this *ContractUpdateTransaction) ClearStakedAccountID() *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.stakedAccountID = &AccountID{Account: 0}
	return this
}

func (this *ContractUpdateTransaction) ClearStakedNodeID() *ContractUpdateTransaction {
	this._RequireNotFrozen()
	*this.stakedNodeID = -1
	return this
}

func (this *ContractUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

// Sign uses the provided privateKey to sign the transaction.
func (this *ContractUpdateTransaction) Sign(
	privateKey PrivateKey,
) *ContractUpdateTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *ContractUpdateTransaction) SignWithOperator(
	client *Client,
) (*ContractUpdateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *ContractUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractUpdateTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *ContractUpdateTransaction) Freeze() (*ContractUpdateTransaction, error) {
	_, err := this.transaction.Freeze()
	return this, err
}

func (this *ContractUpdateTransaction) FreezeWith(client *Client) (*ContractUpdateTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractUpdateTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractUpdateTransaction) SetMaxTransactionFee(fee Hbar) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *ContractUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *ContractUpdateTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractUpdateTransaction.
func (this *ContractUpdateTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractUpdateTransaction.
func (this *ContractUpdateTransaction) SetTransactionMemo(memo string) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *ContractUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractUpdateTransaction.
func (this *ContractUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	ContractUpdateTransaction.
func (this *ContractUpdateTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractUpdateTransaction.
func (this *ContractUpdateTransaction) SetTransactionID(transactionID TransactionID) *ContractUpdateTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountID sets the _Node AccountID for this ContractUpdateTransaction.
func (this *ContractUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractUpdateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *ContractUpdateTransaction) SetMaxRetry(count int) *ContractUpdateTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *ContractUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractUpdateTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *ContractUpdateTransaction) SetMaxBackoff(max time.Duration) *ContractUpdateTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *ContractUpdateTransaction) SetMinBackoff(min time.Duration) *ContractUpdateTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *ContractUpdateTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractUpdateTransaction:%d", timestamp.UnixNano())
}

func (this *ContractUpdateTransaction) SetLogLevel(level LogLevel) *ContractUpdateTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *ContractUpdateTransaction) getName() string {
	return "ContractUpdateTransaction"
}

func (this *ContractUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.contractID != nil {
		if err := this.contractID.ValidateChecksum(client); err != nil {
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

func (this *ContractUpdateTransaction) build() *services.TransactionBody {
	body := &services.ContractUpdateTransactionBody{
		DeclineReward: &wrapperspb.BoolValue{Value: this.declineReward},
	}

	if this.maxAutomaticTokenAssociations != 0 {
		body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: this.maxAutomaticTokenAssociations}
	}

	if this.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*this.expirationTime)
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.adminKey != nil {
		body.AdminKey = this.adminKey._ToProtoKey()
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	if this.autoRenewAccountID != nil {
		body.AutoRenewAccountId = this.autoRenewAccountID._ToProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = this.memo
	} else {
		body.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrapperspb.StringValue{Value: this.memo},
		}
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	return &services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: body,
		},
	}
}

func (this *ContractUpdateTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractUpdateTransactionBody{
		DeclineReward: &wrapperspb.BoolValue{Value: this.declineReward},
	}

	if this.maxAutomaticTokenAssociations != 0 {
		body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: this.maxAutomaticTokenAssociations}
	}

	if this.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*this.expirationTime)
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.adminKey != nil {
		body.AdminKey = this.adminKey._ToProtoKey()
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	if this.autoRenewAccountID != nil {
		body.AutoRenewAccountId = this.autoRenewAccountID._ToProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = this.memo
	} else {
		body.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrapperspb.StringValue{Value: this.memo},
		}
	}

	if this.adminKey != nil {
		body.AdminKey = this.adminKey._ToProtoKey()
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: body,
		},
	}, nil
}

func (this *ContractUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().UpdateContract,
	}
}
