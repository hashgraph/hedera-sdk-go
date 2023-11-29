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

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ContractCreateTransaction which is used to start a new smart contract instance.
// After the instance is created, the ContractID for it is in the receipt, and can be retrieved by the Record or with a GetByKey query.
// The instance will run the bytecode, either stored in a previously created file or in the transaction body itself for
// small contracts.
type ContractCreateTransaction struct {
	transaction
	byteCodeFileID                *FileID
	proxyAccountID                *AccountID
	adminKey                      Key
	gas                           int64
	initialBalance                int64
	autoRenewPeriod               *time.Duration
	parameters                    []byte
	memo                          string
	initcode                      []byte
	autoRenewAccountID            *AccountID
	maxAutomaticTokenAssociations int32
	stakedAccountID               *AccountID
	stakedNodeID                  *int64
	declineReward                 bool
}

// NewContractCreateTransaction creates ContractCreateTransaction which is used to start a new smart contract instance.
// After the instance is created, the ContractID for it is in the receipt, and can be retrieved by the Record or with a GetByKey query.
// The instance will run the bytecode, either stored in a previously created file or in the transaction body itself for
// small contracts.
func NewContractCreateTransaction() *ContractCreateTransaction {
	this := ContractCreateTransaction{
		transaction: _NewTransaction(),
	}

	this.SetAutoRenewPeriod(131500 * time.Minute)
	this._SetDefaultMaxTransactionFee(NewHbar(20))
	this.e = &this

	return &this
}

func _ContractCreateTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *ContractCreateTransaction {
	key, _ := _KeyFromProtobuf(pb.GetContractCreateInstance().GetAdminKey())
	autoRenew := _DurationFromProtobuf(pb.GetContractCreateInstance().GetAutoRenewPeriod())
	stakedNodeID := pb.GetContractCreateInstance().GetStakedNodeId()

	var stakeNodeAccountID *AccountID
	if pb.GetContractCreateInstance().GetStakedAccountId() != nil {
		stakeNodeAccountID = _AccountIDFromProtobuf(pb.GetContractCreateInstance().GetStakedAccountId())
	}

	var autoRenewAccountID *AccountID
	if pb.GetContractCreateInstance().AutoRenewAccountId != nil {
		autoRenewAccountID = _AccountIDFromProtobuf(pb.GetContractCreateInstance().GetAutoRenewAccountId())
	}

	return &ContractCreateTransaction{
		transaction:                   this,
		byteCodeFileID:                _FileIDFromProtobuf(pb.GetContractCreateInstance().GetFileID()),
		adminKey:                      key,
		gas:                           pb.GetContractCreateInstance().Gas,
		initialBalance:                pb.GetContractCreateInstance().InitialBalance,
		autoRenewPeriod:               &autoRenew,
		parameters:                    pb.GetContractCreateInstance().ConstructorParameters,
		memo:                          pb.GetContractCreateInstance().GetMemo(),
		initcode:                      pb.GetContractCreateInstance().GetInitcode(),
		autoRenewAccountID:            autoRenewAccountID,
		maxAutomaticTokenAssociations: pb.GetContractCreateInstance().MaxAutomaticTokenAssociations,
		stakedAccountID:               stakeNodeAccountID,
		stakedNodeID:                  &stakedNodeID,
		declineReward:                 pb.GetContractCreateInstance().GetDeclineReward(),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *ContractCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractCreateTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// SetBytecodeFileID
// If the initcode is large (> 5K) then it must be stored in a file as hex encoded ascii.
func (this *ContractCreateTransaction) SetBytecodeFileID(byteCodeFileID FileID) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.byteCodeFileID = &byteCodeFileID
	this.initcode = nil
	return this
}

// GetBytecodeFileID returns the FileID of the file containing the contract's bytecode.
func (this *ContractCreateTransaction) GetBytecodeFileID() FileID {
	if this.byteCodeFileID == nil {
		return FileID{}
	}

	return *this.byteCodeFileID
}

// SetBytecode
// If it is small then it may either be stored as a hex encoded file or as a binary encoded field as part of the transaction.
func (this *ContractCreateTransaction) SetBytecode(code []byte) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.initcode = code
	this.byteCodeFileID = nil
	return this
}

// GetBytecode returns the bytecode for the contract.
func (this *ContractCreateTransaction) GetBytecode() []byte {
	return this.initcode
}

/**
 * Sets the state of the instance and its fields can be modified arbitrarily if this key signs a transaction
 * to modify it. If this is null, then such modifications are not possible, and there is no administrator
 * that can override the normal operation of this smart contract instance. Note that if it is created with no
 * admin keys, then there is no administrator to authorize changing the admin keys, so
 * there can never be any admin keys for that instance.
 */
func (this *ContractCreateTransaction) SetAdminKey(adminKey Key) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.adminKey = adminKey
	return this
}

// GetAdminKey returns the key that can sign to modify the state of the instance
// and its fields can be modified arbitrarily if this key signs a transaction
func (this *ContractCreateTransaction) GetAdminKey() (Key, error) {
	return this.adminKey, nil
}

// Sets the gas to run the constructor.
func (this *ContractCreateTransaction) SetGas(gas uint64) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.gas = int64(gas)
	return this
}

// GetGas returns the gas to run the constructor.
func (this *ContractCreateTransaction) GetGas() uint64 {
	return uint64(this.gas)
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (this *ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.initialBalance = initialBalance.AsTinybar()
	return this
}

// GetInitialBalance gets the initial number of Hbar in the account
func (this *ContractCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(this.initialBalance)
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *ContractCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (this *ContractCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if this.autoRenewPeriod != nil {
		return *this.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalID account, or is an account that isn't a _Node, then this account is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was not set.
func (this *ContractCreateTransaction) SetProxyAccountID(proxyAccountID AccountID) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.proxyAccountID = &proxyAccountID
	return this
}

// Deprecated
func (this *ContractCreateTransaction) GetProxyAccountID() AccountID {
	if this.proxyAccountID == nil {
		return AccountID{}
	}

	return *this.proxyAccountID
}

// SetConstructorParameters Sets the constructor parameters
func (this *ContractCreateTransaction) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.parameters = params._Build(nil)
	return this
}

// SetConstructorParametersRaw Sets the constructor parameters as their raw bytes.
func (this *ContractCreateTransaction) SetConstructorParametersRaw(params []byte) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.parameters = params
	return this
}

// GetConstructorParameters returns the constructor parameters
func (this *ContractCreateTransaction) GetConstructorParameters() []byte {
	return this.parameters
}

// SetContractMemo Sets the memo to be associated with this contract.
func (this *ContractCreateTransaction) SetContractMemo(memo string) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.memo = memo
	return this
}

// GetContractMemo returns the memo associated with this contract.
func (this *ContractCreateTransaction) GetContractMemo() string {
	return this.memo
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (this *ContractCreateTransaction) SetAutoRenewAccountID(id AccountID) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.autoRenewAccountID = &id
	return this
}

// GetAutoRenewAccountID returns the account to be used at the end of the auto renewal period
func (this *ContractCreateTransaction) GetAutoRenewAccountID() AccountID {
	if this.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *this.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (this *ContractCreateTransaction) SetMaxAutomaticTokenAssociations(max int32) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.maxAutomaticTokenAssociations = max
	return this
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that this contract can be automatically associated
func (this *ContractCreateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return this.maxAutomaticTokenAssociations
}

// SetStakedAccountID sets the account ID of the account to which this contract is staked.
func (this *ContractCreateTransaction) SetStakedAccountID(id AccountID) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.stakedAccountID = &id
	return this
}

// GetStakedAccountID returns the account ID of the account to which this contract is staked.
func (this *ContractCreateTransaction) GetStakedAccountID() AccountID {
	if this.stakedAccountID != nil {
		return *this.stakedAccountID
	}

	return AccountID{}
}

// SetStakedNodeID sets the node ID of the node to which this contract is staked.
func (this *ContractCreateTransaction) SetStakedNodeID(id int64) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.stakedNodeID = &id
	return this
}

// GetStakedNodeID returns the node ID of the node to which this contract is staked.
func (this *ContractCreateTransaction) GetStakedNodeID() int64 {
	if this.stakedNodeID != nil {
		return *this.stakedNodeID
	}

	return 0
}

// SetDeclineStakingReward sets if the contract should decline to pay the account's staking revenue.
func (this *ContractCreateTransaction) SetDeclineStakingReward(decline bool) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.declineReward = decline
	return this
}

// GetDeclineStakingReward returns if the contract should decline to pay the account's staking revenue.
func (this *ContractCreateTransaction) GetDeclineStakingReward() bool {
	return this.declineReward
}

func (this *ContractCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (this *ContractCreateTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *ContractCreateTransaction) Sign(
	privateKey PrivateKey,
) *ContractCreateTransaction {
	this.transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *ContractCreateTransaction) SignWithOperator(
	client *Client,
) (*ContractCreateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *ContractCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractCreateTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *ContractCreateTransaction) Freeze() (*ContractCreateTransaction, error) {
	_,err := this.transaction.Freeze()
	return this, err
}

func (this *ContractCreateTransaction) FreezeWith(client *Client) (*ContractCreateTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractCreateTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractCreateTransaction) SetMaxTransactionFee(fee Hbar) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *ContractCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *ContractCreateTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractCreateTransaction.
func (this *ContractCreateTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractCreateTransaction.
func (this *ContractCreateTransaction) SetTransactionMemo(memo string) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *ContractCreateTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractCreateTransaction.
func (this *ContractCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this	ContractCreateTransaction.
func (this *ContractCreateTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractCreateTransaction.
func (this *ContractCreateTransaction) SetTransactionID(transactionID TransactionID) *ContractCreateTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractCreateTransaction.
func (this *ContractCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *ContractCreateTransaction) SetMaxRetry(count int) *ContractCreateTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *ContractCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractCreateTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *ContractCreateTransaction) SetMaxBackoff(max time.Duration) *ContractCreateTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *ContractCreateTransaction) SetMinBackoff(min time.Duration) *ContractCreateTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *ContractCreateTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractCreateTransaction:%d", timestamp.UnixNano())
}

func (this *ContractCreateTransaction) SetLogLevel(level LogLevel) *ContractCreateTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *ContractCreateTransaction) getName() string {
	return "ContractCreateTransaction"
}

func (this *ContractCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.byteCodeFileID != nil {
		if err := this.byteCodeFileID.ValidateChecksum(client); err != nil {
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

func (this *ContractCreateTransaction) build() *services.TransactionBody {
	body := &services.ContractCreateTransactionBody{
		Gas:                           this.gas,
		InitialBalance:                this.initialBalance,
		ConstructorParameters:         this.parameters,
		Memo:                          this.memo,
		MaxAutomaticTokenAssociations: this.maxAutomaticTokenAssociations,
		DeclineReward:                 this.declineReward,
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.adminKey != nil {
		body.AdminKey = this.adminKey._ToProtoKey()
	}

	if this.byteCodeFileID != nil {
		body.InitcodeSource = &services.ContractCreateTransactionBody_FileID{FileID: this.byteCodeFileID._ToProtobuf()}
	} else if len(this.initcode) != 0 {
		body.InitcodeSource = &services.ContractCreateTransactionBody_Initcode{Initcode: this.initcode}
	}

	if this.autoRenewAccountID != nil {
		body.AutoRenewAccountId = this.autoRenewAccountID._ToProtobuf()
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	pb := services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: body,
		},
	}

	return &pb
}
func (this *ContractCreateTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractCreateTransactionBody{
		Gas:                           this.gas,
		InitialBalance:                this.initialBalance,
		ConstructorParameters:         this.parameters,
		Memo:                          this.memo,
		MaxAutomaticTokenAssociations: this.maxAutomaticTokenAssociations,
		DeclineReward:                 this.declineReward,
	}

	if this.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*this.autoRenewPeriod)
	}

	if this.adminKey != nil {
		body.AdminKey = this.adminKey._ToProtoKey()
	}

	if this.byteCodeFileID != nil {
		body.InitcodeSource = &services.ContractCreateTransactionBody_FileID{FileID: this.byteCodeFileID._ToProtobuf()}
	} else if len(this.initcode) != 0 {
		body.InitcodeSource = &services.ContractCreateTransactionBody_Initcode{Initcode: this.initcode}
	}

	if this.autoRenewAccountID != nil {
		body.AutoRenewAccountId = this.autoRenewAccountID._ToProtobuf()
	}

	if this.stakedAccountID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedAccountId{StakedAccountId: this.stakedAccountID._ToProtobuf()}
	} else if this.stakedNodeID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedNodeId{StakedNodeId: *this.stakedNodeID}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCreateInstance{
			ContractCreateInstance: body,
		},
	}, nil
}

func (this *ContractCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CreateContract,
	}
}