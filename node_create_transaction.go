package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

/**
 * A transaction to create a new node in the network address book.
 * The transaction, once complete, enables a new consensus node
 * to join the network, and requires governing council authorization.
 * <p>
 * This transaction body SHALL be considered a "privileged transaction".
 * <p>
 *
 * - MUST be signed by the governing council.
 * - MUST be signed by the `Key` assigned to the
 *   `admin_key` field.
 * - The newly created node information SHALL be added to the network address
 *   book information in the network state.
 * - The new entry SHALL be created in "state" but SHALL NOT participate in
 *   network consensus and SHALL NOT be present in network "configuration"
 *   until the next "upgrade" transaction (as noted below).
 * - All new address book entries SHALL be added to the active network
 *   configuration during the next `freeze` transaction with the field
 *   `freeze_type` set to `PREPARE_UPGRADE`.
 *
 * ### Record Stream Effects
 * Upon completion the newly assigned `node_id` SHALL be in the transaction
 * receipt.
 */
type NodeCreateTransaction struct {
	*Transaction[*NodeCreateTransaction]
	accountID           *AccountID
	description         string
	gossipEndpoints     []Endpoint
	serviceEndpoints    []Endpoint
	gossipCaCertificate []byte
	grpcCertificateHash []byte
	adminKey            Key
}

func NewNodeCreateTransaction() *NodeCreateTransaction {
	tx := &NodeCreateTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _NodeCreateTransactionFromProtobuf(tx Transaction[*NodeCreateTransaction], pb *services.TransactionBody) NodeCreateTransaction {
	adminKey, err := _KeyFromProtobuf(pb.GetNodeCreate().GetAdminKey())
	if err != nil {
		return NodeCreateTransaction{}
	}

	accountID := _AccountIDFromProtobuf(pb.GetNodeCreate().GetAccountId())
	gossipEndpoints := make([]Endpoint, 0)
	for _, endpoint := range pb.GetNodeCreate().GetGossipEndpoint() {
		gossipEndpoints = append(gossipEndpoints, EndpointFromProtobuf(endpoint))
	}
	serviceEndpoints := make([]Endpoint, 0)
	for _, endpoint := range pb.GetNodeCreate().GetServiceEndpoint() {
		serviceEndpoints = append(serviceEndpoints, EndpointFromProtobuf(endpoint))
	}

	nodeCreateTransaction := NodeCreateTransaction{
		accountID:           accountID,
		description:         pb.GetNodeCreate().GetDescription(),
		gossipEndpoints:     gossipEndpoints,
		serviceEndpoints:    serviceEndpoints,
		gossipCaCertificate: pb.GetNodeCreate().GetGossipCaCertificate(),
		grpcCertificateHash: pb.GetNodeCreate().GetGrpcCertificateHash(),
		adminKey:            adminKey,
	}

	tx.childTransaction = &nodeCreateTransaction
	nodeCreateTransaction.Transaction = &tx
	return nodeCreateTransaction
}

// GetAccountID AccountID of the node
func (tx *NodeCreateTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// SetAccountID get the AccountID of the node
func (tx *NodeCreateTransaction) SetAccountID(accountID AccountID) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetDescription get the description of the node
func (tx *NodeCreateTransaction) GetDescription() string {
	return tx.description
}

// SetDescription set the description of the node
func (tx *NodeCreateTransaction) SetDescription(description string) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.description = description
	return tx
}

// GetServiceEndpoints the list of service endpoints for gossip.
func (tx *NodeCreateTransaction) GetGossipEndpoints() []Endpoint {
	return tx.gossipEndpoints
}

// SetServiceEndpoints the list of service endpoints for gossip.
func (tx *NodeCreateTransaction) SetGossipEndpoints(gossipEndpoints []Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.gossipEndpoints = gossipEndpoints
	return tx
}

// AddGossipEndpoint add an endpoint for gossip to the list of service endpoints for gossip.
func (tx *NodeCreateTransaction) AddGossipEndpoint(endpoint Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.gossipEndpoints = append(tx.gossipEndpoints, endpoint)
	return tx
}

// GetServiceEndpoints the list of service endpoints for gRPC calls.
func (tx *NodeCreateTransaction) GetServiceEndpoints() []Endpoint {
	return tx.serviceEndpoints
}

// SetServiceEndpoints the list of service endpoints for gRPC calls.
func (tx *NodeCreateTransaction) SetServiceEndpoints(serviceEndpoints []Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = serviceEndpoints
	return tx
}

// AddServiceEndpoint the list of service endpoints for gRPC calls.
func (tx *NodeCreateTransaction) AddServiceEndpoint(endpoint Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = append(tx.serviceEndpoints, endpoint)
	return tx
}

// GetGossipCaCertificate the certificate used to sign gossip events.
func (tx *NodeCreateTransaction) GetGossipCaCertificate() []byte {
	return tx.gossipCaCertificate
}

// SetGossipCaCertificate the certificate used to sign gossip events.
// This value MUST be the DER encoding of the certificate presented.
func (tx *NodeCreateTransaction) SetGossipCaCertificate(gossipCaCertificate []byte) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.gossipCaCertificate = gossipCaCertificate
	return tx
}

// GetGrpcCertificateHash the hash of the node gRPC TLS certificate.
func (tx *NodeCreateTransaction) GetGrpcCertificateHash() []byte {
	return tx.grpcCertificateHash
}

// SetGrpcCertificateHash the hash of the node gRPC TLS certificate.
// This value MUST be a SHA-384 hash.
func (tx *NodeCreateTransaction) SetGrpcCertificateHash(grpcCertificateHash []byte) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.grpcCertificateHash = grpcCertificateHash
	return tx
}

// GetAdminKey an administrative key controlled by the node operator.
func (tx *NodeCreateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetAdminKey an administrative key controlled by the node operator.
func (tx *NodeCreateTransaction) SetAdminKey(adminKey Key) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = adminKey
	return tx
}

// ----------- Overridden functions ----------------

func (tx NodeCreateTransaction) getName() string {
	return "NodeCreateTransaction"
}

func (tx NodeCreateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx NodeCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_NodeCreate{
			NodeCreate: tx.buildProtoBody(),
		},
	}
}

func (tx NodeCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_NodeCreate{
			NodeCreate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx NodeCreateTransaction) buildProtoBody() *services.NodeCreateTransactionBody {
	body := &services.NodeCreateTransactionBody{
		Description: tx.description,
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
		body.GossipCaCertificate = tx.gossipCaCertificate
	}

	if tx.grpcCertificateHash != nil {
		body.GrpcCertificateHash = tx.grpcCertificateHash
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	return body
}

func (tx NodeCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().CreateNode,
	}
}

func (tx NodeCreateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx NodeCreateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
