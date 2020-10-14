package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// Transaction contains the protobuf of a prepared transaction which can be signed and executed.
type Query struct {
	pb       *proto.Query
	pbHeader *proto.QueryHeader

	paymentTransactionID        TransactionID
	nodeID                      AccountID
	maxQueryPayment             Hbar
	queryPayment                Hbar
	nextPaymentTransactionIndex int
	nextTransactionIndex        int

	paymentTransactionNodeIDs []AccountID
	paymentTransactions       []*proto.Transaction

	isPaymentRequired bool
}

func newQuery(isPaymentRequired bool, queryHeader *proto.QueryHeader) Query {
	return Query{
		pb:                        &proto.Query{},
		pbHeader:                  queryHeader,
		paymentTransactionID:      TransactionID{},
		nextTransactionIndex:      0,
		paymentTransactions:       make([]*proto.Transaction, 0),
		paymentTransactionNodeIDs: make([]AccountID, 0),
		isPaymentRequired:         isPaymentRequired,
	}
}

func (query *Query) SetNodeId(accountID AccountID) *Query {
	query.nodeID = accountID
	return query
}

func (query *Query) GetNodeId() AccountID {
	return query.nodeID
}

func query_getNodeId(request request, client *Client) AccountID {
	if len(request.query.paymentTransactionNodeIDs) > 0 {
		return request.query.paymentTransactionNodeIDs[request.query.nextPaymentTransactionIndex]
	}

	if request.query.nodeID.isZero() {
		return request.query.nodeID
	} else {
		return client.getNextNode()
	}
}

func (query *Query) SetQueryPayment(queryPayment Hbar) *Query {
	query.queryPayment = queryPayment
	return query
}

func (query *Query) SetMaxQueryPayment(queryMaxPayment Hbar) *Query {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *Query) IsPaymentRequired() bool {
	return true
}

func query_makeRequest(request request) protoRequest {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.pbHeader.Payment = request.query.paymentTransactions[request.query.nextPaymentTransactionIndex]
	}
	request.query.pbHeader.ResponseType = proto.ResponseType_ANSWER_ONLY
	return protoRequest{
		query: request.query.pb,
	}
}

func query_advanceRequest(request request) {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.nextPaymentTransactionIndex++
	}
}

func query_mapResponse(request request, response response, _ AccountID, protoRequest protoRequest) (intermediateResponse, error) {
	return intermediateResponse{
		query: response.query,
	}, nil
}
