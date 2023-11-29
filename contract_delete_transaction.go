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

	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// ContractDeleteTransaction marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
type ContractDeleteTransaction struct {
	transaction
	contractID        *ContractID
	transferContactID *ContractID
	transferAccountID *AccountID
	permanentRemoval  bool
}

// NewContractDeleteTransaction creates ContractDeleteTransaction which marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
func NewContractDeleteTransaction() *ContractDeleteTransaction {
	this := ContractDeleteTransaction{
		transaction: _NewTransaction(),
	}
	this._SetDefaultMaxTransactionFee(NewHbar(2))
	this.e = &this

	return &this
}

func _ContractDeleteTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *ContractDeleteTransaction {
	return &ContractDeleteTransaction{
		transaction:       this,
		contractID:        _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetContractID()),
		transferContactID: _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferContractID()),
		transferAccountID: _AccountIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferAccountID()),
		permanentRemoval:  pb.GetContractDeleteInstance().GetPermanentRemoval(),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *ContractDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractDeleteTransaction {
	this.transaction.SetGrpcDeadline(deadline)
	return this
}

// Sets the contract ID which will be deleted.
func (this *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.contractID = &contractID
	return this
}

// Returns the contract ID which will be deleted.
func (this *ContractDeleteTransaction) GetContractID() ContractID {
	if this.contractID == nil {
		return ContractID{}
	}

	return *this.contractID
}

// Sets the contract ID which will receive all remaining hbars.
func (this *ContractDeleteTransaction) SetTransferContractID(transferContactID ContractID) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transferContactID = &transferContactID
	return this
}

// Returns the contract ID which will receive all remaining hbars.
func (this *ContractDeleteTransaction) GetTransferContractID() ContractID {
	if this.transferContactID == nil {
		return ContractID{}
	}

	return *this.transferContactID
}

// Sets the account ID which will receive all remaining hbars.
func (this *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transferAccountID = &accountID

	return this
}

// Returns the account ID which will receive all remaining hbars.
func (this *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	if this.transferAccountID == nil {
		return AccountID{}
	}

	return *this.transferAccountID
}

// SetPermanentRemoval
// If set to true, means this is a "synthetic" system transaction being used to
// alert mirror nodes that the contract is being permanently removed from the ledger.
// IMPORTANT: User transactions cannot set this field to true, as permanent
// removal is always managed by the ledger itself. Any ContractDeleteTransaction
// submitted to HAPI with permanent_removal=true will be rejected with precheck status
// PERMANENT_REMOVAL_REQUIRES_SYSTEM_INITIATION.
func (this *ContractDeleteTransaction) SetPermanentRemoval(remove bool) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.permanentRemoval = remove

	return this
}

// GetPermanentRemoval returns true if this is a "synthetic" system transaction.
func (this *ContractDeleteTransaction) GetPermanentRemoval() bool {
	return this.permanentRemoval
}

func (this *ContractDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *ContractDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ContractDeleteTransaction {
	this.transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *ContractDeleteTransaction) SignWithOperator(
	client *Client,
) (*ContractDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *ContractDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractDeleteTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *ContractDeleteTransaction) Freeze() (*ContractDeleteTransaction, error) {
	_,err := this.transaction.Freeze()
	return this, err
}

func (this *ContractDeleteTransaction) FreezeWith(client *Client) (*ContractDeleteTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractDeleteTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *ContractDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *ContractDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *ContractDeleteTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this ContractDeleteTransaction.
func (this *ContractDeleteTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractDeleteTransaction.
func (this *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *ContractDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractDeleteTransaction.
func (this *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this ContractDeleteTransaction.
func (this *ContractDeleteTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractDeleteTransaction.
func (this *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractDeleteTransaction.
func (this *ContractDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *ContractDeleteTransaction) SetMaxRetry(count int) *ContractDeleteTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

func (this *ContractDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractDeleteTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *ContractDeleteTransaction) SetMaxBackoff(max time.Duration) *ContractDeleteTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *ContractDeleteTransaction) SetMinBackoff(min time.Duration) *ContractDeleteTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (this *ContractDeleteTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("ContractDeleteTransaction:%d", timestamp.UnixNano())
}

func (this *ContractDeleteTransaction) SetLogLevel(level LogLevel) *ContractDeleteTransaction {
	this.transaction.SetLogLevel(level)
	return this
}

// ----------- overriden functions ----------------

func (this *ContractDeleteTransaction) getName() string {
	return "ContractDeleteTransaction"
}
func (this *ContractDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.contractID != nil {
		if err := this.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if this.transferContactID != nil {
		if err := this.transferContactID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if this.transferAccountID != nil {
		if err := this.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *ContractDeleteTransaction) build() *services.TransactionBody {
	body := &services.ContractDeleteTransactionBody{
		PermanentRemoval: this.permanentRemoval,
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	if this.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: this.transferContactID._ToProtobuf(),
		}
	}

	if this.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: this.transferAccountID._ToProtobuf(),
		}
	}

	pb := services.TransactionBody{
		TransactionFee:           this.transactionFee,
		Memo:                     this.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		TransactionID:            this.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}

	return &pb
}
func (this *ContractDeleteTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractDeleteTransactionBody{
		PermanentRemoval: this.permanentRemoval,
	}

	if this.contractID != nil {
		body.ContractID = this.contractID._ToProtobuf()
	}

	if this.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: this.transferContactID._ToProtobuf(),
		}
	}

	if this.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: this.transferAccountID._ToProtobuf(),
		}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}, nil
}

func (this *ContractDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().DeleteContract,
	}
}