package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// SPDX-License-Identifier: Apache-2.0

/**
 * A transaction to modify address book node attributes.
 *
 * - This transaction SHALL enable the node operator, as identified by the
 *   `admin_key`, to modify operational attributes of the node.
 * - This transaction MUST be signed by the active `admin_key` for the node.
 * - If this transaction sets a new value for the `admin_key`, then both the
 *   current `admin_key`, and the new `admin_key` MUST sign this transaction.
 * - This transaction SHALL NOT change any field that is not set (is null) in
 *   this transaction body.
 * - This SHALL create a pending update to the node, but the change SHALL NOT
 *   be immediately applied to the active configuration.
 * - All pending node updates SHALL be applied to the active network
 *   configuration during the next `freeze` transaction with the field
 *   `freeze_type` set to `PREPARE_UPGRADE`.
 *
 * ### Record Stream Effects
 * Upon completion the `node_id` for the updated entry SHALL be in the
 * transaction receipt.
 */
type NodeUpdateTransaction struct {
	*Transaction[*NodeUpdateTransaction]
	nodeID              uint64
	accountID           *AccountID
	description         string
	gossipEndpoints     []Endpoint
	serviceEndpoints    []Endpoint
	gossipCaCertificate []byte
	grpcCertificateHash []byte
	adminKey            Key
}

func NewNodeUpdateTransaction() *NodeUpdateTransaction {
	tx := &NodeUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _NodeUpdateTransactionFromProtobuf(tx Transaction[*NodeUpdateTransaction], pb *services.TransactionBody) NodeUpdateTransaction {
	adminKey, err := _KeyFromProtobuf(pb.GetNodeUpdate().GetAdminKey())
	if err != nil {
		return NodeUpdateTransaction{}
	}

	accountID := _AccountIDFromProtobuf(pb.GetNodeUpdate().GetAccountId())
	gossipEndpoints := make([]Endpoint, 0)
	for _, endpoint := range pb.GetNodeUpdate().GetGossipEndpoint() {
		gossipEndpoints = append(gossipEndpoints, EndpointFromProtobuf(endpoint))
	}
	serviceEndpoints := make([]Endpoint, 0)
	for _, endpoint := range pb.GetNodeUpdate().GetServiceEndpoint() {
		serviceEndpoints = append(serviceEndpoints, EndpointFromProtobuf(endpoint))
	}

	var certificate []byte
	if pb.GetNodeUpdate().GetGossipCaCertificate() != nil {
		certificate = pb.GetNodeUpdate().GetGossipCaCertificate().Value
	}

	var description string
	if pb.GetNodeUpdate().GetDescription() != nil {
		description = pb.GetNodeUpdate().GetDescription().Value
	}

	var certificateHash []byte
	if pb.GetNodeUpdate().GetGrpcCertificateHash() != nil {
		certificateHash = pb.GetNodeUpdate().GetGrpcCertificateHash().Value
	}

	nodeUpdateTransaction := NodeUpdateTransaction{
		nodeID:              pb.GetNodeUpdate().GetNodeId(),
		accountID:           accountID,
		description:         description,
		gossipEndpoints:     gossipEndpoints,
		serviceEndpoints:    serviceEndpoints,
		gossipCaCertificate: certificate,
		grpcCertificateHash: certificateHash,
		adminKey:            adminKey,
	}

	tx.childTransaction = &nodeUpdateTransaction
	nodeUpdateTransaction.Transaction = &tx
	return nodeUpdateTransaction
}

// GetNodeID he consensus node identifier in the network state.
func (tx *NodeUpdateTransaction) GetNodeID() uint64 {
	return tx.nodeID
}

// SetNodeID the consensus node identifier in the network state.
func (tx *NodeUpdateTransaction) SetNodeID(nodeID uint64) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.nodeID = nodeID
	return tx
}

// GetAccountID AccountID of the node
func (tx *NodeUpdateTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// SetAccountID get the AccountID of the node
func (tx *NodeUpdateTransaction) SetAccountID(accountID AccountID) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetDescription get the description of the node
func (tx *NodeUpdateTransaction) GetDescription() string {
	return tx.description
}

// SetDescription set the description of the node
func (tx *NodeUpdateTransaction) SetDescription(description string) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.description = description
	return tx
}

// SetDescription remove the description contents.
func (tx *NodeUpdateTransaction) ClearDescription(description string) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.description = ""
	return tx
}

// GetServiceEndpoints the list of service endpoints for gossip.
func (tx *NodeUpdateTransaction) GetGossipEndpoints() []Endpoint {
	return tx.gossipEndpoints
}

// SetServiceEndpoints the list of service endpoints for gossip.
func (tx *NodeUpdateTransaction) SetGossipEndpoints(gossipEndpoints []Endpoint) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.gossipEndpoints = gossipEndpoints
	return tx
}

// AddGossipEndpoint add an endpoint for gossip to the list of service endpoints for gossip.
func (tx *NodeUpdateTransaction) AddGossipEndpoint(endpoint Endpoint) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.gossipEndpoints = append(tx.gossipEndpoints, endpoint)
	return tx
}

// GetServiceEndpoints the list of service endpoints for gRPC calls.
func (tx *NodeUpdateTransaction) GetServiceEndpoints() []Endpoint {
	return tx.serviceEndpoints
}

// SetServiceEndpoints the list of service endpoints for gRPC calls.
func (tx *NodeUpdateTransaction) SetServiceEndpoints(serviceEndpoints []Endpoint) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = serviceEndpoints
	return tx
}

// AddServiceEndpoint the list of service endpoints for gRPC calls.
func (tx *NodeUpdateTransaction) AddServiceEndpoint(endpoint Endpoint) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = append(tx.serviceEndpoints, endpoint)
	return tx
}

// GetGossipCaCertificate the certificate used to sign gossip events.
func (tx *NodeUpdateTransaction) GetGossipCaCertificate() []byte {
	return tx.gossipCaCertificate
}

// SetGossipCaCertificate the certificate used to sign gossip events.
// This value MUST be the DER encoding of the certificate presented.
func (tx *NodeUpdateTransaction) SetGossipCaCertificate(gossipCaCertificate []byte) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.gossipCaCertificate = gossipCaCertificate
	return tx
}

// GetGrpcCertificateHash the hash of the node gRPC TLS certificate.
func (tx *NodeUpdateTransaction) GetGrpcCertificateHash() []byte {
	return tx.grpcCertificateHash
}

// SetGrpcCertificateHash the hash of the node gRPC TLS certificate.
// This value MUST be a SHA-384 hash.
func (tx *NodeUpdateTransaction) SetGrpcCertificateHash(grpcCertificateHash []byte) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.grpcCertificateHash = grpcCertificateHash
	return tx
}

// GetAdminKey an administrative key controlled by the node operator.
func (tx *NodeUpdateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetAdminKey an administrative key controlled by the node operator.
func (tx *NodeUpdateTransaction) SetAdminKey(adminKey Key) *NodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = adminKey
	return tx
}

// ----------- Overridden functions ----------------

func (tx NodeUpdateTransaction) getName() string {
	return "NodeUpdateTransaction"
}

func (tx NodeUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx NodeUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_NodeUpdate{
			NodeUpdate: tx.buildProtoBody(),
		},
	}
}

func (tx NodeUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_NodeUpdate{
			NodeUpdate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx NodeUpdateTransaction) buildProtoBody() *services.NodeUpdateTransactionBody {
	body := &services.NodeUpdateTransactionBody{
		Description: wrapperspb.String(tx.description),
		NodeId:      tx.nodeID,
	}

	if tx.accountID != nil {
		body.AccountId = tx.accountID._ToProtobuf()
	}

	for _, endpoint := range tx.gossipEndpoints {
		body.GossipEndpoint = append(body.GossipEndpoint, endpoint._ToProtobuf())
	}

	for _, endpoint := range tx.serviceEndpoints {
		body.ServiceEndpoint = append(body.ServiceEndpoint, endpoint._ToProtobuf())
	}

	if tx.gossipCaCertificate != nil {
		body.GossipCaCertificate = wrapperspb.Bytes(tx.gossipCaCertificate)
	}

	if tx.grpcCertificateHash != nil {
		body.GrpcCertificateHash = wrapperspb.Bytes(tx.grpcCertificateHash)
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	return body
}

func (tx NodeUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().UpdateNode,
	}
}

func (tx NodeUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx NodeUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
