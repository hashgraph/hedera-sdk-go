package hedera

import (
	"fmt"
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
	default:
		return false
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

func (query *TransactionReceiptQuery) SetNodeId(accountID AccountID) *TransactionReceiptQuery {
	query.paymentTransactionNodeIDs = make([]AccountID, 0)
	query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, accountID)
	return query
}

func (query *TransactionReceiptQuery) GetNodeId(client *Client) AccountID {
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

	// query.queryPayment = NewHbar(0)
	// query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	// cost := query.queryPayment
	// cost := NewHbar(0)

	// if len(query.paymentTransactionNodeIDs) == 0 {
	// 	size := client.getNumberOfNodesForTransaction()
	// 	for i := 0; i < size; i++ {
	// 		query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, client.getNextNode())
	// 	}
	// }

	// for _, nodeID := range query.paymentTransactionNodeIDs {
	// 	transaction, err := query_makePaymentTransaction(
	// 		query.paymentTransactionID,
	// 		nodeID,
	// 		client.operator,
	// 		cost,
	// 	)
	// 	if err != nil {
	// 		return TransactionReceipt{}, err
	// 	}

	// 	query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, nodeID)
	// 	query.paymentTransactions = append(query.paymentTransactions, transaction)
	// }

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
