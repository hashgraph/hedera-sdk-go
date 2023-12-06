package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use tx file except in compliance with the License.
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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ContractCreateTransaction which is used to start a new smart contract instance.
// After the instance is created, the ContractID for it is in the receipt, and can be retrieved by the Record or with a GetByKey query.
// The instance will run the bytecode, either stored in a previously created file or in the transaction body itself for
// small contracts.
type ContractCreateTransaction struct {
	Transaction
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
	tx := ContractCreateTransaction{
		Transaction: _NewTransaction(),
	}

	tx.SetAutoRenewPeriod(131500 * time.Minute)
	tx._SetDefaultMaxTransactionFee(NewHbar(20))
	tx.e = &tx

	return &tx
}

func _ContractCreateTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *ContractCreateTransaction {
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

	resultTx := &ContractCreateTransaction{
		Transaction:                   tx,
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
	resultTx.e = resultTx
	return resultTx
}

// SetBytecodeFileID
// If the initcode is large (> 5K) then it must be stored in a file as hex encoded ascii.
func (tx *ContractCreateTransaction) SetBytecodeFileID(byteCodeFileID FileID) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.byteCodeFileID = &byteCodeFileID
	tx.initcode = nil
	return tx
}

// GetBytecodeFileID returns the FileID of the file containing the contract's bytecode.
func (tx *ContractCreateTransaction) GetBytecodeFileID() FileID {
	if tx.byteCodeFileID == nil {
		return FileID{}
	}

	return *tx.byteCodeFileID
}

// SetBytecode
// If it is small then it may either be stored as a hex encoded file or as a binary encoded field as part of the transaction.
func (tx *ContractCreateTransaction) SetBytecode(code []byte) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.initcode = code
	tx.byteCodeFileID = nil
	return tx
}

// GetBytecode returns the bytecode for the contract.
func (tx *ContractCreateTransaction) GetBytecode() []byte {
	return tx.initcode
}

/**
 * Sets the state of the instance and its fields can be modified arbitrarily if tx key signs a transaction
 * to modify it. If tx is null, then such modifications are not possible, and there is no administrator
 * that can override the normal operation of tx smart contract instance. Note that if it is created with no
 * admin keys, then there is no administrator to authorize changing the admin keys, so
 * there can never be any admin keys for that instance.
 */
func (tx *ContractCreateTransaction) SetAdminKey(adminKey Key) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = adminKey
	return tx
}

// GetAdminKey returns the key that can sign to modify the state of the instance
// and its fields can be modified arbitrarily if tx key signs a transaction
func (tx *ContractCreateTransaction) GetAdminKey() (Key, error) {
	return tx.adminKey, nil
}

// Sets the gas to run the constructor.
func (tx *ContractCreateTransaction) SetGas(gas uint64) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.gas = int64(gas)
	return tx
}

// GetGas returns the gas to run the constructor.
func (tx *ContractCreateTransaction) GetGas() uint64 {
	return uint64(tx.gas)
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (tx *ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.initialBalance = initialBalance.AsTinybar()
	return tx
}

// GetInitialBalance gets the initial number of Hbar in the account
func (tx *ContractCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(tx.initialBalance)
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

func (tx *ContractCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which tx account is proxy staked. If proxyAccountID is not set,
// is an invalID account, or is an account that isn't a _Node, then tx account is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was not set.
func (tx *ContractCreateTransaction) SetProxyAccountID(proxyAccountID AccountID) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.proxyAccountID = &proxyAccountID
	return tx
}

// Deprecated
func (tx *ContractCreateTransaction) GetProxyAccountID() AccountID {
	if tx.proxyAccountID == nil {
		return AccountID{}
	}

	return *tx.proxyAccountID
}

// SetConstructorParameters Sets the constructor parameters
func (tx *ContractCreateTransaction) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.parameters = params._Build(nil)
	return tx
}

// SetConstructorParametersRaw Sets the constructor parameters as their raw bytes.
func (tx *ContractCreateTransaction) SetConstructorParametersRaw(params []byte) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.parameters = params
	return tx
}

// GetConstructorParameters returns the constructor parameters
func (tx *ContractCreateTransaction) GetConstructorParameters() []byte {
	return tx.parameters
}

// SetContractMemo Sets the memo to be associated with tx contract.
func (tx *ContractCreateTransaction) SetContractMemo(memo string) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetContractMemo returns the memo associated with tx contract.
func (tx *ContractCreateTransaction) GetContractMemo() string {
	return tx.memo
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of tx contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (tx *ContractCreateTransaction) SetAutoRenewAccountID(id AccountID) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &id
	return tx
}

// GetAutoRenewAccountID returns the account to be used at the end of the auto renewal period
func (tx *ContractCreateTransaction) GetAutoRenewAccountID() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that tx contract can be automatically associated
// with (i.e., receive air-drops from).
func (tx *ContractCreateTransaction) SetMaxAutomaticTokenAssociations(max int32) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.maxAutomaticTokenAssociations = max
	return tx
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that tx contract can be automatically associated
func (tx *ContractCreateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return tx.maxAutomaticTokenAssociations
}

// SetStakedAccountID sets the account ID of the account to which tx contract is staked.
func (tx *ContractCreateTransaction) SetStakedAccountID(id AccountID) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.stakedAccountID = &id
	return tx
}

// GetStakedAccountID returns the account ID of the account to which tx contract is staked.
func (tx *ContractCreateTransaction) GetStakedAccountID() AccountID {
	if tx.stakedAccountID != nil {
		return *tx.stakedAccountID
	}

	return AccountID{}
}

// SetStakedNodeID sets the node ID of the node to which tx contract is staked.
func (tx *ContractCreateTransaction) SetStakedNodeID(id int64) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.stakedNodeID = &id
	return tx
}

// GetStakedNodeID returns the node ID of the node to which tx contract is staked.
func (tx *ContractCreateTransaction) GetStakedNodeID() int64 {
	if tx.stakedNodeID != nil {
		return *tx.stakedNodeID
	}

	return 0
}

// SetDeclineStakingReward sets if the contract should decline to pay the account's staking revenue.
func (tx *ContractCreateTransaction) SetDeclineStakingReward(decline bool) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.declineReward = decline
	return tx
}

// GetDeclineStakingReward returns if the contract should decline to pay the account's staking revenue.
func (tx *ContractCreateTransaction) GetDeclineStakingReward() bool {
	return tx.declineReward
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *ContractCreateTransaction) Sign(
	privateKey PrivateKey,
) *ContractCreateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ContractCreateTransaction) SignWithOperator(
	client *Client,
) (*ContractCreateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.SignWithOperator(client)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *ContractCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractCreateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *ContractCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractCreateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ContractCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractCreateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ContractCreateTransaction) Freeze() (*ContractCreateTransaction, error) {
	_, err := tx.Transaction.Freeze()
	return tx, err
}

func (tx *ContractCreateTransaction) FreezeWith(client *Client) (*ContractCreateTransaction, error) {
	_, err := tx.Transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *ContractCreateTransaction) SetMaxTransactionFee(fee Hbar) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ContractCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx ContractCreateTransaction.
func (tx *ContractCreateTransaction) SetTransactionMemo(memo string) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx ContractCreateTransaction.
func (tx *ContractCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx ContractCreateTransaction.
func (tx *ContractCreateTransaction) SetTransactionID(transactionID TransactionID) *ContractCreateTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for tx ContractCreateTransaction.
func (tx *ContractCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ContractCreateTransaction) SetMaxRetry(count int) *ContractCreateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *ContractCreateTransaction) SetMaxBackoff(max time.Duration) *ContractCreateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ContractCreateTransaction) SetMinBackoff(min time.Duration) *ContractCreateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *ContractCreateTransaction) SetLogLevel(level LogLevel) *ContractCreateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *ContractCreateTransaction) getName() string {
	return "ContractCreateTransaction"
}

func (tx *ContractCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.byteCodeFileID != nil {
		if err := tx.byteCodeFileID.ValidateChecksum(client); err != nil {
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

func (tx *ContractCreateTransaction) build() *services.TransactionBody {
	pb := services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: tx.buildProtoBody(),
		},
	}

	return &pb
}
func (tx *ContractCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCreateInstance{
			ContractCreateInstance: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *ContractCreateTransaction) buildProtoBody() *services.ContractCreateTransactionBody {
	body := &services.ContractCreateTransactionBody{
		Gas:                           tx.gas,
		InitialBalance:                tx.initialBalance,
		ConstructorParameters:         tx.parameters,
		Memo:                          tx.memo,
		MaxAutomaticTokenAssociations: tx.maxAutomaticTokenAssociations,
		DeclineReward:                 tx.declineReward,
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.byteCodeFileID != nil {
		body.InitcodeSource = &services.ContractCreateTransactionBody_FileID{FileID: tx.byteCodeFileID._ToProtobuf()}
	} else if len(tx.initcode) != 0 {
		body.InitcodeSource = &services.ContractCreateTransactionBody_Initcode{Initcode: tx.initcode}
	}

	if tx.autoRenewAccountID != nil {
		body.AutoRenewAccountId = tx.autoRenewAccountID._ToProtobuf()
	}

	if tx.stakedAccountID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedAccountId{StakedAccountId: tx.stakedAccountID._ToProtobuf()}
	} else if tx.stakedNodeID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedNodeId{StakedNodeId: *tx.stakedNodeID}
	}

	return body
}

func (tx *ContractCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CreateContract,
	}
}

func (tx *ContractCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
