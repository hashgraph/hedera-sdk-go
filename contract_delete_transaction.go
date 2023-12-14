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
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// ContractDeleteTransaction marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
type ContractDeleteTransaction struct {
	Transaction
	contractID        *ContractID
	transferContactID *ContractID
	transferAccountID *AccountID
	permanentRemoval  bool
}

// NewContractDeleteTransaction creates ContractDeleteTransaction which marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
func NewContractDeleteTransaction() *ContractDeleteTransaction {
	tx := ContractDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _ContractDeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *ContractDeleteTransaction {
	resultTx := &ContractDeleteTransaction{
		Transaction:       tx,
		contractID:        _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetContractID()),
		transferContactID: _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferContractID()),
		transferAccountID: _AccountIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferAccountID()),
		permanentRemoval:  pb.GetContractDeleteInstance().GetPermanentRemoval(),
	}
	return resultTx
}

// Sets the contract ID which will be deleted.
func (tx *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// Returns the contract ID which will be deleted.
func (tx *ContractDeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// Sets the contract ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) SetTransferContractID(transferContactID ContractID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transferContactID = &transferContactID
	return tx
}

// Returns the contract ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) GetTransferContractID() ContractID {
	if tx.transferContactID == nil {
		return ContractID{}
	}

	return *tx.transferContactID
}

// Sets the account ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transferAccountID = &accountID

	return tx
}

// Returns the account ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	if tx.transferAccountID == nil {
		return AccountID{}
	}

	return *tx.transferAccountID
}

// SetPermanentRemoval
// If set to true, means tx is a "synthetic" system transaction being used to
// alert mirror nodes that the contract is being permanently removed from the ledger.
// IMPORTANT: User transactions cannot set tx field to true, as permanent
// removal is always managed by the ledger itself. Any ContractDeleteTransaction
// submitted to HAPI with permanent_removal=true will be rejected with precheck status
// PERMANENT_REMOVAL_REQUIRES_SYSTEM_INITIATION.
func (tx *ContractDeleteTransaction) SetPermanentRemoval(remove bool) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.permanentRemoval = remove

	return tx
}

// GetPermanentRemoval returns true if tx is a "synthetic" system transaction.
func (tx *ContractDeleteTransaction) GetPermanentRemoval() bool {
	return tx.permanentRemoval
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *ContractDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ContractDeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ContractDeleteTransaction) SignWithOperator(
	client *Client,
) (*ContractDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *ContractDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

func (tx *ContractDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ContractDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ContractDeleteTransaction) Freeze() (*ContractDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *ContractDeleteTransaction) FreezeWith(client *Client) (*ContractDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *ContractDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ContractDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx ContractDeleteTransaction.
func (tx *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx ContractDeleteTransaction.
func (tx *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx ContractDeleteTransaction.
func (tx *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for tx ContractDeleteTransaction.
func (tx *ContractDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ContractDeleteTransaction) SetMaxRetry(count int) *ContractDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *ContractDeleteTransaction) SetMaxBackoff(max time.Duration) *ContractDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ContractDeleteTransaction) SetMinBackoff(min time.Duration) *ContractDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *ContractDeleteTransaction) SetLogLevel(level LogLevel) *ContractDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *ContractDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *ContractDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *ContractDeleteTransaction) getName() string {
	return "ContractDeleteTransaction"
}
func (tx *ContractDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.transferContactID != nil {
		if err := tx.transferContactID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.transferAccountID != nil {
		if err := tx.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *ContractDeleteTransaction) build() *services.TransactionBody {
	pb := services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: tx.buildProtoBody(),
		},
	}

	return &pb
}

func (tx *ContractDeleteTransaction) buildProtoBody() *services.ContractDeleteTransactionBody {
	body := &services.ContractDeleteTransactionBody{
		PermanentRemoval: tx.permanentRemoval,
	}

	if tx.contractID != nil {
		body.ContractID = tx.contractID._ToProtobuf()
	}

	if tx.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: tx.transferContactID._ToProtobuf(),
		}
	}

	if tx.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: tx.transferAccountID._ToProtobuf(),
		}
	}

	return body
}

func (tx *ContractDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *ContractDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().DeleteContract,
	}
}

func (tx *ContractDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
