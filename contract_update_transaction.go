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
	Transaction
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
	transaction := ContractUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _ContractUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *ContractUpdateTransaction {
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

	return &ContractUpdateTransaction{
		Transaction:                   transaction,
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
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *ContractUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractUpdateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetContractID sets The Contract ID instance to update (this can't be changed on the contract)
func (transaction *ContractUpdateTransaction) SetContractID(contractID ContractID) *ContractUpdateTransaction {
	transaction.contractID = &contractID
	return transaction
}

func (transaction *ContractUpdateTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

// Deprecated
func (transaction *ContractUpdateTransaction) SetBytecodeFileID(bytecodeFileID FileID) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.bytecodeFileID = &bytecodeFileID
	return transaction
}

// Deprecated
func (transaction *ContractUpdateTransaction) GetBytecodeFileID() FileID {
	if transaction.bytecodeFileID == nil {
		return FileID{}
	}

	return *transaction.bytecodeFileID
}

// SetAdminKey sets the key which can be used to arbitrarily modify the state of the instance by signing a
// ContractUpdateTransaction to modify it. If the admin key was never set then such modifications are not possible,
// and there is no administrator that can overrIDe the normal operation of the smart contract instance.
func (transaction *ContractUpdateTransaction) SetAdminKey(publicKey PublicKey) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAdminKey() (Key, error) {
	return transaction.adminKey, nil
}

// Deprecated
// SetProxyAccountID sets the ID of the account to which this contract is proxy staked. If proxyAccountID is left unset,
// is an invalID account, or is an account that isn't a _Node, then this contract is automatically proxy staked to a _Node
// chosen by the _Network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking,
// or if it is not currently running a _Node, then it will behave as if proxyAccountID was never set.
func (transaction *ContractUpdateTransaction) SetProxyAccountID(proxyAccountID AccountID) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.proxyAccountID = &proxyAccountID
	return transaction
}

// Deprecated
func (transaction *ContractUpdateTransaction) GetProxyAccountID() AccountID {
	if transaction.proxyAccountID == nil {
		return AccountID{}
	}

	return *transaction.proxyAccountID
}

// SetAutoRenewPeriod sets the duration for which the contract instance will automatically charge its account to
// renew for.
func (transaction *ContractUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return *transaction.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetExpirationTime extends the expiration of the instance and its account to the provIDed time. If the time provIDed
// is the current or past time, then there will be no effect.
func (transaction *ContractUpdateTransaction) SetExpirationTime(expiration time.Time) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *ContractUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetContractMemo sets the memo associated with the contract (max 100 bytes)
func (transaction *ContractUpdateTransaction) SetContractMemo(memo string) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	// if transaction.pb.GetMemoWrapper() != nil {
	//	transaction.pb.GetMemoWrapper().Value = memo
	// } else {
	//	transaction.pb.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
	//		MemoWrapper: &wrapperspb.StringValue{Value: memo},
	//	}
	// }

	return transaction
}

// SetAutoRenewAccountID
// An account to charge for auto-renewal of this contract. If not set, or set to an
// account with zero hbar balance, the contract's own hbar balance will be used to
// cover auto-renewal fees.
func (transaction *ContractUpdateTransaction) SetAutoRenewAccountID(id AccountID) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &id
	return transaction
}

func (transaction *ContractUpdateTransaction) GetAutoRenewAccountID() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// SetMaxAutomaticTokenAssociations
// The maximum number of tokens that this contract can be automatically associated
// with (i.e., receive air-drops from).
func (transaction *ContractUpdateTransaction) SetMaxAutomaticTokenAssociations(max int32) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxAutomaticTokenAssociations = max
	return transaction
}

func (transaction *ContractUpdateTransaction) GetMaxAutomaticTokenAssociations() int32 {
	return transaction.maxAutomaticTokenAssociations
}

func (transaction *ContractUpdateTransaction) GetContractMemo() string {
	return transaction.memo
}

func (transaction *ContractUpdateTransaction) SetStakedAccountID(id AccountID) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedAccountID = &id
	return transaction
}

func (transaction *ContractUpdateTransaction) GetStakedAccountID() AccountID {
	if transaction.stakedAccountID != nil {
		return *transaction.stakedAccountID
	}

	return AccountID{}
}

func (transaction *ContractUpdateTransaction) SetStakedNodeID(id int64) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedNodeID = &id
	return transaction
}

func (transaction *ContractUpdateTransaction) GetStakedNodeID() int64 {
	if transaction.stakedNodeID != nil {
		return *transaction.stakedNodeID
	}

	return 0
}

func (transaction *ContractUpdateTransaction) SetDeclineStakingReward(decline bool) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.declineReward = decline
	return transaction
}

func (transaction *ContractUpdateTransaction) GetDeclineStakingReward() bool {
	return transaction.declineReward
}

func (transaction *ContractUpdateTransaction) ClearStakedAccountID() *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.stakedAccountID = &AccountID{Account: 0}
	return transaction
}

func (transaction *ContractUpdateTransaction) ClearStakedNodeID() *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	*transaction.stakedNodeID = -1
	return transaction
}

func (transaction *ContractUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.ValidateChecksum(client); err != nil {
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

func (transaction *ContractUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.ContractUpdateTransactionBody{
		DeclineReward: &wrapperspb.BoolValue{Value: transaction.declineReward},
	}

	if transaction.maxAutomaticTokenAssociations != 0 {
		body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: transaction.maxAutomaticTokenAssociations}
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccountId = transaction.autoRenewAccountID._ToProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = transaction.memo
	} else {
		body.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrapperspb.StringValue{Value: transaction.memo},
		}
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: body,
		},
	}
}

func (transaction *ContractUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractUpdateTransactionBody{
		DeclineReward: &wrapperspb.BoolValue{Value: transaction.declineReward},
	}

	if transaction.maxAutomaticTokenAssociations != 0 {
		body.MaxAutomaticTokenAssociations = &wrapperspb.Int32Value{Value: transaction.maxAutomaticTokenAssociations}
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccountId = transaction.autoRenewAccountID._ToProtobuf()
	}

	if body.GetMemoWrapper() != nil {
		body.GetMemoWrapper().Value = transaction.memo
	} else {
		body.MemoField = &services.ContractUpdateTransactionBody_MemoWrapper{
			MemoWrapper: &wrapperspb.StringValue{Value: transaction.memo},
		}
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.stakedAccountID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedAccountId{StakedAccountId: transaction.stakedAccountID._ToProtobuf()}
	} else if transaction.stakedNodeID != nil {
		body.StakedId = &services.ContractUpdateTransactionBody_StakedNodeId{StakedNodeId: *transaction.stakedNodeID}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: body,
		},
	}, nil
}

func _ContractUpdateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().UpdateContract,
	}
}

func (transaction *ContractUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractUpdateTransaction) Sign(
	privateKey PrivateKey,
) *ContractUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *ContractUpdateTransaction) SignWithOperator(
	client *Client,
) (*ContractUpdateTransaction, error) {
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
func (transaction *ContractUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractUpdateTransaction) Execute(
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
		_ContractUpdateTransactionGetMethod,
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

func (transaction *ContractUpdateTransaction) Freeze() (*ContractUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractUpdateTransaction) FreezeWith(client *Client) (*ContractUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ContractUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *ContractUpdateTransaction) SetMaxTransactionFee(fee Hbar) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *ContractUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *ContractUpdateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetTransactionMemo(memo string) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *ContractUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetTransactionID(transactionID TransactionID) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this ContractUpdateTransaction.
func (transaction *ContractUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *ContractUpdateTransaction) SetMaxRetry(count int) *ContractUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *ContractUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractUpdateTransaction {
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
func (transaction *ContractUpdateTransaction) SetMaxBackoff(max time.Duration) *ContractUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *ContractUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *ContractUpdateTransaction) SetMinBackoff(min time.Duration) *ContractUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *ContractUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *ContractUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractUpdateTransaction:%d", timestamp.UnixNano())
}

func (transaction *ContractUpdateTransaction) SetLogLevel(level LogLevel) *ContractUpdateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
