package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionRecordQuery struct {
	Query
	pb *proto.TransactionGetRecordQuery
}

func NewTransactionRecordQuery() *TransactionRecordQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
	pb := &proto.TransactionGetRecordQuery{Header: &header}

	return &TransactionRecordQuery{
		Query: query,
		pb:    pb,
	}
}

func TransactionRecordQuery_shouldRetry(status Status, response response) bool {
	if status == StatusBusy {
		return true
	}

	fmt.Printf("%+v\n", response.query)

	status = Status(response.query.GetTransactionGetRecord().TransactionRecord.Receipt.Status)

	switch status {
	case StatusBusy:
	case StatusUnknown:
	case StatusOk:
	case StatusReceiptNotFound:
	case StatusRecordNotFound:
		return true
	default:
		return false
	}

	return false
}

func TransactionRecordQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetTransactionGetRecord().Header.NodeTransactionPrecheckCode)
}

func TransactionRecordQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetTxRecordByTxID,
	}
}

func (query *TransactionRecordQuery) SetTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	query.pb.TransactionID = transactionID.toProtobuf()
	return query
}

func (query *TransactionRecordQuery) SetNodeAccountID(accountID AccountID) *TransactionRecordQuery {
	query.paymentTransactionNodeIDs = make([]AccountID, 0)
	query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, accountID)
	return query
}

func (query *TransactionRecordQuery) GetNodeAccountId(client *Client) AccountID {
	if query.paymentTransactionNodeIDs != nil {
		return query.paymentTransactionNodeIDs[query.nextPaymentTransactionIndex]
	}

	if query.nodeID.isZero() {
		return query.nodeID
	} else {
		return client.getNextNode()
	}
}

func (query *TransactionRecordQuery) SetQueryPayment(queryPayment Hbar) *TransactionRecordQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionRecordQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionRecordQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return TransactionRecord{}, errNoClientProvided
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	cost := query.queryPayment

	if len(query.paymentTransactionNodeIDs) == 0 {
		size := client.getNumberOfNodesForTransaction()
		for i := 0; i < size; i++ {
			query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, client.getNextNode())
		}
	}

	for _, nodeID := range query.paymentTransactionNodeIDs {
		transaction, err := query_makePaymentTransaction(
			query.paymentTransactionID,
			nodeID,
			client.operator,
			cost,
		)
		if err != nil {
			return TransactionRecord{}, err
		}

		query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, nodeID)
		query.paymentTransactions = append(query.paymentTransactions, transaction)
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		TransactionRecordQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeId,
		TransactionRecordQuery_getMethod,
		TransactionRecordQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return TransactionRecord{}, err
	}

	return TransactionRecordFromProtobuf(resp.query.GetTransactionGetRecord().TransactionRecord), nil
}
