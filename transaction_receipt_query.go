package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionReceiptQuery struct {
	Query
	pb *proto.TransactionGetReceiptQuery
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
	pb := &proto.TransactionGetReceiptQuery{Header: &header}
	query.pb.Query = &proto.Query_TransactionGetReceipt{
		TransactionGetReceipt: pb,
	}

	return &TransactionReceiptQuery{
		Query: query,
		pb:    pb,
	}
}

func transactionReceiptQuery_shouldRetry(status Status, response response) bool {
	switch status {
	case StatusBusy, StatusUnknown, StatusReceiptNotFound:
		return true
	}

	status = Status(response.query.GetTransactionGetReceipt().Receipt.Status)

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound:
		return true
	default:
		return false
	}
}

func transactionReceiptQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetTransactionGetReceipt().Header.NodeTransactionPrecheckCode)
}

func transactionReceiptQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetTransactionReceipts,
	}
}

func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.pb.TransactionID = transactionID.toProtobuf()
	return query
}

func (query *TransactionReceiptQuery) SetNodeAccountID(accountID AccountID) *TransactionReceiptQuery {
	query.paymentTransactionNodeIDs = make([]AccountID, 0)
	query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, accountID)
	return query
}

func (query *TransactionReceiptQuery) GetNodeAccountId(client *Client) AccountID {
	if query.paymentTransactionNodeIDs != nil {
		return query.paymentTransactionNodeIDs[query.nextPaymentTransactionIndex]
	}

	if query.nodeID.isZero() {
		return query.nodeID
	} else {
		return client.getNextNode()
	}
}

func (query *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil || client.operator == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		transactionReceiptQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeId,
		transactionReceiptQuery_getMethod,
		transactionReceiptQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt().Receipt), nil
}
