package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
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
	*Transaction[*NodeDeleteTransaction]
	nodeID uint64
}

func NewNodeDeleteTransaction() *NodeDeleteTransaction {
	tx := &NodeDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _NodeDeleteTransactionFromProtobuf(pb *services.TransactionBody) *NodeDeleteTransaction {
	return &NodeDeleteTransaction{
		nodeID: pb.GetNodeDelete().NodeId,
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

func (tx *NodeDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *NodeDeleteTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*NodeDeleteTransaction](baseTx)
}
