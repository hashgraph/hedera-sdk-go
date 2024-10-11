package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"
)

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

/**
 * A transaction to delete a node from the network address book.
 *
 * This transaction body SHALL be considered a "privileged transaction".
 *
 * - A transaction MUST be signed by the governing council.
 * - Upon success, the address book entry SHALL enter a "pending delete"
 *   state.
 * - All address book entries pending deletion SHALL be removed from the
 *   active network configuration during the next `freeze` transaction with
 *   the field `freeze_type` set to `PREPARE_UPGRADE`.<br/>
 * - A deleted address book node SHALL be removed entirely from network state.
 * - A deleted address book node identifier SHALL NOT be reused.
 *
 * ### Record Stream Effects
 * Upon completion the "deleted" `node_id` SHALL be in the transaction
 * receipt.
 */
type NodeDeleteTransaction struct {
	Transaction
	nodeID uint64
}

func NewNodeDeleteTransaction() *NodeDeleteTransaction {
	tx := &NodeDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _NodeDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *NodeDeleteTransaction {
	return &NodeDeleteTransaction{
		Transaction: transaction,
		nodeID:      pb.GetNodeDelete().NodeId,
	}
}

// GetNodeID he consensus node identifier in the network state.
func (tx *NodeDeleteTransaction) GetNodeID() uint64 {
	return tx.nodeID
}

// SetNodeID the consensus node identifier in the network state.
func (tx *NodeDeleteTransaction) SetNodeID(nodeID uint64) *NodeDeleteTransaction {
	tx._RequireNotFrozen()
	tx.nodeID = nodeID
	return tx
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *NodeDeleteTransaction) Sign(privateKey PrivateKey) *NodeDeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *NodeDeleteTransaction) SignWithOperator(client *Client) (*NodeDeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *NodeDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NodeDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *NodeDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *NodeDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *NodeDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *NodeDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *NodeDeleteTransaction) Freeze() (*NodeDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *NodeDeleteTransaction) FreezeWith(client *Client) (*NodeDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this NodeDeleteTransaction.
func (tx *NodeDeleteTransaction) SetMaxTransactionFee(fee Hbar) *NodeDeleteTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *NodeDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *NodeDeleteTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this NodeDeleteTransaction.
func (tx *NodeDeleteTransaction) SetTransactionMemo(memo string) *NodeDeleteTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this NodeDeleteTransaction.
func (tx *NodeDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *NodeDeleteTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *NodeDeleteTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this NodeDeleteTransaction.
func (tx *NodeDeleteTransaction) SetTransactionID(transactionID TransactionID) *NodeDeleteTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this NodeDeleteTransaction.
func (tx *NodeDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *NodeDeleteTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *NodeDeleteTransaction) SetMaxRetry(count int) *NodeDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *NodeDeleteTransaction) SetMaxBackoff(max time.Duration) *NodeDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *NodeDeleteTransaction) SetMinBackoff(min time.Duration) *NodeDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *NodeDeleteTransaction) SetLogLevel(level LogLevel) *NodeDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *NodeDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *NodeDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *NodeDeleteTransaction) getName() string {
	return "NodeDeleteTransaction"
}

func (tx *NodeDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	return nil
}

func (tx *NodeDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_NodeDelete{
			NodeDelete: tx.buildProtoBody(),
		},
	}
}

func (tx *NodeDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_NodeDelete{
			NodeDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *NodeDeleteTransaction) buildProtoBody() *services.NodeDeleteTransactionBody {
	return &services.NodeDeleteTransactionBody{
		NodeId: tx.nodeID,
	}
}

func (tx *NodeDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().DeleteNode,
	}
}

func (tx *NodeDeleteTransaction) preFreezeWith(client *Client) {
	// No special actions needed.
}

func (tx *NodeDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
