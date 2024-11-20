package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

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
	*Transaction[*NodeDeleteTransaction]
	nodeID uint64
}

func NewNodeDeleteTransaction() *NodeDeleteTransaction {
	tx := &NodeDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _NodeDeleteTransactionFromProtobuf(tx Transaction[*NodeDeleteTransaction], pb *services.TransactionBody) NodeDeleteTransaction {
	nodeDeleteTransaction := NodeDeleteTransaction{
		nodeID: pb.GetNodeDelete().NodeId,
	}

	tx.childTransaction = &nodeDeleteTransaction
	nodeDeleteTransaction.Transaction = &tx
	return nodeDeleteTransaction
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

// ----------- Overridden functions ----------------

func (tx NodeDeleteTransaction) getName() string {
	return "NodeDeleteTransaction"
}

func (tx NodeDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	return nil
}

func (tx NodeDeleteTransaction) build() *services.TransactionBody {
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

func (tx NodeDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_NodeDelete{
			NodeDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx NodeDeleteTransaction) buildProtoBody() *services.NodeDeleteTransactionBody {
	return &services.NodeDeleteTransactionBody{
		NodeId: tx.nodeID,
	}
}

func (tx NodeDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().DeleteNode,
	}
}

func (tx NodeDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx NodeDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
