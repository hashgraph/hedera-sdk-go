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
	transaction := ContractCreateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(131500 * time.Minute)
	transaction._SetDefaultMaxTransactionFee(NewHbar(20))

	return &transaction
}

func _ContractCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ContractCreateTransaction {
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
		Transaction:                   transaction,
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
func (transaction *ContractCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractCreateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetBytecodeFileID
// If the initcode is large (> 5K) then it must be stored in a file as hex encoded ascii.
func (transaction *ContractCreateTransaction) SetBytecodeFileID(byteCodeFileID FileID) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.byteCodeFileID = &byteCodeFileID
	transaction.initcode = nil
	return transaction
}

// GetBytecodeFileID returns the FileID of the file containing the contract's bytecode.
func (transaction *ContractCreateTransaction) GetBytecodeFileID() FileID {
	if transaction.byteCodeFileID == nil {
		return FileID{}
	}

	return *transaction.byteCodeFileID
}

// SetBytecode
// If it is small then it may either be stored as a hex encoded file or as a binary encoded field as part of the transaction.
func (transaction *ContractCreateTransaction) SetBytecode(code []byte) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.initcode = code
	transaction.byteCodeFileID = nil
	return transaction
}

// GetBytecode returns the bytecode for the contract.
func (transaction *ContractCreateTransaction) GetBytecode() []byte {
	return transaction.initcode
}

/**
 * Sets the state of the instance and its fields can be modified arbitrarily if this key signs a transaction
 * to modify it. If this is null, then such modifications are not possible, and there is no administrator
 * that can override the normal operation of this smart contract instance. Note that if it is created with no
 * admin keys, then there is no administrator to authorize changing the admin keys, so
 * there can never be any admin keys for that instance.
 */
func (transaction *ContractCreateTransaction) SetAdminKey(adminKey Key) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = adminKey
	return transaction
}

// GetAdminKey returns the key that can sign to modify the state of the instance
// and its fields can be modified arbitrarily if this key signs a transaction
func (transaction *ContractCreateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
}

// Sets the gas to run the constructor.
func (transaction *ContractCreateTransaction) SetGas(gas uint64) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.gas = int64(gas)
	return transaction
}

// GetGas returns the gas to run the constructor.
func (transaction *ContractCreateTransaction) GetGas() uint64 {
	return uint64(transaction.gas)
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.initialBalance = initialBalance.AsTinybar()
	return transaction
}

// GetInitialBalance gets the initial number of Hbar in the account
func (transaction *ContractCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(transaction.initialBalance)
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

func (transaction *ContractCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalID account, or is an account that isn't a _Node, then this account is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was not set.
func (transaction *ContractCreateTransaction) SetProxyAccountID(proxyAccountID AccountID) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.proxyAccountID = &proxyAccountID
	return transaction
}

// Deprecated
func (transaction *ContractCreateTransaction) GetProxyAccountID() AccountID {
	if transaction.proxyAccountID == nil {
		return AccountID{}
	}

	return *transaction.proxyAccountID
}

// SetConstructorParameters Sets the constructor parameters
func (transaction *ContractCreateTransaction) SetConstructorParameters(params *ContractFunctionParameters) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.parameters = params._Build(nil)
	return transaction
}

// SetConstructorParametersRaw Sets the constructor parameters as their raw bytes.
func (transaction *ContractCreateTransaction) SetConstructorParametersRaw(params []byte) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.parameters = params
	return transaction
}

// GetConstructorParameters returns the constructor parameters
func (transaction *ContractCreateTransaction) GetConstructorParameters() []byte {
	return transaction.parameters
}

// SetContractMemo Sets the memo to be associated with this contract.
func (transaction *ContractCreateTransaction) SetContractMemo(memo string) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

// GetContractMemo returns the memo associated with this contract.
func (transaction *ContractCreateTransaction) GetContractMemo() string {
	return transaction.memo
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (transaction *ContractCreateTransaction) SetAutoRenewAccountID(id AccountID) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &id
	return transaction
}

// GetAutoRenewAccountID returns the account to be used at the end of the auto renewal period
func (transaction *ContractCreateTransaction) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (transaction *ContractCreateTransaction) SetMaxAutomaticTokenAssociations(max int32) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

// GetMaxAutomaticTokenAssociations returns the maximum number of tokens that this contract can be automatically associated
func (transaction *ContractCreateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return transaction.maxAutomaticTokenAssociations
}

// SetStakedAccountID sets the account ID of the account to which this contract is staked.
func (transaction *ContractCreateTransaction) SetStakedAccountID(id AccountID) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedAccountID = &id
	return transaction
}

// GetStakedAccountID returns the account ID of the account to which this contract is staked.
func (transaction *ContractCreateTransaction) GetStakedAccountID() AccountID {
	if transaction.stakedAccountID != nil {
		return *transaction.stakedAccountID
	}

	return AccountID{}
}

// SetStakedNodeID sets the node ID of the node to which this contract is staked.
func (transaction *ContractCreateTransaction) SetStakedNodeID(id int64) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedNodeID = &id
	return transaction
}

// GetStakedNodeID returns the node ID of the node to which this contract is staked.
func (transaction *ContractCreateTransaction) GetStakedNodeID() int64 {
	if transaction.stakedNodeID != nil {
		return *transaction.stakedNodeID
	}

	return 0
}

// SetDeclineStakingReward sets if the contract should decline to pay the account's staking revenue.
func (transaction *ContractCreateTransaction) SetDeclineStakingReward(decline bool) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.declineReward = decline
	return transaction
}

// GetDeclineStakingReward returns if the contract should decline to pay the account's staking revenue.
func (transaction *ContractCreateTransaction) GetDeclineStakingReward() bool {
	return transaction.declineReward
}

func (transaction *ContractCreateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.byteCodeFileID != nil {
		if err := transaction.byteCodeFileID.ValidateChecksum(client); err != nil {
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

func (transaction *ContractCreateTransaction) _Build() *services.TransactionBody {
	body := &services.ContractCreateTransactionBody{
		Gas:                           transaction.gas,
		InitialBalance:                transaction.initialBalance,
		ConstructorParameters:         transaction.parameters,
		Memo:                          transaction.memo,
		MaxAutomaticTokenAssociations: transaction.maxAutomaticTokenAssociations,
		DeclineReward:                 transaction.declineReward,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.byteCodeFileID != nil {
		body.InitcodeSource = &services.ContractCreateTransactionBody_FileID{FileID: transaction.byteCodeFileID._ToProtobuf()}
	} else if len(transaction.initcode) != 0 {
		body.InitcodeSource = &services.ContractCreateTransactionBody_Initcode{Initcode: transaction.initcode}
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccountId = transaction.autoRenewAccountID._ToProtobuf()
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	pb := services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: body,
		},
	}

	return &pb
}

func (transaction *ContractCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractCreateTransactionBody{
		Gas:                           transaction.gas,
		InitialBalance:                transaction.initialBalance,
		ConstructorParameters:         transaction.parameters,
		Memo:                          transaction.memo,
		MaxAutomaticTokenAssociations: transaction.maxAutomaticTokenAssociations,
		DeclineReward:                 transaction.declineReward,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.byteCodeFileID != nil {
		body.InitcodeSource = &services.ContractCreateTransactionBody_FileID{FileID: transaction.byteCodeFileID._ToProtobuf()}
	} else if len(transaction.initcode) != 0 {
		body.InitcodeSource = &services.ContractCreateTransactionBody_Initcode{Initcode: transaction.initcode}
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccountId = transaction.autoRenewAccountID._ToProtobuf()
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.ContractCreateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCreateInstance{
			ContractCreateInstance: body,
		},
	}, nil
}

func _ContractCreateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().CreateContract,
	}
}

func (transaction *ContractCreateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractCreateTransaction) Sign(
	privateKey PrivateKey,
) *ContractCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ContractCreateTransaction) SignWithOperator(
	client *Client,
) (*ContractCreateTransaction, error) {
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
func (transaction *ContractCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractCreateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractCreateTransaction) Execute(
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
		_ContractCreateTransactionGetMethod,
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

func (transaction *ContractCreateTransaction) Freeze() (*ContractCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractCreateTransaction) FreezeWith(client *Client) (*ContractCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ContractCreateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractCreateTransaction) SetMaxTransactionFee(fee Hbar) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ContractCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ContractCreateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetTransactionMemo(memo string) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ContractCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	ContractCreateTransaction.
func (transaction *ContractCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetTransactionID(transactionID TransactionID) *ContractCreateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractCreateTransaction.
func (transaction *ContractCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ContractCreateTransaction) SetMaxRetry(count int) *ContractCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *ContractCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractCreateTransaction {
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
func (transaction *ContractCreateTransaction) SetMaxBackoff(max time.Duration) *ContractCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ContractCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ContractCreateTransaction) SetMinBackoff(min time.Duration) *ContractCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ContractCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ContractCreateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractCreateTransaction:%d", timestamp.UnixNano())
}

func (transaction *ContractCreateTransaction) SetLogLevel(level LogLevel) *ContractCreateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
